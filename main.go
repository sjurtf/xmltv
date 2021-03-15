package main

import (
	"github.com/go-co-op/gocron"
	"os"
	"time"
	"xmltv-exporter/cmd"
	"xmltv-exporter/xmltv"
)

func main() {
	xmltv.Init(os.Getenv("XMLTV_DOMAIN"))

	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Every("1d").Do(cmd.RefreshToday)
	_, _ = s.Every("3d").Do(cmd.RefreshWeek)
	s.StartAsync()
	cmd.RefreshToday()
	cmd.RefreshWeek()

	cmd.ServeEpg(os.Getenv("XMLTV_PORT"))
}
