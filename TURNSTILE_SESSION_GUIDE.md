# Turnstileセッション管理機能 - 使用ガイド

## 概要

Turnstileセッション管理機能により、一度検証に成功したユーザーは一定期間内または一定回数まで再検証が不要になります。これにより、ユーザー体験が大幅に向上します。

## 実装方法

### クッキーベースのセッション管理
- 検証済みセッション情報を**署名付きクッキー**で管理
- HMAC-SHA256による改ざん防止
- IPアドレスバインディングでセッション横取り防止
- 有効期限と使用回数制限をサポート

## 設定方法

### ServerConfig.yml
```yaml
# Cloudflare Turnstile設定
TurnstileSiteKey: ${TURNSTILE_SITE_KEY}
TurnstileSecretKey: ${TURNSTILE_SECRET_KEY}

# Turnstileセッション管理
TurnstileUseSession: true       # セッション機能の有効/無効
TurnstileSessionDuration: 3600  # セッション有効期間（秒）デフォルト1時間
TurnstileSessionMaxUses: 10     # セッション内の最大利用回数（0=無制限）
```

### 設定パラメータ詳細

#### TurnstileUseSession (デフォルト: true)
- `true`: セッション機能を有効化（推奨）
- `false`: 毎回Turnstile検証を実行（最高セキュリティ）

#### TurnstileSessionDuration (デフォルト: 3600秒 = 1時間)
セッションの有効期間を秒単位で指定。

**推奨値:**
- 一般的なWebサイト: `3600` (1時間)
- セキュリティ重視: `900` (15分)
- 開発/テスト: `86400` (24時間)

#### TurnstileSessionMaxUses (デフォルト: 10回)
セッション内での最大利用回数。

**推奨値:**
- `0`: 無制限（有効期限のみで制御）
- `10`: 適度な制限
- `5`: セキュリティ重視

## 動作フロー

### 初回アクセス
```
1. ユーザーが /contributors にアクセス
2. セッションクッキーなし → Turnstileチャレンジを表示
3. ユーザーが検証を完了
4. サーバーがCloudflare APIで検証
5. 成功 → セッションクッキーを発行 + データ表示
```

### 2回目以降のアクセス（セッション有効期間内）
```
1. ユーザーが /contributors にアクセス
2. セッションクッキーあり → 検証
   - 署名チェック（改ざん検証）
   - IPアドレス確認（横取り防止）
   - 有効期限確認
   - 使用回数確認
3. すべてOK → Turnstileスキップ + データ表示
4. クッキーの使用回数をインクリメント
```

### セッション期限切れ/上限到達時
```
1. ユーザーが /contributors にアクセス
2. セッションクッキーの検証失敗
3. Turnstileチャレンジを再表示
4. ユーザーが再検証
5. 新しいセッションクッキーを発行
```

## クッキーの詳細

### クッキー名
`turnstile_verified`

### クッキー属性
- **HttpOnly**: JavaScriptからアクセス不可（XSS対策）
- **Secure**: HTTPS通信のみ（中間者攻撃対策）
- **SameSite=Strict**: CSRF対策
- **Path**: `/`（全パスで有効）
- **MaxAge**: 設定された有効期間

### クッキー値の構造
```
Base64(IP|Timestamp|Counter|HMAC-SHA256-Signature)
```

**例:**
```
MTI3LjAuMC4xfDE3MzEwNTI4MDB8MHxhYmNkZWYxMjM0NTY3ODkwYWJjZGVm...
```

### セキュリティ機能

1. **HMAC署名**: サーバー秘密鍵による署名で改ざん防止
2. **IPバインディング**: クッキー作成時のIPとリクエストIPを照合
3. **タイムスタンプ**: 作成時刻を記録し有効期限を管理
4. **カウンター**: 使用回数を追跡し上限を管理

## 設定例

### Case 1: 一般的なWebサイト（推奨）
ユーザー体験重視、適度なセキュリティ

```yaml
TurnstileUseSession: true
TurnstileSessionDuration: 3600  # 1時間
TurnstileSessionMaxUses: 0      # 無制限
```

**効果:**
- ✅ 1時間以内なら何度アクセスしても再検証不要
- ✅ ユーザーストレス最小
- ⚠️ セキュリティは中程度

---

### Case 2: セキュリティ重視
重要なデータを扱う場合

```yaml
TurnstileUseSession: true
TurnstileSessionDuration: 900   # 15分
TurnstileSessionMaxUses: 5      # 5回まで
```

**効果:**
- ✅ 15分または5回でセッション期限切れ
- ✅ セキュリティレベル高
- ⚠️ ユーザーは頻繁に再検証が必要

---

### Case 3: 開発/テスト環境
利便性最優先

```yaml
TurnstileUseSession: true
TurnstileSessionDuration: 86400 # 24時間
TurnstileSessionMaxUses: 0      # 無制限
```

