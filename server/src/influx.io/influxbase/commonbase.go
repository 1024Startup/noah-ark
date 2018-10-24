package influxbase

import "log"
import "runtime/debug"
import "encoding/json"
import "reflect"

/**
错误的处理
**/
//简单把error转成bool ，如有错误则输出
func HasError(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

//exception,把error抛出
func ThorwError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

//接收并打印信息，需要使用defer调用，再代码最后执行
func RecoverError() bool {
	if p := recover(); p != nil {
		log.Printf("panic recover! p: %v", p)
		debug.PrintStack()
		return false
	}
	return true
}

type DBConfig struct {
	Host string //"192.168.10.10:27017"mongodb  //"127.0.0.1:6379"redis
}

type ResposeData struct {
	Ret  interface{} //"0" 成功 “1”失败
	Msg  string      //描述
	Data string      //返回数据
}

//主要用于Json解释
func JsonToResposeData(jsonstring string) ResposeData {
	var res ResposeData
	res.Ret = "-1"
	res.Msg = "初始值"
	res.FromJson(jsonstring)
	return res
}

//主要用于Json解释
func (n *ResposeData) IsSuccess() bool {

	switch n.Ret.(type) {
	case int64:
		value, ok := n.Ret.(int)
		if ok {
			return value == 0
		}
		return false

	case string:
		value, ok := n.Ret.(string)
		if ok {
			return value == "0"
		}
		return false

	case float64:
		value, ok := n.Ret.(float64)
		if ok {
			return value == 0
		}
		return false
	default:
		log.Println("%v", reflect.TypeOf(n.Ret))
		return false
	}

	return false
	/*fmt.Println(*n)
	  fmt.Println(err)*/
}

//主要用于Json解释
func (n *ResposeData) FromJson(jsonstring string) bool {

	err := json.Unmarshal([]byte(jsonstring), n)
	/*fmt.Println(*n)
	  fmt.Println(err)*/
	return HasError(err)
}

//主要用于Json解释
func (n ResposeData) ToJson() (string, bool) {
	//var jsonstring []byte
	jsonstring, _ := json.Marshal(n)
	//fmt.Println(jsonstring)
	return string(jsonstring[:]), true
}
