package main

type EventKind = int

const (
	EventCpuid EventKind = iota
	EventRdmsr
	EventWrmsr
	EventIo
	EventMovCr
	EventMovDr
	EventTsc
	EventRdpmc
	EventXsetbv
	EventException
	EventExternalInt
	EventVmcall
)

type EventHandler interface {
	Handle(v *VmxCtx) bool
}

type EventFunc func(v *VmxCtx) bool

func (f EventFunc) Handle(v *VmxCtx) bool { return f(v) }

type EventDispatcher struct {
	handlers map[EventKind]EventHandler
	enabled  map[EventKind]bool
}

func newEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[EventKind]EventHandler),
		enabled:  make(map[EventKind]bool),
	}
}

func (d *EventDispatcher) On(kind EventKind, handler EventHandler) {
	d.handlers[kind] = handler
}

func (d *EventDispatcher) OnFunc(kind EventKind, fn EventFunc) {
	d.handlers[kind] = fn
}

func (d *EventDispatcher) Enable(kind EventKind) {
	d.enabled[kind] = true
}

func (d *EventDispatcher) Disable(kind EventKind) {
	delete(d.enabled, kind)
}

func (d *EventDispatcher) IsEnabled(kind EventKind) bool {
	return d.enabled[kind]
}

func (d *EventDispatcher) Dispatch(kind EventKind, v *VmxCtx) bool {
	if !d.enabled[kind] {
		return false
	}
	if h, ok := d.handlers[kind]; ok {
		return h.Handle(v)
	}
	return false
}

func (d *EventDispatcher) Reset() {
	d.handlers = make(map[EventKind]EventHandler)
	d.enabled = make(map[EventKind]bool)
}

type CpuidHandler struct{}

func (h *CpuidHandler) Handle(v *VmxCtx) bool {
	cpuId := v.Regs.Rax

	switch cpuId {
	case 0:
		v.Regs.Rbx = 0x756E6547
		v.Regs.Rdx = 0x6C65746E
		v.Regs.Rcx = 0x44534147
	case 1:
		v.Regs.Rcx &= ^uint64(1<<31 | 1<<30 | 1<<29 | 1<<28 | 1<<27)
	case 0xA:
		v.Regs.Rax = 0
		v.Regs.Rbx = 0
		v.Regs.Rcx = 0
		v.Regs.Rdx = 0
	}

	return true
}

type MsrHandler struct {
	onRead  func(v *VmxCtx, msr uint32, value uint64) bool
	onWrite func(v *VmxCtx, msr uint32, value uint64) bool
}

func NewMsrHandler() *MsrHandler {
	return &MsrHandler{
		onRead:  defaultMsrRead,
		onWrite: defaultMsrWrite,
	}
}

func (h *MsrHandler) HandleRead(v *VmxCtx) bool {
	msr := uint32(v.Regs.Rcx)
	value := __readmsr(msr)

	switch msr {
	case 0x174:
		value &= ^uint64(1 << 2)
	}

	v.Regs.Rax = value & 0xFFFFFFFF
	v.Regs.Rdx = value >> 32
	return h.onRead(v, msr, value)
}

func (h *MsrHandler) HandleWrite(v *VmxCtx) bool {
	msr := uint32(v.Regs.Rcx)
	value := (v.Regs.Rdx << 32) | v.Regs.Rax

	switch msr {
	case 0x174:
		value &= ^uint64(1 << 2)
	}

	__writemsr(msr, value)
	return h.onWrite(v, msr, value)
}

func defaultMsrRead(_ *VmxCtx, _ uint32, _ uint64) bool  { return true }
func defaultMsrWrite(_ *VmxCtx, _ uint32, _ uint64) bool { return true }

type IoHandler struct {
	onIn  func(v *VmxCtx, port uint16, size uint32) bool
	onOut func(v *VmxCtx, port uint16, size uint32) bool
}

func NewIoHandler() *IoHandler {
	return &IoHandler{
		onIn:  defaultIoIn,
		onOut: defaultIoOut,
	}
}

func (h *IoHandler) Handle(v *VmxCtx) bool {
	qual := v.ExitQualification
	isOut := qual&0x8 != 0
	size := (qual & 0x7) + 1
	port := uint16(qual >> 16)

	if isOut {
		switch size {
		case 1:
			_ = v.Regs.Rax & 0xFF
		case 2:
			_ = v.Regs.Rax & 0xFFFF
		case 4:
			_ = v.Regs.Rax & 0xFFFFFFFF
		}
		return h.onOut(v, port, size)
	}
	return h.onIn(v, port, size)
}

func defaultIoIn(_ *VmxCtx, _ uint16, _ uint32) bool  { return true }
func defaultIoOut(_ *VmxCtx, _ uint16, _ uint32) bool { return true }

type CrAccessHandler struct {
	onChange func(v *VmxCtx, crNum uint32, read, write bool) bool
}

func NewCrAccessHandler() *CrAccessHandler {
	return &CrAccessHandler{
		onChange: defaultCrChange,
	}
}

func (h *CrAccessHandler) Handle(v *VmxCtx) bool {
	qual := v.ExitQualification
	crNum := (qual >> 8) & 0xF
	accessType := (qual >> 4) & 0x3
	lmswOp := qual&0x1 != 0

	read := accessType == 0 || accessType == 2
	write := accessType == 1 || accessType == 2

	if lmswOp && write {
		crNum = 0
	}

	return h.onChange(v, crNum, read, write)
}

func defaultCrChange(_ *VmxCtx, _ uint32, _, _ bool) bool { return true }

type TscHandler struct {
	offset int64
	rdtsc  func() uint64
}

func NewTscHandler() *TscHandler {
	return &TscHandler{
		rdtsc: nativeRdtsc,
	}
}

func (h *TscHandler) Handle(v *VmxCtx) bool {
	tsc := h.rdtsc()
	if h.offset != 0 {
		tsc += uint64(h.offset)
	}
	v.Regs.Rax = tsc & 0xFFFFFFFF
	v.Regs.Rdx = tsc >> 32
	return true
}

func nativeRdtsc() uint64 {
	return rdtscValue()
}

type ExceptionHandler struct {
	onBp  func(v *VmxCtx) bool
	onGp  func(v *VmxCtx) bool
	onUd  func(v *VmxCtx) bool
	onNmi func(v *VmxCtx) bool
}

func NewExceptionHandler() *ExceptionHandler {
	return &ExceptionHandler{
		onBp:  defaultBpHandler,
		onGp:  defaultGpHandler,
		onUd:  defaultUdHandler,
		onNmi: defaultNmiHandler,
	}
}

func (h *ExceptionHandler) HandleBreakpoint(v *VmxCtx) bool {
	return h.onBp(v)
}

func (h *ExceptionHandler) HandleGeneralProtectionFault(v *VmxCtx) bool {
	return h.onGp(v)
}

func (h *ExceptionHandler) HandleUndefinedOpcode(v *VmxCtx) bool {
	return h.onUd(v)
}

func (h *ExceptionHandler) HandleNmi(v *VmxCtx) bool {
	return h.onNmi(v)
}

func defaultBpHandler(_ *VmxCtx) bool  { return false }
func defaultGpHandler(_ *VmxCtx) bool  { return false }
func defaultUdHandler(_ *VmxCtx) bool  { return false }
func defaultNmiHandler(_ *VmxCtx) bool { return false }
