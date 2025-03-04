package ShowroomCGIlib
import (
	"fmt"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)
/*
	イベントがeventとwebentに共通して存在するかチェックする。
*/
func FindHistoricalData(
	eventinflist *[]exsrapi.Event_Inf,
) (
	err error,
) {

	for i, eventinf := range *eventinflist {
		sqls := "select count(*) from event e join wevent we on e.eventid = we.eventid where we.eventid = ?"
		err = srdblib.Db.QueryRow(sqls, eventinf.Event_ID).Scan(&(*eventinflist)[i].Target)
		if err != nil {
			err = fmt.Errorf("QueryRow().Scan(): %w", err)
			return
		}
	}
	return
}