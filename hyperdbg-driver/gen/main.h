#pragma once

#ifndef _MAIN_GEN_H_
#define _MAIN_GEN_H_

#ifdef ENV_WINDOWS
#include <ntifs.h>
#include <ntstrsafe.h>
#endif
#include "builtin_kernel.h"
#include "wdk.h"

// -- Extern type aliases --
typedef int32_t NTSTATUS;
typedef bool BOOLEAN;
typedef uint64_t SIZE_T;
typedef int8_t CHAR;
typedef uint32_t ULONG;
typedef uint16_t USHORT;
typedef uint8_t UCHAR;
typedef uint64_t ULONG_PTR;


// -- Forward declarations --
typedef struct Driver Driver;
typedef struct PagePermission PagePermission;
typedef struct EptEntry EptEntry;
typedef struct MtrrState MtrrState;
typedef struct EventDispatcher EventDispatcher;
typedef struct CpuidHandler CpuidHandler;
typedef struct TscHandler TscHandler;
typedef struct ExecTrapState ExecTrapState;
typedef struct Breakpoint Breakpoint;
typedef struct LbrEntry LbrEntry;
typedef struct BtsEntry BtsEntry;
typedef struct TraceUserConfig TraceUserConfig;
typedef struct LbrConfig LbrConfig;
typedef struct BtsConfig BtsConfig;
typedef struct IoBitmapManager IoBitmapManager;
typedef struct VmmInitRequest VmmInitRequest;
typedef struct MemReadResponse MemReadResponse;
typedef struct RegRequest RegRequest;
typedef struct RegResponse RegResponse;
typedef struct InveptRequest InveptRequest;
typedef struct MtrrRequest MtrrRequest;
typedef struct ChangeCoreRequest ChangeCoreRequest;
typedef struct ProcessCr3Request ProcessCr3Request;
typedef struct CallFunctionRequest CallFunctionRequest;
typedef struct QueryPacketRequest QueryPacketRequest;
typedef struct TransparentModeRequest TransparentModeRequest;
typedef struct TraceStartRequest TraceStartRequest;
typedef struct TraceBufferRequest TraceBufferRequest;
typedef struct SerialConfig SerialConfig;
typedef struct MtrrCapabilities MtrrCapabilities;
typedef struct ModeBasedExecHookState ModeBasedExecHookState;
typedef struct ModeBasedExecEntry ModeBasedExecEntry;
typedef struct MsrBitmapManager MsrBitmapManager;
typedef struct PoolEntry PoolEntry;
typedef struct PoolStats PoolStats;
typedef struct VmmOptimizer VmmOptimizer;
typedef struct VmmOptConfig VmmOptConfig;
typedef struct EptCacheStats EptCacheStats;
typedef struct Spinlock Spinlock;
typedef struct RwLock RwLock;
typedef struct VmxRootSpinlock VmxRootSpinlock;
typedef struct EferHookState EferHookState;
typedef struct GUEST_REGS GUEST_REGS;
typedef struct GUEST_XMM_REGS GUEST_XMM_REGS;
typedef struct CR3_TYPE CR3_TYPE;
typedef struct EPT_PML4E EPT_PML4E;
typedef struct EPT_PDPTE EPT_PDPTE;
typedef struct EPT_PDE EPT_PDE;
typedef struct EPT_PTE EPT_PTE;
typedef struct PML4E PML4E;
typedef struct EPT_POINTER EPT_POINTER;
typedef struct MTRR_RANGE_DESCRIPTOR MTRR_RANGE_DESCRIPTOR;
typedef struct VMXOFF_STATE VMXOFF_STATE;
typedef struct NMI_BROADCAST_STATE NMI_BROADCAST_STATE;
typedef struct MEMORY_MAPPER MEMORY_MAPPER;
typedef struct BUFFER_HEADER BUFFER_HEADER;
typedef struct LOG_BUFFER_INFO LOG_BUFFER_INFO;
typedef struct Pml1Table Pml1Table;
typedef struct EptCacheEntry EptCacheEntry;
typedef struct TransparentModeConfig TransparentModeConfig;
typedef struct LogMessage LogMessage;
typedef struct MtrrRange MtrrRange;
typedef struct EptTable EptTable;
typedef struct HookManager HookManager;
typedef struct TracerState TracerState;
typedef struct KdSerialState KdSerialState;
typedef struct LogBuffer LogBuffer;
typedef struct PoolManager PoolManager;
typedef struct EptCache EptCache;
typedef struct ModifyRegsRequest ModifyRegsRequest;
typedef struct ExecTrapEntry ExecTrapEntry;
typedef struct EptHook EptHook;
typedef struct ProcessContext ProcessContext;
typedef struct EptHookRequest EptHookRequest;
typedef struct EptUnhookRequest EptUnhookRequest;
typedef struct MemReadRequest MemReadRequest;
typedef struct VirtPhysRequest VirtPhysRequest;
typedef struct SwitchProcessRequest SwitchProcessRequest;
typedef struct ReadWriteMemRequest ReadWriteMemRequest;
typedef struct EPT_PAGE_TABLE EPT_PAGE_TABLE;
typedef struct EPT_STATE EPT_STATE;
typedef struct SplitRecord SplitRecord;
typedef struct EvasionState EvasionState;
typedef struct MemoryMtrrState MemoryMtrrState;
typedef struct Logger Logger;
typedef struct VCPU VCPU;
typedef struct MemoryManager MemoryManager;
typedef struct EventHandler EventHandler;
typedef struct MsrHandler MsrHandler;
typedef struct IoHandler IoHandler;
typedef struct CrAccessHandler CrAccessHandler;
typedef struct ExceptionHandler ExceptionHandler;
typedef struct Hypervisor Hypervisor;
// -- Types --
typedef struct Driver {
    DEVICE_OBJECT* device;
    bool handleInUse;
    bool initialized;
} Driver;
typedef struct PagePermission {
    bool Read;
    bool Write;
    bool Exec;
} PagePermission;
typedef struct EptEntry {
    uint64_t raw;
} EptEntry;
typedef so_int Level;
typedef struct MtrrState {
    so_Slice ranges;
    uint8_t defaultType;
    uint32_t numRanges;
    bool fixedRangesInitialized;
} MtrrState;
typedef so_int EventKind;
typedef struct EventDispatcher {
    so_Map* handlers;
    so_Map* enabled;
} EventDispatcher;
typedef struct CpuidHandler {
    char _pad;
} CpuidHandler;
typedef struct TscHandler {
    int64_t offset;
    uint64_t (*rdtsc)();
} TscHandler;
typedef struct ExecTrapState {
    bool Enabled;
    so_Slice TrapList;
    bool MtfPending;
} ExecTrapState;
typedef so_int HookKind;
typedef struct Breakpoint {
    uint64_t Address;
    uint8_t OriginalByte;
    bool Active;
} Breakpoint;
typedef so_int EvasionTechnique;
typedef so_int DebuggerHidingMethod;
typedef so_int TraceMode;
typedef struct LbrEntry {
    uint64_t From;
    uint64_t To;
    uint64_t Info;
    uint64_t Misc;
} LbrEntry;
typedef struct BtsEntry {
    uint64_t LastBranchFrom;
    uint64_t LastBranchTo;
    uint64_t Flags;
} BtsEntry;
typedef struct TraceUserConfig {
    uint64_t LbrFilter;
    uint64_t BtsBufferSize;
    bool CallstackMode;
} TraceUserConfig;
typedef struct LbrConfig {
    bool Enabled;
    uint64_t Filter;
    so_int NumEntries;
    bool CallstackMode;
} LbrConfig;
typedef struct BtsConfig {
    bool Enabled;
    uint64_t BufferSize;
    uint64_t BufferPhysical;
    uintptr_t BufferVirtual;
    bool InterruptOnFull;
    uint32_t BranchType;
    uint64_t BaseMsrValue;
    uint64_t MaskMsrValue;
} BtsConfig;
typedef uint32_t ExitReason;
typedef struct IoBitmapManager {
    so_byte bitmapA[8192];
    so_byte bitmapB[8192];
} IoBitmapManager;
typedef struct VmmInitRequest {
    uint32_t NumCores;
} VmmInitRequest;
typedef struct MemReadResponse {
    uint64_t Address;
    uint32_t Size;
    bool Success;
} MemReadResponse;
typedef struct RegRequest {
    uint32_t RegIndex;
    uint32_t CoreId;
} RegRequest;
typedef struct RegResponse {
    uint64_t Value;
    uint32_t RegIndex;
    uint32_t CoreId;
    bool Success;
} RegResponse;
typedef struct InveptRequest {
    uint64_t Type;
    uint64_t Eptp;
} InveptRequest;
typedef struct MtrrRequest {
    uint64_t BaseAddr;
    uint64_t EndAddr;
} MtrrRequest;
typedef struct ChangeCoreRequest {
    uint32_t TargetCoreId;
} ChangeCoreRequest;
typedef struct ProcessCr3Request {
    uint32_t ProcessId;
} ProcessCr3Request;
typedef struct CallFunctionRequest {
    uint64_t FunctionAddress;
    uint32_t ProcessId;
    uint64_t OptionalParam1;
    uint64_t OptionalParam2;
    uint64_t OptionalParam3;
    uint64_t OptionalParam4;
} CallFunctionRequest;
typedef struct QueryPacketRequest {
    uint32_t PacketType;
} QueryPacketRequest;
typedef struct TransparentModeRequest {
    bool Enable;
    uint32_t Techniques;
    uint32_t DebuggerPid;
    uint32_t KdPid;
} TransparentModeRequest;
typedef struct TraceStartRequest {
    uint32_t Mode;
    uint64_t LbrFilter;
    uint64_t BtsBufferSize;
    bool CallstackMode;
    bool InterruptOnBts;
    uint32_t BranchType;
} TraceStartRequest;
typedef struct TraceBufferRequest {
    uint64_t BufferSize;
    uintptr_t Buffer;
} TraceBufferRequest;
typedef uint32_t SerialBaudRate;
typedef so_int ParityType;
typedef so_int StopBits;
typedef struct SerialConfig {
    uint32_t PortNumber;
    uint32_t BaudRate;
    uint32_t DataBits;
    so_int Parity;
    so_int StopBits;
    bool UseIrq;
    uint32_t IrqNumber;
} SerialConfig;
typedef so_int LogLevel;
typedef uint8_t MemoryType;
typedef struct MtrrCapabilities {
    so_int VarCnt;
    bool FixedSupported;
} MtrrCapabilities;
typedef struct ModeBasedExecHookState {
    bool Enabled;
    so_Slice CsSelectors;
} ModeBasedExecHookState;
typedef struct ModeBasedExecEntry {
    uint16_t Selector;
    bool Hooked;
} ModeBasedExecEntry;
typedef struct MsrBitmapManager {
    so_byte bitmap[8192];
} MsrBitmapManager;
typedef so_int OptimizationLevel;
typedef struct PoolEntry {
    uintptr_t Address;
    uint64_t Size;
    uint32_t Tag;
    PoolEntry* Next;
    PoolEntry* Prev;
} PoolEntry;
typedef struct PoolStats {
    uint64_t TotalAllocated;
    uint64_t TotalFreed;
    so_int CurrentPools;
    so_int MaxPools;
} PoolStats;
typedef struct VmmOptimizer {
    so_int level;
    bool eptCacheEnabled;
    bool msrBitmapCached;
    bool ioBitmapCached;
    bool inveptBatching;
} VmmOptimizer;
typedef struct VmmOptConfig {
    bool EptCache;
    bool MsrBitmapCache;
    bool IoBitmapCache;
    bool InveptBatching;
} VmmOptConfig;
typedef struct EptCacheStats {
    so_int TotalEntries;
    so_int MaxEntries;
    uint64_t TotalHits;
} EptCacheStats;
typedef struct Spinlock {
    uint32_t lock;
} Spinlock;
typedef struct RwLock {
    uint32_t readers;
    uint32_t writer;
} RwLock;
typedef struct VmxRootSpinlock {
    uintptr_t raw;
} VmxRootSpinlock;
typedef struct EferHookState {
    bool Enabled;
    bool UdHandling;
    bool CetSupport;
} EferHookState;
typedef struct GUEST_REGS {
    uint64_t Rax;
    uint64_t Rcx;
    uint64_t Rdx;
    uint64_t Rbx;
    uint64_t Rsp;
    uint64_t Rbp;
    uint64_t Rsi;
    uint64_t Rdi;
    uint64_t R8;
    uint64_t R9;
    uint64_t R10;
    uint64_t R11;
    uint64_t R12;
    uint64_t R13;
    uint64_t R14;
    uint64_t R15;
} GUEST_REGS;
typedef struct GUEST_XMM_REGS {
    uint8_t Xmm0[16];
    uint8_t Xmm1[16];
    uint8_t Xmm2[16];
    uint8_t Xmm3[16];
    uint8_t Xmm4[16];
    uint8_t Xmm5[16];
    uint8_t Xmm6[16];
    uint8_t Xmm7[16];
    uint8_t Xmm8[16];
    uint8_t Xmm9[16];
    uint8_t Xmm10[16];
    uint8_t Xmm11[16];
    uint8_t Xmm12[16];
    uint8_t Xmm13[16];
    uint8_t Xmm14[16];
    uint8_t Xmm15[16];
} GUEST_XMM_REGS;
typedef struct CR3_TYPE {
    uint64_t Flags;
} CR3_TYPE;
typedef struct EPT_PML4E {
    uint64_t AsUInt;
} EPT_PML4E;
typedef struct EPT_PDPTE {
    uint64_t AsUInt;
} EPT_PDPTE;
typedef struct EPT_PDE {
    uint64_t AsUInt;
} EPT_PDE;
typedef struct EPT_PTE {
    uint64_t AsUInt;
} EPT_PTE;
typedef struct PML4E {
    uint64_t AsUInt;
} PML4E;
typedef struct EPT_POINTER {
    uint64_t AsUInt;
} EPT_POINTER;
typedef struct MTRR_RANGE_DESCRIPTOR {
    uintptr_t PhysicalBaseAddress;
    uintptr_t PhysicalEndAddress;
    uint8_t MemoryType;
    bool FixedRange;
} MTRR_RANGE_DESCRIPTOR;
typedef struct VMXOFF_STATE {
    bool Executed;
    uint64_t Rip;
    uint64_t Rsp;
} VMXOFF_STATE;
typedef struct NMI_BROADCAST_STATE {
    int32_t Action;
} NMI_BROADCAST_STATE;
typedef struct MEMORY_MAPPER {
    uint64_t ReadPteAddr;
    uint64_t ReadVirtAddr;
    uint64_t WritePteAddr;
    uint64_t WriteVirtAddr;
} MEMORY_MAPPER;
typedef struct BUFFER_HEADER {
    uint32_t OperationNumber;
    uint32_t BufferLength;
    bool Valid;
} BUFFER_HEADER;
typedef struct LOG_BUFFER_INFO {
    uintptr_t BufferLock;
    uintptr_t BufferLockForNonImmMessage;
    uint64_t BufferForMultipleNonImmediateMessage;
    uint32_t CurrentLengthOfNonImmBuffer;
    uint64_t BufferStartAddress;
    uint64_t BufferEndAddress;
    uint32_t CurrentIndexToSend;
    uint32_t CurrentIndexToWrite;
    uint64_t BufferStartAddressPriority;
    uint64_t BufferEndAddressPriority;
    uint32_t CurrentIndexToSendPriority;
    uint32_t CurrentIndexToWritePriority;
} LOG_BUFFER_INFO;
typedef uint64_t VmcallNumber;
typedef struct Pml1Table {
    EptEntry entries[512];
} Pml1Table;
typedef struct EptCacheEntry {
    uint64_t Gpa;
    uint64_t Pa;
    EptEntry Entry;
    uint64_t LastAccess;
    uint64_t HitCount;
} EptCacheEntry;
typedef struct TransparentModeConfig {
    bool Enabled;
    EvasionTechnique Techniques;
    uint32_t DebuggerPid;
} TransparentModeConfig;
typedef struct LogMessage {
    LogLevel Level;
    uint64_t Time;
    uint32_t CpuId;
    uint8_t Process[15];
    uint8_t Message[256];
    uint32_t Length;
} LogMessage;
typedef struct MtrrRange {
    uint64_t BaseAddr;
    uint64_t EndAddr;
    MemoryType Type;
    bool Valid;
} MtrrRange;
typedef struct EptTable {
    EptEntry pml4[512];
    EptEntry (*pdptTables[511])[512];
    EptEntry* pdTables[512][512];
    so_Slice splits;
    MtrrState* mtrrState;
    Spinlock* lock;
} EptTable;
typedef struct HookManager {
    so_Slice hooks;
    so_int capacity;
    Spinlock* lock;
    so_Slice mtfRestoreList;
} HookManager;
typedef struct TracerState {
    LbrConfig lbrConfig;
    BtsConfig btsConfig;
    bool active;
    Spinlock* lock;
    so_Slice lbrBuffer;
    so_Slice btsBuffer;
} TracerState;
typedef struct KdSerialState {
    SerialConfig config;
    bool connected;
    Spinlock* lock;
    so_Slice txBuffer;
    so_Slice rxBuffer;
    uint32_t txHead;
    uint32_t txTail;
    uint32_t rxHead;
    uint32_t rxTail;
    uint32_t bufferSize;
    uint16_t portBase;
    uint8_t lineControl;
    uint8_t modemControl;
    uint8_t intEnable;
    bool divisorLatch;
} KdSerialState;
typedef struct LogBuffer {
    so_Slice Messages;
    uint32_t ReadIndex;
    uint32_t WriteIndex;
    uint32_t Count;
    uint32_t Capacity;
    Spinlock* Lock;
    uint64_t OverflowCount;
} LogBuffer;
typedef struct PoolManager {
    PoolEntry* poolList;
    uint64_t totalAllocated;
    uint64_t totalFreed;
    Spinlock* lock;
    so_int maxPools;
    so_int currentPools;
} PoolManager;
typedef struct EptCache {
    so_Map* entries;
    Spinlock* lock;
    so_int maxSize;
} EptCache;
typedef struct ModifyRegsRequest {
    GUEST_REGS Regs;
    uint32_t CoreId;
} ModifyRegsRequest;
typedef struct ExecTrapEntry {
    uint64_t Address;
    CR3_TYPE Cr3;
    uint64_t OriginalRip;
    bool Active;
} ExecTrapEntry;
typedef struct EptHook {
    uint64_t PhysAddr;
    uint64_t VirtAddr;
    CR3_TYPE Cr3;
    uint64_t FakePagePa;
    EptEntry* OriginalEntry;
    EptEntry* ModifiedEntry;
    HookKind Kind;
    bool IsActive;
    so_Slice breakpoints;
    so_byte fakePage[4096];
    uint32_t coreId;
} EptHook;
typedef struct ProcessContext {
    uint32_t ProcessId;
    CR3_TYPE Cr3;
    uint64_t BaseAddress;
    int8_t ImageName[256];
} ProcessContext;
typedef struct EptHookRequest {
    uint64_t VirtAddr;
    CR3_TYPE Cr3;
    uint32_t HookType;
    uint32_t ProcessId;
} EptHookRequest;
typedef struct EptUnhookRequest {
    uint64_t VirtAddr;
    CR3_TYPE Cr3;
    uint32_t ProcessId;
} EptUnhookRequest;
typedef struct MemReadRequest {
    uint64_t Address;
    CR3_TYPE Cr3;
    uint32_t Size;
    uint32_t ProcessId;
    bool ReadOrWrite;
} MemReadRequest;
typedef struct VirtPhysRequest {
    uint64_t Address;
    CR3_TYPE Cr3;
    uint32_t ProcessId;
} VirtPhysRequest;
typedef struct SwitchProcessRequest {
    uint32_t ProcessId;
    CR3_TYPE Cr3;
    uint64_t ProcessBaseAddress;
} SwitchProcessRequest;
typedef struct ReadWriteMemRequest {
    uint64_t ReadAddress;
    uint64_t WriteAddress;
    uint32_t ReadSize;
    uint32_t WriteSize;
    bool ReadOrWrite;
    CR3_TYPE Cr3;
    uint32_t ProcessId;
} ReadWriteMemRequest;
typedef struct EPT_PAGE_TABLE {
    EPT_PML4E PML4[512];
    EPT_PDPTE PML3_RSVD[511][512];
    EPT_PDPTE PML3[512];
    EPT_PDE PML2[512][512];
    uint8_t DynamicSplitList[16];
} EPT_PAGE_TABLE;
typedef struct EPT_STATE {
    MTRR_RANGE_DESCRIPTOR MemoryRanges[255];
    uint32_t NumberOfEnabledMemoryRanges;
    uint8_t DefaultMemoryType;
} EPT_STATE;
typedef struct SplitRecord {
    EptEntry* pdEntry;
    Pml1Table* pt;
    so_int refCount;
    uint64_t gpaBase;
} SplitRecord;
typedef struct EvasionState {
    TransparentModeConfig config;
    bool active;
    so_Slice hiddenPids;
    Spinlock* lock;
    uint64_t vmxOriginalCr4;
    so_Map* msrBackup;
    so_Map* ioPortHooks;
    uint64_t tscBase;
    int64_t tscOffset;
    uint64_t (*originalRdtscFn)();
} EvasionState;
typedef struct MemoryMtrrState {
    MemoryType DefaultType;
    MtrrCapabilities Capabilities;
    MtrrRange VariableRanges[255];
} MemoryMtrrState;
typedef struct Logger {
    LogBuffer* buffer;
    LogBuffer* priorityBuffer;
    bool enabled;
    LogLevel minLevel;
    Spinlock* lock;
} Logger;
typedef struct VCPU {
    bool OnVmxRootMode;
    bool IncrementRip;
    bool HasLaunched;
    bool IgnoreMtfUnset;
    bool WaitForImmediateVmexit;
    bool RegisterBreakOnMtf;
    bool IgnoreOneMtf;
    GUEST_REGS* Regs;
    GUEST_XMM_REGS* XmmRegs;
    uint32_t CoreId;
    uint32_t ExitReason;
    uint32_t ExitQualification;
    uint64_t LastVmexitRip;
    uint64_t VmxonRegionPhysicalAddress;
    uint64_t VmxonRegionVirtualAddress;
    uint64_t VmcsRegionPhysicalAddress;
    uint64_t VmcsRegionVirtualAddress;
    uint64_t VmmStack;
    uint64_t MsrBitmapVirtualAddress;
    uint64_t MsrBitmapPhysicalAddress;
    uint64_t IoBitmapVirtualAddressA;
    uint64_t IoBitmapPhysicalAddressA;
    uint64_t IoBitmapVirtualAddressB;
    uint64_t IoBitmapPhysicalAddressB;
    uint32_t QueuedNmi;
    uint32_t PendingExternalInterrupts[64];
    VMXOFF_STATE VmxoffState;
    NMI_BROADCAST_STATE NmiBroadcastingState;
    uintptr_t MtfEptHookRestorePoint;
    EPT_POINTER EptPointer;
    EPT_PAGE_TABLE* EptPageTable;
} VCPU;
typedef struct MemoryManager {
    MemoryMtrrState* mtrr;
} MemoryManager;
typedef struct EventHandler {
    void* self;
    bool (*Handle)(void* self, VCPU* v);
} EventHandler;
typedef bool (*EventFunc)(VCPU*);
typedef struct MsrHandler {
    bool (*onRead)(VCPU* v, uint32_t msr, uint64_t value);
    bool (*onWrite)(VCPU* v, uint32_t msr, uint64_t value);
} MsrHandler;
typedef struct IoHandler {
    bool (*onIn)(VCPU* v, uint16_t port, uint32_t size);
    bool (*onOut)(VCPU* v, uint16_t port, uint32_t size);
} IoHandler;
typedef struct CrAccessHandler {
    bool (*onChange)(VCPU* v, uint32_t crNum, bool read, bool write);
} CrAccessHandler;
typedef struct ExceptionHandler {
    bool (*onBp)(VCPU* v);
    bool (*onGp)(VCPU* v);
    bool (*onUd)(VCPU* v);
    bool (*onNmi)(VCPU* v);
} ExceptionHandler;
typedef VCPU VmxCtx;
typedef struct Hypervisor {
    so_Slice vcpus;
    so_Slice eptTables;
    EventDispatcher* events;
    HookManager* hooks;
    Logger* logger;
    MemoryManager* memory;
    so_Map* processContexts;
    uint32_t activeCoreId;
    bool checkFootprints;
    bool initialized;
    bool paused;
    Spinlock* lock;
} Hypervisor;

