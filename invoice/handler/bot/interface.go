package bot

type Bot interface {
	Start(panicChan chan interface{})
	Stop()
}
