package telegram_helper

import (
	"cmd_helper"
	"encoding/json"
	"fmt"
	"logger_helper"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	MAX_BOTS int = 20
	SLEEP time.Duration = 1 * time.Second
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
	msg_id int
}

var bot_count int = 0
var bots [MAX_BOTS]*TelegramBot


func setParams(token string, user int, bot *TelegramBot) {
	bot.token = token
	bot.user = user
	bot.msg_id = -1
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


func handler(bot *TelegramBot) {
	for true {
		if bot.working < 1 {
			return
		}
		update := Update{}
		var err error
		r, _ := http.Get(fmt.Sprintf(GET_UPDATES, bot.token, bot.msg_id))
		err = parseResponse(r, &update)
		var stdout, stderr string
		if update.Ok && len(update.Result) > 0 &&
			update.Result[0].Message.Chat.Id == bot.user {
			err, stdout, stderr = cmd_helper.ExecCmd(
				update.Result[0].Message.Text)
			if err != nil {
				_, err = http.PostForm(
					fmt.Sprintf(SEND_MSG, bot.token),
					url.Values{
						"chat_id": {strconv.Itoa(
									update.Result[0].Message.Chat.Id)},
						"text":    {stderr},
					})
			}else{
				_, err = http.PostForm(
				fmt.Sprintf(SEND_MSG, bot.token),
					url.Values{
						"chat_id": {strconv.Itoa(
									update.Result[0].Message.Chat.Id)},
						"text":    {stdout},
					})
			}
			if err != nil {
				logger_helper.LogError("Error sending msg")
			} else {
				bot.msg_id = update.Result[0].UpdateId + 1
			}
		}
	}
	time.Sleep(SLEEP)
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
		logger_helper.LogError(fmt.Sprintf("Bot for user %s not found", user))
		return -1
	}
	bots[index].working = 0
	if index < bot_count - 1 {
		bots[index] = bots[bot_count - 1]
	}
	bot_count -= 1
	return 0
}
