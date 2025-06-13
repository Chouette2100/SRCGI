// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bytes"
	"fmt"
	//	"html"
	"log"

	//	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	//	"bufio"
	//	"os"

	//	"runtime"

	//	"encoding/json"

	//	"html/template"
	"net/http"

	//	"database/sql"

	//	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	"github.com/Chouette2100/exsrapi/v2"
	//	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type Color struct {
	Name  string
	Value string
}
type Colormap []Color

// https://www.fukushihoken.metro.tokyo.lg.jp/kiban/machizukuri/kanren/color.files/colorudguideline.pdf
var Colormaplist []Colormap = []Colormap{
	{
		{"#00FFFF", "#00FFFF"},
		{"#FF00FF", "#FF00FF"},
		{"#FFFF00", "#FFFF00"},
		//	-----
		{"#7F7FFF", "#7F7FFF"},
		{"#FF7F7F", "#FF7F7F"},
		{"#7FFF7F", "#7FFF7F"},

		{"#7FBFFF", "#7FBFFF"},
		{"#FF7FBF", "#FF7FBF"},
		{"#BFFF7F", "#BFFF7F"},

		{"#7FFFFF", "#7FFFFF"},
		{"#FF7FFF", "#FF7FFF"},
		{"#FFFF7F", "#FFFF7F"},

		{"#7FFFBF", "#7FFFBF"},
		{"#BF7FFF", "#BF7FFF"},
		{"#FFBF7F", "#FFBF7F"},
		//	-----
		{"#ADADFF", "#ADADFF"},
		{"#FFADAD", "#FFADAD"},
		{"#ADFFAD", "#7FFFAD"},

		{"#ADD6FF", "#ADD6FF"},
		{"#FFADD6", "#FFADD6"},
		{"#D6FFAD", "#D6FFAD"},

		{"#ADFFFF", "#ADFFFF"},
		{"#FFADFF", "#FFADFF"},
		{"#FFFFAD", "#FFFFAD"},

		{"#ADFFD6", "#ADFFD6"},
		{"#D6ADFF", "#D6ADFF"},
		{"#FFD6AD", "#FFD6AD"},
	},
	/*
		{
			{"C0", "#00FFFF"},
			{"M0", "#FF00FF"},
			{"Y0", "#FFFF00"},
			//	-----
			{"C11", "#7F7FFF"},
			{"M11", "#FF7F7F"},
			{"Y11", "#7FFF7F"},

			{"C12", "#7FBFFF"},
			{"M12", "#FF7FBF"},
			{"Y12", "#BFFF7F"},

			{"C13", "#7FFFFF"},
			{"M13", "#FF7FFF"},
			{"Y13", "#FFFF7F"},

			{"C14", "#7FFFBF"},
			{"M14", "#BF7FFF"},
			{"Y14", "#FFBF7F"},
			//	-----
			{"C21", "#ADADFF"},
			{"M21", "#FFADAD"},
			{"Y21", "#7FFFAD"},

			{"C22", "#ADD6FF"},
			{"M22", "#FFADD6"},
			{"Y22", "#D6FFAD"},

			{"C23", "#ADFFFF"},
			{"M23", "#FFADFF"},
			{"Y23", "#FFFFAD"},

			{"C24", "#ADFFD6"},
			{"M24", "#D6ADFF"},
			{"Y24", "#FFD6AD"},
		},
	*/
	{
		{"cyan", "cyan"},
		{"magenta", "magenta"},
		{"yellow", "yellow"},
		{"royalblue", "royalblue"},
		{"coral", "coral"},
		{"khaki", "khaki"},
		{"deepskyblue", "deepskyblue"},
		{"crimson", "crimson"},
		{"orange", "orange"},
		{"lightsteelblue", "lightsteelblue"},
		{"pink", "pink"},
		{"sienna", "sienna"},
		{"springgreen", "springgreen"},
		{"blueviolet", "blueviolet"},
		{"salmon", "salmon"},
		{"lime", "lime"},
		{"red", "red"},
		{"darkorange", "darkorange"},
		{"skyblue", "skyblue"},
		{"lightpink", "lightpink"},
	},
	{
		{"red", "#FF2800"},
		{"yellow", "#FAF500"},
		{"green", "#35A16B"},
		{"blue", "#0041FF"},
		{"skyblue", "#66CCFF"},
		{"lightpink", "#FFD1D1"},
		{"orange", "#FF9900"},
		{"purple", "#9A0079"},
		{"brown", "#663300"},
		{"lightgreen", "#87D7B0"},
		{"white", "#FFFFFF"},
		{"gray", "#77878F"},
	},
}

