package ShowroomCGIlib

import (
	"fmt"
	"log"

	//	"math"
	"sort"
	"strconv"
	"strings"
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

type CntrbH_Header struct {
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

type CntrbHistoryInf struct {
	S_stime     string
	S_etime     string
	Target      int
	Point       int
	Incremental int
	Listener    string
	Lastname    string
}

type CntrbHistory []CntrbHistoryInf

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

func HandlerListCntrbH(w http.ResponseWriter, req *http.Request) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	GetUserInf(req)

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrbH-h.gtpl", "templates/list-cntrbH.gtpl"))
	funcMap := template.FuncMap{
		"sub":   func(i, j int) int { return i - j },
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
	}
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-cntrbH-h.gtpl", "templates/list-cntrbH.gtpl"))

	eventid := req.FormValue("eventid")
	userno, _ := strconv.Atoi(req.FormValue("userno"))
	tlsnid, _ := strconv.Atoi(req.FormValue("tlsnid"))
	ie, _ := strconv.Atoi(req.FormValue("ie"))
	log.Printf(" eventid=%s, userno=%d, tlsnid=%d\n", eventid, userno, tlsnid)

	acqtimelist, _ := SelectAcqTimeList(eventid, userno)

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

	_, _, _, _, _, _, _, _, roomname, roomurlkey, _, _ := GetRoomInfoByAPI(fmt.Sprintf("%d", userno))
	cntrbh_header.Username = roomname
	cntrbh_header.ShortURL = roomurlkey

	tlsnidinflist, _ := SelectTlsnidList(eventid, userno, tlsnid, acqtimelist[len(acqtimelist)-1])

	cntrbh_header.Tlsnid = tlsnid
	cntrbh_header.Listener = tlsnidinflist[1].Listener

	cntrbh_header.Tlsnid_b = tlsnidinflist[0].Tlsnid
	cntrbh_header.Listener_b = tlsnidinflist[0].Listener
	cntrbh_header.Tlsnid_f = tlsnidinflist[2].Tlsnid
	cntrbh_header.Listener_f = tlsnidinflist[2].Listener

	if err := tpl.ExecuteTemplate(w, "list-cntrbH-h.gtpl", cntrbh_header); err != nil {
		log.Println(err)
	}

	cntrbhistory, status := SelectCntrbHistory(eventid, userno, tlsnid, acqtimelist)
	if status != 0 {
		return
	}

	if err := tpl.ExecuteTemplate(w, "list-cntrbH.gtpl", cntrbhistory); err != nil {
		log.Println(err)
	}

}

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

// 指定したリスナーの貢献ポイントの履歴を取得する。
func SelectCntrbHistory(
	eventid string,
	userno int,
	tlsnid int,
	acqtimelist []time.Time, //	配信者の配信時刻のリスト
) (
	cntrbhistory CntrbHistory,
	status int,
) {
	var stmt_er, stmt_tt, stmt_no *sql.Stmt
	//	var rows *sql.Rows

	sql_tt := "select stime, etime, target from timetable "
	sql_tt += " where eventid = ? and userid =? and sampletm2 = ? "
	stmt_tt, srdblib.Dberr = srdblib.Db.Prepare(sql_tt)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbHistory() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt_tt.Close()

	sql_no := "select count(*) from eventrank "
	sql_no += " where eventid = ? and userid =? and t_lsnid = ? and ts = ? "
	stmt_no, srdblib.Dberr = srdblib.Db.Prepare(sql_no)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbHistory() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt_no.Close()

	sql_er := "select point, increment, listner, lastname from eventrank "
	sql_er += " where eventid = ? and userid =? and t_lsnid = ? and ts = ? "
	stmt_er, srdblib.Dberr = srdblib.Db.Prepare(sql_er)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbHistory() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt_er.Close()

	var chi CntrbHistoryInf
	var stime, etime time.Time

	for _, ts := range acqtimelist {

		srdblib.Dberr = stmt_tt.QueryRow(eventid, userno, ts).Scan(&stime, &etime, &chi.Target)
		if srdblib.Dberr != nil {
			log.Printf("%s(%s, %d, %+v)\n", sql_tt, eventid, userno, ts)
			log.Printf("err=[%s]\n", srdblib.Dberr.Error())
			status = -11
		}

		if chi.Target == -1 {
			nfr := (int(etime.Sub(stime).Minutes()) + 45) / 60
			chi.Target, _ = InsertTargetIntoTimtable(eventid, userno, ts, nfr)
		}

		chi.S_stime = stime.Format("01/02 15:04")
		chi.S_etime = etime.Format("01/02 15:04")

		no := 0
		srdblib.Dberr = stmt_no.QueryRow(eventid, userno, tlsnid, ts).Scan(&no)
		if srdblib.Dberr != nil {
			log.Printf("%s(%s, %d, %d, %+v)\n", sql_er, eventid, userno, tlsnid, ts)
			log.Printf("err=[%s]\n", srdblib.Dberr.Error())
			status = -11
		}

		if no != 0 {

			srdblib.Dberr = stmt_er.QueryRow(eventid, userno, tlsnid, ts).Scan(&chi.Point, &chi.Incremental, &chi.Listener, &chi.Lastname)
			if srdblib.Dberr != nil {
				log.Printf("%s(%s, %d, %d, %+v)\n", sql_er, eventid, userno, tlsnid, ts)
				log.Printf("err=[%s]\n", srdblib.Dberr.Error())
				status = -11
			}
			lnl := strings.Split(chi.Lastname, "[")
			if len(lnl) == 2 {
				chi.Lastname = strings.TrimRight(lnl[1], "]")
			}
		} else {
			chi.Point = -1
			chi.Incremental = -1
		}

		cntrbhistory = append(cntrbhistory, chi)
	}

	for i := len(cntrbhistory) - 1; i > 0; i-- {

		if cntrbhistory[i].Listener == cntrbhistory[i-1].Listener {
			//	リスナー名に変化がなければ空白とする。
			cntrbhistory[i].Listener = ""
		}
		/*
			if cntrbhistory[i].Lastname == cntrbhistory[i-1].Lastname {
				cntrbhistory[i].Lastname = ""
			}
		*/

	}

	return
}

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
