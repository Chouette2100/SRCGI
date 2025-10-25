// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Chouette2100/srdblib/v2"
)

// InsertToDoHandler は新規ToDoの追加を行います
func InsertToDoHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// POSTメソッドのみ受け付ける
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// フォームデータを取得
	itype := r.FormValue("itype")
	target := r.FormValue("target")
	issue := r.FormValue("issue")
	solution := r.FormValue("solution")

	// 必須項目のチェック
	if itype == "" || target == "" || issue == "" {
		http.Error(w, "itype, target, issueは必須項目です", http.StatusBadRequest)
		return
	}

	// 新しいTodoを作成
	todo := srdblib.Todo{
		Ts:       time.Now(),
		Itype:    itype,
		Target:   target,
		Issue:    issue,
		Solution: solution,
		Closed:   nil,
	}

	// データベースに挿入
	err := srdblib.Dbmap.Insert(&todo)
	if err != nil {
		log.Printf("InsertToDoHandler: DB insert error: %s\n", err.Error())
		http.Error(w, "データベースエラーが発生しました", http.StatusInternalServerError)
		return
	}

	log.Printf("InsertToDoHandler: Successfully inserted todo (ID: %d)\n", todo.ID)

	// 検索条件を保持してリダイレクト
	redirectURL := "/list-todo"
	query := r.URL.Query()

	// 元の検索条件があればそれを引き継ぐ
	if searchItype := query.Get("itype"); searchItype != "" {
		redirectURL += "?itype=" + searchItype
	}
	if searchTarget := query.Get("target"); searchTarget != "" {
		if redirectURL == "/list-todo" {
			redirectURL += "?target=" + searchTarget
		} else {
			redirectURL += "&target=" + searchTarget
		}
	}
	if searchIssue := query.Get("issue"); searchIssue != "" {
		if redirectURL == "/list-todo" {
			redirectURL += "?issue=" + searchIssue
		} else {
			redirectURL += "&issue=" + searchIssue
		}
	}
	if searchSolution := query.Get("solution"); searchSolution != "" {
		if redirectURL == "/list-todo" {
			redirectURL += "?solution=" + searchSolution
		} else {
			redirectURL += "&solution=" + searchSolution
		}
	}

	// リスト画面にリダイレクト
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
