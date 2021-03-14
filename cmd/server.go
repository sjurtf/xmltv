package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
	"xmltv-exporter/xmltv"
)

func ServeEpg(port string) {

	http.HandleFunc("/channels-norway.xml", ChannelListHandler)
	http.HandleFunc("/", ChannelHandler)

	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func ChannelListHandler(w http.ResponseWriter, r *http.Request) {
	if !isValid(w, r) {
		return
	}

	bytes, err := xmltv.GetChannelList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(bytes)
}

func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	if !isValid(w, r) {
		return
	}

	channelId, err := channelId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := date(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := xmltv.GetSchedule(channelId, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_, _ = w.Write(bytes)
}

/**
Finds channelId from Path
Example input: zebrahd.tv2.no_2021-03-16.xml.gz
*/
func channelId(r *http.Request) (string, error) {
	filename := path.Base(r.URL.Path)
	parts := strings.Split(filename, "_")

	if len(parts) == 0 {
		return "", fmt.Errorf("channel name could not be found from input '%s'", filename)
	}
	return parts[0], nil
}

/**
Finds time.Time from Path
Example input: zebrahd.tv2.no_2021-03-16.xml.gz
*/
func date(r *http.Request) (time.Time, error) {
	filename := path.Base(r.URL.Path)
	parts := strings.Split(filename, "_")

	if len(parts) <= 1 {
		return time.Time{}, fmt.Errorf("date could not be found from input '%s'", filename)
	}

	parts2 := strings.Split(parts[1], ".")

	date, err := time.Parse("2006-01-02", parts2[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("date could not be parsed: '%s'", err)

	}

	return date, nil
}

func isValid(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return false
	}

	if len(r.URL.Path) <= 1 {
		http.Error(w, "Not found", http.StatusNotFound)
		return false
	}

	return true
}
