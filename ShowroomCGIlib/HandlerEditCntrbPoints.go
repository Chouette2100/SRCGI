// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"bytes"
	"fmt"
	//	"html"
	"log"
	//	"os"
	//	"sort"
	"strconv"
	//	"strings"
	"time"

	//	"github.com/PuerkitoBio/goquery"

	//	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	// "database/sql"
	"github.com/Masterminds/sprig/v3"
	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// EditCntrbPointsHandler は枠別リスナー別貢献ポイントの取得対象ルームの編集を行います
func EditCntrbPointsHandler(w http.ResponseWriter, r *http.Request) {

	type CntrbPointsInfo struct {
		Eventid      string
		Eventname    string
		Period       string
		Maxpoint     int
		Gscale       int
		Roominfolist RoomInfoList
	}
	cpi := CntrbPointsInfo{}

	userlist := map[string]string{}
	for i := 0; i < 100; i++ {
		usr := r.FormValue(fmt.Sprintf("usr%d", i))
		// if usr == "" {
		// 	break
		// }
		if usr == "" {
			continue
		}
		userlist[usr] = "1"
	}

	fnc := r.FormValue("func")
	eventid := r.FormValue("eventid")
	log.Printf("      func=%s eventid=%s\n", fnc, eventid)

	userid := r.FormValue("userid")

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

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	switch fnc {
	case "newuser":
		//	新規配信者の追加があるとき

		roominf, status := GetRoomInfoAndPoint(eventid, userid, fmt.Sprintf("%d", eventinf.Nobasis))
		if status == 0 {
			tnow := time.Now().Truncate(time.Second)
			InsertIntoOrUpdateUser(client, tnow, eventid, roominf)
			InsertIntoEventUser(0, eventinf, roominf)
			UpdateEventuserSetPoint(eventid, roominf.ID, roominf.Point)
		} else {
			log.Printf("GetAndUpdateRoomInfoAndPoint() returned %d", status)
		}

	case "deleteuser":
		//	削除ボタンが押されたとき
	case "getAllCntrb":
		//	全員分の貢献ポイントを取得するボタンが押されたとき
		err = GetAllCntrb(eventid)
		if err != nil {
			log.Printf("GetAllCntrb() error=%s\n", err.Error())
			w.Write([]byte(fmt.Sprintf("GetAllCntrb() error=%s\n", err.Error())))
			return
		}
	default:
		// （更新ボタンが押された配信者がいたらそのデータを更新した上で）参加配信者のリストを表示する。
		// if userid != "" {
		// 	UpdateRoomInf(eventid, userid, longname, shortname, istarget, graph, color, iscntrbpoint)
		// }
		for usr := range userlist {
			err = UpdateIsCntrbPoints(eventid, usr, "1")
			if err != nil {
				log.Printf("UpdateEventuserSetCntrbpoints() error=%s\n", err.Error())
				w.Write([]byte(fmt.Sprintf("UpdateEventuserSetCntrbpoints() error=%s\n", err.Error())))
				return
			}
		}
	}

	//	log.Printf(" eventid=%s, userno=%s, longname=%s, shortname=%s, istarget=%s, graph=%s, color=%s\n",
	//		eventid, userno, longname, shortname, istarget, graph, color)

	var roominfolist RoomInfoList

	eventinf, eventname, _ := SelectEventRoomInfList(eventid, &cpi.Roominfolist)
	for i := 0; i < len(roominfolist); i++ {
		switch cpi.Roominfolist[i].Genre {
		case "Voice Actors & Anime":
			cpi.Roominfolist[i].Genre = "VA&A"
		case "Talent Model":
			cpi.Roominfolist[i].Genre = "Tl/Md"
		case "Comedians/Talk Show":
			cpi.Roominfolist[i].Genre = "Cm/TS"
		default:
		}
	}

	cpi.Eventid = eventid
	cpi.Eventname = eventname
	cpi.Period = eventinf.Period

	// テンプレートをパースする
	// funcMap := template.FuncMap{
	// 	"sub":   func(i, j int) int { return i - j },
	// 	"Comma": func(i int) string { return humanize.Comma(int64(i)) },
	// }
	funcMap := sprig.FuncMap()
	funcMap["Comma"] = func(i int) string { return humanize.Comma(int64(i)) }
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/edit-cntrbpoints.gtpl"))

	if err := tpl.ExecuteTemplate(w, "edit-cntrbpoints.gtpl", cpi); err != nil {
		err = fmt.Errorf("tpl.ExecuteTemplate() error=%s", err.Error())
		w.Write([]byte(err.Error()))
		log.Println(err)
	}
}

func UpdateIsCntrbPoints(eventid, suserno, iscntrbpoint string) (err error) {

	userno, _ := strconv.Atoi(suserno)

	if iscntrbpoint == "1" {
		iscntrbpoint = "Y"
	} else {
		iscntrbpoint = "N"
	}

	//	sql = "update eventuser set istarget=?, graph=?, color=? where eventno=? and userno=?"
	//	sql = "update eventuser set istarget=?, graph=?, color=?, iscntrbpoints=? where eventid=? and userno=?"
	sql := "update eventuser set iscntrbpoints = ? where eventid=? and userno=?"

	stmt, er := srdblib.Db.Prepare(sql)
	if er != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", er.Error())
		err = fmt.Errorf("UpdateRoomInf() error(Update/Prepare) err=%s", er.Error())
		return
	}
	defer stmt.Close()

	//	_, err = stmt.Exec(istarget, graph, color, iscntrbpoint, eventid, userno)
	_, err = stmt.Exec(iscntrbpoint, eventid, userno)
	if err != nil {
		log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		err = fmt.Errorf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s", err.Error())
	}

	return

}

/*
func GetAllCntrb(eventid string) (err error) {

	var event srdblib.Event
	var itfc interface{}
	itfc, err = srdblib.Dbmap.Get(&event, eventid)
	if err != nil {
		log.Printf("GetAllCntrb() Dbmap.Get() error=%s\n", err.Error())
		return
	}
	if itfc == nil {
		err = fmt.Errorf("eventid=%s not found", eventid)
		return
	}
	event = *itfc.(*srdblib.Event)

	if event.Eventid == "" {
		err = fmt.Errorf("eventid=%s not found", eventid)
		return
	}

	event.Cmode |= 0x01
	srdblib.Dbmap.Update(&event)

	srdblib.Dbmap.Exec("update eventuser set iscntrbpoints='Y' where eventid=?", eventid)

	return
}
*/
