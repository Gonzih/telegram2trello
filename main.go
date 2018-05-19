package main

import (
	"log"
	"os"
	"regexp"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	payloadKey = "payload"
	stateKey   = "state"
	boardKey   = "board"
	listKey    = "list"
)

var (
	urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	telegramToken,
	trelloAppKey,
	trelloToken,
	trelloUser,
	whitelistedTelegramUser string
	session *sessionStore
)

func init() {
	telegramToken = os.Getenv("TELEGRAM_TOKEN")
	trelloAppKey = os.Getenv("TRELLO_APP_KEY")
	trelloToken = os.Getenv("TRELLO_TOKEN")
	trelloUser = os.Getenv("TRELLO_USER")
	whitelistedTelegramUser = os.Getenv("WHITELISTED_TELEGRAM_USER")

	var err error
	session, err = newSessionStore()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  telegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	b.Handle(tb.OnCallback, cbHandler(b))
	b.Handle(tb.OnText, textHandler(b))
	b.Handle("/reset", resetHandler(b))

	log.Println("Starting")
	b.Start()
}
