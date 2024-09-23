// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"
	"sort"
	"time"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/credential"
	coremodel "github.com/juju/juju/core/model"
	corepermission "github.com/juju/juju/core/permission"
	"github.com/juju/juju/core/user"
	usertesting "github.com/juju/juju/core/user/testing"
	"github.com/juju/juju/domain/access"
	accesserrors "github.com/juju/juju/domain/access/errors"
	modeltesting "github.com/juju/juju/domain/model/state/testing"
	schematesting "github.com/juju/juju/domain/schema/testing"
	loggertesting "github.com/juju/juju/internal/logger/testing"
	"github.com/juju/juju/internal/uuid"
)

type permissionStateSuite struct {
	schematesting.ControllerSuite

	controllerUUID   string
	modelUUID        coremodel.UUID
	defaultModelUUID coremodel.UUID
	debug            bool
}

var _ = gc.Suite(&permissionStateSuite{})

func (s *permissionStateSuite) SetUpTest(c *gc.C) {
	s.ControllerSuite.SetUpTest(c)
	s.controllerUUID = s.SeedControllerUUID(c)

	// Setup to add permissions for user bob on the model

	s.modelUUID = modeltesting.CreateTestModel(c, s.TxnRunnerFactory(), "test-model")
	s.defaultModelUUID = modeltesting.CreateTestModel(c, s.TxnRunnerFactory(), "default-model")

	s.ensureUser(c, "42", "admin", "42", false) // model owner
	s.ensureUser(c, "123", "bob", "42", false)
	s.ensureUser(c, "456", "sue", "42", false)
	s.ensureUser(c, "567", "everyone@external", "42", true)
	s.ensureCloud(c, "987", "test-cloud", "34574", "42")
	s.ensureCloud(c, "654", "another-cloud", "987208634", "42")
}

func (s *permissionStateSuite) TestCreatePermissionModel(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	spec := corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.modelUUID.String(),
				ObjectType: corepermission.Model,
			},
			Access: corepermission.WriteAccess,
		},
	}
	userAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), spec)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(userAccess.UserID, gc.Equals, "123")
	c.Check(userAccess.UserName, gc.Equals, name)
	c.Check(userAccess.Object.ObjectType, gc.Equals, corepermission.Model)
	c.Check(userAccess.Object.Key, gc.Equals, s.modelUUID.String())
	c.Check(userAccess.Access, gc.Equals, corepermission.WriteAccess)
	c.Check(userAccess.DisplayName, gc.Equals, "bob")
	c.Check(userAccess.CreatedBy, gc.Equals, user.AdminUserName)

	s.checkPermissionRow(c, userAccess.UserID, spec)
}

func (s *permissionStateSuite) TestCreatePermissionCloud(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	spec := corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        "test-cloud",
				ObjectType: corepermission.Cloud,
			},
			Access: corepermission.AddModelAccess,
		},
	}
	userAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), spec)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(userAccess.UserID, gc.Equals, "123")
	c.Check(userAccess.UserName, gc.Equals, name)
	c.Check(userAccess.Object.ObjectType, gc.Equals, corepermission.Cloud)
	c.Check(userAccess.Object.Key, gc.Equals, "test-cloud")
	c.Check(userAccess.Access, gc.Equals, corepermission.AddModelAccess)
	c.Check(userAccess.DisplayName, gc.Equals, "bob")
	c.Check(userAccess.CreatedBy, gc.Equals, user.AdminUserName)

	s.checkPermissionRow(c, userAccess.UserID, spec)
}

func (s *permissionStateSuite) TestCreatePermissionController(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	spec := corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.controllerUUID,
				ObjectType: corepermission.Controller,
			},
			Access: corepermission.SuperuserAccess,
		},
	}
	userAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), spec)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(userAccess.UserID, gc.Equals, "123")
	c.Check(userAccess.UserName, gc.Equals, name)
	c.Check(userAccess.Object.ObjectType, gc.Equals, corepermission.Controller)
	c.Check(userAccess.Object.Key, gc.Equals, s.controllerUUID)
	c.Check(userAccess.Access, gc.Equals, corepermission.SuperuserAccess)
	c.Check(userAccess.DisplayName, gc.Equals, "bob")
	c.Check(userAccess.CreatedBy, gc.Equals, user.AdminUserName)

	s.checkPermissionRow(c, userAccess.UserID, spec)
}

func (s *permissionStateSuite) TestCreatePermissionForModelWithBadInfo(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	// model "foo-bar" is not created in this test suite, thus invalid.
	name := usertesting.GenNewName(c, "bob")
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        "foo-bar",
				ObjectType: corepermission.Model,
			},
			Access: corepermission.ReadAccess,
		},
	})
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionTargetInvalid)
}

func (s *permissionStateSuite) TestCreatePermissionForControllerWithBadInfo(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	// The only valid key for an object type of Controller is
	// the controller UUID.
	name := usertesting.GenNewName(c, "bob")
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        "foo-bar",
				ObjectType: corepermission.Controller,
			},
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionTargetInvalid)
}

