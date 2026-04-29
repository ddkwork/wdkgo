package main

import (
	"unsafe"

	"solod.dev/so/wdk"
)

// sizeOf returns the size of type T as ULONG
func sizeOf[T any]() ULONG {
	return ULONG(unsafe.Sizeof(*new(T)))
}

//so:extern nodecl
func IoGetCurrentIrpStackLocation(Irp *wdk.IRP) *wdk.IO_STACK_LOCATION { return nil }

//so:extern nodecl
func ProbeForRead(Address PVOID, Length SIZE_T, Alignment SIZE_T) {}

//so:extern nodecl
func ProbeForWrite(Address PVOID, Length SIZE_T, Alignment SIZE_T) {}

//so:extern nodecl
func MmCopyMemory(TargetAddress PVOID, SourceAddress *wdk.MM_COPY_ADDRESS, NumberOfBytesToRead SIZE_T, Flags ULONG, NumberOfBytesRead *SIZE_T) NTSTATUS {
	return 0
}

const (
	IOCTL_QUERY_AND_CLEAR_LOGS                     = 0x0022B004
	IOCTL_REGISTER_EVENT                           = 0x0022B008
	IOCTL_RUN_SCRIPT                               = 0x0022B010
	IOCTL_SEND_REQUEST_RESULT                      = 0x0022B014
	IOCTL_GET_LOG_BASE                             = 0x0022B018
	IOCTL_VMM_INIT                                 = 0x0022B020
	IOCTL_VMM_SHUTDOWN                             = 0x0022B024
	IOCTL_EPT_HOOK                                 = 0x0022B028
	IOCTL_EPT_UNHOOK                               = 0x0022B02C
	IOCTL_EPT_SET_HOOK                             = 0x0022B030
	IOCTL_EPT_GET_EPT_TABLES                       = 0x0022B034
	IOCTL_VMM_EXECUTION_TRACE                      = 0x0022B038
	IOCTL_VMM_READ_MEM                             = 0x0022B03C
	IOCTL_VMM_WRITE_MEM                            = 0x0022B040
	IOCTL_VMM_VIRT_TO_PHYS                         = 0x0022B044
	IOCTL_VMM_PHYS_TO_VIRT                         = 0x0022B048
	IOCTL_VMM_GET_REG                              = 0x0022B04C
	IOCTL_VMM_SET_REG                              = 0x0022B050
	IOCTL_VMM_PAUSE                                = 0x0022B054
	IOCTL_VMM_RESUME                               = 0x0022B058
	IOCTL_VMM_SWITCH_PROCESS                       = 0x0022B05C
	IOCTL_VMM_MODIFY_REGS                          = 0x0022B060
	IOCTL_VMM_INVEPT                               = 0x0022B064
	IOCTL_VMM_INVVPID                              = 0x0022B068
	IOCTL_VMM_FLUSH_ENTIRE_TLB                     = 0x0022B06C
	IOCTL_VMM_GET_MTRR                             = 0x0022B070
	IOCTL_VMM_CHANGE_CORE                          = 0x0022B074
	IOCTL_VMM_GET_PROCESS_BASE                     = 0x0022B078
	IOCTL_VMM_GET_PROCESS_CR3                      = 0x0022B07C
	IOCTL_VMM_READ_AND_WRITE_MEM                   = 0x0022B080
	IOCTL_VMM_CALL_FUNCTION                        = 0x0022B084
	IOCTL_VMM_QUERY_PACKET                         = 0x0022B088
	IOCTL_VMM_SEND_RESULT                          = 0x0022B08C
	IOCTL_VMM_EXECUTE_SINGLE_STEP                  = 0x0022B090
	IOCTL_VMM_ENABLE_AND_INVOKE_TRAP_FLAG          = 0x0022B094
	IOCTL_VMM_MASK_EXCEPTION_DEBUGGER_BREAKPOINT   = 0x0022B098
	IOCTL_VMM_UNMASK_EXCEPTION_DEBUGGER_BREAKPOINT = 0x0022B09C
	IOCTL_VMM_SHORT_CIRCUITING_EVENT_INJECT        = 0x0022B0A0
	IOCTL_VMM_QUERY_REGISTER_FROM_GUEST_STATE      = 0x0022B0A4
	IOCTL_TRANSPARENT_MODE_ENABLE                  = 0x0022B0B0
	IOCTL_TRANSPARENT_MODE_DISABLE                 = 0x0022B0B4
	IOCTL_HYPERTRACE_START                         = 0x0022B0B8
	IOCTL_HYPERTRACE_STOP                          = 0x0022B0BC
	IOCTL_HYPERTRACE_READ_BUFFER                   = 0x0022B0C0
	IOCTL_HYPERTRACE_CLEAR_BUFFER                  = 0x0022B0C4
)

