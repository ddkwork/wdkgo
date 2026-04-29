package main

import "unsafe"

type HookKind int

const (
	HookExec HookKind = iota
	HookRead
	HookWrite
	HookReadWrite
)

func (k HookKind) String() string {
	switch k {
	case HookExec:
		return "EXEC"
	case HookRead:
		return "READ"
	case HookWrite:
		return "WRITE"
	case HookReadWrite:
		return "RW"
	default:
		return "??"
	}
}

func (k HookKind) Permission() PagePermission {
	switch k {
	case HookExec:
		return PagePermission{Read: true, Write: false, Exec: false}
	case HookRead:
		return PagePermission{Read: false, Write: true, Exec: true}
	case HookWrite:
		return PagePermission{Read: true, Write: false, Exec: true}
	case HookReadWrite:
		return PagePermission{Read: false, Write: false, Exec: false}
	default:
		return PagePermission{}
	}
}

type Breakpoint struct {
	Address      uint64
	OriginalByte uint8
	Active       bool
}

type EptHook struct {
	PhysAddr   uint64
	VirtAddr   uint64
	Cr3        CR3_TYPE
	FakePagePa uint64

	OriginalEntry *EptEntry
	ModifiedEntry *EptEntry

	Kind     HookKind
	IsActive bool

	breakpoints []Breakpoint
	fakePage    [PAGE_SIZE]byte
	coreId      uint32
}

func NewEptHook(va uint64, cr3 CR3_TYPE, kind HookKind) (*EptHook, error) {
	pageBase := va & ^(uint64(PAGE_SIZE - 1))
	pa := VirtToPhys(uintptr(pageBase), cr3)
	if pa == 0 || pa >= uint64(SIZE_512_GB) {
		return nil, ErrInvalidAddress
	}

	hook := &EptHook{
		PhysAddr:    pa,
		VirtAddr:    va,
		Cr3:         cr3,
		Kind:        kind,
		IsActive:    true,
		breakpoints: make([]Breakpoint, 0, MAX_HIDDEN_BREAKPOINTS),
	}

	copyPageToFake(hook, uintptr(pageBase), cr3)

	if kind == HookExec || kind == HookReadWrite {
		offset := va & (uint64(PAGE_SIZE) - 1)
		hook.fakePage[offset] = 0xCC
		hook.AddBreakpoint(va)
	}

	hook.FakePagePa = uint64(uintptr(unsafe.Pointer(&hook.fakePage[0])))
	return hook, nil
}

func copyPageToFake(h *EptHook, pageBase uintptr, cr3 CR3_TYPE) {
	srcVa := VirtToPhys(pageBase, cr3)
	dstVa := uint64(uintptr(unsafe.Pointer(&h.fakePage[0])))

	for i := uint64(0); i < uint64(PAGE_SIZE); i += 8 {
		val := readPhysicalMemory(srcVa + i)
		writePhysicalMemory(dstVa+i, val)
	}
}

func (h *EptHook) AddBreakpoint(addr uint64) bool {
	if len(h.breakpoints) >= int(MAX_HIDDEN_BREAKPOINTS) {
		return false
	}
	offset := addr - (h.VirtAddr & ^(uint64(PAGE_SIZE) - 1))
	if offset >= uint64(PAGE_SIZE) {
		return false
	}

	for _, bp := range h.breakpoints {
		if bp.Address == addr && bp.Active {
			return false
		}
	}

	bp := Breakpoint{
		Address:      addr,
		OriginalByte: h.fakePage[offset],
		Active:       true,
	}
	h.breakpoints = append(h.breakpoints, bp)
	h.fakePage[offset] = 0xCC
	LogDebug("HOOK: Added BP at 0x%X (page=0x%X)", addr, h.PhysAddr)
	return true
}

func (h *EptHook) RemoveBreakpoint(addr uint64) bool {
	for i, bp := range h.breakpoints {
		if bp.Address == addr {
			offset := addr - (h.VirtAddr & ^(uint64(PAGE_SIZE) - 1))
			h.fakePage[offset] = bp.OriginalByte
			h.breakpoints = append(h.breakpoints[:i], h.breakpoints[i+1:]...)
			LogDebug("HOOK: Removed BP at 0x%X", addr)
			return true
		}
	}
	return false
}

func (h *EptHook) FindBreakpoint(addr uint64) *Breakpoint {
	for i := range h.breakpoints {
		if h.breakpoints[i].Address == addr && h.breakpoints[i].Active {
			return &h.breakpoints[i]
		}
	}
	return nil
}

