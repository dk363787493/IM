package svc

import (
	"api/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	MysqlDB *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.MysqlConfig.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:  c,
		MysqlDB: db,
	}
}
