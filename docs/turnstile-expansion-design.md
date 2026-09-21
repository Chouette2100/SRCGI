# Turnstile 拡張設計書

## この文書の目的

この文書は、SRCGI で Cloudflare Turnstile の適用範囲を広げるときに、どこをなぜ修正するのかを整理するための設計書である。

既存の [TURNSTILE_README.md](../TURNSTILE_README.md)、[TURNSTILE_SESSION_DESIGN.md](../TURNSTILE_SESSION_DESIGN.md)、[TURNSTILE_SESSION_GUIDE.md](../TURNSTILE_SESSION_GUIDE.md) は、導入方法やセッション機構の説明が中心である。一方で、実際に保護対象ハンドラーを増やすときには、ハンドラー本体だけでなく、アクセスログ、テンプレート、セッション再利用、data endpoint の扱いまで連動している。この文書はその連動関係を固定するためのものである。

## 基本方針

- Turnstile の追加は「ハンドラーに 1 行足せば終わり」ではない。
- 保護対象の追加は、画面表示、検証状態の記録、再送時の requestid 引き継ぎ、セッションクッキー再利用、アクセス集計の意味付けまで含めて成立する。
- したがって、修正単位は「ハンドラー」「テンプレート」「共通検証処理」「アクセスログ初期化」の 4 層で考える。

## 全体フロー

Turnstile 対応ハンドラーの処理は大まかに次の流れで動く。

1. [main.go](../main.go) の `commonMiddleware()` が requestid を生成し、アクセスログ初期値を作る。
2. 同じ `commonMiddleware()` が、Turnstile 対象ハンドラーだけ `turnstilestatus = 2` を設定する。
3. 対象ハンドラーが [ShowroomCGIlib/TurnstileHandler.go](../ShowroomCGIlib/TurnstileHandler.go) の `CheckTurnstileWithSession()` を呼ぶ。
4. 共通関数は、セッションクッキーが有効なら challenge を省略し、無効なら `cf-turnstile-response` を検証する。
5. challenge が必要なら、ハンドラー固有のテンプレートと [templates/turnstilechallenge.gtpl](../templates/turnstilechallenge.gtpl) を使って画面を返す。
6. 検証成功後にハンドラー本体が継続し、`accesslog.turnstilestatus` を `0` に更新する。
7. challenge 再送で不要になった旧 requestid の行があれば削除し、アクセス集計上の二重計上を避ける。

この流れのどこかを省くと、Turnstile 自体は表示されても、アクセス集計や再送制御の意味が壊れる。

## なぜ middleware 側の修正が必要か

Turnstile を新しいハンドラーへ適用するとき、最初に確認すべき場所は [main.go](../main.go) の `commonMiddleware()` である。

### requestid を使う理由

- requestid は 1 リクエスト単位の識別子である。
- challenge 画面を返した後、同じユーザーが再送すると、challenge 前のアクセスログ行と challenge 後の行が並ぶ。
- この 2 本を区別しないと、Turnstile 失敗や pending の解釈があいまいになる。

そのため、middleware は requestid を context に入れ、各ハンドラーはそれを `RequestID` としてテンプレートや DB 更新に使う。

### turnstilestatus を事前に `2` にする理由

- Turnstile 対象ハンドラーは、実行前の時点では「まだ検証が済んでいない」状態である。
- この状態をアクセスログに残すために、対象ハンドラーだけ `turnstilestatus = 2` で開始する。
- 検証成功後に `0` へ更新することで、「通常通過したアクセス」と「challenge 止まりのアクセス」を分けられる。

したがって、新しく Turnstile を適用するハンドラーは、[main.go](../main.go) の対象ハンドラー switch に追加しなければならない。ここを追加しないと、アクセス集計上は「保護されていない正常アクセス」と見なされる。

## ハンドラー側で修正すべきこと

### 1. ページデータ構造体に Turnstile 用フィールドを足す

代表例は [ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go) と [ShowroomCGIlib/HandlerClosedEvents.go](../ShowroomCGIlib/HandlerClosedEvents.go) である。

最低限必要なのは次の 3 項目である。

- `TurnstileSiteKey`
- `TurnstileError`
- `RequestID`

必要な理由は次のとおり。

