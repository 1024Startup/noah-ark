package model

import (
    _"container/list"
    "influx.io/influxbase"
    "influx.io/mongodb"
    "owdile-notification/lib/base"
    "gopkg.in/mgo.v2/bson"
    _"fmt"
    "log"
    "time"
)



type NotifyEventMongo struct{
    coll *mongodb.MongoDBCollection
    dbname string
    collName string
    dbc influxbase.DBConfig
    mdb * mongodb.MongoDB
}

//通过配置获取实体连接
func GetEventMongoInstance(dbc influxbase.DBConfig) *NotifyEventMongo{
    //log.Printf("dbc     %v\n",dbc)
    if hostMongoMap == nil{
        hostMongoMap = make(map[string] *NotifyEventMongo)
    }
    if hostMongoMap[dbc.Host] == nil {
        var mongo  NotifyEventMongo
        if !mongo.Init(dbc) {
            log.Printf("mongo     null\n")
            return nil
        }
        hostMongoMap[dbc.Host] = &mongo
    }
    return hostMongoMap[dbc.Host]
}



var hostMongoMap map[string] *NotifyEventMongo

func (p* NotifyEventMongo) Init(dbc influxbase.DBConfig) bool{
    p.dbname = "test"
    p.collName = "notify_logs"
    p.dbc = dbc

    p.mdb = mongodb.GetMongoDBInstance(dbc)
    var err error
    p.coll,err = p.mdb.GetCollection(p.dbname,p.collName)

    if influxbase.HasError(err){

        return false
    }
    return true
}

//获取所有通知列表
func (p* NotifyEventMongo)GetUnfinishNotifyEventsFromDB() ([]base.NotifyDetailDb,bool){
    var result []base.NotifyDetailDb

    curtime  :=  time.Now().Unix()
    time := curtime - 3600*24*7//7天前
    conditions := bson.M{"createdtime":bson.M{"$gt":time},
    "lastexcutetime":bson.M{"$lt":curtime},//下次执行时间 小于当前时间
    "excutesstatus":bson.M{"$in":[]base.EXCUTE_STATUS{base.NEW_STATUS,base.ERROR_STATUS}}}

    err := p.coll.Find(conditions).Sort("_id").All(&result)
    if influxbase.HasError(err) {
        return nil,false
    }
    return result,true
}

func (p* NotifyEventMongo)SaveNotifyEventToDB(infoDb base.NotifyDetailDb) bool{
    log.Println("%v",infoDb)

    m := bson.M{"_id":infoDb.Id}
    _,err := p.coll.Upsert(m,infoDb)
    //fmt.Printf("%v",info)
    if influxbase.HasError(err){
        p.reconnectDB()
        return false
    }
    return true
}
func (p* NotifyEventMongo)reconnectDB() bool{

    p.mdb.ReconnectDB()
    p.Init(p.dbc)
    log.Println("reconnect %v",p.dbc)

    return true
}
/*func (p* NotifyEventMongo)SaveNotifyEventsToDB(infoDbs []base.NotifyDetailDb) bool{
    fmt.Printf("\n%v\n",infoDbs)
    //m := bson.M{"listenerinfo.name":infoDb.Name}//event为主键
    //_,err := p.coll.Upsert(m,infoDb)
    err := p.coll.Insert(infoDbs)
    //fmt.Printf("%v",info)
    return !influxbase.HasError(err)
}*/
