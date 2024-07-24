package ShowroomCGIlib

import (
	"fmt"
	//	"log"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

// 指定した条件に該当するイベントのリストを作る。
func SelectEventinflistFromEvent(
	cond int, // 抽出条件	-1:終了したイベント、0: 開催中のイベント、1: 開催予定のイベント
	mode int, // 0: すべて、 1: データ取得中のものに限定
	keyword string, // イベント名検索キーワード
	kwevid string, // イベント検索キーワード
	limit int, //	select文のlimitに使用する値
	offset int, //	select文のoffsetに使用する値
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
	sqls += " we.Fromorder, we.Toorder, we.Resethh, we.Resetmm, we.Nobasis, we.Maxdsp, we.cmap, we.target, we.`rstatus`, we.maxpoint, we.achk "
	sqls += " from wevent we"
	if mode == 1 {
		sqls += " join event e on we.eventid = e.eventid "
	}
	sqls += " where we.achk = 0 "
	switch cond {
	case -1:
		sqls += " and we.endtime < ?"
	case 0:
		sqls += " and we.starttime < ? and we.endtime > ?"
	case 1:
		sqls += " and we.starttime > ?"
	default:
		err = fmt.Errorf("mode=%d is not valid", cond)
		return
	}
	if keyword != "" {
		sqls += " and we.event_name like ? "
	} else if kwevid != "" {
		sqls += " and we.eventid like ? "
	}

	switch cond {
	case -1: //	終了済みイベント
		//	sqls += " order by we.endtime desc, we.starttime desc "
		//	sqls += " order by we.starttime desc, we.endtime desc, we.event_name "
		sqls += " order by we.endtime desc, we.starttime desc, we.event_name "
	case 0: //	開催中のイベント
		sqls += " order by we.endtime, we.starttime, we.event_name "
	case 1: //	開催予定のイベント
		sqls += " order by we.starttime, we.endtime, we.event_name "
	}

	if cond == -1 {
		sqls += " limit ? offset ? "
	}
	//	log.Printf("sql=[%s]\n", sqls)
	var stmts *sql.Stmt
	stmts, srdblib.Dberr = srdblib.Db.Prepare(sqls)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("Prepare(sqls): %w", srdblib.Dberr)
		return
	}
	defer stmts.Close()

	var rows *sql.Rows

	kw := keyword
	if kwevid != "" {
		kw = kwevid
	}

	switch {
	case cond == 0 && kw == "":
		rows, srdblib.Dberr = stmts.Query(tnow, tnow)
	case cond == 0 && kw != "":
		rows, srdblib.Dberr = stmts.Query(tnow, tnow, "%"+kw+"%")
	case cond != 0 && kw == "":
		if cond == -1 {
			rows, srdblib.Dberr = stmts.Query(tnow, limit, offset)
		} else {
			rows, srdblib.Dberr = stmts.Query(tnow)
		}
	case cond != 0 && kw != "":
		if cond == -1 {	
			rows, srdblib.Dberr = stmts.Query(tnow, "%"+kw+"%", limit, offset)
		} else {
			rows, srdblib.Dberr = stmts.Query(tnow, "%"+kw+"%")
		}
	}
	if srdblib.Dberr != nil {
		err = fmt.Errorf("Query(tnow): %w", srdblib.Dberr)
		return
	}
	defer rows.Close()

	eventinflist = make([]exsrapi.Event_Inf, 0)
	eventinf := exsrapi.Event_Inf{}

	for rows.Next() {

		srdblib.Dberr = rows.Scan(
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
		)

		if srdblib.Dberr != nil {
			if srdblib.Dberr.Error() != "sql: no rows in result set" {
				return
			} else {
				err = fmt.Errorf("row.Exec(): %w", srdblib.Dberr)
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

		//	log.Printf("eventinf=[%v]\n", eventinf)

		eventinflist = append(eventinflist, eventinf)

		//	log.Printf("Start_data=%f Dperiod=%f\n", eventinf.Start_date, eventinf.Dperiod)
	}

	return
}
