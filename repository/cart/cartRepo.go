package cart

import (
	"errors"
	"zog/domain/entity"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db}
}

func (cr *CartRepository) Create(userid int) (*entity.Cart, error) {
	cart := &entity.Cart{
		UserId: userid,
	}
	if err := cr.db.Create(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil

}

func (cr *CartRepository) UpdateCart(cart *entity.Cart) error {
	return cr.db.Where("user_id = ?", cart.UserId).Save(&cart).Error
}

func (cr CartRepository) GetByUserID(userid int) (*entity.Cart, error) {
	var cart entity.Cart
	result := cr.db.Where("user_id=?", userid).First(&cart)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, errors.New("cart not found")
	}
	return &cart, nil
}

func (cr CartRepository) GetCartById(userId int) (*entity.Cart, error) {
	var cart entity.Cart
	result := cr.db.Where("user_id=?", userId).First(&cart)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &cart, nil
}

func (cr *CartRepository) CreateCartItem(cartItem *entity.CartItem) error {
	if err := cr.db.Create(&cartItem).Error; err != nil {
		return err
	}
	return nil
}

func (cr CartRepository) UpdateCartItem(cartItem *entity.CartItem) error {
	return cr.db.Save(cartItem).Error
}

func (cr CartRepository) RemoveCartItem(cartItem *entity.CartItem) error {
	return cr.db.Where("product_name=?", cartItem.ProductName).Delete(&cartItem).Error
}

func (cr *CartRepository) GetByName(productName string, cartId int) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	result := cr.db.Where("product_name=? AND cart_id = ?", productName, cartId).First(&cartItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, errors.New("Product not exsisted")
	}
	return &cartItem, nil
}

func (cr *CartRepository) GetAllCartItems(cartId int) ([]entity.CartItem, error) {
	var cartItems []entity.CartItem
	result := cr.db.Where("cart_id=?", cartId).Find(&cartItems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return cartItems, nil
}

func (cr *CartRepository) RemoveCartItems(cartId int) error {
	var cartItems entity.CartItem
	result := cr.db.Where("cart_id=?", cartId).Delete(&cartItems)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil
}

func (cr *CartRepository) GetByType(userId int, addressType string) (*entity.Address, error) {
	var address entity.Address
	result := cr.db.Where(&entity.Address{UserId: userId, Type: addressType}).First(&address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User address not found - add address")
		}
		return nil, result.Error
	}
	return &address, nil
}
func (c *CartRepository) AddTicketToWishlist(ticket *entity.Wishlist) error {
	if err := c.db.Create(ticket).Error; err != nil {
		return err
	}
	return nil
}

func (c *CartRepository) GetTicketFromWishlist(category string, id, userId int) (bool, error) {
	var ticket entity.Wishlist
	result := c.db.Where(&entity.Wishlist{UserId: userId, Category: category, ProductId: id}).First(&ticket)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.New("Error finding ticket")
	}
	return true, nil
}

func (c *CartRepository) AddApparelToWishlist(apparel *entity.Wishlist) error {
	if err := c.db.Create(apparel).Error; err != nil {
		return err
	}
	return nil
}

func (c *CartRepository) GetApparelFromWishlist(category string, id int, userId int) (bool, error) {
	var apparel entity.Wishlist
	result := c.db.Where(&entity.Wishlist{UserId: userId, Category: category, ProductId: id}).First(&apparel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.New("Error finding apparel")
	}
	return true, nil
}

func (c *CartRepository) GetWishlist(userId int) (*[]entity.Wishlist, error) {
	var wishlist []entity.Wishlist
	result := c.db.Where("user_id=?", userId).Find(&wishlist)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &wishlist, nil
}

func (c *CartRepository) RemoveFromWishlist(category string, id, userId int) error {
	product := entity.Wishlist{
		ProductId: id,
		UserId:    userId,
		Category:  category,
	}
	return c.db.Where("user_id=?", userId).Delete(&product).Error
}
