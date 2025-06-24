package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var LogFile_f2b *os.File
var fail2banLogger *log.Logger

func init() {
	var err error
	// fail2ban 用のログファイルを開く
	// 例: /var/log/your-webserver/fail2ban.log
	// ディレクトリが存在しない場合は作成が必要です。
	logFilePath := "/var/log/SRCGI/fail2ban.log"
	logFileDir := "/var/log/SRCGI"
	if err = os.MkdirAll(logFileDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory %s: %v", logFileDir, err)
	}

	LogFile_f2b, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open fail2ban log file %s: %v", logFilePath, err)
	}
	// アプリケーション終了時にファイルを閉じる処理を検討してください
	// defer logFile.Close() // main 関数など適切な場所で

	// タイムスタンプやファイル名などのデフォルトフラグを無効にする
	fail2banLogger = log.New(LogFile_f2b, "", 0)
}

// fail2ban 用のログを出力する関数
func LogForFail2ban(ip, handlerName, botInfo string) {
	// YYYY-MM-DD HH:MM:SS 形式のタイムスタンプを生成
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// ログメッセージを組み立てる
	// ip は必須
	msg := fmt.Sprintf("%s [INFO] ip=%s", timestamp, ip)

	// ハンドラー名があれば追加
	if handlerName != "" {
		msg += fmt.Sprintf(" handler=%s", handlerName)
	}

	// ボット情報があれば追加
	if botInfo != "" {
		msg += fmt.Sprintf(" bot=%s", botInfo)
	}

	// ログファイルに出力
	fail2banLogger.Println(msg)
}