func (s *permissionStateSuite) checkPermissionRow(c *gc.C, userUUID string, spec corepermission.UserAccessSpec) {
	db := s.DB()

	// Find the permission
	row := db.QueryRow(`
SELECT uuid, access_type, object_type, grant_to, grant_on
FROM   v_permission
WHERE  grant_to = ?
AND    grant_on = ?
`, userUUID, spec.Target.Key)
	c.Assert(row.Err(), jc.ErrorIsNil)
	var (
		accessType, objectType, permUuid, grantTo, grantOn string
	)
	err := row.Scan(&permUuid, &accessType, &objectType, &grantTo, &grantOn)
	c.Assert(err, jc.ErrorIsNil)

	// Verify the permission as expected.
	c.Check(permUuid, gc.Not(gc.Equals), "")
	c.Check(accessType, gc.Equals, string(spec.Access))
	c.Check(objectType, gc.Equals, string(spec.Target.ObjectType))
	c.Check(grantTo, gc.Equals, userUUID)
	c.Check(grantOn, gc.Equals, spec.Target.Key)
}

func (s *permissionStateSuite) TestCreatePermissionErrorNoUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	name := usertesting.GenNewName(c, "testme")
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.modelUUID.String(),
				ObjectType: corepermission.Model,
			},
			Access: corepermission.WriteAccess,
		},
	})
	c.Assert(err, jc.ErrorIs, accesserrors.UserNotFound)
}

func (s *permissionStateSuite) TestCreatePermissionErrorDuplicate(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	spec := corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.modelUUID.String(),
				ObjectType: corepermission.Model,
			},
			Access: corepermission.ReadAccess,
		},
	}
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), spec)
	c.Assert(err, jc.ErrorIsNil)

	// Find the permission
	row := s.DB().QueryRow(`
SELECT uuid, access_type_id, object_type_id, grant_to, grant_on
FROM permission
WHERE access_type_id = 0 AND object_type_id = 2
`)
	c.Assert(row.Err(), jc.ErrorIsNil)

	var (
		userUuid, grantTo, grantOn string
		accessTypeID, objectTypeID int
	)
	err = row.Scan(&userUuid, &accessTypeID, &objectTypeID, &grantTo, &grantOn)
	c.Assert(err, jc.ErrorIsNil)

	// Ensure each combination of grant_on and grant_two
	// is unique
	spec.Access = corepermission.WriteAccess
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), spec)
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionAlreadyExists)
	row2 := s.DB().QueryRow(`
SELECT uuid, access_type_id, object_type_id, grant_to, grant_on
FROM permission
WHERE access_type_id = 1 AND object_type_id = 2
`)
	c.Assert(row2.Err(), jc.ErrorIsNil)
	err = row2.Scan(&userUuid, &accessTypeID, &objectTypeID, &grantTo, &grantOn)
	c.Assert(err, jc.ErrorIs, sql.ErrNoRows)
}

func (s *permissionStateSuite) TestDeletePermission(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	target := corepermission.ID{
		Key:        s.modelUUID.String(),
		ObjectType: corepermission.Model,
	}
	spec := corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.ReadAccess,
		},
	}
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), spec)
	c.Assert(err, jc.ErrorIsNil)

	db := s.DB()
	var numRowBefore int
	err = db.QueryRowContext(context.Background(), "SELECT count(*) FROM permission").Scan(&numRowBefore)
	c.Assert(err, jc.ErrorIsNil)

	err = st.DeletePermission(context.Background(), name, target)
	c.Assert(err, jc.ErrorIsNil)

	// Only one row should be deleted.
	var numRowAfterDelete int
	err = db.QueryRowContext(context.Background(), "SELECT count(*) FROM permission").Scan(&numRowAfterDelete)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(numRowBefore-numRowAfterDelete, gc.Equals, 1)
}

func (s *permissionStateSuite) TestDeletePermissionDoesNotExist(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	target := corepermission.ID{
		Key:        s.modelUUID.String(),
		ObjectType: corepermission.Model,
	}

	// Don't fail if the permission does not exist.
	err := st.DeletePermission(context.Background(), usertesting.GenNewName(c, "bob"), target)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *permissionStateSuite) TestReadUserAccessForTarget(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	createUserAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	var (
		userUuid, grantTo, grantOn string
		accessTypeID, objectTypeID int
	)

	row2 := s.DB().QueryRow(`
SELECT uuid, access_type_id, object_type_id, grant_to, grant_on
FROM permission
WHERE grant_to = 123
`)
	c.Assert(row2.Err(), jc.ErrorIsNil)
	err = row2.Scan(&userUuid, &accessTypeID, &objectTypeID, &grantTo, &grantOn)
	c.Assert(err, jc.ErrorIsNil)
	c.Logf("%q, %d, %d to %q, on %q", userUuid, accessTypeID, objectTypeID, grantTo, grantOn)

	readUserAccess, err := st.ReadUserAccessForTarget(context.Background(), name, target)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(createUserAccess, gc.DeepEquals, readUserAccess)
}

