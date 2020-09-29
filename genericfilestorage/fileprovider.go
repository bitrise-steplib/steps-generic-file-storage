package genericfilestorage

import (
	"strings"

	"github.com/bitrise-io/go-utils/sliceutil"
)

// OsEnvFileProvider ...
type OsEnvFileProvider struct {
	osEnvs []string
}

// NewOsEnvFileProvider ...
func NewOsEnvFileProvider(osEnvs []string) OsEnvFileProvider {
	return OsEnvFileProvider{osEnvs: osEnvs}
}

// GetFiles ...
func (provider OsEnvFileProvider) GetFiles() ([]File, error) {
	var files []File
	for _, env := range provider.osEnvs {
		key, value := provider.splitEnv(env)

		if !provider.isGenericKey(key) || provider.isIgnoredKey(key) {
			continue
		}

		file, err := NewFile(value)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (provider OsEnvFileProvider) splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

func (provider OsEnvFileProvider) isGenericKey(key string) bool {
	return strings.HasPrefix(key, "BITRISEIO_") && strings.HasSuffix(key, "_URL")
}

func (provider OsEnvFileProvider) isIgnoredKey(key string) bool {
	ignoredKeys := []string{"BITRISEIO_PULL_REQUEST_REPOSITORY_URL"}
	return sliceutil.IsStringInSlice(key, ignoredKeys)
}
