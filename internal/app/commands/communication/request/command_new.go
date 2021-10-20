package request

import (
	"encoding/json"
	"fmt"

	comm "github.com/denlipov/omp-bot/internal/model/communication"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommunicationRequestCommander) New(inputMsg *tgbotapi.Message) {

	arg := inputMsg.CommandArguments()

	var req comm.Request
	entryJSON := arg
	err := json.Unmarshal([]byte(entryJSON), &req)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Failed to unmarshal user data: %v", err))
		return
	}

	idx, err := c.service.Create(req)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	c.replyBotMsg(inputMsg, fmt.Sprintf("Request created OK; new idx: %d", idx))
}