- `TurnstileSiteKey`: challenge を表示するかどうかのテンプレート条件に使う。
- `TurnstileError`: 検証失敗時に再表示メッセージを出す。
- `RequestID`: challenge 再送時に元のアクセスログ行を追跡する。

### 2. `TurnstileChallengeData` を実装する

[ShowroomCGIlib/TurnstileHandler.go](../ShowroomCGIlib/TurnstileHandler.go) の共通処理は、ハンドラー固有の型に直接依存しない代わりに、`TurnstileChallengeData` を要求する。

そのため、新規対象ハンドラーでは少なくとも次を実装する必要がある。

- `SetTurnstileInfo()`
- `GetTemplatePath()`
- `GetTemplateName()`
- `GetFuncMap()`

必要な理由は次のとおり。

- challenge ページはハンドラー固有のテンプレートへ戻るので、共通処理だけではテンプレートパスが分からない。
- 共通 challenge 部分テンプレートを混ぜて描画するため、テンプレート名と FuncMap も必要になる。

### 3. ハンドラーの早い段階で `CheckTurnstileWithSession()` を呼ぶ

Turnstile 対象ハンドラーでは、重い DB 取得や API 呼び出しの前に `CheckTurnstileWithSession()` を呼ぶ。

理由は 2 つある。

- bot に先に高コスト処理を踏ませないため。
- challenge が必要なとき、後続処理の途中状態を抱えたまま分岐させないため。

代表例:

- [ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go)
- [ShowroomCGIlib/HandlerContributors.go](../ShowroomCGIlib/HandlerContributors.go)
- [ShowroomCGIlib/HandlerEventTop.go](../ShowroomCGIlib/HandlerEventTop.go)

### 4. 検証成功後に `turnstilestatus = 0` へ更新する

challenge を通過して本処理に入ったことをアクセスログへ反映する必要があるため、各ハンドラーは成功後に `accesslog` を更新する。

これが必要なのは、アクセス統計の `legitimate_count` と `turnstile_fail_count` が `turnstilestatus` に依存しているからである。更新しなければ、正常に通った利用者まで失敗側へ数えられる。

### 5. `lastrequestid` を使って古い pending 行を消す

[ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go) や [ShowroomCGIlib/HandlerListCntrbHEx.go](../ShowroomCGIlib/HandlerListCntrbHEx.go) では、challenge 再送時に `lastrequestid` を受け取り、旧 requestid の行を削除している。

これは単なる後始末ではない。理由は次のとおり。

- 初回アクセス時点で pending 行が作られる。
- challenge 通過後は、新しい requestid を持つ成功行が別にできる。
- 旧行を消さないと、同一閲覧が失敗 1 件として残りやすい。

したがって、challenge 後に再送が発生する画面型ハンドラーでは `requestid` hidden フィールドと `lastrequestid` 整理がセットで必要になる。

なお、元リクエストが GET であっても challenge 再送を POST に寄せる構成は有効である。`FormValue()` ベースのハンドラーであれば GET/POST の差異を吸収できるため、必要パラメータを hidden で引き継げば、状態を失わずに challenge 通過後の本処理へ戻せる。

## テンプレート側で修正すべきこと

テンプレートは大きく 2 パターンある。

### パターン 1: 画面内に challenge を直接埋め込む

代表例は [templates/currentevents.gtpl](../templates/currentevents.gtpl) である。

この型では、同じテンプレートの中で通常表示と challenge 表示を切り替える。

必要な要素は次のとおり。

- Turnstile スクリプトの読み込み
- `TurnstileSiteKey` を条件にした challenge 表示分岐
- `RequestID` を hidden で再送すること
- 必要なら元の入力値も hidden で再送すること

これらが必要な理由は次のとおり。

- スクリプトがないと widget 自体が描画されない。
- `TurnstileSiteKey` 条件がないと、通常画面と challenge 画面を分けられない。
- `RequestID` がないと accesslog 上の旧 requestid を整理できない。
- 元の入力値がないと、challenge 後に検索条件や表示条件が失われる。

実務上は、challenge 再送の method を POST で統一するケースが多い。この場合でも、ハンドラーが `FormValue()` で値を読む実装なら、元が GET の画面でも問題なく状態を復元できる。

### パターン 1-補足: 分割テンプレート（h1/h2/body など）

