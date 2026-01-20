package data

import (
	"fmt"
)

func UpdateChampStatData(rm map[string]RoleMap) {
	rolemap, err := GetAllRankStatData()
	if err != nil {
		fmt.Println("Updater error: ", err)
	}

	if rolemap["bronze"].LastUpdated.After(rm["bronze"].LastUpdated) {
		fmt.Printf("Updating rolemap. Current version update: %v, fetched version update: %v\n", rm["bronze"].LastUpdated, rolemap["bronze"].LastUpdated)
		rm = rolemap
	}

}
