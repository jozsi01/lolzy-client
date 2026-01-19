package data

import "time"

type ChampData struct {
	Key  int    `json:"key,string"`
	Name string `json:"name"`
}

type ChampDataResp struct {
	Version string
	Data    map[string]ChampData `json:"data"`
}

type Champ struct {
	Name       string
	WinCount   float64
	MatchCount float64
	Counters   []ChampCounter
	AllMatches float64 // all the matches analyzed by u.gg for pick rate calculation
	Winrate    float64
	Pickrate   float64
}

type ChampCounter struct {
	Name             string
	LostMatchCounter float64
	MatchCount       float64
	LoseRate         float64
}
type RoleMap struct {
	LastUpdated time.Time
	Role        map[string][]Champ
}
