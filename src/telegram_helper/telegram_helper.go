package telegram_helper

import (
	"cmd_helper"
	"encoding/json"
	"fmt"
	"logger_helper"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MAX_BOTS int = 20
	LISTEN time.Duration = 2 * time.Second
	PROCESS time.Duration = 1 * time.Second
	GET_UPDATES string = "https://api.telegram.org/bot%s/getUpdates?offset=%d"
	SEND_MSG string = "https://api.telegram.org/bot%s/sendMessage"
	HELP string = "/help"
	HELP_MSG string = "/list List available scripts.\n"+
					  "/script <script_name> Execute <script_name>.\n"
	EXEC_SCRIPT string = "/script"
	LIST string = "/list"
)

type Update struct {
	Result []Result     `json:"result"`
	Ok bool			`json:"ok"`
}

type Result struct {
	UpdateId int `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int	  `json:"message_id"`
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
}

type TelegramBot struct {
	token string
	user int
	user_name string
	user_pwd string
	user_pass string
	working int
}

type Queue struct {
	msgs []Message
	token string
	running int
	update_id int
	mu sync.Mutex
}

var bot_count int = 0
var bots [MAX_BOTS]*TelegramBot
var queue Queue
var netClient = &http.Client{
	  Timeout: time.Second * 10,
}

func setParams(token string, user int, bot *TelegramBot) {
	bot.token = token
	bot.user = user
}


func spawnBot(token string, user int) *TelegramBot {
	logger_helper.LogInfo("Creating a new bot")
	if bot_count == MAX_BOTS {
		logger_helper.LogError("Max bots reached")
		return nil
	}
	var bot TelegramBot
	bot.working = 0
	setParams(token, user, &bot)
	bots[bot_count] = &bot
	bot_count += 1
	return bots[bot_count - 1]
}


func parseResponse(r *http.Response, update *Update) (error) {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(update); err != nil {
		logger_helper.LogError(
			fmt.Sprintf("could not decode incoming update %s", err.Error()))
		return err
	}
	return nil
}


func start_listening() {
	logger_helper.LogInfo("Start listening...")
	for true {
		if queue.running == 0 {
			return
		}
		update := Update{}
		var err error
		r, err := netClient.Get(fmt.Sprintf(GET_UPDATES, queue.token,
								queue.update_id))

		err = parseResponse(r, &update)
		if err == nil && update.Ok {
			logger_helper.LogInfo("New messages available")
			queue.mu.Lock()
			for _, result := range update.Result {
				queue.msgs = append(queue.msgs, result.Message)
				if queue.update_id <= result.UpdateId {
					queue.update_id = result.UpdateId + 1
				}
			}
			queue.mu.Unlock()
		}
		time.Sleep(LISTEN)
	}
}


func get_msg(user int) (int, *Message) {
	queue.mu.Lock()
	for index := range queue.msgs {
		if queue.msgs[index].Chat.Id == user {
			queue.mu.Unlock()
			return queue.msgs[index].MessageId, &queue.msgs[index]
		}
	}
	queue.mu.Unlock()
	return -1, nil
}


func delete_msg(message_id int) {
	queue.mu.Lock()
	for index := range queue.msgs {
		if queue.msgs[index].MessageId == message_id {
			if index == len(queue.msgs){
				queue.msgs = queue.msgs[:len(queue.msgs)-1]
				queue.mu.Unlock()
				return
			} else {
				queue.msgs = append(queue.msgs[:index],
									queue.msgs[index+1:]...)
				queue.mu.Unlock()
				return
			}
		}
	}
	queue.mu.Unlock()
}


func sendMsg(chat_id int, text string, token string) error {
	var err error
	logger_helper.LogInfo(fmt.Sprintf("Sending msg to %d", chat_id))
	_, err = http.PostForm(
		fmt.Sprintf(SEND_MSG, token),
		url.Values{
		"chat_id": {strconv.Itoa(chat_id)},
		"text": {text},
		})
	if err != nil {
		logger_helper.LogError(fmt.Sprintf("Error sending msg: %s",
										   err.Error()))
	}
	return err
}


func sendHelp(token string, chat_id int) error {
	sendMsg(chat_id, HELP_MSG, token)
	return nil
}


func listScripts(token string, chat_id int, user string, pass string) error {
	var go_root string = os.Getenv("GOPATH")
	var cmd string = fmt.Sprintf("ls %s/scripts", go_root)
	var stdout string
	_, stdout, _ = cmd_helper.ExecCmd(cmd, nil, user, pass)
	sendMsg(chat_id, fmt.Sprintf("Available scripts:\n%s", stdout), token)
	return nil
}


func execCmd(msg string, token string, chat_id int, pass string,
	user_pwd *string, user_name string) error {
	var err error
	if strings.Contains(msg, "sudo") {
		var result string
		result = cmd_helper.ExecSudoCmd(msg, pass, user_pwd, user_name)
		err = sendMsg(chat_id, result, token)
	} else {
		var stdout, stderr string
		err, stdout, stderr = cmd_helper.ExecCmd(msg, user_pwd, user_name,
												 pass)
		if err != nil {
			err = sendMsg(chat_id, stderr, token)
		} else {
			err = sendMsg(chat_id, stdout, token)
		}
	}
	return err
}

func execScript(msg string, token string, user_id int, pass string,
				user_name string) error {
	var err error
	var command, go_root string
	var args []string
	go_root = os.Getenv("GOPATH")
	msg = strings.Replace(msg, "\n", "", 1)
	args = strings.Split(msg, " ")
	if strings.Contains(msg, "sudo") {
		var result string
		command = fmt.Sprintf("%s/scripts/%s %d %s %s", go_root, args[2],
							  user_id, token, args[3])
		result = cmd_helper.ExecSudoScript(command, pass, nil, user_name)
		err = sendMsg(user_id, result, token)
	} else {
		var stdout, stderr string
		command = fmt.Sprintf("%s/scripts/%s %d %s %s", go_root, args[1],
							  user_id, token, args[2])
		err, stdout, stderr = cmd_helper.ExecScript(command, nil, user_name,
													pass)
		if err != nil {
			err = sendMsg(user_id, stderr, token)
		} else {
			err = sendMsg(user_id, stdout, token)
		}
	}
	return err
}


func handler(bot *TelegramBot) {
	for true {
		if bot.working < 1 {
			return
		}
		msg_id, msg := get_msg(bot.user)
		if msg_id > 0 {
			var err error
			if strings.Contains(msg.Text, HELP){
				err = sendHelp(bot.token, msg.Chat.Id)
			}else if strings.Contains(msg.Text, LIST){
				err = listScripts(bot.token, msg.Chat.Id, bot.user_name,
								  bot.user_pass)
			}else if strings.Contains(msg.Text, EXEC_SCRIPT) {
				err = execScript(msg.Text, bot.token, msg.Chat.Id,
								 bot.user_pass, bot.user_name)
			}else{
				err = execCmd(msg.Text, bot.token, msg.Chat.Id,
							  bot.user_pass, &bot.user_pwd, bot.user_name)
			}
			if err == nil {
				delete_msg(msg_id)
			}
		}
	}
	time.Sleep(PROCESS)
}


func getIndexByUser(user int) int {
	var i int = 0
	for i, _ =  range bots {
		if bots[i].user == user {
			return i
		}
	}
	return -1
}

func InitQueue(token string) int {
	queue.mu.Lock()
	if queue.running != 1 {
		queue.running = 1
		queue.token = token
		queue.update_id = 0
		go start_listening()
		queue.mu.Unlock()
		return 0
	}
	logger_helper.LogError("Queue was already running")
	queue.mu.Unlock()
	return -1
}

func InitBot(token string, user_name string, user_id int,
			 user_pass string) int {
	var new_bot *TelegramBot = spawnBot(token, user_id)
	if new_bot == nil {
		logger_helper.LogError(
			fmt.Sprintf("Failed to start a new bot for user %d", user_id))
		return -1
	}
	new_bot.working = 1
	logger_helper.LogInfo(fmt.Sprintf("Starting bot for user %d", user_id))
	new_bot.user_pwd = fmt.Sprintf("/home/%s", user_name)
	new_bot.user_pass = user_pass
	new_bot.user_name = user_name
	go handler(new_bot)
	return 0
}

func DeleteBot(user int) int {
	var index int = getIndexByUser(user)
	if index < 0 {
		logger_helper.LogError(fmt.Sprintf("Bot for user %d not found", user))
		return -1
	}
	bots[index].working = 0
	if index < bot_count - 1 {
		bots[index] = bots[bot_count - 1]
	}
	bot_count -= 1
	return 0
}
