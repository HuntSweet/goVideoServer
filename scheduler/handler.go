package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"video_server/scheduler/dbops"
)

func vidDelRecHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	vid := p.ByName("vid-id")
	if len(vid)==0{
		sendResponse(w,400,"video-id shouldn't be empty")
		return
	}

	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil{
		sendResponse(w,500,"Internal Error")
		return
	}

	sendResponse(w,200,"Delete success")
	return
}
