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
	"github.com/Chouette2100/srdblib"

	"github.com/Chouette2100/exsrapi"
)

const MaxAcq = 5

type CntrbHeader struct {
	Eventid		string
	Eventname	string
	Period		string
	Userno		int
	Username	string
	ShortURL	string
	Ier			int
	Iel			int
	S_stime		[]string
	S_etime		[]string
	Earned		[]int
	Total		[]int
	Target		[]int
	Ifrm		[]int
	Nof			[]int
	Nft			int
	Npb			int
	N1b			int
	N1f			int
	Npf			int
	Nlt			int
}

type CntrbInf struct {
	Ranking      int
	Point        int
	Incremental  []int
	ListenerName string
	LastName     string
	Tlsnid		 int
	Eventid		 string
	Userno		 int
}
//	type	CntrbInfList	[] CntrbInf

/*
        HandlerListCntrb()
		（配信枠別）貢献ポイントランキングを5(MaxAcq)枠分表示する。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。

*/

func HandlerListCntrb(w http.ResponseWriter, req *http.Request) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	GetUserInf(req)

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrb-h1.gtpl","templates/list-cntrb-h2.gtpl","templates/list-cntrb.gtpl"))
    funcMap := template.FuncMap{
        "sub": func(i, j int) int { return i - j },
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
    }
    tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-cntrb-h1.gtpl","templates/list-cntrb-h2.gtpl","templates/list-cntrb.gtpl"))

	

	eventid := req.FormValue("eventid")

	var eventinf exsrapi.Event_Inf
	GetEventInf(eventid, &eventinf)

	userno, _ := strconv.Atoi(req.FormValue("userno"))

	acqtimelist, _ := SelectAcqTimeList(eventid, userno)
	latl := len(acqtimelist)

	sie := req.FormValue("ie")
	ie := 0
	if sie != "" {
		ie, _ = strconv.Atoi(sie)
	} else {
		ie = latl - 1
	}

	log.Printf("***** HandlerListCntrb() called. eventid=%s, userno=%d, ie=%d\n", eventid, userno, ie)

	ib := ie -4
	if ib < 0 {
		ib = 0
	}

	var cntrbheader CntrbHeader

	cntrbheader.Eventid = eventid
	cntrbheader.Eventname = eventinf.Event_name
	cntrbheader.Period = eventinf.Period
	cntrbheader.Userno = userno

	//	戻る側の設定
	if ie < MaxAcq {
		cntrbheader.Nft = -1
		cntrbheader.Npb = -1
		cntrbheader.N1b = -1
	} else {
		cntrbheader.Nft = MaxAcq - 1		//	先頭にもどる
		cntrbheader.Npb = ie - MaxAcq		//	１ページ分戻る
		if cntrbheader.Npb < MaxAcq -1 {
			cntrbheader.Npb = MaxAcq - 1
		}
		cntrbheader.N1b = ie - 1		//	一枠分戻る
	}

	if ie == latl - 1  {
		cntrbheader.Nlt = -1
		cntrbheader.Npf = -1
		cntrbheader.N1f = -1
	} else {
		cntrbheader.Nlt = latl -1			//	最後に進む
		cntrbheader.Npf = ie + MaxAcq		//	１ページ分進む
		if cntrbheader.Npf > latl - 1 {
			cntrbheader.Npf = latl - 1
		}
		cntrbheader.N1f = ie + 1			//	一枠分進む
	}

	_, _, _, _, _, _, _, _, roomname, roomurlkey, _, _  := GetRoomInfoByAPI(fmt.Sprintf("%d", userno)) 
	cntrbheader.Username = roomname
	cntrbheader.ShortURL = roomurlkey

	ts := acqtimelist[ie]

	cntrbinflist, tlsnid2order, status := SelectTlsnid2Order(eventid, userno, ts)
	if status != 0 {
		log.Printf(" SelectCntrbLast() returned %d in HandlerListCntrb()\n", status)
		return
	}


	for  i := ib; i <= ie; i++ {
		ts = acqtimelist[i]
		//	log.Printf(" i=%d ts=%+v\n", i, ts)
		status = SelectCntrb(eventid, userno, ts, &cntrbinflist, tlsnid2order)
		if status != 0 {
			log.Printf(" SelectCntrb() returned %d in HandlerListCntrb()\n", status)
			return
		}
		SelectCntrbHeader(eventid, userno, ts, &cntrbheader)
		cntrbheader.Ifrm[i-ib] = i
		cntrbheader.Nof[i-ib] = i + 1
	}

	if ie == latl -1 {
		cntrbheader.Ier = -1
	} else {
		cntrbheader.Ier = ie + 5
		if cntrbheader.Ier > latl - 1 {
			cntrbheader.Ier = latl -1
		}
	}

	if ie == 0 {
		cntrbheader.Ier = -1
	} else {
		cntrbheader.Iel = ie - 5
		if cntrbheader.Iel < 0 {
			cntrbheader.Iel = 0
		}
	}


	if err := tpl.ExecuteTemplate(w, "list-cntrb-h1.gtpl", cntrbheader); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "list-cntrb-h2.gtpl", cntrbheader); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "list-cntrb.gtpl", cntrbinflist); err != nil {
		log.Println(err)
	}
}

