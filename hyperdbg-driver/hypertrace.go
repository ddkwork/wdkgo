package main

import (
	"unsafe"

	"solod.dev/so/wdk"
)

type TraceMode int

const (
	TraceLbr TraceMode = 1 << iota
	TraceBts
	TracePebs
	TraceAll = 0xFFFFFFFF
)

func (t TraceMode) String() string {
	switch t {
	case TraceLbr:
		return "LBR"
	case TraceBts:
		return "BTS"
	case TracePebs:
		return "PEBS"
	default:
		return "ALL"
	}
}

type LbrEntry struct {
	From uint64
	To   uint64
	Info uint64
	Misc uint64
}

type BtsEntry struct {
	LastBranchFrom uint64
	LastBranchTo   uint64
	Flags          uint64
}

type TraceUserConfig struct {
	LbrFilter     uint64
	BtsBufferSize SIZE_T
	CallstackMode bool
}

type LbrConfig struct {
	Enabled       bool
	Filter        uint64
	NumEntries    int
	CallstackMode bool
}

type BtsConfig struct {
	Enabled         bool
	BufferSize      SIZE_T
	BufferPhysical  uint64
	BufferVirtual   uintptr
	InterruptOnFull bool
	BranchType      uint32
	BaseMsrValue    uint64
	MaskMsrValue    uint64
}

type TracerState struct {
	lbrConfig LbrConfig
	btsConfig BtsConfig
	active    bool
	lock      *Spinlock
	lbrBuffer []LbrEntry
	btsBuffer []BtsEntry
}

func NewTracerState() *TracerState {
	numEntries := 32
	t := &TracerState{
		lbrConfig: LbrConfig{
			NumEntries:    numEntries,
			CallstackMode: true,
		},
		btsConfig: BtsConfig{
			BufferSize:      SIZE_T(SIZE_1_GB),
			InterruptOnFull: true,
			BranchType:      0,
		},
		lock: NewSpinlock(),
	}
	t.lbrBuffer = make([]LbrEntry, numEntries)
	maxBtsEntries := int(SIZE_1_GB / unsafe.Sizeof(BtsEntry{}))
	t.btsBuffer = make([]BtsEntry, maxBtsEntries)
	return t
}

func (t *TracerState) Enable(mode TraceMode, config *TraceUserConfig) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if mode&TraceLbr != 0 {
		if err := t.enableLbr(config); err != nil {
			return err
		}
	}

	if mode&TraceBts != 0 {
		if err := t.enableBts(config); err != nil {
			t.disableLbrLocked()
			return err
		}
	}

	t.active = true
	modeStr := mode.String()
	LogInfo("TRACE: Enabled mode=%s", modeStr)
	return nil
}

func (t *TracerState) Disable() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.active {
		return nil
	}

	if t.lbrConfig.Enabled {
		t.disableLbrLocked()
	}

	if t.btsConfig.Enabled {
		t.disableBtsLocked()
	}

	t.active = false
	LogInfo("TRACE: Disabled")
	return nil
}

func (t *TracerState) IsActive() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.active
}

func (t *TracerState) enableLbr(config *TraceUserConfig) error {
	debugctl := __readmsr(IA32_DEBUGCTL_MSR)

	debugctl |= (1 << 0)
	if config != nil && config.CallstackMode {
		t.lbrConfig.CallstackMode = true
		debugctl |= (1 << 2)
	} else if config != nil {
		t.lbrConfig.CallstackMode = false
		debugctl &= ^(uint64(1) << 2)
	}

	if config != nil {
		t.lbrConfig.Filter = config.LbrFilter
	} else {
		t.lbrConfig.Filter = 0
	}

	__writemsr(IA32_DEBUGCTL_MSR, debugctl)
	t.lbrConfig.Enabled = true
	LogDebug("TRACE: LBR enabled (%d entries)", t.lbrConfig.NumEntries)
	return nil
}

func (t *TracerState) disableLbrLocked() {
	debugctl := __readmsr(IA32_DEBUGCTL_MSR)
	debugctl &= ^uint64((1 << 0) | (1 << 2))
	__writemsr(IA32_DEBUGCTL_MSR, debugctl)
	t.lbrConfig.Enabled = false
	LogDebug("TRACE: LBR disabled")
}

func (t *TracerState) enableBts(config *TraceUserConfig) error {
	bufSize := SIZE_T(SIZE_1_MB)
	if config != nil && config.BtsBufferSize > 0 {
		bufSize = config.BtsBufferSize
	}

	bufVa := wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), uintptr(bufSize), POOLTAG)
	if bufVa == 0 {
		return fmtError("BTS buffer allocation failed")
	}
	wdk.RtlZeroMemory(bufVa, uintptr(bufSize))
	bufPa := VirtToPhys(bufVa, CR3_TYPE{Flags: __readcr3()})

	t.btsConfig.BufferSize = bufSize
	t.btsConfig.BufferPhysical = bufPa
	t.btsConfig.BufferVirtual = bufVa

	baseVal := bufPa | 0x1
	maskVal := ^(bufSize - 1)
	maskVal |= 0x1

	__writemsr(IA32_DS_AREA_MSR, baseVal)
	t.btsConfig.BaseMsrValue = baseVal

	__writemsr(IA32_PEBS_ENABLE_MSR, maskVal)
	t.btsConfig.MaskMsrValue = maskVal

	debugctl := __readmsr(IA32_DEBUGCTL_MSR)
	debugctl |= (1 << 7) | (1 << 8) | (1 << 9)
	__writemsr(IA32_DEBUGCTL_MSR, debugctl)

	t.btsConfig.Enabled = true
	LogDebug("TRACE: BTS enabled (buffer=0x%X bytes, PA=0x%X)", bufSize, bufPa)
	return nil
}

