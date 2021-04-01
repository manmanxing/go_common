package order

import (
	"github.com/jinzhu/gorm"
	"github.com/manmanxing/go_center_common/db/mysql/order"
	"github.com/manmanxing/go_center_common/gorm/log"
	"github.com/manmanxing/go_center_common/gorm/ot"
)

//将 *sql.DB 转为 *gorm.DB

var _DB *gorm.DB

func DB() *gorm.DB {
	return _DB
}

func init() {
	db, err := gorm.Open("mysql", order.MasterDB())
	if err != nil {
		panic(err)
	}
	//使用自定义的日志打印
	db.SetLogger(log.NewLogger())
	db.LogMode(true)
	ot.InitCallback(db, "order_service")
	_DB = db
}