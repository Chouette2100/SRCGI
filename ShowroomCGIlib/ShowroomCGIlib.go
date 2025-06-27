// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	"bytes"
	"fmt"
	"regexp"

	//	"html"
	"log"

	//	"math/rand"
	//	"sort"
	"strconv"
	"strings"
	"time"

	//	"bufio"
	// "os"

	//	"runtime"

	"encoding/json"

	//	"html/template"
	"net/http"

	//	"database/sql"

	//	_ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	"github.com/Chouette2100/exsrapi/v2"
	//	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

/*
0100L1	安定版（～2021.12.26）
0100M0	vscodeでの指摘箇所の修正
0101A0	LinuxとMySQL8.0に対応する。
0101B0	OSとWebサーバに応じた処理を行うようにする。アクセスログを作成する。
0101B1	実行時パラメータをファイルから与えるように変更する。
0101C0	GetRoomInfoByAPI()に配信開始時刻の取得を追加する。
0101D0	詳細なランク情報の導入（Nrank）
0101D1	"Next Live"の表示を追加する。
0101D2	GetScoreEvery5Minutes RU20E4 に適合するバージョン
0101D3	ランクをshow_rank_subdividedからleague_labe + lshow_rank_subdivided にする。
0101E1	環境設定ファイルをyaml形式に変更する。
0101G0	配信枠別貢献ポイントを導入する。
0101G1	list-last.gtplでは維新枠別貢献ポイントの記録があるルームのみリンクを作成する。
0101G2	list-last.gtplにジャンルを追加した。
0101G3	リスナー貢献ポイントの履歴の表示(list-cntrbH)を作成する。
0101G4	一つの貢献ポイントランキングの表示(list-cntrbS)を作成する(リスナー名の突き合わせのチェックが主目的)
0101G5	list-lasth.gtplのリロード予告の表示でデータ取得間隔が5分と固定されていたものを設定値に合わせるように変更する。
0101G6	ModminがIntervalminに対して不適切な値のときは修正して保存する。
0101G7	ランクに関しnext_scoreに加えprev_scoreの表示を追加する。ファンの数の表示を追加する。
0101H0	ファンレベル(/fanlevel)に関する画面を追加する。
0101J0	ファンダム王イベント参加者のファン数ランキングを作成する。
0101J2	終了したイベントについては無条件にルーム詳細情報（ランキング、フォロワ、レベル、ファン数）を出力しない。
0101J2a	"ルーム詳細情報”の説明を追加した。
0101J3	イベントリストのRoom_IDに変えてルーム名を表示する。表示数を6から10にする。
0101J4	NewDocument()をNewDocumentFromReader()に変更する。list-last_h.gtplにルーム情報詳細表示/非表示のボタンを追加する。
0101J5	イベント選択（最近のイベント）にnobasisが0のイベントも表示する（テーブルuserにusernoが0のデータを追加することが必要）
10AA00	枠別貢献ポイントの「目標値(推定)」を追加する。
10AB00	枠別貢献ポイントのポイント、増分の表示でhumanaizeを使用する。リスナー別貢献ポイント履歴に達成状況欄を追加する。
10AB01	WebserverをDbconfig.Webserverに置き換える。枠別貢献ポイントのGTPLを変更する。GTPLのCRLFをLFに変更する。
10AC00	イベント一覧にModminとModsecを追加する。
10AD00	longnameの初期値をusernameに変更する。Apache2LinuxをApache2Ubuntuに訂正する。
10AD01	HandlerGraphPerDay()とHandlerGraphPerSlot()のApache2LinuxをApache2Ubuntuに訂正する。
10AD02	SelectScoreList()で前配信期間、獲得ポイントのデータがないとき上位のデータがコピーして使われないようにする。
10AE00	1. 貢献ポイントランキングの表示を詳細化する。　2. イベント情報ページの内容を表示できるようにする（次回更新準備）
10AE00a	ルーム情報に下位ランクとの差を追加する（GTPLのみの変更）
10AF00	ブロックランキングに対応する（Event_id=30030のみの暫定対応）
10AF01	基準となる配信者のリストのサイズが0のときは対応する処理を行わない（異常終了対策、通常の運用では起きない）
10AG00	イベントリストの取得件数を10から20に変更する。
10AH00	stmt,err := Db.Prepare("...")　と対になる defer stmt.Close() を追加する。
10AJ00	ブロックランキングに仮対応（Event_id=30030以外に拡張）する。イベントリストの表示イベント数を設定可能とする。
10AK00	ブロックランキングのイベント名にblockidを追加する。獲得ポイントの取得時刻の初期値を分散させる。
------------------------------------- 以下公開版 ----------------------------------------------
10AL00	イベント情報にieventid（本来のイベントID、5桁程度の整数）を追加する。
10AL01	イベントの配信者リストを取得するとき、順位にかかわらず獲得ポイントデータを取得する設定とする。
10AL02	イベントの配信者リストを取得するとき、順位にかかわらず獲得ポイントを表示する設定とする。
10AM00	Room_url_keyから取り除く文字列を"/"から"/r/"に変更する。
10AN00	ブロックランキングで貢献ポイントランキングへのリンクを作るときはイベントIDのからブロックIDを取り除く。
10AP00	DBサーバーに接続するときSSHの使用を可能にする。
10AQ00	GetWeightedCnt()で周回数の多い獲得ポイントの採用率が上がるように調整する。
10AQ01	MakePointPerSlot()のperslotの変数宣言をループの中に入れる（毎回初期化されるように）
11AA0l	データベースへのアクセスをsrdblibに移行しつつある。グラフ表示で縮尺の設定を可能とする。
11AA02	intervalmin の値を5固定とする（異常終了に対する緊急対応）
11AA03	intervalminとintervalmin の適正でない入力を排除する。
11AB00	Event_Infの参照先をsrdblibからexsrapiに変更する。
11AB01	データベース保存時、Intervalminが0のときは強制的に5にする。
11AB02	データベース保存時、Intervalminが5でないときは強制的に5にする。
11AC00	開催中イベント一覧の機能を作成し関連箇所を修正する。
11AC01 FindPtPerSlot()でPrepare()に対するdefer Close()の抜けを補う。
11AC02 HandleListCntrb()でボーナスポイントに対する対応を行う。
11AC03 currentevent.gtpl 1行おきに背景色を変える。list-last_h.gtpl 結果が反映される時刻を正す。
11AD00 「SHOWROOMイベント情報ページからDBへのイベント参加ルーム情報の追加と更新」でイベントパラーメータがクリアされる問題を解決する。
11AE00	HandlerEventRoomList()でブロックイベントの参加ルーム一覧も表示できるようにする。
11AF00	開催予定イベント一覧の機能を追加する（HandlerScheduledEvent()）
11AF01	新規イベントの登録ができなくなった問題（＝11AD00の修正で発生したデグレード）に対応する
11AG00	srdblib.SelectFromEvent()の実行前にはsrdblib.Tevent = "event"を行う。 これはSelectFromEvent()の引数とすべき。
11AH00	HandlerCurrentEvent()で全イベント表示、データ取得中イベントのみ表示の切り替えを可能にする。
11AJ00	終了イベント一覧の作成でルームによる絞り込みを可能にする。
11AJ01	終了イベント一覧の作成でルームによる絞り込みを可能にする（不具合の修正）
11AJ02	開催中イベント、終了イベントに関する機能へのリンク切れを解消する。
11AJ03	終了イベントリスト、終了イベントルームリストの表示を改善する。
11AJ04	ページ遷移のレイアウトの共通化を行い、トップ画面を簡素化する。
11AK00	終了イベントでイベントIDとルームIDによる検索を可能にする。
11AL00	画面遷移のためのリンクを新しい機能に合わせる。list-cntrbSで目標値を変更できるようにする。
11AM00	開始前のイベントの登録は開催予定イベントのリストから行い、ルームの登録はイベント開始まで行わない件についてGetAndInsertEventRoomInfo()のフローを変更する。
11AN00	順位に関わりなくデータ取得の対象とするルームの追加でルーム検索を可能とするための準備を行う。
11AN01	api/room/profileでエラーを起きたときエラーの内容をログ出力する。
11AN02	HandlerNewUser() DBにユーザデータが存在しないときlongname、shortnameにAPIで取得した値をセットする。
11AP00	「最近のイベントの獲得ポイント上位のルーム」（HandlerTopRoom()）の機能を追加する。
11AP01	HandlerTopRoom()で日時範囲と表示数の設定を可能にする。
11AP02	GetUserInf()の抜けを補う。
11AQ00	掲示板機能を追加する。
11AQ01	掲示板機能について、HandlerWriteBbs()をDispBbsHandler()に統合し、リモートアドレス、ユーザーエージェントを保存する。
11AQ02	DispBbsHandler()に関して掲示板ページに直接来てもログが表示されるようにする。
11AQ03	終了イベント一覧の表示：51件表示し、50件ずつスクロールする。
11AQ04	ログメッセージを変更する（HandleListCntrb(),HandleListCntrbD(),HandleListCntrbH()）
〃     「(DB登録済み)イベント参加ルーム一覧（確認・編集）」で一覧にないルームを追加した直後の更新の不具合を修正する。
〃	    掲示板の「前ページ」、「次ページ」の操作を終了イベント一覧と同様にする。
11AQ05	Prepare()のあとのdefer stmt.Close()とdefer rows.Close()の抜けを補う。
11AR00	「枠別貢献ポイント一覧表」でリスナーさんの配信枠別貢献ポイントの履歴が表示されないことがある問題の修正。
〃    	ボット等からの接続を拒否（できるように）する。
11AS00	配信枠別貢献ポイントランキングでボット等から適正でないパラメータの要求を検出する。
11AT00	「イベント獲得ポイントランキング」でジャンルの指定を可能にする。
11AT01	MakePointPerDay()のログ出力を間引きする。
11AU00	終了したイベントの検索で、ルーム名、ルームIDで検索したとき、イベントの獲得ポイント上位のリストからイベント情報を見たとき該当ルームがどれかわかりやすくする。
11AV00	HandlerListLast()で確定値が発表されていないルームも表示するようにする。
11AV01	説明書きや表の項目名の修正
11AV02	scheduled-event.gtpl データ取得開始設定の説明を追加する。
11AW00	SelectCurrentScore() stmtを使いまわしているとことを別の変数にする。不具合ではないと思うが誤解を招きそうなので...
11AW01	SelectCurrentScore()の中のdeferでエラーが起きているか否かの検証を行う。
11AW02	説明書きや表の項目名の修正(追加)
11AX00	操作対象のテーブルをsrdblib.Teventで指定する方法から関数の引数とすべき方法に変える。
11AY00	HandlerShowRank()（SHOWランク上位配信者を表示する）を導入する。gorpを導入する。
11AZ00	userテーブルへのINSERTはsrdblib.InsertIntoUser()を用い、userテーブルのPDATEは原則として行わない。
11BA00	Genre, GenreIDの変更にともなう暫定対応（HandlerTopRoom()）+ showrank.gtpl の説明を追加する。
11BB00	未使用の関数GetIsOnliveByAPI()の定義を削除する。グラフ画像ファイル名を生成順の連番とする。
11BB01	過去イベントの検索でルーム名、IDから絞り込む場合は開催中のイベントも検索対象に含める。
11BB02	画像ファイル名はCGIの場合は連番、独立したWebサーバーの場合はPIDの下３桁とする。
11BC00	JSONのデコードが失敗したときのもとデータ（bufstr）のログ出力をやめる（APIが期待する結果を戻さない場合があることがわかっているから）
11BC01	終了済イベントのソート順はendtime descを優先する。
11BD00	UpdateRoomInf()でistargetとiscntrbpointを"N"に設定することを禁止する。
11BD01	獲得ポイント取得対象ルームの範囲を指定しての登録は1〜20に限定する。
11BD02	獲得ポイントの推移のグラフの画面に「表示するルームを選ぶ」というボタンを追加する。
11BD03	獲得ポイントの推移のグラフの画面の「表示するルームを選ぶ」に「グラフの色を変える」を追加する。
〃	    グラフ表示の最大ルーム数のデフォルト値を10から20に変更する。
11BE00	長期間に渡るイベントのグラフの表示方法を調整する。
11BE01	グラフ表示の最大ルーム数のデフォルト値を10から20に変更する（修正）
11BE02	list-last_h.gtplで「このページはブックマーク可能です」の文言を追加する。
11BF00	GraphScore01()でデータが連続していないとき（点になるとき）はcanvas.Circle()で描画する。
11BG00	GetAndInsertEventRoomInfo()でルーム情報の取得をGetEventsRankingByApi()を使う。block_id=0に対応する。
11BH00	HandlerGraphTotal()でグラフ線配色の初期化の機能を追加する。
11BH01	HandlerAddEvent()で起きているエラーの原因を特定するための情報を出力する。
11BH02	GetAndInsertEventRoomInfo()でeregがルーム数より大きいときはeregをルーム数に変更する。
11BH02a	GetAndInsertEventRoomInfo()でeregがルーム数より大きいときはeregをルーム数に変更する。
11BJ00	GetUserInf()でハンドラーが呼ばれたときのパラメータを表示する
11BJ01	top21.gtplで登録できる順位を20から50に拡張する（new-event0.gtplは20のままとする）
11BK00	HandlerEventList()がApiRoomStatus()とApiRoomNext()でエラーを起こしても処理を継続する。
11BM00	HandlerListGiftScore()を作成する
11BN00	HandlerListFanGiftScore()を作成する、HandlerGraphGiftScore()を準備する。
11BN01	HandlerListGiftScore()でGiftid（Grid）の選択を可能にする準備をする。
11BN02	HandlerListGiftScore()でmaxacqとlimitを可変にする。
11BN03	HandlerListFanGiftScore()でmaxacqとlimitを可変にする。
11BN04	DrawLineGraph()を作成する準備をする。
11BN05	list-gs-h1.gtpl, list-vgs-h1.gtpl のレイアウトを調整する。
11BQ00	ギフトランキングのグラフ（HandlerGraphGiftScore()）を作成する。
11BQ01	top.gtpl ギフトランキングのタイトルをより具体的にする
11BQ02	X軸の最小値を10,000から1,000に変更する
11BQ03	Viewerから（本来なかった）Ordernoを削除したことに対しSelectViewerid2Order()を修正する。
11BR00	ギフトランキング貢献ランキング（HandlerGiftScoreCntrb()）を作成する
11BS00	「修羅の道ランキング」（Giftid=13）のために表示の変更（獲得ポイントが取得できないため）
11BS01	ギフトランキング貢献ランキング（HandlerGiftScoreCntrb()）をギフトランキングから呼び出す
11BS02	グラフの凡例（ルーム名）の前に順位を表示する（すべてのルームのデータを表示するわけではないので）
11BS03	グラフ表示にあたらしいカラーマップを追加し、カラーマップの扱い方を変更する
11BT00	localhostからログインしたときは開催予定のイベントのルームを登録できる
11BU00	HandlerAddEvent()を分離し、バグを修正する。
11BT00	参加ルームの登録を行うときpoint==0のルームは除外する
11BW00	〃、獲得ポイント一覧（HandlerListLast()）でレベルイベントは順位のかわりにレベルを表示する
11BW01	HandlerAddEvent()）でレベルイベントの獲得ポイント0で除外したルームがルーム一覧に表示されないようにする。
11BV00	獲得ポイント全データのダウンロード機能（HandlerDlAllPoints()）を追加する。
11BX00	HandlerAddEvent()でweventを使ったあとeventに戻すようにする（獲得ポイントグラフの配色の初期化ができない問題の解決）

	ルームがなくてもイベント登録ができるようにする。

11BY00	HandlerListLast()でのデフォルトの表示を上位１５ルームにする。
11BY01	HandlerListLast()で15ルームのときは「もっと見る」ボタンを表示しない。
11BY02	HandlerListLast()で15ルームのときは「もっと見る」ボタンを表示しない。HadlerNewUser()を分離する。
11BX00	データ取得対象範囲の修正にともなってHandlerAddEvent()の一部機能をメンテナンス中とする。
11BX01	イベント単位で"足切り"を行う
11BZ00	アクセスログをDBに保存する。開催中イベント一覧でアクセス数の多いイベントをマークアップする。
11BZ01	SetThdata()をexsrapiへ移動する。HandlerListLast()でレベルイベントの表示順を検討する。
11BZ02	HandlerCurrentEvents()での強調表示の対象選択を変更する。
11BZ03	SetThdata()をReadThdata()とSetThdata()に分離する（SRGCIと共通）、HadlerNewEvent.goｗ別ファイルとする。
11BZ04	コメント化されたソースを削除する。
11BZ05	獲得ポイントダウンロードでのファイル名の誤り（yyyydddd）を正す（yyyymmdd）
--------- v2.0.0 ---------------
11BZ06	list-last_h.gtplで出力範囲を○位から○位までを○番目から○番目までに変更する。
11CA00	ShowroomCGIlib.goを機能別に分離する。thpoint = max(thinit, thdelta * hh) とする
11CB00	HandlerAddEvent()で、イベント開始後でもエントリー数が30以下のときはすべてのルームを登録する。
11CC00	累積・獲得ポイントの概要(HandlerGraphSum())を追加する
11CD00	累積・獲得ポイントの詳細(HandlerGraphSum2())を追加する
11CE01	HadleerClosedEvent()の一覧に貢献ポイントランキングへのリンクを追加する。
11CE02	Accesslogへの書き込みを非同期化する。
11CF00	貢献ランキングのCSVファイル出力を追加する
11CF01	CSVファイル出力の文字化けに対応する。攻撃的アクセスに対応する。終了イベント一覧に過去のイベントを追加・参照する機能を追加する（作成中）
11CF02	GetUserInf()で単一のIPアドレスから複数のリクエストがあったときは一定値以上のリクエストは拒否する。
11CG00	終了イベント一覧に過去のイベントを追加・参照する機能を追加する
11CG01	イベント結果確定後表示を修正する。
11CG03	HandlerShowRank()で過去のデータを除外する。

	closedevents.gtplでコメント一部が表示されないようにする。
	top.gtplで期間限定の表示を削除する。

11CH00	HandlerOldEvents()を実装する。
11CH01	HandlerOldEvents()でイベント数に矛盾があっても処理を継続する。
11CH02	HandlerClosedEvents()でブロックイベントが2つ出力されないようにする。
11CH03	HandlerClosedEvents()でtop.Limitを51固定とする。
11CH04	oldevents.gtplでの処理結果の表示はやめ、確認ボタンだけを表示する(timeout対策 <== WriteTimeoutが設定されていたのが原因)
11CJ00  github.com/Chouette2100 のパッケージをすべてv2に変更する。
11CK00  GetEventinf()をAPIを使う処理からDBを使う処理に変更する（現在の運用では処理対象のイベントのデータはDBに存在することを前提とできる）
11CK01  HandlerClosedEvents()でパラメータでlimit値を指定できるようにする(暫定対応)、bbsのエスケープ処理を除く、貢献ランキングリストをテーブルに。
11CK02  おそらくSetEventIDofOldEvents()のバグをHandlerShowRank()で暫定対応する
11CL00  closedevents.gtplでのタイトル、注釈の表示を修正する。

	    new-user.gtplでイベントへのリンクの"/"の抜けを修正する。"
		HandlerListLast()の表示をページングしたものをHadlerListLastP()とし、終了イベントの表jに使う。

11CM00	新規に登録されるイベントのEventUrlKeyを取得するをプログラムから非同期で起動する
11CM01	HandlerClosedEvents()での「次ページ」、「前ページ」の処理で limit を変更する処理をやめる、closedevents.gtplについても暫定対応を行う。
11CM02	HandlerTopRoom()での表示件数を30件から50件に変更する。
11CM03	Accesslog書き込み時のチャンネル操作の時間を調べる
11CM04	HandlerShowRank()で B-5ランクトップのユーザーを抽出するSQLの条件に irank != 0 を追加する。
11CM06  HandlerAddEvent()でGetEventInfAndRoomList()を使わず、GetEventQuestRoomsByApi()を使う。
11CN00  ハンドラーの関数名をHandlerXXX()からXXXHandle()に変更する。
11CN01  GetAndInsertEventRoomInfo()のレベルイベントの処理で獲得ポイントが0でないルームの前に0のルームがあるといったイレギュラーなケースに対する対応を行う。
11CN02  list-lastの「枠別貢献」の項目名に「貢献取得」へのリンクを追加する。任意のルームのSHOWランクを表示できるようにする。
11CP00  srdblib.GetFeaturedEvents()の抽出条件にイベント終了日時を追加する。最終結果、獲得ポイント一覧のレイアウトを変更する。
11CP01  終了イベントのルーム名による検索に説明を追加する。EvetRoomListHandler()のルームリスト取得を最新の手法にする準備。Chlogの大きさを調べる。
11CP02  貢献ポイント履歴のテストを行う。
11CQ00  貢献ランキングをAPIで取得して表示するListCntrbEx()を作成する。
11CQ03  貢献ポイントイベント履歴のレイアウトを変更する。
11CQ04  srdblib.Dberr を errに変更する。
11CQ05  go.modを作り直す。if文をswitch文に変更する。
11CQ06  貢献ランキング履歴の表示で現在開催中のイベントの表示は背景を黄色にする。「参加ルーム一覧」を「改修中」とする。
11CR03  ShowroomCGIlib.ServerConfig.LvlBotsを追加し、ボットの排除レベルを設定できるようにする。
11CR04  ShowroomCGIlib.ServerConfig.LvlBots == 3 のときはボットは無条件に排除する、　== 2 のときは特定のハンドラー(entry)のときボットを排除する。
11CR06  ListCntrbHandlerEx()の関数名をListCntrbExHandler()とする、bots.ymlとnotargetentry.ymlのデータを更新する。
11CS00  fail2banのログファイルをログ出力するようにする。GetUserInf()でのウェイト処理をやめる。
11CT00  短時間の連続的なアクセスに対してレート制限を行う。
11CT01  サーバー設定の初期化を行う（MaxChlog: ログ出力待ちチャンネルのバッファ数の定義の追加を含む）
11CT02  Gsum2, Gsum, LPSは枠別獲得ptデータがあるときだけリンクを有効にする。GsumData, Gsumdata1, Gsumdata2は監視対象外とする。

--------------------------------
11----	HandlerGraphOneRoom()を新規に作成する。
*/
const Version = "11CT02"