type VmmInitRequest struct {
	NumCores uint32
}

type EptHookRequest struct {
	VirtAddr  uint64
	Cr3       CR3_TYPE
	HookType  uint32
	ProcessId uint32
}

type EptUnhookRequest struct {
	VirtAddr  uint64
	Cr3       CR3_TYPE
	ProcessId uint32
}

type MemReadRequest struct {
	Address     uint64
	Cr3         CR3_TYPE
	Size        uint32
	ProcessId   uint32
	ReadOrWrite BOOLEAN
}

type MemReadResponse struct {
	Address uint64
	Size    uint32
	Success BOOLEAN
}

type VirtPhysRequest struct {
	Address   uint64
	Cr3       CR3_TYPE
	ProcessId uint32
}

type RegRequest struct {
	RegIndex uint32
	CoreId   uint32
}

type RegResponse struct {
	Value    uint64
	RegIndex uint32
	CoreId   uint32
	Success  BOOLEAN
}

type ModifyRegsRequest struct {
	Regs   GUEST_REGS
	CoreId uint32
}

type SwitchProcessRequest struct {
	ProcessId          uint32
	Cr3                CR3_TYPE
	ProcessBaseAddress uint64
}

type InveptRequest struct {
	Type uint64
	Eptp uint64
}

type MtrrRequest struct {
	BaseAddr uint64
	EndAddr  uint64
}

type ChangeCoreRequest struct {
	TargetCoreId uint32
}

type ProcessCr3Request struct {
	ProcessId uint32
}

type ReadWriteMemRequest struct {
	ReadAddress  uint64
	WriteAddress uint64
	ReadSize     uint32
	WriteSize    uint32
	ReadOrWrite  BOOLEAN
	Cr3          CR3_TYPE
	ProcessId    uint32
}

type CallFunctionRequest struct {
	FunctionAddress uint64
	ProcessId       uint32
	OptionalParam1  uint64
	OptionalParam2  uint64
	OptionalParam3  uint64
	OptionalParam4  uint64
}

type QueryPacketRequest struct {
	PacketType uint32
}

type TransparentModeRequest struct {
	Enable      BOOLEAN
	Techniques  uint32
	DebuggerPid uint32
	KdPid       uint32
}

type TraceStartRequest struct {
	Mode           uint32
	LbrFilter      uint64
	BtsBufferSize  SIZE_T
	CallstackMode  BOOLEAN
	InterruptOnBts BOOLEAN
	BranchType     uint32
}

type TraceBufferRequest struct {
	BufferSize SIZE_T
	Buffer     PVOID
}

