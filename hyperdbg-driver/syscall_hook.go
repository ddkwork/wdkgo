package main

import "unsafe"

type EferHookState struct {
	Enabled    bool
	UdHandling bool
	CetSupport bool
}

var gEferHook = &EferHookState{
	UdHandling: true,
}

//so:extern
func __writecr3(val uint64)

func getHypervisor() *Hypervisor { return gHyp }

func syscallHookConfigureEFER(v *VCPU, enable BOOLEAN) {
	vmxBasic := __readmsr(IA32_VMX_BASIC)
	var vmEntryControls, vmExitControls uint32

	vmRead32(VmcsCtrlVmentryIntrInfo, &vmEntryControls)
	vmRead32(VmcsCtrlPrimaryProcBasedVmExec, &vmExitControls)

	eferMsr := __readmsr(IA32_EFER)

	if enable {
		eferMsr &= ^uint64(1 << 0)

		vmEntryControls |= uint32(IA32_VMX_ENTRY_CTLS_LOAD_IA32_EFER_FLAG)
		vmExitControls |= uint32(IA32_VMX_EXIT_CTLS_SAVE_IA32_EFER_FLAG)

		vmWrite64(VMCS_GUEST_IA32_DEBUGCTL, eferMsr)

		h := getHypervisor()
		if h != nil {
			h.setExceptionBitmap(v, EXCEPTION_VECTOR_UNDEFINED_OPCODE)
		}

		gEferHook.Enabled = true
	} else {
		eferMsr |= (1 << 0)

		vmEntryControls &^= uint32(IA32_VMX_ENTRY_CTLS_LOAD_IA32_EFER_FLAG)
		vmExitControls &^= uint32(IA32_VMX_EXIT_CTLS_SAVE_IA32_EFER_FLAG)

		vmWrite64(VMCS_GUEST_IA32_DEBUGCTL, eferMsr)

		__writemsr(IA32_EFER, eferMsr)

		removeUndefinedInstructionForDisablingSyscallSysret(v)

		gEferHook.Enabled = false
	}

	vmWrite32(VmcsCtrlVmentryIntrInfo, uint32(adjustVmcsControl(vmxBasic, uint64(vmEntryControls))))
	vmWrite32(VmcsCtrlPrimaryProcBasedVmExec, uint32(adjustVmcsControl(vmxBasic, uint64(vmExitControls))))

	status := "DISABLED"
	if enable {
		status = "ENABLED"
	}
	LogDebug("EFER Syscall Hook: %s", status)
}

func syscallHookEmulateSYSCALL(v *VCPU) bool {
	var guestRip, guestRflags uint64
	var instrLen uint32

	vmRead64(VmcsGuestRip, &guestRip)
	vmRead32(VmcsVmexitInstrLength, &instrLen)
	vmRead64(VmcsGuestRflags, &guestRflags)

	lstar := __readmsr(IA32_LSTAR)
	v.Regs.Rcx = guestRip + uint64(instrLen)
	guestRip = lstar
	vmWrite64(VmcsGuestRip, guestRip)

	fmask := __readmsr(IA32_FMASK)
	v.Regs.R11 = guestRflags
	guestRflags &= ^(fmask | X86_FLAGS_RF)
	vmWrite64(VmcsGuestRflags, guestRflags)

	if gEferHook.CetSupport {
		ucet := __readmsr(IA32_U_CET)
		if ucet&(1<<0) != 0 {
			var ssp uint64
			vmRead64(VmcsGuestSsp, &ssp)
			__writemsr(IA32_PL3_SSP, ssp)
			vmWrite64(VmcsGuestSsp, 0)
		}
	}

	star := __readmsr(IA32_STAR)
	csSelector := uint16((star >> 32) & 0xFFFC)
	ssSelector := csSelector + 8

	setSegmentRegister(VmcsGuestCs, csSelector, 0, 0xFFFFFFFF, 0xA09B)
	setSegmentRegister(VmcsGuestSs, ssSelector, 0, 0xFFFFFFFF, 0xC093)

	dispatchEventEferSyscall(v)

	return true
}

