package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"strconv"
)
var Likes = []Like{
	Like{
		id:     0,
		UserID: 1993238,
	},
	Like{
		id:     1,
		UserID: 1984033,
	},
}

type Like struct {
	id int
	UserID int 
}


func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.parameters){
	likeid, err = strconv.Atoi(ps.Byname("likeid"))

	// to handle the error if the request wa made in a wrong way
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	// to handdle the error if the likeid is incorrect or not 
	if likeid < 0 || likeid >= len(Likes) {
		w.WriteHeader(StatusNotFound)
		return 
	}

	// if the likeid is correct then ir do the folloiwng
	Likes := append(Likes[:likeid], Likes[likeid + 1:]...)
}