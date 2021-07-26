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
var Pass string = ""
var Debug string = ""

func init() {
	var filepath string
	os.Mkdir(LogDir, 0777)
	filepath = fmt.Sprintf("%s/%s.log", LogDir,
						   time.Now().Format("01-02-2006_03-04"))
	var rc int = logger_helper.SetLogFile(filepath, Debug)
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	var msg string
	msg = fmt.Sprintf("Version: %s", Version)
	logger_helper.LogInfo(msg)
}


func main() {
	logger_helper.LogInfo("Starting...")
	var cmd string = "sudo ls"
	var result string
	result = cmd_helper.ExecSudoCmd(cmd, Pass)
	var msg string
	msg = fmt.Sprintf("res =\n%s", result)
	logger_helper.LogInfo(msg)
}
