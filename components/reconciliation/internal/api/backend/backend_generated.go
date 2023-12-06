// Code generated by MockGen. DO NOT EDIT.
// Source: backend.go

// Package backend is a generated GoMock package.
package backend

import (
	context "context"
	reflect "reflect"

	service "github.com/formancehq/reconciliation/internal/api/service"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Reconciliation mocks base method.
func (m *MockService) Reconciliation(ctx context.Context, req *service.ReconciliationRequest) (*service.ReconciliationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reconciliation", ctx, req)
	ret0, _ := ret[0].(*service.ReconciliationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reconciliation indicates an expected call of Reconciliation.
func (mr *MockServiceMockRecorder) Reconciliation(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reconciliation", reflect.TypeOf((*MockService)(nil).Reconciliation), ctx, req)
}

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// GetService mocks base method.
func (m *MockBackend) GetService() Service {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetService")
	ret0, _ := ret[0].(Service)
	return ret0
}

// GetService indicates an expected call of GetService.
func (mr *MockBackendMockRecorder) GetService() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetService", reflect.TypeOf((*MockBackend)(nil).GetService))
}
