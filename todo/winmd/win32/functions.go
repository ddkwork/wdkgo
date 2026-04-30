package win32

import (
	"syscall"
	"unsafe"
)

var (
	pMapAndLoad   uintptr
	pUnMapAndLoad uintptr
)

func MapAndLoad(ImageName PSTR, DllPath PSTR, LoadedImage *LOADED_IMAGE, DotDll BOOL, ReadOnly BOOL) (BOOL, WIN32_ERROR) {
	addr := lazyAddr(&pMapAndLoad, libImagehlp, "MapAndLoad")
	ret, _, err := syscall.SyscallN(addr, uintptr(unsafe.Pointer(ImageName)), uintptr(unsafe.Pointer(DllPath)), uintptr(unsafe.Pointer(LoadedImage)), uintptr(DotDll), uintptr(ReadOnly))
	return BOOL(ret), WIN32_ERROR(err)
}

func UnMapAndLoad(LoadedImage *LOADED_IMAGE) (BOOL, WIN32_ERROR) {
	addr := lazyAddr(&pUnMapAndLoad, libImagehlp, "UnMapAndLoad")
	ret, _, err := syscall.SyscallN(addr, uintptr(unsafe.Pointer(LoadedImage)))
	return BOOL(ret), WIN32_ERROR(err)
}

func StrToPstr(str string) PSTR {
	bts := []byte(str)
	return (PSTR)(&bts[0])
}
