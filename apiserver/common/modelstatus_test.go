// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"context"
	"fmt"

	"github.com/juju/errors"
	"github.com/juju/loggo/v2"
	"github.com/juju/names/v5"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/common/mocks"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/core/arch"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/model"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/envcontext"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/internal/storage/provider"
	dummystorage "github.com/juju/juju/internal/storage/provider/dummy"
	"github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/testing/factory"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	statetesting "github.com/juju/juju/state/testing"
)

type modelStatusSuite struct {
	statetesting.StateSuite

	modelStatusAPI *common.ModelStatusAPI
	resources      *common.Resources
	authorizer     apiservertesting.FakeAuthorizer
}

var _ = gc.Suite(&modelStatusSuite{})

func (s *modelStatusSuite) SetUpTest(c *gc.C) {
	// Initial config needs to be set before the StateSuite SetUpTest.
	s.InitialConfig = testing.CustomModelConfig(c, testing.Attrs{
		"name": "controller",
	})
	s.NewPolicy = func(*state.State) state.Policy {
		return statePolicy{}
	}

	s.StateSuite.SetUpTest(c)
	s.resources = common.NewResources()
	s.AddCleanup(func(_ *gc.C) { s.resources.StopAll() })

	s.authorizer = apiservertesting.FakeAuthorizer{
		Tag:      s.Owner,
		AdminTag: s.Owner,
	}

	s.modelStatusAPI = common.NewModelStatusAPI(
		common.NewModelManagerBackend(state.NoopConfigSchemaSource, s.Model, s.StatePool),
		s.authorizer,
		s.authorizer.GetAuthTag().(names.UserTag),
		func(uuid model.UUID) common.ModelConfigService {
			// For these tests, we just need something of the right type. It
			// doesn't actually have to do anything (yet).
			return mocks.NewMockModelConfigService(nil)
		},
	)

	loggo.GetLogger("juju.apiserver.controller").SetLogLevel(loggo.TRACE)
}

func (s *modelStatusSuite) TestModelStatusNonAuth(c *gc.C) {
	// Set up the user making the call.
	user := names.NewUserTag("username")
	anAuthoriser := apiservertesting.FakeAuthorizer{
		Tag: user,
	}

	api := common.NewModelStatusAPI(
		common.NewModelManagerBackend(state.NoopConfigSchemaSource, s.Model, s.StatePool),
		anAuthoriser,
		anAuthoriser.GetAuthTag().(names.UserTag),
		nil,
	)
	controllerModelTag := s.Model.ModelTag().String()

	req := params.Entities{
		Entities: []params.Entity{{Tag: controllerModelTag}},
	}
	result, err := api.ModelStatus(context.Background(), req)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results[0].Error, gc.ErrorMatches, "permission denied")
}

