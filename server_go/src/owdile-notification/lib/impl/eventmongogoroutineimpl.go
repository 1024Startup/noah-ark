package impl

import (
    "owdile-notification/lib/base"
    _"log"
    "time"
    "fmt"
    "influx.io/influxbase"
    "owdile-notification/lib/model"
    "container/list"
)
/**
type ILoopGoroutine interface{
    Init() bool//初始化
    Run() bool//执行方法
    Name() string//唯一名称
    SleepMillisecond() time.Duration
    //Heartbeat() bool //心跳，如果心跳返回false,关闭线程
}
**/

type MongoLoopGoriutine struct{
    name string
    num int
    mongoconf influxbase.DBConfig
    notifyEvents *list.List
}
const mongoLoopGoriutineName = "MongoLoopGoriutine"
//var //base.NotifyDetail//数据列表
func(l MongoLoopGoriutine) Name() string{
    return fmt.Sprintf("%s-%d",l.name,l.num)
}
func(l* MongoLoopGoriutine) Run() bool{
    //查询列表，获取数据
    for element := l.notifyEvents.Front(); element != nil;  {
        pd ,err :=element.Value.(NotifyMongoData)
        if(!err){
            return true
        }
        status := base.SUCCESS_STATUS
        if !pd.Ok {
            status = base.ERROR_STATUS
        }
        //保存到数据库
        ret := saveNotifyEventImpl(l.mongoconf,pd.Url,&pd.NotifyDetail,status)
        enext := element.Next()
        if ret {
            l.notifyEvents.Remove(element)
        }
        element = enext
    }
    return true
}
func(l* MongoLoopGoriutine) Init() bool{
    //查询列表，获取数据
    l.mongoconf = base.GetDbConfigByType(base.CONFIG_MONGO)
    l.notifyEvents = list.New()
    //l.redisconf = base.GetDbConfigByType(base.CONFIG_REDIS)
    return true
}
func(l MongoLoopGoriutine) SleepMillisecond() time.Duration{
    //查询列表，获取数据
    return time.Millisecond *3000//5000毫秒一次
}
func(l MongoLoopGoriutine) Call(para interface{}) bool{
    //保存到内存
    l.notifyEvents.PushBack(para)

    return true

}
type NotifyMongoData struct{
    Url string//
    Ok bool //是否成功
    NotifyDetail  base.NotifyDetail
}

var mongoLoopGoriutine *MongoLoopGoriutine
var MongoLoopGoriutines [maxGoroutineCount]*MongoLoopGoriutine
var initnogo bool

//单一线程
func initMongo() bool{
    if(!initnogo){
        initnogo = true
        mongoLoopGoriutine = &MongoLoopGoriutine{mongoLoopGoriutineName,0,base.GetDbConfigByType(base.CONFIG_MONGO),nil}//new(MongoLoopGoriutine)
        goroutineHelper := GetLoopGoroutineHelperInstance()
        goroutineHelper.StartLoopGoroutine(mongoLoopGoriutine)
    }
    return true;
}

//sliceInt[0], sliceInt[1:]...

func SaveNotifyEvent(strurl string, notifydetail *base.NotifyDetail,ok bool) bool {

    //
    e := initMongo()
    if(e){
        goroutineHelper := GetLoopGoroutineHelperInstance()

        pd := NotifyMongoData{strurl,ok,(*notifydetail)}
        return goroutineHelper.CallLoopGoroutine(mongoLoopGoriutine,pd)
    }

    return false
}
func SaveNotifyEventImpl(dbc influxbase.DBConfig,strurl string, notifydetail *base.NotifyDetail)bool{

    return saveNotifyEventImpl(dbc,strurl,notifydetail,base.NEW_STATUS)
}


func saveNotifyEventImpl(dbc influxbase.DBConfig,strurl string, notifydetail *base.NotifyDetail,status base.EXCUTE_STATUS)bool{
    pNotifyEventMongo := model.GetEventMongoInstance(dbc)

    detailDb := base.GetNotifyDetailDbFromNotifyDetail(strurl,notifydetail)

    detailDb.ExcuteStatus = status

    return pNotifyEventMongo.SaveNotifyEventToDB(*detailDb)
}
