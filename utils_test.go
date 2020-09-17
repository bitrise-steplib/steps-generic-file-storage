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
			name: "Test BITRISEIO_SAMPLE_FILE_URL generic env",
			key:  "BITRISEIO_SAMPLE_FILE_URL",
			want: true,
		},
		{
			name: "Test BITRISEIO_PULL_REQUEST_REPOSITORY_URL generic env",
			key:  "BITRISEIO_PULL_REQUEST_REPOSITORY_URL",
			want: true,
		},
		{
			name: "Test BITRISEIO_SAMPLE_FILE NON generic env",
			key:  "BITRISEIO_SAMPLE_FILE",
			want: false,
		},
		{
			name: "Test BITRISEIO_PULL_REQUEST_REPOSITORY_url NON generic env",
			key:  "BITRISEIO_PULL_REQUEST_REPOSITORY_url",
			want: false,
		},
		{
			name: "Test PULL_REQUEST_REPOSITORY_URL NON generic env",
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

func Test_splitEnv(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		want  string
		want1 string
	}{
		{
			name:  "Test valid env",
			env:   "KEY=value",
			want:  "KEY",
			want1: "value",
		},
		{
			name:  "Test empty env",
			env:   "KEY=",
			want:  "KEY",
			want1: "",
		},
		{
			name:  "Test invalid env",
			env:   "KEY",
			want:  "KEY",
			want1: "",
		},
		{
			name:  "Test valid env with = character in the value",
			env:   "KEY=value=value!=val?",
			want:  "KEY",
			want1: "value=value!=val?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitEnv(tt.env)
			if got != tt.want {
				t.Errorf("splitEnv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitEnv() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
