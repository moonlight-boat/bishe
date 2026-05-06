@echo off
setlocal
chcp 65001 >nul
set "SCRIPT_DIR=%~dp0"
for %%I in ("%SCRIPT_DIR%..") do set "PROJECT_ROOT=%%~fI"
set "WORKFLOW_ROOT=%PROJECT_ROOT%\职达工作流"
set "PYTHON_SCRIPT=%WORKFLOW_ROOT%\codes\main.py"
set "CSV_PATH=%WORKFLOW_ROOT%\auto_job_list.csv"
if not exist "%CSV_PATH%" set "CSV_PATH=%WORKFLOW_ROOT%\for_end_ui\codes\auto_job_list.csv"

if not exist "%CSV_PATH%" (
  echo [ERROR] CSV not found.
  echo Please generate: %WORKFLOW_ROOT%\auto_job_list.csv
  exit /b 1
)

set "DATABASE_PATH=%PROJECT_ROOT%\jobs.db"
set "SERVER_PORT=8080"
set "PYTHON_SCRIPT_PATH=%PYTHON_SCRIPT%"
set "CSV_FILE_PATH=%CSV_PATH%"

if /I "%~1"=="-RunCrawler" (
  echo [STEP] Run crawler first...
  pushd "%WORKFLOW_ROOT%"
  python "%PYTHON_SCRIPT%"
  if errorlevel 1 (
    echo [ERROR] Crawler failed.
    popd
    exit /b 1
  )
  popd
)

echo [STEP] Start backend in new window...
start "job-backend" cmd /k "cd /d "%PROJECT_ROOT%" && set DATABASE_PATH=%DATABASE_PATH% && set SERVER_PORT=%SERVER_PORT% && set PYTHON_SCRIPT_PATH=%PYTHON_SCRIPT_PATH% && set CSV_FILE_PATH=%CSV_FILE_PATH% && go run ."

echo [STEP] Wait backend health...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "$ok=$false; for($i=0;$i -lt 40;$i++){ try{ $h=Invoke-RestMethod -Uri 'http://localhost:8080/health' -TimeoutSec 3; if($h.status -eq 'ok'){ $ok=$true; break } } catch{}; Start-Sleep -Seconds 1 }; if(-not $ok){ exit 2 }"
if errorlevel 1 (
  echo [ERROR] Backend not ready in time.
  exit /b 1
)

echo [STEP] Trigger sync...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "$csv=[uri]::EscapeDataString('%CSV_PATH%'); Invoke-RestMethod -Method Post -Uri ('http://localhost:8080/api/sync?csv_path='+$csv) | ConvertTo-Json -Compress"

echo [STEP] Query jobs...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "Invoke-RestMethod -Uri 'http://localhost:8080/api/jobs?page=1&size=3' | Select-Object total,page,size,total_pages | ConvertTo-Json -Compress"

echo [STEP] Query stats...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "Invoke-RestMethod -Uri 'http://localhost:8080/api/jobs/stats' | Select-Object total | ConvertTo-Json -Compress"

echo.
echo === ONE CLICK DEBUG DONE ===
echo health: http://localhost:8080/health
echo jobs:   http://localhost:8080/api/jobs?page=1^&size=3
echo stats:  http://localhost:8080/api/jobs/stats
endlocal
