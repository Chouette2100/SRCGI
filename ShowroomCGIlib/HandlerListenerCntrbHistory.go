// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/srdblib/v2"
)

// ListenerCntrbHistoryData は貢献ポイント履歴の1レコード
type ListenerCntrbHistoryData struct {
	Point     int       // 貢献ポイント
	Lsnid     int       // リスナーID
	Name      string    // リスナー名
	Userid    int       // リスナーID
	UserName  string    `db:"user_name"` // ルーム名
	Eventid   string    // イベントID
	Ieventid  int       // イベントID
	EventName string    `db:"event_name"` // イベント名
	Starttime time.Time // 開始日時
	Endtime   time.Time // 終了日時
}

// ListenerCntrbHistoryParam はテンプレートに渡すパラメータ
type ListenerCntrbHistoryParam struct {
	EventID    string                     // イベントID
	IeventID   int                        // イベントID
	EventName  string                     // イベント名
	Nmonths    int                        // 過去何ヶ月間
	Minpoint   int                        // 表示する最小ポイント
	Maxnolines int                        // 最大表示数
	Ext        int                        // 0: 参加ルームのみ, 1: 全ルーム
	DataList   []ListenerCntrbHistoryData // 貢献ポイント履歴リスト
	ErrMsg     string                     // エラーメッセージ
}

// ListenerCntrbHistoryHandler はイベント参加ルームのリスナーの貢献ポイント履歴を表示する
func ListenerCntrbHistoryHandler(w http.ResponseWriter, req *http.Request) {

	var param ListenerCntrbHistoryParam
	var err error

	// ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// パラメータを取得
	param.EventID = req.FormValue("eventid")
	if param.EventID == "" {
		param.ErrMsg = "イベントIDが指定されていません"
		renderTemplate(w, param)
		return
	}

	// イベント名を取得
	var itrf1 interface{}
	itrf1, err = srdblib.Dbmap.Get(&srdblib.Event{}, param.EventID)
	if err != nil {
		param.ErrMsg = fmt.Sprintf("イベント情報の取得に失敗しました: %v", err)
		log.Printf("Get Event error: %v\n", err)
		renderTemplate(w, param)
		return
	}
	if itrf1 == nil {
		param.ErrMsg = fmt.Sprintf("イベントID %s が見つかりません", param.EventID)
		renderTemplate(w, param)
		return
	}
	event := itrf1.(*srdblib.Event)
	param.EventName = event.Event_name
	param.IeventID = event.Ieventid

	// nmonths (過去何ヶ月間)
	nmonthsStr := req.FormValue("nmonths")
	if nmonthsStr == "" {
		param.Nmonths = 2 // デフォルト値
	} else {
		param.Nmonths, err = strconv.Atoi(nmonthsStr)
		if err != nil || param.Nmonths < 1 {
			param.Nmonths = 2
			log.Printf("Invalid nmonths: %s, using default: 2\n", nmonthsStr)
		}
	}

	// minpoint (最小ポイント)
	minpointStr := req.FormValue("minpoint")
	if minpointStr == "" {
		param.Minpoint = 200000 // デフォルト値
	} else {
		param.Minpoint, err = strconv.Atoi(minpointStr)
		if err != nil || param.Minpoint < 0 {
			param.Minpoint = 200000
			log.Printf("Invalid minpoint: %s, using default: 200000\n", minpointStr)
		}
	}

	// maxnolines (最大表示数)
	maxnolinesStr := req.FormValue("maxnolines")
	if maxnolinesStr == "" {
		param.Maxnolines = 100 // デフォルト値
	} else {
		param.Maxnolines, err = strconv.Atoi(maxnolinesStr)
		if err != nil || param.Maxnolines < 1 {
			param.Maxnolines = 100
			log.Printf("Invalid maxnolines: %s, using default: 100\n", maxnolinesStr)
		}
	}

	// ext (フィルター)
	// ext=1: 参加ルームのみ表示, ext=0: 全ルーム表示
	extStr := req.FormValue("ext")
	if extStr == "" {
		param.Ext = 1 // デフォルト: 参加ルームのみ
	} else {
		param.Ext, err = strconv.Atoi(extStr)
		if err != nil {
			param.Ext = 1
		}
	}

	// データを取得
	param.DataList, err = selectListenerCntrbHistory(
		param.EventID,
		param.Nmonths,
		param.Minpoint,
		param.Maxnolines,
		param.Ext,
	)
	if err != nil {
		param.ErrMsg = fmt.Sprintf("データ取得エラー: %v", err)
		log.Printf("selectListenerCntrbHistory error: %v\n", err)
		renderTemplate(w, param)
		return
	}

	renderTemplate(w, param)
}

