package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func tgbot(totg chan []neededstruct, fromtg chan string) {
	tgusers := []int64{924377144}
	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			if update.Message != nil { // ignore any non-Message Updates

				fromtg <- update.Message.Text
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "command sended..")
				bot.Send(msg)

			}
		case data := <-totg:
			fmt.Println("SEND NOTIFY", data)
			strok := []string{}
			for _, d := range data {
				strok = append(strok, fmt.Sprintln("*ALARM* for _", d.name, "_\n values range ", d.min, "-", d.max, "\n value now:", d.value))
			}
			for _, uid := range tgusers {
				msg := tgbotapi.NewMessage(uid, strings.Join(strok, "\n"))
				msg.ParseMode = "markdown"
				bot.Send(msg)

			}
		}
	}
}
