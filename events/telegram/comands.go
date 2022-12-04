package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"tg/sitesess.ca/lib/e"
	"tg/sitesess.ca/storage"

)

const (
	RndCmd   = "/rnd"
	HempCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {

	}

	switch text {
	case RndCmd:

	case HempCmd:

	case StartCmd:

	default:

	}

}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) SendRandom(chatId int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}
func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}