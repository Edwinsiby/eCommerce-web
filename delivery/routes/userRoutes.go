package routes

import (
	"zog/delivery/handlers"
	m "zog/delivery/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine, userHandler *handlers.UserHandler) *gin.Engine {

	r.POST("/signup", userHandler.Signup)
	r.POST("/signupwithotp", userHandler.SignupWithOtp)
	r.POST("/signupotpvalidation", userHandler.SignupOtpValidation)
	r.POST("/loginwithotp", userHandler.LoginWithOtp)
	r.POST("/otpvalidation", userHandler.LoginOtpValidation)
	r.POST("/loginwithpassword", userHandler.LoginWithPassword)

	r.GET("/home", m.UserRetriveCookie, userHandler.Home)
	r.POST("/addaddress", m.UserRetriveCookie, userHandler.AddAddress)
	r.GET("/userdetails", m.UserRetriveCookie, userHandler.ShowUserDetails)
	r.PUT("/editprofile", m.UserRetriveCookie, userHandler.EditProfile)
	r.POST("/forgotpassword", m.UserRetriveCookie, userHandler.ChangePassword)
	r.POST("/otpvalidationpassword", m.UserRetriveCookie, userHandler.OtpValidationPassword)
	r.GET("/tickets", m.UserRetriveCookie, userHandler.Tickets)
	r.GET("/searchticket", m.UserRetriveCookie, userHandler.SearchTicket)
	r.GET("/ticketdetails/:ticketid", m.UserRetriveCookie, userHandler.TicketDetails)
	r.GET("/apparels", m.UserRetriveCookie, userHandler.Apparels)
	r.GET("/searchapparel", m.UserRetriveCookie, userHandler.SearchApparels)
	r.GET("/apparelsdetails/:apparelid", m.UserRetriveCookie, userHandler.ApparelDetails)

	r.POST("/addtocart/:category/:productid/:quantity", m.UserRetriveCookie, userHandler.AddToCart)
	r.POST("/addtowishlist/:category/:productid", m.UserRetriveCookie, userHandler.AddToWishlist)
	r.PUT("/increasequantity/:category/:productid", m.UserRetriveCookie, userHandler.IncreaseQuantity)
	r.GET("/usercartlist", m.UserRetriveCookie, userHandler.CartList)
	r.GET("/usercart", m.UserRetriveCookie, userHandler.Cart)
	r.DELETE("/removefromcart/:product/:id", m.UserRetriveCookie, userHandler.RemoveFromCart)
	r.GET("/userwishlist", m.UserRetriveCookie, userHandler.ViewWishlist)
	r.GET("/coupons", m.UserRetriveCookie, userHandler.AvailableCoupons)
	r.POST("/applycoupon/:code", m.UserRetriveCookie, userHandler.ApplyCoupon)
	r.GET("/offer", m.UserRetriveCookie, userHandler.OfferCheck)
	r.POST("/logout", userHandler.Logout)

	return r
}
