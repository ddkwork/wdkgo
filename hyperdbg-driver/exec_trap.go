package main

type ExecTrapState struct {
	Enabled    bool
	TrapList   []ExecTrapEntry
	MtfPending bool
}

type ExecTrapEntry struct {
	Address     uint64
	Cr3         CR3_TYPE
	OriginalRip uint64
	Active      bool
}

var gExecTrap = &ExecTrapState{
	TrapList: make([]ExecTrapEntry, 0),
}

func execTrapEnable(v *VCPU) {
	gExecTrap.Enabled = true
	setMonitorTrapFlag(true)
	LogDebug("Execution trap enabled on core %d", v.CoreId)
}

func execTrapDisable(v *VCPU) {
	gExecTrap.Enabled = false
	setMonitorTrapFlag(false)
	LogDebug("Execution trap disabled on core %d", v.CoreId)
}

func execTrapAddEntry(addr uint64, cr3 CR3_TYPE) int {
	entry := ExecTrapEntry{
		Address: addr,
		Cr3:     cr3,
		Active:  true,
	}

	gExecTrap.TrapList = append(gExecTrap.TrapList, entry)

	return len(gExecTrap.TrapList) - 1
}

func execTrapRemoveEntry(index int) {
	if index >= 0 && index < len(gExecTrap.TrapList) {
		gExecTrap.TrapList[index].Active = false
	}
}

func execTrapHandleMtf(v *VCPU) bool {
	if !gExecTrap.Enabled {
		return false
	}

	var currentRip uint64
	vmRead64(VmcsGuestRip, &currentRip)

	for i := range gExecTrap.TrapList {
		entry := &gExecTrap.TrapList[i]
		if entry.Active && currentRip == entry.Address {
			h := getHypervisor()
			if h != nil {
				h.suppressAdvance(v)
			}

			dispatchEventExecTrap(v, entry.Address, entry.Cr3)

			return true
		}
	}

	return false
}

func execTrapCheckAndHandle(v *VCPU) bool {
	return execTrapHandleMtf(v)
}

func dispatchEventExecTrap(v *VCPU, addr uint64, cr3 CR3_TYPE) {}
