package ShowroomCGIlib

import (
	//	"fmt"
	//	"io/ioutil"
	//	"os"

	//	"gopkg.in/yaml.v2"
)

/*

	20A00	結果をDBで保存する。Excel保存の機能は残存。次に向けての作り込み少々。
	2.0B00		データ取得のタイミングをtimetableから得る。Excelへのデータの保存をやめる。
	2.0B01	timetableの更新で処理が終わっていないものを処理済みにしていた問題を修正する。
	10AJ00	ブロックランキングに仮対応（Event_id=30030以外に拡張）する。イベントリストの表示イベント数を設定可能とする。
	11AA00	データベースへのアクセスをsrdblibに移行しつつある。

*/

const VerFileIOlib = "11AA00"

type ServerConfig struct {
	WebServer string `yaml:"WebServer"`
	HTTPport  string `yaml:"HTTPport"`
	SSLcrt    string `yaml:"SSLcrt"`
	SSLkey    string `yaml:"SSLkey"`
	Dbhost    string `yaml:"Dbhost"`
	Dbport    string `yaml:"Dbport"`
	Dbname    string `yaml:"Dbname"`
	Dbuser    string `yaml:"Dbuser"`
	Dbpw      string `yaml:"Dbpw"`
	UseSSH    bool   `yaml:"UseSSH"`
	NoEvent   int    `yaml:"NoEvent"` //	イベント一覧に表示するイベントの数
}

type SSHConfig struct {
	Hostname   string `yaml:"Hostname"`
	Port       int    `yaml:"Port"`
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	PrivateKey string `yaml:"PrivateKey"`
}
