package ShowroomCGIlib

import (
	"fmt"
	"log"

	//	"math"
	//	"sort"
	"strconv"
	//	"strings"
	"time"

	//	"bufio"
	//	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	//	"github.com/PuerkitoBio/goquery"
	//	svg "github.com/ajstarks/svgo/float"
	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

type CntrbS_Header struct {
	Eventid   string
	Eventname string
	Period    string
	Target    int
	Maxpoint  int
	Gscale    int
	Userno    int
	Username  string
	ShortURL  string
	S_stime   string
	S_etime   string
	Srt       int
	Ie        int
	Ifrm      int
	Ifrm1     int
	Ifrm_b    int
	Ifrm_f    int
}

type CntrbInfS struct {
	Ranking      int
	Point        int
	Incremental  int
	ListenerName string
	LastName     string
	Tlsnid       int
	Eventid      string
	Userno       int
}

func SelectTargetfromTimetable(
	eventid string,
	userno int,
	ts time.Time,
) (
	target int,
	err error,
) {
	sqls := "SELECT target FROM timetable WHERE eventid = ? AND userid = ? AND sampletm2 = ?"
	err = srdblib.Db.QueryRow(sqls, eventid, userno, ts).Scan(&target)
	if err != nil {
		err = fmt.Errorf("QueryRow().Scan()  error: %v", err)
	}
	return
}

func UpdateTimetableSetTarget(
	eventid string,
	userno int,
	ts time.Time,
	target int,
) (
	err error,
) {
	sqlu := "UPDATE timetable SET target = ? WHERE eventid = ? AND userid = ? AND sampletm2 = ?"
	_, err = srdblib.Db.Exec(sqlu, target, eventid, userno, ts)
	if err != nil {
		err = fmt.Errorf("Exec()  error: %v", err)
	}
	return
}

/*
        HandlerListCntrbS()
			一枠分の貢献ポイントランキングを表示する。
			表示内容にはリスナーにニックネームの変更が含まれる。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。

*/

func HandlerListCntrbS(w http.ResponseWriter, req *http.Request) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	GetUserInf(req)

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrb-h.gtpl", "templates/list-cntrb.gtpl"))
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrbS-h.gtpl", "templates/list-cntrbS.gtpl"))
	funcMap := template.FuncMap{
		"sub":   func(i, j int) int { return i - j },
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
	}
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-cntrbS-h.gtpl", "templates/list-cntrbS.gtpl"))

	eventid := req.FormValue("eventid")
	userno, _ := strconv.Atoi(req.FormValue("userno"))
	ifrm, _ := strconv.Atoi(req.FormValue("ifrm"))
	sort := req.FormValue("sort")
	target, _ := strconv.Atoi(req.FormValue("target"))
	ie, _ := strconv.Atoi(req.FormValue("ie"))
	log.Printf(" eventid=%s, userno=%d, ifrm=%d\n", eventid, userno, ifrm)

	acqtimelist, _ := SelectAcqTimeList(eventid, userno)
	ts := acqtimelist[ifrm]

	cntrbinflists, status := SelectCntrbSingle(eventid, userno, ts, sort)
	if status != 0 {
		log.Printf(" SelectCntrbSingle() returned %d in HandlerListCntrbS()\n", status)
		return
	}
	//	log.Printf(" len(cntrbinflists)=%d\n", len(cntrbinflists))

	var eventinf exsrapi.Event_Inf
	GetEventInf(eventid, &eventinf)

	var cntrbs_header CntrbS_Header

	cntrbs_header.Eventid = eventid
	cntrbs_header.Eventname = eventinf.Event_name
	cntrbs_header.Maxpoint = eventinf.Maxpoint
	cntrbs_header.Gscale = eventinf.Gscale
	cntrbs_header.Ie = ie

	_, _, _, _, _, _, _, _, roomname, roomurlkey, _, _ := GetRoomInfoByAPI(fmt.Sprintf("%d", userno))
	cntrbs_header.Userno = userno
	cntrbs_header.Username = roomname
	cntrbs_header.ShortURL = roomurlkey

	var err error
	if target == 0 {
		cntrbs_header.Target, err = SelectTargetfromTimetable(eventid, userno, ts)
		if err != nil {
			log.Printf("SelectTargetfromTimetable() returned %s\n", err)
			return
		}
	} else {
		cntrbs_header.Target = target
		err = UpdateTimetableSetTarget(eventid, userno, ts, target)
		if err != nil {
			log.Printf("SelectTargetfromTimetable() returned %s\n", err)
			return
		}
	}

	var stime, etime time.Time
	sql := "select stime, etime from timetable where eventid = ? and userid = ? and sampletm2 = ? "
	srdblib.Dberr = srdblib.Db.QueryRow(sql, eventid, userno, ts).Scan(&stime, &etime)
	if srdblib.Dberr != nil {
		log.Printf("select stime, etime from timetable where eventid = %s and userid = %d and sampletm2 = %+v\n", eventid, userno, ts)
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -11
		return
	}
	log.Printf(" s=%v e=%v\n", stime, etime)
	cntrbs_header.S_stime = stime.Format("2006-01-02 15:04")
	cntrbs_header.S_etime = etime.Format("15:04")
	cntrbs_header.Ifrm = ifrm
	cntrbs_header.Ifrm1 = ifrm + 1
	cntrbs_header.Ifrm_b = ifrm - 1 //	範囲外になると -1 になる
	cntrbs_header.Ifrm_f = ifrm + 1 //	範囲外になったら -1 にする
	if sort == "D" {
		//	増分順にソートする。
		cntrbs_header.Srt = 1
	} else {
		//	累計順にソートする。
		cntrbs_header.Srt = 0
	}
	if cntrbs_header.Ifrm_f >= len(acqtimelist) {
		cntrbs_header.Ifrm_f = -1
	}

	if err := tpl.ExecuteTemplate(w, "list-cntrbS-h.gtpl", cntrbs_header); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "list-cntrbS.gtpl", cntrbinflists); err != nil {
		log.Println(err)
	}
}

