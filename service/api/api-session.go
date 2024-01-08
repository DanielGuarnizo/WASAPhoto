package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the JSON request body into a User object
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		rt.baseLogger.WithError(err).Warning("Invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user already exists in your system
	existingUser, err := rt.db.GetUserByName(user.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If the user doesn't exist, create a new user and return the identifier
	if existingUser.Username == "" {
		// Generate a new unique ID  for the user
		newUserid := generateUniqueID()

		err := rt.db.CreateUser(newUserid, user.Username)
		if err != nil {
			rt.baseLogger.WithError(err).Warning("Error creating user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return the new user identifier in the response
		response := map[string]string{"user_id": newUserid}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// If the user already exists, return the existing user identifier
	response := map[string]string{"user_id": existingUser.User_ID}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
