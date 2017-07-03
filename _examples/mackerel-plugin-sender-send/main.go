package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/nissy/mackerel-plugin-timeline"
)

const allcount = "allcount"

var (
	replacer = strings.NewReplacer("  ", " ")
)

type SendCount struct {
	count map[string]uint64
}

type send struct {
	time                        time.Time
	result                      string
	statusCode                  string
	transactionID               string
	logHeaderSection            string
	envelopeFromAddress         string
	envelopeToAddress           string
	envelopeToAddressDomain     string
	toIPAddressPort             string
	souceIPAdressPort           string
	smtpSessionConnectionTiming string
	message                     string
}

func parseSend(line string) send {
	sd := send{}
	var i int

	if s := strings.Split(replacer.Replace(line), " "); len(s) == 11 || len(s) == 12 {
		sd.time, _ = time.ParseInLocation(*timeline.TimeLayout, strings.Join(s[0:2], " "), timeline.Location)
		sd.result = s[2]
		sd.statusCode = s[3]
		sd.transactionID = s[4]

		// logHeaderSection off
		if len(s) == 12 {
			sd.logHeaderSection = s[5]
			i++
		}

		sd.envelopeFromAddress = s[5+i]
		sd.envelopeToAddress = s[6+i]
		sd.envelopeToAddressDomain = strings.Split(s[6+i], "@")[1]
		sd.toIPAddressPort = s[7+i]
		sd.souceIPAdressPort = s[8+i]
		sd.smtpSessionConnectionTiming = s[9+i]
		sd.message = s[10+i]
	}

	return sd
}

func (sdc SendCount) ParseTime(line string) time.Time {
	return parseSend(line).time
}

func (sdc SendCount) ToConut(line string) error {
	if send := parseSend(line); timeline.TimeEnd.Unix() < send.time.Unix() {
		return errors.New("End time is small")
	} else if timeline.TimeStart.Unix() <= send.time.Unix() {
		sdc.count[send.result]++
		sdc.count[send.envelopeToAddressDomain]++
		sdc.count[allcount]++
	}

	return nil
}

func (sdc SendCount) ToMetrics() map[string]interface{} {
	m := map[string]interface{}{}

	for _, v := range sdc.ToGraph()["send.result"].Metrics {
		m[v.Name] = sdc.count[v.Label]
	}

	others := sdc.count[allcount]

	for _, v := range sdc.ToGraph()["send.to_address_domain"].Metrics {
		m[v.Name] = sdc.count[v.Label]
		others = others - sdc.count[v.Label]
	}

	m["others"] = others

	return m
}

func (sdc SendCount) ToGraph() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"send.result": {
			Label: "sender send results",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{
					Name:  "sent",
					Label: "sent",
				},
				{
					Name:  "retry",
					Label: "retry",
				},
				{
					Name:  "faild",
					Label: "faild",
				},
			},
		},
		"send.to_address_domain": {
			Label: "sender send to_address_domain",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{
					Name:  "docomo",
					Label: "docomo.ne.jp",
				},
				{
					Name:  "ezweb",
					Label: "ezweb.ne.jp",
				},
				{
					Name:  "softbank",
					Label: "softbank.ne.jp",
				},
				{
					Name:  "isoftbank",
					Label: "i.softbank.jp",
				},
				{
					Name:  "disney",
					Label: "disney.ne.jp",
				},
				{
					Name:  "gmail",
					Label: "gmail.com",
				},
				{
					Name:  "yahoo",
					Label: "yahoo.co.jp",
				},
				{
					Name:  "icloud",
					Label: "icloud.com",
				},
				{
					Name:  "others",
					Label: "others",
				},
			},
		},
	}
}

func main() {
	pl := timeline.NewPlugin(
		SendCount{
			count: make(map[string]uint64),
		},
	)

	os.Exit(exitcode(pl.Run()))
}

func exitcode(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	return 0
}
