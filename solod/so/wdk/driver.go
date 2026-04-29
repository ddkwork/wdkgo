//so:include <ntifs.h>

package wdk

import "unsafe"

//so:extern
type HANDLE = uintptr
//so:extern
type PVOID = unsafe.Pointer
//so:extern
type ULONG = uint32
//so:extern
type LONG = int32
//so:extern
type USHORT = uint16
//so:extern
type SHORT = int16
//so:extern
type UCHAR = uint8
//so:extern
type CHAR = int8
//so:extern
type BOOLEAN bool
//so:extern
type NTSTATUS = int32
//so:extern
type SIZE_T = uint64
//so:extern
type ULONG_PTR = uint64
//so:extern
type PVOID64 = uint64
//so:extern
type ACCESS_MASK = uint32
//so:extern
type KSPIN_LOCK = uintptr
//so:extern
type KIRQL = uint8
//so:extern
type LARGE_INTEGER struct { QuadPart int64 }
//so:extern
type ULARGE_INTEGER struct { QuadPart uint64 }
//so:extern
type PHYSICAL_ADDRESS = LARGE_INTEGER
//so:extern
type PCI_SLOT_NUMBER = uint32
//so:extern
type HRESULT = int32
//so:extern
type UINT64 = uint64
//so:extern
type UINT32 = uint32
//so:extern
type INT32 = int32
//so:extern
type UINT16 = uint16
//so:extern
type INT16 = int16
//so:extern
type UINT8 = uint8
//so:extern
type INT8 = int8

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type DEVICE_OBJECT struct {
	Type int16
	Size int16
	ReferenceCount int32
	DriverObject *DRIVER_OBJECT
	NextDevice *DEVICE_OBJECT
	AttachedDevice *DEVICE_OBJECT
	CurrentIrp *IRP
	Timer uintptr
	Flags uint32
	Characteristics uint32
	Vpb *VPB
	DeviceExtension uintptr
	DeviceType uint32
	StackSize int8
	Queue [36]byte
	AlignmentRequirement uint32
	DeviceQueue [32]byte
	Dpc [48]byte
	ActiveThreadCount uint32
	SecurityDescriptor uintptr
	DeviceLock [24]byte
	SectorSize uint16
	Spare1 uint16
	DeviceObjectExtension *DEVICE_OBJECT_EXTENSION
}

type DRIVER_DISPATCH func(DeviceObject *DEVICE_OBJECT, Irp *IRP) NTSTATUS

type DRIVER_UNLOAD func(DriverObject *DRIVER_OBJECT)

type DRIVER_OBJECT struct {
	Type int16
	Size int16
	DeviceObject *DEVICE_OBJECT
	Flags uint32
	DriverStart uintptr
	DriverSize uint32
	DriverSection uintptr
	DriverExtension *DRIVER_EXTENSION
	DriverName UNICODE_STRING
	HardwareDatabase *UNICODE_STRING
	FastIoDispatch *FAST_IO_DISPATCH
	DriverInit uintptr
	DriverStartIo uintptr
	DriverUnload DRIVER_UNLOAD
	MajorFunction [28]DRIVER_DISPATCH
}

type IRP struct {
	Type int16
	Size uint16
	MdlAddress *MDL
	Flags uint32
	AssociatedIrp IRP_ASSOCIATED_IRP
	ThreadListEntry LIST_ENTRY
	IoStatus IO_STATUS_BLOCK
	RequestorMode int8
	PendingReturned byte
	StackCount int8
	CurrentLocation int8
	Cancel byte
	CancelIrql uint8
	ApcEnvironment int8
	AllocationFlags uint8
	UserIosb *IO_STATUS_BLOCK
	UserEvent uintptr
	Overlay IRP_OVERLAY
	CancelRoutine uintptr
	UserBuffer uintptr
	Tail IRP_TAIL
}

type IRP_ASSOCIATED_IRP struct {
	SystemBuffer uintptr
}

type IRP_OVERLAY struct {
	Data [24]byte
}

type IRP_TAIL struct {
	Data [40]byte
}

type IO_STACK_LOCATION struct {
	MajorFunction byte
	MinorFunction byte
	Flags byte
	Control byte
	Parameters IO_STACK_LOCATION_PARAMETERS
	DeviceObject *DEVICE_OBJECT
	FileObject *FILE_OBJECT
	CompletionRoutine uintptr
	Context uintptr
}

type IO_STACK_LOCATION_PARAMETERS struct {
	DeviceIoControl IO_STACK_LOCATION_DEVICE_IO_CONTROL
}

type IO_STACK_LOCATION_DEVICE_IO_CONTROL struct {
	IoControlCode uint32
	InputBufferLength uint32
	OutputBufferLength uint32
}

