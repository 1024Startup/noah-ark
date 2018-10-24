package main

import
(
	"log"
	"fmt"
    "owdile-notification/lib/manager"
	_"owdile-notification/lib/base"
	_"influx.io/influxbase"
)

func loadNotify(){
     manager.Instance().ReloadNotify()
}
func registerNotifyEvent(){
    var notifyinfojson = "{\"name\":\"puzzle\",\"event\":\"com.puzzle.sendorder\",\"desc\":\"定单通知\",\"type\":\"NormalNotify\"}";
     manager.Instance().RegisterNotify(notifyinfojson)
}

func registerListener(){
    var json = "{\"name\":\"admin.labels\",\"event\":\"com.puzzle.sendorder\",\"notifyurl\":\"http://admin.owdiex.me/execute/callback\",\"type\":\"NormalListener\"}";
     manager.Instance().RegisterListener(json)
}

func sendNotify(){
	var notifyinfojson = "{\"event\":\"com.puzzle.sendorder\",\"data\":\"衣服,100.00\"}";
	 manager.Instance().SendNotify(notifyinfojson)
}
func main() {
	log.Output(1,"sd")
	//defer influxbase.RecoverError()
	log.Println(" loadNotify  ")
	loadNotify()
	//log.Println(" registerNotifyEvent ")
    //registerNotifyEvent()
	//log.Println(" registerListener  ")
	//registerListener()
	var temp string
	for {
		log.Println(" sendNotify  ")
		sendNotify()
		fmt.Scanln(&temp)
		log.Println(temp)
	}

    log.Println(" end  ")
}
