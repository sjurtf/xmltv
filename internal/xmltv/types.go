package xmltv

import (
	"encoding/xml"
)

type Response struct {
	XMLName       xml.Name    `xml:"tv"`
	GeneratorName string      `xml:"generator-info-name,attr"`
	GeneratorUrl  string      `xml:"generator-info-url,attr"`
	ChannelList   []Channel   `xml:"channel,omitempty"`
	ProgrammeList []Programme `xml:"programme,omitempty"`
}

type Channel struct {
	XMLName xml.Name      `xml:"channel"`
	Id      string        `xml:"id,attr"`
	Name    CommonElement `xml:"display-name"`
	BaseUrl string        `xml:"base-url"`
}

type Programme struct {
	XMLName      xml.Name        `xml:"programme"`
	Channel      string          `xml:"channel,attr"`
	Start        string          `xml:"start,attr"`
	Stop         string          `xml:"stop,attr"`
	Titles       []CommonElement `xml:"title"`
	SubTitles    []CommonElement `xml:"sub-title,omitempty"`
	Descriptions []CommonElement `xml:"desc,omitempty"`
	Categories   []CommonElement `xml:"category,omitempty"`
	EpisodeNums  []EpisodeNum    `xml:"episode-num,omitempty"`
	Credits      string          `xml:"credits,omitempty"`
	Date         string          `xml:"date,omitempty"`
	Ratings      []Rating        `xml:"rating,omitempty"`
}

type CommonElement struct {
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

type EpisodeNum struct {
	System     string `xml:"system,attr"`
	EpisodeNum string `xml:",chardata"`
}

type Rating struct {
	Value  string `xml:"value"`
	System string `xml:"system,attr,omitempty"`
}
