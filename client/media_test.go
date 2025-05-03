package client

import (
	"fmt"
	"testing"
	"time"
)

func TestBackupMedia(t *testing.T) {
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
			name: "api key incorrect",
			args: args{
				config: &Config{
					Target:    "media",
					ServiceID: "backup-test",
					Media: MediaConfig{
						APIKey: "incorrectkey",
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
					Target:    "media",
					ServiceID: "backup-test",
					Media: MediaConfig{
						APIKey: "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C",
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

			client := Client{}
			client.Config = tt.args.config

			err := client.BackupMedia(tt.args.baseDir)
			if err != nil {
				fmt.Println(err)
			}
			got := err == nil
			if got != tt.want {
				t.Errorf("backupMedia() = %v, want %v", got, tt.want)
			}
		})
	}
}
