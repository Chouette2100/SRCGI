# Turnstile汎用関数の使用ガイド

## 概要

Turnstile検証とセッション管理を簡単に他のハンドラーに適用できる汎用関数を提供します。わずか数行のコードで、どのハンドラーにもTurnstile保護を追加できます。

---

## 汎用関数の構成

### 新規ファイル
**`ShowroomCGIlib/TurnstileHandler.go`**

提供される関数:
1. `CheckTurnstileWithSession()` - セッション管理付きTurnstile検証（推奨）
2. `SimpleTurnstileCheck()` - セッションなしの単純検証

### インターフェース
```go
type TurnstileChallengeData interface {
    SetTurnstileInfo(siteKey string, errorMsg string)
    GetTemplatePath() string
    GetTemplateName() string
}
```

---

## 使用方法

### Step 1: データ構造体を準備

ハンドラーで使用するデータ構造体に、以下のフィールドを追加:

```go
type YourHandlerData struct {
    // ハンドラー固有のフィールド
    EventID    string
    RoomID     int
    // ... その他のフィールド ...
    
    // Turnstile用のフィールド（必須）
    TurnstileSiteKey string
    TurnstileError   string
}
```

### Step 2: インターフェースを実装

構造体に3つのメソッドを追加:

```go
// TurnstileChallengeDataインターフェースの実装
func (d *YourHandlerData) SetTurnstileInfo(siteKey string, errorMsg string) {
    d.TurnstileSiteKey = siteKey
    d.TurnstileError = errorMsg
}

func (d *YourHandlerData) GetTemplatePath() string {
    return "templates/your_template.gtpl"
}

func (d *YourHandlerData) GetTemplateName() string {
    return "your_template.gtpl"
}
```

### Step 3: ハンドラーで使用

```go
func YourHandler(w http.ResponseWriter, r *http.Request) {
    // 1. データ構造体を準備
    data := YourHandlerData{
        EventID: r.FormValue("eventid"),
        RoomID:  roomid,
        // ... その他の初期化 ...
    }
    
    // 2. Turnstile検証を実行（この1行だけ！）
    result, err := CheckTurnstileWithSession(w, r, &data)
    if result != TurnstileOK {
        // チャレンジページまたはエラーページが表示済み
        return
    }
    
    // 3. 検証OK、通常処理を続行
    // ... データ取得、処理など ...
    
    // 4. テンプレート表示（TurnstileSiteKeyをセット）
    data.TurnstileSiteKey = Serverconfig.TurnstileSiteKey
    // ... テンプレート実行 ...
}
```

---

## 実際の適用例

### 例1: ContributorsHandler（実装済み）

**Before（70行以上）:**
```go
if Serverconfig.TurnstileSiteKey != "" {
    sessionValid, newCookie, sessionErr := VerifyTurnstileSessionCookie(r)
    if sessionValid {
        // ... セッション処理 ...
    } else {
        turnstileToken := r.FormValue("cf-turnstile-response")
        if turnstileToken == "" {
            // ... チャレンジ表示 ...
            return
        }
        // ... 検証処理 ...
    }
}
```

**After（わずか7行）:**
```go
hcntrbinf := HCntrbInf{
    Ieventid:   ieventid,
    Eventid:    event.Eventid,
    Event_name: event.Event_name,
    Period:     event.Period,
    Roomid:     roomid,
}

result, tsErr := CheckTurnstileWithSession(w, r, &hcntrbinf)
if result != TurnstileOK {
    return
}
```

### 例2: GraphTotalHandler への適用

**ステップ1: データ構造体を定義**
```go
type GraphTotalData struct {
    Eventid          string
    Event_name       string
    // ... 他のフィールド ...
    TurnstileSiteKey string
    TurnstileError   string
}

// インターフェース実装
func (d *GraphTotalData) SetTurnstileInfo(siteKey string, errorMsg string) {
    d.TurnstileSiteKey = siteKey
    d.TurnstileError = errorMsg
}

func (d *GraphTotalData) GetTemplatePath() string {
    return "templates/graph-total.gtpl"
}

func (d *GraphTotalData) GetTemplateName() string {
    return "graph-total.gtpl"
}
```

