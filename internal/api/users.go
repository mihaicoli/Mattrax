package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

func Users(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := srv.DB.GetUsers(r.Context())
		if err != nil {
			log.Printf("[GetUsers Error]: %s\n", err)
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

func User(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := srv.DB.GetUser(r.Context(), vars["upn"])
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetUser Error]: %s\n", err)
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