func (t *TracerState) disableBtsLocked() {
	debugctl := __readmsr(IA32_DEBUGCTL_MSR)
	debugctl &= ^uint64((1 << 7) | (1 << 8) | (1 << 9))
	__writemsr(IA32_DEBUGCTL_MSR, debugctl)

	__writemsr(IA32_DS_AREA_MSR, 0)
	__writemsr(IA32_PEBS_ENABLE_MSR, 0)

	if t.btsConfig.BufferVirtual != 0 {
		wdk.ExFreePoolWithTag(t.btsConfig.BufferVirtual, POOLTAG)
		t.btsConfig.BufferVirtual = 0
		t.btsConfig.BufferPhysical = 0
	}

	t.btsConfig.Enabled = false
	LogDebug("TRACE: BTS disabled")
}

func (t *TracerState) CaptureLbr() ([]LbrEntry, int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.lbrConfig.Enabled || !t.active {
		return nil, 0
	}

	count := t.readLbrFromMsr(t.lbrBuffer)
	if count <= 0 {
		return nil, 0
	}

	result := make([]LbrEntry, count)
	copy(result, t.lbrBuffer[:count])
	return result, count
}

func (t *TracerState) CaptureBts() ([]BtsEntry, int) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.btsConfig.Enabled || !t.active {
		return nil, 0
	}

	count := t.readBtsFromMemory(t.btsBuffer)
	if count <= 0 {
		return nil, 0
	}

	result := make([]BtsEntry, count)
	copy(result, t.btsBuffer[:count])
	return result, count
}

func (t *TracerState) readLbrFromMsr(buffer []LbrEntry) int {
	lbrTos := __readmsr(IA32_LBR_TOS_MSR)
	tos := int(lbrTos & 0xF)

	if tos >= t.lbrConfig.NumEntries || tos >= len(buffer) {
		return 0
	}

	numPairs := tos
	if t.lbrConfig.CallstackMode {
		for i := range numPairs {
			idx := (tos - 1 - i + 16) % 16
			fromMsr := IA32_LBR_FROM_0_MSR + uint32(idx)
			toMsr := IA32_LBR_TO_0_MSR + uint32(idx)
			infoMsr := IA32_LBR_INFO_0_MSR + uint32(idx)
			miscMsr := IA32_LBR_MISC_0_MSR + uint32(idx)

			buffer[i].From = __readmsr(fromMsr)
			buffer[i].To = __readmsr(toMsr)
			buffer[i].Info = __readmsr(infoMsr)
			buffer[i].Misc = __readmsr(miscMsr)
		}
	} else {
		for i := range numPairs {
			fromMsr := IA32_LBR_FROM_0_MSR + uint32(i)
			toMsr := IA32_LBR_TO_0_MSR + uint32(i)
			infoMsr := IA32_LBR_INFO_0_MSR + uint32(i)
			miscMsr := IA32_LBR_MISC_0_MSR + uint32(i)

			buffer[i].From = __readmsr(fromMsr)
			buffer[i].To = __readmsr(toMsr)
			buffer[i].Info = __readmsr(infoMsr)
			buffer[i].Misc = __readmsr(miscMsr)
		}
	}

	return numPairs
}

func (t *TracerState) readBtsFromMemory(buffer []BtsEntry) int {
	if t.btsConfig.BufferVirtual == 0 || t.btsConfig.BufferSize == 0 {
		return 0
	}

	btsBase := (*BtsEntry)(unsafe.Pointer(t.btsConfig.BufferVirtual))
	maxEntries := int(t.btsConfig.BufferSize / SIZE_T(unsafe.Sizeof(BtsEntry{})))

	indexMsr := __readmsr(IA32_BTS_INDEX_MSR)
	absoluteIndex := int(indexMsr / uint64(unsafe.Sizeof(BtsEntry{})))

	count := min(min(absoluteIndex, maxEntries), len(buffer))

	for i := range count {
		entry := (*BtsEntry)(unsafe.Pointer(
			uintptr(unsafe.Pointer(btsBase)) + uintptr(i)*unsafe.Sizeof(BtsEntry{})))
		buffer[i] = *entry
	}

	return count
}

const (
	IA32_LBR_TOS_MSR     uint32 = 0x01C9
	IA32_LBR_FROM_0_MSR  uint32 = 0x600
	IA32_LBR_TO_0_MSR    uint32 = 0x660
	IA32_LBR_INFO_0_MSR  uint32 = 0xDC0
	IA32_LBR_MISC_0_MSR  uint32 = 0xDD0
	IA32_DS_AREA_MSR     uint32 = 0x600
	IA32_PEBS_ENABLE_MSR uint32 = 0x3F1
	IA32_BTS_INDEX_MSR   uint32 = 0x60C
)
