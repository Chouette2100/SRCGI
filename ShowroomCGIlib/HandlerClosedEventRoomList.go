/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

package ShowroomCGIlib

import (
	//	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	//	"strings"
	"time"

	"html/template"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
	"github.com/Chouette2100/srapi"
)

type Erl struct {
	Eventid     int
	Eventname   string
	Eventurl    string
	Ib          int
	Ie          int
	Roomlistinf *srapi.RoomListInf
	Msg         string
	Eventlist   []srapi.Event
}

func SelectLastdataFromWeventuser(
	client *http.Client,
	eventurlkey string,
	ib int,
	ie int,
) (
	roomlistinf *srapi.RoomListInf,
	err error,
) {
	roomlistinf = new(srapi.RoomListInf)
	roomlistinf.RoomList = make([]srapi.Room, 0)

	sqlswe := "SELECT  we.userno, wu.userid, wu.longname, we.vld, we.point "
	sqlswe += " FROM weventuser we JOIN wuser wu ON we.userno = wu.userno "
	sqlswe += " WHERE we.eventid = ? ORDER by we.vld "

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = srdblib.Db.Prepare(sqlswe)
	if err != nil  {
		err = fmt.Errorf("prepare: %s", err.Error())
		return
	}

	rows, err = stmt.Query(eventurlkey)
	if err != nil {
		err = fmt.Errorf("query: %s", err.Error())
		return
	}

	var room srapi.Room
	for rows.Next() {
		err = rows.Scan(
			&room.Room_id,
			&room.Room_url_key,
			&room.Room_name,
			&room.Rank,
			&room.Point,
		)
		if err != nil {
			err = fmt.Errorf("scan: %s", err.Error())
			return
		}
		roomlistinf.RoomList = append(roomlistinf.RoomList, room)
	}

	return
}

/*

	ApiEventRoomList() の戻り値を表示する。

	Ver. 0.1.0

*/
// "/ApiEventRoomList()"に対するハンドラー
// http://localhost:8080/apieventroomlist で呼び出される
func HandlerClosedEventRoomList(
	w http.ResponseWriter,
	r *http.Request,
) {

	GetUserInf(r)

	//	cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("XXXXXX")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//	すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":          func(i int) string { return humanize.Comma(int64(i)) },                           //	3桁ごとに","を入れる関数。
		"UnixtimeToTime": func(i int64, tfmt string) string { return time.Unix(int64(i), 0).Format(tfmt) }, //	UnixTimeを時分に変換する関数。
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/closedeventroomlist.gtpl", "templates/footer.gtpl"))

	var erl Erl

	seventid := r.FormValue("eventid")
	eventurlkey := r.FormValue("eventurlkey")
	if seventid == "" {
		/*
			err = errors.New("eventid が設定されていません。URLのあとに\"?eventid=.....\"を追加してください。<br>あるいは「開催中イベント一覧表」から参加者一覧が必要なイベントを指定してください。")
			erl.Msg = err.Error()
			log.Printf("%s\n", erl.Msg)
		*/
		erl.Eventid = 0
		erl.Eventlist, err = srapi.MakeEventListByApi(client)
		if err != nil {
			err = fmt.Errorf("MakeListOfPoints(): %w", err)
			log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
			erl.Msg = err.Error()
		}
		//	erl.Totalcount = len(top.Eventlist)

		//	ソートが必要ないときは次の行とimportの"sort"をコメントアウトする。
		//	無名関数のリターン値でソート条件を変更できます。
		//	ここではエベント終了日時が近い順にソートしています。
		sort.Slice(erl.Eventlist, func(i, j int) bool { return erl.Eventlist[i].Ended_at < erl.Eventlist[j].Ended_at })

	} else {

		erl.Eventurl = eventurlkey
		erl.Eventid, err = strconv.Atoi(seventid)
		if err != nil {
			err = fmt.Errorf("HandlerEventRoomList(): %w", err)
			erl.Msg = err.Error()
			log.Printf("%s\n", erl.Msg)
		} else {

			sib := r.FormValue("ib")
			erl.Ib, err = strconv.Atoi(sib)
			if err != nil {
				erl.Ib = 1
			}

			sie := r.FormValue("ie")
			erl.Ie, err = strconv.Atoi(sie)
			if err != nil {
				erl.Ie = 10
			}

			if erl.Ie < erl.Ib {
				erl.Ie = erl.Ib
			}

			erl.Roomlistinf, err = SelectLastdataFromWeventuser(client, erl.Eventurl, erl.Ib, erl.Ie)
			if err != nil {
				err = fmt.Errorf("SelectLastdataFromWeventuser(): %w", err)
				erl.Msg = err.Error()
			}
		}
	}
	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "closedeventroomlist.gtpl", erl); err != nil {
		err = fmt.Errorf("HandlerEventRoomList(): %w", err)
		log.Printf("%s\n", err.Error())
	}

}
