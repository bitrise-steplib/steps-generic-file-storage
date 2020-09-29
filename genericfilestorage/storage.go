package genericfilestorage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

// FileProvider ...
type FileProvider interface {
	GetFiles() ([]File, error)
}

// BulkFileDownloader ...
type BulkFileDownloader interface {
	DownloadFiles(files []File, targetDir string) error
}

// Storage ...
type Storage struct {
	bulkdownloader BulkFileDownloader
	fileprovider   FileProvider
	envKey         string
}

// NewStorage ...
func NewStorage(bulkdownloader BulkFileDownloader,
	fileprovider FileProvider,
	envKey string) Storage {
	return Storage{
		bulkdownloader: bulkdownloader,
		fileprovider:   fileprovider,
		envKey:         envKey,
	}
}

// DownloadFiles ...
func (storage Storage) DownloadFiles() (string, error) {

	localStorageDir, err := storage.createLocalStorageDir()
	if err != nil {
		return "", err
	}

	files, err := storage.filesToDownload()
	if err != nil {
		return "", err
	}

	fmt.Println()
	log.Infof("Downloading %d files:", len(files))

	started := time.Now()
	if err := storage.bulkdownloader.DownloadFiles(files, localStorageDir); err != nil {
		return "", fmt.Errorf("failed to download files, error: %s", err)
	}

	log.Printf("  Took: %s", time.Since(started))
	log.Donef("- Done")
	return localStorageDir, nil
}

func (storage Storage) createLocalStorageDir() (string, error) {
	log.Infof("Create storage dir:")

	storageDir, err := pathutil.NormalizedOSTempDirPath(storage.envKey)
	if err != nil {
		return "", fmt.Errorf("failed to create storage temp dir, error: %s", err)
	}
	log.Donef("- Done")

	return storageDir, nil
}

func (storage Storage) filesToDownload() ([]File, error) {
	fmt.Println()
	log.Infof("Parsing Generic Storage Files:")

	fs, err := storage.fileprovider.GetFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file list, error: %s", err)
	}

	log.Debugf("Files to download:")
	storage.logDebugPretty(fs)

	if len(fs) > 0 {
		log.Printf("  %s", files(fs).String(storage.envKey))
	}
	log.Donef("- Done")

	return fs, nil
}

func (storage Storage) logDebugPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("%v\n", string(b))
}
