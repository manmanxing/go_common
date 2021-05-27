package sequence

import (
	"fmt"
	"time"

	"github.com/manmanxing/errors"
	"github.com/manmanxing/go_common/db/mysql/order"
	"github.com/manmanxing/go_common/util"
)

const (
	nums = 10000
	sql  = "select xxx.nextval from dual "
)

var orderNoSeqChan = make(chan int64, nums)

func init() {
	go getBatchOrderNoSeqLoop()
}

//批量生产orderNo seq
func getBatchOrderNoSeqLoop() {
	for {
		list, err := getOrderNoSeqByCount(nums)
		if err != nil {
			fmt.Println("get order seq err", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for i := range list {
			orderNoSeqChan <- list[i]
		}
	}
}

func getOrderNoSeqByCount(count int) (list []int64, err error) {
	list = make([]int64, 0, count)
	query := sql + fmt.Sprintf(" where count = %d", count)
	//todo 后期可以改为生成自增序列的数据库
	rows, err := order.MasterDB().Query(query)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("db query err,sql:%s", query))
		return nil, err
	}
	defer func() {
		if rows != nil {
			e := rows.Close()
			if e != nil {
				err = errors.Wrap(e, "mysql rows close err")
				return
			}
		}
	}()
	list = make([]int64, 0, count)
	var n int64
	for rows.Next() {
		if err = rows.Scan(&n); err != nil {
			err = errors.Wrap(err, "rows scan err")
			return nil, err
		}
		n %= 1000000000000
		list = append(list, n)
	}
	return list, rows.Err()
}

//生成单条orderNo seq
func getSingleTradeNoSeq() (n int64, err error) {
	err = order.MasterDB().QueryRow(sql).Scan(&n)
	if err != nil {
		err = errors.Wrap(err, "db query row scan err")
		return
	}
	n %= 1000000000000
	return
}

//从 channel 获取还是从数据库获取 orderNo seq
func getTradeNoSeq() (n int64, err error) {
	select {
	case n = <-orderNoSeqChan:
		return
	default:
		return getSingleTradeNoSeq()
	}
}

//获取 orderNo
//示例：202103110005895719897912
func NewOrderNo(userId int64) (string, error) {
	if userId <= 0 {
		return "", errors.New("userid is invalid")
	}
	seq, err := getTradeNoSeq()
	if err != nil {
		return "", err
	}

	//拼接 orderNo
	//20060102 + 12位数字 + 用户id的后四位
	return time.Now().In(util.BeijingLocation).Format(util.TimeFormatDate) +
		FormatSequenceNo(seq) +
		LowestFourBytesForUserID(userId), nil
}
