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

	//	"github.com/Chouette2100/exsrapi"
	// "github.com/Chouette2100/srapi"
	"github.com/Chouette2100/srdblib"
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

	if ll, ok := LAlog.Load(ra); ok {
		tnow := time.Now()
		na := ll.(time.Time)
		if tnow.After(na) {
			na = tnow.Add(wait)
			LAlog.Store(ra, na)
			log.Printf("  === %20s %s set Nextaccess to %s ( tnow =%s)\n",
				fna[len(fna)-1],ra, na.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
		} else {
			nna := na.Add(wait)
			LAlog.Store(ra, nna)
			log.Printf("  === %20s %s set Nextaccess to %s ( tnow =%s)\n",
				fna[len(fna)-1],ra, nna.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
			time.Sleep(na.Sub(tnow))
		}
	} else {
		tnow := time.Now()
		na := tnow.Add(wait)
		// LAlog.Store(ra, Lastaccess{na})
		LAlog.Store(ra, na)
		log.Printf("  === %20s %s set Nextaccess to %s ( tnow =%s)\n",
			fna[len(fna)-1],ra, na.Format("2006-01-02 15:04:05.000"), tnow.Format("2006-01-02 15:04:05.000"))
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
