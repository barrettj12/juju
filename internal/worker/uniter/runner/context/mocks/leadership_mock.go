// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/worker/uniter/runner/context (interfaces: LeadershipContext)
//
// Generated by this command:
//
//	mockgen -typed -package mocks -destination mocks/leadership_mock.go github.com/juju/juju/internal/worker/uniter/runner/context LeadershipContext
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLeadershipContext is a mock of LeadershipContext interface.
type MockLeadershipContext struct {
	ctrl     *gomock.Controller
	recorder *MockLeadershipContextMockRecorder
}

// MockLeadershipContextMockRecorder is the mock recorder for MockLeadershipContext.
type MockLeadershipContextMockRecorder struct {
	mock *MockLeadershipContext
}

// NewMockLeadershipContext creates a new mock instance.
func NewMockLeadershipContext(ctrl *gomock.Controller) *MockLeadershipContext {
	mock := &MockLeadershipContext{ctrl: ctrl}
	mock.recorder = &MockLeadershipContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLeadershipContext) EXPECT() *MockLeadershipContextMockRecorder {
	return m.recorder
}

// IsLeader mocks base method.
func (m *MockLeadershipContext) IsLeader() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLeader")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLeader indicates an expected call of IsLeader.
func (mr *MockLeadershipContextMockRecorder) IsLeader() *MockLeadershipContextIsLeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLeader", reflect.TypeOf((*MockLeadershipContext)(nil).IsLeader))
	return &MockLeadershipContextIsLeaderCall{Call: call}
}

// MockLeadershipContextIsLeaderCall wrap *gomock.Call
type MockLeadershipContextIsLeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockLeadershipContextIsLeaderCall) Return(arg0 bool, arg1 error) *MockLeadershipContextIsLeaderCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockLeadershipContextIsLeaderCall) Do(f func() (bool, error)) *MockLeadershipContextIsLeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockLeadershipContextIsLeaderCall) DoAndReturn(f func() (bool, error)) *MockLeadershipContextIsLeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LeaderSettings mocks base method.
func (m *MockLeadershipContext) LeaderSettings(arg0 context.Context) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaderSettings", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LeaderSettings indicates an expected call of LeaderSettings.
func (mr *MockLeadershipContextMockRecorder) LeaderSettings(arg0 any) *MockLeadershipContextLeaderSettingsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaderSettings", reflect.TypeOf((*MockLeadershipContext)(nil).LeaderSettings), arg0)
	return &MockLeadershipContextLeaderSettingsCall{Call: call}
}

// MockLeadershipContextLeaderSettingsCall wrap *gomock.Call
type MockLeadershipContextLeaderSettingsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockLeadershipContextLeaderSettingsCall) Return(arg0 map[string]string, arg1 error) *MockLeadershipContextLeaderSettingsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockLeadershipContextLeaderSettingsCall) Do(f func(context.Context) (map[string]string, error)) *MockLeadershipContextLeaderSettingsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockLeadershipContextLeaderSettingsCall) DoAndReturn(f func(context.Context) (map[string]string, error)) *MockLeadershipContextLeaderSettingsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// WriteLeaderSettings mocks base method.
func (m *MockLeadershipContext) WriteLeaderSettings(arg0 context.Context, arg1 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteLeaderSettings", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteLeaderSettings indicates an expected call of WriteLeaderSettings.
func (mr *MockLeadershipContextMockRecorder) WriteLeaderSettings(arg0, arg1 any) *MockLeadershipContextWriteLeaderSettingsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteLeaderSettings", reflect.TypeOf((*MockLeadershipContext)(nil).WriteLeaderSettings), arg0, arg1)
	return &MockLeadershipContextWriteLeaderSettingsCall{Call: call}
}

// MockLeadershipContextWriteLeaderSettingsCall wrap *gomock.Call
type MockLeadershipContextWriteLeaderSettingsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockLeadershipContextWriteLeaderSettingsCall) Return(arg0 error) *MockLeadershipContextWriteLeaderSettingsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockLeadershipContextWriteLeaderSettingsCall) Do(f func(context.Context, map[string]string) error) *MockLeadershipContextWriteLeaderSettingsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockLeadershipContextWriteLeaderSettingsCall) DoAndReturn(f func(context.Context, map[string]string) error) *MockLeadershipContextWriteLeaderSettingsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
