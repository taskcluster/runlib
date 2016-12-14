package win32

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

var (
	shell32 = NewLazyDLL("shell32.dll")

	procCloseDesktop            = user32.NewProc("CloseDesktop")
	procSwitchDesktop           = user32.NewProc("SwitchDesktop")
	procSetPriorityClass        = kernel32.NewProc("SetPriorityClass")
	procCreateEnvironmentBlock  = userenv.NewProc("CreateEnvironmentBlock")
	procDestroyEnvironmentBlock = userenv.NewProc("DestroyEnvironmentBlock")
	procSHSetKnownFolderPath    = shell32.NewProc("SHSetKnownFolderPath")
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
	lpEnvironment *uintptr, // LPVOID*
	hToken syscall.Handle, // HANDLE
	bInherit bool, // BOOL
) (err error) {
	inherit := uint32(0)
	if bInherit {
		inherit = 1
	}
	r1, _, e1 := procCreateEnvironmentBlock.Call(
		uintptr(unsafe.Pointer(lpEnvironment)),
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
	lpEnvironment uintptr, // LPVOID - beware - unlike LPVOID* in CreateEnvironmentBlock!
) (err error) {
	r1, _, e1 := procDestroyEnvironmentBlock.Call(
		uintptr(unsafe.Pointer(lpEnvironment)),
	)
	if r1 == 0 {
		err = os.NewSyscallError("DestroyEnvironmentBlock", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762249(v=vs.85).aspx
func SHSetKnownFolderPath(
	rfid *syscall.GUID, // REFKNOWNFOLDERID
	dwFlags uint32, // DWORD
	hToken syscall.Handle, // HANDLE
	pszPath *uint16, // PCWSTR
) (err error) {
	r1, _, _ := procSHSetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(pszPath)),
	)
	if r1 != 0 {
		err = syscall.Errno(r1)
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
	defer DestroyEnvironmentBlock(logonEnv)
	var varStartOffset uint
	envList := &[]string{}
	for {
		envVar := syscall.UTF16ToString((*[1 << 15]uint16)(unsafe.Pointer(logonEnv + uintptr(varStartOffset)))[:])
		if envVar == "" {
			break
		}
		*envList = append(*envList, envVar)
		// in UTF16, each rune takes two bytes, as does the trailing uint16(0)
		varStartOffset += uint(2 * (utf8.RuneCountInString(envVar) + 1))
	}
	env, err = MergeEnvLists(envList, env)
	if err != nil {
		return
	}
	return ListToEnvironmentBlock(env), nil
}

func MergeEnvLists(envLists ...*[]string) (*[]string, error) {
	mergedEnv := &[]string{}
	mergedEnvMap := map[string]string{}
	for _, envList := range envLists {
		if envList == nil {
			continue
		}
		for _, envSetting := range *envList {
			if utf8.RuneCountInString(envSetting) > 32767 {
				return nil, fmt.Errorf("Env setting is more than 32767 runes: %v", envSetting)
			}
			spl := strings.SplitN(envSetting, "=", 2)
			if len(spl) != 2 {
				return nil, fmt.Errorf("Could not interpret string %q as `key=value`", envSetting)
			}
			mergedEnvMap[spl[0]] = spl[1]
		}
	}
	for k, v := range mergedEnvMap {
		*mergedEnv = append(*mergedEnv, k+"="+v)
	}
	return mergedEnv, nil
}
