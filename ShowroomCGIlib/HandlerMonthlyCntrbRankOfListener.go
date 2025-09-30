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
	Limit                int
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
	Spoint    int
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
	if monthlyCntrbRank.YearMonth == "" {
		// デフォルトは前月
		tnow := time.Now()
		ly := tnow.Year()
		lm := int(tnow.Month()) - 1
		// 現在の年月が1月の場合、前月は前年の12月になる
		if lm == 0 {
			ly--
			lm = 12
		}
		monthlyCntrbRank.YearMonth = fmt.Sprintf("%04d-%02d", ly, lm)
		log.Printf("Default yearmonth: %s\n", monthlyCntrbRank.YearMonth)
	}
	monthlyCntrbRank.Year, monthlyCntrbRank.Month, err = getYearMonth(monthlyCntrbRank.YearMonth)
	if err != nil {
		log.Printf("Invalid yearmonth: %s (%s)\n", monthlyCntrbRank.YearMonth, err.Error())
		monthlyCntrbRank.Year = time.Now().Year()
		monthlyCntrbRank.Month = int(time.Now().Month()) - 1
		if monthlyCntrbRank.Month == 0 {
			monthlyCntrbRank.Month = 12
			monthlyCntrbRank.Year--
		}
		fmt.Fprintf(w, "Invalid yearmonth: %s (%s)\n", monthlyCntrbRank.YearMonth, err.Error())
		return
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
	lpoint := 0
	llsnid := 0
	lroomno := 0
	lieventid := 0
	// for i, v := range monthlyCntrbRank.MonthlyCntrbRankList {
	for i := len(monthlyCntrbRank.MonthlyCntrbRankList) - 1; i >= 0; i-- {
		v := monthlyCntrbRank.MonthlyCntrbRankList[i]
		// w.Write([]byte(fmt.Sprintf("%d: %s %s %d %s %s %s %d\n",
		// 	v.Lsnid, v.Listener, v.Eventid, v.Ieventid, v.Eventname, v.Longname, v.Endtime.Format("2006-01-02 15:04"), v.Point)))
		if v.Point == lpoint && llsnid == v.Lsnid && lroomno == v.Roomno && lieventid == v.Ieventid {
			if i == len(monthlyCntrbRank.MonthlyCntrbRankList)-2 {
				monthlyCntrbRank.MonthlyCntrbRankList = monthlyCntrbRank.MonthlyCntrbRankList[:i]
			} else {
				monthlyCntrbRank.MonthlyCntrbRankList =
					append(monthlyCntrbRank.MonthlyCntrbRankList[:i+1], monthlyCntrbRank.MonthlyCntrbRankList[i+2:]...)
			}
		} else {
			lpoint = v.Point
			llsnid = v.Lsnid
			lroomno = v.Roomno
			lieventid = v.Ieventid
			log.Printf("%d: %s %s %d %s %s %s %d\n",
				v.Lsnid, v.Listener, v.Eventid, v.Ieventid, v.Eventname, v.Longname, v.Endtime.Format("2006-01-02 15:04"), v.Point)
		}
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

// 指定した年月におけるリスナーの貢献ポイント合計のランキングのリストを作成する
func MonthlyCntrbRankLgHandler(w http.ResponseWriter, req *http.Request) {

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
	if monthlyCntrbRank.YearMonth == "" {
		// デフォルトは前月
		tnow := time.Now()
		ly := tnow.Year()
		lm := int(tnow.Month()) - 1
		// 現在の年月が1月の場合、前月は前年の12月になる
		if lm == 0 {
			ly--
			lm = 12
		}
		monthlyCntrbRank.YearMonth = fmt.Sprintf("%04d-%02d", ly, lm)
		log.Printf("Default yearmonth: %s\n", monthlyCntrbRank.YearMonth)
	}
	monthlyCntrbRank.Year, monthlyCntrbRank.Month, err = getYearMonth(monthlyCntrbRank.YearMonth)
	if err != nil {
		log.Printf("Invalid yearmonth: %s (%s)\n", monthlyCntrbRank.YearMonth, err.Error())
		monthlyCntrbRank.Year = time.Now().Year()
		monthlyCntrbRank.Month = int(time.Now().Month()) - 1
		if monthlyCntrbRank.Month == 0 {
			monthlyCntrbRank.Month = 12
			monthlyCntrbRank.Year--
		}
		fmt.Fprintf(w, "Invalid yearmonth: %s (%s)\n", monthlyCntrbRank.YearMonth, err.Error())
		return
	}

	monthlyCntrbRank.Thpoint, err = strconv.Atoi(req.FormValue("thpoint"))
	if err != nil {
		monthlyCntrbRank.Thpoint = 30000 // デフォルト値
		log.Printf("Invalid thpoint: %s (%s)\n", req.FormValue("thpoint"), err.Error())
	}
	if monthlyCntrbRank.Thpoint < 30000 {
		monthlyCntrbRank.Thpoint = 30000 // デフォルト値
	}

	monthlyCntrbRank.Limit, err = strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		monthlyCntrbRank.Limit = 10 // デフォルト値
		log.Printf("Invalid limit: %s (%s)\n", req.FormValue("limit"), err.Error())
	}
	if monthlyCntrbRank.Limit < 10 {
		monthlyCntrbRank.Limit = 10
	}
	if monthlyCntrbRank.Limit > 50 {
		monthlyCntrbRank.Limit = 50
	}

	// 貢献ポイントのリストを取得する
	monthlyCntrbRank.MonthlyCntrbRankList, err = selectMonthlyCntrbRankLg(
		monthlyCntrbRank.Year,
		monthlyCntrbRank.Month,
		monthlyCntrbRank.Thpoint,
		monthlyCntrbRank.Limit,
	)
	if err != nil {
		err = fmt.Errorf("SelectCntrbHistoryEx() error: %w", err)
		log.Printf("%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	lpoint := 0
	llsnid := 0
	lroomno := 0
	lieventid := 0
	// for i, v := range monthlyCntrbRank.MonthlyCntrbRankList {
	for i := len(monthlyCntrbRank.MonthlyCntrbRankList) - 1; i >= 0; i-- {
		v := monthlyCntrbRank.MonthlyCntrbRankList[i]
		// w.Write([]byte(fmt.Sprintf("%d: %s %s %d %s %s %s %d\n",
		// 	v.Lsnid, v.Listener, v.Eventid, v.Ieventid, v.Eventname, v.Longname, v.Endtime.Format("2006-01-02 15:04"), v.Point)))
		if v.Point == lpoint && llsnid == v.Lsnid && lroomno == v.Roomno && lieventid == v.Ieventid {
			if i == len(monthlyCntrbRank.MonthlyCntrbRankList)-2 {
				monthlyCntrbRank.MonthlyCntrbRankList = monthlyCntrbRank.MonthlyCntrbRankList[:i]
			} else {
				monthlyCntrbRank.MonthlyCntrbRankList =
					append(monthlyCntrbRank.MonthlyCntrbRankList[:i+1], monthlyCntrbRank.MonthlyCntrbRankList[i+2:]...)
			}
		} else {
			lpoint = v.Point
			llsnid = v.Lsnid
			lroomno = v.Roomno
			lieventid = v.Ieventid
			log.Printf("%d: %s %s %d %s %s %s %d\n",
				v.Lsnid, v.Listener, v.Eventid, v.Ieventid, v.Eventname, v.Longname, v.Endtime.Format("2006-01-02 15:04"), v.Point)
		}
	}

	// テンプレートをパースする
	funcMap := template.FuncMap{
		"sub":        func(i, j int) int { return i - j },
		"Comma":      func(i int) string { return humanize.Comma(int64(i)) },
		"FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
		"FormatInt":  func(i int, tfmt string) string { return fmt.Sprintf(tfmt, i) },
	}

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/m-cntrbrank-Lg.gtpl"))

	if err = tpl.ExecuteTemplate(w, "m-cntrbrank-Lg.gtpl", monthlyCntrbRank); err != nil {
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
	ly := tnow.Year()
	lm := int(tnow.Month())
	// 現在の年月をYYYYMM形式で求める
	cym := ly*100 + lm

	// 前月の年、月を求める。現在の年月が1月の場合、前月は前年の12月になる
	lm -= 1
	if lm == 0 {
		ly--
		lm = 12
	}

	parts := strings.Split(yearmonth, "-")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid yearmonth format: expected YYYY-MM, got: %s", yearmonth)
		return
	}

	year, err = strconv.Atoi(parts[0])
	if err != nil || year < 2024 || year > 2100 {
		if err != nil {
			err = fmt.Errorf("invalid year: %w: year must be a number", err)
		} else {
			err = fmt.Errorf("rnvalid year: %d: year must be between 2024 and %d", year, ly)
		}
		return
	}

	month, err = strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		if err != nil {
			err = fmt.Errorf("invalid month: %w: month must be a number", err)
		} else {
			err = fmt.Errorf("invalid month: %d: month must be between 1 and 12", month)
		}
		return
	}

	if (year*100+month) >= cym || (year*100+month) < 202410 {
		err = fmt.Errorf("invalid yearmonth: %s ( year-month must be between 2024-10 and %4d-%02d)", yearmonth, ly, lm)
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

// 指定したリスナーの貢献ポイントの合計を取得する。
func selectMonthlyCntrbRankLg(
	year int,
	month int,
	thpoint int,
	limit int, // デフォルトは100件
) (
	monthlyCntrbRankList []MonthlyCntrbRankData,
	err error,
) {

	sqlst := `
WITH RankedEventRank AS (
-- SELECT er.lsnid, SUBSTRING(v.name,1,20) listener,
   SELECT er.lsnid,
                 e.eventid, e.ieventid, e.event_name eventname, e.starttime, e.endtime,
            --   er.userid roomno, u.userid, u.longname,
                 er.userid roomno,
                 er.point, max_ts,
                 ROW_NUMBER() OVER(PARTITION BY e.ieventid, er.userid, er.lsnid ORDER BY e.eventid ) as rn
         FROM eventrank er
         JOIN (
                 SELECT er.eventid, e.ieventid, userid, lsnid, MAX(ts) AS max_ts
                        FROM eventrank er
                        JOIN wevent e ON er.eventid = e.eventid
                        WHERE e.endtime between  ? and ?     -- 年月の初日と翌月の初日
                           AND e.achk = 0
                           AND er.point > ?  -- 貢献ポイントの閾値
                        GROUP BY eventid, userid, lsnid
         ) AS sub
                 ON er.eventid = sub.eventid
        AND er.userid = sub.userid
                 AND er.lsnid = sub.lsnid
                 AND er.ts = sub.max_ts
         JOIN event e
                 ON er.eventid = e.eventid
    --   JOIN viewer v
    --           ON er.lsnid = v.viewerid
    --   JOIN user u
    --           ON er.userid = u.userno
         WHERE er.point > ?     -- 貢献ポイントの閾値
         ORDER BY er.point desc
),

RankOfSum AS (
SELECT  SUM(point) sump, lsnid
         FROM RankedEventRank
         WHERE rn = 1
    --  GROUP BY lsnid, listener
        GROUP BY lsnid
        ORDER BY sump desc
        LIMIT ?    -- 上位何件を取得するか
)

SELECT rs.sump spoint, rs.lsnid , v.name listener, sum(er.point) point, er.userid roomno, u.longname, u.userid
        FROM eventrank er
        JOIN RankOfSum AS rs
                ON rs.lsnid = er.lsnid
        JOIN wevent AS e ON e.eventid = er.eventid
        JOIN RankedEventRank rer
                ON rer.eventid = er.eventid
                        AND rer.roomno = er.userid
                        AND rer.lsnid = er.lsnid
                        AND rer.max_ts = er.ts
                        AND rer.rn = 1
		JOIN viewer v ON er.lsnid = v.viewerid
    	JOIN user u ON er.userid = u.userno
        WHERE e.endtime BETWEEN ? and ? AND e.achk =0 -- 年月の初日と翌月の初日
        GROUP BY spoint, rs.lsnid, listener, roomno, u.longname, u.userid
        ORDER BY spoint DESC, point DESC
`

	// 月の初日と翌月の初日を求める
	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(0, 1, 0)

	_, err = srdblib.Dbmap.Select(&monthlyCntrbRankList, sqlst,
		startTime, endTime, thpoint, thpoint, limit, startTime, endTime)
	if err != nil {
		err = fmt.Errorf("SelectCntrbHistoryEx() error: %w", err)
		log.Printf("%s\n", err.Error())
		return nil, err
	}

	return
}

/*
WITH RankedEventRank AS (
SELECT er.lsnid, SUBSTRING(v.name,1,20) listener,
                 e.eventid, e.ieventid, e.event_name eventname, e.starttime, e.endtime,
                 u.userno roomno, u.userid, u.longname,
                 er.point,
                 ROW_NUMBER() OVER(PARTITION BY e.ieventid, er.userid, er.lsnid ORDER BY e.eventid ) as rn
         FROM eventrank er
         JOIN (
                 SELECT er.eventid, e.ieventid, userid, lsnid, MAX(ts) AS max_ts
                        FROM eventrank er
                        JOIN wevent e ON er.eventid = e.eventid
                        WHERE e.endtime between  ? and ?
                           AND e.achk =0
                           AND er.point > ?
                        GROUP BY eventid, userid, lsnid
         ) AS sub
                 ON er.eventid = sub.eventid
        AND er.userid = sub.userid
                 AND er.lsnid = sub.lsnid
                 AND er.ts = sub.max_ts
         JOIN event e
                 ON er.eventid = e.eventid
         JOIN viewer v
                 ON er.lsnid = v.viewerid
         JOIN user u
                 ON er.userid = u.userno
         WHERE er.point > ?
 -- T.   AND er.lsnid = 9153337
         ORDER BY er.point desc )
 -- A. SELECT  SUM(point) point, lsnid,  listener
 SELECT  SUM(point) point, lsnid,  listener
 -- B. SELECT  TRUNCATE(point+5000, -4) pt, lsnid,  listener, roomno, SUBSTRING(eventname,1,30)
 -- C. SELECT  SUM(point) point, lsnid,  listener, roomno, longname SELECT  SUM(point) point, lsnid,
 -- C.            listener, roomno, longname
         FROM RankedEventRank
         WHERE rn = 1
 -- A.       GROUP BY lsnid, listener
 -- A.       ORDER BY point desc
 -- A.       LIMIT ?;
        GROUP BY lsnid, listener
        ORDER BY point desc
        LIMIT ?;
 -- B.       ORDER BY roomno, point
 -- B.       LIMIT ?;
 -- C.       GROUP BY lsnid, listener, roomno, longname
 -- C.       ORDER BY point desc

 --------------------------------------------------------------

WITH RankedEventRank AS (
SELECT er.lsnid, SUBSTRING(v.name,1,20) listener,
                 e.eventid, e.ieventid, e.event_name eventname, e.starttime, e.endtime,
                 u.userno roomno, u.userid, u.longname,
                 er.point, max_ts,
                 ROW_NUMBER() OVER(PARTITION BY e.ieventid, er.userid, er.lsnid ORDER BY e.eventid ) as rn
         FROM eventrank er
         JOIN (
                 SELECT er.eventid, e.ieventid, userid, lsnid, MAX(ts) AS max_ts
                        FROM eventrank er
                        JOIN wevent e ON er.eventid = e.eventid
                        WHERE e.endtime between  '2025-06-01' and '2025-07-01'
                           AND e.achk =0
                           AND er.point > 30000
                        GROUP BY eventid, userid, lsnid
         ) AS sub
                 ON er.eventid = sub.eventid
        AND er.userid = sub.userid
                 AND er.lsnid = sub.lsnid
                 AND er.ts = sub.max_ts
         JOIN event e
                 ON er.eventid = e.eventid
         JOIN viewer v
                 ON er.lsnid = v.viewerid
         JOIN user u
                 ON er.userid = u.userno
         WHERE er.point > 30000
         ORDER BY er.point desc
),

RankOfSum AS (
SELECT  SUM(point) sump, lsnid
         FROM RankedEventRank
         WHERE rn = 1
        GROUP BY lsnid, listener
        ORDER BY sump desc
        LIMIT 10
)

SELECT rs.sump, rs.lsnid , er.userid, er.point, e.eventid
        FROM eventrank er
        JOIN RankOfSum AS rs
                ON rs.lsnid = er.lsnid
        JOIN wevent AS e
                ON e.eventid = er.eventid
        JOIN RankedEventRank rer
                ON rer.eventid = er.eventid
                        AND rer.max_ts = er.ts
                        AND rer.roomno = er.userid
                        AND rer.lsnid = er.lsnid
                        AND rer.rn = 1
        WHERE e.endtime BETWEEN  '2025-06-01' and '2025-07-01'
                           AND e.achk =0
        ORDER BY sump DESC, point DESC
*/
