// the user Service 
// date: 2012-05-18
// author: liulianxiang<notedit@gmail.com>

package user

import (
    "io"
    "fmt"
    "time"
    "errors"
    "strings"
    "crypto/sha1"
    "database/sql"
    
    "utils"
)


var SQL_GET_ONE_USERINFO string = "SELECT ukey,nickname,pic_small,pic_big,status,introduction,date_create,date_last_login FROM user_info WHERE ukey=$1"
var SQL_GET_MULTI_USERINFO string = "SELECT ukey,nickname,pic_small,pic_big,status,introduction,date_create,date_last_login FROM user_info WHERE ukey IN ($1)"
var SQL_REGISTER_USER string = "INSERT INTO passwd (ukey,email,salt,password) VALUES ($1,$2,$3,$4)"
var SQL_REGISTER_INSERT_NICKNAME string = "INSERT INTO user_info (ukey,nickname) VALUES ($1,$2)"
var SQL_CHECK_NICKNAME_LOWER string = "SELECT ukey,nickname FROM user_info WHERE lower(nickname) = $1"
var SQL_LOGIN_QUERY_UKEY string = "SELECT ukey,salt,password FROM passwd WHERE email=$1"
var SQL_UPDATE_USERINFO string = "UPDATE user_info SET %s  WHERE ukey='%s'"
var SQL_UPDATE_PASSWORD string = "UPDATE passwd SET salt=$1,password=$2 WHERE ukey=$3"

type User struct {
    DB *sql.DB
}

type UserItem struct {
    Ukey string
    Nickname string
    PicSmall    string  "omitempty"
    PicBig      string  "omitempty"
    Status      int 
    Introduction string
    DateCreate  time.Time
    DateLastLogin   time.Time
}

func (u *User)GetUserInfo(ukey *string,user *UserItem)(err error) {
    if len(*ukey) != 6 {
        err = errors.New("UkeyError:"+*ukey)
        return
    }
    rows,err := u.DB.Query(SQL_GET_ONE_USERINFO,*ukey)
    if err != nil {
        err = errors.New("EmptyError:" + err.Error())
        return
    }
    if !rows.Next() {
        err = errors.New("EmptyError")
        return
    }
    err = rows.Scan(&user.Ukey,
                    &user.Nickname,
                    &user.PicSmall,
                    &user.PicBig,
                    &user.Status,
                    &user.Introduction,
                    &user.DateCreate,
                    &user.DateLastLogin)
    if err != nil {
        return
    }
    return
}

type MultiUser struct {
    Users   []UserItem
}
func (u *User)GetMultiUserInfo(ukey_list *[]string,users *MultiUser)(err error) {
    if len(*ukey_list) == 0 {
        err = errors.New("the ukey_list should be list")
        return 
    }
    ukeystr := strings.Join(*ukey_list,",")
    rows,err := u.DB.Query(SQL_GET_MULTI_USERINFO,ukeystr)
    if err != nil {
        return
    }
    defer rows.Close()
    
    for {
        if rows.Next() {
            var user UserItem
            err = rows.Scan(&user.Ukey,
                            &user.Nickname,
                            &user.PicSmall,
                            &user.PicBig,
                            &user.Status,
                            &user.Introduction,
                            &user.DateCreate,
                            &user.DateLastLogin)
            if err != nil {
                return
            }
            users.Users = append(users.Users,user)
        } else {
            if rows.Err() != nil {
                err = rows.Err()
            }
            break
        }
    }
    return
}

// 注册一个用户
type RegisterUserArg struct {
    Nickname    string
    Email       string
    Password    string
}
func (u *User)RegisterUser(arg *RegisterUserArg,ukey *string)(err error){
    arg.Email = strings.Replace(arg.Email," ","",-1)
    arg.Nickname = strings.Replace(arg.Nickname," ","",-1)
    arg.Password = strings.Replace(arg.Password," ","",-1)
    
    if len(arg.Email) > 50 {
        err = errors.New("EmailError:too long")
        return
    }
    if len(arg.Nickname) >  50 {
        err = errors.New("NicknameError:too long")
        return
    }
    if len(arg.Password) > 40 {
        err = errors.New("PasswordError:too long")
        return
    }
    
    h := sha1.New()
    salt := utils.GenSalt(10)
    io.WriteString(h,salt+arg.Password) // salt
    arg.Password = fmt.Sprintf("%#x",string(h.Sum(nil)))

    if !u.checkNicknameLower(arg.Nickname) {
        err = errors.New("NicknameError:nickname dup")
        return 
    }

    for {
        *ukey = utils.GenUkey()
        r,err := u.DB.Exec(SQL_REGISTER_USER,*ukey,arg.Email,salt,arg.Password)
        if err != nil {
            if strings.Contains(err.Error(),"email_key") {
                err = errors.New("EmailError:"+err.Error())
            } else if strings.Contains(err.Error(),"ukey_key") {
                continue
            } else {
                err = errors.New("InternalError:"+err.Error())
            }
            return err
        }
        if n,err := r.RowsAffected(); n != 1 {
            err = errors.New("InternalError:insert user error")
            return err
        }
        break
    }
    _,err = u.DB.Exec(SQL_REGISTER_INSERT_NICKNAME,*ukey,arg.Nickname)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return
    }
    return
}

