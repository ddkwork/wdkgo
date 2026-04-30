package win32

import "syscall"

type PROPERTYKEY struct {
	Fmtid syscall.GUID
	Pid   uint32
}

type HSTRING *uint16

type IInspectable struct {
	lpVtbl *IInspectableVtbl
}

type IInspectableVtbl struct {
	QueryInterface      uintptr
	AddRef              uintptr
	Release             uintptr
	GetIids             uintptr
	GetRuntimeClassName uintptr
	GetTrustLevel       uintptr
}
