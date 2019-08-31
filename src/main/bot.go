package main

import (
	"clarifai_api"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
)

const (
	TOKEN = "966639602:AAHBNcZv6G9otayUb48ZPu617j5l2du3jP4"
	PROXY = "socks5://192.169.152.231:31403"
)

func SetLanguage(args []string) string {
	var reply string
	if len(args) != 1 {
		reply = `Invalid number of arguments. Command expects one argument: language code`
	} else {
		clarifai_api.Language = args[0]
		reply = `Language set successful`
	}
	return reply
}

func main() {
	client := &http.Client{}
	if len(PROXY) > 0 {
		tgProxyURL, err := url.Parse(PROXY)
		if err != nil {
			log.Printf("Failed to parse proxy URL:%s\n", err)
		}
		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			log.Printf("Failed to obtain proxy dialer: %s\n", err)
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}

	bot, err := tgbotapi.NewBotAPIWithClient(TOKEN, client)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		userMsg := update.Message.Text
		replySplit := strings.Split(userMsg, " ")
		var reply string

		switch replySplit[0] {
		case "/start", "/help":
			reply = `This bot can tell you what is on image you sent. The bot uses Clarifai Predict API. Send an image to get response
/set_lang {language code} - set language to use in Clarifai responses (default is English). Available language codes: http://developer-dev.clarifai.com/developer/guide/languages#supported-languages
Code for Russian: ru
`
		case "/set_lang":
			reply = SetLanguage(replySplit[1:])
		default:
			reply = `Available commands: /start (or /help), /set_lang`
		}

		if update.Message.Photo != nil {
			photo := *update.Message.Photo
			var resp tgbotapi.File
			var respErr error
			if len(photo) > 1 {
				resp, respErr = bot.GetFile(tgbotapi.FileConfig{photo[1].FileID})
			} else {
				resp, respErr = bot.GetFile(tgbotapi.FileConfig{photo[0].FileID})
			}
			if respErr != nil {
				reply = "Error: " + err.Error()
			} else {
				photoUrl := "https://api.telegram.org/file/bot" + TOKEN + "/" + resp.FilePath
				reply = clarifai_api.GetClarifaiResp(photoUrl)
			}
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
