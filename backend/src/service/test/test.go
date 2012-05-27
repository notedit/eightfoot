// the test Service  just for test
// date: 2012-05-25
// author: notedit<notedit@gmail.com>

package test

import (
    "database/sql"
)

type Test struct {
    DB  *sql.DB
}

func (t *Test)GetHelloWorld(user *string,reply *string) (err error) {
    *reply = *user
    return
}
