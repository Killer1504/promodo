<#
.SYNOPSIS
    Build MSI installer for Pomodoro Focus Timer.

.DESCRIPTION
    Orchestrates the full MSI build pipeline:
    1. Validates prerequisites (Go, Node, Wails CLI, .NET SDK, WiX)
    2. Parses version from wails.json (or uses -Version override)
    3. Runs 'wails build' to produce the .exe
    4. Runs 'wix build' to package into .msi
    
    Output: build\bin\PomodoroFocusTimer-{version}-x64.msi

.PARAMETER Version
    Optional version override (e.g., "1.2.0"). If not specified, reads from wails.json info.productVersion.

.PARAMETER SkipWailsBuild
    Skip the 'wails build' step — useful when the exe already exists from a prior build.

.EXAMPLE
    .\build-msi.ps1
    .\build-msi.ps1 -Version "2.0.0"
    .\build-msi.ps1 -SkipWailsBuild
#>

param(
    [string]$Version = "",
    [switch]$SkipWailsBuild
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

# ─── Paths ───────────────────────────────────────────────────────────
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RepoRoot  = Resolve-Path (Join-Path $ScriptDir "..\..")
$BuildBin  = Join-Path $RepoRoot "build\bin"
$MsiDir    = $ScriptDir
$WxsFile   = Join-Path $MsiDir "Product.wxs"
$IconPath  = Join-Path $RepoRoot "build\windows\icon.ico"
$WailsJson = Join-Path $RepoRoot "wails.json"
$ExePath   = Join-Path $BuildBin "PomodoroFocusTimer.exe"

# ─── Helpers ─────────────────────────────────────────────────────────
function Write-Step([string]$msg) {
    Write-Host ""
    Write-Host "=== $msg ===" -ForegroundColor Cyan
}

function Assert-Command([string]$cmd, [string]$label) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Host "ERROR: '$label' is not installed or not in PATH." -ForegroundColor Red
        Write-Host "  Install it and try again." -ForegroundColor Yellow
        exit 1
    }
}

# ─── Step 1: Validate Prerequisites ─────────────────────────────────
Write-Step "Validating prerequisites"

Assert-Command "go"     "Go"
Assert-Command "node"   "Node.js"
Assert-Command "wails"  "Wails CLI"
Assert-Command "dotnet" ".NET SDK"
Assert-Command "wix"    "WiX Toolset"

Write-Host "  Go:      $(go version)" -ForegroundColor Green
Write-Host "  Node:    $(node --version)" -ForegroundColor Green
Write-Host "  Wails:   $(wails version 2>&1 | Select-Object -First 1)" -ForegroundColor Green
Write-Host "  .NET:    $(dotnet --version)" -ForegroundColor Green
Write-Host "  WiX:     $(wix --version)" -ForegroundColor Green

# ─── Step 2: Parse Version ───────────────────────────────────────────
Write-Step "Resolving version"

if ([string]::IsNullOrEmpty($Version)) {
    if (-not (Test-Path $WailsJson)) {
        Write-Host "ERROR: wails.json not found at '$WailsJson'" -ForegroundColor Red
        exit 1
    }
    $config = Get-Content $WailsJson -Raw | ConvertFrom-Json
    $Version = $config.info.productVersion
    
    if ([string]::IsNullOrEmpty($Version)) {
        Write-Host "ERROR: 'info.productVersion' not found in wails.json" -ForegroundColor Red
        Write-Host "  Add: `"info`": { `"productVersion`": `"1.0.0`" }" -ForegroundColor Yellow
        exit 1
    }
    Write-Host "  Version from wails.json: $Version" -ForegroundColor Green
} else {
    Write-Host "  Version override: $Version" -ForegroundColor Green
}

# Validate semver format (MAJOR.MINOR.PATCH)
if ($Version -notmatch '^\d+\.\d+\.\d+$') {
    Write-Host "ERROR: Version '$Version' must be in MAJOR.MINOR.PATCH format (e.g., 1.0.0)" -ForegroundColor Red
    exit 1
}

$MsiOutputName = "PomodoroFocusTimer-$Version-x64.msi"
$MsiOutputPath = Join-Path $BuildBin $MsiOutputName

# ─── Step 3: Wails Build ────────────────────────────────────────────
if ($SkipWailsBuild) {
    Write-Step "Skipping Wails build (using existing exe)"
    if (-not (Test-Path $ExePath)) {
        Write-Host "ERROR: Exe not found at '$ExePath'. Run without -SkipWailsBuild first." -ForegroundColor Red
        exit 1
    }
} else {
    Write-Step "Building application with Wails"
    Push-Location $RepoRoot
    try {
        wails build -platform windows/amd64
        if ($LASTEXITCODE -ne 0) {
            Write-Host "ERROR: 'wails build' failed with exit code $LASTEXITCODE" -ForegroundColor Red
            exit 1
        }
    } finally {
        Pop-Location
    }
}

if (-not (Test-Path $ExePath)) {
    Write-Host "ERROR: Expected exe not found at '$ExePath'" -ForegroundColor Red
    exit 1
}

$exeSize = (Get-Item $ExePath).Length / 1MB
Write-Host "  Exe built: $ExePath ($([math]::Round($exeSize, 2)) MB)" -ForegroundColor Green

# ─── Step 4: Build MSI with WiX ─────────────────────────────────────
Write-Step "Packaging MSI with WiX"

# Ensure output directory exists
if (-not (Test-Path $BuildBin)) {
    New-Item -ItemType Directory -Force -Path $BuildBin | Out-Null
}

Push-Location $MsiDir
try {
    $wixArgs = @(
        "build"
        $WxsFile
        "-ext", "WixToolset.UI.wixext"
        "-d", "Version=$Version"
        "-d", "BuildDir=$BuildBin"
        "-d", "IconPath=$IconPath"
        "-arch", "x64"
        "-o", $MsiOutputPath
    )
    
    Write-Host "  Running: wix $($wixArgs -join ' ')" -ForegroundColor DarkGray
    & wix @wixArgs
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: 'wix build' failed with exit code $LASTEXITCODE" -ForegroundColor Red
        exit 1
    }
} finally {
    Pop-Location
}

# ─── Step 5: Report ─────────────────────────────────────────────────
Write-Step "Build Complete"

if (Test-Path $MsiOutputPath) {
    $msiSize = (Get-Item $MsiOutputPath).Length / 1MB
    Write-Host ""
    Write-Host "  MSI Output: $MsiOutputPath" -ForegroundColor Green
    Write-Host "  Version:    $Version" -ForegroundColor Green
    Write-Host "  Size:       $([math]::Round($msiSize, 2)) MB" -ForegroundColor Green
    Write-Host ""
    
    if ($msiSize -gt 20) {
        Write-Host "  WARNING: MSI exceeds 20 MB target size" -ForegroundColor Yellow
    }
} else {
    Write-Host "ERROR: MSI file was not created at '$MsiOutputPath'" -ForegroundColor Red
    exit 1
}
