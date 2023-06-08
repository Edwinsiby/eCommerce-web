package handlers

import (
	"net/http"
	"strconv"
	middlewares "zog/delivery/middlewares"
	"zog/delivery/models"
	_ "zog/docs"
	"zog/domain/entity"
	cartusecase "zog/usecase/cart"
	productusecase "zog/usecase/product"
	usecase "zog/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	_ "gorm.io/gorm"
)

type UserHandler struct {
	UserUsecase    *usecase.UserUsecase
	ProductUsecase *productusecase.ProductUsecase
	CartUsecase    *cartusecase.CartUsecase
}

func NewUserHandler(UserUsecase *usecase.UserUsecase, ProductUsecase *productusecase.ProductUsecase, CartUsecase *cartusecase.CartUsecase) *UserHandler {
	return &UserHandler{UserUsecase, ProductUsecase, CartUsecase}
}

// UserSignup  godoc
//
//	@Summary		signup
//	@Description	Adding new user to the database
//	@Tags			User Authentication
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			userInput	body		models.Signup	true	"User Data"
//	@Success		200			{object}	entity.User
//	@Router			/signup [post]
func (uh *UserHandler) Signup(c *gin.Context) {
	var userInput models.Signup
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	copier.Copy(&user, &userInput)
	newUser, err := uh.UserUsecase.ExecuteSignup(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// UserSignup Otp  godoc
//
//	@Summary		signup with opt validation
//	@Description	Adding new user to the database
//	@Tags			User Authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.Signup	true	"User Data"
//	@Success		200		{object}	entity.User
//	@Router			/signupwithotp [post]
func (uh *UserHandler) SignupWithOtp(c *gin.Context) {
	var user models.Signup
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := uh.UserUsecase.ExecuteSignupWithOtp(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Otp send succesfuly to": user.Phone, "Key": key})
	}
}

// SignupOtpValidation  godoc
//
//	@Summary		Sign Up Otp Validation
//	@Description	Validating user otp for signup
//	@Tags			User Authentication
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			key	formData	string	true	"Twilio Key"
//	@Param			otp	formData	string	true	"Otp"
//	@Success		200	{string}	string	"Success message"
//	@Router			/signupotpvalidation [post]
func (uh *UserHandler) SignupOtpValidation(c *gin.Context) {
	key := c.PostForm("key")
	otp := c.PostForm("otp")
	err := uh.UserUsecase.ExecuteSignupOtpValidation(key, otp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"massage": "user signup succesfull"})
	}

}

// UserLogin  godoc
//
//	@Summary		Login
//	@Description	Login for user with password
//	@Tags			User Authentication
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			phone		formData	string	true	"Phone No"
//	@Param			password	formData	string	true	"Password"
//	@Success		200			{object}	entity.Login
//	@Router			/loginwithpassword [post]
func (uh *UserHandler) LoginWithPassword(c *gin.Context) {

	phone := c.PostForm("phone")
	password := c.PostForm("password")

	userId, err := uh.UserUsecase.ExecuteLoginWithPassword(phone, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		middlewares.CreateJwtCookie(userId, phone, "user", c)
		c.JSON(http.StatusOK, gin.H{"massage": "user loged in succesfully and cookie stored"})
	}

}

// UserLogin  godoc
//
//	@Summary		Login
//	@Description	Login for user with otp
//	@Tags			User Authentication
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			phone	formData	string	true	"Phone No"
//	@Success		200		{object}	entity.Login
//	@Router			/loginwithotp [post]
func (uh *UserHandler) LoginWithOtp(c *gin.Context) {

	phone := c.PostForm("phone")
	key, err := uh.UserUsecase.ExecuteLogin(phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Otp send succesfuly to": phone, "Key": key})
	}

}

// UserOtpValidation  godoc
//
//	@Summary		Otp Validation
//	@Description	Validating user otp for login validation
//	@Tags			User Authentication
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			otp		formData	string	true	"Otp"
//	@Param			key		formData	string	true	"Key"
//	@Param			phone	formData	string	false	"phone"
//	@Param			resend	formData	string	false	"resend"
//	@Success		200		{object}	entity.Login
//	@Router			/otpvalidation [post]
func (uh *UserHandler) LoginOtpValidation(c *gin.Context) {
	otp := c.PostForm("otp")
	key := c.PostForm("key")
	phone := c.PostForm("phone")
	resend := c.PostForm("resend")

	if resend == "resend" {
		key, err := uh.UserUsecase.ExecuteLogin(phone)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"massage": "otp resend successful", "Key": key})
		}
	} else {
		user, err1 := uh.UserUsecase.ExecuteOtpValidation(key, otp)
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err1.Error()})
			return
		}
		middlewares.CreateJwtCookie(user.ID, user.Phone, "user", c)
		c.JSON(http.StatusOK, gin.H{"massage": "user loged in succesfully and cookie stored"})

	}

}

