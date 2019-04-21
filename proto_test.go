package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestProto(t *testing.T) {
	body := []byte("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n")
	array, err := decodeArray(bufio.NewReader(bytes.NewReader(body)))
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	var buf bytes.Buffer
	err = encodeArray(&buf, array)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	fmt.Println(buf.String())
	if buf.String() != string(body) {
		t.Fail()
		return
	}
}
