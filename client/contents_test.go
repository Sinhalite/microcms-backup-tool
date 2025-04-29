package client

import (
	"testing"
)

func TestBackupContents(t *testing.T) {
	type args struct {
		config  *Config
		baseDir string
	}

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
						GetPublishContentsAPIKey: "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C",
						Endpoints:                []string{"missing"},
						RequestUnit:              10,
						ClassifyByStatus:         false,
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
						GetPublishContentsAPIKey: "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C",
						Endpoints:                []string{"test", "test2"},
						RequestUnit:              10,
						ClassifyByStatus:         false,
					},
				},
				baseDir: "../backup/test/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
