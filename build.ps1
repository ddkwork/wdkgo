# HyperDbg 统一工程构建脚本

$ErrorActionPreference = "Stop"

$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$BUILD_DIR = Join-Path $SCRIPT_DIR "build"
$CONFIG = "Release"
$WDK_PATH = "E:\Program Files\Windows Kits\10"
$EWDK_ISO_PATH = "D:\Admin\vs2022VC\EWDK_br_release_28000_251103-1709.iso"
$EWDK_MOUNT_LETTER = "E:"
$EWDK_BUILD_ENV = "$EWDK_MOUNT_LETTER\BuildEnv\SetupBuildEnv.cmd"

Write-Host "=== HyperDbg 统一工程构建脚本 ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "脚本目录: $SCRIPT_DIR" -ForegroundColor Green
Write-Host ""

$SOLOD_PATH = "so"
if (Get-Command so -ErrorAction SilentlyContinue) {
    $SOLOD_PATH = "so"
} elseif (Test-Path "$SCRIPT_DIR\solod\so.exe") {
    $SOLOD_PATH = "$SCRIPT_DIR\solod\so.exe"
} elseif (Test-Path "$SCRIPT_DIR\solod\cmd\so\main.go") {
    Push-Location "$SCRIPT_DIR\solod"
    go build -o so.exe ./cmd/so
    Pop-Location
    $SOLOD_PATH = "$SCRIPT_DIR\solod\so.exe"
}

$SOLOD_TEMP_DIR = $null
$SOLOD_GEN_DIR = Join-Path $SCRIPT_DIR "solod-gen"

if (Test-Path "$SCRIPT_DIR\hyperdbg-driver") {
    Write-Host "=== 使用 Solod 从 Go 源码生成驱动 C 代码 ===" -ForegroundColor Cyan
    $SOLOD_GEN_DIR = Join-Path "$SCRIPT_DIR\hyperdbg-driver" "gen"
    New-Item -ItemType Directory -Force -Path $SOLOD_GEN_DIR | Out-Null
    Push-Location "$SCRIPT_DIR\hyperdbg-driver"
    try {
        & $SOLOD_PATH translate -o "$SOLOD_GEN_DIR" .
        if ($LASTEXITCODE -ne 0) {
            throw "Solod 翻译失败"
        }
        $SOLOD_TEMP_DIR = $SOLOD_GEN_DIR
        Write-Host "Solod 驱动代码生成完成 (目录: $SOLOD_GEN_DIR)" -ForegroundColor Green
    } catch {
        Write-Host "=== Solod 翻译失败，跳过 Go 驱动代码 ===" -ForegroundColor Yellow
        $SOLOD_TEMP_DIR = $null
    }
    Pop-Location
    Write-Host ""
}

# 检查ISO文件是否存在
if (-not (Test-Path $EWDK_ISO_PATH)) {
    Write-Host "错误: 未找到EWDK ISO文件: $EWDK_ISO_PATH" -ForegroundColor Red
    exit 1
}

# 检查是否已挂载
$mounted = $false
try {
    $volume = Get-Volume -DriveLetter $EWDK_MOUNT_LETTER.Replace(":", "") -ErrorAction SilentlyContinue
    if ($volume -and $volume.DriveType -eq "CD-ROM") {
        Write-Host "EWDK ISO 已挂载到 $EWDK_MOUNT_LETTER" -ForegroundColor Green
        $mounted = $true
    }
} catch {}

if (-not $mounted) {
    Write-Host "正在挂载 EWDK ISO..." -ForegroundColor Yellow
    $diskImage = Mount-DiskImage -ImagePath $EWDK_ISO_PATH -PassThru
    $driveLetter = ($diskImage | Get-Volume).DriveLetter
    Write-Host "EWDK ISO 已挂载到 ${driveLetter}:" -ForegroundColor Green
}

# 检查构建环境脚本是否存在
if (-not (Test-Path $EWDK_BUILD_ENV)) {
    Write-Host "错误: 未找到EWDK构建环境脚本: $EWDK_BUILD_ENV" -ForegroundColor Red
    exit 1
}

Write-Host "使用EWDK构建环境: $EWDK_BUILD_ENV" -ForegroundColor Green
Write-Host ""

# 设置WDK环境变量
$env:WDKContentRoot = $WDK_PATH
Write-Host "设置WDKContentRoot: $env:WDKContentRoot" -ForegroundColor Green
Write-Host ""

# 清理并创建构建目录
if (Test-Path $BUILD_DIR) {
    Write-Host "清理构建目录: $BUILD_DIR" -ForegroundColor Yellow
    Remove-Item -Path $BUILD_DIR -Recurse -Force
}
Write-Host "创建构建目录: $BUILD_DIR" -ForegroundColor Yellow
New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null

if ($SOLOD_TEMP_DIR -and (Test-Path $SOLOD_TEMP_DIR)) {
    Write-Host "复制 Solod 生成的驱动代码到构建目录..." -ForegroundColor Yellow
    Copy-Item -Path "$SOLOD_TEMP_DIR\*" -Destination "$BUILD_DIR\hyperdbg-driver" -Recurse -Force
    Write-Host "Solod 驱动代码已复制到: $BUILD_DIR\hyperdbg-driver" -ForegroundColor Green
    Write-Host ""
}

# 生成项目
Write-Host "=== 生成项目 ===" -ForegroundColor Cyan
Push-Location $BUILD_DIR
try {
    cmd /c "`"$EWDK_BUILD_ENV`" && cmake `"$SCRIPT_DIR`""
    if ($LASTEXITCODE -ne 0) {
        throw "生成项目失败"
    }
} catch {
    Write-Host "=== 生成项目失败 ===" -ForegroundColor Red
    Pop-Location
    exit 1
}

Write-Host ""
Write-Host "=== 编译项目 ===" -ForegroundColor Cyan
try {
    cmd /c "`"$EWDK_BUILD_ENV`" && cmake --build . --config $CONFIG 2>&1"
    if ($LASTEXITCODE -ne 0) {
        throw "编译失败"
    }
} catch {
    Write-Host "=== 编译失败 ===" -ForegroundColor Red
    Pop-Location
    exit 1
}

Pop-Location

Write-Host ""
Write-Host "=== 构建完成 ===" -ForegroundColor Green
Write-Host ""
Write-Host "输出目录: $BUILD_DIR\$CONFIG" -ForegroundColor Yellow
Write-Host ""

# 显示生成的文件
$sysFiles = Get-ChildItem -Path "$BUILD_DIR\$CONFIG\*.sys" -ErrorAction SilentlyContinue
$pdbFiles = Get-ChildItem -Path "$BUILD_DIR\$CONFIG\*.pdb" -ErrorAction SilentlyContinue

if ($sysFiles) {
    Write-Host "生成的驱动文件:" -ForegroundColor Yellow
    $sysFiles | ForEach-Object { Write-Host "  $($_.Name)" }
}

if ($pdbFiles) {
    Write-Host "生成的符号文件:" -ForegroundColor Yellow
    $pdbFiles | ForEach-Object { Write-Host "  $($_.Name)" }
}

# ============================================
# 签名驱动
# ============================================
Write-Host ""
Write-Host "=== 签名驱动 ===" -ForegroundColor Cyan

foreach ($sysFile in $sysFiles) {
    signtool sign /fd SHA256 /s My /n "HyperDbgTest" /t http://timestamp.digicert.com $sysFile.FullName
    Write-Host "已签名: $($sysFile.Name)" -ForegroundColor Green
}

Write-Host ""
Write-Host "=== 全部完成 ===" -ForegroundColor Green
Write-Host ""
