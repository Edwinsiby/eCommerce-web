package main

import (
	"fmt"
	"log"
	"net/http"
	"zog/delivery/handlers"
	"zog/delivery/routes"
	_ "zog/docs"
	adminrepository "zog/repository/admin"
	cartrepository "zog/repository/cart"
	infrastructure "zog/repository/infrastructure"
	orderrepository "zog/repository/order"
	productrepository "zog/repository/product"
	repository "zog/repository/user"
	adminusecase "zog/usecase/admin"
	cartusecase "zog/usecase/cart"
	orderusecase "zog/usecase/order"
	productusecase "zog/usecase/product"
	usecase "zog/usecase/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Zog_festiv eCommerce API
//	@version		1.0
//	@description	API for ecommerce website

//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						token

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		www.zogfestiv.store
//	@BasePath	/

//	@schemes	http
func main() {
	db, err := infrastructure.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	adminRepo := adminrepository.NewAdminRepository(db)
	productRepo := productrepository.NewProductRepository(db)
	cartRepo := cartrepository.NewCartRepository(db)
	orderRepo := orderrepository.NewOrderRepository(db)

	userUsecase := usecase.NewUser(userRepo)
	adminUsecase := adminusecase.NewAdmin(adminRepo)
	productUsecase := productusecase.NewProduct(productRepo)
	cartUsecase := cartusecase.NewCart(cartRepo, productRepo)
	orderUsecase := orderusecase.NewOrder(orderRepo, cartRepo, userRepo, productRepo)

	userHandler := handlers.NewUserHandler(userUsecase, productUsecase, cartUsecase)
	adminHandler := handlers.NewAdminHandler(adminUsecase, productUsecase)
	orderHandler := handlers.NewOrderHandler(orderUsecase)

	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.UserRouter(router, userHandler)
	routes.AdminRouter(router, adminHandler)
	routes.OrderRouter(router, orderHandler)

	fmt.Println("Starting server on port 8080...")
	err1 := http.ListenAndServe(":8080", router)
	if err1 != nil {
		log.Fatal(err1)
	}
}
