// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	"fmt"
	"html"
	"log"

	//	"math/rand"
	// "sort"
	"strconv"
	//	"strings"
	"time"
	//	"os"

	// "runtime"

	// "encoding/json"

	"html/template"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// 入力フォーム画面
func TopFormHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles(
	//		"templates/top.gtpl",
	//		"templates/top0.gtpl",
	//		"templates/bbs-2.gtpl",
	//		"templates/top1.gtpl",
	//		"templates/top2.gtpl",
	//	))

	eventid := r.FormValue("eventid")
	suserno := r.FormValue("userno")
	if suserno == "" {
		suserno = "0"
	}
	userno, _ := strconv.Atoi(suserno)
	log.Printf("      eventid=%s userno=%d\n", eventid, userno)

	if eventid == "" {

		// **********************************************
		var bbs BBS

		bbs.Cntlist = []int{1, 2, 3, 4, 5}
		bbs.Cntr = 9

		//      ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
		//	GetUserInf(req)

		/*
			bbs.Limit, _ = strconv.Atoi(r.FormValue("limit"))
			if bbs.Limit == 0 {
				bbs.Limit = 11
			}
		*/
		bbs.Limit = 11
		bbs.Offset, _ = strconv.Atoi(r.FormValue("offset"))

		action := r.FormValue("action")
		switch action {
		case "next":
			bbs.Offset += bbs.Limit - 1
		case "prev.":
			bbs.Offset -= bbs.Limit - 1
			if bbs.Offset < 0 {
				bbs.Offset = 0
			}
		case "再表示(top)":
			bbs.Offset = 0
		}

		from := r.FormValue("from")
		bbs.Manager = r.FormValue("manager")
		if bbs.Manager == "" {
			bbs.Manager = "black"
		}

		if from == "disp-bbs" {
			/*
				for i, v := range []string{"cnt0", "cnt1", "cnt2", "cnt3", "cnt4"} {
					cntv, _ := strconv.Atoi(r.FormValue(v))
					if cntv > 0 {
						bbs.Cntlist[i] = cntv
					} else {
						bbs.Cntlist[i] = -1
					}
				}
			*/
			bbs.Cntr, _ = strconv.Atoi(r.FormValue("cntr"))
		}

		//      テンプレートで使用する関数を定義する
		funcMap := template.FuncMap{
			"htmlEscapeString": func(s string) string { return html.EscapeString(s) },
			"FormatTime":       func(t time.Time, tfmt string) string { return t.Format(tfmt) },
			"CntToName": func(c int) string {
				cntname := []string{"不具合", "要望", "質問", "その他", "お知らせ", "すべて"}
				return cntname[c]
			},
			"Add": func(n int, m int) int { return n + m },
		}
		// テンプレートをパースする
		tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/top.gtpl", "templates/bbs-2.gtpl", "templates/top1.gtpl", "templates/top2.gtpl"))

		// ログを読み出してHTMLを生成 --- (*7)
		err := loadLogs(&bbs) // データを読み出す
		if err != nil {
			err = fmt.Errorf("loadLogs(): %w", err)
			log.Printf("showHandler(): %s\n", err.Error())
		}
		bbs.Nlog = len(bbs.Loglist)
		// **********************************************

		// マップを展開してテンプレートを出力する
		eventlist, _ := SelectLastEventList()
		if err := tpl.ExecuteTemplate(w, "top.gtpl", eventlist); err != nil {
			log.Println(err)
		}

		//	イベントでポイント比較の基準となる配信者（nobasis）のリストを取得する
		userlist, status := SelectUserList()
		if status == 0 {

			userlist[0].Userlongname = "ポイントの基準となる配信者が設定されていない"
			for i := 0; i < len(userlist); i++ {
				if userlist[i].Userno == userno {
					userlist[i].Selected = "Selected"
				} else {
					userlist[i].Selected = ""
				}
			}

			eventlist, _ = SelectEventList(userno)
			for i := 0; i < len(eventlist); i++ {
				if eventlist[i].EventID == eventid {
					eventlist[i].Selected = "Selected"
				} else {
					eventlist[i].Selected = ""
				}
			}
		}
		// マップを展開してテンプレートを出力する
		//		if err := tpl.ExecuteTemplate(w, "top0.gtpl", userlist); err != nil {
		if err := tpl.ExecuteTemplate(w, "bbs-2.gtpl", bbs); err != nil {
			log.Println(err)
		}
		if err := tpl.ExecuteTemplate(w, "top1.gtpl", eventlist); err != nil {
			log.Println(err)
		}
	} else {
		tpl := template.Must(template.ParseFiles(
			"templates/top2.gtpl",
		))

		//	eventinf, _ := SelectEventInf(eventid)
		//	srdblib.Tevent = "event"
		eventinf, err := srdblib.SelectFromEvent("event", eventid)
		if err != nil {
			//	DBの処理でエラーが発生した。
			return
		} else if eventinf == nil {
			//	指定した eventid のイベントが存在しない。
			return
		}
		Event_inf = *eventinf

		if err := tpl.ExecuteTemplate(w, "top2.gtpl", eventinf); err != nil {
			log.Println(err)
		}
	}

}
func SelectLastEventList() (eventlist []Event, status int) {

	var err error
	var stmt *sql.Stmt
	var rows *sql.Rows

	//	sql := "select eventid, event_name, period, starttime, endtime, nobasis, longname from event join user "
	sql := "select eventid, event_name, period, starttime, endtime, nobasis, modmin, modsec, longname, maxpoint from event join user "
	sql += " where nobasis = userno and endtime IS not null order by endtime desc "
	stmt, err = srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query()
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var event Event
	i := 0
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.EventName, &event.Period, &event.Starttime, &event.Endtime, &event.Pntbasis, &event.Modmin, &event.Modsec, &event.Pbname, &event.Maxpoint)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		event.Gscale = event.Maxpoint % 100
		event.Maxpoint = event.Maxpoint - event.Gscale
		eventlist = append(eventlist, event)
		i++
		if i == Serverconfig.NoEvent {
			break
		}
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	tnow := time.Now()
	for i = 0; i < len(eventlist); i++ {
		eventlist[i].S_start = eventlist[i].Starttime.Format("2006-01-02 15:04")
		eventlist[i].S_end = eventlist[i].Endtime.Format("2006-01-02 15:04")

		if eventlist[i].Starttime.After(tnow) {
			eventlist[i].Status = "これから開催"
		} else if eventlist[i].Endtime.Before(tnow) {
			eventlist[i].Status = "終了"
		} else {
			eventlist[i].Status = "開催中"
		}

	}

	return

}
func SelectEventList(userno int) (eventlist []Event, status int) {

	var err error
	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = srdblib.Db.Prepare("select eventid, event_name from event where endtime IS not null and nobasis = ? order by endtime desc")
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(userno)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var event Event
	i := 0
	for rows.Next() {
		err = rows.Scan(&event.EventID, &event.EventName)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		eventlist = append(eventlist, event)
		i++
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return

}
