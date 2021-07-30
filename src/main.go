package main

import (
	"fmt"
	"logger_helper"
	"os"
	"strconv"
	"telegram_helper"
	"time"
)

var LogDir string = "logs"
var Version string = ""
var Pass string = ""
var Debug string = ""
var Scripts string = "scripts"
var TelegramToken string = ""
var User string = ""

func init() {
	var filepath string
	os.Mkdir(LogDir, 0777)
	os.Mkdir(Scripts, 0777)
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
	var user int
	user, _ = strconv.Atoi(User)
	telegram_helper.InitQueue(TelegramToken)
	telegram_helper.InitBot(TelegramToken, user)
	for true {
		time.Sleep(30 * time.Second)
		logger_helper.LogInfo("Running...")
	}
}