`list-cntrb` 系のように、複数の gtpl を順番に `ExecuteTemplate()` する構成では、challenge 分岐の置き場所に注意が必要である。

- `CheckTurnstileWithSession()` が challenge 表示時に描画するのは `GetTemplatePath()` で返したテンプレートを起点にした経路である。
- そのため、challenge 表示を担うテンプレート（多くは h1 側）単体で、challenge 用の HTML を完結させる必要がある。
- 通常表示は `{{else}}` 側へ寄せ、challenge が不要な場合のみ後続テンプレート（h2/body）へ進む形にする。

この前提を外すと、challenge 表示時に HTML が途切れたり、後続テンプレートとの責務分担が崩れたりする。

### パターン 2: 共通 challenge 部分テンプレートを使う

代表例は [templates/turnstilechallenge.gtpl](../templates/turnstilechallenge.gtpl) と、それを使う [ShowroomCGIlib/HandlerContributors.go](../ShowroomCGIlib/HandlerContributors.go) である。

この型では、challenge 描画の共通部分を別テンプレートに分ける。

その場合の注意点は次のとおり。

- ハンドラー側で `ParseFiles()` に本体テンプレートと `turnstilechallenge.gtpl` の両方を渡す。
- 共通 challenge テンプレートが参照するフィールドを、呼び出し元のデータ構造体が持っている必要がある。

前者が欠けると challenge テンプレート自体が解決できない。後者が欠けると hidden 項目や表示文言が空になり、再送や文脈保持が壊れる。

## 保護対象ハンドラーの分類

Turnstile を広げるときは、対象を 3 種類に分けて考えると判断しやすい。

### 1. 通常ページ型

例:

- [ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go)
- [ShowroomCGIlib/HandlerClosedEvents.go](../ShowroomCGIlib/HandlerClosedEvents.go)
- [ShowroomCGIlib/HandlerContributors.go](../ShowroomCGIlib/HandlerContributors.go)

特徴:

- 人間向けの HTML 画面を返す。
- challenge をその場で表示できる。
- requestid 再送とアクセスログ更新が素直に入る。

適用判断:

- bot に閲覧されるだけでも DB/API コストがかかる。
- ランキングやイベント一覧の大量巡回を防ぎたい。

### 2. 編集・高コスト画面型

例:

- [ShowroomCGIlib/HandlerEventTop.go](../ShowroomCGIlib/HandlerEventTop.go)
- [ShowroomCGIlib/HandlerGraphSum.go](../ShowroomCGIlib/HandlerGraphSum.go)
- [ShowroomCGIlib/HandlerGraphSum2.go](../ShowroomCGIlib/HandlerGraphSum2.go)

特徴:

- 画面表示自体に加えて、裏側で重い集計や編集文脈を持つ。
- 通されたあとの負荷や副作用が大きい。

適用判断:

- bot に踏まれると面倒なことになりやすい画面は、一覧系より優先度が高い。

### 3. data endpoint 型

例:

- [ShowroomCGIlib/HandlerGraphSum.go](../ShowroomCGIlib/HandlerGraphSum.go) の data handler
- [ShowroomCGIlib/HandlerGraphSum2.go](../ShowroomCGIlib/HandlerGraphSum2.go) の data handler 群

特徴:

- JSON やデータ片を返す。
- ここで challenge HTML を返すと呼び出し元 JavaScript が壊れる。

適用判断:

- UI 付き challenge は親ページ側で済ませる。
- data endpoint 側はセッションクッキーが有効かだけ確認する。

つまり、通常ページ型と同じ作りをそのまま流用してはいけない。

## data endpoint が例外扱いになる理由

data endpoint に通常の `CheckTurnstileWithSession()` をそのまま当てると、検証失敗時に HTML challenge が返る。その結果、呼び出し元は JSON を期待しているのに HTML を受け取り、画面が壊れる。

そのため、[ShowroomCGIlib/HandlerGraphSum.go](../ShowroomCGIlib/HandlerGraphSum.go) と [ShowroomCGIlib/HandlerGraphSum2.go](../ShowroomCGIlib/HandlerGraphSum2.go) の data handler 群は、通常ページより簡略化した扱いにしている。

- challenge の表示は親ページで行う。
- data endpoint はセッション確認のみ行う。
- セッション無効時は、画面遷移用 HTML ではなく data endpoint に合った失敗応答を返す。