/*
	        SelectCntrbSingle()

			指定したイベント、ユーザー、時刻の貢献ポイントランキングを一枠分だけ取得する。
			リスナー名の突き合わせのチェックを目的とするページを作るために使用する。

	        引数
			eventid			string			イベントID
			userno			int				ユーザーID
			ts				int				ユーザーID
			loc				int				データの格納場所（ 0 だったら先頭）

	        戻り値
	        cntrbinflists	[] CntrbInf		貢献ポイントランキング（最終貢献ポイント順）
			stats			int				== 0 正常終了	!= 0 データベースアクセス時のエラー
*/
func SelectCntrbSingle(
	eventid string,
	userno int,
	ts time.Time,
	sort string,
) (
	cntrbinflists []CntrbInfS,
	status int,
) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	cntrbinflists = make([]CntrbInfS, 0)

	//	貢献ポイントランキングを取得する。
	sql := "select norder, t_lsnid, point, increment, listner, lastname from eventrank "
	if sort == "D" {
		sql += " where eventid = ? and userid =? and ts = ? order by increment desc"
	} else {
		sql += " where eventid = ? and userid =? and ts = ? order by norder"
	}
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sql)

	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbSingle() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(eventid, userno, ts)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbSingle() (6) err=%s\n", srdblib.Dberr.Error())
		status = -6
		return
	}
	defer rows.Close()

	var cntrbinf CntrbInfS
	tlsnid := 0

	for rows.Next() {
		//	Err = rows.Scan(&cntrbinf.Ranking, &tlsnid, &cntrbinf.Point, &cntrbinf.Incremental[0], &cntrbinf.ListenerName, &cntrbinf.LastName)
		srdblib.Dberr = rows.Scan(&cntrbinf.Ranking, &tlsnid, &cntrbinf.Point, &cntrbinf.Incremental, &cntrbinf.ListenerName, &cntrbinf.LastName)
		if srdblib.Dberr != nil {
			log.Printf("SelectCntrbSingle() (7) err=%s\n", srdblib.Dberr.Error())
			status = -7
			return
		}
		//	log.Printf("%+v\n", cntrbinf)
		cntrbinflists = append(cntrbinflists, cntrbinf)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("SelectCntrbSingle() (8) err=%s\n", srdblib.Dberr.Error())
		status = -8
		return
	}

	status = 0

	return

}
