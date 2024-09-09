// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner_test

import (
	"context"
	"fmt"

	"github.com/juju/names/v5"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/facade/facadetest"
	"github.com/juju/juju/apiserver/facades/agent/provisioner"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/environs/tags"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/internal/storage/provider"
	dummystorage "github.com/juju/juju/internal/storage/provider/dummy"
	coretesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/testing/factory"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
)

func (s *withoutControllerSuite) TestProvisioningInfoWithStorage(c *gc.C) {
	registry := storage.ChainedProviderRegistry{
		dummystorage.StorageProviders(),
		provider.CommonStorageProviders(),
	}
	serviceFactoryGetter := s.ServiceFactoryGetter(c)

	st := s.ControllerModel(c).State()
	storageService := serviceFactoryGetter.FactoryForModel(model.UUID(st.ModelUUID())).Storage(registry)
	err := storageService.CreateStoragePool(context.Background(), "static-pool", "static", map[string]any{"foo": "bar"})
	c.Assert(err, jc.ErrorIsNil)

	cons := constraints.MustParse("cores=123 mem=8G")
	template := state.MachineTemplate{
		Base:        state.UbuntuBase("12.10"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
		Constraints: cons,
		Placement:   "valid",
		Volumes: []state.HostVolumeParams{
			{Volume: state.VolumeParams{Size: 1000, Pool: "static-pool"}},
			{Volume: state.VolumeParams{Size: 2000, Pool: "static-pool"}},
		},
	}
	placementMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), template)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: s.machines[0].Tag().String()},
		{Tag: placementMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	serviceFactory := s.ControllerServiceFactory(c)
	controllerCfg, err := serviceFactory.ControllerConfig().ControllerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)

	expected := params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{
			{Result: &params.ProvisioningInfo{
				ControllerConfig: controllerCfg,
				Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
				Jobs:             []model.MachineJob{model.JobHostUnits},
				Tags: map[string]string{
					tags.JujuController: coretesting.ControllerTag.Id(),
					tags.JujuModel:      coretesting.ModelTag.Id(),
					tags.JujuMachine:    "controller-machine-0",
				},
				EndpointBindings: make(map[string]string),
			}},
			{Result: &params.ProvisioningInfo{
				ControllerConfig: controllerCfg,
				Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
				Constraints:      template.Constraints,
				Placement:        template.Placement,
				Jobs:             []model.MachineJob{model.JobHostUnits},
				Tags: map[string]string{
					tags.JujuController: coretesting.ControllerTag.Id(),
					tags.JujuModel:      coretesting.ModelTag.Id(),
					tags.JujuMachine:    "controller-machine-5",
				},
				EndpointBindings: make(map[string]string),
				Volumes: []params.VolumeParams{{
					VolumeTag:  "volume-0",
					Size:       1000,
					Provider:   "static",
					Attributes: map[string]interface{}{"foo": "bar"},
					Tags: map[string]string{
						tags.JujuController: coretesting.ControllerTag.Id(),
						tags.JujuModel:      coretesting.ModelTag.Id(),
					},
					Attachment: &params.VolumeAttachmentParams{
						MachineTag: placementMachine.Tag().String(),
						VolumeTag:  "volume-0",
						Provider:   "static",
					},
				}, {
					VolumeTag:  "volume-1",
					Size:       2000,
					Provider:   "static",
					Attributes: map[string]interface{}{"foo": "bar"},
					Tags: map[string]string{
						tags.JujuController: coretesting.ControllerTag.Id(),
						tags.JujuModel:      coretesting.ModelTag.Id(),
					},
					Attachment: &params.VolumeAttachmentParams{
						MachineTag: placementMachine.Tag().String(),
						VolumeTag:  "volume-1",
						Provider:   "static",
					},
				}},
			}},
		},
	}
	// The order of volumes is not predictable, so we make sure we
	// compare the right ones. This only applies to Results[1] since
	// it is the only result to contain volumes.
	if expected.Results[1].Result.Volumes[0].VolumeTag != result.Results[1].Result.Volumes[0].VolumeTag {
		vols := expected.Results[1].Result.Volumes
		vols[0], vols[1] = vols[1], vols[0]
	}
	c.Assert(result, jc.DeepEquals, expected)
}

