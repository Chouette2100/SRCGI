// Copyright © 2024 chouette.21.00@gmail.com
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
	//	"strconv"
	"os"
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

type Point struct {
	Pnt  int
	Spnt string
	Tpnt string
}

type PointRecord struct {
	Day       string
	Tday      time.Time
	Pointlist []Point
}

type PointPerDay struct {
	Eventid         string
	Eventname       string
	Period          string
	Maxpoint        int
	Gscale          int
	Usernolist      []int
	Longnamelist    []LongName
	Pointrecordlist []PointRecord
}

func GraphPerdayHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-perday.gtpl"))

	eventid := r.FormValue("eventid")
	Event_inf.Event_ID = eventid
	_, sts := SelectEventInfAndRoomList()

	if sts != 0 {
		log.Printf("MakePointPerDay() status of SelectEventInfAndRoomList() =%d\n", sts)
		return
	}

	log.Printf("      called. eventid=%s\n", eventid)

	ppointperday, _ := MakePointPerDay(Event_inf)

	filename, _ := GraphPerDay(eventid, ppointperday)
	if Serverconfig.WebServer == "nginxSakura" {
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	} else if Serverconfig.WebServer == "Apache2Ubuntu" {
		filename = "/public/" + filename
	}

	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
		"maxpoint": fmt.Sprintf("%d", Event_inf.Maxpoint),
		"gscale":   fmt.Sprintf("%d", Event_inf.Gscale),
	}

	if err := tpl.ExecuteTemplate(w, "graph-perday.gtpl", values); err != nil {
		log.Println(err)
	}

}

func ListPerdayHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/list-perday.gtpl"))

	eventid := r.FormValue("eventid")
	Event_inf.Event_ID = eventid
	_, sts := SelectEventInfAndRoomList()

	if sts != 0 {
		log.Printf("MakePointPerDay() status of SelectEventInfAndRoomList() =%d\n", sts)
		return
	}

	log.Printf("      eventid=%s\n", eventid)

	pointperday, _ := MakePointPerDay(Event_inf)

	if err := tpl.ExecuteTemplate(w, "list-perday.gtpl", *pointperday); err != nil {
		log.Println(err)
	}
}

