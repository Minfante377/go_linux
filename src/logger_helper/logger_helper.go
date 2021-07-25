package logger_helper

import(
	"log"
	"os"
	"time"
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
	var dt string = time.Now().String()
	log.Printf("[INFO - %s]: %s\n", dt, msg)
}


func LogError(msg string) {
	var dt string = time.Now().String()
	log.Printf("[ERROR - %s]: %s\n", dt, msg)
}
