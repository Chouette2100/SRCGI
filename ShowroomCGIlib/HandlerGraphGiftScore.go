// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bytes"
	"fmt"
	//	"html"
	"log"

	//	"math/rand"
	//	"sort"
	"strconv"
	"strings"
	"time"

	//	"bufio"
	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"
	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"
	//	"github.com/PuerkitoBio/goquery"
	//	svg "github.com/ajstarks/svgo/float"
	//	"github.com/dustin/go-humanize"
	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"
	//	"github.com/Chouette2100/exsrapi/v2"
	//	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func HandlerGraphGiftScore(w http.ResponseWriter, req *http.Request) {

	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	//	eventid := req.FormValue("eventid")
	//	maxpoint, _ := strconv.Atoi(req.FormValue("maxpoint"))
	campaignid := req.FormValue("campaignid")
	giftid, _ := strconv.Atoi(req.FormValue("giftid"))
	target, _ := strconv.Atoi(req.FormValue("target"))
	nroom, _ := strconv.Atoi(req.FormValue("nroom"))
	if nroom == 0 {
		nroom = 10
	}
	smaxpoint := req.FormValue("maxpoint")
	maxpoint, _ := strconv.Atoi(smaxpoint)
	if maxpoint < 10000 {
		maxpoint = 0
		smaxpoint = "0"
	}
	sgscale := req.FormValue("gscale")
	if sgscale == "" || sgscale == "0" {
		sgscale = "100"
	}
	gscale, _ := strconv.Atoi(sgscale)
	/*
		gschk100 := ""
		gschk90 := ""
		gschk80 := ""
		gschk70 := ""
		switch sgscale {
		case "100":
			gschk100 = "checked"
		case "90":
			gschk90 = "checked"
		case "80":
			gschk80 = "checked"
		case "70":
			gschk70 = "checked"
		default:
			gschk100 = "checked"
		}
	*/

	intrf, err := srdblib.Dbmap.Get(srdblib.Campaign{}, campaignid)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Get(srdblib.Campaign{}, %s) err = %w", campaignid, err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	campaign := intrf.(*srdblib.Campaign)

	intrf, err = srdblib.Dbmap.Get(srdblib.GiftRanking{}, campaignid, giftid)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Get(srdblib.GiftRanking{}, %d) err = %w", giftid, err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	giftranking := intrf.(*srdblib.GiftRanking)

	filename, err := GraphGiftScore(campaign, giftranking, nroom, maxpoint, target, gscale)
	if err != nil {
		err = fmt.Errorf("GraphGiftScore(): Error: %w", err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	if Serverconfig.WebServer == "nginxSakura" {
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	} else if Serverconfig.WebServer == "Apache2Ubuntu" {
		filename = "/public/" + filename
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-gs.gtpl"))

	type Ggsheader struct {
		Filename     string
		Campaignid   string
		Campaignname string
		Url          string
		Grid         int
		Grname       string
		Period       string
		Maxpoint     int
		Gscale       int
	}
	ggsheader := Ggsheader{
		Filename:     filename,
		Campaignid:   campaign.Campaignid,
		Campaignname: campaign.Campaignname,
		Url:          campaign.Url,
		Grid:         giftranking.Grid,
		Grname:       giftranking.Grname,
		Period:       giftranking.Startedat.Format("2006-01-02 15:04") + " 〜 " + giftranking.Endedat.Format("2006-01-02 15:04"),
		Maxpoint:     maxpoint,
		Gscale:       gscale,
	}

	// マップを展開してテンプレートを出力する
	if err := tpl.ExecuteTemplate(w, "graph-gs.gtpl", ggsheader); err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

//	var Nfseq int

func GraphGiftScore(
	campaign *srdblib.Campaign,
	giftranking *srdblib.GiftRanking,
	nroom int,
	maxpoint int,
	target int,
	gscale int,
) (
	filename string,
	err error,
) {

	if Serverconfig.WebServer == "None" {
		filename = fmt.Sprintf("%03d.svg", Nfseq)
		Nfseq = (Nfseq + 1) % 1000
	} else {
		filename = fmt.Sprintf("%03d.svg", os.Getpid()%1000)
		//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//	filename = fmt.Sprintf("%0d.svg", r.Intn(100))
	}

	maxscore := 0

	xydata := make([]Xydata, 0)
	IDlist := make([]int, 0)

	giftscore := make([]srdblib.GiftScore, 0)
	sqlst := "select userno from giftscore "
	sqlst += " where giftid = ? "
	sqlst += " and ts = (select max(ts) from giftscore where giftid = ?) order by score desc limit ? "
	_, err = srdblib.Dbmap.Select(&giftscore, sqlst, giftranking.Grid, giftranking.Grid, nroom)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(&giftscore, %d)(1) err = %w", giftranking.Grid, err)
		return
	}

	ustartedat := Jtruncate(giftranking.Startedat).Unix()
	for _, v := range giftscore {
		IDlist = append(IDlist, v.Userno)

		giftscore = make([]srdblib.GiftScore, 0)
		sqlst = "select ts, score from giftscore where giftid = ? and userno = ? and ts < ? order by ts "
		_, err = srdblib.Dbmap.Select(&giftscore, sqlst, giftranking.Grid, v.Userno, giftranking.Endedat.Add(2*time.Hour))
		if err != nil {
			err = fmt.Errorf("srdblib.Dbmap.Select(&giftscore, %d, %d,%s)(2) err = %w",
				giftranking.Grid, v.Userno, giftranking.Endedat.Add(2*time.Hour).Format("2006-01-02 15:04"), err)
			return
		}

		xy := Xydata{}
		xy.X = make([]float64, 0)
		xy.Y = make([]float64, 0)
		for _, v := range giftscore {
			xy.X = append(xy.X, float64(v.Ts.Unix()-ustartedat)/60/60/24)
			xy.Y = append(xy.Y, float64(v.Score))
			if maxscore < v.Score {
				maxscore = v.Score
			}
		}

		xydata = append(xydata, xy)

	}

	err = DrawLineGraph(
		filename,              //	（パスなし）ファイル名　ex. 000.svg
		campaign.Campaignname, //	ex.	グラフタイトル
		giftranking.Grname,    //	ex. イベント名
		giftranking.Startedat.Format("2006-01-02 15:04")+" 〜 "+giftranking.Endedat.Format("2006-01-02 15:04"), //	ex. 開催期間
		maxscore,              //	データの最大値
		maxpoint,              //	y軸方向グラフ表示範囲を制限する
		target,                //	目標ポイント
		giftranking.Startedat, //	イベント開始時刻 time.Time
		giftranking.Endedat,   //	イベント終了時刻 time.Time
		0, // グラフ描画に使用するカラーマップ
		1.1,                   //	データ間隔がこの時間を超えたら接続しない(day)
		IDlist,
		&xydata,
	)

	if err != nil {
		err = fmt.Errorf(" err = %w", err)
		return "", err
	}

	return
}
