// Copyright © 2024-2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	"bufio"
	// "bytes"
	"fmt"
	// "html"
	"log"

	//	"math/rand"
	// "sort"
	"os"
	"strconv"
	"strings"
	"time"

	// "runtime"

	// "encoding/json"

	"html/template"
	"net/http"

	// "database/sql"

	// _ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func GraphTotalHandler(w http.ResponseWriter, req *http.Request) {

	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	eventid := req.FormValue("eventid")
	//	maxpoint, _ := strconv.Atoi(req.FormValue("maxpoint"))

	//	描画するポイントの上限を指定する（0であれば制限しない）
	smaxpoint := req.FormValue("maxpoint")
	maxpoint, _ := strconv.Atoi(smaxpoint)
	if maxpoint < 10000 {
		maxpoint = 0
		smaxpoint = "0"
	}

	//　縮尺、ブラウザの横幅100%としたときの縮尺、0指定または指定なしは100%とみなす。
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
	//	resetcolor=on であればグラフの描画色をColorlist1またはColorlist2で初期化する
	resetcolor := req.FormValue("resetcolor")

	log.Printf("      eventid=%s maxpoint=%d(%s) resetcolor=[%s]\n", eventid, maxpoint, smaxpoint, resetcolor)

	if resetcolor == "on" {
		//	グラフの描画色を初期化する
		log.Printf("      Resetcolor(): eventid=%s\n", eventid)
		Resetcolor(eventid, -1)
	}

	//		グラフを作成する
	filename, _ := GraphTotalPoints(eventid, maxpoint, gscale)

	//	環境に応じてファイルのパスを決定する（Webサーバーとして起動した場合、パス指定がなければ/publicを参照する）
	switch Serverconfig.WebServer {
	case "nginxSakura":
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	case "Apache2Ubuntu":
		filename = "/public/" + filename
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-total.gtpl"))

	// テンプレートに出力する値をマップにセット
	/*
		values := map[string]string{
			"filename": req.FormValue("FileName"),
		}
	*/
	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
		"maxpoint": smaxpoint,
		"gscale":   sgscale,
	}

	// マップを展開してテンプレートを出力する
	if err := tpl.ExecuteTemplate(w, "graph-total.gtpl", values); err != nil {
		log.Println(err)
	}
}

