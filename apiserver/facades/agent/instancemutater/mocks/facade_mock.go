// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facade (interfaces: Context,Resources,Authorizer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	facade "github.com/juju/juju/apiserver/facade"
	cache "github.com/juju/juju/core/cache"
	leadership "github.com/juju/juju/core/leadership"
	lease "github.com/juju/juju/core/lease"
	permission "github.com/juju/juju/permission"
	state "github.com/juju/juju/state"
	names_v2 "gopkg.in/juju/names.v2"
	reflect "reflect"
)

// MockContext is a mock of Context interface
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// Auth mocks base method
func (m *MockContext) Auth() facade.Authorizer {
	ret := m.ctrl.Call(m, "Auth")
	ret0, _ := ret[0].(facade.Authorizer)
	return ret0
}

// Auth indicates an expected call of Auth
func (mr *MockContextMockRecorder) Auth() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockContext)(nil).Auth))
}

// Controller mocks base method
func (m *MockContext) Controller() *cache.Controller {
	ret := m.ctrl.Call(m, "Controller")
	ret0, _ := ret[0].(*cache.Controller)
	return ret0
}

// Controller indicates an expected call of Controller
func (mr *MockContextMockRecorder) Controller() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Controller", reflect.TypeOf((*MockContext)(nil).Controller))
}

// Dispose mocks base method
func (m *MockContext) Dispose() {
	m.ctrl.Call(m, "Dispose")
}

// Dispose indicates an expected call of Dispose
func (mr *MockContextMockRecorder) Dispose() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispose", reflect.TypeOf((*MockContext)(nil).Dispose))
}

// Hub mocks base method
func (m *MockContext) Hub() facade.Hub {
	ret := m.ctrl.Call(m, "Hub")
	ret0, _ := ret[0].(facade.Hub)
	return ret0
}

// Hub indicates an expected call of Hub
func (mr *MockContextMockRecorder) Hub() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hub", reflect.TypeOf((*MockContext)(nil).Hub))
}

// ID mocks base method
func (m *MockContext) ID() string {
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockContextMockRecorder) ID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockContext)(nil).ID))
}

// LeadershipChecker mocks base method
func (m *MockContext) LeadershipChecker() (leadership.Checker, error) {
	ret := m.ctrl.Call(m, "LeadershipChecker")
	ret0, _ := ret[0].(leadership.Checker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipChecker indicates an expected call of LeadershipChecker
func (mr *MockContextMockRecorder) LeadershipChecker() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipChecker", reflect.TypeOf((*MockContext)(nil).LeadershipChecker))
}

// LeadershipClaimer mocks base method
func (m *MockContext) LeadershipClaimer(arg0 string) (leadership.Claimer, error) {
	ret := m.ctrl.Call(m, "LeadershipClaimer", arg0)
	ret0, _ := ret[0].(leadership.Claimer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipClaimer indicates an expected call of LeadershipClaimer
func (mr *MockContextMockRecorder) LeadershipClaimer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipClaimer", reflect.TypeOf((*MockContext)(nil).LeadershipClaimer), arg0)
}

// LeadershipPinner mocks base method
func (m *MockContext) LeadershipPinner(arg0 string) (leadership.Pinner, error) {
	ret := m.ctrl.Call(m, "LeadershipPinner", arg0)
	ret0, _ := ret[0].(leadership.Pinner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipPinner indicates an expected call of LeadershipPinner
func (mr *MockContextMockRecorder) LeadershipPinner(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipPinner", reflect.TypeOf((*MockContext)(nil).LeadershipPinner), arg0)
}

// LeadershipReader mocks base method
func (m *MockContext) LeadershipReader(arg0 string) (leadership.Reader, error) {
	ret := m.ctrl.Call(m, "LeadershipReader", arg0)
	ret0, _ := ret[0].(leadership.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeadershipReader indicates an expected call of LeadershipReader
func (mr *MockContextMockRecorder) LeadershipReader(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeadershipReader", reflect.TypeOf((*MockContext)(nil).LeadershipReader), arg0)
}

// Presence mocks base method
func (m *MockContext) Presence() facade.Presence {
	ret := m.ctrl.Call(m, "Presence")
	ret0, _ := ret[0].(facade.Presence)
	return ret0
}

// Presence indicates an expected call of Presence
func (mr *MockContextMockRecorder) Presence() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Presence", reflect.TypeOf((*MockContext)(nil).Presence))
}

// Resources mocks base method
func (m *MockContext) Resources() facade.Resources {
	ret := m.ctrl.Call(m, "Resources")
	ret0, _ := ret[0].(facade.Resources)
	return ret0
}

// Resources indicates an expected call of Resources
func (mr *MockContextMockRecorder) Resources() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resources", reflect.TypeOf((*MockContext)(nil).Resources))
}

// SingularClaimer mocks base method
func (m *MockContext) SingularClaimer() (lease.Claimer, error) {
	ret := m.ctrl.Call(m, "SingularClaimer")
	ret0, _ := ret[0].(lease.Claimer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SingularClaimer indicates an expected call of SingularClaimer
func (mr *MockContextMockRecorder) SingularClaimer() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SingularClaimer", reflect.TypeOf((*MockContext)(nil).SingularClaimer))
}

// State mocks base method
func (m *MockContext) State() *state.State {
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(*state.State)
	return ret0
}

// State indicates an expected call of State
func (mr *MockContextMockRecorder) State() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockContext)(nil).State))
}

// StatePool mocks base method
func (m *MockContext) StatePool() *state.StatePool {
	ret := m.ctrl.Call(m, "StatePool")
	ret0, _ := ret[0].(*state.StatePool)
	return ret0
}

// StatePool indicates an expected call of StatePool
func (mr *MockContextMockRecorder) StatePool() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StatePool", reflect.TypeOf((*MockContext)(nil).StatePool))
}

