package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

type Like struct{
	id: int 
	userId: int 
}
type JSONErrorMsg struct {
	Message string 
}

Likes = []Like{
	{
		id: 1
		userId: 1993238
	},{
		id:2
		userId: 1984033
	},
}


// add a like in the post of another user 
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ps are the parameter passed by the URL path, the package httprouter will retrieve those values in the path 

	// this specify the content-type the server will return to the client
	w.Header().Set("content-type", "application/json")

	// create a varible of type like in which we will parse the data passed in the request body 
	var like Like

	// read and parse the JSON data from the request body into a Like object. 
	// the .Decode method parse the data retrieve in the object memory address &like
	err := json.NewDecoder(r.Body).Decode(&like)

	if err != nil {
    	// Handle error (e.g., invalid JSON format)
    	http.Error(w, "Invalid JSON format", http.StatusBadRequest)
    	return
	}

	new_id := len(Likes)

	var new_like = Like{
		id: new_id
		userid: like.userId
	}

	Likes = append(Likes,new_like)

	
	// we use the json package to assoaciate to an encoder the HTTP response writer w ot it 
	// in such a way to return to the client the json object we specify early 
	err:= json.NewEncoder(w).Encode(new_like)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("like returnd an error in encodng error")
		w.WriteHeader(http.StatusInternalServerError)
		// here we can ignore the error because it comes from the writer so we don't care
		_ = json.NewDecoder(w).Encode(JSONErrorMsg{Message: "internal server error or bad request body"})
	}
}