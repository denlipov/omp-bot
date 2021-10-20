package request

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommunicationRequestCommander) Delete(inputMsg *tgbotapi.Message) {

	arg := inputMsg.CommandArguments()

	idx, err := strconv.Atoi(arg)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Wrong arg: %s", arg))
		return
	}

	_, err = c.service.Remove(uint64(idx))
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Fail to get request with idx %d: %v", idx, err))
		return
	}
	c.replyBotMsg(inputMsg, fmt.Sprintf("Request removed OK"))
}
