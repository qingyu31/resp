package main

import (
	"flag"
	"fmt"
	"github.com/qingyu31/resp"
	"log"
	"net"
	"strconv"
)

var port = flag.Int("p", 6379, "port")

func main() {
	flag.Parse()
	addr := ":"+strconv.Itoa(*port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("listen %s failed %v\n", addr, err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listen %s failed %v\n", addr, err)
			return
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	requestReader := resp.NewRequestReader(conn)
	responseWriter := resp.NewResponseWriter(conn)
	for {
		req, err := requestReader.Next()
		if err != nil {
			log.Printf("connection %s %s\n", conn.RemoteAddr().String(), err)
			return
		}
		err = handleRequest(responseWriter, req)
		if err != nil {
			log.Printf("connection %s %s\n", conn.RemoteAddr().String(), err)
			return
		}
	}

}

var hash = make(map[string][]byte, 1024*1024)

func handleRequest(writer *resp.ResponseWriter, request *resp.Request) error {
	switch request.GetCommand() {
	case "get":
		if len(request.GetArguments()) != 1 {
			e := resp.NewError(fmt.Errorf("ERR wrong number of arguments for '%s' command", request.GetCommand()))
			return writer.Write(e)
		}
		val, ok := hash[string(request.GetArguments()[0])]
		if !ok {
			return writer.Write(resp.NewBulkStrings(nil))
		}
		return writer.Write(resp.NewBulkStringsWithBytes(val))
	case "set":
		if len(request.GetArguments()) != 2 {
			e := resp.NewError(fmt.Errorf("ERR wrong number of arguments for '%s' command", request.GetCommand()))
			return writer.Write(e)
		}
		hash[string(request.GetArguments()[0])] = request.GetArguments()[1]
		return writer.Write(resp.NewSimpleStrings("OK"))
	case "del":
		if len(request.GetArguments()) != 1 {
			e := resp.NewError(fmt.Errorf("ERR wrong number of arguments for '%s' command", request.GetCommand()))
			return writer.Write(e)
		}
		hash[string(request.GetArguments()[0])] = nil
		return writer.Write(resp.NewInteger(1))
	case "info":
		return writer.Write(resp.NewBulkStringsWithBytes([]byte("A redis-like server implement by resp.\r\nThanks for using.\n")))
	default:
		e := resp.NewError(fmt.Errorf("ERR unknown command `%s`", request.GetCommand()))
		return writer.Write(e)
	}
	return nil
}