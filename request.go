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

type RequestReader struct {
	reader *bufio.Reader
}

func NewRequestReader(r io.Reader) *RequestReader {
	rr := new(RequestReader)
	rr.reader = bufio.NewReader(r)
	return rr
}

func (rr RequestReader) Next() (*Request, error) {
	msg, err := ReadMessage(rr.reader)
	if err != nil {
		return nil, err
	}
	array, ok := msg.(*Array)
	if !ok {
		return nil, errIllegalProto
	}
	if len(array.Value()) == 0 {
		return nil, errIllegalRequest
	}
	value := array.Value()
	bulk, ok := value[0].(*BulkStrings)
	if !ok {
		return nil, errIllegalRequest
	}
	if bulk.Value() == nil {
		return nil, errIllegalRequest
	}
	r := new(Request)
	r.command = *bulk.Value()
	for idx, a := range value {
		if idx == 0 {
			continue
		}
		bulk, ok = a.(*BulkStrings)
		if !ok {
			return nil, errIllegalRequest
		}
		if bulk.Value() == nil {
			return nil, errIllegalRequest
		}
		r.arguments = append(r.arguments, []byte(*bulk.Value()))
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
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	msgs := make([]Message, 0, len(r.arguments)+1)
	msgs = append(msgs, NewBulkStrings(&r.command))
	for _, arg := range r.arguments {
		msgs = append(msgs, NewBulkStringsWithBytes(arg))
	}
	a := NewArray(msgs)
	WriteMessage(writer, a)
	writer.Flush()
	return buf.Bytes()
}

type RequestWriter struct {
	writer *bufio.Writer
}

func NewRequestWriter(w io.Writer) *RequestWriter {
	rw := new(RequestWriter)
	rw.writer = bufio.NewWriter(w)
	return rw
}

func (rw RequestWriter) Write(reqs ...*Request) error {
	for _, req := range reqs {
		_, err := rw.writer.Write(req.Bytes())
		if err != nil {
			return err
		}
	}
	return rw.writer.Flush()
}
