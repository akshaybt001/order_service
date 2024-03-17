package adapter

import (
	"fmt"

	"github.com/akshaybt001/order_service/entitties"
	helperstruct "github.com/akshaybt001/order_service/helper_struct"
	"gorm.io/gorm"
)

type OrderAdapter struct {
	DB *gorm.DB
}



func NewOrderAdapter(db *gorm.DB) *OrderAdapter {
	return &OrderAdapter{
		DB: db,
	}
}

func (order *OrderAdapter) OrderAll(items []helperstruct.OrderAll, userId uint) (int, error) {
	var orderId int
	tx := order.DB.Begin()

	quary := "INSERT INTO orders (user_id,payment_type_id,order_status_id,address_id,total) VALUES ($1,$2,$3,$4,0) RETURNING id"
	if err := tx.Raw(quary, userId, 1, 1, 1).Scan(&orderId).Error; err != nil {
		tx.Rollback()
		return -1, err
	}
	if orderId == 0 {
		return -1, fmt.Errorf("order not found")
	}
	for _, item := range items {
		queryItemsInsert := "INSERT INTO order_items (product_id,quantity,total,order_id) VALUES($1,$2,$3,$4)"
		if err := tx.Exec(queryItemsInsert, item.ProductId, item.Quantity, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return -1, err
		}
		queryUpdateTotal := "UPDATE orders SET total = total + $1 WHERE id = $2"
		if err := tx.Exec(queryUpdateTotal, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return -1, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return -1, fmt.Errorf("error while transaction")
	}
	return orderId, nil
}

func (order *OrderAdapter) CancelOrder(orderId uint) error {
	tx := order.DB.Begin()
	queryDelete := "DELETE FROM order_items WHERE order_id = ?"
	if err := tx.Exec(queryDelete, orderId).Error; err != nil {
		tx.Rollback()
		return err
	}
	deleteOrder := "UPDATE orders SET order_status_id = $1 WHERE id = $2"
	if err := tx.Exec(deleteOrder, 5, orderId).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (order *OrderAdapter) ChangeOrderStatus(orderId int, orderStatusId int) error {
	queryUpdate := "UPDATE orders SET order_status_id = $1 WHERE id = $2"
	if err := order.DB.Exec(queryUpdate, orderStatusId, orderId).Error; err != nil {
		return err
	}
	return nil
}

func (order *OrderAdapter) GetAllOrdersUser(userId int) ([]helperstruct.GetAllOrder, error) {
	tx := order.DB.Begin()
	var orders []entitties.Order
	var res []helperstruct.GetAllOrder
	queryGetAll := "SELECT * FROM orders WHERE user_id = ?"
	if err := tx.Raw(queryGetAll, userId).Scan(&orders).Error; err != nil {
		return []helperstruct.GetAllOrder{}, err
	}
	for _, order := range orders {
		var orderItems []entitties.OrderItems
		queryGetAllItems := "SELECT * FROM order_items WHERE order_id = ?"
		if err := tx.Raw(queryGetAllItems, order.Id).Scan(&orderItems).Error; err != nil {
			return []helperstruct.GetAllOrder{}, err
		}
		response := helperstruct.GetAllOrder{
			OrderId:       uint(order.Id),
			AddressId:     uint(order.AddressId),
			PaymentTypeId: uint(order.PaymentTypeId),
			OrderStatusId: uint(order.OrderStatusId),
			OrderItems:    orderItems,
		}
		res = append(res, response)
	}
	if err := tx.Commit().Error; err != nil {
		return []helperstruct.GetAllOrder{}, err
	}
	return res, nil
}

func (order *OrderAdapter) GetAllOrders() ([]helperstruct.GetAllOrder, error) {
	var res []helperstruct.GetAllOrder
	tx := order.DB.Begin()
	var orders []entitties.Order

	orderQuery := "SELECT * FROM orders"
	if err := tx.Raw(orderQuery).Scan(&orders).Error; err != nil {
		tx.Rollback()
		return []helperstruct.GetAllOrder{}, err
	}
	for _, order := range orders {
		var orderItems []entitties.OrderItems
		queryOrderItems := "SELECT * FROM order_items WHERE order_id = ?"
		if err := tx.Raw(queryOrderItems, order.Id).Scan(&orderItems).Error; err != nil {
			tx.Rollback()
			return []helperstruct.GetAllOrder{}, err
		}
		response := helperstruct.GetAllOrder{
			OrderId:       uint(order.Id),
			AddressId:     uint(order.AddressId),
			PaymentTypeId: uint(order.PaymentTypeId),
			OrderStatusId: uint(order.OrderStatusId),
			OrderItems:    orderItems,
		}
		res = append(res, response)
	}
	if err := tx.Commit().Error; err != nil {
		return []helperstruct.GetAllOrder{}, err
	}
	return res, nil
}


func (order *OrderAdapter) GetOrder(orderId int) (helperstruct.GetAllOrder, error) {
	tx := order.DB.Begin()
	var orderData entitties.Order
	queryOrder := "SELECT * FROM orders WHERE id = ?"
	if err := tx.Raw(queryOrder, orderId).Scan(&orderData).Error; err != nil {
		tx.Rollback()
		return helperstruct.GetAllOrder{}, err
	}
	var orderItems []entitties.OrderItems
	queryOrderItem := "SELECT * FROM order_items WHERE order_id = ?"
	if err := tx.Raw(queryOrderItem, orderId).Scan(&orderItems).Error; err != nil {
		return helperstruct.GetAllOrder{}, err
	}
	res := helperstruct.GetAllOrder{
		OrderId:       uint(orderData.Id),
		AddressId:     uint(orderData.AddressId),
		PaymentTypeId: uint(orderData.PaymentTypeId),
		OrderStatusId: uint(orderData.OrderStatusId),
		OrderItems:    orderItems,
	}
	return res, nil
}