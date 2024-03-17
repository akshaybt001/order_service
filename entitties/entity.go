package entitties

type Order struct{
	Id uint `gorm:"PrimaryKey"`
	UserId uint
	PaymentTypeId uint
	AddressId uint
	OrderStatusId uint
	OrderStatus `gorm:"ForeignKey:OrderStatusid"`
	Total float64
}

type OrderItems struct{
	Id uint
	OrderId uint
	Order `gorm:"ForeignKey:OrderId"`
	ProductId uint
	Quantity uint
	Total float64
}

type OrderStatus struct{
	Id uint
	Status string
}