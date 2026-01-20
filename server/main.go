package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/data"
	"slices"
	"strconv"
	"strings"
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
	resultChamps := GetChampStatData(role, top, allChamps)

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

	res := findChampCounters(champ, role, all)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	return

}

func findChampCounters(champion, queriedRole string, all bool) map[string][]data.ChampCounter {
	result := make(map[string][]data.ChampCounter)
	for role, champs := range champdata.Role {
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

func GetChampStatData(role string, top int, allChamps bool) []data.Champ {
	resultChamps := champdata.Role[role]
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

var champdata data.RoleMap

func main() {
	var err error
	champdata, err = data.GetChampStatData()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("GET /api/{role}/meta/", handleMetaQueries)
	http.HandleFunc("GET /api/{champ}/counter/", handleCounterQueries)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
