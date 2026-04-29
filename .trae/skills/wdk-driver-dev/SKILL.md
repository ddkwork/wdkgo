---
name: "wdk-driver-dev"
description: "Go-based Windows kernel driver development workflow using Solod transpiler. Invoke when working on WDK bindings, so translate, EWDK compilation, or hyperdbg-driver."
---

# WDK Driver Development: Go → C → .sys

## ⚠️ GOLDEN RULES

### 🔴 NEVER EDIT GENERATED FILES
**`gen/` 目录下的所有文件都是自动生成的，绝对不能手动编辑！**

包括但不限于：
- `gen/main.c`, `gen/main.h`, `gen/atomic.h`, `gen/slices.h` — Solod 转译输出
- `gen/builtin_kernel.h`, `gen/builtin_kernel.c` — 运行时支持
- `gen/wdk.c`, `gen/wdk.h` — WDK 桥接文件
- `gen/CMakeLists.txt`, `gen/FindWdk.cmake`, `gen/build.ps1` — 构建配置
- `gen/*.asm` — 汇编文件

### 🟢 WHEN ERRORS OCCUR — FIX THE GENERATOR
遇到编译/链接错误时，**必须修改生成器源码**，然后重新生成：

| 错误类型 | 修改位置 | 重新生成命令 |
|---------|---------|-------------|
| Include 找不到 (C1083) | `solod/internal/compiler/builtin.go` 或 Solod 转译器 | `go build && so translate` |
| 类型定义错误 (C2059/C2061) | `solod/internal/compiler/builtin/builtin_kernel.h` | `go build && so translate` |
| 函数名不匹配 | `solod/internal/compiler/builtin.go` (fixAsmFunctions) | `go build && so translate` |
| CMake 配置错误 | `solod/internal/compiler/builtin/CMakeLists.txt` 或 `FindWdk.cmake` | `go build && so translate` |

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
│    ├── main.c / main.h / atomic.h / slices.h            │  ← Solod 从 Go 转译
│    ├── builtin_kernel.h / builtin_kernel.c              │  ← 从 solod/builtin/ 嵌入输出
│    ├── wdk.c / wdk.h                                    │  ← 从 build/hyperdbg-driver/wdk/ 复制(桥接)
│    └── *.asm                                             │  ← 从 hyperdbg-driver/asm/ 复制
└──────────────────────┬──────────────────────────────────┘
                       ▼ EWDK (cmake + MSVC)
                    HyperDbgDriver.sys
```

## Generator Source Files (这些才是你该编辑的)

### 层级1: Solod 转译器核心 (`solod/`)
| 文件 | 作用 |
|------|------|
| `internal/compiler/translate.go` | 主转译流程：Go → C |
| `internal/compiler/builtin.go` | **后处理**: 写 builtin、复制 wdk/asm、修复函数名、修复 include |

### 层级2: 嵌入模板 (`solod/internal/compiler/builtin/`, go:embed)
| 文件 | 输出到 gen/ | 当前状态 |
|------|-----------|---------|
| `builtin_kernel.h` | `gen/builtin_kernel.h` | ✅ 已加 `<stdint.h>` 修复 typedef |
| `builtin_kernel.c` | `gen/builtin_kernel.c` | ✅ 用 `<ntifs.h>` 替代 `<stdarg.h>` |
| `CMakeLists.txt` | `gen/CMakeLists.txt` | ✅ 使用 `wdk_add_driver()` |
| `FindWdk.cmake` | `gen/FindWdk.cmake` | ✅ 移除无效 `WDK::kmLib` |
| `build.ps1` | `gen/build.ps1` | 构建脚本模板 |

### 层级3: WDK Go 绑定生成器 (`winmd/go-winapi-gen/cmd/wdkgen-solod/`)
| 文件 | 输出位置 | 作用 |
|------|---------|------|
| `main.go` | → `solod/so/wdk/driver.go` + `types.go` + `constants.go` + `functions.go` | **从 winmd 元数据生成 Go WDK 绑定包** |
| **输入**: `assets/Windows.Wdk.winmd` | | Windows WDK 元数据 |
| **触发**: `go build -o wdkgen-solod.exe . && ./wdkgen-solod.exe` | | 重新生成 Go 绑定 |

> ⚠️ `wdkgen-solod` 只生成 **Go 代码**到 `solod/so/wdk/` 包。它不直接生成 C。
> C 代码由 Solod 转译器从 Go（包括 wdk 包）转译而来。

### 层级4: WDK 桥接文件 (非嵌入，运行时复制)
| 源位置 | 复制到 | 内容 |
|--------|-------|------|
| `build/hyperdbg-driver/wdk/wdk.h` | `gen/wdk.h` | `#include <ntifs.h>` + `#include "builtin_kernel.h"` |
| `build/hyperdbg-driver/wdk/wdk.c` | `gen/wdk.c` | `#include "wdk.h"` |
| 由 `writeWdkBindings()` 处理 | | |

