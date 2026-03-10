package handler

import "github.com/reflect-homini/stora/internal/provider"

type Handlers struct {
	Auth *AuthHandler
}

func ProvideHandlers(services *provider.Services) *Handlers {
	return &Handlers{
		NewAuthHandler(services.Auth, services.OAuth, services.Session),
	}
}
