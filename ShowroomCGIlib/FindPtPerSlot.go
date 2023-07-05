package ShowroomCGIlib

import (
	"fmt"

	"database/sql"

	"github.com/Chouette2100/srapi"
	"github.com/Chouette2100/srdblib"
)
func FindPtPerSlot(eventid string,roomlist *[]srapi.Room )(
	err error,
) {

	var stmts *sql.Stmt

sqls := "select count(*) from timetable where eventid = ? and userid = ?"
stmts, err = srdblib.Db.Prepare(sqls)
if err != nil {
	err = fmt.Errorf("Prepare(): %w", err)
	return
}
defer stmts.Close()

for i, room := range *roomlist {

	nrow := 0
	err = stmts.QueryRow(eventid, room.Room_id).Scan(&nrow)
	if err != nil {
		err = fmt.Errorf("QueryRow().Scan(): %w", err)
		return
	}

	if nrow != 0 {
		//	timetableに1レコードでもあれば枠別貢献ポイントが記録されている。
		(*roomlist)[i].Isofficial = true	// Ifofficialは目的外使用 w
	} else {
		(*roomlist)[i].Isofficial = false
	}
}
return
}
