package webserver

import (
	"net/http"

	_ "biz-hub-auth-service/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start(token *jwtauth.JWTAuth, expiresIn int) {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.WithValue("jwt", token))
	s.Router.Use(middleware.WithValue("JwtExperesIn", expiresIn))
	s.Router.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.Router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/docs/doc.json")))

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	if origin == "http://example.com" {
		return true
	}
	return true
}
