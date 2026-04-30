package win32

import (
	"syscall"
	"unsafe"
)

var (
	pRoActivateInstance     uintptr
	pRoGetActivationFactory uintptr
)

func RoActivateInstance(activatableClassId HSTRING, instance **IInspectable) HRESULT {
	addr := lazyAddr(&pRoActivateInstance, libKernel32, "RoActivateInstance")
	ret, _, _ := syscall.SyscallN(addr, uintptr(unsafe.Pointer(activatableClassId)), uintptr(unsafe.Pointer(instance)))
	return HRESULT(ret)
}

func RoGetActivationFactory(activatableClassId HSTRING, iid *syscall.GUID, factory unsafe.Pointer) HRESULT {
	addr := lazyAddr(&pRoGetActivationFactory, libKernel32, "RoGetActivationFactory")
	ret, _, _ := syscall.SyscallN(addr, uintptr(unsafe.Pointer(activatableClassId)), uintptr(unsafe.Pointer(iid)), uintptr(factory))
	return HRESULT(ret)
}
