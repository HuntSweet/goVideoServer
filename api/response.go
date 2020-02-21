package main

import (
	"encoding/json"
	"io"
	"net/http"
	"video_server/api/defs"
)

func sendErrorResponse(w http.ResponseWriter,errRep defs.ErrorResponse)  {
	w.WriteHeader(errRep.HttpSC)

	resStr,_ := json.Marshal(&errRep.Error)
	io.WriteString(w,string(resStr))
}

func sendNormalResponse(w http.ResponseWriter,resp string,sc int)  {
	w.WriteHeader(sc)
	io.WriteString(w,resp)

}