// true 的时候说明这个昵称是可以用的
// false 的时候说明这个昵称不是可用的
func (u *User)checkNicknameLower(nickname string) bool {
    rows,_:= u.DB.Query(SQL_CHECK_NICKNAME_LOWER,nickname)
    if !rows.Next() {
        return true
    }
    return false
}

type LoginArg struct{
    Email       string
    Password    string
}
func (u *User)Login(arg *LoginArg,ukey *string)(err error) {
    // for login
    arg.Email = strings.Replace(arg.Email," ","",-1)
    arg.Password = strings.Replace(arg.Password," ","",-1)

    rows,err := u.DB.Query(SQL_LOGIN_QUERY_UKEY,arg.Email)
    if err != nil {
        return
    }
    if !rows.Next() {
        if rows.Err() != nil {
            return rows.Err()
        }
        err = errors.New("EmailError: email does not existe")
        return 
    }
    var salt,password string
    err = rows.Scan(ukey,&salt,&password)
    if err != nil {
        err = errors.New("InternalError:" + err.Error())
        return
    }
    h := sha1.New()
    io.WriteString(h,salt+arg.Password)
    pass := fmt.Sprintf("%#x",string(h.Sum(nil)))
    if pass != password {
        err = errors.New("PasswordError: password is unvalid")
        return
    }
    return 
}

// true: exist  false: do not exist
func (u *User)IsNicknameExist(nickname *string,exist *bool)(err error) {
    rows,err := u.DB.Query(SQL_CHECK_NICKNAME_LOWER,nickname)
    if !rows.Next() {
        *exist = false
        return 
    }
    *exist = true
    return
}


// modify the user_info
type SetUserInfoArg struct {
    Ukey    string
    Info    map[string]interface{}
}
func (u *User)SetUserInfo(arg *SetUserInfoArg,ukey *string)(err error){
    if len(arg.Ukey) != 6 {
        err = errors.New("UkeyError:ukey is unvalid")
        return err
    }
    fields := []string{"nickname","pic_small","pic_big","status","introduction","date_last_login"}
    for field,_ := range arg.Info {
        ok := false
        for _,name := range fields {
            if field == name {
                ok = true
                break
            }
        }
        if !ok {
            err = errors.New("FieldError:" + field)
            return
        }
    }
    fmt.Println(arg.Ukey)
    var i int = 1
    var upPara []interface{}
    var upSQL []string

    for field,value := range arg.Info {
        if field == "nickname" {
            fmt.Println(value.(string))
            if !u.checkNicknameLower(value.(string)) {
                err = errors.New("NicknameError:nickname exist")
                return
            }
        }
        upSQL = append(upSQL,fmt.Sprintf("%s=$%d",field,i))
        upPara = append(upPara,value)
        i += 1
    }
    
    upstr := strings.Join(upSQL,",")
    sql := fmt.Sprintf(SQL_UPDATE_USERINFO,upstr,arg.Ukey)
    fmt.Println(sql)
    fmt.Printf("%#v\n",upPara)
    r,err := u.DB.Exec(sql,upPara...)
    if err != nil {
        fmt.Println(err.Error())
        err = errors.New("InternalError:"+err.Error())
        return
    }
    fmt.Println("affected:")
    if n,err := r.RowsAffected(); n != 1 {
        err = errors.New("InternalError:update user error")
        return err
    }
    return
}

type UpdatePasswordArg struct {
    Ukey        string
    Password    string
}
func (u *User)UpdatePassword(arg *UpdatePasswordArg,ukey *string) (err error){
    if len(arg.Ukey) != 6 {
        err = errors.New("UkeyError:"+arg.Ukey)
        return
    }
    if len(arg.Password) == 0 {
        err = errors.New("PasswordError: unvalid password")
        return
    }
    salt := utils.GenSalt(10)
    h := sha1.New()
    io.WriteString(h,salt+arg.Password)
    pass := fmt.Sprintf("%#x",string(h.Sum(nil)))
    _,err = u.DB.Exec(SQL_UPDATE_PASSWORD,salt,pass,arg.Ukey)
    if err != nil {
        err = errors.New("InternalError:"+err.Error())
        return 
    }
    *ukey = arg.Ukey
    return
}

// 生成一个安全的token, 用于验证邮箱 或者找回密码的验证
func (u *User)GetSecurityToken(ukey *string,stoken *string)(err error){
    //todo  暂时不做这么复杂
    return nil
}

func (u *User)QuerySecurityToken(token *string,ukey *string)(err error){
    // todo  同上
    return nil
}


