package helperstruct

import "github.com/akshaybt001/order_service/entitties"

type OrderAll struct {
	ProductId uint
	Quantity  float64
	Total     uint
}

type GetAllOrder struct {
	OrderId       uint
	AddressId     uint
	PaymentTypeId uint
	OrderStatusId uint
	OrderItems    []entitties.OrderItems
}