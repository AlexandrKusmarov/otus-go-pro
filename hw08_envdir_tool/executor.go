package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	success = 0
	failure = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	com := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	// Установить переменные окружения
	com.Env = os.Environ()
	for k, v := range env {
		com.Env = append(com.Env, fmt.Sprintf("%s=%s", k, v.Value))
	}

	com.Stderr = os.Stderr
	com.Stdout = os.Stdout
	com.Stdin = os.Stdin

	// Запустить команду
	if err := com.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "error running command: %v\n", err)
		return failure
	}
	return success
}
