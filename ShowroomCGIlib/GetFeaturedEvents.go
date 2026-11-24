package ShowroomCGIlib

import (
	"fmt"
	"log"

	"github.com/go-gorp/gorp"
)

func GetFeaturedEvents(
	dbmap *gorp.DbMap,
	mode string,
	num int,
	lmct int,
) (
	eventmap map[string]int,
	err error,
) {
	eventmap = make(map[string]int)

	// sqlst := "select eventid, count(*) ct from accesslog where ts > SUBDATE(now(),INTERVAL ? hour) AND is_bot = 0 "
	// sqlst += " group by eventid order by ct desc limit ? "

	sqlst := "select eventid, ct from featuredevent "
	switch mode {
	case "scheduled":
		sqlst += " where starttime > Now() "
	case "current":
		sqlst += " where Now() between starttime and endtime "
	case "closed":
		sqlst += " where endtime <= Now() "
	default:
		err = fmt.Errorf("GetFeaturedEvents: invalid mode %s", mode)
		log.Printf("Error: %v", err)
		return
	}
	sqlst += " order by ct desc limit ? "

	type event struct {
		Eventid string
		Ct      int
	}
	var events []event

	_, err = dbmap.Select(&events, sqlst, num)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	for _, v := range events {
		if v.Ct < lmct {
			break
		}
		eventmap[v.Eventid] = v.Ct
	}
	return

}