**ステップ2: ハンドラーに適用**
```go
func GraphTotalHandler(w http.ResponseWriter, r *http.Request) {
    eventid := r.FormValue("eventid")
    
    // データ準備
    data := GraphTotalData{
        Eventid: eventid,
    }
    
    // Turnstile検証（この1行！）
    result, _ := CheckTurnstileWithSession(w, r, &data)
    if result != TurnstileOK {
        return
    }
    
    // 通常処理
    // ... グラフデータ取得など ...
    
    // テンプレート表示
    data.TurnstileSiteKey = Serverconfig.TurnstileSiteKey
    // ... テンプレート実行 ...
}
```

**ステップ3: テンプレート更新**
```html
<!-- templates/graph-total.gtpl -->
<!DOCTYPE html>
<html>
<head>
{{if .TurnstileSiteKey}}
<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
{{end}}
</head>
<body>
{{if and .TurnstileSiteKey (not .GraphData)}}
<!-- Turnstileチャレンジ -->
<div>
<h3>セキュリティチェック</h3>
{{if .TurnstileError}}
<p style="color: red;">{{.TurnstileError}}</p>
{{end}}
<form method="POST" action="graph-total">
<input type="hidden" name="eventid" value="{{.Eventid}}">
<div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}"></div>
<button type="submit">確認して続行</button>
</form>
</div>
{{else}}
<!-- 通常のコンテンツ -->
<!-- ... グラフ表示など ... -->
{{end}}
</body>
</html>
```

### 例3: ListLastHandler への適用

```go
type ListLastData struct {
    Eventid          string
    EventUsers       []EventUser
    TurnstileSiteKey string
    TurnstileError   string
}

func (d *ListLastData) SetTurnstileInfo(siteKey string, errorMsg string) {
    d.TurnstileSiteKey = siteKey
    d.TurnstileError = errorMsg
}

func (d *ListLastData) GetTemplatePath() string {
    return "templates/list-last.gtpl"
}

func (d *ListLastData) GetTemplateName() string {
    return "list-last.gtpl"
}

func ListLastHandler(w http.ResponseWriter, r *http.Request) {
    data := ListLastData{
        Eventid: r.FormValue("eventid"),
    }
    
    // Turnstile検証（1行で完了）
    result, _ := CheckTurnstileWithSession(w, r, &data)
    if result != TurnstileOK {
        return
    }
    
    // データ取得
    // ... DB処理など ...
    
    // テンプレート表示
    data.TurnstileSiteKey = Serverconfig.TurnstileSiteKey
    // ... テンプレート実行 ...
}
```

---

## セッションなしの検証

高セキュリティが必要な場合や、セッション管理が不要な場合:

```go
// セッション管理なしの検証
result, err := SimpleTurnstileCheck(w, r, &data)
if result != TurnstileOK {
    return
}
```

この場合、毎回Turnstile検証が実行されます。

---

## テンプレートのパターン

### 基本パターン

```html
<!DOCTYPE html>
<html>
<head>
{{if .TurnstileSiteKey}}
<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>
{{end}}
</head>
<body>

{{if and .TurnstileSiteKey (not .DataField)}}
<!-- Turnstileチャレンジ表示 -->
<div class="turnstile-container">
    <h3>セキュリティチェック</h3>
    {{if .TurnstileError}}
    <p class="error">{{.TurnstileError}}</p>
    {{end}}
    <p>このページを表示するには、セキュリティチェックを完了してください。</p>
    <form method="POST" action="current-path">
        <!-- 必要なパラメータを隠しフィールドで送信 -->
        <input type="hidden" name="param1" value="{{.Param1}}">
        <input type="hidden" name="param2" value="{{.Param2}}">
        
        <!-- Turnstileウィジェット -->
        <div class="cf-turnstile" data-sitekey="{{.TurnstileSiteKey}}" data-theme="light"></div>
        <br>
        <button type="submit">確認して続行</button>
    </form>
</div>
{{else}}
<!-- 通常のコンテンツ表示 -->
<div class="content">
    <!-- ... 通常の表示内容 ... -->
</div>
{{end}}

</body>
</html>
```

