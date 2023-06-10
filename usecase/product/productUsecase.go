package usecase

import (
	"errors"
	"zog/domain/entity"
	repository "zog/repository/product"
)

type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProduct(productRepo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (pu ProductUsecase) ExecuteTicketList(page, limit int, location string) ([]entity.Ticket, error) {
	offset := (page - 1) * limit
	if location == "" {
		ticketlist, err := pu.productRepo.GetAllTickets(offset, limit)
		if err != nil {
			return nil, err
		} else {
			return ticketlist, nil
		}
	} else {
		ticketlist, err := pu.productRepo.GetAllTicketsByLocation(offset, limit, location)
		if err != nil {
			return nil, err
		} else {
			return ticketlist, nil
		}
	}

}

func (p *ProductUsecase) ExecuteTicketSearch(page, limit int, search string) ([]entity.Ticket, error) {
	offset := (page - 1) * limit
	ticketList, err := p.productRepo.GetAllTicketsBySearch(offset, limit, search)
	if err != nil {
		return nil, err
	} else {
		return ticketList, nil
	}
}

func (pu ProductUsecase) ExecuteTicketDetails(id int) (*entity.Ticket, *entity.TicketDetails, error) {

	ticket, err := pu.productRepo.GetTicketByID(id)
	if err != nil {
		return nil, nil, err
	}
	ticketDetails, err := pu.productRepo.GetTicketDetailsByID(id)
	if err != nil {
		return nil, nil, err
	}
	return ticket, ticketDetails, nil
}

func (pu ProductUsecase) ExecuteCreateTicket(ticket entity.Ticket) (int, error) {
	err := pu.productRepo.GetByName(ticket.Name)
	if err == nil {
		return 0, errors.New("product already exists")
	}
	newTicket := &entity.Ticket{
		Name:     ticket.Name,
		Price:    ticket.Price,
		Date:     ticket.Date,
		Location: ticket.Location,
		ImageURL: ticket.ImageURL,
		Category: ticket.Category,
		AdminId:  ticket.AdminId,
	}
	ticketId, err := pu.productRepo.CreateTicket(newTicket)
	if err != nil {
		return 0, err
	} else {
		return ticketId, nil
	}

}

func (pu *ProductUsecase) ExecuteCreateTicketDetails(details entity.TicketDetails) error {
	ticketDetails := &entity.TicketDetails{
		TicketId:    details.TicketId,
		Description: details.Description,
		Venue:       details.Venue,
	}
	err := pu.productRepo.CreateTicketDetails(ticketDetails)
	if err != nil {
		return errors.New("Creating Details failed")
	} else {
		return nil
	}

}

func (et ProductUsecase) ExecuteEditTicket(ticket entity.Ticket, id int) error {
	newTicket, err := et.productRepo.GetTicketByID(id)
	if err != nil {
		return err
	}
	newTicket = &entity.Ticket{
		Name:     ticket.Name,
		Price:    ticket.Price,
		Date:     ticket.Date,
		Location: ticket.Location,
		ImageURL: ticket.ImageURL,
		AdminId:  ticket.AdminId,
	}
	err1 := et.productRepo.UpdateTicket(newTicket)
	if err1 != nil {
		return err1
	} else {
		return nil
	}
}

func (dt ProductUsecase) ExecuteDeleteTicket(id int) error {
	result, err := dt.productRepo.GetTicketByID(id)
	if err != nil {
		return err
	}
	result.Removed = !result.Removed
	err1 := dt.productRepo.UpdateTicket(result)
	if err1 != nil {
		return errors.New("Ticket deletion unsuccesfull")
	}
	return nil
}

func (pu *ProductUsecase) ExecuteApperalList(page, limit int, category string) ([]entity.Apparel, error) {
	offset := (page - 1) * limit
	if category == "" {
		apparellist, err := pu.productRepo.GetAllApparels(offset, limit)
		if err != nil {
			return nil, err
		}
		return apparellist, nil
	} else {
		apparelList, err := pu.productRepo.GetAllApparelsByCategory(offset, limit, category)
		if err != nil {
			return nil, err
		}
		return apparelList, nil
	}

}

func (p *ProductUsecase) ExecuteApparelSearch(page, limit int, search string) ([]entity.Apparel, error) {
	offset := (page - 1) * limit
	apparelList, err := p.productRepo.GetAllApparelsBySearch(offset, limit, search)
	if err != nil {
		return nil, err
	} else {
		return apparelList, nil
	}
}
func (pu *ProductUsecase) ExecuteApparelDetails(id int) (*entity.Apparel, error) {
	apparel, err := pu.productRepo.GetApparelByID(id)
	if err != nil {
		return nil, err
	}
	return apparel, nil
}
func (pu ProductUsecase) ExecuteCreateApparel(apparel entity.Apparel) (int, error) {
	err := pu.productRepo.GetByApparelName(apparel.Name)
	if err == nil {
		return 0, errors.New("product already exists")
	}
	newApparel := &entity.Apparel{
		Name:     apparel.Name,
		Price:    apparel.Price,
		ImageURL: apparel.ImageURL,
		Removed:  apparel.Removed,
		AdminId:  apparel.AdminId,
	}
	apparelId, err1 := pu.productRepo.CreateApparel(newApparel)
	if err1 != nil {
		return 0, err1
	} else {
		return apparelId, nil
	}
}

func (et ProductUsecase) ExecuteEditApparel(apparel entity.Apparel, id int) error {
	newApparel, err := et.productRepo.GetApparelByID(id)
	if err != nil {
		return err
	}
	newApparel = &entity.Apparel{
		Name:     apparel.Name,
		Price:    apparel.Price,
		ImageURL: apparel.ImageURL,
		Removed:  apparel.Removed,
		AdminId:  apparel.AdminId,
	}
	err1 := et.productRepo.UpdateApparel(newApparel)
	if err1 != nil {
		return err1
	} else {
		return nil
	}
}
func (p *ProductUsecase) ExecuteCreateApparelDetails(apparelDetails entity.ApparelDetails) error {
	err := p.productRepo.CreateApparelDetails(&apparelDetails)
	if err != nil {
		return errors.New("Creating Details failed")
	} else {
		return nil
	}
}
func (dt ProductUsecase) ExecuteDeleteApparel(id int) error {
	result, err := dt.productRepo.GetApparelByID(id)
	if err != nil {
		return err
	}
	result.Removed = !result.Removed
	err1 := dt.productRepo.UpdateApparel(result)
	if err1 != nil {
		return errors.New("Apparel deletion unsuccesfull")
	}
	return nil
}
func (pu *ProductUsecase) ExecuteCreateInventory(inventory entity.Inventory) error {

	err := pu.productRepo.CreateInventory(&inventory)
	if err != nil {
		return errors.New("Creating inventory failed")
	} else {
		return nil
	}

}

func (pu *ProductUsecase) ExecuteQuantityUpdate(inventory entity.Inventory, method string, quantity int) error {
	if inventory.ProductCategory == "ticket" {
		product, err := pu.productRepo.GetByProductId(inventory.ProductId, inventory.ProductCategory)
		if err != nil {
			return errors.New("product not found in inventory")
		}
		if method == "increase" {
			product.Quantity += quantity
			err := pu.productRepo.UpdateInventory(product)
			if err != nil {
				return errors.New("ticket quantity update failed")
			}
		} else {
			product.Quantity -= quantity
			err := pu.productRepo.UpdateInventory(product)
			if err != nil {
				return errors.New("ticket quantity update failed")
			}
		}

	} else if inventory.ProductCategory == "apparel" {
		product, err := pu.productRepo.GetByProductId(inventory.ProductId, inventory.ProductCategory)
		if err != nil {
			return errors.New("product not found in inventory")
		}
		if method == "increase" {
			product.Quantity += quantity
			err := pu.productRepo.UpdateInventory(product)
			if err != nil {
				return errors.New("ticket quantity update failed")
			}
		} else {
			product.Quantity -= quantity
			err := pu.productRepo.UpdateInventory(product)
			if err != nil {
				return errors.New("ticket quantity update failed")
			}
		}
	}
	return nil
}

func (p *ProductUsecase) ExecuteAddCoupon(coupon *entity.Coupon) error {
	err := p.productRepo.CreateCoupon(coupon)
	if err != nil {
		return errors.New("Creating Coupon failed")
	} else {
		return nil
	}
}

func (p *ProductUsecase) ExecuteAddOffer(offer *entity.Offer) error {
	err := p.productRepo.CreateOffer(offer)
	if err != nil {
		return errors.New("Creating Offer failed")
	} else {
		return nil
	}
}

func (p *ProductUsecase) ExecuteAvailableCoupons() (*[]entity.Coupon, error) {
	coupons, err := p.productRepo.GetAllCoupons()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	availableCoupons := []entity.Coupon{}
	for _, coupon := range *coupons {
		if coupon.UsageLimit != coupon.UsedCount {
			availableCoupons = append(availableCoupons, coupon)
		}
	}
	return &availableCoupons, nil
}
