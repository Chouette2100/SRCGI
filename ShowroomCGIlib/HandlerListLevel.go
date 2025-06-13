// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"fmt"
	"log"
	"strconv"

	//	"strings"
	"time"

	//	"github.com/dustin/go-humanize"

	"html/template"
	"net/http"

	"database/sql"

	//	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type RoomLevel struct {
	User_name string
	Genre     string
	Rank      string
	Nrank     string
	Prank     string
	Level     int
	Followers int
	Fans      int
	Fans_lst  int
	ts        time.Time
	Sts       string
}

type RoomLevelInf struct {
	Userno        int
	User_name     string
	RoomLevelList []RoomLevel
}

func ListLevelHandler(w http.ResponseWriter, req *http.Request) {

	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/list-level.gtpl"))

	userno, _ := strconv.Atoi(req.FormValue("userno"))
	levelonly, _ := strconv.Atoi(req.FormValue("levelonly"))
	log.Printf("  *** HandlerListLevel() called. userno=%d, levelonly=%d\n", userno, levelonly)

	RoomLevelInf, _ := SelectRoomLevel(userno, levelonly)

	if err := tpl.ExecuteTemplate(w, "list-level.gtpl", RoomLevelInf); err != nil {
		log.Println(err)
	}
}
func SelectRoomLevel(userno int, levelonly int) (roomlevelinf RoomLevelInf, status int) {

	var err error
	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	sqlstmt := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts from userhistory where userno = ? order by ts desc"
	stmt, err = srdblib.Db.Prepare(sqlstmt)
	if err != nil {
		log.Printf("SelectRoomLevel() (3) err=%s\n", err.Error())
		status = -3
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(userno)
	if err != nil {
		log.Printf("SelectRoomLevel() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	defer rows.Close()

	var roomlevel RoomLevel

	roomlevelinf.Userno = userno

	lastlevel := 0

	for rows.Next() {
		err = rows.Scan(&roomlevel.User_name, &roomlevel.Genre, &roomlevel.Rank,
			&roomlevel.Nrank,
			&roomlevel.Prank,
			&roomlevel.Level,
			&roomlevel.Followers,
			&roomlevel.Fans,
			&roomlevel.Fans_lst,
			&roomlevel.ts)
		if err != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", err.Error())
			status = -7
			return
		}

		if lastlevel == 0 {
			roomlevelinf.User_name = roomlevel.User_name
		}

		if levelonly == 1 && roomlevel.Level == lastlevel {
			continue
		}
		lastlevel = roomlevel.Level

		//	roomlevel.Sfollowers = humanize.Comma(int64(roomlevel.Followers))
		roomlevel.Sts = roomlevel.ts.Format("2006/01/02 15:04")

		roomlevelinf.RoomLevelList = append(roomlevelinf.RoomLevelList, roomlevel)

	}

	return
}
