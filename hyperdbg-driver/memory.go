package main

import "unsafe"

type MemoryType uint8

const (
	MemUncacheable    MemoryType = 0
	MemWriteCombining MemoryType = 1
	MemWriteThrough   MemoryType = 4
	MemWriteProtected MemoryType = 5
	MemWriteBack      MemoryType = 6
)

func (m MemoryType) String() string {
	switch m {
	case MemUncacheable:
		return "UC"
	case MemWriteCombining:
		return "WC"
	case MemWriteThrough:
		return "WT"
	case MemWriteProtected:
		return "WP"
	case MemWriteBack:
		return "WB"
	default:
		return "??"
	}
}

type MtrrRange struct {
	BaseAddr uint64
	EndAddr  uint64
	Type     MemoryType
	Valid    bool
}

type MtrrCapabilities struct {
	VarCnt         int
	FixedSupported bool
}

type MemoryMtrrState struct {
	DefaultType    MemoryType
	Capabilities   MtrrCapabilities
	VariableRanges [MAX_MTRR_ENTRIES]MtrrRange
}

func newMemoryManager() *MemoryManager {
	return &MemoryManager{
		mtrr: &MemoryMtrrState{
			DefaultType: MemUncacheable,
		},
	}
}

type MemoryManager struct {
	mtrr *MemoryMtrrState
}

func (m *MemoryManager) QueryMtrrForPa(pa uint64) MemoryType {
	for i := range m.mtrr.VariableRanges {
		r := &m.mtrr.VariableRanges[i]
		if r.Valid && pa >= r.BaseAddr && pa < r.EndAddr {
			return r.Type
		}
	}
	return m.mtrr.DefaultType
}

func (m *MemoryManager) SetDefaultMemType(t MemoryType) {
	m.mtrr.DefaultType = t
}

func (m *MemoryManager) AddVariableRange(base, end uint64, t MemoryType) {
	for i := range m.mtrr.VariableRanges {
		if !m.mtrr.VariableRanges[i].Valid {
			m.mtrr.VariableRanges[i] = MtrrRange{
				BaseAddr: base,
				EndAddr:  end,
				Type:     t,
				Valid:    true,
			}
			return
		}
	}
}

func (m *MemoryManager) ClearAllRanges() {
	for i := range m.mtrr.VariableRanges {
		m.mtrr.VariableRanges[i].Valid = false
	}
}

func VirtToPhys(va uintptr, cr3 CR3_TYPE) uint64 {
	pml4Idx := (uint64(va) >> 39) & 0x1FF
	pml3Idx := (uint64(va) >> 30) & 0x1FF
	pml2Idx := (uint64(va) >> 21) & 0x1FF
	pml1Idx := (uint64(va) >> 12) & 0x1FF

	pml4Table := (*[512]PML4E)(unsafe.Pointer(uintptr(cr3.Flags & 0xFFFFFFFFFFFFF000)))
	pml4Entry := pml4Table[pml4Idx]
	if pml4Entry.AsUInt&0x1 == 0 {
		return 0
	}

	pml3Base := uintptr(pml4Entry.AsUInt & 0xFFFFFFFFFF000)
	pml3Table := (*[512]EPT_PDPTE)(unsafe.Pointer(pml3Base))
	pml3Entry := pml3Table[pml3Idx]
	if pml3Entry.AsUInt&0x1 == 0 {
		return 0
	}

	if pml3Entry.AsUInt&(1<<7) != 0 {
		pageBase := pml3Entry.AsUInt & ((1 << 30) - 1)
		offset := uint64(va) & (uint64(SIZE_1_GB) - 1)
		return pageBase + offset
	}

	pml2Base := uintptr(pml3Entry.AsUInt & 0xFFFFFFFFFF000)
	pml2Table := (*[512]EPT_PDE)(unsafe.Pointer(pml2Base))
	pml2Entry := pml2Table[pml2Idx]
	if pml2Entry.AsUInt&0x1 == 0 {
		return 0
	}

	if pml2Entry.AsUInt&(1<<7) != 0 {
		pageBase := pml2Entry.AsUInt & ((1 << 21) - 1)
		offset := uint64(va) & (uint64(SIZE_2_MB) - 1)
		return pageBase + offset
	}

	pml1Base := uintptr(pml2Entry.AsUInt & 0xFFFFFFFFFF000)
	pml1Table := (*[512]EPT_PTE)(unsafe.Pointer(pml1Base))
	pml1Entry := pml1Table[pml1Idx]
	if pml1Entry.AsUInt&0x1 == 0 {
		return 0
	}

	pageBase := pml1Entry.AsUInt & 0xFFFFFFFFFF000
	offset := uint64(va) & (uint64(PAGE_SIZE) - 1)
	return pageBase + offset
}

func PhysToVirt(pa uint64, cr3 CR3_TYPE) uintptr {
	return 0
}

func CopyPage(src, dst uintptr, cr3 CR3_TYPE) {
	srcPa := VirtToPhys(src, cr3)
	dstPa := VirtToPhys(dst, cr3)

	if srcPa == 0 || dstPa == 0 {
		return
	}

	srcPtr := (*[PAGE_SIZE]byte)(unsafe.Pointer(uintptr(srcPa)))
	dstPtr := (*[PAGE_SIZE]byte)(unsafe.Pointer(uintptr(dstPa)))

	copy(dstPtr[:], srcPtr[:])
}

func ReadPhysMem(pa uint64, buf []byte) int {
	ptr := (*byte)(unsafe.Pointer(uintptr(pa)))
	maxLen := uint64(PAGE_SIZE) - (pa % uint64(PAGE_SIZE))
	if int(maxLen) > len(buf) {
		maxLen = uint64(len(buf))
	}

	for i := uint64(0); i < maxLen; i++ {
		buf[i] = *(*byte)(unsafe.Pointer(uintptr(uintptr(unsafe.Pointer(ptr)) + uintptr(i))))
	}
	return int(maxLen)
}

func WritePhysMem(pa uint64, data []byte) int {
	ptr := (*byte)(unsafe.Pointer(uintptr(pa)))
	maxLen := uint64(PAGE_SIZE) - (pa % uint64(PAGE_SIZE))
	if int(maxLen) > len(data) {
		maxLen = uint64(len(data))
	}

	for i := uint64(0); i < maxLen; i++ {
		*(*byte)(unsafe.Pointer(uintptr(uintptr(unsafe.Pointer(ptr)) + uintptr(i)))) = data[i]
	}
	return int(maxLen)
}

func readPhysicalMemory(pa uint64) uint64 {
	ptr := (*uint64)(unsafe.Pointer(uintptr(pa)))
	return *ptr
}

func writePhysicalMemory(pa uint64, val uint64) {
	ptr := (*uint64)(unsafe.Pointer(uintptr(pa)))
	*ptr = val
}

func invalidateEpt() {
	inveptAllContexts()
}
