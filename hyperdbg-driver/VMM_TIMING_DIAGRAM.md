# HyperDbg VMM 启动时序图

## 1. 整体架构概览

```
┌─────────────────────────────────────────────────────────────────────┐
│                        User Mode (Ring3)                            │
│  ┌─────────────┐   IOCTL    ┌──────────────┐                       │
│  │ Debugger.exe │ ─────────→ │ HyperDbg.sys  │                      │
│  └─────────────┘            └──────┬───────┘                        │
└──────────────────────────────────────┼──────────────────────────────┘
                                      │ IRP_MJ_DEVICE_CONTROL
┌──────────────────────────────────────▼──────────────────────────────┐
│                        Kernel Mode (Ring0)                          │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    Hypervisor (Go)                           │   │
│  │                                                              │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │   │
│  │  │ VCPU[0]  │  │ VCPU[1]  │  │ VCPU[2]  │  │ VCPU[N]  │     │   │
│  │  │ VMCS/EPT │  │ VMCS/EPT │  │ VMCS/EPT │  │ VMCS/EPT │     │   │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘     │   │
│  │       │             │             │             │           │   │
│  │  ┌────▼─────────────▼─────────────▼─────────────▼─────┐     │   │
│  │  │              EventDispatcher                       │     │   │
│  │  │  CPUID | MSR | IO | CR | EPT | Exception          │     │   │
│  │  └────────────────────────┬──────────────────────────┘     │   │
│  │                           │                                 │   │
│  │  ┌────────────────────────▼──────────────────────────┐     │   │
│  │  │                HookManager                         │     │   │
│  │  │  EptHook[] → Breakpoint[] → FakePage              │     │   │
│  │  └───────────────────────────────────────────────────┘     │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │
│  │ Logger   │  │ Evasion  │  │ Tracer   │  │ KdSerial         │   │
│  │ RingBuf  │  │ Anti-DBG │  │ LBR/BTS  │  │ COM Port         │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## 2. 驱动加载流程

```
DriverEntry(driverObj, registryPath)
│
├── ExInitializeDriverRuntime(DrvRtPoolNxOptIn)
│
├── InitializeGlobals(numCpus)
│   │
│   ├── NewLogger(MAX_LOG_BUFFERS, MAX_LOG_BUFFERS_PRIO)
│   │   └── LogBuffer{Messages[1000], Read/Write Index, Spinlock}
│   │
│   ├── NewHypervisor(numCpus)
│   │   ├── vcpus[numCpus]
│   │   ├── eptTables[numCpus]
│   │   ├── events = newEventDispatcher()
│   │   ├── hooks = NewHookManager(40)
│   │   ├── logger = newLoggerInternal()
│   │   └── memory = newMemoryManager()
│   │
│   ├── gHyp.Initialize()  ◄═════════════════════════╗
│   │   └── for coreId := 0..numCpus:               │
│   │       └── initializeVcpu(coreId)               │
│   │           │                                   │
│   │           ├─① allocateVmcsRegion()            │
│   │           │   vmcsSize = IA32_VMX_BASIC & 0x1FFF│
│   │           │   ptr = AllocateZeroedPool(4KB)    │
│   │           │   pa = MmGetPhysicalAddress(ptr)    │
│   │           │   return pa                         │
│   │           │                                    │
│   │           ├─② allocateVmxonRegion()           │
│   │           │   revisionId = IA32_VMX_BASIC      │
│   │           │   *ptr = revisionId (写入RevisionID)│
│   │           │   return MmGetPhysicalAddress(ptr)   │
│   │           │                                    │
│   │           ├─③ allocateMsrBitmap()             │
│   │           │   bitmap[0x1000] 全部设为0xFF       │
│   │           │   清除 0xC00~0xC08, 0x174 的位     │
│   │           │   return pa                         │
│   │           │                                    │
│   │           ├─④ allocateIoBitmap() ×2           │
│   │           │   A: 0x0000~0x7FFF 端口范围        │
│   │           │   B: 0x8000~0xFFFF 端口范围        │
│   │           │   return pa                         │
│   │           │                                    │
│   │           ├─⑤ allocateVmmStack()              │
│   │           │   size = 64KB                      │
│   │           │   stackTop = ptr + size - 16       │
│   │           │   *stackTop = 0xDEAD0BADDEAD0BAD   │
│   │           │   return stackTop PA                │
│   │           │                                    │
│   │           ├─⑥ BuildEptPointer()               │
│   │           │   eptTable.NewEptTable()            │
│   │           │   eptTable.initIdentityMap()       │
│   │           │   eptp = PML4_PA | MemType(6)|RWX  │
│   │           │   v.EptPointer.AsUInt = eptp       │
│   │           │                                    │
│   │           ├─⑦ vmxTurnOn(vmxonPa)  ◄══════════╝
│   │           │   cr4 = __readmsr(IA32_CR4)        │
│   │           │   if !(cr4 & VMXE):                 │
│   │           │       __writemsr(IA32_CR4, cr4|VMXE)│
│   │           │   AsmVmxVmxOn(vmxonPhysicalAddr)     │
│   │           │   ┌─────────────────────────────┐   │
│   │           │   │  CPU 进入 VMX Operation     │   │
│   │           │   │  CR4.VMXE=1 ✓              │   │
│   │           │   └─────────────────────────────┘   │
│   │           │                                    │
│   │           ├─⑧ vmClear(vmcsPa)                  │
│   │           │   AsmVmxVmxClear(vmcsPhysicalAddr)   │
│   │           │   ┌─────────────────────────────┐   │
│   │           │   │  VMCS 状态初始化为 Clear    │   │
│   │           │   │  准备接受 VMPTRLD          │   │
│   │           │   └─────────────────────────────┘   │
│   │           │                                    │
│   │           ├─⑨ vmLoad(vmcsPa)                   │
│   │           │   AsmVmxVmxPtrld(vmcsPhysicalAddr)   │
│   │           │   ┌─────────────────────────────┐   │
│   │           │   │  加载当前 VMCS 指针         │   │
│   │           │   │  后续 VMWRITE 写入此 VMCS  │   │
│   │           │   └─────────────────────────────┘   │
│   │           │                                    │
│   │           └─⑩ setupVmcs(v)  ◄═════════════════╝
│   │               │
│   │               ╔══ Guest State Area ════════════╗
│   │               ║                                ║
│   │               ║  vmWrite(GuestCR0, msr(CR0))  ║
│   │               ║  vmWrite(GuestCR3, msr(CR3))  ║
│   │               ║  vmWrite(GuestCR4, msr(CR4))  ║
│   │               ║  vmWrite(GuestDR7, 0x400)      ║
│   │               ║  vmWrite(GuestRSP, 0)          ║
│   │               ║  vmWrite(GuestRIP, 0)          ║
│   │               ║  vmWrite(GuestRFLAGS, 0x002)   ║
│   │               ║                                ║
│   │               ╠══ Segment Registers ════════════╣
│   │               ║                                ║
│   │               ║  CS: Sel=0x8 Base=0 Lim=FFFFFFFF║
│   │               ║      AR=A0FB (Code/Ring0/L)    ║
│   │               ║                                ║
│   │               ║  DS/ES/FS/GS/SS:               ║
│   │               ║      Sel=0x10 Base=0 Lim=FFFFFFFF║
│   │               ║      AR=C0F3 (Data/Ring0/RW)   ║
│   │               ║                                ║
│   │               ║  GDTR.Base = msr(GDTR_BASE)    ║
│   │               ║  IDTR.Base = msr(IDTR_BASE)    ║
│   │               ║  LDTR.AR = 0x82 (Busy)         ║
│   │               ║  TR.AR = 0x8B (Busy 32bit TSS) ║
│   │               ║                                ║
│   │               ╚══ Control Fields ══════════════╝
│   │                   │
│   │                   ├─ EPT Pointer (EPTP)
│   │                   │  vmWrite(0x201A, eptPointer)
│   │                   │  ┌────────────────────────┐
│   │                   │  │  PML4 Physical Address  │
│   │                   │  │  Memory Type = WB (6)   │
│   │                   │  │  Page Walk Length = 4   │
│   │                   │  │  Enable Bit [0][2]      │
│   │                   │  └────────────────────────┘
│   │                   │
│   │                   ├─ Pin-Based Controls
│   │                   │  adjustVmcsControl(msr(0x481), 0x10000016)
│   │                   │  ┌────────────────────────┐
│   │                   │  │ bit0: INTs exiting      │
│   │                   │  │ bit4: NMIs exiting      │
│   │                   │  │ bit20: VT-x PT enabled  │
│   │                   │  └────────────────────────┘
│   │                   │
│   │                   ├─ Primary Proc-Based Controls
│   │                   │  adjustVmcsControl(msr(0x483), 0x04007EFA|BIT31)
│   │                   │  ┌────────────────────────┐
│   │                   │  │ bit2: IO instructions   │
│   │                   │  │ bit7: HLT exiting       │
│   │                   │  │ bit12: RDTSC exiting     │
│   │                   │  │ bit17: CR8 load/store   │
│   │                   │  │ bit20: CR3 load/store    │
│   │                   │  │ bit25: MOV DR exiting    │
│   │                   │  │ bit28: Use I/O Bitmaps   │
│   │                   │  │ bit29: Use MSR Bitmaps   │
│   │                   │  │ bit31: Activate Secondary│
│   │                   │  └────────────────────────┘
│   │                   │
│   │                   ├─ Secondary Proc-Based Controls ⭐
│   │                   │  adjustVmcsControl(msr(0x485), 0x06D8FBFE)
│   │                   │  ┌────────────────────────┐
│   │                   │  │ bit0: EPT enabled ⭐⭐⭐  │
│   │                   │  │ bit1: RDTSCP enabled     │
│   │                   │  │ bit2: x2APIC mode       │
│   │                   │  │ bit3: VPID enabled       │
│   │                   │  │ bit4: WBINVD exiting     │
│   │                   │  │ bit5: Unrestricted Guest │
│   │                   │  │ bit7: INVPCID enabled    │
│   │                   │  │ bit14: VMFUNC enabled    │
│   │                   │  │ bit15: EPT Violation #VE │
│   │                   │  │ bit22: EPT Write/Exec   │
│   │                   │  │ bit23: EPT Read/Exec    │
│   │                   │  │ bit28: PAUSE loop exit   │
│   │                   │  └────────────────────────┘
│   │                   │
│   │                   ├─ Exception Bitmap
│   │                   │  vmWrite(ExceptionBitmap, 0xFFFFFFFF)
│   │                   │  └─ 拦截所有异常到 Hypervisor
│   │                   │
│   │                   ├─ I/O Bitmap A/B Addresses
│   │                   │  vmWrite(IoBitmapA, ioPaA)
│   │                   │  vmWrite(IoBitmapB, ioPaB)
│   │                   │
│   │                   └─ MSR Bitmap Address
│   │                      vmWrite(MsrBitmap, msrBitmapPa)
│   │
│   │               ╔══ Host State Area ══════════════╗
│   │               ║                                ║
│   │               ║  Host CR0/CR3/CR4              ║
│   │               ║    从当前 CPU 直接复制           ║
│   │               ║                                ║
│   │               ║  Host Segment Selectors         ║
│   │               ║    CS=getCs() SS=getSs() ...    ║
│   │               ║    TR=0x40 (TSS Selector)       ║
│   │               ║                                ║
│   │               ║  Host FS/GS Base               ║
│   │               ║    FS = msr(FS_BASE)            ║
│   │               ║    GS = msr(GS_BASE)            ║
│   │               ║                                ║
│   │               ║  Host GDTR/IDTR Base           ║
│   │               ║  Host SYSENTER CS/ESP/EIP       ║
│   │               ║  Host Segment Access Rights     ║
│   │               ║    CS=A0FB SS/DS/ES=C0F3        ║
│   │               ║                                ║
│   │               ║  ★ Host RSP = VMM Stack Top ★  ║
│   │               ║    VM Exit 时切换到此栈         ║
│   │               ║                                ║
│   │               ║  ★★ Host RIP = vmExitHandler ★★ ║
│   │               ║    VM Exit 执行入口地址         ║
│   │               ╚═════════════════════════════════╝
│   │
│   ├── gEvents = newEventDispatcher()
│   │   └── setupDefaultEventHandlers()
│   │       ├── CpuidHandler    {隐藏虚拟化特征}
│   │       ├── MsrHandler      {拦截MSR读写}
│   │       ├── IoHandler       {I/O端口监控}
│   │       ├── CrAccessHandler {CR寄存器保护}
│   │       ├── TscHandler      {TSC偏移}
│   │       └── ExceptionHandler {断点/UD/GP/NMI}
│   │
│   ├── gHooks = NewHookManager(40)
│   ├── gEvasion = NewEvasionState()
│   ├── gTrace = NewTracerState()
│   ├── gSerial = NewKdSerialState()
│   ├── gPoolMgr = NewPoolManager(1024)
│   ├── gOptimizer = NewVmmOptimizer(OptBasic)
│   └── gDrv = &Driver{}
│
├── IoCreateDevice(\Device\HyperDbgDebuggerDevice)
├── 注册 IRP MajorFunctions:
│   ├── IRP_MJ_CREATE  → drvCreate
│   ├── IRP_MJ_CLOSE   → drvClose
│   ├── IRP_MJ_READ    → drvRead  (日志读取)
│   ├── IRP_MJ_WRITE   → drvWrite
│   └── IRP_MJ_DEVICE_CONTROL → drvIoctl (50+ IOCTL处理)
├── IoCreateSymbolicLink(DosDevices\HyperDbgDebuggerDevice)
└── device.Flags |= DO_BUFFERED_IO
```

## 3. VM Exit 处理流程

```
Guest Code Execution
        │
        ▼  VM Exit 触发
   ┌─────────┐
   │  Hardware │ ← CPU 检测到事件
   │  Save    │   自动保存 Guest State 到 VMCS
   │  State   │   切换到 Host State Area
   └────┬─────┘
        │
        ▼
   ┌─────────────┐
   │vmExitHandler│ ◄── Host RIP 入口点 (汇编实现)
   │  (Assembly) │     保存所有通用寄存器到 GUEST_REGS
   └──────┬──────┘
          │
          ▼
   ┌──────────────────┐
   │ HandleVmExit(regs)│ ← Go 实现
   └──────┬───────────┘
          │
          ├─ coreId = getCurrentProcessorNumber()
          ├─ v = &gHyp.vcpus[coreId]
          ├─ v.Regs = regs
          ├─ v.ExitReason = vmRead(ExitReason)
          ├─ v.LastVmexitRip = vmRead(GuestRip)
          └─ v.ExitQualification = vmRead(ExitQual)
          │
          ▼
   ┌──────────────────┐
   │  h.dispatch(v)    │
   └──────┬───────────┘
          │
          ├── ExitExceptionNmi ──→ handleException()
          │                          ├── VectorBP (0xCC)  → gHooks.HandleBreakpoint()
          │                          ├── VectorGP (#GP)   → injectGp / handleGPFault
          │                          ├── VectorUD (#UD)   → injectUd
          │                          └── VectorNMI        → handleNmi
          │
          ├── ExitCpuid (0x22) ──→ gEvents.Dispatch(EventCpuid)
          │                          └── CpuidHandler.Handle()
          │                              ├── Leaf 0: 返回 "GenuineIntel"
          │                              ├── Leaf 1: 清除 VMX 特征位
          │                              └── Leaf 0xA: 清空 PMU 信息
          │
          ├── ExitRdmsr (0x1F) ──→ gEvents.Dispatch(EventRdmsr)
          │                          └── MsrHandler.HandleRead()
          │                              ├── IA32_DEBUGCTL: 清除 BTF 位
          │                              └── 其他: 正常读取
          │
          ├── ExitWrmsr (0x20) ──→ gEvents.Dispatch(EventWrmsr)
          │                          └── MsrHandler.HandleWrite()
          │                              └── 同上过滤逻辑
          │
          ├── ExitIo (0x1E) ──────→ gEvents.Dispatch(EventIo)
          │                          └── IoHandler.Handle()
          │                              ├── IN 操作: onIn(port, size)
          │                              └── OUT操作: onOut(port, size)
          │
          ├── ExitMovCr (0x1C) ────→ gEvents.Dispatch(EventMovCr)
          │                          └── CrAccessHandler.Handle()
          │                              ├── CR0/CR3/CR4 读写检测
          │                              └── LMSW 转换为 CR0 写
          │
          ├── ExitRdtsc/Rdtscp ────→ gEvents.Dispatch(EventTsc)
          │                          └── TscHandler.Handle()
          │                              tsc = rdtsc() + offset
          │
          ├── ExitEptViolation (0x30) ← EPT 权限违规 ⭐
          │   └── handleEptViolation()
          │       ├── read/exec only → HandleExecHook(gpa)
          │       │                     └── 断点命中检查
          │       │                     └── 单步执行原始指令
          │       └── read/write      → HandleReadWriteHook(gpa)
          │                             └── 切换到 FakePage
          │
          ├── ExitEptMisconfig (0x31) ← EPT 配置错误
          │   └── handleEptMisconfig() → 注入 #UD
          │
          ├── ExitVmcall (0x12) ────→ handleVmcall()
          │   ├── 检查 Magic Number:
          │   │   RAX=HYPERDBG_VMCALL_MAGIC_RAX
          │   │   RCX=HYPERDBG_VMCALL_MAGIC_RCX
          │   │   RDX=HYPERDBG_VMCALL_MAGIC_RDX
          │   └── handleHyperdbgVmcall(num, p1-p4)
          │       ├── 1: Launch Debugger
          │       ├── 2: Test Command
          │       └── >0x100000000: 用户自定义
          │
          ├── ExitMtf (0x25) ──────→ handleMonitorTrapFlag() ⭐
          │   └── 单步完成后的恢复逻辑
          │       ├── gHooks.HandleMtfRestore() → 还原EPT权限
          │       ├── breakpointReapplyHook()   → 重置断点字节
          │       └── kdHandleNmiCallback()     → KD串口回调
          │
          └── Default ──────────────→ return false (不处理)
          
          │
          ▼
   ┌──────────────────┐
   │  Advance IP?      │
   │  if !VmxoffExecuted │
   │    && IncrementRip │
   └──────┬───────────┘
          │ Yes
          ▼
   ┌──────────────────┐
   │  h.advanceIp(v)   │
   │  instrLen = vmRead │
   │  (VmexitInstrLen) │
   │  rip += instrLen  │
   │  vmWrite(GuestRip)│
   └──────┬───────────┘
          │
          ▼
   ┌──────────────────┐
   │  Restore & Resume │
   │  v.OnVmxRootMode=false
   └──────┬───────────┘
          │
          ▼
   ┌─────────┐
   │VMRESUME  │ ← 回到 Guest 继续执行
   │(or VMLAUNCH)
   └─────────┘
```

## 4. EPT Hook 工作原理

```
正常执行路径:
  Guest VA → Guest CR3 → Guest Page Table → Physical Page → Execute

Hook 安装后:
  ┌─────────────────────────────────────────────────────┐
  │                    EPT Page Tables                   │
  │                                                     │
  │  PML4[0] → PDPT[0] → PD[0] ──┬→ PT[0] (Original)  │
  │                             │   Permission: RWX    │
  │                             │                      │
  │                             └→ PT[hook_idx] (Hooked)│
  │                                 Permission: R__    │  ← 去掉 Exec
  │                                 PhysAddr: FakePagePA │  ← 指向假页
  │                                                     │
  └─────────────────────────────────────────────────────┘

执行流程:

  ① Guest 执行到 Hook 地址
     │
     ▼
  ② EPT Violation! (Exec permission denied)
     │
     ▼
  ③ handleEptViolation()
     │
     ├── 读取 ExitQualification: read=0 write=0 exec=1
     ├── 读取 GuestPhysicalAddress = hooked_page_PA
     │
     ▼
  ④ HandleExecHook()
     │
     ├── FindBreakpoint(RIP) → 找到匹配的 BP
     ├── v.IncrementRip = false (不推进 IP)
     ├── v.RegisterBreakOnMtf = true
     │
     ├── restoreOriginalAndInjectBp():
     │   ├── 克隆 OriginalEntry (保存修改前的EPT项)
     │   ├── 将 EPT 项还原为 OriginalEntry (临时恢复权限)
     │   └── 在物理内存中写入 0xCC (INT3) 到断点位置
     │
     └── setMonitorTrapFlag(true) ← 启用单步陷阱
     │
     ▼
  ⑤ VMRESUME → Guest 执行 INT3 (0xCC)
     │
     ▼
  ⑥ VM Exit: Exception NMI (Vector=BP)
     │
     ▼
  ⑦ HandleBreakpoint() 或 HandleException()
     │   (通知调试器: 断点命中)
     │
     ▼
  ⑧ VMRESUME → Guest 继续执行下一条指令
     │
     ▼
  ⑨ VM Exit: Monitor Trap Flag (MTF)
     │
     ▼
  ⑩ HandleMtfRestore()
     │
     ├── 检查 MtfEptHookRestorePoint != 0
     ├── 将 EPT 项改回 ModifiedEntry (去掉 Exec 权限)
     └── setMonitorTrapFlag(false)

结果: 下次执行到同一地址时，再次触发 EPT Violation → 循环
```

## 5. IOCTL 命令分发表

```
User Mode                    Kernel Mode
    │                              │
    │ IOCTL_REQUEST                │
    ├─────────────────────────────►│
    │                              ▼
    │                        ┌─────────────┐
    │                        │  drvIoctl()  │
    │                        └──────┬──────┘
    │                               │
    │              ┌────────────────┼────────────────┐
    │              │                │                │
    │              ▼                ▼                ▼
    │     ┌────────────┐  ┌─────────────┐  ┌────────────┐
    │     │VMM Control │  │Memory Access│  │Register    │
    │     │            │  │             │  │Operations  │
    │     ├────────────┤  ├─────────────┤  ├────────────┤
    │     │INIT        │  │READ_MEM     │  │GET_REG     │
    │     │SHUTDOWN    │  │WRITE_MEM    │  │SET_REG     │
    │     │PAUSE       │  │VIRT_TO_PHYS │  │MODIFY_REGS │
    │     │RESUME      │  │PHYS_TO_VIRT │  │GET_PROCESS │
    │     │SWITCH_PROC │  │READ_WRITE   │  │BASE/CR3    │
    │     └────────────┘  └─────────────┘  └────────────┘
    │              │                │                │
    │              ▼                ▼                ▼
    │     ┌────────────┐  ┌─────────────┐  ┌────────────┐
    │     │EPT Hooks   │  │Tracing      │  │Transparent │
    │     │            │  │             │  │Mode        │
    │     ├────────────┤  ├─────────────┤  ├────────────┤
    │     │HOOK        │  │EXEC_TRACE   │  │ENABLE      │
    │     │UNHOOK      │  │TRACE_START  │  │DISABLE     │
    │     │SET_HOOK    │  │TRACE_STOP   │  │HYPERTRACE  │
    │     │GET_EPT_TBL │  │READ_BUFFER  │  │START/STOP  │
    │     └────────────┘  │CLEAR_BUFFER │  └────────────┘
    │                     └─────────────┘
    │              │                │
    │              ▼                ▼
    │     ┌─────────────────────────────────────┐
    │     │         TLB & Cache Management       │
    │     ├─────────────────────────────────────┤
    │     │INVEPT    INVVPID    FLUSH_TLB       │
    │     │GET_MTRR   CHANGE_CORE                │
    │     └─────────────────────────────────────┘
```

## 6. 内存管理架构

```
┌─────────────────────────────────────────────────────────────┐
│                     Physical Memory                          │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              EPT Identity Map (Default)              │    │
│  │                                                      │    │
│  │  GPA 0x00000000 - 0xFFFFFFFFFFFF                     │    │
│  │  └→ PA 0x00000000 - 0xFFFFFFFFFFFF (1:1 Mapping)    │    │
│  │    Permission: Read + Write + Execute                │    │
│  │    Memory Type: Write Back (from MTRR)              │    │
│  │    Large Pages: 2MB (default for most regions)      │    │
│  │                                                      │    │
│  │  ┌─────────────────────────────────────────┐        │    │
│  │  │         Split Region Example            │        │    │
│  │  │                                         │        │    │
│  │  │  Before: PD[5] -> Large Page (2MB)      │        │    │
│  │  │           PA=0x100000 RWX WB           │        │    │
│  │  │                                         │        │    │
│  │  │  After:  PD[5] -> PT_Base (4KB table)   │        │    │
│  │  │           PT[0]  -> PA=0x100000 RWX     │        │    │
│  │  │           PT[1]  -> PA=0x102000 RWX     │        │    │
│  │  │           ...                          │        │    │
│  │  │           PT[N]  -> PA=0x11F000 RW_    │ ← Hooked│    │
│  │  │                     (FakePage PA)       │        │    │
│  │  └─────────────────────────────────────────┘        │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                  Pool Manager                        │    │
│  │                                                      │    │
│  │  Allocate(size, tag='HDBG')                         │    │
│  │  └→ ExAllocatePool2(NonPaged, size, 'HDBG')         │    │
│  │                                                      │    │
│  │  Free(ptr)                                          │    │
│  │  └→ ExFreePoolWithTag(ptr, 'HDBG')                  │    │
│  │                                                      │    │
│  │  CheckAndPerformDeallocation()                      │    │
│  │  └→ 遍历链表，释放所有未释放的池                    │    │
│  │                                                      │    │
│  │  Stats:                                             │    │
│  │    TotalAllocated: XXX MB                           │    │
│  │    TotalFreed: XXX MB                               │    │
│  │    CurrentPools: XX                                  │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## 7. 关键数据结构关系图

```
┌──────────────────────────────────────────────────────────────────┐
│                        Hypervisor                                 │
│                                                                    │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  struct Hypervisor {                                       │  │
│  │    vcpus: []VCPU                                           │  │
│  │    eptTables: []*EptTable                                  │  │
│  │    events: *EventDispatcher                                │  │
│  │    hooks: *HookManager                                     │  │
│  │    processContexts: map[uint32]*ProcessContext             │  │
│  │  }                                                        │  │
│  └────────────────────────────────────────────────────────────┘  │
│         │                                                          │
│         │ 1:N                                                     │
│         ▼                                                          │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  struct VCPU {                                              │  │
│  │    CoreId: uint32                                           │  │
│  │    Regs: *GUEST_REGS                                        │  │
│  │    XmmRegs: *GUEST_XMM_REGS                                 │  │
│  │    ExitReason: uint32                                       │  │
│  │    LastVmexitRip: uint64                                    │  │
│  │                                                             │  │
│  │    VmcsRegionPhysicalAddress: uint64                        │  │
│  │    VmxonRegionPhysicalAddress: uint64                       │  │
│  │    MsrBitmapPhysicalAddress: uint64                         │  │
│  │    IoBitmapPhysicalAddressA/B: uint64                       │  │
│  │    VmmStack: uint64                                        │  │
│  │                                                             │  │
│  │    EptPointer: EPT_POINTER                                  │  │
│  │    EptPageTable: *EPT_PAGE_TABLE                            │  │
│  │                                                             │  │
│  │    OnVmxRootMode: BOOLEAN                                   │  │
│  │    IncrementRip: BOOLEAN                                     │  │
│  │    HasLaunched: BOOLEAN                                     │  │
│  │  }                                                         │  │
│  └────────────────────────────────────────────────────────────┘  │
│         │                                                          │
│         │ 1:1                                                     │
│         ▼                                                          │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  struct EptTable {                                          │  │
│  │    pml4: [512]EptEntry                                      │  │
│  │    pdptTables: [511]*[512]EptEntry                          │  │
│  │    pdTables: [512][512]*EptEntry                            │  │
│  │    splits: []SplitRecord                                    │  │
│  │    mtrrState: *MtrrState                                   │  │
│  │  }                                                         │  │
│  └────────────────────────────────────────────────────────────┘  │
│         │                                                          │
│         │ 1:N (via HookManager)                                  │
│         ▼                                                          │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │  struct EptHook {                                           │  │
│  │    PhysAddr: uint64                                         │  │
│  │    VirtAddr: uint64                                         │  │
│  │    Cr3: CR3_TYPE                                            │  │
│  │    Kind: HookKind (Exec/Read/Write/RW)                     │  │
│  │    IsActive: bool                                           │  │
│  │                                                             │  │
│  │    OriginalEntry: *EptEntry  (备份原始EPT项)               │  │
│  │    ModifiedEntry: *EptEntry   (修改后的EPT项)              │  │
│  │                                                             │  │
│  │    breakpoints: []Breakpoint                                │  │
│  │    fakePage: [4096]byte        (假页内容)                   │  │
│  │    FakePagePa: uint64           (假页物理地址)              │  │
│  │  }                                                         │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

## 8. VMCS 字段速查表

### Guest State Area
| Offset | Field | Description |
|--------|-------|-------------|
| 0x6800 | CR0 | Guest CR0 |
| 0x6802 | CR3 | Guest CR3 (Page Table Base) |
| 0x6804 | CR4 | Guest CR4 |
| 0x681A | DR7 | Debug Register 7 |
| 0x681C | RSP | Guest Stack Pointer |
| 0x681E | RIP | Guest Instruction Pointer |
| 0x6820 | RFLAGS | Guest RFLAGS |
| 0x6826 | EFER | Extended Feature Enable |

### Control Fields
| Offset | Field | Description |
|--------|-------|-------------|
| 0x4000 | PIN_CONTROLS | Pin-Based VM Execution Controls |
| 0x4002 | PRIMARY_PROC | Primary Processor-Based Controls |
| 0x401E | SECONDARY_PROC | Secondary Processor-Based Controls |
| 0x4004 | EXCEPTION_BITMAP | Exception Bitmap |
| 0x400C | IO_BITMAP_A | I/O Bitmap A Address |
| 0x400E | IO_BITMAP_B | I/O Bitmap B Address |
| 0x4010 | MSR_BITMAP | MSR Bitmap Address |
| 0x201A | EPT_POINTER | EPT Pointer (EPTP) |

### Host State Area
| Offset | Field | Description |
|--------|-------|-------------|
| 0x6C00 | HOST_CR0 | Host CR0 |
| 0x6C02 | HOST_CR3 | Host CR3 |
| 0x6C04 | HOST_CR4 | Host CR4 |
| 0x6C02 | HOST_RSP | Host Stack Pointer (VMM Stack) |
| 0x6C16 | HOST_RIP | Host Instruction Pointer (vmExitHandler) |
| 0xC00 | HOST_CS | Host CS Selector |
| 0xC04 | HOST_SS | Host SS Selector |
| 0xC08 | HOST_FS | Host FS Selector |
| 0xC0A | HOST_GS | Host GS Selector |
| 0xC0E | HOST_TR | Host TR Selector |

---

## 9. EPT 分页原理详解 (Extended Page Tables)

### 9.1 为什么需要 EPT？- 两级地址转换的必要性

在没有虚拟化的情况下，CPU 使用一套页表将 **虚拟地址 (VA)** 转换为 **物理地址 (PA)**：

```
┌─────────────────────────────────────────────────────────────┐
│                    无虚拟化 (Native Execution)               │
│                                                             │
│   Guest Application                                         │
│       │                                                     │
│       ▼                                                     │
│   Virtual Address (VA)                                      │
│       │                                                     │
│       ▼                                                     │
│   ┌──────────────────────────────────────────────────┐      │
│   │           Windows Page Tables (CR3)              │      │
│   │                                                   │      │
│   │    PML4 → PDPT → PD → PT → Physical Page        │      │
│   └──────────────────────────────────────────────────┘      │
│       │                                                     │
│       ▼                                                     │
│   Physical Address (PA)                                     │
│       │                                                     │
│       ▼                                                     │
│   RAM (物理内存)                                             │
└─────────────────────────────────────────────────────────────┘
```

**问题：** 在虚拟化环境中，Guest OS 认为自己拥有整个物理内存，但实际上它看到的"物理地址"是 **Guest Physical Address (GPA)**，不是真实的机器地址。因此需要 **第二级地址转换**：

```
┌─────────────────────────────────────────────────────────────┐
│                    有虚拟化 (VMX + EPT)                       │
│                                                             │
│   Guest Application                                         │
│       │                                                     │
│       ▼                                                     │
│   Guest Virtual Address (GVA)                               │
│       │                                                     │
│       ▼                                                     │
│   ┌──────────────────────────────────────────────────┐      │
│   │     Guest Page Tables (Guest CR3)                │      │
│   │         第一级转换: GVA → GPA                     │      │
│   │         (由 Guest OS 管理)                        │      │
│   └──────────────────────────────────────────────────┘      │
│       │                                                     │
│       ▼                                                     │
│   Guest Physical Address (GPA)                              │
│       │                                                     │
│       ▼                                                     │
│   ┌──────────────────────────────────────────────────┐      │
│   │     Extended Page Tables (EPTP in VMCS)          │      │
│   │         第二级转换: GPA → HPA (Host PA)           │      │
│   │         由 Hypervisor 管理                        │      │
│   └──────────────────────────────────────────────────┘      │
│       │                                                     │
│       ▼                                                     │
│   Host Physical Address (HPA / 真实物理地址)                 │
│       │                                                     │
│       ▼                                                     │
│   RAM (真实物理内存)                                         │
└─────────────────────────────────────────────────────────────┘
```

### 9.2 Windows 分页 vs EPT 分页 - 结构对比

#### Windows x86-64 分页结构 (4级)

```
┌─────────────────────────────────────────────────────────────┐
│              Windows x86-64 页表结构 (CR3 指向)              │
│                                                             │
│   48-bit Virtual Address                                    │
│   ┌─────┬─────┬─────┬─────┬───────────────────────────┐     │
│   │PML4 │PDPT │ PD  │ PT  │       Page Offset          │     │
│   │9bit │9bit │9bit │12bit│        12bit              │     │
│   └──┬──┴──┬──┴──┬──┴──┬──┴────────────┬──────────────┘     │
│      │     │     │     │               │                    │
│      ▼     ▼     ▼     ▼               ▼                    │
│   ┌────┐┌────┐┌────┐┌────┐     ┌────────────┐             │
│   │PML4│→│PDPT│→│PD │→│PT │     │Physical    │             │
│   │512 │ │512 │ │512 │ │512 │     │Page 4KB    │             │
│   │Entry│ │Entry│ │Entry│ │Entry│     │            │             │
│   └────┘└────┘└────┘└────┘     └────────────┘             │
│                                                             │
│   每个 Entry 包含:                                          │
│   ├─ Bit 0: Present (1=存在)                                │
│   ├─ Bit 1: Read/Write (0=R, 1=RW)                         │
│   ├─ Bit 2: User/Supervisor (0=Kernel, 1=User)             │
│   ├─ Bits 51:12: 物理页帧号 (PFN)                           │
│   └─ 其他位: NX, Dirty, Accessed, PAT, etc.                │
│                                                             │
│   支持大页:                                                 │
│   ├─ 2MB Large Page: PT 层直接指向物理页 (PS=1)             │
│   └─ 1GB Huge Page: PD 层直接指向物理页 (PS=1)              │
└─────────────────────────────────────────────────────────────┘
```

#### Intel EPT 分页结构 (4级)

```
┌─────────────────────────────────────────────────────────────┐
│              Intel EPT 页表结构 (EPTP 指向)                  │
│                                                             │
│   Guest Physical Address (GPA)                              │
│   ┌─────┬─────┬─────┬─────┬───────────────────────────┐     │
│   │PML4 │PDPT │ PD  │ PT  │       Page Offset          │     │
│   │9bit │9bit │9bit │12bit│        12bit              │     │
│   └──┬──┴──┬──┴──┬──┴──┬──┴────────────┬──────────────┘     │
│      │     │     │     │               │                    │
│      ▼     ▼     ▼     ▼               ▼                    │
│   ┌────┐┌────┐┌────┐┌────┐     ┌────────────┐             │
│   │EPML4│→│EPDPT│→│EPD │→│EPT │     │Host        │             │
│   │512 │ │512 │ │512 │ │512 │     │Physical    │             │
│   │Entry│ │Entry│ │Entry│ │Entry│     │Page 4KB    │             │
│   └────┘└────┘└────┘└────┘     └────────────┘             │
│                                                             │
│   EPT Entry 格式 (128 bit, 与 Windows 完全不同!):           │
│   ┌────────────────────────────────────────────────────┐   │
│   │ 63:12  │ 11:8  │ 7  │ 6  │ 5  │ 3:2  │ 1  │ 0    │   │
│   ├────────┼───────┼────┼────┼────┼─────┼────┼──────┤   │
│   │ Phys   │ Mem   │ Ign│ Ig │ Ig │ IPAT│ MT │ Read │   │
│   │ Addr   │ Type  │    |    |    │     │    │      │   │
│   │ (52b)  │ (3b)  │    |    |    │     │(2b)│ (1b) │   │
│   └────────┴───────┴────┴────┴────┴─────┴────┴──────┘   │
│                                                             │
│   关键差异:                                                  │
│   ├─ Bit 0: Read (R) = 可读权限                            │
│   ├─ Bit 1: Write (W) = 可写权限 (隐含, 见下文)            │
│   ├─ Bit 2: Execute (X) = 可执行权限 (隐含, 见下文)        │
│   ├─ Bits 2:0 = 000 表示此 Entry 未映射 (Not Present)      │
│   ├─ Bits 5:3 = Memory Type (MT):                          │
│   │   000 = Uncacheable (UC)                              │
│   │   001 = Write Combining (WC)                          │
│   │   010 = Write Through (WT)                            │
│   │   011 = Write Protected (WP)                          │
│   │   100 = Write Back (WB) ← 默认值                      │
│   │   101 = Uncacheable (UC-)                             │
│   │   110 = Write Combining (WC-)                         │
│   │   111 = Write Back (WB-)                              │
│   ├─ Bit 6: Ignore PAT (IPAT)                             │
│   └─ Bits 63:12 = 物理地址 (支持到 52 位)                  │
│                                                             │
│   支持大页:                                                 │
│   ├─ 2MB Large Page: EPT 层直接指向物理页                   │
│   └─ 1GB Huge Page: EPDPT 层直接指向物理页                  │
└─────────────────────────────────────────────────────────────┘
```

### 9.3 权限模型对比 - 核心区别

```
┌─────────────────────────────────────────────────────────────┐
│                    权限模型对比                               │
│                                                             │
│  ╔═════════════════════════════════════════════════════╗    │
│  ║           Windows Page Table Entry (PTE)              ║    │
│  ╠═════════════════════════════════════════════════════╣    │
│  ║  R/W/U/S 三维权限:                                   ║    │
│  ║                                                       ║    │
│  ║  ┌─────────────────────────────────────────┐        ║    │
│  ║  │ Bit 0: Present (存在性)                   │        ║    │
│  ║  │ Bit 1: Read/Write (0=只读, 1=读写)        │        ║    │
│  ║  │ Bit 2: User/Supervisor (0=内核, 1=用户)  │        ║    │
│  ║  │ Bit 63: No Execute (NX, 1=不可执行)      │        ║    │
│  ║  └─────────────────────────────────────────┘        ║    │
│  ║                                                       ║    │
│  ║  组合示例:                                            ║    │
│  ║  ├── P=1 RW=0 U=0 NX=0 → 内核代码段 (可读+执行)      ║    │
│  ║  ├── P=1 RW=1 U=0 NX=1 → 内核数据段 (可读+写)        ║    │
│  ║  ├── P=1 RW=0 U=1 NX=0 → 用户代码段 (可读+执行)      ║    │
│  ║  └── P=1 RW=1 U=1 NX=1 → 用户堆栈 (可读+写)          ║    │
│  ╚═════════════════════════════════════════════════════╝    │
│                                                             │
│  ╔═════════════════════════════════════════════════════╗    │
│  ║           EPT Entry (Extended Page Table)             ║    │
│  ╠═════════════════════════════════════════════════════╣    │
│  ║  R/W/X 三维独立权限 (更灵活!):                       ║    │
│  ║                                                       ║    │
│  ║  ┌─────────────────────────────────────────┐        ║    │
│  ║  │ Bit 0: Read (R) - 允许读取               │        ║    │
│  ║  │ Bit 1: Write (W) - 允许写入              │        ║    │
│  ║  │ Bit 2: Execute (X) - 允许执行            │        ║    │
│  ║  │                                               │        ║    │
│  ║  │ 注意: W 和 X 是隐式编码的!               │        ║    │
│  ║  │   R=0,W=0,X=0 → Not Present (未映射)     │        ║    │
│  ║  │   R=1,W=0,X=0 → 只读 (Read Only)         │        ║    │
│  ║  │   R=1,W=1,X=0 → 读写 (Read+Write)        │        ║    │
│  ║  │   R=1,W=0,X=1 → 读+执行 (Execute)        │        ║    │
│  ║  │   R=1,W=1,X=1 → 全部权限 (RWX) ★ 默认    │        ║    │
│  ║  └─────────────────────────────────────────┘        ║    │
│  ╚═════════════════════════════════════════════════════╝    │
│                                                             │
│  ⭐ 关键优势:                                                │
│  Windows 无法单独控制 "读但不可执行" 或 "执行但不可读"       │
│  但 EPT 可以！这就是 HyperDbg 能实现 Hook 的核心原因          │
└─────────────────────────────────────────────────────────────┘
```

### 9.4 完整地址转换流程图

```
┌─────────────────────────────────────────────────────────────────┐
│              完整的两级地址转换流程                              │
│                                                                 │
│  场景: Guest 执行指令 MOV EAX, [0x12345678]                    │
│                                                                 │
│  Step 1: Guest 线性地址生成                                     │
│  ┌──────────────────────────────────────────────────────┐      │
│  │  GVA = Segment Base + Offset (或 Flat Model 下=GVA)   │      │
│  │                                                      │      │
│  │  GVA = 0x00007FF6`A1B20000 (假设)                    │      │
│  └──────────────────────┬───────────────────────────────┘      │
│                         ▼                                      │
│  Step 2: 第一级转换 (Guest CR3 → GPA)                          │
│  ┌──────────────────────────────────────────────────────┐      │
│  │  CR3 (Guest) = 0x1A000000                              │      │
│  │                                                      │      │
│  │  GVA[47:39]=0x000 → PML4[0] @ PA 0x1A000000          │      │
│  │    → PML4E.Present=1, PFN=0x1B0000                    │      │
│  │                                                      │      │
│  │  GVA[38:30]=0x080 → PDPT[8] @ PA 0x1B00040           │      │
│  │    → PDPTE.Present=1, PS=0, PFN=0x1C0000              │      │
│  │                                                      │      │
│  │  GVA[29:21]=0x064 → PD[100] @ PA 0x1C00200           │      │
│  │    → PDE.Present=1, PS=0, PFN=0x500000               │      │
│  │                                                      │      │
│  │  GVA[20:12]=0x200 → PT[512] @ PA 0x501000           │      │
│  │    → PTE.Present=1, PFN=0xABC000                     │      │
│  │                                                      │      │
│  │  GPA = 0xABC000 + 0x000 = 0xABC000                   │      │
│  └──────────────────────┬───────────────────────────────┘      │
│                         ▼                                      │
│  Step 3: 第二级转换 (EPTP → HPA) ⭐                            │
│  ┌──────────────────────────────────────────────────────┐      │
│  │  EPTP (from VMCS) = 0xF8000000                        │      │
│  │                                                      │      │
│  │  GPA[47:39]=0x000 → EPML4[0] @ HPA 0xF8000000        │      │
│  │    → E-PML4E.RWX=111, PFN=0xF900000                  │      │
│  │                                                      │      │
│  │  GPA[38:30]=0x000 → EPDPT[0] @ HPA 0xF9000000        │      │
│  │    → E-PDPTE.RWX=111, PFN=0xFA00000                  │      │
│  │                                                      │      │
│  │  GPA[29:21]=0x0AB → EPD[171] @ HPA 0xFA00558         │      │
│  │    → E-PDE.RWX=110, PS=1 (2MB Page!)                 │      │
│  │    → HPA Base = 0xABC00000                            │      │
│  │                                                      │      │
│  │  HPA = 0xABC00000 + (GPA & 0x1FFFFF)                  │      │
│  │     = 0xABC00000 + 0xBC000                            │      │
│  │     = 0xABCBC000  ← 最终物理地址!                      │      │
│  └──────────────────────┬───────────────────────────────┘      │
│                         ▼                                      │
│  Step 4: 内存访问                                              │
│  ┌──────────────────────────────────────────────────────┐      │
│  │  CPU 从 HPA 0xABCBC000 读取 4 字节                    │      │
│  │  返回给 Guest 寄存器 EAX                               │      │
│  └──────────────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────────────┘
```

### 9.5 EPT Violation 触发条件详解

```
┌─────────────────────────────────────────────────────────────┐
│              EPT Violation 触发条件矩阵                      │
│                                                             │
│  当 Guest 尝试访问某个 GPA 时，硬件会检查 EPT 权限:          │
│                                                             │
│  ┌───────────┬─────────────────────────────────────────┐   │
│  │  操作类型  │  需要的 EPT 权限位                      │   │
│  ├───────────┼─────────────────────────────────────────┤   │
│  │  读取     │  R=1 (Bit 0)                            │   │
│  │  写入     │  R=1 AND W=1 (Bit 0 & Bit 1)           │   │
│  │  取指     │  R=1 AND X=1 (Bit 0 & Bit 2)           │   │
│  │  DMA 读取 │  R=1 (同读取)                           │   │
│  │  DMA 写入 │  R=1 AND W=1 (同写入)                   │   │
│  └───────────┴─────────────────────────────────────────┘   │
│                                                             │
│  ╔═════════════════════════════════════════════════════╗   │
│  ║  Violation 示例: Exec Hook 的原理                    ║   │
│  ╠═════════════════════════════════════════════════════╣   │
│  ║                                                       ║   │
│  ║  原始 EPT Entry:                                       ║   │
│  ║  ┌─────────────────────────────────────────┐        ║   │
│  ║  │ R=1, W=1, X=1 (RWX) → 正常访问          │        ║   │
│  ║  │ PhysAddr = 0xABC000 (原始物理页)          │        ║   │
│  ║  └─────────────────────────────────────────┘        ║   │
│  ║                                                       ║   │
│  ║  Hook 后 EPT Entry (ModifiedEntry):                   ║   │
│  ║  ┌─────────────────────────────────────────┐        ║   │
│  ║  │ R=1, W=1, X=0 (RW_) → 去掉执行权限!     │        ║   │
│  ║  │ PhysAddr = 0xABC000 (仍指向原页)         │        ║   │
║  ║  └─────────────────────────────────────────┘        ║   │
│  ║                                                       ║   │
│  ║  当 CPU 尝试从该页取指时:                             ║   │
│  ║  1. 需要 X=1, 但实际 X=0                            ║   │
│  ║  2. 硬件触发 VM Exit (EPT Violation)                 ║   │
│  ║  3. Hypervisor 接管控制权                            ║   │
│  ║  4. 可以检查是否命中断点、记录日志等                   ║   │
│  ║                                                       ║   │
│  ╚═════════════════════════════════════════════════════╝   │
│                                                             │
│  ╔═════════════════════════════════════════════════════╗   │
│  ║  Misconfig vs Violation 区别                         ║   │
│  ╠═════════════════════════════════════════════════════╣   │
│  ║                                                       ║   │
│  ║  EPT Misconfig (Exit Reason = 31):                    ║   │
│  ║  - EPT 表项本身配置错误                               ║   │
│  ║  - 例如: 中间层 Entry 的 RWX=000 (未映射)            ║   │
│  ║  - 或者: 物理地址超出 MAXPHYADDR                      ║   │
│  ║  - 结果: 注入 #UD 异常到 Guest                        ║   │
│  ║                                                       ║   │
│  ║  EPT Violation (Exit Reason = 30):                    ║   │
│  ║  - 最终页面的权限不足                                 ║   │
│  ║  - 例如: 尝试执行但 X=0                              ║   │
│  ║  - 或者: 尝试写入但 W=0                              ║   │
│  ║  - 结果: Hypervisor 处理 (Hook/监控)                  ║   │
│  ║                                                       ║   │
│  ╚═════════════════════════════════════════════════════╝   │
└─────────────────────────────────────────────────────────────┘
```

### 9.6 大页拆分过程详细图解

```
┌─────────────────────────────────────────────────────────────┐
│              EPT 大页拆分过程 (Split Large Page)             │
│                                                             │
│  场景: 需要在 2MB 区域内的某个 4KB 页上设置 Hook             │
│                                                             │
│  ╔══ 初始状态: 2MB Large Page ═════════════════════════╗   │
│  ║                                                       ║   │
│  ║  GPA Range: 0x40000000 - 0x41FFFFFF (2MB)            ║   │
│  ║                                                       ║   │
│  ┌─────────────────────────────────────────────────┐    ║   │
│  │  EPML4[0]                                        │    ║   │
│  │    └→ EPDPT[0]                                   │    ║   │
│  │         └→ EPD[64]  (PS=1, 2MB Page)            │    ║   │
│  │              R=1 W=1 X=1                         │    ║   │
│  │              PhysAddr = 0xABC00000               │    ║   │
│  │              MemType = WB (Write Back)           │    ║   │
│  └─────────────────────────────────────────────────┘    ║   │
│  ║                                                       ║   │
│  ║  映射关系:                                             ║   │
│  ║  GPA 0x40000000 ~ 0x41FFFFFF                          ║   │
│  ║  └→ HPA 0xABC00000 ~ 0xABDFFFFF (1:1)                ║   │
│  ╚═════════════════════════════════════════════════════╝   │
│                         │                                  │
│                         ▼ SplitLargePage() 调用             │
│                                                             │
│  ╔══ 拆分后状态: 512 × 4KB Pages ══════════════════════╗   │
│  ║                                                       ║   │
│  ┌─────────────────────────────────────────────────┐    ║   │
│  │  EPML4[0]                                        │    ║   │
│  │    └→ EPDPT[0]                                   │    ║   │
│  │         └→ EPD[64]  (PS=0, Points to PT)         │    ║   │
│  │              R=1 W=1 X=1                         │    ║   │
│  │              PhysAddr = New_PT_Base_PA           │    ║   │
│  │                                                   │    ║   │
│  │         ┌─────────────────────────────────┐      │    ║   │
│  │         │  New PT Table (512 Entries)      │      │    ║   │
│  │         │                                   │      │    ║   │
│  │         │  PT[0]:  R=1 W=1 X=1  → 0xABC00000 │      │    ║   │
│  │         │  PT[1]:  R=1 W=1 X=1  → 0xABC01000 │      │    ║   │
│  │         │  ...                               │      │    ║   │
│  │         │  PT[N]:  R=1 W=1 X=0  → 0xABCN0000 │← Hooked│    ║   │
│  │         │  ...                               │      │    ║   │
│  │         │  PT[511]:R=1 W=1 X=1  → 0xABCFF000│      │    ║   │
│  │         │                                   │      │    ║   │
│  │         │  所有项继承原始 Memory Type (WB)    │      │    ║   │
│  │         └─────────────────────────────────┘      │    ║   │
│  └─────────────────────────────────────────────────┘    ║   │
│  ║                                                       ║   │
│  ║  新映射关系:                                           ║   │
│  ║  GPA 0x40000000 ~ 0x40FFFFF                           ║   │
│  ║  └→ HPA 0xABC00000 ~ 0xABCFFFFF (PT[0])               ║   │
│  ║                                                       ║   │
│  ║  GPA 0x40100000 ~ 0x401FFFFF                          ║   │
│  ║  └→ HPA 0xABC01000 ~ 0xABC01FFF (PT[1])               ║   │
│  ║                                                       ║   │
│  ║  GPA 0x40N0000 ~ 0x40NFFFFF  (被 Hook 的页)           ║   │
│  ║  └→ HPA 0xABCN000 ~ 0xABCNFFF (PT[N], X=0!)          ║   │
│  ╚═════════════════════════════════════════════════════╝   │
│                         │                                  │
│                         ▼ 设置 Exec Hook                    │
│                                                             │
│  ╔══ Hook 安装完成 ════════════════════════════════════╗   │
│  ║                                                       ║   │
│  ┌─────────────────────────────────────────────────┐    ║   │
│  │  EptHook 结构体:                                 │    ║   │
│  │                                                   │    ║   │
│  │  PhysAddr    = 0xABCN0000 (目标物理页)            │    ║   │
│  │  VirtAddr    = 0x40Nxxxxx (目标虚拟地址)           │    ║   │
│  │  Cr3         = Guest_CR3                          │    ║   │
│  │  Kind        = HookExec                           │    ║   │
│  │                                                   │    ║   │
│  │  OriginalEntry = {                                │    ║   │
│  │    R=1, W=1, X=1,                                │    ║   │
│  │    PhysAddr = 0xABCN0000,                         │    ║   │
│  │    MemType = WB                                   │    ║   │
│  │  }                                                │    ║   │
│  │                                                   │    ║   │
│  │  ModifiedEntry = {                                │    ║   │
│  │    R=1, W=1, X=0,  ← 去掉执行权限!               │    ║   │
│  │    PhysAddr = 0xABCN0000,  ← 仍然指向原页         │    ║   │
│  │    MemType = WB                                   │    ║   │
│  │  }                                                │    ║   │
│  │                                                   │    ║   │
│  │  FakePage = [00 90 90 90 90 90 ... ]              │    ║   │
│  │  FakePagePa = MmGetPhysicalAddress(FakePage)      │    ║   │
│  └─────────────────────────────────────────────────┘    ║   │
│  ╚═════════════════════════════════════════════════════╝   │
└─────────────────────────────────────────────────────────────┘
```

### 9.7 EPT Pointer (EPTP) 格式详解

```
┌─────────────────────────────────────────────────────────────┐
│              EPT Pointer (EPTP) 格式 - 写入 VMCS offset 201A│
│                                                             │
│  ┌────────────────────────────────────────────────────┐    │
│  │  63:N        N:5   4:3   2    1    0              │    │
│  ├────────────────┼──────┼─────┼────┼────┼─────────────┤    │
│  │  PML4 Physical │ Walk | Mem | Ena| Ena| Reserved    │    │
│  │  Address       | Len  | Type│ ble | bled|             │    │
│  │  (必须4KB对齐)  │      |     |     |     |             │    │
│  └────────────────┴──────┴─────┴────┴────┴─────────────┘    │
│                                                             │
│  各字段含义:                                                 │
│  ┌────────────────────────────────────────────────────┐    │
│  │                                                    │    │
│  │  Bits [63:N]:                                      │    │
│  │    PML4 Table 的物理基地址 (必须 4KB 对齐)          │    │
│  │    N = MAXPHYADDR - 1 (通常为 51, 支持 52位物理地址)│    │
│  │                                                    │    │
│  │  Bits [5:3]: Page-Walk Length                      │    │
│  │    EPT 页表遍历层数减 1                             │    │
│  │    000 = 3 层 (不常用)                              │    │
│  │    001 = 4 层 ← 我们使用的标准值 (PML4/PDPT/PD/PT) │    │
│  │    010~111 = Reserved                               │    │
│  │                                                    │    │
│  │  Bits [4:3]: EPT Memory Type                       │    │
│  │    用于 EPT 页表结构的缓存策略                       │    │
│  │    00 = Uncacheable (UC)                           │    │
│  │    01 = Write Back (WB) ← 推荐, 性能最好            │    │
│  │    10 = Write Through (WT)                          │    │
│  │    11 = Write Protected (WP)                        │    │
│  │                                                    │    │
│  │  Bit 2: Enable Access and Dirty Flags              │    │
│  │    0 = 不启用 A/D 位跟踪 (性能更好)                  │    │
│  │    1 = 启用 A/D 位跟踪                              │    │
│  │    ← 通常设为 0                                     │    │
│  │                                                    │    │
│  │  Bit 1: Enable EPT Execution-Only Translations     │    │
│  │    0 = 禁用 (不支持仅执行转换)                       │    │
│  │    1 = 启用 (允许设置 X-only 页面)                  │    │
│  │    ← 必须设为 1 以支持 Exec Hook!                   │    │
│  │                                                    │    │
│  │  Bit 0: Reserved (Must be 0)                       │    │
│  │                                                    │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
│  HyperDbg 使用的典型 EPTP 值:                                │
│  ┌────────────────────────────────────────────────────┐    │
│  │                                                    │    │
│  │  EPTP = PML4_Phys_Addr                             │    │
│  │       | (Bit 5:3 = 001)  // 4-level walk           │    │
│  │       | (Bit 4:3 = 01)   // WB memory type         │    │
│  │       | (Bit 2 = 0)      // A/D disabled           │    │
│  │       | (Bit 1 = 1)      // Enable exec-only       │    │
│  │       | (Bit 0 = 0)      // Reserved              │    │
│  │                                                    │    │
│  │  Example:                                           │    │
│  │  If PML4 PA = 0xF8000000                           │    │
│  │  Then EPTP = 0xF800000B                            │    │
│  │       (0xF8000000 | 0xB)                           │    │
│  │                                                    │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 9.8 Go 实现 vs ASM 文件的关系

```
┌─────────────────────────────────────────────────────────────┐
│              Go 代码与 .s 汇编文件的分工                      │
│                                                             │
│  ╔═════════════════════════════════════════════════════╗   │
│  ║  Go (.go 文件) - 业务逻辑层                          ║   │
│  ╠═════════════════════════════════════════════════════╣   │
│  ║                                                       ║   │
│  ║  ✅ 负责:                                             ║   │
│  ║  ├── EPT 数据结构定义和操作 (ept.go)                  ║   │
│  ║  │   ├── EptTable struct                             ║   │
│  ║  │   ├── EptEntry methods (SetPermission, etc.)      ║   │
│  ║  │   ├── BuildEptPointer()                           ║   │
│  ║  │   ├── SplitLargePage() / MergeSmallPages()        ║   │
│  ║  │   └── Identity Map 初始化                         ║   │
│  ║  │                                                    ║   │
│  ║  ├── VMCS 配置逻辑 (hypervisor.go)                   ║   │
│  ║  │   ├── setupVmcs() - 填充所有字段                   ║   │
│  ║  │   ├── adjustVmcsControl() - 控制域调整             ║   │
│  ║  │   └── initializeVcpu() - VCPU 初始化序列           ║   │
│  ║  │                                                    ║   │
│  ║  ├── Event Dispatch (events.go)                      ║   │
│  ║  │   └── HandleVmExit() - VM Exit 处理入口           ║   │
│  ║  │                                                    ║   │
│  ║  ├── Hook Management (hooks.go)                      ║   │
│  ║  │   ├── InstallEptHook()                            ║   │
│  ║  │   ├── HandleExecHook()                            ║   │
│  ║  │   └── RestoreOriginalAndInjectBp()                ║   │
│  ║  │                                                    ║   │
│  ║  └── IOCTL Handler (ioctl.go)                        ║   │
│  ║      └── 50+ 用户态命令处理                          ║   │
│  ║                                                       ║   │
│  ╚═════════════════════════════════════════════════════╝   │
│                                                             │
│  ╔═════════════════════════════════════════════════════╗   │
│  ║  Assembly (.s 文件) - 底层指令层                      ║   │
│  ╠═════════════════════════════════════════════════════╣   │
│  ║                                                       ║   │
│  ║  ✅ 负责:                                             ║   │
│  ║  ├── VMX Instructions (vmx.s)                        ║   │
│  ║  │   ├── AsmVmxVmxOn(pa)      → VMXON 指令          ║   │
│  ║  │   ├── AsmVmxVmxOff()       → VMXOFF 指令         ║   │
│  ║  │   ├── AsmVmxVmxClear(pa)   → VMCLEAR 指令       ║   │
│  ║  │   ├── AsmVmxVmxPtrld(pa)   → VMPTRLD 指令       ║   │
│  ║  │   ├── AsmVmxVmread64(field, val) → VMREAD       ║   │
│  ║  │   ├── AsmVmxVmwrite64(field, val)→ VMWRITE      ║   │
│  ║  │   ├── AsmVmxVmlaunch()     → VMLAUNCH 指令      ║   │
│  ║  │   └── AsmVmxVmresume()     → VMRESUME 指令      ║   │
│  ║  │                                                    ║   │
│  ║  ├── Context Save/Restore (context.s)               ║   │
│  ║  │   ├── vmExitHandler()      → VM Exit 入口点     ║   │
│  ║  │   │   - 保存所有通用寄存器到 GUEST_REGS          ║   │
│  ║  │   │   - 调用 Go 的 HandleVmExit(regs)            ║   │
│  ║  │   │   - 恢复寄存器并执行 VMRESUME               ║   │
│  ║  │   │                                               ║   │
│  ║  │   ├── getCs/getSs/getDs... → 读取 Segment Regs  ║   │
│  ║  │   └── __readmsr/__writemsr → MSR 访问           ║   │
│  ║  │                                                    ║   │
│  ║  └── TLB Invalidation (tlb.s)                       ║   │
│  ║      ├── Invept()           → INVEPT 指令           ║   │
│  ║      ├── Invpid()           → INVPID 指令           ║   │
│  ║      └── Invvpid()          → INVVPID 指令          ║   │
│  ║                                                       ║   │
│  ╚═════════════════════════════════════════════════════╝   │
│                                                             │
│  调用关系图:                                                │
│  ┌────────────────────────────────────────────────────┐    │
│  │                                                    │    │
│  │  hypervisor.go: setupVmcs()                       │    │
│  │      │                                           │    │
│  │      ├─ vmWrite64(field, value)                  │    │
│  │      │   └─ AsmVmxVmwrite64(field, value)  ◄─────┼──── .s
│  │      │                                           │    │
│  │      ├─ vmxTurnOn(pa)                            │    │
│  │      │   └─ AsmVmxVmxOn(pa)               ◄───────┼──── .s
│  │      │                                           │    │
│  │      └─ vmClear(pa)                             │    │
│  │          └─ AsmVmxVmxClear(pa)            ◄───────┼──── .s
│  │                                                    │    │
│  │  events.go: HandleVmExit(regs)                   │    │
│  │      │                                           │    │
│  │      ├─ vmRead64(field, &value)                 │    │
│  │      │   └─ AsmVmxVmread64(field, &value)◄──────┼──── .s
│  │      │                                           │    │
│  │      └─ dispatch() → hooks.HandleExecHook()     │    │
│  │          │                                       │    │
│  │          └─ setMonitorTrapFlag(true)             │    │
│  │              └─ vmWrite64(...) → .s             │    │
│  │                                                    │    │
│  │  context.s: vmExitHandler()                      │    │
│  │      │                                           │    │
│  │      ├─ PUSH ALL REGISTERS (保存 Guest State)    │    │
│  │      ├─ CALL HandleVmExit(GUEST_REGS*)  ◄────────┼──── Go
│  │      ├─ POP ALL REGISTERS (恢复 Host State)      │    │
│  │      └─ JMP VMRESUME (返回 Guest)                │    │
│  │                                                    │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 9.9 EPT Identity Map 实现细节

```
┌─────────────────────────────────────────────────────────────┐
│           EPT Identity Map 初始化过程                        │
│                                                             │
│  目标: 建立 GPA == HPA 的 1:1 映射 (默认配置)               │
│                                                             │
│  initIdentityMap() 函数流程:                                 │
│                                                             │
│  Step 1: 分配并清零 PML4 Table                               │
│  ┌────────────────────────────────────────────────────┐    │
│  │  pml4 = AllocateZeroedPool(4096)                    │    │
│  │  memset(pml4, 0, 4096)                              │    │
│  │  pml4_pa = MmGetPhysicalAddress(pml4)               │    │
│  └────────────────────────────────────────────────────┘    │
│                         │                                  │
│                         ▼                                  │
│  Step 2: 创建第一个 PDPT (覆盖低 512GB)                     │
│  ┌────────────────────────────────────────────────────┐    │
│  │  pdpt0 = AllocateZeroedPool(4096)                   │    │
│  │  pdpt0_pa = MmGetPhysicalAddress(pdpt0)             │    │
│  │                                                    │    │
│  │  设置 pml4[0]:                                      │    │
│  │    R=1, W=1, X=1                                   │    │
│  │    PhysAddr = pdpt0_pa >> 12                        │    │
│  │    MemType = WB                                    │    │
│  └────────────────────────────────────────────────────┘    │
│                         │                                  │
│                         ▼                                  │
│  Step 3: 为每个 PDPT Entry 创建 PD (使用 2MB Large Pages)   │
│  ┌────────────────────────────────────────────────────┐    │
│  │  for i := 0; i < 512; i++ {                        │    │
│  │      pd[i] = AllocateZeroedPool(4096)               │    │
│  │      pd_pa = MmGetPhysicalAddress(pd[i])            │    │
│  │                                                    │    │
│  │      设置 pdpt0[i]:                                 │    │
│  │        R=1, W=1, X=1                               │    │
│  │        PhysAddr = pd_pa >> 12                       │    │
│  │        MemType = WB                                │    │
│  │  }                                                 │    │
│  └────────────────────────────────────────────────────┘    │
│                         │                                  │
│                         ▼                                  │
│  Step 4: 为每个 PD Entry 设置 2MB Large Page (1:1 映射)     │
│  ┌────────────────────────────────────────────────────┐    │
│  │  for pd_idx := 0; pd_idx < 512; pd_idx++ {          │    │
│  │      for entry := 0; entry < 512; entry++ {         │    │
│  │          gpa_base = (pd_idx * 512 + entry) * 2MB    │    │
│  │          hpa_base = gpa_base  // 1:1 mapping!       │    │
│  │                                                    │    │
│  │          mem_type = GetMtrrMemoryType(gpa_base)     │    │
│  │          // 根据 MTRR 查询正确的内存类型             │    │
│  │                                                    │    │
│  │          设置 pd[pd_idx][entry]:                    │    │
│  │            R=1, W=1, X=1                            │    │
│  │            PS=1 (Large Page flag)                   │    │
│  │            PhysAddr = hpa_base >> 12                │    │
│  │            MemType = mem_type                       │    │
│  │      }                                             │    │
│  │  }                                                 │    │
│  └────────────────────────────────────────────────────┘    │
│                         │                                  │
│                         ▼                                  │
│  Step 5: 构建 EPTP 并写入 VMCS                             │
│  ┌────────────────────────────────────────────────────┐    │
│  │  eptp = pml4_pa                                     │    │
│  │       | (1 << 3)  // 4-level page walk              │    │
│  │       | (1 << 1)  // WB memory type                │    │
│  │       | (1 << 1)  // Enable exec-only translations │    │
│  │                                                    │    │
│  │  vmWrite64(VmcsCtrlEptPointer, eptp)               │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
│  最终结果:                                                   │
│  ┌────────────────────────────────────────────────────┐    │
│  │                                                    │    │
│  │  GPA 0x00000000`00000000 → HPA 0x00000000`00000000 │    │
│  │  GPA 0x00000000`00200000 → HPA 0x00000000`00200000 │    │
│  │  GPA 0x00000000`00400000 → HPA 0x00000000`00400000 │    │
│  │  ...                                                │    │
│  │  GPA 0x00007FFF`FFFE0000 → HPA 0x00007FFF`FFFE0000 │    │
│  │  (所有 512GB 都 1:1 映射, 权限 RWX)                 │    │
│  │                                                    │    │
│  │  后续 Hook 操作只需修改对应 EPT Entry 的权限即可!    │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 9.10 MTRR 与 EPT Memory Type 的关系

```
┌─────────────────────────────────────────────────────────────┐
│           MTRR (Memory Type Range Registers) 对 EPT 的影响   │
│                                                             │
│  问题: 不同物理内存区域有不同的缓存属性                       │
│  ┌────────────────────────────────────────────────────┐    │
│  │                                                    │    │
│  │  内存区域              缓存类型                     │    │
│  │  ─────────────────────────────────────────────     │    │
│  │  0x00000000-0x0BFFFFFF  DDR RAM (WB)              │    │
│  │  0xFEC00000-0xFEDFFFFF  I/O APIC (UC)              │    │
│  │  0xFEE00000-0xFEE00FFF  Local APIC (UC)            │    │
│  │  0xFF000000-0xFFFFFFFF  BIOS ROM (UC)              │    │
│  │  PCI MMIO Regions         (UC or WC)               │    │
│  │                                                    │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
│  解决方案: EPT Entry 必须使用正确的 Memory Type              │
│                                                             │
│  GetMtrrMemoryType(gpa) 算法:                               │
│  ┌────────────────────────────────────────────────────┐    │
│  │                                                    │    │
│  │  1. 检查 Fixed MTRRs (覆盖前 1MB)                   │    │
│  │     ├── IA32_MTRR_FIX64K_00000 (0x00000-0xFFFF)    │    │
│  │     ├── IA32_MTRR_FIX16K_80000 (0x80000-0xBFFFF)   │    │
│  │     └── IA32_MTRR_FIX4K_C0000 (0xC0000-0xFFFFF)    │    │
│  │                                                    │    │
│  │  2. 检查 Variable MTRRs (最多 16 个范围)            │    │
│  │     for i := 0; i < varCount; i++ {                │    │
│  │         base = IA32_MTRR_PHYSBASE[i]                │    │
│  │         mask = IA32_MTRR_PHYSMASK[i]                │    │
│  │         if gpa in range(base, mask) {               │    │
│  │             return base.type                        │    │
│  │         }                                           │    │
│  │     }                                               │    │
│  │                                                    │    │
│  │  3. 如果都不匹配, 返回 Default Type                  │    │
│  │     return IA32_MTRR_DEF_TYPE.bits.type             │    │
│  │                                                    │    │
│  └────────────────────────────────────────────────────┘    │
│                                                             │
│  ⚠️ 重要: 如果 EPT Memory Type 设置错误会导致:             │
│  - 缓存一致性问题 (Data Corruption)                        │
│  - 性能下降 (应该 UC 的用了 WB)                            │
│  - 硬件异常 (访问 PCI 设备用错误类型)                       │
│                                                             │
│  因此每次拆分大页或创建新的 EPT Entry 时都必须查询 MTRR!     │
└─────────────────────────────────────────────────────────────┘
```

### 9.11 总结: Windows 分页 vs EPT 分页对比表

| 特征 | Windows Page Table | EPT (Intel VT-x) |
|------|-------------------|-------------------|
| **用途** | VA → PA 转换 | GPA → HPA 转换 |
| **基址寄存器** | CR3 | EPTP (in VMCS) |
| **层级** | 4 级 (PML4/PDPT/PD/PT) | 4 级 (EPML4/EPDPT/EPD/EPT) |
| **Entry 大小** | 64 bits (8 bytes) | 128 bits (16 bytes) |
| **权限维度** | R/W/U/S/NX (4 维) | R/W/X (3 维, 更简洁) |
| **权限粒度** | 不能分离 "读但不执行" | ✅ 支持任意 R/W/X 组合 |
| **内存类型** | PAT (Page Attribute Table) | 直接在 Entry 中 (3 bits) |
| **管理者** | OS 内核 (ntoskrnl.exe) | Hypervisor (HyperDbg.sys) |
| **大页支持** | 2MB / 1GB | 2MB / 1GB |
| **TLB** | 普通 TLB | EPT TLB (独立缓存) |
| **失效指令** | INVLPG | INVEPT / INVPID |
| **主要用途** | 进程隔离、内存保护 | 虚拟化、内存监控、Hook |

---

*文档版本: 1.1*
*更新内容: 新增第9章 EPT 分页原理详解*
*基于: Intel SDM Vol. 3C Chapter 28 "Extended Page Tables"*
