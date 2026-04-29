//so:driver
package main

import (
	"unicode/utf16"

	"solod.dev/so/wdk"
)

const SE_DEBUG_PRIVILEGE = 20

type Driver struct {
	device      *wdk.DEVICE_OBJECT
	handleInUse BOOLEAN
	initialized bool
}

func DriverEntry(driverObj *wdk.DRIVER_OBJECT, registryPath *wdk.UNICODE_STRING) wdk.NTSTATUS {
	wdk.UNREFERENCED_PARAMETER(registryPath)
	wdk.ExInitializeDriverRuntime(wdk.DrvRtPoolNxOptIn)

	gDrv = &Driver{}

	numCpus := wdk.KeQueryActiveProcessorCount(nil)
	if err := InitializeGlobals(numCpus); err != nil {
		LogError("Failed to initialize globals: %v", err)
		return wdk.STATUS_UNSUCCESSFUL
	}

	devName := utf16.Encode([]rune("\\Device\\HyperDbgDebuggerDevice"))
	dosName := utf16.Encode([]rune("\\DosDevices\\HyperDbgDebuggerDevice"))
	var name, dos wdk.UNICODE_STRING
	wdk.RtlInitUnicodeString(&name, &devName[0])
	wdk.RtlInitUnicodeString(&dos, &dosName[0])

	status := wdk.IoCreateDevice(driverObj, 0, &name,
		wdk.FILE_DEVICE_UNKNOWN, wdk.FILE_DEVICE_SECURE_OPEN, wdk.FALSE, &gDrv.device)

	if status == wdk.STATUS_SUCCESS {
		for i := range wdk.IRP_MJ_MAXIMUM_FUNCTION {
			driverObj.MajorFunction[i] = drvUnsupported
		}
		driverObj.MajorFunction[wdk.IRP_MJ_CLOSE] = drvClose
		driverObj.MajorFunction[wdk.IRP_MJ_CREATE] = drvCreate
		driverObj.MajorFunction[wdk.IRP_MJ_READ] = drvRead
		driverObj.MajorFunction[wdk.IRP_MJ_WRITE] = drvWrite
		driverObj.MajorFunction[wdk.IRP_MJ_DEVICE_CONTROL] = drvIoctl
		driverObj.DriverUnload = drvUnload
		wdk.IoCreateSymbolicLink(&dos, &name)

		gDrv.device.Flags |= wdk.DO_BUFFERED_IO
		gDrv.initialized = true

		LogInfo("HyperDbg driver loaded successfully")
	}

	return status
}

func drvUnsupported(_ *wdk.DEVICE_OBJECT, irp *wdk.IRP) wdk.NTSTATUS {
	irp.IoStatus.Status = wdk.STATUS_SUCCESS
	irp.IoStatus.Information = 0
	wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
	return wdk.STATUS_SUCCESS
}

func drvRead(_ *wdk.DEVICE_OBJECT, irp *wdk.IRP) wdk.NTSTATUS {
	stack := wdk.IoGetCurrentIrpStackLocation(irp)
	outLen := stack.Parameters.DeviceIoControl.OutputBufferLength

	if outLen == 0 {
		irp.IoStatus.Status = wdk.STATUS_SUCCESS
		irp.IoStatus.Information = 0
		wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
		return wdk.STATUS_SUCCESS
	}

	buffer := irp.AssociatedIrp.SystemBuffer
	count := gLog.FlushToUser(uintptr(PVOID(buffer)), uintptr(outLen))

	irp.IoStatus.Status = wdk.STATUS_SUCCESS
	irp.IoStatus.Information = uint64(uintptr(count) * unsafeSizeofLogMessage())
	wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
	return wdk.STATUS_SUCCESS
}

func drvWrite(_ *wdk.DEVICE_OBJECT, irp *wdk.IRP) wdk.NTSTATUS {
	irp.IoStatus.Status = wdk.STATUS_SUCCESS
	irp.IoStatus.Information = 0
	wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
	return wdk.STATUS_SUCCESS
}

func drvClose(_ *wdk.DEVICE_OBJECT, irp *wdk.IRP) wdk.NTSTATUS {
	gDrv.handleInUse = false
	irp.IoStatus.Status = wdk.STATUS_SUCCESS
	irp.IoStatus.Information = 0
	wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
	return wdk.STATUS_SUCCESS
}

func drvCreate(_ *wdk.DEVICE_OBJECT, irp *wdk.IRP) wdk.NTSTATUS {
	priv := wdk.LUID{LowPart: SE_DEBUG_PRIVILEGE, HighPart: 0}
	if wdk.SeSinglePrivilegeCheck(priv, int8(irp.RequestorMode)) == wdk.FALSE {
		irp.IoStatus.Status = wdk.STATUS_ACCESS_DENIED
		wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
		return wdk.STATUS_ACCESS_DENIED
	}

	if gDrv.handleInUse {
		irp.IoStatus.Status = wdk.STATUS_SUCCESS
		irp.IoStatus.Information = 0
		wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
		return wdk.STATUS_SUCCESS
	}

	gDrv.handleInUse = true
	irp.IoStatus.Status = wdk.STATUS_SUCCESS
	irp.IoStatus.Information = 0
	wdk.IoCompleteRequest(irp, wdk.IO_NO_INCREMENT)
	return wdk.STATUS_SUCCESS
}

func drvUnload(_ *wdk.DRIVER_OBJECT) {
	LogInfo("HyperDbg driver unloading...")

	CleanupGlobals()

	if gDrv != nil && gDrv.device != nil {
		devName := utf16.Encode([]rune("\\DosDevices\\HyperDbgDebuggerDevice"))
		var dos wdk.UNICODE_STRING
		wdk.RtlInitUnicodeString(&dos, &devName[0])
		wdk.IoDeleteSymbolicLink(&dos)
		wdk.IoDeleteDevice(gDrv.device)
	}

	LogInfo("HyperDbg driver unloaded successfully")
}

func unsafeSizeofLogMessage() uintptr { return 288 }
