package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	"bytes"
	"fmt"
	"log"

	//	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"bufio"
	"os"

	"runtime"

	"encoding/json"

	"html/template"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/PuerkitoBio/goquery"

	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
	"github.com/Chouette2100/srdblib"
)

/*

	0100L1	安定版（～2021.12.26）
	0100M0	vscodeでの指摘箇所の修正
	0101A0	LinuxとMySQL8.0に対応する。
	0101B0	OSとWebサーバに応じた処理を行うようにする。アクセスログを作成する。
	0101B1	実行時パラメータをファイルから与えるように変更する。
	0101C0	GetRoomInfoByAPI()に配信開始時刻の取得を追加する。
	0101D0	詳細なランク情報の導入（Nrank）
	0101D1	"Next Live"の表示を追加する。
	0101D2	GetScoreEvery5Minutes RU20E4 に適合するバージョン
	0101D3	ランクをshow_rank_subdividedからleague_labe + lshow_rank_subdivided にする。
	0101E1	環境設定ファイルをyaml形式に変更する。
	0101G0	配信枠別貢献ポイントを導入する。
	0101G1	list-last.gtplでは維新枠別貢献ポイントの記録があるルームのみリンクを作成する。
	0101G2	list-last.gtplにジャンルを追加した。
	0101G3	リスナー貢献ポイントの履歴の表示(list-cntrbH)を作成する。
	0101G4	一つの貢献ポイントランキングの表示(list-cntrbS)を作成する(リスナー名の突き合わせのチェックが主目的)
	0101G5	list-lasth.gtplのリロード予告の表示でデータ取得間隔が5分と固定されていたものを設定値に合わせるように変更する。
	0101G6	ModminがIntervalminに対して不適切な値のときは修正して保存する。
	0101G7	ランクに関しnext_scoreに加えprev_scoreの表示を追加する。ファンの数の表示を追加する。
	0101H0	ファンレベル(/fanlevel)に関する画面を追加する。
	0101J0	ファンダム王イベント参加者のファン数ランキングを作成する。
	0101J2	終了したイベントについては無条件にルーム詳細情報（ランキング、フォロワ、レベル、ファン数）を出力しない。
	0101J2a	"ルーム詳細情報”の説明を追加した。
	0101J3	イベントリストのRoom_IDに変えてルーム名を表示する。表示数を6から10にする。
	0101J4	NewDocument()をNewDocumentFromReader()に変更する。list-last_h.gtplにルーム情報詳細表示/非表示のボタンを追加する。
	0101J5	イベント選択（最近のイベント）にnobasisが0のイベントも表示する（テーブルuserにusernoが0のデータを追加することが必要）
	10AA00	枠別貢献ポイントの「目標値(推定)」を追加する。
	10AB00	枠別貢献ポイントのポイント、増分の表示でhumanaizeを使用する。リスナー別貢献ポイント履歴に達成状況欄を追加する。
	10AB01	WebserverをDbconfig.Webserverに置き換える。枠別貢献ポイントのGTPLを変更する。GTPLのCRLFをLFに変更する。
	10AC00	イベント一覧にModminとModsecを追加する。
	10AD00	longnameの初期値をusernameに変更する。Apache2LinuxをApache2Ubuntuに訂正する。
	10AD01	HandlerGraphPerDay()とHandlerGraphPerSlot()のApache2LinuxをApache2Ubuntuに訂正する。
	10AD02	SelectScoreList()で前配信期間、獲得ポイントのデータがないとき上位のデータがコピーして使われないようにする。
	10AE00	1. 貢献ポイントランキングの表示を詳細化する。　2. イベント情報ページの内容を表示できるようにする（次回更新準備）
	10AE00a	ルーム情報に下位ランクとの差を追加する（GTPLのみの変更）
	10AF00	ブロックランキングに対応する（Event_id=30030のみの暫定対応）
	10AF01	基準となる配信者のリストのサイズが0のときは対応する処理を行わない（異常終了対策、通常の運用では起きない）
	10AG00	イベントリストの取得件数を10から20に変更する。
	10AH00	stmt,err := Db.Prepare("...")　と対になる defer stmt.Close() を追加する。
	10AJ00	ブロックランキングに仮対応（Event_id=30030以外に拡張）する。イベントリストの表示イベント数を設定可能とする。
	10AK00	ブロックランキングのイベント名にblockidを追加する。獲得ポイントの取得時刻の初期値を分散させる。
	------------------------------------- 以下公開版 ----------------------------------------------
	10AL00	イベント情報にieventid（本来のイベントID、5桁程度の整数）を追加する。
	10AL01	イベントの配信者リストを取得するとき、順位にかかわらず獲得ポイントデータを取得する設定とする。
	10AL02	イベントの配信者リストを取得するとき、順位にかかわらず獲得ポイントを表示する設定とする。
	10AM00	Room_url_keyから取り除く文字列を"/"から"/r/"に変更する。
	10AN00	ブロックランキングで貢献ポイントランキングへのリンクを作るときはイベントIDのからブロックIDを取り除く。
	10AP00	DBサーバーに接続するときSSHの使用を可能にする。
	10AQ00	GetWeightedCnt()で周回数の多い獲得ポイントの採用率が上がるように調整する。
	10AQ01	MakePointPerSlot()のperslotの変数宣言をループの中に入れる（毎回初期化されるように）
	11AA0l	データベースへのアクセスをsrdblibに移行しつつある。グラフ表示で縮尺の設定を可能とする。
	11AA02	intervalmin の値を5固定とする（異常終了に対する緊急対応）
	11AA03	intervalminとintervalmin の適正でない入力を排除する。
	11AB00	Event_Infの参照先をsrdblibからexsrapiに変更する。
	11AB01	データベース保存時、Intervalminが0のときは強制的に5にする。
	11AB02	データベース保存時、Intervalminが5でないときは強制的に5にする。
	11AC00	開催中イベント一覧の機能を作成し関連箇所を修正する。
	11AC01 FindPtPerSlot()でPrepare()に対するdefer Close()の抜けを補う。
	11AC02 HandleListCntrb()でボーナスポイントに対する対応を行う。
	11AC03 currentevent.gtpl 1行おきに背景色を変える。list-last_h.gtpl 結果が反映される時刻を正す。
	11AD00 「SHOWROOMイベント情報ページからDBへのイベント参加ルーム情報の追加と更新」でイベントパラーメータがクリアされる問題を解決する。 
	11AE00	HandlerEventRoomList()でブロックイベントの参加ルーム一覧も表示できるようにする。
	11AF00	開催予定イベント一覧の機能を追加する（HandlerScheduledEvent()）
	11AF01	新規イベントの登録ができなくなった問題（＝11AD00の修正で発生したデグレード）に対応する
	11AG00	srdblib.SelectFromEvent()の実行前にはsrdblib.Tevent = "event"を行う。 これはSelectFromEvent()の引数とすべき。


*/

const Version = "11AG00"

/*
type Event_Inf struct {
	Event_ID    string
	I_Event_ID  int
	Event_name  string
	Event_no    int
	MaxPoint    int
	Start_time  time.Time
	Sstart_time string
	Start_date  float64
	End_time    time.Time
	Send_time   string
	Period      string
	Dperiod     float64
	Intervalmin int
	Modmin      int
	Modsec      int
	Fromorder   int
	Toorder     int
	Resethh     int
	Resetmm     int
	Nobasis     int
	Maxdsp      int
	NoEntry     int
	NoRoom      int    //	ルーム数
	EventStatus string //	"Over", "BeingHeld", "NotHeldYet"
	Pntbasis    int
	Ordbasis    int
	League_ids  string
	Cmap        int
	Target      int
	Maxpoint    int
	//	Status		string		//	"Confirmed":	イベント終了日翌日に確定した獲得ポイントが反映されている。
}
*/

type LongName struct {
	Name string
}

type Point struct {
	Pnt  int
	Spnt string
	Tpnt string
}

type PointRecord struct {
	Day       string
	Tday      time.Time
	Pointlist []Point
}

type PointPerDay struct {
	Eventid         string
	Eventname       string
	Period          string
	Usernolist      []int
	Longnamelist    []LongName
	Pointrecordlist []PointRecord
}

type RoomLevel struct {
	User_name string
	Genre     string
	Rank      string
	Nrank     string
	Prank     string
	Level     int
	Followers int
	Fans      int
	Fans_lst  int
	ts        time.Time
	Sts       string
}

type RoomLevelInf struct {
	Userno        int
	User_name     string
	RoomLevelList []RoomLevel
}

type PerSlot struct {
	Timestart time.Time
	Dstart    string
	Tstart    string
	Tend      string
	Point     string
	Ipoint    int
	Tpoint    string
}

type PerSlotInf struct {
	Eventname   string
	Eventid     string
	Period      string
	Roomname    string
	Roomid      int
	Perslotlist []PerSlot
}

type ColorInf struct {
	Color      string
	Colorvalue string
	Selected   string
}

type ColorInfList []ColorInf

type RoomInfo struct {
	Name      string //	ルーム名のリスト
	Longname  string
	Shortname string
	Account   string //	アカウントのリスト、アカウントは配信のURLの最後の部分の英数字です。
	ID        string //	IDのリスト、IDはプロフィールのURLの最後の部分で5～6桁の数字です。
	Userno    int
	//	APIで取得できるデータ(1)
	Genre      string
	Rank       string
	Irank      int
	Nrank      string
	Prank      string
	Followers  int
	Sfollowers string
	Fans       int
	Fans_lst   int
	Level      int
	Slevel     string
	//	APIで取得できるデータ(2)
	Order        int
	Point        int //	イベント終了後12時間〜36時間はイベントページから取得できることもある
	Spoint       string
	Istarget     string
	Graph        string
	Iscntrbpoint string
	Color        string
	Colorvalue   string
	Colorinflist ColorInfList
	Formid       string
	Eventid      string
	Status       string
	Statuscolor  string
}

type RoomInfoList []RoomInfo

// sort.Sort()のための関数三つ
func (r RoomInfoList) Len() int {
	return len(r)
}

func (r RoomInfoList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomInfoList) Choose(from, to int) (s RoomInfoList) {
	s = r[from:to]
	return
}

var SortByFollowers bool

// 降順に並べる
func (r RoomInfoList) Less(i, j int) bool {
	//	return e[i].point < e[j].point
	if SortByFollowers {
		return r[i].Followers > r[j].Followers
	} else {
		return r[i].Point > r[j].Point
	}
}

var Serverconfig *ServerConfig

/*
var Sshconfig *SSHConfig
var Dialer sshql.Dialer
*/

var Event_inf exsrapi.Event_Inf

/*
var Db *sql.DB
var Err error
*/

var OS string

//	var WebServer string

type Color struct {
	Name  string
	Value string
}

// https://www.fukushihoken.metro.tokyo.lg.jp/kiban/machizukuri/kanren/color.files/colorudguideline.pdf
var Colorlist2 []Color = []Color{
	{"red", "#FF2800"},
	{"yellow", "#FAF500"},
	{"green", "#35A16B"},
	{"blue", "#0041FF"},
	{"skyblue", "#66CCFF"},
	{"lightpink", "#FFD1D1"},
	{"orange", "#FF9900"},
	{"purple", "#9A0079"},
	{"brown", "#663300"},
	{"lightgreen", "#87D7B0"},
	{"white", "#FFFFFF"},
	{"gray", "#77878F"},
}

var Colorlist1 []Color = []Color{
	{"cyan", "cyan"},
	{"magenta", "magenta"},
	{"yellow", "yellow"},
	{"royalblue", "royalblue"},
	{"coral", "coral"},
	{"khaki", "khaki"},
	{"deepskyblue", "deepskyblue"},
	{"crimson", "crimson"},
	{"orange", "orange"},
	{"lightsteelblue", "lightsteelblue"},
	{"pink", "pink"},
	{"sienna", "sienna"},
	{"springgreen", "springgreen"},
	{"blueviolet", "blueviolet"},
	{"salmon", "salmon"},
	{"lime", "lime"},
	{"red", "red"},
	{"darkorange", "darkorange"},
	{"skyblue", "skyblue"},
	{"lightpink", "lightpink"},
}

type Event struct {
	EventID   string
	EventName string
	Period    string
	Starttime time.Time
	S_start   string
	Endtime   time.Time
	S_end     string
	Status    string
	Pntbasis  int
	Modmin    int
	Modsec    int
	Pbname    string
	Selected  string
	Maxpoint  int
	Gscale    int
}

type User struct {
	Userno       int
	Userlongname string
	Selected     string
}

type CurrentScore struct {
	Rank      int
	Srank     string
	Userno    int
	Shorturl  string
	Eventid   string
	Username  string
	Roomgenre string
	Roomrank  string
	Roomnrank string
	Roomprank string
	Roomlevel string
	Followers string
	Fans      int
	Fans_lst  int
	NextLive  string
	Point     int
	Spoint    string
	Sdfr      string
	Pstatus   string
	Ptime     string
	Qstatus   string
	Qtime     string
	Bcntrb    bool
}

func GetSerialFromYymmddHhmmss(yymmdd, hhmmss string) (tserial float64) {

	var year, month, day, hh, mm, ss int

	t19000101 := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)

	fmt.Sscanf(yymmdd, "%d/%d/%d", &year, &month, &day)
	fmt.Sscanf(hhmmss, "%d:%d:%d", &hh, &mm, &ss)

	t1 := time.Date(year, time.Month(month), day, hh, mm, ss, 0, time.Local)

	tserial = t1.Sub(t19000101).Minutes() / 60.0 / 24.0

	return
}

func GetUserInfForHistory() (status int) {

	status = 0

	//	select distinct(nobasis) from event
	stmt, err := srdblib.Db.Prepare("select distinct(nobasis) from event")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var roominf RoomInfo
	var roominflist RoomInfoList

	for rows.Next() {
		err := rows.Scan(&roominf.Userno)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		if roominf.Userno != 0 {
			roominf.ID = fmt.Sprintf("%d", roominf.Userno)
			roominflist = append(roominflist, roominf)
		}
	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	eventid := ""

	//	Update user , Insert into userhistory
	for _, roominf := range roominflist {

		sql := "select currentevent from user where userno = ?"
		err := srdblib.Db.QueryRow(sql, roominf.Userno).Scan(&eventid)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
		}

		roominf.Genre, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level,
			roominf.Followers, roominf.Fans, roominf.Fans_lst, roominf.Name, roominf.Account, _, status = GetRoomInfoByAPI(roominf.ID)
		InsertIntoOrUpdateUser(time.Now().Truncate(time.Second), eventid, roominf)
	}

	return
}

func GetEventListByAPI(eventinflist *[]exsrapi.Event_Inf) (status int) {

	status = 0

	last_page := 1
	total_count := 1

	for page := 1; page <= last_page; page++ {

		URL := "https://www.showroom-live.com/api/event/search?page=" + fmt.Sprintf("%d", page)
		log.Printf("GetEventListByAPI() URL=%s\n", URL)

		resp, err := http.Get(URL)
		if err != nil {
			//	一時的にデータが取得できない。
			log.Printf("GetEventListByAPI() err=%s\n", err.Error())
			//		panic(err)
			status = -1
			return
		}
		defer resp.Body.Close()

		//	JSONをデコードする。
		//	次の記事を参考にさせていただいております。
		//		Go言語でJSONに泣かないためのコーディングパターン
		//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

		var result interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			log.Printf("GetEventListByAPI() err=%s\n", err.Error())
			//	panic(err)
			status = -2
			return
		}

		if page == 1 {
			value, _ := result.(map[string]interface{})["last_page"].(float64)
			last_page = int(value)
			value, _ = result.(map[string]interface{})["total_count"].(float64)
			total_count = int(value)
			log.Printf("GetEventListByAPI() total_count=%d, last_page=%d\n", total_count, last_page)
		}

		noroom := 30
		if page == last_page {
			noroom = total_count % 30
			if noroom == 0 {
				noroom = 30
			}
		}

		for i := 0; i < noroom; i++ {
			var eventinf exsrapi.Event_Inf

			tres := result.(map[string]interface{})["event_list"].([]interface{})[i]

			ttres := tres.(map[string]interface{})["league_ids"]
			norec := len(ttres.([]interface{}))
			if norec == 0 {
				continue
			}
			log.Printf("norec =%d\n", norec)
			eventinf.League_ids = ""
			/*
				for j := 0; j < norec; j++ {
					eventinf.League_ids += ttres.([]interface{})[j].(string) + ","
				}
			*/
			eventinf.League_ids = ttres.([]interface{})[norec-1].(string)
			if eventinf.League_ids != "60" {
				continue
			}

			eventinf.Event_ID, _ = tres.(map[string]interface{})["event_url_key"].(string)
			eventinf.Event_name, _ = tres.(map[string]interface{})["event_name"].(string)
			//	log.Printf("id=%s, name=%s\n", eventinf.Event_ID, eventinf.Event_name)

			started_at, _ := tres.(map[string]interface{})["started_at"].(float64)
			eventinf.Start_time = time.Unix(int64(started_at), 0)
			eventinf.Sstart_time = eventinf.Start_time.Format("06/01/02 15:04")
			ended_at, _ := tres.(map[string]interface{})["ended_at"].(float64)
			eventinf.End_time = time.Unix(int64(ended_at), 0)
			eventinf.Send_time = eventinf.End_time.Format("06/01/02 15:04")

			(*eventinflist) = append((*eventinflist), eventinf)

		}

		//	resp.Body.Close()
	}

	return
}

