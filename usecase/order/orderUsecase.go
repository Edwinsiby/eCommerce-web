package order

import (
	"errors"
	"time"
	"zog/domain/entity"
	"zog/domain/utils"
	cartrepository "zog/repository/cart"
	repository "zog/repository/order"
	productrepository "zog/repository/product"
	userrepository "zog/repository/user"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

type OrderUsecase struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *cartrepository.CartRepository
	userRepo    *userrepository.UserRepository
	productRepo *productrepository.ProductRepository
}

func NewOrder(orderRepo *repository.OrderRepository, cartRepo *cartrepository.CartRepository, userRepo *userrepository.UserRepository, productRepo *productrepository.ProductRepository) *OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo, cartRepo: cartRepo, userRepo: userRepo, productRepo: productRepo}
}

func (ou *OrderUsecase) ExecutePurchaseCod(userId int, address int) (*entity.Invoice, error) {
	var orderItems []entity.OrderItem
	cart, err := ou.cartRepo.GetCartById(userId)
	if err != nil {
		return nil, errors.New("Cart  not found")
	}
	cartItems, err1 := ou.cartRepo.GetAllCartItems(int(cart.ID))
	if err1 != nil {
		return nil, errors.New("Cart Items  not found")
	}
	userAddress, err := ou.userRepo.GetAddressById(address)
	if err != nil {
		return nil, errors.New("User address  not found")
	}
	Total := cart.TotalPrice - float64(cart.OfferPrice)
	order := &entity.Order{
		UserID:        cart.UserId,
		AddressId:     userAddress.ID,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "Cod",
		PaymentStatus: "pending",
	}

	OrderID, err2 := ou.orderRepo.Create(order)
	if err2 != nil {
		return nil, errors.New("Order placing failed")
	}
	invoiceData := &entity.Invoice{
		OrderId:     OrderID,
		UserId:      userId,
		AddressType: userAddress.Type,
		Quantity:    cart.TicketQuantity + cart.ApparelQuantity,
		Price:       order.Total,
		Payment:     order.PaymentMethod,
		Status:      order.PaymentStatus,
		PaymentId:   "nil",
		Remark:      "Zog_Festiv",
	}
	invoice, err := ou.orderRepo.CreateInvoice(invoiceData)
	if err != nil {
		return nil, errors.New("Invoice Creating failed")
	}
	for _, cartItem := range cartItems {
		orderItem := entity.OrderItem{
			OrderID:   OrderID,
			ProductID: cartItem.ProductId,
			Category:  cartItem.Category,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Price,
		}
		orderItems = append(orderItems, orderItem)
		inventory := entity.Inventory{
			ProductId:       cartItem.ProductId,
			ProductCategory: cartItem.Category,
			Quantity:        cartItem.Quantity,
		}
		err = ou.productRepo.DecreaseProductQuantity(&inventory)
	}

	err = ou.orderRepo.CreateOrderItems(orderItems)
	if err != nil {
		return nil, errors.New("User cart is empty")
	}

	err = ou.cartRepo.RemoveCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("Delete cart items failed")
	}
	cart.ApparelQuantity = 0
	cart.TicketQuantity = 0
	cart.TotalPrice = 0
	cart.OfferPrice = 0
	err = ou.cartRepo.UpdateCart(cart)
	if err != nil {
		return nil, errors.New("Updating cart failed")
	}
	return invoice, nil
}

func (ou *OrderUsecase) ExecutePurchasePaypal(userId int, address int) (string, error) {

	token, err := utils.GeneratePayPalAccessToken()
	if err != nil {
		return "", errors.New("Access token not created")
	}
	return token, nil

}

