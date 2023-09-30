// Code generated by MockGen. DO NOT EDIT.
// Source: model.go

// Package mock_note is a generated GoMock package.
package mock_note

import (
	context "context"
	note "notes-manager/src/internal/note"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, note *note.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, note)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, note)
}

// Delete mocks base method.
func (m *MockRepository) Delete(ctx context.Context, noteId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, noteId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(ctx, noteId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ctx, noteId)
}

// Fetch mocks base method.
func (m *MockRepository) Fetch(ctx context.Context, noteId uuid.UUID) (*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", ctx, noteId)
	ret0, _ := ret[0].(*note.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockRepositoryMockRecorder) Fetch(ctx, noteId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockRepository)(nil).Fetch), ctx, noteId)
}

// FetchAll mocks base method.
func (m *MockRepository) FetchAll(ctx context.Context, userId uuid.UUID) ([]*note.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchAll", ctx, userId)
	ret0, _ := ret[0].([]*note.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchAll indicates an expected call of FetchAll.
func (mr *MockRepositoryMockRecorder) FetchAll(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchAll", reflect.TypeOf((*MockRepository)(nil).FetchAll), ctx, userId)
}

// Update mocks base method.
func (m *MockRepository) Update(ctx context.Context, note *note.NoteDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, note)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(ctx, note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), ctx, note)
}
