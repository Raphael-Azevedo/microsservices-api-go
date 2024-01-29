package web

import (
	"biz-hub-logger-service/data"
	"encoding/json"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Error struct {
	Message string `json:"message"`
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type WebLoggerHandler struct {
	Models data.Models
}

func NewWebLoggerHandler(
	Models data.Models,

) *WebLoggerHandler {
	return &WebLoggerHandler{
		Models: Models,
	}
}
func (l *WebLoggerHandler) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = l.Models.LogEntry.Insert(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	out, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
