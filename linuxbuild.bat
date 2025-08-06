@echo off
echo Building Go Edge Device for Linux ARM64...

REM Check if Go is installed
go version >nul 2>&1
IF %ERRORLEVEL% NEQ 0 (
    echo ❌ Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Check if go.mod exists
IF NOT EXIST go.mod (
    echo ❌ go.mod file not found. Please run 'go mod init' first.
    pause
    exit /b 1
)

REM Set build environment
SET GOOS=linux
SET GOARCH=arm64

REM Always create .env file
echo Creating .env file...
echo # Edge Device Environment Configuration > .env
echo BASE_URL=https://by2e8b0nwd.execute-api.us-west-2.amazonaws.com/dev/v1 >> .env
echo ENCRYPTION_KEY=71ec9990f25c21a3dd3496c125a319e328c64dc675d5724fd0554aeac434ac68 >> .env
echo DEVICE_ID=GO-TEST-DEVICE >> .env
echo SECRET_KEY=SECRET_KEY >> .env
echo DEVICE_FROM=CRC_EVANS >> .env
echo IOT_ENDPOINT=a11km98y2evm1b-ats.iot.us-west-2.amazonaws.com:8883 >> .env
echo.
echo ⚠️  .env file (re)created. Please update with your actual values!
echo.

REM Check if certs directory exists and has files
IF NOT EXIST certs\ (
    echo Creating certs directory...
    mkdir certs
)

DIR certs\*.pem >nul 2>&1
IF %ERRORLEVEL% NEQ 0 (
    echo ⚠️  Warning: No certificate files found in certs/ directory
    echo    The application may fail to connect to MQTT endpoint without proper certificates
    echo.
)

REM Download dependencies
echo Downloading dependencies...
go mod download
IF %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to download dependencies
    pause
    exit /b 1
)

REM Clean previous build
IF EXIST go-edge-device (
    echo Removing previous build...
    del go-edge-device
)

REM Build the application
echo Building application...
go build -o go-edge-device

REM Check build result
IF %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Build succeeded!
    echo.
    echo 📋 Next steps:
    echo 1. Update the .env file with your actual values
    echo 2. Place certificate files in the certs/ directory
    echo 3. Transfer the binary to your Linux ARM64 device
    echo 4. Run: ./go-edge-device
    echo.
    echo 📁 Build output: go-edge-device
    echo.
) ELSE (
    echo.
    echo ❌ Build failed!
    echo Please check the error messages above.
    echo.
    pause
    exit /b 1
)
