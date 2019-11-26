package ewc

type BasePresenter struct {
	errorsChan chan string
	infoChan   chan string
}