func (s *permissionStateSuite) TestReadUserAccessForTargetExternalUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	jimUserName := usertesting.GenNewName(c, "jim@juju")

	s.ensureUser(c, "777", jimUserName.Name(), "42", true)

	// Add Jim's permissions.
	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	expected, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.LoginAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check Jim's permissions are what we added (everyone@external has no permissions yet).
	userAccess, err := st.ReadUserAccessForTarget(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(userAccess, gc.DeepEquals, expected)

	// Add everyone@external's permissions. These are higher than Jim's.
	everyoneUserAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets the higher level of access given to everyone@external.
	expectedWithEveryoneSuperuser := expected
	expectedWithEveryoneSuperuser.Access = everyoneUserAccess.Access
	expectedWithEveryoneSuperuser.PermissionID = everyoneUserAccess.PermissionID
	userAccess, err = st.ReadUserAccessForTarget(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(userAccess, gc.DeepEquals, expectedWithEveryoneSuperuser)

	// Delete Jim's permissions.
	err = st.DeletePermission(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets everyone@external permissions when he has none
	// himself.
	userAccess, err = st.ReadUserAccessForTarget(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(userAccess, gc.DeepEquals, expectedWithEveryoneSuperuser)
}

func (s *permissionStateSuite) TestReadUserAccessForTargetUserNotFound(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	_, err := st.ReadUserAccessForTarget(context.Background(), usertesting.GenNewName(c, "dave"), target)
	c.Assert(err, jc.ErrorIs, accesserrors.UserNotFound)
}

func (s *permissionStateSuite) TestReadUserAccessForTargetPermissionNotFound(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	// Bob is added in SetUpTest.
	_, err := st.ReadUserAccessForTarget(context.Background(), usertesting.GenNewName(c, "bob"), target)
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionNotFound)
}

func (s *permissionStateSuite) TestReadUserAccessLevelForTarget(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	name := usertesting.GenNewName(c, "bob")
	target := corepermission.ID{
		Key:        "test-cloud",
		ObjectType: corepermission.Cloud,
	}
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: name,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	readUserAccessType, err := st.ReadUserAccessLevelForTarget(context.Background(), name, target)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(readUserAccessType, gc.Equals, corepermission.AddModelAccess)
}

func (s *permissionStateSuite) TestReadUserAccessLevelForTargetExternalUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")
	s.ensureUser(c, "777", jimUserName.Name(), "42", true)

	// Add Jim's permissions.
	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.LoginAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check Jim's permissions are what we added (everyone@external has no permissions yet).
	accessLevel, err := st.ReadUserAccessLevelForTarget(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accessLevel, gc.Equals, corepermission.LoginAccess)

	// Add everyone@external's permissions. These are higher than Jim's.
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets the higher level of access given to everyone@external.
	accessLevel, err = st.ReadUserAccessLevelForTarget(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accessLevel, gc.Equals, corepermission.SuperuserAccess)

	// Delete Jim's permissions.
	err = st.DeletePermission(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets everyone@external permissions when he has none
	// himself.
	accessLevel, err = st.ReadUserAccessLevelForTarget(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accessLevel, gc.DeepEquals, corepermission.SuperuserAccess)
}

func (s *permissionStateSuite) TestEnsureExternalUserIfAuthorized(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")

	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	// Add everyone@external's permissions.
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.LoginAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	err = st.EnsureExternalUserIfAuthorized(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)

	userSt := NewUserState(s.TxnRunnerFactory())
	// Check that jim has now been added as a user with no permissions.
	jim, err := userSt.GetUserByName(context.Background(), jimUserName)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(jim.Name, gc.Equals, jimUserName)
	c.Check(jim.DisplayName, gc.Equals, jimUserName.Name())
	c.Check(jim.UUID, gc.Not(gc.Equals), "")
	c.Check(jim.CreatorUUID, gc.Not(gc.Equals), "")
}

// TestEnsureExternalUserIfAuthorizedNoNewUser checks that no error is returned if the user
// already exists.
func (s *permissionStateSuite) TestEnsureExternalUserIfAuthorizedNoNewUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")
	s.ensureUser(c, "777", jimUserName.Name(), "42", true)

	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}

	err := st.EnsureExternalUserIfAuthorized(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)
}

// TestEnsureExternalUserIfAuthorized checks that no error is returned if the
// user does not exist and does not have access.
func (s *permissionStateSuite) TestEnsureExternalUserIfAuthorizedNoAccess(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")

	target := corepermission.ID{
		Key:        s.controllerUUID,
		ObjectType: corepermission.Controller,
	}
	err := st.EnsureExternalUserIfAuthorized(context.Background(), jimUserName, target)
	c.Assert(err, jc.ErrorIsNil)

	// Check the user has not been added.
	userSt := NewUserState(s.TxnRunnerFactory())
	_, err = userSt.GetUserByName(context.Background(), jimUserName)
	c.Assert(err, jc.ErrorIs, accesserrors.UserNotFound)
}

func (s *permissionStateSuite) TestReadAllUserAccessForUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	s.setupForRead(c, st)

	name := usertesting.GenNewName(c, "bob")
	userAccesses, err := st.ReadAllUserAccessForUser(context.Background(), name)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(userAccesses, gc.HasLen, 4)
	for _, access := range userAccesses {
		c.Assert(access.UserName, gc.Equals, name)
		c.Assert(access.CreatedBy, gc.Equals, user.AdminUserName)
	}
}

