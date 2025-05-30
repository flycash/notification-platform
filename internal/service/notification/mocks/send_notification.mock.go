// Code generated by MockGen. DO NOT EDIT.
// Source: ./send_notification.go
//
// Generated by this command:
//
//	mockgen -source=./send_notification.go -destination=./mocks/send_notification.mock.go -package=notificationmocks -typed SendService
//

// Package notificationmocks is a generated GoMock package.
package notificationmocks

import (
	context "context"
	reflect "reflect"

	domain "gitee.com/flycash/notification-platform/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockSendService is a mock of SendService interface.
type MockSendService struct {
	ctrl     *gomock.Controller
	recorder *MockSendServiceMockRecorder
	isgomock struct{}
}

// MockSendServiceMockRecorder is the mock recorder for MockSendService.
type MockSendServiceMockRecorder struct {
	mock *MockSendService
}

// NewMockSendService creates a new mock instance.
func NewMockSendService(ctrl *gomock.Controller) *MockSendService {
	mock := &MockSendService{ctrl: ctrl}
	mock.recorder = &MockSendServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendService) EXPECT() *MockSendServiceMockRecorder {
	return m.recorder
}

// BatchSendNotifications mocks base method.
func (m *MockSendService) BatchSendNotifications(ctx context.Context, ns ...domain.Notification) (domain.BatchSendResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range ns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BatchSendNotifications", varargs...)
	ret0, _ := ret[0].(domain.BatchSendResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchSendNotifications indicates an expected call of BatchSendNotifications.
func (mr *MockSendServiceMockRecorder) BatchSendNotifications(ctx any, ns ...any) *MockSendServiceBatchSendNotificationsCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, ns...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSendNotifications", reflect.TypeOf((*MockSendService)(nil).BatchSendNotifications), varargs...)
	return &MockSendServiceBatchSendNotificationsCall{Call: call}
}

// MockSendServiceBatchSendNotificationsCall wrap *gomock.Call
type MockSendServiceBatchSendNotificationsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendServiceBatchSendNotificationsCall) Return(arg0 domain.BatchSendResponse, arg1 error) *MockSendServiceBatchSendNotificationsCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendServiceBatchSendNotificationsCall) Do(f func(context.Context, ...domain.Notification) (domain.BatchSendResponse, error)) *MockSendServiceBatchSendNotificationsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendServiceBatchSendNotificationsCall) DoAndReturn(f func(context.Context, ...domain.Notification) (domain.BatchSendResponse, error)) *MockSendServiceBatchSendNotificationsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BatchSendNotificationsAsync mocks base method.
func (m *MockSendService) BatchSendNotificationsAsync(ctx context.Context, ns ...domain.Notification) (domain.BatchSendAsyncResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range ns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BatchSendNotificationsAsync", varargs...)
	ret0, _ := ret[0].(domain.BatchSendAsyncResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchSendNotificationsAsync indicates an expected call of BatchSendNotificationsAsync.
func (mr *MockSendServiceMockRecorder) BatchSendNotificationsAsync(ctx any, ns ...any) *MockSendServiceBatchSendNotificationsAsyncCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, ns...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSendNotificationsAsync", reflect.TypeOf((*MockSendService)(nil).BatchSendNotificationsAsync), varargs...)
	return &MockSendServiceBatchSendNotificationsAsyncCall{Call: call}
}

// MockSendServiceBatchSendNotificationsAsyncCall wrap *gomock.Call
type MockSendServiceBatchSendNotificationsAsyncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendServiceBatchSendNotificationsAsyncCall) Return(arg0 domain.BatchSendAsyncResponse, arg1 error) *MockSendServiceBatchSendNotificationsAsyncCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendServiceBatchSendNotificationsAsyncCall) Do(f func(context.Context, ...domain.Notification) (domain.BatchSendAsyncResponse, error)) *MockSendServiceBatchSendNotificationsAsyncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendServiceBatchSendNotificationsAsyncCall) DoAndReturn(f func(context.Context, ...domain.Notification) (domain.BatchSendAsyncResponse, error)) *MockSendServiceBatchSendNotificationsAsyncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendNotification mocks base method.
func (m *MockSendService) SendNotification(ctx context.Context, n domain.Notification) (domain.SendResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNotification", ctx, n)
	ret0, _ := ret[0].(domain.SendResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendNotification indicates an expected call of SendNotification.
func (mr *MockSendServiceMockRecorder) SendNotification(ctx, n any) *MockSendServiceSendNotificationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotification", reflect.TypeOf((*MockSendService)(nil).SendNotification), ctx, n)
	return &MockSendServiceSendNotificationCall{Call: call}
}

// MockSendServiceSendNotificationCall wrap *gomock.Call
type MockSendServiceSendNotificationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendServiceSendNotificationCall) Return(arg0 domain.SendResponse, arg1 error) *MockSendServiceSendNotificationCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendServiceSendNotificationCall) Do(f func(context.Context, domain.Notification) (domain.SendResponse, error)) *MockSendServiceSendNotificationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendServiceSendNotificationCall) DoAndReturn(f func(context.Context, domain.Notification) (domain.SendResponse, error)) *MockSendServiceSendNotificationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendNotificationAsync mocks base method.
func (m *MockSendService) SendNotificationAsync(ctx context.Context, n domain.Notification) (domain.SendResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNotificationAsync", ctx, n)
	ret0, _ := ret[0].(domain.SendResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendNotificationAsync indicates an expected call of SendNotificationAsync.
func (mr *MockSendServiceMockRecorder) SendNotificationAsync(ctx, n any) *MockSendServiceSendNotificationAsyncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotificationAsync", reflect.TypeOf((*MockSendService)(nil).SendNotificationAsync), ctx, n)
	return &MockSendServiceSendNotificationAsyncCall{Call: call}
}

// MockSendServiceSendNotificationAsyncCall wrap *gomock.Call
type MockSendServiceSendNotificationAsyncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockSendServiceSendNotificationAsyncCall) Return(arg0 domain.SendResponse, arg1 error) *MockSendServiceSendNotificationAsyncCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockSendServiceSendNotificationAsyncCall) Do(f func(context.Context, domain.Notification) (domain.SendResponse, error)) *MockSendServiceSendNotificationAsyncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockSendServiceSendNotificationAsyncCall) DoAndReturn(f func(context.Context, domain.Notification) (domain.SendResponse, error)) *MockSendServiceSendNotificationAsyncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
