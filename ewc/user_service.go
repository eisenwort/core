package ewc

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserService struct {
	ApiService
	saveChan   chan *User
	dbService  *DbUserService
	UserChan   chan *User
	ErrorsChan chan string
}

func NewUserService() *UserService {
	srv := new(UserService)
	srv.dbService = NewDbUserService(driver, connectionString)
	srv.saveChan = make(chan *User, chanSize)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.UserChan = make(chan *User, chanSize)

	go srv.listeners()
	return srv
}

func (srv *UserService) Register(login, password string) {
	user := new(User)
	data, err := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})

	if err != nil {
		log.Println("marshal regiatration data error:", err)
		srv.ErrorsChan <- "Ошибка регистрации"
		return
	}

	srv.post("/users", data, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			errorData := make(map[string]string)

			if err := json.NewDecoder(r.Body).Decode(&errorData); err == nil {
				srv.ErrorsChan <- errorData["error"]
			}

			return
		}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			srv.ErrorsChan <- "Ошибка регистрации"
			return
		}
	})

	if user.ID == 0 {
		return
	}

	srv.UserChan <- user
	userID = user.ID
	userIDHeader = fmt.Sprintf("%d", user.ID)

	srv.saveChan <- user
}

func (srv *UserService) IsLogin() {
	user := srv.dbService.IsLogin()

	if user == nil {
		srv.UserChan <- nil
	} else {
		srv.UserChan <- user
		userID = user.ID
		userIDHeader = fmt.Sprintf("%d", user.ID)
	}
}

func (srv *UserService) Logout() {
	if err := os.Remove(connectionString); err != nil {
		log.Println("delete db file error")
	}
}

func (srv *UserService) Update(item *User) {
	data, _ := json.Marshal(map[string]string{})
	requestUrl := fmt.Sprintf("/users/%d", item.ID)
	user := new(User)

	srv.put(requestUrl, data, func(r *http.Response) {
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			srv.UserChan <- nil
			return
		}
		srv.saveChan <- user
		srv.UserChan <- user
	})
}

func (srv *UserService) Login(login, password string) {
	data, err := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})
	user := new(User)

	if err != nil {
		return
	}

	srv.post("/login", data, func(r *http.Response) {
		if r.StatusCode != http.StatusOK {
			srv.ErrorsChan <- "Неверный логин или пароль"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			srv.ErrorsChan <- "Ошибка входа"
			return
		}
	})
	if user.ID == 0 {
		return
	}
	srv.saveChan <- user
	srv.UserChan <- user
	userID = user.ID
	userIDHeader = fmt.Sprintf("%d", user.ID)
}

func (srv *UserService) Migrate() {
	srv.dbService.Migrate()
}

func (srv *UserService) listeners() {
	for {
		select {
		case user := <-srv.saveChan:
			srv.dbService.Save(user)
		}
	}
}