// idで指定した配信者さんの獲得ポイントを取得する。
// 戻り値は 獲得ポイント、順位、上位とのポイント差（1位の場合は2位とのポイント差）、イベント名
// レベルイベントのときは順位、上位とのポイント差は0がセットされる。
func GetPointsByAPI(id string) (Point, Rank, Gap int, EventID string) {

	//	獲得ポイントなどの配信者情報を得るURL（このURLについては記事参照）
	URL := "https://www.showroom-live.com/api/room/event_and_support?room_id=" + id

	resp, err := http.Get(URL)
	if err != nil {
		//	一時的にデータが取得できない。
		//		panic(err)
		return 0, 0, 0, "**Error** http.Get(URL)"
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		return 0, 0, 0, "**Error** http.Get(URL)"
	}

	//	イベントが終わっている、イベント参加をとりやめた、SHOWROOMをやめた、などの対応
	if result.(map[string]interface{})["event"] == nil {
		return 0, 0, 0, "not held yet./over./not entry."
	}

	if result.(map[string]interface{})["event"].(map[string]interface{})["ranking"] != nil {
		//	ランキングのあるイベントの場合
		//	（順位に応じて特典が与えられるイベント、ただし獲得ポイントに対して特典が与えられる場合でも順位付けがある場合はこちら）

		//	獲得ポイント
		l, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["point"].(float64)
		//	順位
		m, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["rank"].(float64)
		//	ポイント差
		n, _ := result.(map[string]interface{})["event"].(map[string]interface{})["ranking"].(map[string]interface{})["gap"].(float64)

		Point = int(l)
		Rank = int(m)
		Gap = int(n)

		//	イベント名
		EventID, _ = result.(map[string]interface{})["event"].(map[string]interface{})["event_url"].(string)
		EventID = strings.Replace(EventID, "https://www.showroom-live.com/event/", "", -1)

	} else if result.(map[string]interface{})["event"].(map[string]interface{})["quest"] != nil {
		//	レベルイベント（ランキングのないイベント）の場合
		//	（アバ権やステッカーなど獲得ポイントに応じて特典が与えられるイベント、ただし順位付けがある場合は除く）

		//	獲得ポイント
		l, _ := result.(map[string]interface{})["event"].(map[string]interface{})["quest"].(map[string]interface{})["support"].(map[string]interface{})["current_point"].(float64)
		//	順位
		m := 0.0
		//	ポイント差
		n := 0.0

		Point = int(l)
		Rank = int(m)
		Gap = int(n)

		//	イベント名
		EventID, _ = result.(map[string]interface{})["event"].(map[string]interface{})["event_url"].(string)
		EventID = strings.Replace(EventID, "https://www.showroom-live.com/event/", "", -1)

	} else {
		//	上記ランキングイベントでもレベルイベントでもない場合
		//	もしこのようなケースが存在するならJSONを確認して新たにコーディングする
		log.Println(" N/A")
		return 0, 0, 0, "N/A"
	}

	return
}

/*
 */
func GetIsOnliveByAPI(room_id string) (
	isonlive bool, //	true:	配信中
	startedat time.Time, //	配信開始時刻（isonliveがtrueのときだけ意味があります）
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/profile?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	//	配信中か？
	isonlive, _ = result.(map[string]interface{})["is_onlive"].(bool)

	if isonlive {
		//	配信開始時刻の取得
		value, _ := result.(map[string]interface{})["current_live_started_at"].(float64)
		startedat = time.Unix(int64(value), 0).Truncate(time.Second)
		//	log.Printf("current_live_stared_at %f %v\n", value, startedat)
	}

	return

}

func GetAciveFanByAPI(room_id string, yyyymm string) (nofan int) {

	nofan = -1

	url := "https://www.showroom-live.com/api/active_fan/room?room_id=" + room_id + "&ym=" + yyyymm

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		nofan = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		nofan = -2
		return
	}

	value, _ := result.(map[string]interface{})["total_user_count"].(float64)
	nofan = int(value)

	return
}
func GetRoomInfoByAPI(room_id string) (
	genre string,
	rank string,
	nrank string,
	prank string,
	level int,
	followers int,
	fans int,
	fans_lst int,
	roomname string,
	roomurlkey string,
	startedat time.Time,
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/profile?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	value, _ := result.(map[string]interface{})["follower_num"].(float64)
	followers = int(value)

	tnow := time.Now()
	fans = GetAciveFanByAPI(room_id, tnow.Format("200601"))
	yy := tnow.Year()
	mm := tnow.Month() - 1
	if mm < 0 {
		yy -= 1
		mm = 12
	}
	fans_lst = GetAciveFanByAPI(room_id, fmt.Sprintf("%04d%02d", yy, mm))

	genre, _ = result.(map[string]interface{})["genre_name"].(string)

	rank, _ = result.(map[string]interface{})["league_label"].(string)
	ranks, _ := result.(map[string]interface{})["show_rank_subdivided"].(string)
	rank = rank + " | " + ranks

	value, _ = result.(map[string]interface{})["next_score"].(float64)
	nrank = humanize.Comma(int64(value))
	value, _ = result.(map[string]interface{})["prev_score"].(float64)
	prank = humanize.Comma(int64(value))

	value, _ = result.(map[string]interface{})["room_level"].(float64)
	level = int(value)

	roomname, _ = result.(map[string]interface{})["room_name"].(string)

	roomurlkey, _ = result.(map[string]interface{})["room_url_key"].(string)

	//	配信開始時刻の取得
	value, _ = result.(map[string]interface{})["current_live_started_at"].(float64)
	startedat = time.Unix(int64(value), 0).Truncate(time.Second)
	//	log.Printf("current_live_stared_at %f %v\n", value, startedat)

	return

}

func GetNextliveByAPI(room_id string) (
	nextlive string,
	status int,
) {

	status = 0

	//	https://qiita.com/takeru7584/items/f4ba4c31551204279ed2
	url := "https://www.showroom-live.com/api/room/next_live?room_id=" + room_id

	resp, err := http.Get(url)
	if err != nil {
		//	一時的にデータが取得できない。
		//	resp.Body.Close()
		//		panic(err)
		status = -1
		return
	}
	defer resp.Body.Close()

	//	JSONをデコードする。
	//	次の記事を参考にさせていただいております。
	//		Go言語でJSONに泣かないためのコーディングパターン
	//		https://qiita.com/msh5/items/dc524e38073ed8e3831b

	var result interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		//	panic(err)
		status = -2
		return
	}

	nextlive, _ = result.(map[string]interface{})["text"].(string)

	return

}