var Chimgfn chan int
var Chlog chan *srdblib.Accesslog

type LongName struct {
	Name string
}
type ColorInf struct {
	Color      string
	Colorvalue string
	Selected   string
}

type ColorInfList []ColorInf

type RoomInfo struct {
	Name      string //	ルーム名のリスト
	Longname  string
	Shortname string
	Account   string //	アカウントのリスト、アカウントは配信のURLの最後の部分の英数字です。
	ID        string //	IDのリスト、IDはプロフィールのURLの最後の部分で5～6桁の数字です。
	Userno    int
	//	APIで取得できるデータ(1)
	Genre      string
	Rank       string
	Irank      int
	Nrank      string
	Prank      string
	Followers  int
	Sfollowers string
	Fans       int
	Fans_lst   int
	Level      int
	Slevel     string
	//	APIで取得できるデータ(2)
	Order        int
	Point        int //	イベント終了後12時間〜36時間はイベントページから取得できることもある
	Spoint       string
	Istarget     string
	Graph        string
	Iscntrbpoint string
	Color        string
	Colorvalue   string
	Colorinflist ColorInfList
	Formid       string
	Eventid      string
	Status       string
	Statuscolor  string
}

type RoomInfoList []RoomInfo

// sort.Sort()のための関数三つ
func (r RoomInfoList) Len() int {
	return len(r)
}