func (s *permissionStateSuite) TestReadAllUserAccessForUserExternalUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")
	s.ensureUser(c, "777", jimUserName.Name(), "42", true)

	// Add Jim's permissions.
	controllerTarget := corepermission.ID{
		ObjectType: corepermission.Controller,
		Key:        s.controllerUUID,
	}
	jimsControllerAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: controllerTarget,
			Access: corepermission.LoginAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	modelTarget := corepermission.ID{
		ObjectType: corepermission.Model,
		Key:        s.modelUUID.String(),
	}
	jimsModelAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: modelTarget,
			Access: corepermission.ReadAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check Jim's permissions are what we added (everyone@external has no permissions yet).
	expected := []corepermission.UserAccess{jimsModelAccess, jimsControllerAccess}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].PermissionID > expected[j].PermissionID
	})

	accesses, err := st.ReadAllUserAccessForUser(context.Background(), jimUserName)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expected)

	// Add everyone@external's permissions. These are higher than Jim's.
	everyoneController, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: controllerTarget,
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	everyoneModel, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: modelTarget,
			Access: corepermission.AdminAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets the higher level of access given to everyone@external.
	expectedWithExternal := []corepermission.UserAccess{jimsModelAccess, jimsControllerAccess}
	expectedWithExternal[0].Access = everyoneModel.Access
	expectedWithExternal[0].PermissionID = everyoneModel.PermissionID
	expectedWithExternal[1].Access = everyoneController.Access
	expectedWithExternal[1].PermissionID = everyoneController.PermissionID
	sort.Slice(expectedWithExternal, func(i, j int) bool {
		return expectedWithExternal[i].PermissionID > expectedWithExternal[j].PermissionID
	})

	accesses, err = st.ReadAllUserAccessForUser(context.Background(), jimUserName)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expectedWithExternal)

	// Delete one of Jim's permissions (the controller permission).
	err = st.DeletePermission(context.Background(), jimUserName, controllerTarget)
	c.Assert(err, jc.ErrorIsNil)
	err = st.DeletePermission(context.Background(), jimUserName, modelTarget)
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets everyone@external permissions when he has none
	// himself.
	accesses, err = st.ReadAllUserAccessForUser(context.Background(), jimUserName)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expectedWithExternal)
}

func (s *permissionStateSuite) TestReadAllUserAccessForUserUserNotFound(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	_, err := st.ReadAllUserAccessForUser(context.Background(), usertesting.GenNewName(c, "dave"))
	c.Assert(err, jc.ErrorIs, accesserrors.UserNotFound)
}

func (s *permissionStateSuite) TestReadAllUserAccessPermissionNotFound(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	s.ensureUser(c, "777", "dave", "42", true)

	_, err := st.ReadAllUserAccessForUser(context.Background(), usertesting.GenNewName(c, "dave"))
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionNotFound)
}

func (s *permissionStateSuite) TestReadAllUserAccessForTarget(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	s.setupForRead(c, st)
	targetCloud := corepermission.ID{
		Key:        "test-cloud",
		ObjectType: corepermission.Cloud,
	}
	userAccesses, err := st.ReadAllUserAccessForTarget(context.Background(), targetCloud)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(userAccesses, gc.HasLen, 2)
	accessZero := userAccesses[0]
	c.Check(accessZero.Access, gc.Equals, corepermission.AddModelAccess)
	c.Check(accessZero.Object, gc.Equals, targetCloud)
	accessOne := userAccesses[1]
	c.Check(accessOne.Access, gc.Equals, corepermission.AddModelAccess)
	c.Check(accessOne.Object, gc.Equals, targetCloud)

	c.Check(accessZero.UserID, gc.Not(gc.Equals), accessOne.UserID)
}

func (s *permissionStateSuite) TestReadAllUserAccessForTargetExternalUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")
	johnUserName := usertesting.GenNewName(c, "john@juju")
	s.ensureUser(c, "777", jimUserName.Name(), "42", true)
	s.ensureUser(c, "888", johnUserName.Name(), "42", true)

	// Add Jim and John's permissions.
	cloudTarget := corepermission.ID{
		ObjectType: corepermission.Cloud,
		Key:        "test-cloud",
	}
	jimsCloudAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTarget,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	johnsCloudAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: johnUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTarget,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check Jim's permissions are what we added (everyone@external has no permissions yet).
	expected := []corepermission.UserAccess{jimsCloudAccess, johnsCloudAccess}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].PermissionID > expected[j].PermissionID
	})

	accesses, err := st.ReadAllUserAccessForTarget(context.Background(), cloudTarget)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expected)

	// Add everyone@external's permissions. These are higher than Jim's, but not Johns.
	everyoneCloud, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTarget,
			Access: corepermission.AdminAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets the higher level of access given to everyone@external.
	expectedWithExternal := []corepermission.UserAccess{jimsCloudAccess, johnsCloudAccess}
	expectedWithExternal[0].Access = everyoneCloud.Access
	expectedWithExternal[0].PermissionID = everyoneCloud.PermissionID
	expectedWithExternal[1].Access = everyoneCloud.Access
	expectedWithExternal[1].PermissionID = everyoneCloud.PermissionID
	sort.Slice(expectedWithExternal, func(i, j int) bool {
		return expectedWithExternal[i].PermissionID > expectedWithExternal[j].PermissionID
	})

	accesses, err = st.ReadAllUserAccessForTarget(context.Background(), cloudTarget)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expectedWithExternal)

	// Delete Jim and John's own permission.
	err = st.DeletePermission(context.Background(), johnUserName, cloudTarget)
	c.Assert(err, jc.ErrorIsNil)
	err = st.DeletePermission(context.Background(), jimUserName, cloudTarget)
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim and John get everyone@external permissions when they have
	// none themselves.
	accesses, err = st.ReadAllUserAccessForTarget(context.Background(), cloudTarget)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expectedWithExternal)
}

