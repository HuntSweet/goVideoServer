package session

import (

	"log"
	"sync"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)



var sessionMap *sync.Map

func init()  {
	sessionMap = &sync.Map{}

}

func deleteExpiredSession(sid string)  {
	dbops.DeleteSession(sid)
	sessionMap.Delete(sid)
}

func nowInMill() int64 {
	return time.Now().UnixNano()/1000000
}
//写入cache
func LoadSessionFromDB()  {
	r,err := dbops.RetriveAllSession()
	if err != nil{
		log.Printf("LoadSessionFromDB err:%s",err)
		return
	}

	r.Range(func(k, v interface{}) bool {
		//这里为什么使用v.(*defs.SimpleSession)，而不是直接store（k,v)
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k,ss)
		return true
	})
}

func GenerateNewSessionId(un string) string {
	uid,err := utils.NewUUID()
	if err != nil{
		return ""
	}
	ct := time.Now().UnixNano()/1000000 //毫秒
	ttl := ct + 30*60*1000 //过期时间
	ss := &defs.SimpleSession{
		Username: un,
		TTL:      ttl,
	}
	sessionMap.Store(uid,ss)
	dbops.InserSession(uid,ttl,un)

	return uid
}

func IsSessionExpired(sid string) (string,bool) {
	ss,ok := sessionMap.Load(sid)
	ct := nowInMill()
	if ok{
		if ss.(*defs.SimpleSession).TTL < ct{
			//过期则删除
			deleteExpiredSession(sid)
			return "",true
		}
		return ss.(*defs.SimpleSession).Username,false
	} else {
		ss, err := dbops.RetriveSession(sid)
		if err != nil || ss == nil{
			return "",false
		}

		if ss.TTL < ct{
			deleteExpiredSession(sid)
			return "",false
		}

		sessionMap.Store(sid,ss)
		return ss.Username,false
	}

}