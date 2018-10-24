package models

import (
	"fmt"
	"influx.io/influxbase"
	"influx.io/mongodb"
)

type Apps struct {
	id          string
	tag         string
	title       string
	background  string
	description string
	buttons     []AppButton
}
type AppButton struct {
	Name     string `json:"name"`
	TypeName string `json:"type"`
	Url      string `json:"url"`
}

func (apps *Apps) GetList(page, category int) ([]Apps, error) {
	mongo := mongodb.GetMongoDBInstance(influxbase.DBConfig{Host: "192.168.10.10:27017"})

	collection, err := mongo.GetCollection("noah", "apps")

	if err == nil {
		var result []Apps
		err := collection.Find(nil).All(&result)
		if err == nil {
			for _, v := range result {
				fmt.Print(v)
			}

			return result, err
		}
	}

	return nil, err
}

func (apps *Apps) Add(title, background, description, tag string, buttons []AppButton) {

}