### 层级5: ASM 源文件 (运行时复制)
| 源位置 | 复制到 | 由谁处理 |
|--------|-------|---------|
| `hyperdbg-driver/asm/*.asm` | `gen/*.asm` | `copyAsmFiles()` |

## Key Functions in builtin.go (后处理管线)

```go
// translate.go 主流程调用顺序:
// 1. writeBuiltin(outDir, kernelMode)     → 输出 builtin_kernel.h/.c
// 2. writeWdkBindings(outDir)             → 复制 wdk.c/wdk.h
// 3. copyAsmFiles(entry, outDir)          → 复制 *.asm
// 4. fixAsmFunctions(outDir)              → 修复函数名前缀

// writeBuiltin(kernelMode=true):
//   输出 builtin_kernel.h → gen/builtin_kernel.h (保持原名!)
//   输出 builtin_kernel.c → gen/builtin_kernel.c (保持原名!)

// fixAsmFunctions():
//   main.h: main_AsmXxx → AsmXxx, main_ 前缀移除
//   main.c: ASM 函数体替换为 extern 声明
```

## Build Commands

```powershell
# 完整构建流程 (4步) — ⚠️ 必须按顺序执行！

# 步骤0: Go 语法验证和测试 (翻译前必做!)
cd "c:\Users\Admin\Desktop\New folder\hyperdbg-driver"
go fix ./...
go vet ./...
go test ./...          # 单元测试必须通过，否则生成的 C 代码不可靠

# 步骤1: 构建 Solod 转译器
cd "c:\Users\Admin\Desktop\New folder\solod"; go build -o ../solod.exe ./cmd/so

# 步骤2: 运行翻译器 (Go → C)
cd "c:\Users\Admin\Desktop\New folder"; .\solod.exe translate "c:\Users\Admin\Desktop\New folder\hyperdbg-driver"

# 步骤3: EWDK 编译 (C → .sys)
cd "c:\Users\Admin\Desktop\New folder\hyperdbg-driver\gen"; powershell -ExecutionPolicy Bypass -File build.ps1

# 仅重新生成 WDK Go 绑定包 (当 winmd 变更时)
cd "c:\Users\Admin\Desktop\New folder\winmd\go-winapi-gen\cmd\wdkgen-solod"
go build -o wdkgen-solod.exe .; .\wdkgen-solod.exe
# 生成后也要验证 wdk 包:
cd "c:\Users\Admin\Desktop\New folder\solod"; go fix ./so/wdk/...; go vet ./so/wdk/...
```

### ⚠️ 翻译前验证原则
**在运行 `so translate` 之前，必须确保 Go 源码本身是正确的！**
1. `go fix ./...` — 自动修复过时语法
2. `go vet ./...` — 静态分析检查
3. `go test ./...` — 单元测试通过

如果 Go 代码本身有语法错误或测试失败，翻译器生成的 C 代码必然不可靠。
**同样，WDK 绑定生成器 (`wdkgen-solod`) 运行后也必须验证生成的 `so/wdk` 包。**

## Error Fix History (供参考)

| # | 错误 | 根因 | 修复位置 | 状态 |
|---|------|------|---------|------|
| 1 | `WDK::kmLib not found` | CMake 保留目标名 | `FindWdk.cmake` 删除 meta-target | ✅ |
| 2 | `builtin_kernel.h not found` | writeBuiltin 重命名冲突 | `builtin.go` 保持原名输出 | ✅ |
| 3 | C2059 `so_byte` 语法错误 | 缺少 `<stdint.h>` | `builtin_kernel.h` 加 include | ✅ |
| 4 | `atomic.h not found` | Solod 转译生成但缺失 | **待排查** | 🔴 当前 |

## Current Error: `atomic.h` not found

```
error C1083: Cannot open include file: 'atomic.h': No such file or directory
(main.h:11)
```

`atomic.h` 和 `slices.h` 是 Solod 转译时为 Go 包生成的辅助头文件。需要检查：
1. Solod 是否正确生成了这些文件到 `gen/`
2. 如果没有，是否需要在 `builtin.go` 中补充生成逻辑

## Progress Tracker

### Completed
- [x] `writeBuiltin()` 保持 builtin_kernel.h 原名输出
- [x] `FindWdk.cmake` 移除无效 `WDK::kmLib`
- [x] `builtin_kernel.h` 加 `<stdint.h>` 修复 C2059
- [x] `builtin_kernel.c` 用 `<ntifs.h>` 替代 `<stdarg.h>`
- [x] 技能文件建立 Golden Rules（不编辑生成文件）

### 🔴 In Progress — 修复 `atomic.h` / `slices.h` 缺失
- [ ] 排查 Solod 为何未生成 `atomic.h`、`slices.h` 到 `gen/`
- [ ] 修复生成逻辑确保所有依赖头文件存在
- [ ] 验证编译通过

### Pending
- [ ] hyperevade.go, hypertrace.go, kdserial.go, script_eval.go
- [ ] Full EWDK 编译测试
- [ ] Driver 签名和加载测试
