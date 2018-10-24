package impl
import
(
	_"fmt"
    "owdile-notification/lib/base"
	"owdile-notification/lib/model"
	"influx.io/influxbase"
)

type BaseListener struct{
	base.ListenerInfo
    notify base.INotify
}

func (nl  BaseListener) SetNotify(notify base.INotify) {
	nl.notify = notify
}
func (nl  BaseListener) SendNotify(n *base.NotifyDetail) bool{
	db := base.GetDbConfigByType(base.CONFIG_MONGO)
	//保存到数据库
	ret := SaveNotifyEventImpl(db,nl.NotifyUrl ,n)//保存数据到本地，认为已经成功
	if !ret {
		return false
	}
	//发送
	//Post(nl.NotifyUrl ,n)
	return true
    //return true
}
//保存到数据库
func (nl  BaseListener) Save(dbc influxbase.DBConfig) bool{

	lstInfoDb := base.GetListenerInfoDbFromListener(&nl.ListenerInfo)

	model := model.GetListenerInstance(dbc)

	return model.SaveNotifyInfosToDB(*lstInfoDb)
}

func (nl  BaseListener) load() bool{
    return true
}

type NormalListener struct{
    BaseListener
}
