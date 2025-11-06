# アクセス集計表機能

## 概要
Webサーバーへのアクセスログを以下の3つの観点から集計表示する機能です：
1. ハンドラー別のアクセス数
2. IPアドレス別のアクセス数（上位50件、暗号化表示）
3. ユーザーエージェント別のアクセス数（上位50件）

## アクセスURL
`http://[server]/accesstable`

## 機能詳細

### 共通仕様
- 表示期間：終了日から遡って7日間（固定）
- デフォルト終了日：アクセス当日
- ボットアクセスを除外（`is_bot = 0`）
- 特定の管理用IPアドレスを除外
- 合計アクセス数の降順でソート

### 集計タイプ

#### 1. ハンドラー別（type=handler）
- 全ハンドラーのアクセス数を表示
- ハンドラー名で識別

#### 2. IPアドレス別（type=ip）
- 期間内アクセス数上位50件のIPアドレスを表示
- IPアドレスはAES-256-GCMで暗号化して表示
- セキュリティとプライバシー保護のため

#### 3. ユーザーエージェント別（type=useragent）
- 期間内アクセス数上位50件のユーザーエージェントを表示
- ブラウザやクライアントの分布を把握

## 使用方法

### 基本的な使い方
1. `/accesstable` にアクセス（デフォルト：ハンドラー別、直近7日間）
2. 画面上部のボタンで集計タイプを切り替え
3. 終了日を変更して「再表示」ボタンをクリック

### URLパラメータ
```
/accesstable?type=[handler|ip|useragent]&end_date=YYYY-MM-DD
```

例：
- `/accesstable` - ハンドラー別、今日まで
- `/accesstable?type=ip` - IPアドレス別、今日まで
- `/accesstable?type=useragent&end_date=2025-11-01` - UA別、2025-11-01まで

## IPアドレスの暗号化・復号化

### 暗号化アルゴリズム
**AES-256-GCM（Galois/Counter Mode）**

- 暗号化キー長：256ビット（32バイト）
- 認証付き暗号化（AEAD）
- ノンス（Nonce）：12バイト、ランダム生成
- エンコーディング：Base64 URL-safe

### 暗号化処理フロー
1. 平文IPアドレス（例：`192.168.1.1`）
2. AES-256ブロック暗号の初期化
3. GCMモードの初期化
4. ランダムノンス生成（12バイト）
5. GCM暗号化（ノンス + 暗号文 + 認証タグ）
6. Base64 URL-safeエンコード

### 復号化処理フロー
1. Base64 URL-safe文字列をデコード
2. ノンス部分（先頭12バイト）と暗号文部分を分離
3. GCM復号化（認証検証を含む）
4. 平文IPアドレスを取得

### セキュリティ上の注意
- 暗号化キーは `ShowroomCGIlib/IpEncryption.go` 内の `encryptionKey` 変数に定義
- **本番環境では必ず以下を実施してください：**
  1. 環境変数または設定ファイルから暗号化キーを読み込む
  2. キーをソースコード内にハードコーディングしない
  3. キーを定期的にローテーションする
  4. キーのバックアップを安全に保管する

### 復号化ツールの作成例

Goでの復号化ツール：

```go
package main

import (
	"fmt"
	"os"
	"github.com/Chouette2100/SRCGI/ShowroomCGIlib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: decrypt <encrypted_ip>")
		os.Exit(1)
	}
	
	encryptedIP := os.Args[1]
	decrypted, err := ShowroomCGIlib.DecryptIP(encryptedIP)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Decrypted IP: %s\n", decrypted)
}
```

## ファイル構成

### 新規作成ファイル
- `ShowroomCGIlib/IpEncryption.go` - IP暗号化/復号化機能
- `ShowroomCGIlib/dbAccessTable.go` - データベースアクセス関数
- `ShowroomCGIlib/HandlerAccessTable.go` - HTTPハンドラー
- `templates/accesstable.gtpl` - HTMLテンプレート

### 変更ファイル
- `ShowroomCGIlib/ShowroomCGIlib.go` - データ構造の追加
- `main.go` - ルーティングの追加

## データベース

### 使用テーブル
`accesslog` テーブル

```sql
CREATE TABLE accesslog (
  handler varchar(45) NOT NULL,
  remoteaddress char(15) NOT NULL,
  useragent varchar(256) NOT NULL,
  referer varchar(256) NOT NULL DEFAULT '',
  formvalues varchar(1024) NOT NULL,
  eventid varchar(100) NOT NULL,
  roomid int NOT NULL,
  ts datetime(3) NOT NULL,
  is_bot tinyint(1) DEFAULT 0,
  PRIMARY KEY (ts, eventid)
);
```

### 除外されるIPアドレス
- `59.166.119.117`
- `10.63.22.1`
- `149.88.103.40`

## 画面構成

### ナビゲーション
- トップ
- 開催中イベント一覧
- 開催予定イベント一覧
- 終了イベント一覧
- 日別アクセス統計
- 時刻別アクセス統計
- **アクセス集計表** ← 新規追加

### 集計タイプ切り替えボタン
- ハンドラー別（緑色ハイライト表示）
- IPアドレス別
- ユーザーエージェント別

### 日付フィルター
- 終了日入力フィールド
- 再表示ボタン

### テーブル表示
| 項目 | 日1 | 日2 | 日3 | 日4 | 日5 | 日6 | 日7 | 合計 |
|------|-----|-----|-----|-----|-----|-----|-----|------|
| 項目1 | 数値 | 数値 | ... | ... | ... | ... | 数値 | 数値 |
| 項目2 | 数値 | 数値 | ... | ... | ... | ... | 数値 | 数値 |

- 日付ヘッダー：YY-MM-DD形式（例：25-11-06）
- 合計列：黄色背景でハイライト
- ホバー時：行全体が強調表示

### 統計カード
- 件数：表示されている項目の数
- 総アクセス数：7日間の合計アクセス数

## 保守性の考慮

### 単一ハンドラー設計
3つの集計機能を1つのハンドラー（`AccessTableHandler`）に統合：
- コードの重複を最小化
- 共通ロジックの一元管理
- テンプレートの再利用

### 汎用データ取得関数
`buildTableRows` 関数でジェネリックな処理を実装：
- ハンドラー、IP、UAの各集計で共通処理を再利用
- 保守性と拡張性の向上

### テンプレート条件分岐
単一テンプレートで3つの表示モードを制御：
- 集計タイプに応じた表示切り替え
- 一貫したUI/UX

## 今後の拡張案
- 表示日数の可変化（7日固定→ユーザー指定）
- 開始日の直接指定
- CSV/Excelエクスポート機能
- グラフ表示の追加
- リファラー別集計の追加
