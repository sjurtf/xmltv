package cmd

import (
	"log"
	"time"
	"xmltv-exporter/tv2"
	"xmltv-exporter/xmltv"
)

func MapEgp() {
	today := time.Now()
	futureOneWeek := time.Now().Add(time.Hour * 24 * 3)

	var t = today
	for t.Before(futureOneWeek) {
		channels := tv2.FetchEpg(t)

		for _, channel := range channels {
			xmltv.BuildCache(t, channel)
		}

		t = t.Add(time.Hour * 24)
	}
	log.Println("Refreshed epg cache")
}
