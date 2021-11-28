package request

import (
	"encoding/json"
	"log"

	"github.com/denlipov/omp-bot/internal/app/path"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CallbackListData struct {
	Offset uint64 `json:"offset"`
}

const maxEntriesToList = 5

func (c *CommunicationRequestCommander) CallbackList(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {

	nextPageOffsetData := CallbackListData{}
	err := json.Unmarshal([]byte(callbackPath.CallbackData), &nextPageOffsetData)
	if err != nil {
		log.Printf("CommunicationRequestCommander.CallbackList: "+
			"error reading json data for type CallbackListData from "+
			"input string %v - %v", callbackPath.CallbackData, err)
		return
	}

	requests, _ := c.service.List(nextPageOffsetData.Offset, maxEntriesToList)
	nReqs := len(requests)
	outputMsgText := ""
	if nReqs > 0 {
		outputMsgText = "Requests list:\n\n"
	} else {
		outputMsgText = "No more requests found"
	}
	for _, req := range requests {
		outputMsgText += req.String() + "\n"
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, outputMsgText)

	if nReqs == maxEntriesToList {
		prepareCallbackMsg(&msg, uint64(nReqs)+nextPageOffsetData.Offset)
	}

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("CommunicationRequestCommander.CallbackList: error sending reply message to chat - %v", err)
	}
}

func (c *CommunicationRequestCommander) List(inputMsg *tgbotapi.Message) {

	requests, _ := c.service.List(0, maxEntriesToList)
	nReqs := len(requests)
	outputMsgText := ""
	if nReqs > 0 {
		outputMsgText = "Requests list:\n\n"
	} else {
		outputMsgText = "No requests found"
	}

	for _, req := range requests {
		outputMsgText += req.String() + "\n"
	}

	log.Println(outputMsgText)

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, outputMsgText)

	if nReqs == maxEntriesToList {
		prepareCallbackMsg(&msg, uint64(nReqs))
	}

	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("CommunicationRequestCommander.List: error sending reply message to chat - %v", err)
	}
}

func prepareCallbackMsg(msg *tgbotapi.MessageConfig, offset uint64) {
	serializedData, _ := json.Marshal(CallbackListData{
		Offset: offset,
	})

	callbackPath := path.CallbackPath{
		Domain:       "communication",
		Subdomain:    "request",
		CallbackName: "list",
		CallbackData: string(serializedData),
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Next page", callbackPath.String()),
		),
	)
}
