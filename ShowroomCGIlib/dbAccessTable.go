// Copyright © 2024-2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/Chouette2100/srdblib/v2"
)

// GetAccessByHandler はハンドラーごとの日別アクセス数を取得する
func GetAccessByHandler(endDate time.Time, days int) ([]AccessTableRow, []string, error) {
	startDate := endDate.AddDate(0, 0, -(days - 1))

	sql := `
		SELECT 
			handler,
			DATE(ts) AS access_date,
			COUNT(*) AS access_count
		FROM accesslog
		WHERE DATE(ts) BETWEEN ? AND ?
		  AND is_bot = 0
		  AND remoteaddress != '59.166.119.117'
		  AND remoteaddress != '10.63.22.1'
		  AND remoteaddress != '149.88.103.40'
		GROUP BY handler, access_date
		ORDER BY handler, access_date
	`

	type AccessRow struct {
		Handler     string `db:"handler"`
		AccessDate  string `db:"access_date"`
		AccessCount int    `db:"access_count"`
	}

	var rows []AccessRow
	_, err := srdblib.Dbmap.Select(&rows, sql,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	if err != nil {
		return nil, nil, fmt.Errorf("GetAccessByHandler: %w", err)
	}

	tableRows, headers := buildTableRows(rows, startDate, days, func(row AccessRow) string { return row.Handler })
	return tableRows, headers, nil
}

// GetAccessByIP はIPアドレスごとの日別アクセス数を取得する（上位50件のみ）
func GetAccessByIP(endDate time.Time, days int, dkey string) ([]AccessTableRow, []string, error) {
	startDate := endDate.AddDate(0, 0, -(days - 1))

	// まず、期間全体でのアクセス数上位50件のIPアドレスを取得
	sqlTop := `
		SELECT remoteaddress, COUNT(*) AS total_count
		FROM accesslog
		WHERE DATE(ts) BETWEEN ? AND ?
		  AND is_bot = 0
		  AND remoteaddress != '59.166.119.117'
		  AND remoteaddress != '10.63.22.1'
		  AND remoteaddress != '149.88.103.40'
		GROUP BY remoteaddress
		ORDER BY total_count DESC
		LIMIT 50
	`

	type TopIP struct {
		Remoteaddress string `db:"remoteaddress"`
		TotalCount    int    `db:"total_count"`
	}

	var topIPs []TopIP
	_, err := srdblib.Dbmap.Select(&topIPs, sqlTop,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	if err != nil {
		return nil, nil, fmt.Errorf("GetAccessByIP (top): %w", err)
	}

	if len(topIPs) == 0 {
		return []AccessTableRow{}, generateDateHeaders(startDate, days), nil
	}

	// 上位50件のIPアドレスの日別アクセス数を取得
	// ipList := make([]interface{}, len(topIPs))
	ipList := make([]string, len(topIPs))
	for i, ip := range topIPs {
		ipList[i] = ip.Remoteaddress
	}

	sql := `
		SELECT 
			remoteaddress,
			DATE(ts) AS access_date,
			COUNT(*) AS access_count
		FROM accesslog
		-- WHERE DATE(ts) BETWEEN ? AND ?
		WHERE DATE(ts) BETWEEN :Sdate AND :Edate
		  AND is_bot = 0
		--  AND remoteaddress IN ( ? )
		  AND remoteaddress IN ( :Ips )
		GROUP BY remoteaddress, access_date
		ORDER BY remoteaddress, access_date
	`

	type AccessRow struct {
		Remoteaddress string `db:"remoteaddress"`
		AccessDate    string `db:"access_date"`
		AccessCount   int    `db:"access_count"`
	}

	var rows []AccessRow
	_, err = srdblib.Dbmap.Select(&rows, sql,
		// startDate.Format("2006-01-02"),
		// endDate.Format("2006-01-02"),
		//	ipList)
		// map[string]interface{}{"Ips": ipList})
		map[string]interface{}{"Ips": ipList,
			"Sdate": startDate.Format("2006-01-02"), "Edate": endDate.Format("2006-01-02")})
	if err != nil {
		return nil, nil, fmt.Errorf("GetAccessByIP: %w", err)
	}

	// IPアドレスを暗号化
	tableRows, headers := buildTableRows(rows, startDate, days, func(row AccessRow) string { return row.Remoteaddress })

	// encryptionKey []byteを16進数表現にして、16進表現のdkey stringと比較する
	encryptionKeyHex := fmt.Sprintf("%X", encryptionKey)
	if dkey != encryptionKeyHex {

		for i := range tableRows {
			encrypted, err := EncryptIP(tableRows[i].Key)
			if err != nil {
				log.Printf("Warning: failed to encrypt IP %s: %v", tableRows[i].Key, err)
				encrypted = "[暗号化エラー]"
			}
			tableRows[i].Key = encrypted
		}
	}

	return tableRows, headers, nil
}

// GetAccessByUserAgent はユーザーエージェントごとの日別アクセス数を取得する（上位50件のみ）
func GetAccessByUserAgent(endDate time.Time, days int) ([]AccessTableRow, []string, error) {
	startDate := endDate.AddDate(0, 0, -(days - 1))

	// まず、期間全体でのアクセス数上位50件のユーザーエージェントを取得
	sqlTop := `
		SELECT useragent, COUNT(*) AS total_count
		FROM accesslog
		WHERE DATE(ts) BETWEEN ? AND ?
		  AND is_bot = 0
		  AND remoteaddress != '59.166.119.117'
		  AND remoteaddress != '10.63.22.1'
		  AND remoteaddress != '149.88.103.40'
		GROUP BY useragent
		ORDER BY total_count DESC
		LIMIT 50
	`

	type TopUA struct {
		Useragent  string `db:"useragent"`
		TotalCount int    `db:"total_count"`
	}

	var topUAs []TopUA
	_, err := srdblib.Dbmap.Select(&topUAs, sqlTop,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))
	if err != nil {
		return nil, nil, fmt.Errorf("GetAccessByUserAgent (top): %w", err)
	}

	if len(topUAs) == 0 {
		return []AccessTableRow{}, generateDateHeaders(startDate, days), nil
	}

	// 上位50件のユーザーエージェントの日別アクセス数を取得
	// uaList := make([]interface{}, len(topUAs))
	uaList := make([]string, len(topUAs))
	for i, ua := range topUAs {
		uaList[i] = ua.Useragent
	}

	sql := `
		SELECT 
			useragent,
			DATE(ts) AS access_date,
			COUNT(*) AS access_count
		FROM accesslog
		-- WHERE DATE(ts) BETWEEN ? AND ?
		WHERE DATE(ts) BETWEEN :Sdate AND :Edate
		  AND is_bot = 0
		--   AND useragent IN ( ? )
		  AND useragent IN ( :Ual )
		GROUP BY useragent, access_date
		ORDER BY useragent, access_date
	`

	type AccessRow struct {
		Useragent   string `db:"useragent"`
		AccessDate  string `db:"access_date"`
		AccessCount int    `db:"access_count"`
	}

	var rows []AccessRow
	_, err = srdblib.Dbmap.Select(&rows, sql,
		// startDate.Format("2006-01-02"),
		// endDate.Format("2006-01-02"),
		// uaList)
		map[string]interface{}{"Ual": uaList,
			"Sdate": startDate.Format("2006-01-02"), "Edate": endDate.Format("2006-01-02")})
	if err != nil {
		return nil, nil, fmt.Errorf("GetAccessByUserAgent: %w", err)
	}

	tableRows, headers := buildTableRows(rows, startDate, days, func(row AccessRow) string { return row.Useragent })
	return tableRows, headers, nil
}

// buildTableRows はクエリ結果からAccessTableRowのスライスを構築する
func buildTableRows[T any](rows []T, startDate time.Time, days int, keyFunc func(T) string) ([]AccessTableRow, []string) {
	dateHeaders := generateDateHeaders(startDate, days)

	// キーごとにデータを集約
	dataMap := make(map[string]map[string]int) // map[key]map[date]count

	for _, row := range rows {
		key := keyFunc(row)

		var date string
		var count int

		// リフレクションを使用してフィールドを取得
		v := reflect.ValueOf(row)

		// AccessDateフィールドを取得
		accessDateField := v.FieldByName("AccessDate")
		if accessDateField.IsValid() && accessDateField.Kind() == reflect.String {
			dateStr := accessDateField.String()
			// タイムスタンプ形式から日付部分のみを抽出（YYYY-MM-DD形式）
			if len(dateStr) >= 10 {
				date = dateStr[:10]
			} else {
				date = dateStr
			}
		}

		// AccessCountフィールドを取得
		accessCountField := v.FieldByName("AccessCount")
		if accessCountField.IsValid() && accessCountField.Kind() == reflect.Int {
			count = int(accessCountField.Int())
		}

		if _, exists := dataMap[key]; !exists {
			dataMap[key] = make(map[string]int)
		}
		dataMap[key][date] = count
	}

	// AccessTableRowのスライスを構築
	var tableRows []AccessTableRow
	for key, dateCounts := range dataMap {
		row := AccessTableRow{
			Key:         key,
			DailyCounts: make([]DailyAccessCount, days),
			Total:       0,
		}

		for i := 0; i < days; i++ {
			date := startDate.AddDate(0, 0, i).Format("2006-01-02")
			count := dateCounts[date]
			row.DailyCounts[i] = DailyAccessCount{
				Date:  date,
				Count: count,
			}
			row.Total += count
		}

		tableRows = append(tableRows, row)
	}

	// 合計の降順でソート
	for i := 0; i < len(tableRows)-1; i++ {
		for j := i + 1; j < len(tableRows); j++ {
			if tableRows[i].Total < tableRows[j].Total {
				tableRows[i], tableRows[j] = tableRows[j], tableRows[i]
			}
		}
	}

	return tableRows, dateHeaders
}

// generateDateHeaders は日付ヘッダーのスライスを生成する
func generateDateHeaders(startDate time.Time, days int) []string {
	headers := make([]string, days)
	for i := 0; i < days; i++ {
		date := startDate.AddDate(0, 0, i)
		headers[i] = date.Format("06-01-02") // YY-MM-DD形式
	}
	return headers
}
