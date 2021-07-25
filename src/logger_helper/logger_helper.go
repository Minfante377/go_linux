package logger_helper

import(
	"log"
	"os"
	"time"
)

const(
	layout string = "01-02-2006_03:04"
)


func SetLogFile(filepath string) int {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return -1
	}
	log.SetOutput(file)
	return 0
}


func LogInfo(msg string) {
	var dt string = time.Now().Format(layout)
	log.Printf("[INFO - %s]: %s\n", dt, msg)
}


func LogError(msg string) {
	var dt string = time.Now().Format(layout)
	log.Printf("[ERROR - %s]: %s\n", dt, msg)
}
