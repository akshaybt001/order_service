package adapter

import helperstruct "github.com/akshaybt001/order_service/helper_struct"

type AdapterInterface interface {
	OrderAll(items []helperstruct.OrderAll,userId uint) (int,error)
	CancelOrder(orderId uint)error
	ChangeOrderStatus(orderId, orderStatusId int)error
	GetAllOrdersUser(userId int)([]helperstruct.GetAllOrder,error)
	GetAllOrders()([]helperstruct.GetAllOrder,error)
	GetOrder(orderId int) (helperstruct.GetAllOrder, error)
}