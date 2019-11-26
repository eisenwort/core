package ewc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
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
	srv.dbService = NewDbUserService(driver, connectionString)
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

	claims := JwtClaims{}
	_, err = jwt.ParseWithClaims(tokenData.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	userID = claims.Id
	userIDHeader = fmt.Sprintf("%d", claims.Id)

	srv.getChan <- claims.Id
}

func (srv *UserService) IsLogin() {
	user := srv.dbService.IsLogin()

	if user == nil {
		// srv.UserChan <- nil
		srv.LoginChan <- false
	} else {
		srv.LoginChan <- true
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
	user := User{}

	srv.put(requestUrl, data, func(r *http.Response) {
		if r.StatusCode == http.StatusConflict {
			srv.ErrorsChan <- "Логин уже занят"
			return
		}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			srv.UserChan <- nil
			return
		}
		srv.saveChan <- user
		srv.UserChan <- &user
	})
}

func (srv *UserService) Login(login, password string) {
	data, err := json.Marshal(map[string]string{
		"login":    login,
		"password": password,
	})
	tokenData := TokenData{}

	if err != nil {
		return
	}

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

	srv.LoginChan <- true
	claims := JwtClaims{}
	_, err = jwt.ParseWithClaims(tokenData.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	userID = claims.Id
	userIDHeader = fmt.Sprintf("%d", claims.Id)

	srv.getChan <- claims.Id
}

func (srv *UserService) getUser(id int64) {
	requestUrl := fmt.Sprintf("/users/%d", id)
	user := User{}

	srv.get(requestUrl, func(r *http.Response) {
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return
		}
		srv.saveChan <- user
	})
}

func (srv *UserService) Migrate() {
	srv.dbService.Migrate()
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
