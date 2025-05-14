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
	"strings"
	"time"

	"html/template"
	"net/http"

	//	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type ShowRank struct {
	HdrMsg   string
	Userlist *[]srdblib.User
	ErrMsg   string
}

// SHOWランク上位ルームを抽出する
func SelectShowRank(
	client *http.Client,
	limit int,
) (
	userlist *[]srdblib.User,
	err error,
) {

	userlist = new([]srdblib.User)

	// sqltr := " select * from user where irank between 0 and ? order by irank "
	// sqltr := " select * from user where irank between 0 and ? and ts > ? order by irank "
	sqltr := " select * from user where irank between 0 and ? and ts > ? and fanpower > 0 order by irank "

	var ul []interface{}
	ul, err = srdblib.Dbmap.Select(srdblib.User{}, sqltr, limit, time.Now().Add(-time.Hour*25))
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(): %w", err)
		return
	}

	for _, v := range ul {
		*userlist = append(*userlist, *(v.(*srdblib.User)))
	}

	return
}

/*

	HandlerShowRank()
		SHOWランク上位配信者を表示する

	Ver. 0.1.0

*/
// http://localhost:8080/showrank で呼び出される
func HandlerShowRank(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	showrank := &ShowRank{}

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
		"Comma":      func(i int) string { return humanize.Comma(int64(i)) }, //	3桁ごとに","を入れる関数。
		"FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
		"Add":        func(a int, b int) int { return a + b },
		"Showrank": func(s string) string {
			sa := strings.Split(s, " | ")
			return sa[len(sa)-1]
		},
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/showrank.gtpl"))

	//	showrank.Userlist, err = SelectShowRank(client, 260000000)	//	SS-5〜A-5
	//	showrank.Userlist, err = SelectShowRank(client, 300000000)	//	SS-5〜A-1
	var user1 srdblib.User
	// FIXME: irank != 0 の条件が必要な理由を明確にすること(2025-05-14)
	// err = srdblib.Dbmap.SelectOne(&user1, "select * from user where irank = (select min(irank) from user where `rank` = 'B-5')")
	err = srdblib.Dbmap.SelectOne(&user1, "select * from user where irank = (select min(irank) from user where `rank` = 'B-5' and irank != 0) ")
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.SelectOne(): %w", err)
		log.Printf("HandlerShowRank(): %s\n", err.Error())
		return
	}

	showrank.Userlist, err = SelectShowRank(client, user1.Irank+1) //	SS-5〜A-1とB-5のトップ
	if err != nil {
		err = fmt.Errorf("SelectShowRank(): %w", err)
		log.Printf("HandlerShowRank(): %s\n", err.Error())
		return
	}
	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "showrank.gtpl", showrank); err != nil {
		err = fmt.Errorf("Handler(): %w", err)
		log.Printf("%s\n", err.Error())
	}

}
