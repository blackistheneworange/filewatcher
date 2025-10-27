//go:build windows
// +build windows

package process

import (
	"errors"
	"golang.org/x/sys/windows"
	"fmt"
	"unsafe"
)

const PROCESS_STILL_ACTIVE_CODE = 259

type WindowsProcess struct {
	Process
	execCmd		string
	jobHandle	windows.Handle
	startupInfo	windows.StartupInfo
	processInfo	windows.ProcessInformation
}

func _CreateProcess(execCmd string) *WindowsProcess {
	process := &WindowsProcess{ execCmd: execCmd }
	return process
}

func (process *WindowsProcess) StartProcess() error {
	return process.resumeProcess()
}

func (process *WindowsProcess) StopProcess() error {
	// no action if process is not running
	if isRunning, _ := process.isRunning(); !isRunning {
		return nil
	}

	if err := windows.TerminateJobObject(process.jobHandle, 1); err != nil {
		return err
	}
	process.closeHandles()

	return nil
}

// resumeProcess resumes the primary thread of the process which is created in suspended state
func (process *WindowsProcess) resumeProcess() error {
	if isRunning, isRunningErr := process.isRunning(); isRunning || isRunningErr != nil {
		if isRunningErr != nil {
			return isRunningErr
		}
		process.StopProcess()
		*process = *_CreateProcess(process.execCmd)
	}
	
	if err := process.createSuspendedProcess(); err != nil {
		return errors.New(fmt.Sprintf("windows process error: %s", err.Error()))
	}

	if err := process.assignProcessToJobObject(); err != nil {
		return errors.New(fmt.Sprintf("windows process error: %s", err.Error()))
	}

	if process.processInfo.Thread == 0 {
		return errors.New("windows process error: No process found to resume")
	}

	if _, err := windows.ResumeThread(process.processInfo.Thread); err != nil {
		return err
	}

	return nil
}

func (process *WindowsProcess) isRunning() (bool, error) {
    var exitCode uint32
	if process.processInfo.Process == 0 {
		return false, nil
	}
    err := windows.GetExitCodeProcess(process.processInfo.Process, &exitCode)
    if err != nil {
        return false, err
    }

    return exitCode == PROCESS_STILL_ACTIVE_CODE, nil
}

func (process *WindowsProcess) assignProcessToJobObject() error {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	// JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE - Causes all processes associated with the job to terminate when the last handle to the job is closed.
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	_, err = windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if err != nil {
		return err
	}

	process.jobHandle = job
	return windows.AssignProcessToJobObject(job, process.processInfo.Process)
}

func (process *WindowsProcess) createSuspendedProcess() error {
	if process.execCmd == "" {
		return errors.New("No executable path provided")
	}

	cmdLine := windows.StringToUTF16Ptr(process.execCmd)

	//https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createprocessa
	err := windows.CreateProcess(
		nil,						// lpApplicationName
		cmdLine,              		// lpCommandLine
        nil,                  		// lpProcessAttributes
        nil,                  		// lpThreadAttributes
        false,               		// bInheritHandles
        windows.CREATE_SUSPENDED, 	// dwCreationFlags (optional)
        nil,                  		// lpEnvironment
        nil,                  		// lpCurrentDirectory
        &process.startupInfo,       // lpStartupInfo
        &process.processInfo,       // lpProcessInformation
	)

	return err
}

func (process *WindowsProcess) closeHandles() {
	//cleanup
	if process.processInfo.Thread != 0 {
		windows.CloseHandle(process.processInfo.Thread)
	}
	if process.processInfo.Process != 0 {
		windows.CloseHandle(process.processInfo.Process)
	}
	if process.jobHandle != 0 {
		windows.CloseHandle(process.jobHandle)
	}
}