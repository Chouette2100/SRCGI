package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	//	"io"

	// "golang.org/x/crypto/ssh/terminal"

	// "golang.org/x/term"

	"context"
	// "io"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"

	"crypto/tls"

	//	"html/template"
	"net/http"

	// "net/http/cgi"

	//	"github.com/dustin/go-humanize"

	"github.com/go-gorp/gorp"
	"github.com/google/uuid"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srcom"
	"github.com/Chouette2100/srdblib/v2"
	"github.com/Chouette2100/srhandler/v2"

	"SRCGI/ShowroomCGIlib"
)

/*
	0100a0	トップ画面をイベント選択/新規イベント登録の2つに絞った最初のバージョン（2021.04.02）
	0100a1	graph-totalのグラフの描画が残り2点になったときの処理の誤りを修正（2021.04.02）
	0100b0	データベース間の移行をスムーズにするためeventnoと(timeacq.)idxを使わず、eventidとtsを使うように修正した。(2021.04.22)
	0100c0	eventnoと(timeacq.)idxのinsertを行わないようにした。この2つのカラムの削除を前提としている(2021.04.23)
	0100d0	eventnoと(timeacq.)idxの2つのカラムを削除して動作を確認したもの。PerSlotのレイアウトの変更(2021.04.27)
	0100e0	ユーザ情報が変化したときだけhistoryに保存する(2021.05.06)
	0100f0	確定データを結果に反映するようにした。直近獲得ポイントのリストを詳細化した(2021.06.06)
	0100g0	グラフのオートスケーリングの変更、イベント登録直後の直近獲得ポイントリストでのpanic()対応(2021.06.06)
	0100g0	最終結果の前回獲得ポイント更新をConf.だけでなくProv.の場合も行うようにした(2021.06.13)
	0100h0	獲得ポイントの推移のグラフの表示範囲を設定できるようにした(2021.07.10)
	0100i0	トップ画面でユーザーの選択を行う(2021.07.11)
	0100j0	イベントパラメータの基準配信者をリスト化した(2021.07.11)
	0100k0	累積獲得ポイントの表示を追加（Perslot, Perday)(2021.07.14)
	0100l0	獲得ポイントリストでデータを取得していない配信者がいた場合の空白行挿入の処理を変更した(2021.11.23)
	0100l1	獲得ポイントリストのリロードを5秒遅らせる。10秒後から15秒後になるので50ルーム取得あたりまでOKなはず(2021.11.25)

	0100L1	安定版（～2021.12.26）
	0100M0	ライブラリ（ShowroomCGIlib.go）のバージョンも表示するようにした。
	0101A0	Linux、MySQL8.0、ローカルのLinux Mintの環境に対応する。
	0101B0	OSとWebサーバに応じた処理を行うようにする。
	0101C0	SSL対応
	0101C1	実行時パラメータをファイルから与えるように変更する。
	0101D0	ShowroomCGI 0101D2, GetScoreEvery5Minutes RU20E4, GetContPoints 20B02 適合バージョン
	0101F0	ShowroomCGIlibをサブディレクトリに移動する。
	0101F1	環境設定ファイルをyaml形式に変更する。
	0101G0	配信枠別貢献ポイントを実装する。
	0101H0	配信枠別貢献ポイントを実装する。
	0101J0	ファンダム王イベント参加者のファン数ランキングを作成する。 (Ver.1.0.0)
	0101K0	DBサーバーに接続するときSSHの使用を可能にする。
	0200A0	データベースへのアクセスをsrdblibに移行しつつある。
	0200A1	バージョンの表記にsrdblibのバージョンを追加する。
	0200A2	データ取得の間隔を0に設定できる問題に暫定的に対応する（5以外の入力を禁止する）
	0200B0	トップページに「イベント新規登録」へのページ内リンクを追加する。
	0201A0	srdblib.OpenDB()のインターフェース変更に対応する。
	0202A0	開催中イベント一覧の機能を作成し関連箇所を修正する。
	0202A1	rootpath($SCRIPT_NAME)とWebserverの設定の整合性をチェックする。
	00AA00	配信中ルーム一覧の機能（HandlerCurrentDistributions()）を追加する。
	00AA00	終了イベント一覧の機能（HandlerClosedEvents()）を追加する。
	00AB00	終了イベント一覧にイベント名の検索機能を追加する。各ページの上部にリンクボタンを追加する。
	00AC00	開催予定イベント一覧の機能(HadleScheduledEvents())を追加する。
	00AD00	srhandler.HandlerT008topForm(),srhandler.HandlerT008topForm()の呼び出しを追加する。
	00AE00	「最近のイベントの獲得ポイント上位のルーム」（HandlerTopRoom()）の機能を追加する。
	00AF00	掲示板機能を追加する。
	00AF01	掲示板機能について、HandlerWriteBbs()をDispBbsHandler()に統合し、リモートアドレス、ユーザーエージェントを保存する。
	00AG00	「枠別貢献ポイント一覧表」でリスナーさんの配信枠別貢献ポイントの履歴が表示されないことがある問題の修正。
			ボット等からの接続を拒否（できるように）する。
	00AG01	DenyIp.txtに関するログ出力を削除する。
	00AH00	ログファイル名を毎日午前0時に更新する。
	00AJ00	設定の追加　SetMaxOpenConns(8), SetMaxIdleConns(8), SetConnMaxLifetime(time.Second * 10)
	00AK00	ログファイル名変更のタイミングを（間違った午前9時から）午前0時に変更する。
	00AK01	SetConnMaxLifetime()に関するコメントを追加する。
	00AK02	SetConnMaxLifetime()の設定を10秒から20秒に変更する（HandlerTopRoom()のタイムアウト対策）
	00AK03	SetMaxOpenConns(8), SetMaxIdleConns(12), SetConnMaxLifetime(time.Minute * 5),SetConnMaxIdolTime(time.Minute * 5)
	00AL00	SHOWランクの一覧を表示できるようにする。gorpを導入する。
	00XX00	メンテナンス用
	00AM00	メンテナンス用の取り込み SRCGI.go でコメントに Maintenance とあるところを変更する
			SHOWROOMの2024年6月の仕様変更に /currentdistrb の機能をあわせる
	00AM01	メンテナンス用を取り込んだときbbs-1.gtplからbbs-1_org.gtplへの変更を忘れたところを修正する。
	00AM02	通常とメンテナンスの切り替えを ShowroomCGIlib.Serverconfig.Maintenance で行う。
	11BH00	HandlerGraphTotal()でグラフ線配色の初期化の機能を追加する。
	11BL00	srdblib.UpinsUserSetProperty()に対する srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")を追加する
	11BM00	HandlerListGiftScore()を作成する
	11BN00	HandlerListFanGiftScore()を作成する、HandlerGraphGiftScore()を準備する。
	11BN01	HandlerListGiftScore()でGiftid（Grid）の選択を可能にする。
	11BP00	旧URL（https/chouette2100.com:8443/cgi-bin/SRCGI/top）に対応する
	11BQ00	ギフトランキングのグラフ（HandlerGraphGiftScore()）を作成する。
	11BS00	「修羅の道ランキング」（Giftid=13）のために表示の変更（獲得ポイントが取得できないため）
	11BV00	獲得ポイント全データのダウンロード機能（HandlerDlAllPoints()）を追加する。
	11BZ00	アクセスログをDBに保存する
	11CC00	累積・獲得ポイントの概要(HandlerGraphSum())を追加する
	11CD00	累積・獲得ポイントの詳細(HandlerGraphSum2())を追加する
	11CE00	グラフ画像のファイル名の連番の発行はチャンネルを介して行う。
	11CE02	Accesslogへの書き込みを非同期化する。
	11CF00	貢献ランキングのCSVファイル出力を追加する
	11CF01	CSVファイル出力の文字化けに対応する。攻撃的アクセスに対応する。終了イベント一覧に過去のイベントを追加・参照する機能を追加する（作成中）
	11CF02	終了イベント一覧に過去のイベントを追加・参照する機能を追加する
	11CG00	commonMiddleware()を導入し、コンテクストとグレースフルシャットダウンを導入する。
	11CG02	サーバー起動パラメータを変更する。
	11CH00  指定した配信者の過去イベントの一覧を取り込む機能の調整をする。
	11CH01  main()のファイル名をmain.goと変更する。
	11CH05  http.Server.WriteTimeoutの設定を10秒から30秒に変更する。
	11CJ00  github.com/Chouette2100 のパッケージをすべてv2に変更する。
	11CL00  closedevents.gtplでのタイトル、注釈の表示を修正する。
        new-user.gtplでイベントへのリンクの"/"の抜けを修正する。"
		HandlerListLast()の表示をページングしたものをHadlerListLastP()とし、終了イベントの表jに使う。
	11CM02  タイムアウトの設定を30秒から90秒に変更する。
	11CM04  my_script.envを導入する
	11CM05  srdblibのバージョンをv2.3.2に変更する
	11CN00  ハンドラーの関数名をHandlerXXX()からXXXHandle()に変更する。
	11CN03  メンテナンスモードのエントリを"/"のみにする。
	11CP02  貢献ポイント履歴のテストを行う。
	11CQ00  貢献ランキングをAPIで取得して表示するListCntrbExHandler()を作成する。
	11CR00-2  特定useragentに対してメンテナンス中のレスポンスを返すようにする。
	11CR03  ShowroomCGIlib.ServerConfig.LvlBotsを追加し、ボットの排除レベルを設定できるようにする。
	11CR04  ShowroomCGIlib.ServerConfig.LvlBots == 3 のときはボットは無条件に排除する、　== 2 のときは特定のハンドラー(entry)のときボットを排除する。
	11CR05  排除したボットアクセス情報にハンドラー名を追加する。
	11CR06  ListCntrbHandlerEx()の関数名をListCntrbExHandler()とする、bots.ymlとnotargetentry.ymlのデータを更新する。
	11CS00  fail2banのログファイルをログ出力するようにする。GetUserInf()でのウェイト処理をやめる。
	11CS01  ボットとボットでない場合ののログ出力をunknownとnotabotに分ける。
	11CT00  短時間の連続的なアクセスに対してレート制限を行う。
	11CT01  サーバー設定の初期化を行う（MaxChlog: ログ出力待ちチャンネルのバッファ数の定義の追加を含む）
	11CT02  Gsum2, Gsum, LPSは枠別獲得ptデータがあるときだけリンクを有効にする。GsumData, Gsumdata1, Gsumdata2は監視対象外とする。
	11CT03  ログ出力をハンドラーがどう取り扱われたかわかりやすくする。　00:レートリミット、10:ボット、20:ハンドラー実行
	11CT04  リクエストの解析で抜けていたr.ParseForm()を追加する
	11CT07  commonMiddleware()のログ出力でボットに対してはPH-16、ボットでないものに対してはPH-20のログ出力を行う。
	11CU00  月別イベント・リスナー貢献ポイントランキングを新たに作成する。

	11CU03  トップ画面BBSの欄にバージョンを表示する。
	11CU04  MonthlyCntrbRankLgHandler()を新しく作る。ListCntrbHExHandler()でブロックイベントの結果が重複して表示される問題を修正する。
	11CU07  プロキシが使われている場合は、X-Forwarded-ForヘッダーからIPアドレスを取得するようにする。
	11CV00  Event_Infの使用をやめる
	200100  バージョンをGithubと共通化する
	200101  srdblibのバージョンをv2.4.1に変更する
	200102  アクセスログにrefererを追加する
	200300  Flutter/DataTableで使うJSONデータを出力するHandler(SearchUsersHandler())を追加する。
	200301  Flutter/DataTableで作った機能（search-users）を追加する。
	200400  EditCntrbPointsHandler() 枠別リスナー別貢献ポイントの一括登録機能を追加する。
	200500  experimental, todo のページを追加する(表示はされない)
	200501  list-cnrbHex.gtpl の表にグラフへのリンクを追加する。
	200600  ListToDoHandler(), EditToDoHandler(), InsertToDoHandler(), ExperimentalHandler()を追加する。
	200700  ClosedEventsHandler()のUIを改善する。closedevents.gtplのスタイルをcurrentevents.gtpl, scheduled-events.gtplにも適用する。
	200800  accessstatsのエンドポイントとそのハンドラー(AccessStatsHandler())を追加する。
	200802  ContributorsHandler()の使用を一時禁止する。
	200803  ContributorsHandler()の使用を取得済みデータについては開放する。
	200804  accesslogのtsをミリ秒単位で保存する。GraphSum2等のEventidの指定がない場合のエラーチェックを追加する。
	200900  アクセス集計表（ハンドラー別、ユーザーエージェント別、IPアドレス別）を追加する。
	200902  /ApiRoomStatus ハンドラーを悪意のあるアクセスから保護する。
	201000  TopFormHandler()をTopHandler()とEventTopHandler()に分割する。
	201100  Cloudflare Turnstileによるボット対策を導入する。
	201101  Turnstile関連の変数とeventid、roomidはFomrmvaluesに保存しない
	201102  Turnstileの連続した実行を避けるためセッション管理を導入する。
	201103  セッション管理、Turnstile検証処理を繁用関数化する。
	201105  FormValuesのログへの格納方法を変更する
	201106  requestidを用いてTurnstile検証の失敗をログに残す
	201112  Region(国・地域)をアクセスログに保存し、国外からのアクセスを遮断する機能を追加する。
	201113  Turnstile検証処理のあるハンドラーは実行前にal.Turnstilestatus = 2とし、実行時検証できたら =0 とする。
	201119  ログの時刻表示をμ秒単位に変更する。
	201200 ListenerCntrbHistoryHandler(), RoomCntrbHistoryHandler()を追加する。
	201201 ログ出力にCreateLogfile3()を使用するように変更する。LogWorker()のpanic()対応を追加する。
	201202 ログ出力切り替えタイミングの計算に使うtnowの値を現在時にする（バグの修正）
	201203 イベント一覧の注目のイベント表示にHighlightedフィールドを使用する。
	201204 イベント終了後の獲得ポイント一覧の終了の表示を明確にする
	201205 top画面の機能説明を修正する
	201211 新たにTurnstile認証を行うハンドラーをcommonMiddleware()のswitch文に追加する。
	201213 開催予定イベントの参加ルーム一覧を再作成する(完了)
*/

