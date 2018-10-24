//model公用操作
package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"reflect"
	"strings"
)

//Ormer对象
//var O orm.Ormer

//注意一下，varchar最多能存储65535个字符

//以下是数据表model对象
var (
	ModelUser = new(User) //
)

//以下是数据表
var (
	TableUser = GetTable("user")
)

//以下是表字段查询
var Fields = map[string]map[string]string{
	TableUser: StringSliceToMap(GetFields(ModelUser)),
}

//初始化数据库注册
func Init() {
	//初始化数据库
	RegisterDB()
	runmode := beego.AppConfig.String("runmode")
	if runmode == "prod" {
		orm.Debug = false
		orm.RunSyncdb("default", false, false)
	} else {
		orm.Debug = true
		orm.RunSyncdb("default", false, true)
	}

	//安装初始数据
	//install()

	//全局变量赋值
	//ModelConfig.UpdateGlobal() //配置文件全局变量更新
	//ModelSys.UpdateGlobal()    //更新系统配置的全局变量

}

//注册数据库
func RegisterDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	models := []interface{}{
		ModelUser,
	}

	orm.RegisterModelWithPrefix(beego.AppConfig.String("mysql_prefix"), models...)
	db_user := beego.AppConfig.String("mysql_user")
	db_password := beego.AppConfig.String("mysql_password")
	if envpass := os.Getenv("MYSQL_PASSWORD"); envpass != "" {
		db_password = envpass
	}
	db_database := beego.AppConfig.String("mysql_database")
	if envdatabase := os.Getenv("MYSQL_DATABASE"); envdatabase != "" {
		db_database = envdatabase
	}

	db_charset := beego.AppConfig.String("mysql_charset")
	db_host := beego.AppConfig.String("mysql_url")
	if envhost := os.Getenv("MYSQL_HOST"); envhost != "" {
		db_host = envhost
	}

	db_port := beego.AppConfig.String("mysql_port")
	if envport := os.Getenv("MYSQL_PORT"); envport != "" {
		db_port = envport
	}
	dblink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%v", db_user, db_password, db_host, db_port, db_database, db_charset, "Asia%2FShanghai")
	//下面两个参数后面要放到app.conf提供用户配置使用
	// (可选)设置最大空闲连接
	maxIdle := beego.AppConfig.DefaultInt("mysql_max_idle", 50)
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := beego.AppConfig.DefaultInt("mysql_max_conn", 300)
	if err := orm.RegisterDataBase("default", "mysql", dblink, maxIdle, maxConn); err != nil {
		panic(err)
	}

}

//获取指定Strut的字段
//@param            tableObj        Strut结构对象，引用传递
//@return           fields          返回字段数组
func GetFields(tableObj interface{}) (fields []string) {
	elem := reflect.ValueOf(tableObj).Elem()
	for i := 0; i < elem.NumField(); i++ {
		fields = append(fields, elem.Type().Field(i).Name)
	}
	return fields
}

//获取带表前缀的数据表
//@param            table               数据表
func GetTable(table string) string {
	prefix := beego.AppConfig.String("mysql_prefix")
	return prefix + strings.TrimPrefix(table, prefix)
}

//将字符串切片数组转成map
func StringSliceToMap(slice []string) (maps map[string]string) {
	maps = make(map[string]string)
	for _, v := range slice {
		maps[v] = v
	}
	return maps
}