func SelectRoomInf(
	userno int,
) (
	roominf RoomInfo,
	status int,
) {

	status = 0

	sql := "select distinct u.userno, userid, user_name, longname, shortname, genre, nrank, prank, level, followers, fans, fans_lst, e.istarget,e.graph, e.color, e.iscntrbpoints, e.point "
	sql += " from user u join eventuser e "
	//	sql += " where u.userno = e.userno and u.userno = " + fmt.Sprintf("%d", userno)
	sql += " where u.userno = e.userno and u.userno = ? "

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("SelectRoomInf() Prepare() err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(userno).Scan(&roominf.Userno,
		&roominf.Account,
		&roominf.Name,
		&roominf.Longname,
		&roominf.Shortname,
		&roominf.Genre,
		&roominf.Nrank,
		&roominf.Prank,
		&roominf.Level,
		&roominf.Followers,
		&roominf.Fans,
		&roominf.Fans_lst,
		&roominf.Istarget,
		&roominf.Graph,
		&roominf.Color,
		&roominf.Iscntrbpoint,
		&roominf.Point,
	)
	if err != nil {
		log.Printf("SelectRoomInf() Query() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	if roominf.Istarget == "Y" {
		roominf.Istarget = "Checked"
	} else {
		roominf.Istarget = ""
	}
	if roominf.Graph == "Y" {
		roominf.Graph = "Checked"
	} else {
		roominf.Graph = ""
	}
	if roominf.Iscntrbpoint == "Y" {
		roominf.Iscntrbpoint = "Checked"
	} else {
		roominf.Iscntrbpoint = ""
	}
	roominf.Slevel = humanize.Comma(int64(roominf.Level))
	roominf.Sfollowers = humanize.Comma(int64(roominf.Followers))
	roominf.Spoint = humanize.Comma(int64(roominf.Point))
	roominf.Name = strings.ReplaceAll(roominf.Name, "'", "’")

	return
}

func SelectEventRoomInfList(
	eventid string,
	roominfolist *RoomInfoList,
) (
	eventname string,
	status int,
) {

	status = 0

	//	eventno := 0
	//	eventno, eventname, _ = SelectEventNoAndName(eventid)
	//	Event_inf, _ = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	eventno := Event_inf.Event_no
	eventname = Event_inf.Event_name

	sql := "select distinct u.userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, e.istarget, e.graph, e.color, e.iscntrbpoints, e.point "
	sql += " from user u join eventuser e "
	sql += " where u.userno = e.userno and e.eventid= ?"
	if Event_inf.Start_time.After(time.Now()) {
		sql += " order by followers desc"
	} else {
		sql += " order by e.point desc"
	}

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("SelectEventRoomInfList() Prepare() err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(eventid)
	if err != nil {
		log.Printf("SelectRoomIn() Query() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	defer rows.Close()

	ColorlistA := Colorlist2
	ColorlistB := Colorlist1
	if Event_inf.Cmap == 1 {
		ColorlistA = Colorlist1
		ColorlistB = Colorlist2
	}

	colormap := make(map[string]int)

	for i := 0; i < len(ColorlistA); i++ {
		colormap[ColorlistA[i].Name] = i
	}

	var roominf RoomInfo

	i := 0
	for rows.Next() {
		err := rows.Scan(&roominf.Userno,
			&roominf.Account,
			&roominf.Name,
			&roominf.Longname,
			&roominf.Shortname,
			&roominf.Genre,
			&roominf.Rank,
			&roominf.Nrank,
			&roominf.Prank,
			&roominf.Level,
			&roominf.Followers,
			&roominf.Fans,
			&roominf.Fans_lst,
			&roominf.Istarget,
			&roominf.Graph,
			&roominf.Color,
			&roominf.Iscntrbpoint,
			&roominf.Point,
		)

		ci := 0
		for ; ci < len(ColorlistA); ci++ {
			if ColorlistA[ci].Name == roominf.Color {
				roominf.Colorvalue = ColorlistA[ci].Value
				break
			}
		}
		if ci == len(ColorlistA) {
			ci := 0
			for ; ci < len(ColorlistB); ci++ {
				if ColorlistB[ci].Name == roominf.Color {
					roominf.Colorvalue = ColorlistB[ci].Value
					break
				}
			}
			if ci == len(ColorlistB) {
				roominf.Colorvalue = roominf.Color
			}
		}

		if roominf.Istarget == "Y" {
			roominf.Istarget = "Checked"
		} else {
			roominf.Istarget = ""
		}
		if roominf.Graph == "Y" {
			roominf.Graph = "Checked"
		} else {
			roominf.Graph = ""
		}
		if roominf.Iscntrbpoint == "Y" {
			roominf.Iscntrbpoint = "Checked"
		} else {
			roominf.Iscntrbpoint = ""
		}
		roominf.Slevel = humanize.Comma(int64(roominf.Level))
		roominf.Sfollowers = humanize.Comma(int64(roominf.Followers))
		if roominf.Point < 0 {
			roominf.Spoint = ""
		} else {
			roominf.Spoint = humanize.Comma(int64(roominf.Point))
		}
		roominf.Formid = "Form" + fmt.Sprintf("%d", i)
		roominf.Eventid = eventid
		roominf.Name = strings.ReplaceAll(roominf.Name, "'", "’")
		if err != nil {
			log.Printf("SelectEventRoomInfList() Scan() err=%s\n", err.Error())
			status = -7
			return
		}
		//	var colorinf ColorInf
		colorinflist := make([]ColorInf, len(ColorlistA))

		for i := 0; i < len(ColorlistA); i++ {
			colorinflist[i].Color = ColorlistA[i].Name
			colorinflist[i].Colorvalue = ColorlistA[i].Value
		}

		roominf.Colorinflist = colorinflist
		if cidx, ok := colormap[roominf.Color]; ok {
			roominf.Colorinflist[cidx].Selected = "Selected"
		}
		*roominfolist = append(*roominfolist, roominf)

		i++
	}

	if err = rows.Err(); err != nil {
		log.Printf("SelectEventRoomInfList() rows err=%s\n", err.Error())
		status = -8
		return
	}

	if Event_inf.Start_time.After(time.Now()) {
		SortByFollowers = true
	} else {
		SortByFollowers = false
	}
	sort.Sort(*roominfolist)

	/*
		for i := 0; i < len(*roominfolist); i++ {

			sql = "select max(point) from points where "
			sql += " user_id = " + fmt.Sprintf("%d", (*roominfolist)[i].Userno)
			//	sql += " and event_id = " + fmt.Sprintf("%d", eventno)
			sql += " and event_id = " + eventid

			err = Db.QueryRow(sql).Scan(&(*roominfolist)[i].Point)
			(*roominfolist)[i].Spoint = humanize.Comma(int64((*roominfolist)[i].Point))

			if err == nil {
				continue
			} else {
				log.Printf("err=[%s]\n", err.Error())
				if err.Error() != "sql: no rows in result set" {
					eventno = -2
					continue
				} else {
					(*roominfolist)[i].Point = -1
					(*roominfolist)[i].Spoint = ""
				}
			}
		}
	*/

	return
}

func UpdateRoomInf(eventid, suserno, longname, shortname, istarget, graph, color, iscntrbpoint string) (status int) {

	status = 0

	userno, _ := strconv.Atoi(suserno)

	sql := "update user set longname=?, shortname=? where userno = ?"

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(longname, shortname, userno)

	if err != nil {
		log.Printf("UpdateRoomInf() error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		status = -2
		return
	}

	//	eventno, _, _ := SelectEventNoAndName(eventid)

	if istarget == "1" {
		istarget = "Y"
	} else {
		istarget = "N"
	}

	if graph == "1" {
		graph = "Y"
	} else {
		graph = "N"
	}

	if iscntrbpoint == "1" {
		iscntrbpoint = "Y"
	} else {
		iscntrbpoint = "N"
	}

	//	sql = "update eventuser set istarget=?, graph=?, color=? where eventno=? and userno=?"
	sql = "update eventuser set istarget=?, graph=?, color=?, iscntrbpoints=? where eventid=? and userno=?"

	stmt, err = srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateRoomInf() error(Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(istarget, graph, color, iscntrbpoint, eventid, userno)

	if err != nil {
		log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return

}

func UpdateEventuserSetPoint(eventid, userid string, point int) (status int) {
	status = 0

	//	eventno, _, _ := SelectEventNoAndName(eventid)
	userno, _ := strconv.Atoi(userid)

	sql := "update eventuser set point=? where eventid = ? and userno = ?"
	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdateEventuserSetPoint() error (Update/Prepare) err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(point, eventid, userno)

	if err != nil {
		log.Printf("error(UpdateEventuserSetPoint() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return
}
func GetRoomInfoAndPoint(
	eventid string,
	roomid string,
	idbasis string,
) (
	roominf RoomInfo,
	status int,
) {

	status = 0

	roominf.ID = roomid
	roominf.Userno, _ = strconv.Atoi(roomid)

	//	Event_inf, _ = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	roominf.Genre, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers,
		roominf.Fans,
		roominf.Fans_lst,
		roominf.Name, roominf.Account, _, status =
		GetRoomInfoByAPI(roomid)

	point, _, _, peventid := GetPointsByAPI(roominf.ID)
	if peventid == Event_inf.Event_ID {
		roominf.Point = point
		UpdateEventuserSetPoint(peventid, roominf.ID, point)
	} else {
		log.Printf(" %s %s %d\n", Event_inf.Event_ID, peventid, point)
	}

	/*
		if (*roominfolist)[i].ID == idbasis {
			(*eventinfo).Pntbasis = point
			(*eventinfo).Ordbasis = i
		}
	*/

	//	log.Printf(" followers=<%d> level=<%d> nrank=<%s> genre=<%s> point=%d\n",
	//	(*roominfolist)[i].Followers,
	//	(*roominfolist)[i].Level,
	//	(*roominfolist)[i].Nrank,
	//	(*roominfolist)[i].Genre,
	//	(*roominfolist)[i].Point)

	return
}

func GetAndInsertEventRoomInfo(
	client *http.Client,
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	starttimeafternow bool,
	status int,
) {

	log.Println("GetAndInsertEventRoomInfo() Called.")
	log.Println(*eventinfo)

	status = 0
	starttimeafternow = false

	//	イベントに参加しているルームの一覧を取得します。
	//	ルーム名、ID、URLを取得しますが、イベント終了直後の場合の最終獲得ポイントが表示されている場合はそれも取得します。

	if strings.Contains(eventid, "?") {
		status = GetEventInfAndRoomListBR(client, eventid, breg, ereg, eventinfo, roominfolist)
		eia := strings.Split(eventid, "?")
		bka := strings.Split(eia[1], "=")
		eventinfo.Event_name = eventinfo.Event_name + "(" + bka[1] + ")"
	} else {
		status = GetEventInfAndRoomList(eventid, breg, ereg, eventinfo, roominfolist)
	}

	if status != 0 {
		log.Printf("GetEventInfAndRoomList() returned %d\n", status)
		return
	}

	//	各ルームのジャンル、ランク、レベル、フォロワー数を取得します。
	for i := 0; i < (*eventinfo).NoRoom; i++ {
		(*roominfolist)[i].Genre, (*roominfolist)[i].Rank,
			(*roominfolist)[i].Nrank,
			(*roominfolist)[i].Prank,
			(*roominfolist)[i].Level,
			(*roominfolist)[i].Followers,
			(*roominfolist)[i].Fans,
			(*roominfolist)[i].Fans_lst,
			_, _, _, _ = GetRoomInfoByAPI((*roominfolist)[i].ID)

	}

	//	各ルームの獲得ポイントを取得します。
	for i := 0; i < (*eventinfo).NoRoom; i++ {
		point, _, _, eventid := GetPointsByAPI((*roominfolist)[i].ID)
		if eventid == (*eventinfo).Event_ID {
			(*roominfolist)[i].Point = point
			UpdateEventuserSetPoint(eventid, (*roominfolist)[i].ID, point)
			if point < 0 {
				(*roominfolist)[i].Spoint = ""
			} else {
				(*roominfolist)[i].Spoint = humanize.Comma(int64(point))
			}
		} else {
			log.Printf(" %s %s %d\n", (*eventinfo).Event_ID, eventid, point)
		}

		if (*roominfolist)[i].ID == fmt.Sprintf("%d", (*eventinfo).Nobasis) {
			(*eventinfo).Pntbasis = point
			(*eventinfo).Ordbasis = i
		}

		//	log.Printf(" followers=<%d> level=<%d> nrank=<%s> genre=<%s> point=%d\n",
		//	(*roominfolist)[i].Followers,
		//	(*roominfolist)[i].Level,
		//	(*roominfolist)[i].Nrank,
		//	(*roominfolist)[i].Genre,
		//	(*roominfolist)[i].Point)
	}

	if (*eventinfo).Start_time.After(time.Now()) {
		SortByFollowers = true
		sort.Sort(*roominfolist)
		if ereg > len(*roominfolist) {
			ereg = len(*roominfolist)
		}
		r := (*roominfolist).Choose(breg-1, ereg)
		roominfolist = &r
		starttimeafternow = true
	}

	log.Printf(" GetEventRoomInfo() len(*roominfolist)=%d\n", len(*roominfolist))

	log.Println("GetAndInsertEventRoomInfo() before InsertEventIinf()")
	log.Println(*eventinfo)
	status = InsertEventInf(eventinfo)

	if status == 1 {
		log.Println("InsertEventInf() returned 1.")
		UpdateEventInf(eventinfo)
		status = 0
	}
	log.Println("GetAndInsertEventRoomInfo() after InsertEventIinf() or UpdateEventInf")
	log.Println(*eventinfo)

	_, _, status = SelectEventNoAndName(eventid)

	if status == 0 {
		//	InsertRoomInf(eventno, eventid, roominfolist)
		InsertRoomInf(eventid, roominfolist)
	}

	return
}

func InsertEventInf(eventinf *exsrapi.Event_Inf) (
	status int,
) {

	if _, _, status = SelectEventNoAndName((*eventinf).Event_ID); status != 0 {
		sql := "INSERT INTO event(eventid, ieventid, event_name, period, starttime, endtime, noentry,"
		sql += " intervalmin, modmin, modsec, "
		sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, Cmap, target, maxpoint "
		sql += ") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		log.Printf("db.Prepare(sql)\n")
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error InsertEventInf() (INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		if eventinf.Intervalmin != 5 {	//	緊急対応
			log.Printf(" Intervalmin isn't 5. (%dm)\n",eventinf.Intervalmin)
			eventinf.Intervalmin = 5
		}

		log.Printf("row.Exec()\n")
		_, err = stmt.Exec(
			(*eventinf).Event_ID,
			(*eventinf).I_Event_ID,
			(*eventinf).Event_name,
			(*eventinf).Period,
			(*eventinf).Start_time,
			(*eventinf).End_time,
			(*eventinf).NoEntry,
			(*eventinf).Intervalmin,
			(*eventinf).Modmin,
			(*eventinf).Modsec,
			(*eventinf).Fromorder,
			(*eventinf).Toorder,
			(*eventinf).Resethh,
			(*eventinf).Resetmm,
			(*eventinf).Nobasis,
			(*eventinf).Maxdsp,
			(*eventinf).Cmap,
			(*eventinf).Target,
			(*eventinf).Maxpoint+eventinf.Gscale,
		)

		if err != nil {
			log.Printf("error InsertEventInf() (INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
	} else {
		status = 1
	}

	return
}

func UpdateEventInf(eventinf *exsrapi.Event_Inf) (
	status int,
) {

	if _, _, status = SelectEventNoAndName((*eventinf).Event_ID); status == 0 {
		sql := "Update event set "
		sql += " ieventid=?,"
		sql += " event_name=?,"
		sql += " period=?,"
		sql += " starttime=?,"
		sql += " endtime=?,"
		sql += " noentry=?,"
		sql += " intervalmin=?,"
		sql += " modmin=?,"
		sql += " modsec=?,"
		sql += " Fromorder=?,"
		sql += " Toorder=?,"
		sql += " Resethh=?,"
		sql += " Resetmm=?,"
		sql += " Nobasis=?,"
		sql += " Target=?,"
		sql += " Maxdsp=?, "
		sql += " cmap=?, "
		sql += " maxpoint=? "
		//	sql += " where eventno = ?"
		sql += " where eventid = ?"
		log.Printf("db.Prepare(sql)\n")

		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("UpdateEventInf() error (Update/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		if eventinf.Intervalmin != 5 {	//	緊急対応
			log.Printf(" Intervalmin isn't 5. (%dm)\n",eventinf.Intervalmin)
			eventinf.Intervalmin = 5
		}

		log.Printf("row.Exec()\n")
		_, err = stmt.Exec(
			(*eventinf).I_Event_ID,
			(*eventinf).Event_name,
			(*eventinf).Period,
			(*eventinf).Start_time,
			(*eventinf).End_time,
			(*eventinf).NoEntry,
			(*eventinf).Intervalmin,
			(*eventinf).Modmin,
			(*eventinf).Modsec,
			(*eventinf).Fromorder,
			(*eventinf).Toorder,
			(*eventinf).Resethh,
			(*eventinf).Resetmm,
			(*eventinf).Nobasis,
			(*eventinf).Target,
			(*eventinf).Maxdsp,
			(*eventinf).Cmap,
			(*eventinf).Maxpoint+eventinf.Gscale,
			(*eventinf).Event_ID,
		)

		if err != nil {
			log.Printf("error UpdateEventInf() (update/Exec) err=%s\n", err.Error())
			status = -2
		}
	} else {
		status = 1
	}

	return
}

func InsertRoomInf(eventid string, roominfolist *RoomInfoList) {

	log.Printf("***** InsertRoomInf() ***********  NoRoom=%d\n", len(*roominfolist))
	tnow := time.Now().Truncate(time.Second)
	for i := 0; i < len(*roominfolist); i++ {
		log.Printf("   ** InsertRoomInf() ***********  i=%d\n", i)
		InsertIntoOrUpdateUser(tnow, eventid, (*roominfolist)[i])
		status := InsertIntoEventUser(i, eventid, (*roominfolist)[i])
		if status == 0 {
			(*roominfolist)[i].Status = "更新"
			(*roominfolist)[i].Statuscolor = "black"
		} else if status == 1 {
			(*roominfolist)[i].Status = "新規"
			(*roominfolist)[i].Statuscolor = "green"
		} else {
			(*roominfolist)[i].Status = "エラー"
			(*roominfolist)[i].Statuscolor = "red"
		}
	}
	log.Printf("***** end of InsertRoomInf() ***********\n")
}

func InsertIntoOrUpdateUser(tnow time.Time, eventid string, roominf RoomInfo) (status int) {

	status = 0

	isnew := false

	userno, _ := strconv.Atoi(roominf.ID)
	log.Printf("  *** InsertIntoOrUpdateUser() *** userno=%d\n", userno)

	nrow := 0
	err := srdblib.Db.QueryRow("select count(*) from user where userno =" + roominf.ID).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	name := ""
	genre := ""
	rank := ""
	nrank := ""
	prank := ""
	level := 0
	followers := 0
	fans := -1
	fans_lst := -1

	if nrow == 0 {

		isnew = true

		log.Printf("insert into userhistory(*new*) userno=%d rank=<%s> nrank=<%s> prank=<%s> level=%d, followers=%d, fans=%d, fans_lst=%d\n",
			userno, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers, fans, fans_lst)

		sql := "INSERT INTO user(userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts, currentevent)"
		sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

		//	log.Printf("sql=%s\n", sql)
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("InsertIntoOrUpdateUser() error() (INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		lenid := len(roominf.ID)
		_, err = stmt.Exec(
			userno,
			roominf.Account,
			roominf.Name,
			//	roominf.ID,
			roominf.Name,
			roominf.ID[lenid-2:lenid],
			roominf.Genre,
			roominf.Rank,
			roominf.Nrank,
			roominf.Prank,
			roominf.Level,
			roominf.Followers,
			roominf.Fans,
			roominf.Fans_lst,
			tnow,
			eventid,
		)

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = stmt.Exec(
				userno,
				roominf.Account,
				roominf.Account,
				roominf.ID,
				roominf.ID[lenid-2:lenid],
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
				eventid,
			)
			if err != nil {
				log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
				status = -2
			}
		}
	} else {

		sql := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from user where userno = ?"
		err = srdblib.Db.QueryRow(sql, userno).Scan(&name, &genre, &rank, &nrank, &prank, &level, &followers, &fans, &fans_lst)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
		}
		//	log.Printf("current userno=%d name=%s, nrank=%s, prank=%s level=%d, followers=%d\n", userno, name, nrank, prank, level, followers)

		if roominf.Genre != genre ||
			roominf.Rank != rank ||
			//	roominf.Nrank != nrank ||
			//	roominf.Prank != prank ||
			roominf.Level != level ||
			roominf.Followers != followers ||
			roominf.Fans != fans {

			isnew = true

			log.Printf("insert into userhistory(*changed*) userno=%d level=%d, followers=%d, fans=%d\n",
				userno, roominf.Level, roominf.Followers, roominf.Fans)
			sql := "update user set userid=?,"
			sql += "user_name=?,"
			sql += "genre=?,"
			sql += "`rank`=?,"
			sql += "nrank=?,"
			sql += "prank=?,"
			sql += "level=?,"
			sql += "followers=?,"
			sql += "fans=?,"
			sql += "fans_lst=?,"
			sql += "ts=?,"
			sql += "currentevent=? "
			sql += "where userno=?"
			stmt, err := srdblib.Db.Prepare(sql)

			if err != nil {
				log.Printf("InsertIntoOrUpdateUser() error(Update/Prepare) err=%s\n", err.Error())
				status = -1
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(
				roominf.Account,
				roominf.Name,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
				eventid,
				roominf.ID,
			)

			if err != nil {
				log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
				status = -2
			}
		}
		/* else {
			//	log.Printf("not insert into userhistory(*same*) userno=%d level=%d, followers=%d\n", userno, roominf.Level, roominf.Followers)
		}
		*/

	}

	if isnew {
		sql := "INSERT INTO userhistory(userno, user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts)"
		sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?)"
		//	log.Printf("sql=%s\n", sql)
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT into userhistory/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			userno,
			roominf.Name,
			roominf.Genre,
			roominf.Rank,
			roominf.Nrank,
			roominf.Prank,
			roominf.Level,
			roominf.Followers,
			roominf.Fans,
			roominf.Fans_lst,
			tnow,
		)

		if err != nil {
			log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = stmt.Exec(
				userno,
				roominf.Account,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
			)
			if err != nil {
				log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
				status = -2
			}
		}

	}

	return

}
func InsertIntoEventUser(i int, eventid string, roominf RoomInfo) (status int) {

	status = 0

	userno, _ := strconv.Atoi(roominf.ID)

	nrow := 0
	/*
		sql := "select count(*) from eventuser where "
		sql += "userno =" + roominf.ID + " and "
		//	sql += "eventno = " + fmt.Sprintf("%d", eventno)
		sql += "eventid = " + eventid
		//	log.Printf("sql=%s\n", sql)
		err := Db.QueryRow(sql).Scan(&nrow)
	*/
	sql := "select count(*) from eventuser where userno =? and eventid = ?"
	err := srdblib.Db.QueryRow(sql, roominf.ID, eventid).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	Colorlist := Colorlist2
	if Event_inf.Cmap == 1 {
		Colorlist = Colorlist1
	}

	if nrow == 0 {
		sql := "INSERT INTO eventuser(eventid, userno, istarget, graph, color, iscntrbpoints, point) VALUES(?,?,?,?,?,?,?)"
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		//	if i < 10 {
		_, err = stmt.Exec(
			eventid,
			userno,
			"Y",
			"Y",
			Colorlist[i%len(Colorlist)].Name,
			"N",
			roominf.Point,
		)
		/*
			} else {
				_, err = stmt.Exec(
					eventid,
					userno,
					"Y",	//	"N"から変更する＝順位に関わらず獲得ポイントデータを取得する。
					"N",
					Colorlist[i%len(Colorlist)].Name,
					"N",
					roominf.Point,
				)
			}
		*/

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
		status = 1
	}
	return

}

func GetEventInfAndRoomList(
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
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
	maxroom := ereg

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
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}

		content, _ := doc.Find("head > meta:nth-child(6)").Attr("content")
		content_div := strings.Split(content, "/")
		(*eventinfo).Event_ID = content_div[len(content_div)-1]

	} else {
		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventidorfilename
		/*
			doc, err = goquery.NewDocument(_url)
		*/
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
			status = 1
			return
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		//	bufstr := buf.String()
		//	log.Printf("%s\n", bufstr)

		//	doc, error = goquery.NewDocumentFromReader(resp.Body)
		doc, error = goquery.NewDocumentFromReader(buf)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			status = 1
			return
		}

		(*eventinfo).Event_ID = eventidorfilename
	}
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	cevent_id, exists := doc.Find("#eventDetail").Attr("data-event-id")
	if !exists {
		log.Printf("data-event-id not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	eventinfo.I_Event_ID, _ = strconv.Atoi(cevent_id)

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
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
	//	参加ルーム数と表示されているルームの数は違うので、ここで取得したルームの数を以下の処理で使うわけではありません。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)
	log.Printf("breg=%d ereg=%d\n", breg, ereg)

	//	eventno, _, _ := SelectEventNoAndName(eventidorfilename)
	//	log.Printf(" eventno=%d\n", eventno)
	//	(*eventinfo).Event_no = eventno

	//	抽出したルームすべてに対して処理を繰り返す(が、イベント開始後の場合の処理はルーム数がbreg、eregの範囲に限定）
	//	イベント開始前のときはすべて取得し、ソートしたあてで範囲を限定する）
	doc.Find(".listcardinfo").EachWithBreak(func(i int, s *goquery.Selection) bool {
		//	log.Printf("i=%d\n", i)
		if (*eventinfo).Start_time.Before(time.Now()) {
			if i < breg-1 {
				return true
			}
			if i == maxroom {
				return false
			}
		}

		var roominfo RoomInfo

		roominfo.Name = s.Find(".listcardinfo-main-text").Text()

		spoint1 := strings.Split(s.Find(".listcardinfo-sub-single-right-text").Text(), ": ")

		var point int64
		if spoint1[0] != "" {
			spoint2 := strings.Split(spoint1[1], "pt")
			fmt.Sscanf(spoint2[0], "%d", &point)

		} else {
			point = -1
		}
		roominfo.Point = int(point)

		ReplaceString := ""

		selection_c := s.Find(".listcardinfo-menu")

		account, _ := selection_c.Find(".room-url").Attr("href")
		if inputmode == "file" {
			ReplaceString = "https://www.showroom-live.com/"
		} else {
			ReplaceString = "/r/"
		}
		roominfo.Account = strings.Replace(account, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, ReplaceString, "/", -1)

		roominfo.ID, _ = selection_c.Find(".js-follow-btn").Attr("data-room-id")
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		*roominfolist = append(*roominfolist, roominfo)

		//	log.Printf("%11s %-20s %-10s %s\n",
		//		humanize.Comma(int64(roominfo.Point)), roominfo.Account, roominfo.ID, roominfo.Name)
		return true

	})

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}

func GetEventInfAndRoomListBR(
	client *http.Client,
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	status int,
) {

	status = 0

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid

	status = 0

	//	URLからドキュメントを作成します
	_url := "https://www.showroom-live.com/event/" + eventidorfilename
	/*
		doc, err = goquery.NewDocument(_url)
	*/
	resp, error := http.Get(_url)
	if error != nil {
		log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
		status = 1
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bufstr := buf.String()

	log.Printf("%s\n", bufstr)

	//	doc, error = goquery.NewDocumentFromReader(resp.Body)
	doc, error = goquery.NewDocumentFromReader(buf)
	if error != nil {
		log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
		status = 1
		return
	}

	(*eventinfo).Event_ID = eventidorfilename
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	cevent_id, exists := doc.Find("#eventDetail").Attr("data-event-id")
	if !exists {
		log.Printf("data-event-id not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	eventinfo.I_Event_ID, _ = strconv.Atoi(cevent_id)
	event_id := eventinfo.I_Event_ID

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
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
	//	参加ルーム数と表示されているルームの数は違うので、ここで取得したルームの数を以下の処理で使うわけではありません。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)
	log.Printf("breg=%d ereg=%d\n", breg, ereg)

	//	eventno, _, _ := SelectEventNoAndName(eventidorfilename)
	//	log.Printf(" eventno=%d\n", eventno)
	//	(*eventinfo).Event_no = eventno

	eia := strings.Split(eventid, "?")
	bia := strings.Split(eia[1], "=")
	blockid, _ := strconv.Atoi(bia[1])

	/*
		event_id := 30030
		event_id := 31947
	*/

	ebr, err := srapi.GetEventBlockRanking(client, event_id, blockid, breg, ereg)
	if err != nil {
		log.Printf("GetEventBlockRanking() err=%s\n", err.Error())
		status = 1
		return
	}

	ReplaceString := "/r/"

	for _, br := range ebr.Block_ranking_list {

		var roominfo RoomInfo

		roominfo.ID = br.Room_id
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		roominfo.Account = strings.Replace(br.Room_url_key, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, "/", "", -1)

		roominfo.Name = br.Room_name

		*roominfolist = append(*roominfolist, roominfo)

	}

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}

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

	/*
		_, _, status := SelectEventNoAndName(eventidorfilename)
		log.Printf(" status=%d\n", status)
		if status != 0 {
			return
		}
		(*eventinfo).Event_no = eventno
	*/

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
		/*
			doc, err = goquery.NewDocument(_url)
		*/
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

func SelectEventNoAndName(eventid string) (
	eventname string,
	period string,
	status int,
) {

	status = 0

	err := srdblib.Db.QueryRow("select event_name, period from event where eventid ='"+eventid+"'").Scan(&eventname, &period)

	if err == nil {
		return
	} else {
		log.Printf("err=[%s]\n", err.Error())
		if err.Error() != "sql: no rows in result set" {
			status = -2
			return
		}
	}

	status = -1
	return
}

func SelectUserName(userno int) (
	longname string,
	shortname string,
	genre string,
	rank string,
	nrank string,
	prank string,
	level int,
	followers int,
	fans int,
	fans_lst int,
	status int,
) {

	status = 0

	sql := "select longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from user where userno = ?"

	err := srdblib.Db.QueryRow(sql, userno).Scan(&longname, &shortname, &genre, &rank, &nrank, &prank, &level, &followers, &fans, &fans_lst)

	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
	}

	return
}

func SelectUserColor(userno int, eventid string) (
	color string,
	colorvalue string,
	status int,
) {

	Colorlist := Colorlist2
	if Event_inf.Cmap == 1 {
		Colorlist = Colorlist1
	}

	status = 0

	//	sql := "select color from eventuser where userno = ? and eventno = ?"
	sql := "select color from eventuser where userno = ? and eventid = ?"

	err := srdblib.Db.QueryRow(sql, userno, eventid).Scan(&color)

	i := 0
	for ; i < len(Colorlist); i++ {
		if Colorlist[i].Name == color {
			colorvalue = Colorlist[i].Value
			break
		}
	}
	if i == len(Colorlist) {
		colorvalue = color
	}

	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
	}

	return
}

func SelectRoomLevel(userno int, levelonly int) (roomlevelinf RoomLevelInf, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	status = 0

	sqlstmt := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts from userhistory where userno = ? order by ts desc"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sqlstmt)
	if srdblib.Dberr != nil {
		log.Printf("SelectRoomLevel() (3) err=%s\n", srdblib.Dberr.Error())
		status = -3
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(userno)
	if srdblib.Dberr != nil {
		log.Printf("SelectRoomLevel() (6) err=%s\n", srdblib.Dberr.Error())
		status = -6
		return
	}
	defer rows.Close()

	/*
	   type RoomLevel struct {
	   	User_name  string
	   	Genre      string
	   	Rank       string
	   	Nrank       string
	   	Level      int
	   	Followeres int
	   	Sts        string
	   }

	   type RoomLevelInf struct {
	   	Userno        int
	   	User_name      string
	   	RoomLevelList []RoomLevel
	   }
	*/

	var roomlevel RoomLevel

	roomlevelinf.Userno = userno

	lastlevel := 0

	for rows.Next() {
		srdblib.Dberr = rows.Scan(&roomlevel.User_name, &roomlevel.Genre, &roomlevel.Rank,
			&roomlevel.Nrank,
			&roomlevel.Prank,
			&roomlevel.Level,
			&roomlevel.Followers,
			&roomlevel.Fans,
			&roomlevel.Fans_lst,
			&roomlevel.ts)
		if srdblib.Dberr != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", srdblib.Dberr.Error())
			status = -7
			return
		}

		if lastlevel == 0 {
			roomlevelinf.User_name = roomlevel.User_name
		}

		if levelonly == 1 && roomlevel.Level == lastlevel {
			continue
		}
		lastlevel = roomlevel.Level

		//	roomlevel.Sfollowers = humanize.Comma(int64(roomlevel.Followers))
		roomlevel.Sts = roomlevel.ts.Format("2006/01/02 15:04")

		roomlevelinf.RoomLevelList = append(roomlevelinf.RoomLevelList, roomlevel)

	}

	return
}

func SelectCurrentScore(eventid string) (gtime time.Time, eventname string, period string, scorelist []CurrentScore, status int) {

	status = 0

	//	Event_inf, status = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	eventno = Event_inf.Event_no
	eventname = Event_inf.Event_name
	period = Event_inf.Period

	nrow := 0
	sql := "select count(*) from points where eventid = ?"
	srdblib.Dberr = srdblib.Db.QueryRow(sql, eventid).Scan(&nrow)

	if srdblib.Dberr != nil {
		log.Printf("select max(point) from eventuser where eventid = '%s'\n", Event_inf.Event_ID)
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -11
		return
	}
	if nrow == 0 {
		log.Printf("no data in points(where eventid=%s).\n", eventid)
		status = -12
		return
	}

	//	---------------------------------------------------
	//	sql := "select t.idx, t.t from timeacq t join points p where t.idx = p.idx and t.idx = ( select max(idx) from points where event_id = ? )"
	//	sql := "select distinct t.idx, t.t from timeacq t join points p where t.idx = p.idx and t.t = ( select max(t) from points p join timeacq t where p.idx = t.idx and event_id = ? )"
	sql = "select distinct max(ts) from points where eventid = ?"
	//	sql := "select distinct COALESCE(max(ts), ?) from points where eventid = ?"
	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("GetCurrentScore() (3) err=%s\n", err.Error())
		status = -3
		return
	}
	defer stmt.Close()

	//	idx := 0
	//	Err = stmt.QueryRow(time.Now().Add(time.Hour), eventid).Scan(&gtime)
	srdblib.Dberr = stmt.QueryRow(eventid).Scan(&gtime)
	if srdblib.Dberr != nil {
		log.Printf("GetCurrentScore() (4) err=%s\n", srdblib.Dberr.Error())
		status = -4
		return
	}
	log.Printf("gtime=%s\n", gtime.Format("2006/01/02 15:04:06"))
	/*
		if gtime.After(time.Now()) {
			status = -10
			return
		}
	*/

	//	---------------------------------------------------
	//	stmt, err = Db.Prepare("select user_id, `rank`, point, pstatus, ptime, qstatus, qtime from points where eventid = ? and ts = ? order by point desc")
	sql = "select p.user_id, u.userid, p.rank, p.point, p.pstatus, p.ptime, p.qstatus, p.qtime "
	sql += " from points p join user u where p.eventid = ? and p.ts = ? and p.user_id = u.userno order by p.point desc"
	stmt, err = srdblib.Db.Prepare(sql)

	if err != nil {
		log.Printf("GetCurrentScore() (5) err=%s\n", err.Error())
		status = -5
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(eventid, gtime)
	if err != nil {
		log.Printf("GetCurrentScore() (6) err=%s\n", err.Error())
		status = -6
		return
	}
	defer rows.Close()

	//	var score, bscore CurrentScore
	var bscore CurrentScore
	point_bs := 0
	i := 0
	//	shift := 1
	nextrank := 1
	for rows.Next() {
		var score CurrentScore
		err := rows.Scan(&score.Userno, &score.Shorturl, &score.Rank, &score.Point, &score.Pstatus, &score.Ptime, &score.Qstatus, &score.Qtime)
		if err != nil {
			log.Printf("GetCurrentScore() (7) err=%s\n", err.Error())
			status = -7
			return
		}
		if score.Userno == Event_inf.Nobasis {
			point_bs = score.Point
			log.Printf(" Nobasis=%d  point_bs=%d\n", Event_inf.Nobasis, point_bs)
		}
		score.Spoint = humanize.Comma(int64(score.Point))
		username, _, roomgenre, roomrank, roomnrank, roomprank, roomlevel, followers, fans, fans_lst, sts := SelectUserName(score.Userno)
		score.Username = username
		if sts != 0 {
			score.Username = fmt.Sprintf("%d", score.Userno)
		}
		score.Roomgenre = roomgenre
		score.Roomrank = roomrank
		score.Roomnrank = roomnrank
		score.Roomprank = roomprank
		score.Roomlevel = humanize.Comma(int64(roomlevel))
		score.Followers = humanize.Comma(int64(followers))
		score.Fans = fans
		score.Fans_lst = fans_lst

		/*
			nroomlevel := 0
			nfollowers := 0
			score.Roomgenre, score.Roomrank, score.Roomnrank, score.Roomprank, nroomlevel,
				nfollowers, score.Fans, score.Fans_lst, _, _, _, status = GetRoomInfoByAPI(fmt.Sprintf("%d", score.Userno))
			score.Roomlevel = humanize.Comma(int64(nroomlevel))
			score.Followers = humanize.Comma(int64(nfollowers))
			/* */
		/*
			if	score.Roomrank != roomrank ||
				score.Roomnrank != roomnrank ||
				nfollowers != followers ||
				nroomlevel != roomlevel ||
				score.Fans != fans {
				UpdateRoomRankInf (score, nroomlevel, nfollowers)

			}
			/* */

		if score.Rank != 0 {
			score.Srank = fmt.Sprintf("%d", score.Rank)
		} else {
			score.Srank = ""
		}
		//	if score.Rank > i+shift {
		if score.Rank > nextrank {
			//	bscore.Srank = fmt.Sprintf("%d", i+shift)
			bscore.Srank = "-"
			scorelist = append(scorelist, bscore)
			//	shift++
		}
		nextrank = score.Rank + 1

		score.NextLive, _ = GetNextliveByAPI(fmt.Sprintf("%d", score.Userno))
		score.Eventid = eventid

		acqtimelist, _ := SelectAcqTimeList(eventid, score.Userno)
		lenatl := len(acqtimelist)
		log.Printf(" eventid = %s userno = %d len(acqtimelist=%d\n", eventid, score.Userno, lenatl)
		if lenatl != 0 {
			score.Bcntrb = true
		} else {
			score.Bcntrb = false
		}

		scorelist = append(scorelist, score)
		i++
		/*
			if i == 10 {
				break
			}
		*/
	}
	if err = rows.Err(); err != nil {
		log.Printf("GetCurrentScore() (8) err=%s\n", err.Error())
		status = -8
		return
	}

	if point_bs > 0 {
		for i, score := range scorelist {
			if score.Point != 0 {
				scorelist[i].Sdfr = humanize.Comma(int64(score.Point - point_bs))
			}
		}
	}

	return

}

func SelectUserList() (userlist []User, status int) {

	status = 0

	sql := "select distinct(e.nobasis),u.longname "
	sql += " from event e join user u on e.nobasis=u.userno "
	sql += " where e.nobasis != 0 "
	sql += " order by e.nobasis"

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var user User
	i := 0

	user.Userno = 0
	user.Userlongname = ""
	userlist = append(userlist, user)
	i++

	for rows.Next() {
		err := rows.Scan(&user.Userno, &user.Userlongname)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		userlist = append(userlist, user)
		i++
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return

}

func SelectEventuserList(eventid string) (userlist []User, status int) {

	status = 0

	sql := "select e.userno,u.longname "
	sql += " from eventuser e join user u on e.userno=u.userno "
	sql += " where e.eventid = ? "
	sql += " order by e.userno"

	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(eventid)
	if err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	var user User
	i := 0

	user.Userno = 0
	user.Userlongname = "ポイント差は不要"
	userlist = append(userlist, user)
	i++

	for rows.Next() {
		err := rows.Scan(&user.Userno, &user.Userlongname)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		userlist = append(userlist, user)
		i++
	}
	if err = rows.Err(); err != nil {
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return

}

func SelectEventList(userno int) (eventlist []Event, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	/*
		if userno != 0 {
			stmt, Err = Db.Prepare("select eventid, event_name from event where endtime IS not null and nobasis = ? order by endtime desc")
		} else {
			stmt, Err = Db.Prepare("select eventid, event_name from event where endtime IS not null order by endtime desc")
		}
	*/

	stmt, srdblib.Dberr = srdblib.Db.Prepare("select eventid, event_name from event where endtime IS not null and nobasis = ? order by endtime desc")
	if srdblib.Dberr != nil {
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt.Close()

	/*
		if userno != 0 {
			rows, Err = stmt.Query(userno)
		} else {
			rows, Err = stmt.Query()
		}
	*/
	rows, srdblib.Dberr = stmt.Query(userno)
	if srdblib.Dberr != nil {
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer rows.Close()

	var event Event
	i := 0
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&event.EventID, &event.EventName)
		if srdblib.Dberr != nil {
			log.Printf("err=[%s]\n", srdblib.Dberr.Error())
			status = -1
			return
		}
		eventlist = append(eventlist, event)
		i++
		/*
			if i == 10 {
				break
			}
		*/
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}

	return

}

func SelectLastEventList() (eventlist []Event, status int) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	//	sql := "select eventid, event_name, period, starttime, endtime, nobasis, longname from event join user "
	sql := "select eventid, event_name, period, starttime, endtime, nobasis, modmin, modsec, longname, maxpoint from event join user "
	sql += " where nobasis = userno and endtime IS not null order by endtime desc "
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sql)
	if srdblib.Dberr != nil {
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query()
	if srdblib.Dberr != nil {
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}
	defer rows.Close()

	var event Event
	i := 0
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&event.EventID, &event.EventName, &event.Period, &event.Starttime, &event.Endtime, &event.Pntbasis, &event.Modmin, &event.Modsec, &event.Pbname, &event.Maxpoint)
		if srdblib.Dberr != nil {
			log.Printf("err=[%s]\n", srdblib.Dberr.Error())
			status = -1
			return
		}
		event.Gscale = event.Maxpoint % 100
		event.Maxpoint = event.Maxpoint - event.Gscale
		eventlist = append(eventlist, event)
		i++
		if i == Serverconfig.NoEvent {
			break
		}
	}
	if srdblib.Dberr = rows.Err(); srdblib.Dberr != nil {
		log.Printf("err=[%s]\n", srdblib.Dberr.Error())
		status = -1
		return
	}

	tnow := time.Now()
	for i = 0; i < len(eventlist); i++ {
		eventlist[i].S_start = eventlist[i].Starttime.Format("2006-01-02 15:04")
		eventlist[i].S_end = eventlist[i].Endtime.Format("2006-01-02 15:04")

		if eventlist[i].Starttime.After(tnow) {
			eventlist[i].Status = "これから開催"
		} else if eventlist[i].Endtime.Before(tnow) {
			eventlist[i].Status = "終了"
		} else {
			eventlist[i].Status = "開催中"
		}

	}

	return

}

/*
func OpenDb() (status int) {

	status = 0

	if (*Dbconfig).Dbhost == "" {
		(*Dbconfig).Dbhost = "localhost"
	}
	if (*Dbconfig).Dbport == "" {
		(*Dbconfig).Dbport = "3306"
	}
	cnc := "@tcp"
	if Dbconfig.UseSSH {
		Dialer.Hostname = Sshconfig.Hostname
		Dialer.Port = Sshconfig.Port
		Dialer.Username = Sshconfig.Username
		Dialer.Password = Sshconfig.Password
		Dialer.PrivateKey = Sshconfig.PrivateKey

		mysqldrv.New(&Dialer).RegisterDial("ssh+tcp")
		cnc = "@ssh+tcp"
	}
	cnc += "(" + Dbconfig.Dbhost + ":" + Dbconfig.Dbport + ")"
	Db, Err = sql.Open("mysql", Dbconfig.Dbuser+":"+Dbconfig.Dbpw+cnc+"/"+Dbconfig.Dbname+"?parseTime=true&loc=Asia%2FTokyo")

	if Err != nil {
		status = -1
	}
	return
}
*/

func SelectEventInfAndRoomList() (IDlist []int, status int) {

	status = 0

	/*
		//	sql := "select eventno, event_name, period, starttime, endtime from event where eventid ='"+Event_inf.Event_ID+"'"
		sql := "select eventno, event_name, period, starttime, endtime from event where eventid = ?"
		err := Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.Event_no, &Event_inf.Event_name, &Event_inf.Period, &Event_inf.Start_time, &Event_inf.End_time)

		if err != nil {
			log.Printf("select eventno, starttime, endtime from event where eventid ='%s'\n", Event_inf.Event_ID)
			log.Printf("err=[%s]\n", err.Error())
			//	if err.Error() != "sql: no rows in result set" {
			status = -1
			return
			//	}
		}
	*/

	//	Event_inf, _ = SelectEventInf(Event_inf.Event_ID)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(Event_inf.Event_ID)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	log.Printf("eventno=%d\n", Event_inf.Event_no)

	start_date := Event_inf.Start_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.Start_time.Hour()) * time.Hour)
	end_date := Event_inf.End_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

	//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

	Event_inf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
	Event_inf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - Event_inf.Start_date

	//	log.Printf("Start_data=%f Dperiod=%f\n", Event_inf.Start_date, Event_inf.Dperiod)

	//	err = Db.QueryRow("select max(point) from points where event_id = '" + fmt.Sprintf("%d", Event_inf.Event_no) + "'").Scan(&Event_inf.MaxPoint)
	//	sql := "select max(point) from eventuser where eventno = ? and graph = 'Y'"
	sql := "select max(point) from eventuser where eventid = ? and graph = 'Y'"
	err = srdblib.Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.MaxPoint)
	//	err = srdblib.Db.QueryRow(sql, Event_inf.Event_ID).Scan(&Event_inf.Maxpoint)

	if err != nil {
		log.Printf("select max(point) from eventuser where eventid = '%s'\n", Event_inf.Event_ID)
		log.Printf("err=[%s]\n", err.Error())
		status = -2
		return
	}

	//	log.Printf("MaxPoint=%d\n", Event_inf.MaxPoint)

	//	-------------------------------------------------------------------
	//	sql := "select user_id from points where event_id = ? and idx = ( select max(idx) from points where event_id = ? ) order by point desc"
	sql = " select userno from eventuser "
	sql += " where graph = 'Y' "
	//	sql += " and eventno = ? "
	sql += " and eventid = ? "
	sql += " order by point desc"
	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(Event_inf.Event_ID)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}
	defer rows.Close()

	id := 0
	i := 0
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}
		IDlist = append(IDlist, id)
		i++
		if i == Event_inf.Maxdsp {
			break
		}
	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		status = -1
		return
	}

	return
}

/*
func SelectEventInf(eventid string) (eventinf Event_Inf, status int) {

	status = 0

	sql := "select eventid,ieventid,event_name, period, starttime, endtime, noentry, intervalmin, modmin, modsec, "
	sql += " Fromorder, Toorder, Resethh, Resetmm, Nobasis, Maxdsp, cmap, target, maxpoint "
	sql += " from event where eventid = ?"
	err := Db.QueryRow(sql, eventid).Scan(
		&eventinf.Event_ID,
		&eventinf.I_Event_ID,
		&eventinf.Event_name,
		&eventinf.Period,
		&eventinf.Start_time,
		&eventinf.End_time,
		&eventinf.NoEntry,
		&eventinf.Intervalmin,
		&eventinf.Modmin,
		&eventinf.Modsec,
		&eventinf.Fromorder,
		&eventinf.Toorder,
		&eventinf.Resethh,
		&eventinf.Resetmm,
		&eventinf.Nobasis,
		&eventinf.Maxdsp,
		&eventinf.Cmap,
		&eventinf.Target,
		&eventinf.Maxpoint,
	)

	if err != nil {
		log.Printf("%s\n", sql)
		log.Printf("err=[%s]\n", err.Error())
		//	if err.Error() != "sql: no rows in result set" {
		status = -1
		return
		//	}
	}

	//	log.Printf("eventno=%d\n", Event_inf.Event_no)

	start_date := eventinf.Start_time.Truncate(time.Hour).Add(-time.Duration(eventinf.Start_time.Hour()) * time.Hour)
	end_date := eventinf.End_time.Truncate(time.Hour).Add(-time.Duration(eventinf.End_time.Hour())*time.Hour).AddDate(0, 0, 1)

	//	log.Printf("start_t=%v\nstart_d=%v\nend_t=%v\nend_t=%v\n", Event_inf.Start_time, start_date, Event_inf.End_time, end_date)

	eventinf.Start_date = float64(start_date.Unix()) / 60.0 / 60.0 / 24.0
	eventinf.Dperiod = float64(end_date.Unix())/60.0/60.0/24.0 - Event_inf.Start_date

	//	log.Printf("Start_data=%f Dperiod=%f\n", eventinf.Start_date, eventinf.Dperiod)

	return
}
*/

func SelectPointList(userno int, eventid string) (norow int, tp *[]time.Time, pp *[]int) {

	norow = 0

	//	log.Printf("SelectPointList() userno=%d eventid=%s\n", userno, eventid)
	stmt1, err := srdblib.Db.Prepare("SELECT count(*) FROM points where user_id = ? and eventid = ?")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt1.Close()

	//	var norow int
	err = stmt1.QueryRow(userno, eventid).Scan(&norow)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	//	fmt.Println(norow)

	//	----------------------------------------------------

	//	stmt1, err = Db.Prepare("SELECT max(t.t) FROM timeacq t join points p where t.idx=p.idx and user_id = ? and event_id = ?")
	stmt1, err = srdblib.Db.Prepare("SELECT max(ts) FROM points where user_id = ? and eventid = ?")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt1.Close()

	var tfinal time.Time
	err = stmt1.QueryRow(userno, eventid).Scan(&tfinal)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	islastdata := false
	if tfinal.After(Event_inf.End_time.Add(time.Duration(-Event_inf.Intervalmin) * time.Minute)) {
		islastdata = true
	}
	//	fmt.Println(norow)

	//	----------------------------------------------------

	t := make([]time.Time, norow)
	point := make([]int, norow)
	if islastdata {
		t = make([]time.Time, norow+1)
		point = make([]int, norow+1)
	}

	tp = &t
	pp = &point

	if norow == 0 {
		return
	}

	//	----------------------------------------------------

	//	stmt2, err := Db.Prepare("select t.t, p.point from points p join timeacq t on t.idx = p.idx where user_id = ? and event_id = ? order by t.t")
	stmt2, err := srdblib.Db.Prepare("select ts, point from points where user_id = ? and eventid = ? order by ts")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt2.Close()

	rows, err := stmt2.Query(userno, eventid)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(&t[i], &point[i])
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			//	status = -1
			return
		}
		i++

	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}

	if islastdata {
		t[norow] = t[norow-1].Add(15 * time.Minute)
		point[norow] = point[norow-1]
	}

	tp = &t
	pp = &point

	return
}

func MakePointPerDay(eventid string) (p_pointperday *PointPerDay, status int) {

	status = 0

	Event_inf.Event_ID = eventid
	_, sts := SelectEventInfAndRoomList()

	if sts != 0 {
		log.Printf("MakePointPerDay() status of SelectEventInfAndRoomList() =%d\n", sts)
		status = sts
		return
	}

	dstart := Event_inf.Start_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.Start_time.Hour()) * time.Hour)
	if Event_inf.Start_time.Hour()*60+Event_inf.Start_time.Minute() > Event_inf.Resethh*60+Event_inf.Resetmm {
		dstart = dstart.AddDate(0, 0, 1)
	}

	dend := Event_inf.End_time.Truncate(time.Hour).Add(-time.Duration(Event_inf.End_time.Hour()) * time.Hour)
	if Event_inf.End_time.Hour()*60+Event_inf.End_time.Minute() > Event_inf.Resethh*60+Event_inf.Resetmm {
		dend = dend.AddDate(0, 0, 1)
	}

	days := int(dend.Sub(dstart).Hours() / 24)
	dstart = dstart.Add(time.Duration(Event_inf.Resethh*60+Event_inf.Resetmm) * time.Minute)

	log.Printf(" dstart=%s dend=%s days=%d\n", dstart.Format("2006/01/02 15:04:05"), dend.Format("2006/01/02 15:04:05"), days)

	var pointperday PointPerDay
	pointperday.Pointrecordlist = make([]PointRecord, days+1)
	pointperday.Eventname = Event_inf.Event_name
	pointperday.Eventid = eventid
	pointperday.Period = Event_inf.Period

	var roominfolist RoomInfoList
	_, _ = SelectEventRoomInfList(eventid, &roominfolist)
	log.Printf(" no of rooms. = %d\n", len(roominfolist))

	iu := 0 //	リスト作成の対象となるルームのインデックス

	for i := 0; i < len(roominfolist); i++ {

		log.Printf(" Room=%s Graph=%s\n", roominfolist[i].Longname, roominfolist[i].Graph)
		if roominfolist[i].Graph != "Checked" {
			continue
		}

		pointperday.Longnamelist = append(pointperday.Longnamelist, LongName{roominfolist[i].Longname})
		pointperday.Usernolist = append(pointperday.Usernolist, roominfolist[i].Userno)
		for k := 0; k < days+1; k++ {
			pointperday.Pointrecordlist[k].Day = dstart.AddDate(0, 0, k-1).Format("2006/01/02")
			pointperday.Pointrecordlist[k].Tday = dstart.AddDate(0, 0, k)
			if pointperday.Pointrecordlist[k].Tday.After(time.Now()) {
				pointperday.Pointrecordlist[k].Tday = time.Now().Truncate(time.Second)
			}
			if pointperday.Pointrecordlist[k].Tday.After(Event_inf.End_time) {
				pointperday.Pointrecordlist[k].Tday = Event_inf.End_time
			}
			pointperday.Pointrecordlist[k].Pointlist = append(pointperday.Pointrecordlist[k].Pointlist, Point{0, "", ""})
		}

		norow, tp, pp := SelectPointList(roominfolist[i].Userno, Event_inf.Event_ID)

		log.Printf(" no of point data=%d\n", norow)
		if norow == 0 {
			continue
		}

		d := dstart
		k := 0

		for ; ; k++ {
			if (*tp)[0].Before(d.AddDate(0, 0, k)) {
				break
			}
		}

		lastpoint := 0
		prvpoint := 0
		for j := 0; j < len(*tp); j++ {
			if (*tp)[j].After(d.AddDate(0, 0, k)) {
				log.Printf("i(room)=%d, j(time)=%d(%s), k(day)=%d\n", i, j, (*tp)[j].Format("01/02 15:04"), k)
				log.Printf("pointperday.Pointrecordlist[k].Pointlist=%v\n", pointperday.Pointrecordlist[k].Pointlist)
				if (*tp)[j].Sub(d.AddDate(0, 0, k)) < 30*time.Minute || j == 0 || (*pp)[j] == (*pp)[j-1] {
					pointperday.Pointrecordlist[k].Pointlist[iu].Pnt = lastpoint - prvpoint
					pointperday.Pointrecordlist[k].Pointlist[iu].Spnt = humanize.Comma(int64(lastpoint - prvpoint))
					pointperday.Pointrecordlist[k].Pointlist[iu].Tpnt = humanize.Comma(int64(lastpoint))
				} else {
					//	欠測が発生したと思われる場合
					pointperday.Pointrecordlist[k].Pointlist[iu].Pnt = lastpoint - prvpoint
					pointperday.Pointrecordlist[k].Pointlist[iu].Spnt = ""
				}
				prvpoint = lastpoint
				k++
			}
			lastpoint = (*pp)[j]
		}
		pointperday.Pointrecordlist[k].Pointlist[iu].Pnt = lastpoint - prvpoint
		pointperday.Pointrecordlist[k].Pointlist[iu].Spnt = humanize.Comma(int64(lastpoint - prvpoint))
		pointperday.Pointrecordlist[k].Pointlist[iu].Tpnt = humanize.Comma(int64(lastpoint))

		iu++
	}

	//	日々の獲得ポイントが空白の場合は次の日の獲得ポイントは無意味であるので空白とする。
	for k := days - 1; k >= 0; k-- {
		for i := 0; i < iu; i++ {
			if pointperday.Pointrecordlist[k].Pointlist[i].Spnt == "" {
				pointperday.Pointrecordlist[k+1].Pointlist[i].Spnt = ""
				pointperday.Pointrecordlist[k+1].Pointlist[i].Pnt = 0
			}
		}
	}

	p_pointperday = &pointperday

	return
}

