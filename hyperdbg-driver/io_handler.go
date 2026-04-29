package main

const (
	IO_BITMAP_SIZE = 0x2000
)

type IoBitmapManager struct {
	bitmapA [IO_BITMAP_SIZE]byte
	bitmapB [IO_BITMAP_SIZE]byte
}

var gIoBitmap = &IoBitmapManager{}

func initIoBitmap() {
	for i := range IO_BITMAP_SIZE {
		gIoBitmap.bitmapA[i] = 0
		gIoBitmap.bitmapB[i] = 0
	}
}

func ioHandlePerformIoBitmapChange(v *VCPU, portMask uint32) {
	if v.IoBitmapVirtualAddressA == 0 || v.IoBitmapVirtualAddressB == 0 {
		LogError("IO bitmap not initialized")
		return
	}

	for port := range uint32(0x10000) {
		byteOffset := port / 8
		bitOffset := port % 8

		if portMask&(1<<uint(port%32)) != 0 {
			if byteOffset < IO_BITMAP_SIZE {
				if port < 0x8000 {
					gIoBitmap.bitmapA[byteOffset] |= (1 << bitOffset)
				} else {
					idx := byteOffset - 0x1000
					if idx < IO_BITMAP_SIZE {
						gIoBitmap.bitmapB[idx] |= (1 << bitOffset)
					}
				}
			}
		} else {
			if byteOffset < IO_BITMAP_SIZE {
				if port < 0x8000 {
					gIoBitmap.bitmapA[byteOffset] &^= (1 << bitOffset)
				} else {
					idx := byteOffset - 0x1000
					if idx < IO_BITMAP_SIZE {
						gIoBitmap.bitmapB[idx] &^= (1 << bitOffset)
					}
				}
			}
		}
	}

	paA := VirtToPhys(uintptr(v.IoBitmapVirtualAddressA), CR3_TYPE{})
	paB := VirtToPhys(uintptr(v.IoBitmapVirtualAddressB), CR3_TYPE{})

	WritePhysMem(paA, gIoBitmap.bitmapA[:])
	WritePhysMem(paB, gIoBitmap.bitmapB[:])

	cpuBasedControls := uint64(0)
	vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls)

	if portMask != 0xFFFFFFFF && portMask != 0 {
		cpuBasedControls |= (1 << 24)
	} else if portMask == 0xFFFFFFFF {
		cpuBasedControls |= (1 << 24)
	} else {
		cpuBasedControls &= ^(uint64(1 << 24))
	}

	vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls)
}

func ioHandleEnableOrDisableIoPortExiting(v *VCPU, enable BOOLEAN) {
	if enable {
		initIoBitmap()
		ioHandlePerformIoBitmapChange(v, 0xFFFF)
	} else {
		ioHandlePerformIoBitmapChange(v, 0x0000)
	}
}

func handleIoRead(v *VCPU) bool {
	port := uint16(v.Regs.Rdx & 0xFFFF)
	size := (v.Regs.Rax >> 16) & 0x7

	var value uint64 = 0

	switch size {
	case 0:
		value = uint64(inb(port))
	case 1:
		value = uint64(inw(port))
	case 2:
		value = uint64(ind(port))
	default:
		LogError("Invalid IO read size: %d", size)
		return false
	}

	v.Regs.Rax = value

	h := getHypervisor()
	if h != nil {
		h.advanceIp(v)
	}
	return true
}

func handleIoWrite(v *VCPU) bool {
	port := uint16(v.Regs.Rdx & 0xFFFF)
	size := (v.Regs.Rax >> 16) & 0x7
	value := v.Regs.Rax & 0xFFFF

	switch size {
	case 0:
		outb(port, uint8(value))
	case 1:
		outw(port, uint16(value))
	case 2:
		outd(port, uint32(value))
	default:
		LogError("Invalid IO write size: %d", size)
		return false
	}

	h := getHypervisor()
	if h != nil {
		h.advanceIp(v)
	}
	return true
}

func inb(port uint16) uint8        { return 0 }
func inw(port uint16) uint16       { return 0 }
func ind(port uint16) uint32       { return 0 }
func outb(port uint16, val uint8)  {}
func outw(port uint16, val uint16) {}
func outd(port uint16, val uint32) {}
