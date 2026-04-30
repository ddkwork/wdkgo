# wdkgo

Go-based Windows kernel driver development using Solod transpiler.

## 项目结构

```
wdkgo/
├── hyperdbg-driver/     # HyperDbg 内核驱动 (Go 版本)
│   ├── asm/             # 汇编文件 (源)
│   ├── gen/             # 生成的 C 代码
│   ├── driver.go        # 驱动入口
│   ├── hypervisor.go    # 虚拟机管理
│   ├── ept.go           # EPT 钩子
│   ├── events.go        # 事件分发
│   ├── hooks.go         # 钩子管理
│   ├── ioctl.go         # IOCTL 处理
│   ├── memory.go        # 内存管理
│   ├── vmx.go           # VMX 操作
│   └── ...
├── solod/               # Go → C 转译器
│   ├── clang/           # 代码生成器
│   ├── compiler/        # 编译器核心
│   │   ├── builtin/    # 内置 C 头文件
│   │   ├── translate.go
│   │   └── ...
│   └── so/              # Go 兼容层 (wdk)
│       └── wdk/         # Windows Driver Kit 绑定
├── .github/workflows/   # CI/CD
└── build.bat            # 构建脚本
```

## 快速开始

### 1. 运行翻译

```powershell
go test -v -run TestTranslate -count=1 -timeout 300s .
```

### 2. 编译驱动

```powershell
cd hyperdbg-driver/gen
.\build.bat
```

### 3. 签名驱动

驱动的签名需要在 CI/CD 或启用了测试签名的环境中完成。

## 开发指南

### 代码组织

- `hyperdbg-driver/*.go` - 驱动的主要 Go 源代码
- `solod/compiler/builtin/` - 生成的 C 代码需要的内置头文件
- `solod/so/wdk/` - Windows Driver Kit 的 Go 绑定

### 添加新功能

1. 在 `hyperdbg-driver/` 中编写 Go 代码
2. 运行 `go test -v -run TestTranslate` 生成 C 代码
3. 在 `gen/` 目录中编译并测试

## 依赖

- Go 1.26+
- EWDK (Enterprise Windows Driver Kit)
- Visual Studio Build Tools

## 参考

- [HyperDbg](https://github.com/HyperDbg/hyperdbg) - 原始 C 实现
- [Solod](https://github.com/solod-dev/solod) - Go → C 转译器
- [WDK](https://docs.microsoft.com/en-us/windows-hardware/drivers/) - Windows Driver Kit
