package todos

import (
	"fmt"
	"net/http"
)

func JSONError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
