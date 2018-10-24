package models

import (
	"gopkg.in/mgo.v2/bson"
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

func getMongoCollection() (*mongodb.MongoDBCollection, error) {
	mongo := mongodb.GetMongoDBInstance(influxbase.DBConfig{Host: "192.168.10.10:27017"})

	collection, err := mongo.GetCollection("noah", "apps")

	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (apps *Apps) GetList(page, category string) ([]bson.M, error) {
	collection, err := getMongoCollection()

	if err == nil {
		var result []bson.M
		var query bson.M = nil

		if category != "0" {
			query = bson.M{"category_id": category}
		}

		err := collection.Find(query).All(&result)
		if err == nil {
			return result, err
		}
	}

	return nil, err
}

func (apps *Apps) Add(title, background, description, tag string, buttons []AppButton) bool {
	collection, err := getMongoCollection()

	if err == nil {
		err := collection.Insert(Apps{
			tag:         tag,
			title:       title,
			background:  background,
			description: description,
			buttons:     buttons,
		})
		if err == nil {
			return true
		}
	}

	return false
}
