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

    "service/tag"
)

var SQL_GET_ONE_CONTENT string = "SELECT id,title,author_ukey,last_modify_ukey,last_reply_ukey,body,recommend_count,date_create,date_last_reply,show,disable_reply FROM content WHERE id=$1 AND show=true"
var SQL_ADD_ONE_CONTENT string = "INSERT INTO content (title,author_ukey,last_modify_ukey,last_reply_ukey,body) VALUES ($1,$2,$3,$4,$5)"
var SQL_ADD_TAG_CONTENT string = ""
var SQL_DEL_ONE_CONTENT string = "UPDATE content SET show=false WHERE id = $1"
var SQL_LATEST_CONTENT string = "SELECT id,title,author_ukey,last_modify_ukey,last_reply_ukey,body,recommend_count,date_create,date_last_reply,show,disable_reply FROM content WHERE %s ORDER BY date_create DESC LIMIT $1 OFFSET $2"


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
    TagId           []int
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
                    &content.DisableReply)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    return
}

// 获取最新的主题
type LatestContentArg struct {
    Offset          int
    Limit           int
    TagId           int // 可以为0 为0的话取全部的
}
type LaestContentRep struct {
    Content     []ContentItem
}
func (c *Content)GetLatestContent(arg *LatestContentArg,rep *LatestContentRep)(err error){
    tagstr := "1=1 "
    if arg.TagId != 0 {
        tagstr = tagstr+fmt.Sprintf("tag_id=%d ",arg.TagId)
    }
    rows,err := c.DB.Query(fmt.Sprintf(SQL_LATEST_CONTENT,tagstr),arg.Offset,arg.Limit)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    if rows.Err() != nil {
        err = errors.New("InternalError:"+rows.Err().Error())
        return
    }
    for {
        if rows.Next() {
            var conI ContentItm
            err = rows.Scan(&conI.Id,
                            &conI.Title,
                            &conI.AuthorUkey,
                            &conI.LastModifyUkey,
                            &conI.LastReplyUkey,
                            &conI.Body,
                            &conI.RecommendCount,
                            &conI.DateCreate,
                            &conI.DateLastReply,
                            &conI.Show,
                            &conI.DisableReply,
                            &conI.TagId)
            if err != nil {
                return err
            }
            rep.Content.append(rep.Content,conI)
        } else {
            break
        }
    }
    return
}

func (c *Content)AddOneContent(content *ContentItem,cid *int)(err error){
    if len(content.Title) == 0 {
        err = errors.New("ParamError:content title should not be empty")
        return
    }
    r,err := c.DB.Exec(SQL_ADD_ONE_CONTENT,content.Title,content.AuthorUkey,content.LastModifyUkey,
                        content.LastReplyUkey,content.Body)
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



