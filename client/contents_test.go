package client

import (
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

func TestBackupContents(t *testing.T) {
	type args struct {
		config  *Config
		baseDir string
	}

	publishAPIKey := os.Getenv("PUBLISH_API_KEY")
	allStatusAPIKey := os.Getenv("ALL_STATUS_API_KEY")
	metaDataAPIKey := os.Getenv("META_DATA_API_KEY")

	if publishAPIKey == "" || allStatusAPIKey == "" || metaDataAPIKey == "" {
		t.Skip("API keys not set in environment variables")
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
