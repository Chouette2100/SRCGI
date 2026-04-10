/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

package ShowroomCGIlib

import (
	"fmt"
	"log"
	"strconv"
	"time"

	// "os"
	//	"io" //　ログ出力設定用。必要に応じて。
	//	"sort" //	ソート用。必要に応じて。

	"html/template"
	"net/http"

	// "github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v3"
	//	"github.com/Chouette2100/srdblib/v3"
	//	"github.com/Chouette2100/srapi/v2"
)

/*

	開催中のイベント一覧を作るためのハンドラー

	Ver. 0.1.0

*/

/*
	type T008top struct {
		TimeNow    int64
		Totalcount int
		ErrMsg     string
		Eventlist  []srapi.Event
	}
*/
type Room struct {
	Userno    int
	User_name string
}
type T999Dtop struct {
	TimeNow      int64
	Totalcount   int
	ErrMsg       string
	Mode         int    // 0: すべて、 1: データ取得中のものに限定
	Path         int    //	どの検索方法が使われているか？（詳細は HandlerCloesedEvnets()および関連関数を参照）
	Keywordev    string //	検索文字列:イベント名
	Keywordrm    string //	検索文字列:ルーム名
	Kwevid       string //	検索文字列:イベントID
	Userno       int    //	絞り込み対象のルームID
	Limit        int    //	データ取得数
	Offset       int    //	データ取得開始位置
	Eventinflist []exsrapi.Event_Inf
	Roomlist     *[]Room
	// Turnstile導入用(1) ------------------------
	Action           string
	TurnstileSiteKey string
	TurnstileError   string
	RequestID        string
}

// TurnstileChallengeDataインターフェースの実装
func (h *T999Dtop) SetTurnstileInfo(siteKey string, errorMsg string) {
	h.TurnstileSiteKey = siteKey
	h.TurnstileError = errorMsg
}

func (h *T999Dtop) GetTemplatePath() string {
	return "templates/currentevents.gtpl"
}

func (h *T999Dtop) GetTemplateName() string {
	return "currentevents.gtpl"
}

func (h *T999Dtop) GetFuncMap() *template.FuncMap {
	return &CommonFuncMap
}

// -------------------------------------------

// 開催中イベントの一覧を表示するハンドラー
func CurrentEventsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	/*
		client, cookiejar, err := exsrapi.CreateNewClient("")
		if err != nil {
			log.Printf("exsrapi.CeateNewClient(): %s", err.Error())
			return //	エラーがあれば、ここで終了
		}
		defer cookiejar.Save()
	*/

	//      テーブルは"w"で始まるものを操作の対象とする。
	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	/*
		//	テンプレートで使用する関数を定義する
		funcMap := template.FuncMap{
			"Comma":         func(i int) string { return humanize.Comma(int64(i)) },                       //	3桁ごとに","を入れる関数。
			"UnixTimeToStr": func(i int64) string { return time.Unix(int64(i), 0).Format("01-02 15:04") }, //	UnixTimeを年月日時分に変換する関数。
			"TimeToString":  func(t time.Time) string { return t.Format("01-02 15:04") },
			"Divide": func(a, b int) int {
				if b == 0 {
					return 0 // ゼロ除算を避ける
				}
				return a / b
			},
			"Mod": func(a, b int) int {
				if b == 0 {
					return 0 // ゼロ除算を避ける
				}
				return a % b
			},
		}
	*/

	// テンプレートに埋め込むデータ（ポイントやランク）を作成する
	top := new(T999Dtop)
	top.TimeNow = time.Now().Unix()
	top.Mode, _ = strconv.Atoi(r.FormValue("mode"))

	top.Limit = 50
	soffset := r.FormValue("offset")
	if soffset == "" {
		top.Offset = 0
	} else {
		top.Offset, _ = strconv.Atoi(soffset)
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
		result, err := Dbmap0.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", top.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
	} else {
		result, err := Dbmap0.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", top.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
		result, err = Dbmap0.Exec(
			"DELETE FROM accesslog WHERE requestid = ?", lastrequestid)
		log.Printf("  delete from accesslog where lastrequestid = %s result=%+v, err=%+v\n",
			lastrequestid, result, err)
	}
	// -------------------------------------------

	var err error
	cond := 0 // 抽出条件	-1:終了したイベント、0: 開催中のイベント、1: 開催予定のイベント
	top.Eventinflist, err = SelectEventinflistFromEvent(cond, top.Mode, "", "", 0, 0)
	if err != nil {
		err = fmt.Errorf("MakeListOfPoints(): %w", err)
		log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
		top.ErrMsg = err.Error()
	}
	top.Totalcount = len(top.Eventinflist)

	// 参照回数の多いイベントを取得する
	var emap map[string]int
	emap, err = srdblib.GetFeaturedEvents(Dbmap0, "current", 48, 18, 14)
	if err != nil {
		err = fmt.Errorf("GetFeaturedEvents(): %w", err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	for i, v := range top.Eventinflist {
		if _, ok := emap[v.Event_ID]; ok {
			top.Eventinflist[i].Highlighted = 1
		} else {
			top.Eventinflist[i].Highlighted = 0
		}
	}
	err = FindHistoricalData(&top.Eventinflist)
	if err != nil {
		err = fmt.Errorf("FindHistoricalData(): %w", err)
		log.Printf("FindHistoricalData() returned error %s\n", err.Error())
		top.ErrMsg = err.Error()
		return
	}

	//	ソートが必要ないときは次の行とimportの"sort"をコメントアウトする。
	//	無名関数のリターン値でソート条件を変更できます。
	/*
		sort.Slice(top.Eventinflist, func(i, j int) bool {
			if top.Eventinflist[i].End_time.After(top.Eventinflist[j].End_time) {
				return true
			} else if top.Eventinflist[i].End_time == top.Eventinflist[j].End_time {
				return top.Eventinflist[i].Start_time.Before(top.Eventinflist[j].Start_time)
			} else {
				return false
			}
		})
	*/

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(CommonFuncMap).ParseFiles("templates/currentevents.gtpl"))

	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "currentevents.gtpl", top); err != nil {
		log.Printf("tpl.ExecuteTemplate() returned error: %s\n", err.Error())
	}

}
