package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func backupMedia(option Config, baseDir string) {
	log.Println("メディアのバックアップを開始します")
	const requestUnit = 50
	totalCount, err := getTotalCount(option)
	if err != nil {
		log.Println(err)
		log.Fatal("合計件数の取得でエラーが発生しました")
	}
	requiredRequestCount := (totalCount/requestUnit + 1)

	mediaAry, err := getMediaAry(option, requiredRequestCount, requestUnit)
	if err != nil {
		log.Println(err)
		log.Fatal("メディア一覧の取得でエラーが発生しました")
	}
	err = saveMedia(mediaAry, option, totalCount, baseDir)
	if err != nil {
		log.Println(err)
		log.Fatal("メディアの保存でエラーが発生しました")
	}
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

	response := &ManagementAPIMediaResponse{}
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

		response := &ManagementAPIMediaResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return nil, err
		}

		ary = append(ary, response.Media...)
	}

	return ary, nil
}

func saveMedia(mediaAry []Media, option Config, totalCount int, baseDir string) error {
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
		err = os.MkdirAll(baseDir+"media/"+fileDirectory, os.ModePerm)
		if err != nil {
			return err
		}

		file, err := os.Create(baseDir + "media/" + fileDirectory + "/" + fileName)
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
