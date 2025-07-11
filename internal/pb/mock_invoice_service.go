// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/LamThanhNguyen/yoyo-store-backend/internal/pb (interfaces: InvoiceServiceClient)
//
// Generated by this command:
//
//	mockgen -package pb -destination internal/pb/mock_invoice_service.go github.com/LamThanhNguyen/yoyo-store-backend/internal/pb InvoiceServiceClient
//

// Package pb is a generated GoMock package.
package pb

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockInvoiceServiceClient is a mock of InvoiceServiceClient interface.
type MockInvoiceServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockInvoiceServiceClientMockRecorder
	isgomock struct{}
}

// MockInvoiceServiceClientMockRecorder is the mock recorder for MockInvoiceServiceClient.
type MockInvoiceServiceClientMockRecorder struct {
	mock *MockInvoiceServiceClient
}

// NewMockInvoiceServiceClient creates a new mock instance.
func NewMockInvoiceServiceClient(ctrl *gomock.Controller) *MockInvoiceServiceClient {
	mock := &MockInvoiceServiceClient{ctrl: ctrl}
	mock.recorder = &MockInvoiceServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInvoiceServiceClient) EXPECT() *MockInvoiceServiceClientMockRecorder {
	return m.recorder
}

// CreateAndSendInvoice mocks base method.
func (m *MockInvoiceServiceClient) CreateAndSendInvoice(ctx context.Context, in *CreateInvoiceRequest, opts ...grpc.CallOption) (*CreateInvoiceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAndSendInvoice", varargs...)
	ret0, _ := ret[0].(*CreateInvoiceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAndSendInvoice indicates an expected call of CreateAndSendInvoice.
func (mr *MockInvoiceServiceClientMockRecorder) CreateAndSendInvoice(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAndSendInvoice", reflect.TypeOf((*MockInvoiceServiceClient)(nil).CreateAndSendInvoice), varargs...)
}
