@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "ROOT_DIR=%~dp0"
cd /d "%ROOT_DIR%"

set "APP_NAME=metacritic-harvester"
set "ENTRY=./cmd/metacritic-harvester"
set "OUTPUT_ROOT=%ROOT_DIR%output"
set "RELEASE_ROOT=%OUTPUT_ROOT%\releases"
set "STAGE_ROOT=%RELEASE_ROOT%\_stage"
set "VERSION="
set "LDFLAGS=-s -w -buildid="

for /f "delims=" %%i in ('git describe --tags --always --dirty 2^>NUL') do set "VERSION=%%i"
if not defined VERSION set "VERSION=dev"

if not exist "%OUTPUT_ROOT%" mkdir "%OUTPUT_ROOT%"
if exist "%RELEASE_ROOT%" rmdir /s /q "%RELEASE_ROOT%"
mkdir "%RELEASE_ROOT%"
mkdir "%STAGE_ROOT%"

echo Building release artifacts for %APP_NAME% (%VERSION%)
echo Output: %RELEASE_ROOT%
echo.

call :build_target windows amd64 .exe || goto :fail
call :build_target windows arm64 .exe || goto :fail
call :build_target linux amd64 "" || goto :fail
call :build_target linux arm64 "" || goto :fail
call :build_target darwin amd64 "" || goto :fail
call :build_target darwin arm64 "" || goto :fail

powershell -NoProfile -Command ^
  "$ErrorActionPreference='Stop';" ^
  "$files = Get-ChildItem -Path '%RELEASE_ROOT%' -Filter '*.zip' | Sort-Object Name;" ^
  "$lines = foreach ($file in $files) { $hash = (Get-FileHash -Algorithm SHA256 $file.FullName).Hash.ToLowerInvariant(); '{0} *{1}' -f $hash, $file.Name };" ^
  "Set-Content -Path '%RELEASE_ROOT%\SHA256SUMS.txt' -Value $lines -Encoding ASCII"
if errorlevel 1 goto :fail

if exist "%STAGE_ROOT%" rmdir /s /q "%STAGE_ROOT%"

echo.
echo Release artifacts created successfully:
dir /b "%RELEASE_ROOT%"
exit /b 0

:build_target
set "TARGET_GOOS=%~1"
set "TARGET_GOARCH=%~2"
set "TARGET_EXT=%~3"
set "PACKAGE_BASENAME=%APP_NAME%_%VERSION%_%TARGET_GOOS%_%TARGET_GOARCH%"
set "PACKAGE_DIR=%STAGE_ROOT%\%PACKAGE_BASENAME%"
set "BINARY_PATH=%PACKAGE_DIR%\%APP_NAME%%TARGET_EXT%"
set "ARCHIVE_PATH=%RELEASE_ROOT%\%PACKAGE_BASENAME%.zip"

if exist "%PACKAGE_DIR%" rmdir /s /q "%PACKAGE_DIR%"
mkdir "%PACKAGE_DIR%"

echo [%TARGET_GOOS%/%TARGET_GOARCH%] go build
set "CGO_ENABLED=0"
set "GOOS=%TARGET_GOOS%"
set "GOARCH=%TARGET_GOARCH%"
go build -trimpath -ldflags "%LDFLAGS%" -o "%BINARY_PATH%" "%ENTRY%"
if errorlevel 1 exit /b 1

copy /y "README.md" "%PACKAGE_DIR%\README.md" >NUL
copy /y "LICENSE" "%PACKAGE_DIR%\LICENSE" >NUL

powershell -NoProfile -Command ^
  "$ErrorActionPreference='Stop';" ^
  "Compress-Archive -Path '%PACKAGE_DIR%' -DestinationPath '%ARCHIVE_PATH%' -Force"
if errorlevel 1 exit /b 1

exit /b 0

:fail
echo.
echo Build failed.
exit /b 1
