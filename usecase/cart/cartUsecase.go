package cart

import (
	"errors"
	"zog/domain/entity"
	repository "zog/repository/cart"
	productrepository "zog/repository/product"
)

type CartUsecase struct {
	cartRepo    *repository.CartRepository
	productRepo *productrepository.ProductRepository
}

func NewCart(cartRepo *repository.CartRepository, productRepo *productrepository.ProductRepository) *CartUsecase {
	return &CartUsecase{cartRepo: cartRepo, productRepo: productRepo}
}

func (cu *CartUsecase) ExecuteAddToCart(product string, id int, quantity int, userid int) error {
	var userCart *entity.Cart
	var cartId int
	userCart, err := cu.cartRepo.GetByUserID(userid)

	if err != nil {
		cartid, err1 := cu.cartRepo.Create(userid)
		if err1 != nil {
			return errors.New("Failed to create user cart")
		}
		cartId = cartid
	} else {
		cartId = int(userCart.ID)
	}

	if product == "ticket" {
		ticket, err := cu.productRepo.GetTicketByID(id)
		if err != nil {
			return errors.New("Ticket not found")
		}

		cartItem := &entity.CartItem{
			CartId:      cartId,
			ProductId:   int(ticket.ID),
			Category:    "ticket",
			Quantity:    quantity,
			ProductName: ticket.Name,
			Price:       float64(ticket.Price),
		}
		existingTicket, err := cu.cartRepo.GetByName(ticket.Name, cartId)
		if err != nil {
			err = cu.cartRepo.CreateCartItem(cartItem)
			if err != nil {
				return errors.New("Adding new ticket to cart item failed")
			}
		} else {
			existingTicket.Quantity += quantity
			err := cu.cartRepo.UpdateCartItem(existingTicket)
			if err != nil {
				return errors.New("error updating existing cartitem")
			}
		}
		if userCart.TotalPrice == 0 {
			userCart.TotalPrice = cartItem.Price * float64(quantity)
			userCart.TicketQuantity = quantity
		} else {
			userCart.TotalPrice += cartItem.Price * float64(quantity)
			userCart.TicketQuantity += quantity
		}

	} else if product == "apparel" {
		apparel, err := cu.productRepo.GetApparelByID(id)
		if err != nil {
			return errors.New("Apparel not found")
		}
		cartItem := &entity.CartItem{
			CartId:      cartId,
			ProductId:   int(apparel.ID),
			Category:    "apparel",
			Quantity:    quantity,
			ProductName: apparel.Name,
			Price:       float64(apparel.Price),
		}
		existingApparel, err := cu.cartRepo.GetByName(apparel.Name, cartId)
		if err != nil {
			err = cu.cartRepo.CreateCartItem(cartItem)
			if err != nil {
				return errors.New("Adding new ticket to cart item failed")
			}
		} else {
			existingApparel.Quantity += quantity
			err := cu.cartRepo.UpdateCartItem(existingApparel)
			if err != nil {
				return errors.New("error updating existing cartitem")
			}
		}
		if userCart.TotalPrice == 0 {
			userCart.TotalPrice = cartItem.Price * float64(quantity)
			userCart.ApparelQuantity = quantity
		} else {
			userCart.TotalPrice += cartItem.Price * float64(quantity)
			userCart.ApparelQuantity += quantity
		}
	}
	err1 := cu.cartRepo.UpdateCart(userCart)
	if err1 != nil {
		return errors.New("Cart price updation failed")
	}

	return nil
}

func (cu *CartUsecase) ExecuteCart(userId int) (*entity.Cart, error) {
	userCart, err := cu.cartRepo.GetByUserID(userId)
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	} else {
		return userCart, nil
	}

}

func (cu *CartUsecase) ExecuteCartList(userId int) ([]entity.CartItem, error) {
	userCart, err := cu.cartRepo.GetByUserID(userId)
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	}
	cartItems, err := cu.cartRepo.GetAllCartItems(int(userCart.ID))
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (cu *CartUsecase) ExecuteRemoveFromCart(product string, id int, userId int) error {
	userCart, err := cu.cartRepo.GetByUserID(userId)
	if err != nil {
		return errors.New("Failed to find user cart")
	}

	if product == "ticket" {
		ticket, err := cu.productRepo.GetTicketByID(id)
		if err != nil {
			return errors.New("Ticket not found")
		}
		existingTicket, err1 := cu.cartRepo.GetByName(ticket.Name, int(userCart.ID))
		if err1 != nil {
			return errors.New("Removing ticket from cart failed")

		}
		if existingTicket.Quantity == 1 {
			err := cu.cartRepo.RemoveCartItem(existingTicket)
			if err != nil {
				return errors.New("Removing apparel from cart failed")
			}
		} else {
			existingTicket.Quantity -= 1
			err := cu.cartRepo.UpdateCartItem(existingTicket)
			if err != nil {
				return errors.New("error updating existing cartitem")
			}
		}
		userCart.TotalPrice -= float64(ticket.Price)
		userCart.TicketQuantity -= 1
	} else if product == "apparel" {
		apparel, err := cu.productRepo.GetApparelByID(id)
		if err != nil {
			return errors.New("Apparel not found")
		}
		existingApparel, err1 := cu.cartRepo.GetByName(apparel.Name, int(userCart.ID))
		if err1 != nil {
			return errors.New("Removing apparel from cart failed")
		}
		if existingApparel.Quantity == 1 {
			err := cu.cartRepo.RemoveCartItem(existingApparel)
			if err != nil {
				return errors.New("Removin apparel from cart failed")
			}
		} else {
			existingApparel.Quantity -= 1
			err := cu.cartRepo.UpdateCartItem(existingApparel)
			if err != nil {
				return errors.New("error updating existing cartitem")
			}
		}
		userCart.TotalPrice -= float64(apparel.Price)
		userCart.ApparelQuantity -= 1

	}
	if userCart.OfferPrice > 0 {
		userCart.OfferPrice = 0
	}
	err1 := cu.cartRepo.UpdateCart(userCart)
	if err1 != nil {
		return errors.New("Remove from cart failed")
	}

	return nil
}

