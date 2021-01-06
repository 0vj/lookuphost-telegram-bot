package main

import (
	"log"
	"net"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

func main() {
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil { // ignore any non-Message Updates
		log.Fatal(err)
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Text != "/start" {
			validhost := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]{0,61}[a-zA-Z0-9]{0,1}\.([a-zA-Z]{1,6}|[a-zA-Z0-9-]{1,30}\.[a-zA-Z]{2,3})$`)
			if validhost.MatchString(update.Message.Text) {
				lookupback(update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "invalid host. Please do not spam ‍and send valid host🙃 😉")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to lookup bot.\nThis bot is developed for simple loockups with Go language.\n We will be happy if you help develop it:\n https://github.com/mr-tafreshi/lookuphost-telegram-bot")
		        msg.ReplyToMessageID = update.Message.MessageID
                        bot.Send(msg)
		}
	}
}

func lookupback(chatid int64, replyto int, host string) {
	var lookuplist []string
	ns, err := net.LookupNS(host)
	lookuplist = append(lookuplist, "ns (name server) :")
	if err != nil {
		lookuplist = append(lookuplist, "error")
	} else {
		for i := 0; i < len(ns); i++ {
			lookuplist = append(lookuplist, ns[i].Host)
		}
	}
	ip, err := net.LookupIP(host)
	lookuplist = append(lookuplist, "ip :")
	if err != nil {
		lookuplist = append(lookuplist, "error")
	} else {
		for i := 0; i < len(ip); i++ {
			lookuplist = append(lookuplist, ip[i].String())
		}
	}
	cname, err := net.LookupCNAME(host)
	lookuplist = append(lookuplist, "cname :")
	if err != nil {
		lookuplist = append(lookuplist, "error")
	} else {
		lookuplist = append(lookuplist, cname)
	}
	mx, err := net.LookupMX(host)
	lookuplist = append(lookuplist, "mx :")
	if err != nil {
		lookuplist = append(lookuplist, "error")
	} else {
		for i := 0; i < len(mx); i++ {
			lookuplist = append(lookuplist, mx[i].Host)
		}
	}
	msg := tgbotapi.NewMessage(chatid, strings.Join(lookuplist, "\n"))
	msg.ReplyToMessageID = replyto
	bot.Send(msg)
	return
}
