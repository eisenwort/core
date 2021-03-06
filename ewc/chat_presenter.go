package ewc

type ChatPresenter struct {
	BasePresenter
	view        ChatView
	chatService *ChatService
}

func NewChatPresenter(view ChatView) *ChatPresenter {
	pr := new(ChatPresenter)
	pr.view = view
	pr.errorsChan = make(chan string, chanSize)
	pr.chatService = NewChatService()

	go pr.listeners()
	return pr
}

/*func (pr *ChatPresenter) CreateGroup(chat *Chat) {
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
}*/

func (pr *ChatPresenter) Delete(chatJson string) {
	chat := Chat{}
	deserialize(chatJson, &chat)

	if chat.ID == 0 {
		return
	}
	if chat.OwnerID != userID {
		pr.chatService.ErrorsChan <- "Вы не можете удалить чат владельцем которого не являетесь"
		return
	}
	go pr.chatService.Delete(chat)
}

func (pr *ChatPresenter) Exit(chatJson string) {
	chat := Chat{}
	deserialize(chatJson, &chat)

	if chat.ID == 0 {
		return
	}
	if chat.Personal {
		go pr.chatService.Delete(chat)
	} else {
		go pr.chatService.Exit(chat)
	}
}

func (pr *ChatPresenter) Get(id int64) {
	if id <= 0 {
		return
	}
	go pr.chatService.Get(id)
}

func (pr *ChatPresenter) GetList() {
	go pr.chatService.GetChats()
}

func (pr *ChatPresenter) Create(userLogin string) {
	go pr.chatService.Create(userLogin)
}

func (pr *ChatPresenter) Clean(chatJson string) {
	chat := Chat{}
	deserialize(chatJson, &chat)

	if chat.ID == 0 {
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
		case chats := <-pr.chatService.ChatListChan:
			pr.view.DidGetChats(chats)
		}
	}
}
