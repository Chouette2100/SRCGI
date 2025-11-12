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
	// Turnstile導入用(1) ------------------------
	TurnstileSiteKey string
	TurnstileError   string
	RequestID        string
}

type GraphSum2Inf struct {
	HGSinf
}

// TurnstileChallengeDataインターフェースの実装
func (h *GraphSum2Inf) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *GraphSum2Inf) GetTemplatePath() string {
	return "templates/graph-sum2.gtpl"
}

func (h *GraphSum2Inf) GetTemplateName() string {
	return "graph-sum2.gtpl"
}

func (h *GraphSum2Inf) GetFuncMap() *template.FuncMap {
	return nil
}

// -------------------------------------------

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
	graphSum2Inf := GraphSum2Inf{
		HGSinf: HGSinf{
			Title:     "獲得ポイントと累積ポイントの概要",
			Eventid:   eventid,
			Ieventid:  intf.(*srdblib.Event).Ieventid,
			Eventname: intf.(*srdblib.Event).Event_name,
			Period:    intf.(*srdblib.Event).Period,
			Roomid:    roomid,
			Maxpoint:  intf.(*srdblib.Event).Maxpoint,
			// Gscale:    intf.(*srdblib.Event).Gscale,
			Gscale: 100,
		},
	}

	intf, _ = srdblib.Dbmap.Get(srdblib.User{}, roomid)
	graphSum2Inf.Roomname = intf.(*srdblib.User).User_name

	// Turnstile導入用(2) ------------------------
	// Turnstile検証（セッション管理込み）
	// Turnstile検証要求後の状態を管理する
	lastrequestid := ""
	requestid := r.FormValue("requestid")
	if requestid != "" {
		// 最初のパス、検証のための場合と、すでにクッキーを持っている場合と両方ある
		lastrequestid = requestid
	}
	graphSum2Inf.RequestID = r.Context().Value("requestid").(string)
	// ------

	result, tsErr := CheckTurnstileWithSession(w, r, &graphSum2Inf)
	if result != TurnstileOK {
		// チャレンジページまたはエラーページが表示済みなので終了
		if tsErr != nil {
			log.Printf("Turnstile check error: %v\n", tsErr)
		}
		return
	}

	log.Printf(" hcntbinf.RequestID = %s, lastrequestid = %s\n", graphSum2Inf.RequestID, lastrequestid)
	if lastrequestid == "" {
		result, err := srdblib.Dbmap.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", graphSum2Inf.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
	} else {
		//srdblib.Dbmap.Exec("DELETE FROM accesslog WHERE requestid = ?", requestid)
		result, err := srdblib.Dbmap.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", graphSum2Inf.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
		// result, err = srdblib.Dbmap.Exec(
		//      "UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", lastrequestid)
		// log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
		result, err = srdblib.Dbmap.Exec(
			"DELETE FROM accesslog WHERE requestid = ?", lastrequestid)
		log.Printf("  delete from accesslog where lastrequestid = %s result=%+v, err=%+v\n",
			lastrequestid, result, err)
	}
	// -------------------------------------------

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-sum2.gtpl"))

	if err := tpl.ExecuteTemplate(w, "graph-sum2.gtpl", &graphSum2Inf); err != nil {
		// if err := tpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

type Sumdata struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

func GraphSumData1Handler(w http.ResponseWriter, r *http.Request) {

	// Turnstile検証: セッションクッキーの検証
	sessionValid, newCookie, sessionErr := VerifyTurnstileSessionCookie(r)
	if !sessionValid {
		// 検証失敗時はエラーをJSONで返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := "Turnstile verification required"
		if sessionErr != nil {
			errorMsg = fmt.Sprintf("Turnstile verification failed: %v", sessionErr)
			log.Printf("GraphSumData1Handler() Turnstile verification error: %v\n", sessionErr)
		}
		json.NewEncoder(w).Encode(map[string]string{
			"error": errorMsg,
		})
		return
	}

	// セッションクッキーを更新
	if newCookie != nil {
		http.SetCookie(w, newCookie)
	}

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

	// Turnstile検証: セッションクッキーの検証
	sessionValid, newCookie, sessionErr := VerifyTurnstileSessionCookie(r)
	if !sessionValid {
		// 検証失敗時はエラーをJSONで返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := "Turnstile verification required"
		if sessionErr != nil {
			errorMsg = fmt.Sprintf("Turnstile verification failed: %v", sessionErr)
			log.Printf("GraphSumData2Handler() Turnstile verification error: %v\n", sessionErr)
		}
		json.NewEncoder(w).Encode(map[string]string{
			"error": errorMsg,
		})
		return
	}

	// セッションクッキーを更新
	if newCookie != nil {
		http.SetCookie(w, newCookie)
	}

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
