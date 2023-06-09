package handlers

import (
	"net/http"
	"strconv"
	middlewares "zog/delivery/middlewares"
	_ "zog/docs"
	"zog/domain/entity"
	usecase "zog/usecase/admin"
	product "zog/usecase/product"

	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
)

type AdminHandler struct {
	AdminUsecase   *usecase.AdminUsecase
	ProductUsecase *product.ProductUsecase
}

func NewAdminHandler(AdminUsecase *usecase.AdminUsecase, ProductUsecase *product.ProductUsecase) *AdminHandler {
	return &AdminHandler{AdminUsecase, ProductUsecase}
}

// Admin Register  godoc
//
//	@Summary		registering new admin
//	@Description	Adding new admin to the database
//	@Tags			Admin Authentication
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.Admin	true	"Admin Data"
//	@Success		200		{object}	entity.Admin
//	@Router			/registeradmin [post]
func (ac *AdminHandler) RegisterAdmin(c *gin.Context) {
	var admin entity.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser, err := ac.AdminUsecase.ExecuteAdminCreate(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// Admin Login With Password godoc
//
//	@Summary		Admin Login with password
//	@Description	Admin login with password and phone number
//	@Tags			Admin Authentication
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.Login	true	"Admin Data"
//	@Success		200		{object}	entity.Login
//	@Router			/adminloginpassword [post]
func (uh *AdminHandler) AdminLoginWithPassword(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	phone, _ := payload["phone"].(string)
	password, _ := payload["password"].(string)

	adminId, err := uh.AdminUsecase.ExecuteLoginWithPassword(phone, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		middlewares.CreateJwtCookie(adminId, phone, "admin", c)
		c.JSON(http.StatusOK, gin.H{"massage": "admin loged in succesfully and cookie stored"})
	}

}

// Admin Login  godoc
//
//	@Summary		Admin login with otp
//	@Description	Admin login with otp
//	@Tags			Admin Authentication
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.Login	true	"Admin Data"
//	@Success		200		{object}	entity.Login
//	@Router			/adminlogin [post]
func (al *AdminHandler) Login(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	phone, _ := payload["phone"].(string)
	err := al.AdminUsecase.ExecuteAdminLogin(phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Otp send succesfully to": phone})
	}
}

// Admin Otp Validation  godoc
//
//	@Summary		Otp validation
//	@Description	Otp Validation for admin login
//	@Tags			Admin Authentication
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.Login	true	"Admin Data"
//	@Success		200		{object}	entity.Login
//	@Router			/adminotpvalidation [post]
func (ah *AdminHandler) LoginOtpValidation(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	phone, _ := payload["phone"].(string)
	otp, _ := payload["otp"].(string)
	resend, _ := payload["resend"].(string)
	if resend == "resend" {
		err := ah.AdminUsecase.ExecuteAdminLogin(phone)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"massage": "otp resend successful"})
		}
	} else {
		admin, err1 := ah.AdminUsecase.ExecuteOtpValidation(phone, otp)
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err1.Error()})
			return
		}
		middlewares.CreateJwtCookie(int(admin.ID), admin.Phone, "admin", c)
		c.JSON(http.StatusOK, gin.H{"massage": "admin loged in succesfully and cookie stored"})

	}

}

// Admin Home  godoc
//
//	@Summary		Admin dashbord
//	@Description	Admin dashbord
//	@Tags			Admin Authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.AdminDashboard
//	@Router			/adminhome [get]
func (ah *AdminHandler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"options": "Sales Report - User Management - Product Management - Order Management"})
	dashboardResponse, err := ah.AdminUsecase.ExecuteAdminDashboard()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"Dashboard": dashboardResponse})
}

