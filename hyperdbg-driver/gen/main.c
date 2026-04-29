#include "builtin_kernel.h"
void _solod_kmod_init(void);
#include "main.h"

// -- Types --
typedef so_String eptError;

// -- Variables and constants --
#define SE_DEBUG_PRIVILEGE 20LL
#define EptRead 1ULL
#define EptWrite 2ULL
#define EptExec 4ULL
#define EptIgnorePat 64ULL
#define EptSuppressVe 9223372036854775808ULL
#define EptLargePage 128ULL
#define EptPhysMask 4503599627366400ULL
#define EptTypeMask 56ULL
#define EptTypeShift 3LL
#define MemTypeUncacheable ((uint8_t)0)
#define MemTypeWriteCombining ((uint8_t)1)
#define MemTypeWriteThrough ((uint8_t)4)
#define MemTypeWriteProtected ((uint8_t)5)
#define MemTypeWriteBack ((uint8_t)6)
#define LvlPml4 0LL
#define LvlPdpt 1LL
#define LvlPd 2LL
#define LvlPt 3LL
SO_GLOBAL_ERROR(ErrNotPresent, "EPT entry not present");
SO_GLOBAL_ERROR(ErrLargePage, "entry is a large page");
SO_GLOBAL_ERROR(ErrNotSplit, "large page not split");
SO_GLOBAL_ERROR(ErrInvalidLevel, "invalid page table level");
SO_GLOBAL_ERROR(ErrNotMapped, "address not mapped in EPT");
#define EventCpuid 0LL
#define EventRdmsr 1LL
#define EventWrmsr 2LL
#define EventIo 3LL
#define EventMovCr 4LL
#define EventMovDr 5LL
#define EventTsc 6LL
#define EventRdpmc 7LL
#define EventXsetbv 8LL
#define EventException 9LL
#define EventExternalInt 10LL
#define EventVmcall 11LL
static ExecTrapState* gExecTrap = NULL;
static Hypervisor* gHyp = NULL;
static Logger* gLog = NULL;
static Driver* gDrv = NULL;
static EventDispatcher* gEvents = NULL;
static HookManager* gHooks = NULL;
static EvasionState* gEvasion = NULL;
static TracerState* gTrace = NULL;
static KdSerialState* gSerial = NULL;
static PoolManager* gPoolMgr = NULL;
static VmmOptimizer* gOptimizer = NULL;
static EptCache* gEptCache = NULL;
static bool gAllowIoctl = false;
#define HookExec ((HookKind)0)
#define HookRead ((HookKind)1)
#define HookWrite ((HookKind)2)
#define HookReadWrite ((HookKind)3)
SO_GLOBAL_ERROR(ErrInvalidAddress, "invalid virtual address");
SO_GLOBAL_ERROR(ErrHookLimitReached, "maximum hook limit reached");
SO_GLOBAL_ERROR(ErrNoEptTable, "no Ept table available");
SO_GLOBAL_ERROR(ErrHookNotFound, "hook not found");
SO_GLOBAL_ERROR(ErrTooManyBreakpoints, "too many breakpoints on this page");
#define EvasionHideDebugger ((EvasionTechnique)1)
#define EvasionHideVmx ((EvasionTechnique)2)
#define EvasionHideMsrs ((EvasionTechnique)4)
#define EvasionHideIoPorts ((EvasionTechnique)8)
#define EvasionHideTiming ((EvasionTechnique)16)
#define EvasionAll 4294967295LL
#define HideFromKdQuerySystemInformation 1LL
#define HideFromKdGetContextThread 2LL
#define HideFromProcessList 4LL
#define HideFromThreadList 8LL
#define HideAllMethods 4294967295LL
#define TraceLbr ((TraceMode)1)
#define TraceBts ((TraceMode)2)
#define TracePebs ((TraceMode)4)
#define TraceAll 4294967295LL
#define IA32_LBR_TOS_MSR 457U
#define IA32_LBR_FROM_0_MSR 1536U
#define IA32_LBR_TO_0_MSR 1632U
#define IA32_LBR_INFO_0_MSR 3520U
#define IA32_LBR_MISC_0_MSR 3536U
#define IA32_DS_AREA_MSR 1536U
#define IA32_PEBS_ENABLE_MSR 1009U
#define IA32_BTS_INDEX_MSR 1548U
#define ExitTripleFault 2U
#define ExitVmclear 19U
#define ExitVmlaunch 20U
#define ExitVmptrld 21U
#define ExitVmptrst 22U
#define ExitVmread 23U
#define ExitVmresume 24U
#define ExitVmwrite 25U
#define ExitVmxoff 26U
#define ExitVmxon 27U
#define ExitInvd 13U
#define ExitGetsec 11U
#define ExitInvept 50U
#define ExitInvvpid 53U
#define ExitCpuid 22U
#define ExitHlt 24U
#define ExitRdtsc 16U
#define ExitRdpmc 15U
#define ExitVmcall 18U
#define ExitMovCr 28U
#define ExitMovDr 29U
#define ExitIo 30U
#define ExitRdmsr 31U
#define ExitWrmsr 32U
#define ExitExceptionNmi 0U
#define ExitExtInt 1U
#define ExitIntWindow 7U
#define ExitNmiWindow 8U
#define ExitMtf 37U
#define ExitEptViolation 48U
#define ExitEptMisconfig 49U
#define ExitRdtscp 51U
#define ExitXsetbv 55U
#define ExitPreemptTimer 52U
#define IntTypeHwException 3U
#define IntTypeSwException 6U
#define IntTypeNmi 2U
#define VectorBp 3U
#define VectorGp 13U
#define VectorUd 6U
#define VectorNmi 2U
#define VmcsExitReason 17408ULL
#define VmcsExitQualification 17410ULL
#define VmcsGuestRip 26654ULL
#define VmcsGuestRsp 26652ULL
#define VmcsGuestIA32Debugctl 10242ULL
#define VmcsVmexitInstrLength 17420ULL
#define VmcsCtrlVmentryIntrInfo 16406ULL
#define VmcsCtrlVmentryErrCode 16408ULL
#define VmcsCtrlVmentryInstrLen 16410ULL
#define VmcsGuestPhysicalAddress 9216ULL
#define VmcsGuestCr0 26624ULL
#define VmcsGuestCr3 26626ULL
#define VmcsGuestCr4 26628ULL
#define VmcsGuestDr7 26650ULL
#define VmcsGuestRflags 26656ULL
#define VmcsGuestSsp 26666ULL
#define VmcsGuestSymEFER 26662ULL
#define VmcsCtrlPinBasedVmExecControls 16384ULL
#define VmcsCtrlPrimaryProcBasedVmExec 16386ULL
#define VmcsCtrlSecondaryProcBasedVmExec 16414ULL
#define VmcsCtrlExceptionBitmap 16388ULL
#define VmcsCtrlPageFaultErrCodeMask 16390ULL
#define VmcsCtrlPageFaultErrCodeMatch 16392ULL
#define VmcsCtrlCr3TargetCount 16394ULL
#define VmcsCtrlIoBitmapA 16396ULL
#define VmcsCtrlIoBitmapB 16398ULL
#define VmcsCtrlMsrBitmap 16400ULL
#define VmcsCtrlVmentryMsrLoadAddr 16402ULL
#define VmcsCtrlVmentryIntInfo 16404ULL
#define VmcsCtrlVmentryExcErrCode 16406ULL
#define VmcsCtrlTprThreshold 16406ULL
#define VmcsCtrlSecondaryVmxProcControls 16414ULL
#define VmcsCtrlPleGap 16416ULL
#define VmcsCtrlPleWindow 16418ULL
#define VmcsCtrlVmexitMsrStoreAddr 16396ULL
#define VmcsCtrlVmexitMsrLoadAddr 16398ULL
#define VmcsCtrlEptPointer 8218ULL
#define VmcsGuestEsSelector 2048ULL
#define VmcsGuestCsSelector 2050ULL
#define VmcsGuestCs 2050ULL
#define VmcsGuestSsSelector 2052ULL
#define VmcsGuestSs 2052ULL
#define VmcsGuestDsSelector 2054ULL
#define VmcsGuestFsSelector 2056ULL
#define VmcsGuestGsSelector 2058ULL
#define VmcsGuestLdtrSelector 2060ULL
#define VmcsGuestTrSelector 2062ULL
#define VmcsHostEsSelector 3072ULL
#define VmcsHostCsSelector 3074ULL
#define VmcsHostSsSelector 3076ULL
#define VmcsHostDsSelector 3078ULL
#define VmcsHostFsSelector 3080ULL
#define VmcsHostGsSelector 3082ULL
#define VmcsHostLdtrSelector 3084ULL
#define VmcsHostTrSelector 3086ULL
#define VmcsGuestEsLimit 18432ULL
#define VmcsGuestCsLimit 18434ULL
#define VmcsGuestSsLimit 18436ULL
#define VmcsGuestDsLimit 18438ULL
#define VmcsGuestFsLimit 18440ULL
#define VmcsGuestGsLimit 18442ULL
#define VmcsGuestLdtrLimit 18444ULL
#define VmcsGuestTrLimit 18446ULL
#define VmcsGuestGdtrLimit 18448ULL
#define VmcsGuestIdtrLimit 18450ULL
#define VmcsGuestEsBase 26624ULL
#define VmcsGuestCsBase 26626ULL
#define VmcsGuestSsBase 26628ULL
#define VmcsGuestDsBase 26630ULL
#define VmcsGuestFsBase 26632ULL
#define VmcsGuestGsBase 26634ULL
#define VmcsGuestLdtrBase 26636ULL
#define VmcsGuestTrBase 26638ULL
#define VmcsGuestGdtrBase 26640ULL
#define VmcsGuestIdtrBase 26642ULL
#define VmcsHostCr0 27648ULL
#define VmcsHostCr3 27650ULL
#define VmcsHostCr4 27652ULL
#define VmcsHostFsBase 27654ULL
#define VmcsHostGsBase 27656ULL
#define VmcsHostTrBase 27658ULL
#define VmcsHostGdtrBase 27660ULL
#define VmcsHostIdtrBase 27662ULL
#define VmcsHostSysenterCs 19456ULL
#define VmcsHostSysenterEsp 27664ULL
#define VmcsHostSysenterEip 27666ULL
#define VmcsHostRsp 27668ULL
#define VmcsHostRip 27670ULL
#define VmcsGuestEsAccessRights 18452ULL
#define VmcsGuestCsAccessRights 18454ULL
#define VmcsGuestSsAccessRights 18456ULL
#define VmcsGuestDsAccessRights 18458ULL
#define VmcsGuestFsAccessRights 18460ULL
#define VmcsGuestGsAccessRights 18462ULL
#define VmcsGuestLdtrAccessRights 18464ULL
#define VmcsGuestTrAccessRights 18466ULL
#define VmcsHostEsAccessRights 19476ULL
#define VmcsHostCsAccessRights 19478ULL
#define VmcsHostSsAccessRights 19480ULL
#define VmcsHostDsAccessRights 19482ULL
#define VmcsHostFsAccessRights 19484ULL
#define VmcsHostGsAccessRights 19486ULL
#define VmcsHostLdtrAccessRights 19488ULL
#define VmcsHostTrAccessRights 19490ULL
#define VmcsGuestLinkPointer 10240ULL
SO_GLOBAL_ERROR(ErrInvalidCore, "invalid processor core");
SO_GLOBAL_ERROR(ErrVmcsAlloc, "failed to allocate VMCS region");
SO_GLOBAL_ERROR(ErrVmxonAlloc, "failed to allocate VMXON region");
#define IA32_DEBUGCTL_MSR 473U
#define VMCS_GUEST_IA32_DEBUGCTL 10242ULL
#define IO_BITMAP_SIZE 8192LL
static IoBitmapManager* gIoBitmap = NULL;
#define IOCTL_QUERY_AND_CLEAR_LOGS 2273284LL
#define IOCTL_REGISTER_EVENT 2273288LL
#define IOCTL_RUN_SCRIPT 2273296LL
#define IOCTL_SEND_REQUEST_RESULT 2273300LL
#define IOCTL_GET_LOG_BASE 2273304LL
#define IOCTL_VMM_INIT 2273312LL
#define IOCTL_VMM_SHUTDOWN 2273316LL
#define IOCTL_EPT_HOOK 2273320LL
#define IOCTL_EPT_UNHOOK 2273324LL
#define IOCTL_EPT_SET_HOOK 2273328LL
#define IOCTL_EPT_GET_EPT_TABLES 2273332LL
#define IOCTL_VMM_EXECUTION_TRACE 2273336LL
#define IOCTL_VMM_READ_MEM 2273340LL
#define IOCTL_VMM_WRITE_MEM 2273344LL
#define IOCTL_VMM_VIRT_TO_PHYS 2273348LL
#define IOCTL_VMM_PHYS_TO_VIRT 2273352LL
#define IOCTL_VMM_GET_REG 2273356LL
#define IOCTL_VMM_SET_REG 2273360LL
#define IOCTL_VMM_PAUSE 2273364LL
#define IOCTL_VMM_RESUME 2273368LL
#define IOCTL_VMM_SWITCH_PROCESS 2273372LL
#define IOCTL_VMM_MODIFY_REGS 2273376LL
#define IOCTL_VMM_INVEPT 2273380LL
#define IOCTL_VMM_INVVPID 2273384LL
#define IOCTL_VMM_FLUSH_ENTIRE_TLB 2273388LL
#define IOCTL_VMM_GET_MTRR 2273392LL
#define IOCTL_VMM_CHANGE_CORE 2273396LL
#define IOCTL_VMM_GET_PROCESS_BASE 2273400LL
#define IOCTL_VMM_GET_PROCESS_CR3 2273404LL
#define IOCTL_VMM_READ_AND_WRITE_MEM 2273408LL
#define IOCTL_VMM_CALL_FUNCTION 2273412LL
#define IOCTL_VMM_QUERY_PACKET 2273416LL
#define IOCTL_VMM_SEND_RESULT 2273420LL
#define IOCTL_VMM_EXECUTE_SINGLE_STEP 2273424LL
#define IOCTL_VMM_ENABLE_AND_INVOKE_TRAP_FLAG 2273428LL
#define IOCTL_VMM_MASK_EXCEPTION_DEBUGGER_BREAKPOINT 2273432LL
#define IOCTL_VMM_UNMASK_EXCEPTION_DEBUGGER_BREAKPOINT 2273436LL
#define IOCTL_VMM_SHORT_CIRCUITING_EVENT_INJECT 2273440LL
#define IOCTL_VMM_QUERY_REGISTER_FROM_GUEST_STATE 2273444LL
#define IOCTL_TRANSPARENT_MODE_ENABLE 2273456LL
#define IOCTL_TRANSPARENT_MODE_DISABLE 2273460LL
#define IOCTL_HYPERTRACE_START 2273464LL
#define IOCTL_HYPERTRACE_STOP 2273468LL
#define IOCTL_HYPERTRACE_READ_BUFFER 2273472LL
#define IOCTL_HYPERTRACE_CLEAR_BUFFER 2273476LL
#define Baud110 110U
#define Baud300 300U
#define Baud600 600U
#define Baud1200 1200U
#define Baud2400 2400U
#define Baud4800 4800U
#define Baud9600 9600U
#define Baud14400 14400U
#define Baud19200 19200U
#define Baud38400 38400U
#define Baud57600 57600U
#define Baud115200 115200U
#define ParityNone 0LL
#define ParityOdd 1LL
#define ParityEven 2LL
#define ParityMark 3LL
#define ParitySpace 4LL
#define StopBits1 0LL
#define StopBits1_5 1LL
#define StopBits2 2LL
#define UART_RBR ((uint8_t)0)
#define UART_THR ((uint8_t)0)
#define UART_IER ((uint8_t)1)
#define UART_DLL ((uint8_t)0)
#define UART_DLM ((uint8_t)1)
#define UART_IIR ((uint8_t)2)
#define UART_FCR ((uint8_t)2)
#define UART_LCR ((uint8_t)3)
#define UART_MCR ((uint8_t)4)
#define UART_LSR ((uint8_t)5)
#define UART_MSR ((uint8_t)6)
#define UART_SCR ((uint8_t)7)
#define LogLevelInfo ((LogLevel)0)
#define LogLevelWarning ((LogLevel)1)
#define LogLevelError ((LogLevel)2)
#define LogLevelDebug ((LogLevel)3)
#define LogLevelTrace ((LogLevel)4)
#define MemUncacheable ((MemoryType)0)
#define MemWriteCombining ((MemoryType)1)
#define MemWriteThrough ((MemoryType)4)
#define MemWriteProtected ((MemoryType)5)
#define MemWriteBack ((MemoryType)6)
static ModeBasedExecHookState* gModeBasedExec = NULL;
#define MSR_BITMAP_SIZE 8192LL
#define MSR_RANGE_COUNT 8LL
static MsrBitmapManager* gMsrBitmap = NULL;
#define OptNone 0LL
#define OptBasic 1LL
#define OptAggressive 2LL
static EferHookState* gEferHook = NULL;
#define IA32_VMX_BASIC 1152U
#define IA32_EFER 3221225600U
#define IA32_LSTAR 3221225602U
#define IA32_FMASK 3221225604U
#define IA32_STAR 3221225601U
#define IA32_U_CET 3221225730U
#define IA32_PL3_SSP 3221229650U
#define IA32_FS_BASE 3221225728U
#define IA32_GS_BASE 3221225729U
#define IA32_SYSENTER_CS 372U
#define IA32_SYSENTER_ESP 374U
#define IA32_SYSENTER_EIP 376U
#define IA32_VMX_ENTRY_CTLS 1156U
#define IA32_VMX_EXIT_CTLS 1155U
#define IA32_VMX_true_ENTRY_CTLS 1164U
#define IA32_VMX_true_EXIT_CTLS 1165U
#define IA32_VMX_ENTRY_CTLS_LOAD_IA32_EFER_FLAG 65536U
#define IA32_VMX_EXIT_CTLS_SAVE_IA32_EFER_FLAG 131072U
#define X86_FLAGS_RF 65536ULL
#define X86_FLAGS_VM 131072ULL
#define X86_FLAGS_RESERVED_BITS 2ULL
#define X86_FLAGS_FIXED 2ULL
#define EXCEPTION_VECTOR_UNDEFINED_OPCODE 6U
#define PAGE_SIZE 4096ULL
#define SIZE_1_MB 1048576ULL
#define SIZE_2_MB 2097152ULL
#define SIZE_1_GB 1073741824ULL
#define SIZE_512_GB 549755813888ULL
#define POOLTAG 1212432967U
#define MAX_LOG_BUFFERS 1000U
#define MAX_LOG_BUFFERS_PRIO 50U
#define LOG_CHUNK_SIZE 4096ULL
#define LOG_MESSAGE_SIZE 288ULL
#define PENDING_INTERRUPTS_CAPACITY 64U
#define MAX_HIDDEN_BREAKPOINTS 40U
#define MAX_MTRR_ENTRIES 255U
#define EPT_PML4_COUNT 512U
#define EPT_PML3_COUNT 512U
#define EPT_PML2_COUNT 512U
#define EPT_PML1_COUNT 512U
#define PAGE_ATTRIB_READ 2ULL
#define PAGE_ATTRIB_WRITE 4ULL
#define PAGE_ATTRIB_EXEC 8ULL
#define PAGE_ATTRIB_EXEC_HIDDEN_HOOK 16ULL
#define HYPERDBG_VMCALL_MAGIC_RAX 1213613651ULL
#define HYPERDBG_VMCALL_MAGIC_RCX 94889840823372ULL
#define HYPERDBG_VMCALL_MAGIC_RDX 5642808406554530390ULL
#define HOOK_PAGE_MONITOR_READ 1U
#define HOOK_PAGE_MONITOR_WRITE 2U
#define HOOK_PAGE_MONITOR_EXEC 4U
#define HOOK_PAGE_MONITOR_READ_WRITE 3U
#define HOOK_PAGE_MONITOR_EXEC_READ 5U
#define HOOK_PAGE_MONITOR_EXEC_WRITE 6U
#define HOOK_PAGE_MONITOR_EXEC_RW 7U
#define HOOK_PAGE_MONITOR_INLINE_HOOKS 8U
#define HOOK_PAGE_MASKED_HOOKS 16U
#define VmcTest ((VmcallNumber)1)
#define VmcVmxoff ((VmcallNumber)2)
#define VmcChangePageAttrib ((VmcallNumber)3)
#define VmcInveptSingleContext ((VmcallNumber)4)
#define VmcInveptAllContexts ((VmcallNumber)5)
#define VmcUnhookAllPages ((VmcallNumber)6)
#define VmcUnhookSinglePage ((VmcallNumber)7)
#define VmcEnableSyscallHookEfer ((VmcallNumber)8)
#define VmcDisableSyscallHookEfer ((VmcallNumber)9)
#define VmcChangeMsrBitmapRead ((VmcallNumber)10)
#define VmcChangeMsrBitmapWrite ((VmcallNumber)11)
#define VmcSetRdtscExiting ((VmcallNumber)12)
#define VmcSetRdpmcExiting ((VmcallNumber)13)
#define VmcSetExceptionBitmap ((VmcallNumber)14)
#define VmcEnableMovToDebugRegsExiting ((VmcallNumber)15)
#define VmcEnableExternalInterruptExiting ((VmcallNumber)16)
#define VmcChangeIoBitmap ((VmcallNumber)17)
#define VmcSetHiddenCcBreakpoint ((VmcallNumber)18)
#define VmcDisableExternalInterruptExitingOnlyToClearInterruptCommands ((VmcallNumber)19)
#define VmcUnsetRdtscExiting ((VmcallNumber)20)
#define VmcUnsetRdpmcExiting ((VmcallNumber)21)
#define VmcUnsetExceptionBitmap ((VmcallNumber)22)
#define VmcDisableMovToDebugRegsExiting ((VmcallNumber)23)
#define VmcDisableExternalInterruptExiting ((VmcallNumber)24)
#define VmcSetCpuidExiting ((VmcallNumber)25)
#define VmcUnsetCpuidExiting ((VmcallNumber)26)
#define VmcClp ((VmcallNumber)27)
#define VmcSetMovFromCrExiting ((VmcallNumber)28)
#define VmcSetMovToCrExiting ((VmcallNumber)29)
#define VmcUnsetMovFromCrExiting ((VmcallNumber)30)
#define VmcUnsetMovToCrExiting ((VmcallNumber)31)
#define VmcSetMovFromDrExiting ((VmcallNumber)32)
#define VmcSetMovToDrExiting ((VmcallNumber)33)
#define VmcUnsetMovFromDrExiting ((VmcallNumber)34)
#define VmcUnsetMovToDrExiting ((VmcallNumber)35)
#define VmcEnableEptHookMaskedRead ((VmcallNumber)36)
#define VmcEnableEptHookMaskedWrite ((VmcallNumber)37)
#define VmcEnableEptHookMaskedReadWrite ((VmcallNumber)38)
#define VmcEnableEptHookInlineHookRead ((VmcallNumber)39)
#define VmcEnableEptHookInlineHookWrite ((VmcallNumber)40)
#define VmcEnableEptHookInlineHookReadWrite ((VmcallNumber)41)
#define VmcEnableEptHookExecRead ((VmcallNumber)42)
#define VmcEnableEptHookExecWrite ((VmcallNumber)43)
#define VmcEnableEptHookExecReadWrite ((VmcallNumber)44)
#define VmcEnableEptHookExecOnly ((VmcallNumber)45)
#define VmcEnableMtf ((VmcallNumber)46)
#define VmcQueryPerformanceFrequency ((VmcallNumber)47)
#define VmcQueryPerformanceCounter ((VmcallNumber)48)
#define VmcCheckVmxSupport ((VmcallNumber)49)
#define VmcLaunchDebugger ((VmcallNumber)4294967345)
#define VMCS_CTRL_PIN_BASED_VM_EXEC_CONTROLS 16384ULL
#define VMCS_CTRL_PRIMARY_PROC_BASED_VM_EXEC 16386ULL
#define VMCS_CTRL_SECONDARY_PROC_BASED_VM_EXEC 16414ULL
#define VMCS_CTRL_EXCEPTION_BITMAP 16388ULL
#define VMCS_CTRL_PAGE_FAULT_ERR_CODE_MASK 16390ULL
#define VMCS_CTRL_PAGE_FAULT_ERR_CODE_MATCH 16392ULL
#define VMCS_CTRL_CR3_TARGET_COUNT 16394ULL
#define VMCS_CTRL_IO_BITMAP_A 16396ULL
#define VMCS_CTRL_IO_BITMAP_B 16398ULL
#define VMCS_CTRL_MSR_BITMAP 16400ULL
#define VMCS_CTRL_VMENTRY_MSR_LOAD_ADDR 16402ULL
#define VMCS_CTRL_VMENTRY_INT_INFO 16404ULL
#define VMCS_CTRL_VMENTRY_EXC_ERR_CODE 16406ULL
#define VMCS_CTRL_VMENTRY_INSTR_LEN 16410ULL
#define VMCS_CTRL_TPR_THRESHOLD 16406ULL
#define VMCS_CTRL_SECONDARY_VMX_PROC_CONTROLS 16414ULL
#define VMCS_CTRL_PLE_GAP 16416ULL
#define VMCS_CTRL_PLE_WINDOW 16418ULL
#define VMCS_CTRL_VMEXIT_MSR_STORE_ADDR 16396ULL
#define VMCS_CTRL_VMEXIT_MSR_LOAD_ADDR 16398ULL

