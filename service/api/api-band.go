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

	userid := ps.ByName("userid")
	banisher, err := rt.db.GetName(userid)
	if err != nil {
		rt.baseLogger.Warning("when we try to retrieve the name of a banisher by userid there was an error")
	}

	banished := ps.ByName("username")

	// to handle the error if the request wa made in a wrong way
	if banisher == "" || banished == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or username in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Authentication
	authorized := Authentication(w, r, userid)
	if authorized == false {
		return
	}

	// Save into the database the baned relationship
	err = rt.db.BandUser(banisher, banished, userid)
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

	userid := ps.ByName("userid")
	banisher, err := rt.db.GetName(userid)
	if err != nil {
		return
	}

	banished := ps.ByName("username")

	// to handle the error if the request wa made in a wrong way
	if banisher == "" || banished == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or username in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Authentication
	authorized := Authentication(w, r, userid)
	if authorized == false {
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
	rt.baseLogger.Warning("enter in getUserBans")
	// retrieve the esential information to perfomr the operation
	//userid := ps.ByName("userid")
	username := ps.ByName("username")

	// Authentication
	// authorized := Authentication(w, r, userid)
	// if authorized == false {
	// 	return
	// }

	// Get the list of users fromdatabase
	UserList, err := rt.db.GetBans(username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error retreaving user list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// UserList = "sto provando a salvare qualcosa"
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
