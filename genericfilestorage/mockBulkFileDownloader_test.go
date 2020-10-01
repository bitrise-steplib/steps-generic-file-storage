package genericfilestorage

import (
	"github.com/stretchr/testify/mock"
)

// MockBulkFileDownloader ...
type MockBulkFileDownloader struct {
	mock.Mock
}

// DownloadFiles ...
func (m *MockBulkFileDownloader) DownloadFiles(files []File, targetDir string) error {
	args := m.Called(files, targetDir)
	return args.Error(0)
}

// GivenDownloadFilesFail ...
func (m *MockBulkFileDownloader) GivenDownloadFilesFail(reason error) *MockBulkFileDownloader {
	m.On("DownloadFiles", mock.Anything, mock.Anything).Return(reason)
	return m
}

// GivenDownloadFilesSucceed ...
func (m *MockBulkFileDownloader) GivenDownloadFilesSucceed() *MockBulkFileDownloader {
	m.On("DownloadFiles", mock.Anything, mock.Anything).Return(nil)
	return m
}
