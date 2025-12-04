# イベント参加ルームのリスナーの貢献ポイント履歴

## 概要
このハンドラーは、指定されたイベントに参加したリスナーの過去の貢献ポイント履歴を表示します。

## URL
```
/listener-cntrb-history
```

## パラメータ

| パラメータ名 | 型 | 必須 | デフォルト値 | 説明 |
|------------|---|-----|------------|------|
| eventid | string | ✓ | - | イベントID |
| nmonths | int | | 3 | 過去何ヶ月間のデータを表示するか |
| minpoint | int | | 5000 | 表示する最小貢献ポイント |
| maxnolines | int | | 200 | 最大表示行数 |
| ext | int | | 1 | 0: 全ルームの貢献を表示, 1: イベント参加ルームの貢献のみ表示 |

## 使用例

### 基本的な呼び出し
```
/listener-cntrb-history?eventid=EVENT_ID
```

### パラメータを指定した呼び出し
```
/listener-cntrb-history?eventid=EVENT_ID&nmonths=6&minpoint=10000&maxnolines=100&ext=0
```

## 実装ファイル

- ハンドラー: `ShowroomCGIlib/HandlerListenerCntrbHistory.go`
- テンプレート: `templates/listener-cntrb-history.gtpl`
- ルート登録: `main.go` (line 761)

## 使用しているテーブル

- `eventrank`: 貢献ポイントデータ
- `event`: イベント情報
- `user`: ルーム情報
- `viewer`: リスナー情報
- `eventuser`: イベント参加ルーム情報

## 機能

1. イベントに参加したリスナーの過去の貢献履歴を表示
2. 期間（過去何ヶ月間）を指定可能
3. 表示する最小ポイントと最大表示数を設定可能
4. イベント参加ルームのみ/全ルームの切り替え可能

## 注意事項

- データは貢献ポイントの降順で表示されます
- 表示されるのはイベント終了後のデータのみです（`er.ts > e.endtime`）
- リスナー名、ルーム名、イベント名は最大20文字まで表示されます