func MakePointPerSlot(eventid string) (perslotinflist []PerSlotInf, status int) {

	var perslotinf PerSlotInf

	status = 0

	Event_inf.Event_ID = eventid
	//	eventno, eventname, period := SelectEventNoAndName(eventid)
	eventname, period, _ := SelectEventNoAndName(eventid)

	var roominfolist RoomInfoList

	_, sts := SelectEventRoomInfList(eventid, &roominfolist)

	if sts != 0 {
		log.Printf("status of SelectEventRoomInfList() =%d\n", sts)
		status = sts
		return
	}

	for i := 0; i < len(roominfolist); i++ {

		if roominfolist[i].Graph != "Checked" {
			continue
		}

		var perslot PerSlot

		userid := roominfolist[i].Userno

		perslotinf.Eventname = eventname
		perslotinf.Eventid = eventid
		perslotinf.Period = period

		perslotinf.Roomname = roominfolist[i].Name
		perslotinf.Roomid = userid
		perslotinf.Perslotlist = make([]PerSlot, 0)

		norow, tp, pp := SelectPointList(userid, eventid)

		if norow == 0 {
			continue
		}

		sameaslast := true
		plast := (*pp)[0]
		pprv := (*pp)[0]
		tdstart := ""
		tstart := time.Now().Truncate(time.Second)

		for i, t := range *tp {
			//	if (*pp)[i] != plast && sameaslast {
			if (*pp)[i] != plast {
				tstart = t
				/*
					if i != 0 {
						log.Printf("(1) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, (*tp)[i-1]=%s\n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"), (*tp)[i-1].Format("01/02 15:04"))
					} else {
						log.Printf("(1) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, \n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"))
					}
				*/
				if sameaslast {
					//	これまで変化しなかった獲得ポイントが変化し始めた
					pdstart := t.Add(time.Duration(-Event_inf.Modmin) * time.Minute).Format("2006/01/02")
					if pdstart != tdstart {
						perslot.Dstart = pdstart
						tdstart = pdstart
					} else {
						perslot.Dstart = ""
					}
					perslot.Timestart = t.Add(time.Duration(-Event_inf.Modmin) * time.Minute)
					//	perslot.Tstart = t.Add(time.Duration(-Event_inf.Modmin) * time.Minute).Format("15:04")
					if t.Sub((*tp)[i-1]) < 31*time.Minute {
						perslot.Tstart = perslot.Timestart.Format("15:04")
					} else {
						perslot.Tstart = "n/a"
					}
					//	perslot.Tstart = perslot.Timestart.Format("15:04")

					sameaslast = false
					//	} else if (*pp)[i] == plast && !sameaslast && (*tp)[i].Sub((*tp)[i-1]) > 11*time.Minute {
				}
			} else if (*pp)[i] == plast {
				//	if !sameaslast && (*tp)[i].Sub((*tp)[i-1]) > 16*time.Minute {
				if !sameaslast && t.Sub(tstart) > 11*time.Minute {
					//	if !sameaslast {
					/*
						if i != 0 {
							log.Printf("(2) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, (*tp)[i-1]=%s\n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"), (*tp)[i-1].Format("01/02 15:04"))
						} else {
							log.Printf("(2) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, \n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"))
						}
					*/
					if perslot.Tstart != "n/a" {
						perslot.Tend = (*tp)[i-1].Add(time.Duration(-Event_inf.Modmin) * time.Minute).Format("15:04")
					} else {
						perslot.Tend = "n/a"
					}
					perslot.Ipoint = plast - pprv
					perslot.Point = humanize.Comma(int64(plast - pprv))
					perslot.Tpoint = humanize.Comma(int64(plast))
					sameaslast = true
					perslotinf.Perslotlist = append(perslotinf.Perslotlist, perslot)
					pprv = plast
				}
				//	sameaslast = true
			}
			/* else
			{
					if i != 0 {
						log.Printf("(3) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, (*tp)[i-1]=%s\n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"), (*tp)[i-1].Format("01/02 15:04"))
					} else {
						log.Printf("(3) (*pp)[i]=%d, plast=%d, sameaslast=%v, (*tp)[i]=%s, \n", (*pp)[i], plast, sameaslast, (*tp)[i].Format("01/02 15:04"))
					}
			}
			*/
			plast = (*pp)[i]
		}

		if len(perslotinf.Perslotlist) != 0 {
			perslotinflist = append(perslotinflist, perslotinf)
		}

		UpdatePointsSetQstatus(eventid, userid, perslot.Tstart, perslot.Tend, perslot.Point)

	}

	return
}

