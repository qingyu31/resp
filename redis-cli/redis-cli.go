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
		_, err = conn.Write(req.Bytes())
		if err != nil {
			log.Printf("send command:%v", err)
			return
		}
		ret, err := resp.ParseResponse(conn)
		if err != nil {
			log.Printf("read response: %v", err)
			return
		} else {
			fmt.Print(toString(ret, 0))
		}
		fmt.Printf("%s:%d> ", *host, *port)
	}
}

func toString(v interface{}, depth int) string {
	var builder strings.Builder
	for idx := 0; idx < depth; idx++ {
		builder.WriteRune(' ')
	}
	switch val := v.(type) {
	case resp.Integer:
		builder.WriteString("(integer) ")
		builder.WriteString(strconv.Itoa(int(val)))
	case resp.SimpleString:
		builder.WriteString(string(val))
	case resp.Error:
		builder.WriteString("(error) ")
		builder.WriteString(val.Error())
	case resp.Bulk:
		builder.WriteString(fmt.Sprintf("\"%s\"", string(val)))
	case resp.Array:
		for idx, item := range val {
			builder.WriteString(strconv.Itoa(idx) + ") ")
			builder.WriteString(toString(item, depth+1))
		}
	}
	builder.WriteRune('\n')
	return builder.String()
}
