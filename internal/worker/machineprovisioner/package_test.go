// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package machineprovisioner_test

import (
	stdtesting "testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package machineprovisioner_test -destination provisioner_mock_test.go github.com/juju/juju/internal/worker/machineprovisioner ControllerAPI,MachinesAPI
//go:generate go run go.uber.org/mock/mockgen -typed -package machineprovisioner_test -destination dependency_mock_test.go github.com/juju/worker/v4/dependency Getter
//go:generate go run go.uber.org/mock/mockgen -typed -package machineprovisioner_test -destination watcher_mock_test.go github.com/juju/juju/core/watcher StringsWatcher
//go:generate go run go.uber.org/mock/mockgen -typed -package machineprovisioner_test -destination base_mock_test.go github.com/juju/juju/api/base APICaller
func TestPackage(t *stdtesting.T) {
	gc.TestingT(t)
}