func UpdatePointsSetQstatus(
	eventid string,
	userno int,
	tstart string,
	tend string,
	point string,
) (status int) {
	status = 0

	log.Printf("  *** UpdatePointsSetQstatus() *** eventid=%s userno=%d\n", eventid, userno)

	nrow := 0
	//	err := Db.QueryRow("select count(*) from points where eventid = ? and user_id = ? and pstatus = 'Conf.'", eventid, userno).Scan(&nrow)
	sql := "select count(*) from points where eventid = ? and user_id = ? and ( pstatus = 'Conf.' or pstatus = 'Prov.' )"
	err := srdblib.Db.QueryRow(sql, eventid, userno).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	if nrow != 1 {
		return
	}

	log.Printf("  *** UpdatePointsSetQstatus() Update!\n")

	sql = "update points set qstatus =?,"
	sql += "qtime=? "
	//	sql += "where user_id=? and eventid = ? and pstatus = 'Conf.'"
	sql += "where user_id=? and eventid = ? and ( pstatus = 'Conf.' or pstatus = 'Prov.' )"
	stmt, err := srdblib.Db.Prepare(sql)
	if err != nil {
		log.Printf("UpdatePointsSetQstatus() Update/Prepare err=%s\n", err.Error())
		status = -1
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(point, tstart+"--"+tend, userno, eventid)

	if err != nil {
		log.Printf("error(UpdatePointsSetQstatus() Update/Exec) err=%s\n", err.Error())
		status = -2
	}

	return
}

func SelectScoreList(user_id int) (x *[]float64, y *[]float64) {

	stmt1, err := srdblib.Db.Prepare("SELECT count(*) FROM points where user_id = ? and eventid = ?")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt1.Close()

	var norow int
	err = stmt1.QueryRow(user_id, Event_inf.Event_ID).Scan(&norow)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	//	fmt.Println(norow)

	//	----------------------------------------------------

	tu := make([]float64, norow)
	point := make([]float64, norow)

	//	----------------------------------------------------

	//	stmt2, err := Db.Prepare("select t.t, p.point from points p join timeacq t on t.idx = p.idx where user_id = ? and event_id = ? order by t.t")
	stmt2, err := srdblib.Db.Prepare("select ts, point from points where user_id = ? and eventid = ? order by ts")
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer stmt2.Close()

	rows, err := stmt2.Query(user_id, Event_inf.Event_ID)
	if err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}
	defer rows.Close()
	i := 0
	var t time.Time
	for rows.Next() {
		err := rows.Scan(&t, &point[i])
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			//	status = -1
			return
		}
		if t.Before(Event_inf.Start_time) {
			t = Event_inf.Start_time
		}
		tu[i] = float64(t.Unix())/60.0/60.0/24.0 - Event_inf.Start_date
		//	log.Printf("t=%v tu[%d]=%f\n", t, i, tu[i])
		i++

	}
	if err = rows.Err(); err != nil {
		//	log.Fatal(err)
		log.Printf("err=[%s]\n", err.Error())
		//	status = -1
		return
	}

	x = &tu
	y = &point

	return
}

