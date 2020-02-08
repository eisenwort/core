package ewc

type MessagePresenter struct {
	BasePresenter
	view           MessageView
	messageService *MessageService
}

func NewMessagePresenter(view MessageView) *MessagePresenter {
	pr := new(MessagePresenter)
	pr.view = view
	pr.messageService = NewMessageService()
	pr.errorsChan = make(chan string, chanSize)
	pr.infoChan = make(chan string, chanSize)

	go pr.listeners()
	return pr
}

func (pr *MessagePresenter) Send(msg string, text string) {
	if text == "" {
		pr.errorsChan <- "Невозможно отправить пустое сообщение"
		return
	}
	go pr.messageService.Send(msg, text)
}

func (pr *MessagePresenter) Delete(jsonData string) {
	msg := Message{}
	deserialize(jsonData, &msg)

	if msg.UserID != userID {
		pr.errorsChan <- "Невозможно удалить чужое сообщение"
		return
	}
	go pr.messageService.Delete(msg)
}

func (pr *MessagePresenter) GetByChat(chatID int64, page int) {
	if chatID == 0 {
		return
	}
	if page <= 0 {
		return
	}
	go pr.messageService.GetByChat(chatID, page)
}

func (pr *MessagePresenter) listeners() {
	for {
		select {
		case msg := <-pr.messageService.MessageChan:
			pr.view.DidGetMessage(msg)
		case id := <-pr.messageService.MessageDeleteChan:
			pr.view.DidDeleteMessage(id)
		case errorString := <-pr.messageService.ErrorsChan:
			pr.view.ShowError(errorString)
		case infoString := <-pr.messageService.InfoChan:
			pr.view.ShowInfo(infoString)
		case infoString := <-pr.infoChan:
			pr.view.ShowInfo(infoString)
		case errorString := <-pr.errorsChan:
			pr.view.ShowError(errorString)
		}
	}
}
