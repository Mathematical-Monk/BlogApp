package utils

import (
	"blogapi/models"
	"encoding/json"
	"net/http"
)




func RespondWithJson(w http.ResponseWriter, status int, payload models.HttpResponse) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	msg, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Write(msg)
	return nil
}