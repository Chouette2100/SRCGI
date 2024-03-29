package ShowroomCGIlib

import (
	"fmt"
	"log"
	"time"

	//	"math"
	//	"sort"
	//	"strconv"
	//	"strings"
	//	"time"

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
	//	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

/*
        HandlerListCntrbH()
		指定されたイベント、配信者、リスナーの貢献ポイントの履歴を表示する。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。
	0101J1	

*/

func HandlerFlRanking(w http.ResponseWriter, req *http.Request) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if ! isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}



	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/flranking.gtpl",
	))

	eventid := req.FormValue("eventid")
	if eventid == "" {
		eventid = "funlevel"
	}

	log.Printf("/flranking eventid=%s\n", eventid)

	eventandrankinginf, status := SelectFromNoOfFan(eventid)
	if status != 0 {
		log.Printf("SelectFromFluser() returned %d\n", status)
		return
	}

	if err := tpl.ExecuteTemplate(w, "flranking.gtpl", eventandrankinginf); err != nil {
		log.Println(err)
	}

}

type RankingInf struct {
	Room_id   int
	Room_name string
	Srank     string
	Irank     int
	Irorder   int
	Iorder    int
	Fans      int
	Fans_lst  int
}

type EventAndRankingInf struct {
	Eventid        string
	Eventname      string
	Period         string
	Ts_lst         string
	Ts_nxt         string
	RankingInfList []RankingInf
}

func SelectFromNoOfFan(eventid string) (eventandrankinginf EventAndRankingInf, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	var maxts time.Time
	sql := "select max(ts) from nooffan where eventid = ? "
	Err := srdblib.Db.QueryRow(sql, eventid).Scan(&maxts)
	if Err != nil {
		log.Printf("%s\n", sql)
		log.Printf("** SelectFromNoOfFan() err=[%s]\n", Err.Error())
		//	if err.Error() != "sql: no rows in result set" {
		status = -1
		return
		//	}
	}

	eventandrankinginf.Ts_lst = maxts.Format("2006-01-02 15:04")
	eventandrankinginf.Ts_nxt = maxts.Add(3 * time.Hour).Format("2006-01-02 15:04")
	log.Printf(" maxts=%s\n", eventandrankinginf.Ts_lst)

	eventandrankinginf.RankingInfList = make([]RankingInf, 0)

	var eventinf exsrapi.Event_Inf

	status = GetEventInf(eventid, &eventinf)
	log.Printf("eventinf = %+v\n", eventinf)

	eventandrankinginf.Eventid = eventid
	eventandrankinginf.Eventname = eventinf.Event_name
	eventandrankinginf.Period = eventinf.Period

	sql = "select roomid, roomname, irank, srank, iorder, fans, fans_lst from nooffan where ts = ? and eventid = ? order by irank desc, fans desc"
	stmt, Err = srdblib.Db.Prepare(sql)
	if Err != nil {
		log.Printf("** SelectFromNoOfFan() err=[%s]\n", Err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, Err = stmt.Query(maxts, eventid)
	if Err != nil {
		log.Printf("** SelectFromNoOfFan() err=[%s]\n", Err.Error())
		status = -2
		return
	}
	defer rows.Close()

	var rankinginf RankingInf

	irorder := 0
	lstfans := -1
	lstrank := 9
	for rows.Next() {
		Err = rows.Scan(&rankinginf.Room_id, &rankinginf.Room_name, &rankinginf.Irank, &rankinginf.Srank, &rankinginf.Iorder, &rankinginf.Fans, &rankinginf.Fans_lst)
		if Err != nil {
			log.Printf("** SelectFromNoOfFan() err=[%s]\n", Err.Error())
			status = -3
			return
		}
		if rankinginf.Irank != lstrank {
			irorder = 0
			lstrank = rankinginf.Irank
		}
		irorder++
		if rankinginf.Fans != lstfans {
			rankinginf.Irorder = irorder
			lstfans = rankinginf.Fans
		} else {
			rankinginf.Irorder = -1
		}

		eventandrankinginf.RankingInfList = append(eventandrankinginf.RankingInfList, rankinginf)
	}
	if Err = rows.Err(); Err != nil {
		log.Printf("** SelectFromNoOfFan() err=[%s]\n", Err.Error())
		status = -4
		return
	}

	return

}