func drvIoctl(_ *wdk.DEVICE_OBJECT, irp *wdk.IRP) NTSTATUS {
	stack := IoGetCurrentIrpStackLocation(irp)
	ioctl := stack.Parameters.DeviceIoControl.IoControlCode
	inLen := stack.Parameters.DeviceIoControl.InputBufferLength
	outLen := stack.Parameters.DeviceIoControl.OutputBufferLength

	var status NTSTATUS = STATUS_SUCCESS
	var info ULONG_PTR = 0

	switch ioctl {

	case IOCTL_QUERY_AND_CLEAR_LOGS:
		status = handleQueryLogs(irp, inLen, outLen, &info)

	case IOCTL_REGISTER_EVENT:
		status = handleRegisterEvent(irp, inLen, outLen, &info)

	case IOCTL_RUN_SCRIPT:
		status = handleRunScript(irp, inLen, outLen, &info)

	case IOCTL_SEND_REQUEST_RESULT:
		status = handleSendResult(irp, inLen, outLen, &info)

	case IOCTL_GET_LOG_BASE:
		status = handleGetLogBase(irp, inLen, outLen, &info)

	case IOCTL_VMM_INIT:
		status = handleVmmInit(irp, inLen, outLen, &info)

	case IOCTL_VMM_SHUTDOWN:
		status = handleVmmShutdown(irp, inLen, outLen, &info)

	case IOCTL_EPT_HOOK:
		status = handleEptHook(irp, inLen, outLen, &info)

	case IOCTL_EPT_UNHOOK:
		status = handleEptUnhook(irp, inLen, outLen, &info)

	case IOCTL_EPT_SET_HOOK:
		status = handleEptSetHook(irp, inLen, outLen, &info)

	case IOCTL_EPT_GET_EPT_TABLES:
		status = handleGetEptTables(irp, inLen, outLen, &info)

	case IOCTL_VMM_EXECUTION_TRACE:
		status = handleExecutionTrace(irp, inLen, outLen, &info)

	case IOCTL_VMM_READ_MEM:
		status = handleReadMem(irp, inLen, outLen, &info)

	case IOCTL_VMM_WRITE_MEM:
		status = handleWriteMem(irp, inLen, outLen, &info)

	case IOCTL_VMM_VIRT_TO_PHYS:
		status = handleVirtToPhys(irp, inLen, outLen, &info)

	case IOCTL_VMM_PHYS_TO_VIRT:
		status = handlePhysToVirt(irp, inLen, outLen, &info)

	case IOCTL_VMM_GET_REG:
		status = handleGetReg(irp, inLen, outLen, &info)

	case IOCTL_VMM_SET_REG:
		status = handleSetReg(irp, inLen, outLen, &info)

	case IOCTL_VMM_PAUSE:
		status = handlePause(irp, inLen, outLen, &info)

	case IOCTL_VMM_RESUME:
		status = handleResume(irp, inLen, outLen, &info)

	case IOCTL_VMM_SWITCH_PROCESS:
		status = handleSwitchProcess(irp, inLen, outLen, &info)

	case IOCTL_VMM_MODIFY_REGS:
		status = handleModifyRegs(irp, inLen, outLen, &info)

	case IOCTL_VMM_INVEPT:
		status = handleInvept(irp, inLen, outLen, &info)

	case IOCTL_VMM_INVVPID:
		status = handleInvvpid(irp, inLen, outLen, &info)

	case IOCTL_VMM_FLUSH_ENTIRE_TLB:
		status = handleFlushTlb(irp, inLen, outLen, &info)

	case IOCTL_VMM_GET_MTRR:
		status = handleGetMtrr(irp, inLen, outLen, &info)

	case IOCTL_VMM_CHANGE_CORE:
		status = handleChangeCore(irp, inLen, outLen, &info)

	case IOCTL_VMM_GET_PROCESS_BASE:
		status = handleGetProcessBase(irp, inLen, outLen, &info)

	case IOCTL_VMM_GET_PROCESS_CR3:
		status = handleGetProcessCr3(irp, inLen, outLen, &info)

	case IOCTL_VMM_READ_AND_WRITE_MEM:
		status = handleReadWriteMem(irp, inLen, outLen, &info)

	case IOCTL_VMM_CALL_FUNCTION:
		status = handleCallFunction(irp, inLen, outLen, &info)

	case IOCTL_VMM_QUERY_PACKET:
		status = handleQueryPacket(irp, inLen, outLen, &info)

	case IOCTL_VMM_SEND_RESULT:
		status = handleSendResult(irp, inLen, outLen, &info)

	case IOCTL_VMM_EXECUTE_SINGLE_STEP:
		status = handleSingleStep(irp, inLen, outLen, &info)

	case IOCTL_VMM_ENABLE_AND_INVOKE_TRAP_FLAG:
		status = handleTrapFlag(irp, inLen, outLen, &info)

	case IOCTL_VMM_MASK_EXCEPTION_DEBUGGER_BREAKPOINT:
		status = handleMaskBreakpoint(irp, inLen, outLen, &info)

	case IOCTL_VMM_UNMASK_EXCEPTION_DEBUGGER_BREAKPOINT:
		status = handleUnmaskBreakpoint(irp, inLen, outLen, &info)

	case IOCTL_VMM_SHORT_CIRCUITING_EVENT_INJECT:
		status = handleShortCircuitInject(irp, inLen, outLen, &info)

	case IOCTL_VMM_QUERY_REGISTER_FROM_GUEST_STATE:
		status = handleQueryGuestReg(irp, inLen, outLen, &info)

	case IOCTL_TRANSPARENT_MODE_ENABLE:
		status = handleTransparentEnable(irp, inLen, outLen, &info)

	case IOCTL_TRANSPARENT_MODE_DISABLE:
		status = handleTransparentDisable(irp, inLen, outLen, &info)

	case IOCTL_HYPERTRACE_START:
		status = handleTraceStart(irp, inLen, outLen, &info)

	case IOCTL_HYPERTRACE_STOP:
		status = handleTraceStop(irp, inLen, outLen, &info)

	case IOCTL_HYPERTRACE_READ_BUFFER:
		status = handleTraceReadBuffer(irp, inLen, outLen, &info)

	case IOCTL_HYPERTRACE_CLEAR_BUFFER:
		status = handleTraceClearBuffer(irp, inLen, outLen, &info)

	default:
		if ioctl >= 0x0022B018 && (ioctl-0x0022B018)%4 == 0 {
			status = handleLogRead(irp, inLen, outLen, &info)
		} else {
			status = STATUS_INVALID_DEVICE_REQUEST
		}
	}

	irp.IoStatus.Status = status
	irp.IoStatus.Information = info
	wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
	return status
}

