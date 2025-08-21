package utils

import (
	"net/http"
	"encoding/json"
)


func RespondWithJson(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	msg, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Write(msg)
	return nil
}