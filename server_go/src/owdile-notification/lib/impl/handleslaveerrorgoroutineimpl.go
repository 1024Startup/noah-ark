package impl
import (
    _"net/http"
    "owdile-notification/lib/model"
    "owdile-notification/lib/base"
    _"log"
    "time"
    "fmt"
    _"io/ioutil"
    "influx.io/influxbase"
    _"math/rand"
    _"gopkg.in/mgo.v2/bson"
    _"encoding/json"
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

type SlaveErrorLoopGoriutine struct{
    name string
    num int
    millsecond time.Duration//int64
    mongoconf influxbase.DBConfig

}

const SlaveLoopGoriutineName = "SlaveErrorLoopGoriutine"
func(l SlaveErrorLoopGoriutine) Name() string{
    return fmt.Sprintf("%s-%d",l.name,l.num)
}
func(l SlaveErrorLoopGoriutine) Run() bool{
    //查询列表，获取数据，处理出错历史数据

    return true
}
func(l* SlaveErrorLoopGoriutine) Init() bool{
    //查询列表，获取数据
    //l.mongoconf = base.GetDbConfigByType(base.CONFIG_MONGO)
    return true
}
func(l SlaveErrorLoopGoriutine) SleepMillisecond() time.Duration{
    //查询列表，获取数据
    return time.Millisecond *l.millsecond//num毫秒一次
}
func(l SlaveErrorLoopGoriutine) Call(para interface{}) bool{

    return true;
}
func(l SlaveErrorLoopGoriutine) handleUnfinishNotify() bool{
    //查找出数据列表，
    //7天前（处理7天内的数据）
    //数据搬迁，从其他数据库搬到master数据库。保证数据必达
    slaveEvent:= model.GetEventMongoInstance(l.mongoconf)

    detailDBList,ok :=slaveEvent.GetUnfinishNotifyEventsFromDB()
    if ok {
        return false
    }
    masterEvent:= model.GetEventMongoInstance(base.GetDbConfigByType(base.CONFIG_MONGO))
    for _, v := range detailDBList {

        //执行
        /*ret := PostImpl(v.Url,&v.NotifyDetail)

        expireIndex := v.ExcuteTimes
        if expireIndex >= len(expire) {
            expireIndex = len(expire) -1
        }


        v.LastExcuteTime = time.Now().Unix()+expire[expireIndex]
        v.ExcuteTimes += 1
        v.ExcuteStatus = base.ERROR_STATUS
        if ret{
            v.ExcuteStatus = base.SUCCESS_STATUS
        }*/
        //更新数据
        masterEvent.SaveNotifyEventToDB(v)
        //删除数据
        v.ExcuteStatus = base.DELETE_STATUS
        slaveEvent.SaveNotifyEventToDB(v)

    }
    //删除
    return true
}
var slaveLoopGoriutines []*SlaveErrorLoopGoriutine

//重新生成错误处理线程
func ReInitSlaveErrorOrgoroutine() bool{
    goroutineHelper := GetLoopGoroutineHelperInstance()
    if(slaveLoopGoriutines != nil){
        for _, v := range  slaveLoopGoriutines{
            goroutineHelper.StopLoopGoroutine(v)
        }
    }
    //读取数据库配置，每个配置开启一个协程
    dbcs := base.GetSlaveMongoDbConfigs()
    slaveLoopGoriutines = make([]*SlaveErrorLoopGoriutine,len(dbcs))
    for i,v := range dbcs{
        slaveLoopGoriutine :=  &SlaveErrorLoopGoriutine{SlaveLoopGoriutineName,i,5000,v}
        slaveLoopGoriutines[i] = slaveLoopGoriutine
        goroutineHelper.StartLoopGoroutine(slaveLoopGoriutine)
    }
    return true;
}