func (s *modelStatusSuite) TestModelStatusOwnerAllowed(c *gc.C) {
	// Set up the user making the call.
	owner := names.NewUserTag("owner")
	anAuthoriser := apiservertesting.FakeAuthorizer{
		Tag: owner,
	}
	st := s.Factory.MakeModel(c, &factory.ModelParams{Owner: owner})
	defer st.Close()
	api := common.NewModelStatusAPI(
		common.NewModelManagerBackend(state.NoopConfigSchemaSource, s.Model, s.StatePool),
		anAuthoriser,
		anAuthoriser.GetAuthTag().(names.UserTag),
		nil,
	)

	model, err := st.Model()
	c.Assert(err, jc.ErrorIsNil)
	req := params.Entities{
		Entities: []params.Entity{{Tag: model.ModelTag().String()}},
	}
	_, err = api.ModelStatus(context.Background(), req)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *modelStatusSuite) TestModelStatus(c *gc.C) {
	ownerTag := names.NewUserTag("owner")
	otherSt := s.Factory.MakeModel(c, &factory.ModelParams{
		Name:  "dummytoo",
		Owner: ownerTag,
		ConfigAttrs: testing.Attrs{
			"controller": false,
		},
	})
	defer otherSt.Close()

	eight := uint64(8)
	s.Factory.MakeMachine(c, &factory.MachineParams{
		Jobs:            []state.MachineJob{state.JobManageModel},
		Characteristics: &instance.HardwareCharacteristics{CpuCores: &eight},
		InstanceId:      "id-4",
		DisplayName:     "snowflake",
		Volumes: []state.HostVolumeParams{{
			Volume: state.VolumeParams{
				Pool: "modelscoped",
				Size: 123,
			},
		}},
	})
	s.Factory.MakeMachine(c, &factory.MachineParams{
		Jobs:       []state.MachineJob{state.JobHostUnits},
		InstanceId: "id-5",
		Filesystems: []state.HostFilesystemParams{{
			Filesystem: state.FilesystemParams{
				Pool: "modelscoped",
				Size: 123,
			},
		}, {
			Filesystem: state.FilesystemParams{
				Pool: "machinescoped",
				Size: 123,
			},
		}},
	})
	s.Factory.MakeApplication(c, &factory.ApplicationParams{
		Charm: s.Factory.MakeCharm(c, nil),
	})

	otherFactory := factory.NewFactory(otherSt, s.StatePool, testing.FakeControllerConfig())
	otherFactory.MakeMachine(c, &factory.MachineParams{InstanceId: "id-8"})
	otherFactory.MakeMachine(c, &factory.MachineParams{InstanceId: "id-9"})
	otherFactory.MakeApplication(c, &factory.ApplicationParams{
		Charm: otherFactory.MakeCharm(c, nil),
	})

	otherModel, err := otherSt.Model()
	c.Assert(err, jc.ErrorIsNil)

	controllerModelTag := s.Model.ModelTag().String()
	hostedModelTag := otherModel.ModelTag().String()

	req := params.Entities{
		Entities: []params.Entity{{Tag: controllerModelTag}, {Tag: hostedModelTag}},
	}
	results, err := s.modelStatusAPI.ModelStatus(context.Background(), req)
	c.Assert(err, jc.ErrorIsNil)

	arch := arch.DefaultArchitecture
	mem := uint64(64 * 1024 * 1024 * 1024)
	stdHw := &params.MachineHardware{
		Arch: &arch,
		Mem:  &mem,
	}
	c.Assert(results.Results, jc.DeepEquals, []params.ModelStatus{
		{
			ModelTag:           controllerModelTag,
			HostedMachineCount: 1,
			ApplicationCount:   1,
			OwnerTag:           s.Owner.String(),
			Life:               life.Alive,
			Type:               string(state.ModelTypeIAAS),
			Machines: []params.ModelMachineInfo{
				{Id: "0", Hardware: &params.MachineHardware{Cores: &eight}, InstanceId: "id-4", DisplayName: "snowflake", Status: "pending", WantsVote: true},
				{Id: "1", Hardware: stdHw, InstanceId: "id-5", Status: "pending"},
			},
			Applications: []params.ModelApplicationInfo{
				{Name: "mysql"},
			},
			Volumes: []params.ModelVolumeInfo{{
				Id: "0", Status: "pending", Detachable: true,
			}},
			Filesystems: []params.ModelFilesystemInfo{{
				Id: "0", Status: "pending", Detachable: true,
			}, {
				Id: "1/1", Status: "pending", Detachable: false,
			}},
		},
		{
			ModelTag:           hostedModelTag,
			HostedMachineCount: 2,
			ApplicationCount:   1,
			OwnerTag:           ownerTag.String(),
			Life:               life.Alive,
			Type:               string(state.ModelTypeIAAS),
			Machines: []params.ModelMachineInfo{
				{Id: "0", Hardware: stdHw, InstanceId: "id-8", Status: "pending"},
				{Id: "1", Hardware: stdHw, InstanceId: "id-9", Status: "pending"},
			},
			Applications: []params.ModelApplicationInfo{
				{Name: "mysql"},
			},
		},
	})
}

