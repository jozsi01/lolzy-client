package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var hostname string = os.Getenv("MONGODB_HOST")

var coll *mongo.Collection

func ConnectToDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if hostname == "" {
		hostname = "localhost"
	}
	var connectionString string = fmt.Sprintf("mongodb://%s:27017", hostname)
	db, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	coll = db.Database("lolzy").Collection("champ_data")

	if err := db.Ping(ctx, nil); err != nil {
		log.Fatal("Could not ping MongoDB:", err)
		return err
	}

	fmt.Println("Connected to MongoDB!")
	return nil

}
func GetChampData() (ChampDataResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	champdata := GetChampionData()
	_, err := coll.UpdateOne(ctx, bson.M{"type": "champion"}, bson.M{"$set": champdata}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		fmt.Println(err)
		return ChampDataResp{}, err
	}
	return champdata, nil
}

func GetChampStat() (RoleMap, error) {
	chresp, err := GetChampData()
	if err != nil {
		fmt.Println("Error getting champData: ", err)
		return RoleMap{}, err
	}
	roles, lastUpdated, err := GetChampStatData(chresp)
	if err != nil {
		fmt.Println("Error getting champStat: ", err)
		return RoleMap{}, err
	}
	for _, champs := range roles {
		sort.Slice(champs, func(i, j int) bool {
			return champs[i].Winrate > champs[j].Winrate
		})
	}
	return RoleMap{
		LastUpdated: lastUpdated,
		Role:        roles,
	}, nil

}
