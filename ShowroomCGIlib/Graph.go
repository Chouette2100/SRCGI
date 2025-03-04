// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package ShowroomCGIlib

import (
	//	"SRCGI/ShowroomCGIlib"
	//	"bufio"
	// "bytes"
	//	"fmt"
	// "html"
	//	"log"

	//	"math/rand"
	// "sort"
	//	"strconv"
	//	"strings"
	//	"time"
	//	"os"


	// "runtime"

	// "encoding/json"

	//	"html/template"
	//	"net/http"

	// "database/sql"

	// _ "github.com/go-sql-driver/mysql"

	// "github.com/PuerkitoBio/goquery"

	svg "github.com/ajstarks/svgo/float"

	//	"github.com/dustin/go-humanize"

	//	"github.com/goark/sshql"
	//	"github.com/goark/sshql/mysqldrv"

	//	"github.com/Chouette2100/exsrapi/v2"
	// "github.com/Chouette2100/srapi/v2"
	//	"github.com/Chouette2100/srdblib/v2"
)

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
		{1000, 50, 5},
		{1500, 50, 5},
		{2000, 100, 5},
		{3000, 100, 5},
		{5000, 200, 5},
		{7000, 200, 5},
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
		{3.1, 12, 1},
		{10.1, 8, 1},
		{20.1, 4, 1},
		{40.1, 2, 3},
		{80.1, 1, 4},
		{160.1, -3, 4},
		{350.1, -7, 4},
		{700.1, -14, 4},
		{1400.1, -35, 4},
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
		xscaled = -70
		xscalet = 4
	}

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

