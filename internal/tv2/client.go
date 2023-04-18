package tv2

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// https://rest.tv2.no/epg-dw-rest/epg/program/2021/03/15/
const (
	rootUrl   = "https://rest.tv2.no/epg-dw-rest/epg/program/"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:86.0) Gecko/20100101 Firefox/86.0"
)

func FetchEpg(date time.Time) []Channel {
	day := date.Day()
	month := date.Month()
	year := date.Year()

	datePath := fmt.Sprintf("%d/%d/%d", year, int(month), day)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, rootUrl+datePath, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("unable to get data: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("got unhappy statuscode: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	r := &Response{}
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Fatalf("unable to decode data: %s", err)
	}

	return r.Channels
}