// User Home  godoc
//
//	@Summary		User Home
//	@Description	User home with the next navigations
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Success message"
//	@Router			/home [get]
func (uh *UserHandler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"options": "Logout - Add_Address - Cart - Orders - Tickets - Apparels"})
}

// User Add Address  godoc
//
//	@Summary		Add Address
//	@Description	Add new address to the database with user id
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			user	body		entity.Address	true	"User Address"
//	@Success		200		{string}	string			"Success message"
//	@Router			/addaddress [post]
func (uh *UserHandler) AddAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	var address entity.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userAddress := entity.Address{
		UserId: userId,
	}
	err := copier.Copy(&userAddress, &address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = uh.UserUsecase.ExecuteAddAddress(&userAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"massage": "user address added succesfully"})
	}
}

// User Details    godoc
//
//	@Summary		User Details
//	@Description	User profile with adress and user details
//	@Tags			User Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User
//	@Router			/userdetails [get]
func (uh *UserHandler) ShowUserDetails(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	userDetails, address, err1 := uh.UserUsecase.ExecuteShowUserDetails(userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"User Details": userDetails, "address": address})
}

// Edit Profile    godoc
//
//	@Summary		Edit User Details
//	@Description	Edit User details including address
//	@Tags			User Authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.EditUser	true	"User Data"
//	@Success		200		{string}	string			"Success message"
//	@Router			/editprofile [put]
func (uh *UserHandler) EditProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	var userInput models.EditUser
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	copier.Copy(&user, &userInput)
	err1 := uh.UserUsecase.ExecuteEditProfile(user, userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"massage": "user details updated succesfully"})
}

// Change Password   godoc
//
//	@Summary		Change password
//	@Description	Option for changing password from user side
//	@Tags			User Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Success message"
//	@Router			/forgotpassword [post]
func (uh *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	key, err := uh.UserUsecase.ExecuteChangePassword(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"massage": "Otp send seccesfuly", "key": key})
	}

}

// OtpValidationPassword godoc
//
//	@Summary		Otp validation for password
//	@Description	Otp validation for changing password
//	@Tags			User Authentication
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			password	formData	string	true	"New Password"
//	@Param			otp			formData	string	true	"Otp"
//	@Success		200			{string}	string	"Success message"
//	@Router			/otpvalidationpassword [post]
func (uh *UserHandler) OtpValidationPassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	password := c.PostForm("password")
	otp := c.PostForm("otp")
	err := uh.UserUsecase.ExecuteOtpValidationPassword(password, otp, userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"massage": "password changed succesfuly"})

}

// Tickets       godoc
//
//	@Summary		Tickets List
//	@Description	Showing the available tickets in the site
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"page no"
//	@Param			limit	query		int				false	"limit no"
//	@Param			sort	query		string			false	"Sort by location"
//	@Success		200		{object}	entity.Ticket	"Ticket Data"
//	@Router			/tickets [get]
func (uh *UserHandler) Tickets(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")
	location := c.Query("sort")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	ticketList, err := uh.ProductUsecase.ExecuteTicketList(page, limit, location)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ticket list not found"})
		return
	}
	responseList := make([]entity.Ticket, len(ticketList))
	for i, ticket := range ticketList {
		responseList[i] = entity.Ticket{
			ID:          ticket.ID,
			Name:        ticket.Name,
			Price:       ticket.Price,
			Date:        ticket.Date,
			Location:    ticket.Location,
			ImageURL:    ticket.ImageURL,
			SubCategory: ticket.SubCategory,
		}
	}
	c.JSON(http.StatusOK, gin.H{"tickets": responseList})
}

