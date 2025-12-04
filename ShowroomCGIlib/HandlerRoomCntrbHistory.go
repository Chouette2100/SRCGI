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

// RoomCntrbHistoryData はルーム別貢献ポイント履歴の1レコード
type RoomCntrbHistoryData struct {
	Point     int       // 貢献ポイント
	Lsnid     int       // リスナーID
	Name      string    // リスナー名
	EventID   string    // イベントID
	IeventID  int       // イベントID
	EventName string    `db:"event_name"` // イベント名
	Starttime time.Time // 開始日時
	Endtime   time.Time // 終了日時
}

// RoomCntrbHistoryParam はテンプレートに渡すパラメータ
type RoomCntrbHistoryParam struct {
	Userid     int                    // ルームID
	UserName   string                 // ルーム名
	Nmonths    int                    // 過去何ヶ月間
	Minpoint   int                    // 表示する最小ポイント
	Maxnolines int                    // 最大表示数
	DataList   []RoomCntrbHistoryData // 貢献ポイント履歴リスト
	ErrMsg     string                 // エラーメッセージ
}

// RoomCntrbHistoryHandler はルーム別リスナーの貢献ポイント履歴を表示する
func RoomCntrbHistoryHandler(w http.ResponseWriter, req *http.Request) {

	var param RoomCntrbHistoryParam
	var err error

	// ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// パラメータを取得
	useridStr := req.FormValue("userid")
	if useridStr == "" {
		param.ErrMsg = "ルームIDが指定されていません"
		renderRoomCntrbTemplate(w, param)
		return
	}

	param.Userid, err = strconv.Atoi(useridStr)
	if err != nil {
		param.ErrMsg = fmt.Sprintf("無効なルームID: %s", useridStr)
		renderRoomCntrbTemplate(w, param)
		return
	}

	// ルーム名を取得
	var itrf interface{}
	itrf, err = srdblib.Dbmap.Get(&srdblib.User{}, param.Userid)
	if err != nil {
		param.ErrMsg = fmt.Sprintf("ルーム情報の取得に失敗しました: %v", err)
		log.Printf("Get User error: %v\n", err)
		renderRoomCntrbTemplate(w, param)
		return
	}
	if itrf == nil {
		param.ErrMsg = fmt.Sprintf("ルームID %d が見つかりません", param.Userid)
		renderRoomCntrbTemplate(w, param)
		return
	}
	user := itrf.(*srdblib.User)
	param.UserName = user.User_name

	// nmonths (過去何ヶ月間)
	nmonthsStr := req.FormValue("nmonths")
	if nmonthsStr == "" {
		param.Nmonths = 3 // デフォルト値
	} else {
		param.Nmonths, err = strconv.Atoi(nmonthsStr)
		if err != nil || param.Nmonths < 1 {
			param.Nmonths = 3
			log.Printf("Invalid nmonths: %s, using default: 3\n", nmonthsStr)
		}
	}

	// minpoint (最小ポイント)
	minpointStr := req.FormValue("minpoint")
	if minpointStr == "" {
		param.Minpoint = 5000 // デフォルト値
	} else {
		param.Minpoint, err = strconv.Atoi(minpointStr)
		if err != nil || param.Minpoint < 0 {
			param.Minpoint = 5000
			log.Printf("Invalid minpoint: %s, using default: 5000\n", minpointStr)
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

	// データを取得
	param.DataList, err = selectRoomCntrbHistory(
		param.Userid,
		param.Nmonths,
		param.Minpoint,
		param.Maxnolines,
	)
	if err != nil {
		param.ErrMsg = fmt.Sprintf("データ取得エラー: %v", err)
		log.Printf("selectRoomCntrbHistory error: %v\n", err)
		renderRoomCntrbTemplate(w, param)
		return
	}

	renderRoomCntrbTemplate(w, param)
}

// selectRoomCntrbHistory は指定条件でルーム別貢献ポイント履歴を取得する
func selectRoomCntrbHistory(
	userid int,
	nmonths int,
	minpoint int,
	maxnolines int,
) (
	dataList []RoomCntrbHistoryData,
	err error,
) {

	// 現在日の00時00分のnmonths月前の日時
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate := startOfToday.AddDate(0, -nmonths, 0)
	endDate := startOfToday

	sqlst := `
SELECT er.point,
        er.lsnid, v.name,
        e.eventid, e.ieventid, e.event_name,
        e.starttime, e.endtime
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
       AND er.userid = ?
     ORDER BY er.point DESC
     LIMIT ?
`

	_, err = srdblib.Dbmap.Select(
		&dataList,
		sqlst,
		minpoint,
		startDate,
		endDate,
		userid,
		maxnolines,
	)

	if err != nil {
		err = fmt.Errorf("Dbmap.Select() error: %w", err)
		log.Printf("%s\n", err.Error())
		return nil, err
	}

	return dataList, nil
}

// renderRoomCntrbTemplate はテンプレートを実行してレスポンスを返す
func renderRoomCntrbTemplate(w http.ResponseWriter, param RoomCntrbHistoryParam) {
	funcMap := sprig.FuncMap() // https://masterminds.github.io/sprig/

	funcMap["Comma"] = func(i int) string { return humanize.Comma(int64(i)) }
	funcMap["FormatTime"] = func(t time.Time, layout string) string {
		return t.Format(layout)
	}
	funcMap["iscurrent"] = func(t time.Time) bool { return t.After(time.Now()) }

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/room-cntrb-history.gtpl"))

	if err := tpl.ExecuteTemplate(w, "room-cntrb-history.gtpl", param); err != nil {
		log.Printf("Template execution error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
