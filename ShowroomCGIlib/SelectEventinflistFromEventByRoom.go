package ShowroomCGIlib

import (
	"fmt"
	"log"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// 指定した条件に該当するイベントのリストを作る。
func SelectEventinflistFromEventByRoom(
	cond int, // 抽出条件	-1:終了したイベント、0: 開催中のイベント、1: 開催予定のイベント
	mode int, // 0: すべて、 1: データ取得中のものに限定
	userno int, // イベント名検索キーワード
	limit *int, // ページング制限
	offset int, // ページングオフセット
) (
	eventinflist []exsrapi.Event_Inf,
	err error,
) {

	//      テーブルは"w"で始まるものを操作の対象とする。
	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	tnow := time.Now().Truncate(time.Second)

	sqls := "select we.eventid,we.ieventid,we.event_name, we.period, we.starttime, we.endtime, we.noentry, we.intervalmin, we.modmin, we.modsec, "
	sqls += " we.Fromorder, we.Toorder, we.Resethh, we.Resetmm, we.Nobasis, we.Maxdsp, we.cmap, we.target, we.`rstatus`, we.maxpoint, we.achk , we.aclr "
	sqls += " from wevent we"
	if mode == 1 {
		sqls += " join event e on we.eventid = e.eventid "
	}
	sqls += " where we.eventid in "
	sqls += " (select weu.eventid from weventuser weu join wevent we on weu.eventid  = we.eventid "
	//	sqls += " where weu.userno = ? and we.endtime < ? "
	sqls += " where weu.userno = ? and we.starttime < ? "
	sqls += " union select eu.eventid from eventuser eu join event e on eu.eventid = e.eventid "
	//	sqls += " where eu.userno = ? and e.endtime < ? )  "
	sqls += " where eu.userno = ? and e.starttime < ? )  "
	sqls += " order by starttime desc, endtime desc limit ? offset ?"

	//	log.Printf("sql=[%s]\n", sqls)
	var stmts *sql.Stmt
	stmts, err = srdblib.Db.Prepare(sqls)
	if err != nil {
		err = fmt.Errorf("Prepare(sqls): %w", err)
		return
	}
	defer stmts.Close()

	var rows *sql.Rows
	rows, err = stmts.Query(userno, tnow, userno, tnow, limit, offset)
	if err != nil {
		err = fmt.Errorf("Query(userno, userno): %w", err)
		return
	}
	defer rows.Close()

	eventinflist = make([]exsrapi.Event_Inf, 0)
	eventinf := exsrapi.Event_Inf{}

	i := 0
	lastieid := -1
	for rows.Next() {

		err = rows.Scan(
			&eventinf.Event_ID,
			&eventinf.I_Event_ID,
			&eventinf.Event_name,
			&eventinf.Period,
			&eventinf.Start_time,
			&eventinf.End_time,
			&eventinf.NoEntry,
			&eventinf.Intervalmin,
			&eventinf.Modmin,
			&eventinf.Modsec,
			&eventinf.Fromorder,
			&eventinf.Toorder,
			&eventinf.Resethh,
			&eventinf.Resetmm,
			&eventinf.Nobasis,
			&eventinf.Maxdsp,
			&eventinf.Cmap,
			&eventinf.Target,
			&eventinf.Rstatus,
			&eventinf.Maxpoint,
			&eventinf.Achk,
			&eventinf.Aclr,
		)

		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return
			} else {
				err = fmt.Errorf("row.Exec(): %w", err)
				return
			}
		}

		//	log.Printf("eventno=%d\n", Event_inf.Event_no)

		start_date := eventinf.Start_time.Truncate(time.Hour).Add(-time.Duration(eventinf.Start_time.Hour()) * time.Hour)
		end_date := eventinf.End_time.Truncate(time.Hour).Add(-time.Duration(eventinf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

		//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

		eventinf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
		eventinf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - eventinf.Start_date

		eventinf.Gscale = eventinf.Maxpoint % 1000
		eventinf.Maxpoint = eventinf.Maxpoint - eventinf.Gscale

		log.Printf("eventinf=[%v]\n", eventinf)

		if eventinf.I_Event_ID == lastieid {
			if eventinf.Achk == 0 {
				eventinflist[i-1] = eventinf
			}
			// TODO: 2025-04-29 limitを変更することは妥当か？
			// *limit--
			// ----- ----------
		} else {
			eventinflist = append(eventinflist, eventinf)
			lastieid = eventinf.I_Event_ID
			i++
		}

		//	log.Printf("Start_data=%f Dperiod=%f\n", eventinf.Start_date, eventinf.Dperiod)
	}

	return
}
