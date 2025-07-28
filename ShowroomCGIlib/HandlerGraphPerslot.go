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

type PerSlot struct {
	Timestart time.Time
	Dstart    string
	Tstart    string
	Tend      string
	Point     string
	Ipoint    int
	Tpoint    string
}

type PerSlotInf struct {
	Eventname   string
	Eventid     string
	Period      string
	Roomname    string
	Roomid      int
	Perslotlist []PerSlot
}

func GraphPerslotHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-perslot.gtpl"))

	eventid := r.FormValue("eventid")
	log.Printf("      eventid=%s\n", eventid)

	_, perslotinflist, _ := MakePointPerSlot(eventid, 0)

	filename, eventinf, _ := GraphPerSlot(eventid, &perslotinflist)
	switch Serverconfig.WebServer {
	case "nginxSakura":
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	case "Apache2Ubuntu":
		filename = "/public/" + filename
	}

	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
		"maxpoint": fmt.Sprintf("%d", eventinf.Maxpoint),
		"gscale":   fmt.Sprintf("%d", eventinf.Gscale),
	}

	if err := tpl.ExecuteTemplate(w, "graph-perslot.gtpl", values); err != nil {
		log.Println(err)
	}

}

func ListPerslotHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/list-perslot1.gtpl",
		"templates/list-perslot2.gtpl",
	))

	roomid := 0
	sroomid := r.FormValue("roomid")
	if sroomid != "" {
		roomid, _ = strconv.Atoi(sroomid)
	}

	eventid := r.FormValue("eventid")
	//	Event_inf, _ = SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		return
	}
	// Event_inf = *eventinf

	log.Printf("      eventid=%s\n", eventid)

	if err := tpl.ExecuteTemplate(w, "list-perslot1.gtpl", eventinf); err != nil {
		log.Println(err)
	}

	_, perslotinflist, _ := MakePointPerSlot(eventid, roomid)

	if err := tpl.ExecuteTemplate(w, "list-perslot2.gtpl", perslotinflist); err != nil {
		log.Println(err)
	}

}

func MakePointPerSlot(eventid string, roomid int) (eventinf *exsrapi.Event_Inf, perslotinflist []PerSlotInf, status int) {

	var perslotinf PerSlotInf
	var err error

	status = 0

	eventinf = &exsrapi.Event_Inf{}
	// eventinf.Event_ID = eventid
	eventinf, err = srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		log.Printf("MakePointPerSlot() srdblib.SelectFromEvent() err=%s\n", err.Error())
		status = -1
		return
	}
	//	eventno, eventname, period := SelectEventNoAndName(eventid)
	eventname, period, _ := SelectEventNoAndName(eventid)

	var roominfolist RoomInfoList

	if roomid == 0 {
		// eventinf, _, sts := SelectEventRoomInfList(eventid, &roominfolist)
		_, _, sts := SelectEventRoomInfList(eventid, &roominfolist)
		if sts != 0 {
			log.Printf("status of SelectEventRoomInfList() =%d\n", sts)
			status = sts
			return
		}
	} else {
		roominfolist = make(RoomInfoList, 1)
		intf, _ := srdblib.Dbmap.Get(srdblib.User{}, roomid)
		roominfolist[0] = RoomInfo{
			Userno: roomid,
			Name:   intf.(*srdblib.User).User_name,
			Graph:  "Checked",
		}
	}

	for i := 0; i < len(roominfolist); i++ {

		if roominfolist[i].Graph != "Checked" {
			continue
		}

		var perslot PerSlot

		userid := roominfolist[i].Userno

		perslotinf.Eventname = eventname
		perslotinf.Eventid = eventid
		perslotinf.Period = period

		perslotinf.Roomname = roominfolist[i].Name
		perslotinf.Roomid = userid
		perslotinf.Perslotlist = make([]PerSlot, 0)

		norow, tp, pp := SelectPointList(userid, eventinf)

		if norow == 0 {
			continue
		}

		sameaslast := true
		plast := (*pp)[0]
		pprv := (*pp)[0]
		tdstart := ""
		tstart := time.Now().Truncate(time.Second)

		for i, t := range *tp {
			//	if (*pp)[i] != plast && sameaslast {
			if (*pp)[i] != plast {
				tstart = t
				/*
					if i != 0 {
						log.Printf("(1) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, (*tp)[i-1]=%s\n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"), (*tp)[i-1].Format("01/02 15:04"))
					} else {
						log.Printf("(1) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, \n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"))
					}
				*/
				if sameaslast {
					//	これまで変化しなかった獲得ポイントが変化し始めた
					pdstart := t.Add(time.Duration(-eventinf.Modmin) * time.Minute).Format("2006/01/02")
					if pdstart != tdstart {
						perslot.Dstart = pdstart
						tdstart = pdstart
					} else {
						perslot.Dstart = ""
					}
					perslot.Timestart = t.Add(time.Duration(-eventinf.Modmin) * time.Minute)
					//	perslot.Tstart = t.Add(time.Duration(-Event_inf.Modmin) * time.Minute).Format("15:04")
					if t.Sub((*tp)[i-1]) < 31*time.Minute {
						perslot.Tstart = perslot.Timestart.Format("15:04")
					} else {
						perslot.Tstart = "n/a"
					}
					//	perslot.Tstart = perslot.Timestart.Format("15:04")

					sameaslast = false
					//	} else if (*pp)[i] == plast && !sameaslast && (*tp)[i].Sub((*tp)[i-1]) > 11*time.Minute {
				}
			} else if (*pp)[i] == plast {
				//	if !sameaslast && (*tp)[i].Sub((*tp)[i-1]) > 16*time.Minute {
				if !sameaslast && t.Sub(tstart) > 11*time.Minute {
					//	if !sameaslast {
					/*
						if i != 0 {
							log.Printf("(2) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, (*tp)[i-1]=%s\n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"), (*tp)[i-1].Format("01/02 15:04"))
						} else {
							log.Printf("(2) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, \n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"))
						}
					*/
					if perslot.Tstart != "n/a" {
						perslot.Tend = (*tp)[i-1].Add(time.Duration(-eventinf.Modmin) * time.Minute).Format("15:04")
					} else {
						perslot.Tend = "n/a"
					}
					perslot.Ipoint = plast - pprv
					perslot.Point = humanize.Comma(int64(plast - pprv))
					perslot.Tpoint = humanize.Comma(int64(plast))
					sameaslast = true
					perslotinf.Perslotlist = append(perslotinf.Perslotlist, perslot)
					pprv = plast
				}
				//	sameaslast = true
			}
			/* else
			{
					if i != 0 {
						log.Printf("(3) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, (*tp)[i-1]=%s\n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"), (*tp)[i-1].Format("01/02 15:04"))
					} else {
						log.Printf("(3) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, \n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"))
					}
			}
			*/
			plast = (*pp)[i]
		}

		if len(perslotinf.Perslotlist) != 0 {
			perslotinflist = append(perslotinflist, perslotinf)
		}

		UpdatePointsSetQstatus(eventid, userid, perslot.Tstart, perslot.Tend, perslot.Point)

	}

	return
}

