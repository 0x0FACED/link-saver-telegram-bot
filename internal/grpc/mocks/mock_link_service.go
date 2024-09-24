// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/0x0FACED/proto-files/link_service/gen (interfaces: LinkServiceClient)
//
// Generated by this command:
//
//	mockgen -destination=internal/grpc/mocks/mock_link_service.go -package=mocks github.com/0x0FACED/proto-files/link_service/gen LinkServiceClient
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gen "github.com/0x0FACED/proto-files/link_service/gen"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockLinkServiceClient is a mock of LinkServiceClient interface.
type MockLinkServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockLinkServiceClientMockRecorder
}

// MockLinkServiceClientMockRecorder is the mock recorder for MockLinkServiceClient.
type MockLinkServiceClientMockRecorder struct {
	mock *MockLinkServiceClient
}

// NewMockLinkServiceClient creates a new mock instance.
func NewMockLinkServiceClient(ctrl *gomock.Controller) *MockLinkServiceClient {
	mock := &MockLinkServiceClient{ctrl: ctrl}
	mock.recorder = &MockLinkServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinkServiceClient) EXPECT() *MockLinkServiceClientMockRecorder {
	return m.recorder
}

// DeleteLink mocks base method.
func (m *MockLinkServiceClient) DeleteLink(arg0 context.Context, arg1 *gen.DeleteLinkRequest, arg2 ...grpc.CallOption) (*gen.DeleteLinkResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteLink", varargs...)
	ret0, _ := ret[0].(*gen.DeleteLinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteLink indicates an expected call of DeleteLink.
func (mr *MockLinkServiceClientMockRecorder) DeleteLink(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLink", reflect.TypeOf((*MockLinkServiceClient)(nil).DeleteLink), varargs...)
}

// GetAllLinks mocks base method.
func (m *MockLinkServiceClient) GetAllLinks(arg0 context.Context, arg1 *gen.GetAllLinksRequest, arg2 ...grpc.CallOption) (*gen.GetAllLinksResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllLinks", varargs...)
	ret0, _ := ret[0].(*gen.GetAllLinksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllLinks indicates an expected call of GetAllLinks.
func (mr *MockLinkServiceClientMockRecorder) GetAllLinks(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllLinks", reflect.TypeOf((*MockLinkServiceClient)(nil).GetAllLinks), varargs...)
}

// GetLink mocks base method.
func (m *MockLinkServiceClient) GetLink(arg0 context.Context, arg1 *gen.GetLinkRequest, arg2 ...grpc.CallOption) (*gen.GetLinkResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLink", varargs...)
	ret0, _ := ret[0].(*gen.GetLinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLink indicates an expected call of GetLink.
func (mr *MockLinkServiceClientMockRecorder) GetLink(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLink", reflect.TypeOf((*MockLinkServiceClient)(nil).GetLink), varargs...)
}

// GetLinks mocks base method.
func (m *MockLinkServiceClient) GetLinks(arg0 context.Context, arg1 *gen.GetLinksRequest, arg2 ...grpc.CallOption) (*gen.GetLinksResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLinks", varargs...)
	ret0, _ := ret[0].(*gen.GetLinksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinks indicates an expected call of GetLinks.
func (mr *MockLinkServiceClientMockRecorder) GetLinks(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinks", reflect.TypeOf((*MockLinkServiceClient)(nil).GetLinks), varargs...)
}

// SaveLink mocks base method.
func (m *MockLinkServiceClient) SaveLink(arg0 context.Context, arg1 *gen.SaveLinkRequest, arg2 ...grpc.CallOption) (*gen.SaveLinkResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SaveLink", varargs...)
	ret0, _ := ret[0].(*gen.SaveLinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveLink indicates an expected call of SaveLink.
func (mr *MockLinkServiceClientMockRecorder) SaveLink(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveLink", reflect.TypeOf((*MockLinkServiceClient)(nil).SaveLink), varargs...)
}
