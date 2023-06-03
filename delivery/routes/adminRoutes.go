package routes

import (
	"zog/delivery/handlers"
	m "zog/delivery/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine, adminHandler *handlers.AdminHandler) *gin.Engine {

	r.POST("/registeradmin", adminHandler.RegisterAdmin)
	r.POST("/adminloginpassword", adminHandler.AdminLoginWithPassword)
	r.POST("/adminlogin", adminHandler.Login)
	r.POST("/adminotpvalidation", adminHandler.LoginOtpValidation)
	r.GET("/adminhome", m.AdminRetriveCookie, adminHandler.Home)

	r.GET("/usermanagement", m.AdminRetriveCookie, adminHandler.UserList)
	r.GET("/sortuser", m.AdminRetriveCookie, adminHandler.SortUserByPermission)
	r.GET("/searchuser", m.AdminRetriveCookie, adminHandler.SearchUser)
	r.POST("/userpermission/:id", m.AdminRetriveCookie, adminHandler.TogglePermission)

	r.POST("/addticket", m.AdminRetriveCookie, adminHandler.CreateTicket)
	r.PUT("/editticket/:id", m.AdminRetriveCookie, adminHandler.EditTicket)
	r.DELETE("/deleteticket/:id", m.AdminRetriveCookie, adminHandler.DeleteTicket)

	r.POST("/addapparel", m.AdminRetriveCookie, adminHandler.CreateApparel)
	r.PUT("/editappaerl/:id", m.AdminRetriveCookie, adminHandler.EditApparel)
	r.DELETE("/deleteapparel/:id", m.AdminRetriveCookie, adminHandler.DeleteApparel)

	r.POST("/addcoupon", m.AdminRetriveCookie, adminHandler.AddCoupon)
	r.POST("/addoffer", m.AdminRetriveCookie, adminHandler.AddOffer)

	return r
}
