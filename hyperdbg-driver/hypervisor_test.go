package main

import (
	"testing"
)

func TestVmcsFieldConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint64
		expected uint64
	}{
		{"VMCS_GUEST_RIP", VmcsGuestRip, 0x0000681E},
		{"VMCS_GUEST_RSP", VmcsGuestRsp, 0x0000681C},
		{"VMCS_GUEST_CR0", VmcsGuestCr0, 0x00006800},
		{"VMCS_GUEST_CR3", VmcsGuestCr3, 0x00006802},
		{"VMCS_GUEST_CR4", VmcsGuestCr4, 0x00006804},
		{"VMCS_GUEST_RFLAGS", VmcsGuestRflags, 0x00006820},
		{"VMCS_GUEST_DR7", VmcsGuestDr7, 0x0000681A},
		{"VMCS_GUEST_DEBUGCTL", VmcsGuestIA32Debugctl, 0x00002802},
		{"VMCS_GUEST_EFER", VmcsGuestSymEFER, 0x00006826},
		{"VMCS_GUEST_SSP", VmcsGuestSsp, 0x0000682A},
		{"VMCS_GUEST_PHYS_ADDR", VmcsGuestPhysicalAddress, 0x00002400},

		{"VMCS_GUEST_GDTR_BASE", VmcsGuestGdtrBase, 0x00006810},
		{"VMCS_GUEST_GDTR_LIMIT", VmcsGuestGdtrLimit, 0x00004810},
		{"VMCS_GUEST_IDTR_BASE", VmcsGuestIdtrBase, 0x00006812},
		{"VMCS_GUEST_IDTR_LIMIT", VmcsGuestIdtrLimit, 0x00004812},

		{"VMCS_GUEST_CS_SELECTOR", VmcsGuestCsSelector, 0x00000802},
		{"VMCS_GUEST_DS_SELECTOR", VmcsGuestDsSelector, 0x00000806},
		{"VMCS_GUEST_ES_SELECTOR", VmcsGuestEsSelector, 0x00000800},
		{"VMCS_GUEST_FS_SELECTOR", VmcsGuestFsSelector, 0x00000808},
		{"VMCS_GUEST_GS_SELECTOR", VmcsGuestGsSelector, 0x0000080A},
		{"VMCS_GUEST_SS_SELECTOR", VmcsGuestSsSelector, 0x00000804},
		{"VMCS_GUEST_TR_SELECTOR", VmcsGuestTrSelector, 0x0000080E},
		{"VMCS_GUEST_LDTR_SELECTOR", VmcsGuestLdtrSelector, 0x0000080C},

		{"VMCS_GUEST_CS_BASE", VmcsGuestCsBase, 0x00006802},
		{"VMCS_GUEST_DS_BASE", VmcsGuestDsBase, 0x00006806},
		{"VMCS_GUEST_ES_BASE", VmcsGuestEsBase, 0x00006800},
		{"VMCS_GUEST_FS_BASE", VmcsGuestFsBase, 0x00006808},
		{"VMCS_GUEST_GS_BASE", VmcsGuestGsBase, 0x0000680A},
		{"VMCS_GUEST_SS_BASE", VmcsGuestSsBase, 0x00006804},
		{"VMCS_GUEST_TR_BASE", VmcsGuestTrBase, 0x0000680E},
		{"VMCS_GUEST_LDTR_BASE", VmcsGuestLdtrBase, 0x0000680C},

		{"VMCS_GUEST_CS_LIMIT", VmcsGuestCsLimit, 0x00004802},
		{"VMCS_GUEST_DS_LIMIT", VmcsGuestDsLimit, 0x00004806},
		{"VMCS_GUEST_ES_LIMIT", VmcsGuestEsLimit, 0x00004800},
		{"VMCS_GUEST_FS_LIMIT", VmcsGuestFsLimit, 0x00004808},
		{"VMCS_GUEST_GS_LIMIT", VmcsGuestGsLimit, 0x0000480A},
		{"VMCS_GUEST_SS_LIMIT", VmcsGuestSsLimit, 0x00004804},
		{"VMCS_GUEST_TR_LIMIT", VmcsGuestTrLimit, 0x0000480E},
		{"VMCS_GUEST_LDTR_LIMIT", VmcsGuestLdtrLimit, 0x0000480C},

		{"VMCS_GUEST_CS_ACCESS_RIGHTS", VmcsGuestCsAccessRights, 0x00004816},
		{"VMCS_GUEST_DS_ACCESS_RIGHTS", VmcsGuestDsAccessRights, 0x0000481A},
		{"VMCS_GUEST_ES_ACCESS_RIGHTS", VmcsGuestEsAccessRights, 0x00004814},
		{"VMCS_GUEST_FS_ACCESS_RIGHTS", VmcsGuestFsAccessRights, 0x0000481C},
		{"VMCS_GUEST_GS_ACCESS_RIGHTS", VmcsGuestGsAccessRights, 0x0000481E},
		{"VMCS_GUEST_SS_ACCESS_RIGHTS", VmcsGuestSsAccessRights, 0x00004818},
		{"VMCS_GUEST_TR_ACCESS_RIGHTS", VmcsGuestTrAccessRights, 0x00004822},
		{"VMCS_GUEST_LDTR_ACCESS_RIGHTS", VmcsGuestLdtrAccessRights, 0x00004820},

		{"VMCS_HOST_CR0", VmcsHostCr0, 0x00006C00},
		{"VMCS_HOST_CR3", VmcsHostCr3, 0x00006C02},
		{"VMCS_HOST_CR4", VmcsHostCr4, 0x00006C04},
		{"VMCS_HOST_RIP", VmcsHostRip, 0x00006C16},
		{"VMCS_HOST_RSP", VmcsHostRsp, 0x00006C14},

		{"VMCS_HOST_CS_SELECTOR", VmcsHostCsSelector, 0x00000C02},
		{"VMCS_HOST_DS_SELECTOR", VmcsHostDsSelector, 0x00000C06},
		{"VMCS_HOST_ES_SELECTOR", VmcsHostEsSelector, 0x00000C00},
		{"VMCS_HOST_FS_SELECTOR", VmcsHostFsSelector, 0x00000C08},
		{"VMCS_HOST_GS_SELECTOR", VmcsHostGsSelector, 0x00000C0A},
		{"VMCS_HOST_SS_SELECTOR", VmcsHostSsSelector, 0x00000C04},
		{"VMCS_HOST_TR_SELECTOR", VmcsHostTrSelector, 0x00000C0E},
		{"VMCS_HOST_LDTR_SELECTOR", VmcsHostLdtrSelector, 0x00000C0C},

		{"VMCS_HOST_GDTR_BASE", VmcsHostGdtrBase, 0x00006C0C},
		{"VMCS_HOST_IDTR_BASE", VmcsHostIdtrBase, 0x00006C0E},

		{"VMCS_HOST_SYSENTER_CS", VmcsHostSysenterCs, 0x00004C00},
		{"VMCS_HOST_SYSENTER_ESP", VmcsHostSysenterEsp, 0x00006C10},
		{"VMCS_HOST_SYSENTER_EIP", VmcsHostSysenterEip, 0x00006C12},

		{"VMCS_HOST_FS_BASE", VmcsHostFsBase, 0x00006C06},
		{"VMCS_HOST_GS_BASE", VmcsHostGsBase, 0x00006C08},

		{"VMCS_CTRL_PIN_BASED_VM_EXEC_CONTROLS", VmcsCtrlPinBasedVmExecControls, 0x00004000},
		{"VMCS_CTRL_PRIMARY_PROC_BASED_VM_EXEC", VmcsCtrlPrimaryProcBasedVmExec, 0x00004002},
		{"VMCS_CTRL_SECONDARY_PROC_BASED_VM_EXEC", VmcsCtrlSecondaryProcBasedVmExec, 0x0000401E},
		{"VMCS_CTRL_EXCEPTION_BITMAP", VmcsCtrlExceptionBitmap, 0x00004004},
		{"VMCS_CTRL_IO_BITMAP_A", VmcsCtrlIoBitmapA, 0x0000400C},
		{"VMCS_CTRL_IO_BITMAP_B", VmcsCtrlIoBitmapB, 0x0000400E},
		{"VMCS_CTRL_MSR_BITMAP", VmcsCtrlMsrBitmap, 0x00004010},
		{"VMCS_CTRL_EPT_POINTER", VmcsCtrlEptPointer, 0x0000201A},
		{"VMCS_CTRL_PAGE_FAULT_ERR_CODE_MASK", VmcsCtrlPageFaultErrCodeMask, 0x00004006},
		{"VMCS_CTRL_PAGE_FAULT_ERR_CODE_MATCH", VmcsCtrlPageFaultErrCodeMatch, 0x00004008},
		{"VMCS_CTRL_CR3_TARGET_COUNT", VmcsCtrlCr3TargetCount, 0x0000400A},
		{"VMCS_CTRL_VMENTRY_MSR_LOAD_ADDR", VmcsCtrlVmentryMsrLoadAddr, 0x00004012},
		{"VMCS_CTRL_VMENTRY_INT_INFO", VmcsCtrlVmentryIntInfo, 0x00004014},
		{"VMCS_CTRL_VMENTRY_EXC_ERR_CODE", VmcsCtrlVmentryExcErrCode, 0x00004016},
		{"VMCS_CTRL_TPR_THRESHOLD", VmcsCtrlTprThreshold, 0x00004016},
		{"VMCS_CTRL_SECONDARY_VMX_PROC_CONTROLS", VmcsCtrlSecondaryVmxProcControls, 0x0000401E},
		{"VMCS_CTRL_PLE_GAP", VmcsCtrlPleGap, 0x00004020},
		{"VMCS_CTRL_PLE_WINDOW", VmcsCtrlPleWindow, 0x00004022},
		{"VMCS_CTRL_VMEXIT_MSR_STORE_ADDR", VmcsCtrlVmexitMsrStoreAddr, 0x0000400C},
		{"VMCS_CTRL_VMEXIT_MSR_LOAD_ADDR", VmcsCtrlVmexitMsrLoadAddr, 0x0000400E},

		{"VMCS_GUEST_LINK_POINTER", VmcsGuestLinkPointer, 0x00002800},
		{"VMCS_EXIT_REASON", VmcsExitReason, 0x00004400},
		{"VMCS_EXIT_QUALIFICATION", VmcsExitQualification, 0x00004402},
		{"VMCS_EXIT_INSTR_LENGTH", VmcsVmexitInstrLength, 0x0000440C},
		{"VMCS_CTRL_VMENTRY_INTR_INFO", VmcsCtrlVmentryIntrInfo, 0x00004016},
		{"VMCS_CTRL_VMENTRY_ERR_CODE", VmcsCtrlVmentryErrCode, 0x00004018},
		{"VMCS_CTRL_VMENTRY_INSTR_LEN", VmcsCtrlVmentryInstrLen, 0x0000401A},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestExitReasonConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      ExitReason
		expected ExitReason
	}{
		{"ExitTripleFault", ExitTripleFault, 2},
		{"ExitVmclear", ExitVmclear, 0x13},
		{"ExitVmlaunch", ExitVmlaunch, 0x14},
		{"ExitVmptrld", ExitVmptrld, 0x15},
		{"ExitVmptrst", ExitVmptrst, 0x16},
		{"ExitVmread", ExitVmread, 0x17},
		{"ExitVmresume", ExitVmresume, 0x18},
		{"ExitVmwrite", ExitVmwrite, 0x19},
		{"ExitVmxoff", ExitVmxoff, 0x1A},
		{"ExitVmxon", ExitVmxon, 0x1B},
		{"ExitInvd", ExitInvd, 0x0D},
		{"ExitGetsec", ExitGetsec, 0x0B},
		{"ExitInvept", ExitInvept, 0x32},
		{"ExitInvvpid", ExitInvvpid, 0x35},
		{"ExitCpuid", ExitCpuid, 22},
		{"ExitHlt", ExitHlt, 24},
		{"ExitRdtsc", ExitRdtsc, 16},
		{"ExitRdpmc", ExitRdpmc, 15},
		{"ExitVmcall", ExitVmcall, 18},
		{"ExitMovCr", ExitMovCr, 28},
		{"ExitMovDr", ExitMovDr, 29},
		{"ExitIo", ExitIo, 30},
		{"ExitRdmsr", ExitRdmsr, 31},
		{"ExitWrmsr", ExitWrmsr, 32},
		{"ExitExceptionNmi", ExitExceptionNmi, 0},
		{"ExitExtInt", ExitExtInt, 1},
		{"ExitIntWindow", ExitIntWindow, 7},
		{"ExitNmiWindow", ExitNmiWindow, 8},
		{"ExitMtf", ExitMtf, 37},
		{"ExitEptViolation", ExitEptViolation, 48},
		{"ExitEptMisconfig", ExitEptMisconfig, 49},
		{"ExitRdtscp", ExitRdtscp, 51},
		{"ExitXsetbv", ExitXsetbv, 55},
		{"ExitPreemptTimer", ExitPreemptTimer, 52},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestVmcallConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      VmcallNumber
		expected VmcallNumber
	}{
		{"VmcTest", VmcTest, 0x1},
		{"VmcVmxoff", VmcVmxoff, 0x2},
		{"VmcChangePageAttrib", VmcChangePageAttrib, 0x3},
		{"VmcInveptSingleContext", VmcInveptSingleContext, 0x4},
		{"VmcInveptAllContexts", VmcInveptAllContexts, 0x5},
		{"VmcUnhookAllPages", VmcUnhookAllPages, 0x6},
		{"VmcUnhookSinglePage", VmcUnhookSinglePage, 0x7},
		{"VmcEnableSyscallHookEfer", VmcEnableSyscallHookEfer, 0x8},
		{"VmcDisableSyscallHookEfer", VmcDisableSyscallHookEfer, 0x9},
		{"VmcChangeMsrBitmapRead", VmcChangeMsrBitmapRead, 0xA},
		{"VmcChangeMsrBitmapWrite", VmcChangeMsrBitmapWrite, 0xB},
		{"VmcSetRdtscExiting", VmcSetRdtscExiting, 0xC},
		{"VmcSetRdpmcExiting", VmcSetRdpmcExiting, 0xD},
		{"VmcSetExceptionBitmap", VmcSetExceptionBitmap, 0xE},
		{"VmcEnableMovToDebugRegsExiting", VmcEnableMovToDebugRegsExiting, 0xF},
		{"VmcEnableExternalInterruptExiting", VmcEnableExternalInterruptExiting, 0x10},
		{"VmcChangeIoBitmap", VmcChangeIoBitmap, 0x11},
		{"VmcQueryPerformanceFrequency", VmcQueryPerformanceFrequency, 0x2F},
		{"VmcQueryPerformanceCounter", VmcQueryPerformanceCounter, 0x30},
		{"VmcCheckVmxSupport", VmcCheckVmxSupport, 0x31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestEptConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint64
		expected uint64
	}{
		{"EptRead", EptRead, 1 << 0},
		{"EptWrite", EptWrite, 1 << 1},
		{"EptExec", EptExec, 1 << 2},
		{"EptPhysMask", EptPhysMask, 0xFFFFFFFFFF000},
		{"EptTypeShift", EptTypeShift, 3},
		{"EptTypeMask", EptTypeMask, 0x38},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestMsrConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint32
		expected uint32
	}{
		{"IA32_VMX_BASIC", IA32_VMX_BASIC, 0x480},
		{"IA32_FS_BASE", IA32_FS_BASE, 0xC0000100},
		{"IA32_GS_BASE", IA32_GS_BASE, 0xC0000101},
		{"IA32_SYSENTER_CS", IA32_SYSENTER_CS, 0x174},
		{"IA32_SYSENTER_ESP", IA32_SYSENTER_ESP, 0x176},
		{"IA32_SYSENTER_EIP", IA32_SYSENTER_EIP, 0x178},
		{"IA32_EFER", IA32_EFER, 0xC0000080},
		{"IA32_DEBUGCTL", IA32_DEBUGCTL_MSR, 0x1D9},
		{"IA32_VMX_ENTRY_CTLS", IA32_VMX_ENTRY_CTLS, 0x484},
		{"IA32_VMX_EXIT_CTLS", IA32_VMX_EXIT_CTLS, 0x483},
		{"IA32_VMX_true_ENTRY_CTLS", IA32_VMX_true_ENTRY_CTLS, 0x48C},
		{"IA32_VMX_true_EXIT_CTLS", IA32_VMX_true_EXIT_CTLS, 0x48D},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestInterruptConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint32
		expected uint32
	}{
		{"IntTypeHwException", IntTypeHwException, 3},
		{"IntTypeSwException", IntTypeSwException, 6},
		{"IntTypeNmi", IntTypeNmi, 2},
		{"VectorBp", VectorBp, 3},
		{"VectorGp", VectorGp, 13},
		{"VectorUd", VectorUd, 6},
		{"VectorNmi", VectorNmi, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestMsrBitmapConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      uint64
		expected uint64
	}{
		{"MSR_BITMAP_SIZE", MSR_BITMAP_SIZE, 0x2000},
		{"MSR_RANGE_COUNT", MSR_RANGE_COUNT, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%X, want 0x%X", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestEptEntryOperations(t *testing.T) {
	t.Run("EptEntry_SetAndGetRaw", func(t *testing.T) {
		entry := &EptEntry{}
		entry.SetRaw(0x123456789ABCDEF0)
		if entry.Raw() != 0x123456789ABCDEF0 {
			t.Errorf("Raw() = 0x%X, want 0x123456789ABCDEF0", entry.Raw())
		}
	})

	t.Run("EptEntry_IsPresent", func(t *testing.T) {
		entry := &EptEntry{}
		if entry.IsPresent() != false {
			t.Error("IsPresent() should be false for zero entry")
		}
		entry.SetRaw(EptRead)
		if entry.IsPresent() != true {
			t.Error("IsPresent() should be true when EptRead is set")
		}
	})

	t.Run("EptEntry_PhysAddr", func(t *testing.T) {
		entry := &EptEntry{}
		entry.SetRaw(0x123456789000)
		if entry.PhysAddr() != 0x123456789000 {
			t.Errorf("PhysAddr() = 0x%X, want 0x123456789000", entry.PhysAddr())
		}
	})

	t.Run("EptEntry_SetPhysAddr", func(t *testing.T) {
		entry := &EptEntry{}
		entry.SetRaw(0xFFFFFFFFFFFFF000)
		entry.SetPhysAddr(0x123456789000)
		if entry.PhysAddr() != 0x123456789000 {
			t.Errorf("PhysAddr() = 0x%X, want 0x123456789000", entry.PhysAddr())
		}
	})

	t.Run("EptEntry_MemoryType", func(t *testing.T) {
		entry := &EptEntry{}
		entry.SetMemoryType(MemTypeWriteBack)
		if entry.MemoryType() != MemTypeWriteBack {
			t.Errorf("MemoryType() = %d, want %d", entry.MemoryType(), MemTypeWriteBack)
		}
	})

	t.Run("EptEntry_Permission", func(t *testing.T) {
		entry := &EptEntry{}
		perm := PagePermission{Read: true, Write: true, Exec: false}
		entry.SetPermission(perm)
		p := entry.Permission()
		if p.Read != true || p.Write != true || p.Exec != false {
			t.Errorf("Permission() = {Read:%v, Write:%v, Exec:%v}, want {Read:true, Write:true, Exec:false}", p.Read, p.Write, p.Exec)
		}
	})

	t.Run("EptEntry_IsLargePage", func(t *testing.T) {
		entry := &EptEntry{}
		if entry.IsLargePage() != false {
			t.Error("IsLargePage() should be false for zero entry")
		}
		entry.SetLargePage(true)
		if entry.IsLargePage() != true {
			t.Error("IsLargePage() should be true when EptLargePage is set")
		}
	})

	t.Run("EptEntry_Clone", func(t *testing.T) {
		entry := &EptEntry{}
		entry.SetRaw(0x123456789ABCDEF0)
		clone := entry.Clone()
		if clone.Raw() != entry.Raw() {
			t.Errorf("Clone().Raw() = 0x%X, want 0x%X", clone.Raw(), entry.Raw())
		}
		clone.SetRaw(0xFFFFFFFFFFFFFFFF)
		if entry.Raw() == clone.Raw() {
			t.Error("Clone should be independent of original")
		}
	})
}

func TestPagePermission(t *testing.T) {
	t.Run("PagePermission_ToRaw_ReadOnly", func(t *testing.T) {
		p := PagePermission{Read: true, Write: false, Exec: false}
		if p.ToRaw() != EptRead {
			t.Errorf("ToRaw() = 0x%X, want 0x%X", p.ToRaw(), EptRead)
		}
	})

	t.Run("PagePermission_ToRaw_WriteOnly", func(t *testing.T) {
		p := PagePermission{Read: false, Write: true, Exec: false}
		if p.ToRaw() != EptWrite {
			t.Errorf("ToRaw() = 0x%X, want 0x%X", p.ToRaw(), EptWrite)
		}
	})

	t.Run("PagePermission_ToRaw_ExecOnly", func(t *testing.T) {
		p := PagePermission{Read: false, Write: false, Exec: true}
		if p.ToRaw() != EptExec {
			t.Errorf("ToRaw() = 0x%X, want 0x%X", p.ToRaw(), EptExec)
		}
	})

	t.Run("PagePermission_ToRaw_ReadWrite", func(t *testing.T) {
		p := PagePermission{Read: true, Write: true, Exec: false}
		if p.ToRaw() != (EptRead | EptWrite) {
			t.Errorf("ToRaw() = 0x%X, want 0x%X", p.ToRaw(), EptRead|EptWrite)
		}
	})

	t.Run("PagePermission_ToRaw_ReadWriteExec", func(t *testing.T) {
		p := PagePermission{Read: true, Write: true, Exec: true}
		if p.ToRaw() != (EptRead | EptWrite | EptExec) {
			t.Errorf("ToRaw() = 0x%X, want 0x%X", p.ToRaw(), EptRead|EptWrite|EptExec)
		}
	})

	t.Run("PagePermissionFromRaw", func(t *testing.T) {
		raw := EptRead | EptWrite | EptExec
		p := PagePermissionFromRaw(raw)
		if p.Read != true || p.Write != true || p.Exec != true {
			t.Errorf("PagePermissionFromRaw() = {Read:%v, Write:%v, Exec:%v}, want all true", p.Read, p.Write, p.Exec)
		}
	})
}

func TestEptTableWalk(t *testing.T) {
	t.Run("EptTable_pml4Index_Basic", func(t *testing.T) {
		ept := &EptTable{}
		if ept.pml4Index(0) != 0 {
			t.Error("pml4Index(0) should be 0")
		}
	})

	t.Run("EptTable_pdptIndex_Basic", func(t *testing.T) {
		ept := &EptTable{}
		if ept.pdptIndex(0) != 0 {
			t.Error("pdptIndex(0) should be 0")
		}
	})

	t.Run("EptTable_pdIndex_Basic", func(t *testing.T) {
		ept := &EptTable{}
		if ept.pdIndex(0) != 0 {
			t.Error("pdIndex(0) should be 0")
		}
	})

	t.Run("EptTable_ptIndex_Basic", func(t *testing.T) {
		ept := &EptTable{}
		if ept.ptIndex(0) != 0 {
			t.Error("ptIndex(0) should be 0")
		}
	})
}

func TestSpinlock(t *testing.T) {
	t.Run("Spinlock_New", func(t *testing.T) {
		lock := NewSpinlock()
		if lock == nil {
			t.Fatal("NewSpinlock() returned nil")
		}
	})

	t.Run("Spinlock_LockUnlock", func(t *testing.T) {
		lock := NewSpinlock()
		lock.Lock()
		lock.Unlock()
	})

	t.Run("Spinlock_Nested", func(t *testing.T) {
		outer := NewSpinlock()
		outer.Lock()
		inner := NewSpinlock()
		inner.Lock()
		inner.Unlock()
		outer.Unlock()
	})
}

func TestLogger(t *testing.T) {
	t.Run("Logger_New", func(t *testing.T) {
		logger := NewLogger(1024, 1024)
		if logger == nil {
			t.Fatal("NewLogger() returned nil")
		}
	})

	t.Run("Logger_GetBufferSize", func(t *testing.T) {
		logger := NewLogger(1024, 1024)
		size := logger.GetBufferSize()
		if size == 0 {
			t.Error("GetBufferSize() should return non-zero")
		}
	})

	t.Run("Logger_GetOverflowCount", func(t *testing.T) {
		logger := NewLogger(1024, 1024)
		count := logger.GetOverflowCount()
		if count != 0 {
			t.Errorf("GetOverflowCount() = %d, want 0", count)
		}
	})
}

func TestHookManager(t *testing.T) {
	t.Run("HookManager_New", func(t *testing.T) {
		hm := NewHookManager(10)
		if hm == nil {
			t.Fatal("NewHookManager() returned nil")
		}
	})

	t.Run("HookManager_Count", func(t *testing.T) {
		hm := NewHookManager(10)
		count := hm.Count()
		if count != 0 {
			t.Errorf("Count() = %d, want 0 for new manager", count)
		}
	})

	t.Run("HookManager_ActiveCount", func(t *testing.T) {
		hm := NewHookManager(10)
		count := hm.ActiveCount()
		if count != 0 {
			t.Errorf("ActiveCount() = %d, want 0", count)
		}
	})
}

func TestEventDispatcher(t *testing.T) {
	t.Run("EventDispatcher_New", func(t *testing.T) {
		ed := newEventDispatcher()
		if ed == nil {
			t.Fatal("newEventDispatcher() returned nil")
		}
	})

	t.Run("EventDispatcher_IsEnabled", func(t *testing.T) {
		ed := newEventDispatcher()
		enabled := ed.IsEnabled(0)
		if enabled != false {
			t.Error("IsEnabled() should return false for unregistered events")
		}
	})

	t.Run("EventDispatcher_Reset", func(t *testing.T) {
		ed := newEventDispatcher()
		ed.Reset()
	})
}

func TestMemoryManager(t *testing.T) {
	t.Run("MemoryManager_New", func(t *testing.T) {
		mm := newMemoryManager()
		if mm == nil {
			t.Fatal("newMemoryManager() returned nil")
		}
	})
}

func TestVMCSFieldEncoding(t *testing.T) {
	t.Run("VMCSFieldEncoding_BitShifts", func(t *testing.T) {
		encoding := uint32(0x0000681E)
		_ = encoding
		width := (encoding >> 13) & 0x3
		index := (encoding >> 10) & 0x3FF
		type_ := (encoding >> 8) & 0x3

		_ = width
		_ = index
		_ = type_
	})
}

func TestGuestRegs(t *testing.T) {
	t.Run("GuestRegs_Default", func(t *testing.T) {
		regs := GUEST_REGS{}
		if regs.Rax != 0 || regs.Rcx != 0 {
			t.Error("Default GUEST_REGS should be zero")
		}
	})

	t.Run("GuestRegs_SetValues", func(t *testing.T) {
		regs := &GUEST_REGS{}
		regs.Rax = 0x123456789ABCDEF0
		regs.Rcx = 0xFEDCBA9876543210
		regs.R8 = 0x1111111111111111
		regs.R15 = 0xFFFFFFFFFFFFFFFF

		if regs.Rax != 0x123456789ABCDEF0 {
			t.Errorf("Rax = 0x%X, want 0x123456789ABCDEF0", regs.Rax)
		}
		if regs.Rcx != 0xFEDCBA9876543210 {
			t.Errorf("Rcx = 0x%X, want 0xFEDCBA9876543210", regs.Rcx)
		}
		if regs.R8 != 0x1111111111111111 {
			t.Errorf("R8 = 0x%X, want 0x1111111111111111", regs.R8)
		}
		if regs.R15 != 0xFFFFFFFFFFFFFFFF {
			t.Errorf("R15 = 0x%X, want 0xFFFFFFFFFFFFFFFF", regs.R15)
		}
	})
}

func TestGuestXmmRegs(t *testing.T) {
	t.Run("GuestXmmRegs_Size", func(t *testing.T) {
		xmm := GUEST_XMM_REGS{}
		_ = xmm.Xmm15
	})
}

func TestCR3_TYPE(t *testing.T) {
	t.Run("CR3_TYPE_Default", func(t *testing.T) {
		cr3 := CR3_TYPE{}
		if cr3.Flags != 0 {
			t.Errorf("Default CR3_TYPE.Flags = 0x%X, want 0", cr3.Flags)
		}
	})

	t.Run("CR3_TYPE_SetFlags", func(t *testing.T) {
		cr3 := CR3_TYPE{}
		cr3.Flags = 0x12345678
		if cr3.Flags != 0x12345678 {
			t.Errorf("CR3_TYPE.Flags = 0x%X, want 0x12345678", cr3.Flags)
		}
	})
}

func TestEptPointer(t *testing.T) {
	t.Run("EptPointer_Default", func(t *testing.T) {
		ep := EPT_POINTER{}
		if ep.AsUInt != 0 {
			t.Errorf("Default EPT_POINTER.AsUInt = 0x%X, want 0", ep.AsUInt)
		}
	})

	t.Run("EptPointer_SetMemoryType", func(t *testing.T) {
		ep := EPT_POINTER{}
		ep.AsUInt = 0x24
		if ep.AsUInt != 0x24 {
			t.Errorf("AsUInt = 0x%X, want 0x24", ep.AsUInt)
		}
	})
}

func TestMtrrRangeDescriptor(t *testing.T) {
	t.Run("MtrrRangeDescriptor_Default", func(t *testing.T) {
		mrd := MTRR_RANGE_DESCRIPTOR{}
		if mrd.PhysicalBaseAddress != 0 {
			t.Error("Default PhysicalBaseAddress should be 0")
		}
	})

	t.Run("MtrrRangeDescriptor_SetValues", func(t *testing.T) {
		mrd := &MTRR_RANGE_DESCRIPTOR{}
		mrd.PhysicalBaseAddress = 0x100000
		mrd.PhysicalEndAddress = 0x200000
		mrd.MemoryType = 6
		mrd.FixedRange = false

		if mrd.PhysicalBaseAddress != 0x100000 {
			t.Errorf("PhysicalBaseAddress = 0x%X, want 0x100000", mrd.PhysicalBaseAddress)
		}
		if mrd.MemoryType != 6 {
			t.Errorf("MemoryType = %d, want 6 (WriteBack)", mrd.MemoryType)
		}
	})
}

func TestVmxoffState(t *testing.T) {
	t.Run("VmxoffState_Default", func(t *testing.T) {
		vs := VMXOFF_STATE{}
		if vs.Executed != false {
			t.Error("Default Executed should be false")
		}
		if vs.Rip != 0 {
			t.Errorf("Default Rip = 0x%X, want 0", vs.Rip)
		}
	})

	t.Run("VmxoffState_SetValues", func(t *testing.T) {
		vs := &VMXOFF_STATE{}
		vs.Executed = true
		vs.Rip = 0x1234567890ABCDEF
		vs.Rsp = 0xFEDCBA0987654321

		if vs.Executed != true {
			t.Error("Executed should be true")
		}
		if vs.Rip != 0x1234567890ABCDEF {
			t.Errorf("Rip = 0x%X, want 0x1234567890ABCDEF", vs.Rip)
		}
	})
}

func TestConstants(t *testing.T) {
	t.Run("NtStatusConstants", func(t *testing.T) {
		if STATUS_SUCCESS != 0x00000000 {
			t.Errorf("STATUS_SUCCESS = 0x%X, want 0x00000000", STATUS_SUCCESS)
		}
		if STATUS_UNSUCCESSFUL != -1073741823 {
			t.Errorf("STATUS_UNSUCCESSFUL = 0x%X", STATUS_UNSUCCESSFUL)
		}
	})

	t.Run("SizeConstants", func(t *testing.T) {
		if PAGE_SIZE != 4096 {
			t.Errorf("PAGE_SIZE = %d, want 4096", PAGE_SIZE)
		}
		if SIZE_1_MB != 256*4096 {
			t.Errorf("SIZE_1_MB = %d, want %d", SIZE_1_MB, 256*4096)
		}
		if SIZE_2_MB != 512*4096 {
			t.Errorf("SIZE_2_MB = %d, want %d", SIZE_2_MB, 512*4096)
		}
		if SIZE_1_GB != 512*SIZE_2_MB {
			t.Errorf("SIZE_1_GB = %d, want %d", SIZE_1_GB, 512*SIZE_2_MB)
		}
	})

	t.Run("HookConstants", func(t *testing.T) {
		if HOOK_PAGE_MONITOR_READ != 0x1 {
			t.Errorf("HOOK_PAGE_MONITOR_READ = 0x%X, want 0x1", HOOK_PAGE_MONITOR_READ)
		}
		if HOOK_PAGE_MONITOR_WRITE != 0x2 {
			t.Errorf("HOOK_PAGE_MONITOR_WRITE = 0x%X, want 0x2", HOOK_PAGE_MONITOR_WRITE)
		}
		if HOOK_PAGE_MONITOR_EXEC != 0x4 {
			t.Errorf("HOOK_PAGE_MONITOR_EXEC = 0x%X, want 0x4", HOOK_PAGE_MONITOR_EXEC)
		}
		if HOOK_PAGE_MONITOR_READ_WRITE != 0x3 {
			t.Errorf("HOOK_PAGE_MONITOR_READ_WRITE = 0x%X, want 0x3", HOOK_PAGE_MONITOR_READ_WRITE)
		}
	})

	t.Run("VmcallMagicConstants", func(t *testing.T) {
		if HYPERDBG_VMCALL_MAGIC_RAX != 0x48564653 {
			t.Errorf("HYPERDBG_VMCALL_MAGIC_RAX = 0x%X, want 0x48564653", HYPERDBG_VMCALL_MAGIC_RAX)
		}
		if HYPERDBG_VMCALL_MAGIC_RCX != 0x564D43414C4C {
			t.Errorf("HYPERDBG_VMCALL_MAGIC_RCX = 0x%X, want 0x564D43414C4C", HYPERDBG_VMCALL_MAGIC_RCX)
		}
		if HYPERDBG_VMCALL_MAGIC_RDX != 0x4E4F485950455256 {
			t.Errorf("HYPERDBG_VMCALL_MAGIC_RDX = 0x%X, want 0x4E4F485950455256", HYPERDBG_VMCALL_MAGIC_RDX)
		}
	})

	t.Run("EptCounts", func(t *testing.T) {
		if EPT_PML4_COUNT != 512 {
			t.Errorf("EPT_PML4_COUNT = %d, want 512", EPT_PML4_COUNT)
		}
		if EPT_PML3_COUNT != 512 {
			t.Errorf("EPT_PML3_COUNT = %d, want 512", EPT_PML3_COUNT)
		}
		if EPT_PML2_COUNT != 512 {
			t.Errorf("EPT_PML2_COUNT = %d, want 512", EPT_PML2_COUNT)
		}
		if EPT_PML1_COUNT != 512 {
			t.Errorf("EPT_PML1_COUNT = %d, want 512", EPT_PML1_COUNT)
		}
	})

	t.Run("PageAttribConstants", func(t *testing.T) {
		if PAGE_ATTRIB_READ != 0x2 {
			t.Errorf("PAGE_ATTRIB_READ = 0x%X, want 0x2", PAGE_ATTRIB_READ)
		}
		if PAGE_ATTRIB_WRITE != 0x4 {
			t.Errorf("PAGE_ATTRIB_WRITE = 0x%X, want 0x4", PAGE_ATTRIB_WRITE)
		}
		if PAGE_ATTRIB_EXEC != 0x8 {
			t.Errorf("PAGE_ATTRIB_EXEC = 0x%X, want 0x8", PAGE_ATTRIB_EXEC)
		}
		if PAGE_ATTRIB_EXEC_HIDDEN_HOOK != 0x10 {
			t.Errorf("PAGE_ATTRIB_EXEC_HIDDEN_HOOK = 0x%X, want 0x10", PAGE_ATTRIB_EXEC_HIDDEN_HOOK)
		}
	})
}

func TestVmxConstants(t *testing.T) {
	t.Run("VmcsEncodings", func(t *testing.T) {
		encoding := uint32(0x1000)
		_ = encoding
	})
}

func TestIoBitmap(t *testing.T) {
	t.Run("IoBitmapSize", func(t *testing.T) {
		if IO_BITMAP_SIZE != 0x2000 {
			t.Errorf("IO_BITMAP_SIZE = 0x%X, want 0x2000", IO_BITMAP_SIZE)
		}
	})
}

func TestNewFunctions(t *testing.T) {
	t.Run("NewSpinlock", func(t *testing.T) {
		lock := NewSpinlock()
		if lock == nil {
			t.Fatal("NewSpinlock() returned nil")
		}
	})

	t.Run("NewRwLock", func(t *testing.T) {
		lock := NewRwLock()
		if lock == nil {
			t.Fatal("NewRwLock() returned nil")
		}
	})

	t.Run("NewVmxRootSpinlock", func(t *testing.T) {
		lock := NewVmxRootSpinlock()
		if lock == nil {
			t.Fatal("NewVmxRootSpinlock() returned nil")
		}
	})

	t.Run("NewPoolManager", func(t *testing.T) {
		pm := NewPoolManager(1024)
		if pm == nil {
			t.Fatal("NewPoolManager() returned nil")
		}
	})

	t.Run("NewVmmOptimizer", func(t *testing.T) {
		opt := NewVmmOptimizer(OptBasic)
		if opt == nil {
			t.Fatal("NewVmmOptimizer() returned nil")
		}
	})

	t.Run("NewEptCache", func(t *testing.T) {
		cache := NewEptCache(4096)
		if cache == nil {
			t.Fatal("NewEptCache() returned nil")
		}
	})

	t.Run("NewLogBuffer", func(t *testing.T) {
		buf := NewLogBuffer(100)
		if buf == nil {
			t.Fatal("NewLogBuffer() returned nil")
		}
	})

	t.Run("NewTracerState", func(t *testing.T) {
		ts := NewTracerState()
		if ts == nil {
			t.Fatal("NewTracerState() returned nil")
		}
	})

	t.Run("NewKdSerialState", func(t *testing.T) {
		ss := NewKdSerialState()
		if ss == nil {
			t.Fatal("NewKdSerialState() returned nil")
		}
	})

	t.Run("NewEvasionState", func(t *testing.T) {
		es := NewEvasionState()
		if es == nil {
			t.Fatal("NewEvasionState() returned nil")
		}
	})

	t.Run("NewMsrHandler", func(t *testing.T) {
		mh := NewMsrHandler()
		if mh == nil {
			t.Fatal("NewMsrHandler() returned nil")
		}
	})

	t.Run("NewIoHandler", func(t *testing.T) {
		ih := NewIoHandler()
		if ih == nil {
			t.Fatal("NewIoHandler() returned nil")
		}
	})

	t.Run("NewCrAccessHandler", func(t *testing.T) {
		cah := NewCrAccessHandler()
		if cah == nil {
			t.Fatal("NewCrAccessHandler() returned nil")
		}
	})

	t.Run("NewTscHandler", func(t *testing.T) {
		th := NewTscHandler()
		if th == nil {
			t.Fatal("NewTscHandler() returned nil")
		}
	})

	t.Run("NewExceptionHandler", func(t *testing.T) {
		eh := NewExceptionHandler()
		if eh == nil {
			t.Fatal("NewExceptionHandler() returned nil")
		}
	})

	t.Run("NewEptEntry", func(t *testing.T) {
		entry := NewEptEntry()
		if entry == nil {
			t.Fatal("NewEptEntry() returned nil")
		}
	})

	t.Run("NewMtrrState", func(t *testing.T) {
		ms := NewMtrrState()
		if ms == nil {
			t.Fatal("NewMtrrState() returned nil")
		}
	})
}

func TestHypervisorNew(t *testing.T) {
	t.Run("Hypervisor_Fields", func(t *testing.T) {
		h := &Hypervisor{}
		_ = h.vcpus
		_ = h.eptTables
		_ = h.events
		_ = h.hooks
		_ = h.logger
		_ = h.initialized
		_ = h.paused
		_ = h.activeCoreId
	})
}

func TestVcpuOperations(t *testing.T) {
	t.Run("VCPU_Default", func(t *testing.T) {
		v := &VCPU{}
		if v.CoreId != 0 {
			t.Errorf("Default CoreId = %d, want 0", v.CoreId)
		}
		_ = v.IncrementRip
	})

	t.Run("VCPU_SetExitReason", func(t *testing.T) {
		v := &VCPU{}
		v.ExitReason = 1
		if v.ExitReason != 1 {
			t.Errorf("ExitReason = %d, want 1", v.ExitReason)
		}
	})
}

func TestProcessContext(t *testing.T) {
	t.Run("ProcessContext_Default", func(t *testing.T) {
		pc := &ProcessContext{}
		if pc.ProcessId != 0 {
			t.Errorf("Default ProcessId = %d, want 0", pc.ProcessId)
		}
	})

	t.Run("ProcessContext_SetValues", func(t *testing.T) {
		pc := &ProcessContext{}
		pc.ProcessId = 1234
		pc.Cr3.Flags = 0x1000
		pc.BaseAddress = 0x1234567890

		if pc.ProcessId != 1234 {
			t.Errorf("ProcessId = %d, want 1234", pc.ProcessId)
		}
		if pc.Cr3.Flags != 0x1000 {
			t.Errorf("Cr3.Flags = 0x%X, want 0x1000", pc.Cr3.Flags)
		}
		if pc.BaseAddress != 0x1234567890 {
			t.Errorf("BaseAddress = 0x%X, want 0x1234567890", pc.BaseAddress)
		}
	})
}

func TestIoctlRequestStructs(t *testing.T) {
	t.Run("VmmInitRequest", func(t *testing.T) {
		req := &VmmInitRequest{NumCores: 4}
		if req.NumCores != 4 {
			t.Errorf("NumCores = %d, want 4", req.NumCores)
		}
	})

	t.Run("EptHookRequest", func(t *testing.T) {
		req := &EptHookRequest{
			VirtAddr:  0x12345678,
			Cr3:       CR3_TYPE{Flags: 0x1000},
			HookType:  1,
			ProcessId: 1234,
		}
		if req.VirtAddr != 0x12345678 {
			t.Errorf("VirtAddr = 0x%X, want 0x12345678", req.VirtAddr)
		}
	})

	t.Run("EptUnhookRequest", func(t *testing.T) {
		req := &EptUnhookRequest{VirtAddr: 0x12345678}
		if req.VirtAddr != 0x12345678 {
			t.Errorf("VirtAddr = 0x%X, want 0x12345678", req.VirtAddr)
		}
	})

	t.Run("MemReadRequest", func(t *testing.T) {
		req := &MemReadRequest{
			Address:     0x12345678,
			Cr3:         CR3_TYPE{Flags: 0x1000},
			Size:        8,
			ProcessId:   1234,
			ReadOrWrite: false,
		}
		if req.Size != 8 {
			t.Errorf("Size = %d, want 8", req.Size)
		}
	})

	t.Run("MemReadResponse", func(t *testing.T) {
		resp := &MemReadResponse{}
		if resp == nil {
			t.Fatal("MemReadResponse should not be nil")
		}
	})

	t.Run("VirtPhysRequest", func(t *testing.T) {
		req := &VirtPhysRequest{
			Address:   0x12345678,
			Cr3:       CR3_TYPE{Flags: 0x1000},
			ProcessId: 1234,
		}
		if req.Address != 0x12345678 {
			t.Errorf("Address = 0x%X, want 0x12345678", req.Address)
		}
	})

	t.Run("RegRequest", func(t *testing.T) {
		req := &RegRequest{RegIndex: 0, CoreId: 0}
		if req.RegIndex != 0 {
			t.Errorf("RegIndex = %d, want 0", req.RegIndex)
		}
	})

	t.Run("RegResponse", func(t *testing.T) {
		resp := &RegResponse{Value: 0x12345678}
		if resp.Value != 0x12345678 {
			t.Errorf("Value = 0x%X, want 0x12345678", resp.Value)
		}
	})

	t.Run("ModifyRegsRequest", func(t *testing.T) {
		req := &ModifyRegsRequest{
			Regs:   GUEST_REGS{},
			CoreId: 0,
		}
		if req.CoreId != 0 {
			t.Errorf("CoreId = %d, want 0", req.CoreId)
		}
	})

	t.Run("SwitchProcessRequest", func(t *testing.T) {
		req := &SwitchProcessRequest{ProcessId: 1234}
		if req.ProcessId != 1234 {
			t.Errorf("ProcessId = %d, want 1234", req.ProcessId)
		}
	})

	t.Run("InveptRequest", func(t *testing.T) {
		req := &InveptRequest{Type: 1}
		if req.Type != 1 {
			t.Errorf("Type = %d, want 1", req.Type)
		}
	})

	t.Run("MtrrRequest", func(t *testing.T) {
		req := &MtrrRequest{BaseAddr: 0x1000, EndAddr: 0x2000}
		if req.BaseAddr != 0x1000 {
			t.Errorf("BaseAddr = 0x%X, want 0x1000", req.BaseAddr)
		}
	})

	t.Run("ChangeCoreRequest", func(t *testing.T) {
		req := &ChangeCoreRequest{TargetCoreId: 2}
		if req.TargetCoreId != 2 {
			t.Errorf("TargetCoreId = %d, want 2", req.TargetCoreId)
		}
	})

	t.Run("ProcessCr3Request", func(t *testing.T) {
		req := &ProcessCr3Request{ProcessId: 1234}
		if req.ProcessId != 1234 {
			t.Errorf("ProcessId = %d, want 1234", req.ProcessId)
		}
	})

	t.Run("ReadWriteMemRequest", func(t *testing.T) {
		req := &ReadWriteMemRequest{
			ReadAddress:  0x1000,
			WriteAddress: 0x2000,
			ReadSize:     8,
			WriteSize:    8,
			ReadOrWrite:  false,
			Cr3:          CR3_TYPE{Flags: 0x1000},
			ProcessId:    1234,
		}
		if req.ReadSize != 8 {
			t.Errorf("ReadSize = %d, want 8", req.ReadSize)
		}
	})

	t.Run("CallFunctionRequest", func(t *testing.T) {
		req := &CallFunctionRequest{
			FunctionAddress: 0x12345678,
			ProcessId:       1234,
		}
		if req.FunctionAddress != 0x12345678 {
			t.Errorf("FunctionAddress = 0x%X, want 0x12345678", req.FunctionAddress)
		}
	})

	t.Run("QueryPacketRequest", func(t *testing.T) {
		req := &QueryPacketRequest{}
		if req == nil {
			t.Fatal("QueryPacketRequest should not be nil")
		}
	})

	t.Run("TransparentModeRequest", func(t *testing.T) {
		req := &TransparentModeRequest{
			Enable:      true,
			Techniques:  1,
			DebuggerPid: 1234,
			KdPid:       5678,
		}
		if req.Enable != true {
			t.Errorf("Enable = %v, want true", req.Enable)
		}
	})

	t.Run("TraceStartRequest", func(t *testing.T) {
		req := &TraceStartRequest{Mode: 1}
		if req.Mode != 1 {
			t.Errorf("Mode = %d, want 1", req.Mode)
		}
	})

	t.Run("TraceBufferRequest", func(t *testing.T) {
		req := &TraceBufferRequest{BufferSize: 4096, Buffer: 0}
		if req.BufferSize != 4096 {
			t.Errorf("BufferSize = %d, want 4096", req.BufferSize)
		}
	})
}

func TestEventHandlers(t *testing.T) {
	t.Run("MsrHandler_Fields", func(t *testing.T) {
		mh := &MsrHandler{}
		_ = mh.onRead
		_ = mh.onWrite
	})

	t.Run("IoHandler_Fields", func(t *testing.T) {
		ih := &IoHandler{}
		_ = ih.onIn
		_ = ih.onOut
	})

	t.Run("CrAccessHandler_Fields", func(t *testing.T) {
		cah := &CrAccessHandler{}
		_ = cah
	})

	t.Run("TscHandler_Fields", func(t *testing.T) {
		th := &TscHandler{}
		_ = th
	})

	t.Run("ExceptionHandler_Fields", func(t *testing.T) {
		eh := &ExceptionHandler{}
		_ = eh
	})
}

func TestMemoryTypes(t *testing.T) {
	t.Run("MemoryType_Values", func(t *testing.T) {
		if MemTypeUncacheable != 0 {
			t.Errorf("MemTypeUncacheable = %d, want 0", MemTypeUncacheable)
		}
		if MemTypeWriteCombining != 1 {
			t.Errorf("MemTypeWriteCombining = %d, want 1", MemTypeWriteCombining)
		}
		if MemTypeWriteThrough != 4 {
			t.Errorf("MemTypeWriteThrough = %d, want 4", MemTypeWriteThrough)
		}
		if MemTypeWriteProtected != 5 {
			t.Errorf("MemTypeWriteProtected = %d, want 5", MemTypeWriteProtected)
		}
		if MemTypeWriteBack != 6 {
			t.Errorf("MemTypeWriteBack = %d, want 6", MemTypeWriteBack)
		}
	})
}

func TestSpinlocks(t *testing.T) {
	t.Run("Spinlock_Fields", func(t *testing.T) {
		lock := &Spinlock{}
		_ = lock
	})

	t.Run("VmxRootSpinlock_Fields", func(t *testing.T) {
		lock := &VmxRootSpinlock{}
		_ = lock
	})
}

func TestPagePermissions(t *testing.T) {
	t.Run("PagePermission_Default", func(t *testing.T) {
		p := PagePermission{}
		_ = p.Read
		_ = p.Write
		_ = p.Exec
	})

	t.Run("PagePermission_ReadWrite", func(t *testing.T) {
		p := PagePermission{Read: true, Write: true, Exec: false}
		if !p.Read || !p.Write || p.Exec {
			t.Error("ReadWrite permission not set correctly")
		}
	})

	t.Run("PagePermission_ReadExec", func(t *testing.T) {
		p := PagePermission{Read: true, Write: false, Exec: true}
		if !p.Read || p.Write || !p.Exec {
			t.Error("ReadExec permission not set correctly")
		}
	})

	t.Run("PagePermission_All", func(t *testing.T) {
		p := PagePermission{Read: true, Write: true, Exec: true}
		if !p.Read || !p.Write || !p.Exec {
			t.Error("All permissions not set correctly")
		}
	})

	t.Run("PagePermission_ToRaw", func(t *testing.T) {
		p := PagePermission{Read: true, Write: true, Exec: true}
		raw := p.ToRaw()
		if raw == 0 {
			t.Error("ToRaw should return non-zero value")
		}
	})

	t.Run("PagePermissionFromRaw", func(t *testing.T) {
		p := PagePermissionFromRaw(7)
		if !p.Read || !p.Write || !p.Exec {
			t.Error("PagePermissionFromRaw should parse all permissions")
		}
	})
}

func TestEptEntryFields(t *testing.T) {
	t.Run("EptEntry_Default", func(t *testing.T) {
		e := &EptEntry{}
		if e.Raw() != 0 {
			t.Errorf("Default Raw = 0x%X, want 0", e.Raw())
		}
	})

	t.Run("EptEntry_SetLargePage", func(t *testing.T) {
		e := &EptEntry{}
		e.SetLargePage(true)
		if !e.IsLargePage() {
			t.Error("IsLargePage() should return true after SetLargePage(true)")
		}
	})

	t.Run("EptEntry_SetMemoryType", func(t *testing.T) {
		e := &EptEntry{}
		e.SetMemoryType(MemTypeWriteBack)
		if e.MemoryType() != MemTypeWriteBack {
			t.Errorf("MemoryType() = %d, want %d", e.MemoryType(), MemTypeWriteBack)
		}
	})
}

func TestOptimizationLevel(t *testing.T) {
	t.Run("OptimizationLevel_Values", func(t *testing.T) {
		if OptNone != 0 {
			t.Errorf("OptNone = %d, want 0", OptNone)
		}
		if OptBasic != 1 {
			t.Errorf("OptBasic = %d, want 1", OptBasic)
		}
	})
}

func TestRwLock(t *testing.T) {
	t.Run("RwLock_Operations", func(t *testing.T) {
		lock := NewRwLock()
		lock.RLock()
		lock.RUnlock()
		lock.Lock()
		lock.Unlock()
	})
}

func TestLogBuffer(t *testing.T) {
	t.Run("LogBuffer_Operations", func(t *testing.T) {
		buf := NewLogBuffer(10)
		if buf == nil {
			t.Fatal("NewLogBuffer() returned nil")
		}
	})
}

func TestTracerState(t *testing.T) {
	t.Run("TracerState_Default", func(t *testing.T) {
		ts := &TracerState{}
		_ = ts
	})
}

func TestKdSerialState(t *testing.T) {
	t.Run("KdSerialState_Default", func(t *testing.T) {
		ss := &KdSerialState{}
		_ = ss
	})
}

func TestEvasionState(t *testing.T) {
	t.Run("EvasionState_Default", func(t *testing.T) {
		es := &EvasionState{}
		_ = es
	})
}

func TestMsrHandler(t *testing.T) {
	t.Run("MsrHandler_Default", func(t *testing.T) {
		mh := &MsrHandler{}
		_ = mh
	})
}

func TestIoHandler_Type(t *testing.T) {
	t.Run("IoHandler_Default", func(t *testing.T) {
		ih := &IoHandler{}
		_ = ih
	})
}

func TestCrAccessHandler(t *testing.T) {
	t.Run("CrAccessHandler_Default", func(t *testing.T) {
		cah := &CrAccessHandler{}
		_ = cah
	})
}

func TestTscHandler(t *testing.T) {
	t.Run("TscHandler_Default", func(t *testing.T) {
		th := &TscHandler{}
		_ = th
	})
}

func TestExceptionHandler(t *testing.T) {
	t.Run("ExceptionHandler_Default", func(t *testing.T) {
		eh := &ExceptionHandler{}
		_ = eh
	})
}

func TestMtrrState(t *testing.T) {
	t.Run("MtrrState_Default", func(t *testing.T) {
		ms := &MtrrState{}
		_ = ms
	})
}

func TestEptCache(t *testing.T) {
	t.Run("EptCache_Operations", func(t *testing.T) {
		cache := NewEptCache(100)
		if cache == nil {
			t.Fatal("NewEptCache() returned nil")
		}
	})
}

func TestVmmOptimizer(t *testing.T) {
	t.Run("VmmOptimizer_GetConfig", func(t *testing.T) {
		opt := NewVmmOptimizer(OptBasic)
		cfg := opt.GetConfig()
		_ = cfg
	})
}

func TestDriver(t *testing.T) {
	t.Run("Driver_Default", func(t *testing.T) {
		drv := &Driver{}
		if drv == nil {
			t.Fatal("&Driver{} returned nil")
		}
	})
}

func TestEferHookState(t *testing.T) {
	t.Run("EferHookState_Default", func(t *testing.T) {
		ehs := &EferHookState{}
		if ehs.Enabled {
			t.Error("Default Enabled should be false")
		}
	})
}

func TestExecTrapState(t *testing.T) {
	t.Run("ExecTrapState_Default", func(t *testing.T) {
		ets := &ExecTrapState{}
		_ = ets
	})
}

func TestExecTrapEntry(t *testing.T) {
	t.Run("ExecTrapEntry_Default", func(t *testing.T) {
		ete := &ExecTrapEntry{}
		_ = ete
	})
}

func TestModeBasedExecHookState(t *testing.T) {
	t.Run("ModeBasedExecHookState_Default", func(t *testing.T) {
		mbehs := &ModeBasedExecHookState{}
		_ = mbehs
	})
}

func TestModeBasedExecEntry(t *testing.T) {
	t.Run("ModeBasedExecEntry_Default", func(t *testing.T) {
		mbee := &ModeBasedExecEntry{}
		_ = mbee
	})
}

func TestTypes(t *testing.T) {
	t.Run("EptIndexFunctions", func(t *testing.T) {
		if eptPml1Offset(0x1234) != 0x234 {
			t.Error("eptPml1Offset failed")
		}
		if eptPml1Index(0x1234) != 1 {
			t.Error("eptPml1Index failed")
		}
		if eptPml2Index(0x200000) != 1 {
			t.Error("eptPml2Index failed")
		}
		if eptPml3Index(0x40000000) != 1 {
			t.Error("eptPml3Index failed")
		}
		if eptPml4Index(0x8000000000) != 1 {
			t.Error("eptPml4Index failed")
		}
	})

	t.Run("LogBufferSize", func(t *testing.T) {
		if logBufferSize() == 0 {
			t.Error("logBufferSize returned 0")
		}
		if logBufferSizePrio() == 0 {
			t.Error("logBufferSizePrio returned 0")
		}
	})
}

func TestDriverConstants(t *testing.T) {
	t.Run("SeDebugPrivilege", func(t *testing.T) {
		if SE_DEBUG_PRIVILEGE != 20 {
			t.Errorf("SE_DEBUG_PRIVILEGE = %d, want 20", SE_DEBUG_PRIVILEGE)
		}
	})
}
