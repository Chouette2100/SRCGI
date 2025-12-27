package ShowroomCGIlib

//	"fmt"
//	"io/ioutil"
//	"os"

// "gopkg.in/yaml.v2"

/*

	20A00	結果をDBで保存する。Excel保存の機能は残存。次に向けての作り込み少々。
	2.0B00		データ取得のタイミングをtimetableから得る。Excelへのデータの保存をやめる。
	2.0B01	timetableの更新で処理が終わっていないものを処理済みにしていた問題を修正する。
	10AJ00	ブロックランキングに仮対応（Event_id=30030以外に拡張）する。イベントリストの表示イベント数を設定可能とする。
	11AA00	データベースへのアクセスをsrdblibに移行しつつある。
	SRCGI.00AM02	通常とメンテナンスの切り替えを ShowroomCGIlib.Serverconfig.Maintenance で行う。

*/

const VerFileIOlib = "11AA00"

type ServerConfig struct {
	WebServer                string `yaml:"WebServer"`
	HTTPport                 string `yaml:"HTTPport"`
	SSLcrt                   string `yaml:"SSLcrt"`
	SSLkey                   string `yaml:"SSLkey"`
	NoEvent                  int    `yaml:"NoEvent"` //	イベント一覧に表示するイベントの数
	Maintenance              bool   `yaml:"Maintenance"`
	LvlBots                  int    `yaml:"LvlBots"`                  //	Bot排除のレベル、0:なし、1:低、2:中、3:高
	AccessLimit              int    `yaml:"AccessLimit"`              //	時間枠のアクセス回数上限
	TimeWindow               int    `yaml:"TimeWindow"`               //	アクセス回数制限の時間枠（秒単位）
	MaxChlog                 int    `yaml:"MacChlog"`                 //	ログ出力待ちチャンネルのバッファ数（＝同時実行ハンドラー数）
	TurnstileSiteKey         string `yaml:"TurnstileSiteKey"`         //	Cloudflare Turnstileのサイトキー
	TurnstileSecretKey       string `yaml:"TurnstileSecretKey"`       //	Cloudflare Turnstileのシークレットキー
	TurnstileUseSession      bool   `yaml:"TurnstileUseSession"`      //	Turnstileセッション機能の有効/無効
	TurnstileSessionDuration int    `yaml:"TurnstileSessionDuration"` //	Turnstileセッション有効期間（秒）
	TurnstileSessionMaxUses  int    `yaml:"TurnstileSessionMaxUses"`  //	Turnstileセッション内の最大利用回数
	DenyNonJP                bool   `yaml:"DenyNonJP"`                //	日本国外からのアクセスを拒否するかどうか
	GWURL                    string `yaml:"GWURL"`                    //  AddNewUserの GatewayのURL
}

type SSHConfig struct {
	Hostname   string `yaml:"Hostname"`
	Port       int    `yaml:"Port"`
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	PrivateKey string `yaml:"PrivateKey"`
}
