package ewc

type ContactsPresenter struct {
	BasePresenter
	view        ContactsView
	userService *UserService
}

func NewContactsPresenter(view ContactsView) *ContactsPresenter {
	pr := new(ContactsPresenter)
	pr.view = view
	pr.errorsChan = make(chan string, chanSize)
	pr.userService = NewUserService()

	go pr.listeners()
	return pr
}

func (pr *ContactsPresenter) GetList() {
	go pr.userService.GetFriends()
}

func (pr *ContactsPresenter) AddFriend(login string) {
	if login == "" {
		pr.errorsChan <- "Логин не может быть пустым"
		return
	}
	go pr.userService.AddFriend(login)
}

func (pr *ContactsPresenter) DeleteFriend(id int64) {
	if id <= 0 {
		pr.errorsChan <- "Неверный id"
		return
	}
	go pr.DeleteFriend(id)
}

func (pr *ContactsPresenter) listeners() {
	for {
		select {
		case errorString := <-pr.errorsChan:
			pr.view.ShowError(errorString)
		case errorString := <-pr.userService.ErrorsChan:
			pr.view.ShowError(errorString)
		case friends := <-pr.userService.UserListChan:
			pr.view.DidGetUserList(friends)
		case friend := <-pr.userService.UserChan:
			pr.view.DidGetUser(friend)
		case id := <-pr.userService.DeleteChan:
			pr.view.DidDeleteUser(id)
		}
	}
}