// Search Tickets       godoc
//
//	@Summary		Search Result
//	@Description	Showing the available tickets as per user search
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"page no"
//	@Param			limit	query		int				false	"limit no"
//	@Param			search	query		string			false	"Search By Name"
//	@Success		200		{object}	entity.Ticket	"Ticket Data"
//	@Router			/searchticket [get]
func (u *UserHandler) SearchTicket(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")
	search := c.Query("search")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	ticketList, err := u.ProductUsecase.ExecuteTicketSearch(page, limit, search)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ticket list not found"})
		return
	}
	responseList := make([]entity.Ticket, len(ticketList))
	for i, ticket := range ticketList {
		responseList[i] = entity.Ticket{
			ID:       ticket.ID,
			Name:     ticket.Name,
			Price:    ticket.Price,
			Date:     ticket.Date,
			Location: ticket.Location,
			ImageURL: ticket.ImageURL,
		}
	}
	c.JSON(http.StatusOK, gin.H{"tickets": responseList})
}

// Ticket Details  godoc
//
//	@Summary		Details of a Ticket
//	@Description	Showing details of a single product and option to adding cart
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			ticketid	path		string					true	"Ticket ID"
//	@Success		200			{object}	entity.TicketDetails	"Ticket Details"
//	@Router			/ticketdetails/{ticketid} [get]
func (uh *UserHandler) TicketDetails(c *gin.Context) {

	id := c.Param("ticketid")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	ticket, ticketDetails, err1 := uh.ProductUsecase.ExecuteTicketDetails(Id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Ticket": ticket, "Details": ticketDetails})
}

// Apparels       godoc
//
//	@Summary		Apparel List
//	@Description	Showing the available Apparels in the site
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"page no"
//	@Param			limit	query		int				false	"limit no"
//	@Param			sort	query		string			false	"Sort by Category"
//	@Success		200		{object}	entity.Apparel	"Apparel List"
//	@Router			/apparels [get]
func (uh *UserHandler) Apparels(c *gin.Context) {
	category := c.Query("sort")
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	apparelList, err := uh.ProductUsecase.ExecuteApperalList(page, limit, category)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Apparel list not found"})
		return
	}
	responseList := make([]entity.Apparel, len(apparelList))
	for i, apparel := range apparelList {
		responseList[i] = entity.Apparel{
			ID:          apparel.ID,
			Name:        apparel.Name,
			Price:       apparel.Price,
			ImageURL:    apparel.ImageURL,
			SubCategory: apparel.SubCategory,
		}
	}
	c.JSON(http.StatusOK, gin.H{"Apperals": responseList})
}

// Search Apparels       godoc
//
//	@Summary		Search Result
//	@Description	Showing the available apparels as per user search
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"page no"
//	@Param			limit	query		int				false	"limit no"
//	@Param			search	query		string			false	"Search By Name"
//	@Success		200		{object}	entity.Apparel	"Apparel Data"
//	@Router			/searchapparel [get]
func (u *UserHandler) SearchApparels(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")
	search := c.Query("search")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	apparelList, err := u.ProductUsecase.ExecuteApparelSearch(page, limit, search)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Apparel list not found"})
		return
	}
	responseList := make([]entity.Apparel, len(apparelList))
	for i, apparel := range apparelList {
		responseList[i] = entity.Apparel{
			ID:          apparel.ID,
			Name:        apparel.Name,
			Price:       apparel.Price,
			ImageURL:    apparel.ImageURL,
			SubCategory: apparel.SubCategory,
		}
	}
	c.JSON(http.StatusOK, gin.H{"tickets": responseList})
}