//		グラフの色設定を初期化する。
//	 指定されたカラーマップを1位から順番にわりあてる。
func Resetcolor(
	eventid string, //	色設定初期化の対象となるイベントのイベントID
	cmap int, //	初期化に使うカラーマップ番号、 −1のときはEvent情報からカラーマップ番号を取得する
) error {

	if cmap < 0 {
		erow, err := srdblib.Dbmap.Get(srdblib.Event{}, eventid)
		if err != nil {
			err = fmt.Errorf("Resetcolor(): %w", err)
			return err
		}
		cmap = erow.(*srdblib.Event).Cmap
		if cmap == 0 {
			log.Printf("%+v\n", erow.(*srdblib.Event))
		}
	}
	log.Printf("      Resetcolor(): eventid=%s cmap=%d\n", eventid, cmap)
	clm := Colormaplist[cmap]
	lclm := len(clm)
	rows, err := srdblib.Dbmap.Select(srdblib.Eventuser{},
		"select "+clmlist["eventuser"]+" from eventuser where eventid = ? order by point desc", eventid)
	if err != nil {
		err = fmt.Errorf("Resetcolor(): %w", err)
		return err
	}
	for i, row := range rows {
		eu := row.(*srdblib.Eventuser)
		eu.Color = clm[i%lclm].Name
		_, err = srdblib.Dbmap.Update(eu)
		if err != nil {
			err = fmt.Errorf("Resetcolor(): %w", err)
			return err
		}
	}
	return nil
}
func GraphTotalPoints(eventid string, maxpoint int, gscale int) (filename string, status int) {

	var err error
	var eventinf *exsrapi.Event_Inf
	// eventinf = &exsrapi.Event_Inf{}

	status = 0

	// eventinf.Event_ID = eventid
	eventinf, err = srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		log.Printf("MakePointPerSlot() srdblib.SelectFromEvent() err=%s\n", err.Error())
		status = -1
		return
	}

	idandranklist, sts := SelectEventInfAndRoomList(eventinf)

	if sts != 0 {
		log.Printf("status of SelectEventInfAndRoomList() =%d\n", sts)
		status = sts
		return
	}

	_, period, _ := SelectEventNoAndName(eventid)

	eventinf.Maxpoint = maxpoint
	eventinf.Gscale = gscale
	UpdateEventInf(eventinf)

	if Serverconfig.WebServer == "None" {
		//	Webサーバーとして起動するときは、起動した直後を0とする連番（の下3桁）とする
		filename = fmt.Sprintf("%03d.svg", <-Chimgfn)
		Nfseq = (Nfseq + 1) % 1000
	} else {
		//	CGIのときはプロセスID（の下3桁）とする。
		filename = fmt.Sprintf("%03d.svg", os.Getpid()%1000)
		//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//	filename = fmt.Sprintf("%0d.svg", r.Intn(100))

	}
	//	グラフを描画する
	// GraphScore01(filename, eventinf, idandranklist, period, maxpoint)
	// GraphScore01(filename, eventinf, idandranklist, period, eventinf.MaxPoint)
	GraphScore01(filename, eventinf, idandranklist, period, eventinf.Maxpoint)

	/*
		fmt.Printf("Content-type:text/html\n\n")
		fmt.Printf("<!DOCTYPE html>\n")
		fmt.Printf("<html lang=\"ja\">\n")
		fmt.Printf("<head>\n")
		fmt.Printf("  <meta charset=\"UTF-8\">\n")
		fmt.Printf("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
		//	fmt.Printf("  <meta http-equiv=\"refresh\" content=\"30; URL=\">\n")
		fmt.Printf("  <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n")
		fmt.Printf("  <title></title>\n")
		fmt.Printf("</head>\n")
		fmt.Printf("<body>\n")
		fmt.Printf("<img src=\"test.svg\" alt=\"\" width=\"100%%\">")
		fmt.Printf("</body>\n")
		fmt.Printf("</html>\n")
	*/

	return
}

