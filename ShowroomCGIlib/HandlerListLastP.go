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
	//	"sort"
	"strconv"
	"strings"
	"time"

	//	"bufio"
	//	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"

	//	"database/sql"

	//	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi/v2"
	//	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func ListLastPHandler(w http.ResponseWriter, req *http.Request) {

	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	status := 0

	var list_last struct {
		Eventid   string
		Ieventid  string
		Userno    string
		Detail    string
		Isover    string
		Page      int
		Maxrooms  int
		NoRooms   int
		Roomid    int
		Scorelist []CurrentScore
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrb-h1.gtpl","templates/list-cntrb-h2.gtpl","templates/list-cntrb.gtpl"))
	funcMap := template.FuncMap{
		//	3桁ごとに","を挿入する
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
		//	イベントIDがブロックIDを含む場合はそれを取り除く。
		"DelBlockID": func(eid string) string {
			eia := strings.Split(eid, "?")
			if len(eia) == 2 {
				return eia[0]
			} else {
				return eid
			}
		},
		"Add": func(a1, a2 int) int { return a1 + a2 },
		// 整数除算用の関数を追加
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			// return (a + b - 1) / b
			return a / b
		},
	}
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-lastP.gtpl", "templates/list-lastP_h.gtpl"))

	eventid := req.FormValue("eventid")
	list_last.Eventid = eventid
	userno := req.FormValue("userno")
	list_last.Userno = userno
	list_last.Roomid, _ = strconv.Atoi(req.FormValue("roomid"))
	list_last.Detail = req.FormValue("detail")

	// ページ番号を取得（デフォルトは1）
	pageStr := req.FormValue("page")
	if pageStr == "" {
		list_last.Page = 1
	} else {
		list_last.Page, _ = strconv.Atoi(pageStr)
		if list_last.Page < 1 {
			list_last.Page = 1
		}
	}

	// 1ページあたりの表示数を20件に設定
	list_last.Maxrooms = 20

	// データ取得のオフセットを計算
	offset := (list_last.Page - 1) * list_last.Maxrooms

	log.Printf("      eventid=%s, detail=%s, page=%d\n", eventid, list_last.Detail, list_last.Page)
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
	// Event_inf = *eventinf
	list_last.Ieventid = strconv.Itoa(eventinf.I_Event_ID)

	tdata, eventname, period, scorelist, totalCount, status := SelectCurrentScoreP(eventid, list_last.Maxrooms, offset)
	list_last.NoRooms = totalCount
	list_last.Scorelist = scorelist

	// 結果が0件の場合、前のページに戻る（ただしページ1より小さくならないようにする）
	if len(scorelist) == 0 && list_last.Page > 1 {
		list_last.Page--
		offset = (list_last.Page - 1) * list_last.Maxrooms
		// tdata, eventname, period, scorelist, totalCount, status = SelectCurrentScoreP(eventid, list_last.Maxrooms, offset)
		tdata, eventname, period, scorelist, totalCount, status = SelectCurrentScoreP(eventid, list_last.Maxrooms, offset)
		list_last.Scorelist = scorelist
		list_last.NoRooms = totalCount
	}

	for i := 0; i < len(scorelist); i++ {
		switch scorelist[i].Roomgenre {
		case "Voice Actors & Anime":
			scorelist[i].Roomgenre = "VA&A"
		case "Talent Model":
			scorelist[i].Roomgenre = "Tl/Md"
		case "Comedians/Talk Show":
			scorelist[i].Roomgenre = "Cm/TS"
		default:
		}
	}

	//	tnext := tdata.Add(5 * time.Minute)
	tnext := tdata.Add(time.Duration(eventinf.Intervalmin) * time.Minute) //	0101G5
	//	treload := tnext.Add(5 * time.Second)
	treload := tnext.Add(10 * time.Second)

	values := map[string]string{
		"Eventid":         eventid,
		"Ieventid":        list_last.Ieventid,
		"userno":          userno,
		"UpdateTime":      "データ取得時刻：　" + tdata.Format("2006/01/02 15:04:05"),
		"NextTime":        "次のデータ取得は　" + tnext.Format("15:04:05") + "　に予定されています。",
		"ReloadTime":      "画面のリロードが　" + treload.Format("15:04:05") + "　頃に行われます。",
		"SecondsToReload": fmt.Sprintf("%d", int(time.Until(treload).Seconds()+5)),
		"EventName":       eventname,
		"Period":          period,
		"Detail":          list_last.Detail,
		"Page":            fmt.Sprintf("%d", list_last.Page),
		"Maxpoint":        fmt.Sprintf("%d", eventinf.Maxpoint),
		"Gscale":          fmt.Sprintf("%d", eventinf.Gscale),
	}

	// if time.Since(tdata) > 5*time.Minute {
	// 	log.Printf("Application stopped or the event is over. status = %d\n", status)
	// 	values["NextTime"] = "表示されているデータは最新ではありません。"
	// 	values["ReloadTime"] = "もうしわけありませんがデータ取得が復旧するまでしばらくお待ちください。"
	// 	values["SecondsToReload"] = "300"
	// }
	if status != 0 {
		log.Printf("GetCurrentScore() returned %d.\n", status)
		values["UpdateTime"] = "データが取得できませんでした。"
		values["NextTime"] = "もうしわけありませんがしばらくお待ち下さい。"
		values["ReloadTime"] = ""
		values["SecondsToReload"] = "300"
	}
	if time.Now().After(eventinf.End_time) {
		// log.Printf("Application stopped or the event is over. status = %d\n", status)
		values["NextTime"] = "イベントは終了しています。"
		values["ReloadTime"] = ""
		values["SecondsToReload"] = "3600"

		list_last.Isover = "1"
	}
	if time.Now().Before(eventinf.Start_time) {
		values["NextTime"] = "イベントはまだ始まっていません。"
		values["ReloadTime"] = ""
	}
	//	log.Printf("Values=%v", values)
	if err := tpl.ExecuteTemplate(w, "list-last_h", values); err != nil {
		log.Println(err)
	}
	if status != 0 {
		fmt.Fprintf(w, "</body>\n</html>\n")
		return
	}
	if err := tpl.ExecuteTemplate(w, "list-lastP", list_last); err != nil {
		log.Println(err)
	}
}
func SelectCurrentScoreP(
	eventid string,
	limit int,
	offset int,
) (
	gtime time.Time,
	eventname string,
	period string,
	scorelist []CurrentScore,
	totalCount int,
	status int,
) {

	status = 0

	//	Event_inf, status = SelectEventInf(eventid)
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
	// Event_inf = *eventinf

	//	eventno = Event_inf.Event_no
	eventname = eventinf.Event_name
	period = eventinf.Period

	// 総件数を取得
	totalCount = 0
	// sql0 := "select count(*) from points where eventid = ?"
	sql0 := "select count(*) from points where eventid = ?"
	// sql0 += " and ts = (select max(ts) from points where eventid = ? and ts > ?)" ? // TODO: 最後の条件は？
	sql0 += " and ts = (select max(ts) from points where eventid = ? ) "
	err = srdblib.Db.QueryRow(sql0, eventid, eventid).Scan(&totalCount)

	if err != nil {
		log.Printf("select max(point) from eventuser where eventid = '%s'\n", eventinf.Event_ID)
		log.Printf("err=[%s]\n", err.Error())
		status = -11
		return
	}
	if totalCount == 0 {
		log.Printf("no data in points(where eventid=%s).\n", eventid)
		status = -12
		return
	}

	//	---------------------------------------------------
	//	sql := "select t.idx, t.t from timeacq t join points p where t.idx = p.idx and t.idx = ( select max(idx) from points where event_id = ? )"
	//	sql := "select distinct t.idx, t.t from timeacq t join points p where t.idx = p.idx and t.t = ( select max(t) from points p join timeacq t where p.idx = t.idx and event_id = ? )"
	sql1 := "select distinct max(ts) from points where eventid = ?"
	//	sql := "select distinct COALESCE(max(ts), ?) from points where eventid = ?"
	stmt1, err := srdblib.Db.Prepare(sql1)
	if err != nil {
		log.Printf("GetCurrentScore() (3) err=%s\n", err.Error())
		status = -3
		return
	}
	//	defer stmt1.Close()
	defer func() {
		err := stmt1.Close()
		if err != nil {
			log.Printf("stmt1.Close() err=%s\n", err.Error())
		}
	}()

	//	idx := 0
	//	Err = stmt.QueryRow(time.Now().Add(time.Hour), eventid).Scan(&gtime)
	err = stmt1.QueryRow(eventid).Scan(&gtime)
	if err != nil {
		log.Printf("GetCurrentScore() (4) err=%s\n", err.Error())
		status = -4
		return
	}
	log.Printf("gtime=%s\n", gtime.Format("2006/01/02 15:04:06"))
	/*
		if gtime.After(time.Now()) {
			status = -10
			return
		}
	*/

	//	---------------------------------------------------
	//	stmt, err = Db.Prepare("select user_id, `rank`, point, pstatus, ptime, qstatus, qtime from points where eventid = ? and ts = ? order by point desc")
	sql2 := "SELECT p.user_id, u.userid, p.rank, p.point, p.pstatus, p.ptime, p.qstatus, p.qtime "
	sql2 += " FROM points p JOIN user u where p.eventid = ? AND p.user_id = u.userno "
	sql2 += " AND (p.user_id , p.ts) IN (SELECT user_id, MAX(ts) FROM points WHERE eventid = ? AND ts > ? GROUP BY user_id) "

	// HACK: 本来はランキングイベントか否かで処理をわけるべきところだし、別のソートキー（API結果での出現順？）を用意すべき。
	if strings.Contains(eventid, "mattari_fireworks") {
		sql2 += " ORDER BY p.`rank` desc, p.point desc "
	} else {
		sql2 += " ORDER BY p.point desc "
	}
	// HACK: --------------------
	if limit > 0 {
		sql2 += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit)
	}

	stmt2, err := srdblib.Db.Prepare(sql2)

	if err != nil {
		log.Printf("GetCurrentScore() (5) err=%s\n", err.Error())
		status = -5
		return
	}
	//	defer stmt2.Close()
	defer func() {
		err := stmt2.Close()
		if err != nil {
			log.Printf("stmt2.Close() err=%s\n", err.Error())
		}
	}()

	//	rows, err := stmt.Query(eventid, gtime)
	rows2, err := stmt2.Query(eventid, eventid, gtime.Add(-2*time.Minute))
	if err != nil {
		log.Printf("GetCurrentScore() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	defer rows2.Close()

	//	var score, bscore CurrentScore
	// var bscore CurrentScore
	point_bs := 0
	i := 0
	//	shift := 1
	// nextrank := 1
	for rows2.Next() {
		var score CurrentScore
		err := rows2.Scan(&score.Userno, &score.Shorturl, &score.Rank, &score.Point, &score.Pstatus, &score.Ptime, &score.Qstatus, &score.Qtime)
		if err != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", err.Error())
			status = -7
			return
		}
		if score.Userno == eventinf.Nobasis {
			point_bs = score.Point
			log.Printf(" Nobasis=%d  point_bs=%d\n", eventinf.Nobasis, point_bs)
		}
		score.Spoint = humanize.Comma(int64(score.Point))
		username, _, roomgenre, roomrank, roomnrank, roomprank, roomlevel, followers, fans, fans_lst, sts := SelectUserName(score.Userno)
		score.Username = username
		if sts != 0 {
			score.Username = fmt.Sprintf("%d", score.Userno)
		}
		score.Roomgenre = roomgenre
		score.Roomrank = roomrank
		score.Roomnrank = roomnrank
		score.Roomprank = roomprank
		score.Roomlevel = humanize.Comma(int64(roomlevel))
		score.Followers = humanize.Comma(int64(followers))
		score.Fans = fans
		score.Fans_lst = fans_lst

		/*
			nroomlevel := 0
			nfollowers := 0
			score.Roomgenre, score.Roomrank, score.Roomnrank, score.Roomprank, nroomlevel,
				nfollowers, score.Fans, score.Fans_lst, _, _, _, status = GetRoomInfoByAPI(fmt.Sprintf("%d", score.Userno))
			score.Roomlevel = humanize.Comma(int64(nroomlevel))
			score.Followers = humanize.Comma(int64(nfollowers))
			/* */
		/*
			if	score.Roomrank != roomrank ||
				score.Roomnrank != roomnrank ||
				nfollowers != followers ||
				nroomlevel != roomlevel ||
				score.Fans != fans {
				UpdateRoomRankInf (score, nroomlevel, nfollowers)

			}
			/* */

		if score.Rank != 0 {
			score.Srank = fmt.Sprintf("%d", score.Rank)
		} else {
			score.Srank = ""
		}
		// //	if score.Rank > i+shift {
		// 順番が飛んだ飛んだとき空白行を入れる処理はしない。
		// if score.Rank > nextrank {
		// 	//	bscore.Srank = fmt.Sprintf("%d", i+shift)
		// 	bscore.Srank = "-"
		// 	scorelist = append(scorelist, bscore)
		// 	//	shift++
		// }
		// nextrank = score.Rank + 1

		//	score.NextLive, _ = GetNextliveByAPI(fmt.Sprintf("%d", score.Userno))
		score.NextLive, _ = GetNextliveByAPI(strconv.Itoa(score.Userno))
		score.Eventid = eventid

		acqtimelist, _ := SelectAcqTimeList(eventid, score.Userno)
		score.Ncntrb = len(acqtimelist)
		//	log.Printf(" eventid = %s userno = %d len(acqtimelist=%d\n", eventid, score.Userno, lenatl)

		_, perslotinflist, _ := MakePointPerSlot(eventid, score.Userno)
		score.Nperslot = len(perslotinflist)

		scorelist = append(scorelist, score)
		i++
		/*
			if i == 10 {
				break
			}
		*/
	}
	if err = rows2.Err(); err != nil {
		log.Printf("GetCurrentScore() (8) err=%s\n", err.Error())
		status = -8
		return
	}

	if point_bs > 0 {
		for i, score := range scorelist {
			if score.Point != 0 {
				scorelist[i].Sdfr = humanize.Comma(int64(score.Point - point_bs))
			}
		}
	}

	return

}
