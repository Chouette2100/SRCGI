package ShowroomCGIlib

import (
	"fmt"
	"log"
	"time"

	//	"math"
	//	"sort"
	"strconv"
	//	"strings"
	//	"time"

	//	"bufio"
	//	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//	"github.com/PuerkitoBio/goquery"
	//	svg "github.com/ajstarks/svgo/float"
	//	"github.com/dustin/go-humanize"

	"github.com/Chouette2100/srdblib/v2"
)

/*
        HandlerListCntrbH()
		指定されたイベント、配信者、リスナーの貢献ポイントの履歴を表示する。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。
	0101J1	

*/

func HandlerFanLevel(w http.ResponseWriter, req *http.Request) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if ! isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}



	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/fanlevel-user.gtpl",
		"templates/fanlevel-room.gtpl",
		"templates/fanlevel-lfu.gtpl",
		"templates/fanlevel-lfr.gtpl",
	))

	userid := 0
	value := req.FormValue("userid")
	if value == "" {
		userid = 0
	} else {
		userid, _ = strconv.Atoi(value)
	}

	roomid := 0
	value = req.FormValue("roomid")
	if value == "" {
		roomid = 0
	} else {
		roomid, _ = strconv.Atoi(value)
	}

	log.Printf("/fanlevel userid=%d, roomid=%d\n", userid, roomid)

	tnow := time.Now()
	yyyy := tnow.Year()
	mm := int(tnow.Month())

	if userid == 0 && roomid == 0 {

		userinflist, status := SelectFromFluser(0)
		if status != 0 {
			log.Printf("SelectFromFluser() returned %d\n", status)
			return
		}

		if err := tpl.ExecuteTemplate(w, "fanlevel-user.gtpl", userinflist); err != nil {
			log.Println(err)
		}

		roominflist, status := SelectFromFlroom(0)
		if status != 0 {
			log.Printf("SelectFromFlrom() returned %d\n", status)
			return
		}

		if err := tpl.ExecuteTemplate(w, "fanlevel-room.gtpl", roominflist); err != nil {
			log.Println(err)
		}
	} else if userid != 0 {
		lfu, status := SelectLevelForUser(userid, yyyy, mm)
		if status != 0 {
			log.Printf("SelectFromFlrom() returned %d\n", status)
			return
		}

		if err := tpl.ExecuteTemplate(w, "fanlevel-lfu.gtpl", lfu); err != nil {
			log.Println(err)
		}

	} else {
		lfrw, status := SelectLevelForRoom(roomid, yyyy, mm)
		if status != 0 {
			log.Printf("SelectFromFlrom() returned %d\n", status)
			return
		}

		if err := tpl.ExecuteTemplate(w, "fanlevel-lfr.gtpl", lfrw); err != nil {
			log.Println(err)
		}

	}

}

type UserInf struct {
	User_id   int
	User_name string
	Level     int
}

func SelectFromFluser(userid int) (userinflist []UserInf, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	userinflist = make([]UserInf, 0)

	if userid == 0 {
		stmt, srdblib.Dberr = srdblib.Db.Prepare("select user_id, user_name from fluser")
	} else {
		stmt, srdblib.Dberr = srdblib.Db.Prepare("select user_id, user_name from fluser where user_id = ?")
	}
	if srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt.Close()

	if userid == 0 {
		rows, srdblib.Dberr = stmt.Query()
	} else {
		rows, srdblib.Dberr = stmt.Query(userid)
	}
	if srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -2
		return
	}
	defer rows.Close()

	var userinf UserInf
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&userinf.User_id, &userinf.User_name)
		if srdblib.Dberr != nil {
			log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
			status = -3
			return
		}
		userinflist = append(userinflist, userinf)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -4
		return
	}

	return

}

type RoomInf struct {
	Room_id   int
	Room_name string
}

func SelectFromFlroom(roomid int) (roominflist []RoomInf, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	roominflist = make([]RoomInf, 0)

	if roomid == 0 {
		stmt, srdblib.Dberr = srdblib.Db.Prepare("select room_id, room_name from flroom")
	} else {
		stmt, srdblib.Dberr = srdblib.Db.Prepare("select room_id, room_name from flroom where room_id = ?")
	}
	if srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt.Close()

	if roomid == 0 {
		rows, srdblib.Dberr = stmt.Query()
	} else {
		rows, srdblib.Dberr = stmt.Query(roomid)
	}
	if srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -2
		return
	}
	defer rows.Close()

	var roominf RoomInf
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&roominf.Room_id, &roominf.Room_name)
		if srdblib.Dberr != nil {
			log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
			status = -3
			return
		}
		roominflist = append(roominflist, roominf)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -4
		return
	}

	return

}

type LevelForUser struct {
	Room_id   int
	Room_name string
	Level     int
	Level_lst int
}

type LevelForUserW struct {
	Userid    int
	Username  string
	Levellist []LevelForUser
}

