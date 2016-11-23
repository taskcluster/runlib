package win32

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	procCloseDesktop            = user32.NewProc("CloseDesktop")
	procSwitchDesktop           = user32.NewProc("SwitchDesktop")
	procSetPriorityClass        = kernel32.NewProc("SetPriorityClass")
	procCreateEnvironmentBlock  = userenv.NewProc("CreateEnvironmentBlock")
	procDestroyEnvironmentBlock = userenv.NewProc("DestroyEnvironmentBlock")
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

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762270(v=vs.85).aspx
func CreateEnvironmentBlock(
	lpEnvironment *uintptr, // LPVOID
	hToken syscall.Handle, // HANDLE
	bInherit bool, // BOOL
) (err error) {
	inherit := uint32(0)
	if bInherit {
		inherit = 1
	}
	r1, _, e1 := procCreateEnvironmentBlock.Call(
		uintptr(unsafe.Pointer(&lpEnvironment)),
		uintptr(hToken),
		uintptr(inherit),
	)
	if r1 == 0 {
		err = os.NewSyscallError("CreateEnvironmentBlock", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762274(v=vs.85).aspx
func DestroyEnvironmentBlock(
	lpEnvironment *uintptr, // LPVOID
) (err error) {
	r1, _, e1 := procDestroyEnvironmentBlock.Call(
		uintptr(unsafe.Pointer(&lpEnvironment)),
	)
	if r1 == 0 {
		err = os.NewSyscallError("DestroyEnvironmentBlock", e1)
	}
	return
}

// CreateEnvironment returns an environment block, suitable for use with the
// CreateProcessAsUser system call. The default environment variables of hUser
// are overlayed with values in env.
func CreateEnvironment(env *[]string, hUser syscall.Handle) (envBlock *uint16, err error) {
	var logonEnv uintptr
	err = CreateEnvironmentBlock(&logonEnv, hUser, false)
	if err != nil {
		return
	}
	defer DestroyEnvironmentBlock(&logonEnv)
	for len := 0; len < 10000; len++ {
		x := (*uint32)(unsafe.Pointer(logonEnv + uintptr(len)))
		if *x == 0 {
			break
		}
		fmt.Printf("%v: %X\n", len, *x)
	}
	return ListToEnvironmentBlock(env), nil
}