// func MakePointPerDay(eventid string) (p_pointperday *PointPerDay, status int) {
func MakePointPerDay(Event_inf exsrapi.Event_Inf) (p_pointperday *PointPerDay, status int) {

	status = 0

	dstart := Event_inf.Start_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.Start_time.Hour()) * time.Hour)
	if Event_inf.Start_time.Hour()*60+Event_inf.Start_time.Minute() > Event_inf.Resethh*60+Event_inf.Resetmm {
		dstart = dstart.AddDate(0, 0, 1)
	}

	dend := Event_inf.End_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.End_time.Hour()) * time.Hour)
	if Event_inf.End_time.Hour()*60+Event_inf.End_time.Minute() > Event_inf.Resethh*60+Event_inf.Resetmm {
		dend = dend.AddDate(0, 0, 1)
	}

	days := int(dend.Sub(dstart).Hours() / 24)
	dstart = dstart.Add(time.Duration(Event_inf.Resethh*60+Event_inf.Resetmm) * time.Minute)

	log.Printf(" dstart=%s dend=%s days=%d\n", dstart.Format("2006/01/02 15:04:05"), dend.Format("2006/01/02 15:04:05"), days)

	var pointperday PointPerDay
	pointperday.Pointrecordlist = make([]PointRecord, days+1)
	pointperday.Eventname = Event_inf.Event_name
	pointperday.Eventid = Event_inf.Event_ID
	pointperday.Period = Event_inf.Period

	var roominfolist RoomInfoList
	_, _ = SelectEventRoomInfList(Event_inf.Event_ID, &roominfolist)
	log.Printf(" no of rooms. = %d\n", len(roominfolist))

	iu := 0 //	リスト作成の対象となるルームのインデックス

	for i := 0; i < len(roominfolist); i++ {

		if roominfolist[i].Graph != "Checked" {
			continue
		}
		log.Printf(" Room=%s Graph=%s\n", roominfolist[i].Longname, roominfolist[i].Graph)

		pointperday.Longnamelist = append(pointperday.Longnamelist, LongName{roominfolist[i].Longname})
		pointperday.Usernolist = append(pointperday.Usernolist, roominfolist[i].Userno)
		for k := 0; k < days+1; k++ {
			pointperday.Pointrecordlist[k].Day = dstart.AddDate(0, 0, k-1).Format("2006/01/02")
			pointperday.Pointrecordlist[k].Tday = dstart.AddDate(0, 0, k)
			if pointperday.Pointrecordlist[k].Tday.After(time.Now()) {
				pointperday.Pointrecordlist[k].Tday = time.Now().Truncate(time.Second)
			}
			if pointperday.Pointrecordlist[k].Tday.After(Event_inf.End_time) {
				pointperday.Pointrecordlist[k].Tday = Event_inf.End_time
			}
			pointperday.Pointrecordlist[k].Pointlist = append(pointperday.Pointrecordlist[k].Pointlist, Point{0, "", ""})
		}

		norow, tp, pp := SelectPointList(roominfolist[i].Userno, Event_inf.Event_ID)

		log.Printf(" no of point data=%d\n", norow)
		if norow == 0 {
			continue
		}

		d := dstart
		k := 0

		for ; ; k++ {
			if (*tp)[0].Before(d.AddDate(0, 0, k)) {
				break
			}
		}

		lastpoint := 0
		prvpoint := 0
		for j := 0; j < len(*tp); j++ {
			if (*tp)[j].After(d.AddDate(0, 0, k)) {
				//	log.Printf("i(room)=%d, j(time)=%d(%s), k(day)=%d\n", i, j, (*tp)[j].Format("01/02 15:04"), k)
				//	log.Printf("pointperday.Pointrecordlist[k].Pointlist=%v\n", pointperday.Pointrecordlist[k].Pointlist)
				if (*tp)[j].Sub(d.AddDate(0, 0, k)) < 30*time.Minute || j == 0 || (*pp)[j] == (*pp)[j-1] {
					pointperday.Pointrecordlist[k].Pointlist[iu].Pnt = lastpoint - prvpoint
					pointperday.Pointrecordlist[k].Pointlist[iu].Spnt = humanize.Comma(int64(lastpoint - prvpoint))
					pointperday.Pointrecordlist[k].Pointlist[iu].Tpnt = humanize.Comma(int64(lastpoint))
				} else {
					//	欠測が発生したと思われる場合
					pointperday.Pointrecordlist[k].Pointlist[iu].Pnt = lastpoint - prvpoint
					pointperday.Pointrecordlist[k].Pointlist[iu].Spnt = ""
				}
				prvpoint = lastpoint
				k++
			}
			lastpoint = (*pp)[j]
		}
		pointperday.Pointrecordlist[k].Pointlist[iu].Pnt = lastpoint - prvpoint
		pointperday.Pointrecordlist[k].Pointlist[iu].Spnt = humanize.Comma(int64(lastpoint - prvpoint))
		pointperday.Pointrecordlist[k].Pointlist[iu].Tpnt = humanize.Comma(int64(lastpoint))

		iu++
	}

	//	日々の獲得ポイントが空白の場合は次の日の獲得ポイントは無意味であるので空白とする。
	for k := days - 1; k >= 0; k-- {
		for i := 0; i < iu; i++ {
			if pointperday.Pointrecordlist[k].Pointlist[i].Spnt == "" {
				pointperday.Pointrecordlist[k+1].Pointlist[i].Spnt = ""
				pointperday.Pointrecordlist[k+1].Pointlist[i].Pnt = 0
			}
		}
	}

	p_pointperday = &pointperday

	return
}
func GraphPerDay(
	eventid string,
	pointperday *PointPerDay,
) (
	filename string,
	status int,
) {

	status = 0

	//	Event_inf, status = SelectEventInf(eventid)
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

	//	描画領域を決定する
	width := 3840.0
	height := 2160.0
	lwmargin := width / 24.0
	rwmargin := width / 6.0
	uhmargin := height / 7.5
	lhmargin := height / 15.0
	bstroke := width / 800.0

	vwidth := width - lwmargin - rwmargin
	vheight := height - uhmargin - lhmargin

	xorigin := lwmargin
	yorigin := height - lhmargin

	//	SVG出力ファイルを設定し、背景色を決める。
	//	filename = fmt.Sprintf("%0d.svg", os.Getpid()%100)
	if Serverconfig.WebServer == "None" {
		filename = fmt.Sprintf("%03d.svg", Nfseq)
		Nfseq = (Nfseq + 1) % 1000
	} else {
		filename = fmt.Sprintf("%03d.svg", os.Getpid()%1000)
	}
	file, err := os.OpenFile("public/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//	panic(err)
		return
	}

	bw := bufio.NewWriter(file)

	canvas := svg.New(bw)

	//	canvas := svg.New(os.Stdout)

	canvas.Start(width, height)
	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	maxpoint := 0
	for _, pointrecord := range (*pointperday).Pointrecordlist {
		for _, point := range pointrecord.Pointlist {
			if point.Pnt > maxpoint && point.Spnt != "" {
				maxpoint = point.Pnt
			}
		}
	}
	yupper, yscales, yscalel, _ := DetYaxScale(maxpoint)

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "配信日毎の獲得ポイント",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, Event_inf.Event_name,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), Event_inf.Period,
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

	xupper := Event_inf.Dperiod
	xscale := vwidth / float64(xupper)
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	//	一目盛の表示長さ
	dxl := 0.0
	if xscaled > 0 {
		dxl = 1.0 / float64(xscaled) * xscale
	} else {
		dxl = -1.0 * float64(xscaled) * xscale
	}

	tval := Event_inf.Start_time
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

	colorlist := make([]string, len((*pointperday).Usernolist))
	for i, userno := range (*pointperday).Usernolist {
		_, colorlist[i], _ = SelectUserColor(userno, Event_inf.Event_ID)
		longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(userno)
		if sts != 0 {
			longname = fmt.Sprintf("%d", userno)
		}
		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(i)
		//	canvas.Circle(xln+rwmargin/4.0, yln, 10.0, "stroke:"+colorlist[i]+";fill:"+colorlist[i])
		Mark(i, canvas, xln+rwmargin/4.0, yln, 10.0, colorlist[i])

		canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
	}

	//	日毎の獲得ポイントデータを描画する
	for _, pointrecord := range (*pointperday).Pointrecordlist {
		x := (float64(pointrecord.Tday.Unix())/60.0/60.0/24.0-Event_inf.Start_date)*xscale + xorigin
		for i, point := range pointrecord.Pointlist {
			if point.Spnt == "" {
				continue
			}
			y := float64(point.Pnt)*yscale + yorigin
			//	log.Printf("t=%7.3f, p=%8d, x=%7.2f, y=%7.2f\n",
			//		float64(pointrecord.Tday.Unix())/60.0/60.0/24.0-Event_inf.Start_date,
			//		point.Pnt, x, y)
			//	canvas.Circle(x, y, 10.0, "stroke:"+colorlist[i]+";fill:"+colorlist[i])
			Mark(i, canvas, x, y, 10.0, colorlist[i])
		}
	}

	canvas.End()

	bw.Flush()
	file.Close()

	return

}
