package main

import (
	"testing"
)

func Test_getFiles(t *testing.T) {
	tests := []struct {
		name    string
		envs    []string
		want    int
		wantErr bool
	}{
		{
			name:    "Test empty env list",
			envs:    []string{},
			want:    0,
			wantErr: false,
		},
		{
			name: "Test env list",
			envs: []string{
				"BITRISEIO_SAMPLE_FILE_URL=https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh",
				"SAMPLE_FILE_URL=",
				"BITRISEIO_SAMPLE_FILE_URL=https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh",
				"BITRISEIO_SAMPLE_FILE_URL=https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "Test env list with ignorable env",
			envs: []string{
				"BITRISEIO_SAMPLE_FILE_URL=https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh",
				"SAMPLE_FILE_URL=",
				"BITRISEIO_PULL_REQUEST_REPOSITORY_URL=",
				"BITRISEIO_SAMPLE_FILE_URL=https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh",
				"BITRISEIO_SAMPLE_FILE_URL=https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh",
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFiles(tt.envs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				var pths []string
				for _, f := range got {
					pths = append(pths, f.Name)
				}
				t.Errorf("getFiles() returned a wrong sized file list = %v, want %v\nReturned files:\n%v", len(got), tt.want, got)
			}
		})
	}
}
