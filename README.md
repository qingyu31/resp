## RESP
An implement of REdis Serialization Protocol. 
Designed using to implement redis-like server or client.

### Usage
#### encoding
1. request
```go
req := NewRequest(command,arguments...)
req.Bytes()
```
2. response
```go
//todo
```
#### decoding
1. request
```go
req, err := ParseRequest(r)
```
2. response
```go
resp, err := ParseResponse(r)
//resp is one of resp.SimpleString, resp.Integer, resp.Error, resp.Bulk or resp.Array.
```

### redis-cli
An implement of go console client with resp.