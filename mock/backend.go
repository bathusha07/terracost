// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cycloidio/terracost/backend (interfaces: Backend)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	price "github.com/cycloidio/terracost/price"
	product "github.com/cycloidio/terracost/product"
	gomock "github.com/golang/mock/gomock"
)

// Backend is a mock of Backend interface.
type Backend struct {
	ctrl     *gomock.Controller
	recorder *BackendMockRecorder
}

// BackendMockRecorder is the mock recorder for Backend.
type BackendMockRecorder struct {
	mock *Backend
}

// NewBackend creates a new mock instance.
func NewBackend(ctrl *gomock.Controller) *Backend {
	mock := &Backend{ctrl: ctrl}
	mock.recorder = &BackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Backend) EXPECT() *BackendMockRecorder {
	return m.recorder
}

// Prices mocks base method.
func (m *Backend) Prices() price.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prices")
	ret0, _ := ret[0].(price.Repository)
	return ret0
}

// Prices indicates an expected call of Prices.
func (mr *BackendMockRecorder) Prices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prices", reflect.TypeOf((*Backend)(nil).Prices))
}

// Products mocks base method.
func (m *Backend) Products() product.Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Products")
	ret0, _ := ret[0].(product.Repository)
	return ret0
}

// Products indicates an expected call of Products.
func (mr *BackendMockRecorder) Products() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Products", reflect.TypeOf((*Backend)(nil).Products))
}
