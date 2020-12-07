package genericfilestorage

import (
	"net/url"
	"strings"

	"github.com/bitrise-io/go-utils/log"

	"github.com/bitrise-io/go-utils/sliceutil"
)

const (
	genericKeyPrefix      = "BITRISEIO_"
	genericKeySuffix      = "_URL"
	pullRequestIgnoredKey = "BITRISEIO_PULL_REQUEST_REPOSITORY_URL"
	awsHostPart           = "concrete-userfiles-production"
	awsPathPart           = "project_file_storage_documents"
)

var ignoredKeys = []string{pullRequestIgnoredKey}

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

		if !provider.shouldHandle(key, value) {
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
	return strings.HasPrefix(key, genericKeyPrefix) && strings.HasSuffix(key, genericKeySuffix)
}

func (provider OsEnvFileProvider) isIgnoredKey(key string) bool {
	return sliceutil.IsStringInSlice(key, ignoredKeys)
}

func (provider OsEnvFileProvider) shouldHandle(key, url string) bool {
	validKey := provider.isGenericKey(key) && !provider.isIgnoredKey(key)
	if !validKey {
		return false
	}
	return provider.isValidURL(key, url)
}

func (provider OsEnvFileProvider) isValidURL(key, rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	containsHost := strings.Contains(parsedURL.Host, awsHostPart)
	if !containsHost {
		log.Warnf("URL: %s missing required host part.", key)
		return false
	}

	conatinsPath := strings.Contains(parsedURL.Path, awsPathPart)
	if !conatinsPath {
		log.Warnf("URL: %s missing required path part.", key)
		return false
	}

	return true
}