func (s *withoutControllerSuite) TestProvisioningInfoRootDiskVolume(c *gc.C) {
	registry := storage.ChainedProviderRegistry{
		dummystorage.StorageProviders(),
		provider.CommonStorageProviders(),
	}
	serviceFactoryGetter := s.ServiceFactoryGetter(c)

	st := s.ControllerModel(c).State()
	storageService := serviceFactoryGetter.FactoryForModel(model.UUID(st.ModelUUID())).Storage(registry)
	err := storageService.CreateStoragePool(context.Background(), "static-pool", "static", map[string]any{"foo": "bar"})
	c.Assert(err, jc.ErrorIsNil)
	template := state.MachineTemplate{
		Base:        state.UbuntuBase("12.10"),
		Constraints: constraints.MustParse("root-disk-source=static-pool"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
	}
	machine, err := st.AddOneMachine(s.InstancePrechecker(c, st), template)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: machine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results[0].Error, gc.IsNil)
	c.Assert(result.Results[0].Result, gc.NotNil)

	c.Assert(result.Results[0].Result.RootDisk, jc.DeepEquals, &params.VolumeParams{
		Provider:   "static",
		Attributes: map[string]interface{}{"foo": "bar"},
	})
}

func (s *withoutControllerSuite) TestProvisioningInfoWithMultiplePositiveSpaceConstraints(c *gc.C) {
	s.addSpacesAndSubnets(c)

	st := s.ControllerModel(c).State()
	cons := constraints.MustParse("cores=123 mem=8G spaces=space1,space2")
	template := state.MachineTemplate{
		Base:        state.UbuntuBase("12.10"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
		Constraints: cons,
		Placement:   "valid",
	}
	placementMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), template)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: placementMachine.Tag().String()},
	}}

	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.IsNil)

	serviceFactory := s.ControllerServiceFactory(c)
	controllerCfg, err := serviceFactory.ControllerConfig().ControllerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)

	expected := &params.ProvisioningInfo{
		ControllerConfig: controllerCfg,
		Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
		Constraints:      template.Constraints,
		Placement:        template.Placement,
		Jobs:             []model.MachineJob{model.JobHostUnits},
		Tags: map[string]string{
			tags.JujuController: coretesting.ControllerTag.Id(),
			tags.JujuModel:      coretesting.ModelTag.Id(),
			tags.JujuMachine:    "controller-machine-5",
		},
		EndpointBindings: make(map[string]string),
		ProvisioningNetworkTopology: params.ProvisioningNetworkTopology{
			SubnetAZs: map[string][]string{
				"subnet-0": {"zone0"},
				"subnet-1": {"zone1"},
				"subnet-2": {"zone2"},
			},
			SpaceSubnets: map[string][]string{
				"space1": {"subnet-0"},
				"space2": {"subnet-1", "subnet-2"},
			},
		},
	}

	res := result.Results[0].Result
	c.Assert(res.SubnetAZs, jc.DeepEquals, expected.SubnetAZs)
	c.Assert(res.SpaceSubnets, gc.HasLen, 2)
	c.Assert(res.SpaceSubnets["space1"], jc.SameContents, expected.SpaceSubnets["space1"])
	c.Assert(res.SpaceSubnets["space2"], jc.SameContents, expected.SpaceSubnets["space2"])
	expected.ProvisioningNetworkTopology = params.ProvisioningNetworkTopology{}
	res.ProvisioningNetworkTopology = params.ProvisioningNetworkTopology{}
	c.Assert(res, jc.DeepEquals, expected)
}

