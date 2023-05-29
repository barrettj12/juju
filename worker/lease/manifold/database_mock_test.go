// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/database (interfaces: TxnRunner)

// Package manifold_test is a generated GoMock package.
package manifold_test

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	sqlair "github.com/canonical/sqlair"
	gomock "github.com/golang/mock/gomock"
)

// MockTxnRunner is a mock of TxnRunner interface.
type MockTxnRunner struct {
	ctrl     *gomock.Controller
	recorder *MockTxnRunnerMockRecorder
}

// MockTxnRunnerMockRecorder is the mock recorder for MockTxnRunner.
type MockTxnRunnerMockRecorder struct {
	mock *MockTxnRunner
}

// NewMockTxnRunner creates a new mock instance.
func NewMockTxnRunner(ctrl *gomock.Controller) *MockTxnRunner {
	mock := &MockTxnRunner{ctrl: ctrl}
	mock.recorder = &MockTxnRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTxnRunner) EXPECT() *MockTxnRunnerMockRecorder {
	return m.recorder
}

// StdTxn mocks base method.
func (m *MockTxnRunner) StdTxn(arg0 context.Context, arg1 func(context.Context, *sql.Tx) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StdTxn", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// StdTxn indicates an expected call of StdTxn.
func (mr *MockTxnRunnerMockRecorder) StdTxn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StdTxn", reflect.TypeOf((*MockTxnRunner)(nil).StdTxn), arg0, arg1)
}

// Txn mocks base method.
func (m *MockTxnRunner) Txn(arg0 context.Context, arg1 func(context.Context, *sqlair.TX) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Txn", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Txn indicates an expected call of Txn.
func (mr *MockTxnRunnerMockRecorder) Txn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Txn", reflect.TypeOf((*MockTxnRunner)(nil).Txn), arg0, arg1)
}
