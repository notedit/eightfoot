// the utils pkg
// author: notedit<notedit@gmail.com>
// date: 20120525

package utils

import (
    "strings"
    "strconv"
    "math/rand"
    "encoding/hex"
)


// generate ukey 
func GenUkey()(ukey string) {
    randint := int(rand.Int31())
    return Int2Ukey(randint)
}

func Int2Ukey(randint int)(ukey string){
    bs := make([]byte,6)
    var i,mod int
    for i,_ = range bs {
        randint,mod = randint/36,randint%36
        if mod < 10 {
            bs[5-i] = byte(mod + 48)
        } else {
            bs[5-i] = byte(mod + 97 - 10)
        }
    }
    return string(bs)
}

// generate salt
func GenSalt(n int) (salt string) {
    bs := make([]byte,n)
    var i int
    for i,_ = range bs {
        bs[i] = byte(rand.Intn(26)+97)
    }
    return string(bs)
}

// join []int
func JoinInt(ints []int,sep string) string {
    t := make([]string,len(ints))
    for i := range ints {
        t[i] = strconv.Itoa(ints[i])
    }
    return strings.Join(t,sep)
}


// http://play.golang.org/p/pKPEeNWsjD
// 用来escape 特殊字符  "Hello, '世界'" => E'Hello, \'\u4e16\u754c\''
func HexBuffer1(input string) string {
    output := make([]byte,0,3+len(input))
    output = append(output,"E'"...)
    for _,r := range input {
        var s string
        switch r {
        case '\'':
            s = `\'`
        case '\\':
            s = `\\`
        default:
            s = strconv.QuoteRuneToASCII(r)
            // get rid of surrounding single quote
            s = s[1:len(s)-1]
        }
        output = append(output,s...)
    }
    return string(append(output,'\''))
}

// 用来escape 特殊字符串  
// "Hello, '世界'" => E'\x48\x65\x6c\x6c\x6f\x2c\x20\x27\xe4\xb8\x96\xe7\x95\x8c\x27'
func HexBuffer2(input string) string {
    L := len(input)
    payload_size := 4 * L
    output := make([]byte,payload_size+3)
    copy(output,"E'")
    output[len(output)-1] = '\''
    payload := output[2 : len(output)-1]
    hexdump := payload[len(payload)-2*L:]
    hex.Encode(hexdump, []byte(input))
    for i := 0; i < L; i++ {
        j := i * 4
        k := i * 2
        copy(payload[j:j+4], `\x`+string(hexdump[k:k+2]))
    }
    return string(output)
}
