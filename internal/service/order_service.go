package service

import (
	"context"
	"database/sql"
	"fmt"
	operatingsystem "os"
	"runtime/debug"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
	jwtentity "github.com/arthurhzna/Golang_gRPC/internal/entity/jwt"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/order"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
}

type orderService struct {
	db                *sql.DB
	orderRepository   repository.IOrderRepository
	productRepository repository.IProductRepository
}

func NewOrderService(db *sql.DB, orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		db:                db,
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}

func (os *orderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := os.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := recover(); e != nil {
			if tx != nil {
				tx.Rollback()
			}
			debug.PrintStack()
			panic(e)
		}
	}()

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	orderRepository := os.orderRepository.WithTransaction(tx)
	productRepository := os.productRepository.WithTransaction(tx)

	numbering, err := orderRepository.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(req.Products))
	for i, product := range req.Products {
		productIds[i] = product.Id
	}

	products, err := productRepository.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*entity.Product)
	for i := range products {
		productMap[products[i].Id] = products[i]
	}

	var total float64 = 0
	for _, p := range req.Products {
		if productMap[p.Id] == nil {
			return &order.CreateOrderResponse{
				Base: utils.NotFoundResponse(fmt.Sprintf("Product %s not found", p.Id)),
			}, nil
		}
		total += productMap[p.Id].Price * float64(p.Quantity)
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)
	orderEntity := entity.Order{
		Id:              uuid.NewString(),
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
		UserId:          claims.Subject,
		OrderStatusCode: entity.OrderStatusCodeUnpaid,
		UserFullName:    claims.FullName,
		Address:         req.Address,
		PhoneNumber:     req.PhoneNumber,
		Notes:           &req.Notes,
		Total:           total,
		ExpiredAt:       &expiredAt,
		CreatedAt:       now,
		CreatedBy:       claims.FullName,
	}
	invoiceItems := make([]xendit.InvoiceItem, 0)
	for _, p := range req.Products {
		prod := productMap[p.Id]
		if prod == nil {
			invoiceItems = append(invoiceItems, xendit.InvoiceItem{
				Name:     prod.Name,
				Price:    prod.Price,
				Quantity: int(p.Quantity),
			})

		}
	}

	XenditInvoice, XenditErr := invoice.CreateWithContext(ctx, &invoice.CreateParams{
		ExternalID: orderEntity.Id,
		Amount:     total,
		Customer: xendit.InvoiceCustomer{
			GivenNames: claims.FullName,
		},
		Currency:           "IDR",
		SuccessRedirectURL: fmt.Sprintf("%s/checkout/%s/success", operatingsystem.Getenv("FE_BASE_URL"), orderEntity.Id),
		Items:              invoiceItems,
	})

	if XenditErr != nil {
		err = XenditErr
		return nil, err
	}

	orderEntity.XenditInvoiceId = &XenditInvoice.ID
	orderEntity.XenditInvoiceUrl = &XenditInvoice.InvoiceURL

	err = orderRepository.CreateOrder(ctx, &orderEntity)
	if err != nil {
		return nil, err
	}

	for _, p := range req.Products {
		var orderItem = entity.OrderItem{
			Id:                   uuid.NewString(),
			ProductId:            p.Id,
			ProductName:          productMap[p.Id].Name,
			ProductImageFileName: productMap[p.Id].ImageFileName,
			ProductPrice:         productMap[p.Id].Price,
			Quantity:             p.Quantity,
			OrderId:              orderEntity.Id,
			CreatedAt:            now,
			CreatedBy:            claims.FullName,
		}
		err = orderRepository.CreateOrderItem(ctx, &orderItem)
		if err != nil {
			return nil, err
		}
	}

	numbering.Number++

	err = orderRepository.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Base: utils.SuccessResponse("Order created successfully"),
		Id:   orderEntity.Id,
	}, nil
}
