// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	// "strings"
	//	"io" //　ログ出力設定用。必要に応じて。
	//	"sort" //	ソート用。必要に応じて。

	// "html/template"
	"html/template"
	"net/http"

	// "github.com/dustin/go-humanize"

	//	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
	"github.com/Chouette2100/srdblib"
)

/*
終了イベント一覧を作るためのハンドラー

Ver. 0.1.0
*/

func HandlerOldEvents(
	w http.ResponseWriter,
	r *http.Request,
) {
	var err error
	_, _, isallow := GetUserInf(r)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	seventid := r.FormValue("ieventid")
	ieventid, _ := strconv.Atoi(seventid)

        var intf []interface{}
	intf, err = srdblib.Dbmap.Select(&srdblib.Wevent{}, "SELECT * FROM wevent WHERE ieventid = ?", ieventid)
        if err != nil {
		err = fmt.Errorf("srdblib.Dbmap.Select(): %s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
        event := *(intf[0].(*srdblib.Wevent))

	sroomid := r.FormValue("roomid")
	roomid, _ := strconv.Atoi(sroomid)

	client, cookiejar, err := exsrapi.CreateNewClient("")
	if err != nil {
		err = fmt.Errorf("exsrapi.CeateNewClient(): %s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))

		return //       エラーがあれば、ここで終了
	}
	defer cookiejar.Save()

	hcntrbinf := HCntrbInf{
		Ieventid:   ieventid,
		Eventid:    event.Eventid,
		Event_name: event.Event_name,
		Period:     event.Period,
		Roomid:     roomid,
	}

	if err = MakeFileOfContributors(client, &hcntrbinf); err != nil {
		err = fmt.Errorf("MakeFileOfContributors(): %w", err)
                log.Printf("err=%s\n", err.Error())
                w.Write([]byte(err.Error()))
		return
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/contributors.gtpl",
	))

	if err := tpl.ExecuteTemplate(w, "contributors.gtpl", hcntrbinf); err != nil {
		err = fmt.Errorf("tpl.ExceuteTemplate(w,\"contributors.gtpl\", hcntrbinf) err=%s", err.Error())
                log.Printf("err=%s\n", err.Error())
                w.Write([]byte(err.Error()))
	}

}

func GetAndSaveOldEvents(
	client *http.Client,
	roomid int,
) (
	err error,
) {


	// 過去のイベント一覧を取得する
	var rpe *srapi.RoomsPastevents
	rpe, err = srapi.GetRoomsPasteventsByApi(client, roomid)
	if err != nil {
		err = fmt.Errorf("GetRoomsPasteventsByApi(): %w", err)
		return
	}

	// 取得したデータにあるイベント数と実際のイベント数が一致するか確認する
	// ただし、一致しない場合でもエラーとはしない（原因はわからないが、たまに一致しないことがある）
	log.Printf("GetRoomsPasteventsByApi(): rpe.Count=%d, len(rpe.Pastevents)=%d", rpe.TotalEntries, len(rpe.Events))

	// 取得した過去のイベント一覧のうち未保存のものをイベントテーブルに格納する
	for _, pe := range rpe.Events {

		// すでにデータがあるかどうかを確認する
		err = srdblib.Dbmap.SelectOne(&srdblib.Wevent{}, "SELECT * FROM wevent WHERE ieventid = ?", pe.EventID)
		if err == nil {
			log.Printf("GetAndSaveOldEvents(): ieventid=%d is already saved(%s)\n", pe.EventID, pe.EventName)
			continue
		}

		// データが存在しないときはイベント参加ユーザをユーザテーブルに格納した上でイベント情報をイベントテーブルに格納する
		log.Printf("  ieventid=%d\n", pe.EventID)
		
		// イベントに参加したルームの一覧を取得する
		var rli *srapi.RoomListInf
		rli, err = srapi.GetRoominfFromEventByApi(client, pe.EventID, 1, 2000)
		if err != nil {
			err = fmt.Errorf("GetRoominfFromEventByApi(): %w", err)
			return
		}

		// ユーザ情報をユーザテーブルに格納する
		tnow := time.Now().Truncate(time.Second)
		for _, v := range rli.RoomList {
			suid := fmt.Sprintf("%09d", v.Room_id)
			wuser := srdblib.Wuser {
				Userno: v.Room_id,
				Userid: v.Room_url_key,
				User_name: v.Room_name,
				Longname: v.Room_name,
				Shortname: suid[len(suid)-2:],
				Ts: tnow,
			}
			err = srdblib.UpinsWuserSetProperty(client , tnow, &wuser, 1440, 1000)
			if err != nil {
				err = fmt.Errorf("UpinsWuserSetProperty() err=%w", err)
				return
			}
		}


		wevent := srdblib.Wevent{
			Ieventid:   pe.EventID,
		}

		if err = srdblib.Dbmap.Insert(&wevent); err != nil {
			err = fmt.Errorf("Insert(): %w", err)
			return
		}



		// テストのため、格納は一回だけにする
		break
	}

	return

}

// イベントの貢献ポイントランキングをファイルに書き出す
func MakeFileOfOldEvents(
	client *http.Client,
	hc *HCntrbInf,
) (
	err error,
) {
	// 貢献ポイントデータを取得する
	var result []Result
	result, err = GetAndSaveContributors(client, hc.Ieventid, hc.Roomid)
	if err != nil {
		err = fmt.Errorf("GetAndSaveContributors(): %w", err)
		return
	}

	hc.Result = result

	// 貢献ポイントデータをファイルに書き出す
	//	ファイル名の決定
	uqn := ""
	if Serverconfig.WebServer == "None" {
		uqn = time.Now().Format("05.000")
		uqn = strings.Replace(uqn, ".", "", -1)
	} else {
		uqn = fmt.Sprintf("%03d", os.Getpid()%1000)
	}
	hc.Filename = fmt.Sprintf("%s_%d_%d_%s.csv", hc.Eventid, hc.Ieventid, hc.Roomid, uqn)

	file, err := os.OpenFile("public/"+hc.Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		err = fmt.Errorf("os.OpenFile(public/%s) err=%w", hc.Filename, err)
		return
	}
	defer file.Close()

	// ファイルに書き出す
	fmt.Fprintf(file, "%d, %s, \"%s\"\n", hc.Ieventid, hc.Eventid, hc.Event_name)
	fmt.Fprintf(file, ",Period, \"%s\"\n", hc.Period)
	fmt.Fprintf(file, ",Roomid, %d\n", hc.Roomid)
	fmt.Fprintf(file, "Rank,Viewerid,Viewername,Point\n")
	for _, r := range hc.Result {
		fmt.Fprintf(file, "%d,%d,%s,%d\n", r.Irank, r.Viewerid, r.Name, r.Point)
	}

	return
}
