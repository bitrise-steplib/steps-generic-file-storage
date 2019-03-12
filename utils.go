package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/log"
)

func isGenericKey(key string) bool {
	return strings.HasPrefix(key, "BITRISEIO_") && strings.HasSuffix(key, "_URL")
}

func isIgnoredKey(key string) bool {
	ignoredKeys := []string{"BITRISEIO_PULL_REQUEST_REPOSITORY_URL"}

	for _, ignoredKey := range ignoredKeys {
		if key == ignoredKey {
			return true
		}
	}
	return false
}

func splitEnv(env string) (string, string) {
	e := strings.Split(env, "=")
	return e[0], strings.Join(e[1:], "=")
}

func logDebugPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Debugf("%v\n", string(b))
}
