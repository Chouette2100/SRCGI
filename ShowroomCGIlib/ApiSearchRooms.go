/*!
Copyright © 2025 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/

package ShowroomCGIlib

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ApiSearchRoomsHandler はルーム名検索APIのハンドラー
// /api/search-rooms に対応
func ApiSearchRoomsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	_, _, isallow := GetUserInf(r)
	if !isallow {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	// クエリパラメータからルーム名を取得
	keyword := r.FormValue("keyword")
	if keyword == "" {
		http.Error(w, "keyword parameter is required", http.StatusBadRequest)
		return
	}

	// limit と offset のパラメータを取得
	limit := 50
	offset := 0

	if limitStr := r.FormValue("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.FormValue("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// ルーム検索を実行
	roomlist, err := SelectUsernoAndName(keyword, limit, offset)
	if err != nil {
		log.Printf("ApiSearchRoomsHandler: SelectUsernoAndName() error: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(roomlist); err != nil {
		log.Printf("ApiSearchRoomsHandler: json.Encode() error: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
