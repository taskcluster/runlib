package win32

import (
	"os"
	"syscall"
)

var (
	procCloseDesktop     = user32.NewProc("CloseDesktop")
	procSwitchDesktop    = user32.NewProc("SwitchDesktop")
	procSetPriorityClass = kernel32.NewProc("SetPriorityClass")
)

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686219(v=vs.85).aspx
	ABOVE_NORMAL_PRIORITY_CLASS   = 0x00008000
	BELOW_NORMAL_PRIORITY_CLASS   = 0x00004000
	HIGH_PRIORITY_CLASS           = 0x00000080
	IDLE_PRIORITY_CLASS           = 0x00000040
	NORMAL_PRIORITY_CLASS         = 0x00000020
	PROCESS_MODE_BACKGROUND_BEGIN = 0x00100000
	PROCESS_MODE_BACKGROUND_END   = 0x00200000
	REALTIME_PRIORITY_CLASS       = 0x00000100
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686347(v=vs.85).aspx
func SwitchDesktop(
	hDesktop Hdesk, // HDESK
) (err error) {
	r1, _, e1 := procSwitchDesktop.Call(
		uintptr(hDesktop),
	)
	if r1 == 0 {
		err = os.NewSyscallError("SwitchDesktop", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682024(v=vs.85).aspx
func CloseDesktop(
	hDesktop Hdesk, // HDESK
) (err error) {
	r1, _, e1 := procCloseDesktop.Call(
		uintptr(hDesktop),
	)
	if r1 == 0 {
		err = os.NewSyscallError("CloseDesktop", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686219(v=vs.85).aspx
func SetPriorityClass(
	hProcess syscall.Handle, // HANDLE
	dwPriorityClass uint32, // DWORD
) (err error) {
	r1, _, e1 := procSetPriorityClass.Call(
		uintptr(hProcess),
		uintptr(dwPriorityClass),
	)
	if r1 == 0 {
		err = os.NewSyscallError("SetPriorityClass", e1)
	}
	return
}
