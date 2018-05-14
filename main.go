package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-steputils/tools"
)

const envKey = "GENERIC_FILE_STORAGE"

type file struct {
	Name string
	URL  string
}

type files []file

func (fs files) String() string {
	fileNames := []string{}
	for _, file := range fs {
		fileNames = append(fileNames, fmt.Sprintf("$%s/%s", envKey, file.Name))
	}
	return strings.Join(fileNames, "\n  ")
}

func main() {
	if os.Getenv("ENABLE_DEBUG") == "true" {
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

	if err := tools.ExportEnvironmentWithEnvman("GENERIC_FILE_STORAGE", storageDir); err != nil {
		failf("Failed to export GENERIC_FILE_STORAGE, error: %s", err)
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

func isGenericKey(key string) bool {
	return strings.HasPrefix(key, "BITRISEIO_") && strings.HasSuffix(key, "_URL")
}

func getStorageTempDirPath() (string, error) {
	return pathutil.NormalizedOSTempDirPath("_GENERIC_FILE_STORAGE_")
}

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

func urlToFile(urlStr string) (file, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return file{}, err
	}
	return file{URL: urlStr, Name: filepath.Base(url.Path)}, nil
}

func getFiles() ([]file, error) {
	files := []file{}

	for _, env := range os.Environ() {
		key, value := splitEnv(env)

		if !isGenericKey(key) {
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

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil {
			err = cerr
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

func logDebugPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("%+v\n", string(b))
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}
