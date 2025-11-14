// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"html/template"
	"net/http"

	"github.com/Chouette2100/srdblib/v2"
)

// Turnstile導入用(1) ------------------------

type GraphSumInf struct {
	HGSinf
}

// TurnstileChallengeDataインターフェースの実装
func (h *GraphSumInf) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *GraphSumInf) GetTemplatePath() string {
	return "templates/graph-sum.gtpl"
}

func (h *GraphSumInf) GetTemplateName() string {
	return "graph-sum.gtpl"
}

func (h *GraphSumInf) GetFuncMap() *template.FuncMap {
	return nil
}

// -------------------------------------------

func GraphSumHandler(w http.ResponseWriter, r *http.Request) {

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
	top := GraphSumInf{
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
	top.Roomname = intf.(*srdblib.User).User_name

	// Turnstile導入用(2) ------------------------
	// Turnstile検証（セッション管理込み）
	// Turnstile検証要求後の状態を管理する
	lastrequestid := ""
	requestid := r.FormValue("requestid")
	if requestid != "" {
		// 最初のパス、検証のための場合と、すでにクッキーを持っている場合と両方ある
		lastrequestid = requestid
	}
	top.RequestID = r.Context().Value("requestid").(string)
	// ------

	result, tsErr := CheckTurnstileWithSession(w, r, &top)
	if result != TurnstileOK {
		// チャレンジページまたはエラーページが表示済みなので終了
		if tsErr != nil {
			log.Printf("Turnstile check error: %v\n", tsErr)
		}
		return
	}

	log.Printf(" hcntbinf.RequestID = %s, lastrequestid = %s\n", top.RequestID, lastrequestid)
	if lastrequestid == "" {
		result, err := srdblib.Dbmap.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", top.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
	} else {
		//srdblib.Dbmap.Exec("DELETE FROM accesslog WHERE requestid = ?", requestid)
		result, err := srdblib.Dbmap.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", top.RequestID)
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
	tpl := template.Must(template.ParseFiles("templates/graph-sum.gtpl"))

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 外部のHTMLテンプレートを読み込む
			tmpl := template.Must(template.ParseFiles("template.gtpl"))
			tmpl.Execute(w, nil)
		})
	*/

	if err := tpl.ExecuteTemplate(w, "graph-sum.gtpl", &top); err != nil {
		// if err := tpl.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func GraphSumDataHandler(w http.ResponseWriter, r *http.Request) {

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

	time.Sleep(1 * time.Second) // INSERTに時間がかかる場合があるため待機
	requestid := r.Context().Value("requestid").(string)
	result, err := srdblib.Dbmap.Exec(
		"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", requestid)
	log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)

	eventid := r.FormValue("eventid")
	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	_, perslotinflist, _ := MakePointPerSlot(eventid, roomid)
	if len(perslotinflist) == 0 {
		err := fmt.Errorf("GraphSumDataHandler() - MakePointPerSlot() Data not found")
		log.Printf("GraphSumDataHandler() - %s\n", err.Error())
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

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
