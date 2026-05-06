param(
  [switch]$RunCrawler = $false,
  [int]$Port = 8080
)

$ErrorActionPreference = 'Stop'

function Step([string]$msg) { Write-Host ('[STEP] ' + $msg) -ForegroundColor Cyan }
function Ok([string]$msg) { Write-Host ('[OK] ' + $msg) -ForegroundColor Green }
function Warn([string]$msg) { Write-Host ('[WARN] ' + $msg) -ForegroundColor Yellow }

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$workflowRoot = Join-Path $projectRoot '职达工作流'

$pythonScriptPath = Join-Path $workflowRoot 'codes\main.py'
$csvPath = Join-Path $workflowRoot 'auto_job_list.csv'
if (-not (Test-Path $csvPath)) {
  $fallbackCsv = Join-Path $workflowRoot 'for_end_ui\codes\auto_job_list.csv'
  if (Test-Path $fallbackCsv) {
    $csvPath = $fallbackCsv
    Warn ('use fallback csv: ' + $csvPath)
  }
}
if (-not (Test-Path $csvPath)) { throw 'csv not found' }

$databasePath = Join-Path $projectRoot 'jobs.db'
$serverUrl = ('http://localhost:' + $Port)
$healthUrl = ($serverUrl + '/health')
$syncUrl = ($serverUrl + '/api/sync?csv_path=' + [System.Uri]::EscapeDataString($csvPath))
$jobsUrl = ($serverUrl + '/api/jobs?page=1') + [char]38 + 'size=3'
$statsUrl = ($serverUrl + '/api/jobs/stats')

Step 'set env vars'
$env:DATABASE_PATH = $databasePath
$env:SERVER_PORT = [string]$Port
$env:PYTHON_SCRIPT_PATH = $pythonScriptPath
$env:CSV_FILE_PATH = $csvPath

if ($RunCrawler) {
  Step 'run crawler first'
  Push-Location $workflowRoot
  try {
    python $pythonScriptPath
    Ok 'crawler done'
  } finally {
    Pop-Location
  }
}

Step 'start backend window'
Start-Process -FilePath 'powershell' `
  -WorkingDirectory $projectRoot `
  -ArgumentList '-NoExit', '-Command', 'go run .' | Out-Null

Step 'wait health'
$healthy = $false
for ($i = 0; $i -lt 40; $i++) {
  try {
    $h = Invoke-RestMethod -Uri $healthUrl -Method Get -TimeoutSec 3
    if ($h.status -eq 'ok') { $healthy = $true; break }
  } catch {
    Start-Sleep -Seconds 1
  }
}
if (-not $healthy) { throw 'backend not ready in time' }
Ok ('backend ready: ' + $healthUrl)

Step 'trigger sync'
$syncRes = Invoke-RestMethod -Uri $syncUrl -Method Post -TimeoutSec 10
Ok ('sync response: ' + ($syncRes | ConvertTo-Json -Compress))

Step 'query jobs/stats'
$jobsRes = Invoke-RestMethod -Uri $jobsUrl -Method Get -TimeoutSec 10
$statsRes = Invoke-RestMethod -Uri $statsUrl -Method Get -TimeoutSec 10
Ok ('jobs total: ' + $jobsRes.total)
Ok ('jobs sample size: ' + $jobsRes.jobs.Count)
Ok ('stats total: ' + $statsRes.total)

Write-Host ''
Write-Host '=== ONE CLICK DEBUG DONE ===' -ForegroundColor Green
Write-Host ('server: ' + $serverUrl)
Write-Host ('csv: ' + $csvPath)
Write-Host ('1) ' + $healthUrl)
Write-Host ('2) ' + $jobsUrl)
Write-Host ('3) ' + $statsUrl)
