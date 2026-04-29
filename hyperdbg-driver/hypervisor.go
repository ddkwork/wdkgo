package main

import (
	"unsafe"

	"solod.dev/so/wdk"
)

type ExitReason = uint32

const (
	ExitTripleFault  ExitReason = 2
	ExitVmclear      ExitReason = 0x13
	ExitVmlaunch     ExitReason = 0x14
	ExitVmptrld      ExitReason = 0x15
	ExitVmptrst      ExitReason = 0x16
	ExitVmread       ExitReason = 0x17
	ExitVmresume     ExitReason = 0x18
	ExitVmwrite      ExitReason = 0x19
	ExitVmxoff       ExitReason = 0x1A
	ExitVmxon        ExitReason = 0x1B
	ExitInvd         ExitReason = 0x0D
	ExitGetsec       ExitReason = 0x0B
	ExitInvept       ExitReason = 0x32
	ExitInvvpid      ExitReason = 0x35
	ExitCpuid        ExitReason = 22
	ExitHlt          ExitReason = 24
	ExitRdtsc        ExitReason = 16
	ExitRdpmc        ExitReason = 15
	ExitVmcall       ExitReason = 18
	ExitMovCr        ExitReason = 28
	ExitMovDr        ExitReason = 29
	ExitIo           ExitReason = 30
	ExitRdmsr        ExitReason = 31
	ExitWrmsr        ExitReason = 32
	ExitExceptionNmi ExitReason = 0
	ExitExtInt       ExitReason = 1
	ExitIntWindow    ExitReason = 7
	ExitNmiWindow    ExitReason = 8
	ExitMtf          ExitReason = 37
	ExitEptViolation ExitReason = 48
	ExitEptMisconfig ExitReason = 49
	ExitRdtscp       ExitReason = 51
	ExitXsetbv       ExitReason = 55
	ExitPreemptTimer ExitReason = 52
)

const (
	IntTypeHwException uint32 = 3
	IntTypeSwException uint32 = 6
	IntTypeNmi         uint32 = 2

	VectorBp  uint32 = 3
	VectorGp  uint32 = 13
	VectorUd  uint32 = 6
	VectorNmi uint32 = 2
)

const (
	VmcsExitReason           uint64 = 0x00004400
	VmcsExitQualification    uint64 = 0x00004402
	VmcsGuestRip             uint64 = 0x0000681E
	VmcsGuestRsp             uint64 = 0x0000681C
	VmcsGuestIA32Debugctl    uint64 = 0x00002802
	VmcsVmexitInstrLength    uint64 = 0x0000440C
	VmcsCtrlVmentryIntrInfo  uint64 = 0x00004016
	VmcsCtrlVmentryErrCode   uint64 = 0x00004018
	VmcsCtrlVmentryInstrLen  uint64 = 0x0000401A
	VmcsGuestPhysicalAddress uint64 = 0x00002400

	VmcsGuestCr0     uint64 = 0x00006800
	VmcsGuestCr3     uint64 = 0x00006802
	VmcsGuestCr4     uint64 = 0x00006804
	VmcsGuestDr7     uint64 = 0x0000681A
	VmcsGuestRflags  uint64 = 0x00006820
	VmcsGuestSsp     uint64 = 0x0000682A
	VmcsGuestSymEFER uint64 = 0x00006826
)

type ProcessContext struct {
	ProcessId   uint32
	Cr3         CR3_TYPE
	BaseAddress uint64
	ImageName   [256]CHAR
}

type Hypervisor struct {
	vcpus           []VCPU
	eptTables       []*EptTable
	events          *EventDispatcher
	hooks           *HookManager
	logger          *Logger
	memory          *MemoryManager
	processContexts map[uint32]*ProcessContext
	activeCoreId    uint32
	checkFootprints bool
	initialized     bool
	paused          bool
	lock            *Spinlock
}

func NewHypervisor(numCpus uint32) *Hypervisor {
	h := &Hypervisor{
		vcpus:           make([]VCPU, numCpus),
		eptTables:       make([]*EptTable, numCpus),
		events:          newEventDispatcher(),
		hooks:           NewHookManager(int(MAX_HIDDEN_BREAKPOINTS)),
		logger:          NewLogger(uint32(LOG_CHUNK_SIZE), uint32(LOG_CHUNK_SIZE)),
		memory:          newMemoryManager(),
		processContexts: make(map[uint32]*ProcessContext),
		lock:            NewSpinlock(),
	}

	for i := range numCpus {
		h.vcpus[i].CoreId = i
		h.vcpus[i].IncrementRip = true
		h.eptTables[i] = NewEptTable()
	}

	return h
}

func (h *Hypervisor) Initialize() error {
	if h.initialized {
		return nil
	}

	count := wdk.KeQueryActiveProcessorCount(nil)
	for coreId := range count {
		if err := h.initializeVcpu(coreId); err != nil {
			LogError("Failed to init VCPU %d: %v", coreId, err)
			return err
		}
	}

	h.initialized = true
	LogInfo("Hypervisor initialized with %d VCPUs", count)
	return nil
}

