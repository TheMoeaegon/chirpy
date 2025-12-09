package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Moee1149/chirpy/internal/auth"
)

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileServerHits.Load()
	html := fmt.Sprintf("<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>", hits)
	w.Write([]byte(html))
}

func (cfg *apiConfig) handleReset(platform string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Store(0)
		if strings.ToLower(platform) != "dev" {
			responsdWithError(w, 403, "Forbidden")
			return
		}
		cfg.dbQueries.DropTable(r.Context())
		type response struct {
			Message string `json:"message"`
		}
		msg := response{
			Message: "Delete User Table",
		}
		respondWithJSON(w, 200, msg)
	}
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responsdWithError(w, 400, err.Error())
	}
	refresh_token, err := cfg.dbQueries.GetToken(r.Context(), refreshToken)
	if err == sql.ErrNoRows {
		responsdWithError(w, 401, "Unauthorized: token not found")
		return
	}
	if err != nil {
		responsdWithError(w, 500, "Internal Server Error")
		return
	}
	if time.Now().After(refresh_token.ExpiresAt) {
		responsdWithError(w, 401, "Unauthorized: token expired")
		return
	}
	token, err := auth.MakeJWT(refresh_token.UserID, cfg.jwtKey, 3600*time.Second)
	if err != nil {
		responsdWithError(w, 500, "Internal Server Error")
		return
	}
	resp := struct {
		TOKEN string `json:"token"`
	}{
		TOKEN: token,
	}
	respondWithJSON(w, 200, resp)
}

func (cfg *apiConfig) hanldeRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responsdWithError(w, 400, err.Error())
	}
	_, err = cfg.dbQueries.RevokeToken(r.Context(), refreshToken)
	if err == sql.ErrNoRows {
		responsdWithError(w, 401, "Unauthorized: token not found")
		return
	}
	if err != nil {
		responsdWithError(w, 500, "Internal Server Error")
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
