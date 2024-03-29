// Code generated by MockGen. DO NOT EDIT.
// Source: document.go
//
// Generated by this command:
//
//	mockgen -source=document.go -destination=./mocks/document.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	domain "web-studio-backend/internal/app/domain"

	gomock "go.uber.org/mock/gomock"
)

// MockDocumentRepository is a mock of DocumentRepository interface.
type MockDocumentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentRepositoryMockRecorder
}

// MockDocumentRepositoryMockRecorder is the mock recorder for MockDocumentRepository.
type MockDocumentRepositoryMockRecorder struct {
	mock *MockDocumentRepository
}

// NewMockDocumentRepository creates a new mock instance.
func NewMockDocumentRepository(ctrl *gomock.Controller) *MockDocumentRepository {
	mock := &MockDocumentRepository{ctrl: ctrl}
	mock.recorder = &MockDocumentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocumentRepository) EXPECT() *MockDocumentRepositoryMockRecorder {
	return m.recorder
}

// AddDocumentToProject mocks base method.
func (m *MockDocumentRepository) AddDocumentToProject(ctx context.Context, docID, projectID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDocumentToProject", ctx, docID, projectID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDocumentToProject indicates an expected call of AddDocumentToProject.
func (mr *MockDocumentRepositoryMockRecorder) AddDocumentToProject(ctx, docID, projectID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDocumentToProject", reflect.TypeOf((*MockDocumentRepository)(nil).AddDocumentToProject), ctx, docID, projectID)
}

// CreateDocument mocks base method.
func (m *MockDocumentRepository) CreateDocument(ctx context.Context, doc *domain.Document) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDocument", ctx, doc)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDocument indicates an expected call of CreateDocument.
func (mr *MockDocumentRepositoryMockRecorder) CreateDocument(ctx, doc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDocument", reflect.TypeOf((*MockDocumentRepository)(nil).CreateDocument), ctx, doc)
}

// DeleteDocument mocks base method.
func (m *MockDocumentRepository) DeleteDocument(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDocument", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDocument indicates an expected call of DeleteDocument.
func (mr *MockDocumentRepositoryMockRecorder) DeleteDocument(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDocument", reflect.TypeOf((*MockDocumentRepository)(nil).DeleteDocument), ctx, id)
}

// GetDocument mocks base method.
func (m *MockDocumentRepository) GetDocument(ctx context.Context, id int32) (*domain.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDocument", ctx, id)
	ret0, _ := ret[0].(*domain.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDocument indicates an expected call of GetDocument.
func (mr *MockDocumentRepositoryMockRecorder) GetDocument(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDocument", reflect.TypeOf((*MockDocumentRepository)(nil).GetDocument), ctx, id)
}

// GetProjectDocuments mocks base method.
func (m *MockDocumentRepository) GetProjectDocuments(ctx context.Context, projectID int32) ([]domain.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectDocuments", ctx, projectID)
	ret0, _ := ret[0].([]domain.Document)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectDocuments indicates an expected call of GetProjectDocuments.
func (mr *MockDocumentRepositoryMockRecorder) GetProjectDocuments(ctx, projectID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectDocuments", reflect.TypeOf((*MockDocumentRepository)(nil).GetProjectDocuments), ctx, projectID)
}

// RemoveDocumentFromProject mocks base method.
func (m *MockDocumentRepository) RemoveDocumentFromProject(ctx context.Context, docID, projectID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDocumentFromProject", ctx, docID, projectID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDocumentFromProject indicates an expected call of RemoveDocumentFromProject.
func (mr *MockDocumentRepositoryMockRecorder) RemoveDocumentFromProject(ctx, docID, projectID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDocumentFromProject", reflect.TypeOf((*MockDocumentRepository)(nil).RemoveDocumentFromProject), ctx, docID, projectID)
}
