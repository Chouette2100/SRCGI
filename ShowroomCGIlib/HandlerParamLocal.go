package ShowroomCGIlib

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HandlerParamLocal(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/param-local.gtpl"))
	values := map[string]string{
		"function": "イベントパラメータの設定",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "param-local.gtpl", values); err != nil {
		log.Println(err)
	}

}

func ParamGlobalHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/param-global.gtpl"))
	values := map[string]string{
		"function": "共通パラメータの設定",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "param-global.gtpl", values); err != nil {
		log.Println(err)
	}

}

func CsvTotalHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/csv-total.gtpl"))
	values := map[string]string{
		"function": "獲得ポイントの推移（CSV）",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "csv-total.gtpl", values); err != nil {
		log.Println(err)
	}

}

func GraphDfrHandler(w http.ResponseWriter, r *http.Request) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-dfr.gtpl"))
	values := map[string]string{
		"function": "獲得ポイントの差の推移（グラフ）",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "graph-dfr.gtpl", values); err != nil {
		log.Println(err)
	}

}
