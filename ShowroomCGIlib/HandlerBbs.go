/*!
Copyright © 2024 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

/*
	SRCGI.00AM02	通常とメンテナンスの切り替えを ShowroomCGIlib.Serverconfig.Maintenance で行う。
*/

// 参考	100行未満かつGo標準ライブラリだけで作る掲示板
// https://news.mynavi.jp/techplus/article/gogogo-9/
package ShowroomCGIlib

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	"html/template"
	"net/http"

	"database/sql"

	"github.com/Chouette2100/srdblib/v2"
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
	Ra    string    //	リモートアドレス
	Ua    string    //	ユーザーエージェント
}

type BBS struct {
	Manager string //	コメント表示に使用する色
	Cntr    int    //	1: 不具合、	2: 要望、3: 質問、4: その他、5: お知らせ、9: すべて
	Cntlist []int  //	ラジオボタンの制御（1: 不具合、	2: 要望、3: 質問、4: その他、5: お知らせ、9: すべて）
	Offset  int    //	表示開始位置
	Limit   int    //	表示投稿数
	Nlog    int    //	読みこんだ投稿の数
	Loglist []Logm //	ログメッセージ
}

// リクエストの内容によって投稿を書き込み、あるいは投稿一覧を表示する
func DispBbsHandler(w http.ResponseWriter, r *http.Request) {

	var bbs BBS
	var logm Logm

	bbs.Cntlist = []int{1, 2, 3, 4, 5}

	//      ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	ra, ua, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	//      ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	//	GetUserInf(req)

	prc := r.FormValue("prc")
	if prc == "write" {
		//	書き込みが行われた
		logm.Title = r.FormValue("title")
		logm.Name = r.FormValue("name")
		logm.Body = r.FormValue("body")
		logm.Cntw, _ = strconv.Atoi(r.FormValue("cntw"))
		logm.Color = r.FormValue("color")

		logm.CTime = time.Now()

		raa := strings.Split(ra, ":")
		logm.Ra = raa[0]
		logm.Ua = ua

		err := saveLog(&logm)
		if err != nil {
			err = fmt.Errorf("saveLogs(): %w", err)
			log.Printf("writeHandler(): %s\n", err.Error())
		}
		//	ログ書き込み後のログの表示はジャンルを限定せず表示する。
		bbs.Cntr = 9
	} else {
		//	読み込みが行われた。前回と同じジャンルの投稿を表示する。
		//	URL直打ちのときは前回がないので、actionの処理のとき別に設定する。
		bbs.Cntr, _ = strconv.Atoi(r.FormValue("cntr"))
	}

	//	一度に表示する投稿の数
	/*
		bbs.Limit, _ = strconv.Atoi(r.FormValue("limit"))
		if bbs.Limit == 0 {
			bbs.Limit = 11
		}
	*/
	bbs.Limit = 11

	//	何番目の投稿から表示するか？
	//	上のlimitとあわせ  select ........... limit , offset に用いる。
	bbs.Offset, _ = strconv.Atoi(r.FormValue("offset"))

	action := r.FormValue("action")
	if action == "next" {
		//	次ページを表示する。
		bbs.Offset += bbs.Limit - 1
	} else if action == "prev" {
		//	前ページを表示する。
		bbs.Offset -= bbs.Limit - 1
		if bbs.Offset < 0 {
			bbs.Offset = 0
		}
	} else if action == "再表示(top)" {
		//	投稿を最初（＝投稿順としては最後）から表示する。
		bbs.Offset = 0
	} else {
		//	actionが指定されていないとき＝URL(/bbs-disp)直打ちのときは全ログを表示する。
		bbs.Cntr = 9
	}

	//	メンバー名を取得する。managerに色名を指定すると投稿がその色で表示されるとともに
	//	管理者としてのお知らせの投稿ができる。
	bbs.Manager = r.FormValue("manager")
	if bbs.Manager == "" {
		bbs.Manager = "black"
	}

	//      テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"htmlEscapeString": func(s string) string { return html.EscapeString(s) }, //	必要か（もっとかんたんな方法がないか）確認のこと
		"FormatTime":       func(t time.Time, tfmt string) string { return t.Format(tfmt) },
		"CntToName": func(c int) string {
			cntname := []string{"不具合", "要望", "質問", "その他", "お知らせ", "すべて"}
			return cntname[c]
		},
		"Add": func(n int, m int) int { return n + m },
	}
	// テンプレートをパースする
	var tpl *template.Template
	if !Serverconfig.Maintenance {
		tpl = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/bbs-1_org.gtpl", "templates/bbs-2.gtpl", "templates/bbs-3.gtpl"))
	} else {
		/* Maintenance */
		tpl = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/bbs-1_maint.gtpl", "templates/bbs-2.gtpl", "templates/bbs-3.gtpl"))

	}

	// ログを読み出す
	err := loadLogs(&bbs) // データを読み出す
	if err != nil {
		err = fmt.Errorf("loadLogs(): %w", err)
		log.Printf("showHandler(): %s\n", err.Error())
	}

	bbs.Nlog = len(bbs.Loglist)

	if !Serverconfig.Maintenance {
		if err := tpl.ExecuteTemplate(w, "bbs-1_org.gtpl", bbs); err != nil {
			log.Println(err)
		}
	} else {
		/* Maintenance */
		if err := tpl.ExecuteTemplate(w, "bbs-1_maint.gtpl", bbs); err != nil {
			log.Println(err)
		}

	}
	if err := tpl.ExecuteTemplate(w, "bbs-2.gtpl", bbs); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "bbs-3.gtpl", bbs); err != nil {
		log.Println(err)
	}

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
		bbs.Offset -= bbs.Limit - 1
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
	defer stmt.Close()

	if bbs.Cntr != 9 {
		rows, err = stmt.Query(bbs.Cntr, bbs.Offset, bbs.Limit)
	} else {
		rows, err = stmt.Query(bbs.Offset, bbs.Limit)
	}
	if err != nil {
		err = fmt.Errorf("query(): %w", err)
		return
	}
	defer rows.Close()

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
		if err.Error() == "sql: Scan error on column index 0, name \"max(id)\": converting NULL to int is unsupported" {
			//	DBに1件も投稿がないときmax(id)がNULLとなることへの対策
			log.Printf("savelog():<%s>\n", err.Error())
			nrow = 0
		} else {
			err = fmt.Errorf("Db.QueryRow(): %w", err)
			log.Printf("savelog():%s\n", err.Error())
			return
		}
	}

	//	投稿連番
	logm.ID = nrow + 1

	sqlip := "insert into bbslog (id, cnt, title, name, body, ctime, color, ra, ua ) values(?,?,?,?,?,?,?,?,?)"
	_, srdblib.Dberr = srdblib.Db.Exec(sqlip, logm.ID, logm.Cntw, logm.Title, logm.Name, logm.Body, logm.CTime, logm.Color, logm.Ra, logm.Ua)
	if srdblib.Dberr != nil {
		err := fmt.Errorf("Db.Exec(sqlip,...): %w", srdblib.Dberr)
		log.Printf("err=[%s]\n", err.Error())
	}
	return
}
