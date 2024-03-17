package db

import (
	"github.com/akshaybt001/order_service/entitties"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(connect string) (*gorm.DB,error){
	db,err:=gorm.Open(postgres.Open(connect),&gorm.Config{})
	if err!=nil{
		return nil,err
	}
	db.AutoMigrate(&entitties.Order{})
	db.AutoMigrate(&entitties.OrderItems{})
	db.AutoMigrate(&entitties.OrderStatus{})

	return db,nil
}