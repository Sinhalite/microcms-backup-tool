package main

import (
	"log"

	"github.com/Sinhalite/microcms-backup-tool/client"
)

func main() {
	client := &client.Client{Config: &client.Config{}}
	err := client.LoadConfig("config.json")
	if err != nil {
		log.Fatal("正常にオプションをセットできませんでした")
	}

	baseDir, err := client.MakeBackupDir()
	if err != nil {
		log.Fatal("正常にバックアップディレクトリを作成できませんでした")
	}

	err = client.StartBackup(baseDir)
	if err != nil {
		log.Printf("バックアップに失敗しました: %v", err)
		log.Fatal("正常にバックアップを処理できませんでした")
	}
}
