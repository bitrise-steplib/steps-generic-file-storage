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

func (fs files) String() string {
	fileNames := []string{}
	for _, file := range fs {
		fileNames = append(fileNames, fmt.Sprintf("$%s/%s", envKey, file.Name))
	}
	return strings.Join(fileNames, "\n  ")
}

func urlToFile(urlStr string) (file, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return file{}, err
	}

	return file{
		URL:  urlStr,
		Name: filepath.Base(url.Path),
	}, nil
}