// Apparel Details  godoc
//
//	@Summary		Details of a Apparel
//	@Description	Showing details of a single product and option to adding cart
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			apparelid	path		string					true	"Apparel ID"
//	@Success		200			{object}	entity.ApparelDetails	"Apparel Details"
//	@Router			/appareldetails/{apparelid} [get]
func (uh *UserHandler) ApparelDetails(c *gin.Context) {
	id := c.Param("apparelid")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	apparel, err1 := uh.ProductUsecase.ExecuteApparelDetails(Id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Apparel": apparel})
}

// Add to cart  godoc
//
//	@Summary		Add product to cart
//	@Description	Adding product with quantity to cart with product id
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			category	path		string	true	"Ticket/Apparel"
//	@Param			productid	path		string	true	"Product ID"
//	@Param			quantity	path		string	true	"Product Quantity"
//	@Success		200			{string}	string	"Success message"
//	@Router			/addtocart/{category}/{productid}/{quantity} [post]
func (uh *UserHandler) AddToCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	product := c.Param("category")
	strId := c.Param("productid")
	strQuantity := c.Param("quantity")
	Id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	quantity, err := strconv.Atoi(strQuantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	err = uh.CartUsecase.ExecuteAddToCart(product, Id, quantity, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}

// Add to cart  godoc
//
//	@Summary		Add product to wishlist
//	@Description	Adding single product to wishlist with product id
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			category	path		string	true	"Ticket/Apparel"
//	@Param			productid	path		string	true	"Product ID"
//	@Success		200			{string}	string	"Success message"
//	@Router			/addtowishlist/{category}/{productid} [post]
func (u *UserHandler) AddToWishlist(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	product := c.Param("category")
	strId := c.Param("productid")
	Id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	err = u.CartUsecase.ExecuteAddToWishlist(product, Id, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to wishlist successfully"})
}

// Increase quantity
//
//	@Summary		Increase quantity of existing product in cart
//	@Description	Increasing quantity one by one in the cart
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			category	path		string	true	"Product Category"
//	@Param			productid	path		string	true	"Product Id"
//	@Success		200			{string}	string	"Added"
//	@Router			/increasequantity/{category}/{productid} [put]
func (uh *UserHandler) IncreaseQuantity(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	product := c.Param("category")
	strId := c.Param("productid")
	Id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	err1 := uh.CartUsecase.ExecuteAddToCart(product, Id, 1, userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "Quantity increased"})
	}

}

// Cart     godoc
//
//	@Summary		User Cart
//	@Description	Showing user cart
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Cart	"User Cart"
//	@Router			/usercart [get]
func (uh *UserHandler) Cart(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	var userCartResponse entity.Cart
	userCart, err1 := uh.CartUsecase.ExecuteCart(userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	copier.Copy(&userCartResponse, &userCart)
	c.JSON(http.StatusOK, gin.H{"User Cart": userCartResponse})
}

// Cart List    godoc
//
//	@Summary		Cart List
//	@Description	Showing the products in user cart
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.CartItem	"Cart List"
//	@Router			/usercartlist [get]
func (uh *UserHandler) CartList(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	cartItems, err1 := uh.CartUsecase.ExecuteCartList(userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Cart List": cartItems})
}

// RemoveFromCart godoc
//
//	@Summary		Remove Product from cart
//	@Description	Removing product from the cart for unique and decrese quantity for existing product
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			product	path		string	true	"Product Name"
//	@Param			id		path		int		true	"Product ID"
//	@Success		200		{string}	string	"Success message"
//	@Router			/removefromcart/{product}/{id} [delete]
func (uh *UserHandler) RemoveFromCart(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	id := c.Param("id")
	product := c.Param("product")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
		return
	}
	err1 := uh.CartUsecase.ExecuteRemoveFromCart(product, Id, userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart successfully"})
}

// Wishlist    godoc
//
//	@Summary		Wish List
//	@Description	Showing the products in user wishlist
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Wishlist	"Wishlist"
//	@Router			/userwishlist [get]
func (u *UserHandler) ViewWishlist(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	wishList, err1 := u.CartUsecase.ExecuteViewWishlist(userId)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Wish List": wishList})
}

// Available Coupon  godoc
//
//	@Summary		checking coupon availability
//	@Description	showing the available coupons and eligibility
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.Coupon	"Available coupons"
//	@Router			/coupons [get]
func (u *UserHandler) AvailableCoupons(c *gin.Context) {

	couponList, err := u.ProductUsecase.ExecuteAvailableCoupons()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Available coupons are ": couponList})
}

// Apply Coupon  godoc
//
//	@Summary		checking coupon availability and adding offer amount
//	@Description	applying coupon offer for user cart amount
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"coupon code"
//	@Success		200		{string}	string	"total amount"
//	@Router			/applycoupon/{code} [post]
func (u *UserHandler) ApplyCoupon(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	code := c.Param("code")

	totalOffer, err := u.CartUsecase.ExecuteApplyCoupon(userId, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Offer price for coupon is ": totalOffer, "Offer": "Applied succesfuly"})
	}
}

// Offer Check godoc
//
//	@Summary		checking offer availability
//	@Description	finding and showing offer for user with respect to user cart
//	@Tags			User Shopping
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"offers"
//	@Router			/offer [get]
func (u *UserHandler) OfferCheck(c *gin.Context) {
	userID, _ := c.Get("userID")
	userId := userID.(int)
	offerList, err := u.CartUsecase.ExecuteOfferCheck(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Available offers are ": offerList})
	}
}

// LogOut     godoc
//
//	@Summary		logout
//	@Description	Deleting cookie from the browser while logout
//	@Tags			User Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Success message"
//	@Router			/logout [post]
func (uh *UserHandler) Logout(c *gin.Context) {
	err := middlewares.DeleteCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user cookie deletion failed"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
	}
}
