package tv2

import (
	"time"
)

type Response struct {
	Date     string    `json:"date"`
	Channels []Channel `json:"channel"`
}

type Channel struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"shortName"`
	Category  string    `json:"category"`
	Score     int       `json:"score"`
	Tv2Id     int       `json:"tv2id"`
	Programs  []Program `json:"program"`
}

type Program struct {
	Id              int       `json:"id"`
	SeriesId        string    `json:"srsid"`
	ProgramId       string    `json:"programId"`
	House           string    `json:"house"`
	Title           string    `json:"title"`
	Start           time.Time `json:"start"`
	Stop            time.Time `json:"stop"`
	Episode         int       `json:"epnr"`
	EpisodeTotal    int       `json:"eptot"`
	Season          int       `json:"season"`
	SeriesSynopsis  string    `json:"srsyn"`
	EpisodeSynopsis string    `json:"epsyn"`

	Nationality    string `json:"natio"`
	Genre          string `json:"genre"`
	ProductionYear int    `json:"pyear"`

	IsReplay bool `json:"isrepl"`
	IsLive   bool `json:"islive"`
	IsMovie  bool `json:"ismov"`
	IsSyn    bool `json:"issyn"`
}
