package main

import (
	"github.com/adlio/trello"
	tb "gopkg.in/tucnak/telebot.v2"
)

func replyWith(sender *tb.User, b *tb.Bot, options map[string]string, label string) error {
	inlineBtns := make([][]tb.InlineButton, len(options))

	i := 0
	for id, name := range options {
		inlineBtns[i] = []tb.InlineButton{
			tb.InlineButton{
				Unique: id,
				Text:   name,
			},
		}
		i++
	}

	_, err := b.Send(sender, label, &tb.ReplyMarkup{
		InlineKeyboard: inlineBtns,
	})

	return err
}

func replyWithBoards(sender *tb.User, b *tb.Bot) error {
	client := trello.NewClient(trelloAppKey, trelloToken)

	member, err := client.GetMember(trelloUser, trello.Defaults())
	if err != nil {
		return err
	}

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		return err
	}

	options := make(map[string]string, len(boards))

	for _, board := range boards {
		options[board.ID] = board.Name
	}

	return replyWith(sender, b, options, "Which board?")
}

func replyWithLists(sender *tb.User, b *tb.Bot, selectedBoard string) error {
	client := trello.NewClient(trelloAppKey, trelloToken)

	board, err := client.GetBoard(selectedBoard, trello.Defaults())
	if err != nil {
		return err
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return err
	}

	options := make(map[string]string, len(lists))

	for _, list := range lists {
		options[list.ID] = list.Name
	}

	return replyWith(sender, b, options, "Which list?")
}

func replyWithControls(message *tb.Message, b *tb.Bot, text string) error {
	replyBtn := tb.ReplyButton{Text: "/reset"}
	replyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{replyBtn},
	}

	_, err := b.Reply(message, text, &tb.ReplyMarkup{
		ResizeReplyKeyboard: true,
		ForceReply:          false,
		ReplyKeyboard:       replyKeys,
	})

	return err
}
