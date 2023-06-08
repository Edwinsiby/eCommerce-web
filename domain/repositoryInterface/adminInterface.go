package repository

import "zog/domain/entity"

type AdminInterface interface {
	GetByID(id uint) (*entity.Admin, error)
	GetByEmail(email string) (*entity.Admin, error)
	Create(admin *entity.Admin) error
	Update(admin *entity.Admin) error
	Delete(id uint) error
}
