@echo off

REM Set the custom binary name
SET BINARY_NAME=pixelita.exe

REM Build the Go source code with optimizations and without console window
go build -o "%BINARY_NAME%" -ldflags "-s -w -H=windowsgui" -gcflags "all=-N -l" .\cmd\main.go

REM Compress binary with UPX
upx %BINARY_NAME%

REM Check if the build was successful
IF %ERRORLEVEL% EQU 0 (
    echo Build successful! The binary '%BINARY_NAME%' is ready.
) ELSE (
    echo Build failed.
)
