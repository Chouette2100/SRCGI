/*!
Copyright © 2024 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

// 参考	100行未満かつGo標準ライブラリだけで作る掲示板
// https://news.mynavi.jp/techplus/article/gogogo-9/
package ShowroomCGIlib

import (
	"fmt"
	"html"
	"log"
	//	"os"
	"strconv"
	"time"

	//	"encoding/json"

	"html/template"
	"net/http"

	"database/sql"

	//	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

// 投稿の内容
type Logm struct {
	ID    int       //	連番
	Cntw  int       //	1: 不具合、	2: 要望、3: 質問、4: その他、5: お知らせ
	Title string    //	タイトル
	Name  string    //	投稿者名
	Body  string    //	投稿本文
	CTime time.Time //	投稿日時
	Color string    //	表示色
}

type BBS struct {
	Manager string //	コメント表示に使用する色
	Cntr    int    //	1: 不具合、	2: 要望、3: 質問、4: その他、5: お知らせ、9: すべて
	Cntlist [] int    //	ラジオボタンの制御（1: 不具合、	2: 要望、3: 質問、4: その他、5: お知らせ、9: すべて）
	Offset  int    //	表示開始位置
	Limit   int    //	表示投稿数
	Loglist []Logm
}

// 投稿一覧を表示する
func HandlerDispBbs(w http.ResponseWriter, r *http.Request) {

	var bbs BBS

	bbs.Cntlist = []int{1, 2, 3, 4, 5}
	bbs.Cntr = 9

	//      ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	GetUserInf(r)

	bbs.Limit, _ = strconv.Atoi(r.FormValue("limit"))
	if bbs.Limit == 0 {
		bbs.Limit = 10
	}
	bbs.Offset, _ = strconv.Atoi(r.FormValue("offset"))

	action := r.FormValue("action")
	if action == "next" {
		bbs.Offset += bbs.Limit
	} else if action == "prev." {
		bbs.Offset -= bbs.Limit
		if bbs.Offset < 0 {
			bbs.Offset = 0
		}
	} else if action == "再表示(top)" {
		bbs.Offset = 0
	}

	from := r.FormValue("from")
	bbs.Manager = r.FormValue("manager")
	if bbs.Manager == "" {
		bbs.Manager = "black"
	}

	if from == "disp-bbs" {
		/*
			for i, v := range []string{"cnt0", "cnt1", "cnt2", "cnt3", "cnt4"} {
				cntv, _ := strconv.Atoi(r.FormValue(v))
				if cntv > 0 {
					bbs.Cntlist[i] = cntv
				} else {
					bbs.Cntlist[i] = -1
				}
			}
		*/
		bbs.Cntr, _ = strconv.Atoi(r.FormValue("cntr"))
	}

	//      テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"htmlEscapeString": func(s string) string { return html.EscapeString(s) },
		"FormatTime":       func(t time.Time, tfmt string) string { return t.Format(tfmt) },
		"CntToName": func(c int) string {
			cntname := []string{"不具合", "要望", "質問", "その他", "お知らせ", "すべて"}
			return cntname[c]
		},
		"Add": func(n int, m int) int { return n + m },
	}
	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/bbs-1.gtpl", "templates/bbs-2.gtpl", "templates/bbs-3.gtpl"))

	// ログを読み出してHTMLを生成 --- (*7)
	err := loadLogs(&bbs) // データを読み出す
	if err != nil {
		err = fmt.Errorf("loadLogs(): %w", err)
		log.Printf("showHandler(): %s\n", err.Error())
	}

	if err := tpl.ExecuteTemplate(w, "bbs-1.gtpl", bbs); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "bbs-2.gtpl", bbs); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "bbs-3.gtpl", bbs); err != nil {
		log.Println(err)
	}

}

