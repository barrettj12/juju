// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/bootstrap (interfaces: ControllerConfigService,FlagService,ObjectStoreGetter,LegacyState,HTTPClient)
//
// Generated by this command:
//
//	mockgen -package bootstrap -destination bootstrap_mock_test.go github.com/juju/juju/internal/worker/bootstrap ControllerConfigService,FlagService,ObjectStoreGetter,LegacyState,HTTPClient
//

// Package bootstrap is a generated GoMock package.
package bootstrap

import (
	context "context"
	http "net/http"
	reflect "reflect"

	controller "github.com/juju/juju/controller"
	objectstore "github.com/juju/juju/core/objectstore"
	binarystorage "github.com/juju/juju/state/binarystorage"
	gomock "go.uber.org/mock/gomock"
)

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

// MockFlagService is a mock of FlagService interface.
type MockFlagService struct {
	ctrl     *gomock.Controller
	recorder *MockFlagServiceMockRecorder
}

// MockFlagServiceMockRecorder is the mock recorder for MockFlagService.
type MockFlagServiceMockRecorder struct {
	mock *MockFlagService
}

// NewMockFlagService creates a new mock instance.
func NewMockFlagService(ctrl *gomock.Controller) *MockFlagService {
	mock := &MockFlagService{ctrl: ctrl}
	mock.recorder = &MockFlagServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlagService) EXPECT() *MockFlagServiceMockRecorder {
	return m.recorder
}

// GetFlag mocks base method.
func (m *MockFlagService) GetFlag(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlag", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlag indicates an expected call of GetFlag.
func (mr *MockFlagServiceMockRecorder) GetFlag(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlag", reflect.TypeOf((*MockFlagService)(nil).GetFlag), arg0, arg1)
}

// SetFlag mocks base method.
func (m *MockFlagService) SetFlag(arg0 context.Context, arg1 string, arg2 bool, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFlag", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFlag indicates an expected call of SetFlag.
func (mr *MockFlagServiceMockRecorder) SetFlag(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFlag", reflect.TypeOf((*MockFlagService)(nil).SetFlag), arg0, arg1, arg2, arg3)
}

// MockObjectStoreGetter is a mock of ObjectStoreGetter interface.
type MockObjectStoreGetter struct {
	ctrl     *gomock.Controller
	recorder *MockObjectStoreGetterMockRecorder
}

// MockObjectStoreGetterMockRecorder is the mock recorder for MockObjectStoreGetter.
type MockObjectStoreGetterMockRecorder struct {
	mock *MockObjectStoreGetter
}

// NewMockObjectStoreGetter creates a new mock instance.
func NewMockObjectStoreGetter(ctrl *gomock.Controller) *MockObjectStoreGetter {
	mock := &MockObjectStoreGetter{ctrl: ctrl}
	mock.recorder = &MockObjectStoreGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectStoreGetter) EXPECT() *MockObjectStoreGetterMockRecorder {
	return m.recorder
}

// GetObjectStore mocks base method.
func (m *MockObjectStoreGetter) GetObjectStore(arg0 context.Context, arg1 string) (objectstore.ObjectStore, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObjectStore", arg0, arg1)
	ret0, _ := ret[0].(objectstore.ObjectStore)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectStore indicates an expected call of GetObjectStore.
func (mr *MockObjectStoreGetterMockRecorder) GetObjectStore(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectStore", reflect.TypeOf((*MockObjectStoreGetter)(nil).GetObjectStore), arg0, arg1)
}

// MockLegacyState is a mock of LegacyState interface.
type MockLegacyState struct {
	ctrl     *gomock.Controller
	recorder *MockLegacyStateMockRecorder
}

// MockLegacyStateMockRecorder is the mock recorder for MockLegacyState.
type MockLegacyStateMockRecorder struct {
	mock *MockLegacyState
}

// NewMockLegacyState creates a new mock instance.
func NewMockLegacyState(ctrl *gomock.Controller) *MockLegacyState {
	mock := &MockLegacyState{ctrl: ctrl}
	mock.recorder = &MockLegacyStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLegacyState) EXPECT() *MockLegacyStateMockRecorder {
	return m.recorder
}

// ControllerModelUUID mocks base method.
func (m *MockLegacyState) ControllerModelUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerModelUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ControllerModelUUID indicates an expected call of ControllerModelUUID.
func (mr *MockLegacyStateMockRecorder) ControllerModelUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerModelUUID", reflect.TypeOf((*MockLegacyState)(nil).ControllerModelUUID))
}

// ToolsStorage mocks base method.
func (m *MockLegacyState) ToolsStorage(arg0 objectstore.ObjectStore) (binarystorage.StorageCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToolsStorage", arg0)
	ret0, _ := ret[0].(binarystorage.StorageCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToolsStorage indicates an expected call of ToolsStorage.
func (mr *MockLegacyStateMockRecorder) ToolsStorage(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToolsStorage", reflect.TypeOf((*MockLegacyState)(nil).ToolsStorage), arg0)
}

// MockHTTPClient is a mock of HTTPClient interface.
type MockHTTPClient struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPClientMockRecorder
}

// MockHTTPClientMockRecorder is the mock recorder for MockHTTPClient.
type MockHTTPClientMockRecorder struct {
	mock *MockHTTPClient
}

// NewMockHTTPClient creates a new mock instance.
func NewMockHTTPClient(ctrl *gomock.Controller) *MockHTTPClient {
	mock := &MockHTTPClient{ctrl: ctrl}
	mock.recorder = &MockHTTPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPClient) EXPECT() *MockHTTPClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockHTTPClient) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockHTTPClientMockRecorder) Do(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHTTPClient)(nil).Do), arg0)
}
