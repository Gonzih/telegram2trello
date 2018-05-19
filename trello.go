package main

import (
	"github.com/adlio/trello"
	"github.com/badoux/goscraper"
)

const (
	titleLimit = 70
)

func extractUrl(input string) string {
	return urlRegexp.FindString(input)
}

func generateCardName(title, desc string) string {
	if len(title) < titleLimit && len(desc) > 0 {
		lenLeft := titleLimit - len(title)

		if lenLeft > len(desc) {
			title += " - " + desc
		} else {
			title += " - " + desc[0:lenLeft]
		}

		if lenLeft < len(desc) {
			title += "..."
		}
	}

	return title
}

func storeMessageInTrello(listID string, payload string) error {
	client := trello.NewClient(trelloAppKey, trelloToken)

	card := trello.Card{}
	card.Desc = payload

	if urlRegexp.MatchString(payload) {
		url := extractUrl(payload)
		scrape, err := goscraper.Scrape(url, 10)

		if err != nil {
			return err
		}

		card.Name = generateCardName(scrape.Preview.Title, scrape.Preview.Description)
	} else {
		card.Name = payload
	}

	if len(card.Name) == 0 {
		card.Name = payload
	}

	list, err := client.GetList(listID, trello.Defaults())
	if err != nil {
		return err
	}

	return list.AddCard(&card, trello.Defaults())
}
