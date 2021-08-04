package main

import (
	"db_helper"
	"fmt"
	"logger_helper"
	"os"
	"server_helper"
	"telegram_helper"
	"time"

	"github.com/joho/godotenv"
)

var LogDir string = "logs"
var Version string = ""
var Debug string = ""
var Scripts string = "scripts"
var TelegramToken string = ""


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
	db_helper.InitDb("users.db", "users")
	telegram_helper.InitQueue(TelegramToken)
	usernames, user_ids := db_helper.GetUsers("users.db", "users")
	var go_root string = os.Getenv("GOPATH")
	godotenv.Load(fmt.Sprintf("%s/.secrets", go_root))
	for i := range usernames {
		telegram_helper.InitBot(TelegramToken, usernames[i], user_ids[i],
							    os.Getenv(usernames[i]))
	}
	go server_helper.InitServer(":8080", "users.db", "users", TelegramToken)
	for true {
		time.Sleep(30 * time.Second)
		logger_helper.LogInfo("Running...")
	}
}