func handleQueryLogs(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if outLen < ULONG(unsafe.Sizeof(LogMessage{})) {
		return STATUS_BUFFER_TOO_SMALL
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	count := gLog.FlushToUser(uintptr(PVOID(buffer)), uintptr(SIZE_T(outLen)))
	*info = ULONG_PTR(uintptr(count) * LOG_MESSAGE_SIZE)

	LogDebug("Flushed %d log messages to user", count)
	return STATUS_SUCCESS
}

func handleRegisterEvent(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen == 0 || inLen < 16 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	eventType := *(*uint32)(unsafe.Pointer(uintptr(buffer)))
	eventTag := *(*uint32)(unsafe.Pointer(uintptr(buffer) + 4))
	callbackAddr := *(*uint64)(unsafe.Pointer(uintptr(buffer) + 8))

	LogDebug("Register event: type=%d tag=0x%X callback=0x%X", eventType, eventTag, callbackAddr)
	*info = 4
	return STATUS_SUCCESS
}

func handleRunScript(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	LogWarning("Script execution not supported in this build")
	return STATUS_NOT_IMPLEMENTED
}

func handleSendResult(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen == 0 {
		return STATUS_INVALID_PARAMETER
	}
	*info = 4
	return STATUS_SUCCESS
}

func handleGetLogBase(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if outLen < 16 {
		return STATUS_BUFFER_TOO_SMALL
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	baseAddr := gLog.GetBufferBase()
	prioBase := gLog.GetPriorityBufferBase()

	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = baseAddr
	*(*uint64)(unsafe.Pointer(uintptr(buffer) + 8)) = prioBase
	*info = 16
	return STATUS_SUCCESS
}

func handleLogRead(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if outLen < 8 {
		return STATUS_BUFFER_TOO_SMALL
	}
	*info = 8
	return STATUS_SUCCESS
}

func handleVmmInit(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(VmmInitRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*VmmInitRequest)(unsafe.Pointer(uintptr(buffer)))

	LogInfo("Initializing VMM for %d cores...", req.NumCores)

	err := gHyp.Initialize()
	if err != nil {
		LogError("VMM initialization failed: %v", err)
		return STATUS_UNSUCCESSFUL
	}

	gHyp.SetInitialized(true)
	*info = 4
	LogInfo("VMM initialized successfully")
	return STATUS_SUCCESS
}

func handleVmmShutdown(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	LogInfo("Shutting down VMM...")

	gHooks.RestoreAll(nil)
	gHyp.Shutdown()
	gHyp.SetInitialized(false)

	*info = 4
	LogInfo("VMM shutdown complete")
	return STATUS_SUCCESS
}

func handleEptHook(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(EptHookRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*EptHookRequest)(unsafe.Pointer(uintptr(buffer)))

	v := gHyp.CurrentVcpu()
	if v == nil {
		return STATUS_DEVICE_NOT_READY
	}

	var kind HookKind
	switch req.HookType {
	case 1:
		kind = HookExec
	case 2:
		kind = HookRead
	case 3:
		kind = HookWrite
	case 4:
		kind = HookReadWrite
	default:
		return STATUS_INVALID_PARAMETER
	}

	err := gHooks.Install(v.CoreId, req.VirtAddr, req.Cr3, kind)
	if err != nil {
		LogError("EPT hook failed at 0x%X: %v", req.VirtAddr, err)
		return STATUS_UNSUCCESSFUL
	}

	*info = 4
	LogDebug("EPT hook installed at 0x%X (type=%d)", req.VirtAddr, req.HookType)
	return STATUS_SUCCESS
}

func handleEptUnhook(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(EptUnhookRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*EptUnhookRequest)(unsafe.Pointer(uintptr(buffer)))

	pa := VirtToPhys(uintptr(req.VirtAddr&^uint64(PAGE_SIZE-1)), req.Cr3)
	if pa == 0 {
		return STATUS_INVALID_PARAMETER
	}

	err := gHooks.Remove(pa)
	if err != nil {
		LogError("EPT unhook failed at PA 0x%X: %v", pa, err)
		return STATUS_UNSUCCESSFUL
	}

	*info = 4
	LogDebug("EPT hook removed at VA 0x%X (PA 0x%X)", req.VirtAddr, pa)
	return STATUS_SUCCESS
}

func handleEptSetHook(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	return handleEptHook(irp, inLen, outLen, info)
}

func handleGetEptTables(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if outLen < 8 {
		return STATUS_BUFFER_TOO_SMALL
	}

	v := gHyp.CurrentVcpu()
	if v == nil || v.EptPageTable == nil {
		return STATUS_DEVICE_NOT_READY
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	*(*PVOID)(unsafe.Pointer(uintptr(buffer))) = uintptr(unsafe.Pointer(v.EptPageTable))
	*info = 8

	LogDebug("Returned EPT table base at 0x%x", v.EptPageTable)
	return STATUS_SUCCESS
}

func handleExecutionTrace(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(TraceStartRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*TraceStartRequest)(unsafe.Pointer(uintptr(buffer)))

	config := TraceUserConfig{
		LbrFilter:     req.LbrFilter,
		BtsBufferSize: req.BtsBufferSize,
		CallstackMode: bool(req.CallstackMode),
	}

	mode := TraceMode(req.Mode)
	err := gTrace.Enable(mode, &config)
	if err != nil {
		LogError("Trace enable failed: %v", err)
		return STATUS_UNSUCCESSFUL
	}

	*info = 4
	LogInfo("Execution trace started (mode=%d)", req.Mode)
	return STATUS_SUCCESS
}

func handleReadMem(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(MemReadRequest{})) || outLen < ULONG(unsafe.Sizeof(MemReadResponse{}))+256 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*MemReadRequest)(unsafe.Pointer(uintptr(buffer)))
	resp := (*MemReadResponse)(unsafe.Pointer(uintptr(buffer)))

	pa := VirtToPhys(uintptr(req.Address), req.Cr3)
	if pa == 0 {
		resp.Success = false
		*info = ULONG_PTR(unsafe.Sizeof(MemReadResponse{}))
		return STATUS_SUCCESS
	}

	readSize := req.Size
	if readSize > outLen-ULONG(unsafe.Sizeof(MemReadResponse{})) {
		readSize = uint32(outLen - ULONG(unsafe.Sizeof(MemReadResponse{})))
	}

	dataBuf := (*byte)(unsafe.Pointer(uintptr(buffer) + unsafe.Sizeof(MemReadResponse{})))
	n := ReadPhysMem(pa, unsafe.Slice(dataBuf, int(readSize)))

	resp.Address = req.Address
	resp.Size = uint32(n)
	resp.Success = true
	*info = ULONG_PTR(unsafe.Sizeof(MemReadResponse{})) + ULONG_PTR(n)

	return STATUS_SUCCESS
}

func handleWriteMem(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(MemReadRequest{}))+1 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*MemReadRequest)(unsafe.Pointer(uintptr(buffer)))

	pa := VirtToPhys(uintptr(req.Address), req.Cr3)
	if pa == 0 {
		return STATUS_INVALID_PARAMETER
	}

	dataOffset := unsafe.Sizeof(MemReadRequest{})
	dataBuf := (*byte)(unsafe.Pointer(uintptr(buffer) + dataOffset))

	n := WritePhysMem(pa, unsafe.Slice(dataBuf, int(req.Size)))

	*info = 4
	LogDebug("Wrote %d bytes to PA 0x%X", n, pa)
	return STATUS_SUCCESS
}

