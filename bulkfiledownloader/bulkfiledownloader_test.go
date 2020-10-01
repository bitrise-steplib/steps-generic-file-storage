package bulkfiledownloader

import (
	"errors"
	"testing"

	"github.com/bitrise-steplib/steps-generic-file-storage/genericfilestorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_givenFileDownloaderFails_whenDowloadFilesCalled_thenExpectError(t *testing.T) {
	// Given
	expectedError := givenError()
	mockFileDowloader := givenMockFileDownloader().
		GivenGetFails(expectedError)
	bulkdownloader := HTTPBulkFileDownloader{mockFileDowloader}
	files := givenFiles(t, 1)
	target := "whatever/path"

	// When
	actualErr := bulkdownloader.DownloadFiles(files, target)

	// Then
	assert.EqualError(t, expectedError, actualErr.Error())
}

func Test_givenFileDownloaderFailsForSecondCall_whenDowloadFilesCalled_thenExpectError(t *testing.T) {
	// Given
	expectedError := givenError()
	mockFileDowloader := givenMockFileDownloader().
		GivenGetFailsForSecondCall(expectedError)
	bulkdownloader := HTTPBulkFileDownloader{mockFileDowloader}
	files := givenFiles(t, 2)
	target := "whatever/path"

	// When
	actualErr := bulkdownloader.DownloadFiles(files, target)

	// Then
	assert.EqualError(t, expectedError, actualErr.Error())
	mockFileDowloader.AssertNumberOfCalls(t, "Get", 2)
}

func Test_givenFileDownloaderSucceeds_whenDowloadFilesCalled_thenExpectNoError(t *testing.T) {
	// Given
	mockFileDowloader := givenMockFileDownloader().
		GivenGetSucceed()
	bulkdownloader := HTTPBulkFileDownloader{mockFileDowloader}
	files := givenFiles(t, 1)
	target := "whatever/path"

	// When
	actualErr := bulkdownloader.DownloadFiles(files, target)

	// Then
	assert.NoError(t, actualErr)
}

func givenMockFileDownloader() *MockFileDownloader {
	return new(MockFileDownloader)
}

func givenFiles(t *testing.T, size int) []genericfilestorage.File {
	files := make([]genericfilestorage.File, size)
	for i := range files {
		files[i] = givenFile(t, "http://whatever.com/file.txt")
	}
	return files
}

func givenFile(t *testing.T, url string) genericfilestorage.File {
	file, err := genericfilestorage.NewFile(url)
	require.NoError(t, err)

	return file
}

func givenError() error {
	return errors.New("sad error :(")
}
