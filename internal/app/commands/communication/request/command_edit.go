package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/denlipov/com-request-api/pkg/com-request-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (c *CommunicationRequestCommander) Edit(inputMsg *tgbotapi.Message) {

	cmdLine := inputMsg.CommandArguments()

	args := strings.SplitN(cmdLine, " ", 2)
	if len(args) != 2 {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Wrong command args: %s", args))
		return
	}

	idx, _ := strconv.ParseUint(args[0], 10, 64)
	entryJSON := args[1]

	var req pb.Request
	err := json.Unmarshal([]byte(entryJSON), &req)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Failed to unmarshal user data: %v", err))
		return
	}

	err = c.service.Update(idx, req)
	if err != nil {
		c.replyBotMsg(inputMsg, fmt.Sprintf("Failed to update request of idx %d: %v", idx, err))
		return
	}

	c.replyBotMsg(inputMsg, fmt.Sprintf("Request with id %d updated OK", idx))
}
