## RESP
An implement of REdis Serialization Protocol. 
Designed using to implement redis-like server or client.

### Usage
#### client side
1. send request
```go
requestWriter := resp.NewRequestWriter(conn)
request := resp.NewRequest(command,arguments...)
requestWriter.Write(request)
```
2. receive response
```go
responseReader := resp.NewResponseReader(conn)
msg, err := responseReader.Next()
```
#### server side
1. receive request
```go
requestReader := resp.NewRequestReader(conn)
request, err := requestReader.Next()
```
2. send response
```go
responseWriter := resp.NewResponseWriter(conn)
responseWriter.Write(msg)
```

### redis-cli
An implement of go console client with resp.