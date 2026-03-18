package worker

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/reflect-homini/stora/internal/adapters/worker/scheduler"
	"github.com/reflect-homini/stora/internal/core/logger"
	"github.com/reflect-homini/stora/internal/provider"
)

type Worker struct {
	*scheduler.Scheduler
	shutdownFunc func() error
}

func Setup() (*Worker, error) {
	providers, err := provider.All()
	if err != nil {
		return nil, err
	}

	sched, err := scheduler.Setup(providers)
	if err != nil {
		if e := providers.Shutdown(); e != nil {
			logger.Errorf("error shutting down resources: %v", e)
		}
		return nil, err
	}

	return &Worker{sched, providers.Shutdown}, nil
}

func (w *Worker) Run() {
	logger.Info("starting worker...")
	w.Start()
	logger.Info("worker started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	logger.Info("stopping worker...")
	w.Stop()
	logger.Info("worker stopped")

	if err := w.shutdownFunc(); err != nil {
		logger.Error(err)
	}
}