func (h *EptHook) BreakpointCount() int {
	count := 0
	for _, bp := range h.breakpoints {
		if bp.Active {
			count++
		}
	}
	return count
}

func (h *EptHook) IsHiddenBp() bool {
	return h.Kind == HookExec || h.Kind == HookReadWrite
}

func (h *EptHook) ContainsAddress(va uint64) bool {
	pageBase := va & ^(uint64(PAGE_SIZE) - 1)
	hookBase := h.VirtAddr & ^(uint64(PAGE_SIZE) - 1)
	return pageBase == hookBase
}

type HookManager struct {
	hooks          []*EptHook
	capacity       int
	lock           *Spinlock
	mtfRestoreList []*EptHook
}

func NewHookManager(capacity int) *HookManager {
	return &HookManager{
		hooks:          make([]*EptHook, 0, capacity),
		capacity:       capacity,
		lock:           NewSpinlock(),
		mtfRestoreList: make([]*EptHook, 0, 16),
	}
}

func (m *HookManager) Lock()   { m.lock.Lock() }
func (m *HookManager) Unlock() { m.lock.Unlock() }

func (m *HookManager) Install(coreId uint32, va uint64, cr3 CR3_TYPE, kind HookKind) error {
	m.Lock()
	defer m.Unlock()

	existing := m.findByVirtAddr(va)
	if existing != nil {
		return m.updateExistingHook(existing, va, kind)
	}

	if len(m.hooks) >= m.capacity {
		return ErrHookLimitReached
	}

	hook, err := NewEptHook(va, cr3, kind)
	if err != nil {
		return err
	}
	hook.coreId = coreId

	table := gHyp.eptTable(coreId)
	if table == nil {
		return ErrNoEptTable
	}

	if err := table.SplitLargePage(hook.PhysAddr); err != nil {
		return err
	}

	pte, err := table.Lookup(hook.PhysAddr)
	if err != nil || pte == nil {
		return ErrNotMapped
	}

	hook.OriginalEntry = pte.Clone()
	hook.ModifiedEntry = &EptEntry{}
	hook.ModifiedEntry.SetRaw(pte.Raw())
	hook.ModifiedEntry.SetPhysAddr(hook.FakePagePa)
	hook.ModifiedEntry.SetPermission(kind.Permission())

	*pte = *hook.ModifiedEntry
	table.Invalidate()

	m.hooks = append(m.hooks, hook)
	kindStr := kind.String()
	LogInfo("HOOK: Installed %s hook at VA=0x%X PA=0x%X", kindStr, va, hook.PhysAddr)
	return nil
}

func (m *HookManager) updateExistingHook(existing *EptHook, va uint64, kind HookKind) error {
	if kind == HookExec || kind == HookReadWrite {
		if !existing.AddBreakpoint(va) {
			return ErrTooManyBreakpoints
		}
	}
	return nil
}

func (m *HookManager) Remove(physAddr uint64) error {
	m.Lock()
	defer m.Unlock()

	for i, hook := range m.hooks {
		if hook.PhysAddr == physAddr {
			table := gHyp.eptTable(hook.coreId)
			if table != nil && hook.OriginalEntry != nil {
				pte, err := table.Lookup(physAddr)
				if err == nil && pte != nil {
					*pte = *hook.OriginalEntry
				}
			}

			m.hooks = append(m.hooks[:i], m.hooks[i+1:]...)
			invalidateEpt()
			LogInfo("HOOK: Removed hook at PA=0x%X", physAddr)
			return nil
		}
	}
	return ErrHookNotFound
}

func (m *HookManager) RemoveByVirtAddr(va uint64) error {
	m.Lock()
	defer m.Unlock()

	for i, hook := range m.hooks {
		if hook.ContainsAddress(va) {
			table := gHyp.eptTable(hook.coreId)
			if table != nil && hook.OriginalEntry != nil {
				pte, err := table.Lookup(hook.PhysAddr)
				if err == nil && pte != nil {
					*pte = *hook.OriginalEntry
				}
			}

			m.hooks = append(m.hooks[:i], m.hooks[i+1:]...)
			invalidateEpt()
			LogInfo("HOOK: Removed hook at VA=0x%X", va)
			return nil
		}
	}
	return ErrHookNotFound
}

func (m *HookManager) findByPhysAddr(pa uint64) *EptHook {
	for _, hook := range m.hooks {
		if hook.PhysAddr == pa {
			return hook
		}
	}
	return nil
}