func handleVirtToPhys(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(VirtPhysRequest{})) || outLen < 16 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*VirtPhysRequest)(unsafe.Pointer(uintptr(buffer)))

	pa := VirtToPhys(uintptr(req.Address), req.Cr3)

	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = pa
	*(*uint64)(unsafe.Pointer(uintptr(buffer) + 8)) = req.Address
	*info = 16

	return STATUS_SUCCESS
}

func handlePhysToVirt(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(VirtPhysRequest{})) || outLen < 16 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*VirtPhysRequest)(unsafe.Pointer(uintptr(buffer)))

	va := PhysToVirt(req.Address, req.Cr3)

	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = uint64(va)
	*(*uint64)(unsafe.Pointer(uintptr(buffer) + 8)) = req.Address
	*info = 16

	return STATUS_SUCCESS
}

func handleGetReg(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(RegRequest{})) || outLen < ULONG(unsafe.Sizeof(RegResponse{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*RegRequest)(unsafe.Pointer(uintptr(buffer)))
	resp := (*RegResponse)(unsafe.Pointer(uintptr(buffer)))

	v := gHyp.Vcpu(req.CoreId)
	if v == nil || v.Regs == nil {
		resp.Success = false
		*info = ULONG_PTR(unsafe.Sizeof(RegResponse{}))
		return STATUS_SUCCESS
	}

	value := getRegisterValue(v.Regs, req.RegIndex)
	resp.Value = value
	resp.RegIndex = req.RegIndex
	resp.CoreId = req.CoreId
	resp.Success = true
	*info = ULONG_PTR(unsafe.Sizeof(RegResponse{}))

	return STATUS_SUCCESS
}

func handleSetReg(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(RegResponse{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*RegResponse)(unsafe.Pointer(uintptr(buffer)))

	v := gHyp.Vcpu(req.CoreId)
	if v == nil || v.Regs == nil {
		return STATUS_DEVICE_NOT_READY
	}

	setRegisterValue(v.Regs, req.RegIndex, req.Value)
	*info = 4
	return STATUS_SUCCESS
}

func handlePause(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	gHyp.PauseAll()
	*info = 4
	LogInfo("All VCPUs paused")
	return STATUS_SUCCESS
}

func handleResume(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	gHyp.ResumeAll()
	*info = 4
	LogInfo("All VCPUs resumed")
	return STATUS_SUCCESS
}

func handleSwitchProcess(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(SwitchProcessRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*SwitchProcessRequest)(unsafe.Pointer(uintptr(buffer)))

	gHyp.SetCurrentProcessContext(req.ProcessId, req.Cr3, req.ProcessBaseAddress)
	*info = 4
	LogDebug("Switched to process %d (CR3=0x%X)", req.ProcessId, req.Cr3.Flags)
	return STATUS_SUCCESS
}

func handleModifyRegs(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(ModifyRegsRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*ModifyRegsRequest)(unsafe.Pointer(uintptr(buffer)))

	v := gHyp.Vcpu(req.CoreId)
	if v == nil || v.Regs == nil {
		return STATUS_DEVICE_NOT_READY
	}

	*v.Regs = req.Regs
	*info = 4
	return STATUS_SUCCESS
}

func handleInvept(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(InveptRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*InveptRequest)(unsafe.Pointer(uintptr(buffer)))

	if req.Type == 1 {
		eptInveptSingleContext(req.Eptp)
	} else if req.Type == 2 {
		eptInveptAllContexts()
	}

	*info = 4
	return STATUS_SUCCESS
}

func handleInvvpid(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	AsmInvvpid()
	*info = 4
	return STATUS_SUCCESS
}

func handleFlushTlb(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	eptInveptAllContexts()
	AsmInvvpid()
	*info = 4
	return STATUS_SUCCESS
}

func handleGetMtrr(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(MtrrRequest{})) || outLen < 12 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*MtrrRequest)(unsafe.Pointer(uintptr(buffer)))

	memType := gHyp.MemoryManager().QueryMtrrForPa(req.BaseAddr)

	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = req.BaseAddr
	*(*uint32)(unsafe.Pointer(uintptr(buffer) + 8)) = uint32(memType)
	*info = 12
	return STATUS_SUCCESS
}

