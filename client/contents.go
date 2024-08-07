package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (c Client) BackupContents(baseDir string) error {
	log.Println("コンテンツのバックアップを開始します")

	for _, endpoint := range c.Config.Endpoints {
		log.Printf("%sのバックアップを開始します\n", endpoint)

		// 1:ステータスごとの分類を行う場合
		if c.Config.ClassifyByStatus {
			fmt.Println("コンテンツの処理を開始しました")
			// 全コンテンツの合計件数を取得
			allCotentsCount, err := c.getContentsTotalCount(endpoint, c.Config.DraftAndClosedAPIKey)
			if err != nil {
				return fmt.Errorf("全コンテンツの合計件数の取得でエラーが発生しました: %w", err)
			}

			// 必要なリクエスト回数を計算
			requiredRequestCount := (allCotentsCount/c.Config.RequestUnit + 1)

			// 全コンテンツの取得した後、ステータスごとにデータを振り分けて保存する
			err = c.saveContentsWithStatus(endpoint, requiredRequestCount, baseDir)
			if err != nil {
				return fmt.Errorf("コンテンツの保存でエラーが発生しました: %w", err)
			}
		} else {
			// 2:ステータスごとの分類を行わない場合
			totalCount, err := c.getContentsTotalCount(endpoint, c.Config.APIKey)
			if err != nil {
				return fmt.Errorf("コンテンツの合計件数の取得でエラーが発生しました: %w", err)
			}
			requiredRequestCount := (totalCount/c.Config.RequestUnit + 1)

			err = c.saveContents(endpoint, requiredRequestCount, baseDir, c.Config.APIKey, "PUBLISH")
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

func (c Client) getContent(endpoint string, apiKey string, contentId string) (map[string]interface{}, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s.microcms.io/api/v1/%s/%s", c.Config.ServiceID, endpoint, contentId),
		nil)
	req.Header.Set("X-MICROCMS-API-KEY", apiKey)

	client := new(http.Client)
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

	response := &map[string]interface{}{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return *response, err
}

func (c Client) saveContents(endpoint string, requiredRequestCount int, baseDir string, apiKey string, status string) error {
	for i := 0; i < requiredRequestCount; i++ {
		client := new(http.Client)
		requestURL := fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.RequestUnit, c.Config.RequestUnit*i)
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

		// JSONのフォーマット
		json, err := formatJson(resp)
		if err != nil {
			return err
		}

		// 保存用ディレクトリの作成
		dir, err := makeSaveDir(baseDir, endpoint, status, "")
		if err != nil {
			return err
		}

		f, err := os.Create(fmt.Sprintf("%s/%d.json", dir, i+1))
		if err != nil {
			return err
		}

		_, err = f.WriteString(json)
		if err != nil {
			return err
		}

		defer f.Close()

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, requiredRequestCount, requestURL)
	}

	return nil
}

func (c Client) saveContentsWithStatus(endpoint string, requiredRequestCount int, baseDir string) error {
	for i := 0; i < requiredRequestCount; i++ {
		// コンテンツAPIから取得
		client := new(http.Client)
		requestURL := fmt.Sprintf("https://%s.microcms.io/api/v1/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.RequestUnit, c.Config.RequestUnit*i)
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("X-MICROCMS-API-KEY", c.Config.DraftAndClosedAPIKey)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", resp.StatusCode)
		}
		defer resp.Body.Close()

		// マネジメントAPIから取得
		mRequestURL := fmt.Sprintf("https://%s.microcms-management.io/api/v1/contents/%s?limit=%d&offset=%d", c.Config.ServiceID, endpoint, c.Config.RequestUnit, c.Config.RequestUnit*i)
		mReq, _ := http.NewRequest("GET", mRequestURL, nil)
		mReq.Header.Set("X-MICROCMS-API-KEY", c.Config.ManagementAPIKey)
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

		// JSONをデコードする
		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// レスポンスボディを読み込む
		mbody, err := io.ReadAll(mResp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		// JSONをデコードする
		var mdata map[string]interface{}
		if err := json.Unmarshal(mbody, &mdata); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// 振り分け処理
		contents := data["contents"].([]interface{})
		mContents := mdata["contents"].([]interface{})

		for j, content := range contents {
			item := content.(map[string]interface{})
			id := item["id"].(string)
			mid := mContents[j].(map[string]interface{})["id"].(string)

			if id != mid {
				return fmt.Errorf("コンテンツIDが一致しませんでした:%s,%s", id, mid)
			}

			status := mContents[j].(map[string]interface{})["status"].([]interface{})[0].(string)

			number := i*c.Config.RequestUnit + j + 1

			fmt.Println(number, status)

			switch status {
			case "PUBLISH", "DRAFT", "CLOSED":
				c.writeContentsWithStatus(item, baseDir, endpoint, number, status, "")
			case "PUBLISH_AND_DRAFT":
				// 公開中かつ下書き中の時は、公開中のデータと下書きのデータをそれぞれ保存する
				// 1:下書きのデータを保存
				c.writeContentsWithStatus(item, baseDir, endpoint, number, "DRAFT", "PUBLISH_AND_DRAFT")

				// 2:公開中のデータを取得、保存
				// 公開中のデータを取得するたびに、「下書き全取得」が付与されていないAPIキーを利用する
				publishItem, err := c.getContent(endpoint, c.Config.APIKey, id)
				if err != nil {
					log.Fatalf("公開中かつ下書き中コンテンツにおいて、公開中のコンテンツの取得に失敗しました: %v", err)
				}
				c.writeContentsWithStatus(publishItem, baseDir, endpoint, number, "PUBLISH", "")
			default:
				fmt.Println("未知のステータスです")
			}
		}
	}
	return nil
}

func (c Client) writeContentsWithStatus(item map[string]interface{}, baseDir string, endpoint string, number int, status string, draftStatusDetail string) error {
	jsonData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	dir, err := makeSaveDir(baseDir, endpoint, status, draftStatusDetail)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%d.json", dir, number))
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	defer f.Close()
	return nil
}

func formatJson(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(body), "", "  ")
	return buf.String(), err
}

func makeSaveDir(baseDir string, endpoint string, status string, draftStatusDetail string) (string, error) {
	contentsDir := baseDir + "contents/" + endpoint + "/" + status + "/" + draftStatusDetail
	err := os.MkdirAll(contentsDir, os.ModePerm)
	return contentsDir, err
}
