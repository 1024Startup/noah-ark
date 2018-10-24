package model

import (
    _"container/list"
    "influx.io/influxbase"
    "influx.io/mongodb"
    "owdile-notification/lib/base"
    "gopkg.in/mgo.v2/bson"
    "fmt"
)

//通过配置获取实体连接
func GetNotifyInstance(dbc influxbase.DBConfig) *NotifyMongo{
    fmt.Printf("dbc     %v\n",dbc)
    if hostNotifyMap == nil{
        hostNotifyMap = make(map[string] *NotifyMongo)
    }
    if hostNotifyMap[dbc.Host] == nil {
        var notifymongo  NotifyMongo
        if !notifymongo.Init(dbc) {
            fmt.Printf("notifymongo     null\n")
            return nil
        }
        hostNotifyMap[dbc.Host] = &notifymongo
    }
    return hostNotifyMap[dbc.Host]
}


var hostNotifyMap map[string] *NotifyMongo

type NotifyMongo struct{
    coll *mongodb.MongoDBCollection
    dbname string
    collName string
}

func (p* NotifyMongo) Init(dbc influxbase.DBConfig) bool{
    p.dbname = "test"
    p.collName = "notify_info"
    mongo := mongodb.GetMongoDBInstance(dbc)
    var err error
    p.coll,err = mongo.GetCollection(p.dbname,p.collName)

    return !influxbase.HasError(err)
}

//获取所有通知列表
func (p* NotifyMongo)GetNotifyInfosFromDB() ([]base.NotifyInfoDb,bool){
    var result []base.NotifyInfoDb

    err := p.coll.Find(nil).Sort("_id").All(&result)
    if influxbase.HasError(err) {
        return nil,false
    }
    return result,true
}

//
func (p* NotifyMongo)SaveNotifyInfosToDB(notiInfoDb base.NotifyInfoDb) bool{
    fmt.Printf("\n%v\n",notiInfoDb)
    m := bson.M{"notifyinfo.event":notiInfoDb.Event}//event为主键
    _,err := p.coll.Upsert(m,notiInfoDb)
    //err := p.coll.Insert(notiInfoDb)
    //fmt.Printf("%v",info)
    return !influxbase.HasError(err)
}
