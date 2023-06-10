package repository

import (
	"errors"
	"time"
	"zog/domain/entity"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (pr *ProductRepository) GetAllTickets(offset, limit int) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := pr.db.Offset(offset).Limit(limit).Where("removed = ?", false).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}
func (pr *ProductRepository) GetAllTicketsByLocation(offset, limit int, location string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := pr.db.Where("location = ?", location).Offset(offset).Limit(limit).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (p *ProductRepository) GetAllTicketsBySearch(offset, limit int, search string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := p.db.Where("name LIKE ?", search+"%").Offset(offset).Limit(limit).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (gt *ProductRepository) GetTicketByID(id int) (*entity.Ticket, error) {
	var ticket entity.Ticket
	result := gt.db.First(&ticket, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &ticket, nil
}

func (pr *ProductRepository) GetByName(name string) error {
	var existingTicket entity.Ticket
	result := pr.db.Where(&entity.Ticket{Name: name}).First(&existingTicket)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil
}

func (ct *ProductRepository) CreateTicket(ticket *entity.Ticket) (int, error) {
	if err := ct.db.Create(ticket).Error; err != nil {
		return 0, err
	}
	return ticket.ID, nil
}

func (ut *ProductRepository) UpdateTicket(ticket *entity.Ticket) error {
	return ut.db.Save(ticket).Error
}

func (dt *ProductRepository) DeleteTicket(ticket *entity.Ticket) error {
	return dt.db.Delete(ticket).Error
}

func (pr *ProductRepository) CreateTicketDetails(details *entity.TicketDetails) error {
	return pr.db.Create(details).Error
}

func (pr *ProductRepository) GetTicketDetailsByID(id int) (*entity.TicketDetails, error) {
	var ticketDetails entity.TicketDetails
	result := pr.db.Where("ticket_id=?", id).First(&ticketDetails)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &ticketDetails, nil
}
func (pr *ProductRepository) CreateInventory(inventory *entity.Inventory) error {
	return pr.db.Create(inventory).Error
}

func (pr *ProductRepository) GetByProductId(productId int, category string) (*entity.Inventory, error) {
	var existingProduct entity.Inventory
	result := pr.db.Where(&entity.Inventory{ProductCategory: category, ProductId: productId}).First(&existingProduct)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &existingProduct, nil
}

func (pr *ProductRepository) IncreaseProductQuantity(product *entity.Inventory) error {
	return pr.db.Save(product).Error
}

func (pr *ProductRepository) DecreaseProductQuantity(product *entity.Inventory) error {
	existingProduct := &entity.Inventory{}
	err := pr.db.Where("product_category = ? AND product_id =?", product.ProductCategory, product.ProductId).First(existingProduct).Error
	if err != nil {
		return err
	}
	newQuantity := existingProduct.Quantity - product.Quantity
	err = pr.db.Model(existingProduct).Update("Quantity", newQuantity).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) GetAllApparels(offset, limit int) ([]entity.Apparel, error) {
	var apparels []entity.Apparel
	err := pr.db.Offset(offset).Limit(limit).Where("removed = ?", false).Find(&apparels).Error
	if err != nil {
		return nil, err
	}
	return apparels, nil
}

func (p *ProductRepository) GetAllApparelsBySearch(offset, limit int, search string) ([]entity.Apparel, error) {
	var apparels []entity.Apparel
	err := p.db.Where("name LIKE ?", search+"%").Offset(offset).Limit(limit).Find(&apparels).Error
	if err != nil {
		return nil, err
	}
	return apparels, nil
}

func (p *ProductRepository) GetAllApparelsByCategory(offset, limit int, category string) ([]entity.Apparel, error) {
	var apparels []entity.Apparel
	err := p.db.Offset(offset).Limit(limit).Where("removed = ? AND sub_category = ?", false, category).Find(&apparels).Error
	if err != nil {
		return nil, err
	}
	return apparels, nil
}

func (gt *ProductRepository) GetApparelByID(id int) (*entity.Apparel, error) {
	var apparel entity.Apparel
	result := gt.db.First(&apparel, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &apparel, nil
}

func (pr *ProductRepository) GetByApparelName(apparelName string) error {

	var existingApparel entity.Ticket
	result := pr.db.Where(&entity.Apparel{Name: apparelName}).First(&existingApparel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil

}

func (ct *ProductRepository) CreateApparel(apparel *entity.Apparel) (int, error) {
	if err := ct.db.Create(apparel).Error; err != nil {
		return 0, err
	}
	return apparel.ID, nil
}
func (ct *ProductRepository) CreateApparelDetails(apparelDetails *entity.ApparelDetails) error {
	if err := ct.db.Create(apparelDetails).Error; err != nil {
		return err
	}
	return nil
}
func (ut *ProductRepository) UpdateApparel(apparel *entity.Apparel) error {
	return ut.db.Save(apparel).Error
}
func (dt *ProductRepository) DeleteApparel(apparel *entity.Apparel) error {
	return dt.db.Delete(apparel).Error
}
func (pr *ProductRepository) UpdateInventory(product *entity.Inventory) error {
	return pr.db.Delete(product).Error
}

func (p *ProductRepository) CreateCoupon(coupon *entity.Coupon) error {
	if err := p.db.Create(coupon).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) GetAllCoupons() (*[]entity.Coupon, error) {
	var coupns []entity.Coupon
	currentTime := time.Now()
	err := p.db.Where("valid_until >= ?", currentTime).Find(&coupns).Error
	if err != nil {
		return nil, err
	}
	return &coupns, nil
}

func (p *ProductRepository) GetCouponByCode(code string) (*entity.Coupon, error) {
	coupon := &entity.Coupon{}
	err := p.db.Where("code = ?", code).First(coupon).Error
	if err != nil {
		return nil, err
	}
	return coupon, nil
}

func (p *ProductRepository) UpdateCouponCount(coupon *entity.Coupon) error {
	return p.db.Save(coupon).Error
}

func (p *ProductRepository) UpdateCouponUsage(usedCoupon *entity.UsedCoupon) error {
	if err := p.db.Create(usedCoupon).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) CheckCouponUsage(usedCoupon *entity.UsedCoupon) error {

	return nil
}
func (p *ProductRepository) CreateOffer(offer *entity.Offer) error {
	if err := p.db.Create(offer).Error; err != nil {
		return err
	}
	return nil
}

func (p *ProductRepository) GetOfferByPrice(price int) (*[]entity.Offer, error) {
	offers := &[]entity.Offer{}
	err := p.db.Where("min_price <= ?", price).Find(offers).Error
	if err != nil {
		return nil, err
	} else if offers == nil {
		return nil, err
	}
	return offers, nil
}
