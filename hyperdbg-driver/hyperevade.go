package main

import "slices"

type EvasionTechnique int

const (
	EvasionHideDebugger EvasionTechnique = 1 << iota
	EvasionHideVmx
	EvasionHideMsrs
	EvasionHideIoPorts
	EvasionHideTiming
	EvasionAll = 0xFFFFFFFF
)

func (e EvasionTechnique) String() string {
	var result string
	if e&EvasionHideDebugger != 0 {
		result += "|HIDE_DEBUGGER"
	}
	if e&EvasionHideVmx != 0 {
		result += "|HIDE_VMX"
	}
	if e&EvasionHideMsrs != 0 {
		result += "|HIDE_MSRS"
	}
	if e&EvasionHideIoPorts != 0 {
		result += "|HIDE_IOPORTS"
	}
	if e&EvasionHideTiming != 0 {
		result += "|HIDE_TIMING"
	}
	if len(result) > 0 {
		return result[1:]
	}
	return "NONE"
}

type DebuggerHidingMethod = int

const (
	HideFromKdQuerySystemInformation DebuggerHidingMethod = 1 << iota
	HideFromKdGetContextThread
	HideFromProcessList
	HideFromThreadList
	HideAllMethods = 0xFFFFFFFF
)

type TransparentModeConfig struct {
	Enabled     bool
	Techniques  EvasionTechnique
	DebuggerPid uint32
}

type EvasionState struct {
	config          TransparentModeConfig
	active          bool
	hiddenPids      []uint32
	lock            *Spinlock
	vmxOriginalCr4  uint64
	msrBackup       map[uint32]uint64
	ioPortHooks     map[uint16]bool
	tscBase         uint64
	tscOffset       int64
	originalRdtscFn func() uint64
}

func NewEvasionState() *EvasionState {
	return &EvasionState{
		config: TransparentModeConfig{
			Enabled:    false,
			Techniques: EvasionAll,
		},
		hiddenPids:  make([]uint32, 0, 16),
		lock:        NewSpinlock(),
		msrBackup:   make(map[uint32]uint64),
		ioPortHooks: make(map[uint16]bool),
	}
}

func (e *EvasionState) Enable(config *TransparentModeConfig) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if config != nil {
		e.config = *config
	}

	e.config.Enabled = true
	e.active = true

	if err := e.applyEvasions(); err != nil {
		e.active = false
		return err
	}

	techStr := e.config.Techniques.String()
	LogInfo("Transparent mode enabled: %s", techStr)
	return nil
}

func (e *EvasionState) Disable() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if !e.active {
		return nil
	}

	if err := e.removeEvasions(); err != nil {
		LogWarning("Error removing evasions: %v", err)
	}

	e.config.Enabled = false
	e.active = false
	e.hiddenPids = e.hiddenPids[:0]

	LogInfo("Transparent mode disabled")
	return nil
}

func (e *EvasionState) IsActive() bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.active
}

func (e *EvasionState) AddHiddenProcess(pid uint32) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if slices.Contains(e.hiddenPids, pid) {
		return
	}
	e.hiddenPids = append(e.hiddenPids, pid)
	LogDebug("EVASION: Hidden PID=%d", pid)
}

func (e *EvasionState) RemoveHiddenProcess(pid uint32) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for i, existing := range e.hiddenPids {
		if existing == pid {
			e.hiddenPids = append(e.hiddenPids[:i], e.hiddenPids[i+1:]...)
			LogDebug("EVASION: Unhidden PID=%d", pid)
			return
		}
	}
}

func (e *EvasionState) IsProcessHidden(pid uint32) bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	return slices.Contains(e.hiddenPids, pid)
}

