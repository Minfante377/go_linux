package main

import(
	"cmd_helper"
	"fmt"
	"logger_helper"
	"os"
	"time"
)

var LogDir string = "logs"
var Version string = ""

func init() {
	var filepath string
	os.Mkdir(LogDir, 0777)
	filepath = fmt.Sprintf("%s/%s.log", LogDir, time.Now().String())
	var rc int = logger_helper.SetLogFile(filepath)
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	var msg string
	msg = fmt.Sprintf("Version: %s", Version)
	logger_helper.LogInfo(msg)
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
