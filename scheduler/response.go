package main

import (
	"io"
	"net/http"
)

func sendResponse(w http.ResponseWriter,sc int,rsp string)  {
	w.WriteHeader(sc)
	io.WriteString(w,rsp)

}
