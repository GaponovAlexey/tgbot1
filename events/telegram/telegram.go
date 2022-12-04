package telegram

import (
	"errors"

	"tg/sitesess.ca/client/telegram"
	"tg/sitesess.ca/events"
	"tg/sitesess.ca/lib/e"
	"tg/sitesess.ca/storage"

)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}
type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return nil, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		p.processMessage(event)
	default:
		return e.Wrap("can't process message", errors.New("unknown event"))
	}
}
func (p *Processor) processMessage(event events.Event) error {
	meta, err:= meta(event)
	if err!= nil {
		return e.Wrap("can't process message", err)
	}
	
}

func meta(event events.Event) (Meta, error) {
	res, ok:= event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get Meta", errors.New("unknown meta type"))
	}
	return res, nil
}
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.UserName,
		}
	}

	//chatId username

}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message

}