func (s *permissionStateSuite) TestReadAllAccessForUserAndObjectTypeCloud(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	s.setupForRead(c, st)

	name := usertesting.GenNewName(c, "bob")
	users, err := st.ReadAllAccessForUserAndObjectType(context.Background(), name, corepermission.Cloud)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(users, gc.HasLen, 2)

	var foundTestCloud, foundAnotherCloud bool
	for _, userAccess := range users {
		c.Check(userAccess.UserName, gc.Equals, name)
		c.Check(userAccess.CreatedBy, gc.Equals, user.AdminUserName)
		c.Check(userAccess.UserID, gc.Equals, "123")
		c.Check(userAccess.Access, gc.Equals, corepermission.AddModelAccess)
		c.Check(userAccess.Object.ObjectType, gc.Equals, corepermission.Cloud)
		if userAccess.Object.Key == "test-cloud" {
			foundTestCloud = true
		}
		if userAccess.Object.Key == "another-cloud" {
			foundAnotherCloud = true
		}
	}
	c.Check(foundTestCloud && foundAnotherCloud, jc.IsTrue, gc.Commentf("%+v", users))
}

func (s *permissionStateSuite) TestReadAllAccessForUserAndObjectTypeModel(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	s.setupForRead(c, st)

	name := usertesting.GenNewName(c, "bob")
	users, err := st.ReadAllAccessForUserAndObjectType(context.Background(), name, corepermission.Model)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(users, gc.HasLen, 2)

	var admin, write bool
	for _, userAccess := range users {
		c.Check(userAccess.UserName, gc.Equals, name)
		c.Check(userAccess.CreatedBy, gc.Equals, user.AdminUserName)
		c.Check(userAccess.UserID, gc.Equals, "123")
		c.Check(userAccess.Object.ObjectType, gc.Equals, corepermission.Model)
		if userAccess.Access == corepermission.WriteAccess {
			write = true
			c.Check(userAccess.Object.Key, gc.Equals, s.defaultModelUUID.String())
		}
		if userAccess.Access == corepermission.AdminAccess {
			admin = true
			c.Check(userAccess.Object.Key, gc.Equals, s.modelUUID.String())
		}
	}
	c.Assert(admin && write, jc.IsTrue, gc.Commentf("%+v", users))
}

func (s *permissionStateSuite) TestReadAllAccessForUserAndObjectTypeController(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	s.setupForRead(c, st)

	name := usertesting.GenNewName(c, "admin")
	users, err := st.ReadAllAccessForUserAndObjectType(context.Background(), name, corepermission.Controller)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(users, gc.HasLen, 1)
	userAccess := users[0]
	c.Check(userAccess.UserName, gc.Equals, name, gc.Commentf("%+v", users))
	c.Check(userAccess.CreatedBy, gc.Equals, user.AdminUserName, gc.Commentf("%+v", users))
	c.Check(userAccess.UserID, gc.Equals, "42", gc.Commentf("%+v", users))
	c.Check(userAccess.Access, gc.Equals, corepermission.SuperuserAccess, gc.Commentf("%+v", users))
}

func (s *permissionStateSuite) TestReadAllAccessForUserAndObjectTypeNotFound(c *gc.C) {
	st := NewState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	jimUserName := usertesting.GenNewName(c, "jim@juju")
	s.ensureUser(c, "777", jimUserName.Name(), "42", true)

	_, err := st.ReadAllAccessForUserAndObjectType(context.Background(), usertesting.GenNewName(c, "bob"), corepermission.Cloud)
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionNotFound)

	_, err = st.ReadAllAccessForUserAndObjectType(context.Background(), jimUserName, corepermission.Cloud)
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionNotFound)
}

