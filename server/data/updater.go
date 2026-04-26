package data

import (
	"log"
)

func UpdateChampStatData(rm map[string]RoleMap) (map[string]RoleMap, error) {
	rolemap, err := GetAllRankStatData()
	log.Printf("updating attempt\n")
	if err != nil {
		log.Println("Updater error: ", err)
		return rm, err
	}
	if rolemap["bronze"].LastUpdated.After(rm["bronze"].LastUpdated) {
		log.Printf("Updating rolemap. Current version update: %v, fetched version update: %v\n", rm["bronze"].LastUpdated, rolemap["bronze"].LastUpdated)
		return rolemap, nil
	}

	return rm, nil
}
