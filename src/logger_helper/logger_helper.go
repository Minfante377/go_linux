package logger_helper

import(
	"fmt"
	"log"
	"os"
	"time"
)

const(
	layout string = "01-02-2006_03:04"
)

var debug_mode string = ""
var log_chan chan string = make(chan string, 100)


func write() {
	if len(log_chan) > 0{
		log.Print(<-log_chan)
	}
}


func SetLogFile(filepath string, debug string) int {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY,
							 0666)
	if err != nil {
		return -1
	}
	log.SetOutput(file)
	debug_mode = debug
	return 0
}


func LogInfo(msg string) {
	var dt string = time.Now().Format(layout)
	var log_msg string = fmt.Sprintf("[INFO - %s]: %s\n", dt, msg)
	log_chan <- log_msg
	go write()
	if debug_mode == "true" {
		fmt.Printf("[INFO - %s]: %s\n", dt, msg)
	}
}


func LogError(msg string) {
	var dt string = time.Now().Format(layout)
	var log_msg string = fmt.Sprintf("[INFO - %s]: %s\n", dt, msg)
	log_chan <- log_msg
	go write()
	if debug_mode == "true" {
		fmt.Printf("[ERROR - %s]: %s\n", dt, msg)
	}
}


func LogTestStep(step string) {
	log.Printf(fmt.Sprintf("[TEST STEP]: %s---------------\n\n", step))
	if debug_mode == "true" {
		fmt.Printf(fmt.Sprintf("[TEST STEP]: %s---------------\n\n", step))
	}
}
