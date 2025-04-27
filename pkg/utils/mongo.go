package utils

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvToObj(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}

func LocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func ConvertStringTimeToTime(t string) time.Time {
	layout := "2006-01-02T15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	result = result.In(loc)
	return result

}
