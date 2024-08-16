// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package fanconfigurer

import (
	"context"

	"github.com/juju/juju/api/base"
	apiwatcher "github.com/juju/juju/api/watcher"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/rpc/params"
)

// Option is a function that can be used to configure a Client.
type Option = base.Option

// WithTracer returns an Option that configures the Client to use the
// supplied tracer.
var WithTracer = base.WithTracer

// Facade provides access to the FanConfigurer API facade.
type Facade struct {
	caller base.FacadeCaller
}

// TODO(nvinuesa): This entire facade should be removed since it's not used
// anymore (the fanconfigurer worker has been removed).
// NewFacade creates a new client-side FanConfigu	er facade.
func NewFacade(caller base.APICaller, options ...Option) *Facade {
	return &Facade{
		caller: base.NewFacadeCaller(caller, "FanConfigurer", options...),
	}
}

// WatchForFanConfigChanges return a NotifyWatcher waiting for the
// fan configuration to change.
func (f *Facade) WatchForFanConfigChanges(ctx context.Context) (watcher.NotifyWatcher, error) {
	var result params.NotifyWatchResult
	err := f.caller.FacadeCall(ctx, "WatchForFanConfigChanges", nil, &result)
	if err != nil {
		return nil, apiservererrors.RestoreError(err)
	}
	return apiwatcher.NewNotifyWatcher(f.caller.RawAPICaller(), result), nil
}

// FanConfig returns the current fan configuration.
func (f *Facade) FanConfig(ctx context.Context) (network.FanConfig, error) {
	var result params.FanConfigResult
	err := f.caller.FacadeCall(ctx, "FanConfig", nil, &result)
	if err != nil {
		return nil, apiservererrors.RestoreError(err)
	}
	return params.FanConfigResultToFanConfig(result)
}
