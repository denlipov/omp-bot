package communication

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/denlipov/omp-bot/internal/app/commands/communication/request"
	"github.com/denlipov/omp-bot/internal/app/path"
)

type Commander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath)
}

type CommunicationCommander struct {
	bot *tgbotapi.BotAPI
	requestCommander Commander
}

func NewCommunicationCommander(
	bot *tgbotapi.BotAPI,
) *CommunicationCommander {
	return &CommunicationCommander{
		bot: bot,
		// subdomainCommander
		requestCommander: request.NewCommunicationRequestCommander(bot),
	}
}

func (c *CommunicationCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.Subdomain {
	case "request":
		c.requestCommander.HandleCallback(callback, callbackPath)
	default:
		log.Printf("CommunicationCommander.HandleCallback: unknown subdomain - %s", callbackPath.Subdomain)
	}
}

func (c *CommunicationCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.Subdomain {
	case "request":
		c.requestCommander.HandleCommand(msg, commandPath)
	default:
		log.Printf("CommunicationCommander.HandleCommand: unknown subdomain - %s", commandPath.Subdomain)
	}
}
