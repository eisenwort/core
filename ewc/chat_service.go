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
	dbUserService  *DbUserService
	saveListChat   chan []Chat
	ChatChan       chan string
	ChatListChan   chan string
	ChatDeleteChan chan bool
	ChatCleanChan  chan bool
}

func NewChatService() *ChatService {
	srv := new(ChatService)
	srv.dbService = NewDbChatService()
	srv.dbUserService = NewDbUserService()
	srv.saveListChat = make(chan []Chat, chanSize)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.ChatChan = make(chan string, chanSize)
	srv.ChatDeleteChan = make(chan bool, chanSize)
	srv.ChatCleanChan = make(chan bool, chanSize)
	srv.ChatListChan = make(chan string, chanSize)

	go srv.listeners()
	return srv
}

func (srv *ChatService) GetChats() {
	chats, err := srv.dbService.GetForUser(userID)

	if err == nil {
		srv.ChatListChan <- serialize(chats)
	}
	srv.get("/chats", func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка получения чатов"
			return
		}

		jsonData := getBodyString(r.Body)
		srv.ChatListChan <- jsonData

		var chats []Chat
		deserialize(jsonData, &chats)
		srv.saveListChat <- chats
	})
}

func (srv *ChatService) Get(id int64) {
	chat, err := srv.dbService.Get(id, []string{"messages"})

	if err == nil {
		srv.ChatChan <- serialize(chat)
	}

	requestUrl := fmt.Sprintf("/chats/%d?include=messages", id)
	srv.get(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка получения чата"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(chat); err != nil {
			srv.ErrorsChan <- "Ошибка получения чата"
			return
		}
		srv.ChatChan <- serialize(chat)
		srv.dbService.Update(chat)
	})
}

func (srv *ChatService) CreatePersonalChat(login string) {
	self := srv.dbUserService.Get(userID)
	user := User{}

	if self.Reseted {
		srv.ErrorsChan <- "Ошибка создания чата"
		return
	}
	srv.get("/users/login/"+login, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- fmt.Sprintf("Пользователя с ником %s не существует", login)
			return
		}

		body := getBodyString(r.Body)

		//if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		if err := json.Unmarshal([]byte(body), &user); err != nil {
			srv.ErrorsChan <- "Ошибка создания чата"
		}
	})
	if user.ID == 0 {
		return
	}

	chat := Chat{
		Users:    []User{user},
		OwnerID:  userID,
		Personal: true,
	}
	srv.post("/chats", chat, func(r *http.Response) {
		if r.StatusCode != http.StatusCreated {
			srv.ErrorsChan <- "Ошибка создания чата"
			return
		}

		data := getBodyString(r.Body)

		if err := json.Unmarshal([]byte(data), &chat); err != nil {
			srv.ErrorsChan <- "Ошибка создания чата"
			return
		}

		srv.ChatChan <- data

		deserialize(data, &chat)
		srv.dbService.Create(chat)
	})
}

func (srv *ChatService) Create(userLogin string) {
	chat := Chat{}

	chat.OwnerID = userID
	chat.Personal = true
	self := srv.dbUserService.Get(userID)

	if self.Reseted {
		srv.ErrorsChan <- "Ошибка создания чата"
		return
	}
	if chat.Personal {
		user := srv.dbUserService.GetByLogin(userLogin)

		if user.ID == 0 {
			user := User{}
			srv.get("/users/login/"+userLogin, func(r *http.Response) {
				if r.StatusCode != http.StatusOK {
					srv.ErrorsChan <- fmt.Sprintf("Пользователя с ником %s не существует", userLogin)
					return
				}
				if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
					srv.ErrorsChan <- "Ошибка создания чата"
				}
			})
			if user.ID == 0 {
				srv.ErrorsChan <- fmt.Sprintf("Пользователя с ником %s не существует", userLogin)
				return
			}
		} else {
			chat.Users[0] = user
		}
	}

	srv.post("/chats", chat, func(r *http.Response) {
		if r.StatusCode != http.StatusCreated {
			srv.ErrorsChan <- "Ошибка создания чата"
			return
		}

		srv.ChatChan <- getBodyString(r.Body)
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
	for {
		select {
		case chats := <-srv.saveListChat:
			srv.dbService.SaveList(chats)
		}
	}
}
