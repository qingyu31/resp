package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestSimpleStrings(t *testing.T) {
	origins := []string{"hello", ""}
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	for _, origin := range origins {
		ss := new(SimpleStrings)
		ss.value = origin
		ss.WriteTo(writer)
	}
	writer.Flush()
	fmt.Println(buf.String())
	reader := bufio.NewReader(buf)
	for _, origin := range origins {
		msg, err := ReadMessage(reader)
		if err != nil {
			t.Error(err)
			return
		}
		result := msg.(*SimpleStrings).value
		if result != origin {
			t.Errorf("origin=%s result=%s\n", origin, result)
			return

		}
	}
}

func TestIntegerStrings(t *testing.T) {
	origins := []int{1111, -100, math.MaxInt32, 0, math.MinInt32}
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	for _, origin := range origins {
		ss := new(Integer)
		ss.value = origin
		ss.WriteTo(writer)
	}
	writer.Flush()
	fmt.Println(buf.String())
	reader := bufio.NewReader(buf)
	for _, origin := range origins {
		msg, err := ReadMessage(reader)
		if err != nil {
			t.Error(err)
			return
		}
		result := msg.(*Integer).value
		if result != origin {
			t.Errorf("origin=%d result=%d\n", origin, result)
			return

		}
	}
}

func TestBulkStrings(t *testing.T) {
	origins := []string{"hello world", "", "hello\nworld", "hello\r\nworld"}
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	for _, origin := range origins {
		ss := new(BulkStrings)
		ss.value = []byte(origin)
		ss.WriteTo(writer)
	}
	ss := new(BulkStrings)
	ss.WriteTo(writer)
	writer.Flush()
	fmt.Println(buf.String())
	reader := bufio.NewReader(buf)
	for idx := 0; idx < len(origins)+1; idx++ {
		msg, err := ReadMessage(reader)
		if err != nil {
			t.Error(err)
			return
		}
		value := msg.(*BulkStrings).value
		if idx == len(origins) {
			if value != nil {
				t.Errorf("result=%v is not nil", value)
				return
			}
			continue
		}
		origin := origins[idx]
		result := string(value)
		if result != origin {
			t.Errorf("origin=%s result=%s\n", origin, result)
			return

		}
	}
}

func TestArray(t *testing.T)  {
	origins := []string{
		"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n",
		"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n",
		"*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n",
}
	builder := new(bytes.Buffer)
	for _, origin := range origins {
		builder.WriteString(origin)
	}
	reader := bufio.NewReader(builder)
	for _, origin := range origins {
		msg, err := ReadMessage(reader)
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}
		buf := new(bytes.Buffer)
		writer := bufio.NewWriter(buf)
		err = msg.WriteTo(writer)
		writer.Flush()
		if err != nil {
			t.Log(err)
			t.Fail()
			return
		}
		result := buf.String()
		fmt.Println(result)
		if origin != result {
			t.Errorf("origin=%s result=%s\n", origin, result)
			return
		}
	}
}
