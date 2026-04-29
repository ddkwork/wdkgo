package main

import (
	"unsafe"

	"solod.dev/so/wdk"
)

type OptimizationLevel = int

const (
	OptNone OptimizationLevel = iota
	OptBasic
	OptAggressive
)

type PoolManager struct {
	poolList       *PoolEntry
	totalAllocated SIZE_T
	totalFreed     SIZE_T
	lock           *Spinlock
	maxPools       int
	currentPools   int
}

type PoolEntry struct {
	Address PVOID
	Size    SIZE_T
	Tag     ULONG
	Next    *PoolEntry
	Prev    *PoolEntry
}

func NewPoolManager(maxEntries int) *PoolManager {
	return &PoolManager{
		lock:     NewSpinlock(),
		maxPools: maxEntries,
	}
}

func (m *PoolManager) Allocate(size SIZE_T, tag ULONG) PVOID {
	m.lock.Lock()
	defer m.lock.Unlock()

	ptr := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), uintptr(size), uint32(tag))
	if ptr == 0 {
		return 0
	}

	entry := AllocatePool(unsafe.Sizeof(PoolEntry{}), POOLTAG)
	if entry == 0 {
		wdk.ExFreePoolWithTag(ptr, POOLTAG)
		return 0
	}

	entryPtr := (*PoolEntry)(unsafe.Pointer(entry))
	entryPtr.Address = uintptr(ptr)
	entryPtr.Size = size
	entryPtr.Tag = tag

	if m.poolList != nil {
		m.poolList.Prev = entryPtr
		entryPtr.Next = m.poolList
	}
	m.poolList = entryPtr

	m.totalAllocated += size
	m.currentPools++

	return uintptr(ptr)
}

func (m *PoolManager) Free(ptr PVOID) {
	if ptr == 0 {
		return
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	entry := m.findEntry(ptr)
	if entry == nil {
		return
	}

	if entry.Next != nil {
		entry.Next.Prev = entry.Prev
	}
	if entry.Prev != nil {
		entry.Prev.Next = entry.Next
	} else if m.poolList == entry {
		m.poolList = entry.Next
	}

	m.totalFreed += entry.Size
	m.currentPools--

	wdk.ExFreePoolWithTag(entry.Address, uint32(entry.Tag))
	FreePool(uintptr(unsafe.Pointer(entry)), POOLTAG)
}

func (m *PoolManager) findEntry(ptr PVOID) *PoolEntry {
	current := m.poolList
	for current != nil {
		if current.Address == ptr {
			return current
		}
		current = current.Next
	}
	return nil
}

func (m *PoolManager) CheckAndPerformDeallocation() int {
	freedCount := 0

	m.lock.Lock()
	current := m.poolList
	for current != nil {
		next := current.Next

		if current.Address != 0 && wdk.MmIsAddressValid(current.Address) {

			wdk.ExFreePoolWithTag(current.Address, uint32(current.Tag))
			m.totalFreed += current.Size
			m.currentPools--
			freedCount++
		}

		current = next
	}
	m.lock.Unlock()

	return freedCount
}

func (m *PoolManager) GetStats() PoolStats {
	m.lock.Lock()
	defer m.lock.Unlock()

	return PoolStats{
		TotalAllocated: m.totalAllocated,
		TotalFreed:     m.totalFreed,
		CurrentPools:   m.currentPools,
		MaxPools:       m.maxPools,
	}
}

type PoolStats struct {
	TotalAllocated SIZE_T
	TotalFreed     SIZE_T
	CurrentPools   int
	MaxPools       int
}

func (s PoolStats) String() string {
	return ""
}

type VmmOptimizer struct {
	level           OptimizationLevel
	eptCacheEnabled bool
	msrBitmapCached bool
	ioBitmapCached  bool
	inveptBatching  bool
}

func NewVmmOptimizer(level OptimizationLevel) *VmmOptimizer {
	o := &VmmOptimizer{level: level}
	o.applyDefaults()
	return o
}

func (o *VmmOptimizer) applyDefaults() {
	switch o.level {
	case OptNone:
		o.eptCacheEnabled = false
		o.msrBitmapCached = false
		o.ioBitmapCached = false
		o.inveptBatching = false
	case OptBasic:
		o.eptCacheEnabled = true
		o.msrBitmapCached = true
		o.ioBitmapCached = false
		o.inveptBatching = false
	case OptAggressive:
		o.eptCacheEnabled = true
		o.msrBitmapCached = true
		o.ioBitmapCached = true
		o.inveptBatching = true
	}
}

func (o *VmmOptimizer) EnableEptCache(enable bool) {
	o.eptCacheEnabled = enable
}

func (o *VmmOptimizer) EnableMsrBitmapCache(enable bool) {
	o.msrBitmapCached = enable
}

func (o *VmmOptimizer) EnableIoBitmapCache(enable bool) {
	o.ioBitmapCached = enable
}

func (o *VmmOptimizer) EnableInveptBatching(enable bool) {
	o.inveptBatching = enable
}

func (o *VmmOptimizer) GetConfig() VmmOptConfig {
	return VmmOptConfig{
		EptCache:       o.eptCacheEnabled,
		MsrBitmapCache: o.msrBitmapCached,
		IoBitmapCache:  o.ioBitmapCached,
		InveptBatching: o.inveptBatching,
	}
}

type VmmOptConfig struct {
	EptCache       bool
	MsrBitmapCache bool
	IoBitmapCache  bool
	InveptBatching bool
}

type EptCache struct {
	entries map[uint64]*EptCacheEntry
	lock    *Spinlock
	maxSize int
}

type EptCacheEntry struct {
	Gpa        uint64
	Pa         uint64
	Entry      EptEntry
	LastAccess uint64
	HitCount   uint64
}

func NewEptCache(maxSize int) *EptCache {
	c := &EptCache{
		entries: make(map[uint64]*EptCacheEntry),
		lock:    NewSpinlock(),
		maxSize: maxSize,
	}
	return c
}

func (c *EptCache) Lookup(gpa uint64) (*EptEntry, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry, ok := c.entries[gpa]
	if !ok {
		return nil, false
	}

	entry.LastAccess = rdtscValue()
	entry.HitCount++
	return &entry.Entry, true
}

func (c *EptCache) Store(gpa uint64, pa uint64, entry EptEntry) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.entries) >= c.maxSize {
		c.evictOldest()
	}

	cacheEntry := &EptCacheEntry{
		Gpa:        gpa,
		Pa:         pa,
		Entry:      entry,
		LastAccess: rdtscValue(),
		HitCount:   1,
	}
	c.entries[gpa] = cacheEntry
}

func (c *EptCache) Invalidate(gpa uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.entries, gpa)
}

func (c *EptCache) InvalidateAll() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.entries = make(map[uint64]*EptCacheEntry)
}

func (c *EptCache) evictOldest() {
	var oldestGpa uint64
	var oldestTime uint64 = ^uint64(0)

	for gpa, entry := range c.entries {
		if entry.LastAccess < oldestTime {
			oldestTime = entry.LastAccess
			oldestGpa = gpa
		}
	}

	if oldestTime != ^uint64(0) {
		delete(c.entries, oldestGpa)
	}
}

func (c *EptCache) GetStats() EptCacheStats {
	c.lock.Lock()
	defer c.lock.Unlock()

	stats := EptCacheStats{
		TotalEntries: len(c.entries),
		MaxEntries:   c.maxSize,
	}

	for _, entry := range c.entries {
		stats.TotalHits += entry.HitCount
	}

	return stats
}

type EptCacheStats struct {
	TotalEntries int
	MaxEntries   int
	TotalHits    uint64
}
