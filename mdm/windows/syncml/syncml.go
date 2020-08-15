package syncml

import (
	"fmt"
	"net/http"

	"github.com/mattrax/xml"
	"github.com/rs/zerolog/log"
)

// MaxRequestBodySize is the maximum amount of data that is allowed in a single request
const MaxRequestBodySize = 524288

// Read safely decodes a SyncML request from the HTTP body into a struct
func Read(r *http.Request, w http.ResponseWriter) (Message, bool) {
	if r.ContentLength > MaxRequestBodySize {
		log.Debug().Int64("length", r.ContentLength).Int("max-length", MaxRequestBodySize).Msg("Request body larger than supported value")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return Message{}, true
	}

	var v Message
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)
	if err := xml.NewDecoder(r.Body).Decode(&v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Str("caller", fmt.Sprintf("%T", v)).Err(err).Msg("Error decoding request")
		return Message{}, true
	}

	return v, false
}
