package ewc

type MessagePresenter struct {
	BasePresenter
	view           ChatView
	messageService *MessageService
}

func NewMessagePresenter(view ChatView) *MessagePresenter {
	pr := new(MessagePresenter)
	pr.view = view
	pr.messageService = NewMessageService()

	go pr.listeners()
	return pr
}

func (pr *MessagePresenter) Send(msg *Message) {
	if msg.Text == "" {
		return
	}
	go pr.messageService.Send(msg)
}

func (pr *MessagePresenter) Delete(msg *Message) {
	if msg.UserID != userID {
		pr.messageService.ErrorsChan <- "Невозможно удалить чужое сообщение"
		return
	}
	go pr.Delete(msg)
}

func (pr *MessagePresenter) listeners() {
	for {
		select {
		case msg := <-pr.messageService.MessageChan:
			pr.view.DidGetMessage(msg)
		case isDeleted := <-pr.messageService.MessageDeleteChan:
			pr.view.DidDeleteMessage(isDeleted)
		case errorString := <-pr.messageService.ErrorsChan:
			pr.view.ShowError(errorString)
		case infoString := <-pr.messageService.InfoChan:
			pr.view.ShowInfo(infoString)
		}
	}
}
