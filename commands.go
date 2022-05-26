package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	command struct {
		name    string
		handler commandHandler
	}
	commandHandler func(update tgbotapi.Update, bot *tgbotapi.BotAPI) string
)

var commands = [...]command{
	{"ping", ping},
	{"chatid", chatid},
	{"syrok", syrok},
}

func execCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {

	for _, cmd := range commands {
		if cmd.name == update.Message.Command() {
			return cmd.handler(update, bot)
		}
	}
	return ""
}

func ping(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {
	return "still alive"
}

func chatid(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {

	return fmt.Sprintf("%d", update.Message.Chat.ID)
}

func syrok(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {
	var compressFile bool
	var file string
	if update.Message.Text == "" {
		return "kek"
	}
	reply := update.Message.ReplyToMessage
	if reply == nil || (reply.Photo == nil && reply.Document == nil) {
		return "lol"
	}
	compressFile = reply.Document != nil

	if !compressFile {
		file = reply.Photo[len(reply.Photo)-1].FileID
	} else {
		file = reply.Document.FileID
	}

	link, err := bot.GetFileDirectURL(file)
	if err != nil {
		return err.Error()
	}

	resp, err := http.Get(link)
	if err != nil {
		return err.Error()
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	mode := 0
	mode, _ = strconv.Atoi(update.Message.CommandArguments())

	data, err := syrokImage(bytes, mode)
	if err != nil {
		return err.Error()
	}
	resp.Body.Close()

	if !compressFile {
		photo := tgbotapi.NewPhoto(update.FromChat().ChatConfig().ChatID,
			tgbotapi.FileBytes{
				Bytes: data,
				Name:  "resultlmao.png"})
		photo.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
		bot.Send(photo)
	} else {
		doc := tgbotapi.NewDocument(update.FromChat().ChatConfig().ChatID,
			tgbotapi.FileBytes{
				Bytes: data,
				Name:  "resultlmao.png"})
		doc.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
		bot.Send(doc)
	}

	return ""
}