// グラフを描画する＝SVGを作成する
func GraphScore01(
	filename string,
	eventinf *exsrapi.Event_Inf, //	イベント情報
	idandranklist []IdAndRank,
	period string, maxpoint int,
) {

	//	描画領域を決定する
	width := 3840.0
	height := 2160.0
	lwmargin := width / 24.0
	rwmargin := width / 6.0
	uhmargin := height / 7.5
	lhmargin := height / 15.0
	bstroke := width / 800.0

	//  描画範囲の大きさ
	vwidth := width - lwmargin - rwmargin
	vheight := height - uhmargin - lhmargin

	xorigin := lwmargin
	yorigin := height - lhmargin

	//	SVG出力ファイルを設定し、背景色を決める。
	file, err := os.OpenFile("public/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//	panic(err)
		return
	}

	bw := bufio.NewWriter(file)

	canvas := svg.New(bw)

	//	canvas := svg.New(os.Stdout)

	canvas.Start(width, height)

	/*
		canvas.Circle(width/2, height/2, 100)
		canvas.Text(width/2, height/2, "ポケGO", "text-anchor:middle;font-size:30px;fill:white;")
	*/

	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	yupper := 0
	yscales := 0
	yscalel := 0

	if maxpoint != 0 {
		yupper, yscales, yscalel, _ = DetYaxScale(maxpoint - 1)
	} else if eventinf.Target > eventinf.MaxPoint {
		//	} else if Event_inf.Target > Event_inf.Maxpoint {
		yupper, yscales, yscalel, _ = DetYaxScale(eventinf.Target - 1)
	} else {
		yupper, yscales, yscalel, _ = DetYaxScale(eventinf.MaxPoint)
		//	yupper, yscales, yscalel, _ = DetYaxScale(Event_inf.Maxpoint)
	}

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "獲得ポイントの推移",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, eventinf.Event_name,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), period,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	//	y軸（ポイント軸）を描画する

	dyl := float64(yscales) * yscale
	value := int64(0)
	yl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%yscalel == 0 {
			wstr = 0.3

			canvas.Text(xorigin-bstroke*5.0, yorigin+yl+bstroke*2.5, humanize.Comma(value),
				"text-anchor:end;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		}
		canvas.Line(xorigin, yorigin+yl, xorigin+vwidth, yorigin+yl, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		yl += dyl
		if -yl > vheight+10 {
			break
		}
		value += int64(yscales)

	}

	//	------------------------------------------

	//	x軸（時間軸）を描画する

	//	x軸に表示する値の上限値
	xupper := eventinf.Dperiod
	//	x軸に表示する値と座標幅の比（表示値１が座標のいくらに相当するか？）
	xscale := vwidth / float64(xupper)
	// xscaled > 0 のとき　一目盛を1日の何分の一にするか？
	// xscaled < 0 のとき、一目盛を1日の何倍にするか？
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	//	一目盛の表示長さ
	dxl := 0.0
	if xscaled > 0 {
		dxl = 1.0 / float64(xscaled) * xscale
	} else {
		dxl = -1.0 * float64(xscaled) * xscale
	}
	tval := eventinf.Start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if xscaled > 0 && i%xscaled == 0 {
			//	xscaled > 0 のときは1日ごとに（00時に）表示を太くする
			wstr = 0.3
		}
		if xscaled > 0 && i%(xscaled*xscalet) == 0 || xscaled < 0 && i%xscalet == 0 {
			//	xscaled > 0 のときはxscalet日ごとに日付を表示する
			//	xscaled < 0 のときはxscaled * xscalet日ごとに日付を表示する
			if xscaled > -10 {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("1/2"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
			} else {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("06/01/02"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
			}

			if xscaled > 0 {
				tval = tval.AddDate(0, 0, xscalet)
			} else {
				tval = tval.AddDate(0, 0, -xscaled*xscalet)
			}
		}

		canvas.Line(xorigin+xl, yorigin, xorigin+xl, yorigin-vheight, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		xl += dxl
		if xl > vwidth+10 {
			break
		}
	}

	//	ターゲットラインを描画する
	if eventinf.Target != 0 {
		x1 := xorigin + (float64(eventinf.Start_time.Unix())/60.0/60.0/24.0-eventinf.Start_date)*xscale
		x2 := xorigin + (float64(eventinf.End_time.Unix())/60.0/60.0/24.0-eventinf.Start_date)*xscale
		y1 := yorigin
		y2 := yorigin + float64(eventinf.Target)*yscale

		log.Printf("Target (x1, y1) %10.2f,%10.2f (x2, y2) %10.2f,%10.2f xorgin, yorigin, vheight %10.2f, %10.2f %10.2f\n",
			x1, y1, x2, y2, xorigin, yorigin, vheight)

		if y2 < yorigin-vheight {
			x2 = (x2-xorigin)*vheight/(yorigin-y2) + xorigin
			y2 = yorigin - vheight
		}

		log.Printf("Target (x1, y1) %10.2f,%10.2f (x2, y2) %10.2f,%10.2f xorgin, yorigin, vheight %10.2f, %10.2f %10.2f\n",
			x1, y1, x2, y2, xorigin, yorigin, vheight)

		canvas.Line(x1, y1, x2, y2, `stroke="white" stroke-width="`+fmt.Sprintf("%.2f", bstroke*0.5)+`" stroke-dasharray="20,10"`)
	}

	//	獲得ポイントデータを描画する

	j := 0
	for _, iar := range idandranklist {

		_, cvalue, _ := SelectUserColor(iar.Userno, eventinf)

		x, y := SelectScoreList(eventinf, iar.Userno)
		maxp := 20

		//	no := len(*x)

		xo := make([]float64, maxp)
		yo := make([]float64, maxp)
		tl := 999.0
		yl := -1000000.0
		k := 0
		for i := 0; i < len(*x); i++ {
			//	fmt.Printf("(%7.1f,%10.1f)\n", (*x)[i], (*y)[i])
			xt := xorigin + (*x)[i]*xscale
			yt := yorigin + (*y)[i]*yscale
			//	fmt.Printf("(*x).[i]=%.3f tl=%.3f (*x)[i]-tl=%.3f\n", (*x)[i], tl, (*x)[i]-tl)
			if (*x)[i]-tl > 0.011 && (*y)[i]-yl > 1.0 {
				//	次のデータとの間に欠測があり、かつ欠測の前後でデータが同一でないときはその部分の描画は行わない。
				canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				xo[0] = xt
				yo[0] = yt
				tl = (*x)[i]
				yl = (*y)[i]
				k = 1
				if yt < yorigin-vheight {
					break
				}
				continue
			}

			if yt < yorigin-vheight {
				if k != 0 {
					xo[k] = (xt-xo[k-1])*(yo[k-1]-(yorigin-vheight))/(yo[k-1]-yt) + xo[k-1]
					yo[k] = yorigin - vheight
					k++
				}
				break
			} else {
				xo[k] = xt
				yo[k] = yt
			}

			tl = (*x)[i]
			yl = (*y)[i]
			k++
			if k == maxp {
				//	一定数のデータずつまとめて描画する。SVGファイルの可読性を高める。
				canvas.Polyline(xo, yo, "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				xo[0] = xt
				yo[0] = yt
				k = 1
			}
		}
		if k > 1 {
			canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		} else {
			canvas.Circle(xo[0], yo[0], bstroke*1.5, "fill=\""+cvalue+"\" stroke=\""+cvalue+"\"")
		}

		//	凡例
		xln := xorigin + vwidth + bstroke*10.0
		yln := yorigin - vheight + bstroke*10*float64(j)

		canvas.Line(xln, yln, xln+bstroke*20.0, yln, "stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		//	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, fmt.Sprintf("%d", IDlist[j]),
		//		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
		longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(idandranklist[j].Userno)
		if sts != 0 {
			longname = fmt.Sprintf("%d", idandranklist[j].Userno)
		}
		srank := "  -."
		if idandranklist[j].Rank > 0 {
			srank = fmt.Sprintf("%3d.", idandranklist[j].Rank)

		}
		canvas.Text(xln+bstroke*25.0, yln+bstroke*2.5, srank+longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		j++
	}

	//  ターゲットラインの凡例
	xln := xorigin + vwidth + bstroke*10.0
	yln := yorigin - vheight + bstroke*10*float64(j)

	canvas.Line(xln, yln, xln+bstroke*20.0, yln, `stroke="white" stroke-width="`+fmt.Sprintf("%.2f", bstroke*0.5)+`" stroke-dasharray="20,10"`)
	//	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, fmt.Sprintf("%d", IDlist[j]),
	//		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
	canvas.Text(xln+bstroke*25.0, yln+bstroke*2.5, "Target",
		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

	canvas.End()

	bw.Flush()
	file.Close()
}

func SelectScoreList(eventinf *exsrapi.Event_Inf, user_id int) (x *[]float64, y *[]float64) {

	stmt1, err := srdblib.Db.Prepare("SELECT count(*) FROM points where user_id = ? and eventid = ?")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt1.Close()

	var norow int
	err = stmt1.QueryRow(user_id, eventinf.Event_ID).Scan(&norow)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	//	fmt.Println(norow)

	//	----------------------------------------------------

	tu := make([]float64, norow)
	point := make([]float64, norow)

	//	----------------------------------------------------

	//	stmt2, err := Db.Prepare("select t.t, p.point from points p join timeacq t on t.idx = p.idx where user_id = ? and event_id = ? order by t.t")
	stmt2, err := srdblib.Db.Prepare("select ts, point from points where user_id = ? and eventid = ? order by ts")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt2.Close()

	rows, err := stmt2.Query(user_id, eventinf.Event_ID)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer rows.Close()
	i := 0
	var t time.Time
	for rows.Next() {
		err := rows.Scan(&t, &point[i])
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			//	status = -1
			return
		}
		if t.Before(eventinf.Start_time) {
			t = eventinf.Start_time
		}
		tu[i] = float64(t.Unix())/60.0/60.0/24.0 - eventinf.Start_date
		//	log.Printf("t=%v tu[%d]=%f\n", t, i, tu[i])
		i++

	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}

	x = &tu
	y = &point

	return
}
