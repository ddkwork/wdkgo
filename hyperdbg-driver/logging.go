package main

import (
	"fmt"
	"unsafe"

	"solod.dev/so/wdk"
)

type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelWarning
	LogLevelError
	LogLevelDebug
	LogLevelTrace
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	default:
		return "????"
	}
}

type LogMessage struct {
	Level   LogLevel
	Time    uint64
	CpuId   uint32
	Process [15]uint8
	Message [256]uint8
	Length  uint32
}

type LogBuffer struct {
	Messages      []LogMessage
	ReadIndex     uint32
	WriteIndex    uint32
	Count         uint32
	Capacity      uint32
	Lock          *Spinlock
	OverflowCount uint64
}

func NewLogBuffer(capacity uint32) *LogBuffer {
	return &LogBuffer{
		Messages: make([]LogMessage, capacity),
		Capacity: capacity,
		Lock:     NewSpinlock(),
	}
}

func (b *LogBuffer) Push(msg *LogMessage) bool {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	if b.Count >= b.Capacity {
		b.OverflowCount++
		return false
	}

	idx := (b.WriteIndex + b.Count) % b.Capacity
	b.Messages[idx] = *msg
	b.Count++
	return true
}

func (b *LogBuffer) Pop() (*LogMessage, bool) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	if b.Count == 0 {
		return nil, false
	}

	msg := &b.Messages[b.ReadIndex]
	b.ReadIndex = (b.ReadIndex + 1) % b.Capacity
	b.Count--
	return msg, true
}

func (b *LogBuffer) Peek() (*LogMessage, bool) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	if b.Count == 0 {
		return nil, false
	}
	return &b.Messages[b.ReadIndex], true
}

func (b *LogBuffer) Clear() {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	b.ReadIndex = 0
	b.WriteIndex = 0
	b.Count = 0
	b.OverflowCount = 0
}

func (b *LogBuffer) IsEmpty() bool {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	return b.Count == 0
}

func (b *LogBuffer) IsFull() bool {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	return b.Count >= b.Capacity
}

func (b *LogBuffer) Available() uint32 {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	return b.Capacity - b.Count
}

func (b *LogBuffer) GetCount() uint32 {
	b.Lock.Lock()
	defer b.Lock.Unlock()
	return b.Count
}

type Logger struct {
	buffer         *LogBuffer
	priorityBuffer *LogBuffer
	enabled        bool
	minLevel       LogLevel
	lock           *Spinlock
}

func NewLogger(bufferSize, prioSize uint32) *Logger {
	l := &Logger{
		buffer:         NewLogBuffer(bufferSize),
		priorityBuffer: NewLogBuffer(prioSize),
		enabled:        true,
		minLevel:       LogLevelInfo,
		lock:           NewSpinlock(),
	}
	return l
}

func (l *Logger) Initialize() error {
	if l.buffer != nil && l.priorityBuffer != nil {
		l.enabled = true
		return nil
	}
	return fmtError("logger initialization failed")
}

func (l *Logger) Uninitialize() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.buffer != nil {
		l.buffer.Clear()
	}
	if l.priorityBuffer != nil {
		l.priorityBuffer.Clear()
	}
	l.enabled = false
}

func (l *Logger) SetMinLevel(level LogLevel) {
	l.minLevel = level
}

func (l *Logger) Enable() {
	l.enabled = true
}

func (l *Logger) Disable() {
	l.enabled = false
}

