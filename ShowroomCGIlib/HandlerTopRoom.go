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
	"strconv"
	"time"

	"html/template"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

type TopRoom struct {
	Room_id       int
	Room_url_key  string
	Room_name     string
	Rank          int
	Genre         string
	Point         int
	Event_id      string
	Event_name    string
	Event_endtime time.Time
}

type Top struct {
	Olim        int
	From        time.Time
	To          time.Time
	ChkGenre	map[string]bool
	TopRoomList []TopRoom
}

func SelectTopRoom(
	client *http.Client,
	olim int,
	fromtime time.Time,
	totime time.Time,
) (
	top *Top,
	err error,
) {

	top = &Top{}
	top.TopRoomList = make([]TopRoom, 0)

	top.Olim = olim
	top.From = fromtime
	top.To = totime

	sqltr := "select p.point, p.`rank`, u.userno,u.genre, e.endtime, u.user_name, p.eventid, e.event_name from points p, user u, event e "
	sqltr += " where p.user_id = u.userno and e.eventid = p.eventid and p.pstatus = 'Conf.'  and e.endtime > ? and e.endtime < ? order by p.point desc limit ?; "

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = srdblib.Db.Prepare(sqltr)
	if err != nil {
		err = fmt.Errorf("prepare(): %w", err)
		return
	}

	rows, err = stmt.Query(fromtime, totime, olim)
	if err != nil {
		err = fmt.Errorf("query(): %w", err)
		return
	}

	var troom TopRoom
	for rows.Next() {
		err = rows.Scan(
			&troom.Point,
			&troom.Rank,
			&troom.Room_id,
			&troom.Genre,
			&troom.Event_endtime,
			&troom.Room_name,
			&troom.Event_id,
			&troom.Event_name,
		)
		if err != nil {
			err = fmt.Errorf("scan(): %w", err)
			return
		}
		top.TopRoomList = append(top.TopRoomList, troom)
	}

	return
}

/*

	SelectTopRoom() の戻り値を表示する。

	Ver. 0.1.0

*/
// http://localhost:8080/toproom で呼び出される
func HandlerTopRoom(
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

	olim, _ := strconv.Atoi(r.FormValue("olim"))
	if olim == 0 {
		olim = 50
	}

	from := r.FormValue("from")
	if from == "" {
		from = "2023-10-01"
	}
	from += " +0900"
	fromtime, err := time.Parse("2006-01-02 -0700", from)
	if err != nil {
		log.Printf("HandlerTopRoom(): time.Parse(): %s", err.Error())
		fromtime = time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local)
	}

	totimelimit := time.Now().Truncate(24 * time.Hour).Add(-48 * time.Hour)

	to := r.FormValue("to")
	if to == "" {
		to = time.Now().Truncate(24 * time.Hour).Add(-48 * time.Hour).Format("2006-01-02")
	}
	to += " +0900"
	totime, err := time.Parse("2006-01-02 -0700", to)
	if err != nil {
		log.Printf("HandlerTopRoom(): time.Parse(): %s", err.Error())
		totime = totimelimit
	}
	if totime.After(totimelimit) {
		totime = totimelimit
	}

	log.Printf("from: %v, to: %v olim: %d\n", fromtime, totime, olim)


	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":      func(i int) string { return humanize.Comma(int64(i)) }, //	3桁ごとに","を入れる関数。
		"FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/toproom.gtpl"))

	top, err := SelectTopRoom(client, olim, fromtime, totime)
	if err != nil {
		err = fmt.Errorf("HandlerTopRoom(): %w", err)
		log.Printf("SelectTopRoom(): %s\n", err.Error())
		return
	}
	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "toproom.gtpl", top); err != nil {
		err = fmt.Errorf("HandlerTopRoom(): %w", err)
		log.Printf("%s\n", err.Error())
	}

}
