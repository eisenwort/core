package ewc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MessageService struct {
	BaseService
	api               *ApiService
	dbService         *DbMessageService
	dbUserService     *DbUserService
	MessageListChan   chan *MessageCollection
	MessageChan       chan *Message
	MessageDeleteChan chan int64
}

func NewMessageService() *MessageService {
	srv := new(MessageService)
	srv.api = NewApiService()
	srv.dbService = NewDbMessageService()
	srv.dbUserService = NewDbUserService()
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.MessageChan = make(chan *Message, chanSize)
	srv.MessageListChan = make(chan *MessageCollection, chanSize)
	srv.MessageDeleteChan = make(chan int64, chanSize)

	go srv.listeners()
	return srv
}

func (srv *MessageService) Send(msg *Message, text string) {
	user := srv.dbUserService.Get(userID)

	if user.Reseted {
		srv.ErrorsChan <- "Произошла неизвестная ошибка"
		return
	}

	srv.api.post("/messages", msg, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка отправки сообщения"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			srv.ErrorsChan <- "Ошибка отправки сообщения"
			return
		}
		srv.MessageChan <- msg
		srv.dbService.Create(msg)
	})
}

func (srv *MessageService) Delete(msg *Message) {
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

func (srv *MessageService) GetByChat(chatID int64) {
	user := srv.dbUserService.Get(userID)

	if user.Reseted {
		srv.ErrorsChan <- "Произошла неизвестная ошибка"
		return
	}

	messages := srv.dbService.GetByChat(chatID)
	col := NewMessageCollection()

	if messages != nil {
		col.s = messages
		srv.MessageListChan <- col
	}

	requestUrl := fmt.Sprintf("/chats/%d/messages", chatID)
	srv.api.get(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка олучения сообщений"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&messages); err != nil {
			srv.ErrorsChan <- "Ошибка олучения сообщений"
			return
		}
		col.s = messages
		srv.MessageListChan <- col
	})
}

func (srv *MessageService) listeners() {
	for {
		select {
		case _ = <-srv.api.RequestErrorChan:
			srv.ErrorsChan <- "Ошибка подключения сети"
		}
	}
}
