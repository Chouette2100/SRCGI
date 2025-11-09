# Turnstile導入の変更サマリー

## 追加されたファイル

1. **ShowroomCGIlib/Turnstile.go** (新規)
   - `VerifyTurnstile()`: Cloudflare APIを使用してトークン検証
   - `TurnstileResponse`: APIレスポンスの構造体

2. **TURNSTILE_README.md** (新規)
   - セットアップ手順
   - 使用方法
   - トラブルシューティング

## 変更されたファイル

### 1. ServerConfig.yml
```yaml
# Cloudflare Turnstile設定
TurnstileSiteKey: ${TURNSTILE_SITE_KEY}
TurnstileSecretKey: ${TURNSTILE_SECRET_KEY}
```

### 2. ShowroomCGIlib/FileIOlib.go
ServerConfig構造体に追加:
```go
TurnstileSiteKey    string `yaml:"TurnstileSiteKey"`
TurnstileSecretKey  string `yaml:"TurnstileSecretKey"`
```

### 3. ShowroomCGIlib/HandlerContributors.go
- HCntrbInf構造体に`TurnstileSiteKey`と`TurnstileError`フィールドを追加
- Turnstile検証ロジックを追加:
  - トークンなし → チャレンジページ表示
  - トークンあり → 検証実行
  - 検証成功 → データ表示
  - 検証失敗 → エラー付きチャレンジページ再表示

### 4. templates/contributors.gtpl
- Turnstileスクリプトの読み込みを追加
- チャレンジフォームを追加:
  - Turnstileウィジェット
  - 隠しフィールド（ieventid, roomid）
  - 送信ボタン
- 条件分岐でチャレンジページとデータページを切り替え

### 5. my_script.env
環境変数設定例を追加:
```bash
export TURNSTILE_SITE_KEY=your_site_key_here
export TURNSTILE_SECRET_KEY=your_secret_key_here
```

## 使用方法

### 本番環境でのセットアップ

1. 環境変数を設定:
   ```bash
   vi my_script.env
   # TURNSTILE_SITE_KEYとTURNSTILE_SECRET_KEYに実際の値を設定
   ```

2. 環境変数を読み込んでビルド:
   ```bash
   source ./my_script.env
   go build -v .
   ```

3. サーバーを起動:
   ```bash
   ./SRCGI
   ```
   または
   ```bash
   systemctl --user restart SRCGI
   ```

### Turnstileを無効化する場合

環境変数を空にするか、コメントアウト:
```bash
# export TURNSTILE_SITE_KEY=
# export TURNSTILE_SECRET_KEY=
```

この場合、従来通りの動作(Turnstileなし)となります。

## 動作確認

1. ブラウザで `/contributors?ieventid=XXX&roomid=YYY` にアクセス
2. Turnstileチャレンジが表示されることを確認
3. チェックボックスをクリックして検証
4. 「確認して続行」ボタンをクリック
5. 貢献ランキングが表示されることを確認

## セキュリティ強化ポイント

1. **ボット対策**: 分散したIPからの自動アクセスをTurnstileで検出
2. **偽装対策**: User-Agent偽装だけでは突破できない
3. **段階的展開**: `/contributors`で実績を確認後、他のエンドポイントにも適用可能
4. **無効化可能**: 環境変数で簡単にON/OFF切り替え

## 今後の展開

必要に応じて以下のハンドラーにも適用検討:
- `/list-last`: 直近の獲得ポイント
- `/graph-total`: 獲得ポイントグラフ
- その他、負荷が高いエンドポイント

適用方法は`TURNSTILE_README.md`の「他のハンドラーへの適用」セクションを参照。
