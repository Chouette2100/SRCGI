# Handler中心の構造整理

作成日: 2026-09-22
対象: SRCGI

---

## 1. 結論

このリポジトリの中心は「URL と ハンドラーの対応」であり、各ハンドラーが最終的に HTML や JSON を出力する構造です。

重要なのは、.gtpl は「機能の本体」ではなく、handler の出力先であるという視点です。

- ルーティング: URL を handler に割り当てる
- handler: DB 参照、集計、データ組み立て
- template (.gtpl): そのデータを描画する
- 最終結果: HTTP 応答

つまり、実際の流れは次のように理解できます。

```text
main()
└─ ルーティング登録
   └─ handler
      ├─ データ取得
      ├─ 集計/生成
      └─ ExecuteTemplate(... .gtpl ...)
         └─ HTML/JSON を返す
```

---

## 2. 親子関係をハンドラー中心で見ると

```text
main
├─ ルーティング
│  ├─ /top → TopHandler
│  ├─ /eventtop → EventTopHandler
│  ├─ /list-level → ListLevelHandler
│  ├─ /list-last → ListLastHandler
│  ├─ /graph-total → GraphTotalHandler
│  ├─ /graph-sum → GraphSumHandler
│  └─ ...
│
└─ ハンドラ群
   ├─ 画面表示系
   │  ├─ TopHandler
   │  ├─ EventTopHandler
   │  ├─ ListLevelHandler
   │  ├─ ListLastHandler
   │  └─ ...
   │
   ├─ グラフ系
   │  ├─ GraphTotalHandler
   │  ├─ GraphPerdayHandler
   │  ├─ GraphPerslotHandler
   │  ├─ GraphSumHandler
   │  └─ ...
   │
   ├─ 編集系
   │  ├─ AddEventHandler
   │  ├─ EditUserHandler
   │  ├─ EditCntrbPointsHandler
   │  └─ ...
   │
   └─ データ補助系
      ├─ GraphSumDataHandler
      ├─ GraphSumData1Handler
      ├─ GraphSumData2Handler
      ├─ AccessStatsHandler
      └─ ...
```

- ここで重要なのは「関数の詳細」ではなく、URL と handler の親子関係です。
- .gtpl はその下位の表示部品として捉えるのが自然です。

---

## 3. 入口は main.go

ルーティングの実体は [main.go](../main.go#L746-L830) にあります。

例えば、ここで各 URL が handler に結びついています。

- [main.go](../main.go#L793-L818)
  - `/top` → `TopHandler`
  - `/eventtop` → `EventTopHandler`
  - `/graph-sum` → `GraphSumHandler`
  - `/list-last` → `ListLastHandler`
  - `/graph-total` → `GraphTotalHandler`

この部分は「どの URL がどのページを返すか」を定義する入口です。

---

## 4. ハンドラーの共通制御

すべてのハンドラーの前では、[main.go](../main.go#L288-L430) の `commonMiddleware` が処理を受けます。

ここが横断的な親役です。

- IP と地域の判定
- Rate limit
- Bot 判定
- fail2ban ログ
- accesslog の非同期保存

この middleware は「ページ本体」ではなく、handler を守る前段として動きます。

```text
URL
└─ commonMiddleware
   ├─ IP / region check
   ├─ rate limit
   ├─ bot check
   └─ handlerへ進む
```

---

## 5. ハンドラーの役割

ハンドラー本体は [ShowroomCGIlib/ShowroomCGIlib.go](../ShowroomCGIlib/ShowroomCGIlib.go#L360-L620) や各ファイルに分かれています。

典型的なハンドラーの役割は三つです。

1. データを集める
2. 画面用の構造体を作る
3. .gtpl を実行して HTML を返す

例:

- [ShowroomCGIlib/HandlerListLevel.go](../ShowroomCGIlib/HandlerListLevel.go#L52-L61)
  - `ListLevelHandler`
  - `templates/list-level.gtpl` を読み込む
  - `ExecuteTemplate(w, "list-level.gtpl", ...)` を実行する

- [ShowroomCGIlib/HandlerAccessStats.go](../ShowroomCGIlib/HandlerAccessStats.go#L101-L104)
  - `AccessStatsHandler`
  - `templates/accessstats.gtpl` を読み込む
  - `ExecuteTemplate(w, "accessstats.gtpl", ...)`

- [ShowroomCGIlib/HandlerListGiftScore.go](../ShowroomCGIlib/HandlerListGiftScore.go#L125-L150)
  - `list-gs-h1.gtpl`, `list-gs-h2.gtpl`, `list-gs.gtpl` をまとめて使う

ここが「gtpl から呼ばれている」という見立てと少し違っていて、実際は handler が template を呼んでいます。

---

## 6. gtpl との関係

よく混乱しやすいポイントです。

### 正しい認識

```text
handler
└─ template (.gtpl)
   └─ 画面に表示
```

### 逆に見えがちだが実際ではないもの

```text
gtpl
└─ handler を呼ぶ
```

これは実際にはほぼありません。

.gtplt はページのテンプレートであり、handler がテンプレートにデータを渡す役割を持ちます。

例え話すると、

- handler = 料理人
- .gtpl = 皿や盛り付け
- データ = 食材

のような関係です。

---

## 7. どこが親で、どこが子か

ハンドラーだけに集中すると、親子関係はこう整理できます。

```text
main
└─ ルーティング
   └─ handler
      ├─ 画面系
      │  └─ .gtpl 表示
      ├─ 一覧系
      │  └─ .gtpl 表示
      ├─ グラフ系
      │  └─ .gtpl / SVG 表示
      └─ 補助系
         └─ JSON / CSV / データ返却
```

この構成では、

- 親: `main` と ルーティング
- 子: handler 群
- 孫: テンプレート表示や API 応答

という感じです。

---

## 8. まとめ

今回の整理で重要なのは次の 3 点です。

1. URL を登録するのは [main.go](../main.go#L746-L830)
2. 全ハンドラーの前にかかるのは [main.go](../main.go#L288-L430) の middleware
3. 実際の機能本体は handler で、.gtpl はその描画先

この観点で見ると、ハンドラー中心の構造はかなり理解しやすくなります。

---

## 9. 追加で見ておくと良いもの

もし次の整理を進めるなら、以下の順番が自然です。

1. handler の一覧化
   - 「画面系」「グラフ系」「編集系」などへ分類
2. handler と template の対応表
   - どの handler がどの .gtpl を使うか
3. 役割の命名の整理
   - 「表示」「データ取得」「補助API」の境界を明確にする

現時点では、実装変更というより「整理と命名の整理」が中心です。

特に、今の整理の価値は「gtpl を見て handler を探す」のではなく、

- URL から handler を辿る
- handler から template を確認する

という見方を揃えることです。