func (l *Logger) Log(level LogLevel, format string, args ...any) {
	if !l.enabled || level < l.minLevel {
		return
	}

	var msg LogMessage
	msg.Level = level
	msg.Time = rdtscValue()
	var cpuNum uint32
	wdk.KeGetCurrentProcessorNumberEx(&cpuNum)
	msg.CpuId = cpuNum

	procName := CurrentProcessName()
	copy(msg.Process[:], procName)

	formatted := fmt.Sprintf(format, args...)
	length := min(len(formatted), len(msg.Message))
	copy(msg.Message[:length], formatted)
	msg.Length = uint32(length)

	if level <= LogLevelWarning {
		l.priorityBuffer.Push(&msg)
	} else {
		l.buffer.Push(&msg)
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.Log(LogLevelInfo, format, args...)
}

func (l *Logger) Warning(format string, args ...any) {
	l.Log(LogLevelWarning, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.Log(LogLevelError, format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	l.Log(LogLevelDebug, format, args...)
}

func (l *Logger) Trace(format string, args ...any) {
	l.Log(LogLevelTrace, format, args...)
}

func (l *Logger) FlushToUser(buffer uintptr, size uintptr) uint32 {
	var count uint32
	remaining := size

	for remaining > 0 {
		msg, ok := l.buffer.Pop()
		if !ok {
			break
		}

		msgSize := unsafe.Sizeof(LogMessage{})
		if uintptr(msgSize) > remaining {
			break
		}

		dst := (*LogMessage)(unsafe.Pointer(uintptr(buffer) + uintptr(size-remaining)))
		*dst = *msg
		remaining -= uintptr(msgSize)
		count++
	}

	return count
}

func (l *Logger) FlushPriorityToUser(buffer uintptr, size uintptr) uint32 {
	var count uint32
	remaining := size

	for remaining > 0 {
		msg, ok := l.priorityBuffer.Pop()
		if !ok {
			break
		}

		msgSize := unsafe.Sizeof(LogMessage{})
		if uintptr(msgSize) > remaining {
			break
		}

		dst := (*LogMessage)(unsafe.Pointer(uintptr(buffer) + uintptr(size-remaining)))
		*dst = *msg
		remaining -= uintptr(msgSize)
		count++
	}

	return count
}

func (l *Logger) GetBufferSize() uint64 {
	return uint64(l.buffer.Capacity)
}

func (l *Logger) GetCurrentBufferSize() uint64 {
	return uint64(l.buffer.GetCount())
}

func (l *Logger) GetOverflowCount() uint64 {
	return l.buffer.OverflowCount
}

func (l *Logger) GetBufferBase() uint64 {
	if l.buffer == nil || len(l.buffer.Messages) == 0 {
		return 0
	}
	return uint64(uintptr(unsafe.Pointer(&l.buffer.Messages[0])))
}

func (l *Logger) GetPriorityBufferBase() uint64 {
	if l.priorityBuffer == nil || len(l.priorityBuffer.Messages) == 0 {
		return 0
	}
	return uint64(uintptr(unsafe.Pointer(&l.priorityBuffer.Messages[0])))
}

func rdtscValue() uint64 {
	// Use nativeRdtsc from events.go
	return nativeRdtsc()
}

func LogInfo(format string, args ...any) {
	if gLog != nil {
		gLog.Info(format, args...)
	}
}

func LogError(format string, args ...any) {
	if gLog != nil {
		gLog.Error(format, args...)
	}
}

func LogWarning(format string, args ...any) {
	if gLog != nil {
		gLog.Warning(format, args...)
	}
}

func LogDebug(format string, args ...any) {
	if gLog != nil {
		gLog.Debug(format, args...)
	}
}

func LogTrace(format string, args ...any) {
	if gLog != nil {
		gLog.Trace(format, args...)
	}
}

func LogCallbackSendBuffer(operation uint32, data string, length uint32, isImmediate bool) {
	if gLog != nil {
		gLog.Log(LogLevelInfo, "Send buffer: op=%d, len=%d", operation, length)
	}
}

func LogRegisterIrpBasedNotification(irp unsafe.Pointer) {
	gAllowIoctl = true
	LogInfo("Registered IRP-based notification")
}

func LogRegisterEventBasedNotification(irp unsafe.Pointer) bool {
	gAllowIoctl = true
	LogInfo("Registered event-based notification")
	return true
}
