// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	"fmt"
	//	"html"
	"log"

	//	"math/rand"
	// "sort"
	"strconv"
	//	"strings"
	//	"time"
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

func ParamEventHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/param-event0.gtpl",
		"templates/param-event1.gtpl",
		"templates/param-event2.gtpl",
	))

	eventid := r.FormValue("eventid")

	log.Printf("      eventid=%s\n", eventid)

	//	eventinf, _ := SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, _ := srdblib.SelectFromEvent("event", eventid)
	// Event_inf = *eventinf

	userlist, _ := SelectEventuserList(eventid)
	for i := 0; i < len(userlist); i++ {
		if userlist[i].Userno == eventinf.Nobasis {
			userlist[i].Selected = "Selected"
		} else {
			userlist[i].Selected = ""
		}
	}

	if err := tpl.ExecuteTemplate(w, "param-event0.gtpl", eventinf); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "param-event1.gtpl", userlist); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "param-event2.gtpl", eventinf); err != nil {
		log.Println(err)
	}

}

func ParamEventCHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/param-eventc.gtpl"))
	eventid := r.FormValue("eventid")
	log.Printf("      eventid=%s\n", eventid)

	//	eventinf, _ := SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, _ := srdblib.SelectFromEvent("event", eventid)
	// Event_inf = *eventinf

	//	log.Println(eventinf)

	eventinf.Fromorder, _ = strconv.Atoi(r.FormValue("fromorder"))
	eventinf.Toorder, _ = strconv.Atoi(r.FormValue("toorder"))
	eventinf.Modmin, _ = strconv.Atoi(r.FormValue("modmin"))
	eventinf.Modsec, _ = strconv.Atoi(r.FormValue("modsec"))

	intervalmin, _ := strconv.Atoi(r.FormValue("intervalmin"))
	switch intervalmin {
	case 5, 6, 10, 15, 20, 30, 60:
		eventinf.Intervalmin = intervalmin
	default:
		eventinf.Intervalmin = 5
	}
	eventinf.Modmin = eventinf.Modmin % eventinf.Intervalmin //	不適切な入力に対する修正
	eventinf.Modsec = eventinf.Modsec % 60

	eventinf.Resethh, _ = strconv.Atoi(r.FormValue("resethh"))
	eventinf.Resetmm, _ = strconv.Atoi(r.FormValue("resetmm"))
	eventinf.Nobasis, _ = strconv.Atoi(r.FormValue("nobasis"))
	eventinf.Target, _ = strconv.Atoi(r.FormValue("target"))
	eventinf.Maxdsp, _ = strconv.Atoi(r.FormValue("maxdsp"))
	ncmap, _ := strconv.Atoi(r.FormValue("cmap"))
	if eventinf.Cmap != ncmap {
		Resetcolor(eventinf.Event_ID, ncmap)
		eventinf.Cmap = ncmap
	}

	//	UpdateEventInf(&eventinf)
	UpdateEventInf(eventinf)
	//	log.Println(eventinf)

	if err := tpl.ExecuteTemplate(w, "param-eventc.gtpl", eventinf); err != nil {
		log.Println(err)
	}

}
func SelectEventuserList(eventid string) (userlist []User, status int) {

	status = 0

	sql := "select e.userno,u.longname "
	sql += " from eventuser e join user u on e.userno=u.userno "
	sql += " where e.eventid = ? "
	sql += " order by e.userno"

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(eventid)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var user User
	i := 0

	user.Userno = 0
	user.Userlongname = "ポイント差は不要"
	userlist = append(userlist, user)
	i++

	for rows.Next() {
		err := rows.Scan(&user.Userno, &user.Userlongname)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		userlist = append(userlist, user)
		i++
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return

}
