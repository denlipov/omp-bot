package request

import (
	comm "github.com/denlipov/omp-bot/internal/model/communication"
	service "github.com/denlipov/omp-bot/internal/service/communication/request"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
        "github.com/denlipov/omp-bot/internal/app/path"
        "log"
        "encoding/json"
        "fmt"
        "strconv"
        "strings"
)

type RequestCommander interface {
	Help(inputMsg *tgbotapi.Message)
	Get(inputMsg *tgbotapi.Message)
	List(inputMsg *tgbotapi.Message)
	Delete(inputMsg *tgbotapi.Message)

	New(inputMsg *tgbotapi.Message)  // return error not implemented
	Edit(inputMsg *tgbotapi.Message) // return error not implemented
}

type CommunicationRequestCommander struct {
	bot     *tgbotapi.BotAPI
	service service.RequestService
}

func NewCommunicationRequestCommander(bot *tgbotapi.BotAPI) *CommunicationRequestCommander {
	service := service.NewDummyRequestService()

	return &CommunicationRequestCommander{
		bot:     bot,
		service: service,
	}
}


func (c *CommunicationRequestCommander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	switch callbackPath.CallbackName {
	case "list":
		c.CallbackList(callback, callbackPath)
	default:
		log.Printf("CommunicationRequestCommander.HandleCallback: unknown callback name: %s", callbackPath.CallbackName)
	}
}

func (c *CommunicationRequestCommander) HandleCommand(message *tgbotapi.Message, commandPath path.CommandPath) {
	switch commandPath.CommandName {
	case "help":
		c.Help(message)
	case "get":
		c.Get(message)
	case "list":
		c.List(message)
	case "delete":
		c.Delete(message)
	case "new":
		c.New(message)
	case "edit":
		c.Edit(message)
	default:
		log.Printf("CommunicationRequestCommander.HandleCommand: unknown command name: %s", commandPath.CommandName)
	}
}


// XXX перенести
type CallbackListData struct {
	Offset uint64 `json:"offset"`
}

// работает с конопкой
// метод, не покрывающий интерфейс RequestCommander

func (c *CommunicationRequestCommander) CallbackList(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {

	parsedData := CallbackListData{}
	err := json.Unmarshal([]byte(callbackPath.CallbackData), &parsedData)
	if err != nil {
		log.Printf("CommunicationRequestCommander.CallbackList: "+
			"error reading json data for type CallbackListData from "+
			"input string %v - %v", callbackPath.CallbackData, err)
		return
	}

	requests, _ := c.service.List(parsedData.Offset, maxEntriesToList)
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
                prepareCallbackMsg(&msg, uint64(nReqs) + parsedData.Offset)
        }

	_, err = c.bot.Send(msg)
	if err != nil {
		log.Printf("CommunicationRequestCommander.CallbackList: error sending reply message to chat - %v", err)
	}
}


// методы интерфейса RequestCommander

func (c *CommunicationRequestCommander) replyBotMsg(inputMsg *tgbotapi.Message, replyMsg string) {

	msg := tgbotapi.NewMessage(inputMsg.Chat.ID, replyMsg)
        _, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending reply message to chat ('%s') - %v", replyMsg, err)
	}
        log.Println(replyMsg)
}

const maxEntriesToList = 5

func (c *CommunicationRequestCommander) Help(inputMsg *tgbotapi.Message) {

        const helpMsgText = `/help__communication__request - показывает текущую справку
/get__communication__request <id> - получить Request по id
/list__communication__requrest - показать список имеющихся Request'ов
/edit__comunication__request <id> <request> - редактировать Request с данным id; формат: id {"user":"имя","desc":"текст"}
/new__communication__request <request> - создать новый Request; формат: {"user":"имя","desc":"текст"}
/delete__communication__request <id> - удалить Request по данному id
`
        c.replyBotMsg(inputMsg, helpMsgText)
}


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


func (c *CommunicationRequestCommander) Edit(inputMsg *tgbotapi.Message) {

        cmdLine := inputMsg.CommandArguments()

        args := strings.SplitN(cmdLine, " ", 2)
        if len(args) != 2 {
                c.replyBotMsg(inputMsg, fmt.Sprintf("Wrong command args: %s", args))
                return
        }

        idx, _ := strconv.ParseUint(args[0], 10, 64)
        entryJSON := args[1]

        var req comm.Request
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

        c.replyBotMsg(inputMsg, fmt.Sprintf("Request created OK; new idx: %d", idx))
}