func (h *Hypervisor) initializeVcpu(coreId uint32) error {
	v := &h.vcpus[coreId]

	vmcsPa := allocateVmcsRegion()
	if vmcsPa == 0 {
		return ErrVmcsAlloc
	}
	v.VmcsRegionPhysicalAddress = vmcsPa
	v.VmcsRegionVirtualAddress = vmcsPa

	vmxonPa := allocateVmxonRegion()
	if vmxonPa == 0 {
		return ErrVmxonAlloc
	}
	v.VmxonRegionPhysicalAddress = vmxonPa
	v.VmxonRegionVirtualAddress = vmxonPa

	msrBitmapPa := allocateMsrBitmap()
	if msrBitmapPa == 0 {
		return fmtError("MSR bitmap alloc failed")
	}
	v.MsrBitmapPhysicalAddress = msrBitmapPa
	v.MsrBitmapVirtualAddress = msrBitmapPa

	ioBitmapPaA := allocateIoBitmap()
	ioBitmapPaB := allocateIoBitmap()
	v.IoBitmapPhysicalAddressA = ioBitmapPaA
	v.IoBitmapVirtualAddressA = ioBitmapPaA
	v.IoBitmapPhysicalAddressB = ioBitmapPaB
	v.IoBitmapVirtualAddressB = ioBitmapPaB

	stackPa := allocateVmmStack()
	if stackPa == 0 {
		return fmtError("stack alloc failed")
	}
	v.VmmStack = stackPa

	if h.eptTables[coreId] != nil {
		eptPtr := h.eptTables[coreId].BuildEptPointer()
		v.EptPointer.AsUInt = eptPtr
	} else {
		h.eptTables[coreId] = NewEptTable()
		eptPtr := h.eptTables[coreId].BuildEptPointer()
		v.EptPointer.AsUInt = eptPtr
	}

	if err := vmxTurnOn(v.VmxonRegionPhysicalAddress); err != nil {
		return fmtError("VMXON failed on core %d: %v", coreId, err)
	}

	if err := vmClear(v.VmcsRegionPhysicalAddress); err != nil {
		return fmtError("VMCLEAR failed on core %d: %v", coreId, err)
	}

	if err := vmLoad(v.VmcsRegionPhysicalAddress); err != nil {
		return fmtError("VMPTRLD failed on core %d: %v", coreId, err)
	}

	if err := h.setupVmcs(v); err != nil {
		return fmtError("VMCS setup failed on core %d: %v", coreId, err)
	}

	LogInfo("VCPU %d initialized: VMCS=0x%X VMXON=0x%X EPTP=0x%X",
		coreId, vmcsPa, vmxonPa, v.EptPointer.AsUInt)

	return nil
}

func (h *Hypervisor) Shutdown() {
	h.lock.Lock()
	defer h.lock.Unlock()

	if !h.initialized {
		return
	}

	for _, v := range h.vcpus {
		gHooks.RestoreAllForCore(v.CoreId)
		if v.OnVmxRootMode {
			vmClear(v.VmcsRegionPhysicalAddress)
		}
	}

	vmxTurnOff()
	h.initialized = false
	LogInfo("Hypervisor shutdown complete")
}

func (h *Hypervisor) SetInitialized(val bool) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.initialized = val
}

func (h *Hypervisor) IsInitialized() bool {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.initialized
}

func (h *Hypervisor) eptTable(coreId uint32) *EptTable {
	if int(coreId) >= len(h.eptTables) {
		return nil
	}
	return h.eptTables[coreId]
}

func (h *Hypervisor) HandleVmExit(regs *GUEST_REGS) bool {
	coreId := getCurrentProcessorNumber()
	if int(coreId) >= len(h.vcpus) {
		return false
	}

	v := &h.vcpus[coreId]
	v.OnVmxRootMode = true
	v.Regs = regs
	v.XmmRegs = (*GUEST_XMM_REGS)(unsafe.Pointer(
		uintptr(unsafe.Pointer(regs)) + unsafe.Sizeof(GUEST_REGS{})))

	exitReason := readExitReason()
	v.ExitReason = uint32(exitReason)
	v.IncrementRip = true

	v.LastVmexitRip = readGuestRip()
	var rsp uint64
	vmRead64(VmcsGuestRsp, &rsp)
	v.Regs.Rsp = rsp
	v.ExitQualification = readExitQualification()

	_ = h.dispatch(v)

	if !v.VmxoffState.Executed && v.IncrementRip {
		if h.checkFootprints {
			h.checkTrapFlag(v)
		}
		h.advanceIp(v)
	}

	done := v.VmxoffState.Executed
	v.OnVmxRootMode = false
	return done
}

