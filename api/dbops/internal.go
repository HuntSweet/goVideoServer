package dbops

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"strconv"
	"sync"
	"video_server/api/defs"
)

func InserSession(sid string,ttl int64,uname string) error {
	ttlstr := strconv.FormatInt(ttl,10)
	stmtIns, err := dbConn.Prepare("INSERT INTO sessions(session_id,TTL,login_name) VALUES (?,?,?)")
	if err != nil{
		return err
	}
	defer stmtIns.Close()

	_,err = stmtIns.Exec(sid,ttlstr,uname)
	if err != nil{
		return err
	}

	return nil
}

func RetriveSession(sid string) (*defs.SimpleSession,error) {
	ss := &defs.SimpleSession{}
	stmtOut,err := dbConn.Prepare("SELECT TTL,login_name FROM sessions WHERE session_id = ?")
	if err != nil{
		return nil,err
	}
	defer stmtOut.Close()
	
	var ttl string
	var uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl,&uname)
	if err != nil && err != sql.ErrNoRows{
		return nil,err
	}
	
	if res,err := strconv.ParseInt(ttl,10,64);err == nil{
		ss.TTL = res
		ss.Username = uname
	} else {
		return nil,err
	}
	
	return ss,nil
	
}

func RetriveAllSession() (*sync.Map,error) {
	m := &sync.Map{}
	stmtOut,err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil{
		log.Printf("%s",err)
		return nil,err
	}
	defer stmtOut.Close()

	rows,err := stmtOut.Query()
	if err != nil{
		log.Printf("%s",err)
		return nil,err
	}
	for rows.Next(){
		var id string
		var ttlstr string
		var login_name string
		if err := rows.Scan(&id,&ttlstr,&login_name);err != nil{
			log.Printf("retrive session error:%s",err)
			return nil,err
		}

		if ttl,err := strconv.ParseInt(ttlstr,10,64);err == nil{
			ss := &defs.SimpleSession{Username:login_name,TTL:ttl}
			m.Store(id,ss)
			log.Printf("session id:%s,ttl:%d",id,ss.TTL)
		}
	}

	return m,nil

}

func DeleteSession(sid string) error {
	stmtDel,err := dbConn.Prepare("DELETE FROM sessions WHERE id = ?")
	if err != nil{
		return err
	}
	defer stmtDel.Close()

	_,err = stmtDel.Query(sid)
	if err != nil{
		return err
	}

	return nil
}