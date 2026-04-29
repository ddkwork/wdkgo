package main

type VmcallNumber uint64

const (
	VmcTest                                                        VmcallNumber = 0x1
	VmcVmxoff                                                      VmcallNumber = 0x2
	VmcChangePageAttrib                                            VmcallNumber = 0x3
	VmcInveptSingleContext                                         VmcallNumber = 0x4
	VmcInveptAllContexts                                           VmcallNumber = 0x5
	VmcUnhookAllPages                                              VmcallNumber = 0x6
	VmcUnhookSinglePage                                            VmcallNumber = 0x7
	VmcEnableSyscallHookEfer                                       VmcallNumber = 0x8
	VmcDisableSyscallHookEfer                                      VmcallNumber = 0x9
	VmcChangeMsrBitmapRead                                         VmcallNumber = 0xA
	VmcChangeMsrBitmapWrite                                        VmcallNumber = 0xB
	VmcSetRdtscExiting                                             VmcallNumber = 0xC
	VmcSetRdpmcExiting                                             VmcallNumber = 0xD
	VmcSetExceptionBitmap                                          VmcallNumber = 0xE
	VmcEnableMovToDebugRegsExiting                                 VmcallNumber = 0xF
	VmcEnableExternalInterruptExiting                              VmcallNumber = 0x10
	VmcChangeIoBitmap                                              VmcallNumber = 0x11
	VmcSetHiddenCcBreakpoint                                       VmcallNumber = 0x12
	VmcDisableExternalInterruptExitingOnlyToClearInterruptCommands VmcallNumber = 0x13
	VmcUnsetRdtscExiting                                           VmcallNumber = 0x14
	VmcUnsetRdpmcExiting                                           VmcallNumber = 0x15
	VmcUnsetExceptionBitmap                                        VmcallNumber = 0x16
	VmcDisableMovToDebugRegsExiting                                VmcallNumber = 0x17
	VmcDisableExternalInterruptExiting                             VmcallNumber = 0x18
	VmcSetCpuidExiting                                             VmcallNumber = 0x19
	VmcUnsetCpuidExiting                                           VmcallNumber = 0x1A
	VmcClp                                                         VmcallNumber = 0x1B
	VmcSetMovFromCrExiting                                         VmcallNumber = 0x1C
	VmcSetMovToCrExiting                                           VmcallNumber = 0x1D
	VmcUnsetMovFromCrExiting                                       VmcallNumber = 0x1E
	VmcUnsetMovToCrExiting                                         VmcallNumber = 0x1F
	VmcSetMovFromDrExiting                                         VmcallNumber = 0x20
	VmcSetMovToDrExiting                                           VmcallNumber = 0x21
	VmcUnsetMovFromDrExiting                                       VmcallNumber = 0x22
	VmcUnsetMovToDrExiting                                         VmcallNumber = 0x23
	VmcEnableEptHookMaskedRead                                     VmcallNumber = 0x24
	VmcEnableEptHookMaskedWrite                                    VmcallNumber = 0x25
	VmcEnableEptHookMaskedReadWrite                                VmcallNumber = 0x26
	VmcEnableEptHookInlineHookRead                                 VmcallNumber = 0x27
	VmcEnableEptHookInlineHookWrite                                VmcallNumber = 0x28
	VmcEnableEptHookInlineHookReadWrite                            VmcallNumber = 0x29
	VmcEnableEptHookExecRead                                       VmcallNumber = 0x2A
	VmcEnableEptHookExecWrite                                      VmcallNumber = 0x2B
	VmcEnableEptHookExecReadWrite                                  VmcallNumber = 0x2C
	VmcEnableEptHookExecOnly                                       VmcallNumber = 0x2D
	VmcEnableMtf                                                   VmcallNumber = 0x2E
	VmcQueryPerformanceFrequency                                   VmcallNumber = 0x2F
	VmcQueryPerformanceCounter                                     VmcallNumber = 0x30
	VmcCheckVmxSupport                                             VmcallNumber = 0x31
	VmcLaunchDebugger                                              VmcallNumber = 0x100000000 + iota
)

