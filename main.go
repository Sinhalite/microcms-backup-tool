package main

import (
	"flag"
	"log"

	"github.com/Sinhalite/microcms-backup-tool/client"
)

func main() {
	client := &client.Client{Config: &client.Config{}}

	// コマンドライン引数の取得
	modeFlag := flag.String("mode", "", "mode value")
	flag.Parse()

	err := client.InitOption(*modeFlag)
	if err != nil {
		log.Fatal("正常にオプションをセットできませんでした")
	}

	baseDir, err := client.MakeBackupDir()
	if err != nil {
		log.Fatal("正常にバックアップディレクトリを作成できませんでした")
	}

	err = client.StartBackup(baseDir)
	if err != nil {
		log.Fatal("正常にバックアップを処理できませんでした")
	}
}