// User Management  godoc
//
//	@Summary		User list
//	@Description	Showing user list for management by admin
//	@Tags			Admin User Management
//	@Accept			json
//	@Produce		json
//	@Param			page	query		string	false	"page no"
//	@Param			limit	query		string	false	"limit no"
//	@Success		200		{object}	entity.User
//	@Router			/usermanagement [get]
func (ul *AdminHandler) UserList(c *gin.Context) {
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
	userList, err1 := ul.AdminUsecase.ExecuteUserList(page, limit)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User list not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Users": userList})
}

// Sort User  godoc
//
//	@Summary		User list by permission
//	@Description	Showing user list for management by sorting with permisiion
//	@Tags			Admin User Management
//	@Accept			json
//	@Produce		json
//	@Param			page		query		string	false	"page no"
//	@Param			limit		query		string	false	"limit no"
//	@Param			permission	query		string	true	"true/false"
//	@Success		200			{object}	entity.User
//	@Router			/sortuser [get]
func (a *AdminHandler) SortUserByPermission(c *gin.Context) {
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
	permission := c.Query("permission")
	userList, err1 := a.AdminUsecase.ExecuteUserListByPermission(page, limit, permission)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User list not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Users": userList})
}

// Search User  godoc
//
//	@Summary		Search user by id or name
//	@Description	Showing user list for management by searching with name or id
//	@Tags			Admin User Management
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"User Name"
//	@Param			id		query		string	false	"User Id"
//	@Success		200		{object}	entity.User
//	@Router			/searchuser [get]
func (a *AdminHandler) SearchUser(c *gin.Context) {
	name := c.DefaultQuery("name", " ")
	userIdStr := c.DefaultQuery("id", "0")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	if userId != 0 {
		userList, err1 := a.AdminUsecase.ExecuteSearchUserById(userId)
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User list not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Users": userList})
	} else if name != " " {
		userList, err1 := a.AdminUsecase.ExecuteSearchUserByName(name)
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User list not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Users": userList})
	} else {
		c.JSON(http.StatusNotAcceptable, gin.H{"Error": "invalid entry"})
	}
}

// Toggle User Permission  godoc
//
//	@Summary		block/unblock user
//	@Description	Toggling user permission for block/unblock
//	@Tags			Admin User Management
//	@Accept			json
//	@Produce		json
//	@param			id	path		string	true	"User ID"
//	@Success		200	{object}	entity.Admin
//	@Router			/userpermission/{id} [post]
func (tp *AdminHandler) TogglePermission(c *gin.Context) {
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	err1 := tp.AdminUsecase.ExecuteTogglePermission(Id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	}
	c.JSON(http.StatusOK, gin.H{"success": "user permission toggled"})
}

// Add Ticket  godoc
//
//	@Summary		Adding new product
//	@Description	Adding new product of category ticket in database
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Produce		json
//	@param			ticket	body		entity.TicketInput	true	"Ticket Data"
//	@Success		200		{object}	entity.Ticket
//	@Router			/addticket [post]
func (ct *AdminHandler) CreateTicket(c *gin.Context) {
	adminID, _ := c.Get("userID")
	adminId := adminID.(int)
	var input entity.TicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticket := entity.Ticket{
		Name:        input.Name,
		Price:       input.Price,
		Date:        input.Date,
		Location:    input.Location,
		ImageURL:    input.ImageURL,
		Category:    "ticket",
		SubCategory: input.Category,
		AdminId:     adminId,
	}
	ticketId, err := ct.ProductUsecase.ExecuteCreateTicket(ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ticketDetails := entity.TicketDetails{
			TicketId:    ticketId,
			Description: input.Description,
			Venue:       input.Venue,
		}
		err := ct.ProductUsecase.ExecuteCreateTicketDetails(ticketDetails)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		inventory := entity.Inventory{
			ProductId:       ticketId,
			ProductCategory: "ticket",
			Quantity:        input.Quantity,
		}
		err = ct.ProductUsecase.ExecuteCreateInventory(inventory)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": "ticket added succesfully"})
}

// Edit Ticket  godoc
//
//	@Summary		Edit existing product data
//	@Description	Edit data of a product in category ticket
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.Ticket	true	"Ticket Data"
//	@Success		200		{object}	entity.Ticket
//	@Router			/editticket/{id} [put]
func (et *AdminHandler) EditTicket(c *gin.Context) {
	var ticket entity.Ticket
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err1 := et.ProductUsecase.ExecuteEditTicket(ticket, Id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "ticket edited succesfully"})
}

