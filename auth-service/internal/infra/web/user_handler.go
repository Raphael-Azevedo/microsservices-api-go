package web

import (
	"biz-hub-auth-service/internal/dto"
	"biz-hub-auth-service/internal/entity"
	"biz-hub-auth-service/internal/usecase"
	"biz-hub-auth-service/pkg/events"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type WebUserHandler struct {
	EventDispatcher events.EventDispatcherInterface
	UserRepository  entity.UserRepositoryInterface
	UserCreated     events.EventInterface
}

func NewWebUserHandler(
	EventDispatcher events.EventDispatcherInterface,
	UserRepository entity.UserRepositoryInterface,
	UserCreated events.EventInterface,
) *WebUserHandler {
	return &WebUserHandler{
		EventDispatcher: EventDispatcher,
		UserRepository:  UserRepository,
		UserCreated:     UserCreated,
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

	createUser := usecase.NewCreateUserUseCase(h.UserRepository, h.UserCreated, h.EventDispatcher)
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

	// err = h.logRequest("authentication", fmt.Sprintf("user created: %s", dto.Email))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
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

	userByEmail := usecase.NewFindUserByEmailUseCase(h.UserRepository, h.UserCreated, h.EventDispatcher)
	user, err := userByEmail.Execute(dtoInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	// err = h.logRequest("authentication", fmt.Sprintf("%s logged in", dtoInput.Email))
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

// TODO logService should .env
func (h *WebUserHandler) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://project-logger-service-1/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
