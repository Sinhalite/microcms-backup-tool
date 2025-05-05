# microcms-backup-tool

<img width="640" alt="Screenshot" src="https://user-images.githubusercontent.com/16186206/212473375-8df10b91-27f5-488c-a579-60edf4a59fa3.png">

## 概要

microCMS で管理しているコンテンツとメディア(画像・ファイル)を取得し、保存するツールです。

## 注意事項

- 非公式ツールです。利用にあたっては、自己責任にてお願いいたします。
- 一部は動作保証のないベータ版の機能であるマネジメントAPI (https://document.microcms.io/management-api/get-media) を利用しています。
- 利用するAPIキーには、あらかじめ適切な権限付与が必要です。詳しくは API キーのドキュメント (https://document.microcms.io/content-api/x-microcms-api-key) を確認してください。
- APIキーの秘匿等の考慮はされていないため、取り扱いにはご注意ください。

## 利用方法

1. ルートディレクトリに、`config.json`を作成し、必要情報を設定してください。
2. バックアップ対象のサービスにおいて、適切なAPIキーの権限付与を行います。
3. ルートディレクトリにて、`go run .`を実行します。
4. `backup`フォルダの中に、指定したデータのバックアップファイルが保存されます。

## 設定ファイル

`config.json`
```json
{
  "target": "all",
  "serviceId": "xxxxxxxxxx",
  "contents": {
    "getPublishContentsAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "getAllStatusContentsAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "getContentsMetaDataAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "endpoints": ["hoge", "fuga"],
    "requestUnit": 100,
    "classifyByStatus": true,
    "saveAsCSV": false
  },
  "media": {
    "apiKey": "xxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

設定されたサービスに対してバックアップを実施します。

#### targetの説明
`target`は、以下の 3 項目より選択してください。

- `all` : コンテンツとメディア
- `contents` : コンテンツのみ
- `media` : メディアのみ

#### APIキーの説明

`contents.getPublishContentsAPIKey`
- 公開中のみのコンテンツのGET権限を付与してください
- 公開中コンテンツの取得に使用

`contents.getAllStatusContentsAPIKey`
- 公開中、下書き、公開終了のコンテンツのGET権限を付与してください
- 全ステータスのコンテンツ取得に使用

`contents.getContentsMetaDataAPIKey`
- コンテンツのメタデータのGET権限を付与してください
- コンテンツのステータス情報取得に使用

`media.getMediaAPIKey`
- メディアのGET権限を付与してください
- メディアファイルの取得に使用

#### コンテンツの保存形式

### 1. ステータス別分類（`classifyByStatus: true`）

コンテンツはステータスごとに分類されて保存されます。

#### JSON形式（`saveAsCSV: false`）
- コンテンツは個別のJSONファイルとして保存されます

#### CSV形式（`saveAsCSV: true`）
- コンテンツは1つのCSVファイルとして保存されます
- ネストされたJSONオブジェクトや配列は文字列として保存されます
- ファイル名は`contents.csv`となります

### 2. ステータス別分類なし（`classifyByStatus: false`）

コンテンツは1つのファイルとして保存されます。

#### JSON形式（`saveAsCSV: false`）
- コンテンツは個別のJSONファイルとして保存されます

#### CSV形式（`saveAsCSV: true`）
- コンテンツは1つのCSVファイルとして保存されます
- ネストされたJSONオブジェクトや配列は文字列として保存されます
- ファイル名は`contents.csv`となります
