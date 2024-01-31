package usecase

import (
	"biz-hub-auth-service/internal/dto"
	"biz-hub-auth-service/internal/entity"
)

type CreateUserUseCase struct {
	UserRepository entity.UserRepositoryInterface
}

func NewCreateUserUseCase(
	UserRepository entity.UserRepositoryInterface,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: UserRepository,
	}
}

func (c *CreateUserUseCase) Execute(input dto.CreateUserInput) (entity.User, error) {
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return entity.User{}, err
	}

	if err := c.UserRepository.Create(user); err != nil {
		return entity.User{}, err
	}

	return *user, nil
}
