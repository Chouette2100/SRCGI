// Copyright © 2024-2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// AccessTableHandler はハンドラー、IPアドレス、ユーザーエージェント別のアクセス数を集計表示するハンドラー
func AccessTableHandler(w http.ResponseWriter, r *http.Request) {
	var data AccessTableData

	// 現在時刻を設定
	data.TimeNow = time.Now()

	// クエリパラメータから集計タイプを取得（デフォルト: handler）
	typeParam := r.FormValue("type")
	if typeParam == "" {
		typeParam = "handler"
	}
	data.Type = typeParam

	dkey := r.FormValue("dkey")
	data.Dkey = dkey

	eIp := r.FormValue("eip")
	data.EIp = eIp

	// 終了日を取得（デフォルト: 今日）
	endDateStr := r.FormValue("end_date")
	var endDate time.Time
	var err error

	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			log.Printf("AccessTableHandler() endDate parse error: %v", err)
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
	}
	data.EndDate = endDate.Format("2006-01-02")

	// 固定で7日間のデータを取得
	days := 7

	// タイプに応じてデータを取得
	switch typeParam {
	case "handler":
		data.Rows, data.DateHeaders, err = GetAccessByHandler(endDate, days)
	case "ip":
		data.Rows, data.DateHeaders, err = GetAccessByIP(endDate, days, dkey)
	case "useragent":
		data.Rows, data.DateHeaders, err = GetAccessByUserAgent(endDate, days)
	default:
		log.Printf("AccessTableHandler() invalid type: %s", typeParam)
		http.Error(w, "Invalid type parameter", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("AccessTableHandler() data fetch error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// テンプレートを実行
	tpl := template.Must(template.ParseFiles("templates/accesstable.gtpl"))
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("AccessTableHandler() template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}