func (s *withoutControllerSuite) TestProvisioningInfoWithEndpointBindings(c *gc.C) {
	spaces := s.addSpacesAndSubnets(c)

	networkService := s.serviceFactory.Network()
	networkService.AddSubnet(context.Background(), network.SubnetInfo{
		CIDR:              "10.10.4.0/24",
		ProviderId:        "subnet-alpha",
		AvailabilityZones: []string{"zone-alpha"},
		SpaceID:           network.AlphaSpaceId,
		VLANTag:           43,
	})

	st := s.ControllerModel(c).State()
	wordpressMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), state.MachineTemplate{
		Base: state.UbuntuBase("12.10"),
		Jobs: []state.MachineJob{state.JobHostUnits},
	})
	c.Assert(err, jc.ErrorIsNil)

	f, release := s.NewFactory(c, s.ControllerModelUUID())
	defer release()

	// Simulates running `juju deploy --bind "..."`.
	bindings := map[string]string{
		"url": spaces[0].ID, // has both name and provider ID
		"db":  spaces[1].ID, // has only name, no provider ID
	}
	wordpressService := f.MakeApplication(c, &factory.ApplicationParams{
		Charm:            f.MakeCharm(c, &factory.CharmParams{Name: "wordpress"}),
		EndpointBindings: bindings,
	})

	wordpressUnit, err := wordpressService.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = wordpressUnit.AssignToMachine(wordpressMachine)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: wordpressMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	serviceFactory := s.ControllerServiceFactory(c)
	controllerCfg, err := serviceFactory.ControllerConfig().ControllerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)

	expected := params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{{
			Result: &params.ProvisioningInfo{
				ControllerConfig: controllerCfg,
				Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
				Jobs:             []model.MachineJob{model.JobHostUnits},
				Tags: map[string]string{
					tags.JujuController:    coretesting.ControllerTag.Id(),
					tags.JujuModel:         coretesting.ModelTag.Id(),
					tags.JujuMachine:       "controller-machine-5",
					tags.JujuUnitsDeployed: wordpressUnit.Name(),
				},
				// Ensure space names are translated to provider IDs, where
				// possible.
				EndpointBindings: map[string]string{
					"":                network.AlphaSpaceName,
					"admin-api":       network.AlphaSpaceName,
					"cache":           network.AlphaSpaceName,
					"db-client":       network.AlphaSpaceName,
					"logging-dir":     network.AlphaSpaceName,
					"monitoring-port": network.AlphaSpaceName,
					"foo-bar":         network.AlphaSpaceName,
					"db":              "space2",         // just name, no provider ID
					"url":             "first space id", // has provider ID
				},
				ProvisioningNetworkTopology: params.ProvisioningNetworkTopology{
					SubnetAZs: map[string][]string{
						"subnet-0":     {"zone0"},
						"subnet-1":     {"zone1"},
						"subnet-2":     {"zone2"},
						"subnet-alpha": {"zone-alpha"},
					},
					SpaceSubnets: map[string][]string{
						"space1": {"subnet-0"},
						"space2": {"subnet-1", "subnet-2"},
						"alpha":  {"subnet-alpha"},
					},
				},
			},
		}},
	}
	check := jc.NewMultiChecker()
	check.Add(`(*.Results[0].Result).ProvisioningNetworkTopology.SpaceSubnets["space2"]`, jc.SameContents, jc.ExpectedValue)
	c.Assert(result, check, expected)
}

func (s *withoutControllerSuite) TestProvisioningInfoWithEndpointBindingsAndNoAlphaSpace(c *gc.C) {
	spaces := s.addSpacesAndSubnets(c)

	st := s.ControllerModel(c).State()
	wordpressMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), state.MachineTemplate{
		Base: state.UbuntuBase("12.10"),
		Jobs: []state.MachineJob{state.JobHostUnits},
	})
	c.Assert(err, jc.ErrorIsNil)

	f, release := s.NewFactory(c, s.ControllerModelUUID())
	defer release()

	// Simulates running `juju deploy --bind "..."`.
	bindings := map[string]string{
		"url": spaces[0].ID, // has both name and provider ID
		"db":  spaces[1].ID, // has only name, no provider ID
	}
	wordpressService := f.MakeApplication(c, &factory.ApplicationParams{
		Charm:            f.MakeCharm(c, &factory.CharmParams{Name: "wordpress"}),
		EndpointBindings: bindings,
	})

	wordpressUnit, err := wordpressService.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = wordpressUnit.AssignToMachine(wordpressMachine)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: wordpressMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	expected := params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{{
			Error: apiservertesting.ServerError(
				"matching subnets to zones: cannot use space \"alpha\" as deployment target: no subnets"),
		}},
	}
	c.Assert(result, jc.DeepEquals, expected)
}

