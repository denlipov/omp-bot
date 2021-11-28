package request

import (
	"log"

	"github.com/denlipov/omp-bot/internal/app/path"
	service "github.com/denlipov/omp-bot/internal/service/communication/request"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RequestCommander interface {
	Help(inputMsg *tgbotapi.Message)
	Get(inputMsg *tgbotapi.Message)
	List(inputMsg *tgbotapi.Message)
	Delete(inputMsg *tgbotapi.Message)

	New(inputMsg *tgbotapi.Message)  // return error not implemented
	Edit(inputMsg *tgbotapi.Message) // return error not implemented
}

type CommunicationRequestCommander struct {
	bot     *tgbotapi.BotAPI
	service service.RequestService
}

func NewCommunicationRequestCommander(bot *tgbotapi.BotAPI) *CommunicationRequestCommander {
	service := service.NewDummyRequestService()

	return &CommunicationRequestCommander{
		bot:     bot,
		service: service,
	}
}

func (c *CommunicationRequestCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.CallbackList(callback, callbackPath)
	default:
		log.Printf("CommunicationRequestCommander.HandleCallback: unknown callback name: %s", callbackPath.CallbackName)
	}
}

func (c *CommunicationRequestCommander) HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.CommandName {
	case "help":
		c.Help(message)
	case "get":
		c.Get(message)
	case "list":
		c.List(message)
	case "delete":
		c.Delete(message)
	case "new":
		c.New(message)
	case "edit":
		c.Edit(message)
	default:
		log.Printf("CommunicationRequestCommander.HandleCommand: unknown command name: %s", commandPath.CommandName)
	}
}
