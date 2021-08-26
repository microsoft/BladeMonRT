// Code generated by MockGen. DO NOT EDIT.
// Source: ./persistent_key_value_store.go

// Package store is a generated GoMock package.
package store

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPersistentKeyValueStoreInterface is a mock of PersistentKeyValueStoreInterface interface.
type MockPersistentKeyValueStoreInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPersistentKeyValueStoreInterfaceMockRecorder
}

// MockPersistentKeyValueStoreInterfaceMockRecorder is the mock recorder for MockPersistentKeyValueStoreInterface.
type MockPersistentKeyValueStoreInterfaceMockRecorder struct {
	mock *MockPersistentKeyValueStoreInterface
}

// NewMockPersistentKeyValueStoreInterface creates a new mock instance.
func NewMockPersistentKeyValueStoreInterface(ctrl *gomock.Controller) *MockPersistentKeyValueStoreInterface {
	mock := &MockPersistentKeyValueStoreInterface{ctrl: ctrl}
	mock.recorder = &MockPersistentKeyValueStoreInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersistentKeyValueStoreInterface) EXPECT() *MockPersistentKeyValueStoreInterfaceMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockPersistentKeyValueStoreInterface) Clear() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear")
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear.
func (mr *MockPersistentKeyValueStoreInterfaceMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockPersistentKeyValueStoreInterface)(nil).Clear))
}

// GetValue mocks base method.
func (m *MockPersistentKeyValueStoreInterface) GetValue(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValue indicates an expected call of GetValue.
func (mr *MockPersistentKeyValueStoreInterfaceMockRecorder) GetValue(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockPersistentKeyValueStoreInterface)(nil).GetValue), key)
}

// InitTable mocks base method.
func (m *MockPersistentKeyValueStoreInterface) InitTable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitTable")
	ret0, _ := ret[0].(error)
	return ret0
}

// InitTable indicates an expected call of InitTable.
func (mr *MockPersistentKeyValueStoreInterfaceMockRecorder) InitTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitTable", reflect.TypeOf((*MockPersistentKeyValueStoreInterface)(nil).InitTable))
}

// SetValue mocks base method.
func (m *MockPersistentKeyValueStoreInterface) SetValue(key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetValue", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetValue indicates an expected call of SetValue.
func (mr *MockPersistentKeyValueStoreInterfaceMockRecorder) SetValue(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetValue", reflect.TypeOf((*MockPersistentKeyValueStoreInterface)(nil).SetValue), key, value)
}
