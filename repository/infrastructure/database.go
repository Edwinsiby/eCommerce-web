package infrastructre

import (
	"fmt"
	"zog/delivery/models"
	"zog/domain/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() (*gorm.DB, error) {
	dsn := "host=localhost user=edwin dbname=edwin password=acid port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	DB = db
	DB.AutoMigrate(&entity.OtpKey{}, &models.Signup{}, &entity.Admin{}, &entity.User{}, &entity.Ticket{}, &entity.Apparel{}, &entity.CartItem{}, &entity.Cart{}, &entity.Wishlist{}, &entity.Order{}, &entity.OrderItem{}, &entity.Address{}, &entity.Inventory{}, &entity.Invoice{}, &entity.Return{}, &entity.Coupon{}, &entity.Offer{})
	return db, nil
}
