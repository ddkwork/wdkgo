package win32

import "syscall"

type (
	BOOL      = int32
	BOOLEAN   = uint8
	BSTR      = *uint16
	HANDLE    = uintptr
	HINSTANCE = uintptr
	HRESULT   = int32
	HWND      = uintptr
	LPARAM    = uintptr
	PSTR      = *uint8
	PWSTR     = *uint16
	WPARAM    = uintptr
	HRSRC     = uintptr
)

type WIN32_ERROR uint32

const (
	TRUE  BOOL = 1
	FALSE BOOL = 0
)

func (me WIN32_ERROR) Error() string {
	return syscall.Errno(me).Error()
}
