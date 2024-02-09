package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get needed information to perfomr the operation
	userid := ps.ByName("userid")
	banisher, _ := rt.db.GetName(userid)
	id := r.Header.Get("Authorization")

	// Authentication
	is_valid, err := rt.db.Validate(banisher, id)
	if is_valid == false {
		w.WriteHeader(http.StatusUnauthorized)

		// You can include additional information in the response body if needed
		response := map[string]string{
			"error":   "UnauthorizedError",
			"message": "Authentication information is missing or invalid",
		}

		// Convert the response to JSON and write it to the response body
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// read and parse the json data from the request body into an username upadate object
	type Body struct {
		Banished string `json:"banished"`
	}
	var input Body
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Banished == "" {
		// Handle error (e.g., invalid JSON format)
		rt.baseLogger.WithError(err).Warning("invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return
	}

	// defer closing the request body
	defer r.Body.Close()

	// Save into the database the baned relationship
	err = rt.db.BandUser(banisher, input.Banished, userid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving baned user into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Respond with a 204 status code (success, no content)
	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get needed information to perfomr the operation
	userid := ps.ByName("userid")
	banisher, _ := rt.db.GetName(userid)
	id := r.Header.Get("Authorization")
	banished := ps.ByName("banished")

	// Authentication
	is_valid, err := rt.db.Validate(banisher, id)
	if is_valid == false {
		w.WriteHeader(http.StatusUnauthorized)

		// You can include additional information in the response body if needed
		response := map[string]string{
			"error":   "UnauthorizedError",
			"message": "Authentication information is missing or invalid",
		}

		// Convert the response to JSON and write it to the response body
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// Remove the user from the baned collection in the database
	err = rt.db.UnbandUser(banisher, banished)
	if err != nil {
		// Check if the error is of the same type as sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found or not muted
			rt.baseLogger.WithError(err).Warning("user not found or not baned in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "user not found or not muted in the database"})
			return
		}

		// Handle other errors
		rt.baseLogger.WithError(err).Warning("Error removing user from the muted collection in the database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Respond with a 204 status code (success, no content)
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getUserBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// rt.baseLogger.Warning("enter in getUserBans")

	// Get needed information to perfomr the operation
	userid := ps.ByName("userid")
	username, _ := rt.db.GetName(userid)
	id := r.Header.Get("Authorization")
	searchUsername := r.URL.Query().Get("searchUsername")

	// Authentication
	is_valid, err := rt.db.Validate(username, id)
	if is_valid == false {
		w.WriteHeader(http.StatusUnauthorized)

		// You can include additional information in the response body if needed
		response := map[string]string{
			"error":   "UnauthorizedError",
			"message": "Authentication information is missing or invalid",
		}

		// Convert the response to JSON and write it to the response body
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// Get the list of users fromdatabase
	UserList, err := rt.db.GetBans(searchUsername)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error retreaving user list of bans")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the list of users  in to a variable before writing it to the response writer
	response, err := json.Marshal(UserList)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding UserList for response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusOK)
	// rt.baseLogger.Warning(UserList)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
	// rt.baseLogger.Warning("if you see this is suppose the response was written correctly in the response body ")

}
