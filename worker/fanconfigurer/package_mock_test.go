// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/fanconfigurer (interfaces: FanConfigurerFacade)
//
// Generated by this command:
//
//	mockgen -package fanconfigurer -destination package_mock_test.go github.com/juju/juju/worker/fanconfigurer FanConfigurerFacade
//

// Package fanconfigurer is a generated GoMock package.
package fanconfigurer

import (
	reflect "reflect"

	network "github.com/juju/juju/core/network"
	watcher "github.com/juju/juju/core/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockFanConfigurerFacade is a mock of FanConfigurerFacade interface.
type MockFanConfigurerFacade struct {
	ctrl     *gomock.Controller
	recorder *MockFanConfigurerFacadeMockRecorder
}

// MockFanConfigurerFacadeMockRecorder is the mock recorder for MockFanConfigurerFacade.
type MockFanConfigurerFacadeMockRecorder struct {
	mock *MockFanConfigurerFacade
}

// NewMockFanConfigurerFacade creates a new mock instance.
func NewMockFanConfigurerFacade(ctrl *gomock.Controller) *MockFanConfigurerFacade {
	mock := &MockFanConfigurerFacade{ctrl: ctrl}
	mock.recorder = &MockFanConfigurerFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFanConfigurerFacade) EXPECT() *MockFanConfigurerFacadeMockRecorder {
	return m.recorder
}

// FanConfig mocks base method.
func (m *MockFanConfigurerFacade) FanConfig() (network.FanConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FanConfig")
	ret0, _ := ret[0].(network.FanConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FanConfig indicates an expected call of FanConfig.
func (mr *MockFanConfigurerFacadeMockRecorder) FanConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FanConfig", reflect.TypeOf((*MockFanConfigurerFacade)(nil).FanConfig))
}

// WatchForFanConfigChanges mocks base method.
func (m *MockFanConfigurerFacade) WatchForFanConfigChanges() (watcher.NotifyWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchForFanConfigChanges")
	ret0, _ := ret[0].(watcher.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchForFanConfigChanges indicates an expected call of WatchForFanConfigChanges.
func (mr *MockFanConfigurerFacadeMockRecorder) WatchForFanConfigChanges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchForFanConfigChanges", reflect.TypeOf((*MockFanConfigurerFacade)(nil).WatchForFanConfigChanges))
}
