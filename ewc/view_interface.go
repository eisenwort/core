package ewc

type View interface {
	ShowError(errorString string)
	ShowInfo(message string)
}

type LoginView interface {
	View
	DidGetUser(user *User)
	DidLogin(isSuccess bool)
}

type ChatListView interface {
	View
	DidGetChatList(chats *ChatCollection)
	DidGetChat(chat *Chat)
}

type ChatView interface {
	View
	DidGetChat(chat *Chat)
	DidDeleteChan(success bool)
	DidClean(success bool)
	DidGetMessage(message *Message)
	DidGetMessageList(messages *MessageCollection)
	DidDeleteMessage(id int64)
}

type MessageView interface {
	View
	DidGetMessageList(messages *MessageCollection)
}
