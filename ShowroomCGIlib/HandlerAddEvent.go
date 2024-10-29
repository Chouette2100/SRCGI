// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	//	"html"
	"log"
	"strconv"
	//	"strings"
	"sort"
	"time"

	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	//	"database/sql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

// イベントを獲得ポイントデータ取得の対象としてeventテーブルに登録する。
// イベントが開催中であれば指定した順位内のルームを取得対象として登録する。
// イベントが開催予定のものであればルームの登録は行わない。
// イベント開催中、開催予定にかかわらず、取得対象ルームの追加は srAddNewOnes で行われる。
func HandlerAddEvent(w http.ResponseWriter, r *http.Request) {

	ra, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	localhost := false
	if ra == "127.0.0.1" {
		localhost = true
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/add-event1.gtpl",
		"templates/add-event2.gtpl",
		"templates/error.gtpl",
	))

	eventinf := &exsrapi.Event_Inf{}
	var roominfolist RoomInfoList

	eventid := r.FormValue("eventid")
	breg := r.FormValue("breg")
	ereg := r.FormValue("ereg")
	ibreg, _ := strconv.Atoi(breg)
	iereg, _ := strconv.Atoi(ereg)
	log.Printf(" eventid =[%s], ibreg=%d, iereg=%d\n", eventid, ibreg, iereg)

	if r.FormValue("from") != "new-event" {
		//	すでに登録済みのイベントの参加ルームの更新を行うとき
		//	eventinf, _ = SelectEventInf(eventid)
		//	srdblib.Tevent = "event"
		eventinf, _ = srdblib.SelectFromEvent("event", eventid)

		//	w.Write([]byte("Called. not 'from new-event'\n"))
		log.Printf("  Called. not 'from new-event'\n")

		bufstr := fmt.Sprintf("eventinf=%+v\n", eventinf)
		//	w.Write([]byte(bufstr))
		log.Printf("%s\n", bufstr)
	} else {
		//	新規にイベントを登録するとき
		//	eventinf = &exsrapi.Event_Inf{}
		//	srdblib.Tevent = "wevent"
		eventinf, _ = srdblib.SelectFromEvent("wevent", eventid)
		if eventinf == nil {
			log.Printf("[%s] is not found in wevent table\n", eventid)
			values := map[string]string{
				"Msg001":   "入力したイベントID( ",
				"Msg002":   " )をもつイベントは存在しません！",
				"ReturnTo": "top",
				"Eventid":  eventid,
			}
			if err := tpl.ExecuteTemplate(w, "error.gtpl", values); err != nil {
				log.Println(err)
			}
			return
		}

		log.Println("  Called. 'from new-event'")
		eventinf.Modmin, _ = strconv.Atoi(r.FormValue("modmin"))
		eventinf.Modsec, _ = strconv.Atoi(r.FormValue("modsec"))

		intervalmin, _ := strconv.Atoi(r.FormValue("intervalmin"))
		switch intervalmin {
		case 5, 6, 10, 15, 20, 30, 60:
			eventinf.Intervalmin = intervalmin
		default:
			eventinf.Intervalmin = 5
		}
		eventinf.Modmin = eventinf.Modmin % eventinf.Intervalmin //	不適切な入力に対する修正
		eventinf.Modsec = eventinf.Modsec % 60

		eventinf.Resethh, _ = strconv.Atoi(r.FormValue("resethh"))
		eventinf.Resetmm, _ = strconv.Atoi(r.FormValue("resetmm"))
		eventinf.Nobasis, _ = strconv.Atoi(r.FormValue("nobasis"))
		eventinf.Target, _ = strconv.Atoi(r.FormValue("target"))
		eventinf.Maxdsp, _ = strconv.Atoi(r.FormValue("maxdsp"))
		eventinf.Cmap, _ = strconv.Atoi(r.FormValue("cmap"))
	}
	eventinf.Fromorder = ibreg
	eventinf.Toorder = iereg

	inprogress := true
	if eventinf.Start_time.After(time.Now()) {
		//	イベント開催前
		inprogress = false
	}

	status := 0
	if !inprogress && !localhost {
	//	if !inprogress || localhost {
		//	イベントが開催前であり、かつローカルホストからの実行でもないとき
		//	イベント参加ルームの登録はできない
		status = InsertEventInf(localhost, eventinf)
	} else {
		//	イベントが開催中であるかローカルホストでの実行のとき
		//	イベント参加ルームの登録ができる

		Event_inf = *eventinf

		log.Println("before GetAndInsertEventRoomInfo()")
		log.Println(eventinf)

		//      cookiejarがセットされたHTTPクライアントを作る
		client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
		if err != nil {
			log.Printf("CreateNewClient: %s\n", err.Error())
			return
		}
		//      すべての処理が終了したらcookiejarを保存する。
		defer jar.Save()

		var ril *RoomInfoList
		ril, status = GetAndInsertEventRoomInfo(client, localhost, inprogress, eventid, ibreg, iereg, eventinf, &roominfolist)
		roominfolist = *ril
	}
	if status != 0 {

		values := map[string]string{
			"Msg001":   "入力したイベントID( ",
			"Msg002":   " )をもつイベントが存在しないか、開催前イベントか、参加ルームがありません！",
			"ReturnTo": "top",
			"Eventid":  eventid,
		}
		if err := tpl.ExecuteTemplate(w, "error.gtpl", values); err != nil {
			log.Println(err)
		}

	} else {

		if err := tpl.ExecuteTemplate(w, "add-event1.gtpl", eventinf); err != nil {
			log.Println(err)
		}
		if iereg > len(roominfolist) {
			iereg = len(roominfolist)
		}
		if err := tpl.ExecuteTemplate(w, "add-event2.gtpl", roominfolist[ibreg-1:iereg]); err != nil {
			log.Println(err)
		}
	}

}

