package wdk

import "unsafe"

//so:extern
func DbgPrintEx(ComponentId uint32, Level uint32, Format *byte) uint32 { panic(`solod:extern`) }

//so:extern
func DbgPrintReturnControlC(Format *byte) uint32 { panic(`solod:extern`) }

//so:extern
func DbgPrompt(Prompt *byte, Response *byte, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func DbgQueryDebugFilterState(ComponentId uint32, Level uint32) uint32 { panic(`solod:extern`) }

//so:extern
func DbgSetDebugFilterState(ComponentId uint32, Level uint32, State byte) uint32 { panic(`solod:extern`) }

//so:extern
func EtwEventEnabled(RegHandle uintptr, EventDescriptor uintptr) byte { panic(`solod:extern`) }

//so:extern
func HalAcquireDisplayOwnership(ResetDisplayParameters PHAL_RESET_DISPLAY_PARAMETERS) { panic(`solod:extern`) }

//so:extern
func HalAllocateAdapterChannel(AdapterObject uintptr, Wcb WAIT_CONTEXT_BLOCK, NumberOfMapRegisters uint32, ExecutionRoutine DRIVER_CONTROL) uint32 { panic(`solod:extern`) }

//so:extern
func HalAllocateCommonBuffer(AdapterObject uintptr, Length uint32, LogicalAddress *int64, CacheEnabled byte) uintptr { panic(`solod:extern`) }

//so:extern
func HalAllocateCrashDumpRegisters(AdapterObject uintptr, NumberOfMapRegisters *uint32) uintptr { panic(`solod:extern`) }

//so:extern
func HalAllocateHardwareCounters(GroupAffinty uintptr, GroupCount uint32, ResourceList PHYSICAL_COUNTER_RESOURCE_LIST, CounterSetHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func HalAssignSlotResources(RegistryPath *UNICODE_STRING, DriverClassName *UNICODE_STRING, DriverObject DRIVER_OBJECT, DeviceObject DEVICE_OBJECT, BusType INTERFACE_TYPE, BusNumber uint32, SlotNumber uint32, AllocatedResources CM_RESOURCE_LIST) uint32 { panic(`solod:extern`) }

//so:extern
func HalBugCheckSystem(ErrorSource uintptr, ErrorRecord WHEA_ERROR_RECORD) { panic(`solod:extern`) }

//so:extern
func HalDmaAllocateCrashDumpRegistersEx(Adapter uintptr, NumberOfMapRegisters uint32, Type HAL_DMA_CRASH_DUMP_REGISTER_TYPE, MapRegisterBase **unsafe.Pointer, MapRegistersAvailable *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func HalDmaFreeCrashDumpRegistersEx(Adapter uintptr, Type HAL_DMA_CRASH_DUMP_REGISTER_TYPE) uint32 { panic(`solod:extern`) }

//so:extern
func HalFreeCommonBuffer(AdapterObject uintptr, Length uint32, LogicalAddress int64, VirtualAddress uintptr, CacheEnabled byte) { panic(`solod:extern`) }

//so:extern
func HalFreeHardwareCounters(CounterSetHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func HalGetAdapter(DeviceDescription DEVICE_DESCRIPTION, NumberOfMapRegisters *uint32) uintptr { panic(`solod:extern`) }

//so:extern
func HalGetBusData(BusDataType BUS_DATA_TYPE, BusNumber uint32, SlotNumber uint32, Buffer uintptr, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func HalGetBusDataByOffset(BusDataType BUS_DATA_TYPE, BusNumber uint32, SlotNumber uint32, Buffer uintptr, Offset uint32, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func HalGetInterruptVector(InterfaceType INTERFACE_TYPE, BusNumber uint32, BusInterruptLevel uint32, BusInterruptVector uint32, Irql *byte, Affinity *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func HalMakeBeep(Frequency uint32) byte { panic(`solod:extern`) }

//so:extern
func HalReadDmaCounter(AdapterObject uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func HalSetBusData(BusDataType BUS_DATA_TYPE, BusNumber uint32, SlotNumber uint32, Buffer uintptr, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func HalSetBusDataByOffset(BusDataType BUS_DATA_TYPE, BusNumber uint32, SlotNumber uint32, Buffer uintptr, Offset uint32, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func HalTranslateBusAddress(InterfaceType INTERFACE_TYPE, BusNumber uint32, BusAddress int64, AddressSpace *uint32, TranslatedAddress *int64) byte { panic(`solod:extern`) }

//so:extern
func IoFlushAdapterBuffers(AdapterObject uintptr, Mdl MDL, MapRegisterBase uintptr, CurrentVa uintptr, Length uint32, WriteToDevice byte) byte { panic(`solod:extern`) }

//so:extern
func IoFreeAdapterChannel(AdapterObject uintptr) { panic(`solod:extern`) }

//so:extern
func IoFreeMapRegisters(AdapterObject uintptr, MapRegisterBase uintptr, NumberOfMapRegisters uint32) { panic(`solod:extern`) }

//so:extern
func IoMapTransfer(AdapterObject uintptr, Mdl MDL, MapRegisterBase uintptr, CurrentVa uintptr, Length *uint32, WriteToDevice byte) int64 { panic(`solod:extern`) }

//so:extern
func KeFlushWriteBuffer() { panic(`solod:extern`) }

//so:extern
func KeQueryPerformanceCounter(PerformanceFrequency *int64) int64 { panic(`solod:extern`) }

//so:extern
func NtAccessCheckAndAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, DesiredAccess uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtAccessCheckByTypeAndAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, PrincipalSelfSid uintptr, DesiredAccess uint32, AuditType uint32, Flags uint32, ObjectTypeList uintptr, ObjectTypeListLength uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtAccessCheckByTypeResultListAndAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, PrincipalSelfSid uintptr, DesiredAccess uint32, AuditType uint32, Flags uint32, ObjectTypeList uintptr, ObjectTypeListLength uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtAccessCheckByTypeResultListAndAuditAlarmByHandle(SubsystemName *UNICODE_STRING, HandleId uintptr, ClientToken uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, PrincipalSelfSid uintptr, DesiredAccess uint32, AuditType uint32, Flags uint32, ObjectTypeList uintptr, ObjectTypeListLength uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtAdjustGroupsToken(TokenHandle uintptr, ResetToDefault byte, NewState uintptr, BufferLength uint32, PreviousState uintptr, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtAdjustPrivilegesToken(TokenHandle uintptr, DisableAllPrivileges byte, NewState uintptr, BufferLength uint32, PreviousState uintptr, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtAllocateLocallyUniqueId(Luid *LUID) uint32 { panic(`solod:extern`) }

//so:extern
func NtAllocateVirtualMemory(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, ZeroBits uintptr, RegionSize *uintptr, AllocationType uint32, Protect uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtAllocateVirtualMemoryEx(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, RegionSize *uintptr, AllocationType uint32, PageProtection uint32, ExtendedParameters uintptr, ExtendedParameterCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCancelIoFileEx(FileHandle uintptr, IoRequestToCancel *IO_STATUS_BLOCK, IoStatusBlock *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func NtCancelTimer(TimerHandle uintptr, CurrentState *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtClose(Handle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtCloseObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, GenerateOnClose byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtCommitComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtCommitEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtCommitRegistryTransaction(TransactionHandle uintptr, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCommitTransaction(TransactionHandle uintptr, Wait byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateDirectoryObject(DirectoryHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateEnlistment(EnlistmentHandle *uintptr, DesiredAccess uint32, ResourceManagerHandle uintptr, TransactionHandle uintptr, ObjectAttributes OBJECT_ATTRIBUTES, CreateOptions uint32, NotificationMask uint32, EnlistmentKey uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateEvent(EventHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, EventType uintptr, InitialState byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateFile(FileHandle *uintptr, DesiredAccess uintptr, ObjectAttributes OBJECT_ATTRIBUTES, IoStatusBlock *IO_STATUS_BLOCK, AllocationSize *int64, FileAttributes uintptr, ShareAccess uintptr, CreateDisposition NTCREATEFILE_CREATE_DISPOSITION, CreateOptions NTCREATEFILE_CREATE_OPTIONS, EaBuffer uintptr, EaLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateKey(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TitleIndex uint32, Class *UNICODE_STRING, CreateOptions uint32, Disposition *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateKeyTransacted(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TitleIndex uint32, Class *UNICODE_STRING, CreateOptions uint32, TransactionHandle uintptr, Disposition *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateRegistryTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, CreateOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateResourceManager(ResourceManagerHandle *uintptr, DesiredAccess uint32, TmHandle uintptr, RmGuid *GUID, ObjectAttributes OBJECT_ATTRIBUTES, CreateOptions uint32, Description *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateSection(SectionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, MaximumSize *int64, SectionPageProtection uint32, AllocationAttributes uint32, FileHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateSectionEx(SectionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, MaximumSize *int64, SectionPageProtection uint32, AllocationAttributes uint32, FileHandle uintptr, ExtendedParameters uintptr, ExtendedParameterCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateTimer(TimerHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TimerType uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, Uow *GUID, TmHandle uintptr, CreateOptions uint32, IsolationLevel uint32, IsolationFlags uint32, Timeout *int64, Description *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtCreateTransactionManager(TmHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, LogFileName *UNICODE_STRING, CreateOptions uint32, CommitStrength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtDeleteFile(ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtDeleteKey(KeyHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtDeleteObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, GenerateOnClose byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtDeleteValueKey(KeyHandle uintptr, ValueName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtDeviceIoControlFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, IoControlCode uint32, InputBuffer uintptr, InputBufferLength uint32, OutputBuffer uintptr, OutputBufferLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtDisplayString(String *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtDuplicateObject(SourceProcessHandle uintptr, SourceHandle uintptr, TargetProcessHandle uintptr, TargetHandle *uintptr, DesiredAccess uint32, HandleAttributes uint32, Options uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtDuplicateToken(ExistingTokenHandle uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, EffectiveOnly byte, TokenType uintptr, NewTokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtEnumerateKey(KeyHandle uintptr, Index uint32, KeyInformationClass KEY_INFORMATION_CLASS, KeyInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtEnumerateTransactionObject(RootObjectHandle uintptr, QueryType uintptr, ObjectCursor uintptr, ObjectCursorLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtEnumerateValueKey(KeyHandle uintptr, Index uint32, KeyValueInformationClass KEY_VALUE_INFORMATION_CLASS, KeyValueInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtFilterToken(ExistingTokenHandle uintptr, Flags uint32, SidsToDisable uintptr, PrivilegesToDelete uintptr, RestrictedSids uintptr, NewTokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtFlushBuffersFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func NtFlushBuffersFileEx(FileHandle uintptr, Flags uint32, Parameters uintptr, ParametersSize uint32, IoStatusBlock *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func NtFlushKey(KeyHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtFlushVirtualMemory(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, RegionSize *uintptr, IoStatus *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func NtFreeVirtualMemory(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, RegionSize *uintptr, FreeType uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtFsControlFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, FsControlCode uint32, InputBuffer uintptr, InputBufferLength uint32, OutputBuffer uintptr, OutputBufferLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtGetNotificationResourceManager(ResourceManagerHandle uintptr, TransactionNotification uintptr, NotificationLength uint32, Timeout *int64, ReturnLength *uint32, Asynchronous uint32, AsynchronousContext uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtImpersonateAnonymousToken(ThreadHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtLoadDriver(DriverServiceName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtLockFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, ByteOffset *int64, Length *int64, Key uint32, FailImmediately byte, ExclusiveLock byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtMakeTemporaryObject(Handle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtManagePartition(TargetHandle uintptr, SourceHandle uintptr, PartitionInformationClass PARTITION_INFORMATION_CLASS, PartitionInformation uintptr, PartitionInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtMapViewOfSection(SectionHandle uintptr, ProcessHandle uintptr, BaseAddress **unsafe.Pointer, ZeroBits uintptr, CommitSize uintptr, SectionOffset *int64, ViewSize *uintptr, InheritDisposition SECTION_INHERIT, AllocationType uint32, Win32Protect uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtNotifyChangeKey(KeyHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, CompletionFilter uint32, WatchTree byte, Buffer uintptr, BufferSize uint32, Asynchronous byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtNotifyChangeMultipleKeys(MasterKeyHandle uintptr, Count uint32, SubordinateObjects OBJECT_ATTRIBUTES, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, CompletionFilter uint32, WatchTree byte, Buffer uintptr, BufferSize uint32, Asynchronous byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenDirectoryObject(DirectoryHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenEnlistment(EnlistmentHandle *uintptr, DesiredAccess uint32, ResourceManagerHandle uintptr, EnlistmentGuid *GUID, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenEvent(EventHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenFile(FileHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, IoStatusBlock *IO_STATUS_BLOCK, ShareAccess uint32, OpenOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenKey(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenKeyEx(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, OpenOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenKeyTransacted(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TransactionHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenKeyTransactedEx(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, OpenOptions uint32, TransactionHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, ClientToken uintptr, DesiredAccess uint32, GrantedAccess uint32, Privileges *PRIVILEGE_SET, ObjectCreation byte, AccessGranted byte, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenProcess(ProcessHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, ClientId uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenProcessToken(ProcessHandle uintptr, DesiredAccess uint32, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenProcessTokenEx(ProcessHandle uintptr, DesiredAccess uint32, HandleAttributes uint32, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenRegistryTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenResourceManager(ResourceManagerHandle *uintptr, DesiredAccess uint32, TmHandle uintptr, ResourceManagerGuid *GUID, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenSection(SectionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenSymbolicLinkObject(LinkHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenThreadToken(ThreadHandle uintptr, DesiredAccess uint32, OpenAsSelf byte, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenThreadTokenEx(ThreadHandle uintptr, DesiredAccess uint32, OpenAsSelf byte, HandleAttributes uint32, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenTimer(TimerHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, Uow *GUID, TmHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtOpenTransactionManager(TmHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, LogFileName *UNICODE_STRING, TmIdentity *GUID, OpenOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtPowerInformation(InformationLevel uintptr, InputBuffer uintptr, InputBufferLength uint32, OutputBuffer uintptr, OutputBufferLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrePrepareComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrePrepareEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrepareComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrepareEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrivilegeCheck(ClientToken uintptr, RequiredPrivileges *PRIVILEGE_SET, Result *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrivilegeObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ClientToken uintptr, DesiredAccess uint32, Privileges *PRIVILEGE_SET, AccessGranted byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtPrivilegedServiceAuditAlarm(SubsystemName *UNICODE_STRING, ServiceName *UNICODE_STRING, ClientToken uintptr, Privileges *PRIVILEGE_SET, AccessGranted byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtPropagationComplete(ResourceManagerHandle uintptr, RequestCookie uint32, BufferLength uint32, Buffer uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtPropagationFailed(ResourceManagerHandle uintptr, RequestCookie uint32, PropStatus uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryDirectoryFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS, ReturnSingleEntry byte, FileName *UNICODE_STRING, RestartScan byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryDirectoryFileEx(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS, QueryFlags uint32, FileName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryDirectoryObject(DirectoryHandle uintptr, Buffer uintptr, Length uint32, ReturnSingleEntry byte, RestartScan byte, Context *uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryEaFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ReturnSingleEntry byte, EaList uintptr, EaListLength uint32, EaIndex *uint32, RestartScan byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryFullAttributesFile(ObjectAttributes OBJECT_ATTRIBUTES, FileInformation FILE_NETWORK_OPEN_INFORMATION) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationByName(ObjectAttributes OBJECT_ATTRIBUTES, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationEnlistment(EnlistmentHandle uintptr, EnlistmentInformationClass uintptr, EnlistmentInformation uintptr, EnlistmentInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationProcess(ProcessHandle uintptr, ProcessInformationClass PROCESSINFOCLASS, ProcessInformation uintptr, ProcessInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationResourceManager(ResourceManagerHandle uintptr, ResourceManagerInformationClass uintptr, ResourceManagerInformation uintptr, ResourceManagerInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationThread(ThreadHandle uintptr, ThreadInformationClass THREADINFOCLASS, ThreadInformation uintptr, ThreadInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationToken(TokenHandle uintptr, TokenInformationClass uintptr, TokenInformation uintptr, TokenInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationTransaction(TransactionHandle uintptr, TransactionInformationClass uintptr, TransactionInformation uintptr, TransactionInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryInformationTransactionManager(TransactionManagerHandle uintptr, TransactionManagerInformationClass uintptr, TransactionManagerInformation uintptr, TransactionManagerInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryKey(KeyHandle uintptr, KeyInformationClass KEY_INFORMATION_CLASS, KeyInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryMultipleValueKey(KeyHandle uintptr, ValueEntries KEY_VALUE_ENTRY, EntryCount uint32, ValueBuffer uintptr, BufferLength *uint32, RequiredBufferLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryObject(Handle uintptr, ObjectInformationClass OBJECT_INFORMATION_CLASS, ObjectInformation uintptr, ObjectInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryQuotaInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ReturnSingleEntry byte, SidList uintptr, SidListLength uint32, StartSid uintptr, RestartScan byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtQuerySecurityObject(Handle uintptr, SecurityInformation uint32, SecurityDescriptor uintptr, Length uint32, LengthNeeded *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQuerySymbolicLinkObject(LinkHandle uintptr, LinkTarget *UNICODE_STRING, ReturnedLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQuerySystemInformation(SystemInformationClass SYSTEM_INFORMATION_CLASS, SystemInformation uintptr, SystemInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQuerySystemTime(SystemTime *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryTimerResolution(MaximumTime *uint32, MinimumTime *uint32, CurrentTime *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryValueKey(KeyHandle uintptr, ValueName *UNICODE_STRING, KeyValueInformationClass KEY_VALUE_INFORMATION_CLASS, KeyValueInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryVirtualMemory(ProcessHandle uintptr, BaseAddress uintptr, MemoryInformationClass MEMORY_INFORMATION_CLASS, MemoryInformation uintptr, MemoryInformationLength uintptr, ReturnLength *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtQueryVolumeInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FsInformation uintptr, Length uint32, FsInformationClass FS_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func NtReadFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ByteOffset *int64, Key *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtReadOnlyEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtRecoverEnlistment(EnlistmentHandle uintptr, EnlistmentKey uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtRecoverResourceManager(ResourceManagerHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtRecoverTransactionManager(TransactionManagerHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtRegisterProtocolAddressInformation(ResourceManager uintptr, ProtocolId *GUID, ProtocolInformationSize uint32, ProtocolInformation uintptr, CreateOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtRenameKey(KeyHandle uintptr, NewName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtRenameTransactionManager(LogFileName *UNICODE_STRING, ExistingTransactionManagerGuid *GUID) uint32 { panic(`solod:extern`) }

//so:extern
func NtRestoreKey(KeyHandle uintptr, FileHandle uintptr, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtRollbackComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtRollbackEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtRollbackRegistryTransaction(TransactionHandle uintptr, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtRollbackTransaction(TransactionHandle uintptr, Wait byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtRollforwardTransactionManager(TransactionManagerHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtSaveKey(KeyHandle uintptr, FileHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtSaveKeyEx(KeyHandle uintptr, FileHandle uintptr, Format uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetEaFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetEvent(EventHandle uintptr, PreviousState *int32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationEnlistment(EnlistmentHandle uintptr, EnlistmentInformationClass uintptr, EnlistmentInformation uintptr, EnlistmentInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationKey(KeyHandle uintptr, KeySetInformationClass KEY_SET_INFORMATION_CLASS, KeySetInformation uintptr, KeySetInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationResourceManager(ResourceManagerHandle uintptr, ResourceManagerInformationClass uintptr, ResourceManagerInformation uintptr, ResourceManagerInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationThread(ThreadHandle uintptr, ThreadInformationClass THREADINFOCLASS, ThreadInformation uintptr, ThreadInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationToken(TokenHandle uintptr, TokenInformationClass uintptr, TokenInformation uintptr, TokenInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationTransaction(TransactionHandle uintptr, TransactionInformationClass uintptr, TransactionInformation uintptr, TransactionInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationTransactionManager(TmHandle uintptr, TransactionManagerInformationClass uintptr, TransactionManagerInformation uintptr, TransactionManagerInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetInformationVirtualMemory(ProcessHandle uintptr, VmInformationClass VIRTUAL_MEMORY_INFORMATION_CLASS, NumberOfEntries uintptr, VirtualAddresses MEMORY_RANGE_ENTRY, VmInformation uintptr, VmInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetQuotaInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetSecurityObject(Handle uintptr, SecurityInformation uint32, SecurityDescriptor uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetTimer(TimerHandle uintptr, DueTime *int64, TimerApcRoutine PTIMER_APC_ROUTINE, TimerContext uintptr, ResumeTimer byte, Period int32, PreviousState *byte) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetTimerEx(TimerHandle uintptr, TimerSetInformationClass TIMER_SET_INFORMATION_CLASS, TimerSetInformation uintptr, TimerSetInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetValueKey(KeyHandle uintptr, ValueName *UNICODE_STRING, TitleIndex uint32, Type uint32, Data uintptr, DataSize uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtSetVolumeInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FsInformation uintptr, Length uint32, FsInformationClass FS_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func NtSinglePhaseReject(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtTerminateProcess(ProcessHandle uintptr, ExitStatus uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtUnloadDriver(DriverServiceName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func NtUnlockFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, ByteOffset *int64, Length *int64, Key uint32) uint32 { panic(`solod:extern`) }

//so:extern
func NtUnmapViewOfSection(ProcessHandle uintptr, BaseAddress uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func NtWaitForSingleObject(Handle uintptr, Alertable byte, Timeout *int64) uint32 { panic(`solod:extern`) }

//so:extern
func NtWriteFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ByteOffset *int64, Key *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func PfxFindPrefix(PrefixTable PREFIX_TABLE, FullName *STRING) PREFIX_TABLE_ENTRY { panic(`solod:extern`) }

//so:extern
func PfxInitialize(PrefixTable PREFIX_TABLE) { panic(`solod:extern`) }

//so:extern
func PfxInsertPrefix(PrefixTable PREFIX_TABLE, Prefix *STRING, PrefixTableEntry PREFIX_TABLE_ENTRY) byte { panic(`solod:extern`) }

//so:extern
func PfxRemovePrefix(PrefixTable PREFIX_TABLE, PrefixTableEntry PREFIX_TABLE_ENTRY) { panic(`solod:extern`) }

//so:extern
func RtlAbsoluteToSelfRelativeSD(AbsoluteSecurityDescriptor uintptr, SelfRelativeSecurityDescriptor uintptr, BufferLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAddAccessAllowedAce(Acl uintptr, AceRevision uint32, AccessMask uint32, Sid uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAddAccessAllowedAceEx(Acl uintptr, AceRevision uint32, AceFlags uint32, AccessMask uint32, Sid uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAddAce(Acl uintptr, AceRevision uint32, StartingAceIndex uint32, AceList uintptr, AceListLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAllocateAndInitializeSid(IdentifierAuthority uintptr, SubAuthorityCount byte, SubAuthority0 uint32, SubAuthority1 uint32, SubAuthority2 uint32, SubAuthority3 uint32, SubAuthority4 uint32, SubAuthority5 uint32, SubAuthority6 uint32, SubAuthority7 uint32, Sid *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAllocateAndInitializeSidEx(IdentifierAuthority uintptr, SubAuthorityCount byte, SubAuthorities *uint32, Sid *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAllocateHeap(HeapHandle uintptr, Flags uint32, Size uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlAppendStringToString(Destination *STRING, Source *STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAppendUnicodeStringToString(Destination *UNICODE_STRING, Source *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAppendUnicodeToString(Destination *UNICODE_STRING, Source *uint16) uint32 { panic(`solod:extern`) }

//so:extern
func RtlAreBitsClear(BitMapHeader RTL_BITMAP, StartingIndex uint32, Length uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlAreBitsSet(BitMapHeader RTL_BITMAP, StartingIndex uint32, Length uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlAssert(VoidFailedAssertion uintptr, VoidFileName uintptr, LineNumber uint32, MutableMessage *byte) { panic(`solod:extern`) }

//so:extern
func RtlCheckRegistryKey(RelativeTo uint32, Path *uint16) uint32 { panic(`solod:extern`) }

//so:extern
func RtlClearAllBits(BitMapHeader RTL_BITMAP) { panic(`solod:extern`) }

//so:extern
func RtlClearBit(BitMapHeader RTL_BITMAP, BitNumber uint32) { panic(`solod:extern`) }

//so:extern
func RtlClearBits(BitMapHeader RTL_BITMAP, StartingIndex uint32, NumberToClear uint32) { panic(`solod:extern`) }

//so:extern
func RtlCmDecodeMemIoResource(Descriptor CM_PARTIAL_RESOURCE_DESCRIPTOR, Start *uint64) uint64 { panic(`solod:extern`) }

//so:extern
func RtlCmEncodeMemIoResource(Descriptor CM_PARTIAL_RESOURCE_DESCRIPTOR, Type byte, Length uint64, Start uint64) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCompareAltitudes(Altitude1 *UNICODE_STRING, Altitude2 *UNICODE_STRING) int32 { panic(`solod:extern`) }

//so:extern
func RtlCompareMemoryUlong(Source uintptr, Length uintptr, Pattern uint32) uintptr { panic(`solod:extern`) }

//so:extern
func RtlCompareString(String1 *STRING, String2 *STRING, CaseInSensitive byte) int32 { panic(`solod:extern`) }

//so:extern
func RtlCompareUnicodeString(String1 *UNICODE_STRING, String2 *UNICODE_STRING, CaseInSensitive byte) int32 { panic(`solod:extern`) }

//so:extern
func RtlCompareUnicodeStrings(String1 *uint16, String1Length uintptr, String2 *uint16, String2Length uintptr, CaseInSensitive byte) int32 { panic(`solod:extern`) }

//so:extern
func RtlCompressBuffer(CompressionFormatAndEngine uint16, UncompressedBuffer *byte, UncompressedBufferSize uint32, CompressedBuffer *byte, CompressedBufferSize uint32, UncompressedChunkSize uint32, FinalCompressedSize *uint32, WorkSpace uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlContractHashTable(HashTable RTL_DYNAMIC_HASH_TABLE) byte { panic(`solod:extern`) }

//so:extern
func RtlCopyBitMap(Source RTL_BITMAP, Destination RTL_BITMAP, TargetBit uint32) { panic(`solod:extern`) }

//so:extern
func RtlCopyLuid(DestinationLuid *LUID, SourceLuid *LUID) { panic(`solod:extern`) }

//so:extern
func RtlCopySid(DestinationSidLength uint32, DestinationSid uintptr, SourceSid uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCopyString(DestinationString *STRING, SourceString *STRING) { panic(`solod:extern`) }

//so:extern
func RtlCopyUnicodeString(DestinationString *UNICODE_STRING, SourceString *UNICODE_STRING) { panic(`solod:extern`) }

//so:extern
func RtlCreateAcl(Acl uintptr, AclLength uint32, AclRevision uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCreateHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Shift uint32, Flags uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlCreateHashTableEx(HashTable RTL_DYNAMIC_HASH_TABLE, InitialSize uint32, Shift uint32, Flags uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlCreateHeap(Flags uint32, HeapBase uintptr, ReserveSize uintptr, CommitSize uintptr, Lock uintptr, Parameters RTL_HEAP_PARAMETERS) uintptr { panic(`solod:extern`) }

//so:extern
func RtlCreateRegistryKey(RelativeTo uint32, Path *uint16) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCreateSecurityDescriptor(SecurityDescriptor uintptr, Revision uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCreateServiceSid(ServiceName *UNICODE_STRING, ServiceSid uintptr, ServiceSidLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCreateSystemVolumeInformationFolder(VolumeRootPath *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCreateUnicodeString(DestinationString *UNICODE_STRING, SourceString *uint16) byte { panic(`solod:extern`) }

//so:extern
func RtlCreateVirtualAccountSid(Name *UNICODE_STRING, BaseSubAuthority uint32, Sid uintptr, SidLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlCustomCPToUnicodeN(CustomCP CPTABLEINFO, UnicodeString *uint16, MaxBytesInUnicodeString uint32, BytesInUnicodeString *uint32, CustomCPString *byte, BytesInCustomCPString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDecompressBuffer(CompressionFormat uint16, UncompressedBuffer *byte, UncompressedBufferSize uint32, CompressedBuffer *byte, CompressedBufferSize uint32, FinalUncompressedSize *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDecompressBufferEx(CompressionFormat uint16, UncompressedBuffer *byte, UncompressedBufferSize uint32, CompressedBuffer *byte, CompressedBufferSize uint32, FinalUncompressedSize *uint32, WorkSpace uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDecompressFragment(CompressionFormat uint16, UncompressedFragment *byte, UncompressedFragmentSize uint32, CompressedBuffer *byte, CompressedBufferSize uint32, FragmentOffset uint32, FinalUncompressedSize *uint32, WorkSpace uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDelete(Links RTL_SPLAY_LINKS) RTL_SPLAY_LINKS { panic(`solod:extern`) }

//so:extern
func RtlDeleteAce(Acl uintptr, AceIndex uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDeleteElementGenericTable(Table RTL_GENERIC_TABLE, Buffer uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlDeleteElementGenericTableAvl(Table RTL_AVL_TABLE, Buffer uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlDeleteElementGenericTableAvlEx(Table RTL_AVL_TABLE, NodeOrParent uintptr) { panic(`solod:extern`) }

//so:extern
func RtlDeleteHashTable(HashTable RTL_DYNAMIC_HASH_TABLE) { panic(`solod:extern`) }

//so:extern
func RtlDeleteNoSplay(Links RTL_SPLAY_LINKS, Root RTL_SPLAY_LINKS) { panic(`solod:extern`) }

//so:extern
func RtlDeleteRegistryValue(RelativeTo uint32, Path *uint16, ValueName *uint16) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDestroyHeap(HeapHandle uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlDosPathNameToNtPathName_U_WithStatus(DosFileName *uint16, NtFileName *UNICODE_STRING, FilePart **uint16, Reserved uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDowncaseUnicodeChar(SourceCharacter uint16) uint16 { panic(`solod:extern`) }

//so:extern
func RtlDowncaseUnicodeString(DestinationString *UNICODE_STRING, SourceString *UNICODE_STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlDuplicateUnicodeString(Flags uint32, StringIn *UNICODE_STRING, StringOut *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlEndEnumerationHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) { panic(`solod:extern`) }

//so:extern
func RtlEndStrongEnumerationHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) { panic(`solod:extern`) }

//so:extern
func RtlEndWeakEnumerationHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) { panic(`solod:extern`) }

//so:extern
func RtlEnumerateEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) RTL_DYNAMIC_HASH_TABLE_ENTRY { panic(`solod:extern`) }

//so:extern
func RtlEnumerateGenericTable(Table RTL_GENERIC_TABLE, Restart byte) uintptr { panic(`solod:extern`) }

//so:extern
func RtlEnumerateGenericTableAvl(Table RTL_AVL_TABLE, Restart byte) uintptr { panic(`solod:extern`) }

//so:extern
func RtlEnumerateGenericTableLikeADirectory(Table RTL_AVL_TABLE, MatchFunction PRTL_AVL_MATCH_FUNCTION, MatchData uintptr, NextFlag uint32, RestartKey **unsafe.Pointer, DeleteCount *uint32, Buffer uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlEnumerateGenericTableWithoutSplaying(Table RTL_GENERIC_TABLE, RestartKey **unsafe.Pointer) uintptr { panic(`solod:extern`) }

//so:extern
func RtlEnumerateGenericTableWithoutSplayingAvl(Table RTL_AVL_TABLE, RestartKey **unsafe.Pointer) uintptr { panic(`solod:extern`) }

//so:extern
func RtlEqualPrefixSid(Sid1 uintptr, Sid2 uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlEqualSid(Sid1 uintptr, Sid2 uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlEqualString(String1 *STRING, String2 *STRING, CaseInSensitive byte) byte { panic(`solod:extern`) }

//so:extern
func RtlEqualUnicodeString(String1 *UNICODE_STRING, String2 *UNICODE_STRING, CaseInSensitive byte) byte { panic(`solod:extern`) }

//so:extern
func RtlExpandHashTable(HashTable RTL_DYNAMIC_HASH_TABLE) byte { panic(`solod:extern`) }

//so:extern
func RtlExtractBitMap(Source RTL_BITMAP, Destination RTL_BITMAP, TargetBit uint32, NumberOfBits uint32) { panic(`solod:extern`) }

//so:extern
func RtlFindClearBits(BitMapHeader RTL_BITMAP, NumberToFind uint32, HintIndex uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindClearBitsAndSet(BitMapHeader RTL_BITMAP, NumberToFind uint32, HintIndex uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindClearRuns(BitMapHeader RTL_BITMAP, RunArray RTL_BITMAP_RUN, SizeOfRunArray uint32, LocateLongestRuns byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindClosestEncodableLength(SourceLength uint64, TargetLength *uint64) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindLastBackwardRunClear(BitMapHeader RTL_BITMAP, FromIndex uint32, StartingRunIndex *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindLeastSignificantBit(Set uint64) int8 { panic(`solod:extern`) }

//so:extern
func RtlFindLongestRunClear(BitMapHeader RTL_BITMAP, StartingIndex *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindMostSignificantBit(Set uint64) int8 { panic(`solod:extern`) }

//so:extern
func RtlFindNextForwardRunClear(BitMapHeader RTL_BITMAP, FromIndex uint32, StartingRunIndex *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindSetBits(BitMapHeader RTL_BITMAP, NumberToFind uint32, HintIndex uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFindSetBitsAndClear(BitMapHeader RTL_BITMAP, NumberToFind uint32, HintIndex uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFreeHeap(HeapHandle uintptr, Flags uint32, BaseAddress uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlFreeSid(Sid uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlFreeUTF8String(utf8String *STRING) { panic(`solod:extern`) }

//so:extern
func RtlGUIDFromString(GuidString *UNICODE_STRING, Guid *GUID) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGenerate8dot3Name(Name *UNICODE_STRING, AllowExtendedCharacters byte, Context GENERATE_NAME_CONTEXT, Name8dot3 *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetAce(Acl uintptr, AceIndex uint32, Ace **unsafe.Pointer) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetActiveConsoleId() uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetCallersAddress(CallersAddress **unsafe.Pointer, CallersCaller **unsafe.Pointer) { panic(`solod:extern`) }

//so:extern
func RtlGetCompressionWorkSpaceSize(CompressionFormatAndEngine uint16, CompressBufferWorkSpaceSize *uint32, CompressFragmentWorkSpaceSize *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetConsoleSessionForegroundProcessId() uint64 { panic(`solod:extern`) }

//so:extern
func RtlGetDaclSecurityDescriptor(SecurityDescriptor uintptr, DaclPresent *byte, Dacl **uintptr, DaclDefaulted *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetElementGenericTable(Table RTL_GENERIC_TABLE, I uint32) uintptr { panic(`solod:extern`) }

//so:extern
func RtlGetElementGenericTableAvl(Table RTL_AVL_TABLE, I uint32) uintptr { panic(`solod:extern`) }

//so:extern
func RtlGetEnabledExtendedFeatures(FeatureMask uint64) uint64 { panic(`solod:extern`) }

//so:extern
func RtlGetGroupSecurityDescriptor(SecurityDescriptor uintptr, Group *uintptr, GroupDefaulted *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetNextEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Context RTL_DYNAMIC_HASH_TABLE_CONTEXT) RTL_DYNAMIC_HASH_TABLE_ENTRY { panic(`solod:extern`) }

//so:extern
func RtlGetNtProductType(NtProductType uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlGetNtSystemRoot() *uint16 { panic(`solod:extern`) }

//so:extern
func RtlGetOwnerSecurityDescriptor(SecurityDescriptor uintptr, Owner *uintptr, OwnerDefaulted *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetPersistedStateLocation(SourceID *uint16, CustomValue *uint16, DefaultPath *uint16, StateLocationType STATE_LOCATION_TYPE, TargetPath *uint16, BufferLengthIn uint32, BufferLengthOut *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetSaclSecurityDescriptor(SecurityDescriptor uintptr, SaclPresent *byte, Sacl **uintptr, SaclDefaulted *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetSuiteMask() uint32 { panic(`solod:extern`) }

//so:extern
func RtlGetVersion(lpVersionInformation uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlHashUnicodeString(String *UNICODE_STRING, CaseInSensitive byte, HashAlgorithm uint32, HashValue *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIdentifierAuthoritySid(Sid uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlIdnToAscii(Flags uint32, SourceString *uint16, SourceStringLength int32, DestinationString *uint16, DestinationStringLength *int32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIdnToNameprepUnicode(Flags uint32, SourceString *uint16, SourceStringLength int32, DestinationString *uint16, DestinationStringLength *int32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIdnToUnicode(Flags uint32, SourceString *uint16, SourceStringLength int32, DestinationString *uint16, DestinationStringLength *int32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlInitCodePageTable(TableBase *uint16, CodePageTable CPTABLEINFO) { panic(`solod:extern`) }

//so:extern
func RtlInitEnumerationHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) byte { panic(`solod:extern`) }

//so:extern
func RtlInitStrongEnumerationHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) byte { panic(`solod:extern`) }

//so:extern
func RtlInitUTF8String(DestinationString *STRING, SourceString *int8) { panic(`solod:extern`) }

//so:extern
func RtlInitUTF8StringEx(DestinationString *STRING, SourceString *int8) uint32 { panic(`solod:extern`) }

//so:extern
func RtlInitUnicodeStringEx(DestinationString *UNICODE_STRING, SourceString *uint16) uint32 { panic(`solod:extern`) }

//so:extern
func RtlInitWeakEnumerationHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) byte { panic(`solod:extern`) }

//so:extern
func RtlInitializeBitMap(BitMapHeader RTL_BITMAP, BitMapBuffer *uint32, SizeOfBitMap uint32) { panic(`solod:extern`) }

//so:extern
func RtlInitializeGenericTable(Table RTL_GENERIC_TABLE, CompareRoutine PRTL_GENERIC_COMPARE_ROUTINE, AllocateRoutine PRTL_GENERIC_ALLOCATE_ROUTINE, FreeRoutine PRTL_GENERIC_FREE_ROUTINE, TableContext uintptr) { panic(`solod:extern`) }

//so:extern
func RtlInitializeGenericTableAvl(Table RTL_AVL_TABLE, CompareRoutine PRTL_AVL_COMPARE_ROUTINE, AllocateRoutine PRTL_AVL_ALLOCATE_ROUTINE, FreeRoutine PRTL_AVL_FREE_ROUTINE, TableContext uintptr) { panic(`solod:extern`) }

//so:extern
func RtlInitializeSid(Sid uintptr, IdentifierAuthority uintptr, SubAuthorityCount byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlInitializeSidEx(Sid uintptr, IdentifierAuthority uintptr, SubAuthorityCount byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlInsertElementGenericTable(Table RTL_GENERIC_TABLE, Buffer uintptr, BufferSize uint32, NewElement *byte) uintptr { panic(`solod:extern`) }

//so:extern
func RtlInsertElementGenericTableAvl(Table RTL_AVL_TABLE, Buffer uintptr, BufferSize uint32, NewElement *byte) uintptr { panic(`solod:extern`) }

//so:extern
func RtlInsertElementGenericTableFull(Table RTL_GENERIC_TABLE, Buffer uintptr, BufferSize uint32, NewElement *byte, NodeOrParent uintptr, SearchResult TABLE_SEARCH_RESULT) uintptr { panic(`solod:extern`) }

//so:extern
func RtlInsertElementGenericTableFullAvl(Table RTL_AVL_TABLE, Buffer uintptr, BufferSize uint32, NewElement *byte, NodeOrParent uintptr, SearchResult TABLE_SEARCH_RESULT) uintptr { panic(`solod:extern`) }

//so:extern
func RtlInsertEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Entry RTL_DYNAMIC_HASH_TABLE_ENTRY, Signature uintptr, Context RTL_DYNAMIC_HASH_TABLE_CONTEXT) byte { panic(`solod:extern`) }

//so:extern
func RtlInt64ToUnicodeString(Value uint64, Base uint32, String *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIntegerToUnicodeString(Value uint32, Base uint32, String *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIoDecodeMemIoResource(Descriptor IO_RESOURCE_DESCRIPTOR, Alignment *uint64, MinimumAddress *uint64, MaximumAddress *uint64) uint64 { panic(`solod:extern`) }

//so:extern
func RtlIoEncodeMemIoResource(Descriptor IO_RESOURCE_DESCRIPTOR, Type byte, Length uint64, Alignment uint64, MinimumAddress uint64, MaximumAddress uint64) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIsApiSetImplemented(apiSetName *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIsCloudFilesPlaceholder(FileAttributes uint32, ReparseTag uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlIsDosDeviceName_U(DosFileName *uint16) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIsGenericTableEmpty(Table RTL_GENERIC_TABLE) byte { panic(`solod:extern`) }

//so:extern
func RtlIsGenericTableEmptyAvl(Table RTL_AVL_TABLE) byte { panic(`solod:extern`) }

//so:extern
func RtlIsMultiSessionSku() byte { panic(`solod:extern`) }

//so:extern
func RtlIsMultiUsersInSessionSku() byte { panic(`solod:extern`) }

//so:extern
func RtlIsNonEmptyDirectoryReparsePointAllowed(ReparseTag uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlIsNormalizedString(NormForm uint32, SourceString *uint16, SourceStringLength int32, Normalized *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIsPartialPlaceholder(FileAttributes uint32, ReparseTag uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlIsPartialPlaceholderFileHandle(FileHandle uintptr, IsPartialPlaceholder *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIsPartialPlaceholderFileInfo(InfoBuffer uintptr, InfoClass FILE_INFORMATION_CLASS, IsPartialPlaceholder *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlIsStateSeparationEnabled() byte { panic(`solod:extern`) }

//so:extern
func RtlIsUntrustedObject(Handle uintptr, Object uintptr, UntrustedObject *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlLengthRequiredSid(SubAuthorityCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlLengthSecurityDescriptor(SecurityDescriptor uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlLengthSid(Sid uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlLookupElementGenericTable(Table RTL_GENERIC_TABLE, Buffer uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlLookupElementGenericTableAvl(Table RTL_AVL_TABLE, Buffer uintptr) uintptr { panic(`solod:extern`) }

//so:extern
func RtlLookupElementGenericTableFull(Table RTL_GENERIC_TABLE, Buffer uintptr, NodeOrParent **unsafe.Pointer, SearchResult TABLE_SEARCH_RESULT) uintptr { panic(`solod:extern`) }

//so:extern
func RtlLookupElementGenericTableFullAvl(Table RTL_AVL_TABLE, Buffer uintptr, NodeOrParent **unsafe.Pointer, SearchResult TABLE_SEARCH_RESULT) uintptr { panic(`solod:extern`) }

//so:extern
func RtlLookupEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Signature uintptr, Context RTL_DYNAMIC_HASH_TABLE_CONTEXT) RTL_DYNAMIC_HASH_TABLE_ENTRY { panic(`solod:extern`) }

//so:extern
func RtlLookupFirstMatchingElementGenericTableAvl(Table RTL_AVL_TABLE, Buffer uintptr, RestartKey **unsafe.Pointer) uintptr { panic(`solod:extern`) }

//so:extern
func RtlMapGenericMask(AccessMask *uint32, GenericMapping uintptr) { panic(`solod:extern`) }

//so:extern
func RtlMultiByteToUnicodeN(UnicodeString *uint16, MaxBytesInUnicodeString uint32, BytesInUnicodeString *uint32, MultiByteString *byte, BytesInMultiByteString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlMultiByteToUnicodeSize(BytesInUnicodeString *uint32, MultiByteString *byte, BytesInMultiByteString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNormalizeSecurityDescriptor(SecurityDescriptor *uintptr, SecurityDescriptorLength uint32, NewSecurityDescriptor *uintptr, NewSecurityDescriptorLength *uint32, CheckOnly byte) byte { panic(`solod:extern`) }

//so:extern
func RtlNormalizeString(NormForm uint32, SourceString *uint16, SourceStringLength int32, DestinationString *uint16, DestinationStringLength *int32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNtStatusToDosErrorNoTeb(Status uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberGenericTableElements(Table RTL_GENERIC_TABLE) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberGenericTableElementsAvl(Table RTL_AVL_TABLE) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberOfClearBits(BitMapHeader RTL_BITMAP) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberOfClearBitsInRange(BitMapHeader RTL_BITMAP, StartingIndex uint32, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberOfSetBits(BitMapHeader RTL_BITMAP) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberOfSetBitsInRange(BitMapHeader RTL_BITMAP, StartingIndex uint32, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlNumberOfSetBitsUlongPtr(Target uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlOemStringToUnicodeString(DestinationString *UNICODE_STRING, SourceString *STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlOemToUnicodeN(UnicodeString *uint16, MaxBytesInUnicodeString uint32, BytesInUnicodeString *uint32, OemString *byte, BytesInOemString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlPrefixString(String1 *STRING, String2 *STRING, CaseInSensitive byte) byte { panic(`solod:extern`) }

//so:extern
func RtlPrefixUnicodeString(String1 *UNICODE_STRING, String2 *UNICODE_STRING, CaseInSensitive byte) byte { panic(`solod:extern`) }

//so:extern
func RtlQueryPackageIdentity(TokenObject uintptr, PackageFullName *uint16, PackageSize *uintptr, AppId *uint16, AppIdSize *uintptr, Packaged *byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlQueryPackageIdentityEx(TokenObject uintptr, PackageFullName *uint16, PackageSize *uintptr, AppId *uint16, AppIdSize *uintptr, DynamicId *GUID, Flags *uint64) uint32 { panic(`solod:extern`) }

//so:extern
func RtlQueryProcessPlaceholderCompatibilityMode() byte { panic(`solod:extern`) }

//so:extern
func RtlQueryRegistryValueWithFallback(PrimaryHandle uintptr, FallbackHandle uintptr, ValueName *UNICODE_STRING, ValueLength uint32, ValueType *uint32, ValueData uintptr, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlQueryRegistryValues(RelativeTo uint32, Path *uint16, QueryTable RTL_QUERY_REGISTRY_TABLE, Context uintptr, Environment uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlQueryThreadPlaceholderCompatibilityMode() byte { panic(`solod:extern`) }

//so:extern
func RtlQueryValidationRunlevel(ComponentName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRandom(Seed *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRandomEx(Seed *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRealPredecessor(Links RTL_SPLAY_LINKS) RTL_SPLAY_LINKS { panic(`solod:extern`) }

//so:extern
func RtlRealSuccessor(Links RTL_SPLAY_LINKS) RTL_SPLAY_LINKS { panic(`solod:extern`) }

//so:extern
func RtlRemoveEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Entry RTL_DYNAMIC_HASH_TABLE_ENTRY, Context RTL_DYNAMIC_HASH_TABLE_CONTEXT) byte { panic(`solod:extern`) }

//so:extern
func RtlReplaceSidInSd(SecurityDescriptor uintptr, OldSid uintptr, NewSid uintptr, NumChanges *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRunOnceBeginInitialize(RunOnce uintptr, Flags uint32, Context **unsafe.Pointer) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRunOnceComplete(RunOnce uintptr, Flags uint32, Context uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRunOnceExecuteOnce(RunOnce uintptr, InitFn PRTL_RUN_ONCE_INIT_FN, Parameter uintptr, Context **unsafe.Pointer) uint32 { panic(`solod:extern`) }

//so:extern
func RtlRunOnceInitialize(RunOnce uintptr) { panic(`solod:extern`) }

//so:extern
func RtlSecondsSince1970ToTime(ElapsedSeconds uint32, Time *int64) { panic(`solod:extern`) }

//so:extern
func RtlSecondsSince1980ToTime(ElapsedSeconds uint32, Time *int64) { panic(`solod:extern`) }

//so:extern
func RtlSelfRelativeToAbsoluteSD(SelfRelativeSecurityDescriptor uintptr, AbsoluteSecurityDescriptor uintptr, AbsoluteSecurityDescriptorSize *uint32, Dacl uintptr, DaclSize *uint32, Sacl uintptr, SaclSize *uint32, Owner uintptr, OwnerSize *uint32, PrimaryGroup uintptr, PrimaryGroupSize *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlSetAllBits(BitMapHeader RTL_BITMAP) { panic(`solod:extern`) }

//so:extern
func RtlSetBit(BitMapHeader RTL_BITMAP, BitNumber uint32) { panic(`solod:extern`) }

//so:extern
func RtlSetBits(BitMapHeader RTL_BITMAP, StartingIndex uint32, NumberToSet uint32) { panic(`solod:extern`) }

//so:extern
func RtlSetDaclSecurityDescriptor(SecurityDescriptor uintptr, DaclPresent byte, Dacl uintptr, DaclDefaulted byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlSetGroupSecurityDescriptor(SecurityDescriptor uintptr, Group uintptr, GroupDefaulted byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlSetOwnerSecurityDescriptor(SecurityDescriptor uintptr, Owner uintptr, OwnerDefaulted byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlSetProcessPlaceholderCompatibilityMode(Mode byte) byte { panic(`solod:extern`) }

//so:extern
func RtlSetThreadPlaceholderCompatibilityMode(Mode byte) byte { panic(`solod:extern`) }

//so:extern
func RtlSplay(Links RTL_SPLAY_LINKS) RTL_SPLAY_LINKS { panic(`solod:extern`) }

//so:extern
func RtlStringFromGUID(Guid *GUID, GuidString *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlStronglyEnumerateEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) RTL_DYNAMIC_HASH_TABLE_ENTRY { panic(`solod:extern`) }

//so:extern
func RtlSubAuthorityCountSid(Sid uintptr) *byte { panic(`solod:extern`) }

//so:extern
func RtlSubAuthoritySid(Sid uintptr, SubAuthority uint32) *uint32 { panic(`solod:extern`) }

//so:extern
func RtlSubtreePredecessor(Links RTL_SPLAY_LINKS) RTL_SPLAY_LINKS { panic(`solod:extern`) }

//so:extern
func RtlSubtreeSuccessor(Links RTL_SPLAY_LINKS) RTL_SPLAY_LINKS { panic(`solod:extern`) }

//so:extern
func RtlTestBit(BitMapHeader RTL_BITMAP, BitNumber uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlTimeFieldsToTime(TimeFields TIME_FIELDS, Time *int64) byte { panic(`solod:extern`) }

//so:extern
func RtlTimeToSecondsSince1980(Time *int64, ElapsedSeconds *uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlTimeToTimeFields(Time *int64, TimeFields TIME_FIELDS) { panic(`solod:extern`) }

//so:extern
func RtlUTF8StringToUnicodeString(DestinationString *UNICODE_STRING, SourceString *STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUTF8ToUnicodeN(UnicodeStringDestination *uint16, UnicodeStringMaxByteCount uint32, UnicodeStringActualByteCount *uint32, UTF8StringSource *byte, UTF8StringByteCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeStringToCountedOemString(DestinationString *STRING, SourceString *UNICODE_STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeStringToInteger(String *UNICODE_STRING, Base uint32, Value *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeStringToUTF8String(DestinationString *STRING, SourceString *UNICODE_STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeToCustomCPN(CustomCP CPTABLEINFO, CustomCPString *byte, MaxBytesInCustomCPString uint32, BytesInCustomCPString *uint32, UnicodeString *uint16, BytesInUnicodeString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeToMultiByteN(MultiByteString *byte, MaxBytesInMultiByteString uint32, BytesInMultiByteString *uint32, UnicodeString *uint16, BytesInUnicodeString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeToOemN(OemString *byte, MaxBytesInOemString uint32, BytesInOemString *uint32, UnicodeString *uint16, BytesInUnicodeString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUnicodeToUTF8N(UTF8StringDestination *byte, UTF8StringMaxByteCount uint32, UTF8StringActualByteCount *uint32, UnicodeStringSource *uint16, UnicodeStringByteCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeChar(SourceCharacter uint16) uint16 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeString(DestinationString *UNICODE_STRING, SourceString *UNICODE_STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeStringToCountedOemString(DestinationString *STRING, SourceString *UNICODE_STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeStringToOemString(DestinationString *STRING, SourceString *UNICODE_STRING, AllocateDestinationString byte) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeToCustomCPN(CustomCP CPTABLEINFO, CustomCPString *byte, MaxBytesInCustomCPString uint32, BytesInCustomCPString *uint32, UnicodeString *uint16, BytesInUnicodeString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeToMultiByteN(MultiByteString *byte, MaxBytesInMultiByteString uint32, BytesInMultiByteString *uint32, UnicodeString *uint16, BytesInUnicodeString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpcaseUnicodeToOemN(OemString *byte, MaxBytesInOemString uint32, BytesInOemString *uint32, UnicodeString *uint16, BytesInUnicodeString uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlUpperChar(Character byte) byte { panic(`solod:extern`) }

//so:extern
func RtlUpperString(DestinationString *STRING, SourceString *STRING) { panic(`solod:extern`) }

//so:extern
func RtlValidRelativeSecurityDescriptor(SecurityDescriptorInput uintptr, SecurityDescriptorLength uint32, RequiredInformation uint32) byte { panic(`solod:extern`) }

//so:extern
func RtlValidSecurityDescriptor(SecurityDescriptor uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlValidSid(Sid uintptr) byte { panic(`solod:extern`) }

//so:extern
func RtlValidateUnicodeString(Flags uint32, String *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlVerifyVersionInfo(VersionInfo uintptr, TypeMask uint32, ConditionMask uint64) uint32 { panic(`solod:extern`) }

//so:extern
func RtlWalkFrameChain(Callers **unsafe.Pointer, Count uint32, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlWeaklyEnumerateEntryHashTable(HashTable RTL_DYNAMIC_HASH_TABLE, Enumerator RTL_DYNAMIC_HASH_TABLE_ENUMERATOR) RTL_DYNAMIC_HASH_TABLE_ENTRY { panic(`solod:extern`) }

//so:extern
func RtlWriteRegistryValue(RelativeTo uint32, Path *uint16, ValueName *uint16, ValueType uint32, ValueData uintptr, ValueLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func RtlxAnsiStringToUnicodeSize(AnsiString *STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlxOemStringToUnicodeSize(OemString *STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlxUnicodeStringToAnsiSize(UnicodeString *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func RtlxUnicodeStringToOemSize(UnicodeString *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAccessCheckAndAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, DesiredAccess uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAccessCheckByTypeAndAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, PrincipalSelfSid uintptr, DesiredAccess uint32, AuditType uint32, Flags uint32, ObjectTypeList uintptr, ObjectTypeListLength uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAccessCheckByTypeResultListAndAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, PrincipalSelfSid uintptr, DesiredAccess uint32, AuditType uint32, Flags uint32, ObjectTypeList uintptr, ObjectTypeListLength uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAccessCheckByTypeResultListAndAuditAlarmByHandle(SubsystemName *UNICODE_STRING, HandleId uintptr, ClientToken uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, PrincipalSelfSid uintptr, DesiredAccess uint32, AuditType uint32, Flags uint32, ObjectTypeList uintptr, ObjectTypeListLength uint32, GenericMapping uintptr, ObjectCreation byte, GrantedAccess *uint32, AccessStatus *int32, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAdjustGroupsToken(TokenHandle uintptr, ResetToDefault byte, NewState uintptr, BufferLength uint32, PreviousState uintptr, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAdjustPrivilegesToken(TokenHandle uintptr, DisableAllPrivileges byte, NewState uintptr, BufferLength uint32, PreviousState uintptr, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAllocateLocallyUniqueId(Luid *LUID) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAllocateVirtualMemory(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, ZeroBits uintptr, RegionSize *uintptr, AllocationType uint32, Protect uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwAllocateVirtualMemoryEx(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, RegionSize *uintptr, AllocationType uint32, PageProtection uint32, ExtendedParameters uintptr, ExtendedParameterCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCancelIoFileEx(FileHandle uintptr, IoRequestToCancel *IO_STATUS_BLOCK, IoStatusBlock *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCancelTimer(TimerHandle uintptr, CurrentState *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwClose(Handle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCloseObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, GenerateOnClose byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCommitComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCommitEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCommitRegistryTransaction(TransactionHandle uintptr, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCommitTransaction(TransactionHandle uintptr, Wait byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateDirectoryObject(DirectoryHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateEnlistment(EnlistmentHandle *uintptr, DesiredAccess uint32, ResourceManagerHandle uintptr, TransactionHandle uintptr, ObjectAttributes OBJECT_ATTRIBUTES, CreateOptions uint32, NotificationMask uint32, EnlistmentKey uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateEvent(EventHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, EventType uintptr, InitialState byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateFile(FileHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, IoStatusBlock *IO_STATUS_BLOCK, AllocationSize *int64, FileAttributes uint32, ShareAccess uint32, CreateDisposition uint32, CreateOptions uint32, EaBuffer uintptr, EaLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateKey(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TitleIndex uint32, Class *UNICODE_STRING, CreateOptions uint32, Disposition *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateKeyTransacted(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TitleIndex uint32, Class *UNICODE_STRING, CreateOptions uint32, TransactionHandle uintptr, Disposition *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateRegistryTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, CreateOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateResourceManager(ResourceManagerHandle *uintptr, DesiredAccess uint32, TmHandle uintptr, ResourceManagerGuid *GUID, ObjectAttributes OBJECT_ATTRIBUTES, CreateOptions uint32, Description *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateSection(SectionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, MaximumSize *int64, SectionPageProtection uint32, AllocationAttributes uint32, FileHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateSectionEx(SectionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, MaximumSize *int64, SectionPageProtection uint32, AllocationAttributes uint32, FileHandle uintptr, ExtendedParameters uintptr, ExtendedParameterCount uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateTimer(TimerHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TimerType uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, Uow *GUID, TmHandle uintptr, CreateOptions uint32, IsolationLevel uint32, IsolationFlags uint32, Timeout *int64, Description *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwCreateTransactionManager(TmHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, LogFileName *UNICODE_STRING, CreateOptions uint32, CommitStrength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDeleteFile(ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDeleteKey(KeyHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDeleteObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, GenerateOnClose byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDeleteValueKey(KeyHandle uintptr, ValueName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDeviceIoControlFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, IoControlCode uint32, InputBuffer uintptr, InputBufferLength uint32, OutputBuffer uintptr, OutputBufferLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDisplayString(String *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDuplicateObject(SourceProcessHandle uintptr, SourceHandle uintptr, TargetProcessHandle uintptr, TargetHandle *uintptr, DesiredAccess uint32, HandleAttributes uint32, Options uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwDuplicateToken(ExistingTokenHandle uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, EffectiveOnly byte, TokenType uintptr, NewTokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwEnumerateKey(KeyHandle uintptr, Index uint32, KeyInformationClass KEY_INFORMATION_CLASS, KeyInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwEnumerateTransactionObject(RootObjectHandle uintptr, QueryType uintptr, ObjectCursor uintptr, ObjectCursorLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwEnumerateValueKey(KeyHandle uintptr, Index uint32, KeyValueInformationClass KEY_VALUE_INFORMATION_CLASS, KeyValueInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFilterToken(ExistingTokenHandle uintptr, Flags uint32, SidsToDisable uintptr, PrivilegesToDelete uintptr, RestrictedSids uintptr, NewTokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFlushBuffersFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFlushBuffersFileEx(FileHandle uintptr, FLags uint32, Parameters uintptr, ParametersSize uint32, IoStatusBlock *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFlushKey(KeyHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFlushVirtualMemory(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, RegionSize *uintptr, IoStatus *IO_STATUS_BLOCK) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFreeVirtualMemory(ProcessHandle uintptr, BaseAddress **unsafe.Pointer, RegionSize *uintptr, FreeType uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwFsControlFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, FsControlCode uint32, InputBuffer uintptr, InputBufferLength uint32, OutputBuffer uintptr, OutputBufferLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwGetNotificationResourceManager(ResourceManagerHandle uintptr, TransactionNotification uintptr, NotificationLength uint32, Timeout *int64, ReturnLength *uint32, Asynchronous uint32, AsynchronousContext uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwImpersonateAnonymousToken(ThreadHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwLoadDriver(DriverServiceName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwLockFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, ByteOffset *int64, Length *int64, Key uint32, FailImmediately byte, ExclusiveLock byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwMakeTemporaryObject(Handle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwManagePartition(TargetHandle uintptr, SourceHandle uintptr, PartitionInformationClass PARTITION_INFORMATION_CLASS, PartitionInformation uintptr, PartitionInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwMapViewOfSection(SectionHandle uintptr, ProcessHandle uintptr, BaseAddress **unsafe.Pointer, ZeroBits uintptr, CommitSize uintptr, SectionOffset *int64, ViewSize *uintptr, InheritDisposition SECTION_INHERIT, AllocationType uint32, Win32Protect uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwNotifyChangeKey(KeyHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, CompletionFilter uint32, WatchTree byte, Buffer uintptr, BufferSize uint32, Asynchronous byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwNotifyChangeMultipleKeys(MasterKeyHandle uintptr, Count uint32, SubordinateObjects OBJECT_ATTRIBUTES, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, CompletionFilter uint32, WatchTree byte, Buffer uintptr, BufferSize uint32, Asynchronous byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenDirectoryObject(DirectoryHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenEnlistment(EnlistmentHandle *uintptr, DesiredAccess uint32, RmHandle uintptr, EnlistmentGuid *GUID, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenEvent(EventHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenFile(FileHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, IoStatusBlock *IO_STATUS_BLOCK, ShareAccess uint32, OpenOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenKey(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenKeyEx(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, OpenOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenKeyTransacted(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, TransactionHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenKeyTransactedEx(KeyHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, OpenOptions uint32, TransactionHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ObjectTypeName *UNICODE_STRING, ObjectName *UNICODE_STRING, SecurityDescriptor uintptr, ClientToken uintptr, DesiredAccess uint32, GrantedAccess uint32, Privileges *PRIVILEGE_SET, ObjectCreation byte, AccessGranted byte, GenerateOnClose *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenProcess(ProcessHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, ClientId uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenProcessToken(ProcessHandle uintptr, DesiredAccess uint32, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenProcessTokenEx(ProcessHandle uintptr, DesiredAccess uint32, HandleAttributes uint32, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenRegistryTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenResourceManager(ResourceManagerHandle *uintptr, DesiredAccess uint32, TmHandle uintptr, ResourceManagerGuid *GUID, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenSection(SectionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenSymbolicLinkObject(LinkHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenThreadToken(ThreadHandle uintptr, DesiredAccess uint32, OpenAsSelf byte, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenThreadTokenEx(ThreadHandle uintptr, DesiredAccess uint32, OpenAsSelf byte, HandleAttributes uint32, TokenHandle *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenTimer(TimerHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenTransaction(TransactionHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, Uow *GUID, TmHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwOpenTransactionManager(TmHandle *uintptr, DesiredAccess uint32, ObjectAttributes OBJECT_ATTRIBUTES, LogFileName *UNICODE_STRING, TmIdentity *GUID, OpenOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPowerInformation(InformationLevel uintptr, InputBuffer uintptr, InputBufferLength uint32, OutputBuffer uintptr, OutputBufferLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrePrepareComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrePrepareEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrepareComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrepareEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrivilegeCheck(ClientToken uintptr, RequiredPrivileges *PRIVILEGE_SET, Result *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrivilegeObjectAuditAlarm(SubsystemName *UNICODE_STRING, HandleId uintptr, ClientToken uintptr, DesiredAccess uint32, Privileges *PRIVILEGE_SET, AccessGranted byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPrivilegedServiceAuditAlarm(SubsystemName *UNICODE_STRING, ServiceName *UNICODE_STRING, ClientToken uintptr, Privileges *PRIVILEGE_SET, AccessGranted byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPropagationComplete(ResourceManagerHandle uintptr, RequestCookie uint32, BufferLength uint32, Buffer uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwPropagationFailed(ResourceManagerHandle uintptr, RequestCookie uint32, PropStatus uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryDirectoryFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS, ReturnSingleEntry byte, FileName *UNICODE_STRING, RestartScan byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryDirectoryFileEx(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS, QueryFlags uint32, FileName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryDirectoryObject(DirectoryHandle uintptr, Buffer uintptr, Length uint32, ReturnSingleEntry byte, RestartScan byte, Context *uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryEaFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ReturnSingleEntry byte, EaList uintptr, EaListLength uint32, EaIndex *uint32, RestartScan byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryFullAttributesFile(ObjectAttributes OBJECT_ATTRIBUTES, FileInformation FILE_NETWORK_OPEN_INFORMATION) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationByName(ObjectAttributes OBJECT_ATTRIBUTES, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationEnlistment(EnlistmentHandle uintptr, EnlistmentInformationClass uintptr, EnlistmentInformation uintptr, EnlistmentInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationProcess(ProcessHandle uintptr, ProcessInformationClass PROCESSINFOCLASS, ProcessInformation uintptr, ProcessInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationResourceManager(ResourceManagerHandle uintptr, ResourceManagerInformationClass uintptr, ResourceManagerInformation uintptr, ResourceManagerInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationThread(ThreadHandle uintptr, ThreadInformationClass THREADINFOCLASS, ThreadInformation uintptr, ThreadInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationToken(TokenHandle uintptr, TokenInformationClass uintptr, TokenInformation uintptr, TokenInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationTransaction(TransactionHandle uintptr, TransactionInformationClass uintptr, TransactionInformation uintptr, TransactionInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryInformationTransactionManager(TransactionManagerHandle uintptr, TransactionManagerInformationClass uintptr, TransactionManagerInformation uintptr, TransactionManagerInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryKey(KeyHandle uintptr, KeyInformationClass KEY_INFORMATION_CLASS, KeyInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryMultipleValueKey(KeyHandle uintptr, ValueEntries KEY_VALUE_ENTRY, EntryCount uint32, ValueBuffer uintptr, BufferLength *uint32, RequiredBufferLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryObject(Handle uintptr, ObjectInformationClass OBJECT_INFORMATION_CLASS, ObjectInformation uintptr, ObjectInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryQuotaInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ReturnSingleEntry byte, SidList uintptr, SidListLength uint32, StartSid uintptr, RestartScan byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQuerySecurityObject(Handle uintptr, SecurityInformation uint32, SecurityDescriptor uintptr, Length uint32, LengthNeeded *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQuerySymbolicLinkObject(LinkHandle uintptr, LinkTarget *UNICODE_STRING, ReturnedLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQuerySystemInformation(SystemInformationClass SYSTEM_INFORMATION_CLASS, SystemInformation uintptr, SystemInformationLength uint32, ReturnLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQuerySystemTime(SystemTime *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryTimerResolution(MaximumTime *uint32, MinimumTime *uint32, CurrentTime *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryValueKey(KeyHandle uintptr, ValueName *UNICODE_STRING, KeyValueInformationClass KEY_VALUE_INFORMATION_CLASS, KeyValueInformation uintptr, Length uint32, ResultLength *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryVirtualMemory(ProcessHandle uintptr, BaseAddress uintptr, MemoryInformationClass MEMORY_INFORMATION_CLASS, MemoryInformation uintptr, MemoryInformationLength uintptr, ReturnLength *uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwQueryVolumeInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FsInformation uintptr, Length uint32, FsInformationClass FS_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func ZwReadFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ByteOffset *int64, Key *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwReadOnlyEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRecoverEnlistment(EnlistmentHandle uintptr, EnlistmentKey uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRecoverResourceManager(ResourceManagerHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRecoverTransactionManager(TransactionManagerHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRegisterProtocolAddressInformation(ResourceManager uintptr, ProtocolId *GUID, ProtocolInformationSize uint32, ProtocolInformation uintptr, CreateOptions uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRenameKey(KeyHandle uintptr, NewName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRenameTransactionManager(LogFileName *UNICODE_STRING, ExistingTransactionManagerGuid *GUID) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRestoreKey(KeyHandle uintptr, FileHandle uintptr, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRollbackComplete(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRollbackEnlistment(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRollbackRegistryTransaction(TransactionHandle uintptr, Flags uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRollbackTransaction(TransactionHandle uintptr, Wait byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwRollforwardTransactionManager(TransactionManagerHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSaveKey(KeyHandle uintptr, FileHandle uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSaveKeyEx(KeyHandle uintptr, FileHandle uintptr, Format uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetEaFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetEvent(EventHandle uintptr, PreviousState *int32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationEnlistment(EnlistmentHandle uintptr, EnlistmentInformationClass uintptr, EnlistmentInformation uintptr, EnlistmentInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FileInformation uintptr, Length uint32, FileInformationClass FILE_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationKey(KeyHandle uintptr, KeySetInformationClass KEY_SET_INFORMATION_CLASS, KeySetInformation uintptr, KeySetInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationResourceManager(ResourceManagerHandle uintptr, ResourceManagerInformationClass uintptr, ResourceManagerInformation uintptr, ResourceManagerInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationThread(ThreadHandle uintptr, ThreadInformationClass THREADINFOCLASS, ThreadInformation uintptr, ThreadInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationToken(TokenHandle uintptr, TokenInformationClass uintptr, TokenInformation uintptr, TokenInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationTransaction(TransactionHandle uintptr, TransactionInformationClass uintptr, TransactionInformation uintptr, TransactionInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationTransactionManager(TmHandle uintptr, TransactionManagerInformationClass uintptr, TransactionManagerInformation uintptr, TransactionManagerInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetInformationVirtualMemory(ProcessHandle uintptr, VmInformationClass VIRTUAL_MEMORY_INFORMATION_CLASS, NumberOfEntries uintptr, VirtualAddresses MEMORY_RANGE_ENTRY, VmInformation uintptr, VmInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetQuotaInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetSecurityObject(Handle uintptr, SecurityInformation uint32, SecurityDescriptor uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetTimer(TimerHandle uintptr, DueTime *int64, TimerApcRoutine PTIMER_APC_ROUTINE, TimerContext uintptr, ResumeTimer byte, Period int32, PreviousState *byte) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetTimerEx(TimerHandle uintptr, TimerSetInformationClass TIMER_SET_INFORMATION_CLASS, TimerSetInformation uintptr, TimerSetInformationLength uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetValueKey(KeyHandle uintptr, ValueName *UNICODE_STRING, TitleIndex uint32, Type uint32, Data uintptr, DataSize uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSetVolumeInformationFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, FsInformation uintptr, Length uint32, FsInformationClass FS_INFORMATION_CLASS) uint32 { panic(`solod:extern`) }

//so:extern
func ZwSinglePhaseReject(EnlistmentHandle uintptr, TmVirtualClock *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwTerminateProcess(ProcessHandle uintptr, ExitStatus uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwUnloadDriver(DriverServiceName *UNICODE_STRING) uint32 { panic(`solod:extern`) }

//so:extern
func ZwUnlockFile(FileHandle uintptr, IoStatusBlock *IO_STATUS_BLOCK, ByteOffset *int64, Length *int64, Key uint32) uint32 { panic(`solod:extern`) }

//so:extern
func ZwUnmapViewOfSection(ProcessHandle uintptr, BaseAddress uintptr) uint32 { panic(`solod:extern`) }

//so:extern
func ZwWaitForSingleObject(Handle uintptr, Alertable byte, Timeout *int64) uint32 { panic(`solod:extern`) }

//so:extern
func ZwWriteFile(FileHandle uintptr, Event uintptr, ApcRoutine uintptr, ApcContext uintptr, IoStatusBlock *IO_STATUS_BLOCK, Buffer uintptr, Length uint32, ByteOffset *int64, Key *uint32) uint32 { panic(`solod:extern`) }

//so:extern
func VDbgPrintEx(ComponentId uint32, Level uint32, Format *byte, arglist *int8) uint32 { panic(`solod:extern`) }

//so:extern
func VDbgPrintExWithPrefix(Prefix *byte, ComponentId uint32, Level uint32, Format *byte, arglist *int8) uint32 { panic(`solod:extern`) }

