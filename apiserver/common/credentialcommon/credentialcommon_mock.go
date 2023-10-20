// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common/credentialcommon (interfaces: CredentialService)

// Package credentialcommon is a generated GoMock package.
package credentialcommon

import (
	context "context"
	reflect "reflect"

	cloud "github.com/juju/juju/cloud"
	credential "github.com/juju/juju/domain/credential"
	gomock "go.uber.org/mock/gomock"
)

// MockCredentialService is a mock of CredentialService interface.
type MockCredentialService struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialServiceMockRecorder
}

// MockCredentialServiceMockRecorder is the mock recorder for MockCredentialService.
type MockCredentialServiceMockRecorder struct {
	mock *MockCredentialService
}

// NewMockCredentialService creates a new mock instance.
func NewMockCredentialService(ctrl *gomock.Controller) *MockCredentialService {
	mock := &MockCredentialService{ctrl: ctrl}
	mock.recorder = &MockCredentialServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialService) EXPECT() *MockCredentialServiceMockRecorder {
	return m.recorder
}

// CloudCredential mocks base method.
func (m *MockCredentialService) CloudCredential(arg0 context.Context, arg1 credential.ID) (cloud.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential", arg0, arg1)
	ret0, _ := ret[0].(cloud.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockCredentialServiceMockRecorder) CloudCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockCredentialService)(nil).CloudCredential), arg0, arg1)
}

// InvalidateCredential mocks base method.
func (m *MockCredentialService) InvalidateCredential(arg0 context.Context, arg1 credential.ID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateCredential", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateCredential indicates an expected call of InvalidateCredential.
func (mr *MockCredentialServiceMockRecorder) InvalidateCredential(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateCredential", reflect.TypeOf((*MockCredentialService)(nil).InvalidateCredential), arg0, arg1, arg2)
}
