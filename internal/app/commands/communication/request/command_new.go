package request

import (
	"encoding/json"
	"fmt"

	pb "github.com/denlipov/com-request-api/pkg/com-request-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

func (c *CommunicationRequestCommander) New(inputMsg *tgbotapi.Message) {

	arg := inputMsg.CommandArguments()

	var req pb.Request
	entryJSON := arg
	err := json.Unmarshal([]byte(entryJSON), &req)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Failed to unmarshal user data: %v", err))
		return
	}

	log.Info().Msgf("Message received: %+v", req)

	idx, err := c.service.Create(req)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	c.replyBotMsg(inputMsg, fmt.Sprintf("Request created OK; new idx: %d", idx))
}
