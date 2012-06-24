// date: 2012-05-18
// author: notedit<notedit@gmail.com>

package main

import (
    "flag"
    "log"
    "database/sql"
    "msgpack-oocrpc/rpc"
    "service/user"
    "service/content"
    "service/tag"
    "service/test"
    _ "github.com/bmizerany/pq"
)


var port *uint = flag.Uint("port",9090,"the port the server will listen on")
var host *string = flag.String("host","localhost","the host the server will listen on")


func main() {
    flag.Parse()
    db,err := sql.Open("postgres","sslmode=disable user=user port=5432 password=password dbname=database")
    if err != nil {
        log.Fatal("can not connecto to the postgres server:",err)
    }
    log.Println("backend server is starting...")
    server := rpc.NewServer(*host,*port)
    server.Register(&tag.Tag{db})
    server.Register(&content.Content{db})
    server.Register(&test.Test{db})
    server.Register(&user.User{db})
    log.Printf("backend listening %d\n",*port)
    server.Serv()

}
