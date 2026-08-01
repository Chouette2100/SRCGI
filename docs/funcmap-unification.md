# FuncMap統一作業ドキュメント

作成日: 2026-08-01  
対象リポジトリ: SRCGI (`ShowroomCGIlib/`)

---

## 1. 今回の作業で完了したこと

### 1.1 新規追加した関数（ShowroomCGIlib/ShowroomCGIlib.go）

| 関数名 | 内容 |
|---|---|
| `CloneCommonFuncMap()` | `CommonFuncMap` の浅いコピーを返す。独自拡張がいらないハンドラーが使う |
| `MergeCommonFuncMap(extra)` | `CommonFuncMap` のコピーに `extra` の関数を上書き合成して返す。独自関数が必要なハンドラーが使う |

### 1.2 CommonFuncMap に昇格させた関数

以下の関数を `CommonFuncMap` に追加済み（`ShowroomCGIlib/ShowroomCGIlib.go` の `init()` に集約）。

| 関数名 | 旧所在 | 内容 |
|---|---|---|
| `FormatTime(t, layout?)` | 各ハンドラー（旧 `time.Time, string` の2引数固定） | 可変引数対応。1引数時はデフォルト `"2006-01-02 15:04"` |
| `TimeToString(t, layout?)` | 各ハンドラー（旧 `time.Time` の1引数固定） | 可変引数対応。1引数時はデフォルト `"01-02 15:04"` |
| `UnixTimeToYYYYMMDDHHMM` | `HandlerCurrentDistributors.go` | UnixTime → `"2006-01-02 15:04"` |
| `UnixTimeToHHMM` | `HandlerCurrentDistributors.go` | UnixTime → `"15:04"` |
| `t2s` | 各ハンドラー | `time.Time, string` → フォーマット済み文字列（`FormatTime` と同じ挙動） |
| `FormatInt` | `HandlerMonthlyCntrbRankOfListener.go` | `fmt.Sprintf(layout, int)` |
| `Showrank` | `HandlerShowRank.go`, `HandlerTmShowRank.go` | `" \| "` 区切り文字列の末尾要素を返す |
| `iscurrent` | `HandlerRoomCntrbHistory.go`, `HandlerListenerCntrbHistory.go` | `time.Time.After(time.Now())` |
| `div` | `HandlerListLastP.go` | ゼロ除算安全な整数除算 |

### 1.3 各ハンドラーの変更方針（適用済み）

| 旧パターン | 新パターン | 対象ファイル数 |
|---|---|---|
| `sprig.FuncMap()` を直接使用 | `CloneCommonFuncMap()` | 6 |
| `funcMap := CommonFuncMap`（危険：直接変更） | `CloneCommonFuncMap()` | 3 |
| `template.FuncMap{...}` のみ（Commonなし） | `MergeCommonFuncMap(template.FuncMap{...})` | 21 |
| `CommonFuncMap` を直接 `.Funcs()` に渡す | `CloneCommonFuncMap()` | 3 |

---

## 2. 次フェーズの作業項目

### 2.1 【必須】ハンドラーのローカル重複関数の削除

`MergeCommonFuncMap(template.FuncMap{ ... })` の中に、`CommonFuncMap` と**同名・同定義**の関数が残っている。
これらを削除し、`CloneCommonFuncMap()` に切り替える。

#### 削除候補一覧

以下の関数は `CommonFuncMap` に同定義で存在するため、ローカル宣言を削除できる。

