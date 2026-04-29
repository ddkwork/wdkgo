package main

var (
	gHyp        *Hypervisor
	gLog        *Logger
	gDrv        *Driver
	gEvents     *EventDispatcher
	gHooks      *HookManager
	gEvasion    *EvasionState
	gTrace      *TracerState
	gSerial     *KdSerialState
	gPoolMgr    *PoolManager
	gOptimizer  *VmmOptimizer
	gEptCache   *EptCache
	gAllowIoctl bool
)

func InitializeGlobals(numCpus uint32) error {
	var err error

	gLog = NewLogger(MAX_LOG_BUFFERS, MAX_LOG_BUFFERS_PRIO)
	if err = gLog.Initialize(); err != nil {
		return fmtError("logger init failed: %v", err)
	}

	gPoolMgr = NewPoolManager(1024)

	gHyp = NewHypervisor(numCpus)
	if err = gHyp.Initialize(); err != nil {
		return fmtError("hypervisor init failed: %v", err)
	}

	gEvents = newEventDispatcher()
	setupDefaultEventHandlers()

	gHooks = NewHookManager(int(MAX_HIDDEN_BREAKPOINTS))

	gEvasion = NewEvasionState()

	gTrace = NewTracerState()

	gSerial = NewKdSerialState()

	gOptimizer = NewVmmOptimizer(OptBasic)

	if gOptimizer.GetConfig().EptCache {
		gEptCache = NewEptCache(4096)
	}

	gDrv = &Driver{
		device:      nil,
		handleInUse: false,
		initialized: false,
	}

	LogInfo("Global state initialized for %d CPUs", numCpus)
	return nil
}

func CleanupGlobals() {
	if gEptCache != nil {
		gEptCache = nil
	}

	if gOptimizer != nil {
		gOptimizer = nil
	}

	if gSerial != nil && gSerial.IsConnected() {
		gSerial.Uninitialize()
	}
	gSerial = nil

	if gTrace != nil && gTrace.IsActive() {
		gTrace.Disable()
	}
	gTrace = nil

	if gEvasion != nil && gEvasion.IsActive() {
		gEvasion.Disable()
	}
	gEvasion = nil

	if gHooks != nil {
		gHooks.RemoveAll()
	}
	gHooks = nil

	if gEvents != nil {
		gEvents.Reset()
	}
	gEvents = nil

	if gHyp != nil && gHyp.IsInitialized() {
		gHyp.Shutdown()
	}
	gHyp = nil

	if gPoolMgr != nil {
		gPoolMgr.CheckAndPerformDeallocation()
	}
	gPoolMgr = nil

	if gLog != nil {
		gLog.Uninitialize()
	}
	gLog = nil

	if gDrv != nil {
		gDrv.device = nil
		gDrv.initialized = false
	}
	gDrv = nil
}

func setupDefaultEventHandlers() {
	cpuidHandler := &CpuidHandler{}
	ioHandler := NewIoHandler()
	crHandler := NewCrAccessHandler()
	tscHandler := NewTscHandler()

	gEvents.On(EventCpuid, cpuidHandler)
	gEvents.On(EventIo, ioHandler)
	gEvents.On(EventMovCr, crHandler)
	gEvents.On(EventTsc, tscHandler)

	gEvents.Enable(EventCpuid)
	gEvents.Enable(EventIo)
	gEvents.Enable(EventMovCr)
	gEvents.Enable(EventTsc)
}

func handleExceptionDispatch(_ *VCPU, _ *ExceptionHandler) bool { return true }

func GetHypervisor() *Hypervisor           { return gHyp }
func GetLogger() *Logger                   { return gLog }
func GetDriver() *Driver                   { return gDrv }
func GetEventDispatcher() *EventDispatcher { return gEvents }
func GetHookManager() *HookManager         { return gHooks }
func GetEvasionState() *EvasionState       { return gEvasion }
func GetTracerState() *TracerState         { return gTrace }
func GetSerialState() *KdSerialState       { return gSerial }
func GetPoolManager() *PoolManager         { return gPoolMgr }
func GetOptimizer() *VmmOptimizer          { return gOptimizer }
func GetEptCache() *EptCache               { return gEptCache }
