// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/model/modelmigration (interfaces: ModelService,ReadOnlyModelService,UserService,ControllerConfigService)
//
// Generated by this command:
//
//	mockgen -package modelmigration -destination migrations_mock_test.go github.com/juju/juju/domain/model/modelmigration ModelService,ReadOnlyModelService,UserService,ControllerConfigService
//

// Package modelmigration is a generated GoMock package.
package modelmigration

import (
	context "context"
	reflect "reflect"

	controller "github.com/juju/juju/controller"
	model "github.com/juju/juju/core/model"
	user "github.com/juju/juju/core/user"
	model0 "github.com/juju/juju/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockModelService is a mock of ModelService interface.
type MockModelService struct {
	ctrl     *gomock.Controller
	recorder *MockModelServiceMockRecorder
}

// MockModelServiceMockRecorder is the mock recorder for MockModelService.
type MockModelServiceMockRecorder struct {
	mock *MockModelService
}

// NewMockModelService creates a new mock instance.
func NewMockModelService(ctrl *gomock.Controller) *MockModelService {
	mock := &MockModelService{ctrl: ctrl}
	mock.recorder = &MockModelServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelService) EXPECT() *MockModelServiceMockRecorder {
	return m.recorder
}

// CreateModel mocks base method.
func (m *MockModelService) CreateModel(arg0 context.Context, arg1 model0.ModelCreationArgs) (model.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateModel", arg0, arg1)
	ret0, _ := ret[0].(model.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateModel indicates an expected call of CreateModel.
func (mr *MockModelServiceMockRecorder) CreateModel(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateModel", reflect.TypeOf((*MockModelService)(nil).CreateModel), arg0, arg1)
}

// DeleteModel mocks base method.
func (m *MockModelService) DeleteModel(arg0 context.Context, arg1 model.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteModel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteModel indicates an expected call of DeleteModel.
func (mr *MockModelServiceMockRecorder) DeleteModel(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteModel", reflect.TypeOf((*MockModelService)(nil).DeleteModel), arg0, arg1)
}

// ModelType mocks base method.
func (m *MockModelService) ModelType(arg0 context.Context, arg1 model.UUID) (model.ModelType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelType", arg0, arg1)
	ret0, _ := ret[0].(model.ModelType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelType indicates an expected call of ModelType.
func (mr *MockModelServiceMockRecorder) ModelType(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelType", reflect.TypeOf((*MockModelService)(nil).ModelType), arg0, arg1)
}

// MockReadOnlyModelService is a mock of ReadOnlyModelService interface.
type MockReadOnlyModelService struct {
	ctrl     *gomock.Controller
	recorder *MockReadOnlyModelServiceMockRecorder
}

// MockReadOnlyModelServiceMockRecorder is the mock recorder for MockReadOnlyModelService.
type MockReadOnlyModelServiceMockRecorder struct {
	mock *MockReadOnlyModelService
}

// NewMockReadOnlyModelService creates a new mock instance.
func NewMockReadOnlyModelService(ctrl *gomock.Controller) *MockReadOnlyModelService {
	mock := &MockReadOnlyModelService{ctrl: ctrl}
	mock.recorder = &MockReadOnlyModelServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadOnlyModelService) EXPECT() *MockReadOnlyModelServiceMockRecorder {
	return m.recorder
}

// CreateModel mocks base method.
func (m *MockReadOnlyModelService) CreateModel(arg0 context.Context, arg1 model0.ReadOnlyModelCreationArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateModel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateModel indicates an expected call of CreateModel.
func (mr *MockReadOnlyModelServiceMockRecorder) CreateModel(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateModel", reflect.TypeOf((*MockReadOnlyModelService)(nil).CreateModel), arg0, arg1)
}

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// GetUserByName mocks base method.
func (m *MockUserService) GetUserByName(arg0 context.Context, arg1 string) (user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", arg0, arg1)
	ret0, _ := ret[0].(user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockUserServiceMockRecorder) GetUserByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockUserService)(nil).GetUserByName), arg0, arg1)
}

// MockControllerConfigService is a mock of ControllerConfigService interface.
type MockControllerConfigService struct {
	ctrl     *gomock.Controller
	recorder *MockControllerConfigServiceMockRecorder
}

// MockControllerConfigServiceMockRecorder is the mock recorder for MockControllerConfigService.
type MockControllerConfigServiceMockRecorder struct {
	mock *MockControllerConfigService
}

// NewMockControllerConfigService creates a new mock instance.
func NewMockControllerConfigService(ctrl *gomock.Controller) *MockControllerConfigService {
	mock := &MockControllerConfigService{ctrl: ctrl}
	mock.recorder = &MockControllerConfigServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControllerConfigService) EXPECT() *MockControllerConfigServiceMockRecorder {
	return m.recorder
}

// ControllerConfig mocks base method.
func (m *MockControllerConfigService) ControllerConfig(arg0 context.Context) (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig", arg0)
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockControllerConfigServiceMockRecorder) ControllerConfig(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockControllerConfigService)(nil).ControllerConfig), arg0)
}
