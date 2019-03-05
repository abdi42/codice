package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

type ErrorString struct {
	S string
}

func (e *ErrorString) Error() string {
	return e.S
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
