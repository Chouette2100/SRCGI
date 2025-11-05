// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// "strings"
	//	"io" //　ログ出力設定用。必要に応じて。
	//	"sort" //	ソート用。必要に応じて。

	// "html/template"
	"html/template"
	"net/http"

	// "github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

/*
終了イベント一覧を作るためのハンドラー

Ver. 0.1.0
*/

// イベントにおける貢献ポイントデータ
type Contribution struct {
	Ieventid int
	Roomid   int
	Viewerid int
	Irank    int
	Point    int
}

type Result struct {
	Irank    int
	Viewerid int
	Name     string
	Point    int
}

type HCntrbInf struct {
	Ieventid   int
	Eventid    string
	Event_name string
	Period     string
	Roomid     int
	Result     []Result
	Filename   string
}

func ContributorsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var err error
	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	seventid := r.FormValue("ieventid")
	ieventid, _ := strconv.Atoi(seventid)

	var intf []interface{}
	intf, err = srdblib.Dbmap.Select(&srdblib.Wevent{}, "SELECT "+clmlist["wevent"]+" FROM wevent WHERE ieventid = ?", ieventid)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(): %s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	event := *(intf[0].(*srdblib.Wevent))

	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	client, cookiejar, err := exsrapi.CreateNewClient("")
	if err != nil {
		err = fmt.Errorf("exsrapi.CeateNewClient(): %s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))

		return //       エラーがあれば、ここで終了
	}
	defer cookiejar.Save()

	hcntrbinf := HCntrbInf{
		Ieventid:   ieventid,
		Eventid:    event.Eventid,
		Event_name: event.Event_name,
		Period:     event.Period,
		Roomid:     roomid,
	}

	if err = MakeFileOfContributors(client, &hcntrbinf); err != nil {
		err = fmt.Errorf("MakeFileOfContributors(): %w", err)
		log.Printf("err=%s\n", err.Error())
		// w.Write([]byte(err.Error()))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/contributors.gtpl",
	))

	if err := tpl.ExecuteTemplate(w, "contributors.gtpl", hcntrbinf); err != nil {
		err = fmt.Errorf("tpl.ExceuteTemplate(w,\"contributors.gtpl\", hcntrbinf) err=%s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))
	}

}

func GetAndSaveContributors(
	client *http.Client,
	ieventid int,
	roomid int,
) (
	result []Result,
	err error,
) {

	// すでにデータがあるかどうかを確認する
	var count int64
	sqlst := "SELECT COUNT(*) cnt FROM contribution WHERE ieventid = ? AND roomid = ? "
	count, err = srdblib.Dbmap.SelectInt(sqlst, ieventid, roomid)
	if err != nil {
		err = fmt.Errorf("GetAndSaveContributors(): %w", err)
		return
	}
	if count == 0 {

		// しばらくのあいだデータ未取得の場合はこの機能を停止する。
		err = fmt.Errorf("GetAndSaveContributors(): contribution data not found and data fetching is disabled temporarily")
		log.Printf("%s\n", err.Error())
		return

		/*
			var cr *srapi.Contribution_ranking
			cr, err = srapi.ApiEventContribution_ranking(client, ieventid, roomid)
			if err != nil {
				err = fmt.Errorf("GetContribution(): %w", err)
				return
			}
			// 取得した貢献ポイントデータを格納する
			tnow := time.Now().Truncate(time.Second)
			for _, c := range cr.Ranking {

				viewer := srdblib.Viewer{
					Viewerid: c.UserID,
					Name:     c.Name,
					Sname:    c.Name,
					Ts:       tnow,
				}
				if err = srdblib.UpinsViewerSetProperty(client, tnow, &viewer); err != nil {
					err = fmt.Errorf("UpinsViewerSetProperty(): %w", err)
					return
				}

				contribution := Contribution{
					Ieventid: ieventid,
					Roomid:   roomid,
					Viewerid: c.UserID,
					Irank:    c.Rank,
					Point:    c.Point,
				}
				// log.Printf("Insert(): %8d%8d%8d%4d%10d\n",
				// 	ieventid, roomid, contribution.Viewerid, contribution.Irank, contribution.Point)

				if err = srdblib.Dbmap.Insert(&contribution); err != nil {
					err = fmt.Errorf("Insert(): %w", err)
					// log.Printf("Insert(): %s\n", err.Error())
					return
				}

			}
		*/
	}

	sqlst = "SELECT c.irank, c.viewerid, v.name, c.point FROM contribution c "
	sqlst += " JOIN viewer v ON c.viewerid = v.viewerid "
	sqlst += " WHERE c.ieventid = ? and c.roomid = ? ORDER BY c.irank"

	var intf []interface{}
	intf, err = srdblib.Dbmap.Select(&Result{}, sqlst, ieventid, roomid)
	if err != nil {
		err = fmt.Errorf("GetAndSaveContributors(): %w", err)
		return
	}

	result = make([]Result, len(intf))

	for i, v := range intf {
		result[i] = *(v.(*Result))
	}

	return

}

// イベントの貢献ポイントランキングをファイルに書き出す
func MakeFileOfContributors(
	client *http.Client,
	hc *HCntrbInf,
) (
	err error,
) {
	// 貢献ポイントデータを取得する
	var result []Result
	result, err = GetAndSaveContributors(client, hc.Ieventid, hc.Roomid)
	if err != nil {
		err = fmt.Errorf("GetAndSaveContributors(): %w", err)
		return
	}

	hc.Result = result

	// 貢献ポイントデータをファイルに書き出す
	//	ファイル名の決定
	uqn := ""
	if Serverconfig.WebServer == "None" {
		uqn = time.Now().Format("05.000")
		uqn = strings.Replace(uqn, ".", "", -1)
	} else {
		uqn = fmt.Sprintf("%03d", os.Getpid()%1000)
	}
	hc.Filename = fmt.Sprintf("%s_%d_%d_%s", hc.Eventid, hc.Ieventid, hc.Roomid, uqn)

	var file *os.File
	for i := 0; i < 2; i++ {

		// ファイルを開く
		fn := "public/" + hc.Filename
		if i == 0 {
			fn = fn + "-1.csv"
		} else {
			fn = fn + "-2.csv"
		}
		file, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			err = fmt.Errorf("os.OpenFile(public/%s) err=%w", hc.Filename, err)
			return
		}
		defer file.Close()

		if i == 1 {
			file.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
		}

		// ファイルに書き出す
		fmt.Fprintf(file, "%d, %s, \"%s\"\n", hc.Ieventid, hc.Eventid, hc.Event_name)
		fmt.Fprintf(file, ",Period, \"%s\"\n", hc.Period)
		fmt.Fprintf(file, ",Roomid, %d\n", hc.Roomid)
		fmt.Fprintf(file, "Rank,Viewerid,Viewername,Point\n")
		for _, r := range hc.Result {
			fmt.Fprintf(file, "%d,%d,%d,\"%s\"\n", r.Irank, r.Viewerid, r.Point, r.Name)
		}
	}

	return
}
