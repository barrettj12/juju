// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	"github.com/juju/collections/set"
	"github.com/juju/errors"

	"github.com/juju/juju/core/changestream"
	"github.com/juju/juju/core/providertracker"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/eventsource"
)

// WatchableService provides the API for working with external controllers
// and the ability to create watchers.
type WatchableService struct {
	ProviderService
	watcherFactory WatcherFactory
}

// NewWatchableService returns a new watchable service reference wrapping the
// input state and provider.
func NewWatchableService(st State, provider providertracker.ProviderGetter[Provider], watcherFactory WatcherFactory, logger Logger) *WatchableService {
	return &WatchableService{
		ProviderService: ProviderService{
			Service: Service{
				st:     st,
				logger: logger,
			},
			provider: provider,
		},
		watcherFactory: watcherFactory,
	}
}

// Watch returns a watcher that observes changes to subnets and their
// association (fan underlays), filtered based on the provided list of subnets
// to watch.
func (s *WatchableService) WatchSubnets(ctx context.Context, subnetUUIDsToWatch set.Strings) (watcher.StringsWatcher, error) {
	if s.watcherFactory != nil {
		filter := subnetUUIDsFilter(subnetUUIDsToWatch)

		subnetWatcher, err := s.watcherFactory.NewNamespaceMapperWatcher(
			"subnet",
			changestream.All,
			s.st.AllSubnetsQuery,
			eventsource.FilterEvents(filter),
		)
		if err != nil {
			return nil, errors.Trace(err)
		}
		subnetAssociationWatcher, err := s.watcherFactory.NewNamespaceMapperWatcher(
			"subnet_association",
			changestream.All,
			s.st.AllAssociatedSubnetsQuery,
			eventsource.FilterEvents(filter),
		)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return eventsource.NewMultiStringsWatcher(ctx, subnetWatcher, subnetAssociationWatcher)
	}
	return nil, errors.NotYetAvailablef("subnet watcher")
}

// subnetUUIDsFilter filters the returned subnet UUIDs from the changelog
// according to the user-provided list of subnet UUIDs.
// To keep the compatibility with legacy watchers, if the input set of subnets
// is empty then no filtering is applied.
func subnetUUIDsFilter(subnetUUIDsToWatch set.Strings) func(changestream.ChangeEvent) bool {
	if subnetUUIDsToWatch.IsEmpty() {
		return func(_ changestream.ChangeEvent) bool { return true }
	}
	return func(event changestream.ChangeEvent) bool {
		return subnetUUIDsToWatch.Contains(event.Changed())
	}
}
