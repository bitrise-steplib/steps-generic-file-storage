package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/bitrise-step-export-universal-apk/filedownloader"
	"github.com/bitrise-steplib/steps-generic-file-storage/bulkfiledownloader"
	"github.com/bitrise-steplib/steps-generic-file-storage/genericfilestorage"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/bitrise-tools/go-steputils/tools"
)

const genericFileStorageEnv = "GENERIC_FILE_STORAGE"

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

// GenericFileStorage ...
type GenericFileStorage interface {
	DownloadFiles() (string, error)
}

// Config is defining the input arguments required by the Step.
type Config struct {
	EnableDebug string `env:"enable_debug"`
}

func main() {
	var config Config
	if err := stepconf.Parse(&config); err != nil {
		failf("Error: %s \n", err)
	}
	stepconf.Print(config)
	fmt.Println()

	if config.EnableDebug == "true" {
		log.SetEnableDebugLog(true)
	}

	fileDownloader := filedownloader.New(http.DefaultClient)
	bulkfileDownloader := bulkfiledownloader.New(fileDownloader)
	fileprovider := genericfilestorage.NewOsEnvFileProvider(os.Environ())

	var storage GenericFileStorage
	storage = genericfilestorage.NewStorage(
		bulkfileDownloader,
		fileprovider,
		genericFileStorageEnv,
	)

	storageDir, err := storage.DownloadFiles()
	if err != nil {
		failf("%s", "az")
	}

	if err := tools.ExportEnvironmentWithEnvman(genericFileStorageEnv, storageDir); err != nil {
		failf("Failed to export path: %s to the env: %s, error: %s", storageDir, genericFileStorageEnv, err)
	}
	log.Printf("  %s: %s", genericFileStorageEnv, storageDir)
}
