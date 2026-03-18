package config

import (
	"errors"

	"github.com/itsLeonB/ungerr"
	"github.com/kelseyhightower/envconfig"
)

type loadable interface {
	Prefix() string
}

type Config struct {
	App            App
	Auth           Auth
	DB             DB
	LLM            LLM
	Mail           Mail
	OTel           OTel
	Valkey         Valkey
	OAuthProviders OAuthProviders
}

var Global *Config

func Load() error {
	var err error

	app, e := load[App]()
	if e != nil {
		err = errors.Join(err, e)
	}

	auth, e := load[Auth]()
	if e != nil {
		err = errors.Join(err, e)
	}

	db, e := load[DB]()
	if e != nil {
		err = errors.Join(err, e)
	}

	llm, e := load[LLM]()
	if e != nil {
		err = errors.Join(err, e)
	}

	mail, e := load[Mail]()
	if e != nil {
		err = errors.Join(err, e)
	}

	otel, e := load[OTel]()
	if e != nil {
		err = errors.Join(err, e)
	}

	valkey, e := load[Valkey]()
	if e != nil {
		err = errors.Join(err, e)
	}

	oAuthProviders, e := loadOAuthProviderConfig()
	if e != nil {
		err = errors.Join(err, e)
	}

	if err != nil {
		return ungerr.Wrap(err, "error loading configs")
	}

	Global = &Config{
		App:            app,
		Auth:           auth,
		DB:             db,
		LLM:            llm,
		Mail:           mail,
		OTel:           otel,
		Valkey:         valkey,
		OAuthProviders: oAuthProviders,
	}

	return nil
}

func load[T loadable]() (T, error) {
	var cfg T
	err := envconfig.Process(cfg.Prefix(), &cfg)
	return cfg, err
}