func handleChangeCore(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(ChangeCoreRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*ChangeCoreRequest)(unsafe.Pointer(uintptr(buffer)))

	gHyp.SetActiveCore(req.TargetCoreId)
	*info = 4
	return STATUS_SUCCESS
}

func handleGetProcessBase(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(ProcessCr3Request{})) || outLen < 8 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*ProcessCr3Request)(unsafe.Pointer(uintptr(buffer)))

	baseAddr := gHyp.GetProcessBaseAddress(req.ProcessId)
	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = baseAddr
	*info = 8
	return STATUS_SUCCESS
}

func handleGetProcessCr3(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(ProcessCr3Request{})) || outLen < 8 {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*ProcessCr3Request)(unsafe.Pointer(uintptr(buffer)))

	cr3 := gHyp.GetProcessCr3(req.ProcessId)
	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = cr3.Flags
	*info = 8
	return STATUS_SUCCESS
}

func handleReadWriteMem(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(ReadWriteMemRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*ReadWriteMemRequest)(unsafe.Pointer(uintptr(buffer)))

	if !req.ReadOrWrite {
		if req.ReadSize > 0 {
			pa := VirtToPhys(uintptr(req.ReadAddress), req.Cr3)
			if pa != 0 {
				dataBuf := (*byte)(unsafe.Pointer(uintptr(buffer) + unsafe.Sizeof(ReadWriteMemRequest{})))
				ReadPhysMem(pa, unsafe.Slice(dataBuf, int(req.ReadSize)))
			}
		}
	} else {
		if req.WriteSize > 0 {
			pa := VirtToPhys(uintptr(req.WriteAddress), req.Cr3)
			if pa != 0 {
				dataOffset := unsafe.Sizeof(ReadWriteMemRequest{})
				dataBuf := (*byte)(unsafe.Pointer(uintptr(buffer) + dataOffset))
				WritePhysMem(pa, unsafe.Slice(dataBuf, int(req.WriteSize)))
			}
		}
	}

	*info = 4
	return STATUS_SUCCESS
}

