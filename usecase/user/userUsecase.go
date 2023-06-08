package usecase

import (
	"errors"
	"zog/delivery/models"
	"zog/domain/entity"
	"zog/domain/utils"
	repository "zog/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUser(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (us *UserUsecase) ExecuteSignup(user entity.User) (*entity.User, error) {
	email, err := us.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if email != nil {
		return nil, errors.New("user with this email already exists")
	}
	phone, err := us.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, errors.New("error with server")
	}
	if phone != nil {
		return nil, errors.New("user with this phone no already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  string(hashedPassword),
	}

	err1 := us.userRepo.Create(newUser)
	if err1 != nil {
		return nil, err1
	}

	return newUser, nil
}

func (uu *UserUsecase) ExecuteSignupWithOtp(user models.Signup) (string, error) {
	var otpKey entity.OtpKey
	email, err := uu.userRepo.GetByEmail(user.Email)
	if err != nil {
		return "", errors.New("error with server")
	}
	if email != nil {
		return "", errors.New("user with this email already exists")
	}
	phone, err := uu.userRepo.GetByPhone(user.Phone)
	if err != nil {
		return "", errors.New("error with server")
	}
	if phone != nil {
		return "", errors.New("user with this phone no already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)
	key, err := utils.SendOtp(user.Phone)
	if err != nil {
		return "", err
	} else {
		err = uu.userRepo.CreateSignup(&user)
		otpKey.Key = key
		otpKey.Phone = user.Phone
		err = uu.userRepo.CreateOtpKey(&otpKey)
		if err != nil {
			return "", err
		}
		return key, nil
	}
}

func (uu *UserUsecase) ExecuteSignupOtpValidation(key string, otp string) error {
	result, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return err
	}
	user, err := uu.userRepo.GetSignupByPhone(result.Phone)
	if err != nil {
		return err
	}
	err = utils.CheckOtp(result.Phone, otp)
	if err != nil {
		return err
	} else {
		newUser := &entity.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Password:  user.Password,
		}

		err1 := uu.userRepo.Create(newUser)
		if err1 != nil {
			return err1
		} else {
			return nil
		}
	}

}

func (ul *UserUsecase) ExecuteLoginWithPassword(phone, password string) (int, error) {
	user, err := ul.userRepo.GetByPhone(phone)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, errors.New("user with this phone not found")
	}
	permission, err := ul.userRepo.CheckPermission(user)
	if permission == false {
		return 0, errors.New("user permission denied")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, errors.New("Invalid Password")
	} else {
		return user.ID, nil
	}

}

func (u *UserUsecase) ExecuteLogin(phone string) (string, error) {
	var otpKey entity.OtpKey
	result, err := u.userRepo.GetByPhone(phone)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("user with this phone not found")
	}
	permission, err := u.userRepo.CheckPermission(result)
	if permission == false {
		return "", errors.New("user permission denied")
	}
	key, err := utils.SendOtp(phone)
	if err != nil {
		return "", err
	} else {
		otpKey.Key = key
		otpKey.Phone = phone
		err = u.userRepo.CreateOtpKey(&otpKey)
		if err != nil {
			return "", err
		}
		return key, nil
	}

}

func (uu *UserUsecase) ExecuteOtpValidation(key, otp string) (*entity.User, error) {
	result, err := uu.userRepo.GetByKey(key)
	if err != nil {
		return nil, err
	}
	user, err := uu.userRepo.GetByPhone(result.Phone)
	if err != nil {
		return nil, err
	}
	err1 := utils.CheckOtp(result.Phone, otp)
	if err1 != nil {
		return nil, err1
	}
	return user, nil
}

func (uu *UserUsecase) ExecuteAddAddress(address *entity.Address) error {
	err := uu.userRepo.CreateAddress(address)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) ExecuteShowUserDetails(userId int) (*entity.User, *[]entity.Address, error) {
	user, err := uu.userRepo.GetByID(userId)
	if err != nil {
		return nil, nil, err
	}
	address, err1 := uu.userRepo.GetAddressByUserId(userId)
	if err1 != nil {
		return nil, nil, err1
	}
	if user != nil && address != nil {
		return user, address, nil
	} else {
		return nil, nil, errors.New("user with this id not found")
	}

}

func (uu *UserUsecase) ExecuteEditProfile(user entity.User, userid int) error {
	user.ID = userid
	err := uu.userRepo.Update(&user)
	if err != nil {
		return errors.New("User details updation failed")
	}
	return nil
}

func (uu *UserUsecase) ExecuteChangePassword(userId int) (string, error) {
	var otpKey entity.OtpKey
	user, err := uu.userRepo.GetByID(userId)
	if err != nil {
		return "", err
	}
	key, err1 := utils.SendOtp(user.Phone)
	if err1 != nil {
		return "", err
	} else {
		otpKey.Key = key
		otpKey.Phone = user.Phone
		err = uu.userRepo.CreateOtpKey(&otpKey)
		if err != nil {
			return "", err
		}
		return key, nil
	}

}

func (uu *UserUsecase) ExecuteOtpValidationPassword(password string, otp string, userId int) error {
	user, err := uu.userRepo.GetByID(userId)
	if err != nil {
		return err
	}
	err = utils.CheckOtp(user.Phone, otp)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	err = uu.userRepo.Update(user)
	if err != nil {
		return errors.New("Password changing failed")
	}
	return nil

}
