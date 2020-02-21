package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tempVedioId string
func clearTable()  {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M)  {
	clearTable()
	m.Run()
	clearTable()
}

func TestUserWorkFlow(t *testing.T)  {
	t.Run("AddUserCredential",testAddUserCredential)
	t.Run("GetUserCredential",testGetUserCredential)
	t.Run("DeleteUser",testDeleteUser)
	t.Run("RegetUser",testRegetUser)
	
}

func testAddUserCredential(t *testing.T) {
	err := AddUserCredential("aws","123")
	if err != nil{
		t.Errorf("AddUserCredential err:%v",err)
	}
	
}

func testGetUserCredential(t *testing.T) {
	pwd,err := GetUserCredential("aws")
	if pwd != "123" || err != nil{
		t.Error("GetUserCredential failed")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("aws","123")
	if err != nil{
		t.Errorf("DeleteUser err:%v",err)
	}
}

func testRegetUser(t *testing.T)  {
	pwd,err := GetUserCredential("aws")
	if err != nil{
		t.Error(err)
	}
	if pwd != "" {
		t.Error("delete user test failed")
	}
}

func TestVideoWorkFlow(t *testing.T)  {
	clearTable()
	t.Run("Prepare User",testAddUserCredential)
	t.Run("AddNewVideo",testAddNewVideo)
	t.Run("GetVideoInfo",testGetVideoInfo)
	t.Run("DeleteVideoInfo",testDeleteVideoInfo)
	t.Run("RegetVideoInfo",testRegetVideoInfo)
	//clearTable()
}

func testAddNewVideo(t *testing.T) {
	res,err := AddNewVideo(1,"test_video")
	if err != nil{
		t.Errorf("AddNewVideo err：%v",err)
	}
	tempVedioId = res.Id

}

func testGetVideoInfo(t *testing.T) {
	_,err := GetVideoInfo(tempVedioId)
	if err != nil{
		t.Errorf("GetVideoInfo err：%v",err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempVedioId)
	if err != nil{
		t.Errorf("DeleteVideoInfo err：%v",err)
	}
}

func testRegetVideoInfo(t *testing.T)  {
	res,err := GetVideoInfo(tempVedioId)
	if err != nil || res != nil{
		t.Errorf("RegetVideoInfo err：%v",err)
	}
}

func TestCommentWorkFlow(t *testing.T)  {
	clearTable()
	t.Run("AddUser",testAddUserCredential)
	t.Run("AddNewComment",testAddNewComment)
	t.Run("ListComments",testListComments)
}

func testAddNewComment(t *testing.T) {
	vid := "1"
	aid := 1
	content := "i love this video"
	err := AddNewComment(vid,aid,content)
	if err != nil{
		t.Errorf("AddNewComment err:%v",err)
	}

}

func testListComments(t *testing.T) {
	vid := "1"
	from := 1514764800
	to,_ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000,10))

	res,err := ListComments(vid,from,to)
	if err != nil{
		t.Errorf("ListComments err:%v",err)
	}

	for i,v := range res{
		fmt.Printf("comment:%d,%v \n",i,v)
	}
}
