package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"errors"
	"database/sql"
)
func (rt *_router) bandUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	userid:= ps.ByName("userid")
    // to handle the error if the request wa made in a wrong way
    if userid == "" {
        // Handle the case when "likeid" is not present in the request.
        rt.baseLogger.Warning("the userid in the path is empty")
        w.WriteHeader(http.StatusBadRequest)
        return 
    }
	// Authentication 
	authorized := Authentication(w,r,userid)
	if authorized == false{
		return 
	}

	var muteReq MuteRequest
	err := json.NewDecoder(r.Body).Decode(&muteReq)
	if err != nil {
        rt.baseLogger.WithError(err).Warning("invalid JSON format")
        w.WriteHeader(http.StatusBadRequest)
        _ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
        return
    }
    // defer closing the request body
    defer r.Body.Close()

	// Validate muted user ID
	if muteReq.MutedUserID == "" {
		rt.baseLogger.Warning("MutedUserID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Save into the database the muted relationship
	err = rt.db.BandUser(userid, muteReq.MutedUserID)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving muted user into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Respond with the muted user ID
	response := map[string]string{"mutedUserId": muteReq.MutedUserID}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(encodedResponse)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
}


func (rt *_router) unbandUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// load parameter from the path 
	userid := ps.ByName("userid")
	mutedUserid := ps.ByName("mutedUserid")

	// Check if userid or mutedUserid is empty
    if userid == "" || mutedUserid == "" {
        rt.baseLogger.Warning("the userid or mutedUserid in the path is empty")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	// Authentication 
	authorized := Authentication(w,r,userid)
	if authorized == false{
		return 
	}

    // Remove the user from the muted collection in the database
    err := rt.db.UnbandUser(userid, mutedUserid)
    if err != nil {
        // Check if the error is of the same type as sql.ErrNoRows
        if errors.Is(err, sql.ErrNoRows) {
            // Resource not found or not muted
            rt.baseLogger.WithError(err).Warning("user not found or not muted in the database")
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