// Date: 2012-06-23
// Author: notedit<notedit@gmail.com>

package rpc

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Add(args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(args *Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}


func (t *Arith) NError(args *Args, reply *Reply) error {
	return errors.New("normalerror")
}

func (t *Arith)SimpleValue(arg *int,reply *bool) error {
    if *arg == 2 {
        *reply = true
    }
    return nil
}

func startServer() {
	newServer := NewServer("localhost", 9090)
	newServer.Register(new(Arith))
	newServer.Serv()
}

func TestServer(t *testing.T) {
	go startServer()
	client := New("localhost:9090")

	fmt.Println("string....")
	// normal calls
	args := &Args{7, 8}
	reply := new(Reply)
	err := client.Call("Arith.Add", args, reply)
	if err != nil {
		t.Errorf("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
	}

	// Nonexistent method
	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Arith.BadOperation", args, reply)
	if err == nil {
		t.Error("BadOperation: expected errpor")
	} else if !strings.Contains(err.Error(), "method") {
        t.Error("BadOperation: expected none exist method,got:", err.Error())
	}

	// normal error

	err = client.Call("Arith.NError", args, reply)
	if err == nil {
        t.Error("NEerror:expected normal error")
	} else if !strings.Contains(err.Error(), "normalerror") {
        t.Error("NError:expected an normal error, got ", err.Error())
	}

	// Unknown service
	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("Unknow.Arith", args, reply)
	if err == nil {
        t.Error("Unknow.Arith:expected Unknow service error")
	} else if !strings.Contains(err.Error(), "service") {
        t.Error("Unknow.Arith:expected Unknow service error: got ", err.Error())
	}

	// Error test
	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Arith.Div", args, reply)

	if err == nil {
		t.Error("Div: expected error")
	} else if !strings.Contains(err.Error(), "divide by") {
        t.Error("Div:expected divide by zero error detail:", err.Error())
	}

    // SimpleValue
    arg := 2
    var rep bool
    err = client.Call("Arith.SimpleValue",&arg,&rep)
    if err != nil {
        t.Error("SimpleValue:SimpleValue Error")
    } 
    if !rep {
        t.Error("SimpleValue:the rep should be true")
    }

}
