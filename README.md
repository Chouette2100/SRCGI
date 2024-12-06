# SRCGI

SHOWROOMのイベントでの参加者（配信者）の獲得ポイントの推移を表示するためのWebサーバー/CGIのソースです。

どのイベントのどの配信者をデータ取得の対象にするのか、またどの配信者のデータを表示するのか、リスナーの枠別貢献ポイントデータを取得するのか、などの設定をユーザーが自由に行うことができます。

このソースをダウンロード・ビルドして動かすために必要なデータベーススキーマやデータ取得プロセスのソースは別途公開します（[『SHOWROOM イベント 獲得ポイント一覧』関連のソースの公開について(1)](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/4fccae)）

現在 [SHOWROOM イベント 獲得ポイント一覧](https://chouette2100.com:8443/cgi-bin/SC1/SC1/top) で実際に動いているのを見ることができます。

ビルドから実行までの話は今のところ、UNIX/Linuxであれば

[【Unix/Linux】Githubにあるサンプルプログラムの実行方法](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/220e38)

Windowsであれば

[【Windows】かんたんなWebサーバーの作り方](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/c5cab5)

あたりをご覧ください。

そのうちちゃんとしたのを書くつもりです。

以下もご参照ください。

[SHOWROOMのAPI、その使い方とサンプルソース](https://zenn.dev/chouette2100/books/d8c28f8ff426b7)

[SHOWROOMのAPIで何ができるか？](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/84023c)

整合性のある最新バージョンの一覧（2024-12-06現在）
|機能名|種別|Ver.|機能|
|---|---|---|---|
|SRDB|SQL|v2.0.0|データベース作成用スキーマ|
|SRCGI|daemon<br>Webサーバ|v2.0.0|獲得ポイントデータの取得の設定、データの参照<br>（ShowroomCGI は使用しない）|
|SRGSE5M|daemon|v2.0.0|獲得ポイントデータの取得<br>(GetScoreEvery5minutes は使用しない)|
|SRGPC|daemon|v-.-.-|配信枠別リスナー別貢献ポイントの算出<br>（GetPointsCont01 は使用しない）|
|SRGCE|cron|v2.0.0|新規イベントデータ追加、ブロックイベント・イベントボックスの展開<br>srAddNewOnesはこちらに統合した|
|UpdateUserSetProperty|cron|v-.-.-|配信者属性の更新<br>おもにSHOWRANK上位ルームの一覧を作るために使用する<br>（SRUUI/SRUEUI は使用しない）|
|SRGGR|cron|v-.-.-|ライバー王決定戦等の結果データを取得する|
|srapi|library|v2.0.1|SHOWROOMのAPIのラッパー|
|exsrapi|library|v2.0.0|共通ライブラリ|
|srdblib|library|v2.0.0|DBのアクセスのためのライブラリ|
|srhandler|library|v2.0.0|Webサーバー/CGIのハンドラー<br>（一部、SRCGIで使うものは含まない）|

