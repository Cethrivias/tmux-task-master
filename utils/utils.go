package utils

import "os/exec"

func Run(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()

	return string(out), err
}