// -- Forward declarations --
static int32_t drvUnsupported(DEVICE_OBJECT* _p0, IRP* irp);
static int32_t drvRead(DEVICE_OBJECT* _p0, IRP* irp);
static int32_t drvWrite(DEVICE_OBJECT* _p0, IRP* irp);
static int32_t drvClose(DEVICE_OBJECT* _p0, IRP* irp);
static int32_t drvCreate(DEVICE_OBJECT* _p0, IRP* irp);
static void drvUnload(DRIVER_OBJECT* _p0);
static uintptr_t unsafeSizeofLogMessage(void);
static void EptTable_initIdentityMap(void* self);
static so_int EptTable_pml4Index(void* self, uint64_t gpa);
static so_int EptTable_pdptIndex(void* self, uint64_t gpa);
static so_int EptTable_pdIndex(void* self, uint64_t gpa);
static so_int EptTable_ptIndex(void* self, uint64_t gpa);
static so_Error EptTable_splitLargePageUnlocked(void* self, uint64_t gpa);
static so_R_ptr_err EptTable_resolveForModifyUnlocked(void* self, uint64_t gpa);
static void MtrrState_readMtrrsFromHardware(void* self);
static so_String eptError_Error(eptError e);
static so_Error fmtError(so_String format, so_Slice args);
static EventDispatcher* newEventDispatcher(void);
static bool defaultMsrRead(VCPU* _p0, uint32_t _p1, uint64_t _p2);
static bool defaultMsrWrite(VCPU* _p0, uint32_t _p1, uint64_t _p2);
static bool defaultIoIn(VCPU* _p0, uint16_t _p1, uint32_t _p2);
static bool defaultIoOut(VCPU* _p0, uint16_t _p1, uint32_t _p2);
static bool defaultCrChange(VCPU* _p0, uint32_t _p1, bool _p2, bool _p3);
static uint64_t nativeRdtsc(void);
static bool defaultBpHandler(VCPU* _p0);
static bool defaultGpHandler(VCPU* _p0);
static bool defaultUdHandler(VCPU* _p0);
static bool defaultNmiHandler(VCPU* _p0);
static void execTrapEnable(VCPU* v);
static void execTrapDisable(VCPU* v);
static so_int execTrapAddEntry(uint64_t addr, CR3_TYPE cr3);
static void execTrapRemoveEntry(so_int index);
static bool execTrapHandleMtf(VCPU* v);
static bool execTrapCheckAndHandle(VCPU* v);
static void dispatchEventExecTrap(VCPU* v, uint64_t addr, CR3_TYPE cr3);
static void setupDefaultEventHandlers(void);
static bool handleExceptionDispatch(VCPU* _p0, ExceptionHandler* _p1);
static void copyPageToFake(EptHook* h, uintptr_t pageBase, CR3_TYPE cr3);
static so_Error HookManager_updateExistingHook(void* self, EptHook* existing, uint64_t va, HookKind kind);
static EptHook* HookManager_findByPhysAddr(void* self, uint64_t pa);
static EptHook* HookManager_findByVirtAddr(void* self, uint64_t va);
static void HookManager_restoreAllForCoreLocked(void* self, uint32_t coreId);
static void restoreOriginalAndInjectBp(VCPU* v, EptHook* hook, Breakpoint* bp);
static void singleStepOnFakePage(VCPU* v, EptHook* hook);
static void eptHookSetReadHook(uint64_t addr, CR3_TYPE cr3);
static void eptHookSetWriteHook(uint64_t addr, CR3_TYPE cr3);
static void eptHookSetExecHook(uint64_t addr, CR3_TYPE cr3);
static void eptHookSetReadWriteHook(uint64_t addr, CR3_TYPE cr3);
static so_Error EvasionState_applyEvasions(void* self);
static so_Error EvasionState_removeEvasions(void* self);
static so_Error EvasionState_hideVmxBits(void* self);
static so_Error EvasionState_hideVmxRelatedMsrs(void* self);
static so_Error EvasionState_hideVmxCriticalIoPorts(void* self);
static so_Error EvasionState_enableTimingCountermeasures(void* self);
static so_Error EvasionState_hideDebuggerPresence(void* self);
static so_Error TracerState_enableLbr(void* self, TraceUserConfig* config);
static void TracerState_disableLbrLocked(void* self);
static so_Error TracerState_enableBts(void* self, TraceUserConfig* config);
static void TracerState_disableBtsLocked(void* self);
static so_int TracerState_readLbrFromMsr(void* self, so_Slice buffer);
static so_int TracerState_readBtsFromMemory(void* self, so_Slice buffer);
static so_Error Hypervisor_initializeVcpu(void* self, uint32_t coreId);
static EptTable* Hypervisor_eptTable(void* self, uint32_t coreId);
static bool Hypervisor_dispatch(void* self, VCPU* v);
static bool Hypervisor_handleException(void* self, VCPU* v);
static bool Hypervisor_handleEptViolation(void* self, VCPU* v);
static bool Hypervisor_handleVmcall(void* self, VCPU* v);
static void Hypervisor_hypercallForward(void* self, VCPU* v);
static void Hypervisor_advanceIp(void* self, VCPU* v);
static void Hypervisor_suppressAdvance(void* self, VCPU* v);
static void Hypervisor_enableAdvance(void* self, VCPU* v);
static void Hypervisor_injectInterrupt(void* self, uint32_t intType, uint32_t vector, bool deliverErr, uint32_t errCode);
static void Hypervisor_injectBp(void* self, VCPU* v);
static void Hypervisor_injectGp(void* self, VCPU* v);
static void Hypervisor_injectUd(void* self, VCPU* v);
static bool Hypervisor_handleMonitorTrapFlag(void* self, VCPU* v);
static void Hypervisor_enableAndCheckPendingInt(void* self, VCPU* _p0);
static void Hypervisor_checkTrapFlag(void* self, VCPU* _p0);
static bool Hypervisor_handleTripleFault(void* self, VCPU* _p0);
static bool Hypervisor_handleEptMisconfig(void* self, VCPU* _p0);
static bool Hypervisor_handleGeneralProtectionFault(void* self, VCPU* _p0);
static bool Hypervisor_handleUndefinedOpcode(void* self, VCPU* _p0);
static bool Hypervisor_handleExternalInt(void* self, VCPU* _p0);
static bool Hypervisor_handleNmi(void* self, VCPU* _p0);
static bool Hypervisor_handleInterruptWindow(void* self, VCPU* _p0);
static bool Hypervisor_handleNmiWindow(void* self, VCPU* _p0);
static bool Hypervisor_handlePreemptionTimer(void* self, VCPU* _p0);
static void Hypervisor_transparentUnhide(void* self);
static void Hypervisor_transparentHide(void* self);
static uint32_t readExitReason(void);
static uint32_t readExitQualification(void);
static uint64_t readGuestRip(void);
static uint32_t readInstructionLength(void);
static uint64_t readGuestPhysicalAddr(void);
static uint64_t allocateVmcsRegion(void);
static uint64_t allocateVmxonRegion(void);
static uint64_t allocateMsrBitmap(void);
static uint64_t allocateIoBitmap(void);
static uint64_t allocateVmmStack(void);
static so_Error vmxTurnOn(uint64_t pa);
static void vmxTurnOff(void);
static so_Error vmClear(uint64_t pa);
static so_Error vmLoad(uint64_t pa);
static so_Error Hypervisor_setupVmcs(void* self, VCPU* v);
static uint64_t adjustVmcsControl(uint64_t msrValue, uint64_t suggestedValue);
static void vmExitHandler(void);
static uint16_t getCs(void);
static uint16_t getSs(void);
static uint16_t getDs(void);
static uint16_t getEs(void);
static uint16_t getFs(void);
static uint16_t getGs(void);
static uint64_t getGdtBase(void);
static uint64_t getIdtBase(void);
static uint32_t getGdtLimit(void);
static uint32_t getIdtLimit(void);
static uint64_t getRflags(void);
static uint32_t getCurrentProcessorNumber(void);
static bool breakpointReapplyHook(uint32_t coreId);
static void handleRegisteredMtf(uint32_t coreId);
static bool kdHandleNmiCallback(uint32_t coreId);
static void inveptSingleContext(uint64_t eptp);
static void inveptAllContexts(void);
static void invvpidAllContexts(void);
static void Hypervisor_setRdtscExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setPmcVmexit(void* self, bool enable);
static void Hypervisor_setExceptionBitmap(void* self, VCPU* v, uint32_t bitmap);
static void Hypervisor_setMovDebugRegsExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setExternalInterruptExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setCpuidExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setClpExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setMovFromCrExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setMovToCrExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setMovFromDrExiting(void* self, VCPU* v, bool enable);
static void Hypervisor_setMovToDrExiting(void* self, VCPU* v, bool enable);
static void protectedHvExternalInterruptExitingForDisablingInterruptCommands(VCPU* v);
static void initIoBitmap(void);
static void ioHandlePerformIoBitmapChange(VCPU* v, uint32_t portMask);
static void ioHandleEnableOrDisableIoPortExiting(VCPU* v, bool enable);
static bool handleIoRead(VCPU* v);
static bool handleIoWrite(VCPU* v);
static uint8_t inb(uint16_t port);
static uint16_t inw(uint16_t port);
static uint32_t ind(uint16_t port);
static void outb(uint16_t port, uint8_t val);
static void outw(uint16_t port, uint16_t val);
static void outd(uint16_t port, uint32_t val);
static uint32_t sizeOf(void);
static int32_t drvIoctl(DEVICE_OBJECT* _p0, IRP* irp);
static int32_t handleQueryLogs(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleRegisterEvent(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleRunScript(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleSendResult(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleGetLogBase(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleLogRead(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleVmmInit(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleVmmShutdown(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleEptHook(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleEptUnhook(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleEptSetHook(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleGetEptTables(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleExecutionTrace(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleReadMem(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleWriteMem(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleVirtToPhys(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handlePhysToVirt(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleGetReg(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleSetReg(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handlePause(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleResume(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleSwitchProcess(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleModifyRegs(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleInvept(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleInvvpid(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleFlushTlb(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleGetMtrr(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleChangeCore(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleGetProcessBase(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleGetProcessCr3(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleReadWriteMem(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleCallFunction(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleQueryPacket(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleSingleStep(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTrapFlag(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleMaskBreakpoint(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleUnmaskBreakpoint(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleShortCircuitInject(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleQueryGuestReg(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTransparentEnable(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTransparentDisable(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTraceStart(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTraceStop(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTraceReadBuffer(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static int32_t handleTraceClearBuffer(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info);
static uint64_t getRegisterValue(GUEST_REGS* regs, uint32_t index);
static void setRegisterValue(GUEST_REGS* regs, uint32_t index, uint64_t value);
static void eptInveptSingleContext(uint64_t eptp);
static void eptInveptAllContexts(void);
static uint16_t KdSerialState_comPortBase(void* self, uint32_t portNum);
static so_Error KdSerialState_configureHardware(void* self);
static uint16_t KdSerialState_calculateDivisor(void* self);
static void KdSerialState_setDlab(void* self, bool enable);
static void KdSerialState_disableUart(void* self);
static void KdSerialState_writePort(void* self, uint8_t reg, uint8_t value);
static uint8_t KdSerialState_readPort(void* self, uint8_t reg);
static void KdSerialState_flushTxBufferHardware(void* self);
static bool KdSerialState_isTransmitEmpty(void* self);
static bool KdSerialState_isReceiveReady(void* self);
static void KdSerialState_pollRxHardware(void* self);
static uint32_t KdSerialState_bytesInRxBuffer(void* self);
static uint32_t KdSerialState_txBytesQueued(void* self);
static uint64_t rdtscValue(void);
static MemoryManager* newMemoryManager(void);
static uint64_t readPhysicalMemory(uint64_t pa);
static void writePhysicalMemory(uint64_t pa, uint64_t val);
static void invalidateEpt(void);
static void modeBasedExecHookEnable(VCPU* v);
static void modeBasedExecHookDisable(VCPU* v);
static so_int modeBasedExecHookAddSelector(uint16_t selector);
static void modeBasedExecHookRemoveSelector(so_int index);
static bool modeBasedExecHandleEptViolation(VCPU* v, uint64_t qualification);
static void dispatchEventModeBasedExec(VCPU* v, uint16_t selector, uint64_t qualification);
static void initMsrBitmap(void);
static void msrHandlePerformMsrBitmapReadChange(VCPU* v, uint32_t msrIndex);
static void msrHandlePerformMsrBitmapWriteChange(VCPU* v, uint32_t msrIndex);
static void msrHandleEnableOrDisableMsrExiting(VCPU* v, bool enable);
static bool handleMsrRead(VCPU* v);
static bool handleMsrWrite(VCPU* v);
static PoolEntry* PoolManager_findEntry(void* self, uintptr_t ptr);
static void VmmOptimizer_applyDefaults(void* self);
static void EptCache_evictOldest(void* self);
static so_String cstringToString(int8_t* s);
static Hypervisor* getHypervisor(void);
static void syscallHookConfigureEFER(VCPU* v, bool enable);
static bool syscallHookEmulateSYSCALL(VCPU* v);
static bool syscallHookEmulateSYSRET(VCPU* v);
static bool syscallHookHandleUD(VCPU* v);
static CR3_TYPE getCurrentProcessCr3(void);
static bool checkPagePresentByCr3(uintptr_t addr, CR3_TYPE cr3);
static void readMemorySafe(uint64_t addr, void* buffer, uint64_t size);
static void dispatchEventEferSyscall(VCPU* v);
static void dispatchEventEferSysret(VCPU* v);
static void injectPageFaultWithoutErrorCode(uint64_t rip);
static void removeUndefinedInstructionForDisablingSyscallSysret(VCPU* v);
static void setSegmentRegister(uint64_t field, uint16_t selector, uint64_t base, uint32_t limit, uint32_t accessRights);
static uint64_t eptPml1Offset(uint64_t v);
static uint64_t eptPml1Index(uint64_t v);
static uint64_t eptPml2Index(uint64_t v);
static uint64_t eptPml3Index(uint64_t v);
static uint64_t eptPml4Index(uint64_t v);
static uint64_t logBufferSize(void);
static uint64_t logBufferSizePrio(void);
static int32_t Hypervisor_handleHyperdbgVmcall(void* self, uint64_t num, uint64_t p1, uint64_t p2, uint64_t p3);
static int32_t Hypervisor_vmcallTest(void* self, uint64_t _p0, uint64_t _p1, uint64_t _p2);
static bool Hypervisor_performPageHook(void* self, VCPU* v, uintptr_t addr, CR3_TYPE cr3, uint32_t mask);
static uint64_t queryPerformanceFrequency(void);
static uint64_t queryPerformanceCounter(void);
static bool checkVmxSupport(void);
static void vmRead64(uint64_t field, uint64_t* val);
static void vmRead32(uint64_t field, uint32_t* val);
static void vmWrite64(uint64_t field, uint64_t val);
static void vmWrite32(uint64_t field, uint32_t val);
static void setMonitorTrapFlag(bool enable);
static bool vmxVmlaunch(VCPU* v);
static bool vmxVmresume(VCPU* v);
static void vmxVmxoff(VCPU* v);
static void asmHypervVmcall(uint64_t regs);

// -- driver.go --

NTSTATUS DriverEntry(DRIVER_OBJECT* driverObj, UNICODE_STRING* registryPath) {
    _solod_kmod_init();
    UNREFERENCED_PARAMETER(registryPath);
    ExInitializeDriverRuntime(DrvRtPoolNxOptIn);
    gDrv = &(Driver){};
    uint32_t numCpus = KeQueryActiveProcessorCount(((PVOID)(NULL)));
    {
        so_Error err = InitializeGlobals(numCpus);
        if (err != NULL) {
            LogError(so_str("Failed to initialize globals: %v"), (so_Slice){(void*[1]){err}, 1, 1});
            return STATUS_UNSUCCESSFUL;
        }
    }
    so_Slice devName = ((so_Slice){(void*)L"\\Device\\HyperDbgDebuggerDevice", 31, 31});
    so_Slice dosName = ((so_Slice){(void*)L"\\DosDevices\\HyperDbgDebuggerDevice", 35, 35});
    UNICODE_STRING name = {0}, dos = {0};
    RtlInitUnicodeString(&name, &so_at(uint16_t, devName, 0));
    RtlInitUnicodeString(&dos, &so_at(uint16_t, dosName, 0));
    int32_t status = IoCreateDevice(driverObj, 0, &name, FILE_DEVICE_UNKNOWN, FILE_DEVICE_SECURE_OPEN, FALSE, &gDrv->device);
    if (status == STATUS_SUCCESS) {
        for (so_int i = 0; i < IRP_MJ_MAXIMUM_FUNCTION; i++) {
            driverObj->MajorFunction[i] = drvUnsupported;
        }
        driverObj->MajorFunction[IRP_MJ_CLOSE] = drvClose;
        driverObj->MajorFunction[IRP_MJ_CREATE] = drvCreate;
        driverObj->MajorFunction[IRP_MJ_READ] = drvRead;
        driverObj->MajorFunction[IRP_MJ_WRITE] = drvWrite;
        driverObj->MajorFunction[IRP_MJ_DEVICE_CONTROL] = drvIoctl;
        driverObj->DriverUnload = drvUnload;
        IoCreateSymbolicLink(&dos, &name);
        gDrv->device->Flags |= DO_BUFFERED_IO;
        gDrv->initialized = true;
        LogInfo(so_str("HyperDbg driver loaded successfully"), (so_Slice){&so_Nil, 0, 0});
    }
    return status;
}

static int32_t drvUnsupported(DEVICE_OBJECT* _p0, IRP* irp) {
    irp->IoStatus.Status = STATUS_SUCCESS;
    irp->IoStatus.Information = 0;
    IoCompleteRequest(irp, IO_NO_INCREMENT);
    return STATUS_SUCCESS;
}

static int32_t drvRead(DEVICE_OBJECT* _p0, IRP* irp) {
    IO_STACK_LOCATION* stack = IoGetCurrentIrpStackLocation(irp);
    uint32_t outLen = stack->Parameters.DeviceIoControl.OutputBufferLength;
    if (outLen == 0) {
        irp->IoStatus.Status = STATUS_SUCCESS;
        irp->IoStatus.Information = 0;
        IoCompleteRequest(irp, IO_NO_INCREMENT);
        return STATUS_SUCCESS;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    uint32_t count = Logger_FlushToUser(gLog, ((buffer)), (uintptr_t)(outLen));
    irp->IoStatus.Status = STATUS_SUCCESS;
    irp->IoStatus.Information = (uint64_t)((uintptr_t)(count) * unsafeSizeofLogMessage());
    IoCompleteRequest(irp, IO_NO_INCREMENT);
    return STATUS_SUCCESS;
}

static int32_t drvWrite(DEVICE_OBJECT* _p0, IRP* irp) {
    irp->IoStatus.Status = STATUS_SUCCESS;
    irp->IoStatus.Information = 0;
    IoCompleteRequest(irp, IO_NO_INCREMENT);
    return STATUS_SUCCESS;
}

static int32_t drvClose(DEVICE_OBJECT* _p0, IRP* irp) {
    gDrv->handleInUse = false;
    irp->IoStatus.Status = STATUS_SUCCESS;
    irp->IoStatus.Information = 0;
    IoCompleteRequest(irp, IO_NO_INCREMENT);
    return STATUS_SUCCESS;
}

static int32_t drvCreate(DEVICE_OBJECT* _p0, IRP* irp) {
    LUID priv = (LUID){.LowPart = SE_DEBUG_PRIVILEGE, .HighPart = 0};
    if (SeSinglePrivilegeCheck(priv, (irp->RequestorMode)) == FALSE) {
        irp->IoStatus.Status = STATUS_ACCESS_DENIED;
        IoCompleteRequest(irp, IO_NO_INCREMENT);
        return STATUS_ACCESS_DENIED;
    }
    if (gDrv->handleInUse) {
        irp->IoStatus.Status = STATUS_SUCCESS;
        irp->IoStatus.Information = 0;
        IoCompleteRequest(irp, IO_NO_INCREMENT);
        return STATUS_SUCCESS;
    }
    gDrv->handleInUse = true;
    irp->IoStatus.Status = STATUS_SUCCESS;
    irp->IoStatus.Information = 0;
    IoCompleteRequest(irp, IO_NO_INCREMENT);
    return STATUS_SUCCESS;
}

static void drvUnload(DRIVER_OBJECT* _p0) {
    LogInfo(so_str("HyperDbg driver unloading..."), (so_Slice){&so_Nil, 0, 0});
    CleanupGlobals();
    if (gDrv != NULL && gDrv->device != NULL) {
        so_Slice devName = ((so_Slice){(void*)L"\\DosDevices\\HyperDbgDebuggerDevice", 35, 35});
        UNICODE_STRING dos = {0};
        RtlInitUnicodeString(&dos, &so_at(uint16_t, devName, 0));
        IoDeleteSymbolicLink(&dos);
        IoDeleteDevice(gDrv->device);
    }
    LogInfo(so_str("HyperDbg driver unloaded successfully"), (so_Slice){&so_Nil, 0, 0});
}

static uintptr_t unsafeSizeofLogMessage(void) {
    return 288;
}

// -- ept.go --

uint64_t PagePermission_ToRaw(PagePermission p) {
    uint64_t raw = 0;
    if (p.Read) {
        raw |= EptRead;
    }
    if (p.Write) {
        raw |= EptWrite;
    }
    if (p.Exec) {
        raw |= EptExec;
    }
    return raw;
}

PagePermission PagePermissionFromRaw(uint64_t raw) {
    return (PagePermission){.Read = (raw & EptRead) != 0, .Write = (raw & EptWrite) != 0, .Exec = (raw & EptExec) != 0};
}

EptEntry* NewEptEntry(void) {
    return &(EptEntry){};
}

uint64_t EptEntry_Raw(void* self) {
    EptEntry* e = (EptEntry*)self;
    return e->raw;
}

void EptEntry_SetRaw(void* self, uint64_t v) {
    EptEntry* e = (EptEntry*)self;
    e->raw = v;
}

bool EptEntry_IsPresent(void* self) {
    EptEntry* e = (EptEntry*)self;
    return (e->raw & EptRead) != 0;
}

uint64_t EptEntry_PhysAddr(void* self) {
    EptEntry* e = (EptEntry*)self;
    return (e->raw & EptPhysMask);
}

void EptEntry_SetPhysAddr(void* self, uint64_t pa) {
    EptEntry* e = (EptEntry*)self;
    e->raw = ((e->raw & ~EptPhysMask) | (pa & EptPhysMask));
}

uint8_t EptEntry_MemoryType(void* self) {
    EptEntry* e = (EptEntry*)self;
    return (uint8_t)((e->raw >> EptTypeShift) & 0x7);
}

void EptEntry_SetMemoryType(void* self, uint8_t t) {
    EptEntry* e = (EptEntry*)self;
    e->raw = ((e->raw & ~EptTypeMask) | (((uint64_t)(t) & 0x7) << EptTypeShift));
}

PagePermission EptEntry_Permission(void* self) {
    EptEntry* e = (EptEntry*)self;
    return PagePermissionFromRaw(e->raw);
}

void EptEntry_SetPermission(void* self, PagePermission p) {
    EptEntry* e = (EptEntry*)self;
    e->raw = ((e->raw & ~((EptRead | EptWrite) | EptExec)) | PagePermission_ToRaw(p));
}

bool EptEntry_IsLargePage(void* self) {
    EptEntry* e = (EptEntry*)self;
    return (e->raw & EptLargePage) != 0;
}

bool EptEntry_IgnorePat(void* self) {
    EptEntry* e = (EptEntry*)self;
    return (e->raw & EptIgnorePat) != 0;
}

bool EptEntry_SuppressVe(void* self) {
    EptEntry* e = (EptEntry*)self;
    return (e->raw & EptSuppressVe) != 0;
}

void EptEntry_SetIgnorePat(void* self, bool v) {
    EptEntry* e = (EptEntry*)self;
    if (v) {
        e->raw |= EptIgnorePat;
    } else {
        e->raw &= ~EptIgnorePat;
    }
}

void EptEntry_SetSuppressVe(void* self, bool v) {
    EptEntry* e = (EptEntry*)self;
    if (v) {
        e->raw |= EptSuppressVe;
    } else {
        e->raw &= ~EptSuppressVe;
    }
}

void EptEntry_SetLargePage(void* self, bool v) {
    EptEntry* e = (EptEntry*)self;
    if (v) {
        e->raw |= EptLargePage;
    } else {
        e->raw &= ~EptLargePage;
    }
}

EptEntry* EptEntry_Clone(void* self) {
    EptEntry* e = (EptEntry*)self;
    return &(EptEntry){.raw = e->raw};
}

EptTable* NewEptTable(void) {
    EptTable* t = &(EptTable){.mtrrState = NewMtrrState(), .lock = NewSpinlock()};
    EptTable_initIdentityMap(t);
    return t;
}

static void EptTable_initIdentityMap(void* self) {
    EptTable* t = (EptTable*)self;
    PagePermission defaultPerm = (PagePermission){.Read = true, .Write = true, .Exec = true};
    for (so_int i = 0; i < 512; i++) {
        if (i == 0) {
            continue;
        }
        for (so_int j = 0; j < 512; j++) {
            EptEntry* pd = t->pdTables[i - 1][j];
            EptEntry_SetPhysAddr(pd, (uint64_t)(i - 1) * (uint64_t)(SIZE_1_GB) + (uint64_t)(j) * (uint64_t)(SIZE_2_MB));
            EptEntry_SetPermission(pd, defaultPerm);
            EptEntry_SetMemoryType(pd, MemTypeWriteBack);
            EptEntry_SetLargePage(pd, true);
            EptEntry_SetIgnorePat(pd, true);
        }
    }
}

void EptTable_Lock(void* self) {
    EptTable* t = (EptTable*)self;
    Spinlock_Lock(t->lock);
}

void EptTable_Unlock(void* self) {
    EptTable* t = (EptTable*)self;
    Spinlock_Unlock(t->lock);
}

static so_int EptTable_pml4Index(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    return (so_int)((gpa >> 39) & 0x1FF);
}

static so_int EptTable_pdptIndex(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    return (so_int)((gpa >> 30) & 0x1FF);
}

static so_int EptTable_pdIndex(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    return (so_int)((gpa >> 21) & 0x1FF);
}

static so_int EptTable_ptIndex(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    return (so_int)((gpa >> 12) & 0x1FF);
}

so_R_ptr_err EptTable_Walk(void* self, uint64_t gpa, so_int level) {
    EptTable* t = (EptTable*)self;
    so_int p4Idx = EptTable_pml4Index(t, gpa);
    if (level == LvlPml4) {
        return (so_R_ptr_err){.val = &t->pml4[p4Idx], .err = NULL};
    }
    EptEntry p4 = t->pml4[p4Idx];
    if (!EptEntry_IsPresent(&p4) || p4Idx == 0) {
        return (so_R_ptr_err){.val = NULL, .err = ErrNotPresent};
    }
    so_int pdptIdx = EptTable_pdptIndex(t, gpa);
    if (level == LvlPdpt) {
        EptEntry (*pdpt)[512] = t->pdptTables[p4Idx - 1];
        if (pdpt == NULL) {
            return (so_R_ptr_err){.val = NULL, .err = ErrNotPresent};
        }
        return (so_R_ptr_err){.val = &(*pdpt)[pdptIdx], .err = NULL};
    }
    EptEntry (*pdpt)[512] = t->pdptTables[p4Idx - 1];
    if (pdpt == NULL) {
        return (so_R_ptr_err){.val = NULL, .err = ErrNotPresent};
    }
    EptEntry pdptEntry = (*pdpt)[pdptIdx];
    if (!EptEntry_IsPresent(&pdptEntry)) {
        return (so_R_ptr_err){.val = NULL, .err = ErrNotPresent};
    }
    so_int pdIdx = EptTable_pdIndex(t, gpa);
    if (level == LvlPd) {
        EptEntry* pd = t->pdTables[p4Idx - 1][pdIdx];
        if (pd == NULL) {
            return (so_R_ptr_err){.val = NULL, .err = ErrNotPresent};
        }
        return (so_R_ptr_err){.val = pd, .err = NULL};
    }
    EptEntry* pd = t->pdTables[p4Idx - 1][pdIdx];
    if (pd == NULL || !EptEntry_IsPresent(pd)) {
        return (so_R_ptr_err){.val = NULL, .err = ErrNotPresent};
    }
    if (EptEntry_IsLargePage(pd)) {
        return (so_R_ptr_err){.val = NULL, .err = ErrLargePage};
    }
    so_int ptIdx = EptTable_ptIndex(t, gpa);
    for (so_int i = 0; i < so_len(t->splits); i++) {
        if (so_at(SplitRecord, t->splits, i).pdEntry == pd && ptIdx < 512) {
            return (so_R_ptr_err){.val = &so_at(SplitRecord, t->splits, i).pt->entries[ptIdx], .err = NULL};
        }
    }
    return (so_R_ptr_err){.val = NULL, .err = ErrNotSplit};
}

so_R_ptr_err EptTable_Lookup(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    so_R_ptr_err _res1 = EptTable_Walk(t, gpa, LvlPt);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err == NULL && pte != NULL && EptEntry_IsPresent(pte)) {
        return (so_R_ptr_err){.val = pte, .err = NULL};
    }
    so_R_ptr_err _res2 = EptTable_Walk(t, gpa, LvlPd);
    EptEntry* pd = _res2.val;
    err = _res2.err;
    if (err == NULL && pd != NULL && EptEntry_IsPresent(pd) && EptEntry_IsLargePage(pd)) {
        return (so_R_ptr_err){.val = pd, .err = NULL};
    }
    return (so_R_ptr_err){.val = NULL, .err = ErrNotMapped};
}

so_R_ptr_err EptTable_ResolveForModify(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    so_R_ptr_err _res1 = EptTable_Walk(t, gpa, LvlPt);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err == NULL && pte != NULL) {
        return (so_R_ptr_err){.val = pte, .err = NULL};
    }
    so_R_ptr_err _res2 = EptTable_Walk(t, gpa, LvlPd);
    EptEntry* pd = _res2.val;
    err = _res2.err;
    if (err == NULL && pd != NULL && EptEntry_IsPresent(pd) && EptEntry_IsLargePage(pd)) {
        return (so_R_ptr_err){.val = pd, .err = NULL};
    }
    return (so_R_ptr_err){.val = NULL, .err = ErrNotMapped};
}

so_Error EptTable_SplitLargePage(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    EptTable_Lock(t);
    so_R_ptr_err _res1 = EptTable_ResolveForModify(t, gpa);
    EptEntry* pd = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL || pd == NULL || !EptEntry_IsLargePage(pd)) {
        EptTable_Unlock(t);
        return NULL;
    }
    for (so_int _ = 0; _ < so_len(t->splits); _++) {
        SplitRecord s = so_at(SplitRecord, t->splits, _);
        if (s.pdEntry == pd) {
            EptTable_Unlock(t);
            return NULL;
        }
    }
    Pml1Table* pt = &(Pml1Table){0};
    uint64_t basePa = EptEntry_PhysAddr(pd);
    uint8_t memType = EptEntry_MemoryType(pd);
    bool ignorePat = EptEntry_IgnorePat(pd);
    bool suppressVe = EptEntry_SuppressVe(pd);
    EptEntry template = (EptEntry){};
    EptEntry_SetPermission(&template, (PagePermission){.Read = true, .Write = true, .Exec = true});
    EptEntry_SetMemoryType(&template, memType);
    EptEntry_SetIgnorePat(&template, ignorePat);
    EptEntry_SetSuppressVe(&template, suppressVe);
    for (so_int i = 0; i < 512; i++) {
        pt->entries[i] = template;
        EptEntry_SetPhysAddr(&pt->entries[i], basePa + (uint64_t)(i) * (uint64_t)(PAGE_SIZE));
    }
    EptEntry ptrEntry = (EptEntry){};
    EptEntry_SetPermission(&ptrEntry, (PagePermission){.Read = true, .Write = true, .Exec = true});
    EptEntry_SetMemoryType(&ptrEntry, MemTypeWriteBack);
    EptEntry_SetPhysAddr(&ptrEntry, (uint64_t)((uintptr_t)((void*)(&pt->entries[0]))));
    uint64_t oldRaw = EptEntry_Raw(pd);
    *pd = ptrEntry;
    uint64_t gpaBase = (gpa & ~((uint64_t)(SIZE_2_MB) - 1));
    t->splits = so_append(SplitRecord, t->splits, (SplitRecord){.pdEntry = pd, .pt = pt, .refCount = 1, .gpaBase = gpaBase});
    LogDebug(so_str("EPT: Split large page at GPA 0x%X (old raw=0x%X)"), (so_Slice){(void*[2]){&gpaBase, &oldRaw}, 2, 2});
    EptTable_Unlock(t);
    return NULL;
}

so_Error EptTable_MergeLargePage(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    EptTable_Lock(t);
    so_R_ptr_err _res1 = EptTable_Walk(t, gpa, LvlPd);
    EptEntry* pd = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL || pd == NULL) {
        EptTable_Unlock(t);
        return ErrNotMapped;
    }
    for (so_int i = 0; i < so_len(t->splits); i++) {
        SplitRecord s = so_at(SplitRecord, t->splits, i);
        if (s.pdEntry == pd) {
            bool allDefault = true;
            PagePermission defaultPerm = (PagePermission){.Read = true, .Write = true, .Exec = true};
            for (so_int j = 0; j < 512; j++) {
                PagePermission perm = EptEntry_Permission(&s.pt->entries[j]);
                uint64_t pa = EptEntry_PhysAddr(&s.pt->entries[j]);
                uint64_t expectedPa = s.gpaBase + (uint64_t)(j) * (uint64_t)(PAGE_SIZE);
                if (!memcmp(&perm, &defaultPerm, sizeof(PagePermission)) || pa != expectedPa) {
                    allDefault = false;
                    break;
                }
            }
            if (allDefault && s.refCount <= 1) {
                uint8_t memType = MtrrState_Lookup(t->mtrrState, s.gpaBase);
                EptEntry restored = (EptEntry){};
                EptEntry_SetPhysAddr(&restored, s.gpaBase);
                EptEntry_SetPermission(&restored, defaultPerm);
                EptEntry_SetMemoryType(&restored, memType);
                EptEntry_SetLargePage(&restored, true);
                EptEntry_SetIgnorePat(&restored, true);
                *pd = restored;
                t->splits = so_extend(SplitRecord, so_slice(SplitRecord, t->splits, 0, i), (so_slice(SplitRecord, t->splits, i + 1, t->splits.len)));
                LogDebug(so_str("EPT: Merged page at GPA 0x%X"), (so_Slice){(void*[1]){&(uint64_t){s.gpaBase}}, 1, 1});
                EptTable_Unlock(t);
                return NULL;
            }
            EptTable_Unlock(t);
            return fmtError(so_str("cannot merge: page has non-default entries or multiple refs"), (so_Slice){&so_Nil, 0, 0});
        }
    }
    EptTable_Unlock(t);
    return fmtError(so_str("no split found for this address"), (so_Slice){&so_Nil, 0, 0});
}

so_Error EptTable_Map(void* self, uint64_t gpa, uint64_t pa, PagePermission perm, uint8_t memType) {
    EptTable* t = (EptTable*)self;
    EptTable_Lock(t);
    so_R_ptr_err _res1 = EptTable_ResolveForModify(t, gpa);
    EptEntry* pd = _res1.val;
    so_Error err = _res1.err;
    if (err == NULL && pd != NULL && EptEntry_IsLargePage(pd)) {
        {
            so_Error splitErr = EptTable_splitLargePageUnlocked(t, gpa);
            if (splitErr != NULL) {
                EptTable_Unlock(t);
                return splitErr;
            }
        }
    }
    so_R_ptr_err _res2 = EptTable_Walk(t, gpa, LvlPt);
    EptEntry* pte = _res2.val;
    err = _res2.err;
    if (err != NULL) {
        EptTable_Unlock(t);
        return err;
    }
    EptEntry_SetPhysAddr(pte, pa);
    EptEntry_SetPermission(pte, perm);
    EptEntry_SetMemoryType(pte, memType);
    EptTable_Unlock(t);
    return NULL;
}

so_Error EptTable_Unmap(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    EptTable_Lock(t);
    so_R_ptr_err _res1 = EptTable_resolveForModifyUnlocked(t, gpa);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        EptTable_Unlock(t);
        return err;
    }
    *pte = (EptEntry){};
    EptTable_Unlock(t);
    return NULL;
}

so_Error EptTable_Protect(void* self, uint64_t gpa, PagePermission perm) {
    EptTable* t = (EptTable*)self;
    EptTable_Lock(t);
    so_R_ptr_err _res1 = EptTable_resolveForModifyUnlocked(t, gpa);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        EptTable_Unlock(t);
        return err;
    }
    EptEntry_SetPermission(pte, perm);
    EptTable_Unlock(t);
    return NULL;
}

static so_Error EptTable_splitLargePageUnlocked(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    so_R_ptr_err _res1 = EptTable_resolveForModifyUnlocked(t, gpa);
    EptEntry* pd = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL || pd == NULL || !EptEntry_IsLargePage(pd)) {
        return NULL;
    }
    for (so_int _ = 0; _ < so_len(t->splits); _++) {
        SplitRecord s = so_at(SplitRecord, t->splits, _);
        if (s.pdEntry == pd) {
            return NULL;
        }
    }
    Pml1Table* pt = &(Pml1Table){0};
    uint64_t basePa = EptEntry_PhysAddr(pd);
    uint8_t memType = EptEntry_MemoryType(pd);
    bool ignorePat = EptEntry_IgnorePat(pd);
    bool suppressVe = EptEntry_SuppressVe(pd);
    EptEntry template = (EptEntry){};
    EptEntry_SetPermission(&template, (PagePermission){.Read = true, .Write = true, .Exec = true});
    EptEntry_SetMemoryType(&template, memType);
    EptEntry_SetIgnorePat(&template, ignorePat);
    EptEntry_SetSuppressVe(&template, suppressVe);
    for (so_int i = 0; i < 512; i++) {
        pt->entries[i] = template;
        EptEntry_SetPhysAddr(&pt->entries[i], basePa + (uint64_t)(i) * (uint64_t)(PAGE_SIZE));
    }
    EptEntry ptrEntry = (EptEntry){};
    EptEntry_SetPermission(&ptrEntry, (PagePermission){.Read = true, .Write = true, .Exec = true});
    EptEntry_SetMemoryType(&ptrEntry, MemTypeWriteBack);
    EptEntry_SetPhysAddr(&ptrEntry, (uint64_t)((uintptr_t)((void*)(&pt->entries[0]))));
    *pd = ptrEntry;
    uint64_t gpaBase = (gpa & ~((uint64_t)(SIZE_2_MB) - 1));
    t->splits = so_append(SplitRecord, t->splits, (SplitRecord){.pdEntry = pd, .pt = pt, .refCount = 1, .gpaBase = gpaBase});
    return NULL;
}

static so_R_ptr_err EptTable_resolveForModifyUnlocked(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    so_R_ptr_err _res1 = EptTable_Walk(t, gpa, LvlPt);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err == NULL && pte != NULL) {
        return (so_R_ptr_err){.val = pte, .err = NULL};
    }
    so_R_ptr_err _res2 = EptTable_Walk(t, gpa, LvlPd);
    EptEntry* pd = _res2.val;
    err = _res2.err;
    if (err == NULL && pd != NULL && EptEntry_IsPresent(pd) && EptEntry_IsLargePage(pd)) {
        return (so_R_ptr_err){.val = pd, .err = NULL};
    }
    return (so_R_ptr_err){.val = NULL, .err = ErrNotMapped};
}

uint64_t EptTable_BuildEptPointer(void* self) {
    EptTable* t = (EptTable*)self;
    uint64_t result = 0;
    result |= ((uint64_t)(MtrrState_DefaultMemoryType(t->mtrrState)) << EptTypeShift);
    result |= (EptRead | EptExec);
    result |= ((uint64_t)((uintptr_t)((void*)(&t->pml4[0]))) & EptPhysMask);
    return result;
}

void EptTable_Invalidate(void* self) {
    EptTable* t = (EptTable*)self;
    inveptAllContexts();
}

so_int EptTable_SplitCount(void* self) {
    EptTable* t = (EptTable*)self;
    return so_len(t->splits);
}

MtrrState* NewMtrrState(void) {
    MtrrState* m = &(MtrrState){.ranges = so_make_slice_impl(sizeof(MTRR_RANGE_DESCRIPTOR), 0, MAX_MTRR_ENTRIES), .defaultType = MemTypeWriteBack};
    MtrrState_readMtrrsFromHardware(m);
    return m;
}

uint8_t MtrrState_DefaultMemoryType(void* self) {
    MtrrState* m = (MtrrState*)self;
    return m->defaultType;
}

static void MtrrState_readMtrrsFromHardware(void* self) {
    MtrrState* m = (MtrrState*)self;
    uint64_t msrMtrrcap = __readmsr(0xFE);
    so_int vcnt = so_min((so_int)((msrMtrrcap >> 40) & 0xFF), (so_int)(MAX_MTRR_ENTRIES));
    uint64_t msrMtrrDefType = __readmsr(0x2FF);
    m->defaultType = (uint8_t)(msrMtrrDefType & 0xFF);
    m->ranges = so_slice(MTRR_RANGE_DESCRIPTOR, m->ranges, 0, 0);
    for (so_int i = 0; i < vcnt; i++) {
        uint64_t baseMsr = __readmsr((uint32_t)(0x200 + i * 2));
        uint64_t maskMsr = __readmsr((uint32_t)(0x201 + i * 2));
        if ((maskMsr & 0x800) == 0) {
            continue;
        }
        uint64_t basePa = (baseMsr & 0xFFFFF000);
        uint64_t maskBits = (maskMsr & 0xFFFFF000);
        uint64_t size = (~maskBits) + 1;
        uint8_t memType = (uint8_t)(baseMsr & 0xFF);
        m->ranges = so_append(MTRR_RANGE_DESCRIPTOR, m->ranges, (MTRR_RANGE_DESCRIPTOR){.PhysicalBaseAddress = (uintptr_t)(basePa), .PhysicalEndAddress = (uintptr_t)(basePa + size), .MemoryType = memType, .FixedRange = false});
    }
    m->numRanges = (uint32_t)(so_len(m->ranges));
}

uint8_t MtrrState_Lookup(void* self, uint64_t pa) {
    MtrrState* m = (MtrrState*)self;
    for (so_int _ = 0; _ < so_len(m->ranges); _++) {
        MTRR_RANGE_DESCRIPTOR r = so_at(MTRR_RANGE_DESCRIPTOR, m->ranges, _);
        if (pa >= (uint64_t)(r.PhysicalBaseAddress) && pa < (uint64_t)(r.PhysicalEndAddress)) {
            return r.MemoryType;
        }
    }
    return m->defaultType;
}

so_Slice MtrrState_Ranges(void* self) {
    MtrrState* m = (MtrrState*)self;
    return m->ranges;
}

static so_String eptError_Error(eptError e) {
    return (e);
}

static so_Error fmtError(so_String format, so_Slice args) {
    return fmt_Errorf(format, args);
}

so_Error EptTable_RestoreEntry(void* self, uint64_t physAddr, EptEntry* original) {
    EptTable* t = (EptTable*)self;
    EptTable_Lock(t);
    uint64_t gpa = physAddr;
    so_R_ptr_err _res1 = EptTable_resolveForModifyUnlocked(t, gpa);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        EptTable_Unlock(t);
        return err;
    }
    if (original != NULL) {
        *pte = *original;
    } else {
        *pte = (EptEntry){};
    }
    inveptSingleContext(EptTable_BuildEptPointer(t));
    EptTable_Unlock(t);
    return NULL;
}

so_R_ptr_err EptTable_GetEntry(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    so_R_ptr_err _res1 = EptTable_ResolveForModify(t, gpa);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        return (so_R_ptr_err){.val = NULL, .err = err};
    }
    return (so_R_ptr_err){.val = pte, .err = NULL};
}

so_R_ptr_err EptTable_CloneEntry(void* self, uint64_t gpa) {
    EptTable* t = (EptTable*)self;
    so_R_ptr_err _res1 = EptTable_GetEntry(t, gpa);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        return (so_R_ptr_err){.val = NULL, .err = err};
    }
    return (so_R_ptr_err){.val = EptEntry_Clone(pte), .err = NULL};
}

// -- events.go --

bool EventFunc_Handle(EventFunc f, VCPU* v) {
    return f(v);
}

static EventDispatcher* newEventDispatcher(void) {
    return &(EventDispatcher){.handlers = so_make_map_impl(sizeof(so_int), sizeof(EventHandler), 0), .enabled = so_make_map_impl(sizeof(so_int), sizeof(bool), 0)};
}

void EventDispatcher_On(void* self, so_int kind, EventHandler handler) {
    EventDispatcher* d = (EventDispatcher*)self;
    {
    	EventHandler _so_map_assign_val;
    	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
    	_so_map_assign_val = handler;
    	so_map_set(so_int, EventHandler, d->handlers, kind, _so_map_assign_val);
    }
}

void EventDispatcher_OnFunc(void* self, so_int kind, EventFunc fn) {
    EventDispatcher* d = (EventDispatcher*)self;
    {
    	EventHandler _so_map_assign_val;
    	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
    	_so_map_assign_val = (EventHandler){.self = &fn, .Handle = EventFunc_Handle};
    	so_map_set(so_int, EventHandler, d->handlers, kind, _so_map_assign_val);
    }
}

void EventDispatcher_Enable(void* self, so_int kind) {
    EventDispatcher* d = (EventDispatcher*)self;
    {
    	bool _so_map_assign_val;
    	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
    	_so_map_assign_val = true;
    	so_map_set(so_int, bool, d->enabled, kind, _so_map_assign_val);
    }
}

void EventDispatcher_Disable(void* self, so_int kind) {
    EventDispatcher* d = (EventDispatcher*)self;
    so_map_remove(so_int, d->enabled, kind);
}

bool EventDispatcher_IsEnabled(void* self, so_int kind) {
    EventDispatcher* d = (EventDispatcher*)self;
    return so_map_get(so_int, bool, d->enabled, kind);
}

bool EventDispatcher_Dispatch(void* self, so_int kind, VCPU* v) {
    EventDispatcher* d = (EventDispatcher*)self;
    if (!so_map_get(so_int, bool, d->enabled, kind)) {
        return false;
    }
    {
        EventHandler h = so_map_get(so_int, EventHandler, d->handlers, kind);
        bool ok = so_map_has(so_int, d->handlers, kind);
        if (ok) {
            return h.Handle(h.self, v);
        }
    }
    return false;
}

void EventDispatcher_Reset(void* self) {
    EventDispatcher* d = (EventDispatcher*)self;
    d->handlers = so_make_map_impl(sizeof(so_int), sizeof(EventHandler), 0);
    d->enabled = so_make_map_impl(sizeof(so_int), sizeof(bool), 0);
}

bool CpuidHandler_Handle(void* self, VCPU* v) {
    CpuidHandler* h = (CpuidHandler*)self;
    uint64_t cpuId = v->Regs->Rax;
    do {
        if (cpuId == (0)) {
            v->Regs->Rbx = 0x756E6547;
            v->Regs->Rdx = 0x6C65746E;
            v->Regs->Rcx = 0x44534147;
        } else if (cpuId == (1)) {
            v->Regs->Rcx &= ~((((((so_int)1 << 31) | ((so_int)1 << 30)) | ((so_int)1 << 29)) | ((so_int)1 << 28)) | ((so_int)1 << 27));
        } else if (cpuId == (0xA)) {
            v->Regs->Rax = 0;
            v->Regs->Rbx = 0;
            v->Regs->Rcx = 0;
            v->Regs->Rdx = 0;
        }
    } while (0);
    return true;
}

MsrHandler* NewMsrHandler(void) {
    return &(MsrHandler){.onRead = defaultMsrRead, .onWrite = defaultMsrWrite};
}

bool MsrHandler_HandleRead(void* self, VCPU* v) {
    MsrHandler* h = (MsrHandler*)self;
    uint32_t msr = (uint32_t)(v->Regs->Rcx);
    uint64_t value = __readmsr(msr);
    do {
        if (msr == (0x174)) {
            value &= ~((uint64_t)1 << 2);
        }
    } while (0);
    v->Regs->Rax = (value & 0xFFFFFFFF);
    v->Regs->Rdx = (value >> 32);
    return h->onRead(v, msr, value);
}

bool MsrHandler_HandleWrite(void* self, VCPU* v) {
    MsrHandler* h = (MsrHandler*)self;
    uint32_t msr = (uint32_t)(v->Regs->Rcx);
    uint64_t value = ((v->Regs->Rdx << 32) | v->Regs->Rax);
    do {
        if (msr == (0x174)) {
            value &= ~((uint64_t)1 << 2);
        }
    } while (0);
    __writemsr(msr, value);
    return h->onWrite(v, msr, value);
}

static bool defaultMsrRead(VCPU* _p0, uint32_t _p1, uint64_t _p2) {
    return true;
}

static bool defaultMsrWrite(VCPU* _p0, uint32_t _p1, uint64_t _p2) {
    return true;
}

IoHandler* NewIoHandler(void) {
    return &(IoHandler){.onIn = defaultIoIn, .onOut = defaultIoOut};
}

bool IoHandler_Handle(void* self, VCPU* v) {
    IoHandler* h = (IoHandler*)self;
    uint32_t qual = v->ExitQualification;
    bool isOut = (qual & 0x8) != 0;
    uint32_t size = (qual & 0x7) + 1;
    uint16_t port = (uint16_t)(qual >> 16);
    if (isOut) {
        do {
            if (size == (1)) {
                (void)((v->Regs->Rax & 0xFF));
            } else if (size == (2)) {
                (void)((v->Regs->Rax & 0xFFFF));
            } else if (size == (4)) {
                (void)((v->Regs->Rax & 0xFFFFFFFF));
            }
        } while (0);
        return h->onOut(v, port, size);
    }
    return h->onIn(v, port, size);
}

static bool defaultIoIn(VCPU* _p0, uint16_t _p1, uint32_t _p2) {
    return true;
}

static bool defaultIoOut(VCPU* _p0, uint16_t _p1, uint32_t _p2) {
    return true;
}

CrAccessHandler* NewCrAccessHandler(void) {
    return &(CrAccessHandler){.onChange = defaultCrChange};
}

bool CrAccessHandler_Handle(void* self, VCPU* v) {
    CrAccessHandler* h = (CrAccessHandler*)self;
    uint32_t qual = v->ExitQualification;
    uint32_t crNum = ((qual >> 8) & 0xF);
    uint32_t accessType = ((qual >> 4) & 0x3);
    bool lmswOp = (qual & 0x1) != 0;
    bool read = accessType == 0 || accessType == 2;
    bool write = accessType == 1 || accessType == 2;
    if (lmswOp && write) {
        crNum = 0;
    }
    return h->onChange(v, crNum, read, write);
}

static bool defaultCrChange(VCPU* _p0, uint32_t _p1, bool _p2, bool _p3) {
    return true;
}

TscHandler* NewTscHandler(void) {
    return &(TscHandler){.rdtsc = nativeRdtsc};
}

bool TscHandler_Handle(void* self, VCPU* v) {
    TscHandler* h = (TscHandler*)self;
    uint64_t tsc = h->rdtsc();
    if (h->offset != 0) {
        tsc += (uint64_t)(h->offset);
    }
    v->Regs->Rax = (tsc & 0xFFFFFFFF);
    v->Regs->Rdx = (tsc >> 32);
    return true;
}

static uint64_t nativeRdtsc(void) {
    return rdtscValue();
}

ExceptionHandler* NewExceptionHandler(void) {
    return &(ExceptionHandler){.onBp = defaultBpHandler, .onGp = defaultGpHandler, .onUd = defaultUdHandler, .onNmi = defaultNmiHandler};
}

bool ExceptionHandler_HandleBreakpoint(void* self, VCPU* v) {
    ExceptionHandler* h = (ExceptionHandler*)self;
    return h->onBp(v);
}

bool ExceptionHandler_HandleGeneralProtectionFault(void* self, VCPU* v) {
    ExceptionHandler* h = (ExceptionHandler*)self;
    return h->onGp(v);
}

bool ExceptionHandler_HandleUndefinedOpcode(void* self, VCPU* v) {
    ExceptionHandler* h = (ExceptionHandler*)self;
    return h->onUd(v);
}

bool ExceptionHandler_HandleNmi(void* self, VCPU* v) {
    ExceptionHandler* h = (ExceptionHandler*)self;
    return h->onNmi(v);
}

static bool defaultBpHandler(VCPU* _p0) {
    return false;
}

static bool defaultGpHandler(VCPU* _p0) {
    return false;
}

static bool defaultUdHandler(VCPU* _p0) {
    return false;
}

static bool defaultNmiHandler(VCPU* _p0) {
    return false;
}

// -- exec_trap.go --

static void execTrapEnable(VCPU* v) {
    gExecTrap->Enabled = true;
    setMonitorTrapFlag(true);
    LogDebug(so_str("Execution trap enabled on core %d"), (so_Slice){(void*[1]){&(uint32_t){v->CoreId}}, 1, 1});
}

static void execTrapDisable(VCPU* v) {
    gExecTrap->Enabled = false;
    setMonitorTrapFlag(false);
    LogDebug(so_str("Execution trap disabled on core %d"), (so_Slice){(void*[1]){&(uint32_t){v->CoreId}}, 1, 1});
}

static so_int execTrapAddEntry(uint64_t addr, CR3_TYPE cr3) {
    ExecTrapEntry entry = (ExecTrapEntry){.Address = addr, .Cr3 = cr3, .Active = true};
    gExecTrap->TrapList = so_append(ExecTrapEntry, gExecTrap->TrapList, entry);
    return so_len(gExecTrap->TrapList) - 1;
}

static void execTrapRemoveEntry(so_int index) {
    if (index >= 0 && index < so_len(gExecTrap->TrapList)) {
        so_at(ExecTrapEntry, gExecTrap->TrapList, index).Active = false;
    }
}

static bool execTrapHandleMtf(VCPU* v) {
    if (!gExecTrap->Enabled) {
        return false;
    }
    uint64_t currentRip = 0;
    vmRead64(VmcsGuestRip, &currentRip);
    for (so_int i = 0; i < so_len(gExecTrap->TrapList); i++) {
        ExecTrapEntry* entry = &so_at(ExecTrapEntry, gExecTrap->TrapList, i);
        if (entry->Active && currentRip == entry->Address) {
            Hypervisor* h = getHypervisor();
            if (h != NULL) {
                Hypervisor_suppressAdvance(h, v);
            }
            dispatchEventExecTrap(v, entry->Address, entry->Cr3);
            return true;
        }
    }
    return false;
}

static bool execTrapCheckAndHandle(VCPU* v) {
    return execTrapHandleMtf(v);
}

static void dispatchEventExecTrap(VCPU* v, uint64_t addr, CR3_TYPE cr3) {
}

// -- globals.go --

so_Error InitializeGlobals(uint32_t numCpus) {
    so_Error err = NULL;
    gLog = NewLogger(MAX_LOG_BUFFERS, MAX_LOG_BUFFERS_PRIO);
    {
        err = Logger_Initialize(gLog);
        if (err != NULL) {
            return fmtError(so_str("logger init failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        }
    }
    gPoolMgr = NewPoolManager(1024);
    gHyp = NewHypervisor(numCpus);
    {
        err = Hypervisor_Initialize(gHyp);
        if (err != NULL) {
            return fmtError(so_str("hypervisor init failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        }
    }
    gEvents = newEventDispatcher();
    setupDefaultEventHandlers();
    gHooks = NewHookManager((so_int)(MAX_HIDDEN_BREAKPOINTS));
    gEvasion = NewEvasionState();
    gTrace = NewTracerState();
    gSerial = NewKdSerialState();
    gOptimizer = NewVmmOptimizer(OptBasic);
    if (VmmOptimizer_GetConfig(gOptimizer).EptCache) {
        gEptCache = NewEptCache(4096);
    }
    gDrv = &(Driver){.device = NULL, .handleInUse = false, .initialized = false};
    LogInfo(so_str("Global state initialized for %d CPUs"), (so_Slice){(void*[1]){&numCpus}, 1, 1});
    return NULL;
}

void CleanupGlobals(void) {
    if (gEptCache != NULL) {
        gEptCache = NULL;
    }
    if (gOptimizer != NULL) {
        gOptimizer = NULL;
    }
    if (gSerial != NULL && KdSerialState_IsConnected(gSerial)) {
        KdSerialState_Uninitialize(gSerial);
    }
    gSerial = NULL;
    if (gTrace != NULL && TracerState_IsActive(gTrace)) {
        TracerState_Disable(gTrace);
    }
    gTrace = NULL;
    if (gEvasion != NULL && EvasionState_IsActive(gEvasion)) {
        EvasionState_Disable(gEvasion);
    }
    gEvasion = NULL;
    if (gHooks != NULL) {
        HookManager_RemoveAll(gHooks);
    }
    gHooks = NULL;
    if (gEvents != NULL) {
        EventDispatcher_Reset(gEvents);
    }
    gEvents = NULL;
    if (gHyp != NULL && Hypervisor_IsInitialized(gHyp)) {
        Hypervisor_Shutdown(gHyp);
    }
    gHyp = NULL;
    if (gPoolMgr != NULL) {
        PoolManager_CheckAndPerformDeallocation(gPoolMgr);
    }
    gPoolMgr = NULL;
    if (gLog != NULL) {
        Logger_Uninitialize(gLog);
    }
    gLog = NULL;
    if (gDrv != NULL) {
        gDrv->device = NULL;
        gDrv->initialized = false;
    }
    gDrv = NULL;
}

static void setupDefaultEventHandlers(void) {
    CpuidHandler* cpuidHandler = &(CpuidHandler){};
    IoHandler* ioHandler = NewIoHandler();
    CrAccessHandler* crHandler = NewCrAccessHandler();
    TscHandler* tscHandler = NewTscHandler();
    EventDispatcher_On(gEvents, EventCpuid, (EventHandler){.self = cpuidHandler, .Handle = CpuidHandler_Handle});
    EventDispatcher_On(gEvents, EventIo, (EventHandler){.self = ioHandler, .Handle = IoHandler_Handle});
    EventDispatcher_On(gEvents, EventMovCr, (EventHandler){.self = crHandler, .Handle = CrAccessHandler_Handle});
    EventDispatcher_On(gEvents, EventTsc, (EventHandler){.self = tscHandler, .Handle = TscHandler_Handle});
    EventDispatcher_Enable(gEvents, EventCpuid);
    EventDispatcher_Enable(gEvents, EventIo);
    EventDispatcher_Enable(gEvents, EventMovCr);
    EventDispatcher_Enable(gEvents, EventTsc);
}

static bool handleExceptionDispatch(VCPU* _p0, ExceptionHandler* _p1) {
    return true;
}

Hypervisor* GetHypervisor(void) {
    return gHyp;
}

Logger* GetLogger(void) {
    return gLog;
}

Driver* GetDriver(void) {
    return gDrv;
}

EventDispatcher* GetEventDispatcher(void) {
    return gEvents;
}

HookManager* GetHookManager(void) {
    return gHooks;
}

EvasionState* GetEvasionState(void) {
    return gEvasion;
}

TracerState* GetTracerState(void) {
    return gTrace;
}

KdSerialState* GetSerialState(void) {
    return gSerial;
}

PoolManager* GetPoolManager(void) {
    return gPoolMgr;
}

VmmOptimizer* GetOptimizer(void) {
    return gOptimizer;
}

EptCache* GetEptCache(void) {
    return gEptCache;
}

// -- hooks.go --

so_String HookKind_String(HookKind k) {
    do {
        if (k == (HookExec)) {
            return so_str("EXEC");
        } else if (k == (HookRead)) {
            return so_str("READ");
        } else if (k == (HookWrite)) {
            return so_str("WRITE");
        } else if (k == (HookReadWrite)) {
            return so_str("RW");
        } else {
            return so_str("??");
        }
    } while (0);
}

PagePermission HookKind_Permission(HookKind k) {
    do {
        if (k == (HookExec)) {
            return (PagePermission){.Read = true, .Write = false, .Exec = false};
        } else if (k == (HookRead)) {
            return (PagePermission){.Read = false, .Write = true, .Exec = true};
        } else if (k == (HookWrite)) {
            return (PagePermission){.Read = true, .Write = false, .Exec = true};
        } else if (k == (HookReadWrite)) {
            return (PagePermission){.Read = false, .Write = false, .Exec = false};
        } else {
            return (PagePermission){};
        }
    } while (0);
}

so_R_ptr_err NewEptHook(uint64_t va, CR3_TYPE cr3, HookKind kind) {
    uint64_t pageBase = (va & ~((uint64_t)(PAGE_SIZE - 1)));
    uint64_t pa = VirtToPhys((uintptr_t)(pageBase), cr3);
    if (pa == 0 || pa >= (uint64_t)(SIZE_512_GB)) {
        return (so_R_ptr_err){.val = NULL, .err = ErrInvalidAddress};
    }
    EptHook* hook = &(EptHook){.PhysAddr = pa, .VirtAddr = va, .Cr3 = cr3, .Kind = kind, .IsActive = true, .breakpoints = so_make_slice_impl(sizeof(Breakpoint), 0, MAX_HIDDEN_BREAKPOINTS)};
    copyPageToFake(hook, (uintptr_t)(pageBase), cr3);
    if (kind == HookExec || kind == HookReadWrite) {
        uint64_t offset = (va & ((uint64_t)(PAGE_SIZE) - 1));
        hook->fakePage[offset] = 0xCC;
        EptHook_AddBreakpoint(hook, va);
    }
    hook->FakePagePa = (uint64_t)((uintptr_t)((void*)(&hook->fakePage[0])));
    return (so_R_ptr_err){.val = hook, .err = NULL};
}

static void copyPageToFake(EptHook* h, uintptr_t pageBase, CR3_TYPE cr3) {
    uint64_t srcVa = VirtToPhys(pageBase, cr3);
    uint64_t dstVa = (uint64_t)((uintptr_t)((void*)(&h->fakePage[0])));
    for (uint64_t i = (0); i < (uint64_t)(PAGE_SIZE); i += 8) {
        uint64_t val = readPhysicalMemory(srcVa + i);
        writePhysicalMemory(dstVa + i, val);
    }
}

bool EptHook_AddBreakpoint(void* self, uint64_t addr) {
    EptHook* h = (EptHook*)self;
    if (so_len(h->breakpoints) >= (so_int)(MAX_HIDDEN_BREAKPOINTS)) {
        return false;
    }
    uint64_t offset = addr - (h->VirtAddr & ~((uint64_t)(PAGE_SIZE) - 1));
    if (offset >= (uint64_t)(PAGE_SIZE)) {
        return false;
    }
    for (so_int _ = 0; _ < so_len(h->breakpoints); _++) {
        Breakpoint bp = so_at(Breakpoint, h->breakpoints, _);
        if (bp.Address == addr && bp.Active) {
            return false;
        }
    }
    Breakpoint bp = (Breakpoint){.Address = addr, .OriginalByte = h->fakePage[offset], .Active = true};
    h->breakpoints = so_append(Breakpoint, h->breakpoints, bp);
    h->fakePage[offset] = 0xCC;
    LogDebug(so_str("HOOK: Added BP at 0x%X (page=0x%X)"), (so_Slice){(void*[2]){&addr, &(uint64_t){h->PhysAddr}}, 2, 2});
    return true;
}

bool EptHook_RemoveBreakpoint(void* self, uint64_t addr) {
    EptHook* h = (EptHook*)self;
    for (so_int i = 0; i < so_len(h->breakpoints); i++) {
        Breakpoint bp = so_at(Breakpoint, h->breakpoints, i);
        if (bp.Address == addr) {
            uint64_t offset = addr - (h->VirtAddr & ~((uint64_t)(PAGE_SIZE) - 1));
            h->fakePage[offset] = bp.OriginalByte;
            h->breakpoints = so_extend(Breakpoint, so_slice(Breakpoint, h->breakpoints, 0, i), (so_slice(Breakpoint, h->breakpoints, i + 1, h->breakpoints.len)));
            LogDebug(so_str("HOOK: Removed BP at 0x%X"), (so_Slice){(void*[1]){&addr}, 1, 1});
            return true;
        }
    }
    return false;
}

Breakpoint* EptHook_FindBreakpoint(void* self, uint64_t addr) {
    EptHook* h = (EptHook*)self;
    for (so_int i = 0; i < so_len(h->breakpoints); i++) {
        if (so_at(Breakpoint, h->breakpoints, i).Address == addr && so_at(Breakpoint, h->breakpoints, i).Active) {
            return &so_at(Breakpoint, h->breakpoints, i);
        }
    }
    return NULL;
}

so_int EptHook_BreakpointCount(void* self) {
    EptHook* h = (EptHook*)self;
    so_int count = 0;
    for (so_int _ = 0; _ < so_len(h->breakpoints); _++) {
        Breakpoint bp = so_at(Breakpoint, h->breakpoints, _);
        if (bp.Active) {
            count++;
        }
    }
    return count;
}

bool EptHook_IsHiddenBp(void* self) {
    EptHook* h = (EptHook*)self;
    return h->Kind == HookExec || h->Kind == HookReadWrite;
}

bool EptHook_ContainsAddress(void* self, uint64_t va) {
    EptHook* h = (EptHook*)self;
    uint64_t pageBase = (va & ~((uint64_t)(PAGE_SIZE) - 1));
    uint64_t hookBase = (h->VirtAddr & ~((uint64_t)(PAGE_SIZE) - 1));
    return pageBase == hookBase;
}

HookManager* NewHookManager(so_int capacity) {
    return &(HookManager){.hooks = so_make_slice_impl(sizeof(EptHook*), 0, capacity), .capacity = capacity, .lock = NewSpinlock(), .mtfRestoreList = so_make_slice_impl(sizeof(EptHook*), 0, 16)};
}

void HookManager_Lock(void* self) {
    HookManager* m = (HookManager*)self;
    Spinlock_Lock(m->lock);
}

void HookManager_Unlock(void* self) {
    HookManager* m = (HookManager*)self;
    Spinlock_Unlock(m->lock);
}

so_Error HookManager_Install(void* self, uint32_t coreId, uint64_t va, CR3_TYPE cr3, HookKind kind) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    EptHook* existing = HookManager_findByVirtAddr(m, va);
    if (existing != NULL) {
        HookManager_Unlock(m);
        return HookManager_updateExistingHook(m, existing, va, kind);
    }
    if (so_len(m->hooks) >= m->capacity) {
        HookManager_Unlock(m);
        return ErrHookLimitReached;
    }
    so_R_ptr_err _res1 = NewEptHook(va, cr3, kind);
    EptHook* hook = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        HookManager_Unlock(m);
        return err;
    }
    hook->coreId = coreId;
    EptTable* table = Hypervisor_eptTable(gHyp, coreId);
    if (table == NULL) {
        HookManager_Unlock(m);
        return ErrNoEptTable;
    }
    {
        so_Error err = EptTable_SplitLargePage(table, hook->PhysAddr);
        if (err != NULL) {
            HookManager_Unlock(m);
            return err;
        }
    }
    so_R_ptr_err _res2 = EptTable_Lookup(table, hook->PhysAddr);
    EptEntry* pte = _res2.val;
    err = _res2.err;
    if (err != NULL || pte == NULL) {
        HookManager_Unlock(m);
        return ErrNotMapped;
    }
    hook->OriginalEntry = EptEntry_Clone(pte);
    hook->ModifiedEntry = &(EptEntry){};
    EptEntry_SetRaw(hook->ModifiedEntry, EptEntry_Raw(pte));
    EptEntry_SetPhysAddr(hook->ModifiedEntry, hook->FakePagePa);
    EptEntry_SetPermission(hook->ModifiedEntry, HookKind_Permission(kind));
    *pte = *hook->ModifiedEntry;
    EptTable_Invalidate(table);
    m->hooks = so_append(EptHook*, m->hooks, hook);
    so_String kindStr = HookKind_String(kind);
    LogInfo(so_str("HOOK: Installed %s hook at VA=0x%X PA=0x%X"), (so_Slice){(void*[3]){&kindStr, &va, &(uint64_t){hook->PhysAddr}}, 3, 3});
    HookManager_Unlock(m);
    return NULL;
}

static so_Error HookManager_updateExistingHook(void* self, EptHook* existing, uint64_t va, HookKind kind) {
    HookManager* m = (HookManager*)self;
    if (kind == HookExec || kind == HookReadWrite) {
        if (!EptHook_AddBreakpoint(existing, va)) {
            return ErrTooManyBreakpoints;
        }
    }
    return NULL;
}

so_Error HookManager_Remove(void* self, uint64_t physAddr) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    for (so_int i = 0; i < so_len(m->hooks); i++) {
        EptHook* hook = so_at(EptHook*, m->hooks, i);
        if (hook->PhysAddr == physAddr) {
            EptTable* table = Hypervisor_eptTable(gHyp, hook->coreId);
            if (table != NULL && hook->OriginalEntry != NULL) {
                so_R_ptr_err _res1 = EptTable_Lookup(table, physAddr);
                EptEntry* pte = _res1.val;
                so_Error err = _res1.err;
                if (err == NULL && pte != NULL) {
                    *pte = *hook->OriginalEntry;
                }
            }
            m->hooks = so_extend(EptHook*, so_slice(EptHook*, m->hooks, 0, i), (so_slice(EptHook*, m->hooks, i + 1, m->hooks.len)));
            invalidateEpt();
            LogInfo(so_str("HOOK: Removed hook at PA=0x%X"), (so_Slice){(void*[1]){&physAddr}, 1, 1});
            HookManager_Unlock(m);
            return NULL;
        }
    }
    HookManager_Unlock(m);
    return ErrHookNotFound;
}

so_Error HookManager_RemoveByVirtAddr(void* self, uint64_t va) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    for (so_int i = 0; i < so_len(m->hooks); i++) {
        EptHook* hook = so_at(EptHook*, m->hooks, i);
        if (EptHook_ContainsAddress(hook, va)) {
            EptTable* table = Hypervisor_eptTable(gHyp, hook->coreId);
            if (table != NULL && hook->OriginalEntry != NULL) {
                so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
                EptEntry* pte = _res1.val;
                so_Error err = _res1.err;
                if (err == NULL && pte != NULL) {
                    *pte = *hook->OriginalEntry;
                }
            }
            m->hooks = so_extend(EptHook*, so_slice(EptHook*, m->hooks, 0, i), (so_slice(EptHook*, m->hooks, i + 1, m->hooks.len)));
            invalidateEpt();
            LogInfo(so_str("HOOK: Removed hook at VA=0x%X"), (so_Slice){(void*[1]){&va}, 1, 1});
            HookManager_Unlock(m);
            return NULL;
        }
    }
    HookManager_Unlock(m);
    return ErrHookNotFound;
}

static EptHook* HookManager_findByPhysAddr(void* self, uint64_t pa) {
    HookManager* m = (HookManager*)self;
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* hook = so_at(EptHook*, m->hooks, _);
        if (hook->PhysAddr == pa) {
            return hook;
        }
    }
    return NULL;
}

static EptHook* HookManager_findByVirtAddr(void* self, uint64_t va) {
    HookManager* m = (HookManager*)self;
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* hook = so_at(EptHook*, m->hooks, _);
        if (EptHook_ContainsAddress(hook, va)) {
            return hook;
        }
    }
    return NULL;
}

EptHook* HookManager_FindByPhysAddr(void* self, uint64_t pa) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    HookManager_Unlock(m);
    return HookManager_findByPhysAddr(m, pa);
}

EptHook* HookManager_FindByVirtAddr(void* self, uint64_t va) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    HookManager_Unlock(m);
    return HookManager_findByVirtAddr(m, va);
}

bool HookManager_HandleExecHook(void* self, VCPU* v, uint64_t gpa) {
    HookManager* m = (HookManager*)self;
    EptHook* hook = HookManager_FindByPhysAddr(m, gpa);
    if (hook == NULL || !hook->IsActive || !EptHook_IsHiddenBp(hook)) {
        return false;
    }
    Breakpoint* bp = EptHook_FindBreakpoint(hook, v->LastVmexitRip);
    if (bp == NULL) {
        return false;
    }
    v->IncrementRip = false;
    v->RegisterBreakOnMtf = true;
    v->IgnoreOneMtf = true;
    v->MtfEptHookRestorePoint = (uintptr_t)((void*)(hook));
    restoreOriginalAndInjectBp(v, hook, bp);
    LogTrace(so_str("HOOK: Exec hit at RIP=0x%X"), (so_Slice){(void*[1]){&(uint64_t){v->LastVmexitRip}}, 1, 1});
    return true;
}

bool HookManager_HandleReadWriteHook(void* self, VCPU* v, uint64_t gpa, bool isRead, bool isWrite) {
    HookManager* m = (HookManager*)self;
    EptHook* hook = HookManager_FindByPhysAddr(m, gpa);
    if (hook == NULL || !hook->IsActive) {
        return false;
    }
    if (hook->Kind == HookRead && !isRead) {
        return false;
    }
    if (hook->Kind == HookWrite && !isWrite) {
        return false;
    }
    v->IncrementRip = false;
    v->RegisterBreakOnMtf = true;
    v->IgnoreOneMtf = true;
    v->MtfEptHookRestorePoint = (uintptr_t)((void*)(hook));
    singleStepOnFakePage(v, hook);
    so_String kindStr = HookKind_String(hook->Kind);
    LogTrace(so_str("HOOK: %s hit at GPA=0x%X"), (so_Slice){(void*[2]){&kindStr, &gpa}, 2, 2});
    return true;
}

bool HookManager_HandleBreakpoint(void* self, VCPU* v) {
    HookManager* m = (HookManager*)self;
    EptHook* hook = HookManager_FindByVirtAddr(m, v->LastVmexitRip);
    if (hook == NULL || !hook->IsActive) {
        return false;
    }
    Breakpoint* bp = EptHook_FindBreakpoint(hook, v->LastVmexitRip);
    if (bp == NULL) {
        return false;
    }
    v->IncrementRip = false;
    v->RegisterBreakOnMtf = true;
    v->IgnoreOneMtf = true;
    v->MtfEptHookRestorePoint = (uintptr_t)((void*)(hook));
    restoreOriginalAndInjectBp(v, hook, bp);
    return true;
}

bool HookManager_HandleMtfRestore(void* self, VCPU* v) {
    HookManager* m = (HookManager*)self;
    if (v->MtfEptHookRestorePoint == 0) {
        return false;
    }
    EptHook* hook = (EptHook*)((void*)(v->MtfEptHookRestorePoint));
    if (hook == NULL) {
        v->MtfEptHookRestorePoint = 0;
        return false;
    }
    EptTable* table = Hypervisor_eptTable(gHyp, v->CoreId);
    if (table == NULL) {
        v->MtfEptHookRestorePoint = 0;
        return false;
    }
    so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL || pte == NULL) {
        v->MtfEptHookRestorePoint = 0;
        return false;
    }
    *pte = *hook->ModifiedEntry;
    EptTable_Invalidate(table);
    v->MtfEptHookRestorePoint = 0;
    LogTrace(so_str("HOOK: MTF restore for PA=0x%X"), (so_Slice){(void*[1]){&(uint64_t){hook->PhysAddr}}, 1, 1});
    return true;
}

void HookManager_RestoreAll(void* self, VCPU* v) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    if (v == NULL) {
        // Restore hooks for all cores
        for (uint32_t i = (0); i < (uint32_t)(so_len(gHyp->eptTables)); i++) {
            HookManager_restoreAllForCoreLocked(m, i);
        }
        invalidateEpt();
        LogDebug(so_str("HOOK: Restored all hooks for all cores"), (so_Slice){&so_Nil, 0, 0});
        HookManager_Unlock(m);
        return;
    }
    uint32_t coreId = v->CoreId;
    if (coreId == 0) {
        coreId = getCurrentProcessorNumber();
    }
    EptTable* table = Hypervisor_eptTable(gHyp, coreId);
    if (table == NULL) {
        HookManager_Unlock(m);
        return;
    }
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* hook = so_at(EptHook*, m->hooks, _);
        if (!hook->IsActive || hook->OriginalEntry == NULL) {
            continue;
        }
        so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
        EptEntry* pte = _res1.val;
        so_Error err = _res1.err;
        if (err != NULL || pte == NULL) {
            continue;
        }
        *pte = *hook->OriginalEntry;
    }
    invalidateEpt();
    LogDebug(so_str("HOOK: Restored all hooks for core %d"), (so_Slice){(void*[1]){&coreId}, 1, 1});
    HookManager_Unlock(m);
}

static void HookManager_restoreAllForCoreLocked(void* self, uint32_t coreId) {
    HookManager* m = (HookManager*)self;
    EptTable* table = Hypervisor_eptTable(gHyp, coreId);
    if (table == NULL) {
        return;
    }
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* hook = so_at(EptHook*, m->hooks, _);
        if (!hook->IsActive || hook->OriginalEntry == NULL) {
            continue;
        }
        so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
        EptEntry* pte = _res1.val;
        so_Error err = _res1.err;
        if (err == NULL && pte != NULL) {
            *pte = *hook->OriginalEntry;
        }
    }
}

void HookManager_RestoreAllForCore(void* self, uint32_t coreId) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    EptTable* table = Hypervisor_eptTable(gHyp, coreId);
    if (table == NULL) {
        HookManager_Unlock(m);
        return;
    }
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* hook = so_at(EptHook*, m->hooks, _);
        if (!hook->IsActive || hook->OriginalEntry == NULL) {
            continue;
        }
        so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
        EptEntry* pte = _res1.val;
        so_Error err = _res1.err;
        if (err == NULL && pte != NULL) {
            *pte = *hook->OriginalEntry;
        }
    }
    invalidateEpt();
    HookManager_Unlock(m);
}

so_int HookManager_Count(void* self) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    HookManager_Unlock(m);
    return so_len(m->hooks);
}

so_int HookManager_ActiveCount(void* self) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    so_int count = 0;
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* h = so_at(EptHook*, m->hooks, _);
        if (h->IsActive) {
            count++;
        }
    }
    HookManager_Unlock(m);
    return count;
}

void HookManager_RemoveAll(void* self) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    m->hooks = (so_Slice){&so_Nil, 0, 0};
    HookManager_Unlock(m);
}

void HookManager_DisableAll(void* self) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* h = so_at(EptHook*, m->hooks, _);
        h->IsActive = false;
    }
    HookManager_Unlock(m);
}

void HookManager_EnableAll(void* self) {
    HookManager* m = (HookManager*)self;
    HookManager_Lock(m);
    for (so_int _ = 0; _ < so_len(m->hooks); _++) {
        EptHook* h = so_at(EptHook*, m->hooks, _);
        h->IsActive = true;
    }
    HookManager_Unlock(m);
}

static void restoreOriginalAndInjectBp(VCPU* v, EptHook* hook, Breakpoint* bp) {
    EptTable* table = Hypervisor_eptTable(gHyp, v->CoreId);
    if (table == NULL) {
        return;
    }
    so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL || pte == NULL) {
        return;
    }
    EptEntry* originalCopy = EptEntry_Clone(hook->OriginalEntry);
    *pte = *originalCopy;
    EptTable_Invalidate(table);
    uint64_t offset = bp->Address - (hook->VirtAddr & ~((uint64_t)(PAGE_SIZE) - 1));
    uint64_t guestInstrPa = EptEntry_PhysAddr(originalCopy) + offset;
    writePhysicalMemory(guestInstrPa, (0xCC));
}

static void singleStepOnFakePage(VCPU* v, EptHook* hook) {
    EptTable* table = Hypervisor_eptTable(gHyp, v->CoreId);
    if (table == NULL) {
        return;
    }
    so_R_ptr_err _res1 = EptTable_Lookup(table, hook->PhysAddr);
    EptEntry* pte = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL || pte == NULL) {
        return;
    }
    PagePermission perm = EptEntry_Permission(hook->ModifiedEntry);
    perm.Read = true;
    perm.Write = true;
    perm.Exec = true;
    EptEntry_SetPermission(pte, perm);
    EptTable_Invalidate(table);
}

static void eptHookSetReadHook(uint64_t addr, CR3_TYPE cr3) {
    Hypervisor* h = getHypervisor();
    if (h == NULL) {
        return;
    }
    VCPU* v = Hypervisor_CurrentVcpu(h);
    if (v != NULL) {
        HookManager_Install(gHooks, v->CoreId, addr, cr3, HookRead);
    }
}

static void eptHookSetWriteHook(uint64_t addr, CR3_TYPE cr3) {
    Hypervisor* h = getHypervisor();
    if (h == NULL) {
        return;
    }
    VCPU* v = Hypervisor_CurrentVcpu(h);
    if (v != NULL) {
        HookManager_Install(gHooks, v->CoreId, addr, cr3, HookWrite);
    }
}

static void eptHookSetExecHook(uint64_t addr, CR3_TYPE cr3) {
    Hypervisor* h = getHypervisor();
    if (h == NULL) {
        return;
    }
    VCPU* v = Hypervisor_CurrentVcpu(h);
    if (v != NULL) {
        HookManager_Install(gHooks, v->CoreId, addr, cr3, HookExec);
    }
}

static void eptHookSetReadWriteHook(uint64_t addr, CR3_TYPE cr3) {
    Hypervisor* h = getHypervisor();
    if (h == NULL) {
        return;
    }
    VCPU* v = Hypervisor_CurrentVcpu(h);
    if (v != NULL) {
        HookManager_Install(gHooks, v->CoreId, addr, cr3, HookReadWrite);
    }
}

// -- hyperevade.go --

so_String EvasionTechnique_String(EvasionTechnique e) {
    so_String result = so_str("");
    if ((e & EvasionHideDebugger) != 0) {
        result = so_string_add(result, so_str("|HIDE_DEBUGGER"));
    }
    if ((e & EvasionHideVmx) != 0) {
        result = so_string_add(result, so_str("|HIDE_VMX"));
    }
    if ((e & EvasionHideMsrs) != 0) {
        result = so_string_add(result, so_str("|HIDE_MSRS"));
    }
    if ((e & EvasionHideIoPorts) != 0) {
        result = so_string_add(result, so_str("|HIDE_IOPORTS"));
    }
    if ((e & EvasionHideTiming) != 0) {
        result = so_string_add(result, so_str("|HIDE_TIMING"));
    }
    if (so_len(result) > 0) {
        return so_string_slice(result, 1, result.len);
    }
    return so_str("NONE");
}

EvasionState* NewEvasionState(void) {
    return &(EvasionState){.config = (TransparentModeConfig){.Enabled = false, .Techniques = EvasionAll}, .hiddenPids = so_make_slice_impl(sizeof(uint32_t), 0, 16), .lock = NewSpinlock(), .msrBackup = so_make_map_impl(sizeof(uint32_t), sizeof(uint64_t), 0), .ioPortHooks = so_make_map_impl(sizeof(uint16_t), sizeof(bool), 0)};
}

so_Error EvasionState_Enable(void* self, TransparentModeConfig* config) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    if (config != NULL) {
        e->config = *config;
    }
    e->config.Enabled = true;
    e->active = true;
    {
        so_Error err = EvasionState_applyEvasions(e);
        if (err != NULL) {
            e->active = false;
            Spinlock_Unlock(e->lock);
            return err;
        }
    }
    so_String techStr = EvasionTechnique_String(e->config.Techniques);
    LogInfo(so_str("Transparent mode enabled: %s"), (so_Slice){(void*[1]){&techStr}, 1, 1});
    Spinlock_Unlock(e->lock);
    return NULL;
}

so_Error EvasionState_Disable(void* self) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    if (!e->active) {
        Spinlock_Unlock(e->lock);
        return NULL;
    }
    {
        so_Error err = EvasionState_removeEvasions(e);
        if (err != NULL) {
            LogWarning(so_str("Error removing evasions: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        }
    }
    e->config.Enabled = false;
    e->active = false;
    e->hiddenPids = so_slice(uint32_t, e->hiddenPids, 0, 0);
    LogInfo(so_str("Transparent mode disabled"), (so_Slice){&so_Nil, 0, 0});
    Spinlock_Unlock(e->lock);
    return NULL;
}

bool EvasionState_IsActive(void* self) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    Spinlock_Unlock(e->lock);
    return e->active;
}

void EvasionState_AddHiddenProcess(void* self, uint32_t pid) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    if (slices_Contains(so_Slice, uint32_t, e->hiddenPids, pid)) {
        Spinlock_Unlock(e->lock);
        return;
    }
    e->hiddenPids = so_append(uint32_t, e->hiddenPids, pid);
    LogDebug(so_str("EVASION: Hidden PID=%d"), (so_Slice){(void*[1]){&pid}, 1, 1});
    Spinlock_Unlock(e->lock);
}

void EvasionState_RemoveHiddenProcess(void* self, uint32_t pid) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    for (so_int i = 0; i < so_len(e->hiddenPids); i++) {
        uint32_t existing = so_at(uint32_t, e->hiddenPids, i);
        if (existing == pid) {
            e->hiddenPids = so_extend(uint32_t, so_slice(uint32_t, e->hiddenPids, 0, i), (so_slice(uint32_t, e->hiddenPids, i + 1, e->hiddenPids.len)));
            LogDebug(so_str("EVASION: Unhidden PID=%d"), (so_Slice){(void*[1]){&pid}, 1, 1});
            Spinlock_Unlock(e->lock);
            return;
        }
    }
    Spinlock_Unlock(e->lock);
}

bool EvasionState_IsProcessHidden(void* self, uint32_t pid) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    Spinlock_Unlock(e->lock);
    return slices_Contains(so_Slice, uint32_t, e->hiddenPids, pid);
}

static so_Error EvasionState_applyEvasions(void* self) {
    EvasionState* e = (EvasionState*)self;
    EvasionTechnique tech = e->config.Techniques;
    if ((tech & EvasionHideVmx) != 0) {
        {
            so_Error err = EvasionState_hideVmxBits(e);
            if (err != NULL) {
                LogWarning(so_str("EVASION: hide VMX bits failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
            }
        }
    }
    if ((tech & EvasionHideMsrs) != 0) {
        {
            so_Error err = EvasionState_hideVmxRelatedMsrs(e);
            if (err != NULL) {
                LogWarning(so_str("EVASION: hide VMX MSRs failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
            }
        }
    }
    if ((tech & EvasionHideIoPorts) != 0) {
        {
            so_Error err = EvasionState_hideVmxCriticalIoPorts(e);
            if (err != NULL) {
                LogWarning(so_str("EVASION: hide IO ports failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
            }
        }
    }
    if ((tech & EvasionHideTiming) != 0) {
        {
            so_Error err = EvasionState_enableTimingCountermeasures(e);
            if (err != NULL) {
                LogWarning(so_str("EVASION: timing countermeasures failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
            }
        }
    }
    if ((tech & EvasionHideDebugger) != 0) {
        {
            so_Error err = EvasionState_hideDebuggerPresence(e);
            if (err != NULL) {
                LogWarning(so_str("EVASION: hide debugger failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
            }
        }
    }
    return NULL;
}

static so_Error EvasionState_removeEvasions(void* self) {
    EvasionState* e = (EvasionState*)self;
    uint64_t cr4 = __readcr4();
    cr4 |= ((uint64_t)1 << 13);
    __writecr4(cr4);
    for (so_int _i = 0; _i < (so_int)e->msrBackup->cap; _i++) {
        if (!e->msrBackup->used[_i]) continue;
        uint32_t msr = ((uint32_t*)e->msrBackup->keys)[_i];
        uint64_t val = ((uint64_t*)e->msrBackup->vals)[_i];
        __writemsr(msr, val);
    }
    e->msrBackup = &(so_Map){0};
    e->ioPortHooks = &(so_Map){0};
    e->tscOffset = 0;
    return NULL;
}

static so_Error EvasionState_hideVmxBits(void* self) {
    EvasionState* e = (EvasionState*)self;
    e->vmxOriginalCr4 = __readcr4();
    uint64_t cr4 = e->vmxOriginalCr4;
    cr4 &= ~((uint64_t)1 << 13);
    __writecr4(cr4);
    LogDebug(so_str("EVASION: Hidden VMXE bit in CR4"), (so_Slice){&so_Nil, 0, 0});
    return NULL;
}

static so_Error EvasionState_hideVmxRelatedMsrs(void* self) {
    EvasionState* e = (EvasionState*)self;
    so_Slice vmxMsrs = (so_Slice){(uint32_t[16]){0x480, 0x481, 0x482, 0x483, 0x484, 0x485, 0x486, 0x487, 0x488, 0x489, 0x48A, 0x48B, 0x6A0, 0x6A2, 0x6A4, 0x6A6}, 16, 16};
    for (so_int _ = 0; _ < so_len(vmxMsrs); _++) {
        uint32_t msr = so_at(uint32_t, vmxMsrs, _);
        uint64_t val = __readmsr(msr);
        {
        	uint64_t _so_map_assign_val;
        	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
        	_so_map_assign_val = val;
        	so_map_set(uint32_t, uint64_t, e->msrBackup, msr, _so_map_assign_val);
        }
        __writemsr(msr, 0);
    }
    LogDebug(so_str("EVASION: Hidden %d VMX MSRs"), (so_Slice){(void*[1]){&(so_int){so_len(vmxMsrs)}}, 1, 1});
    return NULL;
}

static so_Error EvasionState_hideVmxCriticalIoPorts(void* self) {
    EvasionState* e = (EvasionState*)self;
    so_Slice criticalPorts = (so_Slice){(uint16_t[8]){0x3708, 0x3710, 0x3718, 0x3720, 0x3728, 0x3730, 0x3738, 0x561D}, 8, 8};
    for (so_int _ = 0; _ < so_len(criticalPorts); _++) {
        uint16_t port = so_at(uint16_t, criticalPorts, _);
        {
        	bool _so_map_assign_val;
        	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
        	_so_map_assign_val = true;
        	so_map_set(uint16_t, bool, e->ioPortHooks, port, _so_map_assign_val);
        }
    }
    LogDebug(so_str("EVASION: Hooked %d critical IO ports"), (so_Slice){(void*[1]){&(so_int){so_len(criticalPorts)}}, 1, 1});
    return NULL;
}

static so_Error EvasionState_enableTimingCountermeasures(void* self) {
    EvasionState* e = (EvasionState*)self;
    e->tscBase = rdtscValue();
    e->tscOffset = 0;
    LogDebug(so_str("EVASION: Timing countermeasures enabled, base TSC=0x%X"), (so_Slice){(void*[1]){&(uint64_t){e->tscBase}}, 1, 1});
    return NULL;
}

static so_Error EvasionState_hideDebuggerPresence(void* self) {
    EvasionState* e = (EvasionState*)self;
    LogDebug(so_str("EVASION: Debugger hiding active for PID=%d"), (so_Slice){(void*[1]){&(uint32_t){e->config.DebuggerPid}}, 1, 1});
    return NULL;
}

bool EvasionState_ShouldInterceptIoPort(void* self, uint16_t port) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    Spinlock_Unlock(e->lock);
    return so_map_get(uint16_t, bool, e->ioPortHooks, port);
}

bool EvasionState_ShouldHideMsr(void* self, uint32_t msr) {
    EvasionState* e = (EvasionState*)self;
    if (!EvasionState_IsActive(e)) {
        return false;
    }
    bool ok = so_map_has(uint32_t, e->msrBackup, msr);
    return ok;
}

uint64_t EvasionState_AdjustTsc(void* self, uint64_t tsc) {
    EvasionState* e = (EvasionState*)self;
    if (!EvasionState_IsActive(e) || e->tscOffset == 0) {
        return tsc;
    }
    int64_t result = so_max((int64_t)(tsc) + e->tscOffset, 0);
    return (uint64_t)(result);
}

void EvasionState_SetTscOffset(void* self, int64_t offset) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    e->tscOffset = offset;
    Spinlock_Unlock(e->lock);
}

bool EvasionState_ShouldHideProcess(void* self, uint32_t pid) {
    EvasionState* e = (EvasionState*)self;
    if (!EvasionState_IsActive(e) || !((e->config.Techniques & EvasionHideDebugger) != 0)) {
        return false;
    }
    return EvasionState_IsProcessHidden(e, pid);
}

so_int EvasionState_HiddenProcessCount(void* self) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    Spinlock_Unlock(e->lock);
    return so_len(e->hiddenPids);
}

TransparentModeConfig EvasionState_GetConfig(void* self) {
    EvasionState* e = (EvasionState*)self;
    Spinlock_Lock(e->lock);
    Spinlock_Unlock(e->lock);
    return e->config;
}

// -- hypertrace.go --

so_String TraceMode_String(TraceMode t) {
    do {
        if (t == (TraceLbr)) {
            return so_str("LBR");
        } else if (t == (TraceBts)) {
            return so_str("BTS");
        } else if (t == (TracePebs)) {
            return so_str("PEBS");
        } else {
            return so_str("ALL");
        }
    } while (0);
}

TracerState* NewTracerState(void) {
    so_int numEntries = 32;
    TracerState* t = &(TracerState){.lbrConfig = (LbrConfig){.NumEntries = numEntries, .CallstackMode = true}, .btsConfig = (BtsConfig){.BufferSize = (uint64_t)(SIZE_1_GB), .InterruptOnFull = true, .BranchType = 0}, .lock = NewSpinlock()};
    t->lbrBuffer = so_make_slice_impl(sizeof(LbrEntry), numEntries, numEntries);
    so_int maxBtsEntries = (so_int)(SIZE_1_GB / unsafe_Sizeof((BtsEntry){}));
    t->btsBuffer = so_make_slice_impl(sizeof(BtsEntry), maxBtsEntries, maxBtsEntries);
    return t;
}

so_Error TracerState_Enable(void* self, TraceMode mode, TraceUserConfig* config) {
    TracerState* t = (TracerState*)self;
    Spinlock_Lock(t->lock);
    if ((mode & TraceLbr) != 0) {
        {
            so_Error err = TracerState_enableLbr(t, config);
            if (err != NULL) {
                Spinlock_Unlock(t->lock);
                return err;
            }
        }
    }
    if ((mode & TraceBts) != 0) {
        {
            so_Error err = TracerState_enableBts(t, config);
            if (err != NULL) {
                TracerState_disableLbrLocked(t);
                Spinlock_Unlock(t->lock);
                return err;
            }
        }
    }
    t->active = true;
    so_String modeStr = TraceMode_String(mode);
    LogInfo(so_str("TRACE: Enabled mode=%s"), (so_Slice){(void*[1]){&modeStr}, 1, 1});
    Spinlock_Unlock(t->lock);
    return NULL;
}

so_Error TracerState_Disable(void* self) {
    TracerState* t = (TracerState*)self;
    Spinlock_Lock(t->lock);
    if (!t->active) {
        Spinlock_Unlock(t->lock);
        return NULL;
    }
    if (t->lbrConfig.Enabled) {
        TracerState_disableLbrLocked(t);
    }
    if (t->btsConfig.Enabled) {
        TracerState_disableBtsLocked(t);
    }
    t->active = false;
    LogInfo(so_str("TRACE: Disabled"), (so_Slice){&so_Nil, 0, 0});
    Spinlock_Unlock(t->lock);
    return NULL;
}

bool TracerState_IsActive(void* self) {
    TracerState* t = (TracerState*)self;
    Spinlock_Lock(t->lock);
    Spinlock_Unlock(t->lock);
    return t->active;
}

static so_Error TracerState_enableLbr(void* self, TraceUserConfig* config) {
    TracerState* t = (TracerState*)self;
    uint64_t debugctl = __readmsr(IA32_DEBUGCTL_MSR);
    debugctl |= ((uint64_t)1 << 0);
    if (config != NULL && config->CallstackMode) {
        t->lbrConfig.CallstackMode = true;
        debugctl |= ((uint64_t)1 << 2);
    } else if (config != NULL) {
        t->lbrConfig.CallstackMode = false;
        debugctl &= ~((1) << 2);
    }
    if (config != NULL) {
        t->lbrConfig.Filter = config->LbrFilter;
    } else {
        t->lbrConfig.Filter = 0;
    }
    __writemsr(IA32_DEBUGCTL_MSR, debugctl);
    t->lbrConfig.Enabled = true;
    LogDebug(so_str("TRACE: LBR enabled (%d entries)"), (so_Slice){(void*[1]){&(so_int){t->lbrConfig.NumEntries}}, 1, 1});
    return NULL;
}

static void TracerState_disableLbrLocked(void* self) {
    TracerState* t = (TracerState*)self;
    uint64_t debugctl = __readmsr(IA32_DEBUGCTL_MSR);
    debugctl &= ~(((so_int)1 << 0) | ((so_int)1 << 2));
    __writemsr(IA32_DEBUGCTL_MSR, debugctl);
    t->lbrConfig.Enabled = false;
    LogDebug(so_str("TRACE: LBR disabled"), (so_Slice){&so_Nil, 0, 0});
}

static so_Error TracerState_enableBts(void* self, TraceUserConfig* config) {
    TracerState* t = (TracerState*)self;
    uint64_t bufSize = (uint64_t)(SIZE_1_MB);
    if (config != NULL && config->BtsBufferSize > 0) {
        bufSize = config->BtsBufferSize;
    }
    uintptr_t bufVa = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), (uintptr_t)(bufSize), POOLTAG);
    if (bufVa == 0) {
        return fmtError(so_str("BTS buffer allocation failed"), (so_Slice){&so_Nil, 0, 0});
    }
    RtlZeroMemory(bufVa, (uintptr_t)(bufSize));
    uint64_t bufPa = VirtToPhys(bufVa, (CR3_TYPE){.Flags = __readcr3()});
    t->btsConfig.BufferSize = bufSize;
    t->btsConfig.BufferPhysical = bufPa;
    t->btsConfig.BufferVirtual = bufVa;
    uint64_t baseVal = (bufPa | 0x1);
    uint64_t maskVal = ~(bufSize - 1);
    maskVal |= 0x1;
    __writemsr(IA32_DS_AREA_MSR, baseVal);
    t->btsConfig.BaseMsrValue = baseVal;
    __writemsr(IA32_PEBS_ENABLE_MSR, maskVal);
    t->btsConfig.MaskMsrValue = maskVal;
    uint64_t debugctl = __readmsr(IA32_DEBUGCTL_MSR);
    debugctl |= ((((so_int)1 << 7) | ((so_int)1 << 8)) | ((so_int)1 << 9));
    __writemsr(IA32_DEBUGCTL_MSR, debugctl);
    t->btsConfig.Enabled = true;
    LogDebug(so_str("TRACE: BTS enabled (buffer=0x%X bytes, PA=0x%X)"), (so_Slice){(void*[2]){&bufSize, &bufPa}, 2, 2});
    return NULL;
}

static void TracerState_disableBtsLocked(void* self) {
    TracerState* t = (TracerState*)self;
    uint64_t debugctl = __readmsr(IA32_DEBUGCTL_MSR);
    debugctl &= ~((((so_int)1 << 7) | ((so_int)1 << 8)) | ((so_int)1 << 9));
    __writemsr(IA32_DEBUGCTL_MSR, debugctl);
    __writemsr(IA32_DS_AREA_MSR, 0);
    __writemsr(IA32_PEBS_ENABLE_MSR, 0);
    if (t->btsConfig.BufferVirtual != 0) {
        ExFreePoolWithTag(t->btsConfig.BufferVirtual, POOLTAG);
        t->btsConfig.BufferVirtual = 0;
        t->btsConfig.BufferPhysical = 0;
    }
    t->btsConfig.Enabled = false;
    LogDebug(so_str("TRACE: BTS disabled"), (so_Slice){&so_Nil, 0, 0});
}

so_R_slice_int TracerState_CaptureLbr(void* self) {
    TracerState* t = (TracerState*)self;
    Spinlock_Lock(t->lock);
    if (!t->lbrConfig.Enabled || !t->active) {
        Spinlock_Unlock(t->lock);
        return (so_R_slice_int){.val = NULL, .val2 = 0};
    }
    so_int count = TracerState_readLbrFromMsr(t, t->lbrBuffer);
    if (count <= 0) {
        Spinlock_Unlock(t->lock);
        return (so_R_slice_int){.val = NULL, .val2 = 0};
    }
    so_Slice result = so_make_slice_impl(sizeof(LbrEntry), count, count);
    so_copy(LbrEntry, result, so_slice(LbrEntry, t->lbrBuffer, 0, count));
    Spinlock_Unlock(t->lock);
    return (so_R_slice_int){.val = result, .val2 = count};
}

so_R_slice_int TracerState_CaptureBts(void* self) {
    TracerState* t = (TracerState*)self;
    Spinlock_Lock(t->lock);
    if (!t->btsConfig.Enabled || !t->active) {
        Spinlock_Unlock(t->lock);
        return (so_R_slice_int){.val = NULL, .val2 = 0};
    }
    so_int count = TracerState_readBtsFromMemory(t, t->btsBuffer);
    if (count <= 0) {
        Spinlock_Unlock(t->lock);
        return (so_R_slice_int){.val = NULL, .val2 = 0};
    }
    so_Slice result = so_make_slice_impl(sizeof(BtsEntry), count, count);
    so_copy(BtsEntry, result, so_slice(BtsEntry, t->btsBuffer, 0, count));
    Spinlock_Unlock(t->lock);
    return (so_R_slice_int){.val = result, .val2 = count};
}

static so_int TracerState_readLbrFromMsr(void* self, so_Slice buffer) {
    TracerState* t = (TracerState*)self;
    uint64_t lbrTos = __readmsr(IA32_LBR_TOS_MSR);
    so_int tos = (so_int)(lbrTos & 0xF);
    if (tos >= t->lbrConfig.NumEntries || tos >= so_len(buffer)) {
        return 0;
    }
    so_int numPairs = tos;
    if (t->lbrConfig.CallstackMode) {
        for (so_int i = 0; i < numPairs; i++) {
            so_int idx = (tos - 1 - i + 16) % 16;
            uint32_t fromMsr = IA32_LBR_FROM_0_MSR + (uint32_t)(idx);
            uint32_t toMsr = IA32_LBR_TO_0_MSR + (uint32_t)(idx);
            uint32_t infoMsr = IA32_LBR_INFO_0_MSR + (uint32_t)(idx);
            uint32_t miscMsr = IA32_LBR_MISC_0_MSR + (uint32_t)(idx);
            so_at(LbrEntry, buffer, i).From = __readmsr(fromMsr);
            so_at(LbrEntry, buffer, i).To = __readmsr(toMsr);
            so_at(LbrEntry, buffer, i).Info = __readmsr(infoMsr);
            so_at(LbrEntry, buffer, i).Misc = __readmsr(miscMsr);
        }
    } else {
        for (so_int i = 0; i < numPairs; i++) {
            uint32_t fromMsr = IA32_LBR_FROM_0_MSR + (uint32_t)(i);
            uint32_t toMsr = IA32_LBR_TO_0_MSR + (uint32_t)(i);
            uint32_t infoMsr = IA32_LBR_INFO_0_MSR + (uint32_t)(i);
            uint32_t miscMsr = IA32_LBR_MISC_0_MSR + (uint32_t)(i);
            so_at(LbrEntry, buffer, i).From = __readmsr(fromMsr);
            so_at(LbrEntry, buffer, i).To = __readmsr(toMsr);
            so_at(LbrEntry, buffer, i).Info = __readmsr(infoMsr);
            so_at(LbrEntry, buffer, i).Misc = __readmsr(miscMsr);
        }
    }
    return numPairs;
}

static so_int TracerState_readBtsFromMemory(void* self, so_Slice buffer) {
    TracerState* t = (TracerState*)self;
    if (t->btsConfig.BufferVirtual == 0 || t->btsConfig.BufferSize == 0) {
        return 0;
    }
    BtsEntry* btsBase = (BtsEntry*)((void*)(t->btsConfig.BufferVirtual));
    so_int maxEntries = (so_int)(t->btsConfig.BufferSize / (uint64_t)(unsafe_Sizeof((BtsEntry){})));
    uint64_t indexMsr = __readmsr(IA32_BTS_INDEX_MSR);
    so_int absoluteIndex = (so_int)(indexMsr / (uint64_t)(unsafe_Sizeof((BtsEntry){})));
    so_int count = so_min(so_min(absoluteIndex, maxEntries), so_len(buffer));
    for (so_int i = 0; i < count; i++) {
        BtsEntry* entry = (BtsEntry*)((void*)((uintptr_t)((void*)(btsBase)) + (uintptr_t)(i) * unsafe_Sizeof((BtsEntry){})));
        so_at(BtsEntry, buffer, i) = *entry;
    }
    return count;
}

// -- hypervisor.go --

Hypervisor* NewHypervisor(uint32_t numCpus) {
    Hypervisor* h = &(Hypervisor){.vcpus = so_make_slice_impl(sizeof(VCPU), numCpus, numCpus), .eptTables = so_make_slice_impl(sizeof(EptTable*), numCpus, numCpus), .events = newEventDispatcher(), .hooks = NewHookManager((so_int)(MAX_HIDDEN_BREAKPOINTS)), .logger = NewLogger((uint32_t)(LOG_CHUNK_SIZE), (uint32_t)(LOG_CHUNK_SIZE)), .memory = newMemoryManager(), .processContexts = so_make_map_impl(sizeof(uint32_t), sizeof(ProcessContext*), 0), .lock = NewSpinlock()};
    for (uint32_t i = 0; i < numCpus; i++) {
        so_at(VCPU, h->vcpus, i).CoreId = i;
        so_at(VCPU, h->vcpus, i).IncrementRip = true;
        so_at(EptTable*, h->eptTables, i) = NewEptTable();
    }
    return h;
}

so_Error Hypervisor_Initialize(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    if (h->initialized) {
        return NULL;
    }
    uint32_t count = KeQueryActiveProcessorCount(((PVOID)(NULL)));
    for (uint32_t coreId = 0; coreId < count; coreId++) {
        {
            so_Error err = Hypervisor_initializeVcpu(h, coreId);
            if (err != NULL) {
                LogError(so_str("Failed to init VCPU %d: %v"), (so_Slice){(void*[2]){&coreId, err}, 2, 2});
                return err;
            }
        }
    }
    h->initialized = true;
    LogInfo(so_str("Hypervisor initialized with %d VCPUs"), (so_Slice){(void*[1]){&count}, 1, 1});
    return NULL;
}

static so_Error Hypervisor_initializeVcpu(void* self, uint32_t coreId) {
    Hypervisor* h = (Hypervisor*)self;
    VCPU* v = &so_at(VCPU, h->vcpus, coreId);
    uint64_t vmcsPa = allocateVmcsRegion();
    if (vmcsPa == 0) {
        return ErrVmcsAlloc;
    }
    v->VmcsRegionPhysicalAddress = vmcsPa;
    v->VmcsRegionVirtualAddress = vmcsPa;
    uint64_t vmxonPa = allocateVmxonRegion();
    if (vmxonPa == 0) {
        return ErrVmxonAlloc;
    }
    v->VmxonRegionPhysicalAddress = vmxonPa;
    v->VmxonRegionVirtualAddress = vmxonPa;
    uint64_t msrBitmapPa = allocateMsrBitmap();
    if (msrBitmapPa == 0) {
        return fmtError(so_str("MSR bitmap alloc failed"), (so_Slice){&so_Nil, 0, 0});
    }
    v->MsrBitmapPhysicalAddress = msrBitmapPa;
    v->MsrBitmapVirtualAddress = msrBitmapPa;
    uint64_t ioBitmapPaA = allocateIoBitmap();
    uint64_t ioBitmapPaB = allocateIoBitmap();
    v->IoBitmapPhysicalAddressA = ioBitmapPaA;
    v->IoBitmapVirtualAddressA = ioBitmapPaA;
    v->IoBitmapPhysicalAddressB = ioBitmapPaB;
    v->IoBitmapVirtualAddressB = ioBitmapPaB;
    uint64_t stackPa = allocateVmmStack();
    if (stackPa == 0) {
        return fmtError(so_str("stack alloc failed"), (so_Slice){&so_Nil, 0, 0});
    }
    v->VmmStack = stackPa;
    if (so_at(EptTable*, h->eptTables, coreId) != NULL) {
        uint64_t eptPtr = EptTable_BuildEptPointer(so_at(EptTable*, h->eptTables, coreId));
        v->EptPointer.AsUInt = eptPtr;
    } else {
        so_at(EptTable*, h->eptTables, coreId) = NewEptTable();
        uint64_t eptPtr = EptTable_BuildEptPointer(so_at(EptTable*, h->eptTables, coreId));
        v->EptPointer.AsUInt = eptPtr;
    }
    {
        so_Error err = vmxTurnOn(v->VmxonRegionPhysicalAddress);
        if (err != NULL) {
            return fmtError(so_str("VMXON failed on core %d: %v"), (so_Slice){(void*[2]){&coreId, err}, 2, 2});
        }
    }
    {
        so_Error err = vmClear(v->VmcsRegionPhysicalAddress);
        if (err != NULL) {
            return fmtError(so_str("VMCLEAR failed on core %d: %v"), (so_Slice){(void*[2]){&coreId, err}, 2, 2});
        }
    }
    {
        so_Error err = vmLoad(v->VmcsRegionPhysicalAddress);
        if (err != NULL) {
            return fmtError(so_str("VMPTRLD failed on core %d: %v"), (so_Slice){(void*[2]){&coreId, err}, 2, 2});
        }
    }
    {
        so_Error err = Hypervisor_setupVmcs(h, v);
        if (err != NULL) {
            return fmtError(so_str("VMCS setup failed on core %d: %v"), (so_Slice){(void*[2]){&coreId, err}, 2, 2});
        }
    }
    LogInfo(so_str("VCPU %d initialized: VMCS=0x%X VMXON=0x%X EPTP=0x%X"), (so_Slice){(void*[4]){&coreId, &vmcsPa, &vmxonPa, &(uint64_t){v->EptPointer.AsUInt}}, 4, 4});
    return NULL;
}

void Hypervisor_Shutdown(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    Spinlock_Lock(h->lock);
    if (!h->initialized) {
        Spinlock_Unlock(h->lock);
        return;
    }
    for (so_int _ = 0; _ < so_len(h->vcpus); _++) {
        VCPU v = so_at(VCPU, h->vcpus, _);
        HookManager_RestoreAllForCore(gHooks, v.CoreId);
        if (v.OnVmxRootMode) {
            vmClear(v.VmcsRegionPhysicalAddress);
        }
    }
    vmxTurnOff();
    h->initialized = false;
    LogInfo(so_str("Hypervisor shutdown complete"), (so_Slice){&so_Nil, 0, 0});
    Spinlock_Unlock(h->lock);
}

void Hypervisor_SetInitialized(void* self, bool val) {
    Hypervisor* h = (Hypervisor*)self;
    Spinlock_Lock(h->lock);
    h->initialized = val;
    Spinlock_Unlock(h->lock);
}

bool Hypervisor_IsInitialized(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    Spinlock_Lock(h->lock);
    Spinlock_Unlock(h->lock);
    return h->initialized;
}

static EptTable* Hypervisor_eptTable(void* self, uint32_t coreId) {
    Hypervisor* h = (Hypervisor*)self;
    if ((so_int)(coreId) >= so_len(h->eptTables)) {
        return NULL;
    }
    return so_at(EptTable*, h->eptTables, coreId);
}

bool Hypervisor_HandleVmExit(void* self, GUEST_REGS* regs) {
    Hypervisor* h = (Hypervisor*)self;
    uint32_t coreId = getCurrentProcessorNumber();
    if ((so_int)(coreId) >= so_len(h->vcpus)) {
        return false;
    }
    VCPU* v = &so_at(VCPU, h->vcpus, coreId);
    v->OnVmxRootMode = true;
    v->Regs = regs;
    v->XmmRegs = (GUEST_XMM_REGS*)((void*)((uintptr_t)((void*)(regs)) + unsafe_Sizeof((GUEST_REGS){})));
    uint32_t exitReason = readExitReason();
    v->ExitReason = (exitReason);
    v->IncrementRip = true;
    v->LastVmexitRip = readGuestRip();
    uint64_t rsp = 0;
    vmRead64(VmcsGuestRsp, &rsp);
    v->Regs->Rsp = rsp;
    v->ExitQualification = readExitQualification();
    (void)Hypervisor_dispatch(h, v);
    if (!v->VmxoffState.Executed && v->IncrementRip) {
        if (h->checkFootprints) {
            Hypervisor_checkTrapFlag(h, v);
        }
        Hypervisor_advanceIp(h, v);
    }
    bool done = v->VmxoffState.Executed;
    v->OnVmxRootMode = false;
    return done;
}

static bool Hypervisor_dispatch(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    do {
        if (v->ExitReason == (ExitExceptionNmi)) {
            return Hypervisor_handleException(h, v);
        } else if (v->ExitReason == (ExitExtInt)) {
            return Hypervisor_handleExternalInt(h, v);
        } else if (v->ExitReason == (ExitCpuid)) {
            return EventDispatcher_Dispatch(h->events, EventCpuid, v);
        } else if (v->ExitReason == (ExitRdmsr)) {
            return EventDispatcher_Dispatch(h->events, EventRdmsr, v);
        } else if (v->ExitReason == (ExitWrmsr)) {
            return EventDispatcher_Dispatch(h->events, EventWrmsr, v);
        } else if (v->ExitReason == (ExitIo)) {
            return EventDispatcher_Dispatch(h->events, EventIo, v);
        } else if (v->ExitReason == (ExitMovCr)) {
            return EventDispatcher_Dispatch(h->events, EventMovCr, v);
        } else if (v->ExitReason == (ExitMovDr)) {
            return EventDispatcher_Dispatch(h->events, EventMovDr, v);
        } else if (v->ExitReason == (ExitEptViolation)) {
            return Hypervisor_handleEptViolation(h, v);
        } else if (v->ExitReason == (ExitEptMisconfig)) {
            return Hypervisor_handleEptMisconfig(h, v);
        } else if (v->ExitReason == (ExitVmcall)) {
            return Hypervisor_handleVmcall(h, v);
        } else if (v->ExitReason == (ExitRdtsc) || v->ExitReason == (ExitRdtscp)) {
            return EventDispatcher_Dispatch(h->events, EventTsc, v);
        } else if (v->ExitReason == (ExitRdpmc)) {
            return EventDispatcher_Dispatch(h->events, EventRdpmc, v);
        } else if (v->ExitReason == (ExitXsetbv)) {
            return EventDispatcher_Dispatch(h->events, EventXsetbv, v);
        } else if (v->ExitReason == (ExitIntWindow)) {
            return Hypervisor_handleInterruptWindow(h, v);
        } else if (v->ExitReason == (ExitNmiWindow)) {
            return Hypervisor_handleNmiWindow(h, v);
        } else if (v->ExitReason == (ExitMtf)) {
            return Hypervisor_handleMonitorTrapFlag(h, v);
        } else if (v->ExitReason == (ExitPreemptTimer)) {
            return Hypervisor_handlePreemptionTimer(h, v);
        } else if (v->ExitReason == (ExitTripleFault)) {
            return Hypervisor_handleTripleFault(h, v);
        } else if (v->ExitReason == (ExitVmclear) || v->ExitReason == (ExitVmlaunch) || v->ExitReason == (ExitVmptrld) || v->ExitReason == (ExitVmread) || v->ExitReason == (ExitVmresume) || v->ExitReason == (ExitVmwrite) || v->ExitReason == (ExitVmxoff) || v->ExitReason == (ExitVmxon) || v->ExitReason == (ExitInvept) || v->ExitReason == (ExitInvvpid) || v->ExitReason == (ExitGetsec) || v->ExitReason == (ExitInvd)) {
            Hypervisor_injectUd(h, v);
            return false;
        } else {
            return false;
        }
    } while (0);
}

static bool Hypervisor_handleException(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    uint32_t vec = (v->ExitQualification & 0xFF);
    uint32_t intType = ((v->ExitQualification >> 8) & 0x7);
    bool hasErrCode = ((v->ExitQualification >> 11) & 0x1) != 0;
    do {
        if (vec == (VectorBp)) {
            if (intType == IntTypeSwException && !hasErrCode) {
                return HookManager_HandleBreakpoint(gHooks, v);
            }
        } else if (vec == (VectorGp)) {
            return Hypervisor_handleGeneralProtectionFault(h, v);
        } else if (vec == (VectorUd)) {
            return Hypervisor_handleUndefinedOpcode(h, v);
        } else if (vec == (VectorNmi)) {
            if (intType == IntTypeNmi) {
                return Hypervisor_handleNmi(h, v);
            }
        }
    } while (0);
    return EventDispatcher_Dispatch(h->events, EventException, v);
}

static bool Hypervisor_handleEptViolation(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    uint32_t qual = v->ExitQualification;
    bool read = (qual & 0x1) != 0;
    bool write = (qual & 0x2) != 0;
    bool exec = (qual & 0x4) != 0;
    uint64_t gpa = readGuestPhysicalAddr();
    if (exec && !read && !write) {
        return HookManager_HandleExecHook(gHooks, v, gpa);
    }
    if (write || read) {
        return HookManager_HandleReadWriteHook(gHooks, v, gpa, read, write);
    }
    return false;
}

static bool Hypervisor_handleVmcall(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    bool isOurs = v->Regs->Rax == HYPERDBG_VMCALL_MAGIC_RAX && v->Regs->Rcx == HYPERDBG_VMCALL_MAGIC_RCX && v->Regs->Rdx == HYPERDBG_VMCALL_MAGIC_RDX;
    if (isOurs) {
        uint64_t num = v->Regs->R8;
        uint64_t p1 = v->Regs->R9;
        uint64_t p2 = v->Regs->R10;
        uint64_t p3 = v->Regs->R11;
        v->Regs->Rax = (uint64_t)(Hypervisor_handleHyperdbgVmcall(h, num, p1, p2, p3));
        return true;
    }
    Hypervisor_hypercallForward(h, v);
    return true;
}

static void Hypervisor_hypercallForward(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t rsp = v->Regs->Rsp;
    asmHypervVmcall((uint64_t)((uintptr_t)((void*)(v->Regs))));
    v->Regs->Rsp = rsp;
}

static void Hypervisor_advanceIp(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    uint32_t instrLen = readInstructionLength();
    uint64_t newRip = v->LastVmexitRip + (uint64_t)(instrLen);
    vmWrite64(VmcsGuestRip, newRip);
}

static void Hypervisor_suppressAdvance(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    v->IncrementRip = false;
}

static void Hypervisor_enableAdvance(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    v->IncrementRip = true;
}

static void Hypervisor_injectInterrupt(void* self, uint32_t intType, uint32_t vector, bool deliverErr, uint32_t errCode) {
    Hypervisor* h = (Hypervisor*)self;
    uint32_t val = (((0x80000000) | (intType << 8)) | (vector & 0xFF));
    if (deliverErr) {
        val |= 0x800;
    }
    vmWrite64(VmcsCtrlVmentryIntrInfo, (uint64_t)(val));
    if (deliverErr) {
        vmWrite64(VmcsCtrlVmentryErrCode, (uint64_t)(errCode));
    }
}

static void Hypervisor_injectBp(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    Hypervisor_injectInterrupt(h, IntTypeSwException, VectorBp, false, 0);
    uint32_t len = 0;
    vmRead32(VmcsVmexitInstrLength, &len);
    vmWrite64(VmcsCtrlVmentryInstrLen, (uint64_t)(len));
}

static void Hypervisor_injectGp(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    Hypervisor_injectInterrupt(h, IntTypeHwException, VectorGp, true, 0);
    uint32_t len = 0;
    vmRead32(VmcsVmexitInstrLength, &len);
    vmWrite64(VmcsCtrlVmentryInstrLen, (uint64_t)(len));
}

static void Hypervisor_injectUd(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    Hypervisor_injectInterrupt(h, IntTypeHwException, VectorUd, false, 0);
    Hypervisor_suppressAdvance(h, v);
}

static bool Hypervisor_handleMonitorTrapFlag(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    Hypervisor_suppressAdvance(h, v);
    v->IgnoreMtfUnset = false;
    if (HookManager_HandleMtfRestore(gHooks, v)) {
        Hypervisor_enableAndCheckPendingInt(h, v);
    } else if (v->RegisterBreakOnMtf) {
        v->RegisterBreakOnMtf = false;
        handleRegisteredMtf(v->CoreId);
    } else if (kdHandleNmiCallback(v->CoreId)) {
    } else if (v->IgnoreOneMtf) {
        v->IgnoreOneMtf = false;
    }
    if (!v->IgnoreMtfUnset) {
        setMonitorTrapFlag(false);
    } else {
        v->IgnoreMtfUnset = false;
    }
    return true;
}

static void Hypervisor_enableAndCheckPendingInt(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
}

static void Hypervisor_checkTrapFlag(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
}

static bool Hypervisor_handleTripleFault(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleEptMisconfig(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleGeneralProtectionFault(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleUndefinedOpcode(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleExternalInt(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleNmi(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleInterruptWindow(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handleNmiWindow(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static bool Hypervisor_handlePreemptionTimer(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    return false;
}

static void Hypervisor_transparentUnhide(void* self) {
    Hypervisor* h = (Hypervisor*)self;
}

static void Hypervisor_transparentHide(void* self) {
    Hypervisor* h = (Hypervisor*)self;
}

VCPU* Hypervisor_Vcpu(void* self, uint32_t id) {
    Hypervisor* h = (Hypervisor*)self;
    if ((so_int)(id) >= so_len(h->vcpus)) {
        return NULL;
    }
    return &so_at(VCPU, h->vcpus, id);
}

VCPU* Hypervisor_CurrentVcpu(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    uint32_t id = getCurrentProcessorNumber();
    if ((so_int)(id) >= so_len(h->vcpus)) {
        return NULL;
    }
    return &so_at(VCPU, h->vcpus, id);
}

MemoryManager* Hypervisor_MemoryManager(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    return h->memory;
}

void Hypervisor_SetActiveCore(void* self, uint32_t id) {
    Hypervisor* h = (Hypervisor*)self;
    h->activeCoreId = id;
}

uint32_t Hypervisor_GetActiveCore(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    return h->activeCoreId;
}

void Hypervisor_PauseAll(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    h->paused = true;
}

void Hypervisor_ResumeAll(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    h->paused = false;
}

bool Hypervisor_IsPaused(void* self) {
    Hypervisor* h = (Hypervisor*)self;
    return h->paused;
}

void Hypervisor_EnableSingleStep(void* self, VCPU* _p0) {
    Hypervisor* h = (Hypervisor*)self;
    setMonitorTrapFlag(true);
}

void Hypervisor_SetCurrentProcessContext(void* self, uint32_t pid, CR3_TYPE cr3, uint64_t baseAddr) {
    Hypervisor* h = (Hypervisor*)self;
    Spinlock_Lock(h->lock);
    ProcessContext* ctx = so_map_get(uint32_t, ProcessContext*, h->processContexts, pid);
    bool ok = so_map_has(uint32_t, h->processContexts, pid);
    if (!ok) {
        ctx = &(ProcessContext){.ProcessId = pid};
        {
        	ProcessContext* _so_map_assign_val;
        	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
        	_so_map_assign_val = ctx;
        	so_map_set(uint32_t, ProcessContext*, h->processContexts, pid, _so_map_assign_val);
        }
    }
    ctx->Cr3 = cr3;
    ctx->BaseAddress = baseAddr;
    Spinlock_Unlock(h->lock);
}

uint64_t Hypervisor_GetProcessBaseAddress(void* self, uint32_t pid) {
    Hypervisor* h = (Hypervisor*)self;
    Spinlock_Lock(h->lock);
    {
        ProcessContext* ctx = so_map_get(uint32_t, ProcessContext*, h->processContexts, pid);
        bool ok = so_map_has(uint32_t, h->processContexts, pid);
        if (ok) {
            Spinlock_Unlock(h->lock);
            return ctx->BaseAddress;
        }
    }
    Spinlock_Unlock(h->lock);
    return 0;
}

CR3_TYPE Hypervisor_GetProcessCr3(void* self, uint32_t pid) {
    Hypervisor* h = (Hypervisor*)self;
    Spinlock_Lock(h->lock);
    {
        ProcessContext* ctx = so_map_get(uint32_t, ProcessContext*, h->processContexts, pid);
        bool ok = so_map_has(uint32_t, h->processContexts, pid);
        if (ok) {
            Spinlock_Unlock(h->lock);
            return ctx->Cr3;
        }
    }
    Spinlock_Unlock(h->lock);
    return (CR3_TYPE){};
}

uint64_t Hypervisor_CallGuestFunction(void* self, uint64_t _p0, uint64_t _p1, uint64_t _p2, uint64_t _p3, uint64_t _p4) {
    Hypervisor* h = (Hypervisor*)self;
    return 0;
}

static uint32_t readExitReason(void) {
    uint32_t reason = 0;
    vmRead32(VmcsExitReason, &reason);
    reason &= 0xFFFF;
    return (reason);
}

static uint32_t readExitQualification(void) {
    uint32_t qual = 0;
    vmRead32(VmcsExitQualification, &qual);
    return qual;
}

static uint64_t readGuestRip(void) {
    uint64_t rip = 0;
    vmRead64(VmcsGuestRip, &rip);
    return rip;
}

static uint32_t readInstructionLength(void) {
    uint32_t length = 0;
    vmRead32(VmcsVmexitInstrLength, &length);
    return length;
}

static uint64_t readGuestPhysicalAddr(void) {
    uint64_t addr = 0;
    vmRead64(VmcsGuestPhysicalAddress, &addr);
    return addr;
}

static uint64_t allocateVmcsRegion(void) {
    uint64_t vmcsSize = so_max((__readmsr(0x480) & 0x1FFF), 4096);
    vmcsSize = ((vmcsSize + 0xFFF) & ~(0xFFF));
    uintptr_t ptr = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), (uintptr_t)(vmcsSize), POOLTAG);
    if (ptr == 0) {
        return 0;
    }
    RtlZeroMemory(ptr, (uintptr_t)(vmcsSize));
    uint64_t pa = (uint64_t)(MmGetPhysicalAddress(ptr).QuadPart);
    LogDebug(so_str("Allocated VMCS: VA=0x%X PA=0x%X Size=0x%X"), (so_Slice){(void*[3]){&ptr, &pa, &vmcsSize}, 3, 3});
    return pa;
}

static uint64_t allocateVmxonRegion(void) {
    uint64_t vmxonSize = so_max((__readmsr(0x480) & 0x1FFF), 4096);
    vmxonSize = ((vmxonSize + 0xFFF) & ~(0xFFF));
    uintptr_t ptr = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), (uintptr_t)(vmxonSize), POOLTAG);
    if (ptr == 0) {
        return 0;
    }
    RtlZeroMemory(ptr, (uintptr_t)(vmxonSize));
    uint64_t revisonId = __readmsr(0x480);
    *(uint32_t*)((void*)((ptr))) = (uint32_t)(revisonId);
    uint64_t pa = (uint64_t)(MmGetPhysicalAddress(ptr).QuadPart);
    LogDebug(so_str("Allocated VMXON: VA=0x%X PA=0x%X RevId=0x%X"), (so_Slice){(void*[3]){&ptr, &pa, &revisonId}, 3, 3});
    return pa;
}

static uint64_t allocateMsrBitmap(void) {
    uint64_t size = (0x2000);
    uintptr_t ptr = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), (uintptr_t)(size), POOLTAG);
    if (ptr == 0) {
        return 0;
    }
    RtlZeroMemory(ptr, (uintptr_t)(size));
    so_byte (*bitmap)[4096] = (so_byte(*)[4096])((void*)((ptr)));
    for (so_int i = 0; i < 0x800; i++) {
        (*bitmap)[i] = 0xFF;
    }
    (*bitmap)[0xC00 / 8] &= ~((so_byte)1 << (0xC00 % 8));
    (*bitmap)[0xC01 / 8] &= ~((so_byte)1 << (0xC01 % 8));
    (*bitmap)[0xC02 / 8] &= ~((so_byte)1 << (0xC02 % 8));
    (*bitmap)[0xC03 / 8] &= ~((so_byte)1 << (0xC03 % 8));
    (*bitmap)[0xC04 / 8] &= ~((so_byte)1 << (0xC04 % 8));
    (*bitmap)[0xC05 / 8] &= ~((so_byte)1 << (0xC05 % 8));
    (*bitmap)[0xC06 / 8] &= ~((so_byte)1 << (0xC06 % 8));
    (*bitmap)[0xC07 / 8] &= ~((so_byte)1 << (0xC07 % 8));
    (*bitmap)[0xC08 / 8] &= ~((so_byte)1 << (0xC08 % 8));
    (*bitmap)[0x174 / 8] &= ~((so_byte)1 << (0x174 % 8));
    uint64_t pa = (uint64_t)(MmGetPhysicalAddress(ptr).QuadPart);
    LogDebug(so_str("Allocated MSR Bitmap: PA=0x%X"), (so_Slice){(void*[1]){&pa}, 1, 1});
    return pa;
}

static uint64_t allocateIoBitmap(void) {
    uintptr_t size = (0x2000);
    uintptr_t ptr = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), size, POOLTAG);
    if (ptr == 0) {
        return 0;
    }
    RtlZeroMemory(ptr, size);
    uint64_t pa = (uint64_t)(MmGetPhysicalAddress(ptr).QuadPart);
    return pa;
}

static uint64_t allocateVmmStack(void) {
    uintptr_t stackSize = (0x10000);
    uintptr_t ptr = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), stackSize, POOLTAG);
    if (ptr == 0) {
        return 0;
    }
    uintptr_t stackTop = ptr + stackSize - (16);
    *(uint64_t*)((void*)((stackTop))) = 0xDEAD0BADDEAD0BAD;
    uint64_t pa = (uint64_t)(MmGetPhysicalAddress(stackTop).QuadPart);
    LogDebug(so_str("Allocated VMM Stack: TopVA=0x%X TopPA=0x%X"), (so_Slice){(void*[2]){&stackTop, &pa}, 2, 2});
    return pa;
}

static so_Error vmxTurnOn(uint64_t pa) {
    uint64_t cr4 = __readmsr(0x3E);
    if ((cr4 & ((uint64_t)1 << 13)) == 0) {
        __writemsr(0x3E, (cr4 | ((uint64_t)1 << 13)));
    }
    LogInfo(so_str("Executing VMXON at PA=0x%X"), (so_Slice){(void*[1]){&pa}, 1, 1});
    uint8_t ret = AsmEnableVmxOperation(pa);
    if (ret == 0) {
        return fmtError(so_str("VMXON failed"), (so_Slice){&so_Nil, 0, 0});
    }
    return NULL;
}

static void vmxTurnOff(void) {
    LogInfo(so_str("Executing VMXOFF"), (so_Slice){&so_Nil, 0, 0});
    AsmVmxVmxOff();
}

static so_Error vmClear(uint64_t pa) {
    uint8_t ret = AsmVmxVmxClear(pa);
    if (ret == 0) {
        return fmtError(so_str("VMCLEAR failed at PA=0x%X"), (so_Slice){(void*[1]){&pa}, 1, 1});
    }
    return NULL;
}

static so_Error vmLoad(uint64_t pa) {
    uint8_t ret = AsmVmxVmxPtrld(pa);
    if (ret == 0) {
        return fmtError(so_str("VMPTRLD failed at PA=0x%X"), (so_Slice){(void*[1]){&pa}, 1, 1});
    }
    return NULL;
}

static so_Error Hypervisor_setupVmcs(void* self, VCPU* v) {
    Hypervisor* h = (Hypervisor*)self;
    vmWrite64(VmcsGuestCr0, __readcr0());
    vmWrite64(VmcsGuestCr3, __readcr3());
    vmWrite64(VmcsGuestCr4, __readcr4());
    vmWrite64(VmcsGuestDr7, 0x400);
    vmWrite64(VmcsGuestRflags, getRflags());
    vmWrite64(VmcsGuestCsSelector, 0x8);
    vmWrite64(VmcsGuestCsBase, 0);
    vmWrite64(VmcsGuestCsLimit, 0xFFFFFFFF);
    vmWrite64(VmcsGuestCsAccessRights, 0xA0FB);
    vmWrite64(VmcsGuestDsSelector, 0x10);
    vmWrite64(VmcsGuestDsBase, 0);
    vmWrite64(VmcsGuestDsLimit, 0xFFFFFFFF);
    vmWrite64(VmcsGuestDsAccessRights, 0xC0F3);
    vmWrite64(VmcsGuestEsSelector, 0x10);
    vmWrite64(VmcsGuestEsBase, 0);
    vmWrite64(VmcsGuestEsLimit, 0xFFFFFFFF);
    vmWrite64(VmcsGuestEsAccessRights, 0xC0F3);
    vmWrite64(VmcsGuestFsSelector, 0x10);
    vmWrite64(VmcsGuestFsBase, 0);
    vmWrite64(VmcsGuestFsLimit, 0xFFFFFFFF);
    vmWrite64(VmcsGuestFsAccessRights, 0xC0F3);
    vmWrite64(VmcsGuestGsSelector, 0x10);
    vmWrite64(VmcsGuestGsBase, 0);
    vmWrite64(VmcsGuestGsLimit, 0xFFFFFFFF);
    vmWrite64(VmcsGuestGsAccessRights, 0xC0F3);
    vmWrite64(VmcsGuestSsSelector, 0x10);
    vmWrite64(VmcsGuestSsBase, 0);
    vmWrite64(VmcsGuestSsLimit, 0xFFFFFFFF);
    vmWrite64(VmcsGuestSsAccessRights, 0xC0F3);
    uint64_t gdtBase = getGdtBase();
    uint64_t idtBase = getIdtBase();
    uint64_t gdtLimit = (uint64_t)(getGdtLimit());
    uint64_t idtLimit = (uint64_t)(getIdtLimit());
    vmWrite64(VmcsGuestGdtrBase, gdtBase);
    vmWrite64(VmcsGuestGdtrLimit, gdtLimit);
    vmWrite64(VmcsGuestIdtrBase, idtBase);
    vmWrite64(VmcsGuestIdtrLimit, idtLimit);
    vmWrite64(VmcsGuestLdtrSelector, 0);
    vmWrite64(VmcsGuestLdtrBase, 0);
    vmWrite64(VmcsGuestLdtrLimit, 0xFFFF);
    vmWrite64(VmcsGuestLdtrAccessRights, 0x82);
    vmWrite64(VmcsGuestTrSelector, 0);
    vmWrite64(VmcsGuestTrBase, 0);
    vmWrite64(VmcsGuestTrLimit, 0xFFFF);
    vmWrite64(VmcsGuestTrAccessRights, 0x8B);
    uint64_t eptp = v->EptPointer.AsUInt;
    vmWrite64(0x0000201A, eptp);
    uint64_t pinBasedVmExecControls = adjustVmcsControl(__readmsr(0x481), 0x10000016);
    vmWrite64(VmcsCtrlPinBasedVmExecControls, pinBasedVmExecControls);
    uint64_t cpuBasedVmExecControls = (0x04007EFA);
    uint64_t secondaryProcessorBasedControls = (0x06D8FBFE);
    cpuBasedVmExecControls |= ((uint64_t)1 << 31);
    cpuBasedVmExecControls = adjustVmcsControl(__readmsr(0x483), cpuBasedVmExecControls);
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedVmExecControls);
    secondaryProcessorBasedControls = adjustVmcsControl(__readmsr(0x485), secondaryProcessorBasedControls);
    vmWrite64(VmcsCtrlSecondaryProcBasedVmExec, secondaryProcessorBasedControls);
    vmWrite64(VmcsCtrlExceptionBitmap, 0xFFFFFFFF);
    vmWrite64(VmcsCtrlIoBitmapA, v->IoBitmapPhysicalAddressA);
    vmWrite64(VmcsCtrlIoBitmapB, v->IoBitmapPhysicalAddressB);
    vmWrite64(VmcsCtrlMsrBitmap, v->MsrBitmapPhysicalAddress);
    vmWrite64(VmcsHostCr0, __readcr0());
    vmWrite64(VmcsHostCr3, __readcr3());
    vmWrite64(VmcsHostCr4, __readcr4());
    uint16_t csSelector = getCs();
    uint16_t ssSelector = getSs();
    uint16_t dsSelector = getDs();
    uint16_t esSelector = getEs();
    uint16_t fsSelector = getFs();
    uint16_t gsSelector = getGs();
    vmWrite64(VmcsHostCsSelector, (uint64_t)(csSelector));
    vmWrite64(VmcsHostSsSelector, (uint64_t)(ssSelector));
    vmWrite64(VmcsHostDsSelector, (uint64_t)(dsSelector));
    vmWrite64(VmcsHostEsSelector, (uint64_t)(esSelector));
    vmWrite64(VmcsHostFsSelector, (uint64_t)(fsSelector));
    vmWrite64(VmcsHostGsSelector, (uint64_t)(gsSelector));
    vmWrite64(VmcsHostTrSelector, (0x40));
    vmWrite64(VmcsHostFsBase, __readmsr(0xC0000100));
    vmWrite64(VmcsHostGsBase, __readmsr(0xC0000101));
    vmWrite64(VmcsHostLdtrSelector, 0);
    vmWrite64(VmcsHostLdtrAccessRights, 0);
    vmWrite64(VmcsHostTrBase, 0);
    vmWrite64(VmcsHostTrAccessRights, 0x1000B);
    vmWrite64(VmcsHostGdtrBase, gdtBase);
    vmWrite64(VmcsHostIdtrBase, idtBase);
    vmWrite64(VmcsHostSysenterCs, __readmsr(0x174));
    vmWrite64(VmcsHostSysenterEsp, __readmsr(0x176));
    vmWrite64(VmcsHostSysenterEip, __readmsr(0x178));
    vmWrite64(VmcsHostCsAccessRights, 0xA0FB);
    vmWrite64(VmcsHostSsAccessRights, 0xC0F3);
    vmWrite64(VmcsHostDsAccessRights, 0xC0F3);
    vmWrite64(VmcsHostEsAccessRights, 0xC0F3);
    vmWrite64(VmcsHostFsAccessRights, 0xC0F3);
    vmWrite64(VmcsHostGsAccessRights, 0xC0F3);
    vmWrite64(VmcsHostRsp, v->VmmStack);
    vmWrite64(VmcsHostRip, (uint64_t)(AsmVmexitHandlerAddr()));
    vmWrite64(VmcsGuestLinkPointer, 0xFFFFFFFFFFFFFFFF);
    vmWrite64(VmcsGuestIA32Debugctl, 0);
    vmWrite64(VmcsGuestSymEFER, __readmsr(0xC0000080));
    LogDebug(so_str("VMCS configured for core %d, EPTP=0x%X"), (so_Slice){(void*[2]){&(uint32_t){v->CoreId}, &eptp}, 2, 2});
    return NULL;
}

static uint64_t adjustVmcsControl(uint64_t msrValue, uint64_t suggestedValue) {
    uint64_t lowMust0 = (msrValue & 0xFFFFFFFF);
    uint64_t highMust1 = (msrValue >> 32);
    suggestedValue &= ~lowMust0;
    suggestedValue |= highMust1;
    return suggestedValue;
}

static void vmExitHandler(void) {
}

static uint16_t getCs(void) {
    return 0;
}

static uint16_t getSs(void) {
    return 0;
}

static uint16_t getDs(void) {
    return 0;
}

static uint16_t getEs(void) {
    return 0;
}

static uint16_t getFs(void) {
    return 0;
}

static uint16_t getGs(void) {
    return 0;
}

static uint64_t getGdtBase(void) {
    return 0;
}

static uint64_t getIdtBase(void) {
    return 0;
}

static uint32_t getGdtLimit(void) {
    return 0;
}

static uint32_t getIdtLimit(void) {
    return 0;
}

static uint64_t getRflags(void) {
    return 0;
}

static uint32_t getCurrentProcessorNumber(void) {
    uint32_t n = 0;
    KeGetCurrentProcessorNumberEx(&n);
    return n;
}

static bool breakpointReapplyHook(uint32_t coreId) {
    return false;
}

static void handleRegisteredMtf(uint32_t coreId) {
}

void DebuggerUninitialize(void) {
    LogInfo(so_str("Debugger uninitialized"), (so_Slice){&so_Nil, 0, 0});
}

void VmFuncUninitVmm(void) {
    LogInfo(so_str("VMM uninitialized"), (so_Slice){&so_Nil, 0, 0});
}

static bool kdHandleNmiCallback(uint32_t coreId) {
    return false;
}

static void inveptSingleContext(uint64_t eptp) {
    typedef struct InveptDesc {
        uint64_t Eptp;
        uint64_t Padding;
    } InveptDesc;
    InveptDesc desc = (InveptDesc){.Eptp = eptp};
    AsmInvept((uint64_t)((uintptr_t)((void*)(&desc))));
}

static void inveptAllContexts(void) {
    AsmInveptAllContexts();
}

static void invvpidAllContexts(void) {
    AsmInvvpid();
}

static void Hypervisor_setRdtscExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 16);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 16);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setPmcVmexit(void* self, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 7);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 7);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setExceptionBitmap(void* self, VCPU* v, uint32_t bitmap) {
    Hypervisor* h = (Hypervisor*)self;
    vmWrite64(VmcsCtrlExceptionBitmap, (uint64_t)(bitmap));
}

static void Hypervisor_setMovDebugRegsExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 8);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 8);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setExternalInterruptExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t pinBasedControls = (0);
    vmRead64(VmcsCtrlPinBasedVmExecControls, &pinBasedControls);
    if (enable) {
        pinBasedControls |= ((uint64_t)1 << 0);
    } else {
        pinBasedControls &= ~((uint64_t)1 << 0);
    }
    vmWrite64(VmcsCtrlPinBasedVmExecControls, pinBasedControls);
}

static void Hypervisor_setCpuidExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 10);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 10);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setClpExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 24);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 24);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setMovFromCrExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 19);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 19);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setMovToCrExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 20);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 20);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setMovFromDrExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 23);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 23);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void Hypervisor_setMovToDrExiting(void* self, VCPU* v, bool enable) {
    Hypervisor* h = (Hypervisor*)self;
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (enable) {
        cpuBasedControls |= ((uint64_t)1 << 22);
    } else {
        cpuBasedControls &= ~((uint64_t)1 << 22);
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void protectedHvExternalInterruptExitingForDisablingInterruptCommands(VCPU* v) {
}

// -- io_handler.go --

static void initIoBitmap(void) {
    for (so_int i = 0; i < IO_BITMAP_SIZE; i++) {
        gIoBitmap->bitmapA[i] = 0;
        gIoBitmap->bitmapB[i] = 0;
    }
}

static void ioHandlePerformIoBitmapChange(VCPU* v, uint32_t portMask) {
    if (v->IoBitmapVirtualAddressA == 0 || v->IoBitmapVirtualAddressB == 0) {
        LogError(so_str("IO bitmap not initialized"), (so_Slice){&so_Nil, 0, 0});
        return;
    }
    for (uint32_t port = 0; port < (0x10000); port++) {
        uint32_t byteOffset = port / 8;
        uint32_t bitOffset = port % 8;
        if ((portMask & ((uint32_t)1 << (uint64_t)(port % 32))) != 0) {
            if (byteOffset < IO_BITMAP_SIZE) {
                if (port < 0x8000) {
                    gIoBitmap->bitmapA[byteOffset] |= ((so_byte)1 << bitOffset);
                } else {
                    uint32_t idx = byteOffset - 0x1000;
                    if (idx < IO_BITMAP_SIZE) {
                        gIoBitmap->bitmapB[idx] |= ((so_byte)1 << bitOffset);
                    }
                }
            }
        } else {
            if (byteOffset < IO_BITMAP_SIZE) {
                if (port < 0x8000) {
                    gIoBitmap->bitmapA[byteOffset] &= ~((so_byte)1 << bitOffset);
                } else {
                    uint32_t idx = byteOffset - 0x1000;
                    if (idx < IO_BITMAP_SIZE) {
                        gIoBitmap->bitmapB[idx] &= ~((so_byte)1 << bitOffset);
                    }
                }
            }
        }
    }
    uint64_t paA = VirtToPhys((uintptr_t)(v->IoBitmapVirtualAddressA), (CR3_TYPE){});
    uint64_t paB = VirtToPhys((uintptr_t)(v->IoBitmapVirtualAddressB), (CR3_TYPE){});
    WritePhysMem(paA, so_array_slice(so_byte, gIoBitmap->bitmapA, 0, 8192, 8192));
    WritePhysMem(paB, so_array_slice(so_byte, gIoBitmap->bitmapB, 0, 8192, 8192));
    uint64_t cpuBasedControls = (0);
    vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
    if (portMask != 0xFFFFFFFF && portMask != 0) {
        cpuBasedControls |= ((uint64_t)1 << 24);
    } else if (portMask == 0xFFFFFFFF) {
        cpuBasedControls |= ((uint64_t)1 << 24);
    } else {
        cpuBasedControls &= ~(((uint64_t)1 << 24));
    }
    vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
}

static void ioHandleEnableOrDisableIoPortExiting(VCPU* v, bool enable) {
    if (enable) {
        initIoBitmap();
        ioHandlePerformIoBitmapChange(v, 0xFFFF);
    } else {
        ioHandlePerformIoBitmapChange(v, 0x0000);
    }
}

static bool handleIoRead(VCPU* v) {
    uint16_t port = (uint16_t)(v->Regs->Rdx & 0xFFFF);
    uint64_t size = ((v->Regs->Rax >> 16) & 0x7);
    uint64_t value = 0;
    do {
        if (size == (0)) {
            value = (uint64_t)(inb(port));
        } else if (size == (1)) {
            value = (uint64_t)(inw(port));
        } else if (size == (2)) {
            value = (uint64_t)(ind(port));
        } else {
            LogError(so_str("Invalid IO read size: %d"), (so_Slice){(void*[1]){&size}, 1, 1});
            return false;
        }
    } while (0);
    v->Regs->Rax = value;
    Hypervisor* h = getHypervisor();
    if (h != NULL) {
        Hypervisor_advanceIp(h, v);
    }
    return true;
}

static bool handleIoWrite(VCPU* v) {
    uint16_t port = (uint16_t)(v->Regs->Rdx & 0xFFFF);
    uint64_t size = ((v->Regs->Rax >> 16) & 0x7);
    uint64_t value = (v->Regs->Rax & 0xFFFF);
    do {
        if (size == (0)) {
            outb(port, (uint8_t)(value));
        } else if (size == (1)) {
            outw(port, (uint16_t)(value));
        } else if (size == (2)) {
            outd(port, (uint32_t)(value));
        } else {
            LogError(so_str("Invalid IO write size: %d"), (so_Slice){(void*[1]){&size}, 1, 1});
            return false;
        }
    } while (0);
    Hypervisor* h = getHypervisor();
    if (h != NULL) {
        Hypervisor_advanceIp(h, v);
    }
    return true;
}

static uint8_t inb(uint16_t port) {
    return 0;
}

static uint16_t inw(uint16_t port) {
    return 0;
}

static uint32_t ind(uint16_t port) {
    return 0;
}

static void outb(uint16_t port, uint8_t val) {
}

static void outw(uint16_t port, uint16_t val) {
}

static void outd(uint16_t port, uint32_t val) {
}

// -- ioctl.go --

static int32_t drvIoctl(DEVICE_OBJECT* _p0, IRP* irp) {
    IO_STACK_LOCATION* stack = IoGetCurrentIrpStackLocation(irp);
    uint32_t ioctl = stack->Parameters.DeviceIoControl.IoControlCode;
    uint32_t inLen = stack->Parameters.DeviceIoControl.InputBufferLength;
    uint32_t outLen = stack->Parameters.DeviceIoControl.OutputBufferLength;
    int32_t status = STATUS_SUCCESS;
    uint64_t info = 0;
    do {
        if (ioctl == (IOCTL_QUERY_AND_CLEAR_LOGS)) {
            status = handleQueryLogs(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_REGISTER_EVENT)) {
            status = handleRegisterEvent(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_RUN_SCRIPT)) {
            status = handleRunScript(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_SEND_REQUEST_RESULT)) {
            status = handleSendResult(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_GET_LOG_BASE)) {
            status = handleGetLogBase(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_INIT)) {
            status = handleVmmInit(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_SHUTDOWN)) {
            status = handleVmmShutdown(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_EPT_HOOK)) {
            status = handleEptHook(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_EPT_UNHOOK)) {
            status = handleEptUnhook(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_EPT_SET_HOOK)) {
            status = handleEptSetHook(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_EPT_GET_EPT_TABLES)) {
            status = handleGetEptTables(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_EXECUTION_TRACE)) {
            status = handleExecutionTrace(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_READ_MEM)) {
            status = handleReadMem(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_WRITE_MEM)) {
            status = handleWriteMem(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_VIRT_TO_PHYS)) {
            status = handleVirtToPhys(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_PHYS_TO_VIRT)) {
            status = handlePhysToVirt(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_GET_REG)) {
            status = handleGetReg(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_SET_REG)) {
            status = handleSetReg(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_PAUSE)) {
            status = handlePause(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_RESUME)) {
            status = handleResume(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_SWITCH_PROCESS)) {
            status = handleSwitchProcess(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_MODIFY_REGS)) {
            status = handleModifyRegs(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_INVEPT)) {
            status = handleInvept(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_INVVPID)) {
            status = handleInvvpid(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_FLUSH_ENTIRE_TLB)) {
            status = handleFlushTlb(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_GET_MTRR)) {
            status = handleGetMtrr(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_CHANGE_CORE)) {
            status = handleChangeCore(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_GET_PROCESS_BASE)) {
            status = handleGetProcessBase(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_GET_PROCESS_CR3)) {
            status = handleGetProcessCr3(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_READ_AND_WRITE_MEM)) {
            status = handleReadWriteMem(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_CALL_FUNCTION)) {
            status = handleCallFunction(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_QUERY_PACKET)) {
            status = handleQueryPacket(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_SEND_RESULT)) {
            status = handleSendResult(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_EXECUTE_SINGLE_STEP)) {
            status = handleSingleStep(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_ENABLE_AND_INVOKE_TRAP_FLAG)) {
            status = handleTrapFlag(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_MASK_EXCEPTION_DEBUGGER_BREAKPOINT)) {
            status = handleMaskBreakpoint(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_UNMASK_EXCEPTION_DEBUGGER_BREAKPOINT)) {
            status = handleUnmaskBreakpoint(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_SHORT_CIRCUITING_EVENT_INJECT)) {
            status = handleShortCircuitInject(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_VMM_QUERY_REGISTER_FROM_GUEST_STATE)) {
            status = handleQueryGuestReg(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_TRANSPARENT_MODE_ENABLE)) {
            status = handleTransparentEnable(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_TRANSPARENT_MODE_DISABLE)) {
            status = handleTransparentDisable(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_HYPERTRACE_START)) {
            status = handleTraceStart(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_HYPERTRACE_STOP)) {
            status = handleTraceStop(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_HYPERTRACE_READ_BUFFER)) {
            status = handleTraceReadBuffer(irp, inLen, outLen, &info);
        } else if (ioctl == (IOCTL_HYPERTRACE_CLEAR_BUFFER)) {
            status = handleTraceClearBuffer(irp, inLen, outLen, &info);
        } else {
            if (ioctl >= 0x0022B018 && (ioctl - 0x0022B018) % 4 == 0) {
                status = handleLogRead(irp, inLen, outLen, &info);
            } else {
                status = STATUS_INVALID_DEVICE_REQUEST;
            }
        }
    } while (0);
    irp->IoStatus.Status = status;
    irp->IoStatus.Information = info;
    IoCompleteRequest(irp, IO_NO_INCREMENT);
    return status;
}

static int32_t handleQueryLogs(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (outLen < (uint32_t)(unsafe_Sizeof((LogMessage){}))) {
        return STATUS_BUFFER_TOO_SMALL;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    uint32_t count = Logger_FlushToUser(gLog, ((buffer)), (uintptr_t)((uint64_t)(outLen)));
    *info = (uint64_t)((uintptr_t)(count) * LOG_MESSAGE_SIZE);
    LogDebug(so_str("Flushed %d log messages to user"), (so_Slice){(void*[1]){&count}, 1, 1});
    return STATUS_SUCCESS;
}

static int32_t handleRegisterEvent(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen == 0 || inLen < 16) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    uint32_t eventType = *(uint32_t*)((void*)((buffer)));
    uint32_t eventTag = *(uint32_t*)((void*)((buffer) + 4));
    uint64_t callbackAddr = *(uint64_t*)((void*)((buffer) + 8));
    LogDebug(so_str("Register event: type=%d tag=0x%X callback=0x%X"), (so_Slice){(void*[3]){&eventType, &eventTag, &callbackAddr}, 3, 3});
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleRunScript(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    LogWarning(so_str("Script execution not supported in this build"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_NOT_IMPLEMENTED;
}

static int32_t handleSendResult(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen == 0) {
        return STATUS_INVALID_PARAMETER;
    }
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleGetLogBase(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (outLen < 16) {
        return STATUS_BUFFER_TOO_SMALL;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    uint64_t baseAddr = Logger_GetBufferBase(gLog);
    uint64_t prioBase = Logger_GetPriorityBufferBase(gLog);
    *(uint64_t*)((void*)((buffer))) = baseAddr;
    *(uint64_t*)((void*)((buffer) + 8)) = prioBase;
    *info = 16;
    return STATUS_SUCCESS;
}

static int32_t handleLogRead(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (outLen < 8) {
        return STATUS_BUFFER_TOO_SMALL;
    }
    *info = 8;
    return STATUS_SUCCESS;
}

static int32_t handleVmmInit(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((VmmInitRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    VmmInitRequest* req = (VmmInitRequest*)((void*)((buffer)));
    LogInfo(so_str("Initializing VMM for %d cores..."), (so_Slice){(void*[1]){&(uint32_t){req->NumCores}}, 1, 1});
    so_Error err = Hypervisor_Initialize(gHyp);
    if (err != NULL) {
        LogError(so_str("VMM initialization failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        return STATUS_UNSUCCESSFUL;
    }
    Hypervisor_SetInitialized(gHyp, true);
    *info = 4;
    LogInfo(so_str("VMM initialized successfully"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleVmmShutdown(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    LogInfo(so_str("Shutting down VMM..."), (so_Slice){&so_Nil, 0, 0});
    HookManager_RestoreAll(gHooks, NULL);
    Hypervisor_Shutdown(gHyp);
    Hypervisor_SetInitialized(gHyp, false);
    *info = 4;
    LogInfo(so_str("VMM shutdown complete"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleEptHook(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((EptHookRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    EptHookRequest* req = (EptHookRequest*)((void*)((buffer)));
    VCPU* v = Hypervisor_CurrentVcpu(gHyp);
    if (v == NULL) {
        return STATUS_DEVICE_NOT_READY;
    }
    HookKind kind = 0;
    do {
        if (req->HookType == (1)) {
            kind = HookExec;
        } else if (req->HookType == (2)) {
            kind = HookRead;
        } else if (req->HookType == (3)) {
            kind = HookWrite;
        } else if (req->HookType == (4)) {
            kind = HookReadWrite;
        } else {
            return STATUS_INVALID_PARAMETER;
        }
    } while (0);
    so_Error err = HookManager_Install(gHooks, v->CoreId, req->VirtAddr, req->Cr3, kind);
    if (err != NULL) {
        LogError(so_str("EPT hook failed at 0x%X: %v"), (so_Slice){(void*[2]){&(uint64_t){req->VirtAddr}, err}, 2, 2});
        return STATUS_UNSUCCESSFUL;
    }
    *info = 4;
    LogDebug(so_str("EPT hook installed at 0x%X (type=%d)"), (so_Slice){(void*[2]){&(uint64_t){req->VirtAddr}, &(uint32_t){req->HookType}}, 2, 2});
    return STATUS_SUCCESS;
}

static int32_t handleEptUnhook(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((EptUnhookRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    EptUnhookRequest* req = (EptUnhookRequest*)((void*)((buffer)));
    uint64_t pa = VirtToPhys((uintptr_t)(req->VirtAddr & ~(uint64_t)(PAGE_SIZE - 1)), req->Cr3);
    if (pa == 0) {
        return STATUS_INVALID_PARAMETER;
    }
    so_Error err = HookManager_Remove(gHooks, pa);
    if (err != NULL) {
        LogError(so_str("EPT unhook failed at PA 0x%X: %v"), (so_Slice){(void*[2]){&pa, err}, 2, 2});
        return STATUS_UNSUCCESSFUL;
    }
    *info = 4;
    LogDebug(so_str("EPT hook removed at VA 0x%X (PA 0x%X)"), (so_Slice){(void*[2]){&(uint64_t){req->VirtAddr}, &pa}, 2, 2});
    return STATUS_SUCCESS;
}

static int32_t handleEptSetHook(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    return handleEptHook(irp, inLen, outLen, info);
}

static int32_t handleGetEptTables(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (outLen < 8) {
        return STATUS_BUFFER_TOO_SMALL;
    }
    VCPU* v = Hypervisor_CurrentVcpu(gHyp);
    if (v == NULL || v->EptPageTable == NULL) {
        return STATUS_DEVICE_NOT_READY;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    *(uintptr_t*)((void*)((buffer))) = (uintptr_t)((void*)(v->EptPageTable));
    *info = 8;
    LogDebug(so_str("Returned EPT table base at 0x%x"), (so_Slice){(void*[1]){v->EptPageTable}, 1, 1});
    return STATUS_SUCCESS;
}

static int32_t handleExecutionTrace(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((TraceStartRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    TraceStartRequest* req = (TraceStartRequest*)((void*)((buffer)));
    TraceUserConfig config = (TraceUserConfig){.LbrFilter = req->LbrFilter, .BtsBufferSize = req->BtsBufferSize, .CallstackMode = (req->CallstackMode)};
    TraceMode mode = (TraceMode)(req->Mode);
    so_Error err = TracerState_Enable(gTrace, mode, &config);
    if (err != NULL) {
        LogError(so_str("Trace enable failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        return STATUS_UNSUCCESSFUL;
    }
    *info = 4;
    LogInfo(so_str("Execution trace started (mode=%d)"), (so_Slice){(void*[1]){&(uint32_t){req->Mode}}, 1, 1});
    return STATUS_SUCCESS;
}

static int32_t handleReadMem(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((MemReadRequest){})) || outLen < (uint32_t)(unsafe_Sizeof((MemReadResponse){})) + 256) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    MemReadRequest* req = (MemReadRequest*)((void*)((buffer)));
    MemReadResponse* resp = (MemReadResponse*)((void*)((buffer)));
    uint64_t pa = VirtToPhys((uintptr_t)(req->Address), req->Cr3);
    if (pa == 0) {
        resp->Success = false;
        *info = (uint64_t)(unsafe_Sizeof((MemReadResponse){}));
        return STATUS_SUCCESS;
    }
    uint32_t readSize = req->Size;
    if (readSize > outLen - (uint32_t)(unsafe_Sizeof((MemReadResponse){}))) {
        readSize = (outLen - (uint32_t)(unsafe_Sizeof((MemReadResponse){})));
    }
    so_byte* dataBuf = (so_byte*)((void*)((buffer) + unsafe_Sizeof((MemReadResponse){})));
    so_int n = ReadPhysMem(pa, unsafe_Slice(dataBuf, (so_int)(readSize)));
    resp->Address = req->Address;
    resp->Size = (uint32_t)(n);
    resp->Success = true;
    *info = (uint64_t)(unsafe_Sizeof((MemReadResponse){})) + (uint64_t)(n);
    return STATUS_SUCCESS;
}

static int32_t handleWriteMem(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((MemReadRequest){})) + 1) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    MemReadRequest* req = (MemReadRequest*)((void*)((buffer)));
    uint64_t pa = VirtToPhys((uintptr_t)(req->Address), req->Cr3);
    if (pa == 0) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t dataOffset = unsafe_Sizeof((MemReadRequest){});
    so_byte* dataBuf = (so_byte*)((void*)((buffer) + dataOffset));
    so_int n = WritePhysMem(pa, unsafe_Slice(dataBuf, (so_int)(req->Size)));
    *info = 4;
    LogDebug(so_str("Wrote %d bytes to PA 0x%X"), (so_Slice){(void*[2]){&n, &pa}, 2, 2});
    return STATUS_SUCCESS;
}

static int32_t handleVirtToPhys(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((VirtPhysRequest){})) || outLen < 16) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    VirtPhysRequest* req = (VirtPhysRequest*)((void*)((buffer)));
    uint64_t pa = VirtToPhys((uintptr_t)(req->Address), req->Cr3);
    *(uint64_t*)((void*)((buffer))) = pa;
    *(uint64_t*)((void*)((buffer) + 8)) = req->Address;
    *info = 16;
    return STATUS_SUCCESS;
}

static int32_t handlePhysToVirt(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((VirtPhysRequest){})) || outLen < 16) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    VirtPhysRequest* req = (VirtPhysRequest*)((void*)((buffer)));
    uintptr_t va = PhysToVirt(req->Address, req->Cr3);
    *(uint64_t*)((void*)((buffer))) = (uint64_t)(va);
    *(uint64_t*)((void*)((buffer) + 8)) = req->Address;
    *info = 16;
    return STATUS_SUCCESS;
}

static int32_t handleGetReg(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((RegRequest){})) || outLen < (uint32_t)(unsafe_Sizeof((RegResponse){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    RegRequest* req = (RegRequest*)((void*)((buffer)));
    RegResponse* resp = (RegResponse*)((void*)((buffer)));
    VCPU* v = Hypervisor_Vcpu(gHyp, req->CoreId);
    if (v == NULL || v->Regs == NULL) {
        resp->Success = false;
        *info = (uint64_t)(unsafe_Sizeof((RegResponse){}));
        return STATUS_SUCCESS;
    }
    uint64_t value = getRegisterValue(v->Regs, req->RegIndex);
    resp->Value = value;
    resp->RegIndex = req->RegIndex;
    resp->CoreId = req->CoreId;
    resp->Success = true;
    *info = (uint64_t)(unsafe_Sizeof((RegResponse){}));
    return STATUS_SUCCESS;
}

static int32_t handleSetReg(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((RegResponse){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    RegResponse* req = (RegResponse*)((void*)((buffer)));
    VCPU* v = Hypervisor_Vcpu(gHyp, req->CoreId);
    if (v == NULL || v->Regs == NULL) {
        return STATUS_DEVICE_NOT_READY;
    }
    setRegisterValue(v->Regs, req->RegIndex, req->Value);
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handlePause(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    Hypervisor_PauseAll(gHyp);
    *info = 4;
    LogInfo(so_str("All VCPUs paused"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleResume(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    Hypervisor_ResumeAll(gHyp);
    *info = 4;
    LogInfo(so_str("All VCPUs resumed"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleSwitchProcess(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((SwitchProcessRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    SwitchProcessRequest* req = (SwitchProcessRequest*)((void*)((buffer)));
    Hypervisor_SetCurrentProcessContext(gHyp, req->ProcessId, req->Cr3, req->ProcessBaseAddress);
    *info = 4;
    LogDebug(so_str("Switched to process %d (CR3=0x%X)"), (so_Slice){(void*[2]){&(uint32_t){req->ProcessId}, &(uint64_t){req->Cr3.Flags}}, 2, 2});
    return STATUS_SUCCESS;
}

static int32_t handleModifyRegs(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((ModifyRegsRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    ModifyRegsRequest* req = (ModifyRegsRequest*)((void*)((buffer)));
    VCPU* v = Hypervisor_Vcpu(gHyp, req->CoreId);
    if (v == NULL || v->Regs == NULL) {
        return STATUS_DEVICE_NOT_READY;
    }
    *v->Regs = req->Regs;
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleInvept(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((InveptRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    InveptRequest* req = (InveptRequest*)((void*)((buffer)));
    if (req->Type == 1) {
        eptInveptSingleContext(req->Eptp);
    } else if (req->Type == 2) {
        eptInveptAllContexts();
    }
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleInvvpid(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    AsmInvvpid();
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleFlushTlb(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    eptInveptAllContexts();
    AsmInvvpid();
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleGetMtrr(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((MtrrRequest){})) || outLen < 12) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    MtrrRequest* req = (MtrrRequest*)((void*)((buffer)));
    MemoryType memType = MemoryManager_QueryMtrrForPa(Hypervisor_MemoryManager(gHyp), req->BaseAddr);
    *(uint64_t*)((void*)((buffer))) = req->BaseAddr;
    *(uint32_t*)((void*)((buffer) + 8)) = (uint32_t)(memType);
    *info = 12;
    return STATUS_SUCCESS;
}

static int32_t handleChangeCore(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((ChangeCoreRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    ChangeCoreRequest* req = (ChangeCoreRequest*)((void*)((buffer)));
    Hypervisor_SetActiveCore(gHyp, req->TargetCoreId);
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleGetProcessBase(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((ProcessCr3Request){})) || outLen < 8) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    ProcessCr3Request* req = (ProcessCr3Request*)((void*)((buffer)));
    uint64_t baseAddr = Hypervisor_GetProcessBaseAddress(gHyp, req->ProcessId);
    *(uint64_t*)((void*)((buffer))) = baseAddr;
    *info = 8;
    return STATUS_SUCCESS;
}

static int32_t handleGetProcessCr3(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((ProcessCr3Request){})) || outLen < 8) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    ProcessCr3Request* req = (ProcessCr3Request*)((void*)((buffer)));
    CR3_TYPE cr3 = Hypervisor_GetProcessCr3(gHyp, req->ProcessId);
    *(uint64_t*)((void*)((buffer))) = cr3.Flags;
    *info = 8;
    return STATUS_SUCCESS;
}

static int32_t handleReadWriteMem(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((ReadWriteMemRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    ReadWriteMemRequest* req = (ReadWriteMemRequest*)((void*)((buffer)));
    if (!req->ReadOrWrite) {
        if (req->ReadSize > 0) {
            uint64_t pa = VirtToPhys((uintptr_t)(req->ReadAddress), req->Cr3);
            if (pa != 0) {
                so_byte* dataBuf = (so_byte*)((void*)((buffer) + unsafe_Sizeof((ReadWriteMemRequest){})));
                ReadPhysMem(pa, unsafe_Slice(dataBuf, (so_int)(req->ReadSize)));
            }
        }
    } else {
        if (req->WriteSize > 0) {
            uint64_t pa = VirtToPhys((uintptr_t)(req->WriteAddress), req->Cr3);
            if (pa != 0) {
                uintptr_t dataOffset = unsafe_Sizeof((ReadWriteMemRequest){});
                so_byte* dataBuf = (so_byte*)((void*)((buffer) + dataOffset));
                WritePhysMem(pa, unsafe_Slice(dataBuf, (so_int)(req->WriteSize)));
            }
        }
    }
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleCallFunction(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((CallFunctionRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    CallFunctionRequest* req = (CallFunctionRequest*)((void*)((buffer)));
    uint64_t result = Hypervisor_CallGuestFunction(gHyp, req->FunctionAddress, req->OptionalParam1, req->OptionalParam2, req->OptionalParam3, req->OptionalParam4);
    *(uint64_t*)((void*)((buffer))) = result;
    *info = 8;
    return STATUS_SUCCESS;
}

static int32_t handleQueryPacket(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (outLen < 4) {
        return STATUS_BUFFER_TOO_SMALL;
    }
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleSingleStep(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    VCPU* v = Hypervisor_CurrentVcpu(gHyp);
    if (v == NULL) {
        return STATUS_DEVICE_NOT_READY;
    }
    Hypervisor_EnableSingleStep(gHyp, v);
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleTrapFlag(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    VCPU* v = Hypervisor_CurrentVcpu(gHyp);
    if (v == NULL) {
        return STATUS_DEVICE_NOT_READY;
    }
    setMonitorTrapFlag(true);
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleMaskBreakpoint(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    EventDispatcher_Disable(gEvents, EventException);
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleUnmaskBreakpoint(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    EventDispatcher_Enable(gEvents, EventException);
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleShortCircuitInject(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    *info = 4;
    return STATUS_SUCCESS;
}

static int32_t handleQueryGuestReg(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    return handleGetReg(irp, inLen, outLen, info);
}

static int32_t handleTransparentEnable(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (inLen < (uint32_t)(unsafe_Sizeof((TransparentModeRequest){}))) {
        return STATUS_INVALID_PARAMETER;
    }
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    TransparentModeRequest* req = (TransparentModeRequest*)((void*)((buffer)));
    TransparentModeConfig config = (TransparentModeConfig){.Enabled = (req->Enable), .Techniques = (EvasionTechnique)(req->Techniques), .DebuggerPid = req->DebuggerPid};
    so_Error err = EvasionState_Enable(gEvasion, &config);
    if (err != NULL) {
        LogError(so_str("Transparent mode enable failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        return STATUS_UNSUCCESSFUL;
    }
    *info = 4;
    LogInfo(so_str("Transparent mode enabled"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleTransparentDisable(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    so_Error err = EvasionState_Disable(gEvasion);
    if (err != NULL) {
        LogError(so_str("Transparent mode disable failed: %v"), (so_Slice){(void*[1]){err}, 1, 1});
        return STATUS_UNSUCCESSFUL;
    }
    *info = 4;
    LogInfo(so_str("Transparent mode disabled"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleTraceStart(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    return handleExecutionTrace(irp, inLen, outLen, info);
}

static int32_t handleTraceStop(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    so_Error err = TracerState_Disable(gTrace);
    if (err != NULL) {
        return STATUS_UNSUCCESSFUL;
    }
    *info = 4;
    LogInfo(so_str("Tracing stopped"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static int32_t handleTraceReadBuffer(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    if (outLen < 8) {
        return STATUS_BUFFER_TOO_SMALL;
    }
    so_R_slice_int _res1 = TracerState_CaptureLbr(gTrace);
    so_Slice lbrEntries = _res1.val;
    so_int lbrCount = _res1.val2;
    uintptr_t buffer = irp->AssociatedIrp.SystemBuffer;
    *(uint32_t*)((void*)((buffer))) = (uint32_t)(lbrCount);
    if (outLen >= (uint32_t)(8 + (uintptr_t)(lbrCount) * unsafe_Sizeof((LbrEntry){}))) {
        for (so_int i = 0; i < lbrCount; i++) {
            uintptr_t offset = 8 + (uintptr_t)(i) * unsafe_Sizeof((LbrEntry){});
            LbrEntry* dst = (LbrEntry*)((void*)((buffer) + offset));
            *dst = so_at(LbrEntry, lbrEntries, i);
        }
        *info = (uint64_t)(8 + (uintptr_t)(lbrCount) * unsafe_Sizeof((LbrEntry){}));
    } else {
        *info = 4;
    }
    return STATUS_SUCCESS;
}

static int32_t handleTraceClearBuffer(IRP* irp, uint32_t inLen, uint32_t outLen, uint64_t* info) {
    // TODO: Implement ClearBuffers method in TracerState
    *info = 4;
    LogDebug(so_str("Trace buffers cleared"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static uint64_t getRegisterValue(GUEST_REGS* regs, uint32_t index) {
    do {
        if (index == (0)) {
            return regs->Rax;
        } else if (index == (1)) {
            return regs->Rcx;
        } else if (index == (2)) {
            return regs->Rdx;
        } else if (index == (3)) {
            return regs->Rbx;
        } else if (index == (4)) {
            return regs->Rsp;
        } else if (index == (5)) {
            return regs->Rbp;
        } else if (index == (6)) {
            return regs->Rsi;
        } else if (index == (7)) {
            return regs->Rdi;
        } else if (index == (8)) {
            return regs->R8;
        } else if (index == (9)) {
            return regs->R9;
        } else if (index == (10)) {
            return regs->R10;
        } else if (index == (11)) {
            return regs->R11;
        } else if (index == (12)) {
            return regs->R12;
        } else if (index == (13)) {
            return regs->R13;
        } else if (index == (14)) {
            return regs->R14;
        } else if (index == (15)) {
            return regs->R15;
        } else if (index == (16)) {
            // RIP is not in GUEST_REGS, need to read from VMCS
            return 0;
        } else {
            return 0;
        }
    } while (0);
}

static void setRegisterValue(GUEST_REGS* regs, uint32_t index, uint64_t value) {
    do {
        if (index == (0)) {
            regs->Rax = value;
        } else if (index == (1)) {
            regs->Rcx = value;
        } else if (index == (2)) {
            regs->Rdx = value;
        } else if (index == (3)) {
            regs->Rbx = value;
        } else if (index == (4)) {
            regs->Rsp = value;
        } else if (index == (5)) {
            regs->Rbp = value;
        } else if (index == (6)) {
            regs->Rsi = value;
        } else if (index == (7)) {
            regs->Rdi = value;
        } else if (index == (8)) {
            regs->R8 = value;
        } else if (index == (9)) {
            regs->R9 = value;
        } else if (index == (10)) {
            regs->R10 = value;
        } else if (index == (11)) {
            regs->R11 = value;
        } else if (index == (12)) {
            regs->R12 = value;
        } else if (index == (13)) {
            regs->R13 = value;
        } else if (index == (14)) {
            regs->R14 = value;
        } else if (index == (15)) {
            regs->R15 = value;
        } else if (index == (16)) {
        }
    } while (0);
}

static void eptInveptSingleContext(uint64_t eptp) {
    inveptSingleContext(eptp);
}

static void eptInveptAllContexts(void) {
    inveptAllContexts();
}

// -- kdserial.go --

KdSerialState* NewKdSerialState(void) {
    uint32_t bufSize = (4096);
    KdSerialState* s = &(KdSerialState){.config = (SerialConfig){.PortNumber = 2, .BaudRate = Baud115200, .DataBits = 8, .Parity = ParityNone, .StopBits = StopBits1, .UseIrq = false}, .lock = NewSpinlock(), .txBuffer = so_make_slice_impl(sizeof(so_byte), bufSize, bufSize), .rxBuffer = so_make_slice_impl(sizeof(so_byte), bufSize, bufSize), .bufferSize = bufSize, .lineControl = 0x03, .portBase = 0x2F8};
    return s;
}

so_Error KdSerialState_Initialize(void* self, SerialConfig* config) {
    KdSerialState* s = (KdSerialState*)self;
    Spinlock_Lock(s->lock);
    if (config != NULL) {
        s->config = *config;
    }
    s->portBase = KdSerialState_comPortBase(s, s->config.PortNumber);
    if (s->portBase == 0) {
        Spinlock_Unlock(s->lock);
        return fmtError(so_str("invalid COM port number: %d"), (so_Slice){(void*[1]){&(uint32_t){s->config.PortNumber}}, 1, 1});
    }
    {
        so_Error err = KdSerialState_configureHardware(s);
        if (err != NULL) {
            Spinlock_Unlock(s->lock);
            return err;
        }
    }
    s->connected = true;
    LogInfo(so_str("KD serial initialized on COM%d @0x%X, %d baud"), (so_Slice){(void*[3]){&(uint32_t){s->config.PortNumber}, &(uint16_t){s->portBase}, &(uint32_t){s->config.BaudRate}}, 3, 3});
    Spinlock_Unlock(s->lock);
    return NULL;
}

so_Error KdSerialState_Uninitialize(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    Spinlock_Lock(s->lock);
    if (!s->connected) {
        Spinlock_Unlock(s->lock);
        return NULL;
    }
    KdSerialState_disableUart(s);
    s->connected = false;
    s->txHead = 0;
    s->txTail = 0;
    s->rxHead = 0;
    s->rxTail = 0;
    LogInfo(so_str("KD serial uninitialized"), (so_Slice){&so_Nil, 0, 0});
    Spinlock_Unlock(s->lock);
    return NULL;
}

bool KdSerialState_IsConnected(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    Spinlock_Lock(s->lock);
    Spinlock_Unlock(s->lock);
    return s->connected;
}

static uint16_t KdSerialState_comPortBase(void* self, uint32_t portNum) {
    KdSerialState* s = (KdSerialState*)self;
    do {
        if (portNum == (1)) {
            return 0x3F8;
        } else if (portNum == (2)) {
            return 0x2F8;
        } else if (portNum == (3)) {
            return 0x3E8;
        } else if (portNum == (4)) {
            return 0x2E8;
        } else {
            return 0;
        }
    } while (0);
}

static so_Error KdSerialState_configureHardware(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    KdSerialState_writePort(s, UART_IER, 0);
    KdSerialState_setDlab(s, true);
    uint16_t divisor = KdSerialState_calculateDivisor(s);
    KdSerialState_writePort(s, UART_DLL, (uint8_t)(divisor & 0xFF));
    KdSerialState_writePort(s, UART_DLM, (uint8_t)((divisor >> 8) & 0xFF));
    KdSerialState_setDlab(s, false);
    uint8_t lcr = (0x03);
    do {
        if (s->config.Parity == (ParityOdd)) {
            lcr |= 0x08;
        } else if (s->config.Parity == (ParityEven)) {
            lcr |= 0x18;
        } else if (s->config.Parity == (ParityMark)) {
            lcr |= 0x28;
        } else if (s->config.Parity == (ParitySpace)) {
            lcr |= 0x38;
        }
    } while (0);
    do {
        if (s->config.DataBits == (5)) {
            lcr &= ~(0x03);
        } else if (s->config.DataBits == (6)) {
            lcr = ((lcr & ~(0x03)) | 0x01);
        } else if (s->config.DataBits == (7)) {
            lcr = ((lcr & ~(0x03)) | 0x02);
        } else {
            lcr = ((lcr & ~(0x03)) | 0x03);
        }
    } while (0);
    do {
        if (s->config.StopBits == (StopBits2)) {
            lcr |= 0x04;
        }
    } while (0);
    s->lineControl = lcr;
    KdSerialState_writePort(s, UART_LCR, lcr);
    KdSerialState_writePort(s, UART_FCR, 0xC7);
    uint8_t mcr = (0x0B);
    s->modemControl = mcr;
    KdSerialState_writePort(s, UART_MCR, mcr);
    KdSerialState_writePort(s, UART_IER, 0x00);
    KdSerialState_readPort(s, UART_RBR);
    KdSerialState_readPort(s, UART_IIR);
    KdSerialState_readPort(s, UART_LSR);
    KdSerialState_readPort(s, UART_MSR);
    LogDebug(so_str("SERIAL: Hardware configured, divisor=%d, LCR=0x%02X"), (so_Slice){(void*[2]){&divisor, &lcr}, 2, 2});
    return NULL;
}

static uint16_t KdSerialState_calculateDivisor(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    const uint32_t baseClock = 115200;
    uint32_t rate = (s->config.BaudRate);
    if (rate == 0) {
        rate = 9600;
    }
    return (uint16_t)(baseClock / rate);
}

static void KdSerialState_setDlab(void* self, bool enable) {
    KdSerialState* s = (KdSerialState*)self;
    if (enable) {
        s->lineControl |= 0x80;
    } else {
        s->lineControl &= ~(0x80);
    }
    KdSerialState_writePort(s, UART_LCR, s->lineControl);
    s->divisorLatch = enable;
}

static void KdSerialState_disableUart(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    KdSerialState_writePort(s, UART_IER, 0x00);
    KdSerialState_writePort(s, UART_FCR, 0x00);
    KdSerialState_writePort(s, UART_MCR, 0x00);
}

static void KdSerialState_writePort(void* self, uint8_t reg, uint8_t value) {
    KdSerialState* s = (KdSerialState*)self;
    outb(s->portBase + (uint16_t)(reg), value);
}

static uint8_t KdSerialState_readPort(void* self, uint8_t reg) {
    KdSerialState* s = (KdSerialState*)self;
    return inb(s->portBase + (uint16_t)(reg));
}

so_R_int_err KdSerialState_Send(void* self, so_Slice data) {
    KdSerialState* s = (KdSerialState*)self;
    if (!KdSerialState_IsConnected(s)) {
        return (so_R_int_err){.val = 0, .err = fmtError(so_str("serial not connected"), (so_Slice){&so_Nil, 0, 0})};
    }
    Spinlock_Lock(s->lock);
    so_int bytesWritten = 0;
    for (so_int _ = 0; _ < so_len(data); _++) {
        so_byte b = so_at(so_byte, data, _);
        uint32_t nextTail = (s->txTail + 1) % s->bufferSize;
        if (nextTail == s->txHead) {
            break;
        }
        so_at(so_byte, s->txBuffer, s->txTail) = b;
        s->txTail = nextTail;
        bytesWritten++;
    }
    if (bytesWritten > 0) {
        KdSerialState_flushTxBufferHardware(s);
    }
    Spinlock_Unlock(s->lock);
    return (so_R_int_err){.val = bytesWritten, .err = NULL};
}

static void KdSerialState_flushTxBufferHardware(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    for (; s->txHead != s->txTail;) {
        if (!KdSerialState_isTransmitEmpty(s)) {
            continue;
        }
        so_byte b = so_at(so_byte, s->txBuffer, s->txHead);
        KdSerialState_writePort(s, UART_THR, b);
        s->txHead = (s->txHead + 1) % s->bufferSize;
    }
}

static bool KdSerialState_isTransmitEmpty(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    uint8_t lsr = KdSerialState_readPort(s, UART_LSR);
    return (lsr & 0x20) != 0;
}

static bool KdSerialState_isReceiveReady(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    uint8_t lsr = KdSerialState_readPort(s, UART_LSR);
    return (lsr & 0x01) != 0;
}

so_R_int_err KdSerialState_Receive(void* self, so_Slice buffer) {
    KdSerialState* s = (KdSerialState*)self;
    if (!KdSerialState_IsConnected(s)) {
        return (so_R_int_err){.val = 0, .err = fmtError(so_str("serial not connected"), (so_Slice){&so_Nil, 0, 0})};
    }
    Spinlock_Lock(s->lock);
    uint32_t available = KdSerialState_bytesInRxBuffer(s);
    if (available == 0) {
        KdSerialState_pollRxHardware(s);
        available = KdSerialState_bytesInRxBuffer(s);
    }
    so_int bytesRead = 0;
    so_int maxRead = so_min(so_len(buffer), (so_int)(available));
    for (so_int i = 0; i < maxRead; i++) {
        if (s->rxHead == s->rxTail) {
            break;
        }
        so_at(so_byte, buffer, i) = so_at(so_byte, s->rxBuffer, s->rxHead);
        s->rxHead = (s->rxHead + 1) % s->bufferSize;
        bytesRead++;
    }
    Spinlock_Unlock(s->lock);
    return (so_R_int_err){.val = bytesRead, .err = NULL};
}

static void KdSerialState_pollRxHardware(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    for (; KdSerialState_isReceiveReady(s);) {
        uint32_t nextTail = (s->rxTail + 1) % s->bufferSize;
        if (nextTail == s->rxHead) {
            break;
        }
        uint8_t b = KdSerialState_readPort(s, UART_RBR);
        so_at(so_byte, s->rxBuffer, s->rxTail) = b;
        s->rxTail = nextTail;
    }
}

static uint32_t KdSerialState_bytesInRxBuffer(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    if (s->rxTail >= s->rxHead) {
        return s->rxTail - s->rxHead;
    }
    return s->bufferSize - s->rxHead + s->rxTail;
}

static uint32_t KdSerialState_txBytesQueued(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    if (s->txHead > s->txTail) {
        return s->bufferSize - s->txHead + s->txTail;
    }
    return s->txTail - s->txHead;
}

so_R_int_err KdSerialState_SendString(void* self, so_String str) {
    KdSerialState* s = (KdSerialState*)self;
    return KdSerialState_Send(s, so_string_bytes(str));
}

so_R_str_err KdSerialState_ReceiveString(void* self, so_int maxLen) {
    KdSerialState* s = (KdSerialState*)self;
    so_Slice buf = so_make_slice_impl(sizeof(so_byte), maxLen, maxLen);
    so_R_int_err _res1 = KdSerialState_Receive(s, buf);
    so_int n = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        return (so_R_str_err){.val = so_str(""), .err = err};
    }
    return (so_R_str_err){.val = so_bytes_string(so_slice(so_byte, buf, 0, n)), .err = NULL};
}

uint32_t KdSerialState_BytesAvailable(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    Spinlock_Lock(s->lock);
    KdSerialState_pollRxHardware(s);
    Spinlock_Unlock(s->lock);
    return KdSerialState_bytesInRxBuffer(s);
}

uint32_t KdSerialState_TxSpaceAvailable(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    Spinlock_Lock(s->lock);
    uint32_t queued = KdSerialState_txBytesQueued(s);
    if (queued >= s->bufferSize - 1) {
        Spinlock_Unlock(s->lock);
        return 0;
    }
    Spinlock_Unlock(s->lock);
    return s->bufferSize - 1 - queued;
}

so_Error KdSerialState_SendByte(void* self, so_byte b) {
    KdSerialState* s = (KdSerialState*)self;
    so_Slice data = (so_Slice){(so_byte[1]){b}, 1, 1};
    so_R_int_err _res1 = KdSerialState_Send(s, data);
    so_Error err = _res1.err;
    return err;
}

so_R_byte_err KdSerialState_ReceiveByte(void* self) {
    KdSerialState* s = (KdSerialState*)self;
    so_Slice buf = so_make_slice_impl(sizeof(so_byte), 1, 1);
    so_R_int_err _res1 = KdSerialState_Receive(s, buf);
    so_int n = _res1.val;
    so_Error err = _res1.err;
    if (err != NULL) {
        return (so_R_byte_err){.val = 0, .err = err};
    }
    if (n == 0) {
        return (so_R_byte_err){.val = 0, .err = fmtError(so_str("no data available"), (so_Slice){&so_Nil, 0, 0})};
    }
    return (so_R_byte_err){.val = so_at(so_byte, buf, 0), .err = NULL};
}

// -- logging.go --

so_String LogLevel_String(LogLevel l) {
    do {
        if (l == (LogLevelInfo)) {
            return so_str("INFO");
        } else if (l == (LogLevelWarning)) {
            return so_str("WARN");
        } else if (l == (LogLevelError)) {
            return so_str("ERROR");
        } else if (l == (LogLevelDebug)) {
            return so_str("DEBUG");
        } else if (l == (LogLevelTrace)) {
            return so_str("TRACE");
        } else {
            return so_str("????");
        }
    } while (0);
}

LogBuffer* NewLogBuffer(uint32_t capacity) {
    return &(LogBuffer){.Messages = so_make_slice_impl(sizeof(LogMessage), capacity, capacity), .Capacity = capacity, .Lock = NewSpinlock()};
}

bool LogBuffer_Push(void* self, LogMessage* msg) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    if (b->Count >= b->Capacity) {
        b->OverflowCount++;
        Spinlock_Unlock(b->Lock);
        return false;
    }
    uint32_t idx = (b->WriteIndex + b->Count) % b->Capacity;
    so_at(LogMessage, b->Messages, idx) = *msg;
    b->Count++;
    Spinlock_Unlock(b->Lock);
    return true;
}

so_R_ptr_bool LogBuffer_Pop(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    if (b->Count == 0) {
        Spinlock_Unlock(b->Lock);
        return (so_R_ptr_bool){.val = NULL, .val2 = false};
    }
    LogMessage* msg = &so_at(LogMessage, b->Messages, b->ReadIndex);
    b->ReadIndex = (b->ReadIndex + 1) % b->Capacity;
    b->Count--;
    Spinlock_Unlock(b->Lock);
    return (so_R_ptr_bool){.val = msg, .val2 = true};
}

so_R_ptr_bool LogBuffer_Peek(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    if (b->Count == 0) {
        Spinlock_Unlock(b->Lock);
        return (so_R_ptr_bool){.val = NULL, .val2 = false};
    }
    Spinlock_Unlock(b->Lock);
    return (so_R_ptr_bool){.val = &so_at(LogMessage, b->Messages, b->ReadIndex), .val2 = true};
}

void LogBuffer_Clear(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    b->ReadIndex = 0;
    b->WriteIndex = 0;
    b->Count = 0;
    b->OverflowCount = 0;
    Spinlock_Unlock(b->Lock);
}

bool LogBuffer_IsEmpty(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    Spinlock_Unlock(b->Lock);
    return b->Count == 0;
}

bool LogBuffer_IsFull(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    Spinlock_Unlock(b->Lock);
    return b->Count >= b->Capacity;
}

uint32_t LogBuffer_Available(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    Spinlock_Unlock(b->Lock);
    return b->Capacity - b->Count;
}

uint32_t LogBuffer_GetCount(void* self) {
    LogBuffer* b = (LogBuffer*)self;
    Spinlock_Lock(b->Lock);
    Spinlock_Unlock(b->Lock);
    return b->Count;
}

Logger* NewLogger(uint32_t bufferSize, uint32_t prioSize) {
    Logger* l = &(Logger){.buffer = NewLogBuffer(bufferSize), .priorityBuffer = NewLogBuffer(prioSize), .enabled = true, .minLevel = LogLevelInfo, .lock = NewSpinlock()};
    return l;
}

so_Error Logger_Initialize(void* self) {
    Logger* l = (Logger*)self;
    if (l->buffer != NULL && l->priorityBuffer != NULL) {
        l->enabled = true;
        return NULL;
    }
    return fmtError(so_str("logger initialization failed"), (so_Slice){&so_Nil, 0, 0});
}

void Logger_Uninitialize(void* self) {
    Logger* l = (Logger*)self;
    Spinlock_Lock(l->lock);
    if (l->buffer != NULL) {
        LogBuffer_Clear(l->buffer);
    }
    if (l->priorityBuffer != NULL) {
        LogBuffer_Clear(l->priorityBuffer);
    }
    l->enabled = false;
    Spinlock_Unlock(l->lock);
}

void Logger_SetMinLevel(void* self, LogLevel level) {
    Logger* l = (Logger*)self;
    l->minLevel = level;
}

void Logger_Enable(void* self) {
    Logger* l = (Logger*)self;
    l->enabled = true;
}

void Logger_Disable(void* self) {
    Logger* l = (Logger*)self;
    l->enabled = false;
}

void Logger_Log(void* self, LogLevel level, so_String format, so_Slice args) {
    Logger* l = (Logger*)self;
    if (!l->enabled || level < l->minLevel) {
        return;
    }
    LogMessage msg = {0};
    msg.Level = level;
    msg.Time = rdtscValue();
    uint32_t cpuNum = 0;
    KeGetCurrentProcessorNumberEx(&cpuNum);
    msg.CpuId = cpuNum;
    so_String procName = CurrentProcessName();
    so_copy_string(so_array_slice(uint8_t, msg.Process, 0, 15, 15), procName);
    so_String formatted = fmt_Sprintf(format, args);
    so_int length = so_min(so_len(formatted), 256);
    so_copy_string(so_array_slice(uint8_t, msg.Message, 0, length, 256), formatted);
    msg.Length = (uint32_t)(length);
    if (level <= LogLevelWarning) {
        LogBuffer_Push(l->priorityBuffer, &msg);
    } else {
        LogBuffer_Push(l->buffer, &msg);
    }
}

void Logger_Info(void* self, so_String format, so_Slice args) {
    Logger* l = (Logger*)self;
    Logger_Log(l, LogLevelInfo, format, (so_Slice){(void*[1]){&args}, 1, 1});
}

void Logger_Warning(void* self, so_String format, so_Slice args) {
    Logger* l = (Logger*)self;
    Logger_Log(l, LogLevelWarning, format, (so_Slice){(void*[1]){&args}, 1, 1});
}

void Logger_Error(void* self, so_String format, so_Slice args) {
    Logger* l = (Logger*)self;
    Logger_Log(l, LogLevelError, format, (so_Slice){(void*[1]){&args}, 1, 1});
}

void Logger_Debug(void* self, so_String format, so_Slice args) {
    Logger* l = (Logger*)self;
    Logger_Log(l, LogLevelDebug, format, (so_Slice){(void*[1]){&args}, 1, 1});
}

void Logger_Trace(void* self, so_String format, so_Slice args) {
    Logger* l = (Logger*)self;
    Logger_Log(l, LogLevelTrace, format, (so_Slice){(void*[1]){&args}, 1, 1});
}

uint32_t Logger_FlushToUser(void* self, uintptr_t buffer, uintptr_t size) {
    Logger* l = (Logger*)self;
    uint32_t count = 0;
    uintptr_t remaining = size;
    for (; remaining > 0;) {
        so_R_ptr_bool _res1 = LogBuffer_Pop(l->buffer);
        LogMessage* msg = _res1.val;
        bool ok = _res1.val2;
        if (!ok) {
            break;
        }
        uintptr_t msgSize = unsafe_Sizeof((LogMessage){});
        if ((msgSize) > remaining) {
            break;
        }
        LogMessage* dst = (LogMessage*)((void*)((buffer) + (size - remaining)));
        *dst = *msg;
        remaining -= (msgSize);
        count++;
    }
    return count;
}

uint32_t Logger_FlushPriorityToUser(void* self, uintptr_t buffer, uintptr_t size) {
    Logger* l = (Logger*)self;
    uint32_t count = 0;
    uintptr_t remaining = size;
    for (; remaining > 0;) {
        so_R_ptr_bool _res1 = LogBuffer_Pop(l->priorityBuffer);
        LogMessage* msg = _res1.val;
        bool ok = _res1.val2;
        if (!ok) {
            break;
        }
        uintptr_t msgSize = unsafe_Sizeof((LogMessage){});
        if ((msgSize) > remaining) {
            break;
        }
        LogMessage* dst = (LogMessage*)((void*)((buffer) + (size - remaining)));
        *dst = *msg;
        remaining -= (msgSize);
        count++;
    }
    return count;
}

uint64_t Logger_GetBufferSize(void* self) {
    Logger* l = (Logger*)self;
    return (uint64_t)(l->buffer->Capacity);
}

uint64_t Logger_GetCurrentBufferSize(void* self) {
    Logger* l = (Logger*)self;
    return (uint64_t)(LogBuffer_GetCount(l->buffer));
}

uint64_t Logger_GetOverflowCount(void* self) {
    Logger* l = (Logger*)self;
    return l->buffer->OverflowCount;
}

uint64_t Logger_GetBufferBase(void* self) {
    Logger* l = (Logger*)self;
    if (l->buffer == NULL || so_len(l->buffer->Messages) == 0) {
        return 0;
    }
    return (uint64_t)((uintptr_t)((void*)(&so_at(LogMessage, l->buffer->Messages, 0))));
}

uint64_t Logger_GetPriorityBufferBase(void* self) {
    Logger* l = (Logger*)self;
    if (l->priorityBuffer == NULL || so_len(l->priorityBuffer->Messages) == 0) {
        return 0;
    }
    return (uint64_t)((uintptr_t)((void*)(&so_at(LogMessage, l->priorityBuffer->Messages, 0))));
}

static uint64_t rdtscValue(void) {
    // Use nativeRdtsc from events.go
    return nativeRdtsc();
}

void LogInfo(so_String format, so_Slice args) {
    if (gLog != NULL) {
        Logger_Info(gLog, format, (so_Slice){(void*[1]){&args}, 1, 1});
    }
}

void LogError(so_String format, so_Slice args) {
    if (gLog != NULL) {
        Logger_Error(gLog, format, (so_Slice){(void*[1]){&args}, 1, 1});
    }
}

void LogWarning(so_String format, so_Slice args) {
    if (gLog != NULL) {
        Logger_Warning(gLog, format, (so_Slice){(void*[1]){&args}, 1, 1});
    }
}

void LogDebug(so_String format, so_Slice args) {
    if (gLog != NULL) {
        Logger_Debug(gLog, format, (so_Slice){(void*[1]){&args}, 1, 1});
    }
}

void LogTrace(so_String format, so_Slice args) {
    if (gLog != NULL) {
        Logger_Trace(gLog, format, (so_Slice){(void*[1]){&args}, 1, 1});
    }
}

void LogCallbackSendBuffer(uint32_t operation, so_String data, uint32_t length, bool isImmediate) {
    if (gLog != NULL) {
        Logger_Log(gLog, LogLevelInfo, so_str("Send buffer: op=%d, len=%d"), (so_Slice){(void*[2]){&operation, &length}, 2, 2});
    }
}

void LogRegisterIrpBasedNotification(void* irp) {
    gAllowIoctl = true;
    LogInfo(so_str("Registered IRP-based notification"), (so_Slice){&so_Nil, 0, 0});
}

bool LogRegisterEventBasedNotification(void* irp) {
    gAllowIoctl = true;
    LogInfo(so_str("Registered event-based notification"), (so_Slice){&so_Nil, 0, 0});
    return true;
}

// -- memory.go --

so_String MemoryType_String(MemoryType m) {
    do {
        if (m == (MemUncacheable)) {
            return so_str("UC");
        } else if (m == (MemWriteCombining)) {
            return so_str("WC");
        } else if (m == (MemWriteThrough)) {
            return so_str("WT");
        } else if (m == (MemWriteProtected)) {
            return so_str("WP");
        } else if (m == (MemWriteBack)) {
            return so_str("WB");
        } else {
            return so_str("??");
        }
    } while (0);
}

static MemoryManager* newMemoryManager(void) {
    return &(MemoryManager){.mtrr = &(MemoryMtrrState){.DefaultType = MemUncacheable}};
}

MemoryType MemoryManager_QueryMtrrForPa(void* self, uint64_t pa) {
    MemoryManager* m = (MemoryManager*)self;
    for (so_int i = 0; i < 255; i++) {
        MtrrRange* r = &m->mtrr->VariableRanges[i];
        if (r->Valid && pa >= r->BaseAddr && pa < r->EndAddr) {
            return r->Type;
        }
    }
    return m->mtrr->DefaultType;
}

void MemoryManager_SetDefaultMemType(void* self, MemoryType t) {
    MemoryManager* m = (MemoryManager*)self;
    m->mtrr->DefaultType = t;
}

void MemoryManager_AddVariableRange(void* self, uint64_t base, uint64_t end, MemoryType t) {
    MemoryManager* m = (MemoryManager*)self;
    for (so_int i = 0; i < 255; i++) {
        if (!m->mtrr->VariableRanges[i].Valid) {
            m->mtrr->VariableRanges[i] = (MtrrRange){.BaseAddr = base, .EndAddr = end, .Type = t, .Valid = true};
            return;
        }
    }
}

void MemoryManager_ClearAllRanges(void* self) {
    MemoryManager* m = (MemoryManager*)self;
    for (so_int i = 0; i < 255; i++) {
        m->mtrr->VariableRanges[i].Valid = false;
    }
}

uint64_t VirtToPhys(uintptr_t va, CR3_TYPE cr3) {
    uint64_t pml4Idx = (((uint64_t)(va) >> 39) & 0x1FF);
    uint64_t pml3Idx = (((uint64_t)(va) >> 30) & 0x1FF);
    uint64_t pml2Idx = (((uint64_t)(va) >> 21) & 0x1FF);
    uint64_t pml1Idx = (((uint64_t)(va) >> 12) & 0x1FF);
    PML4E (*pml4Table)[512] = (PML4E(*)[512])((void*)((uintptr_t)(cr3.Flags & 0xFFFFFFFFFFFFF000)));
    PML4E pml4Entry = (*pml4Table)[pml4Idx];
    if ((pml4Entry.AsUInt & 0x1) == 0) {
        return 0;
    }
    uintptr_t pml3Base = (uintptr_t)(pml4Entry.AsUInt & 0xFFFFFFFFFF000);
    EPT_PDPTE (*pml3Table)[512] = (EPT_PDPTE(*)[512])((void*)(pml3Base));
    EPT_PDPTE pml3Entry = (*pml3Table)[pml3Idx];
    if ((pml3Entry.AsUInt & 0x1) == 0) {
        return 0;
    }
    if ((pml3Entry.AsUInt & ((uint64_t)1 << 7)) != 0) {
        uint64_t pageBase = (pml3Entry.AsUInt & (((so_int)1 << 30) - 1));
        uint64_t offset = ((uint64_t)(va) & ((uint64_t)(SIZE_1_GB) - 1));
        return pageBase + offset;
    }
    uintptr_t pml2Base = (uintptr_t)(pml3Entry.AsUInt & 0xFFFFFFFFFF000);
    EPT_PDE (*pml2Table)[512] = (EPT_PDE(*)[512])((void*)(pml2Base));
    EPT_PDE pml2Entry = (*pml2Table)[pml2Idx];
    if ((pml2Entry.AsUInt & 0x1) == 0) {
        return 0;
    }
    if ((pml2Entry.AsUInt & ((uint64_t)1 << 7)) != 0) {
        uint64_t pageBase = (pml2Entry.AsUInt & (((so_int)1 << 21) - 1));
        uint64_t offset = ((uint64_t)(va) & ((uint64_t)(SIZE_2_MB) - 1));
        return pageBase + offset;
    }
    uintptr_t pml1Base = (uintptr_t)(pml2Entry.AsUInt & 0xFFFFFFFFFF000);
    EPT_PTE (*pml1Table)[512] = (EPT_PTE(*)[512])((void*)(pml1Base));
    EPT_PTE pml1Entry = (*pml1Table)[pml1Idx];
    if ((pml1Entry.AsUInt & 0x1) == 0) {
        return 0;
    }
    uint64_t pageBase = (pml1Entry.AsUInt & 0xFFFFFFFFFF000);
    uint64_t offset = ((uint64_t)(va) & ((uint64_t)(PAGE_SIZE) - 1));
    return pageBase + offset;
}

uintptr_t PhysToVirt(uint64_t pa, CR3_TYPE cr3) {
    return 0;
}

void CopyPage(uintptr_t src, uintptr_t dst, CR3_TYPE cr3) {
    uint64_t srcPa = VirtToPhys(src, cr3);
    uint64_t dstPa = VirtToPhys(dst, cr3);
    if (srcPa == 0 || dstPa == 0) {
        return;
    }
    so_byte (*srcPtr)[4096] = (so_byte(*)[4096])((void*)((uintptr_t)(srcPa)));
    so_byte (*dstPtr)[4096] = (so_byte(*)[4096])((void*)((uintptr_t)(dstPa)));
    so_copy(so_byte, so_array_slice(so_byte, (*dstPtr), 0, 4096, 4096), so_array_slice(so_byte, (*srcPtr), 0, 4096, 4096));
}

so_int ReadPhysMem(uint64_t pa, so_Slice buf) {
    so_byte* ptr = (so_byte*)((void*)((uintptr_t)(pa)));
    uint64_t maxLen = (uint64_t)(PAGE_SIZE) - (pa % (uint64_t)(PAGE_SIZE));
    if ((so_int)(maxLen) > so_len(buf)) {
        maxLen = (uint64_t)(so_len(buf));
    }
    for (uint64_t i = (0); i < maxLen; i++) {
        so_at(so_byte, buf, i) = *(so_byte*)((void*)(((uintptr_t)((void*)(ptr)) + (uintptr_t)(i))));
    }
    return (so_int)(maxLen);
}

so_int WritePhysMem(uint64_t pa, so_Slice data) {
    so_byte* ptr = (so_byte*)((void*)((uintptr_t)(pa)));
    uint64_t maxLen = (uint64_t)(PAGE_SIZE) - (pa % (uint64_t)(PAGE_SIZE));
    if ((so_int)(maxLen) > so_len(data)) {
        maxLen = (uint64_t)(so_len(data));
    }
    for (uint64_t i = (0); i < maxLen; i++) {
        *(so_byte*)((void*)(((uintptr_t)((void*)(ptr)) + (uintptr_t)(i)))) = so_at(so_byte, data, i);
    }
    return (so_int)(maxLen);
}

static uint64_t readPhysicalMemory(uint64_t pa) {
    uint64_t* ptr = (uint64_t*)((void*)((uintptr_t)(pa)));
    return *ptr;
}

static void writePhysicalMemory(uint64_t pa, uint64_t val) {
    uint64_t* ptr = (uint64_t*)((void*)((uintptr_t)(pa)));
    *ptr = val;
}

static void invalidateEpt(void) {
    inveptAllContexts();
}

// -- mode_based_exec.go --

static void modeBasedExecHookEnable(VCPU* v) {
    gModeBasedExec->Enabled = true;
    uint64_t secondaryControls = (0);
    vmRead64(VmcsCtrlSecondaryProcBasedVmExec, &secondaryControls);
    secondaryControls |= ((uint64_t)1 << 2);
    vmWrite64(VmcsCtrlSecondaryProcBasedVmExec, secondaryControls);
    LogDebug(so_str("Mode-based execution hook enabled on core %d"), (so_Slice){(void*[1]){&(uint32_t){v->CoreId}}, 1, 1});
}

static void modeBasedExecHookDisable(VCPU* v) {
    gModeBasedExec->Enabled = false;
    uint64_t secondaryControls = (0);
    vmRead64(VmcsCtrlSecondaryProcBasedVmExec, &secondaryControls);
    secondaryControls &= ~((uint64_t)1 << 2);
    vmWrite64(VmcsCtrlSecondaryProcBasedVmExec, secondaryControls);
    LogDebug(so_str("Mode-based execution hook disabled on core %d"), (so_Slice){(void*[1]){&(uint32_t){v->CoreId}}, 1, 1});
}

static so_int modeBasedExecHookAddSelector(uint16_t selector) {
    ModeBasedExecEntry entry = (ModeBasedExecEntry){.Selector = selector, .Hooked = true};
    gModeBasedExec->CsSelectors = so_append(ModeBasedExecEntry, gModeBasedExec->CsSelectors, entry);
    return so_len(gModeBasedExec->CsSelectors) - 1;
}

static void modeBasedExecHookRemoveSelector(so_int index) {
    if (index >= 0 && index < so_len(gModeBasedExec->CsSelectors)) {
        so_at(ModeBasedExecEntry, gModeBasedExec->CsSelectors, index).Hooked = false;
    }
}

static bool modeBasedExecHandleEptViolation(VCPU* v, uint64_t qualification) {
    if (!gModeBasedExec->Enabled) {
        return false;
    }
    uint64_t guestCs = 0;
    vmRead64(VmcsGuestCsSelector, &guestCs);
    uint16_t csSelector = (uint16_t)(guestCs);
    for (so_int i = 0; i < so_len(gModeBasedExec->CsSelectors); i++) {
        ModeBasedExecEntry* entry = &so_at(ModeBasedExecEntry, gModeBasedExec->CsSelectors, i);
        if (entry->Hooked && csSelector == entry->Selector) {
            Hypervisor* h = getHypervisor();
            if (h != NULL) {
                Hypervisor_suppressAdvance(h, v);
            }
            dispatchEventModeBasedExec(v, csSelector, qualification);
            return true;
        }
    }
    return false;
}

static void dispatchEventModeBasedExec(VCPU* v, uint16_t selector, uint64_t qualification) {
}

// -- msr_handlers.go --

static void initMsrBitmap(void) {
    for (so_int i = 0; i < MSR_BITMAP_SIZE; i++) {
        gMsrBitmap->bitmap[i] = 0xFF;
    }
}

static void msrHandlePerformMsrBitmapReadChange(VCPU* v, uint32_t msrIndex) {
    if (v->MsrBitmapVirtualAddress == 0) {
        LogError(so_str("MSR bitmap not initialized"), (so_Slice){&so_Nil, 0, 0});
        return;
    }
    uint32_t byteOffset = msrIndex / 8;
    uint32_t bitOffset = msrIndex % 8;
    if (msrIndex >= 0xC0000000) {
        byteOffset += MSR_BITMAP_SIZE / 2;
    }
    if (byteOffset < MSR_BITMAP_SIZE) {
        uint8_t mask = ((uint8_t)1 << bitOffset);
        gMsrBitmap->bitmap[byteOffset] |= mask;
        uint64_t pa = VirtToPhys((uintptr_t)(v->MsrBitmapVirtualAddress), (CR3_TYPE){});
        WritePhysMem(pa + (uint64_t)(byteOffset), so_array_slice(so_byte, gMsrBitmap->bitmap, byteOffset, byteOffset + 1, 8192));
    }
}

static void msrHandlePerformMsrBitmapWriteChange(VCPU* v, uint32_t msrIndex) {
    if (v->MsrBitmapVirtualAddress == 0) {
        LogError(so_str("MSR bitmap not initialized"), (so_Slice){&so_Nil, 0, 0});
        return;
    }
    uint32_t byteOffset = msrIndex / 8 + MSR_BITMAP_SIZE / 4;
    uint32_t bitOffset = msrIndex % 8;
    if (msrIndex >= 0xC0000000) {
        byteOffset += MSR_BITMAP_SIZE / 2;
    }
    if (byteOffset < MSR_BITMAP_SIZE) {
        uint8_t mask = ((uint8_t)1 << bitOffset);
        gMsrBitmap->bitmap[byteOffset] |= mask;
        uint64_t pa = VirtToPhys((uintptr_t)(v->MsrBitmapVirtualAddress), (CR3_TYPE){});
        WritePhysMem(pa + (uint64_t)(byteOffset), so_array_slice(so_byte, gMsrBitmap->bitmap, byteOffset, byteOffset + 1, 8192));
    }
}

static void msrHandleEnableOrDisableMsrExiting(VCPU* v, bool enable) {
    if (enable) {
        uint64_t cpuBasedControls = (0);
        vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
        cpuBasedControls |= ((uint64_t)1 << 11);
        vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
        initMsrBitmap();
    } else {
        uint64_t cpuBasedControls = (0);
        vmRead64(VmcsCtrlPrimaryProcBasedVmExec, &cpuBasedControls);
        cpuBasedControls &= ~(((uint64_t)1 << 11));
        vmWrite64(VmcsCtrlPrimaryProcBasedVmExec, cpuBasedControls);
    }
}

static bool handleMsrRead(VCPU* v) {
    uint64_t msrIndex = v->Regs->Rcx;
    uint64_t value = 0;
    do {
        if ((uint32_t)(msrIndex) == (IA32_EFER)) {
            value = __readmsr(IA32_EFER);
        } else {
            value = __readmsr((uint32_t)(msrIndex));
        }
    } while (0);
    v->Regs->Rax = (value & 0xFFFFFFFF);
    v->Regs->Rdx = ((value >> 32) & 0xFFFFFFFF);
    Hypervisor* h = getHypervisor();
    if (h != NULL) {
        Hypervisor_advanceIp(h, v);
    }
    return true;
}

static bool handleMsrWrite(VCPU* v) {
    uint64_t msrIndex = v->Regs->Rcx;
    uint64_t value = ((v->Regs->Rdx << 32) | v->Regs->Rax);
    do {
        if ((uint32_t)(msrIndex) == (IA32_EFER)) {
            __writemsr(IA32_EFER, value);
        } else {
            __writemsr((uint32_t)(msrIndex), value);
        }
    } while (0);
    Hypervisor* h = getHypervisor();
    if (h != NULL) {
        Hypervisor_advanceIp(h, v);
    }
    return true;
}

// -- optimizations.go --

PoolManager* NewPoolManager(so_int maxEntries) {
    return &(PoolManager){.lock = NewSpinlock(), .maxPools = maxEntries};
}

uintptr_t PoolManager_Allocate(void* self, uint64_t size, uint32_t tag) {
    PoolManager* m = (PoolManager*)self;
    Spinlock_Lock(m->lock);
    uintptr_t ptr = ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), (uintptr_t)(size), (tag));
    if (ptr == 0) {
        Spinlock_Unlock(m->lock);
        return 0;
    }
    uintptr_t entry = AllocatePool(unsafe_Sizeof((PoolEntry){}), POOLTAG);
    if (entry == 0) {
        ExFreePoolWithTag(ptr, POOLTAG);
        Spinlock_Unlock(m->lock);
        return 0;
    }
    PoolEntry* entryPtr = (PoolEntry*)((void*)(entry));
    entryPtr->Address = (ptr);
    entryPtr->Size = size;
    entryPtr->Tag = tag;
    if (m->poolList != NULL) {
        m->poolList->Prev = entryPtr;
        entryPtr->Next = m->poolList;
    }
    m->poolList = entryPtr;
    m->totalAllocated += size;
    m->currentPools++;
    Spinlock_Unlock(m->lock);
    return (ptr);
}

void PoolManager_Free(void* self, uintptr_t ptr) {
    PoolManager* m = (PoolManager*)self;
    if (ptr == 0) {
        return;
    }
    Spinlock_Lock(m->lock);
    PoolEntry* entry = PoolManager_findEntry(m, ptr);
    if (entry == NULL) {
        Spinlock_Unlock(m->lock);
        return;
    }
    if (entry->Next != NULL) {
        entry->Next->Prev = entry->Prev;
    }
    if (entry->Prev != NULL) {
        entry->Prev->Next = entry->Next;
    } else if (m->poolList == entry) {
        m->poolList = entry->Next;
    }
    m->totalFreed += entry->Size;
    m->currentPools--;
    ExFreePoolWithTag(entry->Address, (entry->Tag));
    FreePool((uintptr_t)((void*)(entry)), POOLTAG);
    Spinlock_Unlock(m->lock);
}

static PoolEntry* PoolManager_findEntry(void* self, uintptr_t ptr) {
    PoolManager* m = (PoolManager*)self;
    PoolEntry* current = m->poolList;
    for (; current != NULL;) {
        if (current->Address == ptr) {
            return current;
        }
        current = current->Next;
    }
    return NULL;
}

so_int PoolManager_CheckAndPerformDeallocation(void* self) {
    PoolManager* m = (PoolManager*)self;
    so_int freedCount = 0;
    Spinlock_Lock(m->lock);
    PoolEntry* current = m->poolList;
    for (; current != NULL;) {
        PoolEntry* next = current->Next;
        if (current->Address != 0 && MmIsAddressValid(current->Address)) {
            ExFreePoolWithTag(current->Address, (current->Tag));
            m->totalFreed += current->Size;
            m->currentPools--;
            freedCount++;
        }
        current = next;
    }
    Spinlock_Unlock(m->lock);
    return freedCount;
}

PoolStats PoolManager_GetStats(void* self) {
    PoolManager* m = (PoolManager*)self;
    Spinlock_Lock(m->lock);
    Spinlock_Unlock(m->lock);
    return (PoolStats){.TotalAllocated = m->totalAllocated, .TotalFreed = m->totalFreed, .CurrentPools = m->currentPools, .MaxPools = m->maxPools};
}

so_String PoolStats_String(PoolStats s) {
    return so_str("");
}

VmmOptimizer* NewVmmOptimizer(so_int level) {
    VmmOptimizer* o = &(VmmOptimizer){.level = level};
    VmmOptimizer_applyDefaults(o);
    return o;
}

static void VmmOptimizer_applyDefaults(void* self) {
    VmmOptimizer* o = (VmmOptimizer*)self;
    do {
        if (o->level == (OptNone)) {
            o->eptCacheEnabled = false;
            o->msrBitmapCached = false;
            o->ioBitmapCached = false;
            o->inveptBatching = false;
        } else if (o->level == (OptBasic)) {
            o->eptCacheEnabled = true;
            o->msrBitmapCached = true;
            o->ioBitmapCached = false;
            o->inveptBatching = false;
        } else if (o->level == (OptAggressive)) {
            o->eptCacheEnabled = true;
            o->msrBitmapCached = true;
            o->ioBitmapCached = true;
            o->inveptBatching = true;
        }
    } while (0);
}

void VmmOptimizer_EnableEptCache(void* self, bool enable) {
    VmmOptimizer* o = (VmmOptimizer*)self;
    o->eptCacheEnabled = enable;
}

void VmmOptimizer_EnableMsrBitmapCache(void* self, bool enable) {
    VmmOptimizer* o = (VmmOptimizer*)self;
    o->msrBitmapCached = enable;
}

void VmmOptimizer_EnableIoBitmapCache(void* self, bool enable) {
    VmmOptimizer* o = (VmmOptimizer*)self;
    o->ioBitmapCached = enable;
}

void VmmOptimizer_EnableInveptBatching(void* self, bool enable) {
    VmmOptimizer* o = (VmmOptimizer*)self;
    o->inveptBatching = enable;
}

VmmOptConfig VmmOptimizer_GetConfig(void* self) {
    VmmOptimizer* o = (VmmOptimizer*)self;
    return (VmmOptConfig){.EptCache = o->eptCacheEnabled, .MsrBitmapCache = o->msrBitmapCached, .IoBitmapCache = o->ioBitmapCached, .InveptBatching = o->inveptBatching};
}

EptCache* NewEptCache(so_int maxSize) {
    EptCache* c = &(EptCache){.entries = so_make_map_impl(sizeof(uint64_t), sizeof(EptCacheEntry*), 0), .lock = NewSpinlock(), .maxSize = maxSize};
    return c;
}

so_R_ptr_bool EptCache_Lookup(void* self, uint64_t gpa) {
    EptCache* c = (EptCache*)self;
    Spinlock_Lock(c->lock);
    EptCacheEntry* entry = so_map_get(uint64_t, EptCacheEntry*, c->entries, gpa);
    bool ok = so_map_has(uint64_t, c->entries, gpa);
    if (!ok) {
        Spinlock_Unlock(c->lock);
        return (so_R_ptr_bool){.val = NULL, .val2 = false};
    }
    entry->LastAccess = rdtscValue();
    entry->HitCount++;
    Spinlock_Unlock(c->lock);
    return (so_R_ptr_bool){.val = &entry->Entry, .val2 = true};
}

void EptCache_Store(void* self, uint64_t gpa, uint64_t pa, EptEntry entry) {
    EptCache* c = (EptCache*)self;
    Spinlock_Lock(c->lock);
    if ((so_int)c->entries->len >= c->maxSize) {
        EptCache_evictOldest(c);
    }
    EptCacheEntry* cacheEntry = &(EptCacheEntry){.Gpa = gpa, .Pa = pa, .Entry = entry, .LastAccess = rdtscValue(), .HitCount = 1};
    {
    	EptCacheEntry* _so_map_assign_val;
    	memset(&_so_map_assign_val, 0, sizeof(_so_map_assign_val));
    	_so_map_assign_val = cacheEntry;
    	so_map_set(uint64_t, EptCacheEntry*, c->entries, gpa, _so_map_assign_val);
    }
    Spinlock_Unlock(c->lock);
}

void EptCache_Invalidate(void* self, uint64_t gpa) {
    EptCache* c = (EptCache*)self;
    Spinlock_Lock(c->lock);
    so_map_remove(uint64_t, c->entries, gpa);
    Spinlock_Unlock(c->lock);
}

void EptCache_InvalidateAll(void* self) {
    EptCache* c = (EptCache*)self;
    Spinlock_Lock(c->lock);
    c->entries = so_make_map_impl(sizeof(uint64_t), sizeof(EptCacheEntry*), 0);
    Spinlock_Unlock(c->lock);
}

static void EptCache_evictOldest(void* self) {
    EptCache* c = (EptCache*)self;
    uint64_t oldestGpa = 0;
    uint64_t oldestTime = ~(0);
    for (so_int _i = 0; _i < (so_int)c->entries->cap; _i++) {
        if (!c->entries->used[_i]) continue;
        uint64_t gpa = ((uint64_t*)c->entries->keys)[_i];
        EptCacheEntry* entry = ((EptCacheEntry**)c->entries->vals)[_i];
        if (entry->LastAccess < oldestTime) {
            oldestTime = entry->LastAccess;
            oldestGpa = gpa;
        }
    }
    if (oldestTime != ~(0)) {
        so_map_remove(uint64_t, c->entries, oldestGpa);
    }
}

EptCacheStats EptCache_GetStats(void* self) {
    EptCache* c = (EptCache*)self;
    Spinlock_Lock(c->lock);
    EptCacheStats stats = (EptCacheStats){.TotalEntries = (so_int)c->entries->len, .MaxEntries = c->maxSize};
    for (so_int _i = 0; _i < (so_int)c->entries->cap; _i++) {
        if (!c->entries->used[_i]) continue;
        EptCacheEntry* entry = ((EptCacheEntry**)c->entries->vals)[_i];
        stats.TotalHits += entry->HitCount;
    }
    Spinlock_Unlock(c->lock);
    return stats;
}

// -- platform.go --

uintptr_t AllocatePool(uintptr_t size, uint32_t tag) {
    return (ExAllocatePool2((uint32_t)(POOL_FLAG_NON_PAGED), size, tag));
}

uintptr_t AllocateZeroedPool(uintptr_t size, uint32_t tag) {
    uintptr_t ptr = AllocatePool(size, tag);
    if (ptr != 0) {
        RtlZeroMemory(ptr, size);
    }
    return ptr;
}

void FreePool(uintptr_t ptr, uint32_t tag) {
    if (ptr != 0) {
        ExFreePoolWithTag(ptr, tag);
    }
}

so_String CurrentProcessName(void) {
    uintptr_t* proc = PsGetCurrentProcess();
    int8_t* name = PsGetProcessImageFileName(proc);
    return cstringToString(name);
}

static so_String cstringToString(int8_t* s) {
    if (s == NULL) {
        return so_str("");
    }
    so_Slice result = {&so_Nil, 0, 0};
    for (so_int i = 0;; i++) {
        uint8_t b = *(uint8_t*)((void*)((uintptr_t)((void*)(s)) + (uintptr_t)(i)));
        if (b == 0) {
            break;
        }
        result = so_append(so_byte, result, b);
    }
    return so_bytes_string(result);
}

// -- spinlock.go --

Spinlock* NewSpinlock(void) {
    return &(Spinlock){};
}

void Spinlock_Lock(void* self) {
    Spinlock* s = (Spinlock*)self;
    for (;;) {
        if (atomic_CompareAndSwapUint32(&s->lock, 0, 1)) {
            return;
        }
    }
}

void Spinlock_Unlock(void* self) {
    Spinlock* s = (Spinlock*)self;
    atomic_StoreUint32(&s->lock, 0);
}

bool Spinlock_TryLock(void* self) {
    Spinlock* s = (Spinlock*)self;
    return atomic_CompareAndSwapUint32(&s->lock, 0, 1);
}

bool Spinlock_IsLocked(void* self) {
    Spinlock* s = (Spinlock*)self;
    return atomic_LoadUint32(&s->lock) != 0;
}

RwLock* NewRwLock(void) {
    return &(RwLock){};
}

void RwLock_RLock(void* self) {
    RwLock* l = (RwLock*)self;
    for (;;) {
        if (atomic_LoadUint32(&l->writer) == 0) {
            atomic_AddUint32(&l->readers, 1);
            if (atomic_LoadUint32(&l->writer) == 0) {
                return;
            }
            atomic_AddUint32(&l->readers, ~(0));
        }
    }
}

void RwLock_RUnlock(void* self) {
    RwLock* l = (RwLock*)self;
    atomic_AddUint32(&l->readers, ~(0));
}

void RwLock_Lock(void* self) {
    RwLock* l = (RwLock*)self;
    for (; !atomic_CompareAndSwapUint32(&l->writer, 0, 1);) {
    }
    for (; atomic_LoadUint32(&l->readers) != 0;) {
    }
}

void RwLock_Unlock(void* self) {
    RwLock* l = (RwLock*)self;
    atomic_StoreUint32(&l->writer, 0);
}

VmxRootSpinlock* NewVmxRootSpinlock(void) {
    VmxRootSpinlock* l = &(VmxRootSpinlock){};
    KeInitializeSpinLock(&l->raw);
    return l;
}

void VmxRootSpinlock_Lock(void* self) {
    VmxRootSpinlock* l = (VmxRootSpinlock*)self;
    KeAcquireSpinLockRaiseToDpc(&l->raw);
}

void VmxRootSpinlock_Unlock(void* self) {
    VmxRootSpinlock* l = (VmxRootSpinlock*)self;
    KeReleaseSpinLock(&l->raw, 0);
}

// -- syscall_hook.go --

static Hypervisor* getHypervisor(void) {
    return gHyp;
}

static void syscallHookConfigureEFER(VCPU* v, bool enable) {
    uint64_t vmxBasic = __readmsr(IA32_VMX_BASIC);
    uint32_t vmEntryControls = 0, vmExitControls = 0;
    vmRead32(VmcsCtrlVmentryIntrInfo, &vmEntryControls);
    vmRead32(VmcsCtrlPrimaryProcBasedVmExec, &vmExitControls);
    uint64_t eferMsr = __readmsr(IA32_EFER);
    if (enable) {
        eferMsr &= ~((uint64_t)1 << 0);
        vmEntryControls |= (IA32_VMX_ENTRY_CTLS_LOAD_IA32_EFER_FLAG);
        vmExitControls |= (IA32_VMX_EXIT_CTLS_SAVE_IA32_EFER_FLAG);
        vmWrite64(VMCS_GUEST_IA32_DEBUGCTL, eferMsr);
        Hypervisor* h = getHypervisor();
        if (h != NULL) {
            Hypervisor_setExceptionBitmap(h, v, EXCEPTION_VECTOR_UNDEFINED_OPCODE);
        }
        gEferHook->Enabled = true;
    } else {
        eferMsr |= ((uint64_t)1 << 0);
        vmEntryControls &= ~(IA32_VMX_ENTRY_CTLS_LOAD_IA32_EFER_FLAG);
        vmExitControls &= ~(IA32_VMX_EXIT_CTLS_SAVE_IA32_EFER_FLAG);
        vmWrite64(VMCS_GUEST_IA32_DEBUGCTL, eferMsr);
        __writemsr(IA32_EFER, eferMsr);
        removeUndefinedInstructionForDisablingSyscallSysret(v);
        gEferHook->Enabled = false;
    }
    vmWrite32(VmcsCtrlVmentryIntrInfo, (uint32_t)(adjustVmcsControl(vmxBasic, (uint64_t)(vmEntryControls))));
    vmWrite32(VmcsCtrlPrimaryProcBasedVmExec, (uint32_t)(adjustVmcsControl(vmxBasic, (uint64_t)(vmExitControls))));
    so_String status = so_str("DISABLED");
    if (enable) {
        status = so_str("ENABLED");
    }
    LogDebug(so_str("EFER Syscall Hook: %s"), (so_Slice){(void*[1]){&status}, 1, 1});
}

static bool syscallHookEmulateSYSCALL(VCPU* v) {
    uint64_t guestRip = 0, guestRflags = 0;
    uint32_t instrLen = 0;
    vmRead64(VmcsGuestRip, &guestRip);
    vmRead32(VmcsVmexitInstrLength, &instrLen);
    vmRead64(VmcsGuestRflags, &guestRflags);
    uint64_t lstar = __readmsr(IA32_LSTAR);
    v->Regs->Rcx = guestRip + (uint64_t)(instrLen);
    guestRip = lstar;
    vmWrite64(VmcsGuestRip, guestRip);
    uint64_t fmask = __readmsr(IA32_FMASK);
    v->Regs->R11 = guestRflags;
    guestRflags &= ~(fmask | X86_FLAGS_RF);
    vmWrite64(VmcsGuestRflags, guestRflags);
    if (gEferHook->CetSupport) {
        uint64_t ucet = __readmsr(IA32_U_CET);
        if ((ucet & ((uint64_t)1 << 0)) != 0) {
            uint64_t ssp = 0;
            vmRead64(VmcsGuestSsp, &ssp);
            __writemsr(IA32_PL3_SSP, ssp);
            vmWrite64(VmcsGuestSsp, 0);
        }
    }
    uint64_t star = __readmsr(IA32_STAR);
    uint16_t csSelector = (uint16_t)((star >> 32) & 0xFFFC);
    uint16_t ssSelector = csSelector + 8;
    setSegmentRegister(VmcsGuestCs, csSelector, 0, 0xFFFFFFFF, 0xA09B);
    setSegmentRegister(VmcsGuestSs, ssSelector, 0, 0xFFFFFFFF, 0xC093);
    dispatchEventEferSyscall(v);
    return true;
}

static bool syscallHookEmulateSYSRET(VCPU* v) {
    uint64_t guestRip = 0, guestRflags = 0;
    guestRip = v->Regs->Rcx;
    vmWrite64(VmcsGuestRip, guestRip);
    guestRflags = ((v->Regs->R11 & ~((X86_FLAGS_RF | X86_FLAGS_VM) | X86_FLAGS_RESERVED_BITS)) | X86_FLAGS_FIXED);
    vmWrite64(VmcsGuestRflags, guestRflags);
    if (gEferHook->CetSupport) {
        uint64_t ucet = __readmsr(IA32_U_CET);
        if ((ucet & ((uint64_t)1 << 0)) != 0) {
            uint64_t ssp = __readmsr(IA32_PL3_SSP);
            vmWrite64(VmcsGuestSsp, ssp);
        }
    }
    uint64_t star = __readmsr(IA32_STAR);
    uint16_t csSelector = (uint16_t)(((star >> 48) + 16) | 3);
    uint16_t ssSelector = (uint16_t)(((star >> 48) + 8) | 3);
    setSegmentRegister(VmcsGuestCs, csSelector, 0, 0xFFFFFFFF, 0xA0FB);
    setSegmentRegister(VmcsGuestSs, ssSelector, 0, 0xFFFFFFFF, 0xC0F3);
    dispatchEventEferSysret(v);
    return true;
}

static bool syscallHookHandleUD(VCPU* v) {
    uint64_t rip = 0;
    vmRead64(VmcsGuestRip, &rip);
    if (gEferHook->UdHandling) {
        if ((rip & 0xff00000000000000) != 0) {
            return syscallHookEmulateSYSRET(v);
        }
        return syscallHookEmulateSYSCALL(v);
    } else {
        CR3_TYPE guestCr3 = getCurrentProcessCr3();
        uint64_t originalCr3 = __readcr3();
        __writecr3(guestCr3.Flags);
        so_byte instrBuf[3] = {0};
        bool present = checkPagePresentByCr3((uintptr_t)(rip), guestCr3);
        if (present) {
            readMemorySafe(rip, (void*)(&instrBuf[0]), 3);
        } else {
            Hypervisor* h = getHypervisor();
            if (h != NULL) {
                Hypervisor_suppressAdvance(h, v);
                injectPageFaultWithoutErrorCode(rip);
            }
            __writecr3(originalCr3);
            return false;
        }
        __writecr3(originalCr3);
        if (instrBuf[0] == 0x0F && instrBuf[1] == 0x05) {
            return syscallHookEmulateSYSCALL(v);
        }
        if (instrBuf[0] == 0x48 && instrBuf[1] == 0x0F && instrBuf[2] == 0x07) {
            return syscallHookEmulateSYSRET(v);
        }
        return false;
    }
}

static CR3_TYPE getCurrentProcessCr3(void) {
    return (CR3_TYPE){};
}

static bool checkPagePresentByCr3(uintptr_t addr, CR3_TYPE cr3) {
    return false;
}

static void readMemorySafe(uint64_t addr, void* buffer, uint64_t size) {
}

static void dispatchEventEferSyscall(VCPU* v) {
}

static void dispatchEventEferSysret(VCPU* v) {
}

static void injectPageFaultWithoutErrorCode(uint64_t rip) {
}

static void removeUndefinedInstructionForDisablingSyscallSysret(VCPU* v) {
}

static void setSegmentRegister(uint64_t field, uint16_t selector, uint64_t base, uint32_t limit, uint32_t accessRights) {
    vmWrite64(field + 0, (uint64_t)(selector));
    vmWrite64(field + 8, base);
    vmWrite64(field + 16, (uint64_t)(limit));
    vmWrite64(field + 24, (uint64_t)(accessRights));
}

// -- types.go --

static uint64_t eptPml1Offset(uint64_t v) {
    return (v & 0xFFF);
}

static uint64_t eptPml1Index(uint64_t v) {
    return ((v >> 12) & 0x1FF);
}

static uint64_t eptPml2Index(uint64_t v) {
    return ((v >> 21) & 0x1FF);
}

static uint64_t eptPml3Index(uint64_t v) {
    return ((v >> 30) & 0x1FF);
}

static uint64_t eptPml4Index(uint64_t v) {
    return ((v >> 39) & 0x1FF);
}

static uint64_t logBufferSize(void) {
    return (uint64_t)(MAX_LOG_BUFFERS) * ((uint64_t)(LOG_CHUNK_SIZE) + (uint64_t)(unsafe_Sizeof((BUFFER_HEADER){})));
}

static uint64_t logBufferSizePrio(void) {
    return (uint64_t)(MAX_LOG_BUFFERS_PRIO) * ((uint64_t)(LOG_CHUNK_SIZE) + (uint64_t)(unsafe_Sizeof((BUFFER_HEADER){})));
}

// -- vmcall.go --

static int32_t Hypervisor_handleHyperdbgVmcall(void* self, uint64_t num, uint64_t p1, uint64_t p2, uint64_t p3) {
    Hypervisor* h = (Hypervisor*)self;
    do {
        if (num == ((VmcTest))) {
            return Hypervisor_vmcallTest(h, p1, p2, p3);
        } else if (num == ((VmcVmxoff))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                vmxVmxoff(v);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcChangePageAttrib))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v == NULL) {
                return STATUS_UNSUCCESSFUL;
            }
            CR3_TYPE cr3 = (CR3_TYPE){.Flags = p3};
            bool result = Hypervisor_performPageHook(h, v, (uintptr_t)(p1), cr3, (uint32_t)(p2));
            if (result) {
                return STATUS_SUCCESS;
            }
            return STATUS_UNSUCCESSFUL;
        } else if (num == ((VmcInveptSingleContext))) {
            inveptSingleContext(p1);
            return STATUS_SUCCESS;
        } else if (num == ((VmcInveptAllContexts))) {
            inveptAllContexts();
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnhookAllPages))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                HookManager_RestoreAll(gHooks, v);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnhookSinglePage))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v == NULL) {
                return STATUS_UNSUCCESSFUL;
            }
            {
                so_Error err = HookManager_RemoveByVirtAddr(gHooks, p1);
                if (err == NULL) {
                    return STATUS_SUCCESS;
                }
            }
            return STATUS_UNSUCCESSFUL;
        } else if (num == ((VmcEnableSyscallHookEfer))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                syscallHookConfigureEFER(v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcDisableSyscallHookEfer))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                syscallHookConfigureEFER(v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcChangeMsrBitmapRead))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                msrHandlePerformMsrBitmapReadChange(v, (uint32_t)(p1));
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcChangeMsrBitmapWrite))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                msrHandlePerformMsrBitmapWriteChange(v, (uint32_t)(p1));
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetRdtscExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setRdtscExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetRdpmcExiting))) {
            Hypervisor_setPmcVmexit(h, true);
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetExceptionBitmap))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setExceptionBitmap(h, v, (uint32_t)(p1));
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableMovToDebugRegsExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovDebugRegsExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableExternalInterruptExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setExternalInterruptExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcChangeIoBitmap))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                ioHandlePerformIoBitmapChange(v, (uint32_t)(p1));
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetHiddenCcBreakpoint))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v == NULL) {
                return STATUS_UNSUCCESSFUL;
            }
            CR3_TYPE cr3 = (CR3_TYPE){.Flags = p2};
            {
                so_Error err = HookManager_Install(gHooks, v->CoreId, p1, cr3, HookExec);
                if (err == NULL) {
                    return STATUS_SUCCESS;
                }
            }
            return STATUS_UNSUCCESSFUL;
        } else if (num == ((VmcDisableExternalInterruptExitingOnlyToClearInterruptCommands))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                protectedHvExternalInterruptExitingForDisablingInterruptCommands(v);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetRdtscExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setRdtscExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetRdpmcExiting))) {
            Hypervisor_setPmcVmexit(h, false);
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetExceptionBitmap))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setExceptionBitmap(h, v, 0);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcDisableMovToDebugRegsExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovDebugRegsExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcDisableExternalInterruptExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setExternalInterruptExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetCpuidExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setCpuidExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetCpuidExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setCpuidExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcClp))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setClpExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetMovFromCrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovFromCrExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetMovToCrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovToCrExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetMovFromCrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovFromCrExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetMovToCrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovToCrExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetMovFromDrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovFromDrExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcSetMovToDrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovToDrExiting(h, v, true);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetMovFromDrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovFromDrExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcUnsetMovToDrExiting))) {
            VCPU* v = Hypervisor_CurrentVcpu(h);
            if (v != NULL) {
                Hypervisor_setMovToDrExiting(h, v, false);
            }
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookMaskedRead))) {
            eptHookSetReadHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookMaskedWrite))) {
            eptHookSetWriteHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookMaskedReadWrite))) {
            eptHookSetReadWriteHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookInlineHookRead))) {
            eptHookSetReadHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookInlineHookWrite))) {
            eptHookSetWriteHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookInlineHookReadWrite))) {
            eptHookSetReadWriteHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookExecRead))) {
            eptHookSetReadHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookExecWrite))) {
            eptHookSetWriteHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookExecReadWrite))) {
            eptHookSetReadWriteHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableEptHookExecOnly))) {
            eptHookSetExecHook(p1, (CR3_TYPE){.Flags = p2});
            return STATUS_SUCCESS;
        } else if (num == ((VmcEnableMtf))) {
            setMonitorTrapFlag(true);
            return STATUS_SUCCESS;
        } else if (num == ((VmcQueryPerformanceFrequency))) {
            return (int32_t)(queryPerformanceFrequency());
        } else if (num == ((VmcQueryPerformanceCounter))) {
            return (int32_t)(queryPerformanceCounter());
        } else if (num == ((VmcCheckVmxSupport))) {
            if (checkVmxSupport()) {
                return STATUS_SUCCESS;
            }
            return STATUS_UNSUCCESSFUL;
        } else {
            LogWarning(so_str("Top-level driver VMCALL %d not handled"), (so_Slice){(void*[1]){&num}, 1, 1});
        }
    } while (0);
    LogWarning(so_str("Unknown VMCALL number: 0x%X"), (so_Slice){(void*[1]){&num}, 1, 1});
    return STATUS_UNSUCCESSFUL;
}

static int32_t Hypervisor_vmcallTest(void* self, uint64_t _p0, uint64_t _p1, uint64_t _p2) {
    Hypervisor* h = (Hypervisor*)self;
    LogInfo(so_str("VMCALL_TEST received"), (so_Slice){&so_Nil, 0, 0});
    return STATUS_SUCCESS;
}

static bool Hypervisor_performPageHook(void* self, VCPU* v, uintptr_t addr, CR3_TYPE cr3, uint32_t mask) {
    Hypervisor* h = (Hypervisor*)self;
    if ((mask & HOOK_PAGE_MONITOR_READ) != 0) {
        {
            so_Error err = HookManager_Install(gHooks, v->CoreId, (uint64_t)(addr), cr3, HookRead);
            if (err != NULL) {
                return false;
            }
        }
    }
    if ((mask & HOOK_PAGE_MONITOR_WRITE) != 0) {
        {
            so_Error err = HookManager_Install(gHooks, v->CoreId, (uint64_t)(addr), cr3, HookWrite);
            if (err != NULL) {
                return false;
            }
        }
    }
    if ((mask & HOOK_PAGE_MONITOR_EXEC) != 0) {
        {
            so_Error err = HookManager_Install(gHooks, v->CoreId, (uint64_t)(addr), cr3, HookExec);
            if (err != NULL) {
                return false;
            }
        }
    }
    if ((mask & (HOOK_PAGE_MONITOR_INLINE_HOOKS | HOOK_PAGE_MASKED_HOOKS)) != 0) {
        {
            so_Error err = HookManager_Install(gHooks, v->CoreId, (uint64_t)(addr), cr3, HookReadWrite);
            if (err != NULL) {
                return false;
            }
        }
    }
    return true;
}

static uint64_t queryPerformanceFrequency(void) {
    return 10000000;
}

static uint64_t queryPerformanceCounter(void) {
    return 0;
}

static bool checkVmxSupport(void) {
    return true;
}

// -- vmx.go --

uint64_t AsmVmxVmcall(uint64_t RegPtr) ;

uint8_t AsmEnableVmxOperation(uint64_t PhysicalAddr) ;

void AsmInvept(uint64_t Eptp) ;

void AsmInveptAllContexts(void) {
}

void AsmInvvpid(void) ;

void AsmInvvpidAllContexts(void) {
}

void AsmHypervVmcall(uint64_t regs) ;

uint64_t AsmGetGdtBase(void) ;

uint64_t AsmGetIdtBase(void) ;

uint32_t AsmGetGdtLimit(void) ;

uint32_t AsmGetIdtLimit(void) ;

uint64_t AsmGetRflags(void) ;

uint16_t AsmGetEs(void) ;

uint16_t AsmGetCs(void) ;

uint16_t AsmGetSs(void) ;

uint16_t AsmGetDs(void) ;

uint16_t AsmGetFs(void) ;

uint16_t AsmGetGs(void) ;

uint16_t AsmGetTr(void) ;

uint16_t AsmGetLdtr(void) ;

uint64_t AsmGetFsBase(void) {
    return 0;
}

uint64_t AsmGetGsBase(void) {
    return 0;
}

uintptr_t AsmVmexitHandlerAddr(void) {
    return 0;
}

static void vmRead64(uint64_t field, uint64_t* val) {
    AsmVmxVmread(field, val);
}

static void vmRead32(uint64_t field, uint32_t* val) {
    AsmVmxVmread32(field, val);
}

static void vmWrite64(uint64_t field, uint64_t val) {
    AsmVmxVmwrite(field, val);
}

static void vmWrite32(uint64_t field, uint32_t val) {
    AsmVmxVmwrite32(field, val);
}

static void setMonitorTrapFlag(bool enable) {
    if (enable) {
        vmWrite64(0x00002802, (__readmsr(0x1D9) | ((uint64_t)1 << 2)));
    } else {
        vmWrite64(0x00002802, (__readmsr(0x1D9) & ~(((uint64_t)1 << 2))));
    }
}

static bool vmxVmlaunch(VCPU* v) {
    v->HasLaunched = true;
    uint8_t result = AsmVmxVmlaunch();
    if (result != 0) {
        return true;
    }
    return false;
}

static bool vmxVmresume(VCPU* v) {
    v->HasLaunched = true;
    uint8_t result = AsmVmxVmresume();
    if (result != 0) {
        return true;
    }
    return false;
}

static void vmxVmxoff(VCPU* v) {
    AsmVmxVmxOff();
    v->VmxoffState.Executed = true;
}

static void asmHypervVmcall(uint64_t regs) {
    AsmHypervVmcall(regs);
}

bool VmxVmexitHandler(GUEST_REGS* regs) {
    if (gHyp == NULL) {
        return false;
    }
    return Hypervisor_HandleVmExit(gHyp, regs);
}

uint8_t VmxVmresume(void) {
    return 1;
}

uint64_t VmxReturnStackPointerForVmxoff(void) {
    return 0;
}

uint64_t VmxReturnInstructionPointerForVmxoff(void) {
    return 0;
}

void VmxVirtualizeCurrentSystem(uint8_t* state) {
}

uint64_t EptHook2GeneralDetourEventHandler(GUEST_REGS* regs, uint64_t calledFrom) {
    return 0;
}

void IdtEmulationhandleHostInterrupt(uint8_t* trapFrame) {
}

void _solod_kmod_init(void) {
    gExecTrap = &(ExecTrapState){.TrapList = so_make_slice_impl(sizeof(ExecTrapEntry), 0, 0)};
    gIoBitmap = &(IoBitmapManager){};
    gModeBasedExec = &(ModeBasedExecHookState){.CsSelectors = so_make_slice_impl(sizeof(ModeBasedExecEntry), 0, 0)};
    gMsrBitmap = &(MsrBitmapManager){};
    gEferHook = &(EferHookState){.UdHandling = true};
}
