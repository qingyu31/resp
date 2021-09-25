package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/qingyu31/resp"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var host = flag.String("h", "127.0.0.1", "hostname")
var port = flag.Int("p", 6379, "port")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Printf("dial tcp:%v", err)
		return
	}
	defer conn.Close()
	requestWriter := resp.NewRequestWriter(conn)
	responseReader := resp.NewResponseReader(conn)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s:%d> ", *host, *port)
	for scanner.Scan() {
		strs := strings.Split(scanner.Text(), " ")
		if strs[0] == "" {
			continue
		}
		if strs[0] == "exit" || strs[0] == "quit" {
			return
		}
		args := make([][]byte, 0, len(strs)-1)
		for idx := 1; idx < len(strs); idx++ {
			args = append(args, []byte(strs[idx]))
		}
		req := resp.NewRequest(string(strs[0]), args...)
		err = requestWriter.Write(req)
		if err != nil {
			log.Printf("send command:%v", err)
			return
		}
		ret, err := responseReader.Next()
		if err != nil {
			log.Printf("read response: %v", err)
			return
		} else {
			fmt.Print(toString(ret, 0))
		}
		fmt.Printf("%s:%d> ", *host, *port)
	}
}

func toString(v resp.Message, depth int) string {
	var builder strings.Builder
	for idx := 0; idx < depth; idx++ {
		builder.WriteRune(' ')
	}
	switch val := v.(type) {
	case *resp.Integer:
		builder.WriteString("(integer) ")
		builder.WriteString(strconv.Itoa(val.Value()))
	case *resp.SimpleStrings:
		builder.WriteString(val.Value())
	case *resp.Error:
		builder.WriteString("(error) ")
		builder.WriteString(val.Value().Error())
	case *resp.BulkStrings:
		if val.Value() == nil {
			builder.WriteString("nil")
			break
		}
		builder.WriteString(fmt.Sprintf("\"%s\"", *val.Value()))
	case *resp.Array:
		for idx, item := range val.Value() {
			builder.WriteString(strconv.Itoa(idx) + ") ")
			builder.WriteString(toString(item, depth+1))
		}
	}
	builder.WriteRune('\n')
	return builder.String()
}
