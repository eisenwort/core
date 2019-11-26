package ewc

type BaseService struct {
	ErrorsChan chan string
	InfoChan   chan string
}
