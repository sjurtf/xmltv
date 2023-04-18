package main

import (
	"os"
	"time"

	"xmltv/internal/xmltv"

	"github.com/go-co-op/gocron"

	"xmltv/cmd"
)

func main() {
	xmltv.Init(os.Getenv("XMLTV_DOMAIN"))

	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Every("1d").Do(cmd.RefreshToday)
	_, _ = s.Every("3d").Do(cmd.RefreshWeek)
	s.StartAsync()
	cmd.RefreshToday()
	cmd.RefreshWeek()

	server := cmd.Server{}
	server.Start()
}
