package order

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/manmanxing/go_center_common/beacon/etcd"
	"github.com/manmanxing/go_center_common/db/mysql"
	"time"
)

const (
	orderDBEtcdAddress = ""
)

var (
	__MasterDB *sql.DB
)

func MasterDB() *sql.DB {
	return __MasterDB
}

func MasterDBClose() error {
	if __MasterDB != nil {
		return __MasterDB.Close()
	}
	return nil
}

//mysql db 基本信息配置
func init() {
	value, err := etcd.GetValue(orderDBEtcdAddress)
	if err != nil {
		panic(err)
	}

	config, err := mysql.Decode([]byte(value))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			e := MasterDBClose()
			if e != nil {
				fmt.Println("master db close err:", e)
			}
			panic(err)
		}
	}()

	__MasterDB, err = sql.Open("mysql", config.Master.DSN)
	if err != nil {
		panic(err)
	}
	__MasterDB.SetMaxOpenConns(config.Master.MaxOpen)
	__MasterDB.SetMaxIdleConns(config.Master.MaxIdle)
	__MasterDB.SetConnMaxLifetime(time.Hour * 1)
	if err = __MasterDB.Ping(); err != nil {
		panic(err)
	}
}
