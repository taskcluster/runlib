package platform

import (
	"github.com/taskcluster/runlib/win32"
)

func (c *ContesterDesktop) Display() error {
	return win32.SwitchDesktop(c.Desktop)
}

func (c *ContesterDesktop) Close() error {
	return win32.CloseDesktop(c.Desktop)
}
