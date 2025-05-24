package shell

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %s: %w", stderr.String(), err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (e *Executor) ExecuteWgCommand(args ...string) (string, error) {
	cmdArgs := append([]string{}, args...)
	return e.Execute("wg", cmdArgs...)
}

func (e *Executor) ExecuteWgQuickCommand(args ...string) (string, error) {
	cmdArgs := append([]string{}, args...)
	return e.Execute("wg-quick", cmdArgs...)
}
