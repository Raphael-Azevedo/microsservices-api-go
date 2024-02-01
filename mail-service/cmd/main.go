package main

import (
	"fmt"
	"log"
	"mail-service/internal/infra/web"
	"mail-service/internal/mailer"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const webPort = "80"

type Config struct {
	Mailer mailer.Mail
}

func main() {
	log.Println("Starting mail service on port", webPort)
	mailer := Config{
		Mailer: createMail(),
	}

	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	webMailerHandler := web.NewWebMailHandler(mailer.Mailer)
	router.Post("/send", webMailerHandler.SendMail)

	fmt.Println("Starting web server on port", webPort)
	http.ListenAndServe(":"+webPort, router)
}

func createMail() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mailer := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return mailer
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	if origin == "http://example.com" {
		return true
	}
	return true
}