func (h *Hypervisor) dispatch(v *VCPU) bool {
	switch v.ExitReason {

	case ExitExceptionNmi:
		return h.handleException(v)

	case ExitExtInt:
		return h.handleExternalInt(v)

	case ExitCpuid:
		return h.events.Dispatch(EventCpuid, v)

	case ExitRdmsr:
		return h.events.Dispatch(EventRdmsr, v)

	case ExitWrmsr:
		return h.events.Dispatch(EventWrmsr, v)

	case ExitIo:
		return h.events.Dispatch(EventIo, v)

	case ExitMovCr:
		return h.events.Dispatch(EventMovCr, v)

	case ExitMovDr:
		return h.events.Dispatch(EventMovDr, v)

	case ExitEptViolation:
		return h.handleEptViolation(v)

	case ExitEptMisconfig:
		return h.handleEptMisconfig(v)

	case ExitVmcall:
		return h.handleVmcall(v)

	case ExitRdtsc, ExitRdtscp:
		return h.events.Dispatch(EventTsc, v)

	case ExitRdpmc:
		return h.events.Dispatch(EventRdpmc, v)

	case ExitXsetbv:
		return h.events.Dispatch(EventXsetbv, v)

	case ExitIntWindow:
		return h.handleInterruptWindow(v)

	case ExitNmiWindow:
		return h.handleNmiWindow(v)

	case ExitMtf:
		return h.handleMonitorTrapFlag(v)

	case ExitPreemptTimer:
		return h.handlePreemptionTimer(v)

	case ExitTripleFault:
		return h.handleTripleFault(v)

	case ExitVmclear, ExitVmlaunch, ExitVmptrld,
		ExitVmread, ExitVmresume, ExitVmwrite, ExitVmxoff,
		ExitVmxon, ExitInvept, ExitInvvpid, ExitGetsec, ExitInvd:
		h.injectUd(v)
		return false

	default:
		return false
	}
}

func (h *Hypervisor) handleException(v *VCPU) bool {
	vec := v.ExitQualification & 0xFF
	intType := (v.ExitQualification >> 8) & 0x7
	hasErrCode := (v.ExitQualification>>11)&0x1 != 0

	switch vec {
	case VectorBp:
		if intType == IntTypeSwException && !hasErrCode {
			return gHooks.HandleBreakpoint(v)
		}
	case VectorGp:
		return h.handleGeneralProtectionFault(v)
	case VectorUd:
		return h.handleUndefinedOpcode(v)
	case VectorNmi:
		if intType == IntTypeNmi {
			return h.handleNmi(v)
		}
	}

	return h.events.Dispatch(EventException, v)
}

func (h *Hypervisor) handleEptViolation(v *VCPU) bool {
	qual := v.ExitQualification
	read := qual&0x1 != 0
	write := qual&0x2 != 0
	exec := qual&0x4 != 0

	gpa := readGuestPhysicalAddr()

	if exec && !read && !write {
		return gHooks.HandleExecHook(v, gpa)
	}

	if write || read {
		return gHooks.HandleReadWriteHook(v, gpa, read, write)
	}

	return false
}

func (h *Hypervisor) handleVmcall(v *VCPU) bool {
	isOurs := v.Regs.Rax == HYPERDBG_VMCALL_MAGIC_RAX &&
		v.Regs.Rcx == HYPERDBG_VMCALL_MAGIC_RCX &&
		v.Regs.Rdx == HYPERDBG_VMCALL_MAGIC_RDX

	if isOurs {
		num := v.Regs.R8
		p1 := v.Regs.R9
		p2 := v.Regs.R10
		p3 := v.Regs.R11
		v.Regs.Rax = uint64(h.handleHyperdbgVmcall(num, p1, p2, p3))
		return true
	}

	h.hypercallForward(v)
	return true
}

func (h *Hypervisor) hypercallForward(v *VCPU) {
	rsp := v.Regs.Rsp
	asmHypervVmcall(uint64(uintptr(unsafe.Pointer(v.Regs))))
	v.Regs.Rsp = rsp
}

func (h *Hypervisor) advanceIp(v *VCPU) {
	instrLen := readInstructionLength()
	newRip := v.LastVmexitRip + uint64(instrLen)
	vmWrite64(VmcsGuestRip, newRip)
}

func (h *Hypervisor) suppressAdvance(v *VCPU) { v.IncrementRip = false }
func (h *Hypervisor) enableAdvance(v *VCPU)   { v.IncrementRip = true }

func (h *Hypervisor) injectInterrupt(intType, vector uint32, deliverErr bool, errCode uint32) {
	val := uint32(0x80000000) | (intType << 8) | (vector & 0xFF)
	if deliverErr {
		val |= 0x800
	}
	vmWrite64(VmcsCtrlVmentryIntrInfo, uint64(val))
	if deliverErr {
		vmWrite64(VmcsCtrlVmentryErrCode, uint64(errCode))
	}
}

func (h *Hypervisor) injectBp(v *VCPU) {
	h.injectInterrupt(IntTypeSwException, VectorBp, false, 0)
	var len uint32
	vmRead32(VmcsVmexitInstrLength, &len)
	vmWrite64(VmcsCtrlVmentryInstrLen, uint64(len))
}

func (h *Hypervisor) injectGp(v *VCPU) {
	h.injectInterrupt(IntTypeHwException, VectorGp, true, 0)
	var len uint32
	vmRead32(VmcsVmexitInstrLength, &len)
	vmWrite64(VmcsCtrlVmentryInstrLen, uint64(len))
}

func (h *Hypervisor) injectUd(v *VCPU) {
	h.injectInterrupt(IntTypeHwException, VectorUd, false, 0)
	h.suppressAdvance(v)
}

func (h *Hypervisor) handleMonitorTrapFlag(v *VCPU) bool {
	h.suppressAdvance(v)
	v.IgnoreMtfUnset = false

	if gHooks.HandleMtfRestore(v) {
		h.enableAndCheckPendingInt(v)
	} else if v.RegisterBreakOnMtf {
		v.RegisterBreakOnMtf = false
		handleRegisteredMtf(v.CoreId)
	} else if kdHandleNmiCallback(v.CoreId) {
	} else if v.IgnoreOneMtf {
		v.IgnoreOneMtf = false
	}

	if !v.IgnoreMtfUnset {
		setMonitorTrapFlag(false)
	} else {
		v.IgnoreMtfUnset = false
	}

	return true
}

