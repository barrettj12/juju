// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package providerservicefactory

import (
	"github.com/juju/errors"
	"github.com/juju/loggo/v2"
	"github.com/juju/worker/v4"
	"gopkg.in/tomb.v2"

	"github.com/juju/juju/core/changestream"
	coremodel "github.com/juju/juju/core/model"
	domainservicefactory "github.com/juju/juju/domain/servicefactory"
	"github.com/juju/juju/internal/servicefactory"
)

// Config is the configuration required for service factory worker.
type Config struct {
	// DBGetter supplies WatchableDB implementations by namespace.
	DBGetter changestream.WatchableDBGetter

	Logger Logger

	NewProviderServiceFactoryGetter ProviderServiceFactoryGetterFn
	NewProviderServiceFactory       ProviderServiceFactoryFn
}

// Validate validates the service factory configuration.
func (config Config) Validate() error {
	if config.DBGetter == nil {
		return errors.NotValidf("nil DBGetter")
	}
	if config.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if config.NewProviderServiceFactory == nil {
		return errors.NotValidf("nil NewProviderServiceFactory")
	}
	if config.NewProviderServiceFactoryGetter == nil {
		return errors.NotValidf("nil NewProviderServiceFactoryGetter")
	}
	return nil
}

// NewWorker returns a new service factory worker, with the given configuration.
func NewWorker(config Config) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Trace(err)
	}

	w := &serviceFactoryWorker{
		factoryGetter: config.NewProviderServiceFactoryGetter(
			config.NewProviderServiceFactory,
			config.DBGetter,
			config.Logger,
		),
	}
	w.tomb.Go(func() error {
		<-w.tomb.Dying()
		return w.tomb.Err()
	})
	return w, nil
}

// serviceFactoryWorker is a worker that holds a reference to a service factory.
// This doesn't actually create them dynamically, it just hands them out
// when asked.
type serviceFactoryWorker struct {
	tomb tomb.Tomb

	factoryGetter servicefactory.ProviderServiceFactoryGetter
}

// FactoryGetter returns the provider service factory getter.
func (w *serviceFactoryWorker) FactoryGetter() servicefactory.ProviderServiceFactoryGetter {
	return w.factoryGetter
}

// Kill kills the service factory worker.
func (w *serviceFactoryWorker) Kill() {
	w.tomb.Kill(nil)
}

// Wait waits for the service factory worker to stop.
func (w *serviceFactoryWorker) Wait() error {
	return w.tomb.Wait()
}

// serviceFactory is a provider service factory type.
type serviceFactory struct {
	servicefactory.ProviderServiceFactory
}

// serviceFactoryGetter is a provider service factory getter that returns a
// provider service factory for the given model uuid. This is late binding,
// so the provider service factory is created on demand.
type serviceFactoryGetter struct {
	newProviderServiceFactory ProviderServiceFactoryFn
	dbGetter                  changestream.WatchableDBGetter
	logger                    Logger
}

// FactoryForModel returns a provider service factory for the given model uuid.
// This will late bind the provider service factory to the actual service
// factory.
func (s *serviceFactoryGetter) FactoryForModel(modelUUID string) servicefactory.ProviderServiceFactory {
	return &serviceFactory{
		ProviderServiceFactory: s.newProviderServiceFactory(
			coremodel.UUID(modelUUID), s.dbGetter, s.logger,
		),
	}
}

// The following loggers are required because the Logger interfaces in the other
// locations have a Child method that returns the same Logger type. As this
// a self referential type, we need to wrap it in a new type to satisfy the
// interface. The solution to this is to return a concrete type, but that is
// not idea either, as it means we can't push the testing logger into the
// tests.

// loggoLogger is a loggo.Logger for the service factory.
type loggoLogger struct {
	loggo.Logger
}

// NewLogger returns a new logger for the service factory.
func NewLogger(ns string) Logger {
	return loggoLogger{
		Logger: loggo.GetLogger(ns),
	}
}

// Child returns a child logger that satisfies the Logger interface.
func (c loggoLogger) Child(name string) Logger {
	return c
}

// serviceFactoryLogger is a Logger for the service factory.
type serviceFactoryLogger struct {
	Logger
}

// Child returns a child logger that satisfies the domainservicefactory.Logger.
func (c serviceFactoryLogger) Child(name string) domainservicefactory.Logger {
	return c
}