| ファイル | 削除できるローカル関数 | CommonFuncMap 側の関数名 |
|---|---|---|
| `HandlerBbs.go` | `htmlEscapeString`, `FormatTime`, `CntToName`, `Add` | 全て同名で CommonFuncMap 済み |
| `HandlerClosedEventRoomList.go` | `Comma`, `UnixtimeToTime` | `Comma`, `UnixtimeToTime` |
| `HandlerCurrentDistributors.go` | `Comma`, `UnixTimeToYYYYMMDDHHMM`, `UnixTimeToHHMM` | 全て昇格済み ※`GidToName` は独自のため残す |
| `HandlerListCntrb.go` | `sub`, `Comma` | `sub`, `Comma` |
| `HandlerListCntrbEx.go` | `sub`, `Comma` | `sub`, `Comma` |
| `HandlerListGiftScore.go` | `sub`, `Comma`, `t2s` | `sub`, `Comma`, `t2s` ※`add` は未昇格 → 要検討 |
| `HandlerListGiftScoreCntrb.go` | `sub`, `Comma`, `t2s` | 全て同名で CommonFuncMap 済み |
| `HandlerListFanGiftScore.go` | `sub`, `Comma`, `t2s` | 全て同名で CommonFuncMap 済み |
| `HandlerListLast.go` | `Comma`, `DelBlockID`, `Add` | 全て同名で CommonFuncMap 済み |
| `HandlerListLastP.go` | `Comma`, `DelBlockID`, `Add`, `div` | 全て昇格済み |
| `HandlerListToDo.go` | `Comma`, `FormatTime` | 同名昇格済み ※`FormatTimePtr` は未昇格 → 要検討 |
| `HandlerMonthlyCntrbRankOfListener.go` (×2) | `sub`, `Comma`, `FormatTime`, `FormatInt` | 全て同名で CommonFuncMap 済み |
| `HandlerOldEvents.go` | `TimeToString`, `TimeToStringY`, `IsTempID` | 全て同名で CommonFuncMap 済み |
| `HandlerScheduledEventsSvr.go` | `Comma`, `UnixTimeToStr` | `Comma`, `UnixTimeToStr` |
| `HandlerShowRank.go` | `Comma`, `FormatTime`, `Add`, `Showrank` | 全て昇格済み |
| `HandlerTmShowRank.go` | `Comma`, `FormatTime`, `Add`, `Showrank` | 全て昇格済み |
| `HandlerTopRoom.go` | `Comma`, `FormatTime` | 全て同名で CommonFuncMap 済み |
| `HandlerEditToDo.go` | `Comma`, `FormatTime` | 同名昇格済み ※`FormatTimePtr` は未昇格 → 要検討 |

#### 削除後の形（例: HandlerTopRoom.go）

```go
// 変更前
funcMap := MergeCommonFuncMap(template.FuncMap{
    "Comma":      func(i int) string { return humanize.Comma(int64(i)) },
    "FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
})

// 変更後
funcMap := CloneCommonFuncMap()
```

---

### 2.2 【要検討】CommonFuncMap への昇格候補

次の関数はいずれかのハンドラーにのみ存在し、他ハンドラーや今後のテンプレートで必要になる可能性がある。

| 関数名 | 現在の所在 | 内容 | 昇格推奨度 |
|---|---|---|---|
| `add` | `HandlerListGiftScore.go` | `i + j` | ★★★ sprig の `add` と同じなので sprig 経由で不要かも |
| `FormatTimePtr` | `HandlerListToDo.go`, `HandlerEditToDo.go` | `*time.Time` 対応フォーマット（nil時は`"-"`） | ★★★ （`FormatTime` に統合できる。現状の `FormatTime` は `*time.Time` 対応済み） |
| `GidToName` | `HandlerCurrentDistributors.go` | ジャンルID→ジャンル名（クロージャ内変数参照） | ★☆☆ 独自データ依存のためCommonFuncMap昇格は不可 |
| `UnixTimeToYYYYMMDDHHMM` | `CommonFuncMap` 済み | `time.Unix(...).Format("2006-01-02 15:04")` | 済み |

**`FormatTimePtr` について**: 現在の `CommonFuncMap["FormatTime"]` は `*time.Time` 対応済み（nilなら`"-"`を返す）。`HandlerListToDo.go` の `FormatTimePtr` と同じ挙動なので、テンプレート側の `FormatTimePtr` 呼び出しを `FormatTime` に統一すれば削除できる。

---

### 2.3 【要対応】Turnstile 経由の GetFuncMap が CommonFuncMap 非取込みのもの

`TurnstileHandler.go` は `challengeData.GetFuncMap()` が返すマップを使ってテンプレートをパースする。
現在の実装では `GetFuncMap()` が `nil` を返すハンドラーのテンプレートには CommonFuncMap が入らず、**告知表示が動かない**。

| ハンドラー | `GetFuncMap()` の戻り値 | 状態 |
|---|---|---|
| `HandlerContributors.go` (`HCntrbInf`) | `nil` | ⚠️ 告知未対応 |
| `HandlerGraphSum.go` (`GraphSumInf`) | `nil` | ⚠️ 告知未対応 |
| `HandlerGraphSum2.go` (`GraphSum2Inf`) | `nil` | ⚠️ 告知未対応 |
| `HandlerListLastC.go` (`ListLastC`) | `nil` | ⚠️ 告知未対応 |
| `HandlerEventTop.go` (`EventTopInf`) | `nil` | ⚠️ 告知未対応 |
| `HandlerCurrentEvents.go` (`T999Dtop`) | `&CommonFuncMap`（直接参照） | ⚠️ コピーにすべき |
| `HandlerClosedEvents.go` (`ClosedEventsInf`) | `ClosedEventsfuncMap`（`= &CommonFuncMap` のポインタ） | ⚠️ コピーにすべき |
| `HandlerListCntrbH.go` (`CntrbH_Header`) | `ListCntrbHfuncMap`（独自、CommonFuncMap非取込み） | ⚠️ 告知未対応 |
| `HandlerListCntrbHEx.go` (`CntrbHEx_Header`) | `ListCntrbHExfuncMap`（独自、CommonFuncMap非取込み） | ⚠️ 告知未対応 |

