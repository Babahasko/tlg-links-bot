package telegram

import (
	"errors"
	"example/tlgbot/clients/telegram"
	"example/tlgbot/events"
	"example/tlgbot/storage"
	"fmt"
)

type TlgProcessor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func New(client *telegram.Client, storage storage.Storage) *TlgProcessor {
	return &TlgProcessor{
		tg:      client,
		offset:  0,
		storage: storage,
	}
}

func (p *TlgProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch events fail: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1

	return res, nil

}

func (p *TlgProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can`t process message: %w", ErrUnknownEventType)
	}
}

func (p *TlgProcessor) processMessage(event events.Event) error{
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("process Message fail: %w", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("doCmd fail: %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok{
		return Meta{}, fmt.Errorf("can`t get meta")
	}
	return res, nil

}

func event(u telegram.Update) events.Event {
	updType := fetchType(u)
	res := events.Event{
		Type: updType,
		Text: fetchText(u),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}
	return res
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}
