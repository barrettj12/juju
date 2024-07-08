// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"fmt"

	"github.com/canonical/sqlair"
	"github.com/juju/errors"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/domain"
)

// State is responsible for accessing the controller/model DB to retrieve the
// controller/model config keys required for the container config.
type State struct {
	*domain.StateBase
}

// NewState creates a new State object.
func NewState(modelFactory database.TxnRunnerFactory) *State {
	return &State{
		StateBase: domain.NewStateBase(modelFactory),
	}
}

// GetModelConfigKeyValues returns the values of the specified model config
// keys from the model database. If a key cannot be found in model config, it
// will be omitted from the result. If no keys are specified, then this method
// returns an empty map.
func (s *State) GetModelConfigKeyValues(
	ctx context.Context,
	keys ...string,
) (map[string]string, error) {
	if len(keys) == 0 {
		return map[string]string{}, nil
	}

	db, err := s.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	input := make(sqlair.S, 0, len(keys))
	for _, key := range keys {
		input = append(input, key)
	}

	stmt, err := s.Prepare(`
SELECT (key, value) AS (&modelConfigRow.*)
FROM model_config
WHERE key in ($S[:])
`, input, modelConfigRow{})

	if err != nil {
		return nil, fmt.Errorf(
			"preparing get model config key values: %w", domain.CoerceError(err),
		)
	}

	result := make([]modelConfigRow, 0, len(keys))
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt, &input).GetAll(&result)
	})

	if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return nil, fmt.Errorf(
			"getting model config key values: %w",
			domain.CoerceError(err),
		)
	}

	rval := make(map[string]string, len(result))
	for _, row := range result {
		rval[row.Key] = row.Value
	}

	return rval, nil
}
