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

func (pr *ChatListPresenter) listeners() {
	for {
		select {
		case chatList := <-pr.chatService.ChatListChan:
			pr.view.DidGetChatList(chatList)
		case errorString := <-pr.errorsChan:
			pr.view.ShowError(errorString)
		}
	}
}
