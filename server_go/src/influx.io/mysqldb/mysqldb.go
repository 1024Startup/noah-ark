package mysqldb

import (
    _"github.com/go-sql-driver/mysql"
    "database/sql"
    "influx.io/influxbase"
    _"errors"
    "fmt"
    "log"
)
//统一连接配置结构，如有增加参数，使用冗余方式加入
type DBConfig struct{
    Host string//"192.168.10.10:27017"mongodb  //"127.0.0.1:6379"redis
}


type IOrmMysql interface{
    ConnectDb(db influxbase.DBConfig) error

}
type MysqlDB struct{
    *sql.DB

}
//连接数据库地址
func (p * MysqlDB)ConnectDb(db influxbase.DBConfig) error{
    var err error
    p.DB, err = sql.Open("mysql", db.Host)//user:password@/dbname //user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname?timeout=90s&collation=utf8mb4_unicode_ci
    if influxbase.HasError(err) {
        return err
    }
    return nil
}


//查询数据
func  query() {
    db, err := sql.Open("mysql", "root:@/shopvisit")
     influxbase.ThorwError(err)

    rows, err := db.Query("SELECT * FROM shopvisit.announcement")
     influxbase.ThorwError(err)

    for rows.Next() {
        columns, _ := rows.Columns()

        scanArgs := make([]interface{}, len(columns))
        values := make([]interface{}, len(columns))

        for i := range values {
            scanArgs[i] = &values[i]
        }

        //将数据保存到 record 字典
        err = rows.Scan(scanArgs...)
        record := make(map[string]string)
        for i, col := range values {
            if col != nil {
                record[columns[i]] = string(col.([]byte))
            }
        }
        fmt.Println(record)
    }
    rows.Close()

}
func query2()  {
    fmt.Println("Query2")
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/shopvisit?charset=utf8")
     influxbase.ThorwError(err)

    rows, err := db.Query("SELECT id,imgUrl,createDate,state FROM announcement")
     influxbase.ThorwError(err)

    for rows.Next(){
        var id int
        var state int
        var imgUrl string
        var createDate string
        //注意这里的Scan括号中的参数顺序，和 SELECT 的字段顺序要保持一致。
        if err := rows.Scan(&id,&imgUrl,&createDate,&state); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s id is %d on %s with state %d\n", imgUrl, id, createDate, state)
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    rows.Close()
}


//插入数据
func insert()  {
    db, err := sql.Open("mysql", "root:@/shopvisit")
     influxbase.ThorwError(err)

    stmt, err := db.Prepare(`INSERT announcement (imgUrl, detailUrl, createDate, state) VALUES (?, ?, ?, ?)`)
     influxbase.ThorwError(err)

    res, err := stmt.Exec("/visitshop/img/ann/cofox1.png",nil,"2017-09-06",0)
     influxbase.ThorwError(err)

    id, err := res.LastInsertId()
     influxbase.ThorwError(err)

    fmt.Println(id)
    stmt.Close()

}

//修改数据
func update() {
    db, err := sql.Open("mysql", "root:@/shopvisit")
     influxbase.ThorwError(err)

    stmt, err := db.Prepare("UPDATE announcement set imgUrl=?, detailUrl=?, createDate=?, state=? WHERE id=?")
     influxbase.ThorwError(err)

    res, err := stmt.Exec("/visitshop/img/ann/cofox2.png", nil, "2017-09-05", 1, 7)
     influxbase.ThorwError(err)

    num, err := res.RowsAffected()
     influxbase.ThorwError(err)

    fmt.Println(num)
    stmt.Close()
}

//删除数据
func remove() {
    db, err := sql.Open("mysql", "root:@/shopvisit")
     influxbase.ThorwError(err)

    stmt, err := db.Prepare("DELETE FROM announcement WHERE id=?")
     influxbase.ThorwError(err)

    res, err := stmt.Exec(7)
     influxbase.ThorwError(err)

    num, err := res.RowsAffected()
     influxbase.ThorwError(err)

    fmt.Println(num)
    stmt.Close()

}
