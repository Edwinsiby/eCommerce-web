package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"zog/domain/entity"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db}
}

func (ac *AdminRepository) Create(admin *entity.Admin) error {
	return ac.db.Create(admin).Error
}

func (ar *AdminRepository) CreateOtpKey(key, phone string) error {
	var otpKey entity.OtpKey
	otpKey.Key = key
	otpKey.Phone = phone
	return ar.db.Create(otpKey).Error
}

func (ar *AdminRepository) GetByPhone(phone string) (*entity.Admin, error) {
	var admin entity.Admin
	result := ar.db.Where(&entity.Admin{Phone: phone}).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &admin, nil
}
func (ar *AdminRepository) GetByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin
	result := ar.db.Where(&entity.Admin{Email: email}).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &admin, nil
}
func (ar *AdminRepository) GetByID(id int) (*entity.User, error) {
	var user entity.User
	result := ar.db.Where(&entity.User{ID: id}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ar *AdminRepository) Update(user *entity.User) error {
	return ar.db.Save(user).Error
}

func (ar *AdminRepository) GetAllUsers(offset, limit int) ([]entity.User, error) {
	var users []entity.User
	err := ar.db.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a AdminRepository) GetAllUsersByPermission(offset, limit int, permission bool) ([]entity.User, error) {
	var users []entity.User
	err := a.db.Offset(offset).Limit(limit).Where("permission = ?", permission).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a *AdminRepository) GetAllUsersById(userId int) ([]entity.User, error) {
	var users []entity.User
	err := a.db.Where("id = ?", userId).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (a *AdminRepository) GetAllUsersByName(name string) ([]entity.User, error) {
	var users []entity.User
	err := a.db.Where("first_name LIKE ?", name).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ar *AdminRepository) GetUsers() (int, int, error) {
	var totalUsers, newUsers int64
	if err := ar.db.Model(&entity.User{}).Count(&totalUsers).Error; err != nil {
		return 0, 0, err
	}

	if err := ar.db.Model(&entity.User{}).Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).Count(&newUsers).Error; err != nil {
		return 0, 0, err
	}
	return int(totalUsers), int(newUsers), nil
}

func (ar *AdminRepository) GetProducts() (int, string, error) {
	var totalApparels, totalTickets int64
	var totalProducts int
	var stocklessCategory string

	if err := ar.db.Model(&entity.Ticket{}).Where("removed = ?", false).Count(&totalTickets).Error; err != nil {
		return 0, "", err
	}

	if err := ar.db.Model(&entity.Apparel{}).Where("removed = ?", false).Count(&totalApparels).Error; err != nil {
		return 0, "", err
	}

	totalProducts = int(totalApparels) + (int(totalTickets))
	query := ar.db.Model(&entity.Inventory{}).
		Select("category").
		Where("category IN (?)", []string{"ticket", "apparel"}).
		Order("quantity").
		Limit(1).
		Pluck("category", &stocklessCategory)
	if query.Error != nil {
		return 0, "", query.Error
	}
	return totalProducts, stocklessCategory, nil

}

func (ar *AdminRepository) GetOrders() (int, int, error) {
	var totalOrders int64
	var totalAmount float64

	if err := ar.db.Model(&entity.Order{}).Count(&totalOrders).Error; err != nil {
		return 0, 0, err
	}

	if err := ar.db.Model(&entity.Order{}).Select("AVG(total)").Row().Scan(&totalAmount); err != nil {
		return 0, 0, err
	}

	return int(totalOrders), int(totalAmount), nil
}

func (ar *AdminRepository) GetOrderByStatus() (int, int, error) {
	var pendingOrders, returnOrders int64

	if err := ar.db.Model(&entity.Order{}).Where("status = ?", "pending").Count(&pendingOrders).Error; err != nil {
		return 0, 0, err
	}

	if err := ar.db.Model(&entity.Order{}).Where("status = ?", "returned").Count(&returnOrders).Error; err != nil {
		return 0, 0, err
	}

	return int(pendingOrders), int(returnOrders), nil
}

func (er *AdminRepository) GetRevenue() (int, error) {
	var totalRevenue float64

	if err := er.db.Model(&entity.Order{}).Select("SUM(total)").Row().Scan(&totalRevenue); err != nil {
		return 0, err
	}

	return int(totalRevenue), nil
}

func (ar *AdminRepository) GetQuery() (int, error) {
	return 0, nil
}
