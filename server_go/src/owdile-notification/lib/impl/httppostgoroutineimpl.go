package impl

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

type PostLoopGoriutine struct{
    name string
    num int
}
const postLoopGoriutineName = "PostLoopGoriutine"
func(l PostLoopGoriutine) Name() string{
    return fmt.Sprintf("%s-%d",l.name,l.num)
}
func(l PostLoopGoriutine) Run() bool{
    //查询列表，获取数据，处理出错历史数据
    log.Printf("run")
    return true
}
func(l PostLoopGoriutine) Init() bool{
    //查询列表，获取数据
    log.Printf("init")
    return true
}
func(l PostLoopGoriutine) SleepMillisecond() time.Duration{
    //查询列表，获取数据
    return time.Millisecond *1000//1000毫秒一次
}
func(l PostLoopGoriutine) Call(para interface{}) bool{
    pd ,err :=para.(postData)
    log.Println("ssssfdfdsfs:%v",para)
    if(!err){
        return false
    }
    //SaveNotifyEvent(pd.Url,&pd.NotifyDetail,ret)
    ret := PostImpl(pd.Url,&pd.NotifyDetail)
    //应该使用观察者合适
    SaveNotifyEvent(pd.Url,&pd.NotifyDetail,ret)
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
//线程池 maxGoroutineCount 个线程
func initPostPool()(*PostLoopGoriutine, bool){
    if(!hasInit){
        hasInit = true

        for i:=0;i<maxGoroutineCount;i++ {
            postLoopGoriutines[i] = &PostLoopGoriutine{postLoopGoriutineName,i}
        }

    }
    count ++
    if count == maxGoroutineCount {
        count = 0
    }

    return postLoopGoriutines[count],true;
}
//新线程
func initPostNew() (*PostLoopGoriutine ,bool){
    newCount ++
    pg := &PostLoopGoriutine{postLoopGoriutineName,newCount}//new(PostLoopGoriutine)
    return pg,true
}



func Post(strurl string, notifydetail *base.NotifyDetail) bool {
    return PostPool(strurl,notifydetail)
    //
    pg,e := initPostNew()
    if(e){
        goroutineHelper := GetLoopGoroutineHelperInstance()

        pd := postData{strurl,(*notifydetail)}
        return goroutineHelper.CallLoopGoroutineOnce(pg,pd)
    }

    return false
}
func PostPool(strurl string, notifydetail *base.NotifyDetail) bool {
    pg,e := initPostPool()
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
