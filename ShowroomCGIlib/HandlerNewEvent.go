// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	//	"database/sql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)
func HandlerNewEvent(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/new-event0.gtpl",
		"templates/new-event1.gtpl",
		"templates/new-event2.gtpl",
	))

	eventid := r.FormValue("eventid")
	suserno := r.FormValue("userno")
	userno, _ := strconv.Atoi(suserno)

	log.Printf("      eventid=%s\n", eventid)

	stm, sts := exsrapi.MakeSampleTime(240, 40)

	values := map[string]string{
		"Eventid":   r.FormValue("eventid"),
		"Eventname": "",
		"Period":    "",
		"Noroom":    "",
		"Msgcolor":  "blue",

		"Stm": fmt.Sprintf("%d", stm),
		"Sts": fmt.Sprintf("%d", sts),

		"Maxcmap": strconv.Itoa(len(Colormaplist)),
	}

	var eventinf exsrapi.Event_Inf

	eia := strings.Split(eventid, "?")
	if len(eia) == 2 {
		eventid = eia[0]
	}

	status := GetEventInf(eventid, &eventinf)
	if status == -1 {
		values["Msg"] = "このイベントはすでに登録されています。"
		values["Submit"] = "hidden"
		values["Msgcolor"] = "red"
		//	Event_inf, _ = SelectEventInf(eventid)
		//	srdblib.Tevent = "event"
		eventinf, _ := srdblib.SelectFromEvent("event", eventid)
		Event_inf = *eventinf

		values["Eventname"] = Event_inf.Event_name
		values["Period"] = Event_inf.Period
	} else if status == -2 {
		values["Msg"] = "指定したIDのイベントは存在しません"
		values["Submit"] = "hidden"
		values["Msgcolor"] = "red"
	} else if status < -2 {
		values["Msg"] = "イベント情報を取得できませんでした（エラーコード＝" + fmt.Sprintf("%d", status) + "）"
		values["Submit"] = "hidden"
		values["Msgcolor"] = "red"
	} else {
		values["Msg"] = "このイベントを登録しますか？"
		values["Submit"] = "submit"
		values["Eventname"] = eventinf.Event_name
		values["Period"] = eventinf.Period
		values["Noroom"] = "　" + humanize.Comma(int64(eventinf.NoEntry))
	}
	/*
		var Eventinflist []Event_Inf
		GetEventListByAPI(&Eventinflist)
	*/

	userlist, _ := SelectUserList()
	userlist[0].Userlongname = "基準とする配信者を設定しない"
	for i := 0; i < len(userlist); i++ {
		if userlist[i].Userno == userno {
			userlist[i].Selected = "Selected"
		} else {
			userlist[i].Selected = ""
		}
	}

	if err := tpl.ExecuteTemplate(w, "new-event0.gtpl", values); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "new-event1.gtpl", userlist); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "new-event2.gtpl", values); err != nil {
		log.Println(err)
	}

}