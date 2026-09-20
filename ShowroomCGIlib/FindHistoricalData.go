package ShowroomCGIlib

import (
	"fmt"
	"strings"

	"github.com/Chouette2100/exsrapi/v2"
)

/*
イベントがeventとwebentに共通して存在するかチェックする。
*/
func FindHistoricalData(
	eventinflist *[]exsrapi.Event_Inf,
) (
	err error,
) {
	if len(*eventinflist) == 0 {
		return
	}

	ids := make([]string, 0, len(*eventinflist))
	seen := make(map[string]struct{}, len(*eventinflist))
	for _, eventinf := range *eventinflist {
		if _, ok := seen[eventinf.Event_ID]; ok {
			continue
		}
		seen[eventinf.Event_ID] = struct{}{}
		ids = append(ids, eventinf.Event_ID)
	}

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	sqls := "select distinct e.eventid from event e join wevent we on e.eventid = we.eventid "
	sqls += "where we.eventid in (" + placeholders + ")"

	rows, err := Db0.Query(sqls, args...)
	if err != nil {
		err = fmt.Errorf("Db0.Query(): %w", err)
		return
	}
	defer rows.Close()

	matched := make(map[string]struct{}, len(ids))
	for rows.Next() {
		var eventID string
		err = rows.Scan(&eventID)
		if err != nil {
			err = fmt.Errorf("rows.Scan(): %w", err)
			return
		}
		matched[eventID] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows.Err(): %w", err)
		return
	}

	for i, eventinf := range *eventinflist {
		if _, ok := matched[eventinf.Event_ID]; ok {
			(*eventinflist)[i].Target = 1
		} else {
			(*eventinflist)[i].Target = 0
		}
	}
	return
}