func (s *withoutControllerSuite) TestConflictingNegativeConstraintWithBindingError(c *gc.C) {
	spaces := s.addSpacesAndSubnets(c)

	networkService := s.serviceFactory.Network()
	networkService.AddSubnet(context.Background(), network.SubnetInfo{
		CIDR:              "10.10.4.0/24",
		ProviderId:        "subnet-alpha",
		AvailabilityZones: []string{"zone-alpha"},
		SpaceID:           network.AlphaSpaceId,
		VLANTag:           43,
	})

	st := s.ControllerModel(c).State()
	cons := constraints.MustParse("spaces=^space1")
	wordpressMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), state.MachineTemplate{
		Base:        state.UbuntuBase("12.10"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
		Constraints: cons,
	})
	c.Assert(err, jc.ErrorIsNil)

	f, release := s.NewFactory(c, s.ControllerModelUUID())
	defer release()

	// Simulates running `juju deploy --bind "..."`.
	bindings := map[string]string{
		"url": spaces[0].ID, // has both name and provider ID
		"db":  spaces[1].ID, // has only name, no provider ID
	}
	wordpressService := f.MakeApplication(c, &factory.ApplicationParams{
		Charm:            f.MakeCharm(c, &factory.CharmParams{Name: "wordpress"}),
		EndpointBindings: bindings,
	})
	wordpressUnit, err := wordpressService.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = wordpressUnit.AssignToMachine(wordpressMachine)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: wordpressMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	expected := params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{{
			Error: apiservertesting.ServerError(
				`negative space constraint "space1" conflicts with wordpress endpoint binding for "url"`),
		}},
	}
	c.Assert(result, jc.DeepEquals, expected)
}

func (s *withoutControllerSuite) TestNoSpaceConstraintsProvidedSpaceTopologyEmpty(c *gc.C) {
	st := s.ControllerModel(c).State()
	wordpressMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), state.MachineTemplate{
		Base: state.UbuntuBase("12.10"),
		Jobs: []state.MachineJob{state.JobHostUnits},
	})
	c.Assert(err, jc.ErrorIsNil)

	f, release := s.NewFactory(c, s.ControllerModelUUID())
	defer release()

	// Simulates running `juju deploy --bind "..."`.
	bindings := map[string]string{
		"url": network.AlphaSpaceId,
		"db":  network.AlphaSpaceId,
	}
	wordpressService := f.MakeApplication(c, &factory.ApplicationParams{
		Charm:            f.MakeCharm(c, &factory.CharmParams{Name: "wordpress"}),
		EndpointBindings: bindings,
	})
	wordpressUnit, err := wordpressService.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = wordpressUnit.AssignToMachine(wordpressMachine)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: wordpressMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.IsNil)
	c.Assert(result.Results[0].Result.ProvisioningNetworkTopology.SubnetAZs, gc.IsNil)
	c.Assert(result.Results[0].Result.ProvisioningNetworkTopology.SpaceSubnets, gc.IsNil)
}

