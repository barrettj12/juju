// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migrations

import (
	"github.com/juju/juju/core/migration"
	lease "github.com/juju/juju/domain/lease/migrations"
)

// Coordinator is the interface that is used to add operations to a migration.
type Coordinator interface {
	// Add adds the given operation to the migration.
	Add(migration.Operation)
}

// Logger is the interface that is used to log messages.
type Logger interface {
	Infof(string, ...any)
	Debugf(string, ...any)
}

// ImportOperations registers the import operations with the given coordinator.
// This is a convenience function that can be used by the main migration package
// to register all the import operations.
func ImportOperations(coordinator Coordinator, logger Logger) {
	// Note: All the import operations are registered here.
	lease.RegisterImport(coordinator, logger)
}