// 投稿を書き込み、投稿の一覧を最初から表示する。
func HandlerWriteBbs(w http.ResponseWriter, r *http.Request) {
	//	r.ParseForm() // フォームを解析 --- (*10)
	var logm Logm

	//      ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	//	GetUserInf(req)

	logm.Title = r.FormValue("title")
	logm.Name = r.FormValue("name")
	logm.Body = r.FormValue("body")
	logm.Cntw, _ = strconv.Atoi(r.FormValue("cntw"))
	logm.Color = r.FormValue("color")

	//	if logm.Name == "" {
	//		logm.Name = "名無し"
	//	}
	logm.CTime = time.Now()

	err := saveLog(&logm)
	if err != nil {
		err = fmt.Errorf("saveLogs(): %w", err)
		log.Printf("writeHandler(): %s\n", err.Error())
	}
	//	http.Redirect(w, r, "/disp-bbs?manager="+logm.Color, http.StatusFound)
	http.Redirect(w, r, "/disp-bbs", http.StatusFound)
}

// ファイルからログファイルの読み込み
func loadLogs(
	bbs *BBS,
) (
	err error,
) {

	nrow := 0
	sqlnr := "select count(*) from bbslog "
	if bbs.Cntr != 9 {
		//	表示するコメント種別が指定されているとき
		sqlnr += " where cnt = ? "
		err = srdblib.Db.QueryRow(sqlnr, bbs.Cntr).Scan(&nrow)
	} else {
		//	すべてのコメントを表示するとき
		err = srdblib.Db.QueryRow(sqlnr).Scan(&nrow)
	}
	if err != nil {
		err = fmt.Errorf("Db.QueryRow(): %w", err)
		log.Printf("savelog():%s\n", err.Error())
		return
	}
	if bbs.Offset >= nrow {
		bbs.Offset -= bbs.Limit
	}

	bbs.Loglist = make([]Logm, 0)

	sqltr := "select id, cnt, title, name, body, ctime, color from bbslog "
	if bbs.Cntr != 9 {
		sqltr += " where cnt = ? "
	}
	sqltr += " order by id desc limit ?, ?;"

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = srdblib.Db.Prepare(sqltr)
	if err != nil {
		err = fmt.Errorf("prepare(): %w", err)
		return
	}

	//	cntlist := make([]interface{}, len(bbs.Cntlist))
	//	for i, v := range bbs.Cntlist {
	//		cntlist[i] = v
	//	}
	if bbs.Cntr != 9 {
		rows, err = stmt.Query(bbs.Cntr, bbs.Offset, bbs.Limit)
	} else {
		rows, err = stmt.Query(bbs.Offset, bbs.Limit)
	}
	if err != nil {
		err = fmt.Errorf("query(): %w", err)
		return
	}

	var logm Logm
	for rows.Next() {
		err = rows.Scan(
			&logm.ID,
			&logm.Cntw,
			&logm.Title,
			&logm.Name,
			&logm.Body,
			&logm.CTime,
			&logm.Color,
		)
		if err != nil {
			err = fmt.Errorf("scan(): %w", err)
			return
		}
		bbs.Loglist = append(bbs.Loglist, logm)
	}

	return

}

// ログファイルの書き込み
func saveLog(logm *Logm) (err error) {

	nrow := 0
	err = srdblib.Db.QueryRow("select max(id) from bbslog").Scan(&nrow)
	if err != nil {
		err = fmt.Errorf("Db.QueryRow(): %w", err)
		log.Printf("savelog():%s\n", err.Error())
		return
	}

	logm.ID = nrow + 1

	sqlip := "insert into bbslog (id, cnt, title, name, body, ctime, color ) values(?,?,?,?,?,?,?)"
	_, srdblib.Dberr = srdblib.Db.Exec(sqlip, logm.ID, logm.Cntw, logm.Title, logm.Name, logm.Body, logm.CTime, logm.Color)
	if srdblib.Dberr != nil {
		err := fmt.Errorf("Db.Exec(sqlip,...): %w", srdblib.Dberr)
		log.Printf("err=[%s]\n", err.Error())
	}
	return
}
