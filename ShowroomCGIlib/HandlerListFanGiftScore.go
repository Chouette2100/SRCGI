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
	"github.com/Chouette2100/srdblib/v2"
	"github.com/dustin/go-humanize"
	//	"github.com/Chouette2100/exsrapi/v2"
)

//	const MaxAcq = 5	//	表示するデータ数、とりあえず HandlerListCntrb()と同じにしてみる
/*
type VgsHeader struct {
	Giftid    int
	Eventname string
	Period    string
	Maxpoint  int
	Maxacq int
	Limit int
	Gscale    int
	Userno    int
	Username  string
	ShortURL  string
	Ier       int
	Iel       int
	Stime     []time.Time
	Earned    []int
	Total     []int
	Target    []int
	Ifrm      []int
	Nof       []int
	Nft       int //	先頭に戻ったときの最後に表示される枠
	Npb       int //	1ページ戻る
	N1b       int //	一枠戻る
	Ncr       int
	N1f       int
	Npf       int
	Nlt       int
	Vgslist   *[]VgsInf
}
*/

type VgsInf struct {
	Viewerid   int
	Viewername string
	Orderno    int
	Score      []int
	Point      int
	LastName   string
	Giftid     int
}

//	type	CntrbInfList	[] CntrbInf

/*
        HandlerListCntrb()
		（配信枠別）貢献ポイントランキングを5(MaxAcq)枠分表示する。

        引数
		w						http.ResponseWriter
		req						*http.Request

        戻り値
        なし



	0101G0	配信枠別貢献ポイントを実装する。

*/

