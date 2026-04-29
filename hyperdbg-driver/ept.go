package main

import (
	"fmt"
	"unsafe"
)

const (
	EptRead       uint64 = 1 << 0
	EptWrite      uint64 = 1 << 1
	EptExec       uint64 = 1 << 2
	EptIgnorePat  uint64 = 1 << 6
	EptSuppressVe uint64 = 1 << 63
	EptLargePage  uint64 = 1 << 7
	EptPhysMask   uint64 = 0xFFFFFFFFFF000
	EptTypeMask   uint64 = 0x38
	EptTypeShift         = 3

	MemTypeUncacheable    uint8 = 0
	MemTypeWriteCombining uint8 = 1
	MemTypeWriteThrough   uint8 = 4
	MemTypeWriteProtected uint8 = 5
	MemTypeWriteBack      uint8 = 6
)

type PagePermission struct {
	Read, Write, Exec bool
}

func (p PagePermission) ToRaw() uint64 {
	var raw uint64
	if p.Read {
		raw |= EptRead
	}
	if p.Write {
		raw |= EptWrite
	}
	if p.Exec {
		raw |= EptExec
	}
	return raw
}

func PagePermissionFromRaw(raw uint64) PagePermission {
	return PagePermission{
		Read:  raw&EptRead != 0,
		Write: raw&EptWrite != 0,
		Exec:  raw&EptExec != 0,
	}
}

type EptEntry struct {
	raw uint64
}

func NewEptEntry() *EptEntry              { return &EptEntry{} }
func (e *EptEntry) Raw() uint64           { return e.raw }
func (e *EptEntry) SetRaw(v uint64)       { e.raw = v }
func (e *EptEntry) IsPresent() bool       { return e.raw&EptRead != 0 }
func (e *EptEntry) PhysAddr() uint64      { return e.raw & EptPhysMask }
func (e *EptEntry) SetPhysAddr(pa uint64) { e.raw = (e.raw & ^EptPhysMask) | (pa & EptPhysMask) }
func (e *EptEntry) MemoryType() uint8     { return uint8((e.raw >> EptTypeShift) & 0x7) }
func (e *EptEntry) SetMemoryType(t uint8) {
	e.raw = (e.raw & ^EptTypeMask) | (uint64(t)&0x7)<<EptTypeShift
}
func (e *EptEntry) Permission() PagePermission { return PagePermissionFromRaw(e.raw) }
func (e *EptEntry) SetPermission(p PagePermission) {
	e.raw = (e.raw & ^(EptRead | EptWrite | EptExec)) | p.ToRaw()
}
func (e *EptEntry) IsLargePage() bool { return e.raw&EptLargePage != 0 }
func (e *EptEntry) IgnorePat() bool   { return e.raw&EptIgnorePat != 0 }
func (e *EptEntry) SuppressVe() bool  { return e.raw&EptSuppressVe != 0 }
func (e *EptEntry) SetIgnorePat(v bool) {
	if v {
		e.raw |= EptIgnorePat
	} else {
		e.raw &= ^EptIgnorePat
	}
}
func (e *EptEntry) SetSuppressVe(v bool) {
	if v {
		e.raw |= EptSuppressVe
	} else {
		e.raw &= ^EptSuppressVe
	}
}
func (e *EptEntry) SetLargePage(v bool) {
	if v {
		e.raw |= EptLargePage
	} else {
		e.raw &= ^EptLargePage
	}
}
func (e *EptEntry) Clone() *EptEntry { return &EptEntry{raw: e.raw} }

type Level = int

const (
	LvlPml4 Level = iota
	LvlPdpt
	LvlPd
	LvlPt
)

type Pml1Table struct {
	entries [512]EptEntry
}

type SplitRecord struct {
	pdEntry  *EptEntry
	pt       *Pml1Table
	refCount int
	gpaBase  uint64
}

type EptTable struct {
	pml4 [512]EptEntry

	pdptTables [511]*[512]EptEntry
	pdTables   [512][512]*EptEntry

	splits []SplitRecord

	mtrrState *MtrrState
	lock      *Spinlock
}

func NewEptTable() *EptTable {
	t := &EptTable{
		mtrrState: NewMtrrState(),
		lock:      NewSpinlock(),
	}
	t.initIdentityMap()
	return t
}

