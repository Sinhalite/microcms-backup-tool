package client

import (
	"fmt"
	"log"
	"os"
	"time"
)

func (c Client) MakeBackupDir() (string, error) {
	// バックアップのディレクトリ作成
	t := time.Now()
	timeDir := t.Format("2006_01_02_15_04_05")
	baseDir := "backup/" + c.Config.ServiceID + "/" + timeDir + "/"

	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	log.Println("バックアップディレクトリを作成しました")
	return baseDir, nil
}

func (c Client) StartBackup(baseDir string) error {
	log.Println("バックアップを開始します")

	switch c.Config.Target {
	case "all":
		err := c.BackupContents(baseDir)
		if err != nil {
			return err
		}
		err = c.BackupMedia(baseDir)
		if err != nil {
			return err
		}
	case "contents":
		err := c.BackupContents(baseDir)
		if err != nil {
			return err
		}
	case "media":
		err := c.BackupMedia(baseDir)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("不明なターゲットが選択されました")
	}
	log.Println("正常にバックアップが終了しました")
	return nil
}
