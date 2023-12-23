package api 

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"math/rand"
	"time"
)

type Like struct {
	Post_ID   string `json:"post_id"`
	User_ID string `json:"user_id"` // Change from userId to UserID
}