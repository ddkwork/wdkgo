package main

const (
	MSR_BITMAP_SIZE = 0x2000
	MSR_RANGE_COUNT = 8
)

type MsrBitmapManager struct {
	bitmap [MSR_BITMAP_SIZE]byte
}

var gMsrBitmap = &MsrBitmapManager{}

func initMsrBitmap() {
	for i := range MSR_BITMAP_SIZE {
		gMsrBitmap.bitmap[i] = 0xFF
	}
}

func msrHandlePerformMsrBitmapReadChange(v *VCPU, msrIndex uint32) {
	if v.MsrBitmapVirtualAddress == 0 {
		LogError("MSR bitmap not initialized")
		return
	}

	byteOffset := msrIndex / 8
	bitOffset := msrIndex % 8

	if msrIndex >= 0xC0000000 {
		byteOffset += MSR_BITMAP_SIZE / 2
	}

	if byteOffset < MSR_BITMAP_SIZE {
		mask := uint8(1 << bitOffset)
		gMsrBitmap.bitmap[byteOffset] |= mask

		pa := VirtToPhys(uintptr(v.MsrBitmapVirtualAddress), CR3_TYPE{})
		WritePhysMem(pa+uint64(byteOffset), gMsrBitmap.bitmap[byteOffset:byteOffset+1])
	}
}

func msrHandlePerformMsrBitmapWriteChange(v *VCPU, msrIndex uint32) {
	if v.MsrBitmapVirtualAddress == 0 {
		LogError("MSR bitmap not initialized")
		return
	}

	byteOffset := msrIndex/8 + MSR_BITMAP_SIZE/4
	bitOffset := msrIndex % 8

	if msrIndex >= 0xC0000000 {
		byteOffset += MSR_BITMAP_SIZE / 2
	}

	if byteOffset < MSR_BITMAP_SIZE {
		mask := uint8(1 << bitOffset)
		gMsrBitmap.bitmap[byteOffset] |= mask

		pa := VirtToPhys(uintptr(v.MsrBitmapVirtualAddress), CR3_TYPE{})
		WritePhysMem(pa+uint64(byteOffset), gMsrBitmap.bitmap[byteOffset:byteOffset+1])
	}
}

func msrHandleEnableOrDisableMsrExiting(v *VCPU, enable BOOLEAN) {
	if enable {
		cpuBasedControls := uint64(0)
		vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
		cpuBasedControls |= (1 << 11)
		vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)

		initMsrBitmap()
	} else {
		cpuBasedControls := uint64(0)
		vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)
		cpuBasedControls &= ^(uint64(1 << 11))
		vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
	}
}

func handleMsrRead(v *VCPU) bool {
	msrIndex := v.Regs.Rcx
	var value uint64 = 0

	switch uint32(msrIndex) {
	case IA32_EFER:
		value = __readmsr(IA32_EFER)
	default:
		value = __readmsr(uint32(msrIndex))
	}

	v.Regs.Rax = value & 0xFFFFFFFF
	v.Regs.Rdx = (value >> 32) & 0xFFFFFFFF

	h := getHypervisor()
	if h != nil {
		h.advanceIp(v)
	}
	return true
}

func handleMsrWrite(v *VCPU) bool {
	msrIndex := v.Regs.Rcx
	value := (v.Regs.Rdx << 32) | v.Regs.Rax

	switch uint32(msrIndex) {
	case IA32_EFER:
		__writemsr(IA32_EFER, value)
	default:
		__writemsr(uint32(msrIndex), value)
	}

	h := getHypervisor()
	if h != nil {
		h.advanceIp(v)
	}
	return true
}
