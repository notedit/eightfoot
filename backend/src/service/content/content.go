// the content Service
// date: 2012-06-02
// author: liulianxiang<notedit@gmail.com>

package content 

import (
    "fmt"
    "time"
    "errors"
    "strings"
    "database/sql"
)

type Content struct {
    DB *sql.DB
}

type ContentItem struct {
    Id              int
    Title           string
    AuthorUkey      string
    LastModifyUkey  string
    LastReplyUkey   string
    Body            string
    RecommendCount  int
    DateCreate      time.Time
    DateLastReply   time.Time
    Show            bool
    DisableReply    bool
    TagId           int
}

func (c *Content)GetOneContent(cid *int,content *ContentItem)(err error){
    if *cid <= 0 {
        err =fmt.Errorf("ParamError:contentid is %v",*cid)
        return
    }
    rows,err := c.DB.Query(SQL_GET_ONE_CONTENT,*cid)
    if err != nil {
        err = errors.New("InternalError:" + err.Error())
        return
    }
    if !rows.Next() {
        err = fmt.Errorf("EmptyError: content %v does not exist",*cid)
        return 
    }
    err = rows.Scan(&content.Id,
                    &content.Title,
                    &content.AuthorUkey,
                    &content.LastModifyUkey,
                    &content.LastReplyUkey,
                    &content.Body,
                    &content.RecommendCount,
                    &content.DateCreate,
                    &content.DateLastReply,
                    &content.Show,
                    &content.DisableReply,
                    &content.TagId)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    return
}

func (c *Content)AddOneContent(content *ContentItem,cid *int)(err error){
    if len(content.Title) == 0 {
        err = errors.New("ParamError:content title should not be empty")
        return
    }
    r,err := c.DB.Exec(SQL_ADD_ONE_CONTENT,content.Title,content,AuthorUkey,content.LastModifyUkey,
                        content.LastReplyUkey,content.Body,content.TagId)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    err = c.DB.QueryRow("SELECT lastval()").Scan(cid)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    return
}

func (c *Content)DelOneContent(cid *int,rcid *int)(err error){
    if *cid <= 0 {
        err = errors.New("ParamError: the content id should be > 0 ")
        return
    }
    _,err = c.DB.Exec(SQL_DEL_ONE_CONTENT,*cid)
    if err != nil {
        return
    }
    *rcid = *cid
    return
}


