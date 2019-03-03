package ewc

import (
	"fmt"
	"net/http"
)

type MessageService struct {
	BaseService
	api               *ApiService
	dbService         *DbMessageService
	MessageListChan   chan *MessageCollection
	MessageChan       chan *Message
	MessageDeleteChan chan bool
}

func NewMessageService() *MessageService {
	srv := new(MessageService)
	srv.api = NewApiService()
	srv.dbService = NewDbMessageService(driver, connectionString)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.MessageChan = make(chan *Message, chanSize)
	srv.MessageListChan = make(chan *MessageCollection, chanSize)
	srv.MessageDeleteChan = make(chan bool, chanSize)

	go srv.listeners()
	return srv
}

func (srv *MessageService) Send(msg *Message) {
	data, _ := json.Marshal(msg)

	srv.api.post("/messages", data, func(r *http.Response) {
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
	requestUrl := fmt.Sprintf("/messages/%d", msg.ID)

	srv.api.delete(requestUrl, func(r *http.Response) {
		isDeleted := r.StatusCode == http.StatusOK
		srv.MessageDeleteChan <- isDeleted

		if isDeleted {
			srv.InfoChan <- "Сообщение удалено"
			srv.dbService.Delete(msg)
		}
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
