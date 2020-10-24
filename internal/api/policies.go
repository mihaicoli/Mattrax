package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

func Policies(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := srv.DB.GetPolicies(r.Context())
		if err != nil {
			log.Printf("[GetPolicies Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Policy(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := srv.DB.GetPolicy(r.Context(), int32(id))
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetPolicy Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