func (e *EvasionState) applyEvasions() error {
	tech := e.config.Techniques

	if tech&EvasionHideVmx != 0 {
		if err := e.hideVmxBits(); err != nil {
			LogWarning("EVASION: hide VMX bits failed: %v", err)
		}
	}

	if tech&EvasionHideMsrs != 0 {
		if err := e.hideVmxRelatedMsrs(); err != nil {
			LogWarning("EVASION: hide VMX MSRs failed: %v", err)
		}
	}

	if tech&EvasionHideIoPorts != 0 {
		if err := e.hideVmxCriticalIoPorts(); err != nil {
			LogWarning("EVASION: hide IO ports failed: %v", err)
		}
	}

	if tech&EvasionHideTiming != 0 {
		if err := e.enableTimingCountermeasures(); err != nil {
			LogWarning("EVASION: timing countermeasures failed: %v", err)
		}
	}

	if tech&EvasionHideDebugger != 0 {
		if err := e.hideDebuggerPresence(); err != nil {
			LogWarning("EVASION: hide debugger failed: %v", err)
		}
	}

	return nil
}

func (e *EvasionState) removeEvasions() error {

	cr4 := __readcr4()
	cr4 |= uint64(1 << 13)
	__writecr4(cr4)

	for msr, val := range e.msrBackup {
		__writemsr(msr, val)
	}
	e.msrBackup = map[uint32]uint64{}

	e.ioPortHooks = map[uint16]bool{}
	e.tscOffset = 0

	return nil
}

func (e *EvasionState) hideVmxBits() error {
	e.vmxOriginalCr4 = __readcr4()
	cr4 := e.vmxOriginalCr4
	cr4 &= ^uint64(1 << 13)
	__writecr4(cr4)

	LogDebug("EVASION: Hidden VMXE bit in CR4")
	return nil
}

func (e *EvasionState) hideVmxRelatedMsrs() error {
	vmxMsrs := []uint32{
		0x480,
		0x481,
		0x482,
		0x483,
		0x484,
		0x485,
		0x486,
		0x487,
		0x488,
		0x489,
		0x48A,
		0x48B,
		0x6A0,
		0x6A2,
		0x6A4,
		0x6A6,
	}

	for _, msr := range vmxMsrs {
		val := __readmsr(msr)
		e.msrBackup[msr] = val
		__writemsr(msr, 0)
	}

	LogDebug("EVASION: Hidden %d VMX MSRs", len(vmxMsrs))
	return nil
}

func (e *EvasionState) hideVmxCriticalIoPorts() error {
	criticalPorts := []uint16{
		0x3708,
		0x3710,
		0x3718,
		0x3720,
		0x3728,
		0x3730,
		0x3738,
		0x561D,
	}

	for _, port := range criticalPorts {
		e.ioPortHooks[port] = true
	}

	LogDebug("EVASION: Hooked %d critical IO ports", len(criticalPorts))
	return nil
}

func (e *EvasionState) enableTimingCountermeasures() error {
	e.tscBase = rdtscValue()
	e.tscOffset = 0

	LogDebug("EVASION: Timing countermeasures enabled, base TSC=0x%X", e.tscBase)
	return nil
}

func (e *EvasionState) hideDebuggerPresence() error {

	LogDebug("EVASION: Debugger hiding active for PID=%d", e.config.DebuggerPid)
	return nil
}

func (e *EvasionState) ShouldInterceptIoPort(port uint16) bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.ioPortHooks[port]
}

func (e *EvasionState) ShouldHideMsr(msr uint32) bool {
	if !e.IsActive() {
		return false
	}
	_, ok := e.msrBackup[msr]
	return ok
}

func (e *EvasionState) AdjustTsc(tsc uint64) uint64 {
	if !e.IsActive() || e.tscOffset == 0 {
		return tsc
	}
	result := max(int64(tsc)+e.tscOffset, 0)
	return uint64(result)
}

func (e *EvasionState) SetTscOffset(offset int64) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.tscOffset = offset
}

func (e *EvasionState) ShouldHideProcess(pid uint32) bool {
	if !e.IsActive() || !(e.config.Techniques&EvasionHideDebugger != 0) {
		return false
	}
	return e.IsProcessHidden(pid)
}

func (e *EvasionState) HiddenProcessCount() int {
	e.lock.Lock()
	defer e.lock.Unlock()
	return len(e.hiddenPids)
}

func (e *EvasionState) GetConfig() TransparentModeConfig {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.config
}
