package dbops

import ("log"
	_ "github.com/go-sql-driver/mysql"
)

func ReadVideoDeletionRecord(count int) ([]string,error) {
	stmtOut,err := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")
	var ids []string
	if err != nil{
		return ids,err
	}
	defer stmtOut.Close()

	rows,err := stmtOut.Query(count)
	if err != nil{
		log.Printf("ReadVideoDeletionRecord err:%v",err)
		return ids,err
	}

	for rows.Next(){
		var id string
		if err := rows.Scan(&id);err != nil{
			return ids,err
		}
		ids = append(ids,id)
	}

	return ids,nil
}

func DelVideoDeletionRecord(vid string) error {
	stmtDel,err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id = ?")
	if err != nil{
		return err
	}
	defer stmtDel.Close()

	_,err = stmtDel.Exec(vid)
	if err != nil{
		log.Printf("DelVideoDeletionRecord err:%v",err)
		return err
	}

	return nil
}