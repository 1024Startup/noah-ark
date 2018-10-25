package models

import (
	"gopkg.in/mgo.v2/bson"
	"influx.io/influxbase"
	"influx.io/mongodb"
)

type Apps struct {
	//id          string
}
type AppsData struct {
	//id          string
	Tag         string `form:"tag" json:"tag" bson:"tag"`
	Title       string `form:"title" json:"title"`
	Background  string `form:"background" json:"background"`
	Category_id  string `form:"category" json:"category_id"`
	Type 		string `form:"type" json:"type"`
	Description string `form:"description" json:"description"`
	Buttons     []AppButton `form:"buttons" json:"buttons"`
}
type AppButton struct {
	Name     string `form:"name" json:"name"`
	Type string `form:"type" json:"type"`
	Url      string `form:"url" json:"url"`
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

func (apps *Apps) Add(data AppsData) bool {
	collection, err := getMongoCollection()

	if err == nil {
		//err := collection.Insert()
		m := bson.M{"title":data.Title}//event为主键
	    _,err := collection.Upsert(m,data)

		if influxbase.HasError(err){

	        return false
	    }
	}

	return false
}
