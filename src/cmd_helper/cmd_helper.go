package cmd_helper

import (
	"bytes"
	"errors"
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
	timeout int64 = 30
)


func execCmdTimeout(cmd *exec.Cmd, timeout int64) error {
	var err chan error = make(chan error, 1)
	var ended chan bool = make(chan bool, 1)

	go func(err chan error, ended chan bool, cmd *exec.Cmd) {
		var cmd_err error
		cmd_err = cmd.Run()
		err <- cmd_err
		ended <- true
	}(err, ended,  cmd)

	end_time := time.Now().Unix() + timeout
	for time.Now().Unix() < end_time {
		if len(ended) != 0 {
			if  <-ended {
				return  <-err
			}
		}
		time.Sleep(time.Second * 1)
	}
	return errors.New("Command timeout")
}


func ExecCmd(command string, pwd_env *string,
			 user string, pass string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var aux string
	aux = command
	if pwd_env != nil{
		command = fmt.Sprintf(`echo %s | su - %s -c "cd %s;%s"`, pass,
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
	err := execCmdTimeout(cmd, timeout)
	if err != nil {
		var msg string
		msg = fmt.Sprintf("Error executing cmd %s: %v", command, err)
		logger_helper.LogError(msg)
		return err, stdout.String(), fmt.Sprint(stderr.String(), "\n",
												err.Error())
	}
	if pwd_env != nil && strings.Contains(aux, "cd") {
		aux_cmd := exec.Command("bash", "-c",
							    fmt.Sprintf(
									`echo %s | su - %s -c "cd %s;%s;pwd"`,
									pass, user, *pwd_env, aux) )
		var path bytes.Buffer
		aux_cmd.Stdout = &path
		aux_cmd.Run()
		*pwd_env = strings.Replace(path.String(), "\n", "", 1)
	}
	return err, stdout.String(), stderr.String()
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
