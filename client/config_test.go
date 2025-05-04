package client

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpFile := "test_config.json"
	defer os.Remove(tmpFile)

	tests := []struct {
		name     string
		content  string
		filePath string
		wantErr  bool
	}{
		{
			name: "正常系: 正しい形式のJSON",
			content: `{
				"target": "contents",
				"serviceId": "test-service",
				"contents": {
					"getPublishContentsAPIKey": "test-key",
					"endpoints": ["test"],
					"requestUnit": 20,
					"classifyByStatus": false,
					"saveAsCSV": false
				}
			}`,
			filePath: tmpFile,
			wantErr:  false,
		},
		{
			name: "正常系: デフォルト値の確認",
			content: `{
				"target": "contents",
				"serviceId": "test-service",
				"contents": {
					"getPublishContentsAPIKey": "test-key",
					"endpoints": ["test"],
					"classifyByStatus": false,
					"saveAsCSV": false
				}
			}`,
			filePath: tmpFile,
			wantErr:  false,
		},
		{
			name:     "異常系: 存在しないファイル",
			content:  "",
			filePath: "non_existent_config.json",
			wantErr:  true,
		},
		{
			name:     "異常系: 不正なJSON形式",
			content:  `{invalid json}`,
			filePath: tmpFile,
			wantErr:  true,
		},
		{
			name: "異常系: 未知のフィールド",
			content: `{
				"target": "contents",
				"serviceId": "test-service",
				"contents": {
					"getPublishContentsAPIKey": "test-key",
					"endpoints": ["test"],
					"unknownField": "value"
				}
			}`,
			filePath: tmpFile,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストファイルの準備
			if tt.content != "" {
				err := os.WriteFile(tt.filePath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("テストファイルの作成に失敗: %v", err)
				}
			}

			// テスト対象のクライアントを作成
			client := &Client{
				Config: &Config{},
			}

			// LoadConfigの実行
			err := client.LoadConfig(tt.filePath)

			// エラーチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 正常系の場合、設定値の検証
			if !tt.wantErr {
				// デフォルト値の確認
				if tt.name == "正常系: デフォルト値の確認" {
					if client.Config.Contents.RequestUnit != 10 {
						t.Errorf("RequestUnit = %v, want %v", client.Config.Contents.RequestUnit, 10)
					}
				}

				// 基本設定の確認
				if client.Config.Target != "contents" {
					t.Errorf("Target = %v, want %v", client.Config.Target, "contents")
				}
				if client.Config.ServiceID != "test-service" {
					t.Errorf("ServiceID = %v, want %v", client.Config.ServiceID, "test-service")
				}

				// Contents設定の確認
				if client.Config.Contents.GetPublishContentsAPIKey != "test-key" {
					t.Errorf("GetPublishContentsAPIKey = %v, want %v", client.Config.Contents.GetPublishContentsAPIKey, "test-key")
				}
				if len(client.Config.Contents.Endpoints) != 1 || client.Config.Contents.Endpoints[0] != "test" {
					t.Errorf("Endpoints = %v, want %v", client.Config.Contents.Endpoints, []string{"test"})
				}
			}
		})
	}
}
