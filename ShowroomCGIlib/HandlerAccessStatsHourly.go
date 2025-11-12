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

// AccessStatsHourlyHandler は時刻単位のアクセス統計を表示するハンドラー
func AccessStatsHourlyHandler(w http.ResponseWriter, r *http.Request) {

	var accessStatsHourlyData AccessStatsHourlyData
	var stats []AccessStatsHourly

	// 現在時刻を設定
	accessStatsHourlyData.TimeNow = time.Now()

	// パラメータから開始日時、終了日時を取得
	startDateTime := r.FormValue("start_datetime")
	endDateTime := r.FormValue("end_datetime")

	// デフォルト値の設定（直近72時間）
	now := time.Now()
	if startDateTime == "" {
		// 現在時刻を次の時刻に切り上げ（例: 16:42 → 17:00）
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		// 72時間前の時刻
		startTime := endTime.Add(-72 * time.Hour)
		startDateTime = startTime.Format("2006-01-02T15:04")
	}
	if endDateTime == "" {
		// 現在時刻を次の時刻に切り上げ（例: 16:42 → 17:00）
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		endDateTime = endTime.Format("2006-01-02T15:04")
	}

	accessStatsHourlyData.StartDateTime = startDateTime
	accessStatsHourlyData.EndDateTime = endDateTime

	// 日時文字列をパース
	startTime, err := time.Parse("2006-01-02T15:04", startDateTime)
	if err != nil {
		log.Printf("AccessStatsHourlyHandler() startDateTime parse error: %v", err)
		http.Error(w, "Invalid start datetime format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04", endDateTime)
	if err != nil {
		log.Printf("AccessStatsHourlyHandler() endDateTime parse error: %v", err)
		http.Error(w, "Invalid end datetime format", http.StatusBadRequest)
		return
	}

	// 終了時刻の1時間後を取得（WHEREの範囲指定のため）
	nextHour := endTime.Add(1 * time.Hour)

	// SQLクエリを実行して時刻単位のアクセス統計を取得（3つのカテゴリ）
	sql := `
		SELECT
			DATE_FORMAT(ts, '%Y-%m-%d %H:00:00') AS access_hour,
			COUNT(*) AS hourly_access_count,
			SUM(CASE WHEN is_bot = 0 AND turnstilestatus = 0 THEN 1 ELSE 0 END) AS legitimate_count,
			SUM(CASE WHEN is_bot = 0 AND turnstilestatus != 0 THEN 1 ELSE 0 END) AS turnstile_fail_count,
			SUM(CASE WHEN is_bot != 0 THEN 1 ELSE 0 END) AS bot_count
		FROM
			accesslog
		WHERE ts >= ? AND ts < ?
		  AND remoteaddress != '59.166.119.117'
		  AND remoteaddress != '10.63.22.1'
		  AND remoteaddress != '149.88.103.40'
		GROUP BY
			access_hour
		ORDER BY
			access_hour;
	`

	// データベースからアクセス統計を取得
	_, err = srdblib.Dbmap.Select(&stats, sql, startTime.Format("2006-01-02 15:04:05"), nextHour.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("AccessStatsHourlyHandler() database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	accessStatsHourlyData.Stats = stats

	// テンプレートを実行
	tpl := template.Must(template.ParseFiles("templates/accessstatshourly.gtpl"))
	err = tpl.Execute(w, accessStatsHourlyData)
	if err != nil {
		log.Printf("AccessStatsHourlyHandler() template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}
