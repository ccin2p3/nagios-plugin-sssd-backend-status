package nagsssdbackend

import "os/exec"

var execFnc = func(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
