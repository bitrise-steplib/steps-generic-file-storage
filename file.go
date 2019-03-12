package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type file struct {
	Name string
	URL  string
}

type files []file

func (fs files) String() (fileNames string) {
	for _, file := range fs {
		fileNames += fmt.Sprintf("$%s/%s\n  ", envKey, file.Name)
	}
	return strings.TrimRight(fileNames, "\n  ")
}

func newFile(urlStr string) (file, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return file{}, err
	}

	return file{
		URL:  urlStr,
		Name: filepath.Base(url.Path),
	}, nil
}
