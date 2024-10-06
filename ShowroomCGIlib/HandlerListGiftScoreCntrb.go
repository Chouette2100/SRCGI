package ShowroomCGIlib

import (
	"fmt"
	"log"

	//	"math"
	//	"sort"
	"strconv"
	//	"strings"
	"time"

	//	"bufio"
	//	"os"

	//	"runtime"

	//	"encoding/json"

	"html/template"
	"net/http"

	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"
	//	svg "github.com/ajstarks/svgo/float"
	"github.com/Chouette2100/srdblib"
	"github.com/dustin/go-humanize"
	//	"github.com/Chouette2100/exsrapi"
)

//	const MaxAcq = 5

type GscHeader struct {
	Campaignid   string
	Campaignname string
	Url          string
	Grid         int
	Grname       string
	Grtype       string
	Eventname    string
	Period       string
	Maxpoint     int
	Maxacq       int
	Limit        int
	Gscale       int
	Userno       int
	Username     string
	ShortURL     string
	Ier          int
	Iel          int
	Stime        []time.Time
	Earned       []int
	Total        []int
	Target       []int
	Ifrm         []int
	Nof          []int
	Nft          int //	先頭に戻ったときの最後に表示される枠
	Npb          int //	1ページ戻る
	N1b          int //	一枠戻る
	Ncr          int
	N1f          int //	一枠進む
	Npf          int //	1ページ進む
	Nlt          int //	最後に進んだとき
	Gsclist      *[]GscInf
	GiftRanking  []srdblib.GiftRanking
	Vgslist      []VgsInf
}

type GscInf struct {
	Giftid   int
	Userno   int
	Viewerid int
	Name     string
	Orderno  int
	Score    []int
}

//	type	CntrbInfList	[] CntrbInf

/*
        HandlerListCntrb()
		（配信枠別）貢献ポイントランキングを5(maxacq)枠分表示する。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。

*/