type Event struct {
	EventID   string
	EventName string
	Period    string
	Starttime time.Time
	S_start   string
	Endtime   time.Time
	S_end     string
	Status    string
	Pntbasis  int
	Modmin    int
	Modsec    int
	Pbname    string
	Selected  string
	Maxpoint  int
	Gscale    int
}

type User struct {
	Userno       int
	Userlongname string
	Selected     string
}

func SelectEventRoomInfList(
	eventid string,
	roominfolist *RoomInfoList,
) (
	eventname string,
	status int,
) {

	status = 0

	//	eventno := 0
	//	eventno, eventname, _ = SelectEventNoAndName(eventid)
	//	Event_inf, _ = SelectEventInf(eventid)
	//	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent("event", eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	eventno := Event_inf.Event_no
	eventname = Event_inf.Event_name

	sql := "select distinct u.userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, e.istarget, e.graph, e.color, e.iscntrbpoints, e.point "
	sql += " from user u join eventuser e "
	sql += " where u.userno = e.userno and e.eventid= ?"
	if Event_inf.Start_time.After(time.Now()) {
		sql += " order by followers desc"
	} else {
		sql += " order by e.point desc"
	}

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("SelectEventRoomInfList() Prepare() err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(eventid)
	if err != nil {
		log.Printf("SelectRoomIn() Query() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	defer rows.Close()

	//	色コードから色名に変換するマップを作る
	//	FIXME: Colormap とは違う、まぎらわしい
	colormap := make(map[string]int)
	cmap := Event_inf.Cmap
	for i := 0; i < len(Colormaplist[cmap]); i++ {
		colormap[Colormaplist[cmap][i].Name] = i
	}

	var roominf RoomInfo

	i := 0
	for rows.Next() {
		err := rows.Scan(&roominf.Userno,
			&roominf.Account,
			&roominf.Name,
			&roominf.Longname,
			&roominf.Shortname,
			&roominf.Genre,
			&roominf.Rank,
			&roominf.Nrank,
			&roominf.Prank,
			&roominf.Level,
			&roominf.Followers,
			&roominf.Fans,
			&roominf.Fans_lst,
			&roominf.Istarget,
			&roominf.Graph,
			&roominf.Color,
			&roominf.Iscntrbpoint,
			&roominf.Point,
		)
		//	FIXME: 色コードでない色名を使えることが問題ではないか？
		//	色名を色コードに変換する
		ci := 0
		for ; ci < len(Colormaplist[cmap]); ci++ {
			if Colormaplist[cmap][ci].Name == roominf.Color {
				roominf.Colorvalue = Colormaplist[cmap][ci].Value
				break
			}
		}
		ii := 0
		if ci == len(Colormaplist[cmap]) {
			var cii int
			for ; ii < len(Colormaplist); ii++ {
				if ii == Event_inf.Cmap {
					continue
				}
				cii = 0
				for ; cii < len(Colormaplist[ii]); cii++ {
					if Colormaplist[ii][cii].Name == roominf.Color {
						roominf.Colorvalue = Colormaplist[ii][cii].Value
						break
					}
				}
				if cii != len(Colormaplist[ii]) {
					break
				}
			}
			if cii == len(Colormaplist[ii]) {
				roominf.Colorvalue = roominf.Color
			}
		}

		if roominf.Istarget == "Y" {
			roominf.Istarget = "Checked"
		} else {
			roominf.Istarget = ""
		}
		if roominf.Graph == "Y" {
			roominf.Graph = "Checked"
		} else {
			roominf.Graph = ""
		}
		if roominf.Iscntrbpoint == "Y" {
			roominf.Iscntrbpoint = "Checked"
		} else {
			roominf.Iscntrbpoint = ""
		}
		roominf.Slevel = humanize.Comma(int64(roominf.Level))
		roominf.Sfollowers = humanize.Comma(int64(roominf.Followers))
		if roominf.Point < 0 {
			roominf.Spoint = ""
		} else {
			roominf.Spoint = humanize.Comma(int64(roominf.Point))
		}
		roominf.Formid = "Form" + fmt.Sprintf("%d", i)
		roominf.Eventid = eventid
		roominf.Name = strings.ReplaceAll(roominf.Name, "'", "’")
		if err != nil {
			log.Printf("SelectEventRoomInfList() Scan() err=%s\n", err.Error())
			status = -7
			return
		}
		//	var colorinf ColorInf
		colorinflist := make([]ColorInf, len(Colormaplist[cmap]))

		for i := 0; i < len(Colormaplist[cmap]); i++ {
			colorinflist[i].Color = Colormaplist[cmap][i].Name
			colorinflist[i].Colorvalue = Colormaplist[cmap][i].Value
		}

		roominf.Colorinflist = colorinflist
		if cidx, ok := colormap[roominf.Color]; ok {
			roominf.Colorinflist[cidx].Selected = "Selected"
		}
		*roominfolist = append(*roominfolist, roominf)

		i++
	}

	if err = rows.Err(); err != nil {
		log.Printf("SelectEventRoomInfList() rows err=%s\n", err.Error())
		status = -8
		return
	}

	if Event_inf.Start_time.After(time.Now()) {
		SortByFollowers = true
	} else {
		SortByFollowers = false
	}
	sort.Sort(*roominfolist)

	return
}

func UpdateEventuserSetPoint(eventid, userid string, point int) (status int) {
	status = 0

	//	eventno, _, _ := SelectEventNoAndName(eventid)
	userno, _ := strconv.Atoi(userid)

	sql := "update eventuser set point=? where eventid = ? and userno = ?"
	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateEventuserSetPoint() error (Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(point, eventid, userno)

	if err != nil {
		log.Printf("error(UpdateEventuserSetPoint() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return
}

func UpdateEventInf(eventinf *exsrapi.Event_Inf) (
	status int,
) {

	if _, _, status = SelectEventNoAndName((*eventinf).Event_ID); status == 0 {
		sql := "Update event set "
		sql += " ieventid=?,"
		sql += " event_name=?,"
		sql += " period=?,"
		sql += " starttime=?,"
		sql += " endtime=?,"
		sql += " noentry=?,"
		sql += " intervalmin=?,"
		sql += " modmin=?,"
		sql += " modsec=?,"
		sql += " Fromorder=?,"
		sql += " Toorder=?,"
		sql += " Resethh=?,"
		sql += " Resetmm=?,"
		sql += " Nobasis=?,"
		sql += " Target=?,"
		sql += " Maxdsp=?, "
		sql += " cmap=?, "
		sql += " maxpoint=? "
		//	sql += " where eventno = ?"
		sql += " where eventid = ?"
		log.Printf("db.Prepare(sql)\n")

		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("UpdateEventInf() error (Update/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		if eventinf.Intervalmin != 5 { //	緊急対応
			log.Printf(" Intervalmin isn't 5. (%dm)\n", eventinf.Intervalmin)
			eventinf.Intervalmin = 5
		}

		log.Printf("row.Exec()\n")
		_, err = stmt.Exec(
			(*eventinf).I_Event_ID,
			(*eventinf).Event_name,
			(*eventinf).Period,
			(*eventinf).Start_time,
			(*eventinf).End_time,
			(*eventinf).NoEntry,
			(*eventinf).Intervalmin,
			(*eventinf).Modmin,
			(*eventinf).Modsec,
			(*eventinf).Fromorder,
			(*eventinf).Toorder,
			(*eventinf).Resethh,
			(*eventinf).Resetmm,
			(*eventinf).Nobasis,
			(*eventinf).Target,
			(*eventinf).Maxdsp,
			(*eventinf).Cmap,
			(*eventinf).Maxpoint+eventinf.Gscale,
			(*eventinf).Event_ID,
		)

		if err != nil {
			log.Printf("error UpdateEventInf() (update/Exec) err=%s\n", err.Error())
			status = -2
		}
	} else {
		status = 1
	}

	return
}

func InsertIntoOrUpdateUser(client *http.Client, tnow time.Time, eventid string, roominf RoomInfo) (status int) {

	status = 0

	//	isnew := false

	userno, _ := strconv.Atoi(roominf.ID)
	log.Printf("  *** InsertIntoOrUpdateUser() *** userno=%d\n", userno)

	nrow := 0
	err := srdblib.Db.QueryRow("select count(*) from user where userno =" + roominf.ID).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	if nrow == 0 {
		// srdblib.InsertIntoUser(client, tnow, userno)
		xu := srdblib.User{Userno: userno}
		srdblib.Env.Waitmsec = 200 // FIXME: 危険
		srdblib.InsertUsertable(client, tnow, &xu)
		srdblib.Env.Waitmsec = 5000
	}

	return

}
func InsertIntoEventUser(i int, eventid string, roominf RoomInfo) (status int) {

	status = 0

	userno, _ := strconv.Atoi(roominf.ID)

	nrow := 0
	sql := "select count(*) from eventuser where userno =? and eventid = ?"
	err := srdblib.Db.QueryRow(sql, roominf.ID, eventid).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	Colorlist := Colormaplist[Event_inf.Cmap]

	if nrow == 0 {
		sql := "INSERT INTO eventuser(eventid, userno, istarget, graph, color, iscntrbpoints, point) VALUES(?,?,?,?,?,?,?)"
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		//	if i < 10 {
		_, err = stmt.Exec(
			eventid,
			userno,
			"Y",
			"Y",
			Colorlist[i%len(Colorlist)].Name,
			"N",
			roominf.Point,
		)

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
		sqlip := "insert into points (ts, user_id, eventid, point, `rank`, gap, pstatus) values(?,?,?,?,?,?,?)"
		_, err = srdblib.Db.Exec(
			sqlip,
			Event_inf.Start_time.Truncate(time.Second),
			userno,
			eventid,
			0,
			1,
			0,
			"=",
		)
		if err != nil {
			err := fmt.Errorf("Db.Exec(sqlip,...): %w", err)
			log.Printf("err=[%s]\n", err.Error())
		}

		status = 1
	}
	return

}

func SelectEventNoAndName(eventid string) (
	eventname string,
	period string,
	status int,
) {

	status = 0

	err := srdblib.Db.QueryRow("select event_name, period from event where eventid ='"+eventid+"'").Scan(&eventname, &period)

	if err == nil {
		return
	} else {
		log.Printf("err=[%s]\n", err.Error())
		if err.Error() != "sql: no rows in result set" {
			status = -2
			return
		}
	}

	status = -1
	return
}

func SelectUserName(userno int) (
	longname string,
	shortname string,
	genre string,
	rank string,
	nrank string,
	prank string,
	level int,
	followers int,
	fans int,
	fans_lst int,
	status int,
) {

	status = 0

	sql := "select longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from user where userno = ?"

	err := srdblib.Db.QueryRow(sql, userno).Scan(&longname, &shortname, &genre, &rank, &nrank, &prank, &level, &followers, &fans, &fans_lst)

	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
	}

	return
}

func SelectUserColor(userno int, eventid string) (
	color string,
	colorvalue string,
	status int,
) {

	Colorlist := Colormaplist[Event_inf.Cmap]

	status = 0

	//	sql := "select color from eventuser where userno = ? and eventno = ?"
	sql := "select color from eventuser where userno = ? and eventid = ?"

	err := srdblib.Db.QueryRow(sql, userno, eventid).Scan(&color)

	i := 0
	for ; i < len(Colorlist); i++ {
		if Colorlist[i].Name == color {
			colorvalue = Colorlist[i].Value
			break
		}
	}
	if i == len(Colorlist) {
		colorvalue = color
	}

	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
	}

	return
}

func SelectUserList() (userlist []User, status int) {

	status = 0

	sql := "select distinct(e.nobasis),u.longname "
	sql += " from event e join user u on e.nobasis=u.userno "
	sql += " where e.nobasis != 0 "
	sql += " order by e.nobasis"

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var user User
	i := 0

	user.Userno = 0
	user.Userlongname = ""
	userlist = append(userlist, user)
	i++

	for rows.Next() {
		err := rows.Scan(&user.Userno, &user.Userlongname)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		userlist = append(userlist, user)
		i++
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return

}

type IdAndRank struct {
	Userno int
	Rank   int
}

func SelectEventInfAndRoomList() (
	idandranklist []IdAndRank,
	status int,
) {

	status = 0

	eventinf, err := srdblib.SelectFromEvent("event", Event_inf.Event_ID)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	log.Printf("eventno=%d\n", Event_inf.Event_no)

	start_date := Event_inf.Start_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.Start_time.Hour()) * time.Hour)
	end_date := Event_inf.End_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

	//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

	Event_inf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
	Event_inf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - Event_inf.Start_date

	//	log.Printf("Start_data=%f Dperiod=%f\n", Event_inf.Start_date, Event_inf.Dperiod)

	sql := "select max(point) from eventuser where eventid = ? and graph = 'Y'"
	err = srdblib.Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.MaxPoint)
	//	err = srdblib.Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.Maxpoint)

	if err != nil {
		log.Printf("select max(point) from eventuser where eventid = '%s'  and graph = 'Y'\n", Event_inf.Event_ID)
		log.Printf("err=[%s]\n", err.Error())
		status = -2
		return
	}

	//	log.Printf("MaxPoint=%d\n", Event_inf.MaxPoint)

	sqlst := "select p.user_id, p.`rank` from points p join eventuser eu on p.user_id = eu.userno and p.eventid = eu.eventid "
	sqlst += " where p.eventid = ? and eu.graph = 'Y' "
	sqlst += " and p.ts = ( select max(ts) from points where eventid = ? ) order by p.point desc "
	stmt, err := srdblib.Db.Prepare(sqlst)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(Event_inf.Event_ID, Event_inf.Event_ID)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	i := 0
	userno := 0
	rank := 0
	for rows.Next() {
		err := rows.Scan(&userno, &rank)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		idandranklist = append(idandranklist, IdAndRank{Userno: userno, Rank: rank})
		i++
		if i == Event_inf.Maxdsp {
			break
		}
	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return
}

func SelectPointList(userno int, eventid string) (norow int, tp *[]time.Time, pp *[]int) {

	norow = 0

	//	log.Printf("SelectPointList() userno=%d eventid=%s\n", userno, eventid)
	stmt1, err := srdblib.Db.Prepare("SELECT count(*) FROM points where user_id = ? and eventid = ?")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt1.Close()

	//	var norow int
	err = stmt1.QueryRow(userno, eventid).Scan(&norow)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	//	fmt.Println(norow)

	//	----------------------------------------------------

	//	stmt1, err = Db.Prepare("SELECT max(t.t) FROM timeacq t join points p where t.idx=p.idx and user_id = ? and event_id = ?")
	stmt1, err = srdblib.Db.Prepare("SELECT max(ts) FROM points where user_id = ? and eventid = ?")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt1.Close()

	var tfinal time.Time
	err = stmt1.QueryRow(userno, eventid).Scan(&tfinal)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	islastdata := false
	if tfinal.After(Event_inf.End_time.Add(time.Duration(-Event_inf.Intervalmin) * time.Minute)) {
		islastdata = true
	}
	//	fmt.Println(norow)

	//	----------------------------------------------------

	t := make([]time.Time, norow)
	point := make([]int, norow)
	if islastdata {
		t = make([]time.Time, norow+1)
		point = make([]int, norow+1)
	}

	tp = &t
	pp = &point

	if norow == 0 {
		return
	}

	//	----------------------------------------------------

	//	stmt2, err := Db.Prepare("select t.t, p.point from points p join timeacq t on t.idx = p.idx where user_id = ? and event_id = ? order by t.t")
	stmt2, err := srdblib.Db.Prepare("select ts, point from points where user_id = ? and eventid = ? order by ts")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt2.Close()

	rows, err := stmt2.Query(userno, eventid)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(&t[i], &point[i])
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			//	status = -1
			return
		}
		i++

	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}

	if islastdata {
		t[norow] = t[norow-1].Add(15 * time.Minute)
		point[norow] = point[norow-1]
	}

	tp = &t
	pp = &point

	return
}
