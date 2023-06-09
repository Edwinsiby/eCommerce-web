package infrastructre

import (
	"fmt"
	"log"
	"os"
	"zog/delivery/models"
	"zog/domain/entity"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dsn string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn = os.Getenv("KEY5")

}
func ConnectToDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	DB = db
	DB.AutoMigrate(&entity.TicketDetails{}, &entity.OtpKey{}, &models.Signup{}, &entity.Admin{}, &entity.User{}, &entity.Ticket{}, &entity.Apparel{}, &entity.CartItem{}, &entity.Cart{}, &entity.Wishlist{}, &entity.Order{}, &entity.OrderItem{}, &entity.Address{}, &entity.Inventory{}, &entity.Invoice{}, &entity.Return{}, &entity.Coupon{}, &entity.Offer{})
	return db, nil
}
