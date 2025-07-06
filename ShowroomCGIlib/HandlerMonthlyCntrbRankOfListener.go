// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
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
	"github.com/dustin/go-humanize"

	// "github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type MonthlyCntrbRank struct {
	YearMonth            string // YYYY-MM形式
	Year                 int
	Month                int
	Thpoint              int
	Thlist               []int
	MonthlyCntrbRankList []MonthlyCntrbRankData
}

type MonthlyCntrbRankData struct {
	Lsnid     int
	Listener  string
	Eventid   string
	Ieventid  int
	Eventname string
	Starttime time.Time
	Endtime   time.Time
	Roomno    int
	Userid    string
	Longname  string
	Point     int
}

// 指定した年月におけるイベント=ルームに対する貢献ポイントのランキングのリストを作成する
func MonthlyCntrbRankOfListenerHandler(w http.ResponseWriter, req *http.Request) {

	var monthlyCntrbRank MonthlyCntrbRank
	var err error

	monthlyCntrbRank.Thlist = []int{900000, 1800000, 2700000}

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// yyyy-mm形式の年月を取得する。
	monthlyCntrbRank.YearMonth = req.FormValue("yearmonth")
	monthlyCntrbRank.Year, monthlyCntrbRank.Month, err = getYearMonth(monthlyCntrbRank.YearMonth)
	if err != nil {
		// fmt.Fprintf(w, "Invalid yearmonth: %s (%s)\n", monthlyCntrbRank.YearMonth, err.Error())
		log.Printf("Invalid yearmonth: %s (%s)\n", monthlyCntrbRank.YearMonth, err.Error())
		monthlyCntrbRank.Year = time.Now().Year()
		monthlyCntrbRank.Month = int(time.Now().Month()) - 1
		if monthlyCntrbRank.Month == 0 {
			monthlyCntrbRank.Month = 12
			monthlyCntrbRank.Year--
		}
		// return
	}

	monthlyCntrbRank.Thpoint, err = strconv.Atoi(req.FormValue("thpoint"))
	if err != nil {
		monthlyCntrbRank.Thpoint = 2700000 // デフォルト値
		log.Printf("Invalid thpoint: %s (%s)\n", req.FormValue("thpoint"), err.Error())
	}

	// 貢献ポイントのリストを取得する
	monthlyCntrbRank.MonthlyCntrbRankList, err = selectMonthlyCntrbRank(
		monthlyCntrbRank.Year,
		monthlyCntrbRank.Month,
		monthlyCntrbRank.Thpoint,
	)
	if err != nil {
		err = fmt.Errorf("SelectCntrbHistoryEx() error: %w", err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	for _, v := range monthlyCntrbRank.MonthlyCntrbRankList {
		// w.Write([]byte(fmt.Sprintf("%d: %s %s %d %s %s %s %d\n",
		// 	v.Lsnid, v.Listener, v.Eventid, v.Ieventid, v.Eventname, v.Longname, v.Endtime.Format("2006-01-02 15:04"), v.Point)))
		log.Printf("%d: %s %s %d %s %s %s %d\n",
			v.Lsnid, v.Listener, v.Eventid, v.Ieventid, v.Eventname, v.Longname, v.Endtime.Format("2006-01-02 15:04"), v.Point)
	}

	// テンプレートをパースする
	funcMap := template.FuncMap{
		"sub":        func(i, j int) int { return i - j },
		"Comma":      func(i int) string { return humanize.Comma(int64(i)) },
		"FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
		"FormatInt":  func(i int, tfmt string) string { return fmt.Sprintf(tfmt, i) },
	}

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/m-cntrbrank-listener.gtpl"))

	if err = tpl.ExecuteTemplate(w, "m-cntrbrank-listener.gtpl", monthlyCntrbRank); err != nil {
		log.Println(err)
	}

}

// GetYearMonth は、YYYY-MM形式の文字列を受け取り、年と月を整数として返す。
// 年は2000年から2100年の範囲で、月は1から12の範囲でなければならない。
// また、現時点でまだ終わっていない年月は無効とする。
// もし無効な形式や値が与えられた場合は、エラーを返す。
// GetYearMonthは、年と月を整数として返す。
// もし無効な形式や値が与えられた場合は、エラーを返す。
// 例: "2023-05" -> (2023, 5, nil)
//
//	"2023-13" -> (0, 0, error)
//	"2025-01" -> (2025, 1, nil)
//	"2025-06" -> (2025, 6, nil)
//	"2025-07" -> (0, 0, error) // 2025年7月はまだ終わっていない年月なので無効(現在＝2025-07-06)
//	"2025-12" -> (0, 0, error) // 2025年12月は未来の年月なので無効
func getYearMonth(
	yearmonth string,
) (year int, month int, err error) {

	tnow := time.Now()
	cym := tnow.Year()*100 + int(tnow.Month())

	parts := strings.Split(yearmonth, "-")
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid yearmonth format")
		return
	}

	year, err = strconv.Atoi(parts[0])
	if err != nil || year < 2000 || year > 2100 {
		err = fmt.Errorf("Invalid year: %w", err)
		return
	}

	month, err = strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		err = fmt.Errorf("Invalid month: %w", err)
		return
	}

	if (year*100 + month) >= cym {
		err = fmt.Errorf("Invalid yearmonth: %s", yearmonth)
		return
	}

	return
}

// 指定したリスナーの貢献ポイントの履歴を取得する。
func selectMonthlyCntrbRank(
	year int,
	month int,
	thpoint int,
) (
	monthlyCntrbRankList []MonthlyCntrbRankData,
	err error,
) {

	sqlst := `
SELECT er.lsnid, SUBSTRING(v.name,1,20) listener,
       e.eventid, e.ieventid, e.event_name eventname, e.starttime, e.endtime,
       u.userno roomno, u.userid, u.longname,
       er.point
        FROM eventrank er
        JOIN (
                SELECT er.eventid, userid, MAX(ts) AS max_ts
                        FROM eventrank er
                        join event e on er.eventid = e.eventid
                        WHERE e.endtime between  ? and ?
                        GROUP BY eventid, userid
        ) AS sub
                ON er.eventid = sub.eventid
                AND er.userid = sub.userid
                AND er.ts = sub.max_ts
        JOIN event e
                ON er.eventid = e.eventid
        JOIN viewer v
                ON er.lsnid = v.viewerid
        JOIN user u
                ON er.userid = u.userno
        WHERE er.point > ?
        ORDER BY er.point desc; `

	// 月の初日と翌月の初日を求める
	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)

	_, err = srdblib.Dbmap.Select(&monthlyCntrbRankList, sqlst, startTime, endTime, thpoint)
	if err != nil {
		err = fmt.Errorf("SelectCntrbHistoryEx() error: %w", err)
		log.Printf("%s\n", err.Error())
		return nil, err
	}

	return
}