この違いを無視すると、見た目には Turnstile を入れたつもりでも、非同期読込画面だけ壊れる。

## セッション管理まわりの前提

[ShowroomCGIlib/TurnstileSession.go](../ShowroomCGIlib/TurnstileSession.go) には、拡張時に意識しておくべき前提がある。

### SiteKey が空なら保護全体が無効になる

`CheckTurnstileWithSession()` は SiteKey が空のとき成功扱いで返る。これは開発時や一時停止時には便利だが、設定ミスでも同じ挙動になる。

したがって、「ハンドラーへ組み込んだか」と「設定で実際に有効化されているか」は別に確認する必要がある。

### IP 不一致は現在は拒否条件ではない

セッションクッキー検証では、IP 不一致はログに残すが即失敗にはしていない。これはプロキシや回線切替への実運用上の配慮である。

したがって、この実装は「同一 IP 固定を強く前提とする厳格セッション」ではなく、「実利用を壊しにくい緩めの再利用制御」である。設計レビューでは、この前提のまま守りたい画面かどうかを確認する。

### セッションの導入は UX 改善であり、保護責務の削除ではない

セッションがあるからといって、ハンドラー側の `turnstilestatus` 更新や requestid 管理が不要になるわけではない。セッションは challenge 頻度を下げるための仕組みであって、アクセスログの意味付けを肩代わりしない。

## FuncMap に関する注意

Turnstile 経由でテンプレートを描画するときは、通常の本処理とは別の経路でテンプレートが解決される。そのため、FuncMap の差分があると「通常表示では動くのに challenge 経由だと壊れる」という形で問題が出る。

この観点は [docs/funcmap-unification.md](./funcmap-unification.md) でも整理されている。

要点は次のとおり。

- 現行実装では、challenge 描画経路は `ShowroomCGIlib/TurnstileHandler.go` で `CloneCommonFuncMap()` をベースに描画する方針で統一している。
- したがって、この設計では「Turnstile challenge 経由画面は CommonFuncMap ベースで描画される」ことを前提にする。
- ハンドラー固有の FuncMap を使いたい場合は、まず CommonFuncMap に必要関数を寄せることを優先し、challenge 経路との乖離を作らない。

したがって、Turnstile 対象ハンドラーを増やすときは、challenge 画面経由でテンプレート関数が欠けないかを確認対象に入れる。

## どのハンドラーを優先して広げるか

優先度判断は、次の 2 軸で行う。

1. bot に踏まれたときのコストや副作用が大きいか。
2. 人間向け HTML 画面で、challenge を自然に出せるか。

この基準で見ると、優先度はおおむね次の順になる。

1. 編集系、高コスト集計系、内部状態を伴う画面
2. DB/API 負荷の大きい一覧系
3. data endpoint 単体

3 は単独で考えず、親ページと一組で導入する。

## 実装チェックリスト

新しいハンドラーに Turnstile を適用するときは、次を上から順に確認する。

- [main.go](../main.go) の対象ハンドラー switch に追加したか。
  理由: `turnstilestatus` を pending で開始しないとアクセス集計の意味が崩れる。
- ページデータ構造体に `TurnstileSiteKey`、`TurnstileError`、`RequestID` を追加したか。
  理由: challenge 表示、エラー表示、再送追跡に必要。
- `TurnstileChallengeData` を実装したか。
  理由: 共通検証処理がテンプレート情報へ到達するために必要。
- 重い処理の前に `CheckTurnstileWithSession()` を呼んでいるか。
  理由: bot に高コスト処理を踏ませないため。
- 検証成功後に `accesslog.turnstilestatus = 0` を更新しているか。
  理由: 正常通過の記録に必要。
- challenge 再送時の `requestid` hidden と `lastrequestid` 後始末があるか。
  理由: pending 行の残骸を防ぐため。
- テンプレートに Turnstile スクリプト、widget、必要な hidden 項目があるか。
  理由: challenge 表示と元状態の再送に必要。
- GET 画面でも challenge 再送を POST で扱う方針にしているか。
  理由: 再送経路を統一しつつ、hidden 項目で状態を保持できるため。