func (s *permissionStateSuite) TestReadAllAccessForUserAndObjectTypeExternalUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))

	jimUserName := usertesting.GenNewName(c, "jim@juju")
	s.ensureUser(c, "777", jimUserName.Name(), "42", true)

	// Add Jim's permissions.
	cloudTargetOne := corepermission.ID{
		ObjectType: corepermission.Cloud,
		Key:        "test-cloud",
	}
	cloudTargetTwo := corepermission.ID{
		ObjectType: corepermission.Cloud,
		Key:        "another-cloud",
	}
	jimsCloudOneAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTargetOne,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, gc.IsNil)
	jimsCloudTwoAccess, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: jimUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTargetTwo,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check Jim's permissions are what we added (everyone@external has no permissions yet).
	expected := []corepermission.UserAccess{jimsCloudOneAccess, jimsCloudTwoAccess}
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].PermissionID > expected[j].PermissionID
	})

	accesses, err := st.ReadAllAccessForUserAndObjectType(context.Background(), jimUserName, corepermission.Cloud)
	c.Assert(err, jc.ErrorIsNil)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})

	// Add everyone@external's permissions. These are higher than Jim's, but not Johns.
	everyoneCloudOne, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTargetOne,
			Access: corepermission.AdminAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	everyoneCloudTwo, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: corepermission.EveryoneUserName,
		AccessSpec: corepermission.AccessSpec{
			Target: cloudTargetTwo,
			Access: corepermission.AdminAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim gets the higher level of access given to everyone@external.
	expectedWithExternal := []corepermission.UserAccess{jimsCloudOneAccess, jimsCloudTwoAccess}
	expectedWithExternal[0].Access = everyoneCloudOne.Access
	expectedWithExternal[0].PermissionID = everyoneCloudOne.PermissionID
	expectedWithExternal[1].Access = everyoneCloudTwo.Access
	expectedWithExternal[1].PermissionID = everyoneCloudTwo.PermissionID
	sort.Slice(expectedWithExternal, func(i, j int) bool {
		return expectedWithExternal[i].PermissionID > expectedWithExternal[j].PermissionID
	})

	accesses, err = st.ReadAllAccessForUserAndObjectType(context.Background(), jimUserName, corepermission.Cloud)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expectedWithExternal)

	// Delete Jim's permissions
	err = st.DeletePermission(context.Background(), jimUserName, cloudTargetOne)
	c.Assert(err, jc.ErrorIsNil)
	err = st.DeletePermission(context.Background(), jimUserName, cloudTargetTwo)
	c.Assert(err, jc.ErrorIsNil)

	// Check that Jim and John get everyone@external permissions when they have
	// none themselves.
	accesses, err = st.ReadAllAccessForUserAndObjectType(context.Background(), jimUserName, corepermission.Cloud)
	sort.Slice(accesses, func(i, j int) bool {
		return accesses[i].PermissionID > accesses[j].PermissionID
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(accesses, gc.DeepEquals, expectedWithExternal)
}

func (s *permissionStateSuite) TestUpdatePermissionGrantNewExternalUser(c *gc.C) {
	st := NewState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	adminName := usertesting.GenNewName(c, "admin")
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: adminName,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.modelUUID.String(),
				ObjectType: corepermission.Model,
			},
			Access: corepermission.AdminAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: adminName,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.controllerUUID,
				ObjectType: corepermission.Controller,
			},
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	target := corepermission.ID{
		ObjectType: corepermission.Model,
		Key:        s.modelUUID.String(),
	}
	tomName := usertesting.GenNewName(c, "tom@external")
	arg := access.UpdatePermissionArgs{
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.WriteAccess,
		},
		Change:  corepermission.Grant,
		Subject: tomName,
	}
	err = st.UpdatePermission(context.Background(), arg)
	c.Assert(err, jc.ErrorIsNil)

	obtainedUserAccess, err := st.ReadUserAccessForTarget(context.Background(), tomName, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(obtainedUserAccess.UserName, gc.Equals, tomName)
	c.Check(obtainedUserAccess.CreatedBy, gc.Equals, corepermission.EveryoneUserName)
	c.Check(obtainedUserAccess.UserID, gc.Not(gc.Equals), "")
	c.Check(obtainedUserAccess.Access, gc.Equals, corepermission.WriteAccess)
	c.Check(obtainedUserAccess.Object.ObjectType, gc.Equals, corepermission.Model)
	c.Check(obtainedUserAccess.Object.Key, gc.Equals, s.modelUUID.String())
}

func (s *permissionStateSuite) TestUpdatePermissionGrantExistingUser(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	// Bob starts with Write access on "default-model"
	s.setupForRead(c, st)

	name := usertesting.GenNewName(c, "bob")
	target := corepermission.ID{
		ObjectType: corepermission.Model,
		Key:        s.defaultModelUUID.String(),
	}
	arg := access.UpdatePermissionArgs{
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.AdminAccess,
		},
		Change:  corepermission.Grant,
		Subject: name,
	}
	err := st.UpdatePermission(context.Background(), arg)
	c.Assert(err, jc.ErrorIsNil)

	obtainedUserAccess, err := st.ReadUserAccessForTarget(context.Background(), name, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(obtainedUserAccess.UserName, gc.Equals, name)
	c.Check(obtainedUserAccess.Access, gc.Equals, corepermission.AdminAccess)
}

func (s *permissionStateSuite) TestUpdatePermissionGrantLessAccess(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	// Bob starts with Write access on "default-model"
	s.setupForRead(c, st)

	name := usertesting.GenNewName(c, "bob")
	target := corepermission.ID{
		ObjectType: corepermission.Model,
		Key:        s.modelUUID.String(),
	}
	arg := access.UpdatePermissionArgs{
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.ReadAccess,
		},
		Change:  corepermission.Grant,
		Subject: name,
	}
	err := st.UpdatePermission(context.Background(), arg)
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionAccessGreater)
}