/*
        SelectAcqTimeList()
		指定したイベント、ユーザーの貢献ランキングを取得した時刻の一覧を取得する。

        引数
		eventid			string			イベントID
		userno			int				ユーザーID

        戻り値
        acqtimelist		[] time.Time	取得時刻一覧

*/
func SelectAcqTimeList(eventid string, userno int) (acqtimelist []time.Time, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	//	貢献ポイントランキングを取得した時刻の一覧を取得する。
	sql := "select sampletm2 from timetable where eventid = ? and userid = ? and status = 1 order by sampletm2"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sql)

	if srdblib.Dberr != nil {
		log.Printf("SelectAcqTimeList() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(eventid, userno)
	if srdblib.Dberr != nil {
		log.Printf("SelectAcqTimeList() (6) err=%s\n", srdblib.Dberr.Error())
		status = -6
		return
	}
	defer rows.Close()

	var ts time.Time

	for rows.Next() {
		srdblib.Dberr = rows.Scan(&ts)
		if srdblib.Dberr != nil {
			log.Printf("SelectAcqTimeList() (7) err=%s\n", srdblib.Dberr.Error())
			status = -7
			return
		}
		//	log.Printf("%+v\n", cntrbinf)
		acqtimelist = append(acqtimelist, ts)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("SelectAcqTimeList() (8) err=%s\n", srdblib.Dberr.Error())
		status = -8
		return
	}

	return

}


/*
		指定したイベント、ユーザー、時刻の貢献ポイントランキングを取得する。
		ここでは順位と累計貢献ポイントは取得しない。

        引数
		eventid			string			イベントID
		userno			int				ユーザーID
		ts				int				ユーザーID
		loc				int				取得データの格納位置
		loc				int				データの格納場所（ 0 だったら先頭）

        戻り値
        cntrbinflist	[] CntrbInf		貢献ポイントランキング（最終貢献ポイント順）
		stats			int				== 0 正常終了	!= 0 データベースアクセス時のエラー

*/
func SelectCntrb(
	eventid string,
	userno int,
	ts time.Time,
	cntrbinflist *[]CntrbInf,
	tlsnid2order map[int]int,
) (
	status int,
) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	//	最後の貢献ポイントランキングを取得する。
	sql := "select t_lsnid, increment from eventrank "
	sql += " where eventid = ? and userid =? and ts = ? order by norder"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sql)

	if srdblib.Dberr != nil {
		log.Printf("SelectCntrb() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(eventid, userno, ts)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrb() (6) err=%s\n", srdblib.Dberr.Error())
		status = -6
		return
	}
	defer rows.Close()

	tlsnid := 0
	incremental := 0

	for i := 0; i < len(*cntrbinflist); i++ {
		(*cntrbinflist)[i].Incremental = append((*cntrbinflist)[i].Incremental, -1) 
	}
	loc := len((*cntrbinflist)[0].Incremental) - 1

	for rows.Next() {
		srdblib.Dberr = rows.Scan(&tlsnid, &incremental)
		if srdblib.Dberr != nil {
			log.Printf("SelectCntrb() (7) err=%s\n", srdblib.Dberr.Error())
			status = -7
			return
		}
		//	log.Printf("%+v\n", cntrbinf)
		i := tlsnid2order[tlsnid]
		(*cntrbinflist)[i].Incremental[loc] = incremental
		//	(*cntrbinflist)[i].Incremental = append((*cntrbinflist)[i].Incremental, incremental) 
		//	log.Printf(" tlsnid=%d i=%d increment[%d]=%d\n", tlsnid, i, loc, incremental)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("SelectCntrb() (8) err=%s\n", srdblib.Dberr.Error())
		status = -8
		return
	}

	status = 0

	return

}

