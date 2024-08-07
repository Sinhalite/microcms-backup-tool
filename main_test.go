package main

import (
	"testing"

	"github.com/Sinhalite/microcms-backup-tool/client"
)

func TestBackupMedia(t *testing.T) {
	type args struct {
		config  *client.Config
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
				config: &client.Config{
					Target:    "media",
					ServiceID: "backup-test",
					APIKey:    "incorrectkey",
				},
				baseDir: "backup/test/",
			},
			want: false,
		},
		{
			name: "normal",
			args: args{
				config: &client.Config{
					Target:    "media",
					ServiceID: "backup-test",
					APIKey:    "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C",
				},
				baseDir: "backup/test/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &client.Client{}
			client.Config = tt.args.config

			err := client.BackupMedia(tt.args.baseDir)
			got := err == nil
			if got != tt.want {
				t.Errorf("backupMedia() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackupContents(t *testing.T) {
	type args struct {
		config  *client.Config
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
				config: &client.Config{
					Target:      "contents",
					ServiceID:   "backup-test",
					APIKey:      "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C",
					Endpoints:   []string{"missing"},
					RequestUnit: 10,
				},
				baseDir: "backup/test/",
			},
			want: false,
		},
		{
			name: "normal",
			args: args{
				config: &client.Config{
					Target:      "contents",
					ServiceID:   "backup-test",
					APIKey:      "5Nw9AZH3BRRyOZS73ohPksRnn5sI49BMx05C",
					Endpoints:   []string{"test", "test2"},
					RequestUnit: 10,
				},
				baseDir: "backup/test/",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &client.Client{}
			client.Config = tt.args.config

			err := client.BackupContents(tt.args.baseDir)
			got := err == nil
			if got != tt.want {
				t.Errorf("backupContents() = %v, want %v", got, tt.want)
			}
		})
	}
}
