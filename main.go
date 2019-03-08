package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-steputils/tools"
)

const envKey = "GENERIC_FILE_STORAGE"

func getStorageTempDirPath() (string, error) {
	return pathutil.NormalizedOSTempDirPath("_GENERIC_FILE_STORAGE_")
}

func getFiles() ([]file, error) {
	var files []file
	for _, env := range os.Environ() {
		key, value := splitEnv(env)

		if !isGenericKey(key) || isIgnoredKey(key) {
			continue
		}

		f, err := urlToFile(value)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer func() {
		cerr := out.Close()
		if cerr != nil {
			log.Errorf("Error during closing the output: %s", cerr)
		}
	}()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			log.Errorf("Error during closing the response body: %s", cerr)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("non successful statusCode: %d\nurl: %s\nbody:\n%s", resp.StatusCode, url, string(body))
	}

	// Write the body to file
	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return err
	}

	return nil
}

func downloadFiles(path string, files []file) error {
	for _, file := range files {
		if err := downloadFile(filepath.Join(path, file.Name), file.URL); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if os.Getenv("enable_debug") == "true" {
		log.SetEnableDebugLog(true)
	}

	log.Infof("Create Storage dir:")

	storageDir, err := getStorageTempDirPath()
	if err != nil {
		failf("Failed to create storage temp dir, error: %s", err)
	}

	if err := tools.ExportEnvironmentWithEnvman(envKey, storageDir); err != nil {
		failf("Failed to export path: %s to the env: %s, error: %s", storageDir, envKey, err)
	}

	log.Printf("  %s: %s", envKey, storageDir)
	log.Donef("- Done")

	fmt.Println()
	log.Infof("Parsing Generic Storage Files:")

	fs, err := getFiles()
	if err != nil {
		failf("Failed to fetch file list, error: %s", err)
	}

	log.Debugf("Files to download:")
	// logDebugPretty(fs) // TODO: use the pretty.go package instead

	if len(fs) > 0 {
		log.Printf("  %s", files(fs))
	}
	log.Donef("- Done")

	fmt.Println()
	log.Infof("Downloading %d files:", len(fs))

	started := time.Now()
	if err := downloadFiles(storageDir, fs); err != nil {
		failf("Failed to download files, error: %s", err)
	}

	log.Printf("  Took: %s", time.Since(started))
	log.Donef("- Done")
}
