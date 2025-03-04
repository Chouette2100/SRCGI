// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	// "fmt"
	//	"html"
	"log"

	//	"math/rand"
	// "sort"
	"strconv"
	"strings"
	"time"
	//	"os"

	"runtime"
	"sync"

	// "encoding/json"

	//	"html/template"
	"net/http"

	// "database/sql"

	"encoding/json"

	// _ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	//	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

/*
ファンクション名とリモートアドレス、ユーザーエージェントを表示する。
*/
//	var Localhost bool
type KV struct {
	K string
	V []string
}

var LAlog sync.Map = sync.Map{}

const wait = time.Millisecond * 1500

func GetUserInf(r *http.Request) (
	ra string,
	ua string,
	isallow bool,
) {

	isallow = true

	pt, _, _, ok := runtime.Caller(1) //	スタックトレースへのポインターを得る。1は一つ上のファンクション。

	var fn string
	if !ok {
		fn = "unknown"
	}

	fn = runtime.FuncForPC(pt).Name()
	fna := strings.Split(fn, ".")

	rap := r.RemoteAddr
	rapa := strings.Split(rap, ":")
	if rapa[0] != "[" {
		ra = rapa[0]
	} else {
		ra = "127.0.0.1"
	}
	ua = r.UserAgent()

	log.Printf("  *** %s() from %s by %s\n", fna[len(fna)-1], ra, ua)
	//	log.Printf("%s() from %s by %s\n", fna[len(fna)-1], ra, ua)

	/*
		if !IsAllowIp(ra) {
			log.Printf("%s is on the Blacklist(%s)", ra, ua)
			isallow = false
			return
		}
	*/

	//	パラメータを表示する
	if err := r.ParseForm(); err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	var al srdblib.Accesslog
	al.Ts = time.Now().Truncate(time.Second)
	al.Handler = fna[len(fna)-1]
	al.Remoteaddress = ra
	al.Useragent = ua

	kvlist := make([]KV, len(r.Form))
	i := 0
	for kvlist[i].K, kvlist[i].V = range r.Form {
		log.Printf("%12v : %v\n", kvlist[i].K, kvlist[i].V)
		switch kvlist[i].K {
		case "eventid":
			al.Eventid = kvlist[i].V[0]
		case "userno", "userid", "user_id", "roomid":
			al.Roomid, _ = strconv.Atoi(kvlist[i].V[0])
		default:
		}
		i++
	}
	jd, err := json.Marshal(kvlist)
	if err != nil {
		log.Printf(" GetUserInf(): %s\n", err.Error())
	}
	al.Formvalues = string(jd)

	Chlog <- &al

	/*
		err = srdblib.Dbmap.Insert(&al)
		if err != nil {
			log.Printf(" GetUserInf(): %s\n", err.Error())
		}
	*/

	twait := wait

	// 	クローラーの場合は待ち時間を10秒とする
	// if strings.Contains(ua, "SemrushBot") || strings.Contains(ua, "Googlebot") {
	if Regexpbots.MatchString(ua) {
		twait = 10 * time.Second
	}

	if ll, ok := LAlog.Load(ra); ok {
		// 待ち時間情報がある場合
		tnow := time.Now()
		na := ll.(time.Time)
		if tnow.After(na) {
			// 待ち時間を過ぎている場合
			na = tnow.Add(twait)
			LAlog.Store(ra, na)
			log.Printf("     === %20s %s set Nextaccess to %s ( tnow =%s)\n",
				fna[len(fna)-1], ra, na.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
		} else {
			// 待ち時間を過ぎていない場合
			nna := na.Add(twait)
			LAlog.Store(ra, nna)
			if nna.Sub(tnow) > twait*6 {
				// 待ち時間が待ち時間単位の6倍を超えた場合は処理を許可しない
				// これは、待ち時間が長くなるのは、処理が終わる前にリクエストが来ていることを意味する。
				log.Printf("     === %20s %s set Nextaccess to %s ( tnow =%s)\n", fna[len(fna)-1],
					ra, nna.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
				isallow = false
				return
			}
			log.Printf("  === %20s %s set Nextaccess to %s ( tnow =%s)\n",
				fna[len(fna)-1], ra, nna.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
			time.Sleep(na.Sub(tnow))
		}
	} else {
		// 待ち時間情報がない場合
		tnow := time.Now()
		na := tnow.Add(twait)
		// LAlog.Store(ra, Lastaccess{na})
		LAlog.Store(ra, na)
		log.Printf("     === %20s %s set Nextaccess to %s ( tnow =%s)\n",
			fna[len(fna)-1], ra, na.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
	}

	return
}

// func logWorker(db *sql.DB, logCh chan string, done chan struct{}) {
func LogWorker() {
	for {
		al := <-Chlog
		if err := srdblib.Dbmap.Insert(al); err != nil {
			log.Printf(" GetUserInf(): %s\n", err.Error())
		}
	}
}
