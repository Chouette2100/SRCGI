// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"html/template"
	"log"
	"net/http"
)

// TurnstileCheckResult はTurnstile検証の結果を表す
type TurnstileCheckResult int

const (
	// TurnstileOK は検証に成功したか、セッションが有効な状態
	TurnstileOK TurnstileCheckResult = iota
	// TurnstileNeedChallenge はチャレンジページを表示する必要がある状態
	TurnstileNeedChallenge
	// TurnstileFailed は検証に失敗した状態
	TurnstileFailed
)

// TurnstileChallengeData はチャレンジページ表示用のデータインターフェース
// 各ハンドラー固有の構造体がこのインターフェースを実装する必要がある
type TurnstileChallengeData interface {
	// SetTurnstileInfo はTurnstile関連の情報をセットする
	SetTurnstileInfo(siteKey string, errorMsg string)
	// GetTemplatePath はテンプレートファイルのパスを返す
	GetTemplatePath() string
	// GetTemplateName はテンプレート名を返す
	GetTemplateName() string
	// GetFuncMap はテンプレートの関数マップを返す
	GetFuncMap() *template.FuncMap
}

// CheckTurnstileWithSession はTurnstileの検証とセッション管理を行う汎用関数
// この関数を呼び出すことで、各ハンドラーでのTurnstile処理を統一できる
//
// 引数:
//   - w: http.ResponseWriter
//   - r: *http.Request
//   - challengeData: チャレンジページ表示用のデータ（TurnstileChallengeDataを実装）
//
// 戻り値:
//   - TurnstileCheckResult: 検証結果
//   - error: エラー（あれば）
//
// 使用例:
//
//	result, err := CheckTurnstileWithSession(w, r, &myData)
//	switch result {
//	case TurnstileOK:
//	    // 検証OK、通常処理を続行
//	case TurnstileNeedChallenge, TurnstileFailed:
//	    // チャレンジページまたはエラーページが表示済み
//	    return
//	}
func CheckTurnstileWithSession(
	w http.ResponseWriter,
	r *http.Request,
	challengeData TurnstileChallengeData,
) (TurnstileCheckResult, error) {
	// Turnstileが無効な場合はOKを返す
	if Serverconfig.TurnstileSiteKey == "" {
		return TurnstileOK, nil
	}

	// まずセッションクッキーを確認
	sessionValid, newCookie, sessionErr := VerifyTurnstileSessionCookie(r)

	if sessionValid {
		// セッションが有効な場合
		log.Printf("Turnstile session valid for IP %s\n", RemoteAddr(r))

		// カウンターを更新したクッキーをセット
		if newCookie != nil {
			http.SetCookie(w, newCookie)
		}
		return TurnstileOK, nil
	}

	// セッションが無効または存在しない場合
	if sessionErr != nil {
		log.Printf("Turnstile session check failed: %v\n", sessionErr)
	}

	// Turnstileトークンを取得
	turnstileToken := r.FormValue("cf-turnstile-response")

	// トークンが空の場合は初回アクセスなので、チャレンジページを表示
	if turnstileToken == "" {
		return showTurnstileChallenge(w, challengeData, "")
	}

	// Turnstile検証を実行
	remoteIP := RemoteAddr(r)
	verified, err := VerifyTurnstile(turnstileToken, remoteIP)
	if err != nil || !verified {
		log.Printf("Turnstile verification failed for IP %s: %v\n", remoteIP, err)
		return showTurnstileChallenge(w, challengeData, "セキュリティチェックに失敗しました。もう一度お試しください。")
	}

	// 検証成功：セッションクッキーを作成
	sessionCookie := CreateTurnstileSessionCookie(remoteIP)
	if sessionCookie != nil {
		http.SetCookie(w, sessionCookie)
		log.Printf("Turnstile session created for IP %s\n", remoteIP)
	}

	return TurnstileOK, nil
}

// showTurnstileChallenge はチャレンジページを表示する内部関数
func showTurnstileChallenge(
	w http.ResponseWriter,
	challengeData TurnstileChallengeData,
	errorMsg string,
) (TurnstileCheckResult, error) {
	// Turnstile情報をセット
	challengeData.SetTurnstileInfo(Serverconfig.TurnstileSiteKey, errorMsg)

	// テンプレートをパース
	var tpl *template.Template
	var err error
	if challengeData.GetFuncMap() != nil {
		tpl = template.Must(template.New("").Funcs(*challengeData.GetFuncMap()).
			ParseFiles(challengeData.GetTemplatePath(), "templates/turnstilechallenge.gtpl"))
	} else {
		tpl, err = template.ParseFiles(challengeData.GetTemplatePath(), "templates/turnstilechallenge.gtpl")
		if err != nil {
			log.Printf("Template parse error: %s\n", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			if errorMsg != "" {
				return TurnstileFailed, err
			}
			return TurnstileNeedChallenge, err
		}
	}

	// テンプレートを実行
	if err := tpl.ExecuteTemplate(w, challengeData.GetTemplateName(), challengeData); err != nil {
		log.Printf("Template execution error: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		if errorMsg != "" {
			return TurnstileFailed, err
		}
		return TurnstileNeedChallenge, err
	}

	if errorMsg != "" {
		return TurnstileFailed, nil
	}
	return TurnstileNeedChallenge, nil
}

// SimpleTurnstileCheck はセッション管理なしのシンプルなTurnstile検証
// レガシーコードとの互換性や、セッション不要な場合に使用
//
// 引数:
//   - w: http.ResponseWriter
//   - r: *http.Request
//   - challengeData: チャレンジページ表示用のデータ
//
// 戻り値:
//   - TurnstileCheckResult: 検証結果
//   - error: エラー（あれば）
func SimpleTurnstileCheck(
	w http.ResponseWriter,
	r *http.Request,
	challengeData TurnstileChallengeData,
) (TurnstileCheckResult, error) {
	// Turnstileが無効な場合はOKを返す
	if Serverconfig.TurnstileSiteKey == "" {
		return TurnstileOK, nil
	}

	// Turnstileトークンを取得
	turnstileToken := r.FormValue("cf-turnstile-response")

	// トークンが空の場合は初回アクセスなので、チャレンジページを表示
	if turnstileToken == "" {
		return showTurnstileChallenge(w, challengeData, "")
	}

	// Turnstile検証を実行
	remoteIP := RemoteAddr(r)
	verified, err := VerifyTurnstile(turnstileToken, remoteIP)
	if err != nil || !verified {
		log.Printf("Turnstile verification failed for IP %s: %v\n", remoteIP, err)
		return showTurnstileChallenge(w, challengeData, "セキュリティチェックに失敗しました。もう一度お試しください。")
	}

	return TurnstileOK, nil
}
