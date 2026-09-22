# Endpoint Navigation と Bot Targeting の整理

作成日: 2026-09-22
対象: SRCGI

---

## 1. 今回の整理軸

今回知りたいのは、

- どのハンドラーがどの画面から起動されるか
- どの画面遷移から最終的に `ListCntrbHExHandler()` に到達するか
- その経路がボットから見て「狙われやすい」構造なのかどうか

という点です。

この観点では、関数の内部呼び出しだけでなく、

- 入口 URL
- 画面遷移
- 最終ハンドラー

という「エンドポイント親子関係」で見るのが適切です。

---

## 2. 入口の定義

URL から handler への接続は [main.go](../main.go#L746-L830) にあります。

例えば、次のような接続が存在します。

- `/top` → `TopHandler`
- `/eventtop` → `EventTopHandler`
- `/list-last` → `ListLastHandler`
- `/graph-sum` → `GraphSumHandler`
- `/list-cntrbex` → `ListCntrbExHandler`
- `/list-cntrbhex` → `ListCntrbHExHandler`

この層は「画面遷移の入口」になります。

---

## 3. 典型的な遷移パターン

あなたが指摘している経路は、実際のコードからも自然な流れです。

```text
TopHandler
  └─ CurrentEventsHandler
      └─ ListLastHandler
          └─ ListCntrbExHandler
              └─ ListCntrbHExHandler
```

### 実コードの対応

- [ShowroomCGIlib/HandlerTop.go](../ShowroomCGIlib/HandlerTop.go#L41-L151)
  - `TopHandler()` はトップ画面を返す
- [ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go#L93-L195)
  - `CurrentEventsHandler()` は開催中イベント一覧を返す
- [ShowroomCGIlib/HandlerListLast.go](../ShowroomCGIlib/HandlerListLast.go#L74-L176)
  - `ListLastHandler()` は最新スコア一覧を返す
- [ShowroomCGIlib/HandlerListCntrbEx.go](../ShowroomCGIlib/HandlerListCntrbEx.go#L103-L220)
  - `ListCntrbExHandler()` は配信者別ランキングを返す
- [ShowroomCGIlib/HandlerListCntrbHEx.go](../ShowroomCGIlib/HandlerListCntrbHEx.go#L115-L260)
  - `ListCntrbHExHandler()` は個別リスナーの履歴や細分化された貢献情報を返す

この流れは単なる関数呼び出しではなく、

- 画面遷移の親子関係
- 入口から最終ページまでの経路

として見るべきです。

---

## 4. 画面遷移としての見方

ここでの見方は次のようなものです。

```text
トップ画面
└─ イベント一覧
   └─ 最新結果一覧
      └─ 参加者/配信者別ランキング
         └─ 個別ハンドラ: ListCntrbHExHandler
```

このとき、各ノードは単なる関数ではなく、「画面またはページの入口」として見る方が実務的です。

つまり、今回の目的は「どの画面から最終画面に到達するか」を可視化することです。

---

## 5. なぜボットが狙うのか

今回の観点では、次のような仮説が重要です。

- `ListCntrbHExHandler()` は DB ベースのデータからページを生成している
- そのページがイベントやユーザー、リスナーの組み合わせごとに広く生成されうる
- 結果として URL の組み合わせが多くなり、総ページ数が非常に大きくなりうる
- そのため、アクセスの対象として「大量に巡回されやすいページ」になりうる

これは単なる「ボットがたまたま来た」のではなく、

- URL が多い
- パラメータが多い
- 生成コストが高い
- DB から組み立てられる

という条件が揃うと、探索型アクセスの対象になりやすい、という見方です。

### 具体的な可能性

`ListCntrbHExHandler()` は `eventid`, `userno`, `tlsnid`, `name` などのクエリパラメータを使っているため、

- イベント数 × 配信者数 × リスナー数
- さらに履歴や期間の組み合わせ

が増える可能性があります。

実際に、これらが DB から動的に作られるなら、

- ページ数が 100 万単位に近づく可能性がある
- その場合、単純な計算上のページ数だけでボットが「狙う」ように見える

という説明は十分にありえます。

> これは「ボットが意図的に狙っている」と断定するよりも、
> 「大量生成可能なページであり、巡回の対象になりやすい」という見方が正確です。

---

## 6. 画面遷移とボット対策の位置づけ

この観点で対策を考えると、次のように分けられます。

### 6.1 ハンドラ自体の対策

`ListCntrbHExHandler()` 自身に対して

- Turnstile
- Bot 判定
- rate limit
- 一部のパラメータに対する制限

を入れるのは当然有効です。

### 6.2 上流経路の対策

しかし、実際には上流の入口も重要です。

- `TopHandler()`
- `CurrentEventsHandler()`
- `ListLastHandler()`
- `ListCntrbExHandler()`

これらの経路から `ListCntrbHExHandler()` に辿り着くケースが多いなら、

- どこから始まっているかを監視する
- 入口レベルでアクセス元の傾向を見て制御する
- そのページ遷移のユースケースを限定する

という対策も効果的です。

### 6.3 重要な視点

「狙われているのは最終ハンドラだけ」ではなく、

- そのハンドラに到達する経路が複数ある
- 上流の画面から大量に生成される URL へ導かれている

という構造の方が、実際の攻撃コストとして大きいです。

---

## 7. 結論

今回の整理で重要なのは、

- まず `main.go` の入口から見て、URL がどこから handler に繋がるかを確認する
- 次に、遷移順を `Top → CurrentEvents → ListLast → ListCntrbEx → ListCntrbHEx` のように整理する
- そのうえで、最終ハンドラだけでなく、その経路の上流にも対策を置く

という視点です。

特に `ListCntrbHExHandler()` は、

- DB から動的生成される
- パラメータの組み合わせが大きくなりうる
- 画面遷移の最終地点になりやすい

という性質があり、ボットから見て「探索しやすい」ページになっている可能性があります。

この仮説は、単なる印象ではなく、コードの構造から十分に説明可能です。

---

## 8. 次の実務アクション候補

次の一歩としては、以下が自然です。

1. 画面遷移図の追加
   - `TopHandler` から `ListCntrbHExHandler` に至る経路を図として書く
2. URL/パラメータの一覧化
   - `eventid`, `userno`, `tlsnid`, `name`, `ie` などの意味を整理する
3. Bot 対策の適用箇所の整理
   - 入口側 vs ハンドラ側 vs middleware 側で何を入れるかを分離する
4. 生成可能範囲の評価
   - どのパラメータが組み合わせ爆発しうるかを評価する

この作業を進めると、対策は「最終ハンドラだけ」に偏らず、

- 入口
- 画面遷移
- 最終生成ページ

の 3 区分で設計できます。

---

## 9. まとめ

今回の整理の本命は、

- 「関数の親子」ではなく
- 「エンドポイントの親子 / 画面遷移の親子」

で見ることです。

そして `ListCntrbHExHandler()` がボットに狙われているように見える理由は、

- 画面遷移の最終地点である
- DB ベースで動的生成される
- 生成パターンが多くなりうる

という構造に由来している可能性が高いです。

この見立てで次の整理を進めるのが最も実用的です。
