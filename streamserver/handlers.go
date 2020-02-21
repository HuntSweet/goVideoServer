package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

)

//以oss播放，直接重定向到相应视频的网址即可
//func ossStreamServer(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
//	log.Println("enter ossStreamServer")
//	targetUrl := "https://aven-videos2.oss-cn-qingdao.aliyunsc.com/videos/" + p.ByName("vid-id")
//	http.Redirect(w,r,targetUrl,301)
//}

func streamHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid

	video,err := os.Open(vl)
	if err != nil{
		log.Printf("open file error:%v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"InternalServerError")
		return
	}

	//让浏览器以mp4格式播放
	w.Header().Set("Content-Type","video/mp4")
	http.ServeContent(w,r,"",time.Now(),video)

	defer video.Close()
	
}

func uploadHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	r.Body = http.MaxBytesReader(w,r.Body,MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil{
		sendErrorResponse(w,http.StatusBadRequest,"File is too big")
	}

	file,_,err := r.FormFile("file")
	if err != nil{
		sendErrorResponse(w,http.StatusInternalServerError,"fromfile error")
		return
	}

	data,err := ioutil.ReadAll(file)
	if err != nil{
		log.Printf("Read file error:%v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"read file error")
		return
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn,data,0666)
	if err != nil{
		log.Printf("WriteFile Error:%v",err)
		sendErrorResponse(w,http.StatusInternalServerError,"WriteFile Error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w,"Uploaded Scuccessfully!")
}

//上传到本地后，再上传到oss
//func ossUploadHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
//	r.Body = http.MaxBytesReader(w,r.Body,MAX_UPLOAD_SIZE)
//	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil{
//		sendErrorResponse(w,http.StatusBadRequest,"File is too big")
//	}
//
//	file,_,err := r.FormFile("file")
//	if err != nil{
//		sendErrorResponse(w,http.StatusInternalServerError,"fromfile error")
//		return
//	}
//
//	data,err := ioutil.ReadAll(file)
//	if err != nil{
//		log.Printf("Read file error:%v",err)
//		sendErrorResponse(w,http.StatusInternalServerError,"read file error")
//		return
//	}
//
//	fn := p.ByName("vid-id")
//	err = ioutil.WriteFile(VIDEO_DIR+fn,data,0666)
//	if err != nil{
//		log.Printf("WriteFile Error:%v",err)
//		sendErrorResponse(w,http.StatusInternalServerError,"WriteFile Error")
//		return
//	}
//
//	//上传oss
//	//oss的绝对路径
//	filename := "videos/" + fn
//	//文件的本地路径
//	path := VIDEO_DIR+fn
//	bucketName := "xxx"
//	if !UploadToOss(filename,path,bucketName){
//		log.Println("UploadToOss Failed")
//		sendErrorResponse(w,500,"UploadToOss Failed")
//		return
//	}
//	os.Remove(path)
//
//	w.WriteHeader(http.StatusCreated)
//	io.WriteString(w,"Uploaded Scuccessfully!")
//}

func testPageHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	t,_ := template.ParseFiles("./videos/upload.html")

	t.Execute(w,nil)
}