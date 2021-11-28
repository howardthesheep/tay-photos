@echo off
SETLOCAL EnableDelayedExpansion

set GoFiles=

cd server
for /f "delims=" %%a in ('dir /b *.go') do (
 set "GoFiles=!GoFiles! %%a"
)

go run %GoFiles%

cd ..
