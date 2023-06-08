package routes

import (
	"zog/delivery/handlers"
	m "zog/delivery/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRouter(r *gin.Engine, orderHandler *handlers.OrderHandler) *gin.Engine {

	r.POST("/placeorder/:addressid/:payment", m.UserRetriveCookie, orderHandler.PlaceOrder)
	r.POST("/paymentverification/:sign/:razorid/:payid", m.UserRetriveCookie, orderHandler.PaymentVerification)
	r.PUT("/cancelorder/:orderid", m.UserRetriveCookie, orderHandler.CancelOrder)
	r.GET("/orderhistory", m.UserRetriveCookie, orderHandler.OrderHistory)
	r.POST("/orderreturn/:orderid", m.UserRetriveCookie, orderHandler.OrderReturn)

	// admin
	r.PUT("/updateorder/:orderid/:status", m.AdminRetriveCookie, orderHandler.AdminOrderUpdate)
	r.POST("/updatereturn/:returnid/:status/:refund", m.AdminRetriveCookie, orderHandler.AdminReturnUpdate)
	r.POST("/refund/:orderid", m.AdminRetriveCookie, orderHandler.AdminRefund)
	r.GET("/salesreportbydate/:start/:end", m.AdminRetriveCookie, orderHandler.SalesReportByDate)
	r.GET("/salesreportbyperiod/:period", m.AdminRetriveCookie, orderHandler.SalesReportByPeriod)
	r.GET("/salesreportbycategory/:category/:period", m.AdminRetriveCookie, orderHandler.SalesReportByCategory)
	r.POST("/sortorders", m.AdminRetriveCookie, orderHandler.SortOrderByStatus)

	return r
}
