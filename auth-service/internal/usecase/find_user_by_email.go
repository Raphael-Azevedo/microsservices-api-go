package usecase

import (
	"biz-hub-auth-service/internal/entity"
	"biz-hub-auth-service/pkg/events"
)

type FindUserByEmailInput struct {
	Email string `json:"email"`
}

type FindUserByEmailUseCase struct {
	UserRepository  entity.UserRepositoryInterface
	FindUserByEmail events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewFindUserByEmailUseCase(
	UserRepository entity.UserRepositoryInterface,
	FindUserByEmail events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *FindUserByEmailUseCase {
	return &FindUserByEmailUseCase{
		UserRepository:  UserRepository,
		FindUserByEmail:     FindUserByEmail,
		EventDispatcher: EventDispatcher,
	}
}

func (c *FindUserByEmailUseCase) Execute(input FindUserByEmailInput) (entity.User, error) {

	user, err := c.UserRepository.FindUserByEmail(input.Email)
	if err != nil {
		return entity.User{}, err
	}

	// c.FindUserByEmail.SetPayload(user)
	// c.EventDispatcher.Dispatch(c.FindUserByEmail)

	return *user, nil
}