- challenge フォームに載せる hidden 項目が、画面状態（検索条件・ページ位置）を維持できる最小集合になっているか。
  理由: `eventid` / `userno(userid)` / `ie` / `nmonths` / `minpoint` / `maxnolines` / `ext` などが欠けると、通過後に表示条件が崩れる。
- 共通 challenge テンプレートを使う場合、`ParseFiles()` に本体と [templates/turnstilechallenge.gtpl](../templates/turnstilechallenge.gtpl) を両方渡しているか。
  理由: challenge 部分テンプレートが解決できないと描画に失敗する。
- 分割テンプレート構成の場合、challenge 分岐側のテンプレート単体で HTML が完結するか。
  理由: challenge 表示時は後続テンプレートが実行されないため。
- data endpoint を持つ画面かどうかを確認したか。
  理由: data endpoint は通常画面と同じ導入方法では壊れることがある。
- challenge 経路が CommonFuncMap ベースで描画される前提と、ハンドラー側テンプレート利用関数が矛盾していないか。
  理由: 通常表示と challenge 表示で使える関数差分が出る事故を防ぐため。

## 実装例リンク

実際の適用例を確認するときは、次の組み合わせを参照すると差分の意図が追いやすい。

- 典型的な一覧系（単一テンプレート）:
  [ShowroomCGIlib/HandlerCurrentEvents.go](../ShowroomCGIlib/HandlerCurrentEvents.go) と [templates/currentevents.gtpl](../templates/currentevents.gtpl)
- 共通 challenge テンプレート利用:
  [ShowroomCGIlib/HandlerContributors.go](../ShowroomCGIlib/HandlerContributors.go) と [templates/contributors.gtpl](../templates/contributors.gtpl)、[templates/turnstilechallenge.gtpl](../templates/turnstilechallenge.gtpl)
- 分割テンプレート（h1/h2/body）で challenge を h1 側に完結させる例:
  [ShowroomCGIlib/HandlerListCntrb.go](../ShowroomCGIlib/HandlerListCntrb.go) と [templates/list-cntrb-h1.gtpl](../templates/list-cntrb-h1.gtpl)、[templates/list-cntrb-h2.gtpl](../templates/list-cntrb-h2.gtpl)、[templates/list-cntrb.gtpl](../templates/list-cntrb.gtpl)
- 分割テンプレート（API貢献ランキング版）:
  [ShowroomCGIlib/HandlerListCntrbEx.go](../ShowroomCGIlib/HandlerListCntrbEx.go) と [templates/list-cntrbex-h1.gtpl](../templates/list-cntrbex-h1.gtpl)、[templates/list-cntrbex-h2.gtpl](../templates/list-cntrbex-h2.gtpl)、[templates/list-cntrbex.gtpl](../templates/list-cntrbex.gtpl)
- GET 画面を challenge 後 POST で再送する例（検索条件 hidden 保持）:
  [ShowroomCGIlib/HandlerListenerCntrbHistory.go](../ShowroomCGIlib/HandlerListenerCntrbHistory.go) と [templates/listener-cntrb-history.gtpl](../templates/listener-cntrb-history.gtpl)
- 同上（ルーム別履歴）:
  [ShowroomCGIlib/HandlerRoomCntrbHistory.go](../ShowroomCGIlib/HandlerRoomCntrbHistory.go) と [templates/room-cntrb-history.gtpl](../templates/room-cntrb-history.gtpl)
- middleware の pending 初期化対象:
  [main.go](../main.go) の `commonMiddleware()` にある `turnstilestatus` 初期化 switch

## 関連文書

- 導入手順や基本的な構成は [TURNSTILE_README.md](../TURNSTILE_README.md)
- セッション設計の背景は [TURNSTILE_SESSION_DESIGN.md](../TURNSTILE_SESSION_DESIGN.md)
- 設定値や運用上の見方は [TURNSTILE_SESSION_GUIDE.md](../TURNSTILE_SESSION_GUIDE.md)
- FuncMap 差分の注意は [docs/funcmap-unification.md](./funcmap-unification.md)

## 補足

この文書は、Turnstile をどの画面へ入れるかの判断そのものを固定するものではない。固定したいのは、適用すると決めたあとに見落としてはいけない責務の位置である。

つまり、保護対象の選定は運用判断だが、選定後の修正箇所は設計として標準化しておくべき、というのがこの文書の前提である。