package request

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

func (c *CommunicationRequestCommander) replyBotMsg(inputMsg *tgbotapi.Message, replyMsg string) {

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, replyMsg)
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Error().Msgf("Error sending reply message to chat ('%s') - %v", replyMsg, err)
	}
	log.Debug().Msgf(replyMsg)
}
