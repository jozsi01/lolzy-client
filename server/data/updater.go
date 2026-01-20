package data

import (
	"fmt"
)

func UpdateChampStatData(rm *RoleMap) {
	rolemap, err := GetChampStatData()
	if err != nil {
		fmt.Println("Updater error: ", err)
	}

	if rolemap.LastUpdated.After(rm.LastUpdated) {
		fmt.Printf("Updating rolemap. Current version update: %v, fetched version update: %v\n", rm.LastUpdated, rolemap.LastUpdated)
		rm = &rolemap
	}

}