func (c *CartUsecase) ExecuteAddToWishlist(catergory string, productId int, userId int) error {
	if catergory == "ticket" {
		ticket, err := c.productRepo.GetTicketByID(productId)
		if err != nil {
			return errors.New("Ticket not found")
		}
		exsisting, err := c.cartRepo.GetTicketFromWishlist(ticket.Category, ticket.ID, userId)
		if err != nil {
			return errors.New("Error finding exsisting product")
		}
		if exsisting == true {
			return errors.New("Product already exsist in wishlist")
		} else {
			wishTicket := &entity.Wishlist{
				UserId:      userId,
				Category:    ticket.Category,
				ProductId:   ticket.ID,
				ProductName: ticket.Name,
				Price:       float64(ticket.Price),
			}
			err = c.cartRepo.AddTicketToWishlist(wishTicket)
			if err != nil {
				return errors.New("Product adding to wishlist failed")
			}
		}
	} else {
		apparel, err := c.productRepo.GetApparelByID(productId)
		if err != nil {
			return errors.New("Apparel not found")
		}
		exsisting, err := c.cartRepo.GetApparelFromWishlist(apparel.Category, apparel.ID)
		if err != nil {
			return errors.New("Error finding exsisting product")
		}
		if exsisting == true {
			return errors.New("Product already exsist in wishlist")
		} else {
			wishApparel := &entity.Wishlist{
				UserId:      userId,
				Category:    apparel.Category,
				ProductId:   apparel.ID,
				ProductName: apparel.Name,
				Price:       float64(apparel.Price),
			}
			err = c.cartRepo.AddApparelToWishlist(wishApparel)
			if err != nil {
				return errors.New("Product adding to wishlist failed")
			}
		}
	}
	return nil
}

func (c *CartUsecase) ExecuteViewWishlist(userId int) (*[]entity.Wishlist, error) {
	wishlist, err := c.cartRepo.GetWishlist(userId)
	if err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (c *CartUsecase) ExecuteApplyCoupon(userId int, code string) (int, error) {
	var totalOffer, totalPrice int
	userCart, err := c.cartRepo.GetByUserID(userId)
	if err != nil {
		return 0, errors.New("Failed to find user cart")
	}
	coupon, err := c.productRepo.GetCouponByCode(code)
	if err != nil {
		return 0, errors.New("Sorry coupon not found")
	}
	cartItems, err := c.cartRepo.GetAllCartItems(int(userCart.ID))
	if err != nil {
		return 0, errors.New("User Cart Items not found")
	}
	for _, cartItem := range cartItems {
		if cartItem.Category == coupon.Category {
			totalPrice += int(cartItem.Price) * cartItem.Quantity
		}
	}
	if totalPrice > 0 {
		if coupon.Type == "percentage" {
			totalOffer = totalPrice / coupon.Amount
		} else {
			totalOffer = coupon.Amount
		}
	} else {
		return 0, errors.New("Add more product from different category")
	}
	if userCart.OfferPrice != 0 {
		return 0, errors.New("User Cart offer already applied")
	} else {
		userCart.OfferPrice = totalOffer
		err = c.cartRepo.UpdateCart(userCart)
		if err != nil {
			return 0, errors.New("User Cart updation failed")
		}
	}

	return totalOffer, nil

}

func (u *CartUsecase) ExecuteOfferCheck(userId int) (*[]entity.Offer, error) {
	userCart, err := u.cartRepo.GetByUserID(userId)
	if err != nil {
		return nil, errors.New("Failed to find user cart")
	}
	offer, err := u.productRepo.GetOfferByPrice(int(userCart.TotalPrice))
	if err != nil {
		return nil, errors.New("No valid offers, Add few more products worth of 500")
	}
	return offer, nil
}
