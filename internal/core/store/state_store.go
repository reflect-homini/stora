package store

import (
	"context"
	"time"

	"github.com/itsLeonB/ungerr"
	"github.com/reflect-homini/stora/internal/adapters/core/store"
	"github.com/reflect-homini/stora/internal/core/config"
)

type StateStore interface {
	Store(ctx context.Context, state string, expiry time.Duration) error
	VerifyAndDelete(ctx context.Context, state string) error
	Shutdown() error
}

func NewStateStore() (StateStore, error) {
	switch config.Global.Auth.StateStore {
	case "inmemory":
		return store.NewInMemoryStateStore(), nil
	case "valkey":
		return store.NewValkeyStateStore(config.Global.Valkey)
	default:
		return nil, ungerr.Unknownf("unimplemented state store: %s", config.Global.Auth.StateStore)
	}
}
