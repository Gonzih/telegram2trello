package main

import (
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func cbHandler(b *tb.Bot) func(*tb.Callback) {
	return func(cb *tb.Callback) {
		payload := strings.Trim(cb.Data, "\f")

		log.Printf("CB %#v", cb)

		if cb.Sender.Username != whitelistedTelegramUser {
			return
		}

		userID := cb.Sender.ID

		state, err := session.Get(userID, stateKey)

		if state == "" || err != nil {
			log.Printf("No state found in session for callback: %v %v", state, err)
		}

		switch state {
		case "waitingForBoard":
			log.Println("Got board response")

			session.Set(userID, map[string]interface{}{
				boardKey: payload,
				stateKey: "waitingForList",
			})

			err := replyWithLists(cb.Sender, b, payload)

			if err != nil {
				b.Send(cb.Sender, err.Error())
			}
		case "waitingForList":
			log.Println("Got list response")
			list := payload

			session.Set(userID, map[string]interface{}{
				listKey:  list,
				stateKey: "done",
			})

			log.Println("geting payload")
			sessionPayload, err := session.Get(userID, payloadKey)

			if err != nil {
				log.Printf("Error getting session payloaad %s", err)
			}

			log.Println("storing message", sessionPayload)
			err = storeMessageInTrello(list, sessionPayload)

			if err != nil {
				log.Printf("Error storing message: %s", err)
			}

			log.Println("responding")
			err = b.Respond(cb, &tb.CallbackResponse{
				Text: "Done!",
			})

			if err != nil {
				log.Printf("Error responding to callback: %s", err)
			}
		}

		b.Respond(cb)
	}
}

func textHandler(b *tb.Bot) func(*tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() || m.Sender.Username != whitelistedTelegramUser {
			b.Reply(m, "I don't want to talk to you, sorry.")
			return
		}

		userID := m.Sender.ID

		state, err := session.Get(userID, stateKey)

		if state == "" || err != nil {
			log.Println("Starting")
			session.Set(userID, map[string]interface{}{
				payloadKey: m.Text,
				stateKey:   "waitingForBoard",
			})
			replyWithBoards(m.Sender, b)
			return
		}

		if state == "done" {
			list, err := session.Get(userID, listKey)

			if err != nil {
				log.Println(err)
				return
			}

			err = storeMessageInTrello(list, m.Text)

			if err != nil {
				log.Println(err)
				return
			}

			err = replyWithControls(m, b, "Done!")

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func resetHandler(b *tb.Bot) func(*tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() || m.Sender.Username != whitelistedTelegramUser {
			b.Reply(m, "I don't want to talk to you, sorry.")
			return
		}

		log.Println("Reseting the session")

		err := session.Clear(m.Sender.ID)
		if err != nil {
			log.Println(err)
			return
		}

		b.Reply(m, "Done!")
	}
}
