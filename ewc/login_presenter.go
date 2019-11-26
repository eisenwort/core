package ewc

import (
	"strings"
)

type LoginPresenter struct {
	BasePresenter
	userService *UserService
	view        LoginView
}

func NewLoginPresenter(view LoginView) *LoginPresenter {
	pr := new(LoginPresenter)
	pr.view = view
	pr.errorsChan = make(chan string, chanSize)
	pr.userService = NewUserService()

	go pr.listeners()
	return pr
}

func (pr *LoginPresenter) Login(login, password string) {
	if login == "" || password == "" {
		pr.errorsChan <- "Логин и пароль обязательны"
		return
	}
	go pr.userService.Login(login, password)
}

func (pr *LoginPresenter) IsLogin() {
	go pr.userService.IsLogin()
}

func (pr *LoginPresenter) Register(login, password string) {
	if login == "" || password == "" {
		pr.errorsChan <- "Логин и пароль обязательны"
		return
	}
	go pr.userService.Register(login, password)
}

func (pr *LoginPresenter) SetDbPath(path string) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	dbPath = path
	connectionString = dbPath + dbName
	pr.userService.Migrate()
}

func (pr *LoginPresenter) listeners() {
	for {
		select {
		case user := <-pr.userService.UserChan:
			pr.view.DidGetUser(user)
		case errorString := <-pr.userService.ErrorsChan:
			pr.view.ShowError(errorString)
		case errorString := <-pr.errorsChan:
			pr.view.ShowError(errorString)
		case isSuccess := <-pr.userService.LoginChan:
			pr.view.DidLogin(isSuccess)
		}
	}
}