func ListFanGiftScoreHandler(w http.ResponseWriter, req *http.Request) {

	//	ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
	_, _, isallow := GetUserInf(req)
	if !isallow {
		fmt.Fprintf(w, "Access Denied\n")
		return
	}

	var gsheader GsHeader

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrb-h1.gtpl","templates/list-cntrb-h2.gtpl","templates/list-cntrb.gtpl"))
	funcMap := template.FuncMap{
		"sub":   func(i, j int) int { return i - j },
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
		"t2s":   func(t time.Time, tfmt string) string { return t.Format(tfmt) },
	}
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"templates/list-vgs-h1.gtpl", "templates/list-vgs-h2.gtpl", "templates/list-vgs.gtpl"))

	//	var eventinf exsrapi.Event_Inf
	//	GetEventInf(eventid, &eventinf)

	gsheader.Campaignid = req.FormValue("campaignid")
	if gsheader.Campaignid == "" {
		gsheader.Campaignid = "kingofliver2024summer-autumn"
		gsheader.Campaignname = "SHOWROOMライバー王決定戦summer/autumn"
		gsheader.Url = "https://campaign.showroom-live.com/kingofliver2024/"
	}

	grid, _ := strconv.Atoi(req.FormValue("giftid"))
	if grid == 0 {
		grid = 206
	}

	limit, _ := strconv.Atoi(req.FormValue("limit"))
	if limit == 0 {
		limit = 60
	}

	maxacq, _ := strconv.Atoi(req.FormValue("maxacq"))
	if maxacq == 0 {
		maxacq = 10
	}

	err := GetGiftRanking(&gsheader, grid, "fan")
	if err != nil {
		err = fmt.Errorf("GetGiftRanking(): error %w", err)
		log.Printf("%s", err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	if grid == 0 {
		grid = gsheader.GiftRanking[0].Grid
	}

	acqtimelist, _ := SelectVgsAcqTimeList(grid)
	if len(acqtimelist) == 0 {
		fmt.Fprintf(w, "HandlerFanGiftScore() No AcqTimeList\n")
		fmt.Fprintf(w, "Check eventid and userno\n")
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

	ib := ie - maxacq + 1
	if ib < 0 {
		ib = 0
		//	ie = maxacq - 1
		ie = latl - 1
	}

	//	gsheader.Eventid = eventid
	gsheader.Grid = grid
	gsheader.Maxacq = maxacq
	gsheader.Limit = limit
	//	gsheader.Eventname = eventinf.Event_name

	//	gsheader.Maxpoint = eventinf.Maxpoint
	//	gsheader.Gscale = eventinf.Gscale

	//	gsheader.Period = eventinf.Period
	//	cntrbheader.Userno = userno

	gsheader.Ncr = ie

	//	戻る側の設定
	if ie < maxacq {
		gsheader.Nft = -1
		gsheader.Npb = -1
		gsheader.N1b = -1
	} else {
		gsheader.Nft = maxacq - 1  //	先頭にもどる
		gsheader.Npb = ie - maxacq //	１ページ分戻る
		if gsheader.Npb < maxacq-1 {
			gsheader.Npb = maxacq - 1
		}
		gsheader.N1b = ie - 1 //	一枠分戻る
	}

	if ie == latl-1 {
		gsheader.Nlt = -1
		gsheader.Npf = -1
		gsheader.N1f = -1
	} else {
		gsheader.Nlt = latl - 1    //	最後に進む
		gsheader.Npf = ie + maxacq //	１ページ分進む
		if gsheader.Npf > latl-1 {
			gsheader.Npf = latl - 1
		}
		gsheader.N1f = ie + 1 //	一枠分進む
	}

	//	_, _, _, _, _, _, _, _, roomname, roomurlkey, _, _ := GetRoomInfoByAPI(fmt.Sprintf("%d", userno))
	//	cntrbheader.Username = roomname
	//	cntrbheader.ShortURL = roomurlkey

	tsie := acqtimelist[ie]

	vgslist, viewerid2order, err := SelectViewerid2Order(grid, tsie, limit)
	if err != nil {
		err = fmt.Errorf("SelectUserno2Order() returned %w", err)
		log.Printf("HandlerListGiftScore(): err = %+v", err)
		fmt.Fprintf(w, "HandlerListGiftScore(): err = %+v", err)
		return
	}

	for i := ib; i <= ie; i++ {
		ts := acqtimelist[i]
		gsheader.Stime = append(gsheader.Stime, ts)
		//	log.Printf(" i=%d ts=%+v\n", i, ts)
		err = SelectVgs(grid, ts, &vgslist, viewerid2order)
		if err != nil {
			err = fmt.Errorf("SelectGs() returned %w", err)
			log.Printf("HandlerListGiftScore(): err = %+v", err)
			fmt.Fprintf(w, "HandlerListGiftScore(): err = %+v", err)
			return
		}
		SelectVgsHeader(grid, ts, &gsheader)
		gsheader.Ifrm[i-ib] = i
		gsheader.Nof[i-ib] = i + 1
	}
	gsheader.Vgslist = vgslist

	if ie == latl-1 {
		gsheader.Ier = -1
	} else {
		gsheader.Ier = ie + 5
		if gsheader.Ier > latl-1 {
			gsheader.Ier = latl - 1
		}
	}

	if ie == 0 {
		gsheader.Ier = -1
	} else {
		gsheader.Iel = ie - 5
		if gsheader.Iel < 0 {
			gsheader.Iel = 0
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

	if err := tpl.ExecuteTemplate(w, "list-vgs-h1.gtpl", gsheader); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "list-vgs-h2.gtpl", gsheader); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "list-vgs.gtpl", gsheader); err != nil {
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
func SelectVgsAcqTimeList(
	giftid int,
) (
	acqtimelist []time.Time,
	err error,
) {

	var rows []interface{}

	//	ギフトランキングを取得した時刻の一覧を取得する。
	sqlst := "select distinct ts from viewergiftscore where giftid = ? order by ts "
	rows, err = srdblib.Dbmap.Select(srdblib.ViewerGiftScore{}, sqlst, giftid)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(ViewerGiftScore{}, giftid=%d)  err=%w", giftid, err)
		return
	}

	for _, v := range rows {
		acqtimelist = append(acqtimelist, v.(*srdblib.ViewerGiftScore).Ts)
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
func SelectVgs(
	giftid int,
	ts time.Time,
	vgslist *[]VgsInf,
	viewerid2order map[int]int,
) (
	err error,
) {

	var rows []interface{}

	//	指定した時刻のギフトランキングを取得する。
	sqlst := "select viewerid, score from viewergiftscore "
	sqlst += " where giftid =? and ts = ? order by orderno "
	rows, err = srdblib.Dbmap.Select(srdblib.ViewerGiftScore{}, sqlst, giftid, ts)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(ViewerGiftScore{}, giftid=%d, ts=%+v)  err=%w", giftid, ts, err)
		return
	}

	for i := range *vgslist {
		(*vgslist)[i].Score = append((*vgslist)[i].Score, 0)
	}
	l := len((*vgslist)[0].Score) - 1

	//	��ー�ー��位を取得する。

	for _, v := range rows {
		if idx, ok := viewerid2order[v.(*srdblib.ViewerGiftScore).Viewerid]; ok {
			(*vgslist)[idx].Score[l] = v.(*srdblib.ViewerGiftScore).Score
		}
	}

	/*
		//	�����ポイントランキングを取得する。

		for i := 0; i < len(*gslist); i++ {
			(*gslist)[i].Incremental = append((*gslist)[i].Incremental, -1)
		}
		loc := len((*gslist)[0].Incremental) - 1

		for rows.Next() {
			srdblib.Dberr = rows.Scan(&tlsnid, &incremental)
			if srdblib.Dberr != nil {
				log.Printf("SelectCntrb() (7) err=%s\n", srdblib.Dberr.Error())
				status = -7
				return
			}
			//	log.Printf("%+v\n", cntrbinf)
			i := tlsnid2order[tlsnid]
			(*gslist)[i].Incremental[loc] = incremental
			//	(*cntrbinflist)[i].Incremental = append((*cntrbinflist)[i].Incremental, incremental)
			//	log.Printf(" tlsnid=%d i=%d increment[%d]=%d\n", tlsnid, i, loc, incremental)
		}
	*/

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
func SelectVgsHeader(
	giftid int,
	ts time.Time,
	gsheader *GsHeader,
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
	gsheader.Ifrm = append(gsheader.Ifrm, 0)
	gsheader.Nof = append(gsheader.Nof, 0)

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
func SelectViewerid2Order(
	giftid int,
	ts time.Time,
	limit int,
) (
	vgslist []VgsInf,
	viewerid2order map[int]int,
	err error,
) {

	viewerid2order = make(map[int]int)

	//	指定された時刻の貢献ポイントランキングを取得する。
	type ViewerAndOrderno struct {
		Viewerid int
		Sname    string
		Orderno  int
	}
	var rows []interface{}
	sqlst := "select v.viewerid, v.sname, vgs.orderno "
	sqlst += " from viewer v join viewergiftscore vgs on v.viewerid = vgs.viewerid "
	sqlst += " where vgs.giftid = ? and vgs.ts = ? order by vgs.orderno limit ? "
	rows, err = srdblib.Dbmap.Select(ViewerAndOrderno{}, sqlst, giftid, ts, limit)
	if err != nil {
		err = fmt.Errorf("Dbmap.Select(Viewer{}, giftid=%d)  err=%w", giftid, err)
		return
	}

	for i, v := range rows {
		vw := v.(*ViewerAndOrderno)
		viewerid2order[vw.Viewerid] = i
		vgslist = append(vgslist, VgsInf{
			Viewerid:   vw.Viewerid,
			Viewername: vw.Sname,
			Orderno:    vw.Orderno,
		})
	}
	/*
		sqlst := "select userno, orderno from giftscore "
		sqlst += " where giftid = ? and ts = ? order by orderno limit 10 "
		rows, err = srdblib.Dbmap.Select(srdblib.GiftScore{}, sqlst, giftid, ts)
		if err != nil {
			err = fmt.Errorf("Dbmap.Select(GiftScore{}, giftid=%d)  err=%w", giftid, err)
			return
		}

		var row interface{}
		for i, v := range rows {
			gs := v.(*srdblib.GiftScore)
			userno2order[gs.Userno] = i
			row, err = srdblib.Dbmap.Get(srdblib.User{}, gs.Userno)
			if err != nil {
				err = fmt.Errorf("Dbmap.Get(User{}, userno=%d)  err=%w", gs.Userno, err)
				return
			}
			if row != nil {
				gslist = append(gslist, GsInf{
					Userno:    gs.Userno,
					User_name: row.(*srdblib.User).User_name,
					Orderno:   gs.Orderno,
				})
			} else {
				gslist = append(gslist, GsInf{
					Userno:    gs.Userno,
					User_name: "n/a",
					Orderno:   gs.Orderno,
				})
			}
		}
	*/

	return

}
