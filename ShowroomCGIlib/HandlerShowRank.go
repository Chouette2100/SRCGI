/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

package ShowroomCGIlib

import (
	//	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"html/template"
	"net/http"
	// "net/http/cookiejar"
	"github.com/juju/persistent-cookiejar"

	// "net/http/cookiejar"

	//	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v3"
)

type Eruser struct {
	srdblib.User
	Tmrank   string
	Tminrank int
	Tmiprank int
	Tmts     time.Time
}

type ShowRank struct {
	HdrMsg string
	// Userlist  *[]srdblib.User

	Userlist  *[]Eruser
	UserlistA *[]Eruser
	ErrMsg    string
}

// SHOWランク上位ルームを抽出する
func SelectShowRank(
	client *http.Client,
	limit int,
) (
	userlist *[]Eruser,
	usermap map[int]*Eruser,

	err error,
) {

	userlist = new([]Eruser)

	// sqltr := " select " + clmlist["user"] + " from user where irank between 0 and ? and ts > ? and fanpower > 0 order by irank "
	sqltr := " select " + clmlist["user"] + " from user where irank between 0 and ? and ts > ? order by irank "

	Dbmap0.AddTableWithName(Eruser{}, "user").SetKeys(false, "Userno")
	defer Dbmap0.AddTableWithName(srdblib.User{}, "user").SetKeys(false, "Userno")
	var ul []interface{}
	ul, err = Dbmap0.Select(Eruser{}, sqltr, limit, time.Now().Add(-time.Hour*25))
	if err != nil {
		err = fmt.Errorf("Dbmap0.Select(): %w", err)
		return
	}

	usermap = make(map[int]*Eruser)
	for _, v := range ul {
		user := v.(*Eruser)
		*userlist = append(*userlist, *user)
		usermap[user.Userno] = user
	}

	return
}

// 月始めのSHOWランク上位ルームを抽出する
func SelectTmShowRank(
	client *http.Client,
) (
	userlist *[]Eruser,
	usermap map[int]*Eruser, // usernoに対するuserhistory
	err error,
) {
	userlist = new([]Eruser)

	// 現在時から年月を求める
	yy, mm, _ := time.Now().Date()
	// 現在の年月の初日を求める
	tmfirst := time.Date(yy, mm, 1, 0, 0, 0, 0, time.Local)
	// データ取得の開始、終了をそれぞれ年月所持つの00時29分、00時35分とする
	// ※ 取得の開始時刻はカテゴリー単位のもので、個別のデータ取得時刻とは異なる
	tb := tmfirst.Add(29 * time.Minute)
	te := tmfirst.Add(35 * time.Minute)

	usermap = make(map[int]*Eruser)
	// sqltr := " select " + clmlist["user"] + " from user where irank between 0 and ? and ts > ? and fanpower > 0 order by irank "
	// sqltr := " select userno, `rank`,nrank, prank, ts from userhistory "
	sqltr := " select "
	sqltr += " userno, user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, "
	sqltr += "ts from userhistory where ts between ? and ? order by ts desc "

	Dbmap0.AddTableWithName(Eruser{}, "userhistory").SetKeys(false, "Userno")
	defer Dbmap0.AddTableWithName(Eruser{}, "userhistory").SetKeys(false, "Userno")
	var ul []interface{}
	ul, err = Dbmap0.Select(Eruser{}, sqltr, tb, te)
	if err != nil {
		err = fmt.Errorf("Dbmap0.Select(): %w", err)
		return
	}

	for _, v := range ul {
		user := v.(*Eruser)
		*userlist = append(*userlist, *user)
		usermap[user.Userno] = user
	}
	return
}

func SelectAddedRooms(nolist []int) (
	pul *[]Eruser,
	err error,
) {

	pul = new([]Eruser)
	// var intf []interface{}
	sqltr := " select " + clmlist["user"] + " from user where userno in (:Users) "
	// intf, err = Dbmap0.Select(srdblib.User{}, sqltr, map[string]interface{}{"Users": nolist})
	_, err = Dbmap0.Select(pul, sqltr, map[string]interface{}{"Users": nolist})
	if err != nil {
		err = fmt.Errorf("Dbmap0.Select(): %w", err)
		log.Printf("SelectAddedRooms(): %s\n", err.Error())
		return
	}

	// for _, v := range intf {
	// 	addedlist = append(addedlist, *(v.(*Eruser)))
	// }

	return

}

