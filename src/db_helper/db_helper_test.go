package db_helper

import (
	"fmt"
	"io/ioutil"
	"logger_helper"
	"os"
	"runtime"
	"testing"
	"time"
)

const (
	log_dir string = "test_logs/db_helper"
	test_db string = "test_db"
)


func init() {
	var go_root string = os.Getenv("GOPATH")
	os.Mkdir(fmt.Sprintf("%s/test_logs", go_root), 0777)
	os.Mkdir(fmt.Sprintf("%s/%s", go_root, log_dir), 0777)
	var filename string
	_, filename, _, _ = runtime.Caller(0)
	var filepath string
	filepath = fmt.Sprintf("%s/%s/db_helper-%s.log", go_root, log_dir,
						   time.Now().Format("01-02-2006_03-04"))
	var rc int = logger_helper.SetLogFile(filepath, "true")
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	logger_helper.LogTestStep(fmt.Sprintf("Testing: %s", filename))
}


func TestInitDb(t *testing.T) {
	es := []struct {
		input          string
	}{
		{test_db},
	}

	logger_helper.LogTestStep(
		"Create db and verify it was created succesfully")
	for _, c := range es {
		logger_helper.LogTestStep(fmt.Sprintf("Create %s", c.input))
		res := InitDb(c.input, "users")
		defer func() {
			logger_helper.LogTestStep("Remove test files")
			os.RemoveAll(c.input)
		}()
		if res != 0 {
			t.Errorf("Error creating db %s", c.input)
		}
		_, err := ioutil.ReadDir(".")
		if err != nil {
			t.Errorf("Error creating db %s", c.input)
		}
	}
}


func TestAddDeleteUser(t *testing.T) {
	es := []struct {
		input []string
		output []string
	}{
		{[]string{test_db, "test_user", "test_id"},
		 []string{"test_user", "test_id"}},
	}

	for _, c := range es {
		logger_helper.LogTestStep("Create test db")
		InitDb(c.input[0], "users")
		defer func() {
			logger_helper.LogTestStep("Remove test files")
			os.RemoveAll(c.input[0])
		}()

		logger_helper.LogTestStep(fmt.Sprintf("Create test user %s",
											  c.input[1]))
		AddUser(c.input[0], "users", c.input[1], c.input[2])

		logger_helper.LogTestStep("Verify user was created with success")
		usernames, user_ids := GetUsers(c.input[0], "users")
		if usernames[0] != c.output[0] || user_ids[0] != c.output[1] {
			t.Errorf("User was not created: (%s, %s) != (%s, %s)", c.output[0],
					 c.output[1], usernames[0], user_ids[0])
		}

		logger_helper.LogTestStep("Delete test user")
		DeleteUser(c.input[0], "users", c.input[2])

		logger_helper.LogTestStep("Verify user was deleted with success")
		usernames, user_ids = GetUsers(c.input[0], "users")
		if len(usernames) > 0 {
			t.Errorf("User was not deleted: (%s, %s)", usernames[0],
					 user_ids[0])
		}
	}
}
