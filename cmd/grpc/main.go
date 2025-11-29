package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/grpcmiddlerware"
	"github.com/arthurhzna/Golang_gRPC/internal/handler"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/arthurhzna/Golang_gRPC/pb/auth"
	"github.com/arthurhzna/Golang_gRPC/pb/cart"
	"github.com/arthurhzna/Golang_gRPC/pb/newsletter"
	"github.com/arthurhzna/Golang_gRPC/pb/order"
	"github.com/arthurhzna/Golang_gRPC/pb/product"
	"github.com/arthurhzna/Golang_gRPC/pkg/database"
	"github.com/joho/godotenv"
	"github.com/xendit/xendit-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	gocache "github.com/patrickmn/go-cache"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	cacheService := gocache.New(time.Hour*24, time.Hour)
	authMiddleware := grpcmiddlerware.NewAuthMiddleware(cacheService)

	db := database.ConnectDb(ctx, os.Getenv("DB_URL"))
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(productRepository, cartRepository)
	cartHandler := handler.NewCartHandler(cartService)

	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(db, orderRepository, productRepository)
	orderHandler := handler.NewOrderHandler(orderService)

	newsletterRepository := repository.NewNewsletterRepository(db)
	newsletterService := service.NewNewsletterService(newsletterRepository)
	newsletterHandler := handler.NewNewsletterHandler(newsletterService)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcmiddlerware.ErrorMiddleware,
		authMiddleware.Middleware,
	))

	if os.Getenv("ENVIRONMENT") == "DEV" {
		reflection.Register(grpcServer)
	}

	auth.RegisterAuthServiceServer(grpcServer, authHandler)
	product.RegisterProductServiceServer(grpcServer, productHandler)
	cart.RegisterCartServiceServer(grpcServer, cartHandler)
	order.RegisterOrderServiceServer(grpcServer, orderHandler)
	newsletter.RegisterNewsletterServiceServer(grpcServer, newsletterHandler)
	grpcServer.Serve(lis)

}
