package genericfilestorage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_givenBulkDownloadFail_whenDownloadFilesCalled_thenExpectError(t *testing.T) {
	// Given
	files := givenFiles()
	mockFileProvider := givenMockFileProvider().
		GivenGetFilesSucceed(files)
	mockBulkdownloader := givenMockBulkFileDownloader().
		GivenDownloadFilesFail(errors.New("sad error"))
	storage := Storage{
		bulkdownloader: mockBulkdownloader,
		fileprovider:   mockFileProvider,
		envKey:         "whatever",
	}

	// When
	path, err := storage.DownloadFiles()

	// Then
	assert.Error(t, err)
	assert.Empty(t, path)
}

func Test_givenBulkDownloadSucceeds_whenDownloadFilesCalled_thenExpectNoError(t *testing.T) {
	// Given
	files := givenFiles()
	mockFileProvider := givenMockFileProvider().
		GivenGetFilesSucceed(files)
	mockBulkdownloader := givenMockBulkFileDownloader().
		GivenDownloadFilesSucceed()
	storage := Storage{
		bulkdownloader: mockBulkdownloader,
		fileprovider:   mockFileProvider,
		envKey:         "whatever",
	}

	// When
	path, err := storage.DownloadFiles()

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.DirExists(t, path)
}

func Test_givenFileProviderFails_whenFilesToDownloadCalled_thenExpectError(t *testing.T) {
	// Given
	mockFileProvider := givenMockFileProvider().
		GivenGetFilesFail(errors.New("sad error"))
	mockBulkdownloader := givenMockBulkFileDownloader()
	storage := Storage{
		bulkdownloader: mockBulkdownloader,
		fileprovider:   mockFileProvider,
		envKey:         "whatever",
	}

	// When
	files, err := storage.filesToDownload()

	// Then
	assert.Error(t, err)
	assert.Nil(t, files)
}

func Test_givenFileProviderSucceeds_whenFilesToDownloadCalled_thenExpectFiles(t *testing.T) {
	// Given
	expectedFiles := givenFiles()
	mockFileProvider := givenMockFileProvider().
		GivenGetFilesSucceed(expectedFiles)
	mockBulkdownloader := givenMockBulkFileDownloader()
	storage := Storage{
		bulkdownloader: mockBulkdownloader,
		fileprovider:   mockFileProvider,
		envKey:         "whatever",
	}

	// When
	actualFiles, err := storage.filesToDownload()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedFiles, actualFiles)
}

func givenMockFileProvider() *MockFileProvider {
	return new(MockFileProvider)
}

func givenMockBulkFileDownloader() *MockBulkFileDownloader {
	return new(MockBulkFileDownloader)
}

func givenFiles() []File {
	return []File{
		File{name: "name", url: "http://host.com/name"},
	}
}
