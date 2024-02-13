package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	// rt.baseLogger.Warning("Enter in the uploadPhoto function")

	// Get userid from the path and hanlde error
	userid := ps.ByName("userid")
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	name, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(name, id)
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

	// create a variable post to parse the data passed in the request body

	var p Post

	// Read and parse the JSON data from the request body into a Post object.
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return
	}
	// rt.baseLogger.Warning("Raw Request Body:", p)

	// defer closing the request body
	defer r.Body.Close()
	// Generate a new unique ID for the like
	postid := generateUniqueID()

	// Save the image to the file system
	imagePath, err := saveImageToFileSystem(postid, p.Image)
	rt.baseLogger.Warning(imagePath)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving image to file system")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	newpost := Post{
		User_ID:  userid,
		Post_ID:  postid,
		Uploaded: p.Uploaded,
		Image:    imagePath,
	}

	// Save the new post in the database
	err = rt.db.UploadPhoto(newpost.PostToDataBase())
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving post into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Encode the newpost to a variable before writing it to the response writer
	response, err := json.Marshal(newpost)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding post for response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get userid and postid from the path and hanlde error
	userid := ps.ByName("userid")
	postid := ps.ByName("postid")
	if userid == "" || postid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	name, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(name, id)
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

	// Remove from the database the post given the postid
	err = rt.db.DeletePhoto(postid)
	if err != nil {
		// check it the err is of the same type of sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("Post not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Post not found in the database"})
			return
		}

		// Handle other errors
		rt.baseLogger.WithError(err).Warning("Error removing post from the database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
