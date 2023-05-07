package main

import (
	//	"fmt"
	"log"

	//	"strconv"
	"time"

	"os"
	"runtime"

	//	"html/template"
	"net/http"

	"net/http/cgi"

	//	"github.com/dustin/go-humanize"

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

*/

const version = "0101J0"

// 入力内容の確認画面
func main() {

	logfilename := version + "_" + ShowroomCGIlib.Version + "_" + time.Now().Format("20060102") + ".txt"
	logfile, err := os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open logfile: " + logfilename + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	// log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//	https://ssabcire.hatenablog.com/entry/2019/02/13/000722
	//	https://konboi.hatenablog.com/entry/2016/04/12/100903

	/*
	crt := ""
	key := ""
	port := ""

	file, err := os.OpenFile("WebServer.txt", os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("Can't open WebServer.txt .\n")
		return
	} else {
		fmt.Fscanln(file, &ShowroomCGIlib.WebServer, &port)
		if ShowroomCGIlib.WebServer == "None" && port == "" {
			log.Printf("No port specified.")
			return
		}
		fmt.Fscanln(file, &crt)
		fmt.Fscanln(file, &key)
		fmt.Fscanln(file, &ShowroomCGIlib.Dbhost)
		fmt.Fscanln(file, &ShowroomCGIlib.Dbname)
		fmt.Fscanln(file, &ShowroomCGIlib.Dbuser)
		fmt.Fscanln(file, &ShowroomCGIlib.Dbpw)
		file.Close()
	}
	*/

	/*
	WebServer string `yaml:"WebServer"`
	HTTPport  string `yaml:"HTTPport"`
	SSLcrt    string `yaml:"SSLcrt"`
	SSLkey    string `yaml:"SSLkey"`
	Dbhost    string `yaml:"Dbhost"`
	Dbname    string `yaml:"Dbname"`
	Dbuser    string `yaml:"Dbuser"`
	Dbpw      string `yaml:"Dbpw"`

	*/
	ShowroomCGIlib.Dbconfig, err = ShowroomCGIlib.LoadConfig("ServerConfig.yml")
	if err != nil {
			log.Printf("err=%s.\n", err.Error())
			panic(err)
	}
	dbconfig := ShowroomCGIlib.Dbconfig

	log.Printf("%+v\n", *dbconfig)

	switch (*dbconfig).WebServer {
	case "nginxSakura":
		fallthrough
	case "Apache2Ubuntu":
		fallthrough
	case "None":
	default:
		log.Printf("Unknown WebServer = <%s> (must be nginxSakura, Apache2Ubuntu or None)\n", (*dbconfig).WebServer)
		return
	}

	ShowroomCGIlib.OS = runtime.GOOS
	rootPath := ""
	if (*dbconfig).WebServer != "None" {
		rootPath = os.Getenv("SCRIPT_NAME")
	}
	log.Printf("\n")
	log.Printf("\n")
	log.Printf("********** WevServer=<%s> port = <%s> OS = <%s> rootPath = <%s>\n", (*dbconfig).WebServer, (*dbconfig).HTTPport, ShowroomCGIlib.OS, rootPath)
	log.Printf("********** crt=<%s> key = <%s>\n", dbconfig.SSLcrt, dbconfig.SSLkey)
	log.Printf("********** Dbhost=<%s> Dbname = <%s> Dbuser = <%s> Dbpw = <%s>\n", (*dbconfig).Dbhost, (*dbconfig).Dbname, (*dbconfig).Dbuser, (*dbconfig).Dbpw)

	status := ShowroomCGIlib.OpenDb()
	if status != 0 {
		log.Printf("Database error.\n")
		return
	}
	defer ShowroomCGIlib.Db.Close()

	if (*dbconfig).WebServer == "None" {
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

	http.HandleFunc(rootPath+"/top", ShowroomCGIlib.HandlerTopForm)

	http.HandleFunc(rootPath+"/list-level", ShowroomCGIlib.HandlerListLevel)

	http.HandleFunc(rootPath+"/list-last", ShowroomCGIlib.HandlerListLast)

	http.HandleFunc(rootPath+"/graph-total", ShowroomCGIlib.HandlerGraphTotal)
	http.HandleFunc(rootPath+"/csv-total", ShowroomCGIlib.HandlerCsvTotal)

	http.HandleFunc(rootPath+"/graph-dfr", ShowroomCGIlib.HandlerGraphDfr)

	http.HandleFunc(rootPath+"/graph-perday", ShowroomCGIlib.HandlerGraphPerday)
	http.HandleFunc(rootPath+"/list-perday", ShowroomCGIlib.HandlerListPerday)

	http.HandleFunc(rootPath+"/graph-perslot", ShowroomCGIlib.HandlerGraphPerslot)
	http.HandleFunc(rootPath+"/list-perslot", ShowroomCGIlib.HandlerListPerslot)

	http.HandleFunc(rootPath+"/add-event", ShowroomCGIlib.HandlerAddEvent)
	http.HandleFunc(rootPath+"/edit-user", ShowroomCGIlib.HandlerEditUser)
	http.HandleFunc(rootPath+"/new-user", ShowroomCGIlib.HandlerNewUser)

	http.HandleFunc(rootPath+"/param-event", ShowroomCGIlib.HandlerParamEvent)
	http.HandleFunc(rootPath+"/param-eventc", ShowroomCGIlib.HandlerParamEventC)

	http.HandleFunc(rootPath+"/new-event", ShowroomCGIlib.HandlerNewEvent)

	http.HandleFunc(rootPath+"/param-global", ShowroomCGIlib.HandlerParamGlobal)

	http.HandleFunc(rootPath+"/list-cntrb", ShowroomCGIlib.HandlerListCntrb)

	http.HandleFunc(rootPath+"/list-cntrbS", ShowroomCGIlib.HandlerListCntrbS)

	http.HandleFunc(rootPath+"/list-cntrbH", ShowroomCGIlib.HandlerListCntrbH)

	http.HandleFunc(rootPath+"/fanlevel", ShowroomCGIlib.HandlerFanLevel)

	http.HandleFunc(rootPath+"/flranking", ShowroomCGIlib.HandlerFlRanking)

	if (*dbconfig).WebServer == "None" {
		//	Webサーバーとして起動
		//	root権限のない（特権昇格ができない）ユーザーで起動した方が安全だと思います。
		//	その場合80や443のポートはlistenできないので、
		//	ポートを変えてルータやOSの設定でport redirectionするか
		//	ケーパビリティを設定してください。
		//	# setcap cap_net_bind_service=+ep ShowroomCGI
		//　（設置したあとこの操作を行うこと）
		//
		if dbconfig.SSLcrt != "" {
			//	証明書があればSSLを使う
			log.Printf("           http.ListenAndServeTLS()\n")
			err := http.ListenAndServeTLS(":"+dbconfig.HTTPport, dbconfig.SSLcrt, dbconfig.SSLkey, nil)
			if err != nil {
				log.Printf("%s\n", err.Error())
			}
		} else {
			log.Printf("           http.ListenAndServe()\n")
			err := http.ListenAndServe(":"+dbconfig.HTTPport, nil)
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
}
