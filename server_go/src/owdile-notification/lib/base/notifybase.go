package base
import "encoding/json"
import "log"
import "influx.io/influxbase"
import "gopkg.in/mgo.v2/bson"
import "time"
//枚举执行状态
type EXCUTE_STATUS int
const (
	NEW_STATUS EXCUTE_STATUS =  iota
	SUCCESS_STATUS
	DELETE_STATUS
	HANDLEING_STATUS
	EXPIRE_STATUS
	ERROR_STATUS
)
type NotifyInfo struct{
	Name string `json:"name"`  //标识发送方的名称
	Type string `json:"type"` //类型，用于确认使用哪种，默认是NormalNotify
	Event string `json:"event"`//事件名，用于标识此类通知
	Desc string `json:"desc"`//描述，用于描述这个通知
}
type NotifyInfoDb struct{
	NotifyInfo
	Status int64 //状态 2删除 0新建
	CreatedTime int64 //创建时间
	Version int64 //版本号，数字方便加载不同的结构
}
const  NotifyInfoDbVersion = 1 //NotifyInfoDb版本号
//通知信息，调用时需要传入的消息结构体
type NotifyDetail struct{
	Event string `json:"event"`//事件名，用于标识此类通知
	Data string `json:"data"`//透传的参数，传输到监听端
}
const  NotifyDetailDbVersion = 1 //NotifyDetailDbVersion版本号
//通知信息数据库相关结构体
type NotifyDetailDb struct
{
	NotifyDetail
	Url string//调用目标
	CreatedTime int64 //创建时间
	LastExcuteTime int64 //最后执行时间
	NextExcuteTime int64 //下镒执行时间
	ExcuteStatus  EXCUTE_STATUS//执行状态
	ExcuteTimes int//执行次数
	LastErrorMsg string//最后错误记录
	Version int64 //版本号，数字方便加载不同的结构
	Id bson.ObjectId `bjson:"_id"`
}
//通知接口，定义通知的方法
type INotify interface{
	SendNotify(n *NotifyDetail) bool //发送通知方法，
	AddListener(name string,l * IListener) bool //增加监听
	Save(dbc influxbase.DBConfig) bool
}

//主要用于Json解释
func (n* NotifyInfo) FromJson(jsonstring string) bool{

    err := json.Unmarshal([]byte(jsonstring), n)
    /*fmt.Println(*n)
    fmt.Println(err)*/
	return influxbase.HasError(err)
}

//主要用于Json解释
func (n* NotifyDetail) FromJson(jsonstring string) bool{

    err := json.Unmarshal([]byte(jsonstring), n)
    /*fmt.Println(*n)
    fmt.Println(err)*/
	return influxbase.HasError(err)
}
//主要用于Json解释
func (n* NotifyDetail) ToUrlValue() (*map[string][]string,bool){
	m := make(map[string][]string)
	m["event"] = []string {n.Event}
	m["data"] = []string {n.Data}
    log.Println("fgfffff,%v",m)
    return &m,true
}

//主要用于Json解释
func (n NotifyInfo) ToJson() (string,bool){
    //var jsonstring []byte
    jsonstring,_ := json.Marshal(n)
    //fmt.Println(jsonstring)
    return string(jsonstring[:]),true
}

//主要用于Json解释
func (n NotifyInfoDb) ToM() (bson.M,bool){
    //var jsonstring []byte

    d := bson.M{"Name": n.Name, "Type": n.Type,"Event":n.Event,"Desc":n.Desc,"Status":n.Status,"CreatedTime":n.CreatedTime}
    //fmt.Println(jsonstring)
    return d,true
}

//构造ListenerInfoDb
func GetNotifyInfoDbFromNotify(l * NotifyInfo) *NotifyInfoDb{
    var ntfInfoDb NotifyInfoDb
	ntfInfoDb.NotifyInfo = *l
	ntfInfoDb.CreatedTime = time.Now().Unix()
	ntfInfoDb.Version = NotifyInfoDbVersion
    return &ntfInfoDb
}

func GetNotifyDetailDbFromNotifyDetail(url string,l * NotifyDetail) *NotifyDetailDb{
    var ntfDetailDb NotifyDetailDb
	ntfDetailDb.NotifyDetail = *l
	ntfDetailDb.CreatedTime = time.Now().Unix()
	ntfDetailDb.LastExcuteTime = ntfDetailDb.CreatedTime
	ntfDetailDb.NextExcuteTime = ntfDetailDb.CreatedTime
	ntfDetailDb.LastErrorMsg ="begin"
	ntfDetailDb.ExcuteStatus = NEW_STATUS
	ntfDetailDb.ExcuteTimes = 0
	ntfDetailDb.Version = NotifyDetailDbVersion
	ntfDetailDb.Url = url
	ntfDetailDb.Id = bson.NewObjectId()
    return &ntfDetailDb
}
