package web

import (
	"biz-hub-auth-service/internal/dto"
	"biz-hub-auth-service/internal/entity"
	"biz-hub-auth-service/internal/event"
	"biz-hub-auth-service/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Error struct {
	Message string `json:"message"`
}

type WebUserHandler struct {
	UserRepository entity.UserRepositoryInterface
	Rabbit         *amqp.Connection
}

func NewWebUserHandler(
	UserRepository entity.UserRepositoryInterface,
	Rabbit *amqp.Connection,
) *WebUserHandler {
	return &WebUserHandler{
		UserRepository: UserRepository,
		Rabbit:         Rabbit,
	}
}

// Create godoc
// @Summary      Create
// @Description  Create
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /user [post]
func (h *WebUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userByEmail := usecase.NewFindUserByEmailUseCase(h.UserRepository)

	email := usecase.FindUserByEmailInput{
		Email: dto.Email,
	}

	user, _ := userByEmail.Execute(email)

	if user.Email == dto.Email {
		http.Error(w, "E-mail already registered", http.StatusConflict)
		return
	}

	createUser := usecase.NewCreateUserUseCase(h.UserRepository)

	output, err := createUser.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go h.logEventViaRabbit("authentication", fmt.Sprintf("user created: %s", dto.Email))
	go h.logEventViaRabbit("mail", dto.Email)
}

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.LoginInput  true  "user request"
// @Success      200
// @Failure      500         {object}  Error
// @Failure      404         {object}  Error
// @Failure      400         {object}  Error
// @Router       /user/login [post]
func (h *WebUserHandler) Login(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExperesIn").(int)
	var userDto dto.LoginInput
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dtoInput := usecase.FindUserByEmailInput{
		Email: userDto.Email,
	}

	userByEmail := usecase.NewFindUserByEmailUseCase(h.UserRepository)

	user, err := userByEmail.Execute(dtoInput)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	if !user.ValidatePassword(userDto.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	accessToken := dto.LoginOutput{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)

	go h.logEventViaRabbit("authentication", fmt.Sprintf("%s logged in", dtoInput.Email))
}

func (h *WebUserHandler) logEventViaRabbit(name, data string) {
	err := h.pushToQueue(name, data)
	if err != nil {
		panic(err)
	}
}

// pushToQueue pushes a message into RabbitMQ
func (h *WebUserHandler) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(h.Rabbit)
	if err != nil {
		return err
	}

	var payload struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	payload.Name = name
	payload.Data = msg

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
