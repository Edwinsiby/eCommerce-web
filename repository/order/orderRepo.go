package order

import (
	"errors"
	"time"
	"zog/domain/entity"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (or *OrderRepository) Create(order *entity.Order) (int, error) {
	if err := or.db.Create(order).Error; err != nil {
		return 0, err
	}
	return int(order.ID), nil
}

func (or *OrderRepository) GetByID(orderId int) (*entity.Order, error) {
	var order entity.Order
	result := or.db.Where("id=?", orderId).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Order not found")
		}
		return nil, errors.New("Order not found")
	}
	return &order, nil
}

func (or *OrderRepository) GetByRazorId(razorId string) (*entity.Order, error) {
	var order entity.Order
	result := or.db.Where("payment_id=?", razorId).First(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Order not found")
		}
		return nil, errors.New("Order not found")
	}
	return &order, nil
}

func (or *OrderRepository) Update(order *entity.Order) error {
	return or.db.Save(&order).Error
}

func (or *OrderRepository) CreateOrderItems(orderItem []entity.OrderItem) error {
	if err := or.db.Create(orderItem).Error; err != nil {
		return err
	}
	return nil
}

func (or *OrderRepository) GetAllOrders(userId, offset, limit int) ([]entity.Order, error) {
	var order []entity.Order
	result := or.db.Offset(offset).Limit(limit).Where("user_id=?", userId).Find(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Order not found")
		}
		return nil, errors.New("Order not found")
	}
	return order, nil
}
func (o *OrderRepository) GetByStatus(offset, limit int, status string) ([]entity.Order, error) {
	var order []entity.Order
	result := o.db.Offset(offset).Limit(limit).Where("status=?", status).Find(&order)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Order not found")
		}
		return nil, errors.New("Order not found")
	}
	return order, nil
}

func (or *OrderRepository) GetByDate(startDate, endDate time.Time) (*entity.SalesReport, error) {
	var Order []entity.Order
	var report entity.SalesReport

	if err := or.db.Model(&Order).Where("created_at BETWEEN ? AND ?", startDate, endDate).Select("SUM(total) as total_sales").Scan(&report).Error; err != nil {
		return nil, err
	}

	if err := or.db.Model(&Order).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&report.TotalOrders).Error; err != nil {
		return nil, err
	}

	if err := or.db.Model(&Order).Where("created_at BETWEEN ? AND ?", startDate, endDate).Select("AVG(total) as average_order").Scan(&report).Error; err != nil {
		return nil, err
	}

	if err := or.db.Model(&Order).Where("created_at BETWEEN ? AND ?", startDate, endDate).Select("payment, COUNT(payment) as payment_method_cnt").Group("payment").Scan(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (or *OrderRepository) GetByCategory(category string, startDate, endDate time.Time) (*entity.SalesReport, error) {
	report := &entity.SalesReport{}

	var orderItems []entity.OrderItem
	if err := or.db.Where("category = ? AND created_at BETWEEN ? AND ?", category, startDate, endDate).Find(&orderItems).Error; err != nil {
		return nil, err
	}

	totalSales := 0.0
	totalOrders := int64(len(orderItems))

	for _, item := range orderItems {
		totalSales += item.Price * float64(item.Quantity)

	}

	report.TotalSales = totalSales
	report.TotalOrders = totalOrders
	report.AverageOrder = totalSales / float64(totalOrders)

	return report, nil
}

func (or *OrderRepository) CreateInvoice(invoice *entity.Invoice) (*entity.Invoice, error) {
	if err := or.db.Create(invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

func (or *OrderRepository) CreateReturn(returnData *entity.Return) error {
	if err := or.db.Create(returnData).Error; err != nil {
		return err
	} else {
		return nil
	}
}
func (or *OrderRepository) GetReturnByID(returnId int) (*entity.Return, error) {
	var returnData entity.Return
	result := or.db.Where("id=?", returnId).First(&returnData)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Return Data not found")
		}
		return nil, errors.New("Return Data not found")
	}
	return &returnData, nil
}
func (or *OrderRepository) GetReturnByOrderID(orderId int) (*entity.Return, error) {
	var returnData entity.Return
	result := or.db.Where("order_id=?", orderId).First(&returnData)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Return Data not found")
		}
		return nil, errors.New("Return Data not found")
	}
	return &returnData, nil
}
func (or *OrderRepository) UpdateReturn(returnData *entity.Return) error {
	return or.db.Save(&returnData).Error
}

func (or *OrderRepository) UpdateUserWallet(user *entity.User) error {
	return or.db.Save(&user).Error
}
