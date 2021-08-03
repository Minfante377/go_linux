package server_helper

import (
	"bytes"
	"db_helper"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger_helper"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

const (
	log_dir string = "test_logs/db_helper"
	test_db string = "test_db"
)


func init() {
	var go_root string = os.Getenv("GOPATH")
	godotenv.Load(fmt.Sprintf("%s/.secrets", go_root))
	os.Mkdir(fmt.Sprintf("%s/test_logs", go_root), 0777)
	os.Mkdir(fmt.Sprintf("%s/%s", go_root, log_dir), 0777)
	var filename string
	_, filename, _, _ = runtime.Caller(0)
	var filepath string
	filepath = fmt.Sprintf("%s/%s/server_helper-%s.log", go_root, log_dir,
						   time.Now().Format("01-02-2006_03-04"))
	var rc int = logger_helper.SetLogFile(filepath, "true")
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	logger_helper.LogTestStep(fmt.Sprintf("Testing: %s", filename))
}


func TestAddDeleteUser(t *testing.T) {
	es := []struct {
		input []string
		output []string
	}{
		{[]string{"test_user", "1234"},
		 []string{"test_user", "1234"}},
	}

	logger_helper.LogTestStep("Init test_db")
	db_helper.InitDb("test_db", "users")
	defer func() {
		logger_helper.LogTestStep("Remove test files")
		os.RemoveAll("test_db")
	}()

	for _, c := range es {
		logger_helper.LogTestStep("Init server")
		var token string
		token = os.Getenv("TELEGRAM_TOKEN")
		var port string = ":8080"
		go InitServer(port, "test_db", "users", token)
		time.Sleep(time.Second * 1)

		logger_helper.LogTestStep(fmt.Sprintf("Create test user %s",
											  c.input[0]))
		id_input, _ := strconv.Atoi(c.input[1])
		var jsonStr = []byte(fmt.Sprintf(`{"username":"%s", "user_id":%d}`,
										 c.input[0], id_input))
		var url string = fmt.Sprintf("http://localhost%s%s", port, ADD_USER)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: time.Second * 5}
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Failed to create user %s", c.input[0])
		}

		logger_helper.LogTestStep("Verify user was created with success")
		url = fmt.Sprintf("http://localhost%s%s", port, GET_USERS)
		req, err = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			t.Errorf("Failed to fetch users")
		}
		var users []User
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &users)
		id_output, _ := strconv.Atoi(c.output[1])
		if users[0].Username != c.output[0] ||
		users[0].UserId != id_output {
			t.Errorf("User was not created: (%s, %s) != (%s, %d)", c.output[0],
					 c.output[1], users[0].Username, users[0].UserId)
		}

		logger_helper.LogTestStep("Delete test user")
		url = fmt.Sprintf("http://localhost%s%s", port, DEL_USER)
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client = &http.Client{}
		res, err = client.Do(req)
		if err != nil {
			t.Errorf("Failed to delete user %s", c.input[0])
		}

		logger_helper.LogTestStep("Verify user was deleted with success")
		url = fmt.Sprintf("http://localhost%s%s", port, GET_USERS)
		req, err = http.NewRequest("GET", url, nil)
		res, err = client.Do(req)
		if err != nil {
			t.Errorf("Failed to fetch users")
		}
		body, _ = ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &users)
		if len(users) > 0 {
			t.Errorf("User was not deleted: (%s, %d)", users[0].Username,
					 users[0].UserId)
		}
	}
}
