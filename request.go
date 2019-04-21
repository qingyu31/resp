package resp

import (
	"bufio"
	"bytes"
	"io"
)

type Request struct {
	command   string
	arguments [][]byte
}

func NewRequest(command string, arguments ...[]byte) *Request {
	r := new(Request)
	r.command = command
	r.arguments = arguments
	return r
}

func ReadRequest(rd io.Reader) ([]byte, error) {
	brd := bufio.NewReader(rd)
	head, err := brd.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if head[0] != '*' {
		return head, nil
	}
	var ret bytes.Buffer
	argCount := 0
	var char byte
	for idx := 0; idx < len(head); idx++ {
		char = head[idx]
		if char < '0' || char > '9' {
			return nil, errIllegalProto
		}
		argCount = argCount*10 + int(head[idx]-'0')
	}
	for idx := 0; idx < argCount*2; idx++ {
		line, err := brd.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		ret.Write(line)
	}
	return ret.Bytes(), nil
}

func ParseRequest(rd io.Reader) (*Request, error) {
	r := new(Request)
	brd := bufio.NewReader(rd)
	array, err := decodeArray(brd)
	if err != nil {
		return nil, err
	}
	if len(array) == 0 {
		return nil, errIllegalProto
	}
	bulk, ok := array[0].(Bulk)
	if !ok {
		return nil, errIllegalProto
	}
	r.command = string(bulk)
	for idx, a := range array {
		if idx == 0 {
			continue
		}
		bulk, ok := a.(Bulk)
		if !ok {
			return nil, errIllegalProto
		}
		r.arguments = append(r.arguments, bulk)
	}
	return r, nil
}

func (r *Request) GetCommand() string {
	return r.command
}

func (r *Request) GetArguments() [][]byte {
	return r.arguments
}

func (r *Request) Bytes() []byte {
	var b bytes.Buffer
	array := make(Array,0,len(r.arguments)+1)
	array = append(array, Bulk(r.command))
	for _, arg := range r.arguments {
		array = append(array, Bulk(arg))
	}
	encodeArray(&b, array)
	return b.Bytes()
}