func syscallHookEmulateSYSRET(v *VCPU) bool {
	var guestRip, guestRflags uint64

	guestRip = v.Regs.Rcx
	vmWrite64(VmcsGuestRip, guestRip)

	guestRflags = (v.Regs.R11 & ^(X86_FLAGS_RF | X86_FLAGS_VM | X86_FLAGS_RESERVED_BITS)) | X86_FLAGS_FIXED
	vmWrite64(VmcsGuestRflags, guestRflags)

	if gEferHook.CetSupport {
		ucet := __readmsr(IA32_U_CET)
		if ucet&(1<<0) != 0 {
			ssp := __readmsr(IA32_PL3_SSP)
			vmWrite64(VmcsGuestSsp, ssp)
		}
	}

	star := __readmsr(IA32_STAR)
	csSelector := uint16(((star >> 48) + 16) | 3)
	ssSelector := uint16(((star >> 48) + 8) | 3)

	setSegmentRegister(VmcsGuestCs, csSelector, 0, 0xFFFFFFFF, 0xA0FB)
	setSegmentRegister(VmcsGuestSs, ssSelector, 0, 0xFFFFFFFF, 0xC0F3)

	dispatchEventEferSysret(v)

	return true
}

func syscallHookHandleUD(v *VCPU) bool {
	var rip uint64
	vmRead64(VmcsGuestRip, &rip)

	if gEferHook.UdHandling {
		if rip&0xff00000000000000 != 0 {
			return syscallHookEmulateSYSRET(v)
		}
		return syscallHookEmulateSYSCALL(v)
	} else {
		guestCr3 := getCurrentProcessCr3()
		originalCr3 := __readcr3()

		__writecr3(guestCr3.Flags)

		var instrBuf [3]byte
		present := checkPagePresentByCr3(uintptr(rip), guestCr3)

		if present {
			readMemorySafe(rip, unsafe.Pointer(&instrBuf[0]), 3)
		} else {
			h := getHypervisor()
			if h != nil {
				h.suppressAdvance(v)
				injectPageFaultWithoutErrorCode(rip)
			}
			__writecr3(originalCr3)
			return false
		}

		__writecr3(originalCr3)

		if instrBuf[0] == 0x0F && instrBuf[1] == 0x05 {
			return syscallHookEmulateSYSCALL(v)
		}
		if instrBuf[0] == 0x48 && instrBuf[1] == 0x0F && instrBuf[2] == 0x07 {
			return syscallHookEmulateSYSRET(v)
		}

		return false
	}
}

const (
	IA32_VMX_BASIC                          uint32 = 0x480
	IA32_EFER                               uint32 = 0xC0000080
	IA32_LSTAR                              uint32 = 0xC0000082
	IA32_FMASK                              uint32 = 0xC0000084
	IA32_STAR                               uint32 = 0xC0000081
	IA32_U_CET                              uint32 = 0xC0000102
	IA32_PL3_SSP                            uint32 = 0xC0001052
	IA32_FS_BASE                            uint32 = 0xC0000100
	IA32_GS_BASE                            uint32 = 0xC0000101
	IA32_SYSENTER_CS                        uint32 = 0x174
	IA32_SYSENTER_ESP                       uint32 = 0x176
	IA32_SYSENTER_EIP                       uint32 = 0x178
	IA32_VMX_ENTRY_CTLS                     uint32 = 0x484
	IA32_VMX_EXIT_CTLS                      uint32 = 0x483
	IA32_VMX_true_ENTRY_CTLS                uint32 = 0x48C
	IA32_VMX_true_EXIT_CTLS                 uint32 = 0x48D
	IA32_VMX_ENTRY_CTLS_LOAD_IA32_EFER_FLAG uint32 = 0x10000
	IA32_VMX_EXIT_CTLS_SAVE_IA32_EFER_FLAG  uint32 = 0x20000
	X86_FLAGS_RF                            uint64 = 0x10000
	X86_FLAGS_VM                            uint64 = 0x20000
	X86_FLAGS_RESERVED_BITS                 uint64 = 0x2
	X86_FLAGS_FIXED                         uint64 = 0x2
	EXCEPTION_VECTOR_UNDEFINED_OPCODE       uint32 = 6
)

func getCurrentProcessCr3() CR3_TYPE                                 { return CR3_TYPE{} }
func checkPagePresentByCr3(addr uintptr, cr3 CR3_TYPE) bool          { return false }
func readMemorySafe(addr uint64, buffer unsafe.Pointer, size SIZE_T) {}
func dispatchEventEferSyscall(v *VCPU)                               {}
func dispatchEventEferSysret(v *VCPU)                                {}
func injectPageFaultWithoutErrorCode(rip uint64)                     {}
func removeUndefinedInstructionForDisablingSyscallSysret(v *VCPU)    {}

func setSegmentRegister(field uint64, selector uint16, base uint64, limit uint32, accessRights uint32) {
	vmWrite64(field+0, uint64(selector))
	vmWrite64(field+8, base)
	vmWrite64(field+16, uint64(limit))
	vmWrite64(field+24, uint64(accessRights))
}
