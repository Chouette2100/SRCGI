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

// TodoListData はテンプレートに渡すデータ構造です
type TodoListData struct {
	Todos    []srdblib.Todo
	Itype    string
	Target   string
	Issue    string
	Solution string
	MaxID    int
	MinID    int
	Limit    int
	HasNext  bool
	HasPrev  bool
	ErrMsg   string
}

// ListToDoHandler はToDoリストの表示と検索を行います
func ListToDoHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// クエリパラメータから検索条件を取得
	itype := r.FormValue("itype")
	target := r.FormValue("target")
	issue := r.FormValue("issue")
	solution := r.FormValue("solution")

	// ページング用のパラメータ
	maxIDStr := r.FormValue("maxid")
	minIDStr := r.FormValue("minid")
	direction := r.FormValue("dir") // "next", "prev", または空文字（初期表示・再表示）

	limit := 50
	var maxID, minID int

	if maxIDStr != "" {
		maxID, _ = strconv.Atoi(maxIDStr)
	}
	if minIDStr != "" {
		minID, _ = strconv.Atoi(minIDStr)
	}

	// SQL構築
	whereConditions := []string{}
	args := []interface{}{}

	// 検索条件の追加
	if itype != "" {
		whereConditions = append(whereConditions, "itype = ?")
		args = append(args, itype)
	}
	if target != "" {
		whereConditions = append(whereConditions, "target = ?")
		args = append(args, target)
	}
	if issue != "" || solution != "" {
		orConditions := []string{}
		if issue != "" {
			orConditions = append(orConditions, "issue LIKE ?")
			args = append(args, "%"+issue+"%")
		}
		if solution != "" {
			orConditions = append(orConditions, "solution LIKE ?")
			args = append(args, "%"+solution+"%")
		}
		if len(orConditions) > 0 {
			whereConditions = append(whereConditions, "("+orConditions[0]+")")
			if len(orConditions) > 1 {
				whereConditions[len(whereConditions)-1] = "(" + orConditions[0] + " OR " + orConditions[1] + ")"
			}
		}
	}

	// ページング条件の追加
	if direction == "next" && minID > 0 {
		// 次ページ: 現在のminIDより小さいIDを取得
		whereConditions = append(whereConditions, "id < ?")
		args = append(args, minID)
	} else if direction == "prev" && maxID > 0 {
		// 前ページ: 現在のmaxIDより大きいIDを取得
		whereConditions = append(whereConditions, "id > ?")
		args = append(args, maxID)
	}

	// WHERE句の構築
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE "
		for i, cond := range whereConditions {
			if i > 0 {
				whereClause += " AND "
			}
			whereClause += cond
		}
	}

	// ORDER BY句（前ページの場合は逆順でソート後、結果を反転）
	orderClause := " ORDER BY id DESC"
	if direction == "prev" {
		orderClause = " ORDER BY id ASC"
	}

	// LIMIT句
	limitClause := fmt.Sprintf(" LIMIT %d", limit+1) // +1は次ページの有無を判定するため

	query := "SELECT id, ts, itype, target, issue, solution, closed FROM todo" + whereClause + orderClause + limitClause

	log.Printf("ListToDoHandler: query=%s, args=%v\n", query, args)

	// データ取得
	rows, err := srdblib.Dbmap.Select(srdblib.Todo{}, query, args...)
	if err != nil {
		log.Printf("ListToDoHandler: DB query error: %s\n", err.Error())
		data := TodoListData{
			ErrMsg:   "データベースエラーが発生しました",
			Itype:    itype,
			Target:   target,
			Issue:    issue,
			Solution: solution,
		}
		renderTodoList(w, data)
		return
	}

	// 結果をTodo型に変換
	todos := make([]srdblib.Todo, 0, len(rows))
	for _, row := range rows {
		todos = append(todos, *row.(*srdblib.Todo))
	}

	// 次ページの有無を判定
	hasNext := false
	hasPrev := false
	if len(todos) > limit {
		hasNext = true
		todos = todos[:limit]
	}

	// 前ページの場合は結果を反転
	switch direction {
	case "prev":
		for i, j := 0, len(todos)-1; i < j; i, j = i+1, j-1 {
			todos[i], todos[j] = todos[j], todos[i]
		}
		hasPrev = maxID > 0 // 実際には追加のクエリで判定すべき
	case "next":
		hasPrev = true
	}

	// MaxID と MinID を計算
	newMaxID := 0
	newMinID := 0
	if len(todos) > 0 {
		newMaxID = todos[0].ID
		newMinID = todos[len(todos)-1].ID
	}

	// 初期表示または再表示の場合、前ページがあるかチェック
	if direction == "" {
		if len(todos) > 0 {
			checkQuery := "SELECT COUNT(*) FROM todo" + whereClause + " AND id > ?"
			checkArgs := append(args, newMaxID)
			var count int
			err := srdblib.Dbmap.Db.QueryRow(checkQuery, checkArgs...).Scan(&count)
			if err == nil && count > 0 {
				hasPrev = true
			}
		}
	}

	data := TodoListData{
		Todos:    todos,
		Itype:    itype,
		Target:   target,
		Issue:    issue,
		Solution: solution,
		MaxID:    newMaxID,
		MinID:    newMinID,
		Limit:    limit,
		HasNext:  hasNext,
		HasPrev:  hasPrev,
	}

	renderTodoList(w, data)
}

func renderTodoList(w http.ResponseWriter, data TodoListData) {
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

	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-todo.gtpl"))

	if err := tpl.ExecuteTemplate(w, "list-todo.gtpl", data); err != nil {
		err = fmt.Errorf("tpl.ExecuteTemplate() error=%s", err.Error())
		w.Write([]byte(err.Error()))
		log.Println(err)
	}
}
