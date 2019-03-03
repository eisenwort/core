package ewc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatService struct {
	BaseService
	ApiService
	dbService      *DbChatService
	ChatChan       chan *Chat
	ChatListChan   chan *ChatCollection
	ChatDeleteChan chan bool
	ChatCleanChan  chan bool
}

func NewChatService() *ChatService {
	srv := new(ChatService)
	srv.dbService = NewDbChatService(driver, connectionString)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.ChatChan = make(chan *Chat, chanSize)
	srv.ChatDeleteChan = make(chan bool, chanSize)
	srv.ChatCleanChan = make(chan bool, chanSize)
	srv.ChatListChan = make(chan *ChatCollection, chanSize)

	go srv.listeners()
	return srv
}

func (srv *ChatService) GetChats() {
	chatCollection := NewChatCollection()
	chats, err := srv.dbService.GetForUser(userID)

	if err == nil {
		chatCollection.s = chats
		srv.ChatListChan <- chatCollection
	}
	srv.get("/chats", func(r *http.Response) {
		chats := make([]*Chat, 0)

		if err := json.NewDecoder(r.Body).Decode(&chats); err != nil {
			srv.ErrorsChan <- "Ошибка получения чатов"
			return
		}

		chatCollection.s = chats
		srv.ChatListChan <- chatCollection
	})
}

func (srv *ChatService) Get(id int64, withMessages bool) {
	chat, err := srv.dbService.Get(id, withMessages)

	if err == nil {
		srv.ChatChan <- chat
	}

	requestUrl := fmt.Sprintf("/chats/%d", id)
	srv.get(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка получения чата"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(chat); err != nil {
			srv.ErrorsChan <- "Ошибка получения чата"
			return
		}
		srv.ChatChan <- chat
		srv.dbService.Update(chat)
	})
}

func (srv *ChatService) Create(chat *Chat) {
	data, _ := json.Marshal(chat)

	srv.post("/chats", data, func(r *http.Response) {
		if r.StatusCode != http.StatusCreated {
			srv.ErrorsChan <- "Ошибка создания чата"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(chat); err != nil {
			srv.ErrorsChan <- "Ошибка создания чата"
			return
		}

		srv.ChatChan <- chat
		srv.dbService.Create(chat)
	})
}

func (srv *ChatService) Delete(chat *Chat) {
	requestUrl := fmt.Sprintf("/chats/%d", chat.ID)

	srv.delete(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка удаления чата"
			return
		}
		srv.ChatDeleteChan <- true
		srv.dbService.Delete(chat)
	})
}

func (srv *ChatService) Exit(chat *Chat) {
	requestUrl := fmt.Sprintf("/chats/%d/exit", chat.ID)

	srv.delete(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка выхода из чата"
			return
		}
		srv.ChatDeleteChan <- true
		srv.dbService.Exit(chat)
	})
}

func (srv *ChatService) Clean(chat *Chat) {
	requestUrl := fmt.Sprintf("/chats/%d/clean", chat.ID)
	result := true

	srv.delete(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			result = false
			return
		}
		if !srv.dbService.Clean(chat) {
			result = false
			return
		}
	})
	if result {
		srv.InfoChan <- "Чат очищен"
	} else {
		srv.ErrorsChan <- "Ошибка очищения чата"
	}
	srv.ChatCleanChan <- result
}

func (srv *ChatService) listeners() {

}