**効果:**
- ✅ 1日中再検証不要
- ✅ 開発効率最大
- ⚠️ 本番環境では非推奨

---

### Case 4: 最高セキュリティ
極めて重要なデータ、または攻撃を受けている場合

```yaml
TurnstileUseSession: false      # セッション機能無効
```

**効果:**
- ✅ 毎回Turnstile検証
- ✅ 最高のセキュリティ
- ❌ ユーザー体験が悪い

---

## 動作確認方法

### 1. 初回アクセス
```bash
curl -v "https://chouette2100.com/contributors?ieventid=XXX&roomid=YYY"
```

期待される動作:
- Turnstileチャレンジが表示される
- セッションクッキーは発行されない

### 2. チャレンジ完了後
ブラウザで検証を完了すると:
- データが表示される
- `Set-Cookie: turnstile_verified=...` ヘッダーが返される

### 3. 2回目のアクセス
同じブラウザから再度アクセス:
- Turnstileチャレンジがスキップされる
- 直接データが表示される

### 4. クッキー確認
ブラウザの開発者ツールで確認:
```
Application > Cookies > https://chouette2100.com
  turnstile_verified: (Base64エンコードされた値)
  HttpOnly: ✓
  Secure: ✓
  SameSite: Strict
```

### 5. ログ確認
サーバーログで動作を確認:
```bash
tail -f XXXXXX_YYYYMMDD.txt | grep Turnstile
```

期待されるログ:
```
Turnstile session created for IP 203.0.113.1
Turnstile session valid for IP 203.0.113.1
Turnstile session check failed: session expired
```

---

## トラブルシューティング

### 問題1: セッションが機能しない
**症状:** 毎回Turnstileチャレンジが表示される

**確認事項:**
1. `TurnstileUseSession: true` になっているか
2. ブラウザがクッキーを受け入れているか
3. HTTPSで接続しているか（HTTP接続ではSecure属性のクッキーは送信されない）

**解決策:**
```bash
# 設定確認
grep TurnstileUseSession ServerConfig.yml

# ブラウザのクッキー設定を確認
# クッキーをブロックしていないか確認
```

---

### 問題2: "IP mismatch" エラー
**症状:** セッションが突然無効になる

**原因:**
- モバイル回線でIPアドレスが変わった
- VPN接続を切り替えた
- プロキシ経由でアクセスしている

**解決策:**
仕様通りの動作です。セキュリティのためIPアドレスが変わったら再検証が必要です。

---

### 問題3: セッションが短すぎる/長すぎる
**症状:** ユーザーから期間の調整要望

**解決策:**
```yaml
# ServerConfig.ymlで調整
TurnstileSessionDuration: 7200  # 2時間に延長
TurnstileSessionMaxUses: 20     # 20回に増加
```

サーバーを再起動:
```bash
systemctl --user restart SRCGI
```

---

### 問題4: クッキーが削除される
**症状:** ブラウザを閉じるとセッションが消える

**原因:**
ブラウザの設定で「終了時にクッキーを削除」が有効

**解決策:**
ユーザーにブラウザ設定の変更を依頼するか、`TurnstileSessionDuration`を長めに設定。

---

## セキュリティに関する注意事項

### 適切な有効期間の設定
- **短すぎる**: ユーザー体験が悪化
- **長すぎる**: セキュリティリスク増加
- **推奨**: 1時間（3600秒）が多くの場合に適切

### IPアドレスバインディング
- モバイルユーザーはIP変更が頻繁
- 正規ユーザーにも再検証を求める可能性がある
- トレードオフを理解した上で運用

### クッキーの署名
- `TurnstileSecretKey`は絶対に公開しない
- 秘密鍵が漏洩すると署名が偽造される
- 定期的な鍵ローテーションも検討

### 使用回数制限
- `0`（無制限）は利便性が高いが、セッション再利用攻撃に弱い
- 適度な制限（5-20回）を推奨

---

## パフォーマンスへの影響

### 初回アクセス
- Cloudflare API呼び出し: 1回
- レイテンシ: +50-200ms（Cloudflareまでの通信時間）

### 2回目以降（セッション有効時）
- Cloudflare API呼び出し: **0回**
- レイテンシ: +1ms未満（HMAC検証のみ）

### セッション機能の効果
100アクセス中90アクセスがセッション有効の場合:
- Cloudflare APIコール削減: 90%
- 平均レイテンシ削減: 80-180ms × 90% = 大幅な改善

---

## まとめ

✅ **セッション機能を有効化することを強く推奨**

理由:
1. ユーザー体験の大幅な改善
2. Cloudflare APIへの負荷削減
3. レスポンス時間の短縮
4. 適切に実装されたセキュリティ

無効化するのは:
- 非常に高いセキュリティが要求される場合のみ
- 現在進行形で攻撃を受けている場合のみ

通常の運用では、推奨設定（1時間、無制限）で十分なセキュリティとユーザー体験のバランスが取れています。
