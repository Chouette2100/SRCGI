// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// Turnstileセッション用のクッキー名
	TurnstileCookieName = "turnstile_verified"
	// HMAC署名用の秘密鍵（本番環境では環境変数等から読み込むべき）
	// ここでは簡易的に固定値を使用（後でServerConfigに追加することも検討）
	TurnstileSecretSalt = "SRCGI-Turnstile-Session-Secret-2025"
)

// TurnstileSession はTurnstile検証済みセッションの情報を保持する
type TurnstileSession struct {
	IP        string // IPアドレス
	Timestamp int64  // 作成時刻（Unix timestamp）
	Counter   int    // 使用回数
	Signature string // HMAC署名
}

// CreateTurnstileSessionCookie はTurnstile検証済みセッション用のクッキーを作成する
func CreateTurnstileSessionCookie(ip string) *http.Cookie {
	if !Serverconfig.TurnstileUseSession {
		return nil
	}

	timestamp := time.Now().Unix()
	counter := 0

	// データを結合
	data := fmt.Sprintf("%s|%d|%d", ip, timestamp, counter)

	// HMAC署名を生成
	signature := generateHMAC(data)

	// クッキー値を作成（Base64エンコード）
	cookieValue := base64.StdEncoding.EncodeToString([]byte(data + "|" + signature))

	// クッキーを作成
	cookie := &http.Cookie{
		Name:     TurnstileCookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   Serverconfig.SSLcrt != "", // HTTPSの場合のみSecure属性を設定
		SameSite: http.SameSiteStrictMode,
		MaxAge:   Serverconfig.TurnstileSessionDuration,
	}

	return cookie
}

// VerifyTurnstileSessionCookie はTurnstileセッションクッキーを検証する
// 戻り値: (検証OK, 新しいクッキー（更新が必要な場合）, エラー)
func VerifyTurnstileSessionCookie(r *http.Request) (bool, *http.Cookie, error) {
	// セッション機能が無効の場合は常にfalse
	if !Serverconfig.TurnstileUseSession {
		return false, nil, nil
	}

	// クッキーを取得
	cookie, err := r.Cookie(TurnstileCookieName)
	if err != nil {
		// クッキーがない場合
		return false, nil, nil
	}

	// Base64デコード
	decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return false, nil, fmt.Errorf("invalid cookie format: %w", err)
	}

	// データを分割
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 4 {
		return false, nil, fmt.Errorf("invalid cookie structure")
	}

	cookieIP := parts[0]
	timestampStr := parts[1]
	counterStr := parts[2]
	signature := parts[3]

	// タイムスタンプをパース
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false, nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	// カウンターをパース
	counter, err := strconv.Atoi(counterStr)
	if err != nil {
		return false, nil, fmt.Errorf("invalid counter: %w", err)
	}

	// 署名を検証
	data := fmt.Sprintf("%s|%d|%d", cookieIP, timestamp, counter)
	expectedSignature := generateHMAC(data)
	if signature != expectedSignature {
		return false, nil, fmt.Errorf("invalid signature")
	}

	// IPアドレスを取得
	currentIP := RemoteAddr(r)

	// IPアドレスが一致するか確認
	if cookieIP != currentIP {
		return false, nil, fmt.Errorf("IP mismatch: cookie=%s, current=%s", cookieIP, currentIP)
	}

	// 有効期限を確認
	now := time.Now().Unix()
	age := now - timestamp
	if age > int64(Serverconfig.TurnstileSessionDuration) {
		return false, nil, fmt.Errorf("session expired: age=%d, max=%d", age, Serverconfig.TurnstileSessionDuration)
	}

	// 使用回数を確認（0は無制限）
	if Serverconfig.TurnstileSessionMaxUses > 0 && counter >= Serverconfig.TurnstileSessionMaxUses {
		return false, nil, fmt.Errorf("max uses exceeded: %d", counter)
	}

	// カウンターをインクリメント
	newCounter := counter + 1
	newData := fmt.Sprintf("%s|%d|%d", cookieIP, timestamp, newCounter)
	newSignature := generateHMAC(newData)
	newCookieValue := base64.StdEncoding.EncodeToString([]byte(newData + "|" + newSignature))

	// 新しいクッキーを作成
	newCookie := &http.Cookie{
		Name:     TurnstileCookieName,
		Value:    newCookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   Serverconfig.SSLcrt != "",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   Serverconfig.TurnstileSessionDuration - int(age), // 残り時間
	}

	return true, newCookie, nil
}

// generateHMAC はHMAC-SHA256署名を生成する
func generateHMAC(data string) string {
	// シークレットキーとソルトを結合して使用
	secret := Serverconfig.TurnstileSecretKey + TurnstileSecretSalt
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
