package ewc

type View interface {
	ShowError(errorString string)
}

type LoginView interface {
	View
	DidGetUser(user *User)
}
