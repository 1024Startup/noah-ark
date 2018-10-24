package impl

import "time"
import "log"
import "influx.io/influxbase"
import _"math/rand"

//消息循环协助类
type ILoopGoroutine interface{
    Init() bool//初始化
    Run() bool//执行方法
    Call(para interface{}) bool//执行方法
    Name() string//唯一名称
    SleepMillisecond() time.Duration
    //Heartbeat() bool //心跳，如果心跳返回false,关闭线程
}
type GoroutineControlPara struct{
    Done bool//为1时结束
    Sync bool//同步
}
type goroutoneTranslationChan struct{
    Para *chan interface{}
    Resp *chan bool
}
/*type GoroutineResult struct{
    ret int//0为成功
    msg string
    data interface{}//值
}*/
type GoroutineParaBase struct{
    GoroutineControlPara
    UserPara interface{}
}
func loopGoroutine(para chan interface{},resp chan bool,run ILoopGoroutine) bool{
    if !run.Init() {
        return false
    }
    defer influxbase.RecoverError()//线程问题，处理
    for {
        select {
        case p := <-para://读取值，

            gpb,s :=p.(GoroutineParaBase)
            if (!s) {
                break;
            }
            if(gpb.Done) {
                log.Println("exiting...%v",run.Name())
                close(resp)
                return true;
            }
            log.Println("threadhame:%s para:v%v",run.Name(),gpb.UserPara)
            ret := run.Call(gpb.UserPara)
            if(gpb.Sync){
                resp <- ret
            }
            break
		default://无值，走此分支
            run.Run()
            millisecond := run.SleepMillisecond()
            time.Sleep( millisecond)
		}

    }
    return true
}

var loopGoroutineHelperInstance LoopGoroutineHelper
func GetLoopGoroutineHelperInstance() LoopGoroutineHelper{
    loopGoroutineHelperInstance.Init()
    return loopGoroutineHelperInstance
}

type LoopGoroutineHelper struct{
    NameParaMap map[string] *goroutoneTranslationChan
    init bool
}
func (lg *LoopGoroutineHelper)Init()bool{
    if !lg.init {
        lg.init = true
        lg.NameParaMap = make(map[string] *goroutoneTranslationChan)//make(map[string]*chan interface{})
    }
    return true
}

func (lg *LoopGoroutineHelper)StartLoopGoroutine( g ILoopGoroutine )bool{
    lg.Init()
    name := g.Name()
    if  lg.NameParaMap[name] == nil {
        p := make(chan interface{}, 1)
        r := make(chan bool, 1)//r在协程内部关闭
        gtChan := goroutoneTranslationChan{&p,&r}
        lg.NameParaMap[name] = &gtChan
        go loopGoroutine(p,r,g)
    }
    return true
}

func (lg *LoopGoroutineHelper)CallLoopGoroutine( g ILoopGoroutine,para interface{} )bool{
    lg.Init()
    name := g.Name()
    if  lg.NameParaMap[name] != nil {//发送到协程
        p := lg.NameParaMap[name].Para
        goroutineParaBase := GoroutineParaBase{GoroutineControlPara{false,false},para}
        (*p) <- goroutineParaBase
        return true
    }
    return false
}
//同步调用
func (lg *LoopGoroutineHelper)CallLoopGoroutineSync( g ILoopGoroutine,para interface{} )bool{
    lg.Init()
    name := g.Name()
    if  lg.NameParaMap[name] != nil {//发送到协程
        p := lg.NameParaMap[name].Para
        goroutineParaBase := GoroutineParaBase{GoroutineControlPara{false,false},para}
        (*p) <- goroutineParaBase
        r := lg.NameParaMap[name].Resp
        ret := <-*r
        return ret
    }
    return false
}

func (lg *LoopGoroutineHelper)StopLoopGoroutine( g ILoopGoroutine )bool{
    lg.Init()
    name := g.Name()
    if  lg.NameParaMap[name] == nil {
        return false
    }else{
        defer influxbase.RecoverError()
        p := lg.NameParaMap[name].Para
        delete(lg.NameParaMap,name)
        goroutineParaBase := GoroutineParaBase{GoroutineControlPara{true,false},nil}
        (*p) <- goroutineParaBase
        close(*p)
        return true
    }
}

func (lg *LoopGoroutineHelper)CallLoopGoroutineOnce( g ILoopGoroutine,para interface{} )bool{
    lg.StartLoopGoroutine(g)
    lg.CallLoopGoroutine(g,para);
    defer lg.StopLoopGoroutine(g)
    return true
}
