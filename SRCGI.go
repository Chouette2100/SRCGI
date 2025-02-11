package main

import (
	//	"fmt"
	//	"io"
	"log"

	//	"strconv"
	"time"

	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"crypto/tls"

	//	"html/template"
	"net/http"

	// "net/http/cgi"

	//	"github.com/dustin/go-humanize"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
	"github.com/Chouette2100/srhandler"

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
	00AF01	掲示板機能について、HandlerWriteBbs()をHandlerDispBbs()に統合し、リモートアドレス、ユーザーエージェントを保存する。
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
*/

const version = "11CG00"

func NewLogfileName(logfile *os.File) {

	var err error

	//	毎日繰り返す
	for {

		tnow := time.Now()

		//	今日の午前9時
		today := tnow.Truncate(24 * time.Hour)

		//	今日の午前0時
		today = today.Add(-9 * time.Hour)
		//	test	today := tnow.Truncate(5 * time.Minute)

		//	次の日の午前0時
		nextday := today.AddDate(0, 0, 1)
		//	test	nextday := today.Add(5 * time.Minute)

		//	日付けが変わるまで待つ
		time.Sleep(nextday.Sub(tnow))

		//	ログファイルを閉じて新しいログファイルを作る
		logfile.Close()

		logfilename := version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + time.Now().Format("20060102") + ".txt"
		//	test	logfilename := version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + time.Now().Format("20060102-1504") + ".txt"

		logfile, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic("cannnot open logfile: " + logfilename + err.Error())
		}
		log.SetOutput(logfile)

		time.Sleep(1 * time.Second)
	}
}

// 共通の処理を行うミドルウェア
// =============================================
func commonMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 共通の処理をここで行う
        log.Println("Common processing")

        // 次のハンドラーを呼び出す
        next(w, r)
    }
}
// =============================================


// =============================================
// サーバー構成
type ServerConfig struct {
	HTTPport string
	SSLcrt   string
	SSLkey   string
}

// =============================================

