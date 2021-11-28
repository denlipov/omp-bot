package request

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommunicationRequestCommander) Get(inputMsg *tgbotapi.Message) {

	arg := inputMsg.CommandArguments()

	idx, err := strconv.Atoi(arg)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Wrong command args: %s", arg))
		return
	}

	req, err := c.service.Describe(uint64(idx))
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Fail to get Request with idx %d: %v", idx, err))
		return
	}

	c.replyBotMsg(inputMsg, req.String())
}