/*
        SelectCntrbHeader()
		貢献ランキング表のヘッダ部分に必要な配信開始・終了時刻を取得する。

        引数
		eventid			string			イベントID
		userno			int				配信者ID
		ts				time.Time		枠を特定する時刻（＝貢献ランキングを取得した時刻）
		cntrbheader		*CntrbHeader	配信開始・終了時刻を格納する構造体

        戻り値
		status			int				終了ステータス（ 0: 正常、　1: DBアクセスでの異常）


*/
func SelectCntrbHeader(
		eventid	string,
		userno	int,
		ts	time.Time,
		cntrbheader *CntrbHeader,
) (
	status int,
) {

		var stime, etime time.Time
		var earned, total int

		status = 0

		sql := "select stime, etime, earnedpoint, totalpoint from timetable where eventid = ? and userid = ? and sampletm2 = ? "
		srdblib.Dberr = srdblib.Db.QueryRow(sql, eventid, userno, ts).Scan(&stime, &etime, &earned, &total)

		if srdblib.Dberr != nil {
			log.Printf("select stime, etime from timetable where eventid = %s and userid = %d and sampletm2 = %+v\n", eventid, userno, ts)
			log.Printf("err=[%s]\n", srdblib.Dberr.Error())
			status = -11
			return
		}

		(*cntrbheader).S_stime = append((*cntrbheader).S_stime, stime.Format("02 15:04"))
		(*cntrbheader).S_etime = append((*cntrbheader).S_etime, etime.Format("02 15:04"))
		(*cntrbheader).Earned = append((*cntrbheader).Earned, earned)
		(*cntrbheader).Total = append((*cntrbheader).Total, total)
		(*cntrbheader).Ifrm = append((*cntrbheader).Ifrm, 0)
		(*cntrbheader).Nof = append((*cntrbheader).Nof, 0)

		return

}

/*
        SelectTlsnid2Order()
		指定したイベント、配信者、枠の仮リスナーIDと貢献ポイントランキングの対応表を作成する。

        引数
		eventid			string			イベントID
		userno			int				配信者ID
		ts				time.Time		枠を特定する時刻（＝貢献ランキングを取得した時刻）

        戻り値
		cntrbinflist	[]CntrbInf		貢献ポイントランキングを格納するための構造体の配列
										ここで累計貢献ポイント、リスナー名が格納される。
		tlsnid2order	map[int]int		仮リスナーIDと貢献ポイントランキングの対応表
		status			int				終了ステータス（ 0: 正常、　1: DBアクセスでの異常）


*/
func SelectTlsnid2Order(
	eventid string,
	userno int,
	ts time.Time,
) (
	cntrbinflist []CntrbInf,
	tlsnid2order map[int]int,
	status int,
) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	tlsnid2order = make(map[int]int)

	//	指定された時刻の貢献ポイントランキングを取得する。
	sql := "select norder, t_lsnid, point, listner, lastname from eventrank "
	sql += " where eventid = ? and userid =? and ts = ? order by norder"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sql)

	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbNow() (5) err=%s\n", srdblib.Dberr.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(eventid, userno, ts)
	if srdblib.Dberr != nil {
		log.Printf("SelectCntrbNow() (6) err=%s\n", srdblib.Dberr.Error())
		status = -6
		return
	}
	defer rows.Close()

	var cntrbinf CntrbInf

	cntrbinf.Eventid = eventid
	cntrbinf.Userno = userno

	i := 0
	for rows.Next() {
		//	Err = rows.Scan(&cntrbinf.Ranking, &cntrbinf.Tlsnid, &cntrbinf.Point, &cntrbinf.Incremental[loc], &cntrbinf.ListenerName, &cntrbinf.LastName)
		srdblib.Dberr = rows.Scan(&cntrbinf.Ranking, &cntrbinf.Tlsnid, &cntrbinf.Point, &cntrbinf.ListenerName, &cntrbinf.LastName)
		if srdblib.Dberr != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", srdblib.Dberr.Error())
			status = -7
			return
		}
		//	log.Printf("%+v\n", cntrbinf)
		cntrbinflist = append(cntrbinflist, cntrbinf)
		if i != 0 && cntrbinflist[i].Point == cntrbinflist[i-1].Point && cntrbinflist[i].Ranking != 999 {
			cntrbinflist[i].Ranking = -1
		}
		tlsnid2order[cntrbinf.Tlsnid] = i
		i++
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("GetCurrentScore() (8) err=%s\n", srdblib.Dberr.Error())
		status = -8
		return
	}

	status = 0

	return

}