const version = "201213"

func NewLogfileName(logfile *os.File) {

	var err error

	//	毎日繰り返す
	for {

		tnow := time.Now()

		//	今日の午前9時
		// today := tnow.Truncate(24 * time.Hour)

		//	今日の午前0時
		// today = today.Add(-9 * time.Hour)
		//	test	today := tnow.Truncate(5 * time.Minute)

		//	次の日の午前0時
		// nextday := today.AddDate(0, 0, 1)
		nextday := time.Date(tnow.Year(), tnow.Month(), tnow.Day()+1, 0, 0, 0, 0, tnow.Location())
		//	test	nextday := today.Add(5 * time.Minute)

		log.Printf(" tnow: %s, nextday: %s\n", tnow.Format("2006-01-02 15:04:05"), nextday.Format("2006-01-02 15:04:05"))

		//	日付けが変わるまで待つ
		time.Sleep(nextday.Sub(tnow) - (1 * time.Minute))
		tnow = time.Now()
		log.Printf(" tnow: %s, nextday: %s\n", tnow.Format("2006-01-02 15:04:05"), nextday.Format("2006-01-02 15:04:05"))
		time.Sleep(nextday.Sub(tnow))

		//	ログファイルを閉じて新しいログファイルを作る
		logfile.Close()

		/*
			// logfilename := version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + time.Now().Format("20060102") + ".txt"
			//	test	logfilename := version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + time.Now().Format("20060102-1504") + ".txt"

			ShowroomCGIlib.VersionOfAll = version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + srapi.Version
			logfilename := ShowroomCGIlib.VersionOfAll + "_" + time.Now().Format("20060102") + ".txt"

			logfile, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				panic("cannnot open logfile: " + logfilename + err.Error())
			}

			// フォアグラウンド（端末に接続されているか）を判定
			isForeground := term.IsTerminal(int(os.Stdout.Fd()))

			if isForeground {
				// フォアグラウンドならログファイル + コンソール
				log.SetOutput(io.MultiWriter(logfile, os.Stdout))
			} else {
				// バックグラウンドならログファイルのみ
				log.SetOutput(logfile)
			}

			// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
			log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
		*/

		logfile, err = srcom.CreateLogfile3(version, ShowroomCGIlib.Version, srdblib.Version, srapi.Version)
		if err != nil {
			panic("cannnot open logfile: " + err.Error())
		}

		time.Sleep(1 * time.Second)
	}
}