func (s *withoutControllerSuite) TestAlphaSpaceConstraintsProvidedExplicitly(c *gc.C) {
	s.addSpacesAndSubnets(c)
	st := s.ControllerModel(c).State()

	networkService := s.serviceFactory.Network()
	networkService.AddSubnet(context.Background(), network.SubnetInfo{
		CIDR:              "10.10.4.0/24",
		ProviderId:        "subnet-alpha",
		AvailabilityZones: []string{"zone-alpha"},
		SpaceID:           network.AlphaSpaceId,
		VLANTag:           43,
	})

	cons := constraints.MustParse("spaces=alpha")
	wordpressMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), state.MachineTemplate{
		Base:        state.UbuntuBase("12.10"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
		Constraints: cons,
	})
	c.Assert(err, jc.ErrorIsNil)

	f, release := s.NewFactory(c, s.ControllerModelUUID())
	defer release()

	// Simulates running `juju deploy --bind "..."`.
	bindings := map[string]string{
		"url": network.AlphaSpaceId,
		"db":  network.AlphaSpaceId,
	}
	wordpressService := f.MakeApplication(c, &factory.ApplicationParams{
		Charm:            f.MakeCharm(c, &factory.CharmParams{Name: "wordpress"}),
		EndpointBindings: bindings,
	})
	wordpressUnit, err := wordpressService.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = wordpressUnit.AssignToMachine(wordpressMachine)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: wordpressMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.IsNil)
	c.Assert(result.Results[0].Result.ProvisioningNetworkTopology.SubnetAZs, gc.DeepEquals, map[string][]string{"subnet-alpha": {"zone-alpha"}})
	c.Assert(result.Results[0].Result.ProvisioningNetworkTopology.SpaceSubnets, gc.DeepEquals, map[string][]string{"alpha": {"subnet-alpha"}})
}
func (s *withoutControllerSuite) addSpacesAndSubnets(c *gc.C) network.SpaceInfos {
	networkService := s.serviceFactory.Network()
	// Add a couple of spaces.
	space1 := network.SpaceInfo{
		Name:       "space1",
		ProviderId: "first space id",
	}
	space2 := network.SpaceInfo{
		Name: "space2",
	}
	sp1ID, err := networkService.AddSpace(context.Background(), space1)
	c.Assert(err, jc.ErrorIsNil)
	sp2ID, err := networkService.AddSpace(context.Background(), space2)
	c.Assert(err, jc.ErrorIsNil)
	// Add 1 subnet into space1, and 2 into space2.
	// Each subnet is in a matching zone (e.g "subnet-#" in "zone#").
	_, err = networkService.AddSubnet(context.Background(), network.SubnetInfo{
		SpaceID:           string(sp1ID),
		CIDR:              "10.0.0.0/24",
		ProviderId:        "subnet-0",
		ProviderNetworkId: "subnet-0",
		VLANTag:           42,
		AvailabilityZones: []string{"zone0"},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = networkService.AddSubnet(context.Background(), network.SubnetInfo{
		SpaceID:           string(sp2ID),
		CIDR:              "10.0.1.0/24",
		ProviderId:        "subnet-1",
		ProviderNetworkId: "subnet-1",
		VLANTag:           42,
		AvailabilityZones: []string{"zone1"},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = networkService.AddSubnet(context.Background(), network.SubnetInfo{
		SpaceID:           string(sp2ID),
		CIDR:              "10.0.2.0/24",
		ProviderId:        "subnet-2",
		ProviderNetworkId: "subnet-2",
		VLANTag:           42,
		AvailabilityZones: []string{"zone2"},
	})
	c.Assert(err, jc.ErrorIsNil)

	return network.SpaceInfos{
		{
			ID:         sp1ID.String(),
			Name:       "space1",
			ProviderId: "first space id",
			Subnets: network.SubnetInfos{
				{
					CIDR:              "10.0.0.0/24",
					ProviderId:        "subnet-0",
					VLANTag:           42,
					AvailabilityZones: []string{"zone0"},
				},
			},
		},
		{
			ID:         sp2ID.String(),
			Name:       "space2",
			ProviderId: "",
			Subnets: network.SubnetInfos{
				{
					CIDR:              "10.0.1.0/24",
					ProviderId:        "subnet-1",
					VLANTag:           42,
					AvailabilityZones: []string{"zone1"},
				},
				{
					CIDR:              "10.0.2.0/24",
					ProviderId:        "subnet-2",
					VLANTag:           42,
					AvailabilityZones: []string{"zone2"},
				},
			},
		},
	}
}

func (s *withoutControllerSuite) TestProvisioningInfoWithUnsuitableSpacesConstraints(c *gc.C) {
	st := s.ControllerModel(c).State()
	// Add an empty space.
	networkService := s.serviceFactory.Network()
	spaceEmpty := network.SpaceInfo{
		Name: "empty",
	}
	_, err := networkService.AddSpace(context.Background(), spaceEmpty)
	c.Assert(err, jc.ErrorIsNil)

	consEmptySpace := constraints.MustParse("cores=123 mem=8G spaces=empty")
	consMissingSpace := constraints.MustParse("cores=123 mem=8G spaces=missing")
	templates := []state.MachineTemplate{{
		Base:        state.UbuntuBase("12.10"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
		Constraints: consEmptySpace,
		Placement:   "valid",
	}, {
		Base:        state.UbuntuBase("12.10"),
		Jobs:        []state.MachineJob{state.JobHostUnits},
		Constraints: consMissingSpace,
		Placement:   "valid",
	}}
	placementMachines, err := st.AddMachines(s.InstancePrechecker(c, st), templates...)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(placementMachines, gc.HasLen, 2)

	args := params.Entities{Entities: []params.Entity{
		{Tag: placementMachines[0].Tag().String()},
		{Tag: placementMachines[1].Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	expectedErrorEmptySpace := `matching subnets to zones: ` +
		`cannot use space "empty" as deployment target: no subnets`
	expectedErrorMissingSpace := `matching subnets to zones: ` +
		`space with name "missing"` // " not found" will be appended by NotFoundError helper below.
	expected := params.ProvisioningInfoResults{Results: []params.ProvisioningInfoResult{
		{Error: apiservertesting.ServerError(expectedErrorEmptySpace)},
		{Error: apiservertesting.NotFoundError(expectedErrorMissingSpace)},
	}}
	c.Assert(result, jc.DeepEquals, expected)
}

func (s *withoutControllerSuite) TestProvisioningInfoWithLXDProfile(c *gc.C) {
	st := s.ControllerModel(c).State()
	profileMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), state.MachineTemplate{
		Base: state.UbuntuBase("12.10"),
		Jobs: []state.MachineJob{state.JobHostUnits},
	})
	c.Assert(err, jc.ErrorIsNil)

	f, release := s.NewFactory(c, s.ControllerModelUUID())
	defer release()
	ch := f.MakeCharm(c, &factory.CharmParams{Name: "lxd-profile"})
	profileService := f.MakeApplication(c, &factory.ApplicationParams{
		Name:  "lxd-profile",
		Charm: ch,
	})
	profileUnit, err := profileService.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = profileUnit.AssignToMachine(profileMachine)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: profileMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	serviceFactory := s.ControllerServiceFactory(c)
	controllerCfg, err := serviceFactory.ControllerConfig().ControllerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)

	mod, err := st.Model()
	c.Assert(err, jc.ErrorIsNil)

	pName := fmt.Sprintf("juju-%s-lxd-profile-%d", mod.Name(), ch.Revision())
	expected := params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{{
			Result: &params.ProvisioningInfo{
				ControllerConfig: controllerCfg,
				Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
				Jobs:             []model.MachineJob{model.JobHostUnits},
				Tags: map[string]string{
					tags.JujuController:    coretesting.ControllerTag.Id(),
					tags.JujuModel:         coretesting.ModelTag.Id(),
					tags.JujuMachine:       "controller-machine-5",
					tags.JujuUnitsDeployed: profileUnit.Name(),
				},
				EndpointBindings: map[string]string{
					"":        network.AlphaSpaceName,
					"another": network.AlphaSpaceName,
					"ubuntu":  network.AlphaSpaceName,
				},
				CharmLXDProfiles: []string{pName},
			},
		}}}
	c.Assert(result, jc.DeepEquals, expected)
}

func (s *withoutControllerSuite) TestStorageProviderFallbackToType(c *gc.C) {
	template := state.MachineTemplate{
		Base:      state.UbuntuBase("12.10"),
		Jobs:      []state.MachineJob{state.JobHostUnits},
		Placement: "valid",
		Volumes: []state.HostVolumeParams{
			{Volume: state.VolumeParams{Size: 1000, Pool: "loop"}},
			{Volume: state.VolumeParams{Size: 1000, Pool: "static"}},
		},
	}
	st := s.ControllerModel(c).State()
	placementMachine, err := st.AddOneMachine(s.InstancePrechecker(c, st), template)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: placementMachine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	serviceFactory := s.ControllerServiceFactory(c)
	controllerCfg, err := serviceFactory.ControllerConfig().ControllerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(result, jc.DeepEquals, params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{
			{Result: &params.ProvisioningInfo{
				ControllerConfig: controllerCfg,
				Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
				Constraints:      template.Constraints,
				Placement:        template.Placement,
				Jobs:             []model.MachineJob{model.JobHostUnits},
				Tags: map[string]string{
					tags.JujuController: coretesting.ControllerTag.Id(),
					tags.JujuModel:      coretesting.ModelTag.Id(),
					tags.JujuMachine:    "controller-machine-5",
				},
				// volume-0 should not be included as it is not managed by
				// the environ provider.
				Volumes: []params.VolumeParams{{
					VolumeTag:  "volume-1",
					Size:       1000,
					Provider:   "static",
					Attributes: nil,
					Tags: map[string]string{
						tags.JujuController: coretesting.ControllerTag.Id(),
						tags.JujuModel:      coretesting.ModelTag.Id(),
					},
					Attachment: &params.VolumeAttachmentParams{
						MachineTag: placementMachine.Tag().String(),
						VolumeTag:  "volume-1",
						Provider:   "static",
					},
				}},
				EndpointBindings: make(map[string]string),
			},
			}},
	})
}

func (s *withoutControllerSuite) TestStorageProviderVolumes(c *gc.C) {
	st := s.ControllerModel(c).State()
	template := state.MachineTemplate{
		Base: state.UbuntuBase("12.10"),
		Jobs: []state.MachineJob{state.JobHostUnits},
		Volumes: []state.HostVolumeParams{
			{Volume: state.VolumeParams{Size: 1000, Pool: "modelscoped"}},
			{Volume: state.VolumeParams{Size: 1000, Pool: "modelscoped"}},
		},
	}
	machine, err := st.AddOneMachine(s.InstancePrechecker(c, st), template)
	c.Assert(err, jc.ErrorIsNil)

	// Provision just one of the volumes, but neither of the attachments.
	sb, err := state.NewStorageBackend(st)
	c.Assert(err, jc.ErrorIsNil)
	err = sb.SetVolumeInfo(names.NewVolumeTag("1"), state.VolumeInfo{
		Pool:       "modelscoped",
		Size:       1000,
		VolumeId:   "vol-ume",
		Persistent: true,
	})
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: machine.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results[0].Error, gc.IsNil)
	c.Assert(result.Results[0].Result, gc.NotNil)

	// volume-0 should be created, as it hasn't yet been provisioned.
	c.Assert(result.Results[0].Result.Volumes, jc.DeepEquals, []params.VolumeParams{{
		VolumeTag: "volume-0",
		Size:      1000,
		Provider:  "modelscoped",
		Tags: map[string]string{
			tags.JujuController: coretesting.ControllerTag.Id(),
			tags.JujuModel:      coretesting.ModelTag.Id(),
		},
		Attachment: &params.VolumeAttachmentParams{
			MachineTag: machine.Tag().String(),
			VolumeTag:  "volume-0",
			Provider:   "modelscoped",
		},
	}})

	// volume-1 has already been provisioned, it just needs to be attached.
	c.Assert(result.Results[0].Result.VolumeAttachments, jc.DeepEquals, []params.VolumeAttachmentParams{{
		MachineTag: machine.Tag().String(),
		VolumeTag:  "volume-1",
		VolumeId:   "vol-ume",
		Provider:   "modelscoped",
	}})
}

func (s *withoutControllerSuite) TestProviderInfoCloudInitUserData(c *gc.C) {
	attrs := map[string]interface{}{"cloudinit-userdata": validCloudInitUserData}
	err := s.serviceFactory.Config().UpdateModelConfig(context.Background(), attrs, nil)
	c.Assert(err, jc.ErrorIsNil)
	template := state.MachineTemplate{
		Base: state.UbuntuBase("12.10"),
		Jobs: []state.MachineJob{state.JobHostUnits},
	}
	st := s.ControllerModel(c).State()
	m, err := st.AddOneMachine(s.InstancePrechecker(c, st), template)
	c.Assert(err, jc.ErrorIsNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: m.Tag().String()},
	}}
	result, err := s.provisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results[0].Result.CloudInitUserData, gc.DeepEquals, map[string]interface{}{
		"packages":        []interface{}{"python-keystoneclient", "python-glanceclient"},
		"preruncmd":       []interface{}{"mkdir /tmp/preruncmd", "mkdir /tmp/preruncmd2"},
		"postruncmd":      []interface{}{"mkdir /tmp/postruncmd", "mkdir /tmp/postruncmd2"},
		"package_upgrade": false})
}

