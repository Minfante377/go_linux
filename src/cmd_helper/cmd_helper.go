package cmd_helper

import (
	"bytes"
	"fmt"
	"logger_helper"
	"os/exec"
)


func ExecCmd(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	var msg string
	msg = fmt.Sprintf("Executing cmd %s", command)
	logger_helper.LogInfo(msg)
	err := cmd.Run()
	if err != nil {
		var msg string
		msg = fmt.Sprintf("Error executing cmd %s: %v", command, err)
		logger_helper.LogError(msg)
	}
	return err, stdout.String(), stderr.String()
}
