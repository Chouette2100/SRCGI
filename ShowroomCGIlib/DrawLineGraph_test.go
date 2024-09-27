// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	"log"
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/srdblib"
)

func TestJtruncate(t *testing.T) {
	type args struct {
		t time.Time
	}

	ts0, _ := time.Parse("2006-01-02 15:04:05 MST", "2024-09-26 00:05:00 JST")
	ts1, _ := time.Parse("2006-01-02 15:04:05 MST", "2024-09-26 08:55:00 JST")
	ts2, _ := time.Parse("2006-01-02 15:04:05 MST", "2024-09-26 09:05:00 JST")
	ts3, _ := time.Parse("2006-01-02 15:04:05 MST", "2024-09-26 23:45:00 JST")

	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "TestJtruncate-0",
			args: args{
				t: ts0,
			},
			want: ts0,
		},
		{
			name: "TestJtruncate-1",
			args: args{
				t: ts1,
			},
			want: ts1,
		},
		{
			name: "TestJtruncate-2",
			args: args{
				t: ts2,
			},
			want: ts2,
		},
		{
			name: "TestJtruncate-3",
			args: args{
				t: ts3,
			},
			want: ts3,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Jtruncate(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Jtruncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDrawLineGraph(t *testing.T) {
	type args struct {
		filename   string
		title0     string
		title1     string
		title2     string
		maxpoint   int
		tmaxpoint  int
		target     int
		start_time time.Time
		end_time   time.Time
		deltax     float64
		IDlist     []int
		xydata     *[]Xydata
	}

	ts, _ := time.Parse("2006-01-02 15:04:05 MST", "2024-09-26 18:00:00 JST")
	te, _ := time.Parse("2006-01-02 15:04:05 MST", "2024-10-02 21:59:00 JST")
	//	tsd := float64(ts.Unix()) / 60 / 60 / 24

	var dbconfig *srdblib.DBConfig
	var err error
	dbconfig, err = srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("Database error. err = %v\n", err)
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}

	srdblib.Db.SetMaxOpenConns(8)
	srdblib.Db.SetMaxIdleConns(12)
	srdblib.Db.SetConnMaxLifetime(time.Minute * 5)
	srdblib.Db.SetConnMaxIdleTime(time.Minute * 5)

	defer srdblib.Db.Close()
	log.Printf("%+v\n", dbconfig)

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db, Dialect: dial, ExpandSliceArgs: true}
	srdblib.Dbmap.AddTableWithName(srdblib.User{}, "user").SetKeys(false, "Userno")
	// srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	// srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	// srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")

	// srdblib.Dbmap.AddTableWithName(srdblib.GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
	// srdblib.Dbmap.AddTableWithName(srdblib.ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	// srdblib.Dbmap.AddTableWithName(srdblib.Viewer{}, "viewer").SetKeys(false, "Viewerid")
	// srdblib.Dbmap.AddTableWithName(srdblib.ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")

	// srdblib.Dbmap.AddTableWithName(srdblib.Campaign{}, "campaign").SetKeys(false, "Campaignid")
	// srdblib.Dbmap.AddTableWithName(srdblib.GiftRanking{}, "giftRanking").SetKeys(false, "Campaignid", "Grid")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestDrawLineGraph",
			args: args{
				filename:   "test.svg",
				title0:     "Test Title 0",
				title1:     "Test Title 1",
				title2:     "Test Title 2",
				maxpoint:   10000,
				tmaxpoint:  0,
				target:     12000,
				start_time: ts,
				end_time:   te,
				deltax:     1.1,
				IDlist:     []int{87911, 111004, 75721},
				xydata: &[]Xydata {
					{
						X: []float64{1.0, 2.0, 5.0},
						Y: []float64{200.0, 1000.0, 2000.0},
					},
					{
						X: []float64{1.0, 4.0, 6.0},
						Y: []float64{7000.0, 8000.0, 9000.0},
					},
					{
						X: []float64{1.5, 2.0, 3.0},
						Y: []float64{0.0, 4000.0, 9999.0},
					},
				},
			},
			wantErr: false,
		},
	}
		// TODO: Add test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DrawLineGraph(
				tt.args.filename,
				tt.args.title0,
				tt.args.title1,
				tt.args.title2,
				tt.args.maxpoint,
				tt.args.tmaxpoint,
				tt.args.target,
				tt.args.start_time,
				tt.args.end_time,
				tt.args.deltax,
				tt.args.IDlist,
				tt.args.xydata); (err != nil) != tt.wantErr {
				t.Errorf("DrawLineGraph() error = %v, wantErr %v", err, tt.wantErr)
		}})
	}
}