func (h *Hypervisor) handleHyperdbgVmcall(num, p1, p2, p3 uint64) NTSTATUS {
	switch num {
	case uint64(VmcTest):
		return h.vmcallTest(p1, p2, p3)

	case uint64(VmcVmxoff):
		v := h.CurrentVcpu()
		if v != nil {
			vmxVmxoff(v)
		}
		return STATUS_SUCCESS

	case uint64(VmcChangePageAttrib):
		v := h.CurrentVcpu()
		if v == nil {
			return STATUS_UNSUCCESSFUL
		}
		cr3 := CR3_TYPE{Flags: p3}
		result := h.performPageHook(v, uintptr(p1), cr3, uint32(p2))
		if result {
			return STATUS_SUCCESS
		}
		return STATUS_UNSUCCESSFUL

	case uint64(VmcInveptSingleContext):
		inveptSingleContext(p1)
		return STATUS_SUCCESS

	case uint64(VmcInveptAllContexts):
		inveptAllContexts()
		return STATUS_SUCCESS

	case uint64(VmcUnhookAllPages):
		v := h.CurrentVcpu()
		if v != nil {
			gHooks.RestoreAll(v)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnhookSinglePage):
		v := h.CurrentVcpu()
		if v == nil {
			return STATUS_UNSUCCESSFUL
		}
		if err := gHooks.RemoveByVirtAddr(p1); err == nil {
			return STATUS_SUCCESS
		}
		return STATUS_UNSUCCESSFUL

	case uint64(VmcEnableSyscallHookEfer):
		v := h.CurrentVcpu()
		if v != nil {
			syscallHookConfigureEFER(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcDisableSyscallHookEfer):
		v := h.CurrentVcpu()
		if v != nil {
			syscallHookConfigureEFER(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcChangeMsrBitmapRead):
		v := h.CurrentVcpu()
		if v != nil {
			msrHandlePerformMsrBitmapReadChange(v, uint32(p1))
		}
		return STATUS_SUCCESS

	case uint64(VmcChangeMsrBitmapWrite):
		v := h.CurrentVcpu()
		if v != nil {
			msrHandlePerformMsrBitmapWriteChange(v, uint32(p1))
		}
		return STATUS_SUCCESS

	case uint64(VmcSetRdtscExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setRdtscExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcSetRdpmcExiting):
		h.setPmcVmexit(true)
		return STATUS_SUCCESS

	case uint64(VmcSetExceptionBitmap):
		v := h.CurrentVcpu()
		if v != nil {
			h.setExceptionBitmap(v, uint32(p1))
		}
		return STATUS_SUCCESS

	case uint64(VmcEnableMovToDebugRegsExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovDebugRegsExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcEnableExternalInterruptExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setExternalInterruptExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcChangeIoBitmap):
		v := h.CurrentVcpu()
		if v != nil {
			ioHandlePerformIoBitmapChange(v, uint32(p1))
		}
		return STATUS_SUCCESS

	case uint64(VmcSetHiddenCcBreakpoint):
		v := h.CurrentVcpu()
		if v == nil {
			return STATUS_UNSUCCESSFUL
		}
		cr3 := CR3_TYPE{Flags: p2}
		if err := gHooks.Install(v.CoreId, p1, cr3, HookExec); err == nil {
			return STATUS_SUCCESS
		}
		return STATUS_UNSUCCESSFUL

	case uint64(VmcDisableExternalInterruptExitingOnlyToClearInterruptCommands):
		v := h.CurrentVcpu()
		if v != nil {
			protectedHvExternalInterruptExitingForDisablingInterruptCommands(v)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetRdtscExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setRdtscExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetRdpmcExiting):
		h.setPmcVmexit(false)
		return STATUS_SUCCESS

	case uint64(VmcUnsetExceptionBitmap):
		v := h.CurrentVcpu()
		if v != nil {
			h.setExceptionBitmap(v, 0)
		}
		return STATUS_SUCCESS

	case uint64(VmcDisableMovToDebugRegsExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovDebugRegsExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcDisableExternalInterruptExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setExternalInterruptExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcSetCpuidExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setCpuidExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetCpuidExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setCpuidExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcClp):
		v := h.CurrentVcpu()
		if v != nil {
			h.setClpExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcSetMovFromCrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovFromCrExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcSetMovToCrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovToCrExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetMovFromCrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovFromCrExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetMovToCrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovToCrExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcSetMovFromDrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovFromDrExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcSetMovToDrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovToDrExiting(v, true)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetMovFromDrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovFromDrExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcUnsetMovToDrExiting):
		v := h.CurrentVcpu()
		if v != nil {
			h.setMovToDrExiting(v, false)
		}
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookMaskedRead):
		eptHookSetReadHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookMaskedWrite):
		eptHookSetWriteHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookMaskedReadWrite):
		eptHookSetReadWriteHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookInlineHookRead):
		eptHookSetReadHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookInlineHookWrite):
		eptHookSetWriteHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookInlineHookReadWrite):
		eptHookSetReadWriteHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookExecRead):
		eptHookSetReadHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookExecWrite):
		eptHookSetWriteHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookExecReadWrite):
		eptHookSetReadWriteHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableEptHookExecOnly):
		eptHookSetExecHook(p1, CR3_TYPE{Flags: p2})
		return STATUS_SUCCESS

	case uint64(VmcEnableMtf):
		setMonitorTrapFlag(true)
		return STATUS_SUCCESS

	case uint64(VmcQueryPerformanceFrequency):
		return NTSTATUS(queryPerformanceFrequency())

	case uint64(VmcQueryPerformanceCounter):
		return NTSTATUS(queryPerformanceCounter())

	case uint64(VmcCheckVmxSupport):
		if checkVmxSupport() {
			return STATUS_SUCCESS
		}
		return STATUS_UNSUCCESSFUL

	default:
		LogWarning("Top-level driver VMCALL %d not handled", num)
	}

	LogWarning("Unknown VMCALL number: 0x%X", num)
	return STATUS_UNSUCCESSFUL
}

func (h *Hypervisor) vmcallTest(_, _, _ uint64) NTSTATUS {
	LogInfo("VMCALL_TEST received")
	return STATUS_SUCCESS
}

func (h *Hypervisor) performPageHook(v *VCPU, addr uintptr, cr3 CR3_TYPE, mask uint32) bool {
	if mask&HOOK_PAGE_MONITOR_READ != 0 {
		if err := gHooks.Install(v.CoreId, uint64(addr), cr3, HookRead); err != nil {
			return false
		}
	}
	if mask&HOOK_PAGE_MONITOR_WRITE != 0 {
		if err := gHooks.Install(v.CoreId, uint64(addr), cr3, HookWrite); err != nil {
			return false
		}
	}
	if mask&HOOK_PAGE_MONITOR_EXEC != 0 {
		if err := gHooks.Install(v.CoreId, uint64(addr), cr3, HookExec); err != nil {
			return false
		}
	}
	if mask&(HOOK_PAGE_MONITOR_INLINE_HOOKS|HOOK_PAGE_MASKED_HOOKS) != 0 {
		if err := gHooks.Install(v.CoreId, uint64(addr), cr3, HookReadWrite); err != nil {
			return false
		}
	}
	return true
}

func queryPerformanceFrequency() uint64 { return 10000000 }
func queryPerformanceCounter() uint64   { return 0 }
func checkVmxSupport() bool            { return true }
