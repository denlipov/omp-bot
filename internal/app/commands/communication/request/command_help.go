package request

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

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
