// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/kasworld/jslog"
)

var refresh int = 3600

var starttime time.Time

var done chan struct{}

var bgExist bool

var lasttime time.Time

func main() {
	starttime = time.Now()
	queryv := GetQuery()

	if refreshQry := strings.TrimSpace(queryv.Get("refresh")); refreshQry != "" {
		value, err := strconv.ParseInt(refreshQry, 0, 64)
		if err == nil {
			refresh = int(value)
		}
	}

	if mvid := queryv.Get("mvid"); mvid != "" {
		setYoutube(mvid)
		bgExist = true
	} else if imgname := queryv.Get("bgimg"); imgname != "" {
		setBGImage(imgname)
		bgExist = true
	}
	displayFrame()
	<-done
}

func GetQuery() url.Values {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		jslog.Errorf("%v", err)
	}
	return u.Query()
}

func jsFrame(js.Value, []js.Value) interface{} {
	displayFrame()
	return nil
}

func displayFrame() {
	defer js.Global().Call("requestAnimationFrame", js.FuncOf(jsFrame))
	thistime := time.Now()

	if thistime.Sub(starttime) > time.Duration(refresh)*time.Second {
		js.Global().Get("location").Call("reload")
	}

	if lasttime.Second() == thistime.Second() {
		return
	}
	lasttime = thistime

	win := js.Global().Get("window")
	winW := win.Get("innerWidth").Float()
	winH := win.Get("innerHeight").Float()

	sizeRef := winW
	if sizeRef > winH {
		sizeRef = winH
	}

	clockFontSize := winW / 4
	if clockFontSize > winH/2.7 {
		clockFontSize = winH / 2.7
	}
	if bgExist {
		clockFontSize /= 1.5
	}
	updateTime(clockFontSize)

	dateFontSize := clockFontSize / 3.5
	updateDate(dateFontSize)

	calendarFontSize := sizeRef / 12
	if bgExist {
		calendarFontSize = sizeRef / 20
	}
	updateCalendar(calendarFontSize)

	updateDebugInfo(calendarFontSize / 2)
}

func setBGImage(imageurl string) {
	str := fmt.Sprintf(`
	<img src="%[1]s" style="width:100%%; height:100%%;">
	`, imageurl)

	jsObj := js.Global().Get("document").Call("getElementById", "bg")
	jsObj.Set("innerHTML", str)
}

func setYoutube(mvid string) {
	str := fmt.Sprintf(`
	<iframe frameborder="0" height="100%%" width="100%%" allow="autoplay"
	src="https://youtube.com/embed/%[1]s?autoplay=1&controls=0&fs=0&loop=1">
	  </iframe>
	`, mvid)

	jsObj := js.Global().Get("document").Call("getElementById", "bg")
	jsObj.Set("innerHTML", str)
}

func updateTime(fontSize float64) {
	jsObj := js.Global().Get("document").Call("getElementById", "time")
	jsObj.Get("style").Set("font-size", fmt.Sprintf("%.1fpx", fontSize))
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v", time.Now().Format("15:04:05"))
	jsObj.Set("innerHTML", buf.String())
}

func updateDebugInfo(fontSize float64) {
	jsObj := js.Global().Get("document").Call("getElementById", "debuginfo")
	jsObj.Get("style").Set("font-size", fmt.Sprintf("%.1fpx", fontSize))
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%.2f", time.Now().Sub(starttime).Seconds())
	jsObj.Set("innerHTML", buf.String())
}

func updateDate(fontSize float64) {
	jsObj := js.Global().Get("document").Call("getElementById", "date")
	jsObj.Get("style").Set("font-size", fmt.Sprintf("%.1fpx", fontSize))
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v", time.Now().Format("2006-01-02 Mon"))
	jsObj.Set("innerHTML", buf.String())
}

func updateCalendar(fontSize float64) {
	jsObj := js.Global().Get("document").Call("getElementById", "calendar")
	jsObj.Get("style").Set("font-size", fmt.Sprintf("%.1fpx", fontSize))
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "<table id=\"t01\">")

	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		fmt.Fprintf(&buf, "<td>%v</td>", weekday.String()[:2])
	}

	fmt.Fprintf(&buf, "<colgroup>")
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		if weekday == 0 {
			fmt.Fprintf(&buf, "<col style=\"width:14%%; \">")
		} else if weekday == 6 {
			fmt.Fprintf(&buf, "<col style=\"width:14%%; \">")
		} else {
			fmt.Fprintf(&buf, "<col style=\"width:14%%; \">")
		}
	}
	fmt.Fprintf(&buf, "<colgroup>")

	today := time.Now()
	now := time.Now()
	now = now.AddDate(0, 0, -int(now.Weekday())-7)
	for week := 0; week < 6; week++ {
		fmt.Fprintf(&buf, "<tr>")
		for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
			if now.Month() != today.Month() {
				fmt.Fprintf(&buf, "<td style=\"color:darkgray;\">%d</td>", now.Day())
			} else {
				if now.Day() != today.Day() {
					if weekday == 0 {
						fmt.Fprintf(&buf, "<td style=\"color:red;\">%d</td>", now.Day())
					} else if weekday == 6 {
						fmt.Fprintf(&buf, "<td style=\"color:SkyBlue;\">%d</td>", now.Day())
					} else {
						fmt.Fprintf(&buf, "<td style=\"color:white;\">%d</td>", now.Day())
					}
				} else {
					fmt.Fprintf(&buf, "<td style=\"color:orange;\">%d</td>", now.Day())
				}
			}
			now = now.AddDate(0, 0, 1)
		}
		fmt.Fprintf(&buf, "</tr>")
	}
	fmt.Fprintf(&buf, "</table>")
	jsObj.Set("innerHTML", buf.String())
}
