package genericfilestorage

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_givenBulkDownloadFail_whenDownloadFilesCalled_thenExpectError(t *testing.T) {
	// Given
	files := givenFiles()
	mockFileProvider := givenMockFileProvider().
		GivenGetFilesSucceed(files)
	expectedError := givenError()
	mockBulkdownloader := givenMockBulkFileDownloader().
		GivenDownloadFilesFail(expectedError)
	storage := Storage{
		bulkdownloader: mockBulkdownloader,
		fileprovider:   mockFileProvider,
		envKey:         "whatever",
	}

	// When
	path, actualErr := storage.DownloadFiles()

	// Then
	assert.EqualError(t, fmt.Errorf("failed to download files, error: %s", expectedError), actualErr.Error())
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
	defer func() {
		err := os.Remove(path)
		require.NoError(t, err)
	}()

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.DirExists(t, path)

}

func Test_givenFileProviderFails_whenFilesToDownloadCalled_thenExpectError(t *testing.T) {
	// Given
	expectedError := givenError()
	mockFileProvider := givenMockFileProvider().
		GivenGetFilesFail(expectedError)
	mockBulkdownloader := givenMockBulkFileDownloader()
	storage := Storage{
		bulkdownloader: mockBulkdownloader,
		fileprovider:   mockFileProvider,
		envKey:         "whatever",
	}

	// When
	files, actualErr := storage.filesToDownload()

	// Then
	assert.EqualError(t, fmt.Errorf("failed to fetch file list, error: %s", expectedError), actualErr.Error())
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
		{name: "name", url: "http://host.com/name"},
	}
}

func givenError() error {
	return errors.New("sad error :(")
}
