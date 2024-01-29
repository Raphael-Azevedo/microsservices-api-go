package usecase

import (
	"biz-hub-auth-service/internal/entity"
	"biz-hub-auth-service/pkg/events"
	"biz-hub-auth-service/internal/dto"
)

type CreateUserUseCase struct {
	UserRepository  entity.UserRepositoryInterface
	UserCreated     events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateUserUseCase(
	UserRepository entity.UserRepositoryInterface,
	UserCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository:  UserRepository,
		UserCreated:     UserCreated,
		EventDispatcher: EventDispatcher,
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

	c.UserCreated.SetPayload(user)
	c.EventDispatcher.Dispatch(c.UserCreated)

	return *user, nil
}
