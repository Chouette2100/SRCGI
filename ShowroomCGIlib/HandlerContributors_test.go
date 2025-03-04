// Copyright © 2025 chouette2100@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php

package ShowroomCGIlib

import (
	"log"

	"net/http"
	"testing"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func TestGetAndSaveContributors(t *testing.T) {

	client, cookiejar, err := exsrapi.CreateNewClient("")
	if err != nil {
		log.Printf("exsrapi.CeateNewClient(): %s", err.Error())
		return //       エラーがあれば、ここで終了
	}
	defer cookiejar.Save()

	//      データベースとの接続をオープンする。
	dbconfig, err := srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("srdblib.OpenDb() error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer srdblib.Db.Close()

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db, Dialect: dial, ExpandSliceArgs: true}
	srdblib.Dbmap.AddTableWithName(Contribution{}, "contribution").SetKeys(false, "Ieventid", "Roomid", "Viewerid")
	srdblib.Dbmap.AddTableWithName(srdblib.Viewer{}, "viewer").SetKeys(false, "Viewerid")
	srdblib.Dbmap.AddTableWithName(srdblib.ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
	// srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
	// srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	type args struct {
		client   *http.Client
		ieventid int
		roomid   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestGetAndSaveContributors",
			args: args{
				client:   client,
				ieventid: 16955,
				roomid:   281387,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetAndSaveContributors(tt.args.client, tt.args.ieventid, tt.args.roomid)
			for _, r := range result {
				log.Printf("irank=%d, viewerid=%d, name=%s, point=%d\n", r.Irank, r.Viewerid, r.Name, r.Point)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAndSaveContributors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
