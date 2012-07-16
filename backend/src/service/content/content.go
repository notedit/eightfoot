// the content Service
// date: 2012-06-02
// author: liulianxiang<notedit@gmail.com>

package content 

import (
    "fmt"
    "time"
    "errors"
    "database/sql"

    "utils"
)

var SQL_GET_ONE_CONTENT string = "SELECT id,title,author_ukey,last_modify_ukey,last_reply_ukey,body,url,atype,recommend_count,date_create,date_last_reply,show,disable_reply FROM content WHERE id=$1 AND show=true"
var SQL_ADD_ONE_CONTENT string = "INSERT INTO content (title,author_ukey,last_modify_ukey,last_reply_ukey,body,url,atype) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
var SQL_DEL_ONE_CONTENT string = "UPDATE content SET show=false WHERE id = $1"
var SQL_GET_CONTENT_COUNT string = "SELECT COUNT(id) FROM content WHERE show=false"
var SQL_LATEST_CONTENT string = "SELECT c.id,c.title,c.author_ukey,c.last_modify_ukey,c.last_reply_ukey,c.body,c.url,c.atype,c.recommend_count,c.date_create,c.date_last_reply,c.show,c.disable_reply FROM content AS c LEFT JOIN tag_map AS tm ON c.id = tm.content_id WHERE %s ORDER BY c.date_create DESC LIMIT $1 OFFSET $2"
var SQL_GET_FOLLOW_CONTENT_COUNT string = "SELECT count(c.id) FROM content AS c WHERE c.tag_id IN (SELECT ft.tag_id FROM follow_tag as ft WHERE ft.ukey =$1) ORDER BY c.date_create DESC"
var SQL_GET_FOLLOW_CONTENT string = "SELECT  FROM content AS c WHERE c.tag_id IN (SELECT ft.tag_id FROM follow_tag as ft WHERE ft.ukey =$1) ORDER BY c.date_create DESC LIMIT $2 OFFSET $3"
var SQL_GET_HOTEST_CONTENT string = "SELECT c.id,c.title,c.author_ukey,c.last_modify_ukey,c.last_reply_ukey,c.body,c.url,c.atype,c.recommend_count,c.date_create,c.date_last_reply,c.show,c.disable_reply FROM content AS c WHERE c.id IN (SELECT r.object_id FROM recommend AS r WHERE r.date_create > $1 GROUP BY r.object_id ORDER BY COUNT(r.object_id) LIMIT $2 OFFSET $3) "

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
    Url             string
    Atype           string
    RecommendCount  int
    DateCreate      time.Time
    DateLastReply   time.Time
    Show            bool
    DisableReply    bool
}


// 不提供tagid   
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
                    &content.Url,
                    &content.Atype,
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

func (c *Content)GetLatestContent(arg *LatestContentArg,cons *[]ContentItem)(err error){
    tagstr := "1=1 "
    if arg.TagId != 0 {
        tagstr = tagstr+fmt.Sprintf(" and tm.tag_id=%d ",arg.TagId)
    }
    fmt.Printf(SQL_LATEST_CONTENT,tagstr)
    rows,err := c.DB.Query(fmt.Sprintf(SQL_LATEST_CONTENT,tagstr),arg.Limit,arg.Offset)
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
            var conI ContentItem
            err = rows.Scan(&conI.Id,
                            &conI.Title,
                            &conI.AuthorUkey,
                            &conI.LastModifyUkey,
                            &conI.LastReplyUkey,
                            &conI.Body,
                            &conI.Url,
                            &conI.Atype,
                            &conI.RecommendCount,
                            &conI.DateCreate,
                            &conI.DateLastReply,
                            &conI.Show,
                            &conI.DisableReply)
            if err != nil {
                return err
            }
            *cons = append(*cons,conI)
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
    err = c.DB.QueryRow(SQL_ADD_ONE_CONTENT,content.Title,content.AuthorUkey,content.LastModifyUkey,
                        content.LastReplyUkey,content.Body,content.Url,content.Atype).Scan(cid)
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

// add @6.30
func (c *Content)GetFollowContentCount(ukey *string,followCount *int)(err error){
    if len(*ukey) != 6 {
        err = errors.New("ParamError: GetFollowContentCount ukey's length should be 6")
        return
    } 
    err = c.DB.QueryRow(SQL_GET_FOLLOW_CONTENT_COUNT,*ukey).Scan(followCount)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
    }
    return
}

// add @7.9 

func (c *Content)GetContentCount(_ *struct{},count *int)(err error){
    err = c.DB.QueryRow(SQL_GET_CONTENT_COUNT).Scan(count)
    if err != nil {
        err = utils.InternalError(err)
    }
    return
}

// add @7.9
type GetFollowContentArg struct{
    Ukey    string
    Offset  int
    Limit   int
}
func (c *Content)GetFollowContent(arg *GetFollowContentArg,cons *[]ContentItem)(err error) {
    if len(arg.Ukey) != 6 {
        err = errors.New("ParamError: GetFollowContent ukey's length should be 6")
        return
    }
    rows,err := c.DB.Query(SQL_GET_FOLLOW_CONTENT,arg.Ukey,arg.Limit,arg.Offset)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return err
    }
    for rows.Next() {
        var conI ContentItem
        err = rows.Scan(&conI.Id,
                        &conI.Title,
                        &conI.AuthorUkey,
                        &conI.LastModifyUkey,
                        &conI.LastReplyUkey,
                        &conI.Body,
                        &conI.Url,
                        &conI.Atype,
                        &conI.RecommendCount,
                        &conI.DateCreate,
                        &conI.DateLastReply,
                        &conI.Show,
                        &conI.DisableReply)
        if err != nil {
            err = errors.New("InternalError:" + err.Error())
            return err
        }
        *cons = append(*cons,conI)
    }
    return
}


// to do 获取最热的文章  todo
type  GetHotestContentArg struct
{
    Offset   int
    Limit    int
}

func (c *Content)GetHotestContent(arg *GetHotestContentArg,cons *[]ContentItem) (err error) {
    // 两天内推荐最多的
    duration := time.Duration(172800) * time.Second
    now := time.Now()
    twodays := now.Add(-duration)
    rows,err := c.DB.Query(SQL_GET_HOTEST_CONTENT,twodays,arg.Limit,arg.Offset)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return err
    }
    for rows.Next() {
        var conI ContentItem
        err = rows.Scan(&conI.Id,
                        &conI.Title,
                        &conI.AuthorUkey,
                        &conI.LastModifyUkey,
                        &conI.LastReplyUkey,
                        &conI.Body,
                        &conI.Url,
                        &conI.Atype,
                        &conI.RecommendCount,
                        &conI.DateCreate,
                        &conI.DateLastReply,
                        &conI.Show,
                        &conI.DisableReply)
        if err != nil {
            err = errors.New("InternalError:" + err.Error())
            return err
        }
        *cons = append(*cons,conI)
    }
    return
}



