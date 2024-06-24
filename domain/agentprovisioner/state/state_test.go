// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	modelconfigstate "github.com/juju/juju/domain/modelconfig/state"
	schematesting "github.com/juju/juju/domain/schema/testing"
	"github.com/juju/juju/environs/config"
)

type suite struct {
	schematesting.ModelSuite
}

var _ = gc.Suite(&suite{})

// TestGetModelConfigKeyValues tests that State.GetModelConfigKeyValues behaves
// as expected:
//   - Requested keys which exist in model config should be returned.
//   - Requested keys which don't exist in model config should not appear in the
//     result, and should not cause an error.
//   - Extra model config keys which are not requested should not be returned.
func (s *suite) TestGetModelConfigKeyValues(c *gc.C) {
	// Set model config in state
	modelConfigState := modelconfigstate.NewState(s.TxnRunnerFactory())
	err := modelConfigState.SetModelConfig(context.Background(), map[string]string{
		config.LXDSnapChannel:                            "5.0/stable",
		config.ContainerImageMetadataURLKey:              "https://images.linuxcontainers.org/",
		config.ContainerImageMetadataDefaultsDisabledKey: "true",
		config.ContainerImageStreamKey:                   "released",
		// Fake keys which will not be requested by the agent provisioner state
		// Hence, they should not show up in the result.
		"key1": "val1",
		"key2": "val2",
	})
	c.Assert(err, jc.ErrorIsNil)

	state := NewState(s.ModelSuite.TxnRunnerFactory())
	modelConfig, err := state.GetModelConfigKeyValues(context.Background(),
		config.LXDSnapChannel,
		config.ContainerImageMetadataURLKey,
		config.ContainerImageMetadataDefaultsDisabledKey,
		config.ContainerImageStreamKey,
		// Fake keys which don't exist in model config, hence they should not
		// show up in the result
		"key3", "key4",
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(modelConfig, jc.DeepEquals, map[string]string{
		config.LXDSnapChannel:                            "5.0/stable",
		config.ContainerImageMetadataURLKey:              "https://images.linuxcontainers.org/",
		config.ContainerImageMetadataDefaultsDisabledKey: "true",
		config.ContainerImageStreamKey:                   "released",
	})
}

// TestGetModelConfigKeyValuesEmptyModelConfig tests that
// State.GetModelConfigKeyValues still works when model config is empty, and
// the sqlair.ErrNoRows is not surfaced.
func (s *suite) TestGetModelConfigKeyValuesEmptyModelConfig(c *gc.C) {
	state := NewState(s.ModelSuite.TxnRunnerFactory())
	modelConfig, err := state.GetModelConfigKeyValues(context.Background(),
		config.LXDSnapChannel,
		config.ContainerImageMetadataURLKey,
		config.ContainerImageMetadataDefaultsDisabledKey,
		config.ContainerImageStreamKey,
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(modelConfig, jc.DeepEquals, map[string]string{})
}

// TestGetModelConfigKeyValuesGetNoKeys tests that if
// State.GetModelConfigKeyValues is called with no requested keys, the
// sqlair.ErrNoRows is not surfaced.
func (s *suite) TestGetModelConfigKeyValuesGetNoKeys(c *gc.C) {
	// Set model config in state
	modelConfigState := modelconfigstate.NewState(s.TxnRunnerFactory())
	err := modelConfigState.SetModelConfig(context.Background(), map[string]string{
		config.LXDSnapChannel:                            "5.0/stable",
		config.ContainerImageMetadataURLKey:              "https://images.linuxcontainers.org/",
		config.ContainerImageMetadataDefaultsDisabledKey: "true",
		config.ContainerImageStreamKey:                   "released",
	})
	c.Assert(err, jc.ErrorIsNil)

	state := NewState(s.ModelSuite.TxnRunnerFactory())
	modelConfig, err := state.GetModelConfigKeyValues(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	c.Check(modelConfig, jc.DeepEquals, map[string]string{})
}
