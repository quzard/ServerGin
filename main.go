package main

import (
	"ServerGin/pkg/config"
	"ServerGin/pkg/log"
	"ServerGin/pkg/orm"
	"ServerGin/pkg/sms"
	"ServerGin/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	sms.NewSNS(config.Config.SecretId, config.Config.SecretKey)

	l, err := log.LoggerToFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	r.Use(l)
	routers.Init(r)
	if err := orm.NewDB("database.db"); err != nil {
		log.Log.WithFields(logrus.Fields{
			"name": "数据库打开失败",
		}).Error(err, "Error")
		return
	}
	err = r.Run(fmt.Sprintf("%s:%s", config.Config.Address, config.Config.Port))
	if err != nil {
	}
}
