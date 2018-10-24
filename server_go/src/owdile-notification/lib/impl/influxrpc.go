package impl
//import "net/rpc"
import (
    "net"
    "fmt"
    "time"
    "influx.io/influxbase"
    "bufio"
    )

//us udp for connet 使用UDP可节省大量开销
type InfluxUDPLoopGoroutine struct{
    buffer []byte
    udpconn *net.UDPConn
}
//唯一标识
func (r InfluxUDPLoopGoroutine) Name() string{
    return "rpc.InfluxUDPLoopGoroutine"
}
//循环休眠时间
func (r InfluxUDPLoopGoroutine) SleepMillisecond() time.Duration{
    return time.Millisecond *10//10毫秒一次
}
//发送返回信息
func (r InfluxUDPLoopGoroutine)sendResponse( addr *net.UDPAddr,ok bool) bool {

    var ret influxbase.ResposeData
    if ok {
        ret.Ret = "0"
        ret.Msg = "success"
    }else{
        ret.Ret = "1"
        ret.Msg = "error"
    }
    json,_ := ret.ToJson()
    _,err := r.udpconn.WriteToUDP([]byte(json), addr)
    /*if ok {
        _,err = r.udpconn.WriteToUDP([]byte("{\"ret\"=\"0\"}"), addr)
    }else{
        _,err = r.udpconn.WriteToUDP([]byte("{\"ret\"=\"1\"}"), addr)
    }*/

    return !influxbase.HasError(err)
}
//初始化
func (r InfluxUDPLoopGoroutine) Init() bool{
    r.buffer = make([]byte, 2048)
    addr := net.UDPAddr{
        Port: 1234,
        IP: net.ParseIP("127.0.0.1"),
    }
    var err error
    r.udpconn, err = net.ListenUDP("udp", &addr)
    if influxbase.HasError(err) {
        return false
    }
    return true
}
//实际工作内容
func (r InfluxUDPLoopGoroutine) Run() bool{
    _,remoteaddr,err := r.udpconn.ReadFromUDP(r.buffer)
    var ok bool
    ok = true
    if influxbase.HasError(err) {
        ok = false
    }
    go r.sendResponse(remoteaddr,ok)
    return true
}

//UPD的客服端
type InfluxUDPClient struct{

}

func (s InfluxUDPClient) Send(data string) bool{
    defer influxbase.RecoverError()
    p :=  make([]byte, 2048)
    conn, err := net.Dial("udp", "127.0.0.1:1234")
    defer conn.Close()
    if influxbase.HasError(err) {
        return false
    }
    fmt.Fprintf(conn, data)
    _, err = bufio.NewReader(conn).Read(p)
    if influxbase.HasError(err) {
        return false
    }

    return true
}
