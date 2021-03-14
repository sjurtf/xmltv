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
		channels := tv2.FetchEpg(today)
		xmltv.UpdateAvailableChannels(channels)

		// TODO make this more efficient
		for _, channel := range channels {
			xmltv.BuildCache(t, channel)
		}

		t = t.Add(time.Hour * 24)
		log.Printf("date: %s", t)
	}
	log.Println("Epg cache refreshed")
}
