# Cloudflare Turnstile導入ガイド

## 概要

このドキュメントでは、SRCGIに導入されたCloudflare Turnstileの機能について説明します。

Turnstileは、Cloudflareが提供する無料のCAPTCHA代替ソリューションで、ボット対策として効果的に機能します。現在、`/contributors`エンドポイントに実装されており、必要に応じて他のハンドラーにも適用できます。

## 実装されたファイル

1. **ServerConfig.yml** - Turnstileの設定項目を追加
2. **ShowroomCGIlib/FileIOlib.go** - ServerConfig構造体にTurnstile設定フィールドを追加
3. **ShowroomCGIlib/Turnstile.go** - Turnstile検証機能の実装
4. **ShowroomCGIlib/HandlerContributors.go** - ContributorsHandlerにTurnstile検証を統合
5. **templates/contributors.gtpl** - Turnstileウィジェット表示用のテンプレート更新
6. **my_script.env** - 環境変数設定の例を追加

## 新機能: セッション管理（2025-11-08追加）

一度検証に成功したユーザーは、一定期間内または一定回数まで再検証が不要になりました。

### セッション機能の特徴
- ✅ ユーザー体験の大幅改善（毎回チャレンジ不要）
- ✅ Cloudflare APIコール数の削減（初回のみ検証）
- ✅ 署名付きクッキーで安全に管理
- ✅ IPアドレスバインディングでセッション横取り防止

### デフォルト設定
- 有効期間: **1時間**
- 使用回数: **10回**（0=無制限も可能）

詳細は [TURNSTILE_SESSION_GUIDE.md](TURNSTILE_SESSION_GUIDE.md) を参照してください。

---

## セットアップ手順

### 1. Cloudflare Turnstileのキーを取得

既にchouette2100.com用のサイトキーとシークレットキーを取得済みの場合は、このステップをスキップできます。

新しいドメインで使用する場合:
1. [Cloudflare Dashboard](https://dash.cloudflare.com/)にログイン
2. Turnstileセクションに移動
3. 新しいサイトを追加
4. サイトキーとシークレットキーを取得

### 2. 環境変数の設定

`my_script.env`ファイルに取得したキーを設定します:

```bash
# Cloudflare Turnstile設定
export TURNSTILE_SITE_KEY=1x00000000000000000000AA  # 実際のサイトキーに置き換え
export TURNSTILE_SECRET_KEY=1x0000000000000000000000000000000AA  # 実際のシークレットキーに置き換え
```

### 3. 環境変数の読み込み

サーバー起動前に環境変数を読み込みます:

```bash
source ./my_script.env
```

### 4. ビルドと起動

```bash
go build -v .
./SRCGI
```

または:

```bash
./run.sh
```

## 動作フロー

### Contributors機能での動作

1. **初回アクセス**: ユーザーが`/contributors?ieventid=XXX&roomid=YYY`にアクセス
   - Turnstileが有効な場合、チャレンジページが表示される
   - ユーザーはCloudflareのセキュリティチェックを完了する必要がある

2. **検証**: ユーザーがチャレンジを完了して送信
   - フォームデータとTurnstileトークンがサーバーに送信される
   - サーバー側でCloudflare APIを使ってトークンを検証
   
3. **結果表示**:
   - 検証成功: 貢献ランキングデータを表示
   - 検証失敗: エラーメッセージと共に再度チャレンジページを表示

### Turnstileを無効化する場合

環境変数を設定しない、または空文字列に設定すると、Turnstile検証はスキップされ、従来通りの動作になります:

```bash
export TURNSTILE_SITE_KEY=
export TURNSTILE_SECRET_KEY=
```

## 他のハンドラーへの適用

現在は`/contributors`のみに実装されていますが、同様の方法で他のハンドラーにも適用できます:

### 適用手順

1. **ハンドラーの更新**:
   ```go
   func YourHandler(w http.ResponseWriter, r *http.Request) {
       // Turnstile検証を追加
       if ShowroomCGIlib.Serverconfig.TurnstileSiteKey != "" {
           turnstileToken := r.FormValue("cf-turnstile-response")
           if turnstileToken == "" {
               // チャレンジページを表示
               // ...
               return
           }
           
           remoteIP := ShowroomCGIlib.RemoteAddr(r)
           verified, err := ShowroomCGIlib.VerifyTurnstile(turnstileToken, remoteIP)
           if err != nil || !verified {
               // エラー処理
               // ...
               return
           }
       }
       
       // 通常の処理
       // ...
   }
   ```

2. **テンプレートの更新**:
   - `<head>`タグ内にTurnstileスクリプトを追加
   - フォーム内にTurnstileウィジェットを配置
   - データ構造体に`TurnstileSiteKey`フィールドを追加

3. **条件分岐**:
   - Turnstileトークンの有無で初回アクセスとデータ送信を判別
   - チャレンジページとデータ表示ページを適切に切り替え

## セキュリティ上の注意

1. **シークレットキーの管理**:
   - シークレットキーは絶対にGitリポジトリにコミットしない
   - 環境変数またはセキュアな設定管理ツールを使用
   - `my_script.env`は`.gitignore`に含める

2. **検証の必須化**:
   - サーバー側での検証は必須
   - クライアント側のチェックのみでは不十分
   - IPアドレスも検証に含めることでセキュリティが向上

3. **エラーログ**:
   - 検証失敗時のログを適切に記録
   - fail2banと連携して異常なアクセスパターンを検出

## トラブルシューティング

### Turnstileウィジェットが表示されない

- ブラウザのJavaScriptが有効になっているか確認
- Cloudflareスクリプトが正しく読み込まれているか確認
- サイトキーが正しく設定されているか確認

### 検証が常に失敗する

- シークレットキーが正しいか確認
- Cloudflare APIへの通信が可能か確認（ファイアウォール設定等）
- サーバーの時刻が正確か確認

### エラーログの確認

```bash
# 最新のログファイルを確認
ls -lt *.txt | head -1
tail -f <最新のログファイル>
```

## テストモード

開発環境でテストする場合、Cloudflareは特別なテストキーを提供しています:

- **常に成功**: `1x00000000000000000000AA` (サイトキー), `1x0000000000000000000000000000000AA` (シークレットキー)
- **常に失敗**: `2x00000000000000000000AB` (サイトキー), `2x0000000000000000000000000000000AB` (シークレットキー)
- **チャレンジ強制**: `3x00000000000000000000FF` (サイトキー), `3x0000000000000000000000000000000FF` (シークレットキー)

## 参考資料

- [Cloudflare Turnstile公式ドキュメント](https://developers.cloudflare.com/turnstile/)
- [Turnstile Server-side validation](https://developers.cloudflare.com/turnstile/get-started/server-side-validation/)
- [Turnstile Testing](https://developers.cloudflare.com/turnstile/troubleshooting/testing/)

## 更新履歴

- 2025-11-08: 初回実装 - `/contributors`にTurnstileを導入
