package main

import (
	"testing"
)

func Test_isIgnoredKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "Not ignored key",
			key:  "BITRISEIO_SAMPLE_FILE_URL",
			want: false,
		},
		{
			name: "Ignored key - BITRISEIO_PULL_REQUEST_REPOSITORY_URL",
			key:  "BITRISEIO_PULL_REQUEST_REPOSITORY_URL",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIgnoredKey(tt.key); got != tt.want {
				t.Errorf("isIgnoredKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isGenericKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{
			name: "Generic_1",
			key:  "BITRISEIO_SAMPLE_FILE_URL",
			want: true,
		},
		{
			name: "Generic_2",
			key:  "BITRISEIO_PULL_REQUEST_REPOSITORY_URL",
			want: true,
		},
		{
			name: "Not Generic_1",
			key:  "BITRISEIO_SAMPLE_FILE",
			want: false,
		},
		{
			name: "NOT Generic_2",
			key:  "BITRISEIO_PULL_REQUEST_REPOSITORY_url",
			want: false,
		},
		{
			name: "NOT Generic_3",
			key:  "PULL_REQUEST_REPOSITORY_URL",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGenericKey(tt.key); got != tt.want {
				t.Errorf("isGenericKey() = %v, want %v", got, tt.want)
			}
		})
	}
}