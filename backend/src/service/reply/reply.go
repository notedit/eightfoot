// the reply Service
// date: 2012-06-02
// author: liulianxiang<notedit@gmail.com>

package reply

import (
    "fmt"
    "errors"
    "time"
    "strings"
    "database/sql"
)


type Reply  struct {
    DB *sql.DB
}


