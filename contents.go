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

func (c *Client) backupContents(baseDir string) error {
	log.Println("コンテンツのバックアップを開始します")

	for _, endpoint := range c.Config.Endpoints {
		log.Printf("%sのバックアップを開始します\n", endpoint)

		totalCount, err := c.getContentsTotalCount(endpoint)
		if err != nil {
			return fmt.Errorf("コンテンツの合計件数の取得でエラーが発生しました: %w", err)
		}
		requiredRequestCount := (totalCount/c.Config.RequestUnit + 1)

		err = c.saveContents(endpoint, requiredRequestCount, baseDir)
		if err != nil {
			return fmt.Errorf("コンテンツの保存でエラーが発生しました: %w", err)
		}
	}
	return nil
}

func (c *Client) getContentsTotalCount(endpoint string) (int, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=0", c.Config.ServiceID, endpoint),
		nil)
	req.Header.Set("X-MICROCMS-API-KEY", c.Config.APIKey)

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

func (c *Client) saveContents(endpoint string, requiredRequestCount int, baseDir string) error {
	for i := 0; i < requiredRequestCount; i++ {
		client := new(http.Client)
		requestURL := fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.RequestUnit, c.Config.RequestUnit*i)
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("X-MICROCMS-API-KEY", c.Config.APIKey)
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
