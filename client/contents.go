package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

func (c Client) BackupContents(baseDir string) error {
	log.Println("コンテンツのバックアップを開始します")

	for _, endpoint := range c.Config.Contents.Endpoints {
		log.Printf("%sのバックアップを開始します\n", endpoint)

		// 1:ステータスごとの分類を行う場合
		if c.Config.Contents.ClassifyByStatus {
			fmt.Println("コンテンツの処理を開始しました")
			// 全コンテンツの合計件数を取得
			allCotentsCount, err := c.getContentsTotalCount(endpoint, c.Config.Contents.GetAllStatusContentsAPIKey)
			if err != nil {
				return fmt.Errorf("全コンテンツの合計件数の取得でエラーが発生しました: %w", err)
			}

			// 必要なリクエスト回数を計算
			requiredRequestCount := (allCotentsCount/c.Config.Contents.RequestUnit + 1)

			// 全コンテンツの取得した後、ステータスごとにデータを振り分けて保存する
			err = c.saveContentsWithStatus(endpoint, requiredRequestCount, baseDir)
			if err != nil {
				return fmt.Errorf("コンテンツの保存でエラーが発生しました: %w", err)
			}
		} else {
			// 2:ステータスごとの分類を行わない場合
			totalCount, err := c.getContentsTotalCount(endpoint, c.Config.Contents.GetPublishContentsAPIKey)
			if err != nil {
				return fmt.Errorf("コンテンツの合計件数の取得でエラーが発生しました: %w", err)
			}
			requiredRequestCount := (totalCount/c.Config.Contents.RequestUnit + 1)

			err = c.saveContents(endpoint, requiredRequestCount, baseDir, c.Config.Contents.GetPublishContentsAPIKey, "PUBLISH")
			if err != nil {
				return fmt.Errorf("コンテンツの保存でエラーが発生しました: %w", err)
			}
		}
	}
	return nil
}

func (c Client) getContentsTotalCount(endpoint string, apiKey string) (int, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=0", c.Config.ServiceID, endpoint),
		nil)
	req.Header.Set("X-MICROCMS-API-KEY", apiKey)

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

func (c Client) saveContents(endpoint string, requiredRequestCount int, baseDir string, apiKey string, status string) error {
	for i := 0; i < requiredRequestCount; i++ {
		client := new(http.Client)
		requestURL := fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.Contents.RequestUnit, c.Config.Contents.RequestUnit*i)
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("X-MICROCMS-API-KEY", apiKey)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer resp.Body.Close()

		// レスポンスボディを読み込む
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// gjsonでcontents配列を取得
		contents := gjson.GetBytes(body, "contents")
		if !contents.IsArray() {
			return fmt.Errorf("contentsが配列ではありません")
		}

		for j, item := range contents.Array() {
			number := i*c.Config.Contents.RequestUnit + j + 1
			// item.Rawで元の順序のままJSON文字列が得られる
			err := c.writeRawJSONWithStatus(item.Raw, baseDir, endpoint, number, status, "")
			if err != nil {
				return err
			}
		}

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, requiredRequestCount, requestURL)
	}

	return nil
}

func (c Client) saveContentsWithStatus(endpoint string, requiredRequestCount int, baseDir string) error {
	for i := 0; i < requiredRequestCount; i++ {
		// コンテンツAPIから取得
		client := new(http.Client)
		requestURL := fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.Contents.RequestUnit, c.Config.Contents.RequestUnit*i)
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("X-MICROCMS-API-KEY", c.Config.Contents.GetAllStatusContentsAPIKey)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer resp.Body.Close()

		// マネジメントAPIから取得
		mRequestURL := fmt.Sprintf("https://%s.microcms-management.io/api/v1/contents/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.Contents.RequestUnit, c.Config.Contents.RequestUnit*i)
		mReq, _ := http.NewRequest("GET", mRequestURL, nil)
		mReq.Header.Set("X-MICROCMS-API-KEY", c.Config.Contents.GetContentsMetaDataAPIKey)
		mResp, err := client.Do(mReq)
		if err != nil {
			return err
		}
		if mResp.StatusCode != http.StatusOK {
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer mResp.Body.Close()

		// レスポンスボディを読み込む
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		mbody, err := io.ReadAll(mResp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// gjsonでcontents配列を取得
		contents := gjson.GetBytes(body, "contents")
		mContents := gjson.GetBytes(mbody, "contents")

		if !contents.IsArray() || !mContents.IsArray() {
			return fmt.Errorf("contentsが配列ではありません")
		}

		for j := 0; j < len(contents.Array()); j++ {
			item := contents.Array()[j]
			mItem := mContents.Array()[j]

			id := item.Get("id").String()
			mid := mItem.Get("id").String()
			if id != mid {
				return fmt.Errorf("コンテンツIDが一致しませんでした:%s,%s", id, mid)
			}

			status := mItem.Get("status.0").String()
			number := i*c.Config.Contents.RequestUnit + j + 1

			fmt.Println(number, status)

			switch status {
			case "PUBLISH", "DRAFT", "CLOSED":
				// item.Rawで元の順序のままJSON文字列が得られる
				c.writeRawJSONWithStatus(item.Raw, baseDir, endpoint, number, status, "")
			case "PUBLISH_AND_DRAFT":
				// 下書き保存
				c.writeRawJSONWithStatus(item.Raw, baseDir, endpoint, number, "DRAFT", "PUBLISH_AND_DRAFT")
				// 公開中データ取得
				publishItem, err := c.getContentWithGJSON(endpoint, c.Config.Contents.GetPublishContentsAPIKey, id)
				if err != nil {
					log.Fatalf("公開中かつ下書き中コンテンツにおいて、公開中のコンテンツの取得に失敗しました: %v", err)
				}
				c.writeRawJSONWithStatus(publishItem.Raw, baseDir, endpoint, number, "PUBLISH", "")
			default:
				fmt.Println("未知のステータスです")
			}
		}
	}
	return nil
}

// itemRawはJSON文字列
func (c Client) writeRawJSONWithStatus(itemRaw string, baseDir, endpoint string, number int, status, draftStatusDetail string) error {
	// JSONを整形
	formattedJson, err := formatJson(itemRaw)
	if err != nil {
		return err
	}

	dir, err := makeSaveDir(baseDir, endpoint, status, draftStatusDetail)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%d.json", dir, number))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(formattedJson)
	return err
}

// 公開中データ取得用
func (c Client) getContentWithGJSON(endpoint, apiKey, contentId string) (gjson.Result, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms.io/api/v1/%s/%s", c.Config.ServiceID, endpoint, contentId),
		nil)
	req.Header.Set("X-MICROCMS-API-KEY", apiKey)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return gjson.Result{}, fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(body), nil
}

func formatJson(rawJson string) (string, error) {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(rawJson), "", "  ")
	return buf.String(), err
}

func makeSaveDir(baseDir string, endpoint string, status string, draftStatusDetail string) (string, error) {
	contentsDir := baseDir + "contents/" + endpoint + "/" + status + "/" + draftStatusDetail
	err := os.MkdirAll(contentsDir, os.ModePerm)
	return contentsDir, err
}
