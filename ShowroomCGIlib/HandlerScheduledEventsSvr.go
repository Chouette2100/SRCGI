/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

package ShowroomCGIlib

import (
	"fmt"
	"html/template"
	//	"io" //　ログ出力設定用。必要に応じて。
	"log"
	"net/http"
	"sort" //	ソート用。必要に応じて。
	"time"

	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

/*

	開催中のイベント一覧を作るためのハンドラー

	Ver. 0.1.0

*/

type T008top struct {
	TimeNow    int64
	Totalcount int
	ErrMsg     string
	Eventlist  []srapi.Event
}

//	"/t008top"に対するハンドラー
//	http://localhost:8080/t008top で呼び出される
func HandlerScheduledEventsSvr(
	w http.ResponseWriter,
	r *http.Request,
) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(r)
	if ! isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}




	client, cookiejar, err := exsrapi.CreateNewClient("")
	if err != nil {
		log.Printf("exsrapi.CeateNewClient(): %s", err.Error())
		return //	エラーがあれば、ここで終了
	}
	defer cookiejar.Save()

	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":         func(i int) string { return humanize.Comma(int64(i)) },                       //	3桁ごとに","を入れる関数。
		"UnixTimeToStr": func(i int64) string { return time.Unix(int64(i), 0).Format("01-02 15:04") }, //	UnixTimeを年月日時分に変換する関数。
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/scheduled-event-svr.gtpl"))

	// テンプレートに埋め込むデータ（ポイントやランク）を作成する
	top := new(T008top)
	top.TimeNow = time.Now().Unix()

	top.Eventlist, err = srapi.MakeEventListByApi(client, 3)
	if err != nil {
		err = fmt.Errorf("MakeListOfPoints(): %w", err)
		log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
		top.ErrMsg = err.Error()
	}
	top.Totalcount = len(top.Eventlist)

	//	ソートが必要ないときは次の行とimportの"sort"をコメントアウトする。
	//	無名関数のリターン値でソート条件を変更できます。
	sort.Slice(top.Eventlist, func(i, j int) bool { return top.Eventlist[i].Started_at < top.Eventlist[j].Started_at })

	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "scheduled-event-svr.gtpl", top); err != nil {
		log.Printf("tpl.ExecuteTemplate() returned error: %s\n", err.Error())
	}

}