func (h *Hypervisor) enableAndCheckPendingInt(_ *VCPU)          {}
func (h *Hypervisor) checkTrapFlag(_ *VCPU)                     {}
func (h *Hypervisor) handleTripleFault(_ *VCPU) bool            { return false }
func (h *Hypervisor) handleEptMisconfig(_ *VCPU) bool           { return false }
func (h *Hypervisor) handleGeneralProtectionFault(_ *VCPU) bool { return false }
func (h *Hypervisor) handleUndefinedOpcode(_ *VCPU) bool        { return false }
func (h *Hypervisor) handleExternalInt(_ *VCPU) bool            { return false }
func (h *Hypervisor) handleNmi(_ *VCPU) bool                    { return false }
func (h *Hypervisor) handleInterruptWindow(_ *VCPU) bool        { return false }
func (h *Hypervisor) handleNmiWindow(_ *VCPU) bool              { return false }
func (h *Hypervisor) handlePreemptionTimer(_ *VCPU) bool        { return false }

func (h *Hypervisor) transparentUnhide() {}
func (h *Hypervisor) transparentHide()   {}

func (h *Hypervisor) Vcpu(id uint32) *VCPU {
	if int(id) >= len(h.vcpus) {
		return nil
	}
	return &h.vcpus[id]
}

func (h *Hypervisor) CurrentVcpu() *VCPU {
	id := getCurrentProcessorNumber()
	if int(id) >= len(h.vcpus) {
		return nil
	}
	return &h.vcpus[id]
}

func (h *Hypervisor) MemoryManager() *MemoryManager { return h.memory }

func (h *Hypervisor) SetActiveCore(id uint32) { h.activeCoreId = id }
func (h *Hypervisor) GetActiveCore() uint32   { return h.activeCoreId }
func (h *Hypervisor) PauseAll()               { h.paused = true }
func (h *Hypervisor) ResumeAll()              { h.paused = false }
func (h *Hypervisor) IsPaused() bool          { return h.paused }

func (h *Hypervisor) EnableSingleStep(_ *VCPU) { setMonitorTrapFlag(true) }

func (h *Hypervisor) SetCurrentProcessContext(pid uint32, cr3 CR3_TYPE, baseAddr uint64) {
	h.lock.Lock()
	defer h.lock.Unlock()

	ctx, ok := h.processContexts[pid]
	if !ok {
		ctx = &ProcessContext{ProcessId: pid}
		h.processContexts[pid] = ctx
	}
	ctx.Cr3 = cr3
	ctx.BaseAddress = baseAddr
}

func (h *Hypervisor) GetProcessBaseAddress(pid uint32) uint64 {
	h.lock.Lock()
	defer h.lock.Unlock()
	if ctx, ok := h.processContexts[pid]; ok {
		return ctx.BaseAddress
	}
	return 0
}

func (h *Hypervisor) GetProcessCr3(pid uint32) CR3_TYPE {
	h.lock.Lock()
	defer h.lock.Unlock()
	if ctx, ok := h.processContexts[pid]; ok {
		return ctx.Cr3
	}
	return CR3_TYPE{}
}

func (h *Hypervisor) CallGuestFunction(_, _, _, _, _ uint64) uint64 { return 0 }

func readExitReason() ExitReason {
	var reason uint32
	vmRead32(VmcsExitReason, &reason)
	reason &= 0xFFFF
	return ExitReason(reason)
}

func readExitQualification() uint32 {
	var qual uint32
	vmRead32(VmcsExitQualification, &qual)
	return qual
}

func readGuestRip() uint64 {
	var rip uint64
	vmRead64(VmcsGuestRip, &rip)
	return rip
}

func readInstructionLength() uint32 {
	var length uint32
	vmRead32(VmcsVmexitInstrLength, &length)
	return length
}

func readGuestPhysicalAddr() uint64 {
	var addr uint64
	vmRead64(VmcsGuestPhysicalAddress, &addr)
	return addr
}

func allocateVmcsRegion() uint64 {
	vmcsSize := max(__readmsr(0x480)&0x1FFF, 4096)
	vmcsSize = (vmcsSize + 0xFFF) & ^uint64(0xFFF)

	ptr := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), uintptr(vmcsSize), POOLTAG)
	if ptr == 0 {
		return 0
	}
	wdk.RtlZeroMemory(ptr, uintptr(vmcsSize))

	pa := uint64(wdk.MmGetPhysicalAddress(ptr).QuadPart)
	LogDebug("Allocated VMCS: VA=0x%X PA=0x%X Size=0x%X", ptr, pa, vmcsSize)
	return pa
}

func allocateVmxonRegion() uint64 {
	vmxonSize := max(__readmsr(0x480)&0x1FFF, 4096)
	vmxonSize = (vmxonSize + 0xFFF) & ^uint64(0xFFF)

	ptr := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), uintptr(vmxonSize), POOLTAG)
	if ptr == 0 {
		return 0
	}
	wdk.RtlZeroMemory(ptr, uintptr(vmxonSize))

	revisonId := __readmsr(0x480)
	*(*uint32)(unsafe.Pointer(uintptr(ptr))) = uint32(revisonId)

	pa := uint64(wdk.MmGetPhysicalAddress(ptr).QuadPart)
	LogDebug("Allocated VMXON: VA=0x%X PA=0x%X RevId=0x%X", ptr, pa, revisonId)
	return pa
}

