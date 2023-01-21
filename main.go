package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func (c *Config) SetConfig(target string, serviceId string, apiKey string, endpoints []string) {
	c.Target = target
	c.ServiceID = serviceId
	c.APIKey = apiKey
	c.Endpoints = endpoints
}

func loadConfig() (*Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	err = d.Decode(&cfg)
	return &cfg, err
}

func getOption(modeFlag string) (*Config, error) {
	// コマンドライン引数で"-mode=auto"がセットされていたら、対話式のメッセージは表示しない
	if modeFlag == "auto" {
		option, err := loadConfig()
		return option, err
	}

	// 対話式でオプションをセット
	option := &Config{}
	var err error

	var target string
	var serviceId string
	var apiKey string
	var endpoints []string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("> モードを選択してください(auto / manual)")
	for scanner.Scan() {
		if scanner.Text() == "auto" {
			option, err = loadConfig()
			break
		}

		if scanner.Text() == "manual" {
			fmt.Println("> 保存する対象を選択してください(all/contents/media)")
			for scanner.Scan() {
				if scanner.Text() != "all" {
					target = scanner.Text()
					break
				}
				if scanner.Text() != "contents" {
					target = scanner.Text()
					break
				}
				if scanner.Text() != "media" {
					target = scanner.Text()
					break
				}
			}

			fmt.Println("> サービスIDを入力してください")
			for scanner.Scan() {
				if scanner.Text() != "" {
					serviceId = scanner.Text()
					break
				}
			}

			fmt.Println("> APIキーを入力してください")
			for scanner.Scan() {
				if scanner.Text() != "" {
					apiKey = scanner.Text()
					break
				}
			}

			if target == "all" || target == "contents" {
				fmt.Println("> エンドポイントの一覧をカンマ区切りで入力してください(hoge,fuga)")
				for scanner.Scan() {
					if scanner.Text() != "" {
						endpointsStr := scanner.Text()
						endpoints = strings.Split(endpointsStr, ",")
						break
					}
				}
			}

			option.SetConfig(target, serviceId, apiKey, endpoints)
			break
		}
	}
	return option, err
}

func main() {
	// コマンドライン因数の取得
	modeFlag := flag.String("mode", "", "mode value")
	flag.Parse()

	option, err := getOption(*modeFlag)
	if err != nil {
		log.Fatal("正常にオプションをセットできませんでした")
	}

	// バックアップのディレクトリ作成
	t := time.Now()
	timeDir := t.Format("2006_01_02_15_04_05")
	baseDir := "backup/" + option.ServiceID + "/" + timeDir + "/"

	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatal("バックアップディレクトリの作成に失敗しました")
	}
	log.Println("バックアップディレクトリを作成しました")
	log.Println("バックアップを開始します")

	switch option.Target {
	case "all":
		backupContents(*option, baseDir)
		backupMedia(*option, baseDir)
	case "contents":
		backupContents(*option, baseDir)
	case "media":
		backupMedia(*option, baseDir)
	default:
		log.Fatal("不明なターゲットが選択されました")
	}
	log.Println("正常にバックアップが終了しました")
}
