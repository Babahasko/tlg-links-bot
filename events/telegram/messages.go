package telegram

const msgHelp = `I can save your pages and can offer you to read them randomly

In order to save pages, send me a link to it

In order to get a random page from your list, send me /rnd
Caution! After that it will be removed!`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedPages = "You have no saved pages"
	msgSaved = "Saved!"
	msgAlreadyExist = "You already have this page"
)