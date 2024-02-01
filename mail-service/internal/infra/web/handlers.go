package web

import (
	"encoding/json"
	"log"
	"mail-service/internal/mailer"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type WebMailHandler struct {
	mailer mailer.Mail
}

func NewWebMailHandler(
	mailer mailer.Mail,
) *WebMailHandler {

	return &WebMailHandler{
		mailer: mailer,
	}
}

func (l *WebMailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload Payload

	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "parameters invalid", http.StatusBadRequest)
		return
	}
	log.Println("Received request payload:", requestPayload)

	msg := mailer.Message{
		From:    "me@example.com",
		To:      requestPayload.Data,
		Subject: "Test email",
		Data:    "Hello world!",
	}

	log.Println("Sending email...")
	err = l.mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.Data,
	}

	out, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "error to create email", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write(out)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
}
