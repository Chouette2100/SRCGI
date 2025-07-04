// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bytes"
	"fmt"
	//	"html"
	"log"

	//	"math/rand"
	//	"sort"
	//	"strconv"
	//	"strings"
	"time"

	"bufio"
	"os"

	//	"runtime"

	//	"encoding/json"

	//	"html/template"
	//	"net/http"

	//	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	//	"github.com/PuerkitoBio/goquery"

	svg "github.com/ajstarks/svgo/float"

	"github.com/dustin/go-humanize"
	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"
	//	"github.com/Chouette2100/exsrapi/v2"
	//	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func Jtruncate(t time.Time) time.Time {
	// return t.Add(9*time.Hour).Truncate(24*time.Hour).Add(-9*time.Hour)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

type Xydata struct {
	X []float64
	Y []float64
}

func DrawLineGraph(
	filename string, //	（パスなし）ファイル名　ex. 000.svg
	title0 string, //	ex.	グラフタイトル
	title1 string, //	ex. イベント名
	title2 string, //	ex. 開催期間
	maxpoint int, //	データの最大値
	tmaxpoint int, //	y軸方向グラフ表示範囲を制限する
	target int, //	目標ポイント
	start_time time.Time, //	イベント開始時刻
	end_time time.Time, //	イベント終了時刻
	cmap int, // グラフの描画に使用するカラーマップ
	deltax float64, //	データ間隔がこの時間を超えたら接続しない(day)
	IDlist []int,
	xydata *[]Xydata,
) (
	err error,
) {

	start_date := float64(Jtruncate(start_time).Unix()) / 60.0 / 60.0 / 24.0
	end_date := float64(Jtruncate(end_time).Add(24*time.Hour).Unix()) / 60.0 / 60.0 / 24.0

	dperiod := end_date - start_date

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
	yupper := 0
	yscales := 0
	yscalel := 0

	if tmaxpoint != 0 {
		yupper, yscales, yscalel, _ = DetYaxScale(tmaxpoint - 1)
	} else if target > maxpoint {
		yupper, yscales, yscalel, _ = DetYaxScale(target - 1)
	} else {
		yupper, yscales, yscalel, _ = DetYaxScale(maxpoint)
	}

	yscale := -vheight / float64(yupper)

	//	log.Printf("yupper=%d yscale=%f dyl=%f\n", yupper, yscale, dyl)

	//	グラフタイトルとイベント情報を出力する
	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5-8*1.5), title0,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*2.5, title1,
		"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*8.0)+"px;fill:white;")

	canvas.Text(lwmargin+vwidth/2.0, uhmargin/2.0+bstroke*(2.5+8*1.5), title2,
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

	//	x軸に表示する値の上限値
	xupper := dperiod
	//	x軸に表示する値と座標幅の比（表示値１が座標のいくらに相当するか？）
	xscale := vwidth / float64(xupper)
	// xscaled > 0 のとき　一目盛を1日の何分の一にするか？
	// xscaled < 0 のとき、一目盛を1日の何倍にするか？
	xscaled, xscalet, _ := DetXaxScale(xupper)
	//	log.Printf("xupper=%f xscale=%f dxl=%f xscalet=%d\n", xupper, xscale, dxl, xscalet)

	//	一目盛の表示長さ
	dxl := 0.0
	if xscaled > 0 {
		dxl = 1.0 / float64(xscaled) * xscale
	} else {
		dxl = -1.0 * float64(xscaled) * xscale
	}
	tval := start_time
	xl := 0.0
	for i := 0; ; i++ {
		wstr := 0.15
		if xscaled > 0 && i%xscaled == 0 {
			//	xscaled > 0 のときは1日ごとに（00時に）表示を太くする
			wstr = 0.3
		}
		if xscaled > 0 && i%(xscaled*xscalet) == 0 || xscaled < 0 && i%xscalet == 0 {
			//	xscaled > 0 のときはxscalet日ごとに日付を表示する
			//	xscaled < 0 のときはxscaled * xscalet日ごとに日付を表示する
			if xscaled > -10 {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("1/2"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
			} else {
				canvas.Text(xorigin+xl, yorigin+bstroke*7.5, tval.Format("06/01/02"),
					"text-anchor:middle;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
			}

			if xscaled > 0 {
				tval = tval.AddDate(0, 0, xscalet)
			} else {
				tval = tval.AddDate(0, 0, -xscaled*xscalet)
			}
		}

		canvas.Line(xorigin+xl, yorigin, xorigin+xl, yorigin-vheight, "stroke=\"white\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*wstr)+"\"")
		xl += dxl
		if xl > vwidth+10 {
			break
		}
	}

	//	ターゲットラインを描画する
	if target != 0 {
		x1 := xorigin + (float64(start_time.Unix())/60.0/60.0/24.0-start_date)*xscale
		x2 := xorigin + (float64(end_time.Unix())/60.0/60.0/24.0-start_date)*xscale
		y1 := yorigin
		y2 := yorigin + float64(target)*yscale

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
	for idx := range IDlist {

		//	_, cvalue, _ := SelectUserColor(id, eventid)
		cvalue := Colormaplist[cmap][idx%len(Colormaplist[cmap])].Value

		//	x, y := SelectScoreList(id)
		x := &(*xydata)[j].X
		y := &(*xydata)[j].Y
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
			if (*x)[i]-tl > deltax && (*y)[i]-yl > 0.9 {
				//	次のデータとの間に欠測があり、かつ欠測の前後でデータが同一でないときはその部分の描画は行わない。
				//	canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				if k > 1 {
					canvas.Polyline(xo[0:k], yo[0:k], "fill=\"none\" stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
				} else {
					canvas.Circle(xo[0], yo[0], bstroke*1.5, "fill=\""+cvalue+"\" stroke=\""+cvalue+"\"")
				}

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
		} else {
			canvas.Circle(xo[0], yo[0], bstroke*1.5, "fill=\""+cvalue+"\" stroke=\""+cvalue+"\"")
		}

		xln := xorigin + vwidth + bstroke*30.0
		yln := yorigin - vheight + bstroke*10*float64(j)

		canvas.Line(xln, yln, xln+rwmargin/4.0, yln, "stroke=\""+cvalue+"\" stroke-width=\""+fmt.Sprintf("%.2f", bstroke*1.0)+"\"")
		//	canvas.Text(xln+rwmargin/3.0, yln+bstroke*2.5, fmt.Sprintf("%d", IDlist[j]),
		//		"text-anchor:start;font-size:"+fmt.Sprintf("%.1f", bstroke*5.0)+"px;fill:white;")
		//	longname, _, _, _, _, _, _, _, _, _, sts := SelectUserName(IDlist[j])
		//	if sts != 0 {
		//		longname = fmt.Sprintf("%d", IDlist[j])
		//	}
		var intrfc interface{}
		intrfc, err = srdblib.Dbmap.Get(srdblib.User{}, IDlist[j])
		if err != nil {
			err = fmt.Errorf("srdblib.Dbmap.Get(srdblib.User{}, %d) err=%w", IDlist[j], err)
			return
		}
		longname := intrfc.(*srdblib.User).Longname
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

	return

}
