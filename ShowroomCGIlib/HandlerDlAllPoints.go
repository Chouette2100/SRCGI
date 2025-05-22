// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	//	"html"
	"log"
	"os"
	"strconv"
	"strings"

	//	"sort"
	"time"

	//	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	//	"database/sql"

	//	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type FapHeader struct {
	Eventinf *srdblib.Event
	Filename string
}

// イベントを獲得ポイントデータ取得の対象としてeventテーブルに登録する。
// イベントが開催中であれば指定した順位内のルームを取得対象として登録する。
// イベントが開催予定のものであればルームの登録は行わない。
// イベント開催中、開催予定にかかわらず、取得対象ルームの追加は srAddNewOnes で行われる。
func DlAllPointsHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/dl-all-points.gtpl",
	))

	eventid := r.FormValue("eventid")
	breg := r.FormValue("breg")
	ereg := r.FormValue("ereg")
	ibreg := 1
	if breg != "" {
		ibreg, _ = strconv.Atoi(breg)
	}
	iereg := 10
	if ereg != "" {
		iereg, _ = strconv.Atoi(ereg)
	}

	hd, err := MakeFileOfAllPoints(eventid, ibreg, iereg)
	if err != nil {
		log.Printf("MakeFileOfAllPoints(%s, %d, %d) err= %s\n", eventid, ibreg, iereg, err.Error())
		w.Write([]byte("MakeFileOfAllPoints(" + eventid + "," + breg + "," + ereg + ") err=" + err.Error()))
		return
	}

	if err := tpl.ExecuteTemplate(w, "dl-all-points.gtpl", hd); err != nil {
		log.Printf("tpl.ExceuteTemplate(w,\"dl-all-points.gtpl\", hd) err=%s\n", err.Error())
		w.Write([]byte("tpl.ExceuteTemplate(w,\"dl-all-points.gtpl\", hd)" + err.Error()))
	}

}

