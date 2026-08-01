package ShowroomCGIlib

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

// AnnouncementHandler は告知入力フォーム表示・投稿を処理するハンドラー
func AnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// フォーム表示
		// displayAnnouncementForm(w, r)
		displayAnnouncementForm(w)
	case http.MethodPost:
		// 告知内容の受け取り・保存
		saveAnnouncement(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// displayAnnouncementForm は告知入力フォームを表示する
func displayAnnouncementForm(w http.ResponseWriter) {
// func displayAnnouncementForm(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("templates/announcement-form.gtpl"))

	data := struct {
		Enabled          bool
		CurrentMessage   string
		CurrentTextColor string
		CurrentBgColor   string
	}{
		Enabled:          CurrentAnnouncement.Enabled,
		CurrentMessage:   CurrentAnnouncement.Message,
		CurrentTextColor: CurrentAnnouncement.TextColor,
		CurrentBgColor:   CurrentAnnouncement.BgColor,
	}

	if err := tpl.ExecuteTemplate(w, "announcement-form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// saveAnnouncement は投稿された告知内容をメモリに保存する
func saveAnnouncement(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	message := r.FormValue("message")
	textColor := r.FormValue("textcolor")
	bgColor := r.FormValue("bgcolor")
	enabled := r.FormValue("enabled") == "on"

	// カラーバリデーション（簡易版）
	if !isValidColor(textColor) {
		textColor = "#000000" // デフォルト：黒
	}
	if !isValidColor(bgColor) {
		bgColor = "#ffffe0" // デフォルト：薄黄
	}

	// 告知を保存
	CurrentAnnouncement = AnnouncementData{
		Enabled:   enabled,
		Message:   template.HTMLEscapeString(message),
		TextColor: textColor,
		BgColor:   bgColor,
	}

	// 成功レスポンス
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>告知を保存しました</title></head>
<body>
<h2>✓ 告知を保存しました</h2>
<p>入力内容:</p>
<ul>
  <li>表示: %v</li>
  <li>メッセージ: %s</li>
  <li>文字色: %s</li>
  <li>背景色: %s</li>
</ul>
<p>このウィンドウを閉じてください。新しいリクエストから告知が反映されます。</p>
</body>
</html>`, enabled, CurrentAnnouncement.Message, textColor, bgColor)
}

// isValidColor は簡易的なカラーバリデーション（#RRGGBB形式）
func isValidColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	validHex := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	return validHex.MatchString(color)
}