func (t *EptTable) initIdentityMap() {
	defaultPerm := PagePermission{Read: true, Write: true, Exec: true}
	for i := range 512 {
		if i == 0 {
			continue
		}
		for j := range 512 {
			pd := t.pdTables[i-1][j]
			pd.SetPhysAddr(uint64(i-1)*uint64(SIZE_1_GB) + uint64(j)*uint64(SIZE_2_MB))
			pd.SetPermission(defaultPerm)
			pd.SetMemoryType(MemTypeWriteBack)
			pd.SetLargePage(true)
			pd.SetIgnorePat(true)
		}
	}
}

func (t *EptTable) Lock()   { t.lock.Lock() }
func (t *EptTable) Unlock() { t.lock.Unlock() }

func (t *EptTable) pml4Index(gpa uint64) int { return int((gpa >> 39) & 0x1FF) }
func (t *EptTable) pdptIndex(gpa uint64) int { return int((gpa >> 30) & 0x1FF) }
func (t *EptTable) pdIndex(gpa uint64) int   { return int((gpa >> 21) & 0x1FF) }
func (t *EptTable) ptIndex(gpa uint64) int   { return int((gpa >> 12) & 0x1FF) }

func (t *EptTable) Walk(gpa uint64, level Level) (*EptEntry, error) {
	p4Idx := t.pml4Index(gpa)
	if level == LvlPml4 {
		return &t.pml4[p4Idx], nil
	}

	p4 := t.pml4[p4Idx]
	if !p4.IsPresent() || p4Idx == 0 {
		return nil, ErrNotPresent
	}

	pdptIdx := t.pdptIndex(gpa)
	if level == LvlPdpt {
		pdpt := t.pdptTables[p4Idx-1]
		if pdpt == nil {
			return nil, ErrNotPresent
		}
		return &(*pdpt)[pdptIdx], nil
	}

	pdpt := t.pdptTables[p4Idx-1]
	if pdpt == nil {
		return nil, ErrNotPresent
	}
	pdptEntry := (*pdpt)[pdptIdx]
	if !pdptEntry.IsPresent() {
		return nil, ErrNotPresent
	}

	pdIdx := t.pdIndex(gpa)
	if level == LvlPd {
		pd := t.pdTables[p4Idx-1][pdIdx]
		if pd == nil {
			return nil, ErrNotPresent
		}
		return pd, nil
	}

	pd := t.pdTables[p4Idx-1][pdIdx]
	if pd == nil || !pd.IsPresent() {
		return nil, ErrNotPresent
	}

	if pd.IsLargePage() {
		return nil, ErrLargePage
	}

	ptIdx := t.ptIndex(gpa)
	for i := range t.splits {
		if t.splits[i].pdEntry == pd && ptIdx < 512 {
			return &t.splits[i].pt.entries[ptIdx], nil
		}
	}
	return nil, ErrNotSplit
}

func (t *EptTable) Lookup(gpa uint64) (*EptEntry, error) {
	pte, err := t.Walk(gpa, LvlPt)
	if err == nil && pte != nil && pte.IsPresent() {
		return pte, nil
	}

	pd, err := t.Walk(gpa, LvlPd)
	if err == nil && pd != nil && pd.IsPresent() && pd.IsLargePage() {
		return pd, nil
	}

	return nil, ErrNotMapped
}

func (t *EptTable) ResolveForModify(gpa uint64) (*EptEntry, error) {
	pte, err := t.Walk(gpa, LvlPt)
	if err == nil && pte != nil {
		return pte, nil
	}

	pd, err := t.Walk(gpa, LvlPd)
	if err == nil && pd != nil && pd.IsPresent() && pd.IsLargePage() {
		return pd, nil
	}

	return nil, ErrNotMapped
}

func (t *EptTable) SplitLargePage(gpa uint64) error {
	t.Lock()
	defer t.Unlock()

	pd, err := t.ResolveForModify(gpa)
	if err != nil || pd == nil || !pd.IsLargePage() {
		return nil
	}

	for _, s := range t.splits {
		if s.pdEntry == pd {
			return nil
		}
	}

	pt := new(Pml1Table)
	basePa := pd.PhysAddr()
	memType := pd.MemoryType()
	ignorePat := pd.IgnorePat()
	suppressVe := pd.SuppressVe()

	template := EptEntry{}
	template.SetPermission(PagePermission{Read: true, Write: true, Exec: true})
	template.SetMemoryType(memType)
	template.SetIgnorePat(ignorePat)
	template.SetSuppressVe(suppressVe)

	for i := range 512 {
		pt.entries[i] = template
		pt.entries[i].SetPhysAddr(basePa + uint64(i)*uint64(PAGE_SIZE))
	}

	ptrEntry := EptEntry{}
	ptrEntry.SetPermission(PagePermission{Read: true, Write: true, Exec: true})
	ptrEntry.SetMemoryType(MemTypeWriteBack)
	ptrEntry.SetPhysAddr(uint64(uintptr(unsafe.Pointer(&pt.entries[0]))))

	oldRaw := pd.Raw()
	*pd = ptrEntry

	gpaBase := gpa & ^(uint64(SIZE_2_MB) - 1)
	t.splits = append(t.splits, SplitRecord{
		pdEntry:  pd,
		pt:       pt,
		refCount: 1,
		gpaBase:  gpaBase,
	})

	LogDebug("EPT: Split large page at GPA 0x%X (old raw=0x%X)", gpaBase, oldRaw)
	return nil
}

