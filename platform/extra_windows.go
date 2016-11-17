package platform

import (
	"github.com/taskcluster/runlib/win32"
)

func (c *ContesterDesktop) Close() error {
	return win32.CloseDesktop(c.Desktop)
}