func DetYaxScale(
	maxpoint int,
) (
	yupper int,
	yscales int,
	yscalel int,
	status int,
) {

	status = 0

	type Yaxis struct {
		Yupper  int
		Yscales int
		Yscalel int
	}

	yaxis := []Yaxis{
		{10000, 500, 5},
		{15000, 500, 5},
		{20000, 1000, 5},
		{30000, 1000, 5},
		{50000, 2000, 5},
		{70000, 2000, 5},
	}

	mlt := 1
	for i := 0; ; i++ {
		if i != 0 && i%len(yaxis) == 0 {
			mlt *= 10
		}
		iy := i % len(yaxis)

		if maxpoint < yaxis[iy].Yupper*mlt {
			yupper = yaxis[iy].Yupper * mlt
			yscales = yaxis[iy].Yscales * mlt
			yscalel = yaxis[iy].Yscalel
			break
		}
	}
	return
}

func DetXaxScale(
	xupper float64,
) (
	xscaled int,
	xscalet int,
	status int,
) {

	status = 0

	type Xaxis struct {
		Xupper  float64
		Xscaled int
		Xscalet int
	}
	xaxis := []Xaxis{
		{1.1, 24, 1},
		{2.1, 12, 1},
		{5.1, 8, 1},
		{10.1, 4, 1},
		{32.1, 2, 1},
		{64.1, 1, 2},
		{128.1, 1, 4},
		{256.1, 1, 7},
	}

	ix := 0
	for ; ix < len(xaxis); ix++ {
		if xupper < xaxis[ix].Xupper {
			xscaled = xaxis[ix].Xscaled
			xscalet = xaxis[ix].Xscalet
			break
		}
	}
	if ix == len(xaxis) {
		xscaled = 1
		xscalet = 28
	}

	return
}

/*
func GraphDfr01(filename string, IDlist []int) {

	//	描画領域を決定する
	width := 3840.0
	height := 2160.0
	lwmargin := width / 24.0
	rwmargin := width / 6.0
	uhmargin := height / 7.5
	lhmargin := height / 15.0
	bstroke := width / 800.0

	vwidth := width - lwmargin - rwmargin
	vheight := height - uhmargin - lhmargin

	xorigin := lwmargin
	yorigin := height - lhmargin

	//	SVG出力ファイルを設定し、背景色を決める。
	file, err := os.OpenFile("public/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//	panic(err)
		return
	}

	bw := bufio.NewWriter(file)

	canvas := svg.New(bw)

	//	canvas := svg.New(os.Stdout)

	canvas.Start(width, height)


	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	yupper, yscales, yscalel, _ := DetYaxScale(Event_inf.MaxPoint)

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "獲得ポイントの推移",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, eventname,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), period,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	//	y軸（ポイント軸）を描画する

	dyl := float64(yscales) * yscale
	value := int64(0)
	yl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%yscalel == 0 {
			wstr = 0.3

			canvas.Text(xorigin-bstroke*5.0, yorigin+yl+bstroke*2.5, humanize.Comma(value),
				"text-anchor:end;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		}
		canvas.Line(xorigin, yorigin+yl, xorigin+vwidth, yorigin+yl, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		yl += dyl
		if -yl > vheight+10 {
			break
		}
		value += int64(yscales)

	}

	//	------------------------------------------

	//	x軸（時間軸）を描画する

	xupper := Event_inf.Dperiod
	xscale := vwidth / float64(xupper)
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	dxl := 1.0 / float64(xscaled) * xscale
	tval := Event_inf.Start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%xscaled == 0 {
			wstr = 0.3
			if i%(xscaled*xscalet) == 0 {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("1/2"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
				tval = tval.AddDate(0, 0, xscalet)
			}

		}
		canvas.Line(xorigin+xl, yorigin, xorigin+xl, yorigin-vheight, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		xl += dxl
		if xl > vwidth+10 {
			break
		}
	}

	//	獲得ポイントデータを描画する

	onemin := 1.0 / 24.0 / 60.0
	xb, yb := SelectScoreList(Event_inf.Nobasis)

	j := 0
	for _, id := range IDlist {

		color, _ := SelectUserColor(id, Event_inf.Event_ID)

		x, y := SelectScoreList(id)
		maxp := 20

		//	no := len(*x)

		xo := make([]float64, maxp)
		yo := make([]float64, maxp)
		tl := 999.0
		yl := -1000000.0
		k := 0

		ib := 0
		flat := false
		if yb[ib+1] == yb[ib] {
			flat = true
		}

		for i := 0; i < len(*x); i++ {
			//	fmt.Printf("(%7.1f,%10.1f)\n", (*x)[i], (*y)[i])

			if math.Abs(y[i]-yb[ib]) < onemin {

			}

			xt := xorigin + (*x)[i]*xscale
			yt := yorigin + (*y)[i]*yscale
			//	fmt.Printf("(*x).[i]=%.3f tl=%.3f (*x)[i]-tl=%.3f\n", (*x)[i], tl, (*x)[i]-tl)
			if (*x)[i]-tl > 0.011 && (*y)[i]-yl > 1.0 {
				canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+color+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				xo[0] = xt
				yo[0] = yt
				tl = (*x)[i]
				yl = (*y)[i]
				k = 1
				continue
			}
			xo[k] = xt
			yo[k] = yt
			tl = (*x)[i]
			yl = (*y)[i]
			k++
			if k == maxp {
				canvas.Polyline(xo, yo, "fill=\"none\" stroke=\""+color+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				xo[0] = xt
				yo[0] = yt
				k = 1
			}
		}
		if k > 1 {
			canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+color+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		}

		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(j)

		canvas.Line(xln, yln, xln+rwmargin/4.0, yln, "stroke=\""+color+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		//	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, fmt.Sprintf("%d", IDlist[j]),
		//		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
		longname, _, _, _, _, sts := SelectUserName(IDlist[j])
		if sts != 0 {
			longname = fmt.Sprintf("%d", IDlist[j])
		}
		canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		j++
	}

	canvas.End()

	bw.Flush()
	file.Close()

	return
}

func GraphDfr(eventid string) (filename string, status int) {

	status = 0

	Event_inf, status = SelectEventInf(eventid)
	if status != 0 {
		return
	}
	eventname := Event_inf.Event_name
	period := Event_inf.Period

	IDlist, sts := SelectEventInfAndRoomList()

	if sts != 0 {
		log.Printf("status of SelectEventInfAndRoomList() =%d\n", sts)
		status = sts
		return
	}

	ibasis := 0
	for ; ibasis < len(IDlist); ibasis++ {
		if IDlist[ibasis] == Event_inf.Nobasis {
			break
		}
	}
	if ibasis == len(IDlist) {
		status = 1
		return
	}

	filename = fmt.Sprintf("%0d.svg", os.Getpid()%100)

	GraphDfr01(filename, IDlist)

	return
}
*/

