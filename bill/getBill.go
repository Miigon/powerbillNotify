package bill

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miigon/powerbillNotify/conf"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type PowerBill struct {
	RoomId      string
	PowerLeft   float64
	TotalUsage  float64
	TotalBought float64
	Time        time.Time
}

const timeFormat = "2006-01-02"

func GetPowerBill(beginTime time.Time, endTime time.Time) (powerBills []PowerBill) {
	enc, err := simplifiedchinese.GB18030.NewEncoder().String(conf.Config.Building)
	if err != nil {
		panic("encoding roomname [" + conf.Config.RoomName + "] error: " + err.Error())
	}
	encodedBuilding := url.QueryEscape(enc)

	data := url.Values{
		"beginTime": {beginTime.Format(timeFormat)},
		"endTime":   {endTime.Format(timeFormat)},
		"type":      {"2"},
		"client":    {conf.Config.ClientField},
		"roomId":    {conf.Config.RoomId},
		"roomName":  {conf.Config.RoomName},
		"building":  {encodedBuilding},
	}

	resp, err := http.PostForm(conf.Config.ServerQueryUrl, data)
	if err != nil {
		panic("failed querying: " + err.Error())
	}
	resdata, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed reading response: " + err.Error())
	}
	return parseResultHTML(string(resdata))
}

const tableStartString = "<table width=\"100%\" summary=\"list of members in EE Studay\""
const tableEndString = "</table>"

var exp = regexp.MustCompile("<tr>[\\n\\s]+" +
	strings.Repeat("<td .*>[\\n\\s]+(.*)[\\n\\s]+</td>[\\n\\s]+", 6) + // 6 fields
	"</tr>")

func parseResultHTML(htmldata string) (powerBills []PowerBill) {
	htmldata = strings.ReplaceAll(htmldata, "\r", "")
	// find table first, not actually necessary since the only thing matching the regexp
	// just happens to be the actual data we want.
	tableStart := strings.Index(htmldata, tableStartString)
	tableEnd := tableStart + strings.Index(htmldata[tableStart:], tableEndString) + len(tableEndString)

	parseFloat := func(str string) float64 {
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			panic(err)
		}
		return f
	}

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	matches := exp.FindAllStringSubmatch(htmldata[tableStart:tableEnd], -1)
	for i, v := range matches {
		if i == 0 {
			continue
		}
		t, err := time.ParseInLocation("2006-01-02 15:04:05.9", v[6], location)
		if err != nil {
			panic(err)
		}
		powerBills = append(powerBills, PowerBill{
			RoomId:      v[2],
			PowerLeft:   parseFloat(v[3]),
			TotalUsage:  parseFloat(v[4]),
			TotalBought: parseFloat(v[5]),
			Time:        t,
		})

	}
	return
}
