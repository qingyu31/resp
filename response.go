package resp

import (
	"bufio"
	"io"
)

type ResponseReader struct {
	reader *bufio.Reader
}

func NewResponseReader(r io.Reader) *ResponseReader {
	rr := new(ResponseReader)
	rr.reader = bufio.NewReader(r)
	return rr
}

func (rr ResponseReader) Next() (Message, error) {
	return ReadMessage(rr.reader)
}

type ResponseWriter struct {
	writer *bufio.Writer
}

func NewResponseWriter(w io.Writer) *ResponseWriter {
	rw := new(ResponseWriter)
	rw.writer = bufio.NewWriter(w)
	return rw
}

func (rw ResponseWriter) Write(messages ...Message) error {
	for _, message := range messages {
		err := WriteMessage(rw.writer, message)
		if err != nil {
			return err
		}
	}
	return rw.writer.Flush()
}
