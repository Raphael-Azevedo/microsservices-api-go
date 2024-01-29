package main

import (
	"biz-hub-auth-service/configs"
	"biz-hub-auth-service/internal/event/handler"
	"biz-hub-auth-service/internal/infra/web/webserver"
	"biz-hub-auth-service/pkg/events"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/streadway/amqp"
)

// @title           Go Auth Service
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Raphael Azevedo
// @contact.url    https://www.linkedin.com/in/raphael-a-neves/
// @contact.email  rfcompanhia@hotmail.com

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Println("started application")
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if configs.Env == "development" {
		err = runMigrations(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("UserCreated", &handler.UserCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webUserHandler := NewWebUserHandler(db, eventDispatcher)
	webserver.AddHandler("/user", webUserHandler.Create)
	webserver.AddHandler("/user/login", webUserHandler.Login)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start(configs.TokenAuth, configs.JwtExperesIn)
	var wg sync.WaitGroup
    wg.Add(1)
    wg.Wait()
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func runMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"users",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