func (m *HookManager) findByVirtAddr(va uint64) *EptHook {
	for _, hook := range m.hooks {
		if hook.ContainsAddress(va) {
			return hook
		}
	}
	return nil
}

func (m *HookManager) FindByPhysAddr(pa uint64) *EptHook {
	m.Lock()
	defer m.Unlock()
	return m.findByPhysAddr(pa)
}

func (m *HookManager) FindByVirtAddr(va uint64) *EptHook {
	m.Lock()
	defer m.Unlock()
	return m.findByVirtAddr(va)
}

func (m *HookManager) HandleExecHook(v *VCPU, gpa uint64) bool {
	hook := m.FindByPhysAddr(gpa)
	if hook == nil || !hook.IsActive || !hook.IsHiddenBp() {
		return false
	}

	bp := hook.FindBreakpoint(v.LastVmexitRip)
	if bp == nil {
		return false
	}

	v.IncrementRip = false
	v.RegisterBreakOnMtf = true
	v.IgnoreOneMtf = true
	v.MtfEptHookRestorePoint = uintptr(unsafe.Pointer(hook))

	restoreOriginalAndInjectBp(v, hook, bp)
	LogTrace("HOOK: Exec hit at RIP=0x%X", v.LastVmexitRip)
	return true
}

func (m *HookManager) HandleReadWriteHook(v *VCPU, gpa uint64, isRead, isWrite bool) bool {
	hook := m.FindByPhysAddr(gpa)
	if hook == nil || !hook.IsActive {
		return false
	}

	if hook.Kind == HookRead && !isRead {
		return false
	}
	if hook.Kind == HookWrite && !isWrite {
		return false
	}

	v.IncrementRip = false
	v.RegisterBreakOnMtf = true
	v.IgnoreOneMtf = true
	v.MtfEptHookRestorePoint = uintptr(unsafe.Pointer(hook))

	singleStepOnFakePage(v, hook)
	kindStr := hook.Kind.String()
	LogTrace("HOOK: %s hit at GPA=0x%X", kindStr, gpa)
	return true
}

func (m *HookManager) HandleBreakpoint(v *VCPU) bool {

	hook := m.FindByVirtAddr(v.LastVmexitRip)
	if hook == nil || !hook.IsActive {
		return false
	}

	bp := hook.FindBreakpoint(v.LastVmexitRip)
	if bp == nil {
		return false
	}

	v.IncrementRip = false
	v.RegisterBreakOnMtf = true
	v.IgnoreOneMtf = true
	v.MtfEptHookRestorePoint = uintptr(unsafe.Pointer(hook))

	restoreOriginalAndInjectBp(v, hook, bp)
	return true
}

func (m *HookManager) HandleMtfRestore(v *VCPU) bool {
	if v.MtfEptHookRestorePoint == 0 {
		return false
	}

	hook := (*EptHook)(unsafe.Pointer(v.MtfEptHookRestorePoint))
	if hook == nil {
		v.MtfEptHookRestorePoint = 0
		return false
	}

	table := gHyp.eptTable(v.CoreId)
	if table == nil {
		v.MtfEptHookRestorePoint = 0
		return false
	}

	pte, err := table.Lookup(hook.PhysAddr)
	if err != nil || pte == nil {
		v.MtfEptHookRestorePoint = 0
		return false
	}

	*pte = *hook.ModifiedEntry
	table.Invalidate()

	v.MtfEptHookRestorePoint = 0
	LogTrace("HOOK: MTF restore for PA=0x%X", hook.PhysAddr)
	return true
}

func (m *HookManager) RestoreAll(v *VCPU) {
	m.Lock()
	defer m.Unlock()

	if v == nil {
		// Restore hooks for all cores
		for i := uint32(0); i < uint32(len(gHyp.eptTables)); i++ {
			m.restoreAllForCoreLocked(i)
		}
		invalidateEpt()
		LogDebug("HOOK: Restored all hooks for all cores")
		return
	}

	coreId := v.CoreId
	if coreId == 0 {
		coreId = getCurrentProcessorNumber()
	}

	table := gHyp.eptTable(coreId)
	if table == nil {
		return
	}

	for _, hook := range m.hooks {
		if !hook.IsActive || hook.OriginalEntry == nil {
			continue
		}

		pte, err := table.Lookup(hook.PhysAddr)
		if err != nil || pte == nil {
			continue
		}

		*pte = *hook.OriginalEntry
	}

	invalidateEpt()
	LogDebug("HOOK: Restored all hooks for core %d", coreId)
}

