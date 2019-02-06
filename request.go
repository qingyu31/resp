package resp

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
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
	head, err := brd.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if head[0] != '*' {
		strs := strings.Split(strings.TrimSuffix(head, LINE_DELIMETER), "")
		r.command = strs[0]
		if len(strs) > 1 {
			r.arguments = make([][]byte, 0, len(strs)-1)
			for _, str := range strs {
				r.arguments = append(r.arguments, []byte(str))
			}
		}
		return r, nil
	}
	argCount, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(head, "*"), LINE_DELIMETER))
	if err != nil {
		return nil, errIllegalProto
	}
	r.arguments = make([][]byte, 0, argCount)
	for idx := 0; idx < argCount; idx++ {
		line, err := brd.ReadString('\n')
		if err != nil {
			return nil, errIllegalProto
		}
		argLen, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(line, "$"), LINE_DELIMETER))
		if err != nil {
			return nil, errIllegalProto
		}
		body, err := brd.ReadBytes('\n')
		body = bytes.TrimSuffix(body, []byte(LINE_DELIMETER))
		if len(body) != argLen {
			return nil, errIllegalProto
		}
		if idx == 0 {
			r.command = string(body)
			continue
		}
		r.arguments = append(r.arguments, body)
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
	b.WriteRune('*')
	b.WriteString(strconv.Itoa(len(r.arguments) + 1))
	b.WriteString(LINE_DELIMETER)
	b.WriteRune('$')
	b.WriteString(strconv.Itoa(len(r.command)))
	b.WriteString(LINE_DELIMETER)
	b.WriteString(r.command)
	b.WriteString(LINE_DELIMETER)
	for _, arg := range r.arguments {
		b.WriteRune('$')
		b.WriteString(strconv.Itoa(len(arg)))
		b.WriteString(LINE_DELIMETER)
		b.Write(arg)
		b.WriteString(LINE_DELIMETER)
	}
	return b.Bytes()
}
