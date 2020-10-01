package bulkfiledownloader

import (
	"path/filepath"

	"github.com/bitrise-steplib/steps-generic-file-storage/genericfilestorage"
)

// FileDownloader represents a type that can download a file.
type FileDownloader interface {
	Get(destination, source string) error
}

// HTTPBulkFileDownloader ...
type HTTPBulkFileDownloader struct {
	downloader FileDownloader
}

// New ...
func New(downloader FileDownloader) HTTPBulkFileDownloader {
	return HTTPBulkFileDownloader{downloader: downloader}
}

// DownloadFiles ...
func (bulkDownloader HTTPBulkFileDownloader) DownloadFiles(files []genericfilestorage.File, targetDir string) error {
	for _, file := range files {
		destination := filepath.Join(targetDir, file.Name())
		if err := bulkDownloader.downloader.Get(destination, file.URL()); err != nil {
			return err
		}
	}

	return nil
}
