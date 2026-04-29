package main

import (
	"unsafe"

	"solod.dev/so/wdk"
)

func AllocatePool(size uintptr, tag uint32) uintptr {
	return uintptr(wdk.ExAllocatePool2(uint32(POOL_FLAG_NON_PAGED), size, tag))
}

func AllocateZeroedPool(size uintptr, tag uint32) uintptr {
	ptr := AllocatePool(size, tag)
	if ptr != 0 {
		wdk.RtlZeroMemory(ptr, size)
	}
	return ptr
}

func FreePool(ptr uintptr, tag uint32) {
	if ptr != 0 {
		wdk.ExFreePoolWithTag(ptr, tag)
	}
}

func CurrentProcessName() string {
	proc := wdk.PsGetCurrentProcess()
	name := wdk.PsGetProcessImageFileName(proc)
	return cstringToString(name)
}

func cstringToString(s *int8) string {
	if s == nil {
		return ""
	}
	var result []byte
	for i := 0; ; i++ {
		b := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(i)))
		if b == 0 {
			break
		}
		result = append(result, b)
	}
	return string(result)
}
