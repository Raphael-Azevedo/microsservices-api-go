//go:build wireinject
// +build wireinject

package main

import (
	"biz-hub-auth-service/internal/entity"
	"biz-hub-auth-service/internal/event"
	"biz-hub-auth-service/internal/infra/database"
	"biz-hub-auth-service/internal/infra/web"
	"biz-hub-auth-service/internal/usecase"
	"biz-hub-auth-service/pkg/events"

	"database/sql"

	"github.com/google/wire"
)

var setUserRepositoryDependency = wire.NewSet(
	database.NewUserRepository,
	wire.Bind(new(entity.UserRepositoryInterface), new(*database.UserRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewUserCreated,
	wire.Bind(new(events.EventInterface), new(*event.UserCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setUserCreatedEvent = wire.NewSet(
	event.NewUserCreated,
	wire.Bind(new(events.EventInterface), new(*event.UserCreated)),
)

var setFindUserByEmailEvent = wire.NewSet(
	event.NewFindUserByEmail,
	wire.Bind(new(events.EventInterface), new(*event.FindUserByEmail)),
)

func NewCreateUserUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateUserUseCase {
	wire.Build(
		setUserRepositoryDependency,
		setUserCreatedEvent,
		usecase.NewCreateUserUseCase,
	)
	return &usecase.CreateUserUseCase{}
}

func NewFindUserByEmailUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.FindUserByEmailUseCase {
	wire.Build(
		setUserRepositoryDependency,
		setFindUserByEmailEvent,
		usecase.NewFindUserByEmailUseCase,
	)
	return &usecase.FindUserByEmailUseCase{}
}

func NewWebUserHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebUserHandler {
	wire.Build(
		setUserRepositoryDependency,
		setUserCreatedEvent,
		web.NewWebUserHandler,
	)
	return &web.WebUserHandler{}
}