func GraphScore01(filename string, IDlist []int, eventname string, period string, maxpoint int) {

	//	描画領域を決定する
	width := 3840.0
	height := 2160.0
	lwmargin := width / 24.0
	rwmargin := width / 6.0
	uhmargin := height / 7.5
	lhmargin := height / 15.0
	bstroke := width / 800.0

	vwidth := width - lwmargin - rwmargin
	vheight := height - uhmargin - lhmargin

	xorigin := lwmargin
	yorigin := height - lhmargin

	//	SVG出力ファイルを設定し、背景色を決める。
	file, err := os.OpenFile("public/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//	panic(err)
		return
	}

	bw := bufio.NewWriter(file)

	canvas := svg.New(bw)

	//	canvas := svg.New(os.Stdout)

	canvas.Start(width, height)

	/*
		canvas.Circle(width/2, height/2, 100)
		canvas.Text(width/2, height/2, "ポケGO", "text-anchor:middle;font-size:30px;fill:white;")
	*/

	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	yupper := 0
	yscales := 0
	yscalel := 0

	if maxpoint != 0 {
		yupper, yscales, yscalel, _ = DetYaxScale(maxpoint - 1)
	} else if Event_inf.Target > Event_inf.MaxPoint {
		//	} else if Event_inf.Target > Event_inf.Maxpoint {
		yupper, yscales, yscalel, _ = DetYaxScale(Event_inf.Target - 1)
	} else {
		yupper, yscales, yscalel, _ = DetYaxScale(Event_inf.MaxPoint)
		//	yupper, yscales, yscalel, _ = DetYaxScale(Event_inf.Maxpoint)
	}

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "獲得ポイントの推移",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, eventname,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), period,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	//	y軸（ポイント軸）を描画する

	dyl := float64(yscales) * yscale
	value := int64(0)
	yl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%yscalel == 0 {
			wstr = 0.3

			canvas.Text(xorigin-bstroke*5.0, yorigin+yl+bstroke*2.5, humanize.Comma(value),
				"text-anchor:end;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		}
		canvas.Line(xorigin, yorigin+yl, xorigin+vwidth, yorigin+yl, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		yl += dyl
		if -yl > vheight+10 {
			break
		}
		value += int64(yscales)

	}

	//	------------------------------------------

	//	x軸（時間軸）を描画する

	xupper := Event_inf.Dperiod
	xscale := vwidth / float64(xupper)
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	dxl := 1.0 / float64(xscaled) * xscale
	tval := Event_inf.Start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%xscaled == 0 {
			wstr = 0.3
			if i%(xscaled*xscalet) == 0 {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("1/2"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
				tval = tval.AddDate(0, 0, xscalet)
			}

		}
		canvas.Line(xorigin+xl, yorigin, xorigin+xl, yorigin-vheight, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		xl += dxl
		if xl > vwidth+10 {
			break
		}
	}

	//	ターゲットラインを描画する
	if Event_inf.Target != 0 {
		x1 := xorigin + (float64(Event_inf.Start_time.Unix())/60.0/60.0/24.0-Event_inf.Start_date)*xscale
		x2 := xorigin + (float64(Event_inf.End_time.Unix())/60.0/60.0/24.0-Event_inf.Start_date)*xscale
		y1 := yorigin
		y2 := yorigin + float64(Event_inf.Target)*yscale

		log.Printf("Target (x1, y1) %10.2f,%10.2f (x2, y2) %10.2f,%10.2f xorgin, yorigin, vheight %10.2f, %10.2f %10.2f\n",
			x1, y1, x2, y2, xorigin, yorigin, vheight)

		if y2 < yorigin-vheight {
			x2 = (x2-xorigin)*vheight/(yorigin-y2) + xorigin
			y2 = yorigin - vheight
		}

		log.Printf("Target (x1, y1) %10.2f,%10.2f (x2, y2) %10.2f,%10.2f xorgin, yorigin, vheight %10.2f, %10.2f %10.2f\n",
			x1, y1, x2, y2, xorigin, yorigin, vheight)

		canvas.Line(x1, y1, x2, y2, `stroke="white" stroke-width="`+fmt.Sprintf("%.2f", bstroke*0.5)+`" stroke-dasharray="20,10"`)
	}

	//	獲得ポイントデータを描画する

	j := 0
	for _, id := range IDlist {

		_, cvalue, _ := SelectUserColor(id, Event_inf.Event_ID)

		x, y := SelectScoreList(id)
		maxp := 20

		//	no := len(*x)

		xo := make([]float64, maxp)
		yo := make([]float64, maxp)
		tl := 999.0
		yl := -1000000.0
		k := 0
		for i := 0; i < len(*x); i++ {
			//	fmt.Printf("(%7.1f,%10.1f)\n", (*x)[i], (*y)[i])
			xt := xorigin + (*x)[i]*xscale
			yt := yorigin + (*y)[i]*yscale
			//	fmt.Printf("(*x).[i]=%.3f tl=%.3f (*x)[i]-tl=%.3f\n", (*x)[i], tl, (*x)[i]-tl)
			if (*x)[i]-tl > 0.011 && (*y)[i]-yl > 1.0 {
				//	次のデータとの間に欠測があり、かつ欠測の前後でデータが同一でないときはその部分の描画は行わない。
				canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				xo[0] = xt
				yo[0] = yt
				tl = (*x)[i]
				yl = (*y)[i]
				k = 1
				if yt < yorigin-vheight {
					break
				}
				continue
			}

			if yt < yorigin-vheight {
				if k != 0 {
					xo[k] = (xt-xo[k-1])*(yo[k-1]-(yorigin-vheight))/(yo[k-1]-yt) + xo[k-1]
					yo[k] = yorigin - vheight
					k++
				}
				break
			} else {
				xo[k] = xt
				yo[k] = yt
			}

			tl = (*x)[i]
			yl = (*y)[i]
			k++
			if k == maxp {
				//	一定数のデータずつまとめて描画する。SVGファイルの可読性を高める。
				canvas.Polyline(xo, yo, "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				xo[0] = xt
				yo[0] = yt
				k = 1
			}
		}
		if k > 1 {
			canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		}

		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(j)

		canvas.Line(xln, yln, xln+rwmargin/4.0, yln, "stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		//	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, fmt.Sprintf("%d", IDlist[j]),
		//		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
		longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(IDlist[j])
		if sts != 0 {
			longname = fmt.Sprintf("%d", IDlist[j])
		}
		canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		j++
	}
	xln := xorigin + vwidth + bstroke*30.0
	yln := yorigin - vheight + bstroke*10*float64(j)

	canvas.Line(xln, yln, xln+rwmargin/4.0, yln, `stroke="white" stroke-width="`+fmt.Sprintf("%.2f", bstroke*0.5)+`" stroke-dasharray="20,10"`)
	//	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, fmt.Sprintf("%d", IDlist[j]),
	//		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, "Target",
		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

	canvas.End()

	bw.Flush()
	file.Close()

}

func GraphTotalPoints(eventid string, maxpoint int, gscale int) (filename string, status int) {

	status = 0

	Event_inf.Event_ID = eventid

	IDlist, sts := SelectEventInfAndRoomList()

	if sts != 0 {
		log.Printf("status of SelectEventInfAndRoomList() =%d\n", sts)
		status = sts
		return
	}

	eventname, period, _ := SelectEventNoAndName(eventid)

	Event_inf.Maxpoint = maxpoint
	Event_inf.Gscale = gscale
	UpdateEventInf(&Event_inf)

	filename = fmt.Sprintf("%0d.svg", os.Getpid()%100)

	GraphScore01(filename, IDlist, eventname, period, maxpoint)

	/*
		fmt.Printf("Content-type:text/html\n\n")
		fmt.Printf("<!DOCTYPE html>\n")
		fmt.Printf("<html lang=\"ja\">\n")
		fmt.Printf("<head>\n")
		fmt.Printf("  <meta charset=\"UTF-8\">\n")
		fmt.Printf("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
		//	fmt.Printf("  <meta http-equiv=\"refresh\" content=\"30; URL=\">\n")
		fmt.Printf("  <meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">\n")
		fmt.Printf("  <title></title>\n")
		fmt.Printf("</head>\n")
		fmt.Printf("<body>\n")
		fmt.Printf("<img src=\"test.svg\" alt=\"\" width=\"100%%\">")
		fmt.Printf("</body>\n")
		fmt.Printf("</html>\n")
	*/

	return
}

func GraphPerSlot(
	eventid string,
	perslotinflist *[]PerSlotInf,
) (
	filename string,
	status int,
) {

	status = 0

	//	Event_inf, status = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	描画領域を決定する
	width := 3840.0
	height := 2160.0
	lwmargin := width / 24.0
	rwmargin := width / 6.0
	uhmargin := height / 7.5
	lhmargin := height / 15.0
	bstroke := width / 800.0

	vwidth := width - lwmargin - rwmargin
	vheight := height - uhmargin - lhmargin

	xorigin := lwmargin
	yorigin := height - lhmargin

	//	SVG出力ファイルを設定し、背景色を決める。
	filename = fmt.Sprintf("%0d.svg", os.Getpid()%100)
	file, err := os.OpenFile("public/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//	panic(err)
		return
	}

	bw := bufio.NewWriter(file)

	canvas := svg.New(bw)

	//	canvas := svg.New(os.Stdout)

	canvas.Start(width, height)

	/*
		canvas.Circle(width/2, height/2, 100)
		canvas.Text(width/2, height/2, "ポケGO", "text-anchor:middle;font-size:30px;fill:white;")
	*/

	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	maxpoint := 0
	for _, perslotinf := range *perslotinflist {
		for _, perslot := range perslotinf.Perslotlist {
			if perslot.Ipoint > maxpoint {
				maxpoint = perslot.Ipoint
			}
		}
	}
	yupper, yscales, yscalel, _ := DetYaxScale(maxpoint)

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "配信枠毎の獲得ポイント",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, Event_inf.Event_name,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), Event_inf.Period,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	//	y軸（ポイント軸）を描画する

	dyl := float64(yscales) * yscale
	value := int64(0)
	yl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%yscalel == 0 {
			wstr = 0.3

			canvas.Text(xorigin-bstroke*5.0, yorigin+yl+bstroke*2.5, humanize.Comma(value),
				"text-anchor:end;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		}
		canvas.Line(xorigin, yorigin+yl, xorigin+vwidth, yorigin+yl, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		yl += dyl
		if -yl > vheight+10 {
			break
		}
		value += int64(yscales)

	}

	//	------------------------------------------

	//	x軸（時間軸）を描画する

	xupper := Event_inf.Dperiod
	xscale := vwidth / float64(xupper)
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	dxl := 1.0 / float64(xscaled) * xscale
	tval := Event_inf.Start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%xscaled == 0 {
			wstr = 0.3
			if i%(xscaled*xscalet) == 0 {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("1/2"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
				tval = tval.AddDate(0, 0, xscalet)
			}

		}
		canvas.Line(xorigin+xl, yorigin, xorigin+xl, yorigin-vheight, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		xl += dxl
		if xl > vwidth+10 {
			break
		}
	}

	//	配信枠毎の獲得ポイントデータを描画する

	for j, perslotinf := range *perslotinflist {
		_, cvalue, _ := SelectUserColor(perslotinf.Roomid, Event_inf.Event_ID)
		for _, perslot := range perslotinf.Perslotlist {
			y := float64(perslot.Ipoint)*yscale + yorigin
			x := (float64(perslot.Timestart.Unix())/60.0/60.0/24.0-Event_inf.Start_date)*xscale + xorigin
			log.Printf("t=%7.3f, p=%8d, x=%7.2f, y=%7.2f\n",
				float64(perslot.Timestart.Unix())/60.0/60.0/24.0-Event_inf.Start_date,
				perslot.Ipoint, x, y)
			//	canvas.Circle(x, y, 10.0, "stroke:"+cvalue+";fill:"+cvalue)
			Mark(j, canvas, x, y, 10.0, cvalue)
		}
		longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(perslotinf.Roomid)
		if sts != 0 {
			longname = fmt.Sprintf("%d", perslotinf.Roomid)
		}
		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(j)
		//	canvas.Circle(xln+rwmargin/4.0, yln, 10.0, "stroke:"+cvalue+";fill:"+cvalue)
		Mark(j, canvas, xln+rwmargin/4.0, yln, 10.0, cvalue)
		canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
	}

	canvas.End()

	bw.Flush()
	file.Close()

	return

}

func GraphPerDay(
	eventid string,
	pointperday *PointPerDay,
) (
	filename string,
	status int,
) {

	status = 0

	//	Event_inf, status = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	//	描画領域を決定する
	width := 3840.0
	height := 2160.0
	lwmargin := width / 24.0
	rwmargin := width / 6.0
	uhmargin := height / 7.5
	lhmargin := height / 15.0
	bstroke := width / 800.0

	vwidth := width - lwmargin - rwmargin
	vheight := height - uhmargin - lhmargin

	xorigin := lwmargin
	yorigin := height - lhmargin

	//	SVG出力ファイルを設定し、背景色を決める。
	filename = fmt.Sprintf("%0d.svg", os.Getpid()%100)
	file, err := os.OpenFile("public/"+filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		//	panic(err)
		return
	}

	bw := bufio.NewWriter(file)

	canvas := svg.New(bw)

	//	canvas := svg.New(os.Stdout)

	canvas.Start(width, height)
	canvas.Rect(1.0, 1.0, width-1.0, height-1.0, "stroke=\"black\" stroke-width=\"0.1\"")

	//	y軸（ポイント軸）の縮尺を決める
	maxpoint := 0
	for _, pointrecord := range (*pointperday).Pointrecordlist {
		for _, point := range pointrecord.Pointlist {
			if point.Pnt > maxpoint && point.Spnt != "" {
				maxpoint = point.Pnt
			}
		}
	}
	yupper, yscales, yscalel, _ := DetYaxScale(maxpoint)

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), "配信日毎の獲得ポイント",
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, Event_inf.Event_name,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), Event_inf.Period,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	//	y軸（ポイント軸）を描画する

	dyl := float64(yscales) * yscale
	value := int64(0)
	yl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%yscalel == 0 {
			wstr = 0.3

			canvas.Text(xorigin-bstroke*5.0, yorigin+yl+bstroke*2.5, humanize.Comma(value),
				"text-anchor:end;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")

		}
		canvas.Line(xorigin, yorigin+yl, xorigin+vwidth, yorigin+yl, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		yl += dyl
		if -yl > vheight+10 {
			break
		}
		value += int64(yscales)

	}

	//	------------------------------------------

	//	x軸（時間軸）を描画する

	xupper := Event_inf.Dperiod
	xscale := vwidth / float64(xupper)
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	dxl := 1.0 / float64(xscaled) * xscale
	tval := Event_inf.Start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if i%xscaled == 0 {
			wstr = 0.3
			if i%(xscaled*xscalet) == 0 {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("1/2"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
				tval = tval.AddDate(0, 0, xscalet)
			}

		}
		canvas.Line(xorigin+xl, yorigin, xorigin+xl, yorigin-vheight, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		xl += dxl
		if xl > vwidth+10 {
			break
		}
	}

	colorlist := make([]string, len((*pointperday).Usernolist))
	for i, userno := range (*pointperday).Usernolist {
		_, colorlist[i], _ = SelectUserColor(userno, Event_inf.Event_ID)
		longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(userno)
		if sts != 0 {
			longname = fmt.Sprintf("%d", userno)
		}
		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(i)
		//	canvas.Circle(xln+rwmargin/4.0, yln, 10.0, "stroke:"+colorlist[i]+";fill:"+colorlist[i])
		Mark(i, canvas, xln+rwmargin/4.0, yln, 10.0, colorlist[i])

		canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, longname,
			"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
	}

	//	日毎の獲得ポイントデータを描画する
	for _, pointrecord := range (*pointperday).Pointrecordlist {
		x := (float64(pointrecord.Tday.Unix())/60.0/60.0/24.0-Event_inf.Start_date)*xscale + xorigin
		for i, point := range pointrecord.Pointlist {
			if point.Spnt == "" {
				continue
			}
			y := float64(point.Pnt)*yscale + yorigin
			//	log.Printf("t=%7.3f, p=%8d, x=%7.2f, y=%7.2f\n",
			//		float64(pointrecord.Tday.Unix())/60.0/60.0/24.0-Event_inf.Start_date,
			//		point.Pnt, x, y)
			//	canvas.Circle(x, y, 10.0, "stroke:"+colorlist[i]+";fill:"+colorlist[i])
			Mark(i, canvas, x, y, 10.0, colorlist[i])
		}
	}

	canvas.End()

	bw.Flush()
	file.Close()

	return

}

func Mark(j int, canvas *svg.SVG, x0, y0, d float64, color string) {

	switch j % 4 {
	case 0:
		canvas.Circle(x0, y0, d, "stroke:"+color+";fill:"+color)
	case 1:
		dyu := d * 1.2
		dyl := dyu * 0.5
		dx := dyu * 0.866
		x := make([]float64, 3)
		y := make([]float64, 3)
		x[0] = x0
		y[0] = y0 - dyu
		x[1] = x0 - dx
		y[1] = y0 + dyl
		x[2] = x0 + dx
		y[2] = y0 + dyl
		canvas.Polygon(x, y, "stroke:"+color+";fill:"+color)
	case 2:
		d = d * 0.9
		x := make([]float64, 4)
		y := make([]float64, 4)
		x[0] = x0 - d
		y[0] = y0 - d
		x[1] = x[0]
		y[1] = y0 + d
		x[2] = x0 + d
		y[2] = y[1]
		x[3] = x[2]
		y[3] = y0 - d
		canvas.Polygon(x, y, "stroke:"+color+";fill:"+color)
	case 3:
		dyu := d * 1.2
		dyl := dyu * 0.5
		dx := dyu * 0.866
		x := make([]float64, 3)
		y := make([]float64, 3)
		x[0] = x0
		y[0] = y0 + dyu
		x[1] = x0 - dx
		y[1] = y0 - dyl
		x[2] = x0 + dx
		y[2] = y0 - dyl
		canvas.Polygon(x, y, "stroke:"+color+";fill:"+color)
	}

}

/*
ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
*/
func GetUserInf(r *http.Request) {

	pt, _, _, ok := runtime.Caller(1) //	スタックトレースへのポインターを得る。1は一つ上のファンクション。

	fn := ""
	if !ok {
		fn = "unknown"
	}

	fn = runtime.FuncForPC(pt).Name()
	fna := strings.Split(fn, ".")

	ra := r.RemoteAddr
	ua := r.UserAgent()

	log.Printf("***** %s() from %s by %s\n", fna[len(fna)-1], ra, ua)
	//	fmt.Printf("%s() from %s by %s\n", fna[len(fna)-1], ra, ua)

}

// 入力フォーム画面
func HandlerTopForm(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/top.gtpl",
		"templates/top0.gtpl",
		"templates/top1.gtpl",
		"templates/top2.gtpl",
	))

	eventid := r.FormValue("eventid")
	suserno := r.FormValue("userno")
	if suserno == "" {
		suserno = "0"
	}
	userno, _ := strconv.Atoi(suserno)
	log.Printf("      eventid=%s userno=%d\n", eventid, userno)

	if eventid == "" {

		// マップを展開してテンプレートを出力する
		eventlist, _ := SelectLastEventList()
		if err := tpl.ExecuteTemplate(w, "top.gtpl", eventlist); err != nil {
			log.Println(err)
		}

		//	イベントでポイント比較の基準となる配信者（nobasis）のリストを取得する
		userlist, status := SelectUserList()
		if status == 0 {

			userlist[0].Userlongname = "ポイントの基準となる配信者が設定されていない"
			for i := 0; i < len(userlist); i++ {
				if userlist[i].Userno == userno {
					userlist[i].Selected = "Selected"
				} else {
					userlist[i].Selected = ""
				}
			}

			eventlist, _ = SelectEventList(userno)
			for i := 0; i < len(eventlist); i++ {
				if eventlist[i].EventID == eventid {
					eventlist[i].Selected = "Selected"
				} else {
					eventlist[i].Selected = ""
				}
			}
		}
		// マップを展開してテンプレートを出力する
		if err := tpl.ExecuteTemplate(w, "top0.gtpl", userlist); err != nil {
			log.Println(err)
		}
		if err := tpl.ExecuteTemplate(w, "top1.gtpl", eventlist); err != nil {
			log.Println(err)
		}
	} else {
		//	eventinf, _ := SelectEventInf(eventid)
		srdblib.Tevent = "event"
		eventinf, err := srdblib.SelectFromEvent(eventid)
		if err != nil {
			//	DBの処理でエラーが発生した。
			return
		} else if eventinf == nil {
			//	指定した eventid のイベントが存在しない。
			return
		}
		Event_inf = *eventinf

		if err := tpl.ExecuteTemplate(w, "top2.gtpl", eventinf); err != nil {
			log.Println(err)
		}
	}

}

func HandlerListLevel(w http.ResponseWriter, req *http.Request) {

	GetUserInf(req)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/list-level.gtpl"))

	userno, _ := strconv.Atoi(req.FormValue("userno"))
	levelonly, _ := strconv.Atoi(req.FormValue("levelonly"))
	log.Printf("***** HandlerListLevel() called. userno=%d, levelonly=%d\n", userno, levelonly)

	RoomLevelInf, _ := SelectRoomLevel(userno, levelonly)

	if err := tpl.ExecuteTemplate(w, "list-level.gtpl", RoomLevelInf); err != nil {
		log.Println(err)
	}
}

func HandlerListLast(w http.ResponseWriter, req *http.Request) {

	GetUserInf(req)

	status := 0

	var list_last struct {
		Detail    string
		Isover    string
		Scorelist []CurrentScore
	}

	// テンプレートをパースする
	//	tpl := template.Must(template.ParseFiles("templates/list-cntrb-h1.gtpl","templates/list-cntrb-h2.gtpl","templates/list-cntrb.gtpl"))
	funcMap := template.FuncMap{
		//	3桁ごとに","を挿入する
		"Comma": func(i int) string { return humanize.Comma(int64(i)) },
		//	イベントIDがブロックIDを含む場合はそれを取り除く。
		"DelBlockID": func(eid string) string {
			eia := strings.Split(eid, "?")
			if len(eia) == 2 {
				return eia[0]
			} else {
				return eid
			}
		},
	}
	tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/list-last.gtpl", "templates/list-last_h.gtpl"))

	eventid := req.FormValue("eventid")
	userno := req.FormValue("userno")
	list_last.Detail = req.FormValue("detail")
	log.Printf("      eventid=%s, detail=%s\n", eventid, list_last.Detail)
	//	Event_inf, _ = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		status = -1
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		status = -2
		return
	}
	Event_inf = *eventinf

	tdata, eventname, period, scorelist, status := SelectCurrentScore(eventid)
	list_last.Scorelist = scorelist
	for i := 0; i < len(scorelist); i++ {
		switch scorelist[i].Roomgenre {
		case "Voice Actors & Anime":
			scorelist[i].Roomgenre = "VA&A"
		case "Talent Model":
			scorelist[i].Roomgenre = "Tl/Md"
		case "Comedians/Talk Show":
			scorelist[i].Roomgenre = "Cm/TS"
		default:
		}
	}

	//	tnext := tdata.Add(5 * time.Minute)
	tnext := tdata.Add(time.Duration(Event_inf.Intervalmin) * time.Minute) //	0101G5
	//	treload := tnext.Add(5 * time.Second)
	treload := tnext.Add(10 * time.Second)

	values := map[string]string{
		"Eventid":         eventid,
		"userno":          userno,
		"UpdateTime":      "データ取得時刻：　" + tdata.Format("2006/01/02 15:04:05"),
		"NextTime":        "次のデータ取得は　" + tnext.Format("15:04:05") + "　に予定されています。",
		"ReloadTime":      "画面のリロードが　" + treload.Format("15:04:05") + "　頃に行われます。",
		"SecondsToReload": fmt.Sprintf("%d", int(time.Until(treload).Seconds()+5)),
		"EventName":       eventname,
		"Period":          period,
		"Detail":          list_last.Detail,
		"Maxpoint":        fmt.Sprintf("%d", Event_inf.Maxpoint),
		"Gscale":          fmt.Sprintf("%d", Event_inf.Gscale),
	}

	if time.Since(tdata) > 5*time.Minute {
		log.Printf("Application stopped or the event is over. status = %d\n", status)
		values["NextTime"] = "表示されているデータは最新ではありません。"
		values["ReloadTime"] = "もうしわけありませんがデータ取得が復旧するまでしばらくお待ちください。"
		values["SecondsToReload"] = "300"
	}
	if status != 0 {
		log.Printf("GetCurrentScore() returned %d.\n", status)
		values["UpdateTime"] = "データが取得できませんでした。"
		values["NextTime"] = "もうしわけありませんがしばらくお待ち下さい。"
		values["ReloadTime"] = ""
		values["SecondsToReload"] = "300"
	}
	if time.Now().After(Event_inf.End_time) {
		log.Printf("Application stopped or the event is over. status = %d\n", status)
		values["NextTime"] = "イベントは終了しています。"
		values["ReloadTime"] = ""
		values["SecondsToReload"] = "3600"

		list_last.Isover = "1"
	}
	if time.Now().Before(Event_inf.Start_time) {
		values["NextTime"] = "イベントはまだ始まっていません。"
		values["ReloadTime"] = ""
	}
	log.Printf("Values=%v", values)
	if err := tpl.ExecuteTemplate(w, "list-last_h", values); err != nil {
		log.Println(err)
	}
	if status != 0 {
		fmt.Fprintf(w, "</body>\n</html>\n")
		return
	}
	if err := tpl.ExecuteTemplate(w, "list-last", list_last); err != nil {
		log.Println(err)
	}
}
func HandlerGraphTotal(w http.ResponseWriter, req *http.Request) {

	GetUserInf(req)

	eventid := req.FormValue("eventid")
	//	maxpoint, _ := strconv.Atoi(req.FormValue("maxpoint"))
	smaxpoint := req.FormValue("maxpoint")
	maxpoint, _ := strconv.Atoi(smaxpoint)
	if maxpoint < 10000 {
		maxpoint = 0
		smaxpoint = "0"
	}
	sgscale := req.FormValue("gscale")
	if sgscale == "" || sgscale == "0" {
		sgscale = "100"
	}
	gscale, _ := strconv.Atoi(sgscale)
	/*
		gschk100 := ""
		gschk90 := ""
		gschk80 := ""
		gschk70 := ""
		switch sgscale {
		case "100":
			gschk100 = "checked"
		case "90":
			gschk90 = "checked"
		case "80":
			gschk80 = "checked"
		case "70":
			gschk70 = "checked"
		default:
			gschk100 = "checked"
		}
	*/

	log.Printf("      eventid=%s maxpoint=%d(%s)\n", eventid, maxpoint, smaxpoint)
	filename, _ := GraphTotalPoints(eventid, maxpoint, gscale)
	if Serverconfig.WebServer == "nginxSakura" {
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	} else if Serverconfig.WebServer == "Apache2Ubuntu" {
		filename = "/public/" + filename
	}

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-total.gtpl"))

	// テンプレートに出力する値をマップにセット
	/*
		values := map[string]string{
			"filename": req.FormValue("FileName"),
		}
	*/
	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
		"maxpoint": smaxpoint,
		"gscale":   sgscale,
	}

	// マップを展開してテンプレートを出力する
	if err := tpl.ExecuteTemplate(w, "graph-total.gtpl", values); err != nil {
		log.Println(err)
	}
}

func HandlerCsvTotal(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/csv-total.gtpl"))
	values := map[string]string{
		"function": "獲得ポイントの推移（CSV）",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "csv-total.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerGraphDfr(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-dfr.gtpl"))
	values := map[string]string{
		"function": "獲得ポイントの差の推移（グラフ）",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "graph-dfr.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerGraphPerday(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-perday.gtpl"))

	eventid := r.FormValue("eventid")

	log.Printf("      called. eventid=%s\n", eventid)

	ppointperday, _ := MakePointPerDay(eventid)

	filename, _ := GraphPerDay(eventid, ppointperday)
	if Serverconfig.WebServer == "nginxSakura" {
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	} else if Serverconfig.WebServer == "Apache2Ubuntu" {
		filename = "/public/" + filename
	}

	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
	}

	if err := tpl.ExecuteTemplate(w, "graph-perday.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerListPerday(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/list-perday.gtpl"))

	eventid := r.FormValue("eventid")

	log.Printf("      eventid=%s\n", eventid)

	pointperday, _ := MakePointPerDay(eventid)

	if err := tpl.ExecuteTemplate(w, "list-perday.gtpl", *pointperday); err != nil {
		log.Println(err)
	}
}

