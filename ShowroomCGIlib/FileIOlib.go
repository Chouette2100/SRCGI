package ShowroomCGIlib

import (
	//	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

/*

	20A00	結果をDBで保存する。Excel保存の機能は残存。次に向けての作り込み少々。
	2.0B00		データ取得のタイミングをtimetableから得る。Excelへのデータの保存をやめる。
	2.0B01	timetableの更新で処理が終わっていないものを処理済みにしていた問題を修正する。
	10AJ00	ブロックランキングに仮対応（Event_id=30030以外に拡張）する。イベントリストの表示イベント数を設定可能とする。

*/

const VerFileIOlib = "10AJ00"

type DBConfig struct {
	WebServer string `yaml:"WebServer"`
	HTTPport  string `yaml:"HTTPport"`
	SSLcrt    string `yaml:"SSLcrt"`
	SSLkey    string `yaml:"SSLkey"`
	Dbhost    string `yaml:"Dbhost"`
	Dbname    string `yaml:"Dbname"`
	Dbuser    string `yaml:"Dbuser"`
	Dbpw      string `yaml:"Dbpw"`
	NoEvent   int    `yaml:"NoEvent"`	//	イベント一覧に表示するイベントの数
}

// 設定ファイルを読み込む
//      以下の記事を参考にさせていただきました。
//              【Go初学】設定ファイル、環境変数から設定情報を取得する
//                      https://note.com/artefactnote/n/n8c22d1ac4b86
//
func LoadConfig(filePath string) (dbconfig *DBConfig, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	content = []byte(os.ExpandEnv(string(content)))

	result := &DBConfig{}
	result.NoEvent = 30
	if err := yaml.Unmarshal(content, result); err != nil {
		return nil, err
	}

	return result, nil
}
