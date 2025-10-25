// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"bytes"
	"fmt"
	//	"html"
	"log"
	//	"os"
	//	"sort"
	"strings"
	"time"

	//	"github.com/PuerkitoBio/goquery"

	//	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	// "database/sql"
	"github.com/Masterminds/sprig/v3"
	"github.com/dustin/go-humanize"
)

// ExperimentalHandler は枠別リスナー別貢献ポイントの取得対象ルームの編集を行います
func ExperimentalHandler(w http.ResponseWriter, r *http.Request) {

	type Experimentl struct {
		Title string
		Date  time.Time
	}
	experimental := Experimentl{}

	funcMap := sprig.FuncMap() // https://masterminds.github.io/sprig/
	funcMap["Comma"] = func(i int) string { return humanize.Comma(int64(i)) }
	funcMap["baseOfEventid"] = func(s string) string { ida := strings.Split(s, "?"); return ida[0] }
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/experimental.gtpl"))

	if err := tpl.ExecuteTemplate(w, "experimental.gtpl", experimental); err != nil {
		err = fmt.Errorf("tpl.ExecuteTemplate() error=%s", err.Error())
		w.Write([]byte(err.Error()))
		log.Println(err)
	}
}
