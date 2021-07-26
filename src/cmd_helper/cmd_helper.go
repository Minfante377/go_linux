package cmd_helper

import (
	"bytes"
	"fmt"
	"logger_helper"
	"os/exec"
	"regexp"
	"time"

	"github.com/google/goexpect"
)

var (
	passRE = regexp.MustCompile("[sudo]")
	promptRE = regexp.MustCompile("%")
	timeout = 30 * time.Second
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


func ExecSudoCmd(command string, pass string) (string) {
	var result string
	e, _, err := expect.Spawn(command, -1)
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Could not spawn process %s",
							               command))
		return ""
	}
	e.Expect(passRE, timeout)
	e.Send(pass + "\n")
	result, _, err = e.Expect(promptRE, timeout)
	return result
}
