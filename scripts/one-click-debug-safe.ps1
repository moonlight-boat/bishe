param(
  [switch]$RunCrawler = $false,
  [int]$Port = 8080
)

$ErrorActionPreference = "Stop"

function LogStep($m) { Write-Host "[STEP] $m" -ForegroundColor Cyan }
function LogOk($m) { Write-Host "[OK] $m" -ForegroundColor Green }
function LogWarn($m) { Write-Host "[WARN] $m" -ForegroundColor Yellow }

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$workflowRoot = Join-Path $projectRoot "职达工作流"

$pythonScriptPath = Join-Path $workflowRoot "codes\main.py"
$csvPath = Join-Path $workflowRoot "auto_job_list.csv"
if (-not (Test-Path $csvPath)) {
  $csvFallback = Join-Path $workflowRoot "for_end_ui\codes\auto_job_list.csv"
  if (Test-Path $csvFallback) {
    $csvPath = $csvFallback
    LogWarn ("CSV fallback: " + $csvPath)
  }
}
if (-not (Test-Path $csvPath)) { throw "CSV not found" }

$dbPath = Join-Path $projectRoot "jobs.db"
$server = "http://localhost:$Port"
$health = "$server/health"
$sync = "$server/api/sync?csv_path=$([uri]::EscapeDataString($csvPath))"
$jobs = "$server/api/jobs?page=1" + [char]38 + "size=3"
$stats = "$server/api/jobs/stats"

LogStep "Set env"
$env:DATABASE_PATH = $dbPath
$env:SERVER_PORT = "$Port"
$env:PYTHON_SCRIPT_PATH = $pythonScriptPath
$env:CSV_FILE_PATH = $csvPath

if ($RunCrawler) {
  LogStep "Run crawler"
  Push-Location $workflowRoot
  try {
    python $pythonScriptPath
    LogOk "Crawler finished"
  } finally {
    Pop-Location
  }
}

LogStep "Start backend window"
Start-Process powershell -WorkingDirectory $projectRoot -ArgumentList "-NoExit", "-Command", "go run ." | Out-Null

LogStep "Wait health"
$ok = $false
for ($i=0; $i -lt 40; $i++) {
  try {
    $h = Invoke-RestMethod -Uri $health -Method Get -TimeoutSec 3
    if ($h.status -eq "ok") { $ok = $true; break }
  } catch {
    Start-Sleep -Seconds 1
  }
}
if (-not $ok) { throw "Backend not ready" }
LogOk ("Backend ready: " + $health)

LogStep "Trigger sync"
$syncRes = Invoke-RestMethod -Uri $sync -Method Post -TimeoutSec 10
LogOk ("Sync response: " + ($syncRes | ConvertTo-Json -Compress))

LogStep "Query APIs"
$jobsRes = Invoke-RestMethod -Uri $jobs -Method Get -TimeoutSec 10
$statsRes = Invoke-RestMethod -Uri $stats -Method Get -TimeoutSec 10
LogOk ("Jobs total: " + $jobsRes.total)
LogOk ("Jobs sample size: " + $jobsRes.jobs.Count)
LogOk ("Stats total: " + $statsRes.total)

Write-Host ""
Write-Host "=== DONE ===" -ForegroundColor Green
Write-Host ("health: " + $health)
Write-Host ("jobs: " + $jobs)
Write-Host ("stats: " + $stats)
