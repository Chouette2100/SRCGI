package ShowroomCGIlib

import (
	// "bytes"
	"context" // contextパッケージを追加
	"encoding/json"
	"fmt"
	"io" // io.ReadAll のために追加
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time" // time.Duration のために追加

	"github.com/Chouette2100/srdblib/v2" // srdblibパッケージのパスは適宜修正
)

// Version はこのクライアントライブラリのバージョン
// ビルド時にリンカーフラグで注入することも可能
// const Version = "100000" // 例: 1.0.0 を 100000 と表現

var useragent string
var mailAddress = "https://chouette2100.com/disp-bbs" // User-Agentのコメント部分に含めるURL

func init() {
	// GoのバージョンとOS/Archを取得
	goVersion := runtime.Version()
	goOS := runtime.GOOS
	goArch := runtime.GOARCH

	var major, minor, patch int
	// Versionが"100000"の場合、major=1, minor=0, patch=0 となるようにパース
	// もしVersionが"1.0.0"のような文字列なら fmt.Sscanf(Version, "%d.%d.%d", ...)
	// %1d%2d%2d は 100000 -> 1, 00, 00 とパースする
	fmt.Sscanf(Version, "%1d%3d%2d", &major, &minor, &patch)

	// RFC 7231 Section 5.5.3 User-Agent Header Field に従い、コメント形式で連絡先を含める
	useragent = fmt.Sprintf("SRCGI/v%d.%d.%d (contact: %s; Go/%s; %s/%s)",
		major, minor, patch, mailAddress, goVersion, goOS, goArch)

	// initでのログ出力は、パッケージロード時に毎回実行されるため、
	// 開発時デバッグ用か、本当に必要な情報に限定するのが良いでしょう。
	log.Printf("ShowroomCGIlib initialized. User-Agent: %s\n", useragent)
}

// AddNewUserRequest はAddNewUser関数に渡すパラメータを構造体で定義します。
// これにより、引数の数が増えても管理しやすくなります。
type AddNewUserRequest struct {
	RoomID        string
	RoomName      string
	IsImmediately bool
	Client        *http.Client  // HTTPクライアント (nilの場合、http.DefaultClientを使用)
	Timeout       time.Duration // リクエストごとのタイムアウト (0の場合、クライアントのデフォルトを使用)
}

func AddNewUser(reqData AddNewUserRequest) (
	user *srdblib.User,
	msg string,
	err error,
) {
	if reqData.Client == nil {
		reqData.Client = http.DefaultClient
	}

	// リクエストごとのタイムアウトを設定するContextを作成
	// reqData.Timeout が 0 の場合は、クライアントのデフォルトタイムアウトに任せる
	ctx := context.Background()
	var cancel context.CancelFunc = func() {} // デフォルトのno-op関数
	if reqData.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), reqData.Timeout)
	}
	defer cancel() // 忘れずにキャンセルを呼び出す

	var turl string
	if reqData.IsImmediately {
		turl = Serverconfig.GWURL + "AddNewUserImmediately"
	} else {
		turl = Serverconfig.GWURL + "AddNewUser"
	}

	u, err := url.Parse(turl)
	if err != nil {
		err = fmt.Errorf("failed to parse URL '%s': %w", turl, err)
		return
	}

	// クエリを組み立て
	values := url.Values{}
	values.Add("roomid", reqData.RoomID) // roomd -> roomid に修正
	if reqData.RoomName != "" {
		values.Add("roomname", reqData.RoomName)
	}

	// Request を生成
	// context.Context を NewRequestWithContext に渡す
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		err = fmt.Errorf("http.NewRequestWithContext() failed for %s: %w", turl, err)
		return
	}

	// 組み立てたクエリを生クエリ文字列に変換して設定
	req.URL.RawQuery = values.Encode()

	// User-Agentを設定
	req.Header.Add("User-Agent", useragent)
	req.Header.Add("Accept-Language", "ja-JP")

	// Doメソッドでリクエストを投げる
	resp, err := reqData.Client.Do(req)
	if err != nil {
		// context.DeadlineExceeded や context.Canceled のエラーもここで捕捉される
		err = fmt.Errorf("client.Do() failed for %s: %w", turl, err)
		return
	}
	defer resp.Body.Close()

	// レスポンスボディを一度読み込む (エラーログのために)
	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		err = fmt.Errorf("failed to read response body from %s: %w", turl, readErr)
		return
	}
	bodyString := string(bodyBytes)

	// HTTPステータスコードのチェック
	if resp.StatusCode >= 400 { // 4xx, 5xx はエラー
		err = fmt.Errorf("server returned error status %d for %s: %s", resp.StatusCode, turl, bodyString)
		return
	}

	if reqData.IsImmediately {
		// AddNewUserImmediately は 200 OK を期待
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("expected status %d OK for immediate request to %s, got %d: %s", http.StatusOK, turl, resp.StatusCode, bodyString)
			return
		}
		if err = json.Unmarshal(bodyBytes, &user); err != nil { // bytes.Bufferを使わず直接Unmarshal
			err = fmt.Errorf("failed to decode JSON response from %s: %w (response body: %s)", turl, err, bodyString)
			return
		}
	} else {
		// AddNewUser は 202 Accepted を期待
		if resp.StatusCode != http.StatusAccepted {
			err = fmt.Errorf("expected status %d Accepted for async request to %s, got %d: %s", http.StatusAccepted, turl, resp.StatusCode, bodyString)
			return
		}
		msg = bodyString
	}
	return
}
