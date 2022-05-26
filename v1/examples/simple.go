package examples

import (
	"net/http"

	"github.com/sachsry/mockhttp/v1/response"
)

// Return success if the length of the path is even
// Return failure if the length of the path is odd
func HandleSimple(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path)%2 != 0 {
		response.Error(w, 500, "length of the path is odd", nil)
		return
	}
	response.Success(w)
}
