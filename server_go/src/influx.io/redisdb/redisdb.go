package redisdb

import (
    "github.com/astaxie/goredis"
    "influx.io/influxbase"
    "errors"
)
//统一连接配置结构，如有增加参数，使用冗余方式加入
type DBConfig struct{
    Host string//"192.168.10.10:27017"mongodb  //"127.0.0.1:6379"redis
}

//Redis
type OrmRedisData struct{
    Key string
    Value string
    Expire int64
}
type IOrmRedis interface{
    ConnectDb(db influxbase.DBConfig) error
    SetData(key string,value []byte,expire int) error
    GetData(key string) ([]byte,error)
}
type RedisDB struct{
    goredis.Client

}
//连接数据库地址
func (p * RedisDB)ConnectDb(db influxbase.DBConfig) error{
    p.Client.Addr = db.Host
    return nil
}
//设置值，同时设置超时时间
func (p * RedisDB)SetData(key string,value []byte,expire int64) error{
    //var err error
    var success bool
    err := p.Client.Set(key,value)
    if influxbase.HasError(err) {
        return err
    }
    success,_ =p.Client.Expire(key,expire)
    if !success {
         p.Client.Del(key)//??是否要删除？
         return errors.New("set Expire error")
    }
    return nil
}
//获取值
func (p * RedisDB)GetData(key string)([]byte,error) {
    return p.Client.Get(key)
}
