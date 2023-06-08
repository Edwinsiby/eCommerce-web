package models

type EditUser struct {
	FirstName string `json:"firstname" bson:"firstname" binding:"required"`
	LastName  string `json:"lastname" bson:"lastname" binding:"required"`
	Email     string `json:"email" bson:"email" binding:"required"`
}

type Signup struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
