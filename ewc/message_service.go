package ewc

type MessageService struct {
	BaseService
	ApiService
	dbService       *DbMessageService
	MessageListChan chan *MessageCollection
	MessageChan     chan *Message
}

func NewMessageService() *MessageService {
	srv := new(MessageService)
	srv.dbService = NewDbMessageService(driver, connectionString)
	srv.ErrorsChan = make(chan string, chanSize)
	srv.InfoChan = make(chan string, chanSize)
	srv.MessageChan = make(chan *Message, chanSize)
	srv.MessageListChan = make(chan *MessageCollection, chanSize)

	go srv.listeners()
	return srv
}

func (srv *MessageService) listeners() {

}
