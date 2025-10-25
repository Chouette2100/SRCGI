package ShowroomCGIlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/srapi/v2"
)

// APIから取得したminifyされたJSONを整形する関数
func formatJSON(rawJSON []byte) (*bytes.Buffer, error) {
	// 1. まずJSONをパースしてデータ構造を確認
	var data interface{}
	if err := json.Unmarshal(rawJSON, &data); err != nil {
		return nil, err
	}

	// 2. 整形してBufferに書き込み
	var prettyJSON bytes.Buffer
	encoder := json.NewEncoder(&prettyJSON)
	encoder.SetIndent("", "  ")  // 2スペースでインデント
	encoder.SetEscapeHTML(false) // HTMLエスケープを無効化

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return &prettyJSON, nil
}

func DisplayJSON20Handler(w http.ResponseWriter, r *http.Request) {
	// jsonData := []byte(`{"name":"John","items":[1,2,3],"active":true}`)

	room_id := r.URL.Query().Get("room_id")
	client := http.Client{}
	jsonData, err := srapi.JsonRoomEventAndSupport(&client, room_id)
	if err != nil {
		err = fmt.Errorf("srapi.JsonRoomEventAndSupport() error=%s", err.Error())
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	prettyJSON, err := formatJSON(jsonData.Bytes())
	if err != nil {
		http.Error(w, "JSONの整形に失敗しました: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonStr := prettyJSON.String()

	// var prettyJSON bytes.Buffer
	// json.Indent(&prettyJSON, []byte(*jsonData), "", "  ")

	/*
		tmpl := template.Must(template.New("json").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>JSONビューア</title>
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/github-dark.min.css">
			<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js"></script>
			<script>hljs.highlightAll();</script>
			<style>
				pre { border-radius: 8px; padding: 20px; }
			</style>
		</head>
		<body>
			<h1>JSONデータ</h1>
			<pre><code class="language-json">{{.JSONData}}</code></pre>
		</body>
		</html>
		`))
	*/

	data := struct {
		JSONData string
	}{
		JSONData: jsonStr,
	}

	/*
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	*/

	funcMap := sprig.FuncMap() // https://masterminds.github.io/sprig/
	funcMap["Comma"] = func(i int) string { return humanize.Comma(int64(i)) }
	funcMap["baseOfEventid"] = func(s string) string { ida := strings.Split(s, "?"); return ida[0] }
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/json20withC.gtpl"))

	if err := tpl.ExecuteTemplate(w, "json20withC.gtpl", data); err != nil {
		err = fmt.Errorf("tpl.ExecuteTemplate() error=%s", err.Error())
		w.Write([]byte(err.Error()))
		log.Println(err)
	}

}