type KV struct {
	K string
	V []string
}

// 共通の処理を行うミドルウェア
// =============================================
func commonMiddleware(limiter *SimpleRateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 共通のo処理をここで行う

		/*
			// remoteaddressを取得
			var ra string
			rap := r.RemoteAddr
			rapa := strings.Split(rap, ":")
			if rapa[0] != "[" {
				ra = rapa[0]
			} else {
				ra = "127.0.0.1"
			}
		*/

		// remoteaddressを取得
		ra := ShowroomCGIlib.RemoteAddr(r)
		rg := ShowroomCGIlib.FindRegionByIP(ra, ShowroomCGIlib.IPRegionList)

		// useragentを取得
		ua := r.UserAgent()

		if rg != "JP" && ShowroomCGIlib.Serverconfig.DenyNonJP {
			// 日本国内からのアクセスでない場合は処理を終了
			log.Printf("  *** PH-00 %s(), %s, \"%s\", region: %s\n", "Non-JP Access", ra, ua, rg)
			return
		}

		var entry string

		// next 関数(ハンドラー)の名前を取得
		// reflect.ValueOf(next).Pointer() で関数のエントリポイントのアドレスを取得
		// runtime.FuncForPC() でそのアドレスに対応する runtime.Func を取得
		funcPtr := reflect.ValueOf(next).Pointer()
		handlerFunc := runtime.FuncForPC(funcPtr)
		entry = "unknown" // デフォルト値

		if handlerFunc != nil {
			// runtime.FuncForPC().Name() は通常 "package/path.FunctionName" の形式を返します。
			// 例: "github.com/your/repo/ShowroomCGIlib.ListCntrbHExHandler"
			fullFuncName := handlerFunc.Name()

			// ユーザーが求めている "ShowroomCGIlib.ListCntrbHExHandler" のような形式は、
			// runtime.FuncForPC().Name() の出力からパス部分を除去した形に一致することが多いです。
			// 最後の '/' の位置を探します。
			lastSlashIndex := strings.LastIndex(fullFuncName, ".")
			if lastSlashIndex != -1 {
				// スラッシュより後ろの部分を取得します。
				// 例: "ShowroomCGIlib.ListCntrbHExHandler"
				entry = fullFuncName[lastSlashIndex+1:]
				// } else {
				// 	// スラッシュがない場合はそのままの名前を使います (mainパッケージなど)。
				// 	// 例: "main.ListCntrbHExHandler"
				// 	handlerName = fullFuncName
			}
		}

		// userAgent := r.Header.Get("User-Agent")
		ignBots := false
		lvlbots := ShowroomCGIlib.Serverconfig.LvlBots
		bmatch := false
		if bmatch = ShowroomCGIlib.Regexpbots.MatchString(ua); (lvlbots == 3 || lvlbots == 2) && bmatch {
			ignBots = true
		}

		// fail2banのログファイルに書き込む
		botInfo := "notabot"
		if bmatch {
			botInfo = "unknown"
		}

		if entry != "GraphSumDataHandler" && entry != "GraphSumData1Handler" && entry != "GraphSumData2Handler" {
			// これら３つのハンドラーは従属的なものなのでチェックの対象から外す
			// FIXME: これらのハンドラーを単独で実行するケースに対応する必要を検討すべき

			LogForFail2ban(ra, entry, botInfo)

			// レートリミッターで許可を判定
			if !limiter.Allow(ra) {
				// 制限を超過した場合、レスポンスを返さずに処理を終了（無視）
				// log.Printf("Ignoring request from rate-limited IP: %s %s %s", ra, r.Method, r.URL.Path)
				log.Printf("  *** PH-00 %s(), %s, %s, \"%s\"\n", entry, ra, rg, ua)
				// ここで何もせず return するのがポイントです。
				// http.ResponseWriter に何も書き込まれないため、net/http はレスポンスを送信しません。
				return
			}
		}

		// 例: 特定のUser-Agent（例えば、古いバージョンのアプリや特定のボットなど）に対してメンテナンスメッセージを返す
		// if ignBots ||
		// 	userAgent == "meta-externalagent/1.1 (+https://developers.facebook.com/docs/sharing/webmasters/crawler)" ||
		// 	userAgent == "Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)" ||
		// 	userAgent == "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Amazonbot/0.1; +https://developer.amazon.com/support/amazonbot) Chrome/119.0.6045.214 Safari/537.36" {
		if _, ok := ShowroomCGIlib.NontargetEntry[entry]; ignBots && (lvlbots == 3 || ok) {
			// log.Printf("Common processing: %s: %s\n", entry, ua)
			// メンテナンス中のステータスコードを設定
			w.WriteHeader(http.StatusServiceUnavailable)

			// メンテナンス中のメッセージをレスポンスボディに書き込む
			// Content-Typeをtext/htmlに設定すると、ブラウザでHTMLとして表示されます
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintln(w, `
<!DOCTYPE html>
<html>
<head>
    <title>メンテナンス中</title>
</head>
<body>
    <h1>現在メンテナンス中です</h1>
    <p>ご迷惑をおかけいたしますが、しばらくお待ちください。</p>
</body>
</html>`)

			log.Printf("  *** PH-10 %s(), %s, \"%s\"\n", entry, ra, ua)
			return // ここで処理を終了
		}
		// } else {
		// 	log.Println("Common processing")
		// }

		// // 通常の処理
		// botInfo := "notabot"
		// if bmatch {
		// 	botInfo = "unknown"
		// }
		// LogForFail2ban(ra, entry, botInfo)

		if bmatch {
			// ボットアクセスのログ出力	}
			log.Printf("  *** PH-16 %s(), %s, \"%s\"\n", entry, ra, ua)
		} else {
			// ボットでないアクセスのログ出力
			log.Printf("  *** PH-20 %s(), %s, \"%s\"\n", entry, ra, ua)
		}

		requestID := uuid.New().String()
		var al srdblib.Accesslog
		al.Ts = time.Now().Truncate(time.Millisecond)
		log.Printf(" Accesslog Ts: %s\n", al.Ts.Format("2006-01-02 15:04:05.000"))
		al.Handler = entry
		al.Remoteaddress = ra
		al.Useragent = ua
		al.Region = rg
		al.Requestid = requestID
		// if entry == "ContributorsHandler" {
		// 	al.Turnstilestatus = 2 // pending
		// } else {
		// 	al.Turnstilestatus = 0 // success or n/a
		// }

		switch entry {
		case "BadRequestHandler",
			"ContributorsHandler",
			"ClosedEventsHandler",
			"EventTopHandler",
			"GraphSumHandler",
			"GraphSumDataHandler",
			"GraphSum2Handler",
			"GraphSumData1Handler",
			"GraphSumData2Handler",
			"ListCntrbHHandler",
			"ListCntrbHExHandler":
			al.Turnstilestatus = 2 // pending => failed
		default:
			al.Turnstilestatus = 0 // success <= pending
		}

		al.Referer = r.Referer()

		//	パラメータを取得する
		if err := r.ParseForm(); err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		kvlist := make([]KV, 0)
		for k, v := range r.Form {
			log.Printf("%12v : %v\n", k, v)
			switch k {
			case "eventid", "ieventid":
				al.Eventid = v[0]
			case "userno", "userid", "user_id", "roomid":
				al.Roomid, _ = strconv.Atoi(v[0])
			case "cf-turnstile-response":
			default:
				kvlist = append(kvlist, KV{K: k, V: v})
			}
		}
		jd, err := json.Marshal(kvlist)
		if err != nil {
			log.Printf(" GetUserInf(): %s\n", err.Error())
		}
		al.Formvalues = string(jd)

		if lenchlog := len(ShowroomCGIlib.Chlog); lenchlog != 0 {
			log.Printf("            length of Chlog = %d\n", lenchlog)
		}

		ShowroomCGIlib.Chlog <- &al

		// 次のハンドラーを呼び出す
		// next(w, r)
		ctx := context.WithValue(r.Context(), "requestid", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// =============================================

// =============================================
// // サーバー構成
// type ServerConfig struct {
// 	HTTPport string
// 	SSLcrt   string
// 	SSLkey   string
// }

// =============================================

// 入力内容の確認画面
func main() {

	/*
		ShowroomCGIlib.VersionOfAll = version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + srapi.Version
		logfilename := ShowroomCGIlib.VersionOfAll + "_" + time.Now().Format("20060102") + ".txt"
		logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic("cannnot open logfile: " + logfilename + err.Error())
		}
		defer logfile.Close()
		defer LogFile_f2b.Close() // LogForFail2ban()で開いたファイルを閉じる

		// フォアグラウンド（端末に接続されているか）を判定
		isForeground := term.IsTerminal(int(os.Stdout.Fd()))

		if isForeground {
			// フォアグラウンドならログファイル + コンソール
			log.SetOutput(io.MultiWriter(logfile, os.Stdout))
		} else {
			// バックグラウンドならログファイルのみ
			log.SetOutput(logfile)
		}

		// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	*/

	logfile, err := srcom.CreateLogfile3(version, ShowroomCGIlib.Version, srdblib.Version, srapi.Version)

	//	ログファイル名を毎日午前0時に更新する
	go NewLogfileName(logfile)

	ShowroomCGIlib.Chlog = make(chan *srdblib.Accesslog, 100)
	defer close(ShowroomCGIlib.Chlog)
	go ShowroomCGIlib.LogWorker()

	//	https://ssabcire.hatenablog.com/entry/2019/02/13/000722
	//	https://konboi.hatenablog.com/entry/2016/04/12/100903

	svconfig := ShowroomCGIlib.ServerConfig{
		HTTPport:    "8080",
		SSLcrt:      "/home/username/.ssh/server.crt",
		SSLkey:      "/home/username/.ssh/server.key",
		WebServer:   "nginxSakura",
		NoEvent:     30,    // イベントの表示数
		Maintenance: false, // メンテナンスモード
		LvlBots:     2,     // ボット排除のレベル、0:なし、1:低、2:中、3:高
		AccessLimit: 3,     // タイムウィンドウあたりのリクエスト上限
		TimeWindow:  1,     // タイムウィンドウの長さ（秒）
		MaxChlog:    10,    // ログ出力待ちチャンネルのバッファ数
		DenyNonJP:   false, // 日本国内以外からのアクセスを拒否するかどうか
		GWURL:       "http://vscode01:8080/",
		// GWURL: "https://gwuu.chouette2100.com/",
	}
	ShowroomCGIlib.Serverconfig = &svconfig
	err = exsrapi.LoadConfig("ServerConfig.yml", ShowroomCGIlib.Serverconfig)
	if err != nil {
		log.Printf("err=%s.\n", err.Error())
		os.Exit(1)
	}
	if svconfig.NoEvent == 0 {
		svconfig.NoEvent = 30
	}
	log.Printf("%+v\n", svconfig)

	ShowroomCGIlib.Regexpbots = ReadBots()
	ReadEntry()

	ShowroomCGIlib.LoadDenyIp("DenyIp.txt")
	//	log.Printf("DenyIp.txt = %v\n", ShowroomCGIlib.DenyIpList)

	// 10秒間に5リクエストまで許可するレートリミッターを作成
	// この値は実際の状況に合わせて調整してください。
	rateLimiter := NewSimpleRateLimiter(svconfig.AccessLimit, time.Duration(svconfig.TimeWindow)*time.Second)

	/*
		var dbconfig *srdblib.DBConfig
		err = exsrapi.LoadConfig("DBConfig.yml", &dbconfig)
		if err != nil {
			log.Printf("err=%s.\n", err.Error())
			os.Exit(1)
		}
		log.Printf("%+v\n", dbconfig)
	*/

	switch svconfig.WebServer {
	case "nginxSakura":
		fallthrough
	case "Apache2Ubuntu":
		fallthrough
	case "None":
	default:
		log.Printf("Unknown WebServer = <%s> (must be nginxSakura, Apache2Ubuntu or None)\n", svconfig.WebServer)
		return
	}

	ShowroomCGIlib.OS = runtime.GOOS
	/*
		rootPath := ""
		if svconfig.WebServer != "None" {
			rootPath = os.Getenv("SCRIPT_NAME")
		}
	*/
	rootPath := os.Getenv("SCRIPT_NAME")
	if rootPath != "" && svconfig.WebServer != "None" {
		log.Printf("**error** rootPath is \"%s\", but WebServer is not \"None\".\n", rootPath)
	} else if rootPath == "" && svconfig.WebServer == "None" {
		log.Printf("**error** rootPath is \"\", but WebServer is \"None\".\n")
	}

	/*	設定ファイルで操作するはず？
		err = os.Setenv("HOME", "/var/www")
		if err != nil {
			log.Printf("os.Setenv(): err=%s.\n", err.Error())
			return
		}
	*/
	home := os.Getenv("HOME")
	log.Printf("\n")
	log.Printf("\n")
	log.Printf("********** WevServer=<%s> port = <%s> OS = <%s> rootPath = <%s> home = <%s>\n",
		svconfig.WebServer, svconfig.HTTPport, ShowroomCGIlib.OS, rootPath, home)
	log.Printf("********** crt=<%s> key = <%s>\n", svconfig.SSLcrt, svconfig.SSLkey)
	//	log.Printf("********** Dbhost=<%s> Dbname = <%s> Dbuser = <%s> Dbpw = <%s>\n", dbconfig.DBhost, dbconfig.DBname, dbconfig.DBuser, dbconfig.DBpswd)

	var dbconfig *srdblib.DBConfig
	dbconfig, err = srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("Database error. err = %v\n", err)
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}

	//  =============================================
	// 画像ファイル名の連番を発行する
	go func() {
		no := 0
		ShowroomCGIlib.Chimgfn = make(chan int)
		for {
			ShowroomCGIlib.Chimgfn <- no
			no++
			if no > 999 {
				no = 0
			}
		}
	}()
	//  =============================================

	//	http://dsas.blog.klab.org/archives/2018-02/configure-sql-db.html
	//	https://qiita.com/hgsgtk/items/770c51559f374b36da3f
	//	http://dsas.blog.klab.org/archives/pixiv-isucon2016-2.html
	//	SetConnMaxLifetime()は必要ないとするものも
	//	https://qiita.com/ichizero/items/36036dbd8a32ce23ca5b
	//	srdblib.Db.SetMaxOpenConns(8)
	//	srdblib.Db.SetMaxIdleConns(8)
	//	srdblib.Db.SetConnMaxLifetime(time.Second * 20)

	//	https://zenn.dev/kouhei_fujii/articles/72ac1f8d4e8a84
	srdblib.Db.SetMaxOpenConns(8)
	srdblib.Db.SetMaxIdleConns(12)

	srdblib.Db.SetConnMaxLifetime(time.Minute * 5)
	srdblib.Db.SetConnMaxIdleTime(time.Minute * 5)

	defer srdblib.Db.Close()
	log.Printf("%+v\n", dbconfig)

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db,
		Dialect:         dial,
		ExpandSliceArgs: true, //スライス引数展開オプションを有効化する
	}
	srdblib.Dbmap.AddTableWithName(srdblib.User{}, "user").SetKeys(false, "Userno")
	srdblib.Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
	srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	srdblib.Dbmap.AddTableWithName(srdblib.Wevent{}, "wevent").SetKeys(false, "Eventid")
	srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")
	srdblib.Dbmap.AddTableWithName(srdblib.Weventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	srdblib.Dbmap.AddTableWithName(srdblib.GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
	srdblib.Dbmap.AddTableWithName(srdblib.ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	srdblib.Dbmap.AddTableWithName(srdblib.Viewer{}, "viewer").SetKeys(false, "Viewerid")
	srdblib.Dbmap.AddTableWithName(srdblib.ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
	srdblib.Dbmap.AddTableWithName(ShowroomCGIlib.Contribution{}, "contribution").SetKeys(false, "Ieventid", "Roomid", "Viewerid")

	srdblib.Dbmap.AddTableWithName(srdblib.Campaign{}, "campaign").SetKeys(false, "Campaignid")
	srdblib.Dbmap.AddTableWithName(srdblib.GiftRanking{}, "giftranking").SetKeys(false, "Campaignid", "Grid")
	srdblib.Dbmap.AddTableWithName(srdblib.Accesslog{}, "accesslog").SetKeys(false, "Ts", "Eventid")

	srdblib.Dbmap.AddTableWithName(srdblib.Todo{}, "todo").SetKeys(false, "ID")

	fileenv := "Env.yml"
	err = exsrapi.LoadConfig(fileenv, &srdblib.Env)
	if err != nil {
		err = fmt.Errorf("exsrapi.Loadconfig(): %w", err)
		log.Printf("%s\n", err.Error())
		return
	}

	if svconfig.WebServer == "None" && !ShowroomCGIlib.Serverconfig.Maintenance {
		// WebServerがNoneの場合はURLにTopがないときpublic（のindex.html）が表示されるようにしておきます。
		http.Handle("/", http.FileServer(http.Dir("public")))
	}
	/*
		else {
			まずWebServer = "None"、つまりこのShowroomCGIをWebサーバーとして使うのがいいのですが、
			レンタルサーバーではこれができないと思います。
			その場合既設のWebサーバーを使うしかないので、その使用条件に合わせて"調整"してください。

			（レンタルサーバーというより自分でApache2をインストールした環境ということになると思いますが）たとえば Apache2 を
				# apt install apache2
				# a2enmod cgid
				# systemctl restart apache2
			とインストールしてCGIを使えるようにした環境だと
				/usr/lib/cgi-bin/
			にCGI（ShorroomCGI）を配置するわけですが、このディレクトリはCGI専用なので/usr/lib/cgi-bin/publicに作成したグラフが表示されません。
			グラフを表示するためには
				# cd /var/www/html
				# ln -s /usr/lib/cg-bin/public
			とhtmlのところにシンボリックリンクを置けばこのプログラムはそのままでグラフを表示することができます。

		}
	*/

	// =============================================
	// メインコンテキストとキャンセル関数を作成
	ctx, cancel := context.WithCancel(context.Background())

	// シグナルチャンネルを作成
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// ゴルーチンでシグナルを待機
	go func() {
		sig := <-sigs
		log.Println("Received signal:", sig)
		cancel()
	}()
	// =============================================

	if !ShowroomCGIlib.Serverconfig.Maintenance {

		// http.HandleFunc(rootPath+"/top", commonMiddleware(rateLimiter, ShowroomCGIlib.TopFormHandler))
		http.HandleFunc(rootPath+"/top", commonMiddleware(rateLimiter, ShowroomCGIlib.TopHandler))

		http.HandleFunc(rootPath+"/eventtop", commonMiddleware(rateLimiter, ShowroomCGIlib.EventTopHandler))

		http.HandleFunc(rootPath+"/list-level", commonMiddleware(rateLimiter, ShowroomCGIlib.ListLevelHandler))

		http.HandleFunc(rootPath+"/list-last", commonMiddleware(rateLimiter, ShowroomCGIlib.ListLastHandler))
		http.HandleFunc(rootPath+"/list-lastP", commonMiddleware(rateLimiter, ShowroomCGIlib.ListLastPHandler))

		http.HandleFunc(rootPath+"/dl-all-points", commonMiddleware(rateLimiter, ShowroomCGIlib.DlAllPointsHandler))

		http.HandleFunc(rootPath+"/graph-total", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphTotalHandler))
		http.HandleFunc(rootPath+"/csv-total", commonMiddleware(rateLimiter, ShowroomCGIlib.CsvTotalHandler))

		http.HandleFunc(rootPath+"/graph-dfr", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphDfrHandler))

		http.HandleFunc(rootPath+"/graph-perday", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphPerdayHandler))
		http.HandleFunc(rootPath+"/list-perday", commonMiddleware(rateLimiter, ShowroomCGIlib.ListPerdayHandler))

		http.HandleFunc(rootPath+"/graph-perslot", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphPerslotHandler))
		http.HandleFunc(rootPath+"/list-perslot", commonMiddleware(rateLimiter, ShowroomCGIlib.ListPerslotHandler))
		http.HandleFunc(rootPath+"/graph-sum", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphSumHandler))
		http.HandleFunc(rootPath+"/graph-sum-data", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphSumDataHandler))
		http.HandleFunc(rootPath+"/graph-sum2", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphSum2Handler))
		http.HandleFunc(rootPath+"/graph-sum-data1", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphSumData1Handler))
		http.HandleFunc(rootPath+"/graph-sum-data2", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphSumData2Handler))

		http.HandleFunc(rootPath+"/add-event", commonMiddleware(rateLimiter, ShowroomCGIlib.AddEventHandler))
		http.HandleFunc(rootPath+"/edit-user", commonMiddleware(rateLimiter, ShowroomCGIlib.EditUserHandler))
		http.HandleFunc(rootPath+"/edit-cntrbpoints", commonMiddleware(rateLimiter, ShowroomCGIlib.EditCntrbPointsHandler))
		http.HandleFunc(rootPath+"/new-user", commonMiddleware(rateLimiter, ShowroomCGIlib.NewUserHandler))

		http.HandleFunc(rootPath+"/param-event", commonMiddleware(rateLimiter, ShowroomCGIlib.ParamEventHandler))
		http.HandleFunc(rootPath+"/param-eventc", commonMiddleware(rateLimiter, ShowroomCGIlib.ParamEventCHandler))

		http.HandleFunc(rootPath+"/new-event", commonMiddleware(rateLimiter, ShowroomCGIlib.NewEventHandler))

		http.HandleFunc(rootPath+"/param-global", commonMiddleware(rateLimiter, ShowroomCGIlib.ParamGlobalHandler))

		http.HandleFunc(rootPath+"/list-cntrb", commonMiddleware(rateLimiter, ShowroomCGIlib.ListCntrbHandler))
		http.HandleFunc(rootPath+"/list-cntrbex", commonMiddleware(rateLimiter, ShowroomCGIlib.ListCntrbExHandler))

		http.HandleFunc(rootPath+"/list-cntrbS", commonMiddleware(rateLimiter, ShowroomCGIlib.ListCntrbSHandler))

		http.HandleFunc(rootPath+"/list-cntrbH", commonMiddleware(rateLimiter, ShowroomCGIlib.ListCntrbHHandler))
		http.HandleFunc(rootPath+"/list-cntrbHEx", commonMiddleware(rateLimiter, ShowroomCGIlib.ListCntrbHExHandler))

		http.HandleFunc(rootPath+"/m-cntrbrank-listener", commonMiddleware(rateLimiter, ShowroomCGIlib.MonthlyCntrbRankOfListenerHandler))
		http.HandleFunc(rootPath+"/m-cntrbrank-Lg", commonMiddleware(rateLimiter, ShowroomCGIlib.MonthlyCntrbRankLgHandler))

		//	イベント参加ルームのリスナーの貢献ポイント履歴
		http.HandleFunc(rootPath+"/listener-cntrb-history", commonMiddleware(rateLimiter, ShowroomCGIlib.ListenerCntrbHistoryHandler))

		//	ルーム別リスナー貢献ポイントランキング
		http.HandleFunc(rootPath+"/room-cntrb-history", commonMiddleware(rateLimiter, ShowroomCGIlib.RoomCntrbHistoryHandler))

		http.HandleFunc(rootPath+"/fanlevel", commonMiddleware(rateLimiter, ShowroomCGIlib.FanLevelHandler))

		http.HandleFunc(rootPath+"/flranking", commonMiddleware(rateLimiter, ShowroomCGIlib.FlRankingHandler))

		http.HandleFunc(rootPath+"/currentdistrb", commonMiddleware(rateLimiter, ShowroomCGIlib.CurrentDistributorsHandler))

		http.HandleFunc(rootPath+"/currentevents", commonMiddleware(rateLimiter, ShowroomCGIlib.CurrentEventsHandler))

		http.HandleFunc(rootPath+"/eventroomlist", commonMiddleware(rateLimiter, ShowroomCGIlib.EventRoomListHandler))

		//	開催予定イベント一覧
		http.HandleFunc(rootPath+"/scheduledevents", commonMiddleware(rateLimiter, ShowroomCGIlib.ScheduledEventsHandler))

		//	開催予定イベント一覧（サーバーから取得）
		http.HandleFunc(rootPath+"/scheduledeventssvr", commonMiddleware(rateLimiter, ShowroomCGIlib.ScheduledEventsSvrHandler))

		//	終了イベント一覧
		http.HandleFunc(rootPath+"/closedevents", commonMiddleware(rateLimiter, ShowroomCGIlib.ClosedEventsHandler))
		http.HandleFunc(rootPath+"/oldevents", commonMiddleware(rateLimiter, ShowroomCGIlib.OldEventsHandler))
		http.HandleFunc(rootPath+"/contributors", commonMiddleware(rateLimiter, ShowroomCGIlib.ContributorsHandler))

		//	API: ルーム検索
		http.HandleFunc(rootPath+"/api/search-rooms", commonMiddleware(rateLimiter, ShowroomCGIlib.ApiSearchRoomsHandler))

		//	イベント最終結果
		http.HandleFunc(rootPath+"/closedeventroomlist", commonMiddleware(rateLimiter, ShowroomCGIlib.ClosedEventRoomListHandler))

		// http.HandleFunc(rootPath+"/apiroomstatus", commonMiddleware(rateLimiter, srhandler.HandlerApiRoomStatus))
		http.HandleFunc(rootPath+"/apiroomstatus", commonMiddleware(rateLimiter, ShowroomCGIlib.BadRequestHandler))

		//	ギフトランキングリスト
		http.HandleFunc(rootPath+"/listgs", commonMiddleware(rateLimiter, ShowroomCGIlib.ListGiftScoreHandler))

		//	ギフトランキンググラフ
		http.HandleFunc(rootPath+"/graphgs", commonMiddleware(rateLimiter, ShowroomCGIlib.GraphGiftScoreHandler))

		//	最強ファンランキングリスト
		http.HandleFunc(rootPath+"/listvgs", commonMiddleware(rateLimiter, ShowroomCGIlib.ListFanGiftScoreHandler))

		//	ギフトランキング貢献ランキングリスト
		http.HandleFunc(rootPath+"/listgsc", commonMiddleware(rateLimiter, ShowroomCGIlib.ListGiftScoreCntrbHandler))

		//	イベント獲得ポイント上位ルーム
		http.HandleFunc(rootPath+"/toproom", commonMiddleware(rateLimiter, ShowroomCGIlib.TopRoomHandler))

		//	SHOWランク上位配信者一覧表
		http.HandleFunc(rootPath+"/showrank", commonMiddleware(rateLimiter, ShowroomCGIlib.ShowRankHandler))

		//	アクセス統計
		http.HandleFunc(rootPath+"/accessstats", commonMiddleware(rateLimiter, ShowroomCGIlib.AccessStatsHandler))

		//	時刻単位アクセス統計
		http.HandleFunc(rootPath+"/accessstatshourly", commonMiddleware(rateLimiter, ShowroomCGIlib.AccessStatsHourlyHandler))

		//	アクセス集計表（ハンドラー、IPアドレス、ユーザーエージェント別）
		http.HandleFunc(rootPath+"/accesstable", commonMiddleware(rateLimiter, ShowroomCGIlib.AccessTableHandler))

		//	掲示板の書き込みと表示、同様の機能が HandlerTopForm()にもある。共通化すべき。
		http.HandleFunc(rootPath+"/disp-bbs", commonMiddleware(rateLimiter, ShowroomCGIlib.DispBbsHandler)) //  資料・実験
		http.HandleFunc(rootPath+"/experimental", commonMiddleware(rateLimiter, ShowroomCGIlib.ExperimentalHandler))

		//	ToDo管理
		http.HandleFunc(rootPath+"/list-todo", commonMiddleware(rateLimiter, ShowroomCGIlib.ListToDoHandler))
		http.HandleFunc(rootPath+"/insert-todo", commonMiddleware(rateLimiter, ShowroomCGIlib.InsertToDoHandler))
		http.HandleFunc(rootPath+"/edit-todo", commonMiddleware(rateLimiter, ShowroomCGIlib.EditToDoHandler))

		http.HandleFunc(rootPath+"/t008top", commonMiddleware(rateLimiter, srhandler.HandlerT008topForm)) //	http://....../t008top で呼び出される。
		http.HandleFunc(rootPath+"/t009top", commonMiddleware(rateLimiter, srhandler.HandlerT009topForm)) //	http://....../t009top で呼び出される。

		// ユーザー検索APIエンドポイント
		http.HandleFunc(rootPath+"/api/users", commonMiddleware(rateLimiter, ShowroomCGIlib.SearchUsersHandler))

		http.HandleFunc(rootPath+"/cgi-bin", commonMiddleware(rateLimiter, CgiBinHandler))
		http.HandleFunc(rootPath+"/cgi-bin/SC1", commonMiddleware(rateLimiter, CgiBinSc1Handler))
		http.HandleFunc(rootPath+"/cgi-bin/SC1/SRCGI", commonMiddleware(rateLimiter, CgiBinSc1SrcgiHandler))
		http.HandleFunc(rootPath+"/cgi-bin/SC1/SRCGI/top", commonMiddleware(rateLimiter, CgiBinSc1SrcgiTopHandler))
		http.HandleFunc(rootPath+"/cgi-bin/test/t009srapi/t008top", commonMiddleware(rateLimiter, T008topHandler))
		http.HandleFunc(rootPath+"/cgi-bin/test/t009srapi/t009top", commonMiddleware(rateLimiter, T009topHandler))
	} else {
		// Maintenance
		http.HandleFunc(rootPath+"/", commonMiddleware(rateLimiter, ShowroomCGIlib.DispBbsHandler))
	}

	// =============================================
	server := &http.Server{Addr: ":8080"}
	// =============================================

	/*
		if svconfig.WebServer == "None" {
			//	Webサーバーとして起動
			//	root権限のない（特権昇格ができない）ユーザーで起動した方が安全だと思います。
			//	その場合80や443のポートはlistenできないので、
			//	ポートを変えてルータやOSの設定でport redirectionするか
			//	ケーパビリティを設定してください。
			//	# setcap cap_net_bind_service=+ep ShowroomCGI
			//　（設置したあとこの操作を行うこと）
			//
			if svconfig.SSLcrt != "" {
				//	証明書があればSSLを使う
				log.Printf("           http.ListenAndServeTLS()\n")
				err := http.ListenAndServeTLS(":"+svconfig.HTTPport, svconfig.SSLcrt, svconfig.SSLkey, nil)
				if err != nil {
					log.Printf("%s\n", err.Error())
				}
			} else {
				log.Printf("           http.ListenAndServe()\n")
				err := http.ListenAndServe(":"+svconfig.HTTPport, nil)
				if err != nil {
					log.Printf("%s\n", err.Error())
				}
			}
		} else { //	CGIとして使う
			log.Printf("           cgi.Serve()\n")
			// CGIを起動
			//	使用するWebServerに応じて設置場所等適宜対応してください。
			cgi.Serve(nil)
		}
	*/

	// =============================================
	go func() {
		//	Webサーバーとして起動
		//	root権限のない（特権昇格ができない）ユーザーで起動した方が安全だと思います。
		//	その場合80や443のポートはlistenできないので、
		//	ポートを変えてルータやOSの設定でport redirectionするか
		//	ケーパビリティを設定してください。
		//	# setcap cap_net_bind_service=+ep ShowroomCGI
		//　（設置したあとこの操作を行うこと）
		//
		if svconfig.SSLcrt != "" {
			//	証明書があればSSLを使う
			log.Printf("           http.ListenAndServeTLS()\n")
			/*
				svconfig := ServerConfig{
					HTTPport: "8080",
					SSLcrt:   "path/to/your/cert.crt",
					SSLkey:   "path/to/your/key.key",
				}
			*/
			// HTTPサーバーを設定
			server := &http.Server{
				Addr:      ":" + svconfig.HTTPport,
				TLSConfig: &tls.Config{
					// 必要に応じてTLS設定を追加
				},
				Handler:      http.DefaultServeMux, // ここでハンドラを指定
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 90 * time.Second,
			}

			// サーバーをTLSで起動
			err := server.ListenAndServeTLS(svconfig.SSLcrt, svconfig.SSLkey)
			if err != nil {
				log.Println("Server error:", err)
			}
			// =============================================
			// err := http.ListenAndServeTLS(":"+svconfig.HTTPport, svconfig.SSLcrt, svconfig.SSLkey, nil)
			// if err != nil {
			// 	log.Printf("%s\n", err.Error())
			// }
		} else {
			log.Printf("           http.ListenAndServe()\n")
			server := &http.Server{Addr: ":" + svconfig.HTTPport}
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Println("Server error:", err)
			}
			// err := http.ListenAndServe(":"+svconfig.HTTPport, nil)
			// 	log.Printf("%s\n", err.Error())
			// }
		}
	}()

	// コンテキストがキャンセルされるのを待つ
	<-ctx.Done()

	// グレースフルシャットダウン
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Println("Server shutdown error:", err)
	}

	log.Println("Server gracefully stopped")
	// =============================================
}
func CgiBinHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("CgiBin called\n"))
}
func CgiBinSc1Handler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("CgiBinSc1 called\n"))
}
func CgiBinSc1SrcgiHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("CgiBinSc1Srcgi called\n"))
}
func CgiBinSc1SrcgiTopHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("<html>"))
	w.Write([]byte("このURLは以下に変更されました<br>"))
	w.Write([]byte("<a href=\"https://chouette2100.com/top\">https://chouette2100.com/top</a>"))
	w.Write([]byte("</html>"))
	//	w.Header().Set("Location", "https://chouette2100.com/top")
	//	w.Write([]byte("302"))
	//	http.Redirect(w, r, "https://chouette2100/top", 0)
}
func T008topHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("<html>"))
	w.Write([]byte("このURLは以下に変更されました<br>"))
	w.Write([]byte("<a href=\"https://chouette2100.com/t008top\">https://chouette2100.com/t008top</a>"))
	w.Write([]byte("</html>"))
}
func T009topHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("<html>"))
	w.Write([]byte("このURLは以下に変更されました<br>"))
	w.Write([]byte("<a href=\"https://chouette2100.com/t009top\">https://chouette2100.com/t009top</a>"))
	w.Write([]byte("</html>"))
}
