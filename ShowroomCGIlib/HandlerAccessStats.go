// Copyright © 2024-2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Chouette2100/srdblib/v2"
)

// AccessStatsHandler はアクセス統計を表示するハンドラー
func AccessStatsHandler(w http.ResponseWriter, r *http.Request) {

	var accessStatsData AccessStatsData
	var stats []AccessStats

	// 現在時刻を設定
	accessStatsData.TimeNow = time.Now()

	// パラメータから開始日、終了日を取得
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")

	// デフォルト値の設定（直近1ヶ月）
	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	accessStatsData.StartDate = startDate
	accessStatsData.EndDate = endDate

	// SQLクエリを実行してアクセス統計を取得
	sql := `
		SELECT
			DATE(ts) AS access_date,
			COUNT(*) AS daily_access_count
		FROM
			accesslog
		WHERE ts between ? AND ?
		  AND is_bot = 0
		  AND remoteaddress != '59.166.119.117'
		  AND remoteaddress != '10.63.22.1'
		  AND remoteaddress != '149.88.103.40'
		GROUP BY
			access_date
		ORDER BY
			access_date;
	`

	// 終了日の次の日を取得（WHEREの範囲指定のため）
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Printf("AccessStatsHandler() endDate parse error: %v", err)
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}
	nextDay := endDateTime.AddDate(0, 0, 1).Format("2006-01-02")

	// データベースからアクセス統計を取得
	_, err = srdblib.Dbmap.Select(&stats, sql, startDate, nextDay)
	if err != nil {
		log.Printf("AccessStatsHandler() database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	accessStatsData.Stats = stats

	// テンプレートを実行
	tpl := template.Must(template.ParseFiles("templates/accessstats.gtpl"))
	err = tpl.Execute(w, accessStatsData)
	if err != nil {
		log.Printf("AccessStatsHandler() template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}
