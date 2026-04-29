# HyperDbg Driver - Module Test Directory Structure

This directory contains isolated test environments for each major module of the HyperDbg driver. Each subdirectory is designed to test the translation and compilation of specific components independently.

## Directory Structure

```
tests/
├── vmm/       # Virtual Machine Monitor core tests
│   ├── main.go
│   └── README.md
├── ept/       # Extended Page Tables tests
│   ├── main.go
│   └── README.md
├── hooks/     # EPT Hook mechanism tests
│   ├── main.go
│   └── README.md
├── syscall/   # Syscall hook (EFER) tests
│   ├── main.go
│   └── README.md
└── vmx/       # VMX instruction tests
    ├── main.go
    └── README.md
```

## Usage

### Test a Single Module

1. Navigate to the module directory:
```powershell
cd tests/vmm
```

2. Create a minimal Go file that imports only the required dependencies:
```go
package main

//so:driver
//so:include <ntddk.h>

func DriverEntry(driver *DRIVER_OBJECT, registryPath *UNICODE_STRING) NTSTATUS {
    return STATUS_SUCCESS
}
```

3. Translate to C:
```powershell
so translate -o ./gen .
```

4. Verify compilation (requires EWDK):
```powershell
cmake --build build --config Debug
```

### Test All Modules

Use the root `build.ps1` script which will:
1. Translate all Go files in `hyperdbg-driver/` using Solod
2. Generate C code in `hyperdbg-driver/gen/`
3. Build with EWDK using CMake
4. Sign the resulting `.sys` driver

## Module Descriptions

### VMM (Virtual Machine Monitor)
- **Files**: hypervisor.go, vmcall.go, events.go
- **Tests**: VMCALL handling, VM-exit dispatch, VCPU management
- **Dependencies**: vmx.s, context.s, tlb.s

### EPT (Extended Page Tables)
- **Files**: ept.go, memory.go
- **Tests**: Page table manipulation, large page splitting, identity map
- **Dependencies**: tlb.s

### Hooks
- **Files**: hooks.go, exec_trap.go, mode_based_exec.go
- **Tests**: EPT hooks, breakpoint management, execution traps
- **Dependencies**: ept.go, vmx.s

### Syscall
- **Files**: syscall_hook.go
- **Tests**: SYSCALL/SYSRET emulation, EFER MSR manipulation
- **Dependencies**: vmx.s, context.s

### VMX
- **Files**: vmx.s, vmx.go
- **Tests**: VMXON/VMOFF, VMLAUNCH/VMRESUME, VMREAD/VMWRITE
- **Dependencies**: None (pure assembly)

## Integration Testing Workflow

1. **Unit Test**: Test each module independently
2. **Integration Test**: Combine modules that depend on each other
3. **Full Build**: Test the complete driver with all modules

### Example: Test EPT + Hooks Integration

```powershell
cd tests/ept-hooks
# Create combined test file
# Translate
# Compile
# Verify
```

## Expected Output

Each test should produce:
- ✅ Successful translation (no errors from `so translate`)
- ✅ Clean compilation (no warnings or errors from MSVC)
- ✅ Valid .sys driver file
- ✅ Successful driver load (optional, requires test machine)

## Troubleshooting

### Translation Errors
- Check for missing `//so:extern` declarations
- Verify type definitions match C struct layouts
- Ensure no Go-specific features are used (generics, interfaces, etc.)

### Compilation Errors
- Check WDK header inclusion paths
- Verify function signatures match extern declarations
- Ensure assembly files (.s) are properly linked

### Runtime Errors
- Check VMCS configuration values
- Verify EPT page table consistency
- Validate memory allocation/deallocation pairs

## Notes

- Assembly files (.s) are NOT translated by Solod; they are linked during EWDK compilation
- The `gen/` directory contains auto-generated C code; DO NOT edit manually
- Test directories can be safely deleted; they are regenerated from source
