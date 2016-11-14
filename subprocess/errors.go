package subprocess

import "github.com/taskcluster/runlib/tools"

const ERR_USER = "HANDS"

func IsUserError(err error) bool {
	return tools.HasAnnotation(err, ERR_USER)
}