func (t *EptTable) MergeLargePage(gpa uint64) error {
	t.Lock()
	defer t.Unlock()

	pd, err := t.Walk(gpa, LvlPd)
	if err != nil || pd == nil {
		return ErrNotMapped
	}

	for i, s := range t.splits {
		if s.pdEntry == pd {
			allDefault := true
			defaultPerm := PagePermission{Read: true, Write: true, Exec: true}
			for j := range s.pt.entries {
				perm := s.pt.entries[j].Permission()
				pa := s.pt.entries[j].PhysAddr()
				expectedPa := s.gpaBase + uint64(j)*uint64(PAGE_SIZE)
				if perm != defaultPerm || pa != expectedPa {
					allDefault = false
					break
				}
			}

			if allDefault && s.refCount <= 1 {
				memType := t.mtrrState.Lookup(s.gpaBase)
				restored := EptEntry{}
				restored.SetPhysAddr(s.gpaBase)
				restored.SetPermission(defaultPerm)
				restored.SetMemoryType(memType)
				restored.SetLargePage(true)
				restored.SetIgnorePat(true)
				*pd = restored

				t.splits = append(t.splits[:i], t.splits[i+1:]...)
				LogDebug("EPT: Merged page at GPA 0x%X", s.gpaBase)
				return nil
			}
			return fmtError("cannot merge: page has non-default entries or multiple refs")
		}
	}
	return fmtError("no split found for this address")
}

func (t *EptTable) Map(gpa uint64, pa uint64, perm PagePermission, memType uint8) error {
	t.Lock()
	defer t.Unlock()

	pd, err := t.ResolveForModify(gpa)
	if err == nil && pd != nil && pd.IsLargePage() {
		if splitErr := t.splitLargePageUnlocked(gpa); splitErr != nil {
			return splitErr
		}
	}

	pte, err := t.Walk(gpa, LvlPt)
	if err != nil {
		return err
	}

	pte.SetPhysAddr(pa)
	pte.SetPermission(perm)
	pte.SetMemoryType(memType)
	return nil
}

func (t *EptTable) Unmap(gpa uint64) error {
	t.Lock()
	defer t.Unlock()

	pte, err := t.resolveForModifyUnlocked(gpa)
	if err != nil {
		return err
	}
	*pte = EptEntry{}
	return nil
}

func (t *EptTable) Protect(gpa uint64, perm PagePermission) error {
	t.Lock()
	defer t.Unlock()

	pte, err := t.resolveForModifyUnlocked(gpa)
	if err != nil {
		return err
	}
	pte.SetPermission(perm)
	return nil
}

func (t *EptTable) splitLargePageUnlocked(gpa uint64) error {
	pd, err := t.resolveForModifyUnlocked(gpa)
	if err != nil || pd == nil || !pd.IsLargePage() {
		return nil
	}

	for _, s := range t.splits {
		if s.pdEntry == pd {
			return nil
		}
	}

	pt := new(Pml1Table)
	basePa := pd.PhysAddr()
	memType := pd.MemoryType()
	ignorePat := pd.IgnorePat()
	suppressVe := pd.SuppressVe()

	template := EptEntry{}
	template.SetPermission(PagePermission{Read: true, Write: true, Exec: true})
	template.SetMemoryType(memType)
	template.SetIgnorePat(ignorePat)
	template.SetSuppressVe(suppressVe)

	for i := range 512 {
		pt.entries[i] = template
		pt.entries[i].SetPhysAddr(basePa + uint64(i)*uint64(PAGE_SIZE))
	}

	ptrEntry := EptEntry{}
	ptrEntry.SetPermission(PagePermission{Read: true, Write: true, Exec: true})
	ptrEntry.SetMemoryType(MemTypeWriteBack)
	ptrEntry.SetPhysAddr(uint64(uintptr(unsafe.Pointer(&pt.entries[0]))))

	*pd = ptrEntry

	gpaBase := gpa & ^(uint64(SIZE_2_MB) - 1)
	t.splits = append(t.splits, SplitRecord{
		pdEntry:  pd,
		pt:       pt,
		refCount: 1,
		gpaBase:  gpaBase,
	})
	return nil
}