#### 対応方針

`GetFuncMap()` を返す全実装を `CloneCommonFuncMap()` または `MergeCommonFuncMap(...)` ベースに変更する。  
ポインタを返す設計なので、関数内でローカルコピーを作ってそのポインタを返す形にする。

```go
// 変更例（HandlerListCntrbH.go）

// 変更前
var ListCntrbHfuncMap = &template.FuncMap{
    "sub":        func(i, j int) int { return i - j },
    "Comma":      func(i int) string { return humanize.Comma(int64(i)) },
    "FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
}

func (h *CntrbH_Header) GetFuncMap() *template.FuncMap {
    return ListCntrbHfuncMap
}

// 変更後（毎回コピーを生成する形）
func (h *CntrbH_Header) GetFuncMap() *template.FuncMap {
    fm := CloneCommonFuncMap()
    // sub, Comma, FormatTime は CommonFuncMap に既に含まれるので不要
    return &fm
}
```

`nil` を返していたハンドラーは次のように変更する：

```go
// 変更前
func (h *HCntrbInf) GetFuncMap() *template.FuncMap {
    return nil
}

// 変更後
func (h *HCntrbInf) GetFuncMap() *template.FuncMap {
    fm := CloneCommonFuncMap()
    return &fm
}
```

あわせて `TurnstileHandler.go` の `nil` 分岐（127行目）を消し、常にFuncMapを使う形に統一することも検討できる。

---

### 2.4 【要確認】コメントアウト状態のローカル FuncMap 残骸の整理

以下のファイルにコメントアウトされたローカル `template.FuncMap{}` ブロックが残っている。  
今後の混乱を避けるために削除することを推奨する。

| ファイル | 行番号 |
|---|---|
| `HandlerClosedEvents.go` | 54〜83行（`/* ... */` ブロック） |
| `HandlerListCntrbH.go` | 122〜125行 |
| `HandlerListCntrbHEx.go` | 125〜130行 |
| `HandlerEditCntrbPoints.go` | 142〜149行 |
| `HandlerScheduledEvents.go` | 83〜89行（`/* ... */` ブロック） |
| `HandlerTop.go` | 122〜131行（`/* ... */` ブロック） |
| `HandlerCurrentEvents.go` | 121〜138行（`/* ... */` ブロック） |

---

### 2.5 【テンプレート側】告知表示がないテンプレート

次のテンプレートには `HasAnnouncement` / `GetAnnouncement` の呼び出しがなく、告知が表示されない。  
ハンドラー側の FuncMap 整備が終わったら、テンプレートにも告知ブロックを追加する。

#### 告知表示のある既存テンプレートの参考実装（再利用可）

```html
{{if (HasAnnouncement)}}
<div style="padding: 15px; margin: 0 0 20px 0; border-radius: 6px;
     background-color: {{(GetAnnouncement).BgColor}};
     color: {{(GetAnnouncement).TextColor}};
     font-size: 16px; font-weight: bold; text-align: center;
     border: 2px solid {{(GetAnnouncement).TextColor}}; word-wrap: break-word;">
  {{(GetAnnouncement).Message}}
</div>
{{end}}
```

#### 未対応テンプレート（告知ブロック追加が必要なもの）

```
templates/add-event2.gtpl
templates/bbs-2.gtpl          ← bbs-1_org.gtpl / bbs-1_maint.gtpl と同時使用
templates/edit-user2.gtpl
templates/fanlevel-room.gtpl
templates/list-cntrb-h2.gtpl
templates/list-cntrb.gtpl
templates/list-cntrbex-h2.gtpl
templates/list-gs-h2.gtpl
templates/list-gs.gtpl
templates/list-gs-h2.gtpl
templates/list-vgs-h2.gtpl
templates/list-vgs.gtpl
templates/new-event1.gtpl
templates/new-event2.gtpl
templates/param-event0.gtpl   ← 管理系・不要の可能性あり
templates/param-event1.gtpl   ← 管理系・不要の可能性あり
templates/param-event2.gtpl   ← 管理系・不要の可能性あり
templates/param-eventc.gtpl   ← 管理系・不要の可能性あり
templates/param-global.gtpl   ← 管理系・不要の可能性あり
templates/param-local.gtpl    ← 管理系・不要の可能性あり
templates/turnstilechallenge.gtpl ← CAPTCHA中間ページのため不要
```

`bbs-3.gtpl`、`footer.gtpl`、`list-cntrb-h2.gtpl` 等の「ヘッダ・フッタ部品テンプレート」は親テンプレート側に告知を入れれば追加不要な場合もある。各テンプレートの構成を確認してから追加すること。

