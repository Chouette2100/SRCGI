// Copyright © 2025 chouette2100@gmail.com
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

	"github.com/Chouette2100/srdblib/v2"
	"github.com/dustin/go-humanize"
)

// TodoEditData は編集画面のテンプレートに渡すデータ構造です
type TodoEditData struct {
	Todo       srdblib.Todo
	IsUpdate   bool // 更新処理の場合true
	ErrMsg     string
	SuccessMsg string
	// 検索条件（リダイレクト用）
	Itype    string
	Target   string
	Issue    string
	Solution string
}

// EditToDoHandler は既存ToDoの編集を行います
func EditToDoHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	if r.Method == http.MethodGet {
		// 編集画面の表示
		handleEditToDoGet(w, r)
	} else if r.Method == http.MethodPost {
		// 更新処理
		handleEditToDoPost(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handleEditToDoGet はGETリクエストで編集画面を表示します
func handleEditToDoGet(w http.ResponseWriter, r *http.Request) {
	// IDをクエリパラメータから取得
	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "IDが指定されていません", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "不正なIDです", http.StatusBadRequest)
		return
	}

	// データベースから該当のToDoを取得
	obj, err := srdblib.Dbmap.Get(srdblib.Todo{}, id)
	if err != nil {
		log.Printf("EditToDoHandler: DB get error: %s\n", err.Error())
		http.Error(w, "データベースエラーが発生しました", http.StatusInternalServerError)
		return
	}

	if obj == nil {
		http.Error(w, "該当するToDoが見つかりません", http.StatusNotFound)
		return
	}

	todo := obj.(*srdblib.Todo)

	// 検索条件を保持（クエリパラメータから取得）
	data := TodoEditData{
		Todo:     *todo,
		IsUpdate: false,
		Itype:    r.FormValue("itype_search"),
		Target:   r.FormValue("target_search"),
		Issue:    r.FormValue("issue_search"),
		Solution: r.FormValue("solution_search"),
	}

	renderEditToDo(w, data)
}

// handleEditToDoPost はPOSTリクエストで更新処理を行います
func handleEditToDoPost(w http.ResponseWriter, r *http.Request) {
	// IDを取得
	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "IDが指定されていません", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "不正なIDです", http.StatusBadRequest)
		return
	}

	// データベースから該当のToDoを取得
	obj, err := srdblib.Dbmap.Get(srdblib.Todo{}, id)
	if err != nil {
		log.Printf("EditToDoHandler: DB get error: %s\n", err.Error())
		http.Error(w, "データベースエラーが発生しました", http.StatusInternalServerError)
		return
	}

	if obj == nil {
		http.Error(w, "該当するToDoが見つかりません", http.StatusNotFound)
		return
	}

	todo := obj.(*srdblib.Todo)

	// フォームデータを取得
	itype := r.FormValue("itype")
	target := r.FormValue("target")
	issue := r.FormValue("issue")
	solution := r.FormValue("solution")
	completed := r.FormValue("completed") // チェックボックスの値

	// 必須項目のチェック
	if itype == "" || target == "" || issue == "" {
		data := TodoEditData{
			Todo:     *todo,
			IsUpdate: true,
			ErrMsg:   "種別、対象、課題は必須項目です",
			Itype:    r.FormValue("itype_search"),
			Target:   r.FormValue("target_search"),
			Issue:    r.FormValue("issue_search"),
			Solution: r.FormValue("solution_search"),
		}
		renderEditToDo(w, data)
		return
	}

	// ToDoを更新
	todo.Itype = itype
	todo.Target = target
	todo.Issue = issue
	todo.Solution = solution

	// 完了チェックボックスの処理
	if completed == "on" {
		// チェックが入っている場合、現在時刻をセット
		now := time.Now()
		todo.Closed = &now
	} else {
		// チェックが入っていない場合はNULL
		todo.Closed = nil
	}

	// データベースを更新
	_, err = srdblib.Dbmap.Update(todo)
	if err != nil {
		log.Printf("EditToDoHandler: DB update error: %s\n", err.Error())
		data := TodoEditData{
			Todo:     *todo,
			IsUpdate: true,
			ErrMsg:   "データベースエラーが発生しました",
			Itype:    r.FormValue("itype_search"),
			Target:   r.FormValue("target_search"),
			Issue:    r.FormValue("issue_search"),
			Solution: r.FormValue("solution_search"),
		}
		renderEditToDo(w, data)
		return
	}

	log.Printf("EditToDoHandler: Successfully updated todo (ID: %d)\n", todo.ID)

	// 検索条件を保持してリダイレクト
	redirectURL := "/list-todo"
	searchItype := r.FormValue("itype_search")
	searchTarget := r.FormValue("target_search")
	searchIssue := r.FormValue("issue_search")
	searchSolution := r.FormValue("solution_search")

	// 検索条件をクエリパラメータに追加
	separator := "?"
	if searchItype != "" {
		redirectURL += separator + "itype=" + searchItype
		separator = "&"
	}
	if searchTarget != "" {
		redirectURL += separator + "target=" + searchTarget
		separator = "&"
	}
	if searchIssue != "" {
		redirectURL += separator + "issue=" + searchIssue
		separator = "&"
	}
	if searchSolution != "" {
		redirectURL += separator + "solution=" + searchSolution
	}

	// リスト画面にリダイレクト
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func renderEditToDo(w http.ResponseWriter, data TodoEditData) {
	funcMap := template.FuncMap{
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
		"FormatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
		"FormatTimePtr": func(t *time.Time) string {
			if t == nil {
				return "-"
			}
			return t.Format("2006-01-02 15:04")
		},
	}

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/edit-todo.gtpl"))

	if err := tpl.ExecuteTemplate(w, "edit-todo.gtpl", data); err != nil {
		err = fmt.Errorf("tpl.ExecuteTemplate() error=%s", err.Error())
		w.Write([]byte(err.Error()))
		log.Println(err)
	}
}
