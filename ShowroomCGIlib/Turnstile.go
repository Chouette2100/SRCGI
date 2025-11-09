// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Cloudflare Turnstileの検証レスポンス構造体
type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
}

// VerifyTurnstile はCloudflare Turnstileトークンを検証する
// token: ブラウザから送信されたトークン
// remoteIP: リクエスト元のIPアドレス(オプション)
// 戻り値: 検証が成功したかどうか、エラー
func VerifyTurnstile(token string, remoteIP string) (bool, error) {
	// シークレットキーが設定されていない場合は検証をスキップ
	if Serverconfig.TurnstileSecretKey == "" {
		return true, nil
	}

	// トークンが空の場合は検証失敗
	if token == "" {
		return false, fmt.Errorf("turnstile token is empty")
	}

	// Cloudflare APIエンドポイント
	verifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	// リクエストボディを作成
	requestBody := map[string]string{
		"secret":   Serverconfig.TurnstileSecretKey,
		"response": token,
	}

	// remoteIPが指定されている場合は追加
	if remoteIP != "" {
		requestBody["remoteip"] = remoteIP
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf("JSON marshal error: %w", err)
	}

	// HTTPクライアントを作成(タイムアウト設定)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// POSTリクエストを送信
	resp, err := client.Post(verifyURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("turnstile API request error: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("response read error: %w", err)
	}

	// JSONをパース
	var turnstileResp TurnstileResponse
	if err := json.Unmarshal(body, &turnstileResp); err != nil {
		return false, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	// 検証結果をチェック
	if !turnstileResp.Success {
		return false, fmt.Errorf("turnstile verification failed: %v", turnstileResp.ErrorCodes)
	}

	return true, nil
}
