package lib

import (
    "net/http"
    "owdile-notification/lib/base"
    "log"
    "time"
    "fmt"
    "io/ioutil"
    "influx.io/influxbase"
    _"math/rand"
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

type NotifyLoopGoriutine struct{
    name string
    num int
}
const notifyLoopGoriutine = "NotifyLoopGoriutine"
func(l NotifyLoopGoriutine) Name() string{
    return fmt.Sprintf("%s-%d",l.name,l.num)
}
func(l NotifyLoopGoriutine) Run() bool{
    //查询列表，获取数据，处理出错历史数据
    log.Printf("run")
    return true
}
func(l NotifyLoopGoriutine) Init() bool{
    //查询列表，获取数据
    log.Printf("init")
    return true
}
func(l NotifyLoopGoriutine) SleepMillisecond() time.Duration{
    //查询列表，获取数据
    return time.Millisecond *1000//1000毫秒一次
}
func(l NotifyLoopGoriutine) Call(para interface{}) bool{
    pd ,err :=para.(postData)
    log.Println("ssssfdfdsfs:%v",para)
    if(!err){
        return false
    }
    //SaveNotifyEvent(pd.Url,&pd.NotifyDetail,ret)
    ret := PostImpl(pd.Url,&pd.NotifyDetail)

    return ret;
}
type postData struct{
    Url string
    NotifyDetail  base.NotifyDetail
}
var goroutineCount int64
const maxGoroutineCount = 20 //最大线程数
var postLoopGoriutine *PostLoopGoriutine
var postLoopGoriutines [maxGoroutineCount]*PostLoopGoriutine
var hasInit bool
var count int
var newCount int
//单一线程
func initPost() bool{
    if(!hasInit){
        hasInit = true
        postLoopGoriutine = &PostLoopGoriutine{postLoopGoriutineName,0}//new(PostLoopGoriutine)
        goroutineHelper := GetLoopGoroutineHelperInstance()
        goroutineHelper.StartLoopGoroutine(postLoopGoriutine)
    }
    return true;
}


func SendNotify(strurl string, notifydetail *base.NotifyDetail) bool {

    pg,e := initPostNew()
    if(e){
        goroutineHelper := GetLoopGoroutineHelperInstance()

        pd := postData{strurl,(*notifydetail)}
        return goroutineHelper.CallLoopGoroutineOnce(pg,pd)
    }

    return false
}

func PostImpl(strurl string, notifydetail *base.NotifyDetail) bool {
    log.Println("sssss:::::%v",notifydetail)
    json,_ := notifydetail.ToUrlValue()
    log.Println("json::%s,",json)
    resp, err := http.PostForm(strurl,*json)


    if influxbase.HasError(err) {
        return false
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if influxbase.HasError(err) {
        return false
    }
    //log.Println("%v",body)
    str := string(body[:])
    log.Println("%v",str)
    res := influxbase.JsonToResposeData(str)
    log.Println("%v",res)
    log.Println("%v",res.IsSuccess())
    return res.IsSuccess()
}
