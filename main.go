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

const genericFileStorageEnv = "GENERIC_FILE_STORAGE"

func getFiles(envs []string) (files, error) {
	var files []file
	for _, env := range envs {
		key, value := splitEnv(env)

		if !isGenericKey(key) || isIgnoredKey(key) {
			continue
		}

		f, err := newFile(value)
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

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

func main() {
	if os.Getenv("enable_debug") == "true" {
		log.SetEnableDebugLog(true)
	}

	log.Infof("Create storage dir:")

	storageDir, err := pathutil.NormalizedOSTempDirPath(genericFileStorageEnv)
	if err != nil {
		failf("Failed to create storage temp dir, error: %s", err)
	}

	if err := tools.ExportEnvironmentWithEnvman(genericFileStorageEnv, storageDir); err != nil {
		failf("Failed to export path: %s to the env: %s, error: %s", storageDir, genericFileStorageEnv, err)
	}

	log.Printf("  %s: %s", genericFileStorageEnv, storageDir)
	log.Donef("- Done")

	fmt.Println()
	log.Infof("Parsing Generic Storage Files:")

	fs, err := getFiles(os.Environ())
	if err != nil {
		failf("Failed to fetch file list, error: %s", err)
	}

	log.Debugf("Files to download:")
	logDebugPretty(fs)

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
