// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

// modelConfigRow represents a single key-value pair in model config.
type modelConfigRow struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}
