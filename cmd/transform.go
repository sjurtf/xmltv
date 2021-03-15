package cmd

import (
	"log"
	"time"
	"xmltv-exporter/tv2"
	"xmltv-exporter/xmltv"
)

func RefreshWeek() {
	days := 7
	tomorrow := time.Now().Add(time.Hour * 24)
	refresh(tomorrow, days)
	log.Printf("Refreshed epg cache for next %d days", days)
}

func RefreshToday() {
	refresh(time.Now(), 1)
	log.Println("Refreshed epg cache for today")
}

func refresh(start time.Time, days int) {
	var t = start
	end := time.Now().Add(time.Hour * time.Duration(24*days))

	for t.Before(end) {
		channels := tv2.FetchEpg(t)
		for _, channel := range channels {
			xmltv.BuildCache(t, channel)
		}
		t = t.Add(time.Hour * 24)
	}
}