func HandlerListGiftScoreCntrb(w http.ResponseWriter, req *http.Request) {

	var gscheader GscHeader

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrb-h1.gtpl","templates/list-cntrb-h2.gtpl","templates/list-cntrb.gtpl"))
	funcMap := template.FuncMap{
		"sub":   func(i, j int) int { return i - j },
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
		"t2s":   func(t time.Time, tfmt string) string { return t.Format(tfmt) },
	}
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-gsc.gtpl"))

	//	var eventinf exsrapi.Event_Inf
	//	GetEventInf(eventid, &eventinf)

	/*
		campaignid, _ := strconv.Atoi(req.FormValue("campaignid"))
		giftranking, grid1, err := GetGiftRanking(campaignid)
		if err != nil {
			fmt.Fprintf(w, "HandlerListGiftScore() No GiftRanking. Check campaignid.\n")
			log.Printf("HandlerListGiftScore() No GiftRanking. Check campaignid.\n")
			return
		}
	*/

	gscheader.Campaignid = req.FormValue("campaignid")
	if gscheader.Campaignid == "" {
		gscheader.Campaignid = "kingofliver2024summer-autumn"
		gscheader.Campaignname = "SHOWROOMライバー王決定戦summer/autumn"
		gscheader.Url = "https://campaign.showroom-live.com/kingofliver2024/"
	}

	grid, _ := strconv.Atoi(req.FormValue("giftid"))
	userno, _ := strconv.Atoi(req.FormValue("userno"))
	/*
		if grid == 0 {
			grid = grid1
		}
	*/
	limit, _ := strconv.Atoi(req.FormValue("limit"))
	if limit == 0 {
		limit = 200
	}

	maxacq, _ := strconv.Atoi(req.FormValue("maxacq"))
	if maxacq == 0 {
		maxacq = 10
	}

	ifgr, _ := srdblib.Dbmap.Get(srdblib.GiftRanking{}, gscheader.Campaignid, grid)
	ifuser, _ := srdblib.Dbmap.Get(srdblib.User{}, userno)

	//	err := GetGiftRanking(&gsheader, grid, "liver")
	//	if err != nil {
	//		err = fmt.Errorf("GetGiftRanking(): error %w", err)
	//		log.Printf("%s", err.Error())
	//		w.Write([]byte(err.Error()))
	//		return
	//	}
	//	if grid == 0 {
	//		grid = gsheader.GiftRanking[0].Grid
	//	}

	//	for i := 0; i < len(gsheader.GiftRanking); i++ {
	//		if gsheader.GiftRanking[i].Grid == grid {
	//			gsheader.Grtype =gsheader.GiftRanking[i].Grtype
	//			break
	//		}
	//	}

	acqtimelist, _ := SelectGscAcqTimeList(grid, userno)
	if len(acqtimelist) == 0 {
		fmt.Fprintf(w, "HandlerListGiftScore() No AcqTimeList\n")
		fmt.Fprintf(w, "Check (campaignid,) grid and userno!\n")
		log.Printf("No AcqTimeList\n")
		return
	}

	latl := len(acqtimelist)

	sie := req.FormValue("ie")
	ie := 0
	if sie != "" {
		ie, _ = strconv.Atoi(sie)
	} else {
		ie = latl - 1
	}

	//	log.Printf(". eventid=%s, userno=%d, ie=%d\n", eventid, userno, ie)

	//	ib := ie - maxacq + 1
	ib := ie - maxacq + 1
	if ib < 0 {
		ib = 0
		//	ie = maxacq - 1
		ie = latl - 1
	}

	//	gsheader.Eventid = eventid
	//	gsheader.GiftRanking = giftranking
	gscheader.Grid = grid
	gscheader.Grname = ifgr.(*srdblib.GiftRanking).Grname
	gscheader.Userno = userno
	gscheader.Username = ifuser.(*srdblib.User).User_name
	gscheader.Maxacq = maxacq
	gscheader.Limit = limit
	//	gsheader.Eventname = eventinf.Event_name

	//	gsheader.Maxpoint = eventinf.Maxpoint
	//	gsheader.Gscale = eventinf.Gscale

	//	gsheader.Period = eventinf.Period
	//	cntrbheader.Userno = userno

	//	戻る側の設定
	if ie < maxacq {
		gscheader.Nft = -1
		gscheader.Npb = -1
		gscheader.N1b = -1
	} else {
		gscheader.Nft = maxacq - 1  //	先頭にもどる
		gscheader.Npb = ie - maxacq //	１ページ分戻る
		if gscheader.Npb < maxacq-1 {
			gscheader.Npb = maxacq - 1
		}
		gscheader.N1b = ie - 1 //	一枠分戻る
	}

	if ie == latl-1 {
		gscheader.Nlt = -1
		gscheader.Npf = -1
		gscheader.N1f = -1
	} else {
		gscheader.Nlt = latl - 1    //	最後に進む
		gscheader.Npf = ie + maxacq //	１ページ分進む
		if gscheader.Npf > latl-1 {
			gscheader.Npf = latl - 1
		}
		gscheader.N1f = ie + 1 //	一枠分進む
	}

	//	_, _, _, _, _, _, _, _, roomname, roomurlkey, _, _ := GetRoomInfoByAPI(fmt.Sprintf("%d", userno))
	//	cntrbheader.Username = roomname
	//	cntrbheader.ShortURL = roomurlkey

	tsie := acqtimelist[ie]

	gsclist, viewerid2order, err := SelectViewer2Order(grid, userno, tsie, limit)
	if err != nil {
		err = fmt.Errorf("SelectViewerid2Order() returned %w", err)
		log.Printf("HandlerListGiftScoreCntrb(): err = %+v", err)
		fmt.Fprintf(w, "HandlerListGiftScoreCntrb(): err = %+v", err)
		return
	}

	for i := ib; i <= ie; i++ {
		ts := acqtimelist[i]
		gscheader.Stime = append(gscheader.Stime, ts)
		//	log.Printf(" i=%d ts=%+v\n", i, ts)
		err = SelectGsc(grid, userno, ts, &gsclist, viewerid2order)
		if err != nil {
			err = fmt.Errorf("SelectGs() returned %w", err)
			log.Printf("HandlerListGiftScore(): err = %+v", err)
			fmt.Fprintf(w, "HandlerListGiftScore(): err = %+v", err)
			return
		}
		SelectGscHeader(grid, ts, &gscheader)
		gscheader.Ifrm[i-ib] = i
		gscheader.Nof[i-ib] = i + 1
	}
	gscheader.Gsclist = &gsclist

	if ie == latl-1 {
		gscheader.Ier = -1
	} else {
		gscheader.Ier = ie + 5
		if gscheader.Ier > latl-1 {
			gscheader.Ier = latl - 1
		}
	}

	if ie == 0 {
		gscheader.Ier = -1
	} else {
		gscheader.Iel = ie - 5
		if gscheader.Iel < 0 {
			gscheader.Iel = 0
		}
	}
	/*
		//	順位のないデータ（＝ボーナスポイント）の個数を求める。
		sqlsc := "select count(*) from eventrank where eventid = ? and userid = ? and ts = ? and nrank = 0"
		norow := 0
		srdblib.Db.QueryRow(sqlsc, giftid, tsie).Scan(&norow)
		if norow != 0 {
			//	ボーナスポイントのデータがある
			for i := range gslist {
				if i < norow {
					//	ボーナスポイント
					gslist[i].Ranking = -1
				} else if gslist[i].Ranking > 0 && gslist[i].Ranking < 999 {
					//	獲得ポイント
					gslist[i].Ranking -= norow
				}
			}
			if gslist[1].Point == 0 {
				//	ボーナスポイント部分の2番目でポイントが0のものは除く
				gslist[1] = gslist[0]
				gslist = gslist[1:]
			}
		}
	*/

	//	gsheader.Gsinflist = &gsinflist

	if err := tpl.ExecuteTemplate(w, "list-gsc.gtpl", gscheader); err != nil {
		w.Write([]byte(err.Error()))
		log.Println(err)
	}
}