func (m *HookManager) restoreAllForCoreLocked(coreId uint32) {
	table := gHyp.eptTable(coreId)
	if table == nil {
		return
	}

	for _, hook := range m.hooks {
		if !hook.IsActive || hook.OriginalEntry == nil {
			continue
		}

		pte, err := table.Lookup(hook.PhysAddr)
		if err == nil && pte != nil {
			*pte = *hook.OriginalEntry
		}
	}
}

func (m *HookManager) RestoreAllForCore(coreId uint32) {
	m.Lock()
	defer m.Unlock()

	table := gHyp.eptTable(coreId)
	if table == nil {
		return
	}

	for _, hook := range m.hooks {
		if !hook.IsActive || hook.OriginalEntry == nil {
			continue
		}

		pte, err := table.Lookup(hook.PhysAddr)
		if err == nil && pte != nil {
			*pte = *hook.OriginalEntry
		}
	}

	invalidateEpt()
}

func (m *HookManager) Count() int {
	m.Lock()
	defer m.Unlock()
	return len(m.hooks)
}

func (m *HookManager) ActiveCount() int {
	m.Lock()
	defer m.Unlock()
	count := 0
	for _, h := range m.hooks {
		if h.IsActive {
			count++
		}
	}
	return count
}

func (m *HookManager) RemoveAll() {
	m.Lock()
	defer m.Unlock()
	m.hooks = nil
}

func (m *HookManager) DisableAll() {
	m.Lock()
	defer m.Unlock()
	for _, h := range m.hooks {
		h.IsActive = false
	}
}

func (m *HookManager) EnableAll() {
	m.Lock()
	defer m.Unlock()
	for _, h := range m.hooks {
		h.IsActive = true
	}
}

func restoreOriginalAndInjectBp(v *VCPU, hook *EptHook, bp *Breakpoint) {
	table := gHyp.eptTable(v.CoreId)
	if table == nil {
		return
	}

	pte, err := table.Lookup(hook.PhysAddr)
	if err != nil || pte == nil {
		return
	}

	originalCopy := hook.OriginalEntry.Clone()
	*pte = *originalCopy
	table.Invalidate()

	offset := bp.Address - (hook.VirtAddr & ^(uint64(PAGE_SIZE) - 1))
	guestInstrPa := originalCopy.PhysAddr() + offset

	writePhysicalMemory(guestInstrPa, uint64(0xCC))
}

func singleStepOnFakePage(v *VCPU, hook *EptHook) {
	table := gHyp.eptTable(v.CoreId)
	if table == nil {
		return
	}

	pte, err := table.Lookup(hook.PhysAddr)
	if err != nil || pte == nil {
		return
	}

	perm := hook.ModifiedEntry.Permission()
	perm.Read = true
	perm.Write = true
	perm.Exec = true
	pte.SetPermission(perm)
	table.Invalidate()
}

var (
	ErrInvalidAddress     = fmtError("invalid virtual address")
	ErrHookLimitReached   = fmtError("maximum hook limit reached")
	ErrNoEptTable         = fmtError("no Ept table available")
	ErrHookNotFound       = fmtError("hook not found")
	ErrTooManyBreakpoints = fmtError("too many breakpoints on this page")
)

func eptHookSetReadHook(addr uint64, cr3 CR3_TYPE) {
	h := getHypervisor()
	if h == nil {
		return
	}
	v := h.CurrentVcpu()
	if v != nil {
		gHooks.Install(v.CoreId, addr, cr3, HookRead)
	}
}

func eptHookSetWriteHook(addr uint64, cr3 CR3_TYPE) {
	h := getHypervisor()
	if h == nil {
		return
	}
	v := h.CurrentVcpu()
	if v != nil {
		gHooks.Install(v.CoreId, addr, cr3, HookWrite)
	}
}

func eptHookSetExecHook(addr uint64, cr3 CR3_TYPE) {
	h := getHypervisor()
	if h == nil {
		return
	}
	v := h.CurrentVcpu()
	if v != nil {
		gHooks.Install(v.CoreId, addr, cr3, HookExec)
	}
}

func eptHookSetReadWriteHook(addr uint64, cr3 CR3_TYPE) {
	h := getHypervisor()
	if h == nil {
		return
	}
	v := h.CurrentVcpu()
	if v != nil {
		gHooks.Install(v.CoreId, addr, cr3, HookReadWrite)
	}
}
