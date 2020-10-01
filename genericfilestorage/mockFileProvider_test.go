package genericfilestorage

import (
	"github.com/stretchr/testify/mock"
)

// MockFileProvider ...
type MockFileProvider struct {
	mock.Mock
}

// GetFiles ...
func (m *MockFileProvider) GetFiles() ([]File, error) {
	args := m.Called()
	return args.Get(0).([]File), args.Error(1)
}

// GivenGetFilesFail ...
func (m *MockFileProvider) GivenGetFilesFail(reason error) *MockFileProvider {
	var files []File
	m.On("GetFiles").Return(files, reason)
	return m
}

// GivenGetFilesSucceed ...
func (m *MockFileProvider) GivenGetFilesSucceed(files []File) *MockFileProvider {
	m.On("GetFiles").Return(files, nil)
	return m
}
