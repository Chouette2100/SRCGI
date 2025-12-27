package ShowroomCGIlib

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	//	"errors"
	// "sort"
	// "html/template"

	"github.com/Chouette2100/srapi/v2"
	// "github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srdblib/v2"
)

func GetEventRoomlist(
	eventid int,
	eventurlkey string,
	starttime time.Time,
	endtime time.Time,
	ib int,
	ie int,
) (
	proomlist *[]RoomInfo,
	err error,
) {

	client := http.DefaultClient

	eventurlbase := eventurlkey
	blockid := -1
	if strings.Contains(eventurlkey, "?") {
		euka := strings.Split(eventurlkey, "?block_id=")
		eventurlbase = euka[0]
		blockid, err = strconv.Atoi(euka[1])
	}

	// if starttime.Before(time.Now()) && blockid != -1 {
	if blockid != -1 {
		// ブロックイベントでイベント開始後の場合
		var ebr *srapi.EventBlockRanking
		ebr, err = srapi.GetEventBlockRanking(client, eventid, blockid, ib, ie)
		if err != nil {
			err = fmt.Errorf("GetEventBlockRanking(): %w", err)
			log.Printf("%s\n", err.Error())
		} else {
			if ebr != nil && len(ebr.Block_ranking_list) > 0 {
				roomlist := make([]RoomInfo, len(ebr.Block_ranking_list))
				for i, v := range ebr.Block_ranking_list {
					roomlist[i].Userno, _ = strconv.Atoi(v.Room_id)
					roomlist[i].Name = v.Room_name
					roomlist[i].Account = v.Room_url_key
				}
				proomlist = &roomlist
				return
			}
		}
	}

	var er *srapi.EventRooms
	er, err = srapi.GetEventRoomsByApi(client, eventurlbase, ib, ie)
	if err != nil {
		err = fmt.Errorf("GetEventRoomsByApi(): %w", err)
		log.Printf("%s\n", err.Error())
		return nil, err
	}
	if er == nil || len(er.Rooms) == 0 {
		err = fmt.Errorf("参加ルームが取得できません。イベントURLやイベントIDが正しいか確認してください。")
		log.Printf("%s\n", err.Error())
		return nil, err
	}

	// userテーブルから既存のユーザー情報を取得する
	roomlist := make([]RoomInfo, len(er.Rooms))
	for i, v := range er.Rooms {
		roomlist[i].Userno = v.RoomID
		roomlist[i].Name = v.RoomName
		roomlist[i].Account = v.RoomURLKey
	}
	return &roomlist, nil
}
