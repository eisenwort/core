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
	ApiService
	saveChan   chan User
	getChan    chan int64
	dbService  *DbUserService
	UserChan   chan *User
	ErrorsChan chan string
	LoginChan  chan bool
}

func NewUserService() *UserService {
	srv := new(UserService)
	srv.dbService = NewDbUserService()
	srv.saveChan = make(chan User, chanSize)
	srv.getChan = make(chan int64, 1)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.UserChan = make(chan *User, chanSize)
	srv.LoginChan = make(chan bool, chanSize)

	go srv.listeners()
	return srv
}

func (srv *UserService) Register(login, password string) {
	tokenData := TokenData{}
	data := map[string]string{
		"login":    login,
		"password": password,
	}

	srv.post("/users", data, func(r *http.Response) {
		if r.StatusCode == http.StatusConflict {
			srv.ErrorsChan <- "Логин уже занят"
			return
		}
		if r.StatusCode != http.StatusOK {
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

	if user == nil {
		srv.LoginChan <- false
		return
	}

	userID = user.ID
	userIDHeader = fmt.Sprintf("%d", user.ID)

	userData := srv.dbService.GetUserData()
	jwtToken = userData.RefreshToken
	requestUrl := fmt.Sprintf("/users/%d/refresh", userID)

	srv.post(requestUrl, nil, func(r *http.Response) {
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

	srv.put(requestUrl, item, func(r *http.Response) {
		if r.StatusCode == http.StatusConflict {
			srv.ErrorsChan <- "Логин уже занят"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			srv.UserChan <- nil
			return
		}
		srv.saveChan <- user
		srv.UserChan <- &user
	})
}

func (srv *UserService) Login(login, password string) {
	data := map[string]string{
		"login":    login,
		"password": password,
	}
	tokenData := TokenData{}

	srv.post("/login", data, func(r *http.Response) {
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
		srv.LoginChan <- false
		return
	}

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

	srv.get(requestUrl, func(r *http.Response) {
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return
		}
		srv.saveChan <- user
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
