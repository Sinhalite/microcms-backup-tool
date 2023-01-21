package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func backupContents(option Config, baseDir string) {
	const requestUnit = 10

	for _, endpoint := range option.Endpoints {
		totalCount, err := getContentsTotalCount(option, endpoint)
		if err != nil {
			log.Println(err)
			log.Fatal("コンテンツの合計件数の取得でエラーが発生しました")
		}
		requiredRequestCount := (totalCount/requestUnit + 1)

		err = saveContents(option, endpoint, requiredRequestCount, requestUnit, baseDir)
		if err != nil {
			log.Println(err)
			log.Fatal("コンテンツの保存でエラーが発生しました")
		}
	}
}

func getContentsTotalCount(option Config, endpoint string) (int, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=0", option.ServiceID, endpoint),
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

	response := &ContentsAPIResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return 0, err
	}

	return response.TotalCount, err
}

func saveContents(option Config, endpoint string, requiredRequestCount int, requestUnit int, baseDir string) error {
	for i := 0; i < requiredRequestCount; i++ {
		client := new(http.Client)
		requestURL := fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=%d&offset=%d", option.ServiceID, endpoint, requestUnit, requestUnit*i)
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("X-MICROCMS-API-KEY", option.APIKey)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// JSONのフォーマット
		var buf bytes.Buffer
		err = json.Indent(&buf, []byte(body), "", "  ")
		if err != nil {
			return err
		}

		contentsDir := baseDir + "contents/" + endpoint
		err = os.MkdirAll(contentsDir, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(fmt.Sprintf("%s/%d.json", contentsDir, i+1))
		if err != nil {
			return err
		}

		_, err = f.WriteString(buf.String())
		if err != nil {
			return err
		}

		defer f.Close()

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, requiredRequestCount, requestURL)
	}

	return nil
}
