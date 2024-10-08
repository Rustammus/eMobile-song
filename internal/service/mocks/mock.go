// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks\mock.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	crud "eMobile/internal/crud"
	dto "eMobile/internal/dto"
	reflect "reflect"

	pgtype "github.com/jackc/pgx/v5/pgtype"
	gomock "go.uber.org/mock/gomock"
)

// MockIAudioService is a mock of IAudioService interface.
type MockIAudioService struct {
	ctrl     *gomock.Controller
	recorder *MockIAudioServiceMockRecorder
}

// MockIAudioServiceMockRecorder is the mock recorder for MockIAudioService.
type MockIAudioServiceMockRecorder struct {
	mock *MockIAudioService
}

// NewMockIAudioService creates a new mock instance.
func NewMockIAudioService(ctrl *gomock.Controller) *MockIAudioService {
	mock := &MockIAudioService{ctrl: ctrl}
	mock.recorder = &MockIAudioServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAudioService) EXPECT() *MockIAudioServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIAudioService) Create(audio *dto.AudioCreate) (pgtype.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", audio)
	ret0, _ := ret[0].(pgtype.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIAudioServiceMockRecorder) Create(audio any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIAudioService)(nil).Create), audio)
}

// Delete mocks base method.
func (m *MockIAudioService) Delete(uuid pgtype.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIAudioServiceMockRecorder) Delete(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIAudioService)(nil).Delete), uuid)
}

// Find mocks base method.
func (m *MockIAudioService) Find(uuid pgtype.UUID) (*dto.AudioRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", uuid)
	ret0, _ := ret[0].(*dto.AudioRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockIAudioServiceMockRecorder) Find(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockIAudioService)(nil).Find), uuid)
}

// FindWithLyric mocks base method.
func (m *MockIAudioService) FindWithLyric(uuid pgtype.UUID) (*dto.AudioReadFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindWithLyric", uuid)
	ret0, _ := ret[0].(*dto.AudioReadFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindWithLyric indicates an expected call of FindWithLyric.
func (mr *MockIAudioServiceMockRecorder) FindWithLyric(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindWithLyric", reflect.TypeOf((*MockIAudioService)(nil).FindWithLyric), uuid)
}

// ListByFilter mocks base method.
func (m *MockIAudioService) ListByFilter(filter *dto.AudioFilter, pag crud.Pagination) ([]dto.AudioRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByFilter", filter, pag)
	ret0, _ := ret[0].([]dto.AudioRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByFilter indicates an expected call of ListByFilter.
func (mr *MockIAudioServiceMockRecorder) ListByFilter(filter, pag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByFilter", reflect.TypeOf((*MockIAudioService)(nil).ListByFilter), filter, pag)
}

// ListPag mocks base method.
func (m *MockIAudioService) ListPag(pag crud.Pagination) ([]dto.AudioRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPag", pag)
	ret0, _ := ret[0].([]dto.AudioRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPag indicates an expected call of ListPag.
func (mr *MockIAudioServiceMockRecorder) ListPag(pag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPag", reflect.TypeOf((*MockIAudioService)(nil).ListPag), pag)
}

// Update mocks base method.
func (m *MockIAudioService) Update(uuid pgtype.UUID, audio *dto.AudioUpdate) (*dto.AudioRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", uuid, audio)
	ret0, _ := ret[0].(*dto.AudioRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIAudioServiceMockRecorder) Update(uuid, audio any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIAudioService)(nil).Update), uuid, audio)
}

// MockILyricService is a mock of ILyricService interface.
type MockILyricService struct {
	ctrl     *gomock.Controller
	recorder *MockILyricServiceMockRecorder
}

// MockILyricServiceMockRecorder is the mock recorder for MockILyricService.
type MockILyricServiceMockRecorder struct {
	mock *MockILyricService
}

// NewMockILyricService creates a new mock instance.
func NewMockILyricService(ctrl *gomock.Controller) *MockILyricService {
	mock := &MockILyricService{ctrl: ctrl}
	mock.recorder = &MockILyricServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockILyricService) EXPECT() *MockILyricServiceMockRecorder {
	return m.recorder
}

// ListByAudioPag mocks base method.
func (m *MockILyricService) ListByAudioPag(uuid pgtype.UUID, pag crud.Pagination) ([]dto.LyricRead, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByAudioPag", uuid, pag)
	ret0, _ := ret[0].([]dto.LyricRead)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByAudioPag indicates an expected call of ListByAudioPag.
func (mr *MockILyricServiceMockRecorder) ListByAudioPag(uuid, pag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByAudioPag", reflect.TypeOf((*MockILyricService)(nil).ListByAudioPag), uuid, pag)
}
