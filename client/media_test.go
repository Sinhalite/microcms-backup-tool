package client

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic("Error getting working directory")
	}

	// Get the project root directory (where go.mod is located)
	projectRoot := filepath.Dir(wd)

	// Load .env file from the project root
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic("Error loading .env file: " + envPath)
	}
}

func TestBackupMedia(t *testing.T) {
	type args struct {
		config  *Config
		baseDir string
	}

	mediaAPIKey := os.Getenv("MEDIA_API_KEY")
	if mediaAPIKey == "" {
		t.Skip("Media API key not set in environment variable")
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
