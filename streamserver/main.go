package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHanlder(r *httprouter.Router,cc int) http.Handler {
	m := &middleWareHandler{
		r: r,
		l:NewConnLimiter(cc),
	}
	return m
}

//这个中间件控制的是全局的连接，并不是video的连接
func (m *middleWareHandler) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	if !m.l.GetConn(){
		sendErrorResponse(w,http.StatusTooManyRequests,"Too many requests")
		return
	}

	m.r.ServeHTTP(w,r)
	defer m.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:vid-id",streamHandler)
	router.POST("/upload/:vid-id",uploadHandler)
	router.GET("/testpage",testPageHandler)

	return router
}

func main()  {
	r := RegisterHandlers()
	mh := NewMiddleWareHanlder(r,1000)
	http.ListenAndServe(":9000",mh)
}