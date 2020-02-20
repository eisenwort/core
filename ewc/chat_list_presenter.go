package ewc

type ChatListPresenter struct {
	BasePresenter
	view        ChatListView
	chatService *ChatService
}

func NewChatListPresenter(view ChatListView) *ChatListPresenter {
	pr := new(ChatListPresenter)
	pr.view = view
	pr.errorsChan = make(chan string, chanSize)
	pr.chatService = NewChatService()

	go pr.listeners()
	return pr
}

func (pr *ChatListPresenter) GetList() {
	go pr.chatService.GetChats()
}

func (pr *ChatListPresenter) CreatePersonalChat(login string) {
	if len(login) == 0 {
		pr.errorsChan <- "Нельзя создать чат без участников"
		return
	}

	pr.chatService.CreatePersonalChat(login)
}

func (pr *ChatListPresenter) CreateChat(friendLogin string) {
	if friendLogin == "" {
		pr.chatService.ErrorsChan <- "Для создания диалога нужно выбрать собеседника"
		return
	}
	go pr.chatService.Create(friendLogin)
}

func (pr *ChatListPresenter) listeners() {
	for {
		select {
		case chatList := <-pr.chatService.ChatListChan:
			pr.view.DidGetChatList(chatList)
		case chat := <-pr.chatService.ChatChan:
			pr.view.DidGetChat(chat)
		case errorString := <-pr.errorsChan:
			pr.view.ShowError(errorString)
		case errorString := <-pr.chatService.ErrorsChan:
			pr.view.ShowError(errorString)
		}
	}
}
