// the tag Service
// date: 2012-05-31
// author: liulianxiang<notedit@gmail.com>

package tag

import (
    "fmt"
    "errors"
    "strings"
    "time"
    "database/sql"

    "utils"
)

var SQL_GET_ONE_TAG string = "SELECT id,name,introduction,date_create,content_count,follower_count,show,author_ukey,url_code from tag where id = $1 and show=true"


type Tag struct {
    DB *sql.DB
}

type TagItem struct {
    Id      int
    Name    string
    Introduction string
    DateCreate  time.Time 
    ContentCount    int
    FollowerCount   int
    Show            bool
    AuthorUkey      string
    UrlCode         string
}


func (t *Tag)GetOneTag(tagid *int,tag *TagItem)(err error){
    if *tagid <= 0 {
        err = fmt.Errorf("ParamError:tagid id is %v",*tagid)
        return
    }
    rows,err := t.DB.Query(SQL_GET_ONE_TAG,*tagid)
    if err != nil {
        err = errors.New("InternalError:" + err.Error())
        return
    }
    if !rows.Next() {
        err = fmt.Errorf("EmptyError:tag %v does not exist",*tagid)
        return
    }
    err = rows.Scan(&tag.Id,
                    &tag.Name,
                    &tag.Introduction,
                    &tag.DateCreate,
                    &tag.ContentCount,
                    &tag.FollowerCount,
                    &tag.Show,
                    &tag.AuthorUkey,
                    &tag.UrlCode)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    return
}


func (t *Tag)AddOneTag(tag *TagItem,tagid *int)(err error){
    if len(tag.Name) == 0 {
        err = errors.New("ParamError: tag name should not be empty")
        return
    }
    r,err := t.DB.Exec(SQL_ADD_ONE_TAG,tag.Name,tag.Introduction,tag.AuthorUkey,tag.UrlCode)
    if err != nil {
        if strings.Contains(err.Error(),"name_key") {
            err = fmt.Errorf("NameError: name %v dup",tag.Name)
            return err
        }
        err = errors.New("InternalError:" + err.Error())
        return
    }
    err = t.DB.QueryRow("SELECT lastval()").Scan(tagid)
    if err != nil {
        err = errors.New("InternaleRROR:"+err.Error())
        return
    }
    return
}

func (t *Tag)DelOneTag(tagid *int,tid *int)(err error){
    if *tagid <= 0 {
        err = errors.New("ParamError:the tagid should be > 0")
        return
    }
    _,err = t.DB.Exec(SQL_DEL_ONE_TAG,*tagid)
    if err != nil {
        return 
    }
    *tid = *tagid
    return
}



