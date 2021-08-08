package logger_helper

import(
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const(
	layout string = "01-02-2006_03:04"
)

var debug_mode string = ""
var log_mu sync.Mutex

func SetLogFile(filepath string, debug string) int {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return -1
	}
	log.SetOutput(file)
	debug_mode = debug
	return 0
}


func LogInfo(msg string) {
	var dt string = time.Now().Format(layout)
	log_mu.Lock()
	log.Printf("[INFO - %s]: %s\n", dt, msg)
	log_mu.Unlock()
	if debug_mode == "true" {
		fmt.Printf("[INFO - %s]: %s\n", dt, msg)
	}
}


func LogError(msg string) {
	var dt string = time.Now().Format(layout)
	log_mu.Lock()
	log.Printf("[ERROR - %s]: %s\n", dt, msg)
	log_mu.Unlock()
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
