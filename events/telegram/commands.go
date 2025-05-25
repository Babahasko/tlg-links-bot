package telegram

import (
	"errors"
	"example/tlgbot/storage"
	"fmt"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd = "/rnd"
	HelpCmd = "/help"
	StartCmd = "/start"
)
// Выполняет комманды из сообщения
func (p *TlgProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("get command '%s' from '%s'", text, username)
	// save page: htpp://...
	// rnd: get:/rnd random page
	// help: /help
	// start: /start: hi + help
	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
		// TODO: AddPage()
	}
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHi(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *TlgProcessor) savePage(chatID int, pageURL string, username string) (err error) {
	page := &storage.Page{
		URL: pageURL,
		UserName: username,
	}
	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return fmt.Errorf("save page fail: %w", err)
	}
	if isExist{
		return p.tg.SendMessage(chatID, msgAlreadyExist)
	}

	if err := p.storage.Save(page); err != nil {
		return fmt.Errorf("save page fail: %w", err)
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return fmt.Errorf("send message fail")
	}
	return nil
}

func (p *TlgProcessor) sendRandom(chatID int, username string) (err error) {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages){
		return fmt.Errorf("can`t pick random: %w", err)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return fmt.Errorf("can`t send page: %w", err)
	}
	
	return p.storage.Remove(page)
}

func (p *TlgProcessor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *TlgProcessor) sendHi(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	if strings.HasPrefix(text,"/") {
		return false
	}
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}