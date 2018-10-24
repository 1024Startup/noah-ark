package impl
import
(
	"log"
	_"time"
    "owdile-notification/lib/base"
	//"container/list"
	"influx.io/influxbase"
	"owdile-notification/lib/model"
)

//传统通知数据结构基础类
type BaseNotify struct{
	NotifyInfo base.NotifyInfo
	listeners map[string]*base.IListener//监听列表
	init bool
	//Init func() bool
}
//初始化方法
func (bn *BaseNotify) Init() bool{

    if !bn.init {
        bn.listeners = make(map[string]*base.IListener)//新建map
        bn.init = true
    }
    return true
}
func (bn* BaseNotify) SendNotify(n *base.NotifyDetail) bool {
	bn.Init()
	//log.Printf("hello,dd world\n")
	//遍历listener，保存到redis
	for _,v := range bn.listeners {
		(*v).SendNotify(n)
	}
	//保存到redis
	//
	return true
}

func (bn *BaseNotify) AddListener(name string,l * base.IListener) bool{
	bn.Init()
	//log.Printf("hello,dd world\n")
	bn.listeners[name] = l
	log.Println("base notify %v",bn)
	(*l).SetNotify(bn)
	return  true
}


func (pn* BaseNotify)Load() bool {//加载方法
	log.Printf("load\n")
	return true
}
/****
*通知信息数据库相关结构体
type NotifyInfoDb struct
{
	NotifyInfo
	CreatedTime int64 //创建时间
	LastExcuteTime int64 //最后执行时间
	ExcuteStatus  EXCUTE_STATUS//执行状态
	LastErrorMsg string//最后错误记录
}
*/
func (pn* BaseNotify) Save(dbc influxbase.DBConfig) bool{//保存方法
	pn.Init()
	//fmt.Printf("save\n")

	ntfInfoDb := base.GetNotifyInfoDbFromNotify(&pn.NotifyInfo)

	log.Printf("pn.notifyInfo: %v\n",pn.NotifyInfo)

	notifymodel := model.GetNotifyInstance(dbc)

	log.Printf("notifymodel: %v\n",notifymodel)
	return notifymodel.SaveNotifyInfosToDB(*ntfInfoDb)

}


func (pn* BaseNotify) updateNotifyResult() (n []int,success bool){//更新方法*
	log.Printf("updateNotifyResult\n")
	var slice = make([]int, 10)
	return slice,true
}

type NormalNotify struct{
	BaseNotify
}
/*func (p NormalNotify) AddListener(l * base.IListener) bool {
	return true
}
*/