// selectListenerCntrbHistory は指定条件で貢献ポイント履歴を取得する
func selectListenerCntrbHistory(
	eventid string,
	nmonths int,
	minpoint int,
	maxnolines int,
	ext int,
) (
	dataList []ListenerCntrbHistoryData,
	err error,
) {

	// 現在日の00時00分のnmonths月前の日時
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate := startOfToday.AddDate(0, -nmonths, 0)
	endDate := startOfToday

	var sqlst string
	if ext == 1 {
		// このイベントに参加するルームに対する貢献のみ表示
		sqlst = `
SELECT er.point, er.lsnid, v.name,
	   er.userid, u.user_name,
	   er.eventid, e.ieventid, e.event_name, e.starttime, e.endtime
	FROM eventrank er
    JOIN event e ON e.eventid = er.eventid
    JOIN user u ON u.userno = er.userid
    JOIN viewer v ON v.viewerid = er.lsnid
	  WHERE er.point > ?                                      -- minpoint
    AND e.starttime BETWEEN ? AND ?      -- tb, te
       AND e.eventid not like '%block_id=0'
       AND er.ts = (
             SELECT MAX(ts)FROM eventrank
              WHERE eventid = e.eventid AND userid = u.userno
           )
       AND er.userid IN (
             SELECT userno AS userid FROM eventuser
              WHERE eventid = ?    -- eventid
           )
     ORDER BY er.point DESC
     LIMIT ?                                                 -- LIMIT
	`

		_, err = srdblib.Dbmap.Select(
			&dataList,
			sqlst,
			minpoint,
			startDate,
			endDate,
			eventid,
			maxnolines,
		)
	} else {
		// このイベントに参加していないルームに対する貢献も表示
		sqlst = `
SELECT er.point, er.lsnid, v.name,
	   er.userid, u.user_name,
	   er.eventid, e.ieventid, e.event_name, e.starttime, e.endtime
	FROM eventrank er
  JOIN event e ON e.eventid = er.eventid
  JOIN user u ON u.userno = er.userid
  JOIN viewer v ON v.viewerid = er.lsnid
  WHERE er.point > ?
    AND e.starttime BETWEEN ? AND ?
       AND e.eventid not like '%block_id=0'
       AND er.ts = (
             SELECT MAX(ts)FROM eventrank
              WHERE eventid = e.eventid AND userid = u.userno
           )
       AND er.userid IN (
             SELECT userno AS userid FROM eventuser
              WHERE eventid = ?
           )
     ORDER BY er.point DESC
     LIMIT ?
`

		_, err = srdblib.Dbmap.Select(
			&dataList,
			sqlst,
			minpoint,
			startDate,
			endDate,
			eventid,
			maxnolines,
		)
	}

	if err != nil {
		err = fmt.Errorf("Dbmap.Select() error: %w", err)
		log.Printf("%s\n", err.Error())
		return nil, err
	}

	return dataList, nil
}

// renderTemplate はテンプレートを実行してレスポンスを返す
func renderTemplate(w http.ResponseWriter, param ListenerCntrbHistoryParam) {
	funcMap := sprig.FuncMap() // https://masterminds.github.io/sprig/

	funcMap["Comma"] = func(i int) string { return humanize.Comma(int64(i)) }
	funcMap["FormatTime"] = func(t time.Time, layout string) string {
		return t.Format(layout)
	}
	// funcMap["now"] = time.Now
	funcMap["iscurrent"] = func(t time.Time) bool { return t.After(time.Now()) }

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/listener-cntrb-history.gtpl"))

	if err := tpl.ExecuteTemplate(w, "listener-cntrb-history.gtpl", param); err != nil {
		log.Printf("Template execution error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