// MockResources is a mock of Resources interface
type MockResources struct {
	ctrl     *gomock.Controller
	recorder *MockResourcesMockRecorder
}

// MockResourcesMockRecorder is the mock recorder for MockResources
type MockResourcesMockRecorder struct {
	mock *MockResources
}

// NewMockResources creates a new mock instance
func NewMockResources(ctrl *gomock.Controller) *MockResources {
	mock := &MockResources{ctrl: ctrl}
	mock.recorder = &MockResourcesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResources) EXPECT() *MockResourcesMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockResources) Get(arg0 string) facade.Resource {
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(facade.Resource)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockResourcesMockRecorder) Get(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockResources)(nil).Get), arg0)
}

// Register mocks base method
func (m *MockResources) Register(arg0 facade.Resource) string {
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockResourcesMockRecorder) Register(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockResources)(nil).Register), arg0)
}

// Stop mocks base method
func (m *MockResources) Stop(arg0 string) error {
	ret := m.ctrl.Call(m, "Stop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockResourcesMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockResources)(nil).Stop), arg0)
}

// MockAuthorizer is a mock of Authorizer interface
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// AuthApplicationAgent mocks base method
func (m *MockAuthorizer) AuthApplicationAgent() bool {
	ret := m.ctrl.Call(m, "AuthApplicationAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthApplicationAgent indicates an expected call of AuthApplicationAgent
func (mr *MockAuthorizerMockRecorder) AuthApplicationAgent() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthApplicationAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthApplicationAgent))
}

// AuthClient mocks base method
func (m *MockAuthorizer) AuthClient() bool {
	ret := m.ctrl.Call(m, "AuthClient")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthClient indicates an expected call of AuthClient
func (mr *MockAuthorizerMockRecorder) AuthClient() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthClient", reflect.TypeOf((*MockAuthorizer)(nil).AuthClient))
}

// AuthController mocks base method
func (m *MockAuthorizer) AuthController() bool {
	ret := m.ctrl.Call(m, "AuthController")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthController indicates an expected call of AuthController
func (mr *MockAuthorizerMockRecorder) AuthController() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthController", reflect.TypeOf((*MockAuthorizer)(nil).AuthController))
}

// AuthMachineAgent mocks base method
func (m *MockAuthorizer) AuthMachineAgent() bool {
	ret := m.ctrl.Call(m, "AuthMachineAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthMachineAgent indicates an expected call of AuthMachineAgent
func (mr *MockAuthorizerMockRecorder) AuthMachineAgent() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthMachineAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthMachineAgent))
}

// AuthOwner mocks base method
func (m *MockAuthorizer) AuthOwner(arg0 names_v2.Tag) bool {
	ret := m.ctrl.Call(m, "AuthOwner", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthOwner indicates an expected call of AuthOwner
func (mr *MockAuthorizerMockRecorder) AuthOwner(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthOwner", reflect.TypeOf((*MockAuthorizer)(nil).AuthOwner), arg0)
}

// AuthUnitAgent mocks base method
func (m *MockAuthorizer) AuthUnitAgent() bool {
	ret := m.ctrl.Call(m, "AuthUnitAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthUnitAgent indicates an expected call of AuthUnitAgent
func (mr *MockAuthorizerMockRecorder) AuthUnitAgent() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthUnitAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthUnitAgent))
}

// ConnectedModel mocks base method
func (m *MockAuthorizer) ConnectedModel() string {
	ret := m.ctrl.Call(m, "ConnectedModel")
	ret0, _ := ret[0].(string)
	return ret0
}

// ConnectedModel indicates an expected call of ConnectedModel
func (mr *MockAuthorizerMockRecorder) ConnectedModel() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectedModel", reflect.TypeOf((*MockAuthorizer)(nil).ConnectedModel))
}

// GetAuthTag mocks base method
func (m *MockAuthorizer) GetAuthTag() names_v2.Tag {
	ret := m.ctrl.Call(m, "GetAuthTag")
	ret0, _ := ret[0].(names_v2.Tag)
	return ret0
}

// GetAuthTag indicates an expected call of GetAuthTag
func (mr *MockAuthorizerMockRecorder) GetAuthTag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthTag", reflect.TypeOf((*MockAuthorizer)(nil).GetAuthTag))
}

// HasPermission mocks base method
func (m *MockAuthorizer) HasPermission(arg0 permission.Access, arg1 names_v2.Tag) (bool, error) {
	ret := m.ctrl.Call(m, "HasPermission", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasPermission indicates an expected call of HasPermission
func (mr *MockAuthorizerMockRecorder) HasPermission(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPermission", reflect.TypeOf((*MockAuthorizer)(nil).HasPermission), arg0, arg1)
}

// UserHasPermission mocks base method
func (m *MockAuthorizer) UserHasPermission(arg0 names_v2.UserTag, arg1 permission.Access, arg2 names_v2.Tag) (bool, error) {
	ret := m.ctrl.Call(m, "UserHasPermission", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserHasPermission indicates an expected call of UserHasPermission
func (mr *MockAuthorizerMockRecorder) UserHasPermission(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserHasPermission", reflect.TypeOf((*MockAuthorizer)(nil).UserHasPermission), arg0, arg1, arg2)
}