// Delete Ticket  godoc
//
//	@Summary		Delete existing product from database
//	@Description	Soft deleting the data of a product from database in category ticket
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Produce		json
//	@param			ProductId	query		int	true	"product id"
//	@Success		200			{object}	entity.Ticket
//	@Router			/deleteticket/{id} [delete]
func (dt *AdminHandler) DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	err1 := dt.ProductUsecase.ExecuteDeleteTicket(Id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ticket not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "ticket deleted succesfully"})
	}

}

// Add Apparel  godoc
//
//	@Summary		Adding new product
//	@Description	Adding new product of category apparel in database
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.ApparelInput	true	"Apparel Data"
//	@Success		200		{object}	entity.Apparel
//	@Router			/addapparel [post]
func (aa *AdminHandler) CreateApparel(c *gin.Context) {
	var input entity.ApparelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	apparelId, err := aa.ProductUsecase.ExecuteCreateApparel(input.Apparel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = aa.ProductUsecase.ExecuteCreateApparelDetails(input.ApparelDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "Apparel added succesfully"})
	input.Inventory.ProductId = apparelId
	err = aa.ProductUsecase.ExecuteCreateInventory(input.Inventory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "apparel inventory created succesfully"})
	}
}

// Edit Apparel  godoc
//
//	@Summary		Edit existing product data
//	@Description	Edit data of a product in category apparel
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		entity.Apparel	true	"Apparel Data"
//	@Success		200		{object}	entity.Apparel
//	@Router			/editapparel/{id} [put]
func (ea *AdminHandler) EditApparel(c *gin.Context) {
	var apparel entity.Apparel
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	if err := c.ShouldBindJSON(&apparel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err1 := ea.ProductUsecase.ExecuteEditApparel(apparel, Id)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "apparel edited succesfully"})
}

// Delete Apparel  godoc
//
//	@Summary		Delete existing product from database
//	@Description	Soft deleting the data of a product from database in category apparel
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Produce		json
//	@param			ProductId	query		int	true	"product id"
//	@Success		200			{object}	entity.Apparel
//	@Router			/deleteapparel/{id} [delete]
func (da *AdminHandler) DeleteApparel(c *gin.Context) {
	id := c.Param("id")
	Id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "str conversion failed"})
	}
	err1 := da.ProductUsecase.ExecuteDeleteApparel(Id)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Apparel not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "apparel deleted succesfully"})
	}
}

// Add Coupon   godoc
//
//	@Summary		Adding coupon by admin
//	@Description	Addig coupon for users, with a unique code
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Param			admin	body		entity.Coupon	true	"coupon"
//	@Success		200		{string}	string			"Success masage"
//	@Router			/addcoupon  [post]
func (ah *AdminHandler) AddCoupon(c *gin.Context) {

	var coupon entity.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := ah.ProductUsecase.ExecuteAddCoupon(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "coupon created succesfully"})
	}
}

// Add Offer   godoc
//
//	@Summary		Adding offer by admin
//	@Description	Addig coupon for users, with a unique code
//	@Tags			Admin Product&Offer Management
//	@Accept			json
//	@Param			admin	body		entity.Offer	true	"offer"
//	@Success		200		{string}	string			"Success masage"
//	@Router			/addoffer  [post]
func (ah *AdminHandler) AddOffer(c *gin.Context) {
	var offer entity.Offer
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := ah.ProductUsecase.ExecuteAddOffer(&offer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "offer created succesfully"})
	}
}
