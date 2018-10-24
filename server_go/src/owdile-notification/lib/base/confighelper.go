package base
import (
    _"log"
    "github.com/c4pt0r/ini"
    "os"
    "path/filepath"
    "fmt"
    _"strings"
    "influx.io/influxbase"
)
//配置表
type DB_ENV_CONFIG_TYPE int
const (
    CONFIG_MONGO DB_ENV_CONFIG_TYPE =  iota
    CONFIG_REDIS
    CONFIG_MYSQL
    //////
	DEV_ENV_CONFIG_MONGO
	DEV_ENV_CONFIG_REDIS
	DEV_ENV_CONFIG_MYSQL
    PRODUCT_ENV_CONFIG_MONGO
    PRODUCT_ENV_CONFIG_REDIS
    PRODUCT_ENV_CONFIG_MYSQL
)


func GetSlaveMongoDbConfigs() []influxbase.DBConfig {
    var dbc  influxbase.DBConfig
    NodeEnv := os.Getenv("NODE_ENV")
    if(NodeEnv == "dev"){
        dbc.Host = "127.0.0.1:27017"
    }
    dbcs  := make([]influxbase.DBConfig,1);
    dbcs[0] = dbc
    return dbcs
}
//获取到配置
func GetDbConfigByType(configtype DB_ENV_CONFIG_TYPE) influxbase.DBConfig {
    var dbc  influxbase.DBConfig
    NodeEnv := os.Getenv("NODE_ENV")
    if(CONFIG_MONGO == configtype && NodeEnv == "dev")||
    configtype == DEV_ENV_CONFIG_MONGO{
        dbc.Host = "192.168.10.10:27017"
        return dbc
    }


    if (CONFIG_REDIS == configtype && NodeEnv == "dev")||
    configtype == DEV_ENV_CONFIG_REDIS{
        dbc.Host = "192.168.10.10:6379"
        return dbc
    }
    return dbc

}


func getCurrentDirectory() string {
    /*dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
    }
    return strings.Replace(dir, "\\", "/", -1)
*/
    execpath,_ := os.Executable()
    return filepath.Join(filepath.Dir(execpath), "./dbconfig.conf")
}
func getConfig(env string,field string, deflautvalue string) *string {
    configfile := getCurrentDirectory()
    fmt.Println(configfile)
    conf := ini.NewConf(configfile)//("owdile-notification\\conf\\app.conf")
    fmt.Println(conf)
    conf.Parse()
    return conf.String(env, field, deflautvalue)

}



func GetDbConfig(field string,deflautvalue string) * string{
    NodeEnv := os.Getenv("NODE_ENV")

    fmt.Println(NodeEnv)
    return getConfig(NodeEnv,field,deflautvalue)
}
