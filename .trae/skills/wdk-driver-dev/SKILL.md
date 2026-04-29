---
name: "wdk-driver-dev"
description: "Go-based Windows kernel driver development workflow using Solod transpiler. Invoke when working on WDK bindings, so translate, EWDK compilation, or hyperdbg-driver."
---

# WDK Driver Development: Go → C → .sys

## ⚠️ GOLDEN RULES

### 🔴 NEVER EDIT GENERATED FILES
**`gen/` 目录下的所有文件都是自动生成的，绝对不能手动编辑！**

包括但不限于：
- `gen/main.c`, `gen/main.h` — Solod 转译输出
- `gen/builtin_kernel.h`, `gen/builtin_kernel.c` — 运行时支持
- `gen/wdk.c`, `gen/wdk.h` — WDK 桥接文件
- `gen/CMakeLists.txt`, `gen/build.bat` — 构建配置
- `gen/*.asm` — 汇编文件

### 🟢 WHEN ERRORS OCCUR — FIX THE GENERATOR
遇到编译/链接错误时，**必须修改生成器源码**，然后重新生成：

| 错误类型 | 修改位置 | 重新生成命令 |
|---------|---------|-------------|
| Include 找不到 (C1083) | `solod/internal/compiler/builtin.go` 或 Solod 转译器 | `go build && so translate` |
| 类型定义错误 (C2059/C2061) | `solod/internal/compiler/builtin/builtin_kernel.h` | `go build && so translate` |
| 函数名不匹配 | `solod/internal/compiler/builtin.go` (fixAsmFunctions) | `go build && so translate` |
| CMake 配置错误 | `solod/internal/compiler/builtin/CMakeLists.txt` | `go build && so translate` |
| WDK 桥接缺失 | `solod/internal/compiler/builtin/wdk.h` 或 `wdk.c` | `go build && so translate` |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│  Go Source (hyperdbg-driver/)                           │
│    ├── driver.go, ioctl.go, loader.go, types.go, ...   │
│    └── import "solod.dev/so/wdk"  ← WDK Go 绑定包      │
└──────────────────────┬──────────────────────────────────┘
                       ▼ so translate (Solod 转译器)
┌─────────────────────────────────────────────────────────┐
│  gen/ (自动生成，不要编辑!)                              │
│    ├── main.c / main.h                                  │  ← Solod 从 Go 转译
│    ├── builtin_kernel.h / builtin_kernel.c              │  ← 从 solod/builtin/ 嵌入输出
│    ├── wdk.c / wdk.h                                    │  ← 从 solod/builtin/ 嵌入输出(桥接)
│    ├── CMakeLists.txt / build.bat                       │  ← 从 solod/builtin/ 嵌入输出
│    └── *.asm                                             │  ← 从 hyperdbg-driver/asm/ 复制
└──────────────────────┬──────────────────────────────────┘
                       ▼ EWDK (ewdk.cmake + km_sys)
                    HyperDbgDriver.sys
```

## Generator Source Files (这些才是你该编辑的)

### 层级1: Solod 转译器核心 (`solod/`)
| 文件 | 作用 |
|------|------|
| `internal/compiler/translate.go` | 主转译流程：Go → C |
| `internal/compiler/builtin.go` | **后处理**: 写 builtin、写 wdk/asm、修复函数名、修复 include |

### 层级2: 嵌入模板 (`solod/internal/compiler/builtin/`, go:embed)
| 文件 | 输出到 gen/ | 当前状态 |
|------|-----------|---------|
| `builtin_kernel.h` | `gen/builtin_kernel.h` | ✅ 已加 `<stdint.h>` 修复 typedef |
| `builtin_kernel.c` | `gen/builtin_kernel.c` | ✅ 用 `<ntifs.h>` 替代 `<stdarg.h>` |
| `CMakeLists.txt` | `gen/CMakeLists.txt` | ✅ 使用 `ewdk.cmake` + `km_sys()` |
| `build.bat` | `gen/build.bat` | ✅ NMake 构建脚本 |
| `wdk.h` | `gen/wdk.h` | ✅ 嵌入模板: `#include <ntifs.h>` + `#include "builtin_kernel.h"` |
| `wdk.c` | `gen/wdk.c` | ✅ 嵌入模板: `#include "wdk.h"` |

### 层级3: WDK Go 绑定生成器 (`winmd/go-winapi-gen/cmd/wdkgen-solod/`)
| 文件 | 输出位置 | 作用 |
|------|---------|------|
| `main.go` | → `solod/so/wdk/driver.go` + `types.go` + `constants.go` + `functions.go` | **从 winmd 元数据生成 Go WDK 绑定包** |
| **输入**: `assets/Windows.Wdk.winmd` | | Windows WDK 元数据 |
| **触发**: `go build -o wdkgen-solod.exe . && ./wdkgen-solod.exe` | | 重新生成 Go 绑定 |

> ⚠️ `wdkgen-solod` 只生成 **Go 代码**到 `solod/so/wdk/` 包。它不直接生成 C。
> C 代码由 Solod 转译器从 Go（包括 wdk 包）转译而来。

### 层级4: ASM 源文件 (运行时复制)
| 源位置 | 复制到 | 由谁处理 |
|--------|-------|---------|
| `hyperdbg-driver/asm/*.asm` | `gen/*.asm` | `copyAsmFiles()` |

## Key Functions in builtin.go (后处理管线)