var validCloudInitUserData = `
packages:
  - 'python-keystoneclient'
  - 'python-glanceclient'
preruncmd:
  - mkdir /tmp/preruncmd
  - mkdir /tmp/preruncmd2
postruncmd:
  - mkdir /tmp/postruncmd
  - mkdir /tmp/postruncmd2
package_upgrade: false
`[1:]

func (s *withoutControllerSuite) TestProvisioningInfoPermissions(c *gc.C) {
	serviceFactory := s.ControllerServiceFactory(c)

	// Login as a machine agent for machine 0.
	anAuthorizer := s.authorizer
	anAuthorizer.Controller = false
	anAuthorizer.Tag = s.machines[0].Tag()
	aProvisioner, err := provisioner.MakeProvisionerAPI(context.Background(), facadetest.ModelContext{
		Auth_:           anAuthorizer,
		State_:          s.ControllerModel(c).State(),
		StatePool_:      s.StatePool(),
		Resources_:      s.resources,
		ServiceFactory_: serviceFactory,
		Logger_:         loggertesting.WrapCheckLog(c),
		ControllerUUID_: coretesting.ControllerTag.Id(),
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(aProvisioner, gc.NotNil)

	args := params.Entities{Entities: []params.Entity{
		{Tag: s.machines[0].Tag().String()},
		{Tag: s.machines[0].Tag().String() + "-lxd-0"},
		{Tag: "machine-42"},
		{Tag: s.machines[1].Tag().String()},
		{Tag: "application-bar"},
	}}

	// Only machine 0 and containers therein can be accessed.
	results, err := aProvisioner.ProvisioningInfo(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)

	controllerCfg, err := serviceFactory.ControllerConfig().ControllerConfig(context.Background())
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(results, jc.DeepEquals, params.ProvisioningInfoResults{
		Results: []params.ProvisioningInfoResult{
			{Result: &params.ProvisioningInfo{
				ControllerConfig: controllerCfg,
				Base:             params.Base{Name: "ubuntu", Channel: "12.10/stable"},
				Jobs:             []model.MachineJob{model.JobHostUnits},
				Tags: map[string]string{
					tags.JujuController: coretesting.ControllerTag.Id(),
					tags.JujuModel:      coretesting.ModelTag.Id(),
					tags.JujuMachine:    "controller-machine-0",
				},
				EndpointBindings: make(map[string]string),
			},
			},
			{Error: apiservertesting.NotFoundError("machine 0/lxd/0")},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
		},
	})
}
