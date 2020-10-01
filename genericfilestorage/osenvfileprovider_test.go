package genericfilestorage

import "testing"

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
	provider := OsEnvFileProvider{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := provider.splitEnv(tt.env)
			if got != tt.want {
				t.Errorf("splitEnv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitEnv() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

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
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
				"SAMPLE_FILE_URL=",
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
				"BITRISEIO_SAMPLE_FILE_URL=https://production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/uploads/59985/file.txt",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "Test env list with ignorable env",
			envs: []string{
				"SAMPLE_FILE_URL=",
				"BITRISEIO_PULL_REQUEST_REPOSITORY_URL=",
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
				"SAMPLE_FILE_URL=",
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
				"BITRISEIO_SAMPLE_FILE_URL=https://concrete-userfiles-production.s3.us-west-2.amazonaws.com/project_file_storage_documents/uploads/59985/file.txt",
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := OsEnvFileProvider{osEnvs: tt.envs}
			got, err := provider.GetFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				var pths []string
				for _, f := range got {
					pths = append(pths, f.Name())
				}
				t.Errorf("getFiles() returned a wrong sized file list = %v, want %v\nReturned files:\n%v", len(got), tt.want, got)
			}
		})
	}
}
