// Date: 2012-06-17
// Author: notedit<notedit@gmail.com>

package tag

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

func SetupTag() *Tag {
    db,err := sql.Open("postgres","sslmode=disable user=user port=5432 password=password dbname=database")
    if err != nil {
        log.Fatal("can not connect to the postgre server:",err)
    }
    tag := &Tag{db}
    return tag
}

func TestGetOneTag(t *testing.T) {
	tag := SetupTag()

	// normal 
	tagid := 1
	tagItem := TagItem{}
	err := tag.GetOneTag(&tagid,&tagItem)
	if err != nil {
		t.Errorf("GetOneTag: expect nil got %v",err)
	} else if tagItem.Id != tagid {
		t.Error("GetOneTag: can not get the correct TagItem")
	}

	// ParamError
	tagid = 0 
	tagItem = TagItem{}
	err = tag.GetOneTag(&tagid,&tagItem)
	if err == nil {
		t.Error("GetOneTag: expect ParamError got nil")
	} else if !strings.Contains(err.Error(),"ParamError") {
		t.Errorf("GetOneTag: expect ParamError got %v",err)
	}

	// EmptyError

	tagid = 1000
	tagItem = TagItem{}
	err = tag.GetOneTag(&tagid,&tagItem)
	if err == nil {
		t.Error("GetOneTag:expect EmptyError got nil")
	} else if !strings.Contains(err.Error(),"EmptyError") {
		t.Errorf("GetOneTag: expect EmptyError got %v",err)
	}
}

func TestAddOneTag(t *testing.T) {
	tag := SetupTag()

	// normal 
	var tagid int 
	tagItem := TagItem{Name:"NormalTag",Introduction:"Introduction",DateCreate:time.Now(),AuthorUkey:"user01",UrlCode:""}
	err := tag.AddOneTag(&tagItem,&tagid)
	if err != nil {
		t.Errorf("AddOneTag: expect nil got %v",err)
	} else if tagid <= 0 {
		t.Error("AddOneTag: can not AddOneTag")
	}

	// ParamError
	tagItem = TagItem{Name:"",Introduction:"Introduction",DateCreate:time.Now(),AuthorUkey:"user01",UrlCode:""}
	tagid = 0
	err = tag.AddOneTag(&tagItem,&tagid)
	if err == nil {
		t.Error("AddOneTag:expect ParamError got nil ")
	} else if !strings.Contains(err.Error(),"ParamError") {
		t.Errorf("AddOneTag:expect ParamError got %v",err)
	}

	// NameError 
	tagItem = TagItem{Name:"tagname01",Introduction:"Introduction",DateCreate:time.Now(),AuthorUkey:"user01",UrlCode:""}
	tagid = 0 
	err = tag.AddOneTag(&tagItem,&tagid)
	if err == nil {
		t.Error("AddOneTag: expect NameError got nil")
	} else if !strings.Contains(err.Error(),"NameError") {
		t.Errorf("AddOneTag: expect NameError got %v",err)
	}
}

func TestDelOneTag(t *testing.T) {
	tag := SetupTag()

	// normal 
	tagid := 50
	var tid int
	err := tag.DelOneTag(&tagid,&tid)
	if err != nil {
		t.Errorf("DelOneTag:expect nil got %v",err)
	} else if tagid != tid {
		t.Error("DelOneTag: can not del one tag")
	}

	// ParamError
	tagid = -1
	tid = 0
	err = tag.DelOneTag(&tagid,&tid)
	if err == nil {
		t.Errorf("DelOneTag: expect ParamError got %v",err)
	} else if !strings.Contains(err.Error(),"ParamError") {
		t.Errorf("DelOneTag: expect ParamError got %v",err)
	}
}

func TestGetLatestUpdateTag(t *testing.T) {
	tag := SetupTag()

	// normal 
	arg := LatestUpdateTagArg{Limit:2,Offset:0}
	rep := LatestUpdateTagRep{}
	err := tag.GetLatestUpdateTag(&arg,&rep)
	if err != nil {
		t.Errorf("GetLatestUpdateTag: expect nil got %v",err)
	} else if len(rep.Tag) != 2 {
		t.Error("GetLatestUpdateTag: can not get normal result")
	}

	// ParamError 
	arg = LatestUpdateTagArg{Limit:-1,Offset:0}
	rep = LatestUpdateTagRep{}
	err = tag.GetLatestUpdateTag(&arg,&rep)
	if err == nil {
		t.Errorf("GetLatestUpdateTag: expect ParamError got %v",err)
	} else if !strings.Contains(err.Error(),"ParamError") {
		t.Errorf("GetLatestUpdateTag: expect ParamError got %v",err)
	}
}

func TestGetContentTag(t *testing.T) {
	tag := SetupTag()

	//normal 
	cid := 2
	rep := GetContentTagRep{}
	err := tag.GetContentTag(&cid,&rep)
	if err != nil {
		t.Errorf("GetContentTag: expect nil got %v",err)
	} else if len(rep.Tag) != 3 {
		t.Error("GetContentTag: can not get correct result")
	}
}

func TestSetContentTag(t *testing.T) {
	tag := SetupTag()

	// normal
	tagname := []string{"tag02","tag03","tag09","tag中文"}
	arg := SetContentTagArg{ContentId:3,TagName:tagname}
	var cid int
	err := tag.SetContentTag(&arg,&cid)
	if err != nil {
		t.Errorf("SetContentTag: expect nil got %v",err)
	} else if cid != arg.ContentId {
		t.Log(cid)
		t.Log(arg.ContentId)
		t.Errorf("SetContentTag: can not set the correct tag")
	}

	// ParamError
	arg = SetContentTagArg{ContentId:-1,TagName:tagname}
	cid = 0
	err = tag.SetContentTag(&arg,&cid)
	if err == nil {
		t.Error("SetContentTag: expect ParamError got nil")
	} else if !strings.Contains(err.Error(),"ParamError") {
		t.Errorf("SetContentTag: expect ParamError got %v",err)
	}

	
}




