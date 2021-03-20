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
	DocHeader           = "<?xml version=\"1.0\" encoding=\"utf-8\"?><!DOCTYPE tv SYSTEM \"xmltv.dtd\">\n"
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

		name := CommonElement{
			Lang:  "en",
			Value: c.Name,
		}

		channel := Channel{
			Id:      xmlChannelIdMap[c.Id],
			Name:    name,
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
		cat := "series"
		if p.IsMovie {
			cat = "movie"
		}

		category := CommonElement{
			Lang:  "en",
			Value: cat,
		}

		titles := CommonElement{
			Lang:  "nb",
			Value: p.Title,
		}

		var subtitles []CommonElement
		var episodeNums []EpisodeNum
		if p.Season != 0 || p.Episode != 0 {
			episode := EpisodeNum{
				System:     XmltvEpisodeStd,
				EpisodeNum: formatEpisode(p.Season, p.Episode, p.EpisodeTotal),
			}
			if !p.IsMovie {
				episodeNums = []EpisodeNum{episode}
				subtitles = []CommonElement{{
					Lang:  "nb",
					Value: formatEpisodeHuman(p.Season, p.Episode),
				},
				}
			}

		} else {
			episodeNums = nil
			subtitles = nil
		}

		var desc string
		if p.EpisodeSynopsis == "" && p.SeriesSynopsis == "" {
			desc = "Beskrivelse ikke tilgjengelig"
		} else if p.EpisodeSynopsis != "" {
			desc = p.EpisodeSynopsis
		} else {
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
			Categories:   []CommonElement{category},
			Descriptions: []CommonElement{description},
			SubTitles:    subtitles,
			EpisodeNums:  episodeNums,
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

	bytes, err := xml.MarshalIndent(resp, "", "  ")
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
	return fmt.Sprintf("%d.%d.0/1", s-1, e-1)
}

func formatEpisodeHuman(s, e int) string {
	return fmt.Sprintf("S%dE%d", s, e)
}

func formatTime(date time.Time) string {
	return date.Format(XmltvDateFormat)
}

func formatCacheKey(date time.Time) string {
	y, m, d := date.Date()
	return fmt.Sprintf("%d-%s-%d", y, m, d)
}
