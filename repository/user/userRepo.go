package repository

import (
	"errors"
	"zog/delivery/models"
	"zog/domain/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) GetByID(id int) (*entity.User, error) {
	var user entity.User
	result := ur.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := ur.db.Where(&entity.User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func (ur *UserRepository) GetByPhone(phone string) (*entity.User, error) {
	var user entity.User
	result := ur.db.Where(&entity.User{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) CheckPermission(user *entity.User) (bool, error) {
	result := ur.db.Where(&entity.User{Phone: user.Phone}).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	permission := user.Permission
	return permission, nil
}

func (ur *UserRepository) GetByKey(key string) (*entity.OtpKey, error) {
	var otpKey entity.OtpKey
	result := ur.db.Where(&entity.OtpKey{Key: key}).First(&otpKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &otpKey, nil
}

func (ur *UserRepository) CreateSignup(user *models.Signup) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) GetSignupByPhone(phone string) (*models.Signup, error) {
	var user models.Signup
	result := ur.db.Where(&models.Signup{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) Create(user *entity.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) Update(user *entity.User) error {
	return ur.db.Updates(user).Error
}

func (ur *UserRepository) Delete(id uint) error {
	return ur.db.Delete(&entity.User{}, id).Error
}

func (ur *UserRepository) CreateAddress(address *entity.Address) error {
	return ur.db.Create(address).Error
}

func (ur *UserRepository) GetAddressByUserId(userid int) (*[]entity.Address, error) {
	var address []entity.Address
	result := ur.db.Where(&entity.Address{UserId: userid}).Find(&address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &address, nil
}
func (cr *UserRepository) GetAddressById(addressId int) (*entity.Address, error) {
	var address entity.Address
	result := cr.db.Where(&entity.Address{ID: addressId}).First(&address)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &address, nil
}

func (ar *UserRepository) CreateOtpKey(otpKey *entity.OtpKey) error {

	return ar.db.Create(otpKey).Error
}

// done
