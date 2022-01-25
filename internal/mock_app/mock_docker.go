// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/app/docker.go

// Package mock_app is a generated GoMock package.
package mock_app

import (
	context "context"
	reflect "reflect"

	types "github.com/docker/docker/api/types"
	filters "github.com/docker/docker/api/types/filters"
	gomock "github.com/golang/mock/gomock"
)

// MockIDockerClient is a mock of IDockerClient interface.
type MockIDockerClient struct {
	ctrl     *gomock.Controller
	recorder *MockIDockerClientMockRecorder
}

// MockIDockerClientMockRecorder is the mock recorder for MockIDockerClient.
type MockIDockerClientMockRecorder struct {
	mock *MockIDockerClient
}

// NewMockIDockerClient creates a new mock instance.
func NewMockIDockerClient(ctrl *gomock.Controller) *MockIDockerClient {
	mock := &MockIDockerClient{ctrl: ctrl}
	mock.recorder = &MockIDockerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDockerClient) EXPECT() *MockIDockerClientMockRecorder {
	return m.recorder
}

// BuildCachePrune mocks base method.
func (m *MockIDockerClient) BuildCachePrune(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildCachePrune", ctx, opts)
	ret0, _ := ret[0].(*types.BuildCachePruneReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildCachePrune indicates an expected call of BuildCachePrune.
func (mr *MockIDockerClientMockRecorder) BuildCachePrune(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildCachePrune", reflect.TypeOf((*MockIDockerClient)(nil).BuildCachePrune), ctx, opts)
}

// ImageList mocks base method.
func (m *MockIDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageList", ctx, options)
	ret0, _ := ret[0].([]types.ImageSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageList indicates an expected call of ImageList.
func (mr *MockIDockerClientMockRecorder) ImageList(ctx, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageList", reflect.TypeOf((*MockIDockerClient)(nil).ImageList), ctx, options)
}

// ImageRemove mocks base method.
func (m *MockIDockerClient) ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageRemove", ctx, imageID, options)
	ret0, _ := ret[0].([]types.ImageDeleteResponseItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageRemove indicates an expected call of ImageRemove.
func (mr *MockIDockerClientMockRecorder) ImageRemove(ctx, imageID, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageRemove", reflect.TypeOf((*MockIDockerClient)(nil).ImageRemove), ctx, imageID, options)
}

// ImagesPrune mocks base method.
func (m *MockIDockerClient) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (types.ImagesPruneReport, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImagesPrune", ctx, pruneFilters)
	ret0, _ := ret[0].(types.ImagesPruneReport)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImagesPrune indicates an expected call of ImagesPrune.
func (mr *MockIDockerClientMockRecorder) ImagesPrune(ctx, pruneFilters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImagesPrune", reflect.TypeOf((*MockIDockerClient)(nil).ImagesPrune), ctx, pruneFilters)
}
