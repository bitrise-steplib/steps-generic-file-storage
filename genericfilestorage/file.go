package genericfilestorage

import (
	"fmt"
	neturl "net/url"
	"path/filepath"
	"strings"
)

// File ...
type File struct {
	name string
	url  string
}

// Name ...
func (file File) Name() string {
	return file.name
}

// URL ...
func (file File) URL() string {
	return file.url
}

// NewFile ...
func NewFile(url string) (File, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return File{}, err
	}

	return File{
		url:  url,
		name: filepath.Base(u.Path),
	}, nil
}

type files []File

// String returns env key and the file name pairs separated by new line
func (files files) String(path string) (fileNames string) {
	for _, file := range files {
		fileNames += fmt.Sprintf("$%s/%s\n  ", path, file.Name())
	}
	return strings.TrimRight(fileNames, "\n  ")
}
