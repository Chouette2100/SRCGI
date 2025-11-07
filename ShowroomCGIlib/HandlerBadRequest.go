// Copyright © 2024-2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"log"

	"net/http"
)

// 望ましくないアクセスに対するハンドラー
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {

	ruk := r.FormValue("room_url_key")

	log.Printf("BadRequestHandler(): Bad Request from %s (%s)\n", r.RemoteAddr, ruk)
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
