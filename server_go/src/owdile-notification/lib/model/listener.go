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
func GetListenerInstance(dbc influxbase.DBConfig) *ListenerMongo{
    fmt.Printf("dbc     %v\n",dbc)
    if hostListenerMap == nil{
        hostListenerMap = make(map[string] *ListenerMongo)
    }
    if hostListenerMap[dbc.Host] == nil {
        var mongo  ListenerMongo
        if !mongo.Init(dbc) {
            fmt.Printf("mongo     null\n")
            return nil
        }
        hostListenerMap[dbc.Host] = &mongo
    }
    return hostListenerMap[dbc.Host]
}


var hostListenerMap map[string] *ListenerMongo

type ListenerMongo struct{
    coll *mongodb.MongoDBCollection
    dbname string
    collName string
}

func (p* ListenerMongo) Init(dbc influxbase.DBConfig) bool{
    p.dbname = "test"
    p.collName = "listener_info"
    mongo := mongodb.GetMongoDBInstance(dbc)
    var err error
    p.coll,err = mongo.GetCollection(p.dbname,p.collName)

    return !influxbase.HasError(err)
}

//获取所有通知列表
func (p* ListenerMongo)GetListenersFromDB() ([]base.ListenerInfoDb,bool){
    var result []base.ListenerInfoDb

    err := p.coll.Find(nil).Sort("_id").All(&result)
    if influxbase.HasError(err) {
        return nil,false
    }
    return result,true
}

func (p* ListenerMongo)SaveNotifyInfosToDB(infoDb base.ListenerInfoDb) bool{
    fmt.Printf("\n%v\n",infoDb)
    m := bson.M{"listenerinfo.name":infoDb.Name}//event为主键
    _,err := p.coll.Upsert(m,infoDb)
    //err := p.coll.Insert(notiInfoDb)
    //fmt.Printf("%v",info)
    return !influxbase.HasError(err)
}
