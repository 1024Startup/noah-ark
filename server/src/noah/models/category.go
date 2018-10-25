package models

import
(
    "influx.io/influxbase"
    "influx.io/mongodb"
    "gopkg.in/mgo.v2/bson"
    "log"
    "encoding/json"
)
type CategoryDB struct{
    Id string `json:"id"`
    Name string `json:"name"`
    Icon string `json:"icon"`
    Url string `json:"url"`
}

type Category struct{
    coll *mongodb.MongoDBCollection
    dbname string
    collName string
    dbc influxbase.DBConfig
    mdb * mongodb.MongoDB
}

//通过配置获取实体连接
func GetCategoryInstance() *Category{
    //log.Printf("dbc     %v\n",dbc)
    var dbc influxbase.DBConfig
    dbc.Host = "192.168.10.10:27017"

    var mongo  Category
    if !mongo.Init(dbc) {
        log.Printf("mongo     null\n")
        return nil
    }
    return &mongo
}
func (p* Category) Init(dbc influxbase.DBConfig) bool{
    p.dbname = "noah"
    p.collName = "categories"
    p.dbc = dbc
    p.mdb = mongodb.GetMongoDBInstance(dbc)
    var err error
    p.coll,err = p.mdb.GetCollection(p.dbname,p.collName)

    if influxbase.HasError(err){

        return false
    }
    return true
}

func(p* Category) GetCategory() (string,error){


    var result []CategoryDB//bson.M
    conditions := bson.M{}//bson.M{"createdtime":bson.M{"$gt":time},
    //"lastexcutetime":bson.M{"$lt":curtime},//下次执行时间 小于当前时间
    //"excutesstatus":bson.M{"$in":[]base.EXCUTE_STATUS{base.NEW_STATUS,base.ERROR_STATUS}}}
    //"

    err := p.coll.Find(conditions).Sort("_id").All(&result)
    jsonstring,_ := json.Marshal(result)
    if influxbase.HasError(err) {
        return "error",err
    }
    log.Println("jsonstring:::",string(jsonstring[:]))
    return string(jsonstring[:]),nil
}

//主要用于Json解释
func (n * CategoryDB) ToJson() (string,bool){
    //var jsonstring []byte
    jsonstring,_ := json.Marshal(*n)
    //fmt.Println(jsonstring)
    return string(jsonstring[:]),true
}
