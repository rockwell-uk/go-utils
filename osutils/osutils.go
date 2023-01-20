package osutils

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// unixOS is the set of GOOS values matched by the "unix" build tag.
var unixOS = map[string]bool{
	"aix":       true,
	"android":   true,
	"darwin":    true,
	"dragonfly": true,
	"freebsd":   true,
	"hurd":      true,
	"illumos":   true,
	"ios":       true,
	"linux":     true,
	"netbsd":    true,
	"openbsd":   true,
	"solaris":   true,
}

func GetOS() string {
	return runtime.GOOS
}

func GetArch() string {
	return runtime.GOARCH
}

func OSIsUnix() bool {
	os := GetOS()

	if _, ok := unixOS[os]; ok {
		return true
	}

	return false
}

func RunCommand(command string, args ...string) (string, string, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func RunCommandSilent(command string, args ...string) (string, string, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func RunShellCommand(shell string, command string) (string, string, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(shell, "-c", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func RunShellCommandSilent(shell string, command string) (string, string, error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(shell, "-c", command)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func CommandExists(cmd string) bool {

	_, err := exec.LookPath(cmd)

	return err == nil
}

func GetULimit() (int, error) {

	if !OSIsUnix() {
		// @TODO figure out what to do on non unix
		return 0, nil
	}

	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return 0, err
	}

	return int(rLimit.Cur), nil
}
