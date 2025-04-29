# microcms-backup-tool

<img width="640" alt="Screenshot" src="https://user-images.githubusercontent.com/16186206/212473375-8df10b91-27f5-488c-a579-60edf4a59fa3.png">

## 概要

microCMS で管理しているコンテンツとメディア(画像・ファイル)を取得し、保存するツールです。コンテンツはステータス（公開中、下書き、公開終了）ごとに分類して保存することができます。

## 注意事項

- 非公式ツールです。利用にあたっては、自己責任にてお願いいたします。
- メディアの取得にあたっては、ベータ版の機能であるマネジメント API (https://document.microcms.io/management-api/get-media) を利用しています。
- 利用する API キーには、あらかじめ適切な権限付与が必要です。詳しくは API キーのドキュメント (https://document.microcms.io/content-api/x-microcms-api-key) を確認してください。
- API キーの秘匿等の考慮はされていないため、取り扱いにはご注意ください。

## 利用方法

1. ルートディレクトリにて、`go run .`を実行します。

## 設定ファイル

あらかじめルートディレクトリに、`config.json`を作成し、必要情報を設定してください。

```json
{
  "target": "all",
  "serviceId": "xxxxxxxxxx",
  "contents": {
    "endpoints": ["hoge", "fuga"],
    "requestUnit": 100,
    "classifyByStatus": true,
    "getPublishContentsAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "getAllStatusContentsAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "getContentsMetaDataAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx"
  },
  "media": {
    "getMediaAPIKey": "xxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

設定されたサービスに対してバックアップを実施します。

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

#### コンテンツのステータス分類

`contents.classifyByStatus`を`true`に設定すると、コンテンツは以下のように分類されて保存されます：

- `PUBLISH`: 公開中のコンテンツ
- `DRAFT`: 下書きのコンテンツ
- `CLOSED`: 公開終了のコンテンツ
- `PUBLISH_AND_DRAFT`: 公開中かつ下書きのコンテンツ（両方の状態で保存）

3. `backup`フォルダの中に、ファイルが保存されます。コンテンツはステータスごとのディレクトリに分類されて保存されます。