func allocateMsrBitmap() uint64 {
	size := SIZE_T(0x2000)
	ptr := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), uintptr(size), POOLTAG)
	if ptr == 0 {
		return 0
	}
	wdk.RtlZeroMemory(ptr, uintptr(size))

	bitmap := (*[0x1000]byte)(unsafe.Pointer(uintptr(ptr)))
	for i := range 0x800 {
		bitmap[i] = 0xFF
	}
	bitmap[0xC00/8] &= ^byte(1 << (0xC00 % 8))
	bitmap[0xC01/8] &= ^byte(1 << (0xC01 % 8))
	bitmap[0xC02/8] &= ^byte(1 << (0xC02 % 8))
	bitmap[0xC03/8] &= ^byte(1 << (0xC03 % 8))
	bitmap[0xC04/8] &= ^byte(1 << (0xC04 % 8))
	bitmap[0xC05/8] &= ^byte(1 << (0xC05 % 8))
	bitmap[0xC06/8] &= ^byte(1 << (0xC06 % 8))
	bitmap[0xC07/8] &= ^byte(1 << (0xC07 % 8))
	bitmap[0xC08/8] &= ^byte(1 << (0xC08 % 8))
	bitmap[0x174/8] &= ^byte(1 << (0x174 % 8))

	pa := uint64(wdk.MmGetPhysicalAddress(ptr).QuadPart)
	LogDebug("Allocated MSR Bitmap: PA=0x%X", pa)
	return pa
}

func allocateIoBitmap() uint64 {
	size := uintptr(0x2000)
	ptr := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), size, POOLTAG)
	if ptr == 0 {
		return 0
	}
	wdk.RtlZeroMemory(ptr, size)
	pa := uint64(wdk.MmGetPhysicalAddress(ptr).QuadPart)
	return pa
}

func allocateVmmStack() uint64 {
	stackSize := uintptr(0x10000)
	ptr := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), stackSize, POOLTAG)
	if ptr == 0 {
		return 0
	}

	stackTop := ptr + stackSize - uintptr(16)
	*(*uint64)(unsafe.Pointer(uintptr(stackTop))) = 0xDEAD0BADDEAD0BAD
	pa := uint64(wdk.MmGetPhysicalAddress(stackTop).QuadPart)
	LogDebug("Allocated VMM Stack: TopVA=0x%X TopPA=0x%X", stackTop, pa)
	return pa
}

func vmxTurnOn(pa uint64) error {
	cr4 := __readmsr(0x3E)
	if cr4&(1<<13) == 0 {
		__writemsr(0x3E, cr4|(1<<13))
	}

	LogInfo("Executing VMXON at PA=0x%X", pa)
	ret := AsmEnableVmxOperation(pa)
	if ret == 0 {
		return fmtError("VMXON failed")
	}
	return nil
}

func vmxTurnOff() {
	LogInfo("Executing VMXOFF")
	AsmVmxVmxOff()
}

func vmClear(pa uint64) error {
	ret := AsmVmxVmxClear(pa)
	if ret == 0 {
		return fmtError("VMCLEAR failed at PA=0x%X", pa)
	}
	return nil
}

func vmLoad(pa uint64) error {
	ret := AsmVmxVmxPtrld(pa)
	if ret == 0 {
		return fmtError("VMPTRLD failed at PA=0x%X", pa)
	}
	return nil
}

