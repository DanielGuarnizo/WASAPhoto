package api

import (
	"net/http"
)




// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes, this is done to each requet we specify in the API file 
	
	rt.router.POST("/likes/", rt.likePhoto)







	//rt.router.GET("/", rt.getHelloWorld)
	//rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	//rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
