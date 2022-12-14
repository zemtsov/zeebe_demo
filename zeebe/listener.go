package zeebe

import (
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"os"
	"os/signal"
	"syscall"
)

type Listener interface {
	AddHandler(string, worker.JobHandler)
	Listen()
}

type listener struct {
	client   zbc.Client
	handlers []worker.JobWorker
}

func NewListener(client zbc.Client) Listener {
	return &listener{
		client: client,
	}
}

func (l *listener) AddHandler(task string, handler worker.JobHandler) {
	l.handlers = append(l.handlers, l.client.NewJobWorker().JobType(task).Handler(handler).Open())
}

func (l *listener) Listen() {

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		l.clean()
		os.Exit(0)
	}
}

func (l *listener) clean() {
	for _, handler := range l.handlers {
		handler.Close()
		handler.AwaitClose()
	}
}