/*
	        SelectGsAcqTimeList()
			指定したgiftidのギフトランキングを取得した時刻の一覧を取得する。

	        引数
			gift			int				ギフトid

	        戻り値
	        acqtimelist		[] time.Time	取得時刻一覧
*/
func SelectGscAcqTimeList(
	grid int,
	userno int,
) (
	acqtimelist []time.Time,
	err error,
) {

	var rows []interface{}

	//	ギフトランキングを取得した時刻の一覧を取得する。
	sqlst := "select distinct ts from giftscorecntrb where giftid = ? and userno = ? order by ts "
	rows, err = srdblib.Dbmap.Select(srdblib.GiftScoreCntrb{}, sqlst, grid, userno)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(GiftScoreCntrb{}, grid=%d, userno=%d)  err=%w", grid, userno, err)
		return
	}

	for _, v := range rows {
		acqtimelist = append(acqtimelist, v.(*srdblib.GiftScoreCntrb).Ts)
	}

	return

}

/*
			指定したイベント、ユーザー、時刻の貢献ポイントランキングを取得する。
			ここでは順位と累計貢献ポイントは取得しない。

	        引数
			eventid			string			イベントID
			userno			int				ユーザーID
			ts				int				ユーザーID
			loc				int				取得データの格納位置
			loc				int				データの格納場所（ 0 だったら先頭）

	        戻り値
	        cntrbinflist	[] CntrbInf		貢献ポイントランキング（最終貢献ポイント順）
			stats			int				== 0 正常終了	!= 0 データベースアクセス時のエラー
*/
func SelectGsc(
	grid int,
	userno int,
	ts time.Time,
	gsclist *[]GscInf,
	viewerid2order map[int]int,
) (
	err error,
) {

	var rows []interface{}

	//	指定した時刻のギフトランキングを取得する。
	sqlst := "select viewerid, score from giftscorecntrb "
	sqlst += " where giftid =? and userno = ? and ts = ? order by orderno "
	rows, err = srdblib.Dbmap.Select(srdblib.GiftScoreCntrb{}, sqlst, grid, userno, ts)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(GiftScoreCntrb{}, grid=%d, ts=%+v)  err=%w", grid, ts, err)
		return
	}

	for i := range *gsclist {
		(*gsclist)[i].Score = append((*gsclist)[i].Score, 0)
	}
	l := len((*gsclist)[0].Score) - 1

	//	��ー�ー��位を取得する。

	for _, v := range rows {
		if idx, ok := viewerid2order[v.(*srdblib.GiftScoreCntrb).Viewerid]; ok {
			(*gsclist)[idx].Score[l] = v.(*srdblib.GiftScoreCntrb).Score
		}
	}
	return

}

