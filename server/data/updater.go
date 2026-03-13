package data

import (
	"log"
)

func UpdateChampStatData(rm map[string]RoleMap) {

	rolemap, err := GetAllRankStatData()
	log.Printf("updating attempt\n")
	if err != nil {
		log.Println("Updater error: ", err)
	}
	if rolemap["bronze"].LastUpdated.After(rm["bronze"].LastUpdated) {
		log.Printf("Updating rolemap. Current version update: %v, fetched version update: %v\n", rm["bronze"].LastUpdated, rolemap["bronze"].LastUpdated)
		rm = rolemap
	}

}
