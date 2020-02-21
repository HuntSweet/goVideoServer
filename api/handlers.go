package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"
	"video_server/api/utils"
)

func CreateUser(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	res,_ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}

	if err := json.Unmarshal(res,ubody);err != nil{
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}

	//写入数据库
	if err := dbops.AddUserCredential(ubody.Username,ubody.Pwd);err != nil{
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}

	//写入cache
	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUP{
		Success:   true,
		SessionId: id,
	}

	if resp,err := json.Marshal(su);err != nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w,string(resp),200)
	}

}

func Login(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	res,_ := ioutil.ReadAll(r.Body)
	log.Printf("%s",res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res,ubody);err != nil{
		log.Printf("%s",err)
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}

	//Validate the request body
	uname := p.ByName("username")
	log.Printf("Login url name :%s",uname)
	log.Printf("Login body name:%s",ubody.Username)
	if uname != ubody.Username{
		sendErrorResponse(w,defs.ErrorNotAuthUser)
		return
	}

	//validate pwd
	pwd,err := dbops.GetUserCredential(uname)
	log.Printf("Login pwd:%s",pwd)
	log.Printf("Login body pwd:%s",ubody.Pwd)
	if err != nil || len(pwd)==0 || pwd != ubody.Pwd{
		sendErrorResponse(w,defs.ErrorNotAuthUser)
		return
	}

	//generate session
	id := session.GenerateNewSessionId(uname)
	si := &defs.SignedIn{
		Success:   true,
		SessionId: id,
	}
	if resp,err := json.Marshal(si);err != nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w,string(resp),200)
	}

}

func GetUserInfo(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if !ValidateUser(w,r){
		log.Printf("Unauthorized user \n")
		return
	}

	uname := p.ByName("username")
	u,err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("Error in GetUserInfo:%s",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}

	ui := &defs.UserInfo{Id:u.Id}
	if resp,err := json.Marshal(ui);err != nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w,string(resp),200)
	}
}

func AddNewVideo(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if !ValidateUser(w,r){
		log.Printf("Unauthorized user \n")
		return
	}

	res,_ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res,nvbody);err != nil{
		log.Printf("%s",err)
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}

	vi,err := dbops.AddNewVideo(nvbody.AuthorId,nvbody.Name)
	if err != nil{
		log.Printf("AddNewVideo err:%v",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}

	if resp,err := json.Marshal(vi);err != nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	}else{
		sendNormalResponse(w,string(resp),201)
	}
}

func ListAllVideos(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if !ValidateUser(w,r){
		return
	}

	uname := p.ByName("username")
	vs,err := dbops.ListVideoInfo(uname,0,utils.GetCurrentTimestampSec())
	if err != nil{
		log.Printf("ListAllVideos error:%v",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{Videos:vs}

	if resp,err := json.Marshal(vsi);err != nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w,string(resp),200)
	}

}

func DeleteVideo(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if !ValidateUser(w,r){
		return
	}

	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil{
		log.Printf("DeleteVideo err :%v",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w,"DeleteVideo Success",204)
}

func PostComment(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if !ValidateUser(w,r){
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody,cbody);err != nil{
		log.Printf("%v",err)
		sendErrorResponse(w,defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := p.ByName("vid-id")
	if err := dbops.AddNewComment(vid,cbody.AuthorId,cbody.Content);err != nil{
		log.Printf("AddNewComment err:%v",err)
		sendErrorResponse(w,defs.ErrorDBError)
	} else {
		sendNormalResponse(w,"AddNewComment ok",201)
	}

}

func ShowComments(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	if !ValidateUser(w,r){
		return
	}

	vid := p.ByName("vid-id")
	cm,err := dbops.ListComments(vid,0,utils.GetCurrentTimestampSec())
	if err != nil{
		log.Printf("ListComments err:%v",err)
		sendErrorResponse(w,defs.ErrorDBError)
		return
	}

	cms := &defs.Comments{Comments:cm}
	if resp,err := json.Marshal(cms);err != nil{
		sendErrorResponse(w,defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w,string(resp),200)
	}

}