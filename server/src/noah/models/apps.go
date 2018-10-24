package models

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"influx.io/influxbase"
	"influx.io/mongodb"
)

type Apps struct {
}
type AppButton struct {
	Name     string `json:"name"`
	TypeName string `json:"type"`
	Url      string `json:"url"`
}

func (app Apps) GetList(page, category int) (string, error) {
	mongo := mongodb.GetMongoDBInstance(influxbase.DBConfig{Host: "192.168.10.10:27017"})

	collection, err := mongo.GetCollection("noah", "apps")

	if err != nil {
		var result []AppButton
		err := collection.Find(bson.M{}).All(&result)
		if err != nil {
			return "", err
		}
		if data, err := json.Marshal(result); err != nil {
			return string(data[:]), nil
		}
	}

	return "", err
}
