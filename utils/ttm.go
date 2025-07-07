package utils

import "os/exec"

func TtmCreate(taskName string) (string, error) {
	out, err := exec.Command("ttm", "create", taskName).Output()

	return string(out), err
}

func TtmAdd(taskName string) (string, error) {
	out, err := exec.Command("ttm", "add", taskName).CombinedOutput()

	return string(out), err
}

func TtmDelete(taskName string) (string, error) {
	cmd := exec.Command("ttm", "delete", taskName)
	pp, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	pp.Write([]byte("y\n"))
	out, err := cmd.CombinedOutput()

	return string(out), err
}
