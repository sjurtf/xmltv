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
	XMLName xml.Name `xml:"channel"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"display-name"`
	BaseUrl string   `xml:"base-url"`
}

type Programme struct {
	XMLName     xml.Name   `xml:"programme"`
	Channel     string     `xml:"channel,attr"`
	Start       string     `xml:"start,attr"`
	Stop        string     `xml:"stop,attr"`
	Title       string     `xml:"title"`
	SubTitle    string     `xml:"sub-title,omitempty"`
	Description string     `xml:"desc,omitempty"`
	EpisodeNum  EpisodeNum `xml:"episode-num,omitempty"`
	Credits     string     `xml:"credits,omitempty"`
	Date        string     `xml:"date,omitempty"`
	Categories  []Category `xml:"category,omitempty"`
	Rating      []Rating   `xml:"rating,omitempty"`
}

type EpisodeNum struct {
	XMLName    xml.Name `xml:"episode-num"`
	System     string   `xml:"system,attr"`
	EpisodeNum string   `xml:",chardata"`
}

type Rating struct {
	System string `xml:"system,attr"`
	Value  string `xml:"value"`
}

type Category struct {
	Value string `xml:",chardata"`
	Lang  string `xml:"lang,attr"`
}
