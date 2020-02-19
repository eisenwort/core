package ewc

import (
	"fmt"
	"net/http"
)

type MessageService struct {
	BaseService
	api               *ApiService
	dbService         *DbMessageService
	dbUserService     *DbUserService
	saveChan          chan Message
	saveListChan      chan []Message
	MessageListChan   chan string
	MessageChan       chan string
	IdChan            chan int64
	MessageDeleteChan chan int64
}

func NewMessageService() *MessageService {
	srv := new(MessageService)
	srv.api = NewApiService()
	srv.dbService = NewDbMessageService()
	srv.dbUserService = NewDbUserService()
	srv.saveChan = make(chan Message, chanSize)
	srv.saveListChan = make(chan []Message, chanSize)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.MessageChan = make(chan string, chanSize)
	srv.MessageListChan = make(chan string, chanSize)
	srv.MessageDeleteChan = make(chan int64, chanSize)
	srv.IdChan = make(chan int64, chanSize)

	go srv.listeners()
	return srv
}

func (srv *MessageService) Send(jsonData string, text string) {
	user := srv.dbUserService.Get(userID)

	if user.Reseted {
		srv.ErrorsChan <- "Произошла неизвестная ошибка"
		return
	}

	msg := Message{}
	deserialize(jsonData, &msg)

	srv.api.post("/messages", msg, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка отправки сообщения"
			return
		}

		jsonData := getBodyString(r.Body)
		srv.MessageChan <- jsonData

		deserialize(jsonData, &msg)
		srv.saveChan <- msg
	})
}

func (srv *MessageService) Delete(msg Message) {
	id := msg.ID
	requestUrl := fmt.Sprintf("/messages/%d", msg.ID)

	srv.api.delete(requestUrl, func(r *http.Response) {
		isDeleted := r.StatusCode == http.StatusOK

		if isDeleted {
			srv.MessageDeleteChan <- id
			srv.InfoChan <- "Сообщение удалено"
			srv.dbService.Delete(msg)
		}
	})
}

func (srv *MessageService) GetByChat(chatID int64, page int) {
	user := srv.dbUserService.Get(userID)

	if user.Reseted {
		srv.ErrorsChan <- "Произошла неизвестная ошибка"
		return
	}

	messages := srv.dbService.GetByChat(chatID, page)

	if messages != nil {
		srv.MessageListChan <- serialize(messages)
	}

	requestUrl := fmt.Sprintf("/chats/%d/messages", chatID)
	srv.api.get(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка олучения сообщений"
			return
		}

		jsonData := getBodyString(r.Body)
		srv.MessageListChan <- jsonData

		deserialize(jsonData, &messages)
		srv.saveListChan <- messages
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
