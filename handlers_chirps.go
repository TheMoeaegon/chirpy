package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Moee1149/chirpy/internal/auth"
	"github.com/Moee1149/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responsdWithError(w, 400, err.Error())
		return
	}

	_, err = auth.ValidateJwt(token, cfg.jwtKey)
	if err != nil {
		responsdWithError(w, 401, "Invalid Token")
		return
	}
	type parameters struct {
		Body    string `json:"body"`
		User_Id string `json:"user_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithJSON(w, 500, fmt.Sprintf("Error decoding json: %v", err))
		return
	}
	if len(params.Body) > 140 {
		responsdWithError(w, 400, "The chirpy is too long")
		return
	}
	//check for profane words
	cleanedBody := validateBadWords(params.Body)
	userId, err := uuid.Parse(params.User_Id)
	if err != nil {
		responsdWithError(w, 400, "Invalid user_id format")
		return
	}
	chirpsParams := database.CreateChirpyParams{
		Body:   cleanedBody,
		UserID: userId,
	}

	chirp, err := cfg.dbQueries.CreateChirpy(r.Context(), chirpsParams)
	if err != nil {
		responsdWithError(w, 500, fmt.Sprintf("Error adding chirps: %v", err))
		return
	}
	resp := chirpSchema{
		ID:        chirp.ID.String(),
		CreateAt:  chirp.CreatedAt.UTC().Format("2006-01-0215:04:05Z"),
		UpdatedAt: chirp.UpdatedAt.UTC().Format("2006-01-0215:04:05Z"),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	respondWithJSON(w, 201, resp)
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	if author_id != "" {
		user_id, err := uuid.Parse(author_id)
		if err != nil {
			responsdWithError(w, 400, err.Error())
			return
		}
		chirps, err := cfg.dbQueries.GetChirpsByAuthorId(r.Context(), user_id)
		if err != nil {
			responsdWithError(w, 500, fmt.Sprintf("Error getting chirps: %v", err))
			return
		}
		respondWithJSON(w, 200, chirps)
	} else {
		chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			responsdWithError(w, 500, fmt.Sprintf("Error getting chirps: %v", err))
			return
		}
		respondWithJSON(w, 200, chirps)
	}
}

func (cfg *apiConfig) handleGetChirpsById(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(pathValue)
	if err != nil {
		responsdWithError(w, 400, "Invalid user_id format")
		return
	}
	chirp, err := cfg.dbQueries.GetChirpsById(r.Context(), chirpId)
	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("Error getting chirps: %v", err))
		return
	}
	resp := chirpSchema{
		ID:        chirp.ID.String(),
		CreateAt:  chirp.CreatedAt.UTC().Format("2006-01-0215:04:05Z"),
		UpdatedAt: chirp.UpdatedAt.UTC().Format("2006-01-0215:04:05Z"),
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	respondWithJSON(w, 200, resp)
}

func (cfg *apiConfig) handleDeleteChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responsdWithError(w, 401, err.Error())
		return
	}
	user_id, err := auth.ValidateJwt(token, cfg.jwtKey)
	if err != nil {
		responsdWithError(w, 401, err.Error())
		return
	}
	pathValue := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(pathValue)
	if err != nil {
		responsdWithError(w, 400, "Invalid chirp_id format")
		return
	}
	chirp, err := cfg.dbQueries.GetChirpsById(r.Context(), chirpId)
	if err != nil {
		if err == sql.ErrNoRows {
			responsdWithError(w, 404, "Chirp not found")
			return
		}
		responsdWithError(w, 500, "An Unknow error occured")
		return
	}
	params := database.DeleteChirpByIdParams{
		ID:     chirp.ID,
		UserID: user_id,
	}
	result, err := cfg.dbQueries.DeleteChirpById(r.Context(), params)
	if err != nil {
		responsdWithError(w, 500, "Internal Server Error")
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		responsdWithError(w, 500, "Internal Server Error")
		return
	}
	if rowsAffected == 0 {
		responsdWithError(w, 403, "")
		return
	}
	resp := struct {
		MESSAGE string `json:"message"`
	}{
		MESSAGE: "chirp deleted successfully",
	}
	respondWithJSON(w, 204, resp)
}