### 条件分岐のポイント

チャレンジページと通常ページの切り替えは、「データが存在するか」で判定:

```html
{{if and .TurnstileSiteKey (not .DataField)}}
    <!-- チャレンジページ -->
{{else}}
    <!-- 通常ページ -->
{{end}}
```

- `.TurnstileSiteKey`: Turnstileが有効
- `(not .DataField)`: データがまだない（初回アクセス）

---

## チェックリスト

新しいハンドラーにTurnstileを適用する際のチェックリスト:

### コード側
- [ ] データ構造体に `TurnstileSiteKey` と `TurnstileError` フィールドを追加
- [ ] `SetTurnstileInfo()` メソッドを実装
- [ ] `GetTemplatePath()` メソッドを実装
- [ ] `GetTemplateName()` メソッドを実装
- [ ] ハンドラーで `CheckTurnstileWithSession()` を呼び出し
- [ ] 戻り値が `TurnstileOK` でない場合は return
- [ ] テンプレート表示前に `TurnstileSiteKey` をセット

### テンプレート側
- [ ] `<head>` にTurnstileスクリプトを追加
- [ ] チャレンジページの条件分岐を追加
- [ ] フォームに必要なパラメータの隠しフィールドを追加
- [ ] Turnstileウィジェット（`cf-turnstile` div）を配置
- [ ] エラーメッセージ表示領域を追加

---

## メリット

### 1. コードの簡潔化
- **70行以上 → 7行**に削減
- 複雑な条件分岐が不要
- 可読性が大幅に向上

### 2. 保守性の向上
- 1箇所の修正で全ハンドラーに反映
- バグ修正が容易
- 動作が統一される

### 3. 適用の容易さ
- 3つのメソッドを実装するだけ
- 1行の関数呼び出しで完了
- テンプレートも共通パターン

### 4. 柔軟性
- セッション管理の有無を選択可能
- ハンドラーごとに動作をカスタマイズ可能
- 既存コードへの影響なし

---

## トラブルシューティング

### 問題1: "does not implement TurnstileChallengeData" エラー

**原因:** 3つのメソッドのいずれかが未実装または署名が違う

**解決策:**
```go
func (d *YourData) SetTurnstileInfo(siteKey string, errorMsg string) {
    d.TurnstileSiteKey = siteKey
    d.TurnstileError = errorMsg
}

func (d *YourData) GetTemplatePath() string {
    return "templates/your.gtpl"
}

func (d *YourData) GetTemplateName() string {
    return "your.gtpl"
}
```

### 問題2: チャレンジページが表示されない

**原因:** テンプレートの条件分岐が正しくない

**解決策:**
データの有無を判定する適切なフィールドを使用:
```html
{{if and .TurnstileSiteKey (not .ActualDataField)}}
```

### 問題3: パラメータが引き継がれない

**原因:** フォームの隠しフィールドが不足

**解決策:**
必要なすべてのパラメータを追加:
```html
<input type="hidden" name="eventid" value="{{.Eventid}}">
<input type="hidden" name="roomid" value="{{.Roomid}}">
```

---

## ベストプラクティス

### 1. データ構造体の設計
```go
type HandlerData struct {
    // 必須パラメータ（検証前に必要）
    EventID string
    RoomID  int
    
    // データ（検証後に取得）
    Results []Result
    
    // Turnstile関連（必須）
    TurnstileSiteKey string
    TurnstileError   string
}
```

### 2. エラーハンドリング
```go
result, err := CheckTurnstileWithSession(w, r, &data)
if result != TurnstileOK {
    if err != nil {
        log.Printf("Turnstile error: %v", err)
    }
    return
}
```

### 3. テンプレートの構造化
- チャレンジページと通常ページを明確に分離
- エラーメッセージ表示を統一
- CSSでスタイリングを統一

---

## まとめ

汎用関数により:

✅ **3つのメソッド実装 + 1行の呼び出し**で完了
✅ 70行以上のコードが7行に
✅ 全ハンドラーで動作が統一される
✅ 保守性が劇的に向上
✅ 新規ハンドラーへの適用が簡単

次のハンドラーへの適用も、このパターンに従えば数分で完了します！
