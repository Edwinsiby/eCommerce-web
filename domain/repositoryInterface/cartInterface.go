package repository

import "zog/domain/entity"

type CartInterface interface {
	GetByID(id uint) (*entity.Cart, error)
	GetByUserID(userID uint) (*entity.Cart, error)
	Create(cart *entity.Cart) error
	Update(cart *entity.Cart) error
	Delete(id uint) error
}