func (ou *OrderUsecase) ExecutePurchaseRazorPay(userId int, address int, c *gin.Context) (string, int, error) {
	var orderItems []entity.OrderItem
	cart, err := ou.cartRepo.GetCartById(userId)
	if err != nil {
		return "", 0, errors.New("Cart  not found")
	}
	cartItems, err1 := ou.cartRepo.GetAllCartItems(int(cart.ID))
	if err1 != nil {
		return "", 0, errors.New("Cart Items  not found")
	}
	userAddress, err := ou.userRepo.GetAddressById(address)
	if err != nil {
		return "", 0, errors.New("User address  not found")
	}
	client := razorpay.NewClient("rzp_test_O6q2DXJHecJBHI", "MU9PWzkhTBSCkPnxEUOAZdYW")

	data := map[string]interface{}{
		"amount":   int(cart.TotalPrice) * 100,
		"currency": "INR",
		"receipt":  "101",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", 0, errors.New("Payment not initiated")
	}
	razorId, _ := body["id"].(string)
	Total := cart.TotalPrice - float64(cart.OfferPrice)
	order := &entity.Order{
		UserID:        cart.UserId,
		AddressId:     userAddress.ID,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "razorpay",
		PaymentStatus: "pending",
		PaymentId:     razorId,
	}
	OrderId, err2 := ou.orderRepo.Create(order)
	if err2 != nil {
		return "", 0, errors.New("Order placing failed")
	}
	for _, cartItem := range cartItems {
		orderItem := entity.OrderItem{
			OrderID:   OrderId,
			ProductID: cartItem.ProductId,
			Category:  cartItem.Category,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	err3 := ou.orderRepo.CreateOrderItems(orderItems)
	if err3 != nil {
		return "", 0, errors.New("User cart is empty")
	}
	return razorId, OrderId, nil
}

func (ou *OrderUsecase) ExecuteRazorPaymentVerification(Signature, razorId, paymentId string) (*entity.Invoice, error) {

	result, err := ou.orderRepo.GetByRazorId(razorId)
	if err != nil {
		return nil, errors.New("Order not found")
	}
	err1 := utils.RazorPaymentVerification(Signature, razorId, paymentId)
	if err1 != nil {
		result.PaymentStatus = "failed"
		err2 := ou.orderRepo.Update(result)
		if err2 != nil {
			return nil, errors.New("payment updation failed")
		}
		return nil, err1
	}
	result.PaymentStatus = "successful"
	err3 := ou.orderRepo.Update(result)
	if err3 != nil {
		return nil, errors.New("payment updation failed")
	}
	userCart, err := ou.cartRepo.GetByUserID(result.UserID)
	if err != nil {
		return nil, errors.New("User cart not found")
	}
	userAddress, err := ou.userRepo.GetAddressById(result.AddressId)
	if err != nil {
		return nil, errors.New("User address  not found")
	}
	Total := userCart.TotalPrice - float64(userCart.OfferPrice)
	invoiceData := &entity.Invoice{
		OrderId:     result.ID,
		UserId:      result.UserID,
		AddressType: userAddress.Type,
		Quantity:    userCart.TicketQuantity + userCart.ApparelQuantity,
		Price:       Total,
		Payment:     "razorpay",
		Status:      "succesful",
		PaymentId:   "nil",
		Remark:      "Zog_Festiv",
	}
	invoice, err := ou.orderRepo.CreateInvoice(invoiceData)
	if err != nil {
		return nil, errors.New("Invoice Creating failed")
	}
	err4 := ou.cartRepo.RemoveCartItems(int(userCart.ID))
	if err4 != nil {
		return nil, errors.New("Delete cart items failed")
	}
	userCart = &entity.Cart{
		OfferPrice:      0,
		TotalPrice:      0,
		TicketQuantity:  0,
		ApparelQuantity: 0,
	}
	err5 := ou.cartRepo.UpdateCart(userCart)
	if err5 != nil {
		return nil, errors.New("Updating cart failed")
	}
	return invoice, nil
}

func (ou *OrderUsecase) ExecutePurchaseWallet(userId int, address int) (*entity.Invoice, error) {
	var orderItems []entity.OrderItem
	user, err := ou.userRepo.GetByID(userId)
	if err != nil {
		return nil, errors.New("User not found")
	}
	cart, err := ou.cartRepo.GetCartById(userId)
	if err != nil {
		return nil, errors.New("Cart  not found")
	}
	if user.Wallet < int(cart.TotalPrice) {
		return nil, errors.New("Wallet have not enough money-add money or choose another method")
	}
	cartItems, err1 := ou.cartRepo.GetAllCartItems(int(cart.ID))
	if err1 != nil {
		return nil, errors.New("Cart Items  not found")
	}
	userAddress, err := ou.userRepo.GetAddressById(address)
	if err != nil {
		return nil, errors.New("User address  not found")
	}
	Total := cart.TotalPrice - float64(cart.OfferPrice)
	order := &entity.Order{
		UserID:        cart.UserId,
		AddressId:     userAddress.ID,
		Total:         Total,
		Status:        "pending",
		PaymentMethod: "wallet",
		PaymentStatus: "succesful",
	}

	OrderID, err2 := ou.orderRepo.Create(order)
	if err2 != nil {
		return nil, errors.New("Order placing failed")
	}
	user.Wallet -= int(order.Total)
	err = ou.orderRepo.UpdateUserWallet(user)
	if err != nil {
		return nil, errors.New("Wallet updation failed")
	}
	invoiceData := &entity.Invoice{
		OrderId:     OrderID,
		UserId:      userId,
		AddressType: userAddress.Type,
		Quantity:    cart.TicketQuantity + cart.ApparelQuantity,
		Price:       order.Total,
		Payment:     order.PaymentMethod,
		Status:      order.PaymentStatus,
		PaymentId:   "nil",
		Remark:      "Zog_Festiv",
	}
	invoice, err := ou.orderRepo.CreateInvoice(invoiceData)
	if err != nil {
		return nil, errors.New("Invoice Creating failed")
	}
	for _, cartItem := range cartItems {
		orderItem := entity.OrderItem{
			OrderID:   OrderID,
			ProductID: cartItem.ProductId,
			Category:  cartItem.Category,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	err = ou.orderRepo.CreateOrderItems(orderItems)
	if err != nil {
		return nil, errors.New("User cart is empty")
	}

	err = ou.cartRepo.RemoveCartItems(int(cart.ID))
	if err != nil {
		return nil, errors.New("Delete cart items failed")
	}
	cart = &entity.Cart{
		TotalPrice:      0,
		TicketQuantity:  0,
		ApparelQuantity: 0,
	}
	err = ou.cartRepo.UpdateCart(cart)
	if err != nil {
		return nil, errors.New("Updating cart failed")
	}
	return invoice, nil
}

func (ou *OrderUsecase) ExecuteCancelOrder(orderId int) error {
	result, err := ou.orderRepo.GetByID(orderId)
	if err != nil {
		return errors.New("Order not found")
	}
	user, err := ou.userRepo.GetByID(result.UserID)
	if err != nil {
		return errors.New("User not found")
	}
	if result.Status != "pending" && result.Status != "confirmed" {
		return errors.New("order cancelation failed- cancel time exceeded")
	}
	if result.PaymentStatus == "successful" {
		result.PaymentStatus = "refund"
		user.Wallet = int(result.Total)
		err = ou.orderRepo.UpdateUserWallet(user)
		if err != nil {
			return errors.New("User wallet updation failed")
		}
	}
	result.Status = "canceled"
	err1 := ou.orderRepo.Update(result)
	if err1 != nil {
		return errors.New("order cancelation failed")
	}
	return nil
}

func (ou *OrderUsecase) ExecuteOrderHistory(userId, page, limit int) ([]entity.Order, error) {
	offset := (page - 1) * limit
	orderList, err := ou.orderRepo.GetAllOrders(userId, offset, limit)
	if err != nil {
		return nil, err
	}
	return orderList, nil
}

func (ou *OrderUsecase) ExecuteReturnOrder(returnData entity.Return) error {
	order, err := ou.orderRepo.GetByID(returnData.OrderId)
	if err != nil {
		return errors.New("Order not found")
	}
	order.Status = "return"
	order.PaymentStatus = "refund"
	err = ou.orderRepo.Update(order)
	if err != nil {
		return errors.New("order updation failed")
	}
	err = ou.orderRepo.CreateReturn(&returnData)
	if err != nil {
		return errors.New("return creation failed")
	}
	return nil
}

func (ou *OrderUsecase) ExecuteReturnUpdate(status string, returnId int) error {
	result, err := ou.orderRepo.GetReturnByID(returnId)
	if err != nil {
		return errors.New("Order not found")
	}
	order, err := ou.orderRepo.GetByID(result.OrderId)
	if err != nil {
		return errors.New("Order not found")
	}
	result.Status = status
	result.Refund = "wallet"
	result.TotalPrice = int(order.Total)
	err = ou.orderRepo.UpdateReturn(result)
	if err != nil {
		return errors.New("return updation failed")
	}
	order.Status = "return"
	err = ou.orderRepo.Update(order)
	if err != nil {
		return errors.New("order updation failed")
	}
	return nil
}

func (ou *OrderUsecase) ExecuteRefund(orderId int) error {
	order, err := ou.orderRepo.GetByID(orderId)
	if err != nil {
		return errors.New("Order not found")
	}

	if order.Status == "return" {
		result, err := ou.orderRepo.GetReturnByOrderID(order.ID)
		if err != nil {
			return errors.New("Return not found")
		}
		user, err := ou.userRepo.GetByID(result.UserId)
		if err != nil {
			return errors.New("Order not found")
		}
		if result.Refund == "wallet" {
			user.Wallet = result.TotalPrice
			err = ou.orderRepo.UpdateUserWallet(user)
			if err != nil {
				return errors.New("User wallet updation failed")
			}
		}
		result.Status = "completed"
		err = ou.orderRepo.UpdateReturn(result)
		if err != nil {
			return errors.New("return updation failed")
		} else {
			return nil
		}
	} else {
		user, err := ou.userRepo.GetByID(order.UserID)
		if err != nil {
			return errors.New("Order not found")
		}
		if order.PaymentStatus == "successful" {
			order.PaymentStatus = "refund"
			user.Wallet = int(order.Total)
			err = ou.orderRepo.UpdateUserWallet(user)
			if err != nil {
				return errors.New("User wallet updation failed")
			}
		}
		order.Status = "canceled"
		err1 := ou.orderRepo.Update(order)
		if err1 != nil {
			return errors.New("order cancelation failed")
		}
		return nil
	}

}

func (ou *OrderUsecase) ExecuteOrderUpdate(orderId int, status string) error {
	result, err := ou.orderRepo.GetByID(orderId)
	if err != nil {
		return errors.New("Order not found")
	}
	result.Status = status
	err = ou.orderRepo.Update(result)
	if err != nil {
		return errors.New("order updation failed")
	}
	return nil
}

func (ou *OrderUsecase) ExecuteSalesReportByDate(startDate, endDate time.Time) (*entity.SalesReport, error) {
	orders, err := ou.orderRepo.GetByDate(startDate, endDate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}

func (ou *OrderUsecase) ExecuteSalesReportByPeriod(period string) (*entity.SalesReport, error) {
	startDate, endDate := utils.CalculatePeriodDates(period)

	orders, err := ou.orderRepo.GetByDate(startDate, endDate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}

func (ou *OrderUsecase) ExecuteSalesReportByCategory(category, period string) (*entity.SalesReport, error) {
	startDate, endDate := utils.CalculatePeriodDates(period)
	orders, err := ou.orderRepo.GetByCategory(category, startDate, endDate)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}

func (o *OrderUsecase) ExecuteSortedOrders(page, limit int, status string) ([]entity.Order, error) {
	offset := (page - 1) * limit
	orders, err := o.orderRepo.GetByStatus(offset, limit, status)
	if err != nil {
		return nil, errors.New("report fetching failed")
	}
	return orders, nil
}
