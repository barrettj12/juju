// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package bootstrap

import (
	"context"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3/dependency"
	dependencytesting "github.com/juju/worker/v3/dependency/testing"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/objectstore"
	"github.com/juju/juju/domain/servicefactory/testing"
)

type manifoldSuite struct {
	baseSuite
}

var _ = gc.Suite(&manifoldSuite{})

func (s *manifoldSuite) TestValidateConfig(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.getConfig()
	c.Check(cfg.Validate(), jc.ErrorIsNil)

	cfg.AgentName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.StateName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.ObjectStoreName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.ServiceFactoryName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.BootstrapGateName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.CharmhubHTTPClientName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.Logger = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.RequiresBootstrap = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.AgentBinaryUploader = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.ControllerCharmUploader = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

}

func (s *manifoldSuite) getConfig() ManifoldConfig {
	return ManifoldConfig{
		AgentName:              "agent",
		ObjectStoreName:        "object-store",
		StateName:              "state",
		BootstrapGateName:      "bootstrap-gate",
		ServiceFactoryName:     "service-factory",
		CharmhubHTTPClientName: "charmhub-http-client",
		Logger:                 s.logger,
		AgentBinaryUploader: func(context.Context, string, BinaryAgentStorageService, objectstore.ObjectStore, Logger) (func(), error) {
			return func() {}, nil
		},
		ControllerCharmUploader: func(context.Context, string, objectstore.ObjectStore, Logger) error {
			return nil
		},
		RequiresBootstrap: func(context.Context, FlagService) (bool, error) {
			return false, nil
		},
	}
}

func (s *manifoldSuite) getContext() dependency.Context {
	resources := map[string]any{
		"agent":                s.agent,
		"state":                s.stateTracker,
		"object-store":         s.objectStore,
		"bootstrap-gate":       s.bootstrapUnlocker,
		"charmhub-http-client": s.httpClient,
		"service-factory":      testing.NewTestingServiceFactory(),
	}
	return dependencytesting.StubContext(nil, resources)
}

var expectedInputs = []string{"agent", "state", "object-store", "bootstrap-gate", "service-factory", "charmhub-http-client"}

func (s *manifoldSuite) TestInputs(c *gc.C) {
	c.Assert(Manifold(s.getConfig()).Inputs, jc.SameContents, expectedInputs)
}

func (s *manifoldSuite) TestStartAlreadyBootstrapped(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.expectGateUnlock()
	s.expectAgentConfig(c)

	_, err := Manifold(s.getConfig()).Start(s.getContext())
	c.Assert(err, jc.ErrorIs, dependency.ErrUninstall)
}
