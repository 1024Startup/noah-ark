package manager

import
(
    "owdile-notification/lib/base"
    "owdile-notification/lib/factory"
    "owdile-notification/lib/model"
    "owdile-notification/lib/impl"
    "influx.io/influxbase"
    //"container/list"
    "log"
)
//消息管理类
/**
管理通知列表
非线程安全
**/
type NotifyManager struct{
    notifies map[string]*base.INotify//通知列表
    init bool
    mongoconf influxbase.DBConfig
    redisconf influxbase.DBConfig
}

var notifyInstance  NotifyManager

func Instance() *NotifyManager{
    notifyInstance.Init()
    return &notifyInstance
}

//初始化方法
func (p* NotifyManager) Init() bool{

    if !p.init {
        p.notifies = make(map[string]*base.INotify)//新建map
        p.mongoconf = base.GetDbConfigByType(base.CONFIG_MONGO)
        p.redisconf = base.GetDbConfigByType(base.CONFIG_REDIS)
        p.init = true
        //impl.ReInitSlaveErrorOrgoroutine()
        impl.ReInitMasterErrorOrgoroutine()
        p.reloadNotify()
    }
    return true
}
//重新加载数据
func (p* NotifyManager) ReloadNotify() bool{
    p.Init()
    return p.reloadNotify()
}
//重新加载数据
func (p* NotifyManager) reloadNotify() bool{
    //加载notify
    nmodel := model.GetNotifyInstance(p.mongoconf)
    notifies,ok := nmodel.GetNotifyInfosFromDB()
    if !ok {
        return false
    }
    for i := 0; i< len(notifies); i++{
        notify := factory.GetNotify(&notifies[i].NotifyInfo)
        p.notifies[notifies[i].Event] = &notify;
    }
    //加载listener
    lmodel := model.GetListenerInstance(p.mongoconf)
    listeners,ok := lmodel.GetListenersFromDB()
    if !ok {
        return false
    }
    for i := 0; i< len(listeners); i++{

        if p.notifies[listeners[i].ListenerInfo.Event] == nil {
            continue //无此事件，监听无用
        }
        //加入到notify中
        listener := factory.GetListener(&listeners[i].ListenerInfo)
        listenerName := listeners[i].ListenerInfo.Name
        (*p.notifies[listeners[i].ListenerInfo.Event]).AddListener(listenerName,&listener)
        //p.listeners[listeners[i].Event] = notify;
    }
    return true
}


//发送通知
func (p* NotifyManager) SendNotify(notifydetailjson string) bool{
    p.Init()
    notifydetail := p.createNotifyDetail(notifydetailjson)
    if p.notifies[notifydetail.Event] == nil {
        return false//没有此事件
    }
    notify := p.notifies[notifydetail.Event]
    log.Println("notify :%v",*notify)
    return (*notify).SendNotify(notifydetail)
}

//注册通知
func (p* NotifyManager) RegisterNotify(notifyinfojson string) bool{
    p.Init()

    //查询是否已经存在，如果已经存在则返回true
    n := p.createNotifyInfo(notifyinfojson)

    //log.Printf("notifyinfojson %s\n n:  %v \n",notifyinfojson,n)

    if p.notifies[n.Event] != nil {
        return true//已经注册
    }

    notify := factory.GetNotify(n)
    p.notifies[n.Event] = &notify

    return notify.Save(p.mongoconf)

}

func (p * NotifyManager) RegisterListener(listenerinfojson string) bool{
    p.Init()
    //可能有几种情况，无此Notify,Listener已经存在，数据库操作错误
    //create listenerinfo
    listenerinfo,success := p.createListenerInfo(listenerinfojson)

    if !success {
        return false
    }
    log.Println("%v",listenerinfo)
    //find noity
    inotify := p.notifies[listenerinfo.Event]
    if inotify == nil {
        //不存在此通知
        return false
    }
    //find if exist

    //create listener
    listener := factory.GetListener(listenerinfo)
    //relationship
    //listener.SetNotify(inotify)
    (*inotify).AddListener(listenerinfo.Name,&listener)
    listener.Save(p.mongoconf)
    return true
}

func (p* NotifyManager) createNotifyInfo(notifyinfojson string) *base.NotifyInfo{
    //转jSON到结构体
    var notifyinfo base.NotifyInfo
    notifyinfo.FromJson(notifyinfojson)
    return &notifyinfo
}

func (p* NotifyManager) createNotifyDetail(notifydetailjson string) *base.NotifyDetail{
    //转jSON到结构体
    var notifydetail base.NotifyDetail
    notifydetail.FromJson(notifydetailjson)
    return &notifydetail
}


func (p* NotifyManager) createListenerInfo(listenerinfojson string) (*base.ListenerInfo,bool){
    //转jSON到结构体
    var listenerinfo base.ListenerInfo
    success := listenerinfo.FromJson(listenerinfojson)
    return &listenerinfo,success
}