// 入力内容の確認画面
func main() {

	logfilename := version + "_" + ShowroomCGIlib.Version + "_" + srdblib.Version + "_" + time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	// log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	go NewLogfileName(logfile)

	ShowroomCGIlib.Chlog = make(chan *srdblib.Accesslog, 100)
	defer close(ShowroomCGIlib.Chlog)
	go ShowroomCGIlib.LogWorker()

	//	https://ssabcire.hatenablog.com/entry/2019/02/13/000722
	//	https://konboi.hatenablog.com/entry/2016/04/12/100903

	svconfig := ShowroomCGIlib.ServerConfig{}
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

	ShowroomCGIlib.LoadDenyIp("DenyIp.txt")
	//	log.Printf("DenyIp.txt = %v\n", ShowroomCGIlib.DenyIpList)

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
	srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")

	srdblib.Dbmap.AddTableWithName(srdblib.GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
	srdblib.Dbmap.AddTableWithName(srdblib.ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	srdblib.Dbmap.AddTableWithName(srdblib.Viewer{}, "viewer").SetKeys(false, "Viewerid")
	srdblib.Dbmap.AddTableWithName(srdblib.ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
	srdblib.Dbmap.AddTableWithName(ShowroomCGIlib.Contribution{}, "contribution").SetKeys(false, "Ieventid", "Roomid", "Viewerid")

	srdblib.Dbmap.AddTableWithName(srdblib.Campaign{}, "campaign").SetKeys(false, "Campaignid")
	srdblib.Dbmap.AddTableWithName(srdblib.GiftRanking{}, "giftranking").SetKeys(false, "Campaignid", "Grid")
	srdblib.Dbmap.AddTableWithName(srdblib.Accesslog{}, "accesslog").SetKeys(false, "Ts", "Eventid")

	if svconfig.WebServer == "None" {
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
		/*	Maintenance */

		http.HandleFunc(rootPath+"/top", commonMiddleware(ShowroomCGIlib.HandlerTopForm))

		http.HandleFunc(rootPath+"/list-level", commonMiddleware(ShowroomCGIlib.HandlerListLevel))

		http.HandleFunc(rootPath+"/list-last", commonMiddleware(ShowroomCGIlib.HandlerListLast))

		http.HandleFunc(rootPath+"/dl-all-points", commonMiddleware(ShowroomCGIlib.HandlerDlAllPoints))

		http.HandleFunc(rootPath+"/graph-total", commonMiddleware(ShowroomCGIlib.HandlerGraphTotal))
		http.HandleFunc(rootPath+"/csv-total", commonMiddleware(ShowroomCGIlib.HandlerCsvTotal))

		http.HandleFunc(rootPath+"/graph-dfr", commonMiddleware(ShowroomCGIlib.HandlerGraphDfr))

		http.HandleFunc(rootPath+"/graph-perday", commonMiddleware(ShowroomCGIlib.HandlerGraphPerday))
		http.HandleFunc(rootPath+"/list-perday", commonMiddleware(ShowroomCGIlib.HandlerListPerday))

		http.HandleFunc(rootPath+"/graph-perslot", commonMiddleware(ShowroomCGIlib.HandlerGraphPerslot))
		http.HandleFunc(rootPath+"/list-perslot", commonMiddleware(ShowroomCGIlib.HandlerListPerslot))
		http.HandleFunc(rootPath+"/graph-sum", commonMiddleware(ShowroomCGIlib.HandlerGraphSum))
		http.HandleFunc(rootPath+"/graph-sum-data", commonMiddleware(ShowroomCGIlib.HandlerGraphSumData))
		http.HandleFunc(rootPath+"/graph-sum2", commonMiddleware(ShowroomCGIlib.HandlerGraphSum2))
		http.HandleFunc(rootPath+"/graph-sum-data1", commonMiddleware(ShowroomCGIlib.HandlerGraphSumData1))
		http.HandleFunc(rootPath+"/graph-sum-data2", commonMiddleware(ShowroomCGIlib.HandlerGraphSumData2))

		http.HandleFunc(rootPath+"/add-event", commonMiddleware(ShowroomCGIlib.HandlerAddEvent))
		http.HandleFunc(rootPath+"/edit-user", commonMiddleware(ShowroomCGIlib.HandlerEditUser))
		http.HandleFunc(rootPath+"/new-user", commonMiddleware(ShowroomCGIlib.HandlerNewUser))

		http.HandleFunc(rootPath+"/param-event", commonMiddleware(ShowroomCGIlib.HandlerParamEvent))
		http.HandleFunc(rootPath+"/param-eventc", commonMiddleware(ShowroomCGIlib.HandlerParamEventC))

		http.HandleFunc(rootPath+"/new-event", commonMiddleware(ShowroomCGIlib.HandlerNewEvent))

		http.HandleFunc(rootPath+"/param-global", commonMiddleware(ShowroomCGIlib.HandlerParamGlobal))

		http.HandleFunc(rootPath+"/list-cntrb", commonMiddleware(ShowroomCGIlib.HandlerListCntrb))

		http.HandleFunc(rootPath+"/list-cntrbS", commonMiddleware(ShowroomCGIlib.HandlerListCntrbS))

		http.HandleFunc(rootPath+"/list-cntrbH", commonMiddleware(ShowroomCGIlib.HandlerListCntrbH))

		http.HandleFunc(rootPath+"/fanlevel", commonMiddleware(ShowroomCGIlib.HandlerFanLevel))

		http.HandleFunc(rootPath+"/flranking", commonMiddleware(ShowroomCGIlib.HandlerFlRanking))

		http.HandleFunc(rootPath+"/currentdistrb", commonMiddleware(ShowroomCGIlib.HandlerCurrentDistributors))

		http.HandleFunc(rootPath+"/currentevents", commonMiddleware(ShowroomCGIlib.HandlerCurrentEvents))

		http.HandleFunc(rootPath+"/eventroomlist", commonMiddleware(ShowroomCGIlib.HandlerEventRoomList))

		//	開催予定イベント一覧
		http.HandleFunc(rootPath+"/scheduledevents", commonMiddleware(ShowroomCGIlib.HandlerScheduledEvents))

		//	開催予定イベント一覧（サーバーから取得）
		http.HandleFunc(rootPath+"/scheduledeventssvr", commonMiddleware(ShowroomCGIlib.HandlerScheduledEventsSvr))

		//	終了イベント一覧
		http.HandleFunc(rootPath+"/closedevents", commonMiddleware(ShowroomCGIlib.HandlerClosedEvents))
		http.HandleFunc(rootPath+"/oldevents", commonMiddleware(ShowroomCGIlib.HandlerOldEvents))
		http.HandleFunc(rootPath+"/contributors", commonMiddleware(ShowroomCGIlib.HandlerContributors))

		//	イベント最終結果
		http.HandleFunc(rootPath+"/closedeventroomlist", commonMiddleware(ShowroomCGIlib.HandlerClosedEventRoomList))

		http.HandleFunc(rootPath+"/apiroomstatus", commonMiddleware(srhandler.HandlerApiRoomStatus))

		//	ギフトランキングリスト
		http.HandleFunc(rootPath+"/listgs", commonMiddleware(ShowroomCGIlib.HandlerListGiftScore))

		//	ギフトランキンググラフ
		http.HandleFunc(rootPath+"/graphgs", commonMiddleware(ShowroomCGIlib.HandlerGraphGiftScore))

		//	最強ファンランキングリスト
		http.HandleFunc(rootPath+"/listvgs", commonMiddleware(ShowroomCGIlib.HandlerListFanGiftScore))

		//	ギフトランキング貢献ランキングリスト
		http.HandleFunc(rootPath+"/listgsc", commonMiddleware(ShowroomCGIlib.HandlerListGiftScoreCntrb))

		//	イベント獲得ポイント上位ルーム
		http.HandleFunc(rootPath+"/toproom", commonMiddleware(ShowroomCGIlib.HandlerTopRoom))

		//	SHOWランク上位配信者一覧表
		http.HandleFunc(rootPath+"/showrank", commonMiddleware(ShowroomCGIlib.HandlerShowRank))

		//	掲示板の書き込みと表示、同様の機能が HandlerTopForm()にもある。共通化すべき。
		http.HandleFunc(rootPath+"/disp-bbs", commonMiddleware(ShowroomCGIlib.HandlerDispBbs))

		http.HandleFunc(rootPath+"/t008top", commonMiddleware(srhandler.HandlerT008topForm)) //	http://....../t008top で呼び出される。
		http.HandleFunc(rootPath+"/t009top", commonMiddleware(srhandler.HandlerT009topForm)) //	http://....../t009top で呼び出される。

		http.HandleFunc(rootPath+"/cgi-bin", commonMiddleware(HandlerCgiBin))
		http.HandleFunc(rootPath+"/cgi-bin/SC1", commonMiddleware(HandlerCgiBinSc1))
		http.HandleFunc(rootPath+"/cgi-bin/SC1/SRCGI", commonMiddleware(HandlerCgiBinSc1Srcgi))
		http.HandleFunc(rootPath+"/cgi-bin/SC1/SRCGI/top", commonMiddleware(HandlerCgiBinSc1SrcgiTop))
		http.HandleFunc(rootPath+"/cgi-bin/test/t009srapi/t008top", commonMiddleware(Handlert008top))
		http.HandleFunc(rootPath+"/cgi-bin/test/t009srapi/t009top", commonMiddleware(Handlert009top))

		/* Maintenance ここまで */
	} else {

		/* Maintenance */
		http.HandleFunc(rootPath+"/top", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-level", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-last", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/graph-total", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/csv-total", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/graph-dfr", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/graph-perday", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-perday", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/graph-perslot", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-perslot", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/add-event", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/edit-user", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/new-user", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/param-event", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/param-eventc", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/new-event", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/param-global", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-cntrb", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-cntrbS", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/list-cntrbH", ShowroomCGIlib.HandlerListCntrbH)
		http.HandleFunc(rootPath+"/fanlevel", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/flranking", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/currentdistrb", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/currentevents", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/eventroomlist", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/scheduledevents", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/scheduledeventssvr", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/closedevents", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/closedeventroomlist", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/apiroomstatus", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/toproom", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/showrank", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/disp-bbs", ShowroomCGIlib.HandlerDispBbs)
		http.HandleFunc(rootPath+"/t008top", ShowroomCGIlib.HandlerDispBbs) //	http://....../t008top で呼び出される。
		http.HandleFunc(rootPath+"/t009top", ShowroomCGIlib.HandlerDispBbs) //	http://....../t009top で呼び出される。

		/*	Maintenance ここまで	*/
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
			svconfig := ServerConfig{
				HTTPport: "8080",
				SSLcrt:   "path/to/your/cert.crt",
				SSLkey:   "path/to/your/key.key",
			}
			// HTTPサーバーを設定
			server := &http.Server{
				Addr:      ":" + svconfig.HTTPport,
				TLSConfig: &tls.Config{
					// 必要に応じてTLS設定を追加
				},
				Handler:      http.DefaultServeMux, // ここでハンドラを指定
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
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
func HandlerCgiBin(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("CgiBin called\n"))
}
func HandlerCgiBinSc1(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("CgiBinSc1 called\n"))
}
func HandlerCgiBinSc1Srcgi(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := ShowroomCGIlib.GetUserInf(r)
	if !isallow {
		w.Write([]byte("Access Denied\n"))
		return
	}
	w.Write([]byte("CgiBinSc1Srcgi called\n"))
}
func HandlerCgiBinSc1SrcgiTop(w http.ResponseWriter, r *http.Request) {

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
func Handlert008top(w http.ResponseWriter, r *http.Request) {

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
func Handlert009top(w http.ResponseWriter, r *http.Request) {

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
