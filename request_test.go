package resp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewRequest(t *testing.T) {
	cmd := "SET"
	args := [][]byte{[]byte("mykey"), []byte("myvalue")}
	bs := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	for i := 0; i < 3; i++ {
		r := NewRequest(cmd, args...)
		ret := string(r.Bytes())
		fmt.Println(ret)
		if ret != bs {
			t.Fail()
			return
		}
	}

}

func TestParseRequest(t *testing.T) {
	cmd := "SET"
	args := [][]byte{[]byte("mykey"), []byte("myvalue")}
	bs := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	for i := 0; i < 3; i++ {
		rr := NewRequestReader(bytes.NewReader([]byte(bs)))
		r, e := rr.Next()
		if e != nil {
			t.Fail()
			return
		}
		if r.GetCommand() != cmd {
			fmt.Printf("cmd:%s\n", r.GetCommand())
			t.Fail()
			return
		}
		for i, arg := range r.GetArguments() {
			if string(arg) != string(args[i]) {
				fmt.Printf("arg[%d] = %s\n ", i, string(arg))
				t.Fail()
				return
			}
		}
	}

}
