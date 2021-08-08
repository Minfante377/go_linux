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
)

var (
	passRE = regexp.MustCompile("[sudo]")
	promptRE = regexp.MustCompile("%")
	timeout = 30 * time.Second
)


func ExecCmd(command string, pwd_env *string,
			 user string, pass string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var aux string
	aux = command
	if pwd_env != nil{
		command = fmt.Sprintf(`echo %s | su - %s -c "cd %s;%s;pwd"`, pass,
							  user, *pwd_env, command)
	} else {
		command = fmt.Sprintf(`echo %s | su - %s -c "%s"`, pass, user, command)
	}
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
		return err, stdout.String(), stderr.String()
	}
	var output []string = strings.Split(stdout.String(), "\n")
	if pwd_env != nil && strings.Contains(aux, "cd") {
		*pwd_env = output[0]
	}
	var res string = strings.Join(output[:len(output) - 2], "\n")
	return err, res, stderr.String()
}


func ExecSudoCmd(command string, pass string, pwd_env *string,
			     user string) (string) {
	logger_helper.LogInfo(fmt.Sprintf("Executing sudo cmd: %s", command))
	command = strings.Replace(command, "sudo", "", 1)
	command = fmt.Sprintf(`echo %s |sudo -S %s`, pass, command)
	err, stdout, stderr := ExecCmd(command, pwd_env, user, pass)
	if err != nil {
		logger_helper.LogError("Failed to execute sudo command")
		return stderr
	}
	logger_helper.LogInfo("Successfully executed sudo command")
	return stdout
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


func ExecScript(path string, pwd_env *string, user string,
				pass string) (error, string, string) {
	var err error
	var stdout, stderr, cmd string
	cmd = fmt.Sprintf("/bin/bash %s", path)
	err, stdout, stderr = ExecCmd(cmd, pwd_env, user, pass)
	return err, stdout, stderr
}


func ExecSudoScript(path string, pass string, pwd_env *string,
					user string) string {
	var cmd, result string
	cmd = fmt.Sprintf("sudo /bin/bash %s", path)
	result = ExecSudoCmd(cmd, pass, pwd_env, user)
	logger_helper.LogInfo(fmt.Sprintf("Result is: %s", result))
	return result
}
