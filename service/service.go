package service

import (
	"context"
	"fmt"
	"io"

	"github.com/akshaybt001/order_service/adapter"
	helperstruct "github.com/akshaybt001/order_service/helper_struct"
	"github.com/akshaybt001/proto_files/pb"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	Tracer     opentracing.Tracer
	CartClient pb.CartServiceClient
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type OrderService struct {
	Adapter adapter.AdapterInterface
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(adapter adapter.AdapterInterface) *OrderService {
	return &OrderService{
		Adapter: adapter,
	}
}

func (order *OrderService) OrderAll(ctx context.Context, req *pb.UserId) (*pb.OrderId, error) {
	span := Tracer.StartSpan("orderall grpc")
	defer span.Finish()

	cartItems, err := CartClient.GetAllCart(context.TODO(), &pb.CartCreate{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get items from cart")
	}
	var cart []helperstruct.OrderAll
	for {
		items, err := cartItems.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		item := helperstruct.OrderAll{
			ProductId: uint(items.ProductId),
			Quantity:  float64(items.Quantity),
			Total:     uint(items.Total),
		}
		cart = append(cart, item)
	}
	if len(cart) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	if _, err := CartClient.TruncateCart(context.TODO(), &pb.CartCreate{
		UserId: req.UserId,
	}); err != nil {
		return nil, err
	}
	orderId, err := order.Adapter.OrderAll(cart, uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.OrderId{OrderId: uint32(orderId)}, nil
}

func (order *OrderService) CancelOrder(ctx context.Context, req *pb.OrderId) (*pb.OrderId, error) {
	err := order.Adapter.CancelOrder(uint(req.OrderId))
	if err != nil {
		return nil, err
	}
	return &pb.OrderId{OrderId: req.OrderId}, nil
}

func (order *OrderService) ChangeOrderStatus(ctx context.Context, req *pb.ChangeStatusRequest) (*pb.OrderId, error) {
	err := order.Adapter.ChangeOrderStatus(int(req.OrderId), int(req.StatusId))
	if err != nil {
		return nil, err
	}
	return &pb.OrderId{OrderId: req.OrderId}, nil
}

func (order *OrderService) GetAllOrdersUser(req *pb.UserId, srv pb.OrderService_GetAllOrdersUserServer) error {
	orders, err := order.Adapter.GetAllOrdersUser(int(req.UserId))
	if err != nil {
		return err
	}
	for _, ordr := range orders {
		var orderItems []*pb.OrderItems
		for _, ordrItem := range ordr.OrderItems {
			itm := &pb.OrderItems{
				OrderId:  uint32(ordrItem.OrderId),
				Id:       uint32(ordrItem.ProductId),
				Quantity: int32(ordrItem.Quantity),
				Price:    ordrItem.Total,
			}
			orderItems = append(orderItems, itm)
		}
		res := &pb.GetAllOrdersResponse{
			OrderId:       uint32(ordr.OrderId),
			AddressId:     uint32(ordr.AddressId),
			PaymentTypeId: uint32(ordr.PaymentTypeId),
			OrderStatusId: uint32(ordr.OrderStatusId),
			OrderItems:    orderItems,
		}
		fmt.Println(ordr.OrderStatusId)
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (order *OrderService) GetAllOrders(empty *pb.NoArg,srv pb.OrderService_GetAllOrdersServer)error{
	orders,err:=order.Adapter.GetAllOrders()
	if err!=nil{
		return err
	}
	for _,ordr:=range orders{
		var orderItems []*pb.OrderItems
		for _,ordrItem:=range ordr.OrderItems{
			itm:=&pb.OrderItems{
				OrderId: uint32(ordrItem.OrderId),
				Id: uint32(ordrItem.ProductId),
				Quantity: int32(ordrItem.Quantity),
				Price: ordrItem.Total,
			}
			orderItems=append(orderItems, itm)
		}
		res:=&pb.GetAllOrdersResponse{
			OrderId: uint32(ordr.OrderId),
			AddressId: uint32(ordr.AddressId),
			PaymentTypeId: uint32(ordr.PaymentTypeId),
			OrderStatusId: uint32(ordr.OrderStatusId),
			OrderItems: orderItems,
		}
		if err:=srv.Send(res);err!=nil{
			return err
		}
	}
	return nil
}

func (order *OrderService) GetOrder(ctx context.Context, req *pb.OrderId) (*pb.GetAllOrdersResponse, error) {
	orderData, err := order.Adapter.GetOrder(int(req.OrderId))
	if err != nil {
		return &pb.GetAllOrdersResponse{}, err
	}
	var orderItems []*pb.OrderItems
	for _, item := range orderData.OrderItems {
		ordrItem := &pb.OrderItems{
			Id:       uint32(item.ProductId),
			OrderId:  uint32(item.OrderId),
			Quantity: int32(item.Quantity),
			Price:    item.Total,
		}
		orderItems = append(orderItems, ordrItem)
	}
	res := &pb.GetAllOrdersResponse{
		OrderId:       req.OrderId,
		AddressId:     uint32(orderData.AddressId),
		PaymentTypeId: uint32(orderData.PaymentTypeId),
		OrderStatusId: uint32(orderData.OrderStatusId),
		OrderItems:    orderItems,
	}
	return res, nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "watching is not supported")
}