var Nfseq int

func GraphPerSlot(
	eventid string,
	perslotinflist *[]PerSlotInf,
) (
	filename string,
	eventinf *exsrapi.Event_Inf,
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
	// Event_inf = *eventinf

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

	/*
		canvas.Circle(width/2, height/2, 100)
		canvas.Text(width/2, height/2, "ポケGO", "text-anchor:middle;font-size:30px;fill:white;")
	*/

	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	maxpoint := 0
	for _, perslotinf := range *perslotinflist {
		for _, perslot := range perslotinf.Perslotlist {
			if perslot.Ipoint > maxpoint {
				maxpoint = perslot.Ipoint
			}
		}
	}
	yupper, yscales, yscalel, _ := DetYaxScale(maxpoint)

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "配信枠毎の獲得ポイント",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, eventinf.Event_name,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), eventinf.Period,
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

	xupper := eventinf.Dperiod
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

	tval := eventinf.Start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%xscaled == 0 {
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

	//	配信枠毎の獲得ポイントデータを描画する

	for j, perslotinf := range *perslotinflist {
		_, cvalue, _ := SelectUserColor(perslotinf.Roomid, eventinf)
		for _, perslot := range perslotinf.Perslotlist {
			y := float64(perslot.Ipoint)*yscale + yorigin
			x := (float64(perslot.Timestart.Unix())/60.0/60.0/24.0-eventinf.Start_date)*xscale + xorigin
			log.Printf("t=%7.3f, p=%8d, x=%7.2f, y=%7.2f\n",
				float64(perslot.Timestart.Unix())/60.0/60.0/24.0-eventinf.Start_date,
				perslot.Ipoint, x, y)
			//	canvas.Circle(x, y, 10.0, "stroke:"+cvalue+";fill:"+cvalue)
			Mark(j, canvas, x, y, 10.0, cvalue)
		}
		longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(perslotinf.Roomid)
		if sts != 0 {
			longname = fmt.Sprintf("%d", perslotinf.Roomid)
		}
		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(j)
		//	canvas.Circle(xln+rwmargin/4.0, yln, 10.0, "stroke:"+cvalue+";fill:"+cvalue)
		Mark(j, canvas, xln+rwmargin/4.0, yln, 10.0, cvalue)
		canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
	}

	canvas.End()

	bw.Flush()
	file.Close()

	return

}

func UpdatePointsSetQstatus(
	eventid string,
	userno int,
	tstart string,
	tend string,
	point string,
) (status int) {
	status = 0

	// log.Printf("  *** UpdatePointsSetQstatus() *** eventid=%s userno=%d\n", eventid, userno)

	nrow := 0
	//	err := Db.QueryRow("select count(*) from points where eventid = ? and user_id = ? and pstatus = 'Conf.'", eventid, userno).Scan(&nrow)
	sql := "select count(*) from points where eventid = ? and user_id = ? and ( pstatus = 'Conf.' or pstatus = 'Prov.' )"
	err := srdblib.Db.QueryRow(sql, eventid, userno).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	if nrow != 1 {
		return
	}

	// log.Printf("  *** UpdatePointsSetQstatus() Update!\n")

	sql = "update points set qstatus =?,"
	sql += "qtime=? "
	//	sql += "where user_id=? and eventid = ? and pstatus = 'Conf.'"
	sql += "where user_id=? and eventid = ? and ( pstatus = 'Conf.' or pstatus = 'Prov.' )"
	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdatePointsSetQstatus() Update/Prepare err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(point, tstart+"--"+tend, userno, eventid)

	if err != nil {
		log.Printf("error(UpdatePointsSetQstatus() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return
}
