package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (c *Client) InitOption(modeFlag string) error {
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
