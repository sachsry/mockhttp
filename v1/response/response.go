package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sends a standard 200 response with a generic body
func Success(w http.ResponseWriter) {
	body := map[string]string{
		"status": "ok",
	}
	SuccessWithBody(w, body)
}

// Sends a 200 response with JSON representation of provided body
func SuccessWithBody(w http.ResponseWriter, body interface{}) {
	res, err := json.Marshal(body)
	if err != nil {
		fmt.Println(fmt.Sprintf("unexpected error encountered marshaling json: %v", err))
		return
	}

	w.Write(res)
}

// Sends an error response without notifying new relic
func Error(w http.ResponseWriter, status int, message string, err error) {
	w.WriteHeader(status)

	body := map[string]string{
		"status": getErrorStatus(status),
	}
	if message != "" {
		body["message"] = message
	}
	if err != nil {
		body["error"] = err.Error()
	}

	res, merr := json.Marshal(body)
	if merr != nil {
		fmt.Println(fmt.Sprintf("unexpected error encountered marshaling json: %v", merr))
		return
	}
	w.Write(res)
}

func getErrorStatus(status int) string {
	switch status {
	case 400:
		return "bad request"
	case 401:
		return "unauthorized"
	case 403:
		return "forbidden"
	case 404:
		return "not found"
	default:
		return "internal error"
	}
}
