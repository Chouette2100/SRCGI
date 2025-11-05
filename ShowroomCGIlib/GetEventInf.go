package ShowroomCGIlib

import (
	"log"
	"time"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func GetEventInf(
	eventid string,
	eventinfo *exsrapi.Event_Inf,
) (
	status int,
) {
	var intf interface{}
	var err error

	status = 0

	//	イベント情報を取得します。
	intf, err = srdblib.Dbmap.Get(&srdblib.Wevent{}, eventid)
	if err != nil || intf == nil {
		//	イベント情報が取得できなかった場合
		if err != nil {
			log.Printf("GetEventInf() srdblib.Dbmap.Get() returned %s\n", err.Error())
		} else {
			log.Printf("GetEventInf() srdblib.Dbmap.Get() returned intf is nil\n")
		}
		status = -1
		return
	}
	event := intf.(*srdblib.Wevent)

	//	イベント情報をコピーします。
	eventinfo.Event_ID = eventid
	eventinfo.I_Event_ID = event.Ieventid
	eventinfo.Event_name = event.Event_name
	eventinfo.Start_time = event.Starttime
	eventinfo.End_time = event.Endtime
	eventinfo.Period = event.Period
	if event.Endtime.Before(time.Now()) {
		eventinfo.EventStatus = "Over"
	} else if event.Starttime.After(time.Now()) {
		eventinfo.EventStatus = "NotHeldYet"
	} else {
		eventinfo.EventStatus = "BeingHeld"
	}

	return
}

/*
func GetEventInf(
	eventid string,
	eventinfo *exsrapi.Event_Inf,
) (
	status int,
) {

	//	画面からのデータ取得部分は次を参考にしました。
	//		はじめてのGo言語：Golangでスクレイピングをしてみた
	//		https://qiita.com/ryo_naka/items/a08d70f003fac7fb0808

	//	_url := "https://www.showroom-live.com/event/" + EventID
	//	_url = "file:///C:/Users/kohei47/Go/src/EventRoomList03/20210128-1143.html"
	//	_url = "file:20210128-1143.html"

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid

	status = 0

	if inputmode == "file" {

		//	ファイルからドキュメントを作成します
		f, e := os.Open(eventidorfilename)
		if e != nil {
			//	log.Fatal(e)
			log.Printf("err=[%s]\n", e.Error())
			status = -1
			return
		}
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			status = -4
			return
		}

		content, _ := doc.Find("head > meta:nth-child(6)").Attr("content")
		content_div := strings.Split(content, "/")
		(*eventinfo).Event_ID = content_div[len(content_div)-1]

	} else {
		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventidorfilename
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
			status = 1
			return
		}
		defer resp.Body.Close()

		doc, error = goquery.NewDocumentFromReader(resp.Body)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			status = 1
			return
		}

		(*eventinfo).Event_ID = eventidorfilename
	}
	value, _ := doc.Find("#eventDetail").Attr("data-event-id")
	(*eventinfo).I_Event_ID, _ = strconv.Atoi(value)

	log.Printf(" eventid=%s (%d)\n", (*eventinfo).Event_ID, (*eventinfo).I_Event_ID)

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -2
		return
	}
	(*eventinfo).Period = selector.Find(".info").Text()
	eventinfo.Period = strings.Replace(eventinfo.Period, "\u202f", " ", -1)
	period := strings.Split((*eventinfo).Period, " - ")
	if inputmode == "url" {
		(*eventinfo).Start_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[1]+" JST")
	} else {
		(*eventinfo).Start_time, _ = time.Parse("2006/01/02 15:04 MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("2006/01/02 15:04 MST", period[1]+" JST")
	}

	(*eventinfo).EventStatus = "BeingHeld"
	if (*eventinfo).Start_time.After(time.Now()) {
		(*eventinfo).EventStatus = "NotHeldYet"
	} else if (*eventinfo).End_time.Before(time.Now()) {
		(*eventinfo).EventStatus = "Over"
	}

	//	イベントに参加しているルームの数を求めます。
	//	参加ルーム数と表示されているルームの数は違うので注意。ここで取得しているのは参加ルーム数。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)

	return
}
*/
