package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/data"
	"slices"
	"strconv"
	"strings"
	"time"
)

func handleMetaQueries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "Method not allowed")
		return
	}
	role := r.PathValue("role")
	allStr := r.URL.Query().Get("all")
	allChamps := false
	if allStr != "" {
		if parsed, err := strconv.ParseBool(allStr); err == nil {
			allChamps = parsed
		}
	}
	var top int
	topStr := r.URL.Query().Get("top")
	if topStr != "" {
		if parsed, err := strconv.Atoi(topStr); err == nil {
			if parsed > 0 {
				top = parsed
			}
		}
	}
	rank := r.URL.Query().Get("rank")
	if rank == "" {
		rank = "overall"
	}
	rank = strings.ToLower(rank)

	resultChamps := GetChampStatData(role, top, allChamps, rank)

	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resultChamps)
}
func handleCounterQueries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, "Method not allowed")
		return
	}
	champ := r.PathValue("champ")
	role := r.URL.Query().Get("role")
	allStr := r.URL.Query().Get("all")
	var all bool
	var err error
	if allStr != "" {
		all, err = strconv.ParseBool(allStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error with boolean argument: ", err)
			return
		}
	}
	rank := r.URL.Query().Get("rank")
	if rank == "" {
		rank = "overall"
	}
	rank = strings.ToLower(rank)
	res := findChampCounters(champ, role, rank, all)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return

}

func findChampCounters(champion, queriedRole, rank string, all bool) map[string][]data.ChampCounter {
	fmt.Printf("egesz: %+v\n", champdata[rank])
	result := make(map[string][]data.ChampCounter)
	for role, champs := range champdata[rank].Role {
		for _, champ := range champs {
			if strings.ToUpper(champ.Name) == strings.ToUpper(champion) {
				for _, counter := range champ.Counters {
					if !all && counter.MatchCount > 10 {
						result[role] = append(result[role], counter)
					} else if all {
						result[role] = append(result[role], counter)
					}

				}

			}
		}

	}
	if queriedRole != "" {
		res := make(map[string][]data.ChampCounter)
		res[queriedRole] = result[queriedRole]
		return res
	}
	return result
}

func GetChampStatData(role string, top int, allChamps bool, rank string) []data.Champ {
	resultChamps := champdata[rank].Role[role]
	if !allChamps {
		resultChamps = slices.Collect(func(yield func(data.Champ) bool) {
			for _, ch := range resultChamps {
				if ch.Pickrate > 0.5 {
					if !yield(ch) {
						return
					}
				}
			}
		})
	}
	if top <= 0 {
		top = len(resultChamps)
	}
	return resultChamps[:top]
}

var champdata map[string]data.RoleMap

func handleUpdater(ticker *time.Ticker, done chan bool) {

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Megáll a gorutin")
				return
			case <-ticker.C:
				data.UpdateChampStatData(champdata)
			}

		}
	}()

}

func getEnv(envName, def string) string {
	res := os.Getenv(envName)
	if res == "" {
		return def
	}
	return res
}

func main() {
	updaterFreq, err := time.ParseDuration(getEnv("UPDATER_FREQ", "12h"))
	if err != nil {
		fmt.Println("Couldnt parse updater frequency. Setting to 12h")
		updaterFreq = 12 * time.Hour
	}
	champdata, err = data.GetAllRankStatData()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("GET /api/{role}/meta/", handleMetaQueries)
	http.HandleFunc("GET /api/{champ}/counter/", handleCounterQueries)
	ticker := time.NewTicker(updaterFreq)
	done := make(chan bool)
	go handleUpdater(ticker, done)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
