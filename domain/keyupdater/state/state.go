// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"fmt"

	"github.com/canonical/sqlair"
	"github.com/juju/errors"

	"github.com/juju/juju/core/database"
	coremachine "github.com/juju/juju/core/machine"
	"github.com/juju/juju/domain"
	machineerrors "github.com/juju/juju/domain/machine/errors"
)

type State struct {
	*domain.StateBase
}

// AllPublicKeysQuery returns a state SQL query for fetching the public keys
// available on a model. This is useful for constructing authorised keys
// watchers.
func (s *State) AllPublicKeysQuery() string {
	return "SELECT public_key FROM user_public_ssh_key"
}

// AuthorisedKeysForMachine returns a list of authorised public ssh keys for a
// machine name. If no machine exists for the given machine name an error
// satisfying [machineerrors.NotFound] will be returned.
func (s *State) AuthorisedKeysForMachine(
	ctx context.Context,
	name coremachine.Name,
) ([]string, error) {
	db, err := s.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	machineArg := machineName{name.String()}
	machineStmt, err := s.Prepare(`
SELECT name AS &machineName.*
FROM machine
WHERE name = $machineName.name
`, machineArg)
	if err != nil {
		return nil, fmt.Errorf(
			"preparing select statement for getting machine %q when determining authorised keys: %w",
			name, err,
		)
	}

	stmt, err := s.Prepare(`
SELECT public_key AS &authorisedKey.*
FROM user_public_ssh_key
`, authorisedKey{})
	if err != nil {
		return nil, fmt.Errorf(
			"preparing select statement for getting machine %q authorised keys: %w",
			name, err,
		)
	}

	authorisedKeys := []authorisedKey{}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		// Because we have two queries to run and two error paths we need to
		// handle the errors inside this TX so we can produce the correct error.
		err := tx.Query(ctx, machineStmt, machineArg).Get(&machineArg)
		if errors.Is(err, sqlair.ErrNoRows) {
			return fmt.Errorf(
				"cannot get authorised keys for machine %q: %w",
				name, machineerrors.NotFound,
			)
		} else if err != nil {
			return fmt.Errorf(
				"cannot get authorised keys for machine %q: %w",
				name, domain.CoerceError(err),
			)
		}

		err = tx.Query(ctx, stmt).GetAll(&authorisedKeys)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return fmt.Errorf(
				"cannot get authorised keys for machine %q: %w",
				name, domain.CoerceError(err),
			)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	rval := make([]string, 0, len(authorisedKeys))
	for _, authKey := range authorisedKeys {
		rval = append(rval, authKey.PublicKey)
	}

	return rval, nil
}

// AllAuthorisedKeys returns all authorised keys for the model.
func (s *State) AllAuthorisedKeys(ctx context.Context) ([]string, error) {
	db, err := s.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	stmt, err := s.Prepare(`
SELECT public_key AS &authorisedKey.*
FROM user_public_ssh_key
`, authorisedKey{})
	if err != nil {
		return nil, fmt.Errorf(
			"preparing select statement for getting authorised keys: %w", err,
		)
	}

	authorisedKeys := []authorisedKey{}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, stmt).GetAll(&authorisedKeys)
	})

	if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return nil, fmt.Errorf("cannot get authorised keys for model: %w", domain.CoerceError(err))
	}

	rval := make([]string, 0, len(authorisedKeys))
	for _, authKey := range authorisedKeys {
		rval = append(rval, authKey.PublicKey)
	}

	return rval, nil
}

// NewState constructs a new state for interacting with the underlying
// authorised keys of a model.
func NewState(factory database.TxnRunnerFactory) *State {
	return &State{
		StateBase: domain.NewStateBase(factory),
	}
}