type ETHREAD struct { _ [0]byte }
type EPROCESS struct { _ [0]byte }
type KAPC struct { _ [0]byte }
type KDPC struct { _ [0]byte }
type KEVENT struct { _ [0]byte }
type KSEMAPHORE struct { _ [0]byte }
type KMUTEX struct { _ [0]byte }
type FAST_MUTEX struct { _ [0]byte }
type ERESOURCE struct { _ [0]byte }
type MDL struct { _ [0]byte }
type FILE_OBJECT struct { _ [0]byte }
type VPB struct { _ [0]byte }
type DRIVER_EXTENSION struct { _ [0]byte }
type DEVICE_OBJECT_EXTENSION struct { _ [0]byte }
type ACCESS_STATE struct { _ [0]byte }
type SECURITY_DESCRIPTOR struct { _ [0]byte }
type ACL struct { _ [0]byte }
type SID struct { _ [0]byte }
type TOKEN struct { _ [0]byte }
type KINTERRUPT struct { _ [0]byte }
type KTIMER struct { _ [0]byte }
type KPROCESS struct { _ [0]byte }
type KTHREAD struct { _ [0]byte }
type IO_WORKITEM struct { _ [0]byte }
type IO_TIMER struct { _ [0]byte }
type FILE_LOCK struct { _ [0]byte }
type SECTION struct { _ [0]byte }
type LOOKASIDE_LIST_EX struct { _ [0]byte }
type CM_RESOURCE_LIST struct { _ [0]byte }
type DMA_ADAPTER struct { _ [0]byte }
type WAIT_CONTEXT_BLOCK struct { _ [0]byte }
type EX_RUNDOWN_REF struct { _ [0]byte }
type EX_PUSH_LOCK struct { _ [0]byte }
type FAST_IO_DISPATCH struct { _ [0]byte }
type CACHE_MANAGER_CALLBACKS struct { _ [0]byte }
type GENERIC_MAPPING struct { _ [0]byte }
type SECURITY_CLIENT_CONTEXT struct { _ [0]byte }
type INTERFACE struct { _ [0]byte }
type BUS_INTERFACE_STANDARD struct { _ [0]byte }
type ADAPTER_OBJECT struct { _ [0]byte }
type DMA_OPERATIONS struct { _ [0]byte }
type SCATTER_GATHER_LIST struct { _ [0]byte }
type SHARED_CACHE_MAP struct { _ [0]byte }
type BCB struct { _ [0]byte }
type MMPTE struct { _ [0]byte }
type MMPFN struct { _ [0]byte }
type WHEA_ERROR_RECORD struct { _ [0]byte }

type LIST_ENTRY struct {
	Flink uintptr
	Blink uintptr
}

type SINGLE_LIST_ENTRY struct {
	Next uintptr
}

type UNICODE_STRING struct {
	Length uint16
	MaximumLength uint16
	Buffer uintptr
}

type ANSI_STRING struct {
	Length uint16
	MaximumLength uint16
	Buffer uintptr
}

type STRING = ANSI_STRING

type IO_STATUS_BLOCK struct {
	Status int32
	Information ULONG_PTR
}

type OBJECT_ATTRIBUTES struct {
	Length uint32
	RootDirectory uintptr
	ObjectName uintptr
	Attributes uint32
	SecurityDescriptor uintptr
	SecurityQualityOfService uintptr
}

type CLIENT_ID struct {
	UniqueProcess uintptr
	UniqueThread uintptr
}

type LUID struct {
	LowPart uint32
	HighPart int32
}

type PRIVILEGE_SET struct {
	PrivilegeCount uint32
	Control uint32
	Privilege [1]LUID_AND_ATTRIBUTES
}

type LUID_AND_ATTRIBUTES struct {
	Luid LUID
	Attributes uint32
}

const (
	STATUS_SUCCESS NTSTATUS = 0x00000000
	STATUS_UNSUCCESSFUL NTSTATUS = -1073741823
	STATUS_ACCESS_DENIED NTSTATUS = -1073741790
	STATUS_INVALID_PARAMETER NTSTATUS = -1073741811
)

const (
	IRP_MJ_CREATE = 0x00
	IRP_MJ_CLOSE = 0x02
	IRP_MJ_READ = 0x03
	IRP_MJ_WRITE = 0x04
	IRP_MJ_DEVICE_CONTROL = 0x0E
	IRP_MJ_INTERNAL_DEVICE_CONTROL = 0x0F
	IRP_MJ_MAXIMUM_FUNCTION = 0x1B
)

const IRP_DEVICE_CONTROL = IRP_MJ_DEVICE_CONTROL

const (
	FILE_DEVICE_UNKNOWN = 0x00000022
	FILE_DEVICE_SECURE_OPEN = 0x00000100
)

