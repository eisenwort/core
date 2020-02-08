package ewc

type View interface {
	ShowError(errorString string)
	ShowInfo(message string)
}

type LoginView interface {
	View
	DidLogin(isSuccess bool)
}

type ChatListView interface {
	View
	DidGetChatList(chats string)
	DidGetChat(chat string)
}

type ChatView interface {
	View
	DidGetChats(chats string)
	DidGetChat(chat string)
	DidDeleteChan(success bool)
	DidClean(success bool)
}

type MessageView interface {
	View
	DidGetMessageList(messages string)
	DidGetMessage(messages string)
	DidDeleteMessage(id int64)
}
