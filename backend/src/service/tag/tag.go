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
var SQL_ADD_ONE_TAG string = "INSERT INTO tag (name,introduction,author_ukey,url_code) VALUES ($1,$2,$3,$4)"
var SQL_DEL_ONE_TAG string = "UPDATE tag SET show=false WHERE id=$1"
var SQL_LATEST_UPDATE_TAG string = "SELECT t.id,t.name,t.introduction,t.date_create,t.content_count,t.follower_count,t.show,t.author_ukey,t.url_code from tag AS t WHERE t.id in (SELECT tm.tag_id from tag_map AS tm GROUP BY tm.tag_id ORDER BY t.date_create DESC LIMIT $1 OFFSET $2)"
var SQL_GET_CONTENT_TAG string = "SELECT t.id,t.name,t.introduction,t.date_create,t.content_count,t.follower_count,t.show,t.author_ukey,t.url_code FROM tag AS t LEFT JOIN tag_map AS tm ON t.id = tm.tag_id WHERE tm.content_id = $1"
var SQL_DEL_CONTENT_TAG string = "DELETE FROM tag_map WHERE content_id = $1"
var SQL_IS_TAG_EXIST string = "SELECT id FROM tag WHERE name=$1"
var SQL_SIMPLE_TAG string = "SELECT id,name FROM tag WHERE id in (%s)"
var SQL_ADD_ONE_TAG_RETURN_ID string = "INSERT INTO tag (name,introduction,author_ukey,url_code) VALUES ($1,$2,$3,$4) RETURNING id"
var SQL_INSERT_TAG_MAP string = "INSERT INTO tag_map (tag_id,content_id) VALUES %s"
var SQL_GET_TAG_LIST string = "SELECT id,name,introduction,date_create,content_count,follower_count,show,author_ukey,url_code FROM tag WHERE show=true ORDER BY follower_count DESC Limit $1 Offset $2"
var SQL_GET_TAG_COUNT string = "SELECT COUNT(id) FROM tag WHERE show=true"

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
    _,err = t.DB.Exec(SQL_ADD_ONE_TAG,tag.Name,tag.Introduction,tag.AuthorUkey,tag.UrlCode)
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
        err = errors.New("InternalError:"+err.Error())
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

// 获取最近有更新的标签
type LatestUpdateTagArg struct {
    Limit           int
    Offset          int
}
type LatestUpdateTagRep struct {
    Tag             []TagItem
}
func (t *Tag)GetLatestUpdateTag(arg *LatestUpdateTagArg,tags *[]TagItem)(err error){
    if arg.Limit <= 0 {
        err = errors.New("ParamError: the limit should > 0")
        return
    }
    rows,err := t.DB.Query(SQL_LATEST_UPDATE_TAG,arg.Limit,arg.Offset)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    for {
        if rows.Next() {
            var tag TagItem
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
                return err
            }
            *tags = append(*tags,tag)
        } else {
            break 
        }
    }
    if rows.Err() != nil {
        err = errors.New("InternalError:"+rows.Err().Error())
        return
    }
    return
}

// 获取内容的tag
type GetContentTagRep struct {
    Tag         []TagItem
}
func (t *Tag)GetContentTag(cid *int,tags *[]TagItem)(err error){
    rows,err := t.DB.Query(SQL_GET_CONTENT_TAG,*cid)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    for {
        if rows.Next() {
            var tag TagItem
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
                return err
            }
            *tags = append(*tags,tag)
        } else {
            break
        }
    }
    if rows.Err() != nil {
        err = errors.New("InternalError:"+rows.Err().Error())
        return
    }
    return
}

// 保存内容的tag
type SetContentTagArg struct {
    ContentId   int
    TagName     []string
}
func (t *Tag)SetContentTag(arg *SetContentTagArg,cid *int)(err error){
    if arg.ContentId <= 0 {
        err = errors.New("ParamError:contentid can not be 0")
        return
    }
    if len(arg.TagName) == 0 {
        *cid = arg.ContentId
        return
    }
    // 删除掉tag_map 中的记录
    _,err = t.DB.Exec(SQL_DEL_CONTENT_TAG,arg.ContentId)
    // 查看有没有已经注册过 注释掉的是一个优化方案(未完成)  先用最简单的方式实现
    //upPara := make(interface{},len(arg.TagName))
    //upSQL := make(string,len(arg.TagName))
    //for i,n := range arg.TagName {
    //    upSQL[i] = fmt.Sprintf("$%d",i)
    //    upPara[i] = n
    //    i += 1
    //}
    //sql := fmt.Sprintf(SQL_SIMPLE_TAG,strings.Join(upSQL,","))
    //r,err := t.DB.Exec(sql,upPara...)
    //if err != nil {
    //    err = errors.New("InternalError:"+err.Error())
    //    return
    //}
    tagids := make([]int,len(arg.TagName))
    for i,n := range arg.TagName {
        var tagid int
        // to do to do
        fmt.Println(SQL_IS_TAG_EXIST)
        err = t.DB.QueryRow(SQL_IS_TAG_EXIST,n).Scan(&tagid)
        if err != nil {
            // 这个标签不存在 添加他
            err = t.DB.QueryRow(SQL_ADD_ONE_TAG_RETURN_ID,n,"","","").Scan(&tagid)
            if err != nil {
                err = errors.New("InternalError:"+err.Error())
                return err
            }
            // 在一次判断tagid
            if tagid == 0 {
                err = errors.New("InternalError:insert a tag error")
                return err
            }
        }
        tagids[i] = tagid
    }
    // 插入 tag_map
    insertsql := make([]string,len(tagids))
    for i,n := range tagids {
        insertsql[i] = fmt.Sprintf("(%d,%d)",n,arg.ContentId)
    }
    sql := fmt.Sprintf(SQL_INSERT_TAG_MAP,strings.Join(insertsql,","))
    _,err = t.DB.Exec(sql)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    *cid = arg.ContentId
    return
}

// 获取tag 的数量

func (t *Tag)GetTagCount(_ *struct{},count *int)(err error){
    err = t.DB.QueryRow(SQL_GET_TAG_COUNT).Scan(count)
    if err != nil {
        err = utils.InternalError(err)
    }
    return
}


// 标签列表页面
type GetTagListArg struct{
    Limit     int
    Offset    int
}
func (t *Tag)GetTagList(arg *GetTagListArg,tags *[]TagItem)(err error){
    rows,err := t.DB.Query(SQL_GET_TAG_LIST,arg.Limit,arg.Offset)
    if err != nil {
        err = utils.InternalError(err)
        return
    }
    for rows.Next() {
        var tag  TagItem
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
                err = utils.InternalError(err)
                return err
            }
            *tags = append(*tags,tag)
    }
    return
}






