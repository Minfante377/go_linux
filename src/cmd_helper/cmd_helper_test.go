package cmd_helper

import (
	"fmt"
	"logger_helper"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

const (
	log_dir string = "test_logs/cmd_helper"
	tmp_dir string = "./tmp"
	test_file string = "test_file"
)

var pwd_user string

func init() {
	var go_root string = os.Getenv("GOPATH")
	os.Mkdir(fmt.Sprintf("%s/test_logs", go_root), 0777)
	os.Mkdir(fmt.Sprintf("%s/%s", go_root, log_dir), 0777)
	var filename string
	_, filename, _, _ = runtime.Caller(0)
	var filepath string
	filepath = fmt.Sprintf("%s/%s/cmd_helper-%s.log", go_root, log_dir,
						   time.Now().Format("01-02-2006_03-04"))
	var rc int = logger_helper.SetLogFile(filepath, "true")
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	godotenv.Load(fmt.Sprintf("%s/.secrets", go_root))
	pwd_user = os.Getenv("PWD")
	logger_helper.LogTestStep(fmt.Sprintf("Testing: %s", filename))
}


func TestExecCmd(t *testing.T) {
	es := []struct {
		input          string
		expectedOutput string
	}{
		{"ls tmp", test_file},
	}
	logger_helper.LogTestStep("Create tmp folder")
	os.Mkdir(tmp_dir, 0777)

	logger_helper.LogTestStep(fmt.Sprintf("Create %s", test_file))
	os.Create(fmt.Sprintf("%s/%s", tmp_dir, test_file))

	defer func() {
		logger_helper.LogTestStep("Remove test files")
		os.RemoveAll(tmp_dir)
	}()

	logger_helper.LogTestStep("Exec command and verify output")
	for _, c := range es {
		var err error
		var stdout, stderr string
		err, stdout, stderr = ExecCmd(c.input, &pwd_user, os.Getenv("USER"),
									  os.Getenv(os.Getenv("USER")))
		if err != nil {
			t.Errorf("Error executing cmd:\n%s",stderr)
		}
		stdout = strings.Trim(stdout, "\n")
		if stdout != c.expectedOutput {
			t.Errorf("incorrect output for `%s`: expected `%s` but got `%s`",
					 c.input, c.expectedOutput, stdout)
		}
	}
}


func TestSudoExecCmd(t *testing.T) {
	es := []struct {
		input          string
		expectedOutput string
	}{
		{"sudo ls tmp", test_file},
	}
	logger_helper.LogTestStep("Create tmp folder")
	os.Mkdir(tmp_dir, 0777)

	logger_helper.LogTestStep(fmt.Sprintf("Create %s", test_file))
	os.Create(fmt.Sprintf("%s/%s", tmp_dir, test_file))

	defer func() {
		logger_helper.LogTestStep("Remove test files")
		os.RemoveAll(tmp_dir)
	}()

	logger_helper.LogTestStep("Exec command and verify output")
	for _, c := range es {
		var res string
		res = ExecSudoCmd(c.input, os.Getenv(os.Getenv("USER")), &pwd_user,
						  os.Getenv("USER"))
		res = strings.Trim(res, "\n")
		if res != c.expectedOutput {
			t.Errorf("incorrect output for `%s`: expected `%s` but got `%s`",
					 c.input, c.expectedOutput, res)
		}
	}
}


func TestSaveScript(t *testing.T) {
	es := []struct {
		input          [2]string
		expectedOutput string
	}{
		{[2]string{fmt.Sprintf("%s/%s.sh", tmp_dir, test_file),
				   "echo \"Hello world\""},
		 "echo \"Hello world\""},
	}

	logger_helper.LogTestStep("Create tmp folder")
	os.Mkdir(tmp_dir, 0777)

	defer func() {
		logger_helper.LogTestStep("Remove test files")
		os.RemoveAll(tmp_dir)
	}()

	for _, c := range es {
		logger_helper.LogTestStep("Save test script")
		var err error
		err = SaveScript(c.input[0], c.input[1])
		if err != nil {
			t.Errorf("Test script could not be saved")
		}

		logger_helper.LogTestStep("Check script contents")
		var stdout string
		err, stdout, _ = ExecCmd(fmt.Sprintf("cat %s/%s.sh", tmp_dir,
											 test_file), &pwd_user,
											 os.Getenv("USER"),
										 	 os.Getenv(os.Getenv("USER")))
		if err != nil || c.expectedOutput != stdout {
			t.Errorf(`Contents of the script %s do not match the expected
					 ones %s`, stdout, c.expectedOutput)
		}
	}
}


func TestExecScript(t *testing.T) {
	es := []struct {
		input          [2]string
		expectedOutput string
	}{
		{[2]string{fmt.Sprintf("%s/%s.sh", tmp_dir, test_file),
				   "echo \"Hello world\""},
		 "Hello world\n"},
	}

	logger_helper.LogTestStep("Create tmp folder")
	os.Mkdir(tmp_dir, 0777)

	defer func() {
		logger_helper.LogTestStep("Remove test files")
		os.RemoveAll(tmp_dir)
	}()

	for _, c := range es {
		logger_helper.LogTestStep("Save test script")
		var err error
		err = SaveScript(c.input[0], c.input[1])

		logger_helper.LogTestStep("Exec script and check output")
		var stdout string
		err, stdout, _ = ExecScript(fmt.Sprintf("%s/%s.sh", tmp_dir,
								    test_file), &pwd_user, os.Getenv("USER"),
									os.Getenv(os.Getenv("USER")))
		if err != nil || c.expectedOutput != stdout {
			t.Errorf("Script output %s do not match the expected one %s",
					 stdout, c.expectedOutput)
		}
	}
}


func TestExecSudoScript(t *testing.T) {
	es := []struct {
		input          [2]string
		expectedOutput string
	}{
		{[2]string{fmt.Sprintf("%s/%s.sh", tmp_dir, test_file),
				   "echo \"Hello world\""},
		 "Hello world"},
	}

	logger_helper.LogTestStep("Create tmp folder")
	os.Mkdir(tmp_dir, 0777)

	defer func() {
		logger_helper.LogTestStep("Remove test files")
		os.RemoveAll(tmp_dir)
	}()

	for _, c := range es {
		logger_helper.LogTestStep("Save test script")
		var err error
		err = SaveScript(c.input[0], c.input[1])

		logger_helper.LogTestStep("Exec script and check output")
		var res string
		res = ExecSudoScript(fmt.Sprintf("%s/%s.sh", tmp_dir, test_file),
							 os.Getenv(os.Getenv("USER")), &pwd_user,
							 os.Getenv("USER"))
		if err != nil || !strings.Contains(res, c.expectedOutput) {
			t.Errorf("Script output %s do not match the expected one %s",
			         res, c.expectedOutput)
		}
	}
}
