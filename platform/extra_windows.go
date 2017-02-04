package platform

import (
	"os"
	"runtime"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/taskcluster/runlib/win32"
)

type Desktop struct {
	Desktop     win32.Hdesk
	DesktopName string
}

func (c *Desktop) Display() error {
	return win32.SwitchDesktop(c.Desktop)
}

func (c *Desktop) Close() error {
	return win32.CloseDesktop(c.Desktop)
}

func CreateDesktop() (origDesktop, newDesktop *Desktop, err error) {
	var winsta win32.Hwinsta
	if winsta, err = win32.GetProcessWindowStation(); err != nil {
		return
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var origDesk win32.Hdesk
	if origDesk, err = win32.GetThreadDesktop(win32.GetCurrentThreadId()); err != nil {
		return
	}

	var origDeskName string
	if origDeskName, err = win32.GetUserObjectName(syscall.Handle(origDesk)); err != nil {
		return
	}

	origDesktop = &Desktop{
		Desktop:     origDesk,
		DesktopName: origDeskName,
	}

	var winstaName string
	var newDesk win32.Hdesk
	var newDeskName string
	if winstaName, err = win32.GetUserObjectName(syscall.Handle(winsta)); err == nil {
		shortName := threadIdName("c")

		newDesk, err = win32.CreateDesktop(
			syscall.StringToUTF16Ptr(shortName),
			nil, 0, 0, syscall.GENERIC_ALL, win32.MakeInheritSa())

		if err == nil {
			newDeskName = winstaName + "\\" + shortName
		}
	}

	if err != nil {
		return
	}

	everyone, err := syscall.StringToSid("S-1-1-0")
	if err == nil {
		if err = win32.AddAceToWindowStation(winsta, everyone); err != nil {
			log.Error(err)
		}
		if err = win32.AddAceToDesktop(newDesk, everyone); err != nil {
			log.Error(err)
		}
	} else {
		err = os.NewSyscallError("StringToSid", err)
	}

	newDesktop = &Desktop{
		Desktop:     newDesk,
		DesktopName: newDeskName,
	}

	return
}