func (s *permissionStateSuite) TestUpdatePermissionRevokeRemovePerm(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	s.setupForRead(c, st)
	// Bob starts with Admin access on "default-model".
	// Revoke of Read yields permission removed on the model.
	target := corepermission.ID{
		ObjectType: corepermission.Model,
		Key:        s.defaultModelUUID.String(),
	}
	name := usertesting.GenNewName(c, "bob")
	arg := access.UpdatePermissionArgs{
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.ReadAccess,
		},
		Change:  corepermission.Revoke,
		Subject: name,
	}
	err := st.UpdatePermission(context.Background(), arg)
	c.Assert(err, jc.ErrorIsNil)

	_, err = st.ReadUserAccessForTarget(context.Background(), name, target)
	c.Assert(err, jc.ErrorIs, accesserrors.PermissionNotFound)
}

func (s *permissionStateSuite) TestUpdatePermissionRevoke(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	// Sue starts with Admin access on "test-cloud".
	// Revoke of Admin yields AddModel on clouds.
	s.setupForRead(c, st)

	target := corepermission.ID{
		ObjectType: corepermission.Cloud,
		Key:        "test-cloud",
	}
	name := usertesting.GenNewName(c, "sue")
	arg := access.UpdatePermissionArgs{
		AccessSpec: corepermission.AccessSpec{
			Target: target,
			Access: corepermission.AdminAccess,
		},
		Change:  corepermission.Revoke,
		Subject: name,
	}
	err := st.UpdatePermission(context.Background(), arg)
	c.Assert(err, jc.ErrorIsNil)

	obtainedUserAccess, err := st.ReadUserAccessForTarget(context.Background(), name, target)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(obtainedUserAccess.UserName, gc.Equals, name)
	c.Check(obtainedUserAccess.Access, gc.Equals, corepermission.AddModelAccess)
}

func (s *permissionStateSuite) TestModelAccessForCloudCredential(c *gc.C) {
	st := NewPermissionState(s.TxnRunnerFactory(), loggertesting.WrapCheckLog(c))
	ctx := context.Background()

	modeltesting.CreateTestModel(c, s.TxnRunnerFactory(), "model-access")
	key := credential.Key{
		Cloud: "model-access",
		Owner: usertesting.GenNewName(c, "test-usermodel-access"),
		Name:  "foobar",
	}

	obtained, err := st.AllModelAccessForCloudCredential(ctx, key)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(obtained, gc.HasLen, 1)
	c.Check(obtained[0].ModelName, gc.DeepEquals, "model-access")
	c.Check(obtained[0].OwnerAccess, gc.DeepEquals, corepermission.AdminAccess)
}

func (s *permissionStateSuite) setupForRead(c *gc.C, st *PermissionState) {
	targetCloud := corepermission.ID{
		Key:        "test-cloud",
		ObjectType: corepermission.Cloud,
	}
	bobName := usertesting.GenNewName(c, "bob")
	sueName := usertesting.GenNewName(c, "sue")
	adminName := usertesting.GenNewName(c, "admin")
	_, err := st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: bobName,
		AccessSpec: corepermission.AccessSpec{
			Target: targetCloud,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: sueName,
		AccessSpec: corepermission.AccessSpec{
			Target: targetCloud,
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: bobName,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.modelUUID.String(),
				ObjectType: corepermission.Model,
			},
			Access: corepermission.AdminAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: bobName,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        "another-cloud",
				ObjectType: corepermission.Cloud,
			},
			Access: corepermission.AddModelAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: bobName,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.defaultModelUUID.String(),
				ObjectType: corepermission.Model,
			},
			Access: corepermission.WriteAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	_, err = st.CreatePermission(context.Background(), uuid.MustNewUUID(), corepermission.UserAccessSpec{
		User: adminName,
		AccessSpec: corepermission.AccessSpec{
			Target: corepermission.ID{
				Key:        s.controllerUUID,
				ObjectType: corepermission.Controller,
			},
			Access: corepermission.SuperuserAccess,
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	if s.debug {
		s.printUsers(c)
		s.printUserAuthentication(c)
		s.printClouds(c)
		s.printPermissions(c)
		s.printRead(c)
		s.printModels(c)
	}
}

func (s *permissionStateSuite) ensureUser(c *gc.C, userUUID, name, createdByUUID string, external bool) {
	err := s.TxnRunner().StdTxn(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO user (uuid, name, display_name, external, removed, created_by_uuid, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, userUUID, name, name, external, false, createdByUUID, time.Now())
		return err
	})
	c.Assert(err, jc.ErrorIsNil)
	err = s.TxnRunner().StdTxn(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO user_authentication (user_uuid, disabled)
			VALUES (?, ?)
		`, userUUID, false)
		return err
	})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *permissionStateSuite) ensureCloud(c *gc.C, cloudUUID, cloudName, credUUID, ownerUUID string) {
	err := s.TxnRunner().StdTxn(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO cloud (uuid, name, cloud_type_id, endpoint, skip_tls_verify)
			VALUES (?, ?, 7, "test-endpoint", true)
		`, cloudUUID, cloudName)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO cloud_auth_type (cloud_uuid, auth_type_id)
			VALUES (?, 0), (?, 2)
		`, cloudUUID, cloudUUID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO cloud_credential (uuid, cloud_uuid, auth_type_id, owner_uuid, name, revoked, invalid)
			VALUES (?, ?, ?, ?, "foobar", false, false)
		`, credUUID, cloudUUID, 0, ownerUUID)
		return err
	})

	c.Assert(err, jc.ErrorIsNil)
}

