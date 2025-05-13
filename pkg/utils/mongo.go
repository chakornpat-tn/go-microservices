package utils

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ConvToObjID(id string) bson.ObjectID {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error: ConvToObjID: %s", err.Error())
	}
	return objectID
}

func LocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func ConvertStringTimeToTime(t string) time.Time {
	layout := "2006-01-02 15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	result = result.In(loc)
	return result

}
