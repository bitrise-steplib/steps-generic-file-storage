package bulkfiledownloader

import (
	"github.com/stretchr/testify/mock"
)

// MockFileDownloader ...
type MockFileDownloader struct {
	mock.Mock
}

// Get ...
func (m *MockFileDownloader) Get(destination, source string) error {
	args := m.Called(destination, source)
	return args.Error(0)
}

// GivenGetFails ...
func (m *MockFileDownloader) GivenGetFails(reason error) *MockFileDownloader {
	m.On("Get", mock.Anything, mock.Anything).Return(reason)
	return m
}

// GivenGetFailsForSecondCall ...
func (m *MockFileDownloader) GivenGetFailsForSecondCall(reason error) *MockFileDownloader {
	m.On("Get", mock.Anything, mock.Anything).
		Once().Return(nil)

	m.On("Get", mock.Anything, mock.Anything).
		Once().Return(reason)
	return m
}

// GivenGetSucceed ...
func (m *MockFileDownloader) GivenGetSucceed() *MockFileDownloader {
	m.On("Get", mock.Anything, mock.Anything).Return(nil)
	return m
}
