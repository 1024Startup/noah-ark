package base

import "encoding/json"
import "time"
import "influx.io/influxbase"

//枚举执行状态
type DBROW_STATUS int64
const (
	DBROW_NEW_STATUS DBROW_STATUS =  iota
	DBROW_SUCCESS_STATUS
	DBROW_DELETE_STATUS
)
//监听信息
type ListenerInfo struct{
    Name string `json:"name"` //标识本监听的名称(全局唯一)
    Event string `json:"event"` //监听的事件名称
    NotifyUrl string `json:"notifyurl"` //通知URL
    Type string `json:"type"` //类型，用于确认使用哪种，默认是NormalListener
}

type ListenerInfoDb struct{
    ListenerInfo
	CreatedTime int64
    LastNotifiedTime int64
    Status DBROW_STATUS
	Version int64 //版本号，数字方便加载不同的结构
}
const  ListenerInfoDbVersion = 1 //ListenerInfoDb版本号

type IListener interface{
    SendNotify(n *NotifyDetail) bool
	SetNotify(notify INotify)
	Save(dbc influxbase.DBConfig) bool
}
//主要用于Json解释
func (n* ListenerInfo) FromJson(jsonstring string) bool{

    err := json.Unmarshal([]byte(jsonstring), n)
	//log.Println("%s,/n %v",jsonstring,n)

    return !influxbase.HasError(err)
}
//构造ListenerInfoDb
func GetListenerInfoDbFromListener(l * ListenerInfo) *ListenerInfoDb{
    var lsdb ListenerInfoDb
    lsdb.ListenerInfo = *l
    lsdb.LastNotifiedTime =  time.Now().Unix()
	lsdb.CreatedTime =  time.Now().Unix()
    lsdb.Status = DBROW_NEW_STATUS
	lsdb.Version = ListenerInfoDbVersion
    return &lsdb
}
