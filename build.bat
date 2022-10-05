@echo off
rem # This script is used to start the mage build workflow. Without additional options it will simply recompile the project
rem #
rem # You might want to call this script with -h to get the mage help output
rem # You can also call the script with -l to list all valide targets
set SRCPIT_DIR=%~dp0%

pushd "%SRCPIT_DIR%build\workflow"
echo "go run %SRCPIT_DIR%build\workflow\mage.go -d %SRCPIT_DIR%build\workflow\magefiles -w %SRCPIT_DIR% %*"
go run "%SRCPIT_DIR%build\workflow\mage.go" -d "%SRCPIT_DIR%build\workflow\magefiles"  -w %SRCPIT_DIR% %*
set build_error_code=%ERRORLEVEL%
if "%build_error_code%" NEQ  "0"  (
    echo "Error: go run ./mage.go exit with error"
)
popd

exit /B %build_error_code%