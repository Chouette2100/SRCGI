// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	"fmt"
	// "html"
	"log"

	//	"math/rand"
	// "sort"
	"strconv"
	//	"strings"
	// "time"
	//	"os"

	// "runtime"

	// "encoding/json"

	"html/template"
	"net/http"

	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// 入力フォーム画面
func EventTopHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles(
	//		"templates/top.gtpl",
	//		"templates/top0.gtpl",
	//		"templates/bbs-2.gtpl",
	//		"templates/top1.gtpl",
	//		"templates/top2.gtpl",
	//	))

	eventid := r.FormValue("eventid")
	suserno := r.FormValue("userno")
	if suserno == "" {
		suserno = "0"
	}
	userno, _ := strconv.Atoi(suserno)
	log.Printf("      eventid=%s userno=%d\n", eventid, userno)

	tpl := template.Must(template.ParseFiles(
		"templates/eventtop.gtpl",
	))

	//	eventinf, _ := SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		return
	}
	// Event_inf = *eventinf

	if err := tpl.ExecuteTemplate(w, "eventtop.gtpl", eventinf); err != nil {
		log.Println(err)
	}

}
