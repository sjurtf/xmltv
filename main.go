package main

import (
	"os"
	"xmltv-exporter/cmd"
	"xmltv-exporter/xmltv"
)

func main() {
	xmltv.Init(os.Getenv("XMLTV_DOMAIN"))

	cmd.MapEgp()
	cmd.ServeEpg(os.Getenv("XMLTV_PORT"))

}
