package main
//
//import (
//	"github.com/aliyun/aliyun-oss-go-sdk/oss"
//	"github.com/rs/zerolog/log"
//	"video_server/streamserver/config"
//
//)
//
//var EP string
//var AK string
//var SK string
//
//func init()  {
//	AK = ""
//	SK = ""
//	EP = config.GetOssAddr()
//}
////bn -> bucketname
//func UploadToOss(filename string,path string,bn string) bool {
//	client,err := oss.New(EP,AK,SK)
//	if err != nil{
//		log.Printf("Init oss service error :%s",err)
//		return false
//	}
//
//	bucket,err := client.Bucket(bn)
//	if err != nil{
//		log.Printf("Get bucket error :%s",err)
//		return false
//	}
//
//	err = bucket.UploadFile(filename,path,500*1024,oss.Routines(3))
//	if err != nil{
//		log.Printf("Bucket UploadFile error :%s",err)
//		return false
//	}
//
//	return true
//}