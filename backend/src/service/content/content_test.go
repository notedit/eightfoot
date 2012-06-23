// Date:2012-06-17
// Author: notedit<notedit@gmail.com>

package content

import (
	"log"
	"strings"
	"database/sql"
	"os/exec"
	"testing"
    "time"
	_ "github.com/bmizerany/pq"
)

func TestSetup(t *testing.T) {
    cmd := exec.Command("/bin/sh","-c","cd ../../../ && ./initdb.sh")
    err := cmd.Run()
    if err != nil {
        t.Fatal("can not init the database, got error:"+err.Error())
    }
}

func SetupContent() *Content {
	db,err := sql.Open("postgres","sslmode=disable user=user port=5432 password=password dbname=database")
    if err != nil {
        log.Fatal("can not connect to the postgre server:",err)
    }
    content := &Content{db}
    return content
}

func TestGetOneContent(t *testing.T) {
	content := SetupContent()

	// first the normal result
	content_id := 1
	content_item := ContentItem{}
	err := content.GetOneContent(&content_id,&content_item)
	if err != nil {
		t.Errorf("GetOneContent: expect nil got %v",err)
	} else if content_item.Id != 1 || content_item.AuthorUkey != "user01" {
		t.Error("GetOneContent: contentItem's id  should be 1 ")
	}

	// paramerror 
	content_id = -1 
	content_item = ContentItem{}
	err = content.GetOneContent(&content_id,&content_item)
	if err == nil {
		t.Error("GetOneContent: expect paramerror got nil")
	} else if !strings.Contains(err.Error(),"ParamError") {
		t.Errorf("GetOneContent: expect ParamError got %v",err)
	}

	// EmptyError 
	content_id = 1000 
	content_item = ContentItem{}
	err = content.GetOneContent(&content_id,&content_item)
	if err == nil {
		t.Error("GetOneContent: expect EmptyError got nil ")
	} else if !strings.Contains(err.Error(),"EmptyError") {
		t.Errorf("GetOneContent: expect EmptyError got %v",err)
	}
}

func TestGetLatestContent(t *testing.T) {
	content := SetupContent()

	// no Tagid
	arg := LatestContentArg{Offset:0,Limit:1}
	rep := LatestContentRep{}
	err := content.GetLatestContent(&arg,&rep)
	if err != nil {
		t.Errorf("GetLatestContent:expect nil  got error %v",err)
	} else if len(rep.Content) != 1 {
		t.Errorf("GetLatestContent: can not get contentItem,  %v",rep.Content)
	} 
	// have Tagid
	arg = LatestContentArg{Offset:0,Limit:1,TagId:1}
	rep = LatestContentRep{}
	err = content.GetLatestContent(&arg,&rep)
	if err != nil {
		t.Errorf("GetLatestContent: expect nil got error %v",err)
	} else if len(rep.Content) != 1 {
		t.Error("GetLatestContent: can not get correct contentItem")
	}

}

func TestAddOneContent(t *testing.T) {
	content := SetupContent()

	// normal
	var cid int 
	content_item := ContentItem{Title:"TestAddOneContent",
								AuthorUkey:"user01",
								LastModifyUkey:"user01",
								LastReplyUkey:"user01",
								Body:"this is a body",
								Url:"http://www.baidu.com",
								Atype:"content",
								DateCreate:time.Now()}
	err := content.AddOneContent(&content_item,&cid)
	if err != nil {
		t.Errorf("AddOneContent:expect nil  got %v",err)
	} else if cid <= 100 {
		t.Error("AddOneContent: the cid should > 0 ")
	}
}

func TestDelOneContent(t *testing.T) {
	content := SetupContent()

	cid := 1
	var rcid int
	err := content.DelOneContent(&cid,&rcid)
	if err != nil {
		t.Errorf("DelOneContent: expect nil got %v",err)
	} else if cid != rcid {
		t.Error("DelOneContent: cid should be equail with the rcid")
	} 
}
