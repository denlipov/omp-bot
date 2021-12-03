package communication

import (
	"github.com/denlipov/omp-bot/internal/app/commands/communication/request"
	"github.com/denlipov/omp-bot/internal/app/path"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

type Commander interface {
	HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath)
	HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath)
}

type CommunicationCommander struct {
	bot              *tgbotapi.BotAPI
	requestCommander Commander
}

func NewCommunicationCommander(bot *tgbotapi.BotAPI) *CommunicationCommander {
	return &CommunicationCommander{
		bot:              bot,
		requestCommander: request.NewCommunicationRequestCommander(bot),
	}
}

func (c *CommunicationCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.Subdomain {
	case "request":
		c.requestCommander.HandleCallback(callback, callbackPath)
	default:
		log.Info().Msgf("CommunicationCommander.HandleCallback: unknown subdomain - %s", callbackPath.Subdomain)
	}
}

func (c *CommunicationCommander) HandleCommand(msg *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.Subdomain {
	case "request":
		c.requestCommander.HandleCommand(msg, commandPath)
	default:
		log.Info().Msgf("CommunicationCommander.HandleCommand: unknown subdomain - %s", commandPath.Subdomain)
	}
}
