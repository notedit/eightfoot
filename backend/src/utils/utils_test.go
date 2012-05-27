// the utils testing
// author: notedit<notedit@gmail.com>
// date: 20120526

package utils

import (
	"testing"
)

func TestGenUkey(t *testing.T) {
	ukeyMap := map[int]string{ 
        10000000:"05yc1s",
    }
	for randint,ukey := range ukeyMap {
		tukey := Int2Ukey(randint)
		if ukey != tukey {
			t.Error("GenUkey Error")
		}
	}
}

func TestSalt(t *testing.T) {
	salt := GenSalt(10)
	if len(salt) != 10 {
		t.Error("GenSalt Error")
	}
	salt = GenSalt(12)
	if len(salt) != 12 {
		t.Error("GenSalt Error")
	}
}
