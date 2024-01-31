//go:build wireinject
// +build wireinject

package main

import (
	"biz-hub-auth-service/internal/entity"
	"biz-hub-auth-service/internal/infra/database"
	"biz-hub-auth-service/internal/infra/web"
	"biz-hub-auth-service/internal/usecase"

	"database/sql"

	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

var setUserRepositoryDependency = wire.NewSet(
	database.NewUserRepository,
	wire.Bind(new(entity.UserRepositoryInterface), new(*database.UserRepository)),
)

func NewCreateUserUseCase(db *sql.DB) *usecase.CreateUserUseCase {
	wire.Build(
		setUserRepositoryDependency,
		usecase.NewCreateUserUseCase,
	)
	return &usecase.CreateUserUseCase{}
}

func NewFindUserByEmailUseCase(db *sql.DB) *usecase.FindUserByEmailUseCase {
	wire.Build(
		setUserRepositoryDependency,
		usecase.NewFindUserByEmailUseCase,
	)
	return &usecase.FindUserByEmailUseCase{}
}

func NewWebUserHandler(db *sql.DB, rabbit *amqp.Connection) *web.WebUserHandler {
	wire.Build(
		setUserRepositoryDependency,
		web.NewWebUserHandler,
	)
	return &web.WebUserHandler{}
}
