package client

import (
	"testing"
	"time"
)

func TestBackupAllTargets(t *testing.T) {
	type args struct {
		config  *Config
		baseDir string
	}

	const (
		publishAPIKey   = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
		allStatusAPIKey = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
		metaDataAPIKey  = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
		mediaAPIKey     = "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C"
	)

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "backup contents only",
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
		{
			name: "backup media only",
			args: args{
				config: &Config{
					Target:    "media",
					ServiceID: "backup-test",
					Media: MediaConfig{
						APIKey: mediaAPIKey,
					},
				},
				baseDir: "../backup/test/",
			},
			want: true,
		},
		{
			name: "backup all targets",
			args: args{
				config: &Config{
					Target:    "all",
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
					Media: MediaConfig{
						APIKey: mediaAPIKey,
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

			err := client.StartBackup(tt.args.baseDir)
			got := err == nil
			if got != tt.want {
				t.Errorf("StartBackup() = %v, want %v", got, tt.want)
			}
		})
	}
}