func (r RoomInfoList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomInfoList) Choose(from, to int) (s RoomInfoList) {
	s = r[from:to]
	return
}

var SortByFollowers bool

// 降順に並べる
func (r RoomInfoList) Less(i, j int) bool {
	//	return e[i].point < e[j].point
	if SortByFollowers {
		return r[i].Followers > r[j].Followers
	} else {
		return r[i].Point > r[j].Point
	}
}

var Serverconfig *ServerConfig

/*
var Sshconfig *SSHConfig
var Dialer sshql.Dialer
*/

var Event_inf exsrapi.Event_Inf

var Regexpbots *regexp.Regexp
var NontargetEntry map[string]int

var OS string

//	var WebServer string

// idで指定した配信者さんの獲得ポイントを取得する。
// 戻り値は 獲得ポイント、順位、上位とのポイント差（1位の場合は2位とのポイント差）、イベント名
// レベルイベントのときは順位、上位とのポイント差は0がセットされる。
func GetPointsByAPI(id string) (Point, Rank, Gap int, EventID string) {

	//	獲得ポイントなどの配信者情報を得るURL（このURLについては記事参照）
	URL := "https://www.showroom-live.com/api/room/event_and_support?room_id=" + id

	resp, err := http.Get(URL)
	if err != nil {
		//	一時的にデータが取得できない。
		//		panic(err)
		return 0, 0, 0, "**Error** http.Get(URL)"
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		return 0, 0, 0, "**Error** http.Get(URL)"
	}

	//	イベントが終わっている、イベント参加をとりやめた、SHOWROOMをやめた、などの対応
	if result.(map[string]interface{})["event"] == nil {
		return 0, 0, 0, "not held yet./over./not entry."
	}

	if result.(map[string]interface{})["event"].(map[string]interface{})["ranking"] != nil {
		//	ランキングのあるイベントの場合
		//	（順位に応じて特典が与えられるイベント、ただし獲得ポイントに対して特典が与えられる場合でも順位付けがある場合はこちら）

		//	獲得ポイント
		l, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["point"].(float64)
		//	順位
		m, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["rank"].(float64)
		//	ポイント差
		n, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["gap"].(float64)

		Point = int(l)
		Rank = int(m)
		Gap = int(n)

		//	イベント名
		EventID, _ = result.(map[string]interface{})["event"].(map[string]interface{})["event_url"].(string)
		EventID = strings.Replace(EventID, "https://www.showroom-live.com/event/", "", -1)

	} else if result.(map[string]interface{})["event"].(map[string]interface{})["quest"] != nil {
		//	レベルイベント（ランキングのないイベント）の場合
		//	（アバ権やステッカーなど獲得ポイントに応じて特典が与えられるイベント、ただし順位付けがある場合は除く）

		//	獲得ポイント
		l, _ := result.(map[string]interface{})["event"].(map[string]interface{})["quest"].(map[string]interface{})["support"].(map[string]interface{})["current_point"].(float64)
		//	順位
		m := 0.0
		//	ポイント差
		n := 0.0

		Point = int(l)
		Rank = int(m)
		Gap = int(n)

		//	イベント名
		EventID, _ = result.(map[string]interface{})["event"].(map[string]interface{})["event_url"].(string)
		EventID = strings.Replace(EventID, "https://www.showroom-live.com/event/", "", -1)

	} else {
		//	上記ランキングイベントでもレベルイベントでもない場合
		//	もしこのようなケースが存在するならJSONを確認して新たにコーディングする
		log.Println(" N/A")
		return 0, 0, 0, "N/A"
	}

	return
}

