package entity

type UserRepositoryInterface interface {
	Create(user *User) error
	FindUserByEmail(email string) (*User, error) 
}
