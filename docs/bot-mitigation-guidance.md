# Bot 対策の指針

作成日: 2026-09-22
対象: SRCGI

---

## 1. 基本方針

今回の観点では、ボット対策を「最終ハンドラ単独でやる」のではなく、

- 入口
- 画面遷移
- 最終ページ生成

の 3 層で分けて設計するのがよいです。

特に `ListCntrbHExHandler()` は、

- DB から動的生成される
- パラメータの組み合わせが多い
- 遷移の最終地点になりやすい

という性質があるため、単一のハンドラ対策だけでは不十分な場合があります。

---

## 2. 重要な経路

対象の経路は、実際のコードから以下の流れが想定されます。

```text
TopHandler
  └─ CurrentEventsHandler
      └─ ListLastHandler
          └─ ListCntrbExHandler
              └─ ListCntrbHExHandler
```

対応ファイル:

- [main.go](../main.go#L746-L830): URL の入口
- [ShowroomCGIlib/HandlerTop.go](../ShowroomCGIlib/HandlerTop.go#L41-L151): `TopHandler()`
- [ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go#L93-L195): `CurrentEventsHandler()`
- [ShowroomCGIlib/HandlerListLast.go](../ShowroomCGIlib/HandlerListLast.go#L74-L176): `ListLastHandler()`
- [ShowroomCGIlib/HandlerListCntrbEx.go](../ShowroomCGIlib/HandlerListCntrbEx.go#L103-L220): `ListCntrbExHandler()`
- [ShowroomCGIlib/HandlerListCntrbHEx.go](../ShowroomCGIlib/HandlerListCntrbHEx.go#L115-L260): `ListCntrbHExHandler()`

---

## 3. 対策のレイヤー分け

### 3.1 入口レイヤー

入口は [main.go](../main.go#L288-L430) の `commonMiddleware` が担っています。

ここでやるべき対策:

- IP 単位の rate limit
- User-Agent 判定
- bot 判定
- fail2ban ログ
- 一般的なアクセスパターンの監視

目的:

- 入口で異常なアクセスを削る
- 上流のパスでボットを止める
- 最終ハンドラに到達する前に遮断できる場合がある

この層が強いと、最終的な `ListCntrbHExHandler()` への呼び出し数そのものを減らせます。

---

### 3.2 画面遷移レイヤー

画面遷移の上流では、特に次のようなアクセスを制御できると有効です。

- その画面から `ListCntrbHExHandler()` に遷移する必要がある人だけを通す
- `TopHandler` / `CurrentEventsHandler` / `ListLastHandler` の中で、次の画面に進む条件を強くする
- 直接的なリンクを経由しないアクセスを弾く

目的:

- リンク経由の正常アクセスと、スクレイピング型アクセスを分離する
- 画面遷移の自然さを担保する

特に、イベント一覧や最新一覧がボットの入口になっているなら、ここでの制御が効果的です。

---

### 3.3 最終ハンドラレイヤー

`ListCntrbHExHandler()` 自体にも対策を入れるべきです。

対象:

- Turnstile
- パラメータ検証
- リクエスト率制限
- 不正パラメータの拒否
- 一時的なキャッシュや結果の再利用

この層は最終防衛線です。

- 入口や遷移で止められなかったアクセス
- 直接叩かれたアクセス
- スクレイピングの末端

に対して効きます。

---

## 4. 具体的な設計指針

### 4.1 「入口で止める」設計

まずは多くの異常アクセスを入口で止めます。

- 一般ユーザーは通常の画面遷移順で来る
- ボットはいきなり最終ハンドラを叩く傾向がある
- その場合、入口の rate limit と bot 判定で止められる

設計意図:

- ボットが最終ページへ到達する前に遮断する
- 正常ユーザーに影響を出しにくくする

---

### 4.2 「遷移の自然さ」を担保する

画面遷移の中で、次のような制約を持たせると効果的です。

- ある画面から次の画面へ進むとき、必要なパラメータが揃っている
- 画面遷移の履歴やセッションを使って、ポストバックでないことを確認する
- 直接 URL での強制アクセスを許容しすぎない

これは、Bot からすると「画面遷移の自然な経路」を辿れなくなるため、無駄な巡回が減ります。

---

### 4.3 「最終ハンドラのコストを抑える」設計

`ListCntrbHExHandler()` は DB からデータを生成するため、コストが高いです。

対策としては以下が有効です。

- 生成結果をキャッシュする
- 同一パラメータの重複アクセスを抑制する
- 特定の組み合わせで大量アクセスを拒否する
- 画面の生成が大きすぎる場合はサマリ化する

今回のような頁数が増えやすい構造では、

- 1 回のリクエストが重い
- 生成数が増えやすい
- 調査型アクセスに弱い

ので、コスト抑制が非常に重要です。

---

## 5. 実際に有効な対策の順番

優先順としては次が合理的です。

1. `commonMiddleware` での bot / rate limit / access pattern 確認
2. 上流画面の遷移制御
3. `ListCntrbHExHandler()` 自体の Turnstile / rate limit / parameter validation
4. 生成コストの抑制（キャッシュ、重複制御、サマリ化）

この順番がよい理由は、

- 最も広い範囲を先に止められる
- 低コストで効果が出る
- 最終ハンドラの保守コストが上がらない

からです。

---

## 6. 重要な判断基準

対策を入れるべきかどうかを判断するポイントは、次の 3 つです。

- そのページが DB から動的生成されているか
- URL/パラメータの組み合わせが増えやすいか
- 一連の画面遷移から到達しやすいか

`ListCntrbHExHandler()` は 3 つとも満たしているため、警戒対象として妥当です。

---

## 7. まとめ

今後の対策の指針としては、

- 入口で止める
- 画面遷移で止める
- 最終ハンドラで止める

という 3 層構造で設計するのが最も合理的です。

`ListCntrbHExHandler()` は、単独の対策だけでなく、

- どこから辿れるか
- どのパラメータが増えやすいか
- どの入口がその画面へ導くか

を見ながら対策を入れるべきです。

特に、DB ベースで動的生成される大量ページの最終地点として捉えると、

- ボットの探索対象になりやすい
- 画面遷移の上流でも対策の余地がある

と理解しやすくなります。

---

## 8. 次の実務的な作業候補

次に進むなら、以下が有効です。

1. `ListCntrbHExHandler()` の URL パラメータ一覧を整理する
2. 各パラメータがどの画面から来るかを可視化する
3. `commonMiddleware` に入れるべき判定を整理する
4. 最終ハンドラに入れるべき対策を優先順位付きで列挙する

これが実際の対策設計に繋がります。
