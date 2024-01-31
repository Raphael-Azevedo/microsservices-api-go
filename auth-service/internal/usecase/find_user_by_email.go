package usecase

import (
	"biz-hub-auth-service/internal/entity"
)

type FindUserByEmailInput struct {
	Email string `json:"email"`
}

type FindUserByEmailUseCase struct {
	UserRepository entity.UserRepositoryInterface
}

func NewFindUserByEmailUseCase(
	UserRepository entity.UserRepositoryInterface,
) *FindUserByEmailUseCase {
	return &FindUserByEmailUseCase{
		UserRepository: UserRepository,
	}
}

func (c *FindUserByEmailUseCase) Execute(input FindUserByEmailInput) (entity.User, error) {

	user, err := c.UserRepository.FindUserByEmail(input.Email)
	if err != nil {
		return entity.User{}, err
	}

	return *user, nil
}
