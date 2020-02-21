package dbops

import ("database/sql"
         _ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err error
)

func init()  {
	dbConn,err = sql.Open("mysql","video_server:TspYktmMLPjPAL3K@/video_server?charset=utf8")
	if err != nil{
		panic(err.Error())
	}
}
