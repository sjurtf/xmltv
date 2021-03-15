package xmltv

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
	"xmltv-exporter/tv2"
)

const (
	XmltvDateFormat     = "20060102150400 -0700"
	XmltvEpisodeStd     = "xmtv_ns"
	GeneratorName       = "xmltv.sjurtf.net"
	DocHeader           = `<?xml version="1.0" encoding="utf-8"?><!DOCTYPE tv SYSTEM "xmltv.dtd">`
	defaultGeneratorUrl = "https://xmltv.sjurtf.net/"
)

var channelCache map[int]tv2.Channel
var channelGuideMap map[string]map[string][]tv2.Program
var generatorUrl string

func Init(url string) {
	generatorUrl = defaultGeneratorUrl
	if url != "" {
		generatorUrl = url
	}

	channelCache = make(map[int]tv2.Channel)
	channelGuideMap = make(map[string]map[string][]tv2.Program)
}

func BuildCache(date time.Time, channel tv2.Channel) {
	if channelGuideMap == nil {
		channelGuideMap = make(map[string]map[string][]tv2.Program)
	}

	dateKey := formatCacheKey(date)
	if channelGuideMap[dateKey] == nil {
		channelGuideMap[dateKey] = make(map[string][]tv2.Program)
	}

	xmlChannelId := xmlChannelIdMap[channel.Id]
	if xmlChannelId == "" {
		log.Printf("channel %s with id %d is not mapped", channel.Name, channel.Id)
	}
	channelGuideMap[dateKey][xmlChannelId] = channel.Programs
	channelCache[channel.Id] = channel
}

func GetChannelList() ([]byte, error) {
	if len(channelCache) == 0 {
		return nil, fmt.Errorf("channeldata unavailable")
	}

	var channels []Channel
	var programs []Programme
	for _, c := range channelCache {
		channel := Channel{
			Id:      xmlChannelIdMap[c.Id],
			Name:    c.Name,
			BaseUrl: generatorUrl,
		}
		channels = append(channels, channel)
	}

	log.Printf("fetched available channels. Num channels %d", len(channels))
	return marshall(channels, programs)
}

func GetSchedule(channelId string, date time.Time) ([]byte, error) {
	return marshall(nil, getProgramsForChannel(channelId, date))
}

func getProgramsForChannel(channelId string, date time.Time) []Programme {
	dateKey := formatCacheKey(date)
	guide := channelGuideMap[dateKey][channelId]

	var programs []Programme
	for _, p := range guide {

		titles := CommonElement{
			Lang:  "nb",
			Value: p.Title,
		}

		subtitle := CommonElement{
			Lang:  "nb",
			Value: p.Title,
		}

		episode := EpisodeNum{
			System:     XmltvEpisodeStd,
			EpisodeNum: formatEpisode(p.Season, p.Episode, p.EpisodeTotal),
		}

		desc := p.EpisodeSynopsis
		if desc == "" {
			desc = p.SeriesSynopsis
		}

		description := CommonElement{
			Lang:  "nb",
			Value: desc,
		}

		programme := Programme{
			Channel:      channelId,
			Start:        formatTime(p.Start),
			Stop:         formatTime(p.Stop),
			Titles:       []CommonElement{titles},
			SubTitles:    []CommonElement{subtitle},
			Descriptions: []CommonElement{description},
			EpisodeNums:  []EpisodeNum{episode},
		}
		programs = append(programs, programme)
	}

	log.Printf("fetched guide: %s - channelId %s. Num programs %d", dateKey, channelId, len(guide))
	return programs
}

func marshall(channels []Channel, programs []Programme) ([]byte, error) {
	resp := Response{
		GeneratorName: GeneratorName,
		GeneratorUrl:  generatorUrl,
		ChannelList:   channels,
		ProgrammeList: programs,
	}

	bytes, err := xml.Marshal(resp)
	if err != nil {
		log.Fatalln("unable to marshal")
	}

	return append([]byte(DocHeader), bytes...), nil
}

/*
<episode-num system="xmltv_ns">s.e.p/t</episode-num>
Where s is the season number minus 1.
Where e is the episode number minus 1.
Where p is the part number minus 1.
Where t to the total parts (do not subtract 1)

so Season 7, Episode 5, Part 1 of 2 would appear as:
<episode-num system="xmltv_ns">6.4.0/2</episode-num>
*/
func formatEpisode(s, e, t int) string {
	return fmt.Sprintf("%d.%d/%d", s, e, t)
}

func formatTime(date time.Time) string {
	return date.Format(XmltvDateFormat)
}

func formatCacheKey(date time.Time) string {
	y, m, d := date.Date()
	return fmt.Sprintf("%d-%s-%d", y, m, d)
}
