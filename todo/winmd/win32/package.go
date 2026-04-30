package win32

import (
	"sync/atomic"

	"golang.org/x/sys/windows"
)

var (
	libImagehlp = windows.NewLazySystemDLL("Imagehlp.dll")
	libKernel32 = windows.NewLazySystemDLL("kernel32.dll")
)

func lazyAddr(pAddr *uintptr, lib *windows.LazyDLL, procName string) uintptr {
	addr := atomic.LoadUintptr(pAddr)
	if addr == 0 {
		addr = lib.NewProc(procName).Addr()
		atomic.StoreUintptr(pAddr, addr)
	}
	return addr
}
