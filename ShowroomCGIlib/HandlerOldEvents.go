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

type OldEvents struct {
	User    *srdblib.Wuser
	Wevents []srdblib.Wevent
	ErrMsg  string
}

// 過去のイベントの一覧を作るためのハンドラー
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

	oe := OldEvents{}
	oe.User = new(srdblib.Wuser)

	suserno := r.FormValue("userno")
	oe.User.Userno, _ = strconv.Atoi(suserno)

	twuser := new(srdblib.Wuser)

	client, cookiejar, err := exsrapi.CreateNewClient("")
	if err != nil {
		err = fmt.Errorf("exsrapi.CeateNewClient(): %s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))

		return //       エラーがあれば、ここで終了
	}
	defer cookiejar.Save()

	oe.Wevents, err = GetAndSaveOldEvents(client, oe.User.Userno)
	if err != nil {
		err = fmt.Errorf("GetAndSaveOldEvents(): %s", err.Error())
		log.Printf("err=%s\n", err.Error())
		oe.ErrMsg = err.Error()
	}

	if len(oe.Wevents) == 0 {
		// 指定した配信者が存在しない、あるいは参加したイベントがない場合
		oe.ErrMsg = "No events or userno is invalid"
	} else {
		// 参加したイベントがある場合はユーザ情報を追加あるいは更新する
		twuser.Userno = oe.User.Userno
		oe.User = twuser

		ru, err := srdblib.UpinsUser(client, time.Now().Truncate(time.Second), twuser, 14400, 5000)
		if err != nil {
			err = fmt.Errorf("UpinsUser(): %s", err.Error())
			log.Printf("err=%s\n", err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		oe.User = *ru
	}

	//	テンプレートで使用する関数を定義する
	funcMap := template.FuncMap{
		"TimeToString":  func(t time.Time) string { return t.Format("01-02 15:04") },
		"TimeToStringY": func(t time.Time) string { return t.Format("06-01-02 15:04") },
		"IsTempID":      func(s string) bool { return strings.HasPrefix(s, "@@@@") },
	}

	// テンプレートをパースする
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"templates/oldevents.gtpl",
	))

	if err := tpl.ExecuteTemplate(w, "oldevents.gtpl", &oe); err != nil {
		err = fmt.Errorf("tpl.ExceuteTemplate(w,\"oldevents.gtpl\", hcntrbinf) err=%s", err.Error())
		log.Printf("err=%s\n", err.Error())
		w.Write([]byte(err.Error()))
	}

}

func GetAndSaveOldEvents(
	client *http.Client,
	roomid int,
) (
	wevents []srdblib.Wevent,
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

	// wuser := srdblib.Wuser{
	// 	Userno: roomid,
	// }

	// srdblib.UpinsWuserSetProperty(client, time.Now().Truncate(time.Second), &wuser, 1440, 5000)

	// 取得した過去のイベント一覧のうち未保存のものをイベントテーブルに格納する
	for _, pe := range rpe.Events {

		wevent := srdblib.Wevent{}

		// すでにデータがあるかどうかを確認する
		is_exist := false
		var itrf []interface{}
		itrf, err = srdblib.Dbmap.Select(&wevent, "SELECT eventid, ieventid, event_name, starttime FROM wevent WHERE ieventid = ?", pe.EventID)
		if err != nil {
			err = fmt.Errorf("Select(): %w", err)
			return
		}
		if len(itrf) != 0 {
			is_exist = true
			log.Printf("GetAndSaveOldEvents(): ieventid=%d is already saved(%s)\n", pe.EventID, pe.EventName)
		} else {
			sqlst := "SELECT eventid, ieventid, event_name, starttime FROM wevent "
			sqlst += " WHERE event_name = ? and starttime = ? "
			itrf, err = srdblib.Dbmap.Select(&wevent, sqlst, pe.EventName, time.Unix(int64(pe.StartedAt), 0))
			if err != nil {
				err = fmt.Errorf("Select(): %w", err)
				return
			}
			if len(itrf) != 0 {
				is_exist = true
				if len(itrf) == 1 {
					wevent = *(itrf[0].(*srdblib.Wevent))
					var itrf1 interface{}
					itrf1, err = srdblib.Dbmap.Get(&srdblib.Wevent{}, wevent.Eventid)
					if err != nil {
						err = fmt.Errorf("Get(): %w", err)
						return
					}
					wevent = *(itrf1.(*srdblib.Wevent))
					wevent.Ieventid = pe.EventID
					_, err = srdblib.Dbmap.Update(&wevent)
					if err != nil {
						err = fmt.Errorf("Update(): %w", err)
						return
					}
					log.Printf("GetAndSaveOldEvents(): ieventid=%d is inserted(%s)\n", pe.EventID, pe.EventName)
				}
			}
		}
		if is_exist {
			// 参加したイベントのデータがすでに保存されている場合
			log.Printf("GetAndSaveOldEvents(): ieventid=%d is already saved(%s)\n", pe.EventID, pe.EventName)
			wevent = *(itrf[0].(*srdblib.Wevent))
			eida := strings.Split(wevent.Eventid, "?")
			if len(eida) != 1 {
				wevent.Eventid = eida[0]
			}
			// イベント参加者として登録されているか？
			var weventuser srdblib.Weventuser
			itrf, err = srdblib.Dbmap.Select(&weventuser,
				"SELECT * FROM weventuser WHERE eventid like ? AND userno = ?", wevent.Eventid+"%", roomid)
			if err != nil {
				err = fmt.Errorf("Select(): %w", err)
				return
			}
			if len(itrf) == 0 {
				// イベント参加者として登録されていない場合
				weventuser = srdblib.Weventuser{
					EventuserBR: srdblib.EventuserBR{
						Eventid:       wevent.Eventid,
						Userno:        roomid,
						Istarget:      "Y",
						Iscntrbpoints: "Y",
						Vld:           0,
						Point:         -1,
					},
					Status: 0,
				}
				srdblib.UpinsEventuserG(&weventuser, time.Now().Truncate(time.Second))
			}

		} else {

			eventid := fmt.Sprintf("@@@@%06d", pe.EventID)
			starttime := time.Unix(int64(pe.StartedAt), 0)
			endtime := time.Unix(int64(pe.EndedAt), 0)
			ssstarttime := starttime.Format("Jan 2, 2006 3:04 PM")
			sendtime := endtime.Format("Jan 2, 2006 3:04 PM")
			wevent = srdblib.Wevent{
				Ieventid:    pe.EventID,
				Eventid:     eventid,
				Event_name:  pe.EventName,
				Starttime:   starttime,
				Endtime:     endtime,
				Period:      ssstarttime + " - " + sendtime,
				Intervalmin: 5,
				Modmin:      4,
				Modsec:      0,
				Fromorder:   1,
				Toorder:     3000,
				Resethh:     4,
				Resetmm:     0,
				Maxdsp:      20,
				Cmap:        2,
			}
			if err = srdblib.Dbmap.Insert(&wevent); err != nil {
				err = fmt.Errorf("Insert(): %w", err)
				return
			}

		}

		wevents = append(wevents, wevent)

		weventuser := srdblib.Weventuser{
			EventuserBR: srdblib.EventuserBR{
				Eventid:       wevent.Eventid,
				Userno:        roomid,
				Istarget:      "Y",
				Iscntrbpoints: "Y",
				Vld:           0,
				Point:         -1,
			},
			Status: 0,
		}

		srdblib.UpinsEventuserG(&weventuser, time.Now().Truncate(time.Second))
		// テストのため、格納は一回だけにする
		// if i == 0 {
		// 	break
		// }
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
