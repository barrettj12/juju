// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner_test

import (
	coretesting "github.com/juju/juju/testing"
	stdtesting "testing"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package mocks -destination mocks/service_mock.go github.com/juju/juju/apiserver/facades/agent/provisioner AgentProvisionerService,KeyUpdaterService

func TestPackage(t *stdtesting.T) {
	coretesting.MgoTestPackage(t)
}
