// Date: 2012-05-21
// Author: notedit<notedit@gmail.com>

package user

import (
    "log"
    "strings"
    "database/sql"
    "os/exec"
    "testing"
    _ "github.com/bmizerany/pq"
)


func TestSetup(t *testing.T) {
    cmd := exec.Command("/bin/sh","-c","cd ../../../ && ./initdb.sh")
    err := cmd.Run()
    if err != nil {
        t.Fatal("can not init the database, got error:"+err.Error())
    }
}

func SetupUser() *User {
    db,err := sql.Open("postgres","sslmode=disable user=user port=5432 password=password dbname=database")
    if err != nil {
        log.Fatal("can not connect to the postgre server:",err)
    }
    user := &User{db}
    return user
}


func TestGetUserInfo(t *testing.T) {
    user := SetupUser()

    ukey := "user01"
    iuser := UserItem{}
    err := user.GetUserInfo(&ukey,&iuser)
    if err != nil {
        t.Errorf("GetUserInfo Error: expect nil got %v",err)
    } else if iuser.Ukey != ukey {
        t.Error("GetUserInfo Error: iuser.Ukey should equal ukey")
    }
    
    ukey = "aaaaaa"
    iuser = UserItem{}
    err = user.GetUserInfo(&ukey,&iuser)
    if err == nil {
        t.Error("GetUserInfo Error: expect EmptyError got nil")
    } else if !strings.Contains(err.Error(),"EmptyError") {
        t.Errorf("GetUserInfo Error: expect EmptyError got %v",err)
    }
}

func TestGetMultiUserInfo(t *testing.T) {
    user := SetupUser()

    ukey_list := []string{"user01","user02","user03","user04"}
    mUser := MultiUser{}
    err := user.GetMultiUserInfo(&ukey_list,&mUser)
    if err != nil {
        t.Errorf("GetMultiUserinfo error: expect nil got %v",err)
    }
}

func TestcheckNicknameLower(t *testing.T) {
    user := SetupUser()

    nickname := "bbbbbbb"
    ok := user.checkNicknameLower(nickname)
    if !ok {
        t.Fatal("checkNicknameLower Error")
    }

    nickname = "user02"
    ok = user.checkNicknameLower(nickname)

    if ok {
        t.Fatal("checkNicknameLower Error")
    }
}

func TestRegisterUser(t *testing.T) {
    user := SetupUser()

    registerArg := RegisterUserArg{Nickname:"bbbbbbb",Email:"bcbbbbbb@gmail.com",Password:"pass01"}
    ukey := ""
    err := user.RegisterUser(&registerArg,&ukey)
    if err != nil {
        t.Error("RegisterUser Error:"+err.Error())
    }
    registerArg = RegisterUserArg{Nickname:"nick01",Email:"user01@gmail.com",Password:"pass01"}
    err = user.RegisterUser(&registerArg,&ukey)
    if err == nil {
        t.Error("RegisterUser Error: there should be a NickError")
    } else if !strings.Contains(err.Error(),"NicknameError") {
        t.Errorf("RegisterUser Error: there should be a NicknameErrror got %v",err)
    }

    registerArg = RegisterUserArg{Nickname:"邮箱重复",Email:"user02@gmail.com",Password:"pass01"}
    err = user.RegisterUser(&registerArg,&ukey)
    if err == nil {
        t.Error("RegisterUser Error: there should be an EmailError")
    } else if !strings.Contains(err.Error(),"EmailError") {
        t.Errorf("RegisterUser Error: there should be an EmailError got %v",err)
    }
}

func TestIsNicknameExist(t *testing.T) {
    user := SetupUser()

    nickname := "nick01"
    var exist bool
    user.IsNicknameExist(&nickname,&exist)
    if exist == false {
        t.Error("IsNicknameExist Error: nick01 exist, it should return true")
    }
    exist = true
    nickname = "fafafafafafa"
    user.IsNicknameExist(&nickname,&exist)
    if exist == true {
        t.Error("IsNicknameExist Error: fafafafa does not exist, it should return false")
    }
}

func TestLogin(t *testing.T) {
    user := SetupUser()

    loginarg := LoginArg{Email:"user05@gmail.com",Password:"pass05"}
    ukey := ""
    err := user.Login(&loginarg,&ukey)
    if err != nil {
        t.Error("Login Error: , got: "+err.Error())
    } else if ukey != "user05" {
        t.Error("Login Error: the ukey should be user05")
    }
    loginarg = LoginArg{Email:"aaa@gmail.com",Password:"aaa"}
    err = user.Login(&loginarg,&ukey)
    if err == nil {
        t.Error("Login Error: there should be an error,got nil")
    } else if !strings.Contains(err.Error(),"EmailError") {
        t.Errorf("Login Error: expect EmailError got %v",err)
    }
    loginarg = LoginArg{Email:"user01@gmail.com",Password:"user02"}
    err = user.Login(&loginarg,&ukey)
    if err == nil {
        t.Error("Login Error: there should be an error,got nil")
    } else if !strings.Contains(err.Error(),"PasswordError") {
        t.Errorf("Login Error: expect PasswordError got %v",err)
    }
}

func TestSetUserInfo(t *testing.T) {
    user := SetupUser()
    
    info := map[string]interface{}{"nickname":"用户01"}
    setArg := SetUserInfoArg{Ukey:"user99",Info:info}
    ukey := ""
    err := user.SetUserInfo(&setArg,&ukey)
    if err != nil {
        t.Error("SetUserInfo Error:"+err.Error())
    }

    info = map[string]interface{}{"nickname":"nick03"}
    setArg = SetUserInfoArg{Ukey:"user02",Info:info}
    err = user.SetUserInfo(&setArg,&ukey)
    if err == nil {
        t.Error("SetUserInfo Error: there should be an NicknameError")
    } else if !strings.Contains(err.Error(),"NicknameError") {
        t.Error("SetUserInfo Error: there should be an NicknameError" + err.Error())
    }
}


func TestUpdatePassword(t *testing.T) {
    user := SetupUser()
    uparg := UpdatePasswordArg{}
    uparg.Ukey = "user98"
    uparg.Password = "pas998"
    ukey := ""
    err := user.UpdatePassword(&uparg,&ukey)
    if err != nil {
        t.Error("UpdatePassword Error:"+err.Error())
    }
}
