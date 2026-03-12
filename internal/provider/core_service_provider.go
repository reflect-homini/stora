package provider

import (
	"github.com/reflect-homini/stora/internal/core/config"
	"github.com/reflect-homini/stora/internal/core/mail"
	"github.com/reflect-homini/stora/internal/core/store"
)

type CoreServices struct {
	Mail  mail.Service
	State store.StateStore
}

func (cs *CoreServices) Shutdown() error {
	return cs.State.Shutdown()
}

func ProvideCoreServices() (*CoreServices, error) {
	store, err := store.NewStateStore()
	if err != nil {
		return nil, err
	}

	return &CoreServices{
		mail.NewMailService(config.Global.Mail),
		store,
	}, nil
}