func (s *modelStatusSuite) TestModelStatusCAAS(c *gc.C) {
	ownerTag := names.NewUserTag("owner")
	otherSt := s.Factory.MakeCAASModel(c, &factory.ModelParams{
		Owner: ownerTag,
		ConfigAttrs: testing.Attrs{
			"controller": false,
		},
	})
	defer otherSt.Close()

	otherFactory := factory.NewFactory(otherSt, s.StatePool, testing.FakeControllerConfig())
	app := otherFactory.MakeApplication(c, &factory.ApplicationParams{
		Charm: otherFactory.MakeCharm(c, &factory.CharmParams{Name: "gitlab-k8s", Series: "focal"}),
	})
	otherFactory.MakeUnit(c, &factory.UnitParams{
		Application: app,
	})

	otherModel, err := otherSt.Model()
	c.Assert(err, jc.ErrorIsNil)

	controllerModelTag := s.Model.ModelTag().String()
	hostedModelTag := otherModel.ModelTag().String()

	req := params.Entities{
		Entities: []params.Entity{{Tag: controllerModelTag}, {Tag: hostedModelTag}},
	}
	results, err := s.modelStatusAPI.ModelStatus(context.Background(), req)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(results.Results, jc.DeepEquals, []params.ModelStatus{
		{
			ModelTag:           controllerModelTag,
			HostedMachineCount: 0,
			ApplicationCount:   0,
			OwnerTag:           s.Owner.String(),
			Life:               life.Alive,
			Type:               string(state.ModelTypeIAAS),
			Applications:       []params.ModelApplicationInfo{},
		},
		{
			ModelTag:           hostedModelTag,
			HostedMachineCount: 0,
			ApplicationCount:   1,
			UnitCount:          1,
			OwnerTag:           ownerTag.String(),
			Life:               life.Alive,
			Type:               string(state.ModelTypeCAAS),
			Applications: []params.ModelApplicationInfo{
				{Name: "gitlab"},
			},
		},
	})
}

func (s *modelStatusSuite) TestModelStatusRunsForAllModels(c *gc.C) {
	req := params.Entities{
		Entities: []params.Entity{
			{Tag: "fail.me"},
			{Tag: s.Model.ModelTag().String()},
		},
	}
	expected := params.ModelStatusResults{
		Results: []params.ModelStatus{
			{
				Error: apiservererrors.ServerError(errors.New(`"fail.me" is not a valid tag`))},
			{
				ModelTag: s.Model.ModelTag().String(),
				Life:     life.Value(s.Model.Life().String()),
				OwnerTag: s.Model.Owner().String(),
				Type:     string(state.ModelTypeIAAS),
			},
		},
	}
	result, err := s.modelStatusAPI.ModelStatus(context.Background(), req)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, jc.DeepEquals, expected)
}

type noopStoragePoolGetter struct{}

func (noopStoragePoolGetter) GetStoragePoolByName(_ context.Context, name string) (*storage.Config, error) {
	return nil, fmt.Errorf("storage pool %q not found%w", name, errors.Hide(storageerrors.PoolNotFoundError))
}

type statePolicy struct{}

func (statePolicy) ConfigValidator() (config.Validator, error) {
	return nil, errors.NotImplementedf("ConfigValidator")
}

func (statePolicy) ConstraintsValidator(envcontext.ProviderCallContext) (constraints.Validator, error) {
	return nil, errors.NotImplementedf("ConstraintsValidator")
}

func (statePolicy) StorageServices() (state.StoragePoolGetter, storage.ProviderRegistry, error) {
	return noopStoragePoolGetter{}, storage.ChainedProviderRegistry{
		dummystorage.StorageProviders(),
		provider.CommonStorageProviders(),
	}, nil
}
