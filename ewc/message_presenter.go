package ewc

type MessagePresenter struct {
	BasePresenter
	view           ChatView
	messageService *MessageService
}
