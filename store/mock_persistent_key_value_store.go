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

// GetConfigValue mocks base method.
func (m *MockPersistentKeyValueStoreInterface) GetConfigValue(configName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfigValue", configName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigValue indicates an expected call of GetConfigValue.
func (mr *MockPersistentKeyValueStoreInterfaceMockRecorder) GetConfigValue(configName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigValue", reflect.TypeOf((*MockPersistentKeyValueStoreInterface)(nil).GetConfigValue), configName)
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

// SetConfigValue mocks base method.
func (m *MockPersistentKeyValueStoreInterface) SetConfigValue(configName, configValue string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConfigValue", configName, configValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConfigValue indicates an expected call of SetConfigValue.
func (mr *MockPersistentKeyValueStoreInterfaceMockRecorder) SetConfigValue(configName, configValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConfigValue", reflect.TypeOf((*MockPersistentKeyValueStoreInterface)(nil).SetConfigValue), configName, configValue)
}
