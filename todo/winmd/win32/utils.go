package win32

import (
	"fmt"
	"log"
	"syscall"
)

func FAILED(hr HRESULT) bool {
	return hr < 0
}

func ASSERT_SUCCEEDED(hr HRESULT) {
	if FAILED(hr) {
		log.Fatal("HRESULT failed:", hr)
	}
}

func GuidToStr(guid *syscall.GUID) (string, error) {
	return fmt.Sprintf("{%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X}",
		guid.Data1, guid.Data2, guid.Data3,
		guid.Data4[0], guid.Data4[1], guid.Data4[2], guid.Data4[3],
		guid.Data4[4], guid.Data4[5], guid.Data4[6], guid.Data4[7]), nil
}
