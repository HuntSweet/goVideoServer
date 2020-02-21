package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter,r *http.Request,ps httprouter.Params)  {
	uname,err1 := r.Cookie("username")
	sid,err2 := r.Cookie("session")

	if err1 != nil && err2 != nil{
		//渲染名字
		p := &HomePage{Name:"sweetHunter"}
		t,err := template.ParseFiles("./templates/home.html")
		if err != nil{
			log.Fatalf("Parsing template failed:%v",err)
			return
		}

		t.Execute(w,p)
		return
	}

	//这里判断是否为登陆用户太简单了。。

	if len(uname.Value) != 0 && len(sid.Value) != 0{
		http.Redirect(w,r,"/userhome",http.StatusFound)
		return
	}

}

func userHomeHandler(w http.ResponseWriter,r *http.Request,ps httprouter.Params){
	//第一种情况:直接访问
	uname,err1 := r.Cookie("username")
	_,err2 := r.Cookie("session")

	if err1 != nil && err2 != nil{
		http.Redirect(w,r,"/",http.StatusFound)
	}

	//第二种情况：从homepage登陆过来
	fname := r.FormValue("username")

	var p *UserPage
	if len(uname.Value) != 0{
		p = &UserPage{Name:uname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name:fname}
	}

	t,e := template.ParseFiles("./templates/userhome.html")
	if e != nil{
		log.Printf("Parsing userhome.html failed:%v",e)
		return
	}

	t.Execute(w,p)
	return
}

func apiHandler(w http.ResponseWriter,r *http.Request,ps httprouter.Params){
	if r.Method != http.MethodPost{
		re,_ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w,string(re))
		return
	}

	res,_ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res,apibody);err != nil{
		re,_ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w,string(re))
		return
	}

	request(apibody,w,r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter,r *http.Request,ps httprouter.Params)  {
	//其实应该把要转发的地址写到配置文件里面
	re,_ := url.Parse("http://127.0.0.1:9000")
	proxy := httputil.NewSingleHostReverseProxy(re)

	proxy.ServeHTTP(w,r)
}