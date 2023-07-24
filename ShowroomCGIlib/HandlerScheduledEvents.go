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
	//	"io" //　ログ出力設定用。必要に応じて。
	//	"sort" //	ソート用。必要に応じて。

	"html/template"
	"net/http"

	"github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
	//	"github.com/Chouette2100/srapi"
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
type T999Dtop struct {
	TimeNow      int64
	Totalcount   int
	ErrMsg       string
	Mode			int
	Keyword			string	//	検索文字列
	Eventinflist []exsrapi.Event_Inf
}
*/

// "/T999Dtop"に対するハンドラー
// http://localhost:8080/T999Dtop で呼び出される
func HandlerScheduledEvents(
	w http.ResponseWriter,
	r *http.Request,
) {

	GetUserInf(r)

	/*
		client, cookiejar, err := exsrapi.CreateNewClient("")
		if err != nil {
			log.Printf("exsrapi.CeateNewClient(): %s", err.Error())
			return //	エラーがあれば、ここで終了
		}
		defer cookiejar.Save()
	*/

	//      テーブルは"w"で始まるものを操作の対象とする。
	srdblib.Tevent = "wevent"
	srdblib.Teventuser = "weventuser"
	srdblib.Tuser = "wuser"
	srdblib.Tuserhistory = "wuserhistory"


	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":         func(i int) string { return humanize.Comma(int64(i)) },                       //	3桁ごとに","を入れる関数。
		"UnixTimeToStr": func(i int64) string { return time.Unix(int64(i), 0).Format("01-02 15:04") }, //	UnixTimeを年月日時分に変換する関数。
		"TimeToString":  func(t time.Time) string { return t.Format("01-02 15:04") },
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/scheduled-events.gtpl"))

	// テンプレートに埋め込むデータ（ポイントやランク）を作成する
	top := new(T999Dtop)
	top.TimeNow = time.Now().Unix()
	top.Mode, _ = strconv.Atoi(r.FormValue("mode"))

	var err error
	cond := 1	//	開催予定
	top.Eventinflist, err = SelectEventinflistFromEvent(cond, top.Mode, "")
	if err != nil {
		err = fmt.Errorf("MakeListOfPoints(): %w", err)
		log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
		top.ErrMsg = err.Error()
	}
	top.Totalcount = len(top.Eventinflist)

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

	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "scheduled-events.gtpl", top); err != nil {
		log.Printf("tpl.ExecuteTemplate() returned error: %s\n", err.Error())
	}

}