---

## 3. 作業の優先順位

```
1. テスト実施 → 動作確認
   （今回の変更で既存ページが壊れていないか確認）

2. 2.3 GetFuncMap() が nil のハンドラーの対応
   （Turnstile経由ページで告知が表示されないバグの修正）

3. 2.1 ローカル重複関数の削除
   （コード整理。安全にできるものから順次）

4. 2.2 CommonFuncMap昇格候補の整理
   （FormatTimePtr → FormatTime への統合を含む）

5. 2.5 テンプレートへの告知ブロック追加
   （全ページへの告知表示の完成）

6. 2.4 コメントアウト残骸の削除
   （最後に行うコードクリーンアップ）
```

---

## 4. CommonFuncMap の現在の全関数一覧（参照用）

`ShowroomCGIlib/ShowroomCGIlib.go` の `init()` より（2026-08-01時点）。

| 関数名 | シグネチャ | 内容 |
|---|---|---|
| sprig 全関数 | — | `sprig.FuncMap()` より取込み済み |
| `Add` | `(int, int) int` | 加算 |
| `baseOfEventid` | `(string) string` | `?` 前の文字列を返す |
| `CntToName` | `(int) string` | 問い合わせ種別番号→日本語名 |
| `Comma` | `(int) string` | 3桁区切りカンマ |
| `DelBlockID` | `(string) string` | ブロックID除去（`?` 以降を削除） |
| `div` | `(int, int) int` | ゼロ除算安全な整数除算 |
| `Divide` | `(int, int) int` | ゼロ除算安全な整数除算（`div` と同じ） |
| `FormatInt` | `(int, string) string` | `fmt.Sprintf` ラッパー |
| `FormatTime` | `(any, string?) string` | `time.Time` / `*time.Time` → フォーマット文字列。引数1個でデフォルト `"2006-01-02 15:04"` |
| `GetAnnouncement` | `() AnnouncementData` | 現在の告知データを返す |
| `HasAnnouncement` | `() bool` | 告知が有効かつ文字列が空でない |
| `htmlEscapeString` | `(string) string` | HTML エスケープ |
| `iscurrent` | `(time.Time) bool` | `t.After(time.Now())` |
| `IsTempID` | `(string) bool` | `"@@@@"` プレフィックス判定 |
| `Mod` | `(int, int) int` | ゼロ除算安全な剰余 |
| `SafeColor` | `(string) string` | `#RRGGBB` バリデーション付き色文字列返却 |
| `Showrank` | `(string) string` | `" \| "` 区切りの末尾部分を返す |
| `sub` | `(int, int) int` | 減算 |
| `t2s` | `(time.Time, string) string` | `time.Time` → フォーマット文字列 |
| `TimeToString` | `(any, string?) string` | `time.Time` / `*time.Time` → フォーマット文字列。引数1個でデフォルト `"01-02 15:04"` |
| `TimeToStringY` | `(time.Time) string` | `"06-01-02 15:04"` 固定フォーマット |
| `UnixtimeToTime` | `(int64, string) string` | UnixTime → フォーマット文字列 |
| `UnixTimeToHHMM` | `(int64) string` | UnixTime → `"15:04"` |
| `UnixTimeToStr` | `(int64) string` | UnixTime → `"01-02 15:04"` |
| `UnixTimeToStrY` | `(int64) string` | UnixTime → `"06-01-02 15:04"` |
| `UnixTimeToYYYYMMDDHHMM` | `(int64) string` | UnixTime → `"2006-01-02 15:04"` |

---

## 5. 注意点・落とし穴

### `FormatTime` の引数変更について

旧来の `"FormatTime": func(t time.Time, tfmt string) string` と新しい可変引数版は**テンプレート側の呼び出しは変わらない**（`{{ FormatTime .Ts "2006-01-02 15:04" }}`）。ただし型が `interface{}` になるため、Goコード上で直接呼び出しているコードがあれば注意。

### `TimeToString` の引数変更について

旧来の一部ハンドラーは `"TimeToString": func(t time.Time) string { return t.Format("01-02 15:04") }` と1引数で定義していた。現在の `CommonFuncMap["TimeToString"]` は可変引数対応で、1引数でも2引数でも呼べる。**テンプレート側の `{{ TimeToString .Ts }}` は変更不要**。

### `CommonFuncMap` の直接変更は禁止

今後いかなる場合も `CommonFuncMap[...] = ...` をランタイム（`init()` 以外）に行ってはならない。`CloneCommonFuncMap()` または `MergeCommonFuncMap(...)` を使うこと。  
`funcMap := CommonFuncMap` は**コピーでなく参照**なので必ず禁止。