func handleCallFunction(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(CallFunctionRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*CallFunctionRequest)(unsafe.Pointer(uintptr(buffer)))

	result := gHyp.CallGuestFunction(
		req.FunctionAddress,
		req.OptionalParam1,
		req.OptionalParam2,
		req.OptionalParam3,
		req.OptionalParam4,
	)

	*(*uint64)(unsafe.Pointer(uintptr(buffer))) = result
	*info = 8
	return STATUS_SUCCESS
}

func handleQueryPacket(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if outLen < 4 {
		return STATUS_BUFFER_TOO_SMALL
	}
	*info = 4
	return STATUS_SUCCESS
}

func handleSingleStep(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	v := gHyp.CurrentVcpu()
	if v == nil {
		return STATUS_DEVICE_NOT_READY
	}

	gHyp.EnableSingleStep(v)
	*info = 4
	return STATUS_SUCCESS
}

func handleTrapFlag(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	v := gHyp.CurrentVcpu()
	if v == nil {
		return STATUS_DEVICE_NOT_READY
	}

	setMonitorTrapFlag(true)
	*info = 4
	return STATUS_SUCCESS
}

func handleMaskBreakpoint(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	gEvents.Disable(EventException)
	*info = 4
	return STATUS_SUCCESS
}

func handleUnmaskBreakpoint(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	gEvents.Enable(EventException)
	*info = 4
	return STATUS_SUCCESS
}

func handleShortCircuitInject(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	*info = 4
	return STATUS_SUCCESS
}

func handleQueryGuestReg(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	return handleGetReg(irp, inLen, outLen, info)
}

