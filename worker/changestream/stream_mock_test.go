// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/changestream (interfaces: ChangeStream,DBGetter,EventMultiplexer,EventMultiplexerWorker,FileNotifyWatcher)

// Package changestream is a generated GoMock package.
package changestream

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	changestream "github.com/juju/juju/core/changestream"
	database "github.com/juju/juju/core/database"
)

// MockChangeStream is a mock of ChangeStream interface.
type MockChangeStream struct {
	ctrl     *gomock.Controller
	recorder *MockChangeStreamMockRecorder
}

// MockChangeStreamMockRecorder is the mock recorder for MockChangeStream.
type MockChangeStreamMockRecorder struct {
	mock *MockChangeStream
}

// NewMockChangeStream creates a new mock instance.
func NewMockChangeStream(ctrl *gomock.Controller) *MockChangeStream {
	mock := &MockChangeStream{ctrl: ctrl}
	mock.recorder = &MockChangeStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChangeStream) EXPECT() *MockChangeStreamMockRecorder {
	return m.recorder
}

// NamespacedEventMux mocks base method.
func (m *MockChangeStream) NamespacedEventMux(arg0 string) (EventMultiplexer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamespacedEventMux", arg0)
	ret0, _ := ret[0].(EventMultiplexer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamespacedEventMux indicates an expected call of NamespacedEventMux.
func (mr *MockChangeStreamMockRecorder) NamespacedEventMux(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamespacedEventMux", reflect.TypeOf((*MockChangeStream)(nil).NamespacedEventMux), arg0)
}

// MockDBGetter is a mock of DBGetter interface.
type MockDBGetter struct {
	ctrl     *gomock.Controller
	recorder *MockDBGetterMockRecorder
}

// MockDBGetterMockRecorder is the mock recorder for MockDBGetter.
type MockDBGetterMockRecorder struct {
	mock *MockDBGetter
}

// NewMockDBGetter creates a new mock instance.
func NewMockDBGetter(ctrl *gomock.Controller) *MockDBGetter {
	mock := &MockDBGetter{ctrl: ctrl}
	mock.recorder = &MockDBGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBGetter) EXPECT() *MockDBGetterMockRecorder {
	return m.recorder
}

// GetDB mocks base method.
func (m *MockDBGetter) GetDB(arg0 string) (database.TxnRunner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDB", arg0)
	ret0, _ := ret[0].(database.TxnRunner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDB indicates an expected call of GetDB.
func (mr *MockDBGetterMockRecorder) GetDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDB", reflect.TypeOf((*MockDBGetter)(nil).GetDB), arg0)
}

// MockEventMultiplexer is a mock of EventMultiplexer interface.
type MockEventMultiplexer struct {
	ctrl     *gomock.Controller
	recorder *MockEventMultiplexerMockRecorder
}

// MockEventMultiplexerMockRecorder is the mock recorder for MockEventMultiplexer.
type MockEventMultiplexerMockRecorder struct {
	mock *MockEventMultiplexer
}

// NewMockEventMultiplexer creates a new mock instance.
func NewMockEventMultiplexer(ctrl *gomock.Controller) *MockEventMultiplexer {
	mock := &MockEventMultiplexer{ctrl: ctrl}
	mock.recorder = &MockEventMultiplexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventMultiplexer) EXPECT() *MockEventMultiplexerMockRecorder {
	return m.recorder
}

// Subscribe mocks base method.
func (m *MockEventMultiplexer) Subscribe(arg0 ...changestream.SubscriptionOption) (changestream.Subscription, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(changestream.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockEventMultiplexerMockRecorder) Subscribe(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockEventMultiplexer)(nil).Subscribe), arg0...)
}

// MockEventMultiplexerWorker is a mock of EventMultiplexerWorker interface.
type MockEventMultiplexerWorker struct {
	ctrl     *gomock.Controller
	recorder *MockEventMultiplexerWorkerMockRecorder
}

// MockEventMultiplexerWorkerMockRecorder is the mock recorder for MockEventMultiplexerWorker.
type MockEventMultiplexerWorkerMockRecorder struct {
	mock *MockEventMultiplexerWorker
}

// NewMockEventMultiplexerWorker creates a new mock instance.
func NewMockEventMultiplexerWorker(ctrl *gomock.Controller) *MockEventMultiplexerWorker {
	mock := &MockEventMultiplexerWorker{ctrl: ctrl}
	mock.recorder = &MockEventMultiplexerWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventMultiplexerWorker) EXPECT() *MockEventMultiplexerWorkerMockRecorder {
	return m.recorder
}

// EventMux mocks base method.
func (m *MockEventMultiplexerWorker) EventMux() EventMultiplexer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventMux")
	ret0, _ := ret[0].(EventMultiplexer)
	return ret0
}

// EventMux indicates an expected call of EventMux.
func (mr *MockEventMultiplexerWorkerMockRecorder) EventMux() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventMux", reflect.TypeOf((*MockEventMultiplexerWorker)(nil).EventMux))
}

// Kill mocks base method.
func (m *MockEventMultiplexerWorker) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockEventMultiplexerWorkerMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockEventMultiplexerWorker)(nil).Kill))
}

// Wait mocks base method.
func (m *MockEventMultiplexerWorker) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockEventMultiplexerWorkerMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockEventMultiplexerWorker)(nil).Wait))
}

// MockFileNotifyWatcher is a mock of FileNotifyWatcher interface.
type MockFileNotifyWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockFileNotifyWatcherMockRecorder
}

// MockFileNotifyWatcherMockRecorder is the mock recorder for MockFileNotifyWatcher.
type MockFileNotifyWatcherMockRecorder struct {
	mock *MockFileNotifyWatcher
}

// NewMockFileNotifyWatcher creates a new mock instance.
func NewMockFileNotifyWatcher(ctrl *gomock.Controller) *MockFileNotifyWatcher {
	mock := &MockFileNotifyWatcher{ctrl: ctrl}
	mock.recorder = &MockFileNotifyWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileNotifyWatcher) EXPECT() *MockFileNotifyWatcherMockRecorder {
	return m.recorder
}

// Changes mocks base method.
func (m *MockFileNotifyWatcher) Changes(arg0 string) (<-chan bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changes", arg0)
	ret0, _ := ret[0].(<-chan bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Changes indicates an expected call of Changes.
func (mr *MockFileNotifyWatcherMockRecorder) Changes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockFileNotifyWatcher)(nil).Changes), arg0)
}
