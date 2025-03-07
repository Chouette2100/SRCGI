package ShowroomCGIlib

import (
	"fmt"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/srdblib/v2"
)
func SelectUsernoAndName(
	Keywordrm	string,
	limit		int,
	offset		int,
) (
	roomlist *[]Room,
	err error,
){

	roomlist = new([]Room)

	kw := "%" + Keywordrm + "%"

	sqlsrl := "select userno, user_name from showroom.user where userno in "
	sqlsrl += " (select userno from showroom.user where user_name like ? "
	sqlsrl += " union select userno from showroom.userhistory where user_name like ? ) "  
	sqlsrl += " order by user_name limit ? offset ?  "

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = srdblib.Db.Prepare(sqlsrl)
	if err != nil {
		err = fmt.Errorf("Prepare(): %w", err)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(kw, kw, limit, offset)
	if err != nil {
		err = fmt.Errorf("Query(): %w", err)
		return
	}
	defer rows.Close()


	for rows.Next() {
		var room Room
		err = rows.Scan(&room.Userno, &room.User_name)
		if err != nil {
			err = fmt.Errorf("Scan(): %w", err)
			return
		}
		room.User_name = "(" + fmt.Sprintf("%d", room.Userno) + ")" + room.User_name
		*roomlist = append(*roomlist, room)
	}

	return
}