package ewc

import (
	"fmt"
	"net/http"
)

type MessageService struct {
	BaseService
	dbService         *DbMessageService
	dbUserService     *DbUserService
	saveChan          chan Message
	saveListChan      chan []Message
	lastMessageId     int64
	MessageListChan   chan string
	MessageChan       chan string
	MessageDeleteChan chan int64
}

func NewMessageService() *MessageService {
	srv := new(MessageService)
	srv.dbService = NewDbMessageService()
	srv.dbUserService = NewDbUserService()
	srv.saveChan = make(chan Message, chanSize)
	srv.saveListChan = make(chan []Message, chanSize)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.MessageChan = make(chan string, chanSize)
	srv.MessageListChan = make(chan string, chanSize)
	srv.MessageDeleteChan = make(chan int64, chanSize)

	go srv.listeners()
	return srv
}

func (srv *MessageService) Send(msg Message) {
	httpPost("/messages", msg, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка отправки сообщения"
			return
		}

		jsonData := getBodyString(r.Body)
		deserialize(jsonData, &msg)
		srv.saveChan <- msg
	})
}

func (srv *MessageService) Delete(msg Message) {
	id := msg.ID
	requestUrl := fmt.Sprintf("/messages/%d", msg.ID)

	httpDelete(requestUrl, func(r *http.Response) {
		isDeleted := r.StatusCode == http.StatusOK

		if isDeleted {
			srv.MessageDeleteChan <- id
			srv.InfoChan <- "Сообщение удалено"
			srv.dbService.Delete(msg)
		}
	})
}

// TODO: refactor when implement websocket
func (srv *MessageService) GetByChat(chatID int64, page int) {
	messages := srv.dbService.GetByChat(chatID, page)
	cnt := len(messages)

	if messages != nil && cnt != 0 {
		if messages[cnt-1].ID > srv.lastMessageId {
			// TODO: decrypt all messages
			srv.MessageListChan <- serialize(messages)
			srv.lastMessageId = messages[cnt-1].ID
		}
	}

	requestUrl := createUrl(fmt.Sprintf("/chats/%d/messages", chatID), map[string]string{
		"page": fmt.Sprintf("%d", page),
	})
	httpGet(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка получения сообщений"
			return
		}

		jsonData := getBodyString(r.Body)
		messages = []Message{}

		deserialize(jsonData, &messages)
		cnt = len(messages)

		if cnt != 0 && messages[cnt-1].ID > srv.lastMessageId {
			// TODO: decrypt all messages
			srv.MessageListChan <- jsonData
			srv.saveListChan <- messages
		}
	})
}

func (srv *MessageService) SetAllIsRead(chatID int64) {
	srv.dbService.SetAllIsRead(chatID)
}

func (srv *MessageService) listeners() {
	for {
		select {
		case messages := <-srv.saveListChan:
			for _, message := range messages {
				srv.dbService.Save(message)
			}
		case msg := <-srv.saveChan:
			srv.dbService.Save(msg)
		}
	}
}
