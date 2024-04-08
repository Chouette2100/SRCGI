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

	終了イベント一覧を作るためのハンドラー

	Ver. 0.1.0

*/
func HandlerClosedEvents(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, _, isallow := GetUserInf(r)
	if ! isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}



	//      テーブルは"w"で始まるものを操作の対象とする。
	srdblib.Tevent = "wevent"
	srdblib.Teventuser = "weventuser"
	srdblib.Tuser = "wuser"
	srdblib.Tuserhistory = "wuserhistory"

	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":          func(i int) string { return humanize.Comma(int64(i)) },                          //	3桁ごとに","を入れる関数。
		"UnixTimeToStr":  func(i int64) string { return time.Unix(int64(i), 0).Format("01-02 15:04") },    //	UnixTimeを月日時分に変換する関数。
		"UnixTimeToStrY": func(i int64) string { return time.Unix(int64(i), 0).Format("06-01-02 15:04") }, //	UnixTimeを年月日時分に変換する関数。

		"TimeToString":  func(t time.Time) string { return t.Format("01-02 15:04") },
		"TimeToStringY": func(t time.Time) string { return t.Format("06-01-02 15:04") },
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/closedevents.gtpl"))

	// テンプレートに埋め込むデータ（ポイントやランク）を作成する
	top := new(T999Dtop)
	top.TimeNow = time.Now().Unix()
	top.Mode, _ = strconv.Atoi(r.FormValue("mode")) // 0: すべて、 1: データ取得中のものに限定
	top.Keywordev = r.FormValue("keywordev")
	top.Keywordrm = r.FormValue("keywordrm")
	top.Kwevid = r.FormValue("kwevid")
	top.Userno, _ = strconv.Atoi(r.FormValue("userno"))

	slimit := r.FormValue("limit")
	if slimit == "" {
		top.Limit = 51
	} else {
		top.Limit, _ = strconv.Atoi(slimit)
	}
	soffset := r.FormValue("offset")
	if soffset == "" {
		top.Offset = 0
	} else {
		top.Offset, _ = strconv.Atoi(soffset)
	}

	//	ページ操作	
	action := r.FormValue("action")
	if action == "next" {
		//	次ページを表示する。
		top.Offset += top.Limit - 1
	} else if action == "prev" {
		//	前ページを表示する。
		top.Offset -= top.Limit - 1
		if top.Offset < 0 {
			top.Offset = 0
		}
	} else if action == "top" {
		//	最初から表示する。
		top.Offset = 0
	}

	top.Path, _ = strconv.Atoi(r.FormValue("path"))
	/*
		0. 最初のパス
		1. イベント名で絞り込む
		2. イベントID(Event_url_key)で絞り込む
		3. ルーム名で絞り込む(ルーム名の入力)
		4. ルーム名で絞り込む(ルーム名の選択)
		5. ユーザ番号で選択する
	*/

	top.Roomlist = &[]Room{}

	var err error

	cond := -1 // 抽出条件	-1:終了したイベント、0: 開催中のイベント、1: 開催予定のイベント
	switch top.Path {
	case 0:
		top.Eventinflist, err = SelectEventinflistFromEvent(cond, top.Mode, "", "", top.Limit, top.Offset)
		if err != nil {
			err = fmt.Errorf("SelectEventinflistFromEvent(): %w", err)
			log.Printf("SelectEventinflistFromEvent() returned error %s\n", err.Error())
			top.ErrMsg = err.Error()
		}
	case 1, 2:
		if top.Path == 1 {
			top.Eventinflist, err = SelectEventinflistFromEvent(cond, top.Mode, top.Keywordev, "", top.Limit, top.Offset)
		} else {
			top.Eventinflist, err = SelectEventinflistFromEvent(cond, top.Mode, "", top.Kwevid, top.Limit, top.Offset)
		}
		if err != nil {
			err = fmt.Errorf("SelectEventinflistFromEvent(): %w", err)
			log.Printf("SelectEventinflistFromEvent() returned error %s\n", err.Error())
			top.ErrMsg = err.Error()
		}
	case 3, 4:
		if top.Keywordrm != "" {
			//	ルーム名による絞り込み、ルームの候補リストを作成する。
			top.Roomlist, err = SelectUsernoAndName(top.Keywordrm, 50, 0)
			if err != nil {
				err = fmt.Errorf("SelectUsernoAndName(): %w", err)
				log.Printf("SelectUsernoAndName() returned error %s\n", err.Error())
				top.ErrMsg = err.Error()
			}
		}
		if top.Path == 4 && top.Userno != 0 {
			//	ルーム名による絞り込み(ルームIDに変換後)
			top.Eventinflist, err = SelectEventinflistFromEventByRoom(cond, top.Mode, top.Userno, top.Limit, top.Offset)
			if err != nil {
				err = fmt.Errorf("MakeListOfPoints(): %w", err)
				log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
				top.ErrMsg = err.Error()
			}
			ri, _ := SelectRoomInf(top.Userno)
			top.Keywordrm = ri.Name
			top.Roomlist = & []Room { {Userno: ri.Userno, User_name: "(" + strconv.Itoa(ri.Userno) + ")" + ri.Name} }
		}
	case 5:
			//	ルームIDによる絞り込み
			top.Eventinflist, err = SelectEventinflistFromEventByRoom(cond, top.Mode, top.Userno, top.Limit, top.Offset)
			if err != nil {
				err = fmt.Errorf("MakeListOfPoints(): %w", err)
				log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
				top.ErrMsg = err.Error()
			}
			ri, _ := SelectRoomInf(top.Userno)
			top.Keywordrm = ri.Name
			top.Roomlist = & []Room { {Userno: ri.Userno, User_name: "(" + strconv.Itoa(ri.Userno) + ")" + ri.Name} }
	default:
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
	if err = tpl.ExecuteTemplate(w, "closedevents.gtpl", top); err != nil {
		log.Printf("tpl.ExecuteTemplate() returned error: %s\n", err.Error())
	}

}
