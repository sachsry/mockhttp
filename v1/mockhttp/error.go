package mockhttp

import "fmt"

type ServerError struct {
	Status       string `json:"status"`
	DebugMessage string `json:"message"`
	Error        string `json:"error"`
}

func ValidateErrors(expected, result ServerError) error {
	if expected.Status != result.Status {
		return fmt.Errorf("expected status: %s, but got %s", expected.Status, result.Status)
	}
	if expected.DebugMessage != result.DebugMessage {
		return fmt.Errorf("expected message: %s, but got %s", expected.DebugMessage, result.DebugMessage)
	}
	if expected.Error != result.Error {
		return fmt.Errorf("expected error: %s, but got %s", expected.Error, result.Error)
	}
	return nil
}
