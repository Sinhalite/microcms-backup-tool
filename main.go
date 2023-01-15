package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	Media      []Media
	TotalCount int
	Limit      int
	Offset     int
}

type Media struct {
	Id     string
	Url    string
	Width  int
	Height int
}

type Config struct {
	ServiceID string `json:"serviceId"`
	APIKey    string `json:"apiKey"`
}

func (c *Config) SetConfig(serviceId string, apiKey string) {
	c.ServiceID = serviceId
	c.APIKey = apiKey
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

func getTotalCount(option Config) (int, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms-management.io/api/v1/media?limit=0", option.ServiceID),
		nil)
	req.Header.Set("X-MICROCMS-API-KEY", option.APIKey)

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return 0, err
	}

	return response.TotalCount, err
}

func getMediaAry(option Config, requiredRequestCount int, requestUnit int) ([]Media, error) {
	var ary []Media

	for i := 0; i < requiredRequestCount; i++ {
		client := new(http.Client)
		req, _ := http.NewRequest(
			"GET",
			fmt.Sprintf("https://%s.microcms-management.io/api/v1/media?limit=%d&offset=%d", option.ServiceID, requestUnit, requestUnit*i),
			nil,
		)
		req.Header.Set("X-MICROCMS-API-KEY", option.APIKey)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		response := &Response{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return nil, err
		}

		ary = append(ary, response.Media...)
	}

	return ary, nil
}

func saveMedia(mediaAry []Media, option Config, totalCount int) error {
	t := time.Now()
	timeDir := t.Format("2006_01_02_15_04_05")
	baseDir := "backup/" + option.ServiceID + "/media/" + timeDir + "/"

	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return err
	}

	for i, media := range mediaAry {
		client := new(http.Client)
		req, _ := http.NewRequest("GET", media.Url, nil)
		req.Header.Set("X-MICROCMS-API-KEY", option.APIKey)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer resp.Body.Close()

		ary := strings.Split(media.Url, "/")
		fileName := ary[len(ary)-1]
		fileDirectory := ary[len(ary)-2]

		// ファイルごとのディレクトリを作成する
		// (同じファイル名でアップロード可能なため、一意となるようなパスが付与されている)
		err = os.Mkdir(baseDir+fileDirectory, os.ModePerm)
		if err != nil {
			return err
		}

		file, err := os.Create(baseDir + fileDirectory + "/" + fileName)
		if err != nil {
			return err

		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, totalCount, media.Url)

	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	option := &Config{}

	var serviceId string
	var apiKey string

	fmt.Println("> モードを選択してください(auto / manual)")
	for scanner.Scan() {
		if scanner.Text() == "auto" {
			var err error
			option, err = loadConfig()
			if err != nil {
				log.Fatal(err)
			}
			break
		}

		if scanner.Text() == "manual" {
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
			option.SetConfig(serviceId, apiKey)
			break
		}
	}

	const requestUnit = 50
	totalCount, err := getTotalCount(*option)
	if err != nil {
		log.Println(err)
		log.Fatal("合計件数の取得でエラーが発生しました")
	}
	requiredRequestCount := (totalCount/requestUnit + 1)

	mediaAry, err := getMediaAry(*option, requiredRequestCount, requestUnit)
	if err != nil {
		log.Println(err)
		log.Fatal("メディア一覧の取得でエラーが発生しました")
	}
	err = saveMedia(mediaAry, *option, totalCount)
	if err != nil {
		log.Println(err)
		log.Fatal("メディアの保存でエラーが発生しました")
	}
}
