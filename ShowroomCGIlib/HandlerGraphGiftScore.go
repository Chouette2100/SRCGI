//	Copyright © 2024 chouette.21.00@gmail.com
//	Released under the MIT license
//	https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bytes"
	"fmt"
	//	"html"
	"log"

	//	"math/rand"
	//	"sort"
	"strconv"
	"strings"
	//	"time"

	//	"bufio"
	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"

	//	"database/sql"

	//	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi"
	//	"github.com/Chouette2100/srapi"
	//	"github.com/Chouette2100/srdblib"
)


func HandlerGraphGs(w http.ResponseWriter, req *http.Request) {

	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	eventid := req.FormValue("eventid")
	//	maxpoint, _ := strconv.Atoi(req.FormValue("maxpoint"))
	smaxpoint := req.FormValue("maxpoint")
	maxpoint, _ := strconv.Atoi(smaxpoint)
	if maxpoint < 10000 {
		maxpoint = 0
		smaxpoint = "0"
	}
	sgscale := req.FormValue("gscale")
	if sgscale == "" || sgscale == "0" {
		sgscale = "100"
	}
	gscale, _ := strconv.Atoi(sgscale)
	/*
		gschk100 := ""
		gschk90 := ""
		gschk80 := ""
		gschk70 := ""
		switch sgscale {
		case "100":
			gschk100 = "checked"
		case "90":
			gschk90 = "checked"
		case "80":
			gschk80 = "checked"
		case "70":
			gschk70 = "checked"
		default:
			gschk100 = "checked"
		}
	*/
	resetcolor := req.FormValue("resetcolor")

	//	log.Printf("      eventid=%s maxpoint=%d(%s) resetcolor=[%s]\n", eventid, maxpoint, smaxpoint, resetcolor)

	if resetcolor == "on" {
		Resetcolor(eventid)
	}

	filename, _ := GraphTotalPoints(eventid, maxpoint, gscale)
	if Serverconfig.WebServer == "nginxSakura" {
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	} else if Serverconfig.WebServer == "Apache2Ubuntu" {
		filename = "/public/" + filename
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-total.gtpl"))

	// テンプレートに出力する値をマップにセット
	/*
		values := map[string]string{
			"filename": req.FormValue("FileName"),
		}
	*/
	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
		"maxpoint": smaxpoint,
		"gscale":   sgscale,
	}

	// マップを展開してテンプレートを出力する
	if err := tpl.ExecuteTemplate(w, "graph-total.gtpl", values); err != nil {
		log.Println(err)
	}
}
//	var Nfseq int

func GraphGiftScore(eventid string, maxpoint int, gscale int) (filename string, status int) {

	status = 0

	Event_inf.Event_ID = eventid

	IDlist, sts := SelectEventInfAndRoomList()

	if sts != 0 {
		log.Printf("status of SelectEventInfAndRoomList() =%d\n", sts)
		status = sts
		return
	}

	eventname, period, _ := SelectEventNoAndName(eventid)

	Event_inf.Maxpoint = maxpoint
	Event_inf.Gscale = gscale
	UpdateEventInf(&Event_inf)

	if Serverconfig.WebServer == "None" {
		filename = fmt.Sprintf("%03d.svg", Nfseq)
		Nfseq = (Nfseq + 1) % 1000
	} else {
		filename = fmt.Sprintf("%03d.svg", os.Getpid()%1000)
		//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
		//	filename = fmt.Sprintf("%0d.svg", r.Intn(100))

	}

	GraphScore01(filename, IDlist, eventname, period, maxpoint)

	/*
		fmt.Printf("Content-type:text/html\n\n")
		fmt.Printf("<!DOCTYPE html>\n")
		fmt.Printf("<html lang=\"ja\">\n")
		fmt.Printf("<head>\n")
		fmt.Printf("  <meta charset=\"UTF-8\">\n")
		fmt.Printf("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
		//	fmt.Printf("  <meta http-equiv=\"refresh\" content=\"30; URL=\">\n")
		fmt.Printf("  <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n")
		fmt.Printf("  <title></title>\n")
		fmt.Printf("</head>\n")
		fmt.Printf("<body>\n")
		fmt.Printf("<img src=\"test.svg\" alt=\"\" width=\"100%%\">")
		fmt.Printf("</body>\n")
		fmt.Printf("</html>\n")
	*/

	return
}
