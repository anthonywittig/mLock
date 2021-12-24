// Code generated by MockGen. DO NOT EDIT.
// Source: lockengine.go

// Package mock_lockengine is a generated GoMock package.
package mock_lockengine

import (
	context "context"
	shared "mlock/lambdas/shared"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockDeviceController is a mock of DeviceController interface.
type MockDeviceController struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceControllerMockRecorder
}

// MockDeviceControllerMockRecorder is the mock recorder for MockDeviceController.
type MockDeviceControllerMockRecorder struct {
	mock *MockDeviceController
}

// NewMockDeviceController creates a new mock instance.
func NewMockDeviceController(ctrl *gomock.Controller) *MockDeviceController {
	mock := &MockDeviceController{ctrl: ctrl}
	mock.recorder = &MockDeviceControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceController) EXPECT() *MockDeviceControllerMockRecorder {
	return m.recorder
}

// AddLockCode mocks base method.
func (m *MockDeviceController) AddLockCode(ctx context.Context, prop shared.Property, device shared.Device, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLockCode", ctx, prop, device, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLockCode indicates an expected call of AddLockCode.
func (mr *MockDeviceControllerMockRecorder) AddLockCode(ctx, prop, device, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLockCode", reflect.TypeOf((*MockDeviceController)(nil).AddLockCode), ctx, prop, device, code)
}

// RemoveLockCode mocks base method.
func (m *MockDeviceController) RemoveLockCode(ctx context.Context, prop shared.Property, device shared.Device, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLockCode", ctx, prop, device, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLockCode indicates an expected call of RemoveLockCode.
func (mr *MockDeviceControllerMockRecorder) RemoveLockCode(ctx, prop, device, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLockCode", reflect.TypeOf((*MockDeviceController)(nil).RemoveLockCode), ctx, prop, device, code)
}

// MockDeviceRepository is a mock of DeviceRepository interface.
type MockDeviceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceRepositoryMockRecorder
}

// MockDeviceRepositoryMockRecorder is the mock recorder for MockDeviceRepository.
type MockDeviceRepositoryMockRecorder struct {
	mock *MockDeviceRepository
}

// NewMockDeviceRepository creates a new mock instance.
func NewMockDeviceRepository(ctrl *gomock.Controller) *MockDeviceRepository {
	mock := &MockDeviceRepository{ctrl: ctrl}
	mock.recorder = &MockDeviceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceRepository) EXPECT() *MockDeviceRepositoryMockRecorder {
	return m.recorder
}

// AppendToAuditLog mocks base method.
func (m *MockDeviceRepository) AppendToAuditLog(ctx context.Context, device shared.Device, managedLockCodes []*shared.DeviceManagedLockCode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AppendToAuditLog", ctx, device, managedLockCodes)
	ret0, _ := ret[0].(error)
	return ret0
}

// AppendToAuditLog indicates an expected call of AppendToAuditLog.
func (mr *MockDeviceRepositoryMockRecorder) AppendToAuditLog(ctx, device, managedLockCodes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendToAuditLog", reflect.TypeOf((*MockDeviceRepository)(nil).AppendToAuditLog), ctx, device, managedLockCodes)
}

// List mocks base method.
func (m *MockDeviceRepository) List(ctx context.Context) ([]shared.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]shared.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDeviceRepositoryMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDeviceRepository)(nil).List), ctx)
}

// Put mocks base method.
func (m *MockDeviceRepository) Put(ctx context.Context, item shared.Device) (shared.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, item)
	ret0, _ := ret[0].(shared.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockDeviceRepositoryMockRecorder) Put(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockDeviceRepository)(nil).Put), ctx, item)
}

// MockEmailService is a mock of EmailService interface.
type MockEmailService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailServiceMockRecorder
}

// MockEmailServiceMockRecorder is the mock recorder for MockEmailService.
type MockEmailServiceMockRecorder struct {
	mock *MockEmailService
}

// NewMockEmailService creates a new mock instance.
func NewMockEmailService(ctrl *gomock.Controller) *MockEmailService {
	mock := &MockEmailService{ctrl: ctrl}
	mock.recorder = &MockEmailServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailService) EXPECT() *MockEmailServiceMockRecorder {
	return m.recorder
}

// SendEamil mocks base method.
func (m *MockEmailService) SendEamil(ctx context.Context, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEamil", ctx, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEamil indicates an expected call of SendEamil.
func (mr *MockEmailServiceMockRecorder) SendEamil(ctx, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEamil", reflect.TypeOf((*MockEmailService)(nil).SendEamil), ctx, subject, body)
}

// MockPropertyRepository is a mock of PropertyRepository interface.
type MockPropertyRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPropertyRepositoryMockRecorder
}

// MockPropertyRepositoryMockRecorder is the mock recorder for MockPropertyRepository.
type MockPropertyRepositoryMockRecorder struct {
	mock *MockPropertyRepository
}

// NewMockPropertyRepository creates a new mock instance.
func NewMockPropertyRepository(ctrl *gomock.Controller) *MockPropertyRepository {
	mock := &MockPropertyRepository{ctrl: ctrl}
	mock.recorder = &MockPropertyRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPropertyRepository) EXPECT() *MockPropertyRepositoryMockRecorder {
	return m.recorder
}

// GetCached mocks base method.
func (m *MockPropertyRepository) GetCached(ctx context.Context, id uuid.UUID) (shared.Property, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCached", ctx, id)
	ret0, _ := ret[0].(shared.Property)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCached indicates an expected call of GetCached.
func (mr *MockPropertyRepositoryMockRecorder) GetCached(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCached", reflect.TypeOf((*MockPropertyRepository)(nil).GetCached), ctx, id)
}
