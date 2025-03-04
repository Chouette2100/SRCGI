// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	"bytes"
	"fmt"
	// "html"
	"log"

	//	"math/rand"
	// "sort"
	"strconv"
	"strings"
	"time"
	//	"os"


	// "runtime"

	"encoding/json"

	//	"html/template"
	"net/http"

	// "database/sql"

	// _ "github.com/go-sql-driver/mysql"

	"github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)
func GetEventInfAndRoomListBR(
	client *http.Client,
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	status int,
) {

	status = 0

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid

	status = 0

	//	URLからドキュメントを作成します
	_url := "https://www.showroom-live.com/event/" + eventidorfilename
	/*
		doc, err = goquery.NewDocument(_url)
	*/
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
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	cevent_id, exists := doc.Find("#eventDetail").Attr("data-event-id")
	if !exists {
		log.Printf("data-event-id not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	eventinfo.I_Event_ID, _ = strconv.Atoi(cevent_id)
	event_id := eventinfo.I_Event_ID

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

	eia := strings.Split(eventid, "?")
	bia := strings.Split(eia[1], "=")
	blockid, _ := strconv.Atoi(bia[1])

	ebr, err := srapi.GetEventBlockRanking(client, event_id, blockid, breg, ereg)
	if err != nil {
		log.Printf("GetEventBlockRanking() err=%s\n", err.Error())
		status = 1
		return
	}

	ReplaceString := "/r/"

	for _, br := range ebr.Block_ranking_list {

		var roominfo RoomInfo

		roominfo.ID = br.Room_id
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		roominfo.Account = strings.Replace(br.Room_url_key, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, "/", "", -1)

		roominfo.Name = br.Room_name

		*roominfolist = append(*roominfolist, roominfo)

	}

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}

func GetEventListByAPI(eventinflist *[]exsrapi.Event_Inf) (status int) {

	status = 0

	last_page := 1
	total_count := 1

	for page := 1; page <= last_page; page++ {

		URL := "https://www.showroom-live.com/api/event/search?page=" + fmt.Sprintf("%d", page)
		log.Printf("GetEventListByAPI() URL=%s\n", URL)

		resp, err := http.Get(URL)
		if err != nil {
			//	一時的にデータが取得できない。
			log.Printf("GetEventListByAPI() err=%s\n", err.Error())
			//		panic(err)
			status = -1
			return
		}
		defer resp.Body.Close()

		//	JSONをデコードする。
		//	次の記事を参考にさせていただいております。
		//		Go言語でJSONに泣かないためのコーディングパターン
		//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

		var result interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			log.Printf("GetEventListByAPI() err=%s\n", err.Error())
			//	panic(err)
			status = -2
			return
		}

		if page == 1 {
			value, _ := result.(map[string]interface{})["last_page"].(float64)
			last_page = int(value)
			value, _ = result.(map[string]interface{})["total_count"].(float64)
			total_count = int(value)
			log.Printf("GetEventListByAPI() total_count=%d, last_page=%d\n", total_count, last_page)
		}

		noroom := 30
		if page == last_page {
			noroom = total_count % 30
			if noroom == 0 {
				noroom = 30
			}
		}

		for i := 0; i < noroom; i++ {
			var eventinf exsrapi.Event_Inf

			tres := result.(map[string]interface{})["event_list"].([]interface{})[i]

			ttres := tres.(map[string]interface{})["league_ids"]
			norec := len(ttres.([]interface{}))
			if norec == 0 {
				continue
			}
			log.Printf("norec =%d\n", norec)
			eventinf.League_ids = ""
			/*
				for j := 0; j < norec; j++ {
					eventinf.League_ids += ttres.([]interface{})[j].(string) + ","
				}
			*/
			eventinf.League_ids = ttres.([]interface{})[norec-1].(string)
			if eventinf.League_ids != "60" {
				continue
			}

			eventinf.Event_ID, _ = tres.(map[string]interface{})["event_url_key"].(string)
			eventinf.Event_name, _ = tres.(map[string]interface{})["event_name"].(string)
			//	log.Printf("id=%s, name=%s\n", eventinf.Event_ID, eventinf.Event_name)

			started_at, _ := tres.(map[string]interface{})["started_at"].(float64)
			eventinf.Start_time = time.Unix(int64(started_at), 0)
			eventinf.Sstart_time = eventinf.Start_time.Format("06/01/02 15:04")
			ended_at, _ := tres.(map[string]interface{})["ended_at"].(float64)
			eventinf.End_time = time.Unix(int64(ended_at), 0)
			eventinf.Send_time = eventinf.End_time.Format("06/01/02 15:04")

			(*eventinflist) = append((*eventinflist), eventinf)

		}

		//	resp.Body.Close()
	}

	return
}

func GetSerialFromYymmddHhmmss(yymmdd, hhmmss string) (tserial float64) {

	var year, month, day, hh, mm, ss int

	t19000101 := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)

	fmt.Sscanf(yymmdd, "%d/%d/%d", &year, &month, &day)
	fmt.Sscanf(hhmmss, "%d:%d:%d", &hh, &mm, &ss)

	t1 := time.Date(year, time.Month(month), day, hh, mm, ss, 0, time.Local)

	tserial = t1.Sub(t19000101).Minutes() / 60.0 / 24.0

	return
}

func GetUserInfForHistory(client *http.Client) (status int) {

	status = 0

	//	select distinct(nobasis) from event
	stmt, err := srdblib.Db.Prepare("select distinct(nobasis) from event")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var roominf RoomInfo
	var roominflist RoomInfoList

	for rows.Next() {
		err := rows.Scan(&roominf.Userno)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		if roominf.Userno != 0 {
			roominf.ID = fmt.Sprintf("%d", roominf.Userno)
			roominflist = append(roominflist, roominf)
		}
	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	eventid := ""

	//	Update user , Insert into userhistory
	for _, roominf := range roominflist {

		sql := "select currentevent from user where userno = ?"
		err := srdblib.Db.QueryRow(sql, roominf.Userno).Scan(&eventid)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
		}

		roominf.Genre, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level,
			roominf.Followers, roominf.Fans, roominf.Fans_lst, roominf.Name, roominf.Account, _, status = GetRoomInfoByAPI(roominf.ID)
		InsertIntoOrUpdateUser(client, time.Now().Truncate(time.Second), eventid, roominf)
	}

	return
}