func GetAciveFanByAPI(room_id string, yyyymm string) (nofan int) {

	nofan = -1

	url := "https://www.showroom-live.com/api/active_fan/room?room_id=" + room_id + "&ym=" + yyyymm

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		nofan = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		nofan = -2
		return
	}

	value, _ := result.(map[string]interface{})["total_user_count"].(float64)
	nofan = int(value)

	return
}
func GetRoomInfoByAPI(room_id string) (
	genre string,
	rank string,
	nrank string,
	prank string,
	level int,
	followers int,
	fans int,
	fans_lst int,
	roomname string,
	roomurlkey string,
	startedat time.Time,
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/profile?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	//	bufstr := buf.String()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(buf)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		//	log.Printf("%s", bufstr)
		log.Printf(" GetRoomInfoByAPI() decoder.Decode returned error %s\n", err.Error())
		status = -2
		return
	}

	value, _ := result.(map[string]interface{})["follower_num"].(float64)
	followers = int(value)

	tnow := time.Now()
	fans = GetAciveFanByAPI(room_id, tnow.Format("200601"))
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	fans_lst = GetAciveFanByAPI(room_id, fmt.Sprintf("%04d%02d", yy, mm))

	genre, _ = result.(map[string]interface{})["genre_name"].(string)

	rank, _ = result.(map[string]interface{})["league_label"].(string)
	ranks, _ := result.(map[string]interface{})["show_rank_subdivided"].(string)
	rank = rank + " | " + ranks

	value, _ = result.(map[string]interface{})["next_score"].(float64)
	nrank = humanize.Comma(int64(value))
	value, _ = result.(map[string]interface{})["prev_score"].(float64)
	prank = humanize.Comma(int64(value))

	value, _ = result.(map[string]interface{})["room_level"].(float64)
	level = int(value)

	roomname, _ = result.(map[string]interface{})["room_name"].(string)

	roomurlkey, _ = result.(map[string]interface{})["room_url_key"].(string)

	//	配信開始時刻の取得
	value, _ = result.(map[string]interface{})["current_live_started_at"].(float64)
	startedat = time.Unix(int64(value), 0).Truncate(time.Second)
	//	log.Printf("current_live_stared_at %f %v\n", value, startedat)

	return

}

func GetNextliveByAPI(room_id string) (
	nextlive string,
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/next_live?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	nextlive, _ = result.(map[string]interface{})["text"].(string)

	return

}

func GetRoomInfoAndPoint(
	eventid string,
	roomid string,
	idbasis string,
) (
	roominf RoomInfo,
	status int,
) {

	status = 0

	roominf.ID = roomid
	roominf.Userno, _ = strconv.Atoi(roomid)

	//	Event_inf, _ = SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	roominf.Genre, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers,
		roominf.Fans,
		roominf.Fans_lst,
		roominf.Name, roominf.Account, _, status =
		GetRoomInfoByAPI(roomid)

	point, _, _, peventid := GetPointsByAPI(roominf.ID)
	if peventid == Event_inf.Event_ID {
		roominf.Point = point
		UpdateEventuserSetPoint(peventid, roominf.ID, point)
	} else {
		log.Printf(" %s %s %d\n", Event_inf.Event_ID, peventid, point)
	}

	return
}
