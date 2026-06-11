/*
!
Copyright © 2024 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package ShowroomCGIlib

import (
	// "fmt"
	//	"io"
	// "log"
	//	"os"
	// "strconv"
	// "strings"
	// "time"

	// "net/http"

	//	"github.com/go-gorp/gorp"
	//      "gopkg.in/gorp.v2"

	// "github.com/dustin/go-humanize"
)

/*

0.0.1 UserでGenreが空白の行のGenreをSHOWROOMのAPIで更新する。
0.0.2 Userでirankが-1の行のランクが空白の行のランク情報をSHOWROOMのAPIで更新する。
0.1.0 DBのアクセスにgorpを導入する。
0.1.1 database/sqlを使った部分（コメント）を削除する
0.2.0 Event.goを追加し、User.goにEvent, Eventuser, Wuserを追加する。

*/

//	const Version = "0.1.1"

/*
type User struct {
	Userno       int
	Userid       string
	User_name    string
	Longname     string
	Shortname    string
	Genre        string
	GenreID      int
	Rank         string
	Nrank        string
	Prank        string
	Irank        int
	Inrank       int
	Iprank       int
	Itrank       int
	Level        int
	Followers    int
	Fans         int
	FanPower     int
	Fans_lst     int
	FanPower_lst int
	Ts           time.Time
	Getp         string
	Graph        string
	Color        string
	Currentevent string
}

type Wuser struct {
	Userno    int
	Userid    string
	User_name string
	Longname  string
	Shortname string
	Genre     string
	// GenreID      int
	Rank  string
	Nrank string
	Prank string
	// Irank        int
	// Inrank       int
	// Iprank       int
	// Itrank       int
	Level        int
	Followers    int
	Fans         int
	FanPower     int
	Fans_lst     int
	FanPower_lst int
	Ts           time.Time
	Getp         string
	Graph        string
	Color        string
	Currentevent string
}
*/


/*
type Wuser User

type Wuser struct {
	Userno       int
	Userid       string
	User_name    string
	Longname     string
	Shortname    string
	Genre        string
	Rank         string
	Nrank        string
	Prank        string
	Level        int
	Followers    int
	Fans         int
	Fans_lst     int
	Ts           time.Time
	Getp         string
	Graph        string
	Color        string
	Currentevent string
}
*/

// Rank情報からランクのソートキーを作る
func MakeSortKeyOfRank(rank string, nextscore int) (
	irank int,
) {
	r2n := map[string]int{
		"SS-5":   100000000,
		"SS-4":   200000000,
		"SS-3":   300000000,
		"SS-2":   400000000,
		"SS-1":   500000000,
		"S-5":    520000000,
		"S-4":    540000000,
		"S-3":    560000000,
		"S-2":    580000000,
		"S-1":    600000000,
		"A-5":    610000000,
		"A-4":    620000000,
		"A-3":    630000000,
		"A-2":    640000000,
		"A-1":    650000000,
		"B-5":    660000000,
		"B-4":    670000000,
		"B-3":    680000000,
		"B-2":    690000000,
		"B-1":    700000000,
		"C-10":   710000000,
		"C-9":    720000000,
		"C-8":    730000000,
		"C-7":    740000000,
		"C-6":    750000000,
		"C-5":    760000000,
		"C-4":    770000000,
		"C-3":    780000000,
		"C-2":    790000000,
		"C-1":    800000000,
		"unknown": 1000000000, //	SHOWROOMのアカウントを削除した配信者さん
		//	888888888: irank 未算出
	}

	if sk, ok := r2n[rank]; ok {
		irank = sk + nextscore
	} else {
		irank = 999999999 //	(アイドルで)SHOWRANKの対象ではない配信者さん
	}

	return
}

