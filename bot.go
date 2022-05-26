package main

import (
	"log"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pborman/getopt"
)

type botConfig struct {
	Token string
}

func loadConfig(path string) botConfig {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	var config botConfig
	_, err = toml.Decode(string(b), &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {

	dir := getopt.StringLong("dir", 'd', "./", "working directory")

	getopt.Parse()

	config := loadConfig(*dir + "config.toml")

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if update.Message.IsCommand() {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = execCommand(update, bot)
			if msg.Text != "" {
				msg.ParseMode = tgbotapi.ModeMarkdown
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(msg); err != nil {

					log.Panic(err)
				}
			}
		}
	}
}
