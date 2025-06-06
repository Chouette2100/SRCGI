// Copyright © 2024 chouette.21.00@gmail.com
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

	//	"database/sql"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func EditUserHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/edit-user1.gtpl",
		"templates/edit-user2.gtpl",
		"templates/edit-user3.gtpl",
	))

	userid := r.FormValue("userid")
	eventid := r.FormValue("eventid")
	longname := r.FormValue("longname")
	shortname := r.FormValue("shortname")
	istarget := r.FormValue("istarget")
	graph := r.FormValue("graph")
	iscntrbpoint := r.FormValue("iscntrbpoint")
	color := r.FormValue("color")

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
	Event_inf = *eventinf

	fnc := r.FormValue("func")

	log.Printf("      func=%s eventid=%s userid=%s\n", fnc, eventid, userid)

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

		roominf, status := GetRoomInfoAndPoint(eventid, userid, fmt.Sprintf("%d", Event_inf.Nobasis))
		if status == 0 {
			tnow := time.Now().Truncate(time.Second)
			InsertIntoOrUpdateUser(client, tnow, eventid, roominf)
			InsertIntoEventUser(0, eventid, roominf)
			UpdateEventuserSetPoint(eventid, roominf.ID, roominf.Point)

		} else {
			log.Printf("GetAndUpdateRoomInfoAndPoint() returned %d", status)
		}

	case "deleteuser":
		//	削除ボタンが押されたとき
	default:
		//	（更新ボタンが押された配信者がいたらそのデータを更新した上で）参加配信者のリストを表示する。
		if userid != "" {
			UpdateRoomInf(eventid, userid, longname, shortname, istarget, graph, color, iscntrbpoint)
		}
	}

	//	log.Printf(" eventid=%s, userno=%s, longname=%s, shortname=%s, istarget=%s, graph=%s, color=%s\n",
	//		eventid, userno, longname, shortname, istarget, graph, color)

	var roominfolist RoomInfoList

	eventname, _ := SelectEventRoomInfList(eventid, &roominfolist)
	for i := 0; i < len(roominfolist); i++ {
		switch roominfolist[i].Genre {
		case "Voice Actors & Anime":
			roominfolist[i].Genre = "VA&A"
		case "Talent Model":
			roominfolist[i].Genre = "Tl/Md"
		case "Comedians/Talk Show":
			roominfolist[i].Genre = "Cm/TS"
		default:
		}
	}

	values := map[string]string{
		"Eventid":   eventid,
		"Eventname": eventname,
		"Period":    Event_inf.Period,
	}

	if err := tpl.ExecuteTemplate(w, "edit-user1.gtpl", values); err != nil {
		log.Println(err)
	}

	if err := tpl.ExecuteTemplate(w, "edit-user2.gtpl", roominfolist); err != nil {
		log.Println(err)
	}

	if err := tpl.ExecuteTemplate(w, "edit-user3.gtpl", values); err != nil {
		log.Println(err)
	}

}

func UpdateRoomInf(eventid, suserno, longname, shortname, istarget, graph, color, iscntrbpoint string) (status int) {

	status = 0

	userno, _ := strconv.Atoi(suserno)

	sql := "update user set longname=?, shortname=? where userno = ?"

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(longname, shortname, userno)

	if err != nil {
		log.Printf("UpdateRoomInf() error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		status = -2
		return
	}

	//	eventno, _, _ := SelectEventNoAndName(eventid)

	if istarget == "1" {
		istarget = "Y"
	} else {
		istarget = "N"
	}

	if graph == "1" {
		graph = "Y"
	} else {
		graph = "N"
	}

	if iscntrbpoint == "1" {
		iscntrbpoint = "Y"
	} else {
		iscntrbpoint = "N"
	}

	//	sql = "update eventuser set istarget=?, graph=?, color=? where eventno=? and userno=?"
	//	sql = "update eventuser set istarget=?, graph=?, color=?, iscntrbpoints=? where eventid=? and userno=?"
	sql = "update eventuser set graph=?, color=?"
	if iscntrbpoint == "Y" {
		sql += ", iscntrbpoints= 'Y'"
	}
	if istarget == "Y" {
		sql += ", istarget= 'Y'"
	}
	sql += " where eventid=? and userno=?"

	stmt, err = srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	//	_, err = stmt.Exec(istarget, graph, color, iscntrbpoint, eventid, userno)
	_, err = stmt.Exec(graph, color, eventid, userno)

	if err != nil {
		log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return

}
