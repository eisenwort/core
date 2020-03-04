package ewc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserService struct {
	saveChan     chan User
	getChan      chan int64
	dbService    *DbUserService
	ErrorsChan   chan string
	LoginChan    chan bool
	UserListChan chan string
	UserChan     chan string
	DeleteChan   chan int64
}

func NewUserService() *UserService {
	srv := new(UserService)
	srv.dbService = NewDbUserService()
	srv.saveChan = make(chan User, chanSize)
	srv.getChan = make(chan int64, 1)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.LoginChan = make(chan bool, chanSize)
	srv.UserListChan = make(chan string, chanSize)
	srv.UserChan = make(chan string, chanSize)
	srv.DeleteChan = make(chan int64, chanSize)

	go srv.listeners()
	return srv
}

func (srv *UserService) Register(login, password, passwordForReset string) {
	tokenData := TokenData{}
	data := map[string]string{
		"login":          login,
		"password":       password,
		"reset_password": passwordForReset,
	}

	httpPost("/registration", data, func(r *http.Response) {
		if r.StatusCode == http.StatusConflict {
			srv.ErrorsChan <- "Логин уже занят"
			return
		}
		if r.StatusCode != http.StatusCreated {
			srv.ErrorsChan <- "Ошибка регистрации"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&tokenData); err != nil {
			srv.ErrorsChan <- "Ошибка регистрации"
			return
		}
	})

	if tokenData.Token == "" {
		srv.LoginChan <- false
		return
	}

	srv.LoginChan <- true

	claims := getClaims(tokenData.Token)
	userID = claims.Id
	userIDHeader = fmt.Sprintf("%d", claims.Id)

	srv.getChan <- claims.Id
	srv.dbService.SaveUserData(UserData{RefreshToken: tokenData.RefreshToken})
}

func (srv *UserService) IsLogin() {
	user := srv.dbService.IsLogin()

	if user.ID == 0 {
		srv.LoginChan <- false
		return
	}

	userID = user.ID
	userIDHeader = fmt.Sprintf("%d", user.ID)

	userData := srv.dbService.GetUserData()
	jwtToken = userData.RefreshToken
	requestUrl := fmt.Sprintf("/users/%d/refresh", userID)

	httpPost(requestUrl, nil, func(r *http.Response) {
		tokenData := TokenData{}

		if r.StatusCode != http.StatusOK {
			srv.LoginChan <- false
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&tokenData); err != nil {
			log.Println("decode refresh token data error:", err)
			return
		}

		jwtToken = tokenData.Token
		srv.LoginChan <- true
		srv.dbService.SaveUserData(UserData{RefreshToken: tokenData.RefreshToken})
	})
}

func (srv *UserService) Logout() {
	if err := os.Remove(connectionString); err != nil {
		log.Println("delete db file error")
	}
}

func (srv *UserService) Update(item *User) {
	requestUrl := fmt.Sprintf("/users/%d", item.ID)
	user := User{}

	httpPut(requestUrl, item, func(r *http.Response) {
		if r.StatusCode == http.StatusConflict {
			srv.ErrorsChan <- "Логин уже занят"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			srv.ErrorsChan <- "Ошибка обновления"
			return
		}
		srv.saveChan <- user
	})
}

func (srv *UserService) Login(login, password string) {
	data := map[string]string{
		"login":    login,
		"password": password,
	}
	tokenData := TokenData{}

	httpPost("/login", data, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Неверный логин или пароль"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&tokenData); err != nil {
			srv.ErrorsChan <- "Ошибка входа"
			return
		}
	})
	if tokenData.Token == "" {
		srv.ErrorsChan <- "Неверный логин или пароль"
		return
	}

	//decryptDb()

	jwtToken = tokenData.Token
	srv.LoginChan <- true

	claims := getClaims(tokenData.Token)
	userID = claims.Id
	userIDHeader = fmt.Sprintf("%d", claims.Id)

	srv.getChan <- claims.Id
	srv.dbService.SaveUserData(UserData{RefreshToken: tokenData.RefreshToken})
}

func (srv *UserService) getUser(id int64) {
	if id == 0 {
		srv.ErrorsChan <- "Не удалось получить пользователя"
		return
	}

	requestUrl := fmt.Sprintf("/users/%d", id)
	user := User{}

	httpGet(requestUrl, func(r *http.Response) {
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return
		}
		srv.saveChan <- user
		currentUser = user
	})
}

func (srv *UserService) GetFriends() {
	friends := srv.dbService.GetFriends(userID)

	if len(friends) != 0 {
		srv.UserListChan <- serialize(friends)
	}

	requestUrl := fmt.Sprintf("/users/%d/friends", userID)
	httpGet(requestUrl, func(r *http.Response) {
		log.Println(r.StatusCode, r.Status)
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка получения контактов"
			return
		}

		body := getBodyString(r.Body)
		srv.UserListChan <- body

		friends := []User{}
		deserialize(body, &friends)
		srv.dbService.SaveFriends(userID, friends)
	})
}

func (srv *UserService) AddFriend(login string) {
	data := map[string]string{
		"login": login,
	}
	requestUrl := fmt.Sprintf("/users/%d/friends", userID)
	httpPost(requestUrl, serializeByte(data), func(r *http.Response) {
		if r.StatusCode == http.StatusNotFound {
			srv.ErrorsChan <- "Пользователь с таким логином не найден"
			return
		}
		if r.StatusCode == http.StatusConflict {
			srv.ErrorsChan <- "Пользователь уже добавлен"
			return
		}
		if r.StatusCode != http.StatusCreated {
			srv.ErrorsChan <- "Ошибка добавления контакта"
			return
		}

		body := getBodyString(r.Body)
		srv.UserChan <- body

		friend := User{}
		deserialize(body, &friend)
		srv.dbService.AddFriend(userID, friend.ID)
	})
}

func (srv *UserService) DeleteFriend(id int64) {
	srv.dbService.DeleteFriend(userID, id)

	requestUrl := fmt.Sprintf("/users/%d/friends/%d", userID, id)
	httpDelete(requestUrl, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Ошибка удаления контакта"
			return
		}
		srv.DeleteChan <- id
	})
}

func (srv *UserService) listeners() {
	for {
		select {
		case user := <-srv.saveChan:
			srv.dbService.Save(&user)
		case id := <-srv.getChan:
			srv.getUser(id)
		}
	}
}