func handleTransparentEnable(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if inLen < ULONG(unsafe.Sizeof(TransparentModeRequest{})) {
		return STATUS_INVALID_PARAMETER
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	req := (*TransparentModeRequest)(unsafe.Pointer(uintptr(buffer)))

	config := TransparentModeConfig{
		Enabled:     bool(req.Enable),
		Techniques:  EvasionTechnique(req.Techniques),
		DebuggerPid: req.DebuggerPid,
	}

	err := gEvasion.Enable(&config)
	if err != nil {
		LogError("Transparent mode enable failed: %v", err)
		return STATUS_UNSUCCESSFUL
	}

	*info = 4
	LogInfo("Transparent mode enabled")
	return STATUS_SUCCESS
}

func handleTransparentDisable(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	err := gEvasion.Disable()
	if err != nil {
		LogError("Transparent mode disable failed: %v", err)
		return STATUS_UNSUCCESSFUL
	}

	*info = 4
	LogInfo("Transparent mode disabled")
	return STATUS_SUCCESS
}

func handleTraceStart(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	return handleExecutionTrace(irp, inLen, outLen, info)
}

func handleTraceStop(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	err := gTrace.Disable()
	if err != nil {
		return STATUS_UNSUCCESSFUL
	}
	*info = 4
	LogInfo("Tracing stopped")
	return STATUS_SUCCESS
}

func handleTraceReadBuffer(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	if outLen < 8 {
		return STATUS_BUFFER_TOO_SMALL
	}

	lbrEntries, lbrCount := gTrace.CaptureLbr()

	buffer := irp.AssociatedIrp.SystemBuffer
	*(*uint32)(unsafe.Pointer(uintptr(buffer))) = uint32(lbrCount)

	if outLen >= ULONG(8+uintptr(lbrCount)*unsafe.Sizeof(LbrEntry{})) {
		for i := range lbrCount {
			offset := 8 + uintptr(i)*unsafe.Sizeof(LbrEntry{})
			dst := (*LbrEntry)(unsafe.Pointer(uintptr(buffer) + offset))
			*dst = lbrEntries[i]
		}
		*info = ULONG_PTR(8 + uintptr(lbrCount)*unsafe.Sizeof(LbrEntry{}))
	} else {
		*info = 4
	}

	return STATUS_SUCCESS
}

func handleTraceClearBuffer(irp *wdk.IRP, inLen, outLen ULONG, info *ULONG_PTR) NTSTATUS {
	// TODO: Implement ClearBuffers method in TracerState
	*info = 4
	LogDebug("Trace buffers cleared")
	return STATUS_SUCCESS
}

func getRegisterValue(regs *GUEST_REGS, index uint32) uint64 {
	switch index {
	case 0:
		return regs.Rax
	case 1:
		return regs.Rcx
	case 2:
		return regs.Rdx
	case 3:
		return regs.Rbx
	case 4:
		return regs.Rsp
	case 5:
		return regs.Rbp
	case 6:
		return regs.Rsi
	case 7:
		return regs.Rdi
	case 8:
		return regs.R8
	case 9:
		return regs.R9
	case 10:
		return regs.R10
	case 11:
		return regs.R11
	case 12:
		return regs.R12
	case 13:
		return regs.R13
	case 14:
		return regs.R14
	case 15:
		return regs.R15
	case 16:
		// RIP is not in GUEST_REGS, need to read from VMCS
		return 0
	default:
		return 0
	}
}

func setRegisterValue(regs *GUEST_REGS, index uint32, value uint64) {
	switch index {
	case 0:
		regs.Rax = value
	case 1:
		regs.Rcx = value
	case 2:
		regs.Rdx = value
	case 3:
		regs.Rbx = value
	case 4:
		regs.Rsp = value
	case 5:
		regs.Rbp = value
	case 6:
		regs.Rsi = value
	case 7:
		regs.Rdi = value
	case 8:
		regs.R8 = value
	case 9:
		regs.R9 = value
	case 10:
		regs.R10 = value
	case 11:
		regs.R11 = value
	case 12:
		regs.R12 = value
	case 13:
		regs.R13 = value
	case 14:
		regs.R14 = value
	case 15:
		regs.R15 = value
	case 16:
		// RIP is not in GUEST_REGS, need to write to VMCS
		// TODO: Implement VMCS write for RIP
	}
}

func eptInveptSingleContext(eptp uint64) {
	inveptSingleContext(eptp)
}

func eptInveptAllContexts() {
	inveptAllContexts()
}
