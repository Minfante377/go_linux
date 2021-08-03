package cmd_helper

import (
	"bytes"
	"fmt"
	"logger_helper"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/google/goexpect"
)

var (
	passRE = regexp.MustCompile("[sudo]")
	promptRE = regexp.MustCompile("%")
	timeout = 30 * time.Second
)


func ExecCmd(command string, pwd_env *string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if strings.Contains(command, "cd") {
		command = fmt.Sprintf("%s;echo -n $PWD", command)
	}
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if pwd_env != nil {
		cmd.Dir = *pwd_env
	}
	var msg string
	msg = fmt.Sprintf("Executing cmd %s", command)
	logger_helper.LogInfo(msg)
	err := cmd.Run()
	if err != nil {
		var msg string
		msg = fmt.Sprintf("Error executing cmd %s: %v", command, err)
		logger_helper.LogError(msg)
		return err, stdout.String(), stderr.String()
	}
	if pwd_env != nil && strings.Contains(command, "cd") {
		*pwd_env = stdout.String()
	}

	return err, stdout.String(), stderr.String()
}


func ExecSudoCmd(command string, pass string, pwd_env *string) (string) {
	var result string
	logger_helper.LogInfo(fmt.Sprintf("Executing sudo cmd: %s", command))
	if pwd_env != nil{
		command = fmt.Sprintf("cd %s; %s", *pwd_env, command)
	}
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


func SaveScript(path string, data string) error {
	logger_helper.LogInfo(fmt.Sprintf("Saving script %s", path))
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		logger_helper.LogError("Error saving script")
		return err
	}

	_, err = file.Write([]byte(data))

	if err != nil {
		logger_helper.LogError("Error writing contents to the file")
		file.Close()
		return err
	}

	file.Close()
	return nil
}


func ExecScript(path string, pwd_env *string) (error, string, string) {
	var err error
	var stdout, stderr, cmd string
	cmd = fmt.Sprintf("/bin/bash %s", path)
	err, stdout, stderr = ExecCmd(cmd, pwd_env)
	return err, stdout, stderr
}


func ExecSudoScript(path string, pass string, pwd_env *string) string {
	var cmd, result string
	cmd = fmt.Sprintf("sudo /bin/bash %s", path)
	result = ExecSudoCmd(cmd, pass, pwd_env)
	logger_helper.LogInfo(fmt.Sprintf("Result is: %s", result))
	return result
}
