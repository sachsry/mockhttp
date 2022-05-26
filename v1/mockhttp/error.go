package mockhttp

import "fmt"

type ServerError struct {
	Status       string `json:"status"`
	DebugMessage string `json:"message"`
	Error        string `json:"error"`
}

func (s *ServerError) toString() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("Status: (%s)\nDebugMessage: (%s)\nError: (%s)", s.Status, s.DebugMessage, s.Error)
}
