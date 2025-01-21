package cicd

import (
	"bytes"
	"os/exec"
)

func ExecuteCommand(command string, args []string, env []string) (string, string, int) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Env = append(cmd.Env, env...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	return stdout.String(), stderr.String(), cmd.ProcessState.ExitCode()
}