func (s *permissionStateSuite) printPermissions(c *gc.C) {
	rows, _ := s.DB().Query(`
SELECT uuid, access_type_id, object_type_id, grant_to, grant_on
FROM permission
`)
	defer func() { _ = rows.Close() }()
	var (
		userUuid, grantTo, grantOn string
		accessTypeID, objectTypeID int
	)

	c.Logf("PERMISSIONS")
	for rows.Next() {
		err := rows.Scan(&userUuid, &accessTypeID, &objectTypeID, &grantTo, &grantOn)
		c.Assert(err, jc.ErrorIsNil)
		c.Logf("%q, %d, %d, %q, %q", userUuid, accessTypeID, objectTypeID, grantTo, grantOn)
	}
}

func (s *permissionStateSuite) printUsers(c *gc.C) {
	rows, _ := s.DB().Query(`
SELECT u.uuid, u.name, u.created_by_uuid,  u.removed
FROM v_user_auth u
`)
	defer func() { _ = rows.Close() }()
	var (
		rowUUID, name string
		creatorUUID   user.UUID
		removed       bool
	)
	c.Logf("USERS")
	for rows.Next() {
		err := rows.Scan(&rowUUID, &name, &creatorUUID, &removed)
		c.Assert(err, jc.ErrorIsNil)
		c.Logf("LINE %q, %q, %q,  %t", rowUUID, name, creatorUUID, removed)
	}
}

func (s *permissionStateSuite) printUserAuthentication(c *gc.C) {
	rows, _ := s.DB().Query(`
SELECT user_uuid, disabled
FROM user_authentication
`)
	defer func() { _ = rows.Close() }()
	var (
		userUUID string
		disabled bool
	)
	c.Logf("USERS AUTHENTICATION")
	for rows.Next() {
		err := rows.Scan(&userUUID, &disabled)
		c.Assert(err, jc.ErrorIsNil)
		c.Logf("LINE %q, %t", userUUID, disabled)
	}
}

func (s *permissionStateSuite) printRead(c *gc.C) {
	q := `
SELECT  p.uuid, p.grant_on, p.grant_to, p.access_type, p.object_type,
        u.uuid, u.name, creator.name
FROM    v_user_auth u
        JOIN user AS creator ON u.created_by_uuid = creator.uuid
        JOIN v_permission p ON u.uuid = p.grant_to
`
	rows, _ := s.DB().Query(q)
	defer func() { _ = rows.Close() }()
	var (
		permUUID, grantOn, grantTo, accessType, objectType string
		userUUID, userName, createName                     string
	)
	c.Logf("READ")
	for rows.Next() {
		err := rows.Scan(&permUUID, &grantOn, &grantTo, &accessType, &objectType, &userUUID, &userName, &createName)
		c.Assert(err, jc.ErrorIsNil)
		c.Logf("LINE: uuid %q, on %q, to %q, access %q, object %q, user uuid %q, user name %q, creator name%q", permUUID, grantOn, grantTo, accessType, objectType, userUUID, userName, createName)
	}
}

func (s *permissionStateSuite) printClouds(c *gc.C) {
	rows, _ := s.DB().Query(`
SELECT uuid, name
FROM cloud
`)
	defer func() { _ = rows.Close() }()
	var (
		rowUUID, name string
	)

	c.Logf("CLOUDS")
	for rows.Next() {
		err := rows.Scan(&rowUUID, &name)
		c.Assert(err, jc.ErrorIsNil)
		c.Logf("LINE: uuid %q, name %q", rowUUID, name)
	}
}

func (s *permissionStateSuite) printModels(c *gc.C) {
	rows, _ := s.DB().Query(`
SELECT uuid, name, cloud_name, cloud_credential_cloud_name, cloud_credential_name, cloud_credential_owner_name
FROM v_model
`)
	defer func() { _ = rows.Close() }()
	var (
		muuid, mname, cname, cccname, ccname, cconame string
	)

	c.Logf("MODELS")
	for rows.Next() {
		err := rows.Scan(&muuid, &mname, &cname, &cccname, &ccname, &cconame)
		c.Assert(err, jc.ErrorIsNil)
		c.Logf("LINE model-uuid: %q, model-name: %q, cloud-name: %q, cloud-cred-cloud-name: %q, cloud-cred-name: %q, cloud-cred-owner-name: %q",
			muuid, mname, cname, cccname, ccname, cconame)
	}
}