/*

	HandlerShowRank()
		SHOWランク上位配信者を表示する

	Ver. 0.1.0

*/
// http://localhost:8080/showrank で呼び出される
func ShowRankHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	showrank := &ShowRank{}

	//	cookiejarがセットされたHTTPクライアントを作る
	var err error
	var client *http.Client
	var jar *cookiejar.Jar
	client, jar, err = exsrapi.CreateNewClient("XXXXXX")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//	すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	// SHOWランクとは無関係にランク状況を知りたいルームを追加したとき
	unlist := r.FormValue("unlist")
	unla := strings.Split(unlist, ",")

	lmin := srdblib.Env.Lmin
	waitmsec := srdblib.Env.Waitmsec
	srdblib.Env.Lmin = 60
	srdblib.Env.Waitmsec = 1000

	// userlist := make([]srdblib.User, 0, len(unla))
	nolist := make([]int, 0, len(unla))
	user := srdblib.User{}
	for _, v := range unla {
		un, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("strconv.Atoi() returned error %s\n", err.Error())
			continue
		}
		user.Userno = un
		_, err = srdblib.UpinsUser(Dbmap0, client, time.Now(), &user)
		if err != nil {
			log.Printf("srdblib.UpinsUser() returned error %s\n", err.Error())
			continue
		}
		// userlist = append(userlist, user)
		nolist = append(nolist, un)
	}
	srdblib.Env.Lmin = lmin
	srdblib.Env.Waitmsec = waitmsec

	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"Comma":      func(i int) string { return humanize.Comma(int64(i)) }, //	3桁ごとに","を入れる関数。
		"FormatTime": func(t time.Time, tfmt string) string { return t.Format(tfmt) },
		"Add":        func(a int, b int) int { return a + b },
		"Showrank": func(s string) string {
			sa := strings.Split(s, " | ")
			return sa[len(sa)-1]
		},
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/showrank.gtpl"))

	//	showrank.Userlist, err = SelectShowRank(client, 260000000)	//	SS-5〜A-5
	//	showrank.Userlist, err = SelectShowRank(client, 300000000)	//	SS-5〜A-1
	var user1 Eruser
	// FIXME: irank != 0 の条件が必要な理由を明確にすること(2025-05-14)
	sqlst := "select " + clmlist["user"] + " from user where irank = (select min(irank) from user where `rank` = 'B-5' and irank != 0) "
	err = Dbmap0.SelectOne(&user1, sqlst)
	if err != nil {
		err = fmt.Errorf("Dbmap0.SelectOne(): %w", err)
		log.Printf("HandlerShowRank(): %s\n", err.Error())
		return
	}

	showrank.Userlist, _, err = SelectShowRank(client, user1.Irank+1) //	SS-5〜A-1とB-5のトップ
	if err != nil {
		err = fmt.Errorf("SelectShowRank(): %w", err)
		log.Printf("HandlerShowRank(): %s\n", err.Error())
		return
	}
	_, usermap, err := SelectTmShowRank(client)
	if err != nil {
		err = fmt.Errorf("SelectTmShowRank(): %w", err)
		log.Printf("HandlerShowRank(): %s\n", err.Error())
		return
	}
	for i := range *showrank.Userlist {
		user := &(*showrank.Userlist)[i]
		tirank := MakeSortKeyOfRank(user.Rank, user.Inrank)
		if tirank != user.Irank {
			log.Printf("WARNING: userno=%d rank=%s nrank=%s irank=%d tirank=%d\n", user.Userno, user.Rank, user.Nrank, user.Irank, tirank)
			user.Irank = tirank
		}
		if uh, ok := usermap[user.Userno]; ok {
			user.Tmrank = uh.Rank
			// 9,999,999形式のランクを整数に変換する
			irank, err := strconv.Atoi(strings.ReplaceAll(uh.Nrank, ",", ""))
			if err != nil {
				log.Printf("strconv.Atoi() returned error %s\n", err.Error())
				continue
			}
			user.Tminrank = irank
			iprank, err := strconv.Atoi(strings.ReplaceAll(uh.Prank, ",", ""))
			if err != nil {
				log.Printf("strconv.Atoi() returned error %s\n", err.Error())
				continue
			}
			user.Tmiprank = iprank
			user.Tmts = uh.Ts
		}
	}

	// showrank.UserlistをIrankの降順にソートする
	sort.Slice(*showrank.Userlist, func(i, j int) bool {
		return (*showrank.Userlist)[j].Irank > (*showrank.Userlist)[i].Irank
	})

	if len(nolist) != 0 {
		showrank.UserlistA, err = SelectAddedRooms(nolist)
		if err != nil {
			err = fmt.Errorf("SelectAddedRooms(): %w", err)
			log.Printf("HandlerShowRank(): %s\n", err.Error())
			w.Write([]byte(fmt.Sprintf("SelectAddedRooms() error=%s\n", err.Error())))
			return
		}

		for i := range *showrank.UserlistA {
			user := &(*showrank.UserlistA)[i]
			tirank := MakeSortKeyOfRank(user.Rank, user.Inrank)
			if tirank != user.Irank {
				log.Printf("WARNING: userno=%d rank=%s nrank=%s irank=%d tirank=%d\n", user.Userno, user.Rank, user.Nrank, user.Irank, tirank)
				user.Irank = tirank
			}
			if uh, ok := usermap[user.Userno]; ok {
				user.Tmrank = uh.Rank
				// 9,999,999形式のランクを整数に変換する
				irank, err := strconv.Atoi(strings.ReplaceAll(uh.Nrank, ",", ""))
				if err != nil {
					log.Printf("strconv.Atoi() returned error %s\n", err.Error())
					continue
				}
				user.Tminrank = irank
				iprank, err := strconv.Atoi(strings.ReplaceAll(uh.Prank, ",", ""))
				if err != nil {
					log.Printf("strconv.Atoi() returned error %s\n", err.Error())
					continue
				}
				user.Tmiprank = iprank
				user.Tmts = uh.Ts
			}
		}

	}
	// テンプレートへのデータの埋め込みを行う
	if err = tpl.ExecuteTemplate(w, "showrank.gtpl", showrank); err != nil {
		err = fmt.Errorf("Handler(): %w", err)
		log.Printf("%s\n", err.Error())
	}

}
