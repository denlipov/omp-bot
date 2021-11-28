package request

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommunicationRequestCommander) replyBotMsg(inputMsg *tgbotapi.Message, replyMsg string) {

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, replyMsg)
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending reply message to chat ('%s') - %v", replyMsg, err)
	}
	log.Println(replyMsg)
}
