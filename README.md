# microcms-backup-tool
<img width="640" alt="Screenshot" src="https://user-images.githubusercontent.com/16186206/212473375-8df10b91-27f5-488c-a579-60edf4a59fa3.png">

## 概要
microCMSで管理しているコンテンツとメディア(画像・ファイル)を取得し、保存するツールです。

## 注意事項
- 非公式ツールです。利用にあたっては、自己責任にてお願いいたします。
- メディアの取得にあたっては、ベータ版の機能であるマネジメントAPI (https://document.microcms.io/management-api/get-media) を利用しています。
- 利用するAPIキーには、あらかじめ`GET`、`メディアの取得`の権限付与が必要です。詳しくはAPIキーのドキュメント (https://document.microcms.io/content-api/x-microcms-api-key) を確認してください。
- APIキーの秘匿等の考慮はされていないため、取り扱いにはご注意ください。

## 利用方法
1. ルートディレクトリにて、`go run .`を実行します。
2. `> モードを選択してください(auto / manual)`と表示されるので、`auto`もしくは`manual`を入力します。

### autoモードを利用する場合

あらかじめルートディレクトリに、`config.json`を作成し、必要情報を設定してください。
```json
{
  "target": "all",
  "serviceId": "xxxxxxxxxx",
  "apiKey": "xxxxxxxxxxxxxxxxxxxxxxxx"
  "endpoints": ["hoge", "fuga"]
}
```
設定されたサービスに対してバックアップを実施します。

`target`は、以下の3項目より選択してください。

- `all` : コンテンツとメディア
- `contents` : コンテンツのみ
- `media` : メディアのみ

### manualモードを利用する場合

対話モードにて、必要な項目を聞かれるので、それぞれ必要な値を入力します。

3. `backup`フォルダの中に、ファイルが保存されます。
