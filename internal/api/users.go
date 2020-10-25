package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mattrax/Mattrax/internal/db"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

func Users(srv *mattrax.Server) http.HandlerFunc {
	type CreateUserRequest struct { // TODO: Merge with DB Struct and Fix sql.NullString issue
		Upn      string `json:"upn"`
		Fullname string `json:"fullname"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
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
			return
		} else if r.Method == http.MethodPost {
			fmt.Println("here")

			var cmd CreateUserRequest
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			fmt.Println(cmd)

			if cmd.Upn == "" || cmd.Fullname == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := srv.DB.CreateUser(r.Context(), db.CreateUserParams{
				Upn:      cmd.Upn,
				Fullname: cmd.Fullname,
				Password: sql.NullString{
					Valid:  true,
					String: cmd.Password,
				},
			}); err != nil {
				log.Printf("[CreateUser Error]: %s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
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
