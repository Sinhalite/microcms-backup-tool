package client

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	// CSVファイルとして保存する場合
	if c.Config.Contents.SaveAsCSV {
		return c.saveContentsAsCSV(endpoint, requiredRequestCount, baseDir, apiKey, status)
	}

	// 従来のJSONファイルとして保存する場合
	for i := 0; i < requiredRequestCount; i++ {
		// 1秒のディレイを追加
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

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

// saveContentsAsCSV はコンテンツをCSVファイルとして保存する関数
func (c Client) saveContentsAsCSV(endpoint string, requiredRequestCount int, baseDir string, apiKey string, status string) error {
	// 保存先ディレクトリを作成
	dir, err := makeSaveDir(baseDir, endpoint, status, "")
	if err != nil {
		return err
	}

	// CSVファイルを作成
	csvFile, err := os.Create(fmt.Sprintf("%s/contents.csv", dir))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	// CSVライターを作成
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// すべてのコンテンツで共通のカラムを収集
	allKeys := make(map[string]bool)
	var allContents []gjson.Result
	var orderedKeys []string

	// まずすべてのコンテンツを取得して、存在するすべてのキーを収集
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

		// 各コンテンツのキーを収集
		for _, item := range contents.Array() {
			allContents = append(allContents, item)
			// 最初のコンテンツからキーの順序を取得
			if len(orderedKeys) == 0 {
				item.ForEach(func(key, value gjson.Result) bool {
					keyStr := key.String()
					if !allKeys[keyStr] {
						orderedKeys = append(orderedKeys, keyStr)
						allKeys[keyStr] = true
					}
					return true
				})
			} else {
				// 2つ目以降のコンテンツでは、新しいキーのみを追加
				item.ForEach(func(key, value gjson.Result) bool {
					keyStr := key.String()
					if !allKeys[keyStr] {
						orderedKeys = append(orderedKeys, keyStr)
						allKeys[keyStr] = true
					}
					return true
				})
			}
		}

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, requiredRequestCount, requestURL)
	}

	// ヘッダー行を書き込む
	if err := writer.Write(orderedKeys); err != nil {
		return err
	}

	// 各コンテンツのデータを書き込む
	for _, item := range allContents {
		row := make([]string, len(orderedKeys))
		for i, key := range orderedKeys {
			value := item.Get(key)
			// 値がオブジェクトや配列の場合はJSON文字列として保存
			if value.IsObject() || value.IsArray() {
				row[i] = value.Raw
			} else {
				row[i] = value.String()
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func (c Client) saveContentsWithStatus(endpoint string, requiredRequestCount int, baseDir string) error {
	// すべてのコンテンツで共通のカラムを収集
	allKeys := make(map[string]bool)
	var allContents []gjson.Result
	var orderedKeys []string

	// まずすべてのコンテンツを取得して、存在するすべてのキーを収集
	for i := 0; i < requiredRequestCount; i++ {
		// 1秒のディレイを追加
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

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

		// 各コンテンツのキーを収集
		for j := 0; j < len(contents.Array()); j++ {
			item := contents.Array()[j]
			allContents = append(allContents, item)
			// 最初のコンテンツからキーの順序を取得
			if len(orderedKeys) == 0 {
				item.ForEach(func(key, value gjson.Result) bool {
					keyStr := key.String()
					if !allKeys[keyStr] {
						orderedKeys = append(orderedKeys, keyStr)
						allKeys[keyStr] = true
					}
					return true
				})
			} else {
				// 2つ目以降のコンテンツでは、新しいキーのみを追加
				item.ForEach(func(key, value gjson.Result) bool {
					keyStr := key.String()
					if !allKeys[keyStr] {
						orderedKeys = append(orderedKeys, keyStr)
						allKeys[keyStr] = true
					}
					return true
				})
			}
		}

		// 進捗状況の表示
		fmt.Printf("[%d / %d] %s\n", i+1, requiredRequestCount, requestURL)
	}

	// ステータスごとにコンテンツを分類
	statusContents := make(map[string][]gjson.Result)
	for i := 0; i < requiredRequestCount; i++ {
		// 1秒のディレイを追加
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

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
			return fmt.Errorf("ステータスコード:%d 正常にレスポンスを取得できませんでした", mResp.StatusCode)
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

			switch status {
			case "PUBLISH", "DRAFT", "CLOSED":
				statusContents[status] = append(statusContents[status], item)
			case "PUBLISH_AND_DRAFT":
				// 下書き保存
				statusContents["DRAFT"] = append(statusContents["DRAFT"], item)
				// 公開中データ取得
				publishItem, err := c.getContentWithGJSON(endpoint, c.Config.Contents.GetPublishContentsAPIKey, id)
				if err != nil {
					log.Fatalf("公開中かつ下書き中コンテンツにおいて、公開中のコンテンツの取得に失敗しました: %v", err)
				}
				statusContents["PUBLISH"] = append(statusContents["PUBLISH"], publishItem)
			default:
				fmt.Println("未知のステータスです")
			}
		}
	}

	// 各ステータスごとにCSVファイルを作成
	for status, contents := range statusContents {
		// 保存先ディレクトリを作成
		dir, err := makeSaveDir(baseDir, endpoint, status, "")
		if err != nil {
			return err
		}

		if c.Config.Contents.SaveAsCSV {
			// CSVファイルを作成
			csvFile, err := os.Create(fmt.Sprintf("%s/contents.csv", dir))
			if err != nil {
				return err
			}
			defer csvFile.Close()

			// CSVライターを作成
			writer := csv.NewWriter(csvFile)
			defer writer.Flush()

			// ヘッダー行を書き込む
			if err := writer.Write(orderedKeys); err != nil {
				return err
			}

			// 各コンテンツのデータを書き込む
			for _, item := range contents {
				row := make([]string, len(orderedKeys))
				for i, key := range orderedKeys {
					value := item.Get(key)
					// 値がオブジェクトや配列の場合はJSON文字列として保存
					if value.IsObject() || value.IsArray() {
						row[i] = value.Raw
					} else {
						row[i] = value.String()
					}
				}
				if err := writer.Write(row); err != nil {
					return err
				}
			}
		} else {
			// JSONファイルとして保存
			for i, item := range contents {
				err := c.writeRawJSONWithStatus(item.Raw, baseDir, endpoint, i+1, status, "")
				if err != nil {
					return err
				}
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
