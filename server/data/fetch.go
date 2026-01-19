package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var championDataUrl string = "https://static.bigbrain.gg/assets/lol/riot_static/16.1.1/data/en_US/champion.json"
var campStatDataURL string = "https://stats2.u.gg/lol/1.5/champion_ranking/world/16_1/ranked_solo_5x5/emerald_plus/1.5.0.json"

func GetChampionData() ChampDataResp {
	resp, err := http.Get(championDataUrl)
	if err != nil {
		fmt.Println("Error with the champion meta data request: ", err)
		return ChampDataResp{}
	}
	defer resp.Body.Close()
	var respData ChampDataResp
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		fmt.Println("Error with the reading of the response body: ", err)
		return ChampDataResp{}
	}
	return respData
}

func GetChampStatData(champMetaData ChampDataResp) (map[string][]Champ, time.Time, error) {
	resp, err := http.Get(campStatDataURL)
	if err != nil {
		fmt.Println("Error with the champion meta data request: ", err)
		return nil, time.Time{}, err
	}
	defer resp.Body.Close()
	bytedata, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, time.Time{}, err
	}
	var raw []json.RawMessage
	err = json.Unmarshal(bytedata, &raw)
	if err != nil {
		fmt.Println(err)
		return nil, time.Time{}, err
	}
	return convertChampStatDataIntoRoleChampMap(raw, champMetaData)
}

func getChampNameFromId(id int, champMetaData ChampDataResp) string {
	for _, v := range champMetaData.Data {
		if id == v.Key {
			return v.Name
		}
	}
	return ""
}

func convertChampStatDataIntoRoleChampMap(raw []json.RawMessage, champMetaData ChampDataResp) (map[string][]Champ, time.Time, error) {
	top := raw[0]
	var lastUpdated time.Time
	err := json.Unmarshal(raw[2], &lastUpdated)
	if err != nil {
		fmt.Println("Baj volta a lastUpdated parsnál: ", err)
		return nil, time.Time{}, err
	}
	var totalMatches float64
	err = json.Unmarshal(raw[3], &totalMatches)
	if err != nil {
		fmt.Println("Baj volta a totalMatches parsnál: ", err)
		return nil, time.Time{}, err
	}
	fmt.Println("Time: ", lastUpdated)
	fmt.Printf("all match: %f\n", totalMatches)
	var rawstruct map[string][][]interface{}
	res := make(map[string][]Champ)
	if err := json.Unmarshal(top, &rawstruct); err != nil {
		fmt.Println("Baj volt a őrasenál: ", err)
		return nil, time.Time{}, err
	}
	for role, champs := range rawstruct {
		fmt.Printf("Role: %s\n", role)
		for _, champDatas := range champs {
			var champ Champ
			champID := champDatas[0].(string)
			champid, err := strconv.Atoi(champID)
			if err != nil {
				return nil, time.Time{}, err
			}
			champName := getChampNameFromId(champid, champMetaData)
			champ.Name = champName
			counters := champDatas[1].([]interface{})
			winCount := champDatas[2].(float64)
			matchesCount := champDatas[3].(float64)
			champ.WinCount = winCount
			champ.MatchCount = matchesCount
			champ.AllMatches = totalMatches
			champ.Winrate = (champ.WinCount / champ.MatchCount) * 100
			champ.Pickrate = (champ.MatchCount / champ.AllMatches) * 100
			var chCounters []ChampCounter

			for idx, val := range counters {
				v := val.([]interface{})
				var chCounter ChampCounter
				if idx >= 3 {
					break
				}
				coutnerName := getChampNameFromId(int(v[0].(float64)), champMetaData)
				chCounter.Name = coutnerName
				chCounter.LostMatchCounter = v[1].(float64)
				chCounter.MatchCount = v[2].(float64)
				chCounter.LoseRate = (chCounter.LostMatchCounter / chCounter.MatchCount) * 100
				chCounters = append(chCounters, chCounter)
			}
			champ.Counters = chCounters
			res[role] = append(res[role], champ)
		}
	}
	return res, lastUpdated, nil
}