const (
	VmcsCtrlPinBasedVmExecControls   uint64 = 0x00004000
	VmcsCtrlPrimaryProcBasedVmExec   uint64 = 0x00004002
	VmcsCtrlSecondaryProcBasedVmExec uint64 = 0x0000401E
	VmcsCtrlExceptionBitmap          uint64 = 0x00004004
	VmcsCtrlPageFaultErrCodeMask     uint64 = 0x00004006
	VmcsCtrlPageFaultErrCodeMatch    uint64 = 0x00004008
	VmcsCtrlCr3TargetCount           uint64 = 0x0000400A
	VmcsCtrlIoBitmapA                uint64 = 0x0000400C
	VmcsCtrlIoBitmapB                uint64 = 0x0000400E
	VmcsCtrlMsrBitmap                uint64 = 0x00004010
	VmcsCtrlVmentryMsrLoadAddr       uint64 = 0x00004012
	VmcsCtrlVmentryIntInfo           uint64 = 0x00004014
	VmcsCtrlVmentryExcErrCode        uint64 = 0x00004016
	VmcsCtrlTprThreshold             uint64 = 0x00004016
	VmcsCtrlSecondaryVmxProcControls uint64 = 0x0000401E
	VmcsCtrlPleGap                   uint64 = 0x00004020
	VmcsCtrlPleWindow                uint64 = 0x00004022
	VmcsCtrlVmexitMsrStoreAddr       uint64 = 0x0000400C
	VmcsCtrlVmexitMsrLoadAddr        uint64 = 0x0000400E
	VmcsCtrlEptPointer               uint64 = 0x0000201A

	VmcsGuestEsSelector   uint64 = 0x00000800
	VmcsGuestCsSelector   uint64 = 0x00000802
	VmcsGuestCs           uint64 = 0x00000802
	VmcsGuestSsSelector   uint64 = 0x00000804
	VmcsGuestSs           uint64 = 0x00000804
	VmcsGuestDsSelector   uint64 = 0x00000806
	VmcsGuestFsSelector   uint64 = 0x00000808
	VmcsGuestGsSelector   uint64 = 0x0000080A
	VmcsGuestLdtrSelector uint64 = 0x0000080C
	VmcsGuestTrSelector   uint64 = 0x0000080E

	VmcsHostEsSelector   uint64 = 0x00000C00
	VmcsHostCsSelector   uint64 = 0x00000C02
	VmcsHostSsSelector   uint64 = 0x00000C04
	VmcsHostDsSelector   uint64 = 0x00000C06
	VmcsHostFsSelector   uint64 = 0x00000C08
	VmcsHostGsSelector   uint64 = 0x00000C0A
	VmcsHostLdtrSelector uint64 = 0x00000C0C
	VmcsHostTrSelector   uint64 = 0x00000C0E

	VmcsGuestEsLimit   uint64 = 0x00004800
	VmcsGuestCsLimit   uint64 = 0x00004802
	VmcsGuestSsLimit   uint64 = 0x00004804
	VmcsGuestDsLimit   uint64 = 0x00004806
	VmcsGuestFsLimit   uint64 = 0x00004808
	VmcsGuestGsLimit   uint64 = 0x0000480A
	VmcsGuestLdtrLimit uint64 = 0x0000480C
	VmcsGuestTrLimit   uint64 = 0x0000480E
	VmcsGuestGdtrLimit uint64 = 0x00004810
	VmcsGuestIdtrLimit uint64 = 0x00004812

	VmcsGuestEsBase   uint64 = 0x00006800
	VmcsGuestCsBase   uint64 = 0x00006802
	VmcsGuestSsBase   uint64 = 0x00006804
	VmcsGuestDsBase   uint64 = 0x00006806
	VmcsGuestFsBase   uint64 = 0x00006808
	VmcsGuestGsBase   uint64 = 0x0000680A
	VmcsGuestLdtrBase uint64 = 0x0000680C
	VmcsGuestTrBase   uint64 = 0x0000680E
	VmcsGuestGdtrBase uint64 = 0x00006810
	VmcsGuestIdtrBase uint64 = 0x00006812

	VmcsHostCr0         uint64 = 0x00006C00
	VmcsHostCr3         uint64 = 0x00006C02
	VmcsHostCr4         uint64 = 0x00006C04
	VmcsHostFsBase      uint64 = 0x00006C06
	VmcsHostGsBase      uint64 = 0x00006C08
	VmcsHostTrBase      uint64 = 0x00006C0A
	VmcsHostGdtrBase    uint64 = 0x00006C0C
	VmcsHostIdtrBase    uint64 = 0x00006C0E
	VmcsHostSysenterCs  uint64 = 0x00004C00
	VmcsHostSysenterEsp uint64 = 0x00006C10
	VmcsHostSysenterEip uint64 = 0x00006C12
	VmcsHostRsp         uint64 = 0x00006C14
	VmcsHostRip         uint64 = 0x00006C16

	VmcsGuestEsAccessRights   uint64 = 0x00004814
	VmcsGuestCsAccessRights   uint64 = 0x00004816
	VmcsGuestSsAccessRights   uint64 = 0x00004818
	VmcsGuestDsAccessRights   uint64 = 0x0000481A
	VmcsGuestFsAccessRights   uint64 = 0x0000481C
	VmcsGuestGsAccessRights   uint64 = 0x0000481E
	VmcsGuestLdtrAccessRights uint64 = 0x00004820
	VmcsGuestTrAccessRights   uint64 = 0x00004822

	VmcsHostEsAccessRights   uint64 = 0x00004C14
	VmcsHostCsAccessRights   uint64 = 0x00004C16
	VmcsHostSsAccessRights   uint64 = 0x00004C18
	VmcsHostDsAccessRights   uint64 = 0x00004C1A
	VmcsHostFsAccessRights   uint64 = 0x00004C1C
	VmcsHostGsAccessRights   uint64 = 0x00004C1E
	VmcsHostLdtrAccessRights uint64 = 0x00004C20
	VmcsHostTrAccessRights   uint64 = 0x00004C22

	VmcsGuestLinkPointer uint64 = 0x00002800
)

