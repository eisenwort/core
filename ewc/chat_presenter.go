package ewc

type ChatPresenter struct {
	BasePresenter
	view           ChatView
	chatService    *ChatService
	messageService *MessageService
}

func NewChatPresenter(view ChatView) *ChatPresenter {
	pr := new(ChatPresenter)
	pr.view = view
	pr.errorsChan = make(chan string, chanSize)
	pr.chatService = NewChatService()
	pr.messageService = NewMessageService()

	go pr.listeners()
	return pr
}

func (pr *ChatPresenter) CreateChat(chat *Chat, friendLogin string) {
	if chat == nil {
		return
	}
	if friendLogin == "" {
		pr.chatService.ErrorsChan <- "Для созлания диалога нужно выбрать собеседника"
		return
	}

	chat.OwnerID = userID
	chat.Personal = true
	go pr.chatService.Create(chat)
}

func (pr *ChatPresenter) CreateGroup(chat *Chat) {
	if chat == nil {
		return
	}
	if chat.Name == "" {
		pr.chatService.ErrorsChan <- "Имя не может быть пустым"
		return
	}
	if len(chat.Users) == 0 {
		pr.chatService.ErrorsChan <- "Необходимо добавить хотя бы одного собеседника"
		return
	}

	chat.OwnerID = userID
	go pr.chatService.Create(chat)
}

func (pr *ChatPresenter) Delete(chat *Chat) {
	if chat == nil {
		return
	}
	if chat.OwnerID != userID {
		pr.chatService.ErrorsChan <- "Вы не можете удалить чат владельцем которого не являетесь"
		return
	}
	go pr.chatService.Delete(chat)
}

func (pr *ChatPresenter) Exit(chat *Chat) {
	if chat == nil {
		return
	}
	if chat.Personal {
		go pr.chatService.Delete(chat)
	} else {
		go pr.chatService.Exit(chat)
	}
}

func (pr *ChatPresenter) Get(id int64, withMessages bool) {
	if id <= 0 {
		return
	}
	go pr.chatService.Get(id, withMessages)
}

func (pr *ChatPresenter) Clean(chat *Chat) {
	if chat == nil {
		return
	}
	if !chat.Personal {
		pr.chatService.ErrorsChan <- "Чистить можно только персональные чаты. Из многопользовательских можно выйти, Ваши сообщение будут удалены"
		return
	}
	go pr.chatService.Clean(chat)
}

func (pr *ChatPresenter) listeners() {
	for {
		select {
		case chat := <-pr.chatService.ChatChan:
			pr.view.DidGetChat(chat)
		case errorString := <-pr.errorsChan:
			pr.view.ShowError(errorString)
		case success := <-pr.chatService.ChatDeleteChan:
			pr.view.DidDeleteChan(success)
		case success := <-pr.chatService.ChatCleanChan:
			pr.view.DidClean(success)
		case message := <-pr.messageService.MessageChan:
			pr.view.DidGetMessage(message)
		case messageList := <-pr.messageService.MessageListChan:
			pr.view.DidGetMessageList(messageList)
		}
	}
}
