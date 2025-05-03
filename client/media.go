package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func (c Client) BackupMedia(baseDir string) error {
	log.Println("メディアのバックアップを開始します")
	const requestUnit = 50
	totalCount, err := c.getTotalCount()
	if err != nil {
		return fmt.Errorf("合計件数の取得でエラーが発生しました: %w", err)
	}
	requiredRequestCount := (totalCount/requestUnit + 1)

	mediaAry, err := c.getAllMedia(requiredRequestCount, requestUnit)
	if err != nil {
		return fmt.Errorf("メディア一覧の取得でエラーが発生しました: %w", err)
	}
	err = c.saveMedia(mediaAry, totalCount, baseDir)
	if err != nil {
		return fmt.Errorf("メディアの保存でエラーが発生しました: %w", err)
	}
	return nil
}

func (c Client) getTotalCount() (int, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms-management.io/api/v2/media?limit=0", c.Config.ServiceID),
		nil)
	req.Header.Set("X-MICROCMS-API-KEY", c.Config.Media.APIKey)

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

func (c Client) getAllMedia(requiredRequestCount int, requestUnit int) ([]Media, error) {
	var ary []Media
	var token string

	for i := 0; i < requiredRequestCount; i++ {
		// 1秒のディレイを追加
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

		client := new(http.Client)
		req, _ := http.NewRequest(
			"GET",
			fmt.Sprintf("https://%s.microcms-management.io/api/v2/media?limit=%d&token=%s", c.Config.ServiceID, requestUnit, token),
			nil,
		)
		req.Header.Set("X-MICROCMS-API-KEY", c.Config.Media.APIKey)
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
		token = response.Token
	}

	return ary, nil
}

func (c Client) saveMedia(medias []Media, totalCount int, baseDir string) error {
	for i, media := range medias {
		// 1秒のディレイを追加
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, totalCount, media.Url)

		client := new(http.Client)
		req, _ := http.NewRequest("GET", media.Url, nil)
		req.Header.Set("X-MICROCMS-API-KEY", c.Config.Media.APIKey)

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
		fileName, err = url.QueryUnescape(fileName)
		if err != nil {
			return err
		}
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
	}
	return nil
}
