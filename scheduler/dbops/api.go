package dbops

import "log"

func AddVideoDeletionRecord(vid string) error {
	//values 后面必须加上括号
	stmtIn, err := dbConn.Prepare("INSERT INTO video_del_rec(video_id) VALUES (?)")

	if err != nil{
		return err
	}
	defer stmtIn.Close()

	_,err = stmtIn.Exec(vid)
	if err != nil{
		log.Printf("AddVideoDeletionRecorde error:%v",err)
	}

	return nil
}