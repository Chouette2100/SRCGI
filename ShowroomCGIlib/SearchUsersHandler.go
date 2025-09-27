// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"errors"
	// "fmt"
	"log"
	// "strconv"
	// "time"
	"encoding/json"

	// "html/template"
	"net/http"

	// "database/sql"

	_ "github.com/go-sql-driver/mysql"

	// "github.com/dustin/go-humanize"

	// "github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// searchUsersHandler はユーザー検索リクエストを処理します
func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーの設定 (開発時のみ、本番ではより厳密に)
	// Flutter Webは通常、異なるポートで動作するため、CORS設定が必要です。
	w.Header().Set("Access-Control-Allow-Origin", "*") // 全てのオリジンを許可
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// OPTIONSリクエストはCORSプリフライトリクエストなので、ここで終了
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータ 'q' (検索文字列) を取得
	query := r.URL.Query().Get("q")
	log.Printf("Received search query: '%s'", query)

	var user []srdblib.User
	var err error

	if query != "" {
		// 検索条件がある場合、部分一致で検索
		// gorp.Select は []interface{} を受け取るので、&user を渡す
		_, err = srdblib.Dbmap.Select(&user, "SELECT userno, user_name FROM user WHERE user_name LIKE ? ORDER BY user_name ASC LIMIT 20", "%"+query+"%")
	} else {
		// 検索条件がない場合、全件取得 (または空リストを返すなど)
		_, err = srdblib.Dbmap.Select(&user, "SELECT userno, user_name FROM user ORDER BY user_name ASC LIMIT 20")
	}

	if err != nil {
		log.Printf("Error querying users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// レスポンスヘッダーを設定し、JSONをエンコードして返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