/*
	        SelectCntrbHeader()
			貢献ランキング表のヘッダ部分に必要な配信開始・終了時刻を取得する。

	        引数
			eventid			string			イベントID
			userno			int				配信者ID
			ts				time.Time		枠を特定する時刻（＝貢献ランキングを取得した時刻）
			cntrbheader		*CntrbHeader	配信開始・終了時刻を格納する構造体

	        戻り値
			status			int				終了ステータス（ 0: 正常、　1: DBアクセスでの異常）
*/
func SelectGscHeader(
	grid int,
	ts time.Time,
	gscheader *GscHeader,
) (
	status int,
) {

	//	var stime, etime time.Time
	//	var earned, total int

	status = 0

	//	sql := "select stime, etime, earnedpoint, totalpoint from timetable where eventid = ? and userid = ? and sampletm2 = ? "
	//	srdblib.Dberr = srdblib.Db.QueryRow(sql, giftid, ts).Scan(&stime, &etime, &earned, &total)

	//	if srdblib.Dberr != nil {
	//		//	log.Printf("select stime, etime from timetable where eventid = %s and userid = %d and sampletm2 = %+v\n", eventid, userno, ts)
	//		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
	//		status = -11
	//		return
	//	}

	//	(*cntrbheader).S_stime = append((*cntrbheader).S_stime, stime.Format("02 15:04"))
	//	(*cntrbheader).S_etime = append((*cntrbheader).S_etime, etime.Format("02 15:04"))
	//	(*cntrbheader).Earned = append((*cntrbheader).Earned, earned)
	//	(*cntrbheader).Total = append((*cntrbheader).Total, total)
	(*gscheader).Ifrm = append((*gscheader).Ifrm, 0)
	(*gscheader).Nof = append((*gscheader).Nof, 0)

	return

}

/*
	        SelectTlsnid2Order()
			指定したイベント、配信者、枠の仮リスナーIDと貢献ポイントランキングの対応表を作成する。

	        引数
			eventid			string			イベントID
			userno			int				配信者ID
			ts				time.Time		枠を特定する時刻（＝貢献ランキングを取得した時刻）

	        戻り値
			cntrbinflist	[]CntrbInf		貢献ポイントランキングを格納するための構造体の配列
											ここで累計貢献ポイント、リスナー名が格納される。
			tlsnid2order	map[int]int		仮リスナーIDと貢献ポイントランキングの対応表
			status			int				終了ステータス（ 0: 正常、　1: DBアクセスでの異常）
*/
func SelectViewer2Order(
	grid int,
	userno int,
	ts time.Time,
	limit int,
) (
	gsclist []GscInf,
	viewerid2order map[int]int,
	err error,
) {

	viewerid2order = make(map[int]int)

	//	指定された時刻の貢献ポイントランキングを取得する。
	type ViewerAndGsc struct {
		Viewerid int
		Sname    string
		Orderno  int
	}
	var rows []interface{}
	sqlst := "select v.viewerid, v.sname, gsc.orderno from viewer v join giftscorecntrb gsc on v.viewerid = gsc.viewerid "
	sqlst += " where gsc.giftid = ? and gsc.userno = ? and gsc.ts = ? order by orderno limit ? "
	rows, err = srdblib.Dbmap.Select(ViewerAndGsc{}, sqlst, grid, userno, ts, limit)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(ViewerAndGsc{}, grid=%d)  err=%w", grid, err)
		return
	}

	for i, x := range rows {
		v := x.(*ViewerAndGsc)
		viewerid2order[v.Viewerid] = i
		gsclist = append(gsclist, GscInf{
			Viewerid:    v.Viewerid,
			Name: v.Sname,
			Orderno:   v.Orderno,
		})
	}
	return
}
/*
func GetGiftScoreCntrb(
	gscheader *GscHeader,
	grid int,
	grtype string,
) (
	err error,
) {
	sqlst := "select grid, grname, grtype from giftranking "
	//	sqlst += " where campaignid = ? and grtype = ? order by norder "
	sqlst += " where campaignid = ? and grtype = ? order by endedat desc, startedat desc, norder "
	rows, err := srdblib.Dbmap.Select(srdblib.GiftRanking{}, sqlst, gscheader.Campaignid, grtype)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(GiftScore{}, campaignid=%s)  err=%w", gscheader.Campaignid, err)
		return err
	}
	gscheader.GiftRanking = make([]srdblib.GiftRanking, 0, len(rows))
	for _, v := range rows {
		vgr := v.(*srdblib.GiftRanking)
		gscheader.GiftRanking = append(gscheader.GiftRanking, *vgr)
		if vgr.Grid == grid {
			gscheader.Grid = vgr.Grid
			gscheader.Grname = vgr.Grname
			gscheader.Grtype = vgr.Grtype
		}
	}

	if grid == 0 {
		gscheader.Grid = gscheader.GiftRanking[0].Grid
		gscheader.Grname = gscheader.GiftRanking[0].Grname
		gscheader.Grtype = gscheader.GiftRanking[0].Grtype
	}

	return
}
*/