func SelectLevelForUser(userid int, yyyy int, mm int) (lfuw LevelForUserW, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	var userinf []UserInf
	userinf, status = SelectFromFluser(userid)
	lfuw.Userid = userid
	lfuw.Username = userinf[0].User_name

	lfuw.Levellist = make([]LevelForUser, 0)

	/*
		sql := "select l.room_id, r.room_name, l.level, l.level_lst from fanlevel l join flroom r "
		sql += " where l.room_id = r.room_id and l.user_id = ? "
		sql += " and yyyymm = ? "
		sql += " order by l.level desc, l.level_lst desc "
	*/
	/*
		if time.Now().Day() > 3 {
			sql += " order by l.level desc "
		} else {
			sql += " order by l.level_lst desc "
		}
	*/

	sqlstmt := "select c.room_id, room_name, c.level, p.level_lst from fanlevel c "
	sqlstmt += " left join fanlevel p on c.user_id = p.user_id and c.room_id = p.room_id and c.yyyymm = p.yyyymm + ? "
	sqlstmt += " join flroom r "
	sqlstmt += " where c.user_id = ? and c.yyyymm = ? and c.room_id = r.room_id order by c.level desc, p.level_lst desc"

	stmt, srdblib.Dberr = srdblib.Db.Prepare(sqlstmt)
	if srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(1, userid, yyyy*100+mm)
	if srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -2
		return
	}
	defer rows.Close()

	/*
	type NullInt64 struct {
		Int64 int64
		Valid bool // Valid is true if Int64 is not NULL
	}
	*/
	var nulllevellst sql.NullInt64
	var level LevelForUser
	for rows.Next() {
		//	Err = rows.Scan(&level.Room_id, &level.Room_name, &level.Level, &level.Level_lst)
		srdblib.Dberr = rows.Scan(&level.Room_id, &level.Room_name, &level.Level, &nulllevellst)
		if srdblib.Dberr != nil {
			log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
			status = -3
			return
		}
		if nulllevellst.Valid {
			level.Level_lst = int(nulllevellst.Int64)
		} else {
			level.Level_lst = -1
		}
		lfuw.Levellist = append(lfuw.Levellist, level)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("** SelectFromFlroom() err=[%s]\n", srdblib.Dberr.Error())
		status = -4
		return
	}

	return

}

type LevelForRoom struct {
	User_id   int
	User_name string
	Level     int
	Level_lst int
}

type LevelForRoomW struct {
	Roomid   int
	Roomname string
	Lfr      []LevelForRoom
}

func SelectLevelForRoom(roomid int, yyyy int, mm int) (lfrw LevelForRoomW, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	//	ルーム名を知るためにルーム情報を取得する。
	roominflist, status := SelectFromFlroom(roomid)

	lfrw.Roomid = roomid
	lfrw.Roomname = roominflist[0].Room_name

	lfrw.Lfr = make([]LevelForRoom, 0)

	/*
	sql := "select l.user_id, u.user_name, l.level, l.level_lst from fanlevel l join fluser u "
	sql += " where l.user_id = u.user_id and l.room_id = ? "
	sql += " and yyyymm = ? "
	sql += " order by l.level desc, l.level_lst desc "
	*/
	/*
		if time.Now().Day() > 3 {
			sql += " order by l.level desc"
		} else {
			sql += " order by l.level_lst desc"
		}
	*/
	sqlstmt := "select c.user_id, u.user_name, c.level, p.level_lst from fanlevel c "
	sqlstmt += " left join fanlevel p on c.user_id = p.user_id and c.room_id = p.room_id and c.yyyymm = p.yyyymm + ? "
	sqlstmt += " join fluser u "
	sqlstmt += " where c.room_id = ? and c.yyyymm = ? and c.user_id = u.user_id order by c.level desc, p.level_lst desc"

	stmt, srdblib.Dberr = srdblib.Db.Prepare(sqlstmt)
	if srdblib.Dberr != nil {
		log.Printf("** SelectLevelForRoom() err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(1, roomid, yyyy*100+mm)
	if srdblib.Dberr != nil {
		log.Printf("** SelectLevelForRoom() err=[%s]\n", srdblib.Dberr.Error())
		status = -2
		return
	}
	defer rows.Close()

	var nulllevellst sql.NullInt64
	var level LevelForRoom
	for rows.Next() {
		//	Err = rows.Scan(&level.User_id, &level.User_name, &level.Level, &level.Level_lst)
		srdblib.Dberr = rows.Scan(&level.User_id, &level.User_name, &level.Level, &nulllevellst)
		if srdblib.Dberr != nil {
			log.Printf("** SelectLevelForRoom() err=[%s]\n", srdblib.Dberr.Error())
			status = -3
			return
		}
		if nulllevellst.Valid {
			level.Level_lst = int(nulllevellst.Int64)
		} else {
			level.Level_lst = -1
		}
		lfrw.Lfr = append(lfrw.Lfr, level)
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("** SelectLevelForRoom() err=[%s]\n", srdblib.Dberr.Error())
		status = -4
		return
	}

	return

}
