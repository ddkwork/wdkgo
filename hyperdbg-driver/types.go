package main

import "unsafe"

//so:extern nodecl
type NTSTATUS = int32

//so:extern nodecl
type BOOLEAN = bool

//so:extern nodecl
type SIZE_T = uint64

//so:extern nodecl
type CHAR = int8

//so:extern nodecl
type PVOID = uintptr

//so:extern nodecl
type ULONG = uint32

//so:extern nodecl
type USHORT = uint16

//so:extern nodecl
type UCHAR = uint8

//so:extern nodecl
type ULONG_PTR = uint64

//so:extern
const (
	TRUE  bool = true
	FALSE bool = false
)

//so:extern
const (
	STATUS_SUCCESS                NTSTATUS = 0x00000000
	STATUS_UNSUCCESSFUL           NTSTATUS = -1073741823
	STATUS_NOT_IMPLEMENTED        NTSTATUS = -1073741822
	STATUS_INVALID_PARAMETER      NTSTATUS = -1073741811
	STATUS_BUFFER_TOO_SMALL       NTSTATUS = -1073741790
	STATUS_ACCESS_DENIED          NTSTATUS = -1073741820
	STATUS_INVALID_DEVICE_REQUEST NTSTATUS = -1073741808
	STATUS_DEVICE_NOT_READY       NTSTATUS = -1073741668

	POOL_FLAG_NON_PAGED uint64 = 0x0000000000000040
)

const (
	PAGE_SIZE   uintptr = 4096
	SIZE_1_MB   uintptr = 256 * PAGE_SIZE
	SIZE_2_MB   uintptr = 512 * PAGE_SIZE
	SIZE_1_GB   uintptr = 512 * SIZE_2_MB
	SIZE_512_GB uintptr = 512 * SIZE_1_GB

	POOLTAG uint32 = 0x48444247

	MAX_LOG_BUFFERS      uint32  = 1000
	MAX_LOG_BUFFERS_PRIO uint32  = 50
	LOG_CHUNK_SIZE       uintptr = PAGE_SIZE
	LOG_MESSAGE_SIZE     uintptr = 288

	PENDING_INTERRUPTS_CAPACITY uint32 = 64
	MAX_HIDDEN_BREAKPOINTS      uint32 = 40
	MAX_MTRR_ENTRIES            uint32 = 255

	EPT_PML4_COUNT uint32 = 512
	EPT_PML3_COUNT uint32 = 512
	EPT_PML2_COUNT uint32 = 512
	EPT_PML1_COUNT uint32 = 512

	PAGE_ATTRIB_READ             uint64 = 0x2
	PAGE_ATTRIB_WRITE            uint64 = 0x4
	PAGE_ATTRIB_EXEC             uint64 = 0x8
	PAGE_ATTRIB_EXEC_HIDDEN_HOOK uint64 = 0x10

	HYPERDBG_VMCALL_MAGIC_RAX uint64 = 0x48564653
	HYPERDBG_VMCALL_MAGIC_RCX uint64 = 0x564d43414c4c
	HYPERDBG_VMCALL_MAGIC_RDX uint64 = 0x4e4f485950455256

	HOOK_PAGE_MONITOR_READ         uint32 = 0x1
	HOOK_PAGE_MONITOR_WRITE        uint32 = 0x2
	HOOK_PAGE_MONITOR_EXEC         uint32 = 0x4
	HOOK_PAGE_MONITOR_READ_WRITE   uint32 = 0x3
	HOOK_PAGE_MONITOR_EXEC_READ    uint32 = 0x5
	HOOK_PAGE_MONITOR_EXEC_WRITE   uint32 = 0x6
	HOOK_PAGE_MONITOR_EXEC_RW      uint32 = 0x7
	HOOK_PAGE_MONITOR_INLINE_HOOKS uint32 = 0x8
	HOOK_PAGE_MASKED_HOOKS         uint32 = 0x10
)

type GUEST_REGS struct {
	Rax, Rcx, Rdx, Rbx, Rsp, Rbp, Rsi, Rdi uint64
	R8, R9, R10, R11, R12, R13, R14, R15   uint64
}

type GUEST_XMM_REGS struct {
	Xmm0  [16]uint8
	Xmm1  [16]uint8
	Xmm2  [16]uint8
	Xmm3  [16]uint8
	Xmm4  [16]uint8
	Xmm5  [16]uint8
	Xmm6  [16]uint8
	Xmm7  [16]uint8
	Xmm8  [16]uint8
	Xmm9  [16]uint8
	Xmm10 [16]uint8
	Xmm11 [16]uint8
	Xmm12 [16]uint8
	Xmm13 [16]uint8
	Xmm14 [16]uint8
	Xmm15 [16]uint8
}