func (h *Hypervisor) setupVmcs(v *VCPU) error {
	vmWrite64(VmcsGuestCr0, __readcr0())
	vmWrite64(VmcsGuestCr3, __readcr3())
	vmWrite64(VmcsGuestCr4, __readcr4())

	vmWrite64(VmcsGuestDr7, 0x400)
	vmWrite64(VmcsGuestRflags, getRflags())

	vmWrite64(VmcsGuestCsSelector, 0x8)
	vmWrite64(VmcsGuestCsBase, 0)
	vmWrite64(VmcsGuestCsLimit, 0xFFFFFFFF)
	vmWrite64(VmcsGuestCsAccessRights, 0xA0FB)

	vmWrite64(VmcsGuestDsSelector, 0x10)
	vmWrite64(VmcsGuestDsBase, 0)
	vmWrite64(VmcsGuestDsLimit, 0xFFFFFFFF)
	vmWrite64(VmcsGuestDsAccessRights, 0xC0F3)

	vmWrite64(VmcsGuestEsSelector, 0x10)
	vmWrite64(VmcsGuestEsBase, 0)
	vmWrite64(VmcsGuestEsLimit, 0xFFFFFFFF)
	vmWrite64(VmcsGuestEsAccessRights, 0xC0F3)

	vmWrite64(VmcsGuestFsSelector, 0x10)
	vmWrite64(VmcsGuestFsBase, 0)
	vmWrite64(VmcsGuestFsLimit, 0xFFFFFFFF)
	vmWrite64(VmcsGuestFsAccessRights, 0xC0F3)

	vmWrite64(VmcsGuestGsSelector, 0x10)
	vmWrite64(VmcsGuestGsBase, 0)
	vmWrite64(VmcsGuestGsLimit, 0xFFFFFFFF)
	vmWrite64(VmcsGuestGsAccessRights, 0xC0F3)

	vmWrite64(VmcsGuestSsSelector, 0x10)
	vmWrite64(VmcsGuestSsBase, 0)
	vmWrite64(VmcsGuestSsLimit, 0xFFFFFFFF)
	vmWrite64(VmcsGuestSsAccessRights, 0xC0F3)

	gdtBase := getGdtBase()
	idtBase := getIdtBase()
	gdtLimit := uint64(getGdtLimit())
	idtLimit := uint64(getIdtLimit())

	vmWrite64(VmcsGuestGdtrBase, gdtBase)
	vmWrite64(VmcsGuestGdtrLimit, gdtLimit)
	vmWrite64(VmcsGuestIdtrBase, idtBase)
	vmWrite64(VmcsGuestIdtrLimit, idtLimit)

	vmWrite64(VmcsGuestLdtrSelector, 0)
	vmWrite64(VmcsGuestLdtrBase, 0)
	vmWrite64(VmcsGuestLdtrLimit, 0xFFFF)
	vmWrite64(VmcsGuestLdtrAccessRights, 0x82)

	vmWrite64(VmcsGuestTrSelector, 0)
	vmWrite64(VmcsGuestTrBase, 0)
	vmWrite64(VmcsGuestTrLimit, 0xFFFF)
	vmWrite64(VmcsGuestTrAccessRights, 0x8B)

	eptp := v.EptPointer.AsUInt
	vmWrite64(0x0000201A, eptp)

	pinBasedVmExecControls := adjustVmcsControl(__readmsr(0x481), 0x10000016)
	vmWrite64(VmcsCtrlPinBasedVmExecControls, pinBasedVmExecControls)

	cpuBasedVmExecControls := uint64(0x04007EFA)
	secondaryProcessorBasedControls := uint64(0x06D8FBFE)

	cpuBasedVmExecControls |= (1 << 31)
	cpuBasedVmExecControls = adjustVmcsControl(__readmsr(0x483), cpuBasedVmExecControls)
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedVmExecControls)

	secondaryProcessorBasedControls = adjustVmcsControl(__readmsr(0x485), secondaryProcessorBasedControls)
	vmWrite64(VmcsCtrlSecondaryProcBasedVmExec, secondaryProcessorBasedControls)

	vmWrite64(VmcsCtrlExceptionBitmap, 0xFFFFFFFF)
	vmWrite64(VmcsCtrlIoBitmapA, v.IoBitmapPhysicalAddressA)
	vmWrite64(VmcsCtrlIoBitmapB, v.IoBitmapPhysicalAddressB)
	vmWrite64(VmcsCtrlMsrBitmap, v.MsrBitmapPhysicalAddress)

	vmWrite64(VmcsHostCr0, __readcr0())
	vmWrite64(VmcsHostCr3, __readcr3())
	vmWrite64(VmcsHostCr4, __readcr4())

	csSelector := getCs()
	ssSelector := getSs()
	dsSelector := getDs()
	esSelector := getEs()
	fsSelector := getFs()
	gsSelector := getGs()

	vmWrite64(VmcsHostCsSelector, uint64(csSelector))
	vmWrite64(VmcsHostSsSelector, uint64(ssSelector))
	vmWrite64(VmcsHostDsSelector, uint64(dsSelector))
	vmWrite64(VmcsHostEsSelector, uint64(esSelector))
	vmWrite64(VmcsHostFsSelector, uint64(fsSelector))
	vmWrite64(VmcsHostGsSelector, uint64(gsSelector))
	vmWrite64(VmcsHostTrSelector, uint64(0x40))

	vmWrite64(VmcsHostFsBase, __readmsr(0xC0000100))
	vmWrite64(VmcsHostGsBase, __readmsr(0xC0000101))

	vmWrite64(VmcsHostLdtrSelector, 0)
	vmWrite64(VmcsHostLdtrAccessRights, 0)

	vmWrite64(VmcsHostTrBase, 0)
	vmWrite64(VmcsHostTrAccessRights, 0x1000B)

	vmWrite64(VmcsHostGdtrBase, gdtBase)
	vmWrite64(VmcsHostIdtrBase, idtBase)

	vmWrite64(VmcsHostSysenterCs, __readmsr(0x174))
	vmWrite64(VmcsHostSysenterEsp, __readmsr(0x176))
	vmWrite64(VmcsHostSysenterEip, __readmsr(0x178))

	vmWrite64(VmcsHostCsAccessRights, 0xA0FB)
	vmWrite64(VmcsHostSsAccessRights, 0xC0F3)
	vmWrite64(VmcsHostDsAccessRights, 0xC0F3)
	vmWrite64(VmcsHostEsAccessRights, 0xC0F3)
	vmWrite64(VmcsHostFsAccessRights, 0xC0F3)
	vmWrite64(VmcsHostGsAccessRights, 0xC0F3)

	vmWrite64(VmcsHostRsp, v.VmmStack)
	vmWrite64(VmcsHostRip, uint64(AsmVmexitHandlerAddr()))

	vmWrite64(VmcsGuestLinkPointer, 0xFFFFFFFFFFFFFFFF)
	vmWrite64(VmcsGuestIA32Debugctl, 0)
	vmWrite64(VmcsGuestSymEFER, __readmsr(0xC0000080))

	LogDebug("VMCS configured for core %d, EPTP=0x%X", v.CoreId, eptp)
	return nil
}

