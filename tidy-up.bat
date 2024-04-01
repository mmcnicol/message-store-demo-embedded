@echo off
cd /d "%~dp0"  REM Change the working directory to the batch file's directory

REM Delete files with .data suffix
del *.data /q

REM Delete files with .data.idx suffix
del *.data.idx /q

REM Display message after completion
echo Files with .data and .data.idx suffixes deleted.
pause
