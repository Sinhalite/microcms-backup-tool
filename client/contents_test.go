package client

import (
	"testing"
	"time"
)

func TestBackupContents(t *testing.T) {
	type args struct {
		config  *Config
		baseDir string
	}

	const (
		publishAPIKey   = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
		allStatusAPIKey = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
		metaDataAPIKey  = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
	)

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "missing api",
			args: args{
				config: &Config{
					Target:    "contents",
					ServiceID: "backup-test",
					Contents: ContentsConfig{
						GetPublishContentsAPIKey: publishAPIKey,
						Endpoints:                []string{"missing"},
						RequestUnit:              10,
						ClassifyByStatus:         false,
						SaveAsCSV:                false,
					},
				},
				baseDir: "../backup/test/",
			},
			want: false,
		},
		{
			name: "normal",
			args: args{
				config: &Config{
					Target:    "contents",
					ServiceID: "backup-test",
					Contents: ContentsConfig{
						GetPublishContentsAPIKey: publishAPIKey,
						Endpoints:                []string{"test", "test2"},
						RequestUnit:              10,
						ClassifyByStatus:         false,
						SaveAsCSV:                false,
					},
				},
				baseDir: "../backup/test/",
			},
			want: true,
		},
		{
			name: "classify by status true, save as csv false",
			args: args{
				config: &Config{
					Target:    "contents",
					ServiceID: "backup-test",
					Contents: ContentsConfig{
						GetPublishContentsAPIKey:   publishAPIKey,
						GetAllStatusContentsAPIKey: allStatusAPIKey,
						GetContentsMetaDataAPIKey:  metaDataAPIKey,
						Endpoints:                  []string{"test", "test2"},
						RequestUnit:                10,
						ClassifyByStatus:           true,
						SaveAsCSV:                  false,
					},
				},
				baseDir: "../backup/test/",
			},
			want: true,
		},
		{
			name: "classify by status false, save as csv true",
			args: args{
				config: &Config{
					Target:    "contents",
					ServiceID: "backup-test",
					Contents: ContentsConfig{
						GetPublishContentsAPIKey: publishAPIKey,
						Endpoints:                []string{"test", "test2"},
						RequestUnit:              10,
						ClassifyByStatus:         false,
						SaveAsCSV:                true,
					},
				},
				baseDir: "../backup/test/",
			},
			want: true,
		},
		{
			name: "classify by status true, save as csv true",
			args: args{
				config: &Config{
					Target:    "contents",
					ServiceID: "backup-test",
					Contents: ContentsConfig{
						GetPublishContentsAPIKey:   publishAPIKey,
						GetAllStatusContentsAPIKey: allStatusAPIKey,
						GetContentsMetaDataAPIKey:  metaDataAPIKey,
						Endpoints:                  []string{"test", "test2"},
						RequestUnit:                10,
						ClassifyByStatus:           true,
						SaveAsCSV:                  true,
					},
				},
				baseDir: "../backup/test/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(5 * time.Second)

			client := &Client{}
			client.Config = tt.args.config

			err := client.BackupContents(tt.args.baseDir)
			got := err == nil
			if got != tt.want {
				t.Errorf("backupContents() = %v, want %v", got, tt.want)
			}
		})
	}
}
