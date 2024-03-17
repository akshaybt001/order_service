package initializer

import (
	"github.com/akshaybt001/order_service/adapter"
	"github.com/akshaybt001/order_service/service"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) *service.OrderService{
	adapter:=adapter.NewOrderAdapter(db)
	service:=service.NewOrderService(adapter)

	return service
}