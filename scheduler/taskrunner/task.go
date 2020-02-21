package taskrunner

import (
	"errors"
	"log"
	"os"
	"sync"
	"video_server/scheduler/dbops"
)

func deleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + vid)
	if err != nil && !os.IsNotExist(err){
		log.Printf("deleteVideo Err :%v",err)
		return err
	}
	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	//这里面读取3个
	res,err := dbops.ReadVideoDeletionRecord(3)
	if err != nil{
		log.Printf("Video clear disptcher error:%v",err)
		return err
	}

	if len(res) == 0{
		return errors.New("all tasks done")
	}

	for _,id := range res{
		dc <- id
	}

	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error

	forloop:
		for{
			select {
			case vid := <- dc:
				go func(id interface{}) {
					if err := deleteVideo(id.(string));err != nil{
						errMap.Store(id,err)
						return
					}
					if err := dbops.DelVideoDeletionRecord(id.(string));err != nil{
						errMap.Store(id,err)
						return
					}
				}(vid)
			default:
				break forloop
			}
		}

	errMap.Range(func(key, value interface{}) bool {
		err = value.(error)
		if err != nil{
			return false
		}
		return true
	})

	return err

}
