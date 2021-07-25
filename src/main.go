package main

import(
	"cmd_helper"
	"fmt"
	"logger_helper"
	"os"
	"time"
)

var log_dir string = "logs"


func init() {
	var filepath string
	os.Mkdir(log_dir, 0777)
	filepath = fmt.Sprintf("%s/%s.log", log_dir, time.Now().String())
	var rc int = logger_helper.SetLogFile(filepath)
	if rc != 0 {
		panic("Could not set logger log file!")
	}
}


func main() {
	logger_helper.LogInfo("Starting...")
	var cmd string = "ls"
	var err error
	var stdout, stderr string
	err, stdout, stderr = cmd_helper.ExecCmd(cmd)
	if err == nil{
		var msg string
		msg = fmt.Sprintf("stdout = %s\n stderr = %s", stdout, stderr)
		logger_helper.LogInfo(msg)
	}
}