func HandlerGraphPerslot(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/graph-perslot.gtpl"))

	eventid := r.FormValue("eventid")
	log.Printf("      eventid=%s\n", eventid)

	perslotinflist, _ := MakePointPerSlot(eventid)

	filename, _ := GraphPerSlot(eventid, &perslotinflist)
	if Serverconfig.WebServer == "nginxSakura" {
		rootPath := os.Getenv("SCRIPT_NAME")
		rootPathFields := strings.Split(rootPath, "/")
		log.Printf("[%s] [%s] [%s]\n", rootPathFields[0], rootPathFields[1], rootPathFields[2])
		filename = "/" + rootPathFields[1] + "/public/" + filename
	} else if Serverconfig.WebServer == "Apache2Ubuntu" {
		filename = "/public/" + filename
	}

	values := map[string]string{
		"filename": filename,
		"eventid":  eventid,
	}

	if err := tpl.ExecuteTemplate(w, "graph-perslot.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerListPerslot(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/list-perslot1.gtpl",
		"templates/list-perslot2.gtpl",
	))

	eventid := r.FormValue("eventid")
	//	Event_inf, _ = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		return
	}
	Event_inf = *eventinf

	log.Printf("      eventid=%s\n", eventid)

	if err := tpl.ExecuteTemplate(w, "list-perslot1.gtpl", Event_inf); err != nil {
		log.Println(err)
	}

	perslotinflist, _ := MakePointPerSlot(eventid)

	if err := tpl.ExecuteTemplate(w, "list-perslot2.gtpl", perslotinflist); err != nil {
		log.Println(err)
	}

}

func HandlerEditUser(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/edit-user1.gtpl",
		"templates/edit-user2.gtpl",
		"templates/edit-user3.gtpl",
	))

	userid := r.FormValue("userid")
	eventid := r.FormValue("eventid")
	longname := r.FormValue("longname")
	shortname := r.FormValue("shortname")
	istarget := r.FormValue("istarget")
	graph := r.FormValue("graph")
	iscntrbpoint := r.FormValue("iscntrbpoint")
	color := r.FormValue("color")

	//	Event_inf, _ = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, err := srdblib.SelectFromEvent(eventid)
	if err != nil {
		//	DBの処理でエラーが発生した。
		return
	} else if eventinf == nil {
		//	指定した eventid のイベントが存在しない。
		return
	}
	Event_inf = *eventinf

	fnc := r.FormValue("func")

	log.Printf("      func=%s eventid=%s userid=%s\n", fnc, eventid, userid)

	switch fnc {
	case "newuser":
		//	新規配信者の追加があるとき

		roominf, status := GetRoomInfoAndPoint(eventid, userid, fmt.Sprintf("%d", Event_inf.Nobasis))
		if status == 0 {
			tnow := time.Now().Truncate(time.Second)
			InsertIntoOrUpdateUser(tnow, eventid, roominf)
			InsertIntoEventUser(0, eventid, roominf)
			UpdateEventuserSetPoint(eventid, roominf.ID, roominf.Point)

		} else {
			log.Printf("GetAndUpdateRoomInfoAndPoint() returned %d", status)
		}

	case "deleteuser":
		//	削除ボタンが押されたとき
	default:
		//	（更新ボタンが押された配信者がいたらそのデータを更新した上で）参加配信者のリストを表示する。
		if userid != "" {
			UpdateRoomInf(eventid, userid, longname, shortname, istarget, graph, color, iscntrbpoint)
		}
	}

	//	log.Printf(" eventid=%s, userno=%s, longname=%s, shortname=%s, istarget=%s, graph=%s, color=%s\n",
	//		eventid, userno, longname, shortname, istarget, graph, color)

	var roominfolist RoomInfoList

	eventname, _ := SelectEventRoomInfList(eventid, &roominfolist)
	for i := 0; i < len(roominfolist); i++ {
		switch roominfolist[i].Genre {
		case "Voice Actors & Anime":
			roominfolist[i].Genre = "VA&A"
		case "Talent Model":
			roominfolist[i].Genre = "Tl/Md"
		case "Comedians/Talk Show":
			roominfolist[i].Genre = "Cm/TS"
		default:
		}
	}

	values := map[string]string{
		"Eventid":   eventid,
		"Eventname": eventname,
		"Period":    Event_inf.Period,
	}

	if err := tpl.ExecuteTemplate(w, "edit-user1.gtpl", values); err != nil {
		log.Println(err)
	}

	if err := tpl.ExecuteTemplate(w, "edit-user2.gtpl", roominfolist); err != nil {
		log.Println(err)
	}

	if err := tpl.ExecuteTemplate(w, "edit-user3.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerNewUser(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/new-user.gtpl"))

	eventid := r.FormValue("eventid")

	log.Printf("      eventid=%s\n", eventid)

	roomid := r.FormValue("roomid")
	log.Printf("eventid=%s, roomid=%s\n", eventid, roomid)

	//	eventno, eventname, period := SelectEventNoAndName(eventid)
	//	log.Printf("eventname=%s, period=%s\n", eventname, period)

	//	Event_inf, _ = SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, _ := srdblib.SelectFromEvent(eventid)
	Event_inf = *eventinf

	log.Printf("eventname=%s, period=%s\n", Event_inf.Event_name, Event_inf.Period)

	genre, rank, nrank, prank, level, followers, fans, fans_lst, roomname, roomurlkey, _, status := GetRoomInfoByAPI(roomid)
	log.Printf("genre=%s, level=%d, followers=%d, fans=%d, fans_lst=%d, roomname=%s, roomurlkey=%s, status=%d\n",
		genre, level, followers, fans, fans_lst, roomname, roomurlkey, status)

	userno, _ := strconv.Atoi(roomid)
	roominf, status := SelectRoomInf(userno)

	longname := roominf.Longname
	shortname := roominf.Shortname
	if status != 0 {
		longname = ""
		shortname = ""
	} else {
		_, _, status = SelectUserColor(userno, Event_inf.Event_ID)
	}

	values := map[string]string{
		"Event_ID":   eventid,
		"Event_name": Event_inf.Event_name,
		"Period":     Event_inf.Period,
		"Roomid":     roomid,
		"Roomname":   roomname,
		"Longname":   longname,
		"Shortname":  shortname,
		"Roomurlkey": roomurlkey,
		"Genre":      genre,
		"Rank":       rank,
		"Nrank":      nrank,
		"Prank":      prank,
		"Level":      fmt.Sprintf("%d", level),
		"Followers":  fmt.Sprintf("%d", followers),
		"Fans":       fmt.Sprintf("%d", fans),
		"Fans_lst":   fmt.Sprintf("%d", fans_lst),
		"Submit":     "submit",
		"Label":      "登録しない",
		"Msg1":       "の参加ルームとして",
		"Msg2":       "を登録しますか？",
		"Msg2color":  "black",
	}

	if status == 0 {
		values["Submit"] = "hidden"
		values["Label"] = "戻る"
		values["Msg1"] = "の参加ルームとして"
		values["Msg2"] = "すでに登録されています"
		values["Msg2color"] = "red"
	} else {

		if roomname == "" {
			values["Roomname"] = ""
			values["Roomurlkey"] = ""
			values["Genre"] = ""
			values["Nrank"] = ""
			values["Prank"] = ""
			values["Level"] = ""
			values["Followers"] = ""
			values["Fans"] = ""
			values["Fans_lst"] = ""
			values["Submit"] = "hidden"
			values["Label"] = "戻る"
			values["Msg1"] = ""
			values["Msg2"] = "指定したルームIDのルームは存在しません。"
			values["Msg2color"] = "red"
		} else {
			_, _, _, peventid := GetPointsByAPI(roomid)
			if peventid != eventid && time.Now().After(Event_inf.Start_time) && time.Now().Before(Event_inf.End_time) {
				values["Submit"] = "hidden"
				values["Label"] = "戻る"
				values["Msg1"] = ""
				values["Msg2"] = "指定したルームはこのイベントに参加していません。"
				values["Msg2color"] = "red"
				log.Printf("GetPointsByAPI() returned %s as eventid and eventid = %s\n", peventid, eventid)
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "new-user.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerParamLocal(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/param-local.gtpl"))
	values := map[string]string{
		"function": "イベントパラメータの設定",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "param-local.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerAddEvent(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/add-event1.gtpl",
		"templates/add-event2.gtpl",
		"templates/error.gtpl",
	))

	var eventinf *exsrapi.Event_Inf
	var roominfolist RoomInfoList

	eventid := r.FormValue("eventid")
	breg := r.FormValue("breg")
	ereg := r.FormValue("ereg")
	ibreg, _ := strconv.Atoi(breg)
	iereg, _ := strconv.Atoi(ereg)

	if r.FormValue("from") != "new-event" {
		//	eventinf, _ = SelectEventInf(eventid)
		srdblib.Tevent = "event"
		eventinf, _ = srdblib.SelectFromEvent(eventid)

		log.Println("***** HandlerAddEvent() Called. not 'from new-event'")
		log.Println(eventinf)
	} else {
		eventinf = &exsrapi.Event_Inf{}
		log.Println("***** HandlerAddEvent() Called. 'from new-event'")
		eventinf.Modmin, _ = strconv.Atoi(r.FormValue("modmin"))
		eventinf.Modsec, _ = strconv.Atoi(r.FormValue("modsec"))

		intervalmin, _ := strconv.Atoi(r.FormValue("intervalmin"))
		switch intervalmin {
		case 5, 6, 10, 15, 20, 30, 60:
			eventinf.Intervalmin = intervalmin
		default:
			eventinf.Intervalmin = 5
		}
		eventinf.Modmin = eventinf.Modmin % eventinf.Intervalmin //	不適切な入力に対する修正
		eventinf.Modsec = eventinf.Modsec % 60

		eventinf.Resethh, _ = strconv.Atoi(r.FormValue("resethh"))
		eventinf.Resetmm, _ = strconv.Atoi(r.FormValue("resetmm"))
		eventinf.Nobasis, _ = strconv.Atoi(r.FormValue("nobasis"))
		eventinf.Target, _ = strconv.Atoi(r.FormValue("target"))
		eventinf.Maxdsp, _ = strconv.Atoi(r.FormValue("maxdsp"))
		eventinf.Cmap, _ = strconv.Atoi(r.FormValue("cmap"))
	}
	eventinf.Fromorder = ibreg
	eventinf.Toorder = iereg

	Event_inf = *eventinf

	log.Println("before GetAndInsertEventRoomInfo()")
	log.Println(eventinf)

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	starttimeafternow, status := GetAndInsertEventRoomInfo(client, eventid, ibreg, iereg, eventinf, &roominfolist)
	if status != 0 {

		values := map[string]string{
			"Msg001":   "入力したイベントID( ",
			"Msg002":   " )をもつイベントは存在しません！",
			"ReturnTo": "top",
			"Eventid":  eventid,
		}
		if err := tpl.ExecuteTemplate(w, "error.gtpl", values); err != nil {
			log.Println(err)
		}

	} else {

		if err := tpl.ExecuteTemplate(w, "add-event1.gtpl", eventinf); err != nil {
			log.Println(err)
		}
		if starttimeafternow {
			if iereg > len(roominfolist) {
				iereg = len(roominfolist)
			}
			if err := tpl.ExecuteTemplate(w, "add-event2.gtpl", roominfolist[ibreg-1:iereg]); err != nil {
				log.Println(err)
			}
		} else {
			if err := tpl.ExecuteTemplate(w, "add-event2.gtpl", roominfolist); err != nil {
				log.Println(err)
			}
		}
	}

}

/*
MakeSampleTime()
獲得ポイントを取得するタイミングをランダムに返す

5分に一回を前提として、240秒±40秒のように設定する。
*/
func MakeSampleTime(
	cval int, // ex. 240
	cvar int, // ex. 40
) (stm, sts int) {

	st := cval + int(time.Now().UnixNano()%int64(cvar*2)) - cvar

	stm = st / 60
	sts = st % 60

	return stm, sts
}
func HandlerNewEvent(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/new-event0.gtpl",
		"templates/new-event1.gtpl",
		"templates/new-event2.gtpl",
	))

	eventid := r.FormValue("eventid")
	suserno := r.FormValue("userno")
	userno, _ := strconv.Atoi(suserno)

	log.Printf("      eventid=%s\n", eventid)

	stm, sts := MakeSampleTime(240, 40)

	values := map[string]string{
		"Eventid":   r.FormValue("eventid"),
		"Eventname": "",
		"Period":    "",
		"Noroom":    "",
		"Msgcolor":  "blue",

		"Stm": fmt.Sprintf("%d", stm),
		"Sts": fmt.Sprintf("%d", sts),
	}

	var eventinf exsrapi.Event_Inf

	eia := strings.Split(eventid, "?")
	if len(eia) == 2 {
		eventid = eia[0]
	}

	status := GetEventInf(eventid, &eventinf)
	if status == -1 {
		values["Msg"] = "このイベントはすでに登録されています。"
		values["Submit"] = "hidden"
		values["Msgcolor"] = "red"
		//	Event_inf, _ = SelectEventInf(eventid)
		srdblib.Tevent = "event"
		eventinf, _ := srdblib.SelectFromEvent(eventid)
		Event_inf = *eventinf

		values["Eventname"] = Event_inf.Event_name
		values["Period"] = Event_inf.Period
	} else if status == -2 {
		values["Msg"] = "指定したIDのイベントは存在しません"
		values["Submit"] = "hidden"
		values["Msgcolor"] = "red"
	} else if status < -2 {
		values["Msg"] = "イベント情報を取得できませんでした（エラーコード＝" + fmt.Sprintf("%d", status) + "）"
		values["Submit"] = "hidden"
		values["Msgcolor"] = "red"
	} else {
		values["Msg"] = "このイベントを登録しますか？"
		values["Submit"] = "submit"
		values["Eventname"] = eventinf.Event_name
		values["Period"] = eventinf.Period
		values["Noroom"] = "　" + humanize.Comma(int64(eventinf.NoEntry))
	}
	/*
		var Eventinflist []Event_Inf
		GetEventListByAPI(&Eventinflist)
	*/

	userlist, _ := SelectUserList()
	userlist[0].Userlongname = "基準とする配信者を設定しない"
	for i := 0; i < len(userlist); i++ {
		if userlist[i].Userno == userno {
			userlist[i].Selected = "Selected"
		} else {
			userlist[i].Selected = ""
		}
	}

	if err := tpl.ExecuteTemplate(w, "new-event0.gtpl", values); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "new-event1.gtpl", userlist); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "new-event2.gtpl", values); err != nil {
		log.Println(err)
	}

}

func HandlerParamEvent(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles(
		"templates/param-event0.gtpl",
		"templates/param-event1.gtpl",
		"templates/param-event2.gtpl",
	))

	eventid := r.FormValue("eventid")

	log.Printf("      eventid=%s\n", eventid)

	//	eventinf, _ := SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, _ := srdblib.SelectFromEvent(eventid)
	Event_inf = *eventinf

	userlist, _ := SelectEventuserList(eventid)
	for i := 0; i < len(userlist); i++ {
		if userlist[i].Userno == eventinf.Nobasis {
			userlist[i].Selected = "Selected"
		} else {
			userlist[i].Selected = ""
		}
	}

	if err := tpl.ExecuteTemplate(w, "param-event0.gtpl", eventinf); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "param-event1.gtpl", userlist); err != nil {
		log.Println(err)
	}
	if err := tpl.ExecuteTemplate(w, "param-event2.gtpl", eventinf); err != nil {
		log.Println(err)
	}

}

func HandlerParamEventC(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/param-eventc.gtpl"))
	eventid := r.FormValue("eventid")
	log.Printf("      eventid=%s\n", eventid)

	//	eventinf, _ := SelectEventInf(eventid)
	srdblib.Tevent = "event"
	eventinf, _ := srdblib.SelectFromEvent(eventid)
	Event_inf = *eventinf

	//	log.Println(eventinf)

	eventinf.Fromorder, _ = strconv.Atoi(r.FormValue("fromorder"))
	eventinf.Toorder, _ = strconv.Atoi(r.FormValue("toorder"))
	eventinf.Modmin, _ = strconv.Atoi(r.FormValue("modmin"))
	eventinf.Modsec, _ = strconv.Atoi(r.FormValue("modsec"))

	intervalmin, _ := strconv.Atoi(r.FormValue("intervalmin"))
	switch intervalmin {
	case 5, 6, 10, 15, 20, 30, 60:
		eventinf.Intervalmin = intervalmin
	default:
		eventinf.Intervalmin = 5
	}
	eventinf.Modmin = eventinf.Modmin % eventinf.Intervalmin //	不適切な入力に対する修正
	eventinf.Modsec = eventinf.Modsec % 60

	eventinf.Resethh, _ = strconv.Atoi(r.FormValue("resethh"))
	eventinf.Resetmm, _ = strconv.Atoi(r.FormValue("resetmm"))
	eventinf.Nobasis, _ = strconv.Atoi(r.FormValue("nobasis"))
	eventinf.Target, _ = strconv.Atoi(r.FormValue("target"))
	eventinf.Maxdsp, _ = strconv.Atoi(r.FormValue("maxdsp"))
	eventinf.Cmap, _ = strconv.Atoi(r.FormValue("cmap"))

	//	UpdateEventInf(&eventinf)
	UpdateEventInf(eventinf)
	//	log.Println(eventinf)

	if err := tpl.ExecuteTemplate(w, "param-eventc.gtpl", eventinf); err != nil {
		log.Println(err)
	}

}

func HandlerParamGlobal(w http.ResponseWriter, r *http.Request) {

	GetUserInf(r)

	// テンプレートをパースする
	tpl := template.Must(template.ParseFiles("templates/param-global.gtpl"))
	values := map[string]string{
		"function": "共通パラメータの設定",
		"comment":  "この機能は現在作成中です。",
	}
	if err := tpl.ExecuteTemplate(w, "param-global.gtpl", values); err != nil {
		log.Println(err)
	}

}
