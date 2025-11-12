// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	"fmt"
	// "html"
	"log"

	//	"math/rand"
	// "sort"
	// "strconv"
	//	"strings"
	// "time"
	//	"os"

	// "runtime"

	// "encoding/json"

	"html/template"
	"net/http"

	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// Turnstile導入用(1) ------------------------
type EventTopInf struct {
	*exsrapi.Event_Inf
	TurnstileSiteKey string
	TurnstileError   string
	RequestID        string
}

// TurnstileChallengeDataインターフェースの実装
func (h *EventTopInf) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *EventTopInf) GetTemplatePath() string {
	return "templates/eventtop.gtpl"
}

func (h *EventTopInf) GetTemplateName() string {
	return "eventtop.gtpl"
}

func (h *EventTopInf) GetFuncMap() *template.FuncMap {
	return nil
}

// -------------------------------------------

// 入力フォーム画面
func EventTopHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles(
	//		"templates/top.gtpl",
	//		"templates/top0.gtpl",
	//		"templates/bbs-2.gtpl",
	//		"templates/top1.gtpl",
	//		"templates/top2.gtpl",
	//	))

	eventid := r.FormValue("eventid")
	// suserno := r.FormValue("userno")
	// if suserno == "" {
	// 	suserno = "0"
	// }
	// userno, _ := strconv.Atoi(suserno)
	// log.Printf("      eventid=%s userno=%d\n", eventid, userno)
	log.Printf("      eventid=%s\n", eventid)

	//	eventinf, _ := SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		return
	}
	// Event_inf = *eventinf

	top := &EventTopInf{
		Event_Inf: eventinf,
	}

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

	result, tsErr := CheckTurnstileWithSession(w, r, top)
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

	tpl := template.Must(template.ParseFiles(
		"templates/eventtop.gtpl",
	))

	if err := tpl.ExecuteTemplate(w, "eventtop.gtpl", top); err != nil {
		log.Println(err)
	}

}
