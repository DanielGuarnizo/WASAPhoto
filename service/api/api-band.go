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

	// Get userid from the path and hanlde error
	userid := ps.ByName("userid")
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	banisher, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(banisher, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in Validate")
		return
	}
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

	// Get userid and banished from the path and hanlde error
	userid := ps.ByName("userid")
	banished := ps.ByName("banished")
	if userid == "" || banished == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	banisher, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(banisher, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in validate")
		return
	}
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

	// Get userid and searchUsername from the path and hanlde error
	userid := ps.ByName("userid")
	searchUsername := r.URL.Query().Get("searchUsername")
	if userid == "" || searchUsername == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	username, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(username, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in Validate")
		return
	}
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
