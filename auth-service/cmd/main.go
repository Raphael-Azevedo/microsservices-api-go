package main

import (
	"biz-hub-auth-service/configs"
	"biz-hub-auth-service/internal/infra/web/webserver"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MenssagerConfig struct {
	Rabbit *amqp.Connection
}

// @title           Go Auth Service
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Raphael Azevedo
// @contact.url    https://www.linkedin.com/in/raphael-a-neves/
// @contact.email  rfcompanhia@hotmail.com

// @host      localhost:8081
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @security MySecurityScheme
// @name cors
// @in header
// @type apiKey

// @schemes http https
func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := MenssagerConfig{
		Rabbit: rabbitConn,
	}

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

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webUserHandler := NewWebUserHandler(db, app.Rabbit)
	webserver.AddHandler("/user", webUserHandler.Create)
	webserver.AddHandler("/user/login", webUserHandler.Login)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start(configs.TokenAuth, configs.JwtExperesIn)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
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

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
