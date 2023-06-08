package usecase

import (
	"errors"
	"zog/domain/entity"
	"zog/domain/utils"
	repository "zog/repository/admin"

	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase struct {
	adminRepo *repository.AdminRepository
}

func NewAdmin(adminRepo *repository.AdminRepository) *AdminUsecase {
	return &AdminUsecase{adminRepo: adminRepo}
}

func (ac *AdminUsecase) ExecuteAdminCreate(admin entity.Admin) (*entity.Admin, error) {
	email, err := ac.adminRepo.GetByEmail(admin.Email)
	if err != nil {
		return nil, err
	}
	if email != nil {
		return nil, errors.New("admin with this email already exists")
	}
	phone, err := ac.adminRepo.GetByPhone(admin.Phone)
	if err != nil {
		return nil, err
	}
	if phone != nil {
		return nil, errors.New("admin with this phone already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newAdmin := &entity.Admin{
		AdminName: admin.AdminName,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Role:      admin.Role,
		Password:  string(hashedPassword),
	}

	err = ac.adminRepo.Create(newAdmin)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}
func (au *AdminUsecase) ExecuteLoginWithPassword(phone, password string) (int, error) {
	admin, err := au.adminRepo.GetByPhone(phone)
	if err != nil {
		return 0, err
	}
	if admin == nil {
		return 0, errors.New("admin with this phone not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return 0, errors.New("Invalid Password")
	} else {
		return int(admin.ID), nil
	}

}

func (au *AdminUsecase) ExecuteAdminLogin(phone string) error {
	result, err := au.adminRepo.GetByPhone(phone)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("admin with this phone not found")
	}
	key, err1 := utils.SendOtp(phone)
	if err1 != nil {
		return err
	} else {
		err = au.adminRepo.CreateOtpKey(key, phone)
		if err != nil {
			return err
		}
		return nil
	}

}
func (au *AdminUsecase) ExecuteOtpValidation(phone, otp string) (*entity.Admin, error) {
	result, err := au.adminRepo.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(phone, otp)
	if err1 != nil {
		return nil, err1
	}
	return result, nil
}

func (au *AdminUsecase) ExecuteAdminDashboard() (*entity.AdminDashboard, error) {

	totalUsers, newUsers, err := au.adminRepo.GetUsers()
	if err != nil {
		return nil, errors.New("fetching user count failed")
	}
	totalProducts, stocklessCategory, err := au.adminRepo.GetProducts()
	if err != nil {
		return nil, errors.New("fetching product count failed")
	}
	totalOrders, averageOrderValue, err := au.adminRepo.GetOrders()
	if err != nil {
		return nil, errors.New("fetching order count failed")
	}
	pendingOrders, returnOrders, err := au.adminRepo.GetOrderByStatus()
	if err != nil {
		return nil, errors.New("fetching order count failed")
	}
	totalrevenue, err := au.adminRepo.GetRevenue()
	if err != nil {
		return nil, errors.New("fetching revenue failed")
	}
	totalquery, err := au.adminRepo.GetQuery()
	if err != nil {
		return nil, errors.New("fetching query count failed")
	}

	dashboardResponse := entity.AdminDashboard{
		TotalUsers:        totalUsers,
		NewUsers:          newUsers,
		TotalProducts:     totalProducts,
		StocklessCategory: stocklessCategory,
		TotalOrders:       totalOrders,
		AverageOrderValue: averageOrderValue,
		PendingOrders:     pendingOrders,
		ReturnOrders:      returnOrders,
		TotalRevenue:      totalrevenue,
		TotalQuery:        totalquery,
	}
	return &dashboardResponse, nil
}

func (ul *AdminUsecase) ExecuteUserList(page, limit int) ([]entity.User, error) {
	offset := (page - 1) * limit
	userlist, err := ul.adminRepo.GetAllUsers(offset, limit)
	if err != nil {
		return nil, err
	}
	return userlist, nil
}

func (a *AdminUsecase) ExecuteUserListByPermission(page, limit int, permission string) ([]entity.User, error) {
	offset := (page - 1) * limit
	var sortby bool
	if permission == "true" {
		sortby = true
	} else {
		sortby = false
	}
	userlist, err := a.adminRepo.GetAllUsersByPermission(offset, limit, sortby)
	if err != nil {
		return nil, err
	}
	return userlist, nil
}

func (a *AdminUsecase) ExecuteSearchUserById(userId int) ([]entity.User, error) {
	userlist, err := a.adminRepo.GetAllUsersById(userId)
	if err != nil {
		return nil, err
	}
	return userlist, nil
}
func (a *AdminUsecase) ExecuteSearchUserByName(name string) ([]entity.User, error) {
	userlist, err := a.adminRepo.GetAllUsersByName(name)
	if err != nil {
		return nil, err
	}
	return userlist, nil
}

func (tp *AdminUsecase) ExecuteTogglePermission(id int) error {
	result, err := tp.adminRepo.GetByID(id)
	if err != nil {
		return err
	}
	result.Permission = !result.Permission
	err1 := tp.adminRepo.Update(result)
	if err1 != nil {
		return errors.New("user permission toggling failed")
	}
	return nil
}
