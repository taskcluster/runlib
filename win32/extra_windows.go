package win32

import (
	"os"
)

var (
	procCloseDesktop = user32.NewProc("CloseDesktop")
)

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
