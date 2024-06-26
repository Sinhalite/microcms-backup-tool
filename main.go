package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c *Client) loadConfig() error {
	// デフォルト値を設定
	c.Config.RequestUnit = 10

	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	d.Decode(c.Config)
	return nil
}

func (c *Client) initOption(modeFlag string) error {
	var err error
	// コマンドライン引数で"-mode=auto"がセットされていたら、対話式のメッセージは表示しない
	if modeFlag == "auto" {
		err := c.loadConfig()
		return err
	}

	// 対話式でオプションをセット
	var target string
	var serviceId string
	var apiKey string
	var endpoints []string
	var requestUnit = 100

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("> モードを選択してください(auto / manual)")
	for scanner.Scan() {
		if scanner.Text() == "auto" {
			err = c.loadConfig()
			break
		}

		if scanner.Text() == "manual" {
			fmt.Println("> 保存する対象を選択してください(all/contents/media)")
			for scanner.Scan() {
				if scanner.Text() == "all" || scanner.Text() == "contents" || scanner.Text() == "media" {
					target = scanner.Text()
					break
				} else {
					fmt.Println("入力値に誤りがあります。再入力してください。")
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

				fmt.Println("> コンテンツの1回あたりの取得件数（デフォルト100件）を調整しますか？(No/Yes)")
				for scanner.Scan() {
					if scanner.Text() == "Yes" {
						fmt.Println("> 取得件数を数値で入力してください")
						for scanner.Scan() {
							if num, _ := strconv.Atoi(scanner.Text()); num > 100 {
								fmt.Println(">100件以下で入力してください")
							} else if len(scanner.Text()) > 0 {
								num, err := strconv.Atoi(scanner.Text())
								if err != nil {
									return err
								}
								requestUnit = num
								break
							}
						}
						break
					} else if scanner.Text() == "No" {
						break
					} else {
						fmt.Println("> 正しい選択肢を入力してください。")
					}
				}
			}

			c.Config = &Config{
				Target:      target,
				ServiceID:   serviceId,
				APIKey:      apiKey,
				Endpoints:   endpoints,
				RequestUnit: requestUnit,
			}
			break
		}
	}
	return err
}

func main() {
	client := &Client{&Config{}}

	// コマンドライン引数の取得
	modeFlag := flag.String("mode", "", "mode value")
	flag.Parse()

	err := client.initOption(*modeFlag)
	if err != nil {
		log.Fatal("正常にオプションをセットできませんでした")
	}

	// バックアップのディレクトリ作成
	t := time.Now()
	timeDir := t.Format("2006_01_02_15_04_05")
	baseDir := "backup/" + client.Config.ServiceID + "/" + timeDir + "/"

	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Fatal("バックアップディレクトリの作成に失敗しました")
	}
	log.Println("バックアップディレクトリを作成しました")
	log.Println("バックアップを開始します")

	switch client.Config.Target {
	case "all":
		err = client.backupContents(baseDir)
		if err != nil {
			log.Fatal(err)
		}
		err = client.backupMedia(baseDir)
		if err != nil {
			log.Fatal(err)
		}
	case "contents":
		err = client.backupContents(baseDir)
		if err != nil {
			log.Fatal(err)
		}
	case "media":
		err = client.backupMedia(baseDir)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("不明なターゲットが選択されました")
	}
	log.Println("正常にバックアップが終了しました")
}
