package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattrax/Mattrax/internal/db"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/authentication"
	"golang.org/x/crypto/bcrypt"
)

func Login(srv *mattrax.Server) http.HandlerFunc {
	type Request struct {
		UPN      string `json:"upn"`
		Password string `json:"password"`
	}

	type Response struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var cmd Request
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := srv.DB.GetUserForLogin(r.Context(), cmd.UPN)
		if err != nil {
			log.Printf("[GetUserForLogin Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.Password.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(cmd.Password)); err == bcrypt.ErrMismatchedHashAndPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("[CompareHashAndPassword Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var audience string
		if user.PermissionLevel == db.UserPermissionLevelAdministrator {
			audience = "dashboard"
		} else {
			audience = "enrollment"
		}

		authToken, _, err := srv.Auth.IssueToken(audience, authentication.AuthClaims{
			Subject:      cmd.UPN,
			FullName:     user.Fullname,
			Organisation: srv.Settings.Get().TenantName,
		})
		if err != nil {
			log.Printf("[IssueToken Error]: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(Response{
			Token: authToken,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
