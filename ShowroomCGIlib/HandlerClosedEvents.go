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
	"strings"
	"time"

	//	"io" //　ログ出力設定用。必要に応じて。
	//	"sort" //	ソート用。必要に応じて。

	"html/template"
	"net/http"

	"github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
	//	"github.com/Chouette2100/srapi/v2"
)

/*
終了イベント一覧を作るためのハンドラー

Ver. 0.1.0
*/
func ClosedEventsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	//      テーブルは"w"で始まるものを操作の対象とする。
	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":          func(i int) string { return humanize.Comma(int64(i)) },                          //	3桁ごとに","を入れる関数。
		"UnixTimeToStr":  func(i int64) string { return time.Unix(int64(i), 0).Format("01-02 15:04") },    //	UnixTimeを月日時分に変換する関数。
		"UnixTimeToStrY": func(i int64) string { return time.Unix(int64(i), 0).Format("06-01-02 15:04") }, //	UnixTimeを年月日時分に変換する関数。

		"TimeToString":  func(t time.Time) string { return t.Format("01-02 15:04") },
		"TimeToStringY": func(t time.Time) string { return t.Format("06-01-02 15:04") },
		"DelBlockID": func(eid string) string {
			eia := strings.Split(eid, "?")
			if len(eia) == 2 {
				return eia[0]
			} else {
				return eid
			}
		},
		"IsTempID": func(s string) bool { return strings.HasPrefix(s, "@@@@") },
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

	// slimit := r.FormValue("limit")
	// if slimit == "" {
	// 	top.Limit = 51
	// } else {
	// 	top.Limit, _ = strconv.Atoi(slimit)
	// }
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
	switch action {
	case "next":
		//	次ページを表示する。
		top.Offset += top.Limit - 1
	case "prev":
		//	前ページを表示する。
		top.Offset -= top.Limit - 1
		if top.Offset < 0 {
			top.Offset = 0
		}
	case "top":
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
			top.Eventinflist, err = SelectEventinflistFromEventByRoom(cond, top.Mode, top.Userno, &top.Limit, top.Offset)
			if err != nil {
				err = fmt.Errorf("MakeListOfPoints(): %w", err)
				log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
				top.ErrMsg = err.Error()
			}
			ri, _ := SelectRoomInf(top.Userno)
			top.Keywordrm = ri.Name
			top.Roomlist = &[]Room{{Userno: ri.Userno, User_name: "(" + strconv.Itoa(ri.Userno) + ")" + ri.Name}}
		}
	case 5:
		//	ルームIDによる絞り込み
		top.Eventinflist, err = SelectEventinflistFromEventByRoom(cond, top.Mode, top.Userno, &top.Limit, top.Offset)
		if err != nil {
			err = fmt.Errorf("MakeListOfPoints(): %w", err)
			log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
			top.ErrMsg = err.Error()
		}
		ri, _ := SelectRoomInf(top.Userno)
		top.Keywordrm = ri.Name
		top.Roomlist = &[]Room{{Userno: ri.Userno, User_name: "(" + strconv.Itoa(ri.Userno) + ")" + ri.Name}}
	default:
	}

	top.Totalcount = len(top.Eventinflist)

	/*
		// 参照回数の多いイベントを取得する
		var emap map[string]int
		emap, err = srdblib.GetFeaturedEvents("closed", 72, 16, 6)
		if err != nil {
			err = fmt.Errorf("GetFeaturedEvents(): %w", err)
			log.Printf("%s\n", err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		for i, v := range top.Eventinflist {
			if _, ok := emap[v.Event_ID]; ok {
				top.Eventinflist[i].Aclr += 2
			}
			// else {
			// 	top.Eventinflist[i].Aclr = 0
			// }
		}
	*/

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

func SelectRoomInf(
	userno int,
) (
	roominf RoomInfo,
	status int,
) {

	status = 0

	sql := "select distinct u.userno, userid, user_name, longname, shortname, genre, nrank, prank, level, followers, fans, fans_lst, e.istarget,e.graph, e.color, e.iscntrbpoints, e.point "
	sql += " from user u join eventuser e "
	//	sql += " where u.userno = e.userno and u.userno = " + fmt.Sprintf("%d", userno)
	sql += " where u.userno = e.userno and u.userno = ? "

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("SelectRoomInf() Prepare() err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(userno).Scan(&roominf.Userno,
		&roominf.Account,
		&roominf.Name,
		&roominf.Longname,
		&roominf.Shortname,
		&roominf.Genre,
		&roominf.Nrank,
		&roominf.Prank,
		&roominf.Level,
		&roominf.Followers,
		&roominf.Fans,
		&roominf.Fans_lst,
		&roominf.Istarget,
		&roominf.Graph,
		&roominf.Color,
		&roominf.Iscntrbpoint,
		&roominf.Point,
	)
	if err != nil {
		log.Printf("SelectRoomInf() Query() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	if roominf.Istarget == "Y" {
		roominf.Istarget = "Checked"
	} else {
		roominf.Istarget = ""
	}
	if roominf.Graph == "Y" {
		roominf.Graph = "Checked"
	} else {
		roominf.Graph = ""
	}
	if roominf.Iscntrbpoint == "Y" {
		roominf.Iscntrbpoint = "Checked"
	} else {
		roominf.Iscntrbpoint = ""
	}
	roominf.Slevel = humanize.Comma(int64(roominf.Level))
	roominf.Sfollowers = humanize.Comma(int64(roominf.Followers))
	roominf.Spoint = humanize.Comma(int64(roominf.Point))
	roominf.Name = strings.ReplaceAll(roominf.Name, "'", "’")

	return
}