const (
	IO_NO_INCREMENT = 0
	DO_BUFFERED_IO = 0x00000004
	DO_DIRECT_IO = 0x00000010
)

const (
	DrvRtPoolNxOptIn = 1
)

const SE_DEBUG_PRIVILEGE = 20

func CTL_CODE(DeviceType uint32, Function uint32, Method uint32, Access uint32) uint32 {
	return ((DeviceType) << 16) | ((Access) << 14) | ((Function) << 2) | (Method)
}

func IoGetCurrentIrpStackLocation(Irp *IRP) *IO_STACK_LOCATION {
	return (*IO_STACK_LOCATION)(unsafe.Pointer(uintptr(unsafe.Pointer(Irp)) + unsafe.Sizeof(IRP{}) - unsafe.Sizeof(IO_STACK_LOCATION{})))
}

//so:extern
func IoCreateDevice(DriverObject *DRIVER_OBJECT, DeviceExtensionSize ULONG, DeviceName *UNICODE_STRING, DeviceType ULONG, DeviceCharacteristics ULONG, Exclusive bool, DeviceObject **DEVICE_OBJECT) NTSTATUS { return 0 }

//so:extern
func IoDeleteDevice(DeviceObject *DEVICE_OBJECT) {}

//so:extern
func IoCreateSymbolicLink(SymbolicLinkName *UNICODE_STRING, DeviceName *UNICODE_STRING) NTSTATUS { return 0 }

//so:extern
func IoDeleteSymbolicLink(SymbolicLinkName *UNICODE_STRING) NTSTATUS { return 0 }

//so:extern
func IoCompleteRequest(Irp *IRP, PriorityBoost int32) {}

//so:extern
func RtlInitUnicodeString(DestinationString *UNICODE_STRING, SourceString *uint16) {}

//so:extern
func ExInitializeDriverRuntime(RuntimeFlags ULONG) {}

//so:extern
func SeSinglePrivilegeCheck(PrivilegeValue LUID, PreviousMode int8) bool { return false }

//so:extern
func UNREFERENCED_PARAMETER(Expression any) {}

//so:extern
func ASSERT(Expression bool) {}

//so:extern nodecl
func DbgPrint(Format string, args ...any) int32 { return 0 }

//so:extern nodecl
func ExAllocatePool2(Flags uint32, NumberOfBytes uintptr, Tag uint32) uintptr { return 0 }

//so:extern
func ExFreePoolWithTag(P uintptr, Tag uint32) {}

//so:extern
func RtlZeroMemory(Destination uintptr, Length uintptr) {}

//so:extern
func RtlCopyMemory(Destination uintptr, Source uintptr, Length uintptr) {}

//so:extern
func RtlMoveMemory(Destination uintptr, Source uintptr, Length uintptr) {}

//so:extern
func RtlFillMemory(Destination uintptr, Length uintptr, Pattern uint8) {}

//so:extern
func MmGetVirtualForPhysical(PhysicalAddress uint64) uintptr { return 0 }

//so:extern
func MmIsAddressValid(VirtualAddress uintptr) bool { return false }

//so:extern nodecl
func MmGetPhysicalAddress(VirtualAddress uintptr) PHYSICAL_ADDRESS { return PHYSICAL_ADDRESS{} }

//so:extern nodecl
func MmMapIoSpace(PhysicalAddress PHYSICAL_ADDRESS, NumberOfBytes uintptr, CacheType int32) uintptr { return 0 }

//so:extern
func MmUnmapIoSpace(BaseAddress uintptr, NumberOfBytes uintptr) {}

//so:extern
func KeInitializeSpinLock(SpinLock *KSPIN_LOCK) {}

//so:extern
func KeAcquireSpinLockRaiseToDpc(SpinLock *KSPIN_LOCK) uint8 { return 0 }

//so:extern
func KeReleaseSpinLock(SpinLock *KSPIN_LOCK, OldIrql uint8) {}

//so:extern
func KeGetCurrentIrql() uint8 { return 0 }

//so:extern nodecl
func KeRaiseIrql(NewIrql uint8, OldIrql *uint8) {}

//so:extern
func KeLowerIrql(NewIrql uint8) {}

//so:extern
func KeQueryActiveProcessorCount(Number *uint32) uint32 { return 0 }

//so:extern
func KeGetCurrentProcessorNumberEx(ProcessorNumber *uint32) uint32 { return 0 }

//so:extern
func KeGetCurrentThread() *uintptr { return nil }

//so:extern
func KeStallExecutionProcessor(Microseconds uint32) {}

//so:extern
func PsGetCurrentProcess() *uintptr { return nil }

//so:extern nodecl
func PsGetProcessImageFileName(Process *uintptr) *int8 { return nil }


const (
	FALSE bool = false
	TRUE bool = true
)

