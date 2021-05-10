package config

import (
	"encoding/json"
	"io/ioutil"
)

var Config config

type config struct {
	Address         string
	Port            string
	SecretId        string
	SecretKey       string
	Template_id     string
	Sms_sdk_id      string
	Sign            string
	Available_times int32
	TokenTime       int64
}

func init() {
	data, err := ioutil.ReadFile("conf/app.json")
	if err != nil {
		panic(err)
		return
	}
	// 将json文件数据解析到struct
	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
		return
	}
}
