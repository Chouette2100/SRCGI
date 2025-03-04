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
	//	"time"

	//	"bufio"
	//	"os"

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

func HandlerNewUser(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/new-user.gtpl"))

	eventid := r.FormValue("eventid")

	log.Printf("      eventid=%s\n", eventid)

	roomid := r.FormValue("roomid")
	userno, _ := strconv.Atoi(roomid)
	log.Printf("eventid=%s, roomid=%s\n", eventid, roomid)

	//	eventno, eventname, period := SelectEventNoAndName(eventid)
	//	log.Printf("eventname=%s, period=%s\n", eventname, period)

	//	Event_inf, _ = SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, _ := srdblib.SelectFromEvent("event", eventid)
	Event_inf = *eventinf

	log.Printf("eventname=%s, period=%s\n", Event_inf.Event_name, Event_inf.Period)

	user := srdblib.User{}
	status_db := -1
	status_api := -1
	status_col := -1

	itfc, _ := srdblib.Dbmap.Get(srdblib.User{}, userno)
	if itfc != nil {
		user = *itfc.(*srdblib.User)
		status_db = 0
	} else {
		user.Genre, user.Rank, user.Nrank, user.Prank, user.Level,
			user.Followers, user.Fans, user.Fans_lst, user.User_name, user.Userid, _, status_api = GetRoomInfoByAPI(roomid)
	}

	log.Printf("genre=%s, level=%d, followers=%d, fans=%d, fans_lst=%d, roomname=%s, roomurlkey=%s, status=%d\n",
		user.Genre, user.Level, user.Followers, user.Fans, user.Fans_lst, user.User_name, user.Userid, status_api)

	//	roominf, status_room := SelectRoomInf(userno)

	if status_db != 0 && status_api != 0 {
		user.User_name = strconv.Itoa(userno)
	}
	if status_db != 0 {
		user.Longname = user.User_name
		user.Shortname = fmt.Sprintf("%d", userno%100)
	}

	//	eventuserから色割当を取得する　=> 取得できなければ未登録、取得できればイベントに登録済み
	_, _, status_col = SelectUserColor(userno, Event_inf.Event_ID)

	values := map[string]string{
		"Event_ID":   eventid,
		"Event_name": Event_inf.Event_name,
		"Period":     Event_inf.Period,
		"Roomid":     roomid,
		"Roomname":   user.User_name,
		"Longname":   user.Longname,
		"Shortname":  user.Shortname,
		"Roomurlkey": user.Userid,
		"Genre":      user.Genre,
		"Rank":       user.Rank,
		"Nrank":      user.Nrank,
		"Prank":      user.Prank,
		"Level":      fmt.Sprintf("%d", user.Level),
		"Followers":  fmt.Sprintf("%d", user.Followers),
		"Fans":       fmt.Sprintf("%d", user.Fans),
		"Fans_lst":   fmt.Sprintf("%d", user.Fans_lst),
		"Submit":     "submit",
		"Label":      "登録しない",
		"Msg1":       "の参加ルームとして",
		"Msg2":       "を登録しますか？（（実害はありませんが）ブロックイベントはblock_idが違っていても登録されるので注意してください）",
		"Msg2color":  "black",
	}

	if status_col == 0 {
		values["Submit"] = "hidden"
		values["Label"] = "戻る"
		values["Msg1"] = "の参加ルームとして"
		values["Msg2"] = "すでに登録されています"
		values["Msg2color"] = "red"
	//	} else if status_room != 0 {
	//		values["Submit"] = "hidden"
	//		values["Label"] = "戻る"
	//		values["Msg1"] = ""
	//		values["Msg2"] = "ルーム情報がDB未登録です"
	//		values["Msg2color"] = "red"
	} else {

		if status_db != 0 && status_api != 0 {
			values["Roomname"] = ""
			values["Roomurlkey"] = ""
			values["Genre"] = ""
			values["Nrank"] = ""
			values["Prank"] = ""
			values["Level"] = ""
			values["Followers"] = ""
			values["Fans"] = ""
			values["Fans_lst"] = ""
			values["Submit"] = "hidden"
			values["Label"] = "戻る"
			values["Msg1"] = ""
			values["Msg2"] = "指定したルームIDのルームは存在しないか、ルーム情報の取得ができません"
			values["Msg2color"] = "red"
		} else {
			_, _, _, peventid := GetPointsByAPI(roomid)
			if strings.Contains(eventid, "?block_id=") {
				eida := strings.Split(eventid, "?")
				if strings.Contains(peventid, eida[0]) {
					//	block_id=0 はブロックイベントの全体を意味するのでいかなるblock_idのルームもこれに属する。
					peventid = eventid
				}
			}

			//	if peventid != eventid && time.Now().After(Event_inf.Start_time) && time.Now().Before(Event_inf.End_time) {
			if peventid != eventid {
				values["Submit"] = "hidden"
				values["Label"] = "戻る"
				values["Msg1"] = ""
				values["Msg2"] = "指定したルームはこのイベントに参加していません(あるいはイベントが始まっていません)"
				values["Msg2color"] = "red"
				log.Printf("GetPointsByAPI() returned %s as eventid and eventid = %s\n", peventid, eventid)
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "new-user.gtpl", values); err != nil {
		log.Println(err)
	}

}
