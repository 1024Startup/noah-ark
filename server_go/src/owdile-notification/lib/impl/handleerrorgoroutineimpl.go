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

type MasterErrorLoopGoriutine struct{
    name string
    num int
    millsecond time.Duration//int64
    mongoconf influxbase.DBConfig

}
var expire = []int64{5,10,20,60,120,300,1800,3600}
const errorLoopGoriutineName = "MasterErrorLoopGoriutine"
func(l MasterErrorLoopGoriutine) Name() string{
    return fmt.Sprintf("%s-%d",l.name,l.num)
}
func(l MasterErrorLoopGoriutine) Run() bool{
    //查询列表，获取数据，处理出错历史数据
    l.handleUnfinishNotify()
    return true
}
func(l* MasterErrorLoopGoriutine) Init() bool{
    //查询列表，获取数据
    //l.mongoconf = base.GetDbConfigByType(base.CONFIG_MONGO)
    return true
}
func(l MasterErrorLoopGoriutine) SleepMillisecond() time.Duration{
    //查询列表，获取数据
    return time.Millisecond *l.millsecond//num毫秒一次
}
func(l MasterErrorLoopGoriutine) Call(para interface{}) bool{

    return true;
}
func(l MasterErrorLoopGoriutine) handleUnfinishNotify() bool{
    //查找出数据列表，
    //7天前（处理7天内的数据）

    event:= model.GetEventMongoInstance(l.mongoconf)
    detailDBList,ok :=event.GetUnfinishNotifyEventsFromDB()
    if !ok {
        return false
    }
    for _, v := range detailDBList {

        //执行
        ret := PostImpl(v.Url,&v.NotifyDetail)

        expireIndex := v.ExcuteTimes
        if expireIndex >= len(expire) {
            expireIndex = len(expire) -1
        }


        v.LastExcuteTime = time.Now().Unix()+expire[expireIndex]
        v.ExcuteTimes += 1
        v.ExcuteStatus = base.ERROR_STATUS
        if ret{
            v.ExcuteStatus = base.SUCCESS_STATUS
        }
        //更新数据
        event.SaveNotifyEventToDB(v)

    }
    return true
}
var masterLoopGoriutine *MasterErrorLoopGoriutine

//重新生成错误处理线程
func ReInitMasterErrorOrgoroutine() bool{
    if(masterLoopGoriutine != nil){
        goroutineHelper := GetLoopGoroutineHelperInstance()
        goroutineHelper.StopLoopGoroutine(masterLoopGoriutine)
    }
    masterLoopGoriutine = &MasterErrorLoopGoriutine{errorLoopGoriutineName,0,5000,base.GetDbConfigByType(base.CONFIG_MONGO)}//new(MongoLoopGoriutine)
    goroutineHelper := GetLoopGoroutineHelperInstance()
    goroutineHelper.StartLoopGoroutine(masterLoopGoriutine)
    return true;
}
