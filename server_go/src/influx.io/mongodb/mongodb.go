package mongodb
import (
	"sync"
	"gopkg.in/mgo.v2"
	"influx.io/influxbase"
	"errors"
	"log"
)


type IOrmMongo interface{
    ConnectDb(db influxbase.DBConfig) error

}

//集合，主要操作对象
type MongoDBCollection struct{
		*mgo.Collection
}
type MongoDB struct {
	*mgo.Session
	lock  sync.Mutex
	dbconfig influxbase.DBConfig
}
//
//通过配置获取实体连接
func GetMongoDBInstance(dbc influxbase.DBConfig) *MongoDB{

    if HostSessionPoll == nil{
        HostSessionPoll = make(map[string] *MongoDB)
    }
    if HostSessionPoll[dbc.Host] == nil {
        var mongo  MongoDB
        err := mongo.ConnectDb(dbc)
		if influxbase.HasError(err){
			return nil
		}
        HostSessionPoll[dbc.Host] = &mongo
    }
    return HostSessionPoll[dbc.Host]
}


var HostSessionPoll map[string] *MongoDB

//连接DB，主要是建立与数据的连接
func (m * MongoDB) ConnectDb(db influxbase.DBConfig) error{
	m.lock.Lock()
	defer m.lock.Unlock()
	var err error
	if m.Session != nil  {//关闭重新连接
		m.Session.Close()
		m.Session = nil
	}
	m.Session, err = mgo.Dial(db.Host)  //连接数据库
	if influxbase.HasError(err) {
		return err
	}
	m.dbconfig = db
	return nil
}
//重新连接，如果断开需要重新连接
func (m * MongoDB) ReconnectDB() error{
	//m.lock.Lock()
	//defer m.lock.Unlock()
	return m.ConnectDb(m.dbconfig)
	if m.Session == nil  {//关闭重新连接
		return errors.New("session is nil")
	}
	err := m.Session.Ping()
	if influxbase.HasError(err){
		log.Println("error %v",err)
		return err
	}
	return m.ConnectDb(m.dbconfig)

}
//获取出Collection
func (m * MongoDB) GetCollection(dbname string , collectionname string) (*MongoDBCollection,error){
	//var db * mgo.Database
	db  := m.Session.DB(dbname)

	//var c * mgo.Collection
	c := db.C(collectionname)

	var mdb MongoDBCollection
	mdb.Collection = c

	return  &mdb,nil
}
