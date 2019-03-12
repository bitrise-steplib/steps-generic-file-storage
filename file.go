package main

import (
	"fmt"
	neturl "net/url"
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

func newFile(url string) (file, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return file{}, err
	}

	return file{
		URL:  url,
		Name: filepath.Base(u.Path),
	}, nil
}
