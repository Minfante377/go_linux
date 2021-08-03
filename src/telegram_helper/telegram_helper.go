package telegram_helper

import (
	"cmd_helper"
	"encoding/json"
	"fmt"
	"logger_helper"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	MAX_BOTS int = 20
	LISTEN time.Duration = 2 * time.Second
	PROCESS time.Duration = 1 * time.Second
	GET_UPDATES string = "https://api.telegram.org/bot%s/getUpdates?offset=%d"
	SEND_MSG string = "https://api.telegram.org/bot%s/sendMessage"
)

type Update struct {
	Result []Result     `json:"result"`
	Ok bool 			`json:"ok"`
}

type Result struct {
	UpdateId int `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int 	  `json:"message_id"`
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
}

type TelegramBot struct {
	token string
	user int
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
		defer r.Body.Close()
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
	for index, _ := range queue.msgs {
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


func handler(bot *TelegramBot) {
	for true {
		if bot.working < 1 {
			return
		}
		msg_id, msg := get_msg(bot.user)
		if msg_id > 0 {
			var stdout, stderr string
			var err error
			err, stdout, stderr = cmd_helper.ExecCmd(msg.Text)
			if err != nil {
				_, err = http.PostForm(
					fmt.Sprintf(SEND_MSG, bot.token),
					url.Values{
					"chat_id": {strconv.Itoa(msg.Chat.Id)},
					"text": {stderr},
					})
			}else{
				_, err = http.PostForm(
				fmt.Sprintf(SEND_MSG, bot.token),
					url.Values{
						"chat_id": {strconv.Itoa(msg.Chat.Id)},
						"text":    {stdout},
					})
			}
			if err != nil {
				logger_helper.LogError("Error sending msg")
			} else {
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
	if queue.running != 1 {
		queue.running = 1
		queue.token = token
		queue.update_id = 0
		go start_listening()
		return 0
	}
	logger_helper.LogError("Queue was already running")
	return -1
}

func InitBot(token string, user int) int {
	var new_bot *TelegramBot = spawnBot(token, user)
	if new_bot == nil {
		logger_helper.LogError(
			fmt.Sprintf("Failed to start a new bot for user %d", user))
		return -1
	}
	new_bot.working = 1
	logger_helper.LogInfo(fmt.Sprintf("Starting bot for user %d", user))
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
