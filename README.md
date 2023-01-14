# microcms-backup-tool
## 概要
microCMSで管理しているメディア(画像・ファイル)を取得し、保存するツールです。

## 注意事項
- 非公式ツールです。利用にあたっては、自己責任にてお願いいたします。
- メディアの取得にあたっては、ベータ版の機能であるマネジメントAPI (https://document.microcms.io/management-api/get-media) を利用しています。
- 利用するAPIキーには、あらかじめ`メディアの取得`の権限付与が必要です。詳しくはAPIキーのドキュメント (https://document.microcms.io/content-api/x-microcms-api-key) を確認してください。
- APIキーの秘匿等の考慮はされていないため、取り扱いにはご注意ください。

## 利用方法
1. `go run main.go`を実行します。
2. `> モードを選択してください(auto / manual)`と表示されるので、`auto`もしくは`manual`を入力します。
- **autoモード**を利用する場合

あらかじめルートディレクトリに、`config.json`を作成し、必要情報を設定してください。
```json
{
  "serviceId": "xxxxxxxxxx",
  "apiKey": "xxxxxxxxxxxxxxxxxxxxxxxx"
}
```
設定されたサービスに対してバックアップを実施します。

- **manualモード**を利用する場合

`> サービスIDを入力してください`、`> APIキーを入力してください`とそれぞれ表示されるので、
それぞれ任意のサービスの値を入力してください。

3. `backup`フォルダの中にファイルが保存されます。
