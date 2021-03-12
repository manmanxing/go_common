package mysql

import (
	"encoding/json"
	"github.com/manmanxing/errors"
)

//定义数据库配置的json格式
type Config struct {
	Master struct {
		DSN     string `json:"dsn"`
		MaxOpen int    `json:"max_open"`
		MaxIdle int    `json:"max_idle"`
	} `json:"master"`
	Slave struct {
		DSN     string `json:"dsn"`
		MaxOpen int    `json:"max_open"`
		MaxIdle int    `json:"max_idle"`
	} `json:"slave"`
}

func Decode(data []byte) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		err = errors.Wrap(err,"db config unmarshal err")
		return nil, err
	}
	return &cfg, nil
}