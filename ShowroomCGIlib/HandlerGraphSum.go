// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Chouette2100/srdblib"
)

type HGSinf struct {
	Title     string
	Eventid   string
	Ieventid  int
	Eventname string
	Period    string
	Roomid    int
	Roomname  string
	Maxpoint  int
	Gscale    int
}

func HandlerGraphSum(w http.ResponseWriter, r *http.Request) {

	/*
		hgsinf := HGSinf{
			Title:     "獲得ポイントと累積ポイントの概要",
			Eventid:   "event2015",
			Ieventid:  99999,
			Eventname: "なんかのいべんと",
			Period:    "2025-01-01 00:00 ~ 2025-01-05 23:59",
			Roomid:    1,
			Roomname:  "謎の配信者",
		}
	*/

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	eventid := r.FormValue("eventid")
	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	intf, _ := srdblib.Dbmap.Get(srdblib.Event{}, eventid)
	hgsinf := HGSinf{
		Title:     "獲得ポイントと累積ポイントの概要",
		Eventid:   eventid,
		Ieventid:  intf.(*srdblib.Event).Ieventid,
		Eventname: intf.(*srdblib.Event).Event_name,
		Period:    intf.(*srdblib.Event).Period,
		Roomid:    roomid,
		Maxpoint:    intf.(*srdblib.Event).Maxpoint,
		// Gscale:    intf.(*srdblib.Event).Gscale,
		Gscale:    100,
	}

	intf, _ = srdblib.Dbmap.Get(srdblib.User{}, roomid)
	hgsinf.Roomname = intf.(*srdblib.User).User_name

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-sum.gtpl"))

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 外部のHTMLテンプレートを読み込む
			tmpl := template.Must(template.ParseFiles("template.gtpl"))
			tmpl.Execute(w, nil)
		})
	*/

	if err := tpl.ExecuteTemplate(w, "graph-sum.gtpl", &hgsinf); err != nil {
		// if err := tpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func HandlerGraphSumData(w http.ResponseWriter, r *http.Request) {

	type Sumdata struct {
		Dtime []string  `json:"dtime"`
		Data1 []float64 `json:"data1"`
		Data2 []float64 `json:"data2"`
	}
	sumdata := Sumdata{}

	/*
		data := struct {
			Dtime []string  `json:"dtime"`
			Data1 []float64 `json:"data1"`
			Data2 []float64 `json:"data2"`
		} {
			// Dtime: []float64{0, 3.0, 3.0, 4.0, 4.0, 5.5, 5.5, 6.0, 6.0},
			Dtime: []string{
				"2025-01-01 00:00",
				"2025-01-01 22:00",
				"2025-01-03 12:00",
				"2025-01-04 00:00",
				"2025-01-05 18:00",
			},
			Data1: []float64{0, 10000, 15000, 22500, 35000},
			Data2: []float64{0, 10000, 5000, 7500, 12500},
		}
	*/

	eventid := r.FormValue("eventid")
	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	perslotinflist, _ := MakePointPerSlot(eventid, roomid)

	sumdata.Dtime = make([]string, len(perslotinflist[0].Perslotlist))
	sumdata.Data1 = make([]float64, len(perslotinflist[0].Perslotlist))
	sumdata.Data2 = make([]float64, len(perslotinflist[0].Perslotlist))

	for i, v := range perslotinflist[0].Perslotlist {
		tm := v.Timestart
		sumdata.Dtime[i] = tm.Format("2006-01-02 15:04")
		tp, _ := strconv.Atoi(strings.Replace(v.Tpoint, ",", "", -1))
		sumdata.Data1[i] = float64(tp)
		sumdata.Data2[i] = float64(v.Ipoint)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&sumdata)
}
