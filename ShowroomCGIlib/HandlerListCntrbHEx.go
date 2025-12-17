package ShowroomCGIlib

import (
	"fmt"
	"log"

	//	"math"
	// "sort"
	"strconv"
	"strings"
	"time"

	//	"bufio"
	//	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"

	// "database/sql"

	_ "github.com/go-sql-driver/mysql"
	//	"github.com/PuerkitoBio/goquery"
	//	svg "github.com/ajstarks/svgo/float"
	// "github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type CntrbHEx_Header struct {
	Eventid    string
	Eventname  string
	Period     string
	Maxpoint   int
	Gscale     int
	Userno     int
	Username   string
	ShortURL   string
	Tlsnid     int
	Listener   string
	Ie         int
	Tlsnid_b   int
	Listener_b string
	Tlsnid_f   int
	Listener_f string
}

type CntrbHistoryExInf struct {
	Point     int
	Roomno    int
	Longname  string
	Eventid   string
	Eventname string
	Starttime time.Time
	Endtime   time.Time
	Stnow     time.Time
}

type CntrbHistoryEx []CntrbHistoryExInf

/*
        HandlerListCntrbH()
		指定されたイベント、配信者、リスナーの貢献ポイントの履歴を表示する。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。

*/

func ListCntrbHExHandler(w http.ResponseWriter, req *http.Request) {

	var err error

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrbH-h.gtpl", "templates/list-cntrbH.gtpl"))
	// funcMap := template.FuncMap{
	// 	"sub":        func(i, j int) int { return i - j },
	// 	"Comma":      func(i int) string { return humanize.Comma(int64(i)) },
	// 	"FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
	// }
	// tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-cntrbHEx-h.gtpl", "templates/list-cntrbHEx.gtpl"))

	eventid := req.FormValue("eventid")
	userno, _ := strconv.Atoi(req.FormValue("userno"))
	tlsnid, _ := strconv.Atoi(req.FormValue("tlsnid"))
	name := req.FormValue("name")
	ie, _ := strconv.Atoi(req.FormValue("ie"))
	log.Printf(" eventid=%s, userno=%d, tlsnid=%d\n", eventid, userno, tlsnid)

	var acqtimelist []time.Time
	if name == "" {
		acqtimelist, _ = SelectAcqTimeList(eventid, userno)
		if len(acqtimelist) == 0 {
			fmt.Fprintf(w, "HandlerListCntrbH() No AcqTimeList\n")
			fmt.Fprintf(w, "Check eventid and userno\n")
			log.Printf("No AcqTimeList\n")
			return
		}
	}

	/*
		type Tlsnidinf struct {
		Norder	int
		Tlsnid	int
		Listener	string
		}
		func SelectTlsnidList(eventid string, userno int, tlsnid int, smplt time.Time) ( tlsnidinflist [3]Tlsnidinf, status int) {
	*/

	var eventinf exsrapi.Event_Inf
	GetEventInf(eventid, &eventinf)

	var cntrbh_header CntrbH_Header

	cntrbh_header.Eventid = eventid
	cntrbh_header.Eventname = eventinf.Event_name
	cntrbh_header.Period = eventinf.Period
	cntrbh_header.Ie = ie

	cntrbh_header.Maxpoint = eventinf.Maxpoint

	cntrbh_header.Userno = userno
	cntrbh_header.Name = name

	// Turnstile導入用(2) ------------------------
	// Turnstile検証（セッション管理込み）
	// Turnstile検証要求後の状態を管理する
	lastrequestid := ""
	requestid := req.FormValue("requestid")
	if requestid != "" {
		// 最初のパス、検証のための場合と、すでにクッキーを持っている場合と両方ある
		lastrequestid = requestid
	}
	cntrbh_header.RequestID = req.Context().Value("requestid").(string)
	// ------

	result, tsErr := CheckTurnstileWithSession(w, req, &cntrbh_header)
	if result != TurnstileOK {
		// チャレンジページまたはエラーページが表示済みなので終了
		if tsErr != nil {
			log.Printf("Turnstile check error: %v\n", tsErr)
		}
		return
	}

	log.Printf(" cntrbh_header.RequestID = %s, lastrequestid = %s\n", cntrbh_header.RequestID, lastrequestid)
	if lastrequestid == "" {
		result, err := srdblib.Dbmap.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", cntrbh_header.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
	} else {
		result, err := srdblib.Dbmap.Exec(
			"UPDATE accesslog SET turnstilestatus= 0 WHERE requestid = ?", cntrbh_header.RequestID)
		log.Printf("  Update accesslog turnstilestatus=0 result=%+v, err=%+v\n", result, err)
		result, err = srdblib.Dbmap.Exec(
			"DELETE FROM accesslog WHERE requestid = ?", lastrequestid)
		log.Printf("  delete from accesslog where lastrequestid = %s result=%+v, err=%+v\n",
			lastrequestid, result, err)
	}
	// -------------------------------------------

	// _, _, _, _, _, _, _, _, roomname, roomurlkey, _, _ := GetRoomInfoByAPI(fmt.Sprintf("%d", userno))
	pu := &srdblib.User{}
	var intf interface{}
	intf, err = srdblib.Dbmap.Get(pu, userno)
	if err != nil {
		err = fmt.Errorf("GetUserInf() error: %w", err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	pu = (intf).(*srdblib.User)

	cntrbh_header.Username = pu.Longname
	cntrbh_header.ShortURL = pu.Userid

	if name == "" {
		tlsnidinflist, _ := SelectTlsnidList(eventid, userno, tlsnid, acqtimelist[len(acqtimelist)-1])

		cntrbh_header.Tlsnid = tlsnid
		cntrbh_header.Listener = tlsnidinflist[1].Listener

		cntrbh_header.Tlsnid_b = tlsnidinflist[0].Tlsnid
		cntrbh_header.Listener_b = tlsnidinflist[0].Listener
		cntrbh_header.Tlsnid_f = tlsnidinflist[2].Tlsnid
		cntrbh_header.Listener_f = tlsnidinflist[2].Listener
	} else {
		cntrbh_header.Tlsnid = tlsnid
		cntrbh_header.Listener = name
		cntrbh_header.Tlsnid_b = -1
		cntrbh_header.Tlsnid_f = -1
	}

	// if err := tpl.ExecuteTemplate(w, "list-cntrbHEx-h.gtpl", cntrbh_header); err != nil {
	// 	log.Println(err)
	// }
	var cntrbhistoryEx *CntrbHistoryEx
	cntrbhistoryEx, err = SelectCntrbHistoryEx(eventid, userno, tlsnid)
	if err != nil {
		err = fmt.Errorf("SelectCntrbHistoryEx() error: %w", err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	lpoint := 0
	lroomno := 0
	lid := ""
	for i := len(*cntrbhistoryEx) - 1; i >= 0; i-- {
		v := (*cntrbhistoryEx)[i]
		vid := strings.Split(v.Eventid, "?")
		if v.Point == lpoint && lroomno == v.Roomno && lid == vid[0] {
			// if i == len(*cntrbhistoryEx)-2 {
			// 	*cntrbhistoryEx = (*cntrbhistoryEx)[:i]
			// } else {
			// *cntrbhistoryEx = append((*cntrbhistoryEx)[:i+1], (*cntrbhistoryEx)[i+2:]...)
			*cntrbhistoryEx = append((*cntrbhistoryEx)[:i], (*cntrbhistoryEx)[i+1:]...)
			// }
		} else {
			lpoint = v.Point
			lroomno = v.Roomno
			lid = vid[0]
		}
	}

	cntrbh_header.CntrbhistoryEx = cntrbhistoryEx

	tpl := template.Must(template.New("").Funcs(*ListCntrbHfuncMap).ParseFiles("templates/list-cntrbHEx.gtpl"))

	// if err = tpl.ExecuteTemplate(w, "list-cntrbHEx.gtpl", cntrbhistoryEx); err != nil {
	if err := tpl.ExecuteTemplate(w, "list-cntrbHEx.gtpl", cntrbh_header); err != nil {
		log.Println(err)
	}

}

/*
type Tlsnidinf struct {
	Norder   int
	Tlsnid   int
	Listener string
}

func SelectTlsnidList(eventid string, userno int, tlsnid int, smplt time.Time) (tlsnidinflist [3]Tlsnidinf, status int) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	tlsnidinflist[0].Tlsnid = -1
	tlsnidinflist[2].Tlsnid = -1
	tlsnidinflist[0].Listener = "n/a"
	tlsnidinflist[2].Listener = "n/a"

	//	指定された時刻の貢献ポイントランキングを取得する。
	sql := "select norder, t_lsnid, listner from eventrank "
	sql += " where eventid = ? and userid =? and ts = ? order by norder"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sql)

	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbNow() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(eventid, userno, smplt)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbNow() (6) err=%s\n", srdblib.Dberr.Error())
		status = -6
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		//	Err = rows.Scan(&cntrbinf.Ranking, &cntrbinf.Tlsnid, &cntrbinf.Point, &cntrbinf.Incremental[loc], &cntrbinf.ListenerName, &cntrbinf.LastName)
		if found {
			srdblib.Dberr = rows.Scan(&tlsnidinflist[2].Norder, &tlsnidinflist[2].Tlsnid, &tlsnidinflist[2].Listener)
			break
		} else {
			srdblib.Dberr = rows.Scan(&tlsnidinflist[1].Norder, &tlsnidinflist[1].Tlsnid, &tlsnidinflist[1].Listener)
		}
		if srdblib.Dberr != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", srdblib.Dberr.Error())
			status = -7
			return
		}
		if tlsnidinflist[1].Tlsnid == tlsnid {
			found = true
		} else {
			tlsnidinflist[0].Norder = tlsnidinflist[1].Norder
			tlsnidinflist[0].Tlsnid = tlsnidinflist[1].Tlsnid
			tlsnidinflist[0].Listener = tlsnidinflist[1].Listener
		}
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("GetCurrentScore() (8) err=%s\n", srdblib.Dberr.Error())
		status = -8
		return
	}
	return
}
*/

// 指定したリスナーの貢献ポイントの履歴を取得する。
func SelectCntrbHistoryEx(
	eventid string,
	userno int,
	tlsnid int,
) (
	cntrbhistoryEx *CntrbHistoryEx,
	err error,
) {

	sqlst := " SELECT point, roomno, longname, eventid, eventname, starttime, endtime, NOW() as stnow "
	sqlst += " FROM ( "
	sqlst += "   SELECT er.point, er.userid as roomno, u.longname, er.eventid, e.event_name as eventname, e.starttime, e.endtime, "
	sqlst += "     ROW_NUMBER() "
	sqlst += "       OVER (PARTITION BY er.eventid, er.userid ORDER BY er.ts DESC) as rn " // eventidとuseridの組み合わせごとに、tsの降順で順位を付ける
	sqlst += "   FROM eventrank er "
	sqlst += "   JOIN event e on e.eventid = er.eventid "
	sqlst += "   JOIN user u on u.userno = er.userid "
	sqlst += "   WHERE lsnid = ? " // まずlsnidで絞り込む
	sqlst += " ) AS ranked "
	sqlst += " WHERE rn = 1 and point > 0 "    // ループの最新行（tsが最大）を選択する、pointが0または-1は無効なデータである
	sqlst += " ORDER BY endtime desc, roomno " // 最後にイベント終了日で、終了日が同じ時はユーザー番号順にソートする

	cntrbhistoryEx = &CntrbHistoryEx{}
	_, err = srdblib.Dbmap.Select(cntrbhistoryEx, sqlst, tlsnid)
	if err != nil {
		err = fmt.Errorf("SelectCntrbHistoryEx() error: %w", err)
		log.Printf("%s\n", err.Error())
		return nil, err
	}

	return
}

/*
type P_c struct {
	Pnt  int
	Cnt  int
	Wcnt int
	Cmp  int
}
type Pclist []P_c

// sort.Sort()のための関数三つ
func (p Pclist) Len() int {
	return len(p)
}
func (p Pclist) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Pclist) Choose(from, to int) (s Pclist) {
	s = p[from:to]
	return
}
func (p Pclist) Less(i, j int) bool {
	return p[i].Wcnt > p[j].Wcnt
}

func InsertTargetIntoTimtable(eventid string, userno int, ts time.Time, nfr int) (target int, status int) {
	status = 0

	var stmt_tg *sql.Stmt
	var rows *sql.Rows

	sql_tg := "select increment,count(*) from eventrank "
	sql_tg += " where eventid = ? and userid = ? and ts = ? "
	sql_tg += " group by increment order by count(*) desc;"
	stmt_tg, srdblib.Dberr = srdblib.Db.Prepare(sql_tg)
	if srdblib.Dberr != nil {
		log.Printf("InsertIntoTimtableTarget() (1) err=%s\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt_tg.Close()

	rows, srdblib.Dberr = stmt_tg.Query(eventid, userno, ts)
	if srdblib.Dberr != nil {
		log.Printf("InsertIntoTimtableTarget() (2) err=%s\n", srdblib.Dberr.Error())
		status = -2
		return
	}
	defer rows.Close()

	pcl := make(Pclist, 0)
	var pc P_c

	for rows.Next() {
		srdblib.Dberr = rows.Scan(&pc.Pnt, &pc.Cnt)
		if srdblib.Dberr != nil {
			log.Printf("InsertIntoTimtableTarget() (3) err=%s\n", srdblib.Dberr.Error())
			status = -3
			return
		}
		pcl = append(pcl, pc)
	}

	GetWeightedCnt(pcl, nfr)

	sort.Sort(pcl)

	log.Printf("%s %d %s\n", eventid, userno, ts.Format("2006/01/02 15:04:05"))
	for i := 0; i < len(pcl); i++ {
		log.Printf("%6d %3d %5d\n", pcl[i].Pnt, pcl[i].Cnt, pcl[i].Wcnt)
	}

	if len(pcl) == 0 {
		target = -1
		return
	}
	target = pcl[0].Pnt

	if target != -1 {
		sql_ud := "update timetable set target = ? where eventid = ? and userid = ? and sampletm2 = ?"
		stmt, err := srdblib.Db.Prepare(sql_ud)
		if err != nil {
			log.Printf("InsertIntoTimtableTarget() Update/Prepare err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(target, eventid, userno, ts)

		if err != nil {
			log.Printf("error(UpdatePointsSetQstatus() Update/Exec) err=%s\n", err.Error())
			status = -2
		}

	}

	return
}

func GetWeightedCnt(pcl Pclist, nfr int) {

	maxpnt := 1842 + (nfr-1)*600

	for i := 0; i < len(pcl); i++ {
		cmp := 0
		pcl[i].Wcnt = pcl[i].Cnt
		if pcl[i].Pnt > 599 && pcl[i].Pnt < 1263 {
			switch pcl[i].Pnt % 600 {
			case 62:
				cmp = 8
			case 50:
				cmp = 6
			case 52:
				cmp = -6
			case 40:
				cmp = -8
			default:
				continue
			}
			pcl[i].Cmp = (pcl[i].Pnt/600)*10 + cmp
			pcl[i].Wcnt = pcl[i].Cnt * pcl[i].Cmp
		} else if pcl[i].Pnt >= 1263 && pcl[i].Pnt < maxpnt {
			switch pcl[i].Pnt % 600 {
			case 42:
				cmp = 8
			case 30:
				cmp = 6
			case 592:
				cmp = -6
			case 580:
				cmp = -8
			default:
				continue
			}
			//	pcl[i].Cmp =  ((pcl[i].Pnt + 20) / 600) * 10 + cmp
			pcl[i].Cmp = ((pcl[i].Pnt+20)/600)*10 + cmp
			//	pcl[i].Wcnt = pcl[i].Cnt * pcl[i].Cmp
			pcl[i].Wcnt = pcl[i].Cnt * 3 * pcl[i].Cmp
		} else if pcl[i].Pnt == 0 {
			pcl[i].Wcnt = 2
		} else if pcl[i].Pnt == -1 {
			pcl[i].Wcnt = 0
		}
	}
}
*/
