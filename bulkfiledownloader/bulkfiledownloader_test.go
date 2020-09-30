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
	mockFileDowloader := givenMockFileDownloader().
		GivenGetFails(errors.New("sad error :("))
	bulkdownloader := HTTPBulkFileDownloader{mockFileDowloader}
	files := givenFiles(t)
	target := "whatever/path"

	// When
	actualErr := bulkdownloader.DownloadFiles(files, target)

	// Then
	assert.Error(t, actualErr)
}

func Test_givenFileDownloaderFailsForSecondCall_whenDowloadFilesCalled_thenExpectError(t *testing.T) {
	// Given
	mockFileDowloader := givenMockFileDownloader().
		GivenGetFailsForSecondCall(errors.New("sad error :("))
	bulkdownloader := HTTPBulkFileDownloader{mockFileDowloader}
	files := givenFiles(t)
	target := "whatever/path"

	// When
	actualErr := bulkdownloader.DownloadFiles(files, target)

	// Then
	assert.Error(t, actualErr)
}

func Test_givenFileDownloaderSucceeds_whenDowloadFilesCalled_thenExpectNoError(t *testing.T) {
	// Given
	mockFileDowloader := givenMockFileDownloader().
		GivenGetSucceed()
	bulkdownloader := HTTPBulkFileDownloader{mockFileDowloader}
	files := givenFiles(t)
	target := "whatever/path"

	// When
	actualErr := bulkdownloader.DownloadFiles(files, target)

	// Then
	assert.NoError(t, actualErr)
}

func givenMockFileDownloader() *MockFileDownloader {
	return new(MockFileDownloader)
}

func givenFiles(t *testing.T) []genericfilestorage.File {
	return []genericfilestorage.File{
		givenFile(t, "http://paht.to/file"),
	}
}

func givenFile(t *testing.T, url string) genericfilestorage.File {
	file, err := genericfilestorage.NewFile(url)
	require.NoError(t, err)

	return file
}