func (t *EptTable) resolveForModifyUnlocked(gpa uint64) (*EptEntry, error) {
	pte, err := t.Walk(gpa, LvlPt)
	if err == nil && pte != nil {
		return pte, nil
	}

	pd, err := t.Walk(gpa, LvlPd)
	if err == nil && pd != nil && pd.IsPresent() && pd.IsLargePage() {
		return pd, nil
	}

	return nil, ErrNotMapped
}

func (t *EptTable) BuildEptPointer() uint64 {
	var result uint64
	result |= uint64(t.mtrrState.DefaultMemoryType()) << EptTypeShift
	result |= EptRead | EptExec
	result |= uint64(uintptr(unsafe.Pointer(&t.pml4[0]))) & EptPhysMask
	return result
}

func (t *EptTable) Invalidate() {
	inveptAllContexts()
}

func (t *EptTable) SplitCount() int { return len(t.splits) }

type MtrrState struct {
	ranges                 []MTRR_RANGE_DESCRIPTOR
	defaultType            uint8
	numRanges              uint32
	fixedRangesInitialized bool
}

func NewMtrrState() *MtrrState {
	m := &MtrrState{
		ranges:      make([]MTRR_RANGE_DESCRIPTOR, 0, MAX_MTRR_ENTRIES),
		defaultType: MemTypeWriteBack,
	}
	m.readMtrrsFromHardware()
	return m
}

func (m *MtrrState) DefaultMemoryType() uint8 { return m.defaultType }

func (m *MtrrState) readMtrrsFromHardware() {
	msrMtrrcap := __readmsr(0xFE)
	vcnt := min(int((msrMtrrcap>>40)&0xFF), int(MAX_MTRR_ENTRIES))

	msrMtrrDefType := __readmsr(0x2FF)
	m.defaultType = uint8(msrMtrrDefType & 0xFF)

	m.ranges = m.ranges[:0]
	for i := range vcnt {
		baseMsr := __readmsr(uint32(0x200 + i*2))
		maskMsr := __readmsr(uint32(0x201 + i*2))

		if maskMsr&0x800 == 0 {
			continue
		}

		basePa := baseMsr & 0xFFFFF000
		maskBits := maskMsr & 0xFFFFF000
		size := (^maskBits) + 1
		memType := uint8(baseMsr & 0xFF)

		m.ranges = append(m.ranges, MTRR_RANGE_DESCRIPTOR{
			PhysicalBaseAddress: uintptr(basePa),
			PhysicalEndAddress:  uintptr(basePa + size),
			MemoryType:          memType,
			FixedRange:          false,
		})
	}
	m.numRanges = uint32(len(m.ranges))
}

func (m *MtrrState) Lookup(pa uint64) uint8 {
	for _, r := range m.ranges {
		if pa >= uint64(r.PhysicalBaseAddress) && pa < uint64(r.PhysicalEndAddress) {
			return r.MemoryType
		}
	}
	return m.defaultType
}

func (m *MtrrState) Ranges() []MTRR_RANGE_DESCRIPTOR { return m.ranges }

var (
	ErrNotPresent   = fmtError("EPT entry not present")
	ErrLargePage    = fmtError("entry is a large page")
	ErrNotSplit     = fmtError("large page not split")
	ErrInvalidLevel = fmtError("invalid page table level")
	ErrNotMapped    = fmtError("address not mapped in EPT")
)

type eptError string

func (e eptError) Error() string { return string(e) }

func fmtError(format string, args ...any) error { return fmt.Errorf(format, args...) }

func (t *EptTable) RestoreEntry(physAddr uint64, original *EptEntry) error {
	t.Lock()
	defer t.Unlock()

	gpa := physAddr
	pte, err := t.resolveForModifyUnlocked(gpa)
	if err != nil {
		return err
	}

	if original != nil {
		*pte = *original
	} else {
		*pte = EptEntry{}
	}

	inveptSingleContext(t.BuildEptPointer())
	return nil
}

func (t *EptTable) GetEntry(gpa uint64) (*EptEntry, error) {
	pte, err := t.ResolveForModify(gpa)
	if err != nil {
		return nil, err
	}
	return pte, nil
}

func (t *EptTable) CloneEntry(gpa uint64) (*EptEntry, error) {
	pte, err := t.GetEntry(gpa)
	if err != nil {
		return nil, err
	}
	return pte.Clone(), nil
}