// -- Variables and constants --
extern so_Error ErrNotPresent;
extern so_Error ErrLargePage;
extern so_Error ErrNotSplit;
extern so_Error ErrInvalidLevel;
extern so_Error ErrNotMapped;
extern so_Error ErrInvalidAddress;
extern so_Error ErrHookLimitReached;
extern so_Error ErrNoEptTable;
extern so_Error ErrHookNotFound;
extern so_Error ErrTooManyBreakpoints;
extern so_Error ErrInvalidCore;
extern so_Error ErrVmcsAlloc;
extern so_Error ErrVmxonAlloc;

// -- Functions and methods --
NTSTATUS DriverEntry(DRIVER_OBJECT* driverObj, UNICODE_STRING* registryPath);
uint64_t PagePermission_ToRaw(PagePermission p);
PagePermission PagePermissionFromRaw(uint64_t raw);
EptEntry* NewEptEntry(void);
uint64_t EptEntry_Raw(void* self);
void EptEntry_SetRaw(void* self, uint64_t v);
bool EptEntry_IsPresent(void* self);
uint64_t EptEntry_PhysAddr(void* self);
void EptEntry_SetPhysAddr(void* self, uint64_t pa);
uint8_t EptEntry_MemoryType(void* self);
void EptEntry_SetMemoryType(void* self, uint8_t t);
PagePermission EptEntry_Permission(void* self);
void EptEntry_SetPermission(void* self, PagePermission p);
bool EptEntry_IsLargePage(void* self);
bool EptEntry_IgnorePat(void* self);
bool EptEntry_SuppressVe(void* self);
void EptEntry_SetIgnorePat(void* self, bool v);
void EptEntry_SetSuppressVe(void* self, bool v);
void EptEntry_SetLargePage(void* self, bool v);
EptEntry* EptEntry_Clone(void* self);
EptTable* NewEptTable(void);
void EptTable_Lock(void* self);
void EptTable_Unlock(void* self);
so_R_ptr_err EptTable_Walk(void* self, uint64_t gpa, so_int level);
so_R_ptr_err EptTable_Lookup(void* self, uint64_t gpa);
so_R_ptr_err EptTable_ResolveForModify(void* self, uint64_t gpa);
so_Error EptTable_SplitLargePage(void* self, uint64_t gpa);
so_Error EptTable_MergeLargePage(void* self, uint64_t gpa);
so_Error EptTable_Map(void* self, uint64_t gpa, uint64_t pa, PagePermission perm, uint8_t memType);
so_Error EptTable_Unmap(void* self, uint64_t gpa);
so_Error EptTable_Protect(void* self, uint64_t gpa, PagePermission perm);
uint64_t EptTable_BuildEptPointer(void* self);
void EptTable_Invalidate(void* self);
so_int EptTable_SplitCount(void* self);
MtrrState* NewMtrrState(void);
uint8_t MtrrState_DefaultMemoryType(void* self);
uint8_t MtrrState_Lookup(void* self, uint64_t pa);
so_Slice MtrrState_Ranges(void* self);
so_Error EptTable_RestoreEntry(void* self, uint64_t physAddr, EptEntry* original);
so_R_ptr_err EptTable_GetEntry(void* self, uint64_t gpa);
so_R_ptr_err EptTable_CloneEntry(void* self, uint64_t gpa);
bool EventFunc_Handle(EventFunc f, VCPU* v);
void EventDispatcher_On(void* self, so_int kind, EventHandler handler);
void EventDispatcher_OnFunc(void* self, so_int kind, EventFunc fn);
void EventDispatcher_Enable(void* self, so_int kind);
void EventDispatcher_Disable(void* self, so_int kind);
bool EventDispatcher_IsEnabled(void* self, so_int kind);
bool EventDispatcher_Dispatch(void* self, so_int kind, VCPU* v);
void EventDispatcher_Reset(void* self);
bool CpuidHandler_Handle(void* self, VCPU* v);
MsrHandler* NewMsrHandler(void);
bool MsrHandler_HandleRead(void* self, VCPU* v);
bool MsrHandler_HandleWrite(void* self, VCPU* v);
IoHandler* NewIoHandler(void);
bool IoHandler_Handle(void* self, VCPU* v);
CrAccessHandler* NewCrAccessHandler(void);
bool CrAccessHandler_Handle(void* self, VCPU* v);
TscHandler* NewTscHandler(void);
bool TscHandler_Handle(void* self, VCPU* v);
ExceptionHandler* NewExceptionHandler(void);
bool ExceptionHandler_HandleBreakpoint(void* self, VCPU* v);
bool ExceptionHandler_HandleGeneralProtectionFault(void* self, VCPU* v);
bool ExceptionHandler_HandleUndefinedOpcode(void* self, VCPU* v);
bool ExceptionHandler_HandleNmi(void* self, VCPU* v);
so_Error InitializeGlobals(uint32_t numCpus);
void CleanupGlobals(void);
Hypervisor* GetHypervisor(void);
Logger* GetLogger(void);
Driver* GetDriver(void);
EventDispatcher* GetEventDispatcher(void);
HookManager* GetHookManager(void);
EvasionState* GetEvasionState(void);
TracerState* GetTracerState(void);
KdSerialState* GetSerialState(void);
PoolManager* GetPoolManager(void);
VmmOptimizer* GetOptimizer(void);
EptCache* GetEptCache(void);
so_String HookKind_String(HookKind k);
PagePermission HookKind_Permission(HookKind k);
so_R_ptr_err NewEptHook(uint64_t va, CR3_TYPE cr3, HookKind kind);
bool EptHook_AddBreakpoint(void* self, uint64_t addr);
bool EptHook_RemoveBreakpoint(void* self, uint64_t addr);
Breakpoint* EptHook_FindBreakpoint(void* self, uint64_t addr);
so_int EptHook_BreakpointCount(void* self);
bool EptHook_IsHiddenBp(void* self);
bool EptHook_ContainsAddress(void* self, uint64_t va);
HookManager* NewHookManager(so_int capacity);
void HookManager_Lock(void* self);
void HookManager_Unlock(void* self);
so_Error HookManager_Install(void* self, uint32_t coreId, uint64_t va, CR3_TYPE cr3, HookKind kind);
so_Error HookManager_Remove(void* self, uint64_t physAddr);
so_Error HookManager_RemoveByVirtAddr(void* self, uint64_t va);
EptHook* HookManager_FindByPhysAddr(void* self, uint64_t pa);
EptHook* HookManager_FindByVirtAddr(void* self, uint64_t va);
bool HookManager_HandleExecHook(void* self, VCPU* v, uint64_t gpa);
bool HookManager_HandleReadWriteHook(void* self, VCPU* v, uint64_t gpa, bool isRead, bool isWrite);
bool HookManager_HandleBreakpoint(void* self, VCPU* v);
bool HookManager_HandleMtfRestore(void* self, VCPU* v);
void HookManager_RestoreAll(void* self, VCPU* v);
void HookManager_RestoreAllForCore(void* self, uint32_t coreId);
so_int HookManager_Count(void* self);
so_int HookManager_ActiveCount(void* self);
void HookManager_RemoveAll(void* self);
void HookManager_DisableAll(void* self);
void HookManager_EnableAll(void* self);
so_String EvasionTechnique_String(EvasionTechnique e);
EvasionState* NewEvasionState(void);
so_Error EvasionState_Enable(void* self, TransparentModeConfig* config);
so_Error EvasionState_Disable(void* self);
bool EvasionState_IsActive(void* self);
void EvasionState_AddHiddenProcess(void* self, uint32_t pid);
void EvasionState_RemoveHiddenProcess(void* self, uint32_t pid);
bool EvasionState_IsProcessHidden(void* self, uint32_t pid);
bool EvasionState_ShouldInterceptIoPort(void* self, uint16_t port);
bool EvasionState_ShouldHideMsr(void* self, uint32_t msr);
uint64_t EvasionState_AdjustTsc(void* self, uint64_t tsc);
void EvasionState_SetTscOffset(void* self, int64_t offset);
bool EvasionState_ShouldHideProcess(void* self, uint32_t pid);
so_int EvasionState_HiddenProcessCount(void* self);
TransparentModeConfig EvasionState_GetConfig(void* self);
so_String TraceMode_String(TraceMode t);
TracerState* NewTracerState(void);
so_Error TracerState_Enable(void* self, TraceMode mode, TraceUserConfig* config);
so_Error TracerState_Disable(void* self);
bool TracerState_IsActive(void* self);
so_R_slice_int TracerState_CaptureLbr(void* self);
so_R_slice_int TracerState_CaptureBts(void* self);
Hypervisor* NewHypervisor(uint32_t numCpus);
so_Error Hypervisor_Initialize(void* self);
void Hypervisor_Shutdown(void* self);
void Hypervisor_SetInitialized(void* self, bool val);
bool Hypervisor_IsInitialized(void* self);
bool Hypervisor_HandleVmExit(void* self, GUEST_REGS* regs);
VCPU* Hypervisor_Vcpu(void* self, uint32_t id);
VCPU* Hypervisor_CurrentVcpu(void* self);
MemoryManager* Hypervisor_MemoryManager(void* self);
void Hypervisor_SetActiveCore(void* self, uint32_t id);
uint32_t Hypervisor_GetActiveCore(void* self);
void Hypervisor_PauseAll(void* self);
void Hypervisor_ResumeAll(void* self);
bool Hypervisor_IsPaused(void* self);
void Hypervisor_EnableSingleStep(void* self, VCPU* _p0);
void Hypervisor_SetCurrentProcessContext(void* self, uint32_t pid, CR3_TYPE cr3, uint64_t baseAddr);
uint64_t Hypervisor_GetProcessBaseAddress(void* self, uint32_t pid);
CR3_TYPE Hypervisor_GetProcessCr3(void* self, uint32_t pid);
uint64_t Hypervisor_CallGuestFunction(void* self, uint64_t _p0, uint64_t _p1, uint64_t _p2, uint64_t _p3, uint64_t _p4);
void DebuggerUninitialize(void);
void VmFuncUninitVmm(void);
KdSerialState* NewKdSerialState(void);
so_Error KdSerialState_Initialize(void* self, SerialConfig* config);
so_Error KdSerialState_Uninitialize(void* self);
bool KdSerialState_IsConnected(void* self);
so_R_int_err KdSerialState_Send(void* self, so_Slice data);
so_R_int_err KdSerialState_Receive(void* self, so_Slice buffer);
so_R_int_err KdSerialState_SendString(void* self, so_String str);
so_R_str_err KdSerialState_ReceiveString(void* self, so_int maxLen);
uint32_t KdSerialState_BytesAvailable(void* self);
uint32_t KdSerialState_TxSpaceAvailable(void* self);
so_Error KdSerialState_SendByte(void* self, so_byte b);
so_R_byte_err KdSerialState_ReceiveByte(void* self);
so_String LogLevel_String(LogLevel l);
LogBuffer* NewLogBuffer(uint32_t capacity);
bool LogBuffer_Push(void* self, LogMessage* msg);
so_R_ptr_bool LogBuffer_Pop(void* self);
so_R_ptr_bool LogBuffer_Peek(void* self);
void LogBuffer_Clear(void* self);
bool LogBuffer_IsEmpty(void* self);
bool LogBuffer_IsFull(void* self);
uint32_t LogBuffer_Available(void* self);
uint32_t LogBuffer_GetCount(void* self);
Logger* NewLogger(uint32_t bufferSize, uint32_t prioSize);
so_Error Logger_Initialize(void* self);
void Logger_Uninitialize(void* self);
void Logger_SetMinLevel(void* self, LogLevel level);
void Logger_Enable(void* self);
void Logger_Disable(void* self);
void Logger_Log(void* self, LogLevel level, so_String format, so_Slice args);
void Logger_Info(void* self, so_String format, so_Slice args);
void Logger_Warning(void* self, so_String format, so_Slice args);
void Logger_Error(void* self, so_String format, so_Slice args);
void Logger_Debug(void* self, so_String format, so_Slice args);
void Logger_Trace(void* self, so_String format, so_Slice args);
uint32_t Logger_FlushToUser(void* self, uintptr_t buffer, uintptr_t size);
uint32_t Logger_FlushPriorityToUser(void* self, uintptr_t buffer, uintptr_t size);
uint64_t Logger_GetBufferSize(void* self);
uint64_t Logger_GetCurrentBufferSize(void* self);
uint64_t Logger_GetOverflowCount(void* self);
uint64_t Logger_GetBufferBase(void* self);
uint64_t Logger_GetPriorityBufferBase(void* self);
void LogInfo(so_String format, so_Slice args);
void LogError(so_String format, so_Slice args);
void LogWarning(so_String format, so_Slice args);
void LogDebug(so_String format, so_Slice args);
void LogTrace(so_String format, so_Slice args);
void LogCallbackSendBuffer(uint32_t operation, so_String data, uint32_t length, bool isImmediate);
void LogRegisterIrpBasedNotification(void* irp);
bool LogRegisterEventBasedNotification(void* irp);
so_String MemoryType_String(MemoryType m);
MemoryType MemoryManager_QueryMtrrForPa(void* self, uint64_t pa);
void MemoryManager_SetDefaultMemType(void* self, MemoryType t);
void MemoryManager_AddVariableRange(void* self, uint64_t base, uint64_t end, MemoryType t);
void MemoryManager_ClearAllRanges(void* self);
uint64_t VirtToPhys(uintptr_t va, CR3_TYPE cr3);
uintptr_t PhysToVirt(uint64_t pa, CR3_TYPE cr3);
void CopyPage(uintptr_t src, uintptr_t dst, CR3_TYPE cr3);
so_int ReadPhysMem(uint64_t pa, so_Slice buf);
so_int WritePhysMem(uint64_t pa, so_Slice data);
PoolManager* NewPoolManager(so_int maxEntries);
uintptr_t PoolManager_Allocate(void* self, uint64_t size, uint32_t tag);
void PoolManager_Free(void* self, uintptr_t ptr);
so_int PoolManager_CheckAndPerformDeallocation(void* self);
PoolStats PoolManager_GetStats(void* self);
so_String PoolStats_String(PoolStats s);
VmmOptimizer* NewVmmOptimizer(so_int level);
void VmmOptimizer_EnableEptCache(void* self, bool enable);
void VmmOptimizer_EnableMsrBitmapCache(void* self, bool enable);
void VmmOptimizer_EnableIoBitmapCache(void* self, bool enable);
void VmmOptimizer_EnableInveptBatching(void* self, bool enable);
VmmOptConfig VmmOptimizer_GetConfig(void* self);
EptCache* NewEptCache(so_int maxSize);
so_R_ptr_bool EptCache_Lookup(void* self, uint64_t gpa);
void EptCache_Store(void* self, uint64_t gpa, uint64_t pa, EptEntry entry);
void EptCache_Invalidate(void* self, uint64_t gpa);
void EptCache_InvalidateAll(void* self);
EptCacheStats EptCache_GetStats(void* self);
uintptr_t AllocatePool(uintptr_t size, uint32_t tag);
uintptr_t AllocateZeroedPool(uintptr_t size, uint32_t tag);
void FreePool(uintptr_t ptr, uint32_t tag);
so_String CurrentProcessName(void);
Spinlock* NewSpinlock(void);
void Spinlock_Lock(void* self);
void Spinlock_Unlock(void* self);
bool Spinlock_TryLock(void* self);
bool Spinlock_IsLocked(void* self);
RwLock* NewRwLock(void);
void RwLock_RLock(void* self);
void RwLock_RUnlock(void* self);
void RwLock_Lock(void* self);
void RwLock_Unlock(void* self);
VmxRootSpinlock* NewVmxRootSpinlock(void);
void VmxRootSpinlock_Lock(void* self);
void VmxRootSpinlock_Unlock(void* self);
extern void __writecr3(uint64_t val);
uint64_t AsmVmxVmcall(uint64_t RegPtr);
extern uint8_t AsmVmxVmread(uint64_t Field, uint64_t* FieldValue);
extern uint8_t AsmVmxVmwrite(uint64_t Field, uint64_t FieldValue);
extern uint8_t AsmVmxVmread32(uint64_t Field, uint32_t* FieldValue);
extern uint8_t AsmVmxVmwrite32(uint64_t Field, uint32_t FieldValue);
extern uint8_t AsmVmxVmlaunch(void);
extern uint8_t AsmVmxVmresume(void);
extern uint8_t AsmVmxVmxClear(uint64_t PhysicalAddr);
extern uint8_t AsmVmxVmxPtrld(uint64_t PhysicalAddr);
uint8_t AsmEnableVmxOperation(uint64_t PhysicalAddr);
extern void AsmVmxVmxOff(void);
void AsmInvept(uint64_t Eptp);
void AsmInveptAllContexts(void);
void AsmInvvpid(void);
void AsmInvvpidAllContexts(void);
void AsmHypervVmcall(uint64_t regs);
uint64_t AsmGetGdtBase(void);
uint64_t AsmGetIdtBase(void);
uint32_t AsmGetGdtLimit(void);
uint32_t AsmGetIdtLimit(void);
uint64_t AsmGetRflags(void);
uint16_t AsmGetEs(void);
uint16_t AsmGetCs(void);
uint16_t AsmGetSs(void);
uint16_t AsmGetDs(void);
uint16_t AsmGetFs(void);
uint16_t AsmGetGs(void);
uint16_t AsmGetTr(void);
uint16_t AsmGetLdtr(void);
uint64_t AsmGetFsBase(void);
uint64_t AsmGetGsBase(void);
uintptr_t AsmVmexitHandlerAddr(void);
extern uint64_t __readcr4(void);
extern uint64_t __readcr3(void);
extern uint64_t __readcr0(void);
extern void __writecr4(uint64_t cr4);
bool VmxVmexitHandler(GUEST_REGS* regs);
uint8_t VmxVmresume(void);
uint64_t VmxReturnStackPointerForVmxoff(void);
uint64_t VmxReturnInstructionPointerForVmxoff(void);
void VmxVirtualizeCurrentSystem(uint8_t* state);
uint64_t EptHook2GeneralDetourEventHandler(GUEST_REGS* regs, uint64_t calledFrom);
void IdtEmulationhandleHostInterrupt(uint8_t* trapFrame);

#endif // _MAIN_GEN_H_
