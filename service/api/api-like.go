package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type JSONErrorMsg struct {
	Message string `json:"message"`
}

// add a like in the post of another user
// ps are the parameters passed by the URL path, the package httprouter will retrieve those values in the path
// ctx is a contect object, is a way of passing data accross API boundaries
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Specify the content-type the server will return to the client
	w.Header().Set("Content-Type", "application/json")

	// Get useri and postid from the path and hanlde error
	userid := ps.ByName("userid")
	postid := ps.ByName("postid")
	if userid == "" || postid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or postid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	liker, err := rt.db.GetName(userid) // name of the userid
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(liker, id)
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

	// Save the new like to the database
	err = rt.db.SetLike(postid, liker, userid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving like into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Respond with a 204 status code (success, no content)
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get userid, likerUsername and postid from the path and hanlde error
	userid := ps.ByName("userid")
	postid := ps.ByName("postid")
	likerUsername := ps.ByName("likerUsername")
	if userid == "" || postid == "" || likerUsername == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid, likerUsername or postid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(likerUsername, id)
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

	// Remove from the database the like given the likeid
	err = rt.db.RemoveLike(userid, postid)
	if err != nil {
		// check it the err is of the same type of sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("like not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "like not found in the database"})
			return
		}

		// Handle other errors
		rt.baseLogger.WithError(err).Warning("Error removing like from the database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt *_router) getLikers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get userid and postid from the path and hanlde error
	userid := ps.ByName("userid")
	postid := ps.ByName("postid")
	if postid == "" || userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the likeid in the path is empty")
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

	// Get the list of likers of a post from the database
	UserList, err := rt.db.GetLikers(postid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error retreaving likers list of a post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the list of likers  in to a variable before writing it to the response writer
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

}
