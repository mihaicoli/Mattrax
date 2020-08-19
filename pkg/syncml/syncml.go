package syncml

import (
	"fmt"
	"net/http"

	"github.com/mattrax/Mattrax/pkg"
	"github.com/mattrax/xml"
)

// MaxRequestBodySize is the maximum amount of data that is allowed in a single request
const MaxRequestBodySize = 524288

// Read safely decodes a SyncML request from the HTTP body into a struct
func Read(r *http.Request, w http.ResponseWriter) (Message, bool) {
	if r.ContentLength > MaxRequestBodySize {
		if pkg.ErrorHandler != nil {
			pkg.ErrorHandler(fmt.Sprintf("Request body of size '%d' is larger than the maximum supported size of '%d'", r.ContentLength, MaxRequestBodySize), nil)
		}
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return Message{}, true
	}

	var v Message
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)
	if err := xml.NewDecoder(r.Body).Decode(&v); err != nil {
		if pkg.ErrorHandler != nil {
			pkg.ErrorHandler(fmt.Sprintf("Error decoding request of type '%T'", v), err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return Message{}, true
	}

	return v, false
}
