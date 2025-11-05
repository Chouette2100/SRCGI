// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"fmt"
	"log"
	"strconv"

	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Chouette2100/srdblib/v2"
)

func GraphSum2Handler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	eventid := r.FormValue("eventid")
	sroomid := r.FormValue("roomid")
	if sroomid == "" {
		log.Printf("GraphSumHandler() roomid is empty\n")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	roomid, _ := strconv.Atoi(sroomid)

	intf, _ := srdblib.Dbmap.Get(srdblib.Event{}, eventid)
	if intf == nil {
		log.Printf("Event ID %s not found\n", eventid)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	hgsinf := HGSinf{
		Title:     "獲得ポイントと累積ポイントの概要",
		Eventid:   eventid,
		Ieventid:  intf.(*srdblib.Event).Ieventid,
		Eventname: intf.(*srdblib.Event).Event_name,
		Period:    intf.(*srdblib.Event).Period,
		Roomid:    roomid,
		Maxpoint:  intf.(*srdblib.Event).Maxpoint,
		// Gscale:    intf.(*srdblib.Event).Gscale,
		Gscale: 100,
	}

	intf, _ = srdblib.Dbmap.Get(srdblib.User{}, roomid)
	hgsinf.Roomname = intf.(*srdblib.User).User_name

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-sum2.gtpl"))

	if err := tpl.ExecuteTemplate(w, "graph-sum2.gtpl", &hgsinf); err != nil {
		// if err := tpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

type Sumdata struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

func GraphSumData1Handler(w http.ResponseWriter, r *http.Request) {

	/*
		data := []struct {
			Timestamp string  `json:"timestamp"`
			Value     float64 `json:"value"`
		}{
			{Timestamp: "2025-01-01 00:00", Value: 0},
			{Timestamp: "2025-01-01 21:00", Value: 0},
			{Timestamp: "2025-01-01 22:00", Value: 10000},
			{Timestamp: "2025-01-03 11:00", Value: 10000},
			{Timestamp: "2025-01-03 12:00", Value: 15000},
			{Timestamp: "2025-01-03 23:00", Value: 15000},
			{Timestamp: "2025-01-04 00:00", Value: 22500},
			{Timestamp: "2025-01-05 17:00", Value: 22500},
			{Timestamp: "2025-01-05 18:00", Value: 35000},
		}
	*/

	eventid := r.FormValue("eventid")
	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	sqlst := "select ts, point from points where eventid = ? and user_id = ? order by ts "
	intc, err := srdblib.Dbmap.Select(srdblib.Points{}, sqlst, eventid, roomid)
	if err != nil {
		err = fmt.Errorf("HandlerGraphSumData1() - select() Error: %w", err)
		w.Write([]byte(err.Error()))
		return
	}

	l := len(intc)
	if l == 0 {
		err = fmt.Errorf("HandlerGraphSumData1() - select() Data not found")
		w.Write([]byte(err.Error()))
		return
	}

	sumdata1 := make([]Sumdata, l)
	for i, v := range intc {
		sumdata1[i].Timestamp = v.(*srdblib.Points).Ts.Format("2006-01-02 15:04")
		sumdata1[i].Value = float64(v.(*srdblib.Points).Point)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sumdata1)
}

func GraphSumData2Handler(w http.ResponseWriter, r *http.Request) {

	/*
		sumdata := []struct {
			Timestamp string  `json:"timestamp"`
			Value     float64 `json:"value"`
		} {
			{Timestamp: "2025-01-01 00:00", Value: 0},
			{Timestamp: "2025-01-01 22:00", Value: 10000},
			{Timestamp: "2025-01-03 12:00", Value: 5000},
			{Timestamp: "2025-01-04 00:00", Value: 7500},
			{Timestamp: "2025-01-05 18:00", Value: 12500},
		} */

	eventid := r.FormValue("eventid")
	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	_, perslotinflist, _ := MakePointPerSlot(eventid, roomid)
	if len(perslotinflist) == 0 {
		err := fmt.Errorf("GraphSumData2Handler() - MakePointPerSlot() Data not found")
		log.Printf("GraphSumData2Handler() - %s\n", err.Error())
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	sumdata2 := make([]Sumdata, len(perslotinflist[0].Perslotlist))

	for i, v := range perslotinflist[0].Perslotlist {
		tm := v.Timestart
		sumdata2[i].Timestamp = tm.Format("2006-01-02 15:04")
		sumdata2[i].Value = float64(v.Ipoint)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sumdata2)
}