```go
// translate.go 主流程调用顺序:
// 1. writeBuiltin(outDir, kernelMode)     → 输出 builtin_kernel.h/.c
// 2. writeWdkBindings(outDir)             → 输出 wdk.c/wdk.h (从嵌入模板)
// 3. copyAsmFiles(entry, outDir)          → 复制 *.asm
// 4. copySoLibHeaders(outDir)             → (预留，当前空实现)
// 5. fixTypeOrder(outDir)                 → 重排类型定义解决前向引用
// 6. fixAsmFunctions(outDir)              → 修复函数名前缀
// 7. writeCMake(outDir)                   → 输出 CMakeLists.txt + build.bat

// writeBuiltin(kernelMode=true):
//   输出 builtin_kernel.h → gen/builtin_kernel.h (保持原名!)
//   输出 builtin_kernel.c → gen/builtin_kernel.c (保持原名!)

// writeWdkBindings():
//   从 cmakeFS 嵌入读取 builtin/wdk.h → gen/wdk.h
//   从 cmakeFS 嵌入读取 builtin/wdk.c → gen/wdk.c

// fixAsmFunctions():
//   main.h: main_AsmXxx → AsmXxx, main_ 前缀移除
//   main.c: ASM 函数体替换为 extern 声明
```

## Build System

### 构建系统: ewdk.cmake
项目使用自定义 `d:/ewdk/ewdk.cmake` 替代传统 `FindWdk.cmake`，提供以下函数：
- **`km_sys(target ...)`** — 编译内核驱动 (.sys)，替代 `wdk_add_driver`
- **`km_lib(target ...)`** — 编译内核静态库，替代 `wdk_add_library`
- **`um_exe(target ...)`** — 编译用户态可执行文件
- **`um_lib(target ...)`** — 编译用户态库

### 内核库命名
| 库名 | 作用 |
|------|------|
| `zydis_kernel` | Zydis 反汇编引擎 |
| `Hyperhv` | 虚拟化管理器 |
| `Hyperlog` | 日志系统 |
| `Kdserial` | 内核调试串口 |
| `Hyperevade` | 反检测 |
| `Hypertrace` | 追踪系统 |
| `Platform` | 平台抽象层 |

## Build Commands

```powershell
# 完整构建流程 (3步) — ⚠️ 必须按顺序执行！

# 步骤0: Go 语法验证和翻译 (翻译前必做!)
cd "d:\New\New folder\wdkgo"
go vet ./...
go test -v -run TestTranslate -count=1 -timeout 300s .

# 步骤1: EWDK 编译 (C → .sys)
cd "d:\New\New folder\wdkgo\hyperdbg-driver\gen"; build.bat

# 仅重新生成 WDK Go 绑定包 (当 winmd 变更时)
cd "d:\New\New folder\wdkgo\winmd\go-winapi-gen\cmd\wdkgen-solod"
go build -o wdkgen-solod.exe .; .\wdkgen-solod.exe
# 生成后也要验证 wdk 包:
cd "d:\New\New folder\wdkgo\solod"; go fix ./so/wdk/...; go vet ./so/wdk/...
```

### ⚠️ 翻译前验证原则
**在运行翻译测试之前，必须确保 Go 源码本身是正确的！**
1. `go vet ./...` — 静态分析检查
2. `go test -run TestTranslate` — 运行翻译测试，自动调用 `compiler.Translate`

如果 Go 代码本身有语法错误，翻译器生成的 C 代码必然不可靠。
**同样，WDK 绑定生成器 (`wdkgen-solod`) 运行后也必须验证生成的 `so/wdk` 包。**

## Error Fix History (供参考)

| # | 错误 | 根因 | 修复位置 | 状态 |
|---|------|------|---------|------|
| 1 | `WDK::kmLib not found` | CMake 保留目标名 | `FindWdk.cmake` 删除 meta-target | ✅ |
| 2 | `builtin_kernel.h not found` | writeBuiltin 重命名冲突 | `builtin.go` 保持原名输出 | ✅ |
| 3 | C2059 `so_byte` 语法错误 | 缺少 `<stdint.h>` | `builtin_kernel.h` 加 include | ✅ |
| 4 | `atomic.h not found` | main.h 不再引用 atomic.h | 转译器已修复 | ✅ |
| 5 | wdk.h/wdk.c 缺失 | 硬编码外部路径不存在 | 改为嵌入模板 `builtin/wdk.h` `wdk.c` | ✅ |
| 6 | CMake 构建系统过时 | `FindWdk.cmake` + `wdk_add_driver` | 改用 `ewdk.cmake` + `km_sys` | ✅ |

## Progress Tracker

### Completed
- [x] `writeBuiltin()` 保持 builtin_kernel.h 原名输出
- [x] `builtin_kernel.h` 加 `<stdint.h>` 修复 C2059
- [x] `builtin_kernel.c` 用 `<ntifs.h>` 替代 `<stdarg.h>`
- [x] 技能文件建立 Golden Rules（不编辑生成文件）
- [x] `atomic.h` 问题已解决（转译器不再生成该引用）
- [x] wdk.h/wdk.c 改为嵌入模板，不再依赖外部路径
- [x] CMakeLists.txt 改用 `ewdk.cmake` + `km_sys`
- [x] build.ps1 → build.bat (NMake 构建脚本)
- [x] 库命名更新 (zydis_kernel, Hyperhv 等)
- [x] 主 CMakeLists.txt 同步 hypedbg 项目（含用户态目标）

### Pending
- [ ] hyperevade.go, hypertrace.go, kdserial.go, script_eval.go
- [ ] Full EWDK 编译测试
- [ ] Driver 签名和加载测试
