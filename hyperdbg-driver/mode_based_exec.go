package main

type ModeBasedExecHookState struct {
	Enabled     bool
	CsSelectors []ModeBasedExecEntry
}

type ModeBasedExecEntry struct {
	Selector uint16
	Hooked   bool
}

var gModeBasedExec = &ModeBasedExecHookState{
	CsSelectors: make([]ModeBasedExecEntry, 0),
}

func modeBasedExecHookEnable(v *VCPU) {
	gModeBasedExec.Enabled = true

	secondaryControls := uint64(0)
	vmRead64(VmcsCtrlSecondaryProcBasedVmExec, &secondaryControls)
	secondaryControls |= (1 << 2)
	vmWrite64(VmcsCtrlSecondaryProcBasedVmExec, secondaryControls)

	LogDebug("Mode-based execution hook enabled on core %d", v.CoreId)
}

func modeBasedExecHookDisable(v *VCPU) {
	gModeBasedExec.Enabled = false

	secondaryControls := uint64(0)
	vmRead64(VmcsCtrlSecondaryProcBasedVmExec, &secondaryControls)
	secondaryControls &^= (1 << 2)
	vmWrite64(VmcsCtrlSecondaryProcBasedVmExec, secondaryControls)

	LogDebug("Mode-based execution hook disabled on core %d", v.CoreId)
}

func modeBasedExecHookAddSelector(selector uint16) int {
	entry := ModeBasedExecEntry{
		Selector: selector,
		Hooked:   true,
	}

	gModeBasedExec.CsSelectors = append(gModeBasedExec.CsSelectors, entry)

	return len(gModeBasedExec.CsSelectors) - 1
}

func modeBasedExecHookRemoveSelector(index int) {
	if index >= 0 && index < len(gModeBasedExec.CsSelectors) {
		gModeBasedExec.CsSelectors[index].Hooked = false
	}
}

func modeBasedExecHandleEptViolation(v *VCPU, qualification uint64) bool {
	if !gModeBasedExec.Enabled {
		return false
	}

	var guestCs uint64
	vmRead64(VmcsGuestCsSelector, &guestCs)
	csSelector := uint16(guestCs)

	for i := range gModeBasedExec.CsSelectors {
		entry := &gModeBasedExec.CsSelectors[i]
		if entry.Hooked && csSelector == entry.Selector {
			h := getHypervisor()
			if h != nil {
				h.suppressAdvance(v)
			}

			dispatchEventModeBasedExec(v, csSelector, qualification)

			return true
		}
	}

	return false
}

func dispatchEventModeBasedExec(v *VCPU, selector uint16, qualification uint64) {}
