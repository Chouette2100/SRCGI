// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	// "bytes"
	"fmt"
	//	"html"
	"log"
	// "os"
	"sort"
	"strconv"
	"strings"
	"time"

	// "github.com/PuerkitoBio/goquery"

	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	//	"database/sql"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// イベントを獲得ポイントデータ取得の対象としてeventテーブルに登録する。
// イベントが開催中であれば指定した順位内のルームを取得対象として登録する。
// イベントが開催予定のものであればルームの登録は行わない。
// イベント開催中、開催予定にかかわらず、取得対象ルームの追加は srAddNewOnes で行われる。
func AddEventHandler(w http.ResponseWriter, r *http.Request) {

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

	bnew := true
	if r.FormValue("from") != "new-event" {
		//	すでに登録済みのイベントの参加ルームの更新を行うとき
		bnew = false
		//	eventinf, _ = SelectEventInf(eventid)
		//	srdblib.Tevent = "event"
		eventinf, _ = srdblib.SelectFromEvent("event", eventid)

		//	w.Write([]byte("Called. not 'from new-event'\n"))
		log.Printf("  Called. not 'from new-event'\n")

		//	bufstr := fmt.Sprintf("eventinf=%+v\n", eventinf)
		//	w.Write([]byte(bufstr))
		//	log.Printf("%s\n", bufstr)
		//	w.Write([]byte("すでに獲得ポイント取得の対象となったイベントに対しては、この機能はメンテナンス中のため使用できません。\n"))
		//	w.Write([]byte("特定のルームを対象として追加したいときは「イベントトップ」=>「(DB登録済み)イベント参加ルーム一覧（確認・編集）」=>「一覧にないルームの追加」をお使いください。\n"))
		//	w.Write([]byte("また「イベントトップ」=>「イベントパラメータの設定」=>「ＤＢに登録する順位の範囲」での対象範囲の変更は有効です\n"))
		//	w.Write([]byte("メンテナンスは掲示板のNo.1239【お知らせ】「イベント獲得ポイントデータ取得範囲について」に関連する処理の修正中であるためです）"))
		//	return
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

		thdata, err := exsrapi.ReadThdata()
		if err != nil {
			err = fmt.Errorf("ReadThdata: %w", err)
			w.Write([]byte(err.Error()))
		}
		err = exsrapi.SetThdata(eventinf, thdata)
		if err != nil {
			err = fmt.Errorf("SetThdata: %w", err)
			w.Write([]byte(err.Error()))
		}
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
		if bnew {
			status = InsertEventInf(localhost, eventinf)
		}
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
		ril, status = GetAndInsertEventRoomInfo(client, localhost, inprogress, bnew, eventid, ibreg, iereg, eventinf, &roominfolist)
		if ril != nil {
			roominfolist = *ril
		} else {
			roominfolist = RoomInfoList{}
		}
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
		sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, Cmap, target, maxpoint, thinit, thdelta "
		sql += ") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
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
			eventinf.Thinit,
			eventinf.Thdelta,
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
	bnew bool,
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

	lenpr := 0
	eid := eventid
	//	ランキングイベントの1〜50位の結果を取得する。
	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
	pranking, err := srdblib.GetEventsRankingByApi(client, eid, 2)
	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	//	if err != nil && !localhost {
	bnoroom := false
	if err != nil {
		log.Printf("GetAndInsertEventRoomInfo() GetEventsRankingByApi(client, %s, 2) err=%s\n", eid, err.Error())
		if !strings.Contains(err.Error(), "has no room") {
			status = 1
			return
		} else {
			bnoroom = true
		}
	}

	if bnoroom {
		lenpr = 0
	} else {
		lenpr = len(pranking.Ranking)
	}

	ignorethpoint := false
	if lenpr > 0 && lenpr <= 30 {
		ignorethpoint = true
	}

	hh := time.Since(eventinfo.Start_time).Hours()
	//	thpoint := eventinfo.Thdelta*(int(hh)) + eventinfo.Thinit
	thpoint := eventinfo.Thdelta * (int(hh))
	if thpoint < eventinfo.Thinit {
		thpoint = eventinfo.Thinit
	}
	log.Printf("%s Starttime=%s Hours=%7.2f\n", eventid, eventinfo.Start_time.Format("2006-01-02 15:04:05"), hh)
	log.Printf("%s hh=%d Thinit=%d Thdelta=%d thpoint=%d\n", eventid, int(hh), eventinfo.Thinit, eventinfo.Thdelta, thpoint)

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
			if !ignorethpoint && rinf.Point < thpoint {
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
		/*
			status = GetEventInfAndRoomList(eventid, breg, ereg, eventinfo, roominfolist)
			if status != 0 {
				log.Printf("GetEventInfAndRoomList() returned %d\n", status)
				return
			}
			lenpr = len(*roominfolist)
		*/
		var eqr *srapi.EventQuestRooms
		eqr, err = srapi.GetEventQuestRoomsByApi(client, eventid, 1, eventinfo.Toorder)
		lenpr = len(eqr.EventQuestLevelRanges[0].Rooms)
		if ereg > lenpr {
			ereg = lenpr
		}
		if breg > ereg {
			breg = 1
		}

		roominfolist = &RoomInfoList{}

		for i := breg; i <= ereg; i++ {
			room := eqr.EventQuestLevelRanges[0].Rooms[i-1]
			point, _, _, eventid := GetPointsByAPI(strconv.Itoa(room.RoomID))
			if eventid == (*eventinfo).Event_ID {
				if point < thpoint {
					ereg = i
					break
				}
				// (*roominfolist)[i].Point = point
				roominf := RoomInfo{
					Name:    room.RoomName,
					ID:      strconv.Itoa(room.RoomID),
					Userno:  room.RoomID,
					Account: "",
					Point:   point,
					Order:   i,
					Irank:   888888888,
				}
				*roominfolist = append(*roominfolist, roominf)
				UpdateEventuserSetPoint(eventid, strconv.Itoa(room.RoomID), point)
			} else {
				log.Printf(" %s %s %d\n", (*eventinfo).Event_ID, eventid, point)
			}
		}

		/*
				for i := breg - 1; i < ereg; i++ {
					point, _, _, eventid := GetPointsByAPI((*roominfolist)[i].ID)
					if eventid == (*eventinfo).Event_ID {
						if point < thpoint {
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
		*/
		// 	レベルイベントの処理で獲得ポイントが0でないルームの前に0のルームがあるといったイレギュラーなケースに対する対応
		nroominfolist := RoomInfoList{}
		lenpr = len(*roominfolist)
		if lenpr != 0 {
			nroominfolist = (*roominfolist)[breg-1 : ereg]
			roominfolist = &nroominfolist
		}
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

	if bnew {
		status = InsertEventInf(localhost, eventinfo)
	} else {
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
func InsertRoomInf(client *http.Client, eventid string, roominfolist *RoomInfoList) {

	log.Printf("  *** InsertRoomInf() ***********  NoRoom=%d\n", len(*roominfolist))
	//	srdblib.Dbmap.AddTableWithName(srdblib.User{}, "user").SetKeys(false, "Userno")
	tnow := time.Now().Truncate(time.Second)
	for i := 0; i < len(*roominfolist); i++ {
		//	log.Printf("   ** InsertRoomInf() ***********  i=%d\n", i)
		user := new(srdblib.User)
		user.Userno = (*roominfolist)[i].Userno
		// err := srdblib.UpinsUserSetProperty(client, tnow, user, 1440*5, 200)
		srdblib.Env.Waitmsec = 200 // FIXME: 危険
		_, err := srdblib.UpinsUser(client, tnow, user)
		srdblib.Env.Waitmsec = 5000
		if err != nil {
			log.Printf("srdblib.UpinsUserSetProperty(): err=%v\n", err)
			return
		}
		//	InsertIntoOrUpdateUser(client, tnow, eventid, (*roominfolist)[i])
		status := InsertIntoEventUser(i, eventid, (*roominfolist)[i])
		switch status {
		case 0:
			(*roominfolist)[i].Status = "更新"
			(*roominfolist)[i].Statuscolor = "black"
		case 1:
			(*roominfolist)[i].Status = "新規"
			(*roominfolist)[i].Statuscolor = "green"

			/* この処理はInsertIntoEventUser()に含まれている
			userno, _ := strconv.Atoi((*roominfolist)[i].ID)
			eventinf, _ := srdblib.SelectFromEvent("event", eventid)
			sqlip := "insert into points (ts, user_id, eventid, point, `rank`, gap, pstatus) values(?,?,?,?,?,?,?)"
			_, srdblib.Dberr = srdblib.Db.Exec(
				sqlip,
				eventinf.Start_time.Truncate(time.Second),
				userno,
				eventid,
				0,
				1,
				0,
				"=",
			)
			if srdblib.Dberr != nil {
				err := fmt.Errorf("Db.Exec(sqlip,...): %w", srdblib.Dberr)
				log.Printf("err=[%s]\n", err.Error())
			}
			*/

		default:
			(*roominfolist)[i].Status = "エラー"
			(*roominfolist)[i].Statuscolor = "red"
		}
	}
	log.Printf("  *** end of InsertRoomInf() ***********\n")
}

/*
func GetEventInfAndRoomList(
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	status int,
) {

	//	画面からのデータ取得部分は次を参考にしました。
	//		はじめてのGo言語：Golangでスクレイピングをしてみた
	//		https://qiita.com/ryo_naka/items/a08d70f003fac7fb0808

	//	_url := "https://www.showroom-live.com/event/" + EventID
	//	_url = "file:///C:/Users/kohei47/Go/src/EventRoomList03/20210128-1143.html"
	//	_url = "file:20210128-1143.html"

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid
	maxroom := ereg

	status = 0

	if inputmode == "file" {

		//	ファイルからドキュメントを作成します
		f, e := os.Open(eventidorfilename)
		if e != nil {
			//	log.Fatal(e)
			log.Printf("err=[%s]\n", e.Error())
			status = -1
			return
		}
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}

		content, _ := doc.Find("head > meta:nth-child(6)").Attr("content")
		content_div := strings.Split(content, "/")
		(*eventinfo).Event_ID = content_div[len(content_div)-1]

	} else {
		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventidorfilename
		/*
			doc, err = goquery.NewDocument(_url)
		-/
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
			status = 1
			return
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		//	bufstr := buf.String()
		//	log.Printf("%s\n", bufstr)

		//	doc, error = goquery.NewDocumentFromReader(resp.Body)
		doc, error = goquery.NewDocumentFromReader(buf)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			status = 1
			return
		}

		(*eventinfo).Event_ID = eventidorfilename
	}
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	cevent_id, exists := doc.Find("#eventDetail").Attr("data-event-id")
	if !exists {
		log.Printf("data-event-id not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	eventinfo.I_Event_ID, _ = strconv.Atoi(cevent_id)

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	(*eventinfo).Period = selector.Find(".info").Text()
	eventinfo.Period = strings.Replace(eventinfo.Period, "\u202f", " ", -1)
	period := strings.Split((*eventinfo).Period, " - ")
	if inputmode == "url" {
		(*eventinfo).Start_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[1]+" JST")
	} else {
		(*eventinfo).Start_time, _ = time.Parse("2006/01/02 15:04 MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("2006/01/02 15:04 MST", period[1]+" JST")
	}

	(*eventinfo).EventStatus = "BeingHeld"
	if (*eventinfo).Start_time.After(time.Now()) {
		(*eventinfo).EventStatus = "NotHeldYet"
	} else if (*eventinfo).End_time.Before(time.Now()) {
		(*eventinfo).EventStatus = "Over"
	}

	//	イベントに参加しているルームの数を求めます。
	//	参加ルーム数と表示されているルームの数は違うので、ここで取得したルームの数を以下の処理で使うわけではありません。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)
	log.Printf("breg=%d ereg=%d\n", breg, ereg)

	//	eventno, _, _ := SelectEventNoAndName(eventidorfilename)
	//	log.Printf(" eventno=%d\n", eventno)
	//	(*eventinfo).Event_no = eventno

	//	抽出したルームすべてに対して処理を繰り返す(が、イベント開始後の場合の処理はルーム数がbreg、eregの範囲に限定）
	//	イベント開始前のときはすべて取得し、ソートしたあてで範囲を限定する）
	doc.Find(".listcardinfo").EachWithBreak(func(i int, s *goquery.Selection) bool {
		//	log.Printf("i=%d\n", i)
		if (*eventinfo).Start_time.Before(time.Now()) {
			if i < breg-1 {
				return true
			}
			if i == maxroom {
				return false
			}
		}

		var roominfo RoomInfo

		roominfo.Name = s.Find(".listcardinfo-main-text").Text()

		spoint1 := strings.Split(s.Find(".listcardinfo-sub-single-right-text").Text(), ": ")

		var point int64
		if spoint1[0] != "" {
			spoint2 := strings.Split(spoint1[1], "pt")
			fmt.Sscanf(spoint2[0], "%d", &point)

		} else {
			point = -1
		}
		roominfo.Point = int(point)

		ReplaceString := ""

		selection_c := s.Find(".listcardinfo-menu")

		account, _ := selection_c.Find(".room-url").Attr("href")
		if inputmode == "file" {
			ReplaceString = "https://www.showroom-live.com/"
		} else {
			ReplaceString = "/r/"
		}
		roominfo.Account = strings.Replace(account, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, ReplaceString, "/", -1)

		roominfo.ID, _ = selection_c.Find(".js-follow-btn").Attr("data-room-id")
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		*roominfolist = append(*roominfolist, roominfo)

		//	log.Printf("%11s %-20s %-10s %s\n",
		//		humanize.Comma(int64(roominfo.Point)), roominfo.Account, roominfo.ID, roominfo.Name)
		return true

	})

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}
*/