type CR3_TYPE struct{ Flags uint64 }

type EPT_PML4E struct{ AsUInt uint64 }
type EPT_PDPTE struct{ AsUInt uint64 }
type EPT_PDE struct{ AsUInt uint64 }
type EPT_PTE struct{ AsUInt uint64 }

type PML4E struct{ AsUInt uint64 }

type EPT_POINTER struct{ AsUInt uint64 }

type MTRR_RANGE_DESCRIPTOR struct {
	PhysicalBaseAddress uintptr
	PhysicalEndAddress  uintptr
	MemoryType          uint8
	FixedRange          bool
}

type VMXOFF_STATE struct {
	Executed bool
	Rip      uint64
	Rsp      uint64
}

type NMI_BROADCAST_STATE struct {
	Action int32
}

type VCPU struct {
	OnVmxRootMode          bool
	IncrementRip           bool
	HasLaunched            bool
	IgnoreMtfUnset         bool
	WaitForImmediateVmexit bool
	RegisterBreakOnMtf     bool
	IgnoreOneMtf           bool

	Regs              *GUEST_REGS
	XmmRegs           *GUEST_XMM_REGS
	CoreId            uint32
	ExitReason        uint32
	ExitQualification uint32
	LastVmexitRip     uint64

	VmxonRegionPhysicalAddress uint64
	VmxonRegionVirtualAddress  uint64
	VmcsRegionPhysicalAddress  uint64
	VmcsRegionVirtualAddress   uint64
	VmmStack                   uint64
	MsrBitmapVirtualAddress    uint64
	MsrBitmapPhysicalAddress   uint64
	IoBitmapVirtualAddressA    uint64
	IoBitmapPhysicalAddressA   uint64
	IoBitmapVirtualAddressB    uint64
	IoBitmapPhysicalAddressB   uint64

	QueuedNmi                 uint32
	PendingExternalInterrupts [PENDING_INTERRUPTS_CAPACITY]uint32
	VmxoffState               VMXOFF_STATE
	NmiBroadcastingState      NMI_BROADCAST_STATE
	MtfEptHookRestorePoint    uintptr

	EptPointer   EPT_POINTER
	EptPageTable *EPT_PAGE_TABLE
}

type EPT_PAGE_TABLE struct {
	PML4             [EPT_PML4_COUNT]EPT_PML4E
	PML3_RSVD        [EPT_PML4_COUNT - 1][EPT_PML3_COUNT]EPT_PDPTE
	PML3             [EPT_PML3_COUNT]EPT_PDPTE
	PML2             [EPT_PML3_COUNT][EPT_PML2_COUNT]EPT_PDE
	DynamicSplitList [16]uint8
}

type EPT_STATE struct {
	MemoryRanges                [MAX_MTRR_ENTRIES]MTRR_RANGE_DESCRIPTOR
	NumberOfEnabledMemoryRanges uint32
	DefaultMemoryType           uint8
}

type MEMORY_MAPPER struct {
	ReadPteAddr   uint64
	ReadVirtAddr  uint64
	WritePteAddr  uint64
	WriteVirtAddr uint64
}

type BUFFER_HEADER struct {
	OperationNumber uint32
	BufferLength    uint32
	Valid           bool
}

type LOG_BUFFER_INFO struct {
	BufferLock                           uintptr
	BufferLockForNonImmMessage           uintptr
	BufferForMultipleNonImmediateMessage uint64
	CurrentLengthOfNonImmBuffer          uint32
	BufferStartAddress                   uint64
	BufferEndAddress                     uint64
	CurrentIndexToSend                   uint32
	CurrentIndexToWrite                  uint32
	BufferStartAddressPriority           uint64
	BufferEndAddressPriority             uint64
	CurrentIndexToSendPriority           uint32
	CurrentIndexToWritePriority          uint32
}

func eptPml1Offset(v uint64) uint64 { return v & 0xFFF }
func eptPml1Index(v uint64) uint64  { return (v >> 12) & 0x1FF }
func eptPml2Index(v uint64) uint64  { return (v >> 21) & 0x1FF }
func eptPml3Index(v uint64) uint64  { return (v >> 30) & 0x1FF }
func eptPml4Index(v uint64) uint64  { return (v >> 39) & 0x1FF }

func logBufferSize() uint64 {
	return uint64(MAX_LOG_BUFFERS) * (uint64(LOG_CHUNK_SIZE) + uint64(unsafe.Sizeof(BUFFER_HEADER{})))
}
func logBufferSizePrio() uint64 {
	return uint64(MAX_LOG_BUFFERS_PRIO) * (uint64(LOG_CHUNK_SIZE) + uint64(unsafe.Sizeof(BUFFER_HEADER{})))
}

type VmxCtx = VCPU