func adjustVmcsControl(msrValue, suggestedValue uint64) uint64 {
	lowMust0 := msrValue & 0xFFFFFFFF
	highMust1 := msrValue >> 32
	suggestedValue &= ^lowMust0
	suggestedValue |= highMust1
	return suggestedValue
}

func vmExitHandler() {}

func getCs() uint16       { return 0 }
func getSs() uint16       { return 0 }
func getDs() uint16       { return 0 }
func getEs() uint16       { return 0 }
func getFs() uint16       { return 0 }
func getGs() uint16       { return 0 }
func getGdtBase() uint64  { return 0 }
func getIdtBase() uint64  { return 0 }
func getGdtLimit() uint32 { return 0 }
func getIdtLimit() uint32 { return 0 }
func getRflags() uint64   { return 0 }

func getCurrentProcessorNumber() uint32 {
	var n uint32
	wdk.KeGetCurrentProcessorNumberEx(&n)
	return n
}

func breakpointReapplyHook(coreId uint32) bool { return false }
func handleRegisteredMtf(coreId uint32)        {}

func DebuggerUninitialize() {
	LogInfo("Debugger uninitialized")
}

func VmFuncUninitVmm() {
	LogInfo("VMM uninitialized")
}
func kdHandleNmiCallback(coreId uint32) bool { return false }

var (
	ErrInvalidCore = fmtError("invalid processor core")
	ErrVmcsAlloc   = fmtError("failed to allocate VMCS region")
	ErrVmxonAlloc  = fmtError("failed to allocate VMXON region")
)

const IA32_DEBUGCTL_MSR uint32 = 0x1D9

//so:extern nodecl
func __vmx_vmread(field uint64, value *uint64) uint8 { return AsmVmxVmread(field, value) }

//so:extern nodecl
func __vmx_vmwrite(field uint64, value uint64) uint8 { return AsmVmxVmwrite(field, value) }

//so:extern nodecl
func __readmsr(msr uint32) uint64 { return 0 }

//so:extern nodecl
func __writemsr(msr uint32, value uint64) {}

func inveptSingleContext(eptp uint64) {
	type InveptDesc struct {
		Eptp, Padding uint64
	}
	desc := InveptDesc{Eptp: eptp}
	AsmInvept(uint64(uintptr(unsafe.Pointer(&desc))))
}

func inveptAllContexts() {
	AsmInveptAllContexts()
}

func invvpidAllContexts() {
	AsmInvvpid()
}

func (h *Hypervisor) setRdtscExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 16)
	} else {
		cpuBasedControls &^= (1 << 16)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setPmcVmexit(enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 7)
	} else {
		cpuBasedControls &^= (1 << 7)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setExceptionBitmap(v *VCPU, bitmap uint32) {
	vmWrite64(VmcsCtrlExceptionBitmap, uint64(bitmap))
}

func (h *Hypervisor) setMovDebugRegsExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 8)
	} else {
		cpuBasedControls &^= (1 << 8)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setExternalInterruptExiting(v *VCPU, enable BOOLEAN) {
	pinBasedControls := uint64(0)
	vmRead64(VmcsCtrlPinBasedVmExecControls, &pinBasedControls)
	if enable {
		pinBasedControls |= (1 << 0)
	} else {
		pinBasedControls &^= (1 << 0)
	}
	vmWrite64(VmcsCtrlPinBasedVmExecControls, pinBasedControls)
}

func (h *Hypervisor) setCpuidExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 10)
	} else {
		cpuBasedControls &^= (1 << 10)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setClpExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 24)
	} else {
		cpuBasedControls &^= (1 << 24)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setMovFromCrExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 19)
	} else {
		cpuBasedControls &^= (1 << 19)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setMovToCrExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 20)
	} else {
		cpuBasedControls &^= (1 << 20)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setMovFromDrExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 23)
	} else {
		cpuBasedControls &^= (1 << 23)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func (h *Hypervisor) setMovToDrExiting(v *VCPU, enable BOOLEAN) {
	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
	if enable {
		cpuBasedControls |= (1 << 22)
	} else {
		cpuBasedControls &^= (1 << 22)
	}
	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func protectedHvExternalInterruptExitingForDisablingInterruptCommands(v *VCPU) {}

const (
	VMCS_GUEST_IA32_DEBUGCTL uint64 = 0x00002802
)
