package main

import (
	"sync/atomic"

	"solod.dev/so/wdk"
)

type Spinlock struct {
	lock uint32
}

func NewSpinlock() *Spinlock {
	return &Spinlock{}
}

func (s *Spinlock) Lock() {
	for {
		if atomic.CompareAndSwapUint32(&s.lock, 0, 1) {
			return
		}
	}
}

func (s *Spinlock) Unlock() {
	atomic.StoreUint32(&s.lock, 0)
}

func (s *Spinlock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&s.lock, 0, 1)
}

func (s *Spinlock) IsLocked() bool {
	return atomic.LoadUint32(&s.lock) != 0
}

type RwLock struct {
	readers uint32
	writer  uint32
}

func NewRwLock() *RwLock {
	return &RwLock{}
}

func (l *RwLock) RLock() {
	for {
		if atomic.LoadUint32(&l.writer) == 0 {
			atomic.AddUint32(&l.readers, 1)
			if atomic.LoadUint32(&l.writer) == 0 {
				return
			}
			atomic.AddUint32(&l.readers, ^uint32(0))
		}
	}
}

func (l *RwLock) RUnlock() {
	atomic.AddUint32(&l.readers, ^uint32(0))
}

func (l *RwLock) Lock() {
	for !atomic.CompareAndSwapUint32(&l.writer, 0, 1) {
	}
	for atomic.LoadUint32(&l.readers) != 0 {
	}
}

func (l *RwLock) Unlock() {
	atomic.StoreUint32(&l.writer, 0)
}

type VmxRootSpinlock struct {
	raw wdk.KSPIN_LOCK
}

func NewVmxRootSpinlock() *VmxRootSpinlock {
	l := &VmxRootSpinlock{}
	wdk.KeInitializeSpinLock(&l.raw)
	return l
}

func (l *VmxRootSpinlock) Lock() {
	wdk.KeAcquireSpinLockRaiseToDpc(&l.raw)
}

func (l *VmxRootSpinlock) Unlock() {
	wdk.KeReleaseSpinLock(&l.raw, 0)
}
