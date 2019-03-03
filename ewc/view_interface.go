package ewc

type View interface {
	ShowError(errorString string)
	ShowInfo(message string)
}

type LoginView interface {
	View
	DidGetUser(user *User)
}

type ChatListView interface {
	View
	DidGetChatList(chats *ChatCollection)
}

type ChatView interface {
	View
	DidGetChat(chat *Chat)
	DidDeleteChan(success bool)
	DidClean(success bool)
	DidGetMessage(message *Message)
	DidGetMessageList(messages *MessageCollection)
	DidDeleteMessage(isDeleted bool)
}

type MessageView interface {
	View
	DidGetMessageList(messages *MessageCollection)
}
