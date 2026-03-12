package provider

import (
	"errors"

	"github.com/reflect-homini/stora/internal/core/logger"
)

type Providers struct {
	*DataSources
	*Repositories
	*CoreServices
	*Services
}

func (p *Providers) Shutdown() error {
	var errs error
	if e := p.DataSources.Shutdown(); e != nil {
		errs = errors.Join(errs, e)
	}
	if e := p.CoreServices.Shutdown(); e != nil {
		errs = errors.Join(errs, e)
	}
	return errs
}

func All() (*Providers, error) {
	dataSources, err := ProvideDataSource()
	if err != nil {
		return nil, err
	}

	repos := ProvideRepositories(dataSources.Gorm)

	coreSvcs, err := ProvideCoreServices()
	if err != nil {
		if e := dataSources.Shutdown(); e != nil {
			logger.Error(e)
		}
		return nil, err
	}

	return &Providers{
		DataSources:  dataSources,
		Repositories: repos,
		CoreServices: coreSvcs,
		Services:     ProvideServices(repos, coreSvcs),
	}, nil
}
