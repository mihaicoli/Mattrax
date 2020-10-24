package api

import (
	"encoding/json"
	"log"
	"net/http"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/authentication"
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

		if !user.Password.Valid || user.Password.String != cmd.Password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authToken, _, err := srv.Auth.IssueToken(authentication.AuthClaims{
			Subject:      cmd.UPN,
			FullName:     user.Fullname,
			Organisation: srv.Settings.Get().TenantName,
		})
		if err != nil {
			// log.Error().Err(err).Str("upn", cmd.UPN).Msg("Failed to sign JWT")
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