// 獲得ポイント全データをファイルに書き出す
// 間引きしたデータは復元する
func MakeFileOfAllPoints(
	eventid string,
	fromorder int,
	toorder int,
) (
	hd *FapHeader,
	err error,
) {

	hd = new(FapHeader)

	//	イベント情報を取得する
	var itfc interface{}
	itfc, err = srdblib.Dbmap.Get(srdblib.Event{}, eventid)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Get(Event{}, %s) err=%w", eventid, err)
		return
	}
	if itfc.(*srdblib.Event) == nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(Eventuser{}, %s) no data", eventid)
		return
	}

	hd.Eventinf = itfc.(*srdblib.Event)

	//	指定した条件を満たすルームの一覧を得る
	sqlstmt := "SELECT userno FROM eventuser WHERE eventid = ? ORDER BY point DESC LIMIT ? OFFSET ? "
	//	type MfapUserlist struct {
	//		Userno int
	//	}
	itfc, err = srdblib.Dbmap.Select(srdblib.Eventuser{}, sqlstmt, eventid, toorder-fromorder+1, fromorder-1)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(Eventuser{}, %s) err=%w", eventid, err)
		return
	}
	if itfc.([]interface{}) == nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(Eventuser{}, %s) no data", eventid)
		return
	}
	//	ul := itfc.([]srdblib.Eventuser)
	var Ul []int
	for _, v := range itfc.([]interface{}) {
		Ul = append(Ul, v.(*srdblib.Eventuser).Userno)
	}

	//	サンプルタイムのリストを作成する
	sqlstmt = "SELECT DISTINCT ts FROM points WHERE eventid = :Eventid AND user_id IN ( :Users ) ORDER BY ts "
	//	sqlstmt = "SELECT DISTINCT ts FROM points WHERE eventid = 'mattari_fireworks201' AND user_id IN ( 429729,431217,417115 ) ORDER BY ts "
	type Timelist struct {
		Ts time.Time
	}
	//	Eventid := hd.Eventinf.Eventid
	//	argmap := map[string]interface{}{"Users": []int{429729,431217,417115}, "Eventid": "mattari_fireworks201"}
	//	log.Printf("argmap = %+v\n", argmap)
	//	itfc, err = srdblib.Dbmap.Select(timelist{}, sqlstmt,
	//	map[string]interface{}{"Users": []int{429729,431217,417115}, "Eventid": "mattari_fireworks201"})
	itfctl, err := srdblib.Dbmap.Select(Timelist{}, sqlstmt, map[string]interface{}{"Users": Ul, "Eventid": hd.Eventinf.Eventid})
	//	itfctl, err := srdblib.Dbmap.Select(Timelist{}, sqlstmt)
	if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(timelist{}, %s) err=%w", eventid, err)
		return nil, err
	}
	// tl := itfctl.([]timelist)
	tl := make([]time.Time, len(itfctl))
	tlmap := map[time.Time]int{}
	for i, t := range itfctl {
		tl[i] = t.(*Timelist).Ts
		tlmap[tl[i]] = i
	}

	//	獲得ポイントデータを取得する
	type Pointdata struct {
		Ts   time.Time
		Rank int
		//	level int
		Point int
	}

	type userdata struct {
		//	username string
		Pd []Pointdata
	}

	udl := make([]userdata, len(Ul))
	for i, u := range Ul {
		udl[i].Pd = make([]Pointdata, len(tl))
		for j := 0; j < len(tl); j++ {
			udl[i].Pd[j].Point = -1
		}

		sqlstmt = "SELECT ts, `rank`, point FROM points WHERE eventid = ? AND user_id = ? ORDER BY ts "
		itfcpd, err := srdblib.Dbmap.Select(Pointdata{}, sqlstmt, eventid, u)
		if err != nil {
			err = fmt.Errorf("srdblib.Dbmap.Select(pointdata{}, %s, %d) err=%w", eventid, u, err)
			return nil, err
		}
		for _, tpd := range itfcpd {
			udl[i].Pd[tlmap[tpd.(*Pointdata).Ts]] = *tpd.(*Pointdata)
		}
		/*
			for j := 2; j < len(tl); j++ {
				if udl[i].Pd[j].Point == -1 {
					udl[i].Pd[j].Point = udl[i].Pd[j-1].Point
					udl[i].Pd[j].Rank = udl[i].Pd[j-1].Rank
				}
			}
		*/
	}

	//	取得した結果をファイルに書き出す

	//	ファイル名の決定
	uqn := ""
	// if Serverconfig.WebServer == "None" {
	uqn = time.Now().Format("20060102-150405.00")
	uqn = strings.Replace(uqn, ".", "", -1)
	// } else {
	// 	uqn = fmt.Sprintf("%03d", os.Getpid()%1000)
	// }
	hd.Filename = fmt.Sprintf("%d_%s", hd.Eventinf.Ieventid, uqn)

	for i := 0; i < 2; i++ {
		fn := "public/" + hd.Filename + "-1.csv"
		if i == 1 {
			fn = "public/" + hd.Filename + "-2.csv"
		}
		var file *os.File
		file, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			err = fmt.Errorf("os.OpenFile(%s) err=%w", fn, err)
			return
		}
		//	ファイルを閉じる
		defer file.Close()

		if i == 1 {
			file.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
		}

		ebuf := ""
		ubufno := ""
		ubufnm := ""
		for j := 0; j < len(Ul); j++ {
			ebuf += ","

			itfc, err = srdblib.Dbmap.Get(srdblib.User{}, Ul[j])

			if err != nil {
				err = fmt.Errorf("srdblib.Dbmap.Get(User{}, %d) err=%w", Ul[j], err)
				return
			}
			ubufnm += ",\"" + itfc.(*srdblib.User).Longname + "\""
			ubufno += fmt.Sprintf(",%d", Ul[j])
		}

		fmt.Fprintf(file, "\"%s\"%s\n", hd.Eventinf.Event_name, ebuf)
		fmt.Fprintf(file, "\"%s\"%s\n", hd.Eventinf.Eventid, ebuf)
		fmt.Fprintf(file, "%d%s\n", hd.Eventinf.Ieventid, ebuf)
		fmt.Fprintf(file, "\"%s\"%s\n", hd.Eventinf.Period, ebuf)
		fmt.Fprintf(file, "\n")
		fmt.Fprintf(file, "%s\n", ubufno)
		fmt.Fprintf(file, "%s\n", ubufnm)

		buf := tl[0].Format("2006/01/02 15:04:05")
		fmt.Fprintf(file, "%s\n", buf+ebuf)
		for i := 1; i < len(tl); i++ {
			buf := tl[i].Format("2006/01/02 15:04:05")
			for j := 0; j < len((Ul)); j++ {
				buf += ","
				if udl[j].Pd[i].Point != -1 {
					buf += fmt.Sprintf("%d", udl[j].Pd[i].Point)
				}
			}
			fmt.Fprintf(file, "%s\n", buf)
		}
		file.Close()
	}

	return
}
