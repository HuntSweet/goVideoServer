package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	r := httprouter.New()

	r.GET("/",homeHandler)

	r.POST("/",homeHandler)

	r.GET("/userhome",userHomeHandler)

	r.POST("/userhome",userHomeHandler)

	r.POST("/api",apiHandler)

	r.POST("/upload/:vid-id",proxyHandler)

	r.ServeFiles("/statics/*filepath",http.Dir("./templates"))

	return r
}

func main()  {
	r := RegisterHandler()

	http.ListenAndServe(":8080",r)
}