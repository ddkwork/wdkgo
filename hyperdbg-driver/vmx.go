package main

func AsmVmxVmcall(RegPtr uint64) uint64 { return 0 }

//so:extern
func AsmVmxVmread(Field uint64, FieldValue *uint64) uint8

//so:extern
func AsmVmxVmwrite(Field uint64, FieldValue uint64) uint8

//so:extern
func AsmVmxVmread32(Field uint64, FieldValue *uint32) uint8

//so:extern
func AsmVmxVmwrite32(Field uint64, FieldValue uint32) uint8

//so:extern
func AsmVmxVmlaunch() uint8

//so:extern
func AsmVmxVmresume() uint8

//so:extern
func AsmVmxVmxClear(PhysicalAddr uint64) uint8

//so:extern
func AsmVmxVmxPtrld(PhysicalAddr uint64) uint8

func AsmEnableVmxOperation(PhysicalAddr uint64) uint8 { return 0 }

//so:extern
func AsmVmxVmxOff()

func AsmInvept(Eptp uint64) {}

func AsmInveptAllContexts() {}

func AsmInvvpid() {}

func AsmInvvpidAllContexts() {}

func AsmHypervVmcall(regs uint64) {}

func AsmGetGdtBase() uint64         { return 0 }
func AsmGetIdtBase() uint64         { return 0 }
func AsmGetGdtLimit() uint32        { return 0 }
func AsmGetIdtLimit() uint32        { return 0 }
func AsmGetRflags() uint64          { return 0 }
func AsmGetEs() uint16              { return 0 }
func AsmGetCs() uint16              { return 0 }
func AsmGetSs() uint16              { return 0 }
func AsmGetDs() uint16              { return 0 }
func AsmGetFs() uint16              { return 0 }
func AsmGetGs() uint16              { return 0 }
func AsmGetTr() uint16              { return 0 }
func AsmGetLdtr() uint16            { return 0 }
func AsmGetFsBase() uint64          { return 0 }
func AsmGetGsBase() uint64          { return 0 }
func AsmVmexitHandlerAddr() uintptr { return 0 }

func vmRead64(field uint64, val *uint64) { AsmVmxVmread(field, val) }
func vmRead32(field uint64, val *uint32) { AsmVmxVmread32(field, val) }
func vmWrite64(field uint64, val uint64) { AsmVmxVmwrite(field, val) }
func vmWrite32(field uint64, val uint32) { AsmVmxVmwrite32(field, val) }

//so:extern
func __readcr4() uint64

//so:extern
func __readcr3() uint64

//so:extern
func __readcr0() uint64

//so:extern
func __writecr4(cr4 uint64)

func setMonitorTrapFlag(enable bool) {
	if enable {
		vmWrite64(0x00002802, __readmsr(0x1D9)|1<<2)
	} else {
		vmWrite64(0x00002802, __readmsr(0x1D9)&^(uint64(1<<2)))
	}
}

func vmxVmlaunch(v *VCPU) bool {
	v.HasLaunched = true
	result := AsmVmxVmlaunch()
	if result != 0 {
		return true
	}
	return false
}

func vmxVmresume(v *VCPU) bool {
	v.HasLaunched = true
	result := AsmVmxVmresume()
	if result != 0 {
		return true
	}
	return false
}

func vmxVmxoff(v *VCPU) {
	AsmVmxVmxOff()
	v.VmxoffState.Executed = true
}

func asmHypervVmcall(regs uint64) { AsmHypervVmcall(regs) }

const (
	VMCS_CTRL_PIN_BASED_VM_EXEC_CONTROLS   uint64 = 0x00004000
	VMCS_CTRL_PRIMARY_PROC_BASED_VM_EXEC   uint64 = 0x00004002
	VMCS_CTRL_SECONDARY_PROC_BASED_VM_EXEC uint64 = 0x0000401E
	VMCS_CTRL_EXCEPTION_BITMAP             uint64 = 0x00004004
	VMCS_CTRL_PAGE_FAULT_ERR_CODE_MASK     uint64 = 0x00004006
	VMCS_CTRL_PAGE_FAULT_ERR_CODE_MATCH    uint64 = 0x00004008
	VMCS_CTRL_CR3_TARGET_COUNT             uint64 = 0x0000400A
	VMCS_CTRL_IO_BITMAP_A                  uint64 = 0x0000400C
	VMCS_CTRL_IO_BITMAP_B                  uint64 = 0x0000400E
	VMCS_CTRL_MSR_BITMAP                   uint64 = 0x00004010
	VMCS_CTRL_VMENTRY_MSR_LOAD_ADDR        uint64 = 0x00004012
	VMCS_CTRL_VMENTRY_INT_INFO             uint64 = 0x00004014
	VMCS_CTRL_VMENTRY_EXC_ERR_CODE         uint64 = 0x00004016
	VMCS_CTRL_VMENTRY_INSTR_LEN            uint64 = 0x0000401A
	VMCS_CTRL_TPR_THRESHOLD                uint64 = 0x00004016
	VMCS_CTRL_SECONDARY_VMX_PROC_CONTROLS  uint64 = 0x0000401E
	VMCS_CTRL_PLE_GAP                      uint64 = 0x00004020
	VMCS_CTRL_PLE_WINDOW                   uint64 = 0x00004022
	VMCS_CTRL_VMEXIT_MSR_STORE_ADDR        uint64 = 0x0000400C
	VMCS_CTRL_VMEXIT_MSR_LOAD_ADDR         uint64 = 0x0000400E
)

func VmxVmexitHandler(regs *GUEST_REGS) bool {
	if gHyp == nil {
		return false
	}
	return gHyp.HandleVmExit(regs)
}

func VmxVmresume() uint8 {
	return 1
}

func VmxReturnStackPointerForVmxoff() uint64 { return 0 }

func VmxReturnInstructionPointerForVmxoff() uint64 { return 0 }

func VmxVirtualizeCurrentSystem(state *uint8) {}

func EptHook2GeneralDetourEventHandler(regs *GUEST_REGS, calledFrom uint64) uint64 { return 0 }

func IdtEmulationhandleHostInterrupt(trapFrame *uint8) {}