// イベントを新規に登録する。
// TODO: gorpを使うべき
func InsertEventInf(localhot bool, eventinf *exsrapi.Event_Inf) (
	status int,
) {

	status = 0
	if _, _, sts := SelectEventNoAndName((*eventinf).Event_ID); sts != 0 {
		sql := "INSERT INTO event(eventid, ieventid, event_name, period, starttime, endtime, noentry,"
		sql += " intervalmin, modmin, modsec, "
		sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, Cmap, target, maxpoint "
		sql += ") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		log.Printf("db.Prepare(sql)\n")
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error InsertEventInf() (INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		if eventinf.Intervalmin != 5 { //	緊急対応
			log.Printf(" Intervalmin isn't 5. (%dm)\n", eventinf.Intervalmin)
			eventinf.Intervalmin = 5
		}

		log.Printf("row.Exec()\n")
		_, err = stmt.Exec(
			(*eventinf).Event_ID,
			(*eventinf).I_Event_ID,
			(*eventinf).Event_name,
			(*eventinf).Period,
			(*eventinf).Start_time,
			(*eventinf).End_time,
			(*eventinf).NoEntry,
			(*eventinf).Intervalmin,
			(*eventinf).Modmin,
			(*eventinf).Modsec,
			(*eventinf).Fromorder,
			(*eventinf).Toorder,
			(*eventinf).Resethh,
			(*eventinf).Resetmm,
			(*eventinf).Nobasis,
			(*eventinf).Maxdsp,
			(*eventinf).Cmap,
			(*eventinf).Target,
			(*eventinf).Maxpoint+eventinf.Gscale,
		)

		if err != nil {
			log.Printf("error InsertEventInf() (INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
	} else {
		status = 1
	}

	return
}

// イベントの参加ルーム情報を取得しeventuserテーブルに格納する。
func GetAndInsertEventRoomInfo(
	client *http.Client,
	localhost bool,
	inprogress bool,
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	ril *RoomInfoList,
	status int,
) {

	log.Println("GetAndInsertEventRoomInfo() Called.")
	log.Println(*eventinfo)

	status = 0

	//	イベントに参加しているルームの一覧を取得します。
	//	ルーム名、ID、URLを取得しますが、イベント終了直後の場合の最終獲得ポイントが表示されている場合はそれも取得します。

	/*	ここから
		if strings.Contains(eventid, "?") {
			status = GetEventInfAndRoomListBR(client, eventid, breg, ereg, eventinfo, roominfolist)
			eia := strings.Split(eventid, "?")
			bka := strings.Split(eia[1], "=")
			eventinfo.Event_name = eventinfo.Event_name + "(" + bka[1] + ")"
		} else {
			status = GetEventInfAndRoomList(eventid, breg, ereg, eventinfo, roominfolist)
		}

		if status != 0 {
			log.Printf("GetEventInfAndRoomList() returned %d\n", status)
			return
		}

		//	各ルームのジャンル、ランク、レベル、フォロワー数を取得します。
		for i := 0; i < (*eventinfo).NoRoom; i++ {
			(*roominfolist)[i].Genre, (*roominfolist)[i].Rank,
				(*roominfolist)[i].Nrank,
				(*roominfolist)[i].Prank,
				(*roominfolist)[i].Level,
				(*roominfolist)[i].Followers,
				(*roominfolist)[i].Fans,
				(*roominfolist)[i].Fans_lst,
				_, _, _, _ = GetRoomInfoByAPI((*roominfolist)[i].ID)

		}

		//	各ルームの獲得ポイントを取得します。
		for i := 0; i < (*eventinfo).NoRoom; i++ {
			point, _, _, eventid := GetPointsByAPI((*roominfolist)[i].ID)
			if eventid == (*eventinfo).Event_ID {
				(*roominfolist)[i].Point = point
				UpdateEventuserSetPoint(eventid, (*roominfolist)[i].ID, point)
				if point < 0 {
					(*roominfolist)[i].Spoint = ""
				} else {
					(*roominfolist)[i].Spoint = humanize.Comma(int64(point))
				}
			} else {
				log.Printf(" %s %s %d\n", (*eventinfo).Event_ID, eventid, point)
			}

			if (*roominfolist)[i].ID == fmt.Sprintf("%d", (*eventinfo).Nobasis) {
				(*eventinfo).Pntbasis = point
				(*eventinfo).Ordbasis = i
			}

			//	log.Printf(" followers=<%d> level=<%d> nrank=<%s> genre=<%s> point=%d\n",
			//	(*roominfolist)[i].Followers,
			//	(*roominfolist)[i].Level,
			//	(*roominfolist)[i].Nrank,
			//	(*roominfolist)[i].Genre,
			//	(*roominfolist)[i].Point)
		}

		ここまで修正対象	*/

	//	ここから

	lenpr := 0
	eid := eventid
	//	if strings.Contains(eventid, "?block_id=0") {
	//		//	Block_id=0 は特殊で、ブロックイベントを構成するイベントの一つということではなくイベント全体を示す
	//		eid = strings.ReplaceAll(eventid, "?block_id=0", "")
	//	}
	//	ランキングイベントの1〜50位の結果を取得する。
	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
	pranking, err := srdblib.GetEventsRankingByApi(client, eid, 2)
	//	if err != nil && !localhost {
	if err != nil {
		log.Printf("GetAndInsertEventRoomInfo() GetEventsRankingByApi(client, %s, 2) err=%s\n", eid, err.Error())
		status = -1
		return
	}
	//	} else {
	//	if localhost {
	//		lenpr = 0
	//	} else {
	lenpr = len(pranking.Ranking)
	//	}
	//	}

	if lenpr != 0 {
		//	for _, rinf := range(pranking.Ranking) {
		if ereg > lenpr {
			ereg = lenpr
		}
		if breg > ereg {
			breg = ereg
		}
		for i := breg; i <= ereg; i++ {
			rinf := pranking.Ranking[i-1]
			if rinf.Point == 0 {
				ereg = i - 1
				break
			}
			roominf := RoomInfo{
				Name:    rinf.Room.Name,
				ID:      strconv.Itoa(rinf.Room.RoomID),
				Userno:  rinf.Room.RoomID,
				Account: "",
				Point:   rinf.Point,
				Order:   rinf.Rank,
				Irank:   888888888,
			}
			*roominfolist = append(*roominfolist, roominf)
		}
	} else {
		//	レベルイベントのとき
		status = GetEventInfAndRoomList(eventid, breg, ereg, eventinfo, roominfolist)
		if status != 0 {
			log.Printf("GetEventInfAndRoomList() returned %d\n", status)
			return
		}
		lenpr = len(*roominfolist)
		if ereg > lenpr {
			ereg = lenpr
		}
		if breg > ereg {
			breg = ereg
		}
		if inprogress {
			for i := breg - 1; i < ereg; i++ {
				point, _, _, eventid := GetPointsByAPI((*roominfolist)[i].ID)
				if eventid == (*eventinfo).Event_ID {
					if point == 0 {
						ereg = i
						break
					}
					(*roominfolist)[i].Point = point
					UpdateEventuserSetPoint(eventid, (*roominfolist)[i].ID, point)
				} else {
					log.Printf(" %s %s %d\n", (*eventinfo).Event_ID, eventid, point)
				}

			}
		}

		nroominfolist := (*roominfolist)[breg-1 : ereg]
		roominfolist = &nroominfolist

	}

	for i, rminf := range *roominfolist {
		if rminf.Point < 0 {
			(*roominfolist)[i].Spoint = ""
		} else {
			(*roominfolist)[i].Spoint = humanize.Comma(int64(rminf.Point))
		}

		if rminf.Userno == eventinfo.Nobasis {
			(*eventinfo).Pntbasis = rminf.Point
			(*eventinfo).Ordbasis = i
		}
	}

	//	ここまで新規作成

	if !inprogress {
		SortByFollowers = true
		sort.Sort(*roominfolist)
		if ereg > len(*roominfolist) {
			ereg = len(*roominfolist)
		}
		r := (*roominfolist).Choose(breg-1, ereg)
		roominfolist = &r
	}

	log.Printf(" GetEventRoomInfo() len(*roominfolist)=%d\n", len(*roominfolist))

	//	srdblib.Tevent = "wevent"
	weventinf, _ := srdblib.SelectFromEvent("wevent", eventid)

	eventinfo.Event_name = weventinf.Event_name

	log.Println("GetAndInsertEventRoomInfo() before InsertEventIinf()")
	log.Println(*eventinfo)
	status = InsertEventInf(localhost, eventinfo)

	if status == 1 {
		log.Println("InsertEventInf() returned 1.")
		UpdateEventInf(eventinfo)
		status = 0
	}
	log.Println("GetAndInsertEventRoomInfo() after InsertEventIinf() or UpdateEventInf")
	log.Println(*eventinfo)

	if !inprogress && !localhost {
		//	イベント開始前で、かつローカルホストからの実行でないときはルームの登録は行わない。
		return
	}

	_, _, status = SelectEventNoAndName(eventid)

	if status == 0 {
		//	InsertRoomInf(eventno, eventid, roominfolist)
		InsertRoomInf(client, eventid, roominfolist)
		for i, rinf := range *roominfolist {
			ifc, _ := srdblib.Dbmap.Get(srdblib.User{}, rinf.Userno)
			if ifc != nil {
				user := ifc.(*srdblib.User)
				(*roominfolist)[i].Account = user.Userid
				(*roominfolist)[i].Genre = user.Genre
				(*roominfolist)[i].Rank = user.Rank
				(*roominfolist)[i].Level = user.Level
				(*roominfolist)[i].Followers = user.Followers
			}

		}
	}

	ril = roominfolist

	return
}